package handlers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/support/internal/domain"
	"middleman/support/supportpb"
	"time"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandlers[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		// Support Channel events
		domain.SupportChannelCreatedEvent,
		domain.SupportChannelSettingsUpdatedEvent,
		domain.SupportChannelClosedEvent,
		domain.SupportChannelReactivatedEvent,
		// Ticket events
		domain.TicketCreatedEvent,
		domain.TicketUpdatedEvent,
		domain.TicketAssignedEvent,
		domain.TicketPriorityUpdatedEvent,
		domain.TicketEscalatedEvent,
		domain.TicketResolvedEvent,
		domain.TicketReopenedEvent,
		domain.TicketClosedEvent,
		domain.TicketsMergedEvent,
		domain.TicketsLinkedEvent,
		// Communication events
		domain.TicketReplyAddedEvent,
		domain.InternalNoteAddedEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling domain event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled domain event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling domain event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	// Support Channel events
	case domain.SupportChannelCreatedEvent:
		return h.onSupportChannelCreated(ctx, event)
	case domain.SupportChannelSettingsUpdatedEvent:
		return h.onSupportChannelSettingsUpdated(ctx, event)
	case domain.SupportChannelClosedEvent:
		return h.onSupportChannelClosed(ctx, event)
	case domain.SupportChannelReactivatedEvent:
		return h.onSupportChannelReactivated(ctx, event)
		
	// Ticket events
	case domain.TicketCreatedEvent:
		return h.onTicketCreated(ctx, event)
	case domain.TicketUpdatedEvent:
		return h.onTicketUpdated(ctx, event)
	case domain.TicketAssignedEvent:
		return h.onTicketAssigned(ctx, event)
	case domain.TicketPriorityUpdatedEvent:
		return h.onTicketPriorityUpdated(ctx, event)
	case domain.TicketEscalatedEvent:
		return h.onTicketEscalated(ctx, event)
	case domain.TicketResolvedEvent:
		return h.onTicketResolved(ctx, event)
	case domain.TicketReopenedEvent:
		return h.onTicketReopened(ctx, event)
	case domain.TicketClosedEvent:
		return h.onTicketClosed(ctx, event)
	case domain.TicketsMergedEvent:
		return h.onTicketsMerged(ctx, event)
	case domain.TicketsLinkedEvent:
		return h.onTicketsLinked(ctx, event)
		
	// Communication events
	case domain.TicketReplyAddedEvent:
		return h.onTicketReplyAdded(ctx, event)
	case domain.InternalNoteAddedEvent:
		return h.onInternalNoteAdded(ctx, event)
	}
	return nil
}

// Support Channel event handlers
func (h domainHandlers[T]) onSupportChannelCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.SupportChannelCreated)
	channelID := event.Metadata().Get(ddd.AggregateIDKey).(string)
	return h.publisher.Publish(ctx, supportpb.SupportChannelAggregateChannel,
		ddd.NewEvent(supportpb.SupportChannelCreatedEvent, &supportpb.SupportChannelCreated{
			Id:          channelID,
			UserId:      payload.UserID,
			BusinessId:  payload.BusinessID,
			ChannelType: mapToProtoChannelType(payload.ChannelType),
			Settings:    mapToProtoSettings(payload.Settings),
			CreatedAt:   toTimestamp(event.OccurredAt()),
		}),
	)
}

func (h domainHandlers[T]) onSupportChannelSettingsUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.SupportChannelSettingsUpdated)
	return h.publisher.Publish(ctx, supportpb.SupportChannelAggregateChannel,
		ddd.NewEvent(supportpb.SupportChannelSettingsUpdatedEvent, &supportpb.SupportChannelSettingsUpdated{
			Id:        event.Metadata().Get(ddd.AggregateIDKey).(string),
			Settings:  mapToProtoSettings(payload.Settings),
			UpdatedAt: toTimestamp(event.OccurredAt()),
		}),
	)
}

func (h domainHandlers[T]) onSupportChannelClosed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.SupportChannelClosed)
	return h.publisher.Publish(ctx, supportpb.SupportChannelAggregateChannel,
		ddd.NewEvent(supportpb.SupportChannelClosedEvent, &supportpb.SupportChannelClosed{
			Id:       event.Metadata().Get(ddd.AggregateIDKey).(string),
			ClosedBy: payload.ClosedBy,
			Reason:   payload.Reason,
			ClosedAt: toTimestamp(event.OccurredAt()),
		}),
	)
}

func (h domainHandlers[T]) onSupportChannelReactivated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.SupportChannelReactivated)
	return h.publisher.Publish(ctx, supportpb.SupportChannelAggregateChannel,
		ddd.NewEvent(supportpb.SupportChannelReactivatedEvent, &supportpb.SupportChannelReactivated{
			Id:            event.Metadata().Get(ddd.AggregateIDKey).(string),
			ReactivatedBy: payload.ReactivatedBy,
			ReactivatedAt: toTimestamp(event.OccurredAt()),
		}),
	)
}

// Ticket event handlers
func (h domainHandlers[T]) onTicketCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketCreated)
	return h.publisher.Publish(ctx, supportpb.TicketAggregateChannel,
		ddd.NewEvent(supportpb.TicketCreatedEvent, &supportpb.TicketCreated{
			Id:          event.Metadata().Get(ddd.AggregateIDKey).(string),
			ChannelId:   payload.ChannelID,
			Title:       payload.Title,
			Description: payload.Description,
			Category:    mapToProtoCategory(payload.Category),
			Priority:    mapToProtoPriority(payload.Priority),
			Tags:        payload.Tags,
			Metadata:    payload.Metadata,
			CreatedBy:   payload.CreatedBy,
			Attachments: mapToProtoAttachments(payload.Attachments),
			CreatedAt:   toTimestamp(event.OccurredAt()),
		}),
	)
}

func (h domainHandlers[T]) onTicketUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketUpdated)
	protoEvent := &supportpb.TicketUpdated{
		Id:        event.Metadata().Get(ddd.AggregateIDKey).(string),
		Tags:      payload.Tags,
		Metadata:  payload.Metadata,
		UpdatedBy: payload.UpdatedBy,
		UpdatedAt: toTimestamp(event.OccurredAt()),
	}
	
	if payload.Title != nil {
		protoEvent.Title = payload.Title
	}
	if payload.Description != nil {
		protoEvent.Description = payload.Description
	}
	if payload.Category != nil {
		category := mapToProtoCategory(*payload.Category)
		protoEvent.Category = &category
	}
	
	return h.publisher.Publish(ctx, supportpb.TicketAggregateChannel,
		ddd.NewEvent(supportpb.TicketUpdatedEvent, protoEvent),
	)
}

func (h domainHandlers[T]) onTicketAssigned(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketAssigned)
	return h.publisher.Publish(ctx, supportpb.TicketAggregateChannel,
		ddd.NewEvent(supportpb.TicketAssignedEvent, &supportpb.TicketAssigned{
			Id:               event.Metadata().Get(ddd.AggregateIDKey).(string),
			AssigneeId:       payload.AssigneeID,
			AssigneeType:     mapToProtoAssigneeType(payload.AssigneeType),
			AssignedBy:       payload.AssignedBy,
			AssignmentReason: payload.AssignmentReason,
			AssignedAt:       toTimestamp(event.OccurredAt()),
		}),
	)
}

func (h domainHandlers[T]) onTicketPriorityUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketPriorityUpdated)
	return h.publisher.Publish(ctx, supportpb.TicketAggregateChannel,
		ddd.NewEvent(supportpb.TicketPriorityUpdatedEvent, &supportpb.TicketPriorityUpdated{
			Id:          event.Metadata().Get(ddd.AggregateIDKey).(string),
			OldPriority: mapToProtoPriority(payload.OldPriority),
			NewPriority: mapToProtoPriority(payload.NewPriority),
			UpdatedBy:   payload.UpdatedBy,
			Reason:      payload.Reason,
			UpdatedAt:   toTimestamp(event.OccurredAt()),
		}),
	)
}

func (h domainHandlers[T]) onTicketEscalated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketEscalated)
	return h.publisher.Publish(ctx, supportpb.TicketAggregateChannel,
		ddd.NewEvent(supportpb.TicketEscalatedEvent, &supportpb.TicketEscalated{
			Id:               event.Metadata().Get(ddd.AggregateIDKey).(string),
			FromTier:         mapToProtoTier(payload.FromTier),
			ToTier:           mapToProtoTier(payload.ToTier),
			EscalatedBy:      payload.EscalatedBy,
			EscalationReason: payload.EscalationReason,
			EscalationNotes:  payload.EscalationNotes,
			EscalatedAt:      toTimestamp(event.OccurredAt()),
		}),
	)
}

func (h domainHandlers[T]) onTicketResolved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketResolved)
	return h.publisher.Publish(ctx, supportpb.TicketAggregateChannel,
		ddd.NewEvent(supportpb.TicketResolvedEvent, &supportpb.TicketResolved{
			Id:               event.Metadata().Get(ddd.AggregateIDKey).(string),
			ResolvedBy:       payload.ResolvedBy,
			Resolution:       payload.Resolution,
			AppliedSolutions: payload.AppliedSolutions,
			ResolvedAt:       toTimestamp(event.OccurredAt()),
		}),
	)
}

func (h domainHandlers[T]) onTicketReopened(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketReopened)
	ticketID := event.Metadata().Get(ddd.AggregateIDKey).(string)
	return h.publisher.Publish(ctx, supportpb.TicketAggregateChannel,
		ddd.NewEvent(supportpb.TicketReopenedEvent, &supportpb.TicketReopened{
			Id:           ticketID,
			ReopenedBy:   payload.ReopenedBy,
			ReopenReason: payload.ReopenReason,
			ReopenCount:  int32(payload.ReopenCount),
			ReopenedAt:   toTimestamp(event.OccurredAt()),
		}),
	)
}

func (h domainHandlers[T]) onTicketClosed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketClosed)
	ticketID := event.Metadata().Get(ddd.AggregateIDKey).(string)
	protoEvent := &supportpb.TicketClosed{
		Id:           ticketID,
		ClosedBy:     payload.ClosedBy,
		ClosureNotes: payload.ClosureNotes,
		ClosedAt:     toTimestamp(event.OccurredAt()),
	}
	
	if payload.SatisfactionRating != nil {
		satisfaction := mapToProtoSatisfaction(*payload.SatisfactionRating)
		protoEvent.SatisfactionRating = satisfaction
	}
	
	return h.publisher.Publish(ctx, supportpb.TicketAggregateChannel,
		ddd.NewEvent(supportpb.TicketClosedEvent, protoEvent),
	)
}

func (h domainHandlers[T]) onTicketsMerged(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketsMerged)
	return h.publisher.Publish(ctx, supportpb.TicketAggregateChannel,
		ddd.NewEvent(supportpb.TicketsMergedEvent, &supportpb.TicketsMerged{
			PrimaryTicketId:    payload.PrimaryTicketID,
			SecondaryTicketIds: payload.SecondaryTicketIDs,
			MergedBy:           payload.MergedBy,
			MergeReason:        payload.MergeReason,
			MergedAt:           toTimestamp(event.OccurredAt()),
		}),
	)
}

func (h domainHandlers[T]) onTicketsLinked(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketsLinked)
	return h.publisher.Publish(ctx, supportpb.TicketAggregateChannel,
		ddd.NewEvent(supportpb.TicketsLinkedEvent, &supportpb.TicketsLinked{
			TicketId:         payload.TicketID,
			RelatedTicketIds: payload.RelatedTicketIDs,
			LinkedBy:         payload.LinkedBy,
			RelationshipType: payload.RelationshipType,
			LinkedAt:         toTimestamp(event.OccurredAt()),
		}),
	)
}

// Communication event handlers
func (h domainHandlers[T]) onTicketReplyAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketReplyAdded)
	return h.publisher.Publish(ctx, supportpb.CommunicationAggregateChannel,
		ddd.NewEvent(supportpb.TicketReplyAddedEvent, &supportpb.TicketReplyAdded{
			Id:          payload.ID,
			TicketId:    payload.TicketID,
			AuthorId:    payload.AuthorID,
			AuthorType:  mapToProtoAuthorType(payload.AuthorType),
			Content:     payload.Content,
			Attachments: mapToProtoAttachments(payload.Attachments),
			IsPublic:    payload.IsPublic,
			CreatedAt:   toTimestamp(event.OccurredAt()),
		}),
	)
}

func (h domainHandlers[T]) onInternalNoteAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.InternalNoteAdded)
	return h.publisher.Publish(ctx, supportpb.CommunicationAggregateChannel,
		ddd.NewEvent(supportpb.InternalNoteAddedEvent, &supportpb.InternalNoteAdded{
			Id:             payload.ID,
			TicketId:       payload.TicketID,
			AuthorId:       payload.AuthorID,
			Content:        payload.Content,
			MentionedUsers: payload.MentionedUsers,
			CreatedAt:      toTimestamp(event.OccurredAt()),
		}),
	)
}

// Helper mapping functions
func mapToProtoChannelType(t domain.SupportChannelType) supportpb.SupportChannelType {
	switch t {
	case domain.ChannelTypeTechnical:
		return supportpb.SupportChannelType_TECHNICAL
	case domain.ChannelTypeBilling:
		return supportpb.SupportChannelType_BILLING
	case domain.ChannelTypeSales:
		return supportpb.SupportChannelType_SALES
	case domain.ChannelTypeProduct:
		return supportpb.SupportChannelType_PRODUCT
	default:
		return supportpb.SupportChannelType_GENERAL
	}
}

func mapToProtoSettings(s domain.SupportChannelSettings) *supportpb.SupportChannelSettings {
	return &supportpb.SupportChannelSettings{
		EmailNotifications:  s.EmailNotifications,
		SmsNotifications:    s.SMSNotifications,
		AutoAssignTickets:   s.AutoAssignTickets,
		PreferredLanguage:   s.PreferredLanguage,
		Timezone:            s.Timezone,
		NotificationEmails:  s.NotificationEmails,
		SlaSettings: &supportpb.SLASettings{
			FirstResponseMinutes:  int32(s.SLASettings.FirstResponseMinutes),
			ResolutionHours:       int32(s.SLASettings.ResolutionHours),
			PriorityResponseTimes: convertIntMapToInt32Map(s.SLASettings.PriorityResponseTimes),
		},
	}
}

func mapToProtoCategory(c domain.TicketCategory) supportpb.TicketCategory {
	switch c {
	case domain.CategoryTechnicalIssue:
		return supportpb.TicketCategory_TECHNICAL_ISSUE
	case domain.CategoryBillingIssue:
		return supportpb.TicketCategory_BILLING_ISSUE
	case domain.CategoryAccountIssue:
		return supportpb.TicketCategory_ACCOUNT_ISSUE
	case domain.CategoryProductQuestion:
		return supportpb.TicketCategory_PRODUCT_QUESTION
	case domain.CategoryFeatureRequest:
		return supportpb.TicketCategory_FEATURE_REQUEST
	case domain.CategoryComplaint:
		return supportpb.TicketCategory_COMPLAINT
	case domain.CategoryRefundRequest:
		return supportpb.TicketCategory_REFUND_REQUEST
	case domain.CategoryOrderIssue:
		return supportpb.TicketCategory_ORDER_ISSUE
	case domain.CategoryShippingIssue:
		return supportpb.TicketCategory_SHIPPING_ISSUE
	default:
		return supportpb.TicketCategory_GENERAL_INQUIRY
	}
}

func mapToProtoPriority(p domain.TicketPriority) supportpb.TicketPriority {
	switch p {
	case domain.PriorityHigh:
		return supportpb.TicketPriority_HIGH
	case domain.PriorityUrgent:
		return supportpb.TicketPriority_URGENT
	case domain.PriorityCritical:
		return supportpb.TicketPriority_CRITICAL
	case domain.PriorityLow:
		return supportpb.TicketPriority_LOW
	default:
		return supportpb.TicketPriority_MEDIUM
	}
}

func mapToProtoAssigneeType(t domain.AssigneeType) supportpb.AssigneeType {
	switch t {
	case domain.AssigneeTypeAI:
		return supportpb.AssigneeType_AI_ASSISTANT
	case domain.AssigneeTypeTeam:
		return supportpb.AssigneeType_TEAM
	default:
		return supportpb.AssigneeType_HUMAN_AGENT
	}
}

func mapToProtoAuthorType(t domain.AuthorType) supportpb.AuthorType {
	switch t {
	case domain.AuthorTypeAgent:
		return supportpb.AuthorType_AGENT
	case domain.AuthorTypeAI:
		return supportpb.AuthorType_AI
	case domain.AuthorTypeSystem:
		return supportpb.AuthorType_SYSTEM
	default:
		return supportpb.AuthorType_CUSTOMER
	}
}

func mapToProtoTier(t domain.SupportTier) supportpb.SupportTier {
	switch t {
	case domain.TierTwo:
		return supportpb.SupportTier_TIER_2
	case domain.TierThree:
		return supportpb.SupportTier_TIER_3
	case domain.TierManagement:
		return supportpb.SupportTier_MANAGEMENT
	default:
		return supportpb.SupportTier_TIER_1
	}
}

func mapToProtoSatisfaction(s domain.CustomerSatisfaction) supportpb.CustomerSatisfaction {
	switch s {
	case domain.SatisfactionVeryDissatisfied:
		return supportpb.CustomerSatisfaction_VERY_DISSATISFIED
	case domain.SatisfactionDissatisfied:
		return supportpb.CustomerSatisfaction_DISSATISFIED
	case domain.SatisfactionNeutral:
		return supportpb.CustomerSatisfaction_NEUTRAL
	case domain.SatisfactionSatisfied:
		return supportpb.CustomerSatisfaction_SATISFIED
	case domain.SatisfactionVerySatisfied:
		return supportpb.CustomerSatisfaction_VERY_SATISFIED
	default:
		return supportpb.CustomerSatisfaction_NEUTRAL
	}
}

func mapToProtoAttachments(attachments []domain.Attachment) []*supportpb.Attachment {
	result := make([]*supportpb.Attachment, len(attachments))
	for i, a := range attachments {
		result[i] = &supportpb.Attachment{
			Id:          a.ID,
			Filename:    a.Filename,
			ContentType: a.ContentType,
			SizeBytes:   a.SizeBytes,
			Url:         a.URL,
			UploadedAt:  toTimestamp(a.UploadedAt),
		}
	}
	return result
}

func toTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func convertIntMapToInt32Map(m map[string]int) map[string]int32 {
	result := make(map[string]int32)
	for k, v := range m {
		result[k] = int32(v)
	}
	return result
}