package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"middleman/assistants/internal/application"
	"middleman/assistants/internal/application/commands"
	"middleman/assistants/internal/application/queries"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/users/userspb"
)

type integrationHandlers[T ddd.Event] struct {
	app application.App
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(reg registry.Registry, app application.App, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		app: app,
	}, zerolog.Logger{}, mws...)
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	_, err = subscriber.Subscribe(userspb.UserAggregateChannel, handlers, am.MessageFilter{
		userspb.UserCreatedEvent,
	}, am.GroupName("assistants-users"))
	if err != nil {
		return err
	}

	return err
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
	case userspb.UserCreatedEvent:
		return h.onUserCreated(ctx, event)
	}

	return nil
}
// onUserCreated handles UserCreated events from the users service.
// It creates an assistant for the new user with the appropriate type based on their role.
func (h integrationHandlers[T]) onUserCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*userspb.UserCreated)
	userID := payload.GetId()
	userRole := payload.GetRole()

	// Check if an assistant already exists for this user
	assistants, err := h.app.GetAssistants(ctx, queries.GetAssistants{UserID: userID})
	if err == nil && len(assistants) > 0 {
		return nil // already exists
	}

	// Log the role mapping for debugging
	zerolog.Ctx(ctx).Info().
		Str("userID", userID).
		Str("userRole", userRole).
		Msg("Creating assistant for new user based on role")

	// Create assistant based on user role
	assistantID := uuid.New().String()
	
	switch userRole {
	case "admin":
		return h.app.CreateAdminAssistant(ctx, commands.CreateAdminAssistant{
			ID:     assistantID,
			UserID: userID,
		})
	case "business":
		return h.app.CreateBusinessAssistant(ctx, commands.CreateBusinessAssistant{
			ID:     assistantID,
			UserID: userID,
		})
	case "support":
		return h.app.CreateSupportAssistant(ctx, commands.CreateSupportAssistant{
			ID:     assistantID,
			UserID: userID,
		})
	case "scheduler":
		return h.app.CreateSchedulerAssistant(ctx, commands.CreateSchedulerAssistant{
			ID:     assistantID,
			UserID: userID,
		})
	default:
		// Default to standard user assistant
		return h.app.CreateUserAssistant(ctx, commands.CreateUserAssistant{
			ID:     assistantID,
			UserID: userID,
		})
	}
}
