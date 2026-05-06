package handlers

import (
	"context"
	"middleman/assistants/assistantspb"
	"middleman/assistants/internal/domain"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
		domain.AssistantCreatedEvent,
		domain.AssistantActivatedEvent,
		domain.AssistantDeactivatedEvent,
		domain.AssistantConfigurationUpdatedEvent,
		domain.AssistantRequestProcessedEvent,
		domain.ConversationCreatedEvent,
		domain.MessageAddedEvent,
		domain.ConversationContextUpdatedEvent,
		domain.ConversationArchivedEvent,
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
	case domain.AssistantCreatedEvent:
		return h.onAssistantCreated(ctx, event)
	case domain.AssistantActivatedEvent:
		return h.onAssistantActivated(ctx, event)
	case domain.AssistantDeactivatedEvent:
		return h.onAssistantDeactivated(ctx, event)
	case domain.AssistantConfigurationUpdatedEvent:
		return h.onAssistantConfigurationUpdated(ctx, event)
	case domain.AssistantRequestProcessedEvent:
		return h.onAssistantRequestProcessed(ctx, event)
	case domain.ConversationCreatedEvent:
		return h.onConversationCreated(ctx, event)
	case domain.MessageAddedEvent:
		return h.onMessageAdded(ctx, event)
	case domain.ConversationContextUpdatedEvent:
		return h.onConversationContextUpdated(ctx, event)
	case domain.ConversationArchivedEvent:
		return h.onConversationArchived(ctx, event)
	}
	return nil
}

// Implement missing event handling methods

func (h domainHandlers[T]) onAssistantCreated(ctx context.Context, event ddd.Event) error {
	assistant := event.Payload().(*domain.Assistant)

	// Convert domain capabilities to protobuf capabilities with deduplication
	protoCapabilities := domainToProtoCapabilities(assistant.Capabilities)

	return h.publisher.Publish(ctx, assistantspb.AssistantAggregateChannel,
		ddd.NewEvent(assistantspb.AssistantCreatedEvent, &assistantspb.AssistantCreated{
			Id:           assistant.ID(),
			Name:         assistant.Name,
			Description:  assistant.Description,
			Temperature:  assistant.Temperature,
			Active:       assistant.Active,
			MaxTokens:    int32(assistant.MaxTokens),
			Capabilities: protoCapabilities,
			UserId:       assistant.UserID,
		}),
	)
}

func (h domainHandlers[T]) onAssistantActivated(ctx context.Context, event ddd.Event) error {
	assistant := event.Payload().(*domain.Assistant)

	return h.publisher.Publish(ctx, assistantspb.AssistantAggregateChannel,
		ddd.NewEvent(assistantspb.AssistantActivatedEvent, &assistantspb.AssistantActivated{
			Id:        assistant.ID(),
			Timestamp: assistant.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}),
	)
}

func (h domainHandlers[T]) onAssistantDeactivated(ctx context.Context, event ddd.Event) error {
	assistant := event.Payload().(*domain.Assistant)

	return h.publisher.Publish(ctx, assistantspb.AssistantAggregateChannel,
		ddd.NewEvent(assistantspb.AssistantDeactivatedEvent, &assistantspb.AssistantDeactivated{
			Id:        assistant.ID(),
			Timestamp: assistant.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}),
	)
}

func (h domainHandlers[T]) onAssistantConfigurationUpdated(ctx context.Context, event ddd.Event) error {
	assistant := event.Payload().(*domain.Assistant)

	return h.publisher.Publish(ctx, assistantspb.AssistantAggregateChannel,
		ddd.NewEvent(assistantspb.AssistantConfigurationUpdatedEvent, &assistantspb.AssistantConfigurationUpdated{
			Id:           assistant.ID(),
			Temperature:  assistant.Temperature,
			MaxTokens:    int32(assistant.MaxTokens),
			SystemPrompt: assistant.SystemPrompt,
			UpdatedAt:    assistant.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}),
	)
}

func (h domainHandlers[T]) onAssistantRequestProcessed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.AssistantRequestProcessed)

	assistantID := event.ID() // This should be the aggregate ID

	// Convert domain AssistantAction to protobuf AssistantAction
	var protoActions []*assistantspb.AssistantAction
	for _, action := range payload.ResponseActions {
		// Convert parameters to map[string]string
		parameters := make(map[string]string)
		for k, v := range action.Parameters {
			if str, ok := v.(string); ok {
				parameters[k] = str
			}
		}

		protoActions = append(protoActions, &assistantspb.AssistantAction{
			Type:        action.Type,
			Endpoint:    action.Endpoint,
			Method:      action.Method,
			Parameters:  parameters,
			Description: action.Description,
		})
	}

	// Convert response data to map[string]string
	responseData := make(map[string]string)
	for k, v := range payload.ResponseData {
		if str, ok := v.(string); ok {
			responseData[k] = str
		}
	}

	// Convert request context to map[string]string
	requestContext := make(map[string]string)
	for k, v := range payload.Context {
		if str, ok := v.(string); ok {
			requestContext[k] = str
		}
	}

	return h.publisher.Publish(ctx, assistantspb.AssistantAggregateChannel,
		ddd.NewEvent(assistantspb.AssistantRequestProcessedEvent, &assistantspb.AssistantRequestProcessed{
			Id:              assistantID,
			RequestId:       payload.RequestID,
			UserId:          payload.UserID,
			Message:         payload.Message,
			Context:         requestContext,
			RequestType:     payload.RequestType,
			ResponseId:      payload.ResponseID,
			ResponseMessage: payload.ResponseMessage,
			ResponseData:    responseData,
			Actions:         protoActions,
			Timestamp:       payload.ResponseTimestamp.Format("2006-01-02T15:04:05Z07:00"),
			Status:          payload.ResponseStatus,
			Confidence:      payload.ResponseConfidence,
		}),
	)
}

func (h domainHandlers[T]) onAssistantLoggedOut(ctx context.Context, event ddd.Event) error {
	assistant := event.Payload().(*domain.Assistant)
	return h.publisher.Publish(ctx, assistantspb.AssistantAggregateChannel,
		ddd.NewEvent(assistantspb.AssistantLoggedOutEvent, &assistantspb.AssistantLoggedOut{
			Id: assistant.ID(),
		}),
	)
}

func (h domainHandlers[T]) onAssistantAuthorized(ctx context.Context, event ddd.Event) error {
	assistant := event.Payload().(*domain.Assistant)
	return h.publisher.Publish(ctx, assistantspb.AssistantAggregateChannel,
		ddd.NewEvent(assistantspb.AssistantAuthorizedEvent, &assistantspb.AssistantAuthorized{
			Id: assistant.ID(),
		}),
	)
}

func (h domainHandlers[T]) onAssistantParticipationEnabled(ctx context.Context, event ddd.Event) error {
	assistant := event.Payload().(*domain.Assistant)
	return h.publisher.Publish(ctx, assistantspb.AssistantAggregateChannel,
		ddd.NewEvent(assistantspb.AssistantParticipatingToggledEvent, &assistantspb.AssistantParticipationToggled{
			Id:            assistant.ID(),
			Participating: true,
		}),
	)
}

func (h domainHandlers[T]) onAssistantParticipationDisabled(ctx context.Context, event ddd.Event) error {
	assistant := event.Payload().(*domain.Assistant)
	return h.publisher.Publish(ctx, assistantspb.AssistantAggregateChannel,
		ddd.NewEvent(assistantspb.AssistantParticipatingToggledEvent, &assistantspb.AssistantParticipationToggled{
			Id:            assistant.ID(),
			Participating: false,
		}),
	)
}

func (h domainHandlers[T]) onConversationCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ConversationCreated)

	// Convert context map to string map
	ctxMap := make(map[string]string)
	for k, v := range payload.Context {
		if str, ok := v.(string); ok {
			ctxMap[k] = str
		}
	}

	return h.publisher.Publish(ctx, assistantspb.ConversationAggregateChannel,
		ddd.NewEvent(assistantspb.ConversationCreatedEvent, &assistantspb.ConversationCreated{
			ConversationId: payload.ConversationID,
			UserId:         payload.UserID,
			AssistantId:    payload.AssistantID,
			CreatedAt:      payload.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Context:        ctxMap,
		}),
	)
}

func (h domainHandlers[T]) onMessageAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.MessageAdded)

	// Convert metadata
	metaMap := make(map[string]string)
	for k, v := range payload.Metadata {
		if str, ok := v.(string); ok {
			metaMap[k] = str
		}
	}

	// Convert actions
	var pbActions []*assistantspb.AssistantAction
	for _, act := range payload.ActionsTaken {
		params := make(map[string]string)
		for k, v := range act.Parameters {
			if str, ok := v.(string); ok {
				params[k] = str
			}
		}
		pbActions = append(pbActions, &assistantspb.AssistantAction{
			Type:        act.Type,
			Endpoint:    act.Endpoint,
			Method:      act.Method,
			Parameters:  params,
			Description: act.Description,
		})
	}

	// Map role
	var role assistantspb.MessageRole
	switch payload.Role {
	case domain.AssistantRole:
		role = assistantspb.MessageRole_ASSISTANT
	case domain.SystemRole:
		role = assistantspb.MessageRole_SYSTEM
	default:
		role = assistantspb.MessageRole_USER
	}

	pbMsg := &assistantspb.ConversationMessage{
		Id:           payload.ID,
		Role:         role,
		Content:      payload.Content,
		Timestamp:    payload.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		Metadata:     metaMap,
		ActionsTaken: pbActions,
	}

	return h.publisher.Publish(ctx, assistantspb.ConversationAggregateChannel,
		ddd.NewEvent(assistantspb.MessageAddedEvent, &assistantspb.MessageAdded{
			ConversationId: payload.ConversationID,
			Message:        pbMsg,
			Timestamp:      payload.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		}),
	)
}

func (h domainHandlers[T]) onConversationContextUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ConversationContextUpdated)

	ctxMap := make(map[string]string)
	for k, v := range payload.Context {
		if str, ok := v.(string); ok {
			ctxMap[k] = str
		}
	}

	return h.publisher.Publish(ctx, assistantspb.ConversationAggregateChannel,
		ddd.NewEvent(assistantspb.ConversationContextUpdatedEvent, &assistantspb.ConversationContextUpdated{
			ConversationId: payload.ConversationID,
			Context:        ctxMap,
			UpdatedAt:      payload.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}),
	)
}

func (h domainHandlers[T]) onConversationArchived(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ConversationArchived)

	return h.publisher.Publish(ctx, assistantspb.ConversationAggregateChannel,
		ddd.NewEvent(assistantspb.ConversationArchivedEvent, &assistantspb.ConversationArchived{
			ConversationId: payload.ConversationID,
			ArchivedAt:     payload.ArchivedAt.Format("2006-01-02T15:04:05Z07:00"),
		}),
	)
}

// domainToProtoCapabilities converts domain capabilities to protobuf capabilities with deduplication
func domainToProtoCapabilities(caps []domain.AssistantCapability) []assistantspb.AssistantCapability {
	seen := make(map[assistantspb.AssistantCapability]bool)
	result := make([]assistantspb.AssistantCapability, 0, len(caps))
	
	for _, cap := range caps {
		var protoCap assistantspb.AssistantCapability
		switch cap {
		case domain.CapabilityAssistantManagement:
			protoCap = assistantspb.AssistantCapability_ASSISTANT_MANAGEMENT
		case domain.CapabilityUserInteraction:
			protoCap = assistantspb.AssistantCapability_USER_INTERACTION
		case domain.CapabilityDataAnalysis:
			protoCap = assistantspb.AssistantCapability_DATA_ANALYSIS
		case domain.CapabilityLocationServices:
			protoCap = assistantspb.AssistantCapability_LOCATION_SERVICES
		case domain.CapabilityAuthentication:
			protoCap = assistantspb.AssistantCapability_AUTHENTICATION
		case domain.CapabilityPublicAPIAccess:
			protoCap = assistantspb.AssistantCapability_PUBLIC_API_ACCESS
		case domain.CapabilityJailbreakResistant:
			protoCap = assistantspb.AssistantCapability_JAILBREAK_RESISTANT
		case domain.CapabilityScopeEnforcement:
			protoCap = assistantspb.AssistantCapability_SCOPE_ENFORCEMENT
		case domain.CapabilityDataRetrieval:
			protoCap = assistantspb.AssistantCapability_DATA_RETRIEVAL
		case domain.CapabilitySearchAndFilter:
			protoCap = assistantspb.AssistantCapability_SEARCH_AND_FILTER
		case domain.CapabilityPrivateAPIAccess:
			protoCap = assistantspb.AssistantCapability_PRIVATE_API_ACCESS
		case domain.CapabilityUserDataAccess:
			protoCap = assistantspb.AssistantCapability_USER_DATA_ACCESS
		case domain.CapabilityTokenManagement:
			protoCap = assistantspb.AssistantCapability_TOKEN_MANAGEMENT
		case domain.CapabilityDataMasking:
			protoCap = assistantspb.AssistantCapability_DATA_MASKING
		case domain.CapabilityAuditLogging:
			protoCap = assistantspb.AssistantCapability_AUDIT_LOGGING
		case domain.CapabilityTextGeneration:
			protoCap = assistantspb.AssistantCapability_TEXT_GENERATION
		case domain.CapabilityCodeGeneration:
			protoCap = assistantspb.AssistantCapability_CODE_GENERATION
		case domain.CapabilityWebSearch:
			protoCap = assistantspb.AssistantCapability_WEB_SEARCH
		default:
			// Skip unknown capabilities instead of defaulting
			continue
		}
		
		if !seen[protoCap] {
			seen[protoCap] = true
			result = append(result, protoCap)
		}
	}
	
	return result
}
