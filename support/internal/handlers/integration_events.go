package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"middleman/assistants/assistantspb"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/ordering/orderingpb"
	"middleman/payments/paymentspb"
	"middleman/support/internal/application"
	"middleman/support/internal/application/commands"
	"middleman/support/internal/application/queries"
	"middleman/support/internal/domain"
	"middleman/users/userspb"
)

type integrationHandlers[T ddd.Event] struct {
	app application.App
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(
	reg registry.Registry,
	app application.App,
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewEventHandler(
		reg,
		integrationHandlers[ddd.Event]{
			app: app,
		},
		zerolog.Logger{},
		mws...,
	)
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	// Handle events from ordering service - customers may need support with orders
	_, err = subscriber.Subscribe(
		orderingpb.OrderAggregateChannel,
		handlers,
		am.MessageFilter{
			orderingpb.OrderCreatedEvent,
			orderingpb.OrderCompletedEvent,
			orderingpb.OrderCanceledEvent,
		},
		am.GroupName("support-orders"),
	)
	if err != nil {
		return err
	}

	// Handle events from payments service - payment confirmations
	_, err = subscriber.Subscribe(
		paymentspb.InvoiceAggregateChannel,
		handlers,
		am.MessageFilter{
			paymentspb.InvoicePaidEvent,
		},
		am.GroupName("support-payments"),
	)
	if err != nil {
		return err
	}

	// Handle events from users service - new users might need onboarding support
	_, err = subscriber.Subscribe(
		userspb.UserAggregateChannel,
		handlers,
		am.MessageFilter{
			userspb.UserCreatedEvent,
		},
		am.GroupName("support-users"),
	)
	if err != nil {
		return err
	}

	// Handle events from assistants service - AI assistant interactions
	_, err = subscriber.Subscribe(
		assistantspb.ConversationAggregateChannel,
		handlers,
		am.MessageFilter{
			assistantspb.ConversationCreatedEvent,
			assistantspb.MessageAddedEvent,
		},
		am.GroupName("support-assistants"),
	)
	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling integration event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled integration event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling integration event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	// Order events
	case orderingpb.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	case orderingpb.OrderCanceledEvent:
		return h.onOrderCanceled(ctx, event)

	// Payment events
	case paymentspb.InvoicePaidEvent:
		return h.onInvoicePaid(ctx, event)

	// User events
	case userspb.UserCreatedEvent:
		return h.onUserRegistered(ctx, event)

		// Assistant events
		//case assistantspb.ConversationCreatedEvent:
		//	return h.onConversationCreated(ctx, event)
		//case assistantspb.MessageAddedEvent:
		//	return h.onMessageAdded(ctx, event)
	}

	return nil
}

// Order event handlers
func (h integrationHandlers[T]) onOrderCreated(ctx context.Context, event ddd.Event) error {
	// Could automatically create support channel for high-value orders
	// or orders with special requirements
	return nil
}

func (h integrationHandlers[T]) onOrderCanceled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*orderingpb.OrderCanceled)

	// Create a support ticket for canceled orders
	// Check if user has a support channel
	channels, err := h.app.GetUserSupportChannels(ctx, queries.GetUserSupportChannels{
		UserID:     payload.UserCustomerId,
		ActiveOnly: true,
		Page:       1,
		Limit:      1,
	})

	if err != nil {
		return err
	}

	var channelID string
	if len(channels) == 0 {
		// Create a support channel for the user
		channelID = uuid.New().String()
		err = h.app.CreateSupportChannel(ctx, commands.CreateSupportChannel{
			ID:          channelID,
			UserID:      payload.UserCustomerId,
			ChannelType: domain.ChannelTypeGeneral,
		})
		if err != nil {
			return err
		}
	} else {
		channelID = channels[0].ID
	}

	// Create a ticket about the order cancellation
	return h.app.CreateTicket(ctx, commands.CreateTicket{
		ID:          uuid.New().String(),
		ChannelID:   channelID,
		Title:       "Order Cancellation: " + payload.Id,
		Description: "Order was canceled",
		Category:    domain.CategoryOrderIssue,
		Priority:    domain.PriorityMedium,
		CreatedBy:   "system",
		Metadata: map[string]string{
			"order_id": payload.Id,
		},
	})
}

// Payment event handlers
func (h integrationHandlers[T]) onInvoicePaid(ctx context.Context, event ddd.Event) error {
	// Currently just log successful payments, could be used for analytics
	// or to auto-close related billing support tickets
	return nil
}

// User event handlers
func (h integrationHandlers[T]) onUserRegistered(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*userspb.UserCreated)

	// Optionally create a welcome support channel for new users
	// This could be configured based on business requirements
	if payload.Role == "business" {
		// Business users get automatic support channel
		return h.app.CreateSupportChannel(ctx, commands.CreateSupportChannel{
			ID:          uuid.New().String(),
			UserID:      payload.Id,
			BusinessID:  payload.Id, // For business users, user ID is business ID
			ChannelType: domain.ChannelTypeGeneral,
			Settings: domain.SupportChannelSettings{
				EmailNotifications: true,
				AutoAssignTickets:  true,
				PreferredLanguage:  "en",
			},
		})
	}

	return nil
}

// Assistant event handlers
func (h integrationHandlers[T]) onConversationCreated(ctx context.Context, event ddd.Event) error {
	// Track new AI assistant conversations that might need support
	return nil
}

func (h integrationHandlers[T]) onMessageAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*assistantspb.MessageAdded)

	// Check if message indicates user needs human support
	// This is a simplified check - in production you'd use NLP or specific keywords
	if payload.Message != nil && containsEscalationKeywords(payload.Message.Content) {
		// Extract user ID from conversation - for now just use assistant ID as proxy
		// In production, you'd look up the conversation to get the actual user ID
		userID := payload.AssistantId

		// Create support ticket for potential escalation
		channels, err := h.app.GetUserSupportChannels(ctx, queries.GetUserSupportChannels{
			UserID:     userID,
			ActiveOnly: true,
			Page:       1,
			Limit:      1,
		})

		if err != nil {
			return err
		}

		var channelID string
		if len(channels) == 0 {
			channelID = uuid.New().String()
			err = h.app.CreateSupportChannel(ctx, commands.CreateSupportChannel{
				ID:          channelID,
				UserID:      userID,
				ChannelType: domain.ChannelTypeGeneral,
			})
			if err != nil {
				return err
			}
		} else {
			channelID = channels[0].ID
		}

		return h.app.CreateTicket(ctx, commands.CreateTicket{
			ID:          uuid.New().String(),
			ChannelID:   channelID,
			Title:       "User Requested Human Support",
			Description: "User message indicates they need human assistance: " + payload.Message.Content,
			Category:    domain.CategoryGeneralInquiry,
			Priority:    domain.PriorityMedium,
			CreatedBy:   "ai_assistant",
			Metadata: map[string]string{
				"conversation_id": payload.ConversationId,
				"message_id":      payload.Message.Id,
			},
		})
	}

	return nil
}

// Helper function to check for escalation keywords
func containsEscalationKeywords(content string) bool {
	// Simple keyword check - in production use NLP
	keywords := []string{"human", "agent", "support", "help", "speak to someone", "real person"}
	lowerContent := strings.ToLower(content)
	for _, keyword := range keywords {
		if strings.Contains(lowerContent, keyword) {
			return true
		}
	}
	return false
}
