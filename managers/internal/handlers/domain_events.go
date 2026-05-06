package handlers

import (
	"context"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/managers/internal/domain"
	"middleman/managers/managerspb"
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
		domain.ManagerCreatedEvent,
		domain.ManagerActivatedEvent,
		domain.ManagerDeactivatedEvent,
		domain.ManagerConfigurationUpdatedEvent,
		domain.ManagerRequestProcessedEvent,
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
	case domain.ManagerCreatedEvent:
		return h.onManagerCreated(ctx, event)
	case domain.ManagerActivatedEvent:
		return h.onManagerActivated(ctx, event)
	case domain.ManagerDeactivatedEvent:
		return h.onManagerDeactivated(ctx, event)
	case domain.ManagerConfigurationUpdatedEvent:
		return h.onManagerConfigurationUpdated(ctx, event)
	case domain.ManagerRequestProcessedEvent:
		return h.onManagerRequestProcessed(ctx, event)
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

func (h domainHandlers[T]) onManagerCreated(ctx context.Context, event ddd.Event) error {
	manager := event.Payload().(*domain.Manager)

	// Convert domain capabilities to protobuf capabilities with deduplication
	protoCapabilities := domainToProtoCapabilities(manager.Capabilities)

	return h.publisher.Publish(ctx, managerspb.ManagerAggregateChannel,
		ddd.NewEvent(managerspb.ManagerCreatedEvent, &managerspb.ManagerCreated{
			Id:           manager.ID(),
			Name:         manager.Name,
			Description:  manager.Description,
			Temperature:  manager.Temperature,
			Active:       manager.Active,
			MaxTokens:    int32(manager.MaxTokens),
			Capabilities: protoCapabilities,
			UserId:       manager.UserID,
		}),
	)
}

func (h domainHandlers[T]) onManagerActivated(ctx context.Context, event ddd.Event) error {
	manager := event.Payload().(*domain.Manager)

	return h.publisher.Publish(ctx, managerspb.ManagerAggregateChannel,
		ddd.NewEvent(managerspb.ManagerActivatedEvent, &managerspb.ManagerActivated{
			Id:        manager.ID(),
			Timestamp: manager.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}),
	)
}

func (h domainHandlers[T]) onManagerDeactivated(ctx context.Context, event ddd.Event) error {
	manager := event.Payload().(*domain.Manager)

	return h.publisher.Publish(ctx, managerspb.ManagerAggregateChannel,
		ddd.NewEvent(managerspb.ManagerDeactivatedEvent, &managerspb.ManagerDeactivated{
			Id:        manager.ID(),
			Timestamp: manager.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}),
	)
}

func (h domainHandlers[T]) onManagerConfigurationUpdated(ctx context.Context, event ddd.Event) error {
	manager := event.Payload().(*domain.Manager)

	return h.publisher.Publish(ctx, managerspb.ManagerAggregateChannel,
		ddd.NewEvent(managerspb.ManagerConfigurationUpdatedEvent, &managerspb.ManagerConfigurationUpdated{
			Id:           manager.ID(),
			Temperature:  manager.Temperature,
			MaxTokens:    int32(manager.MaxTokens),
			SystemPrompt: manager.SystemPrompt,
			UpdatedAt:    manager.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}),
	)
}

func (h domainHandlers[T]) onManagerRequestProcessed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ManagerRequestProcessed)

	managerID := event.ID() // This should be the aggregate ID

	// Convert domain ManagerAction to protobuf ManagerAction
	var protoActions []*managerspb.ManagerAction
	for _, action := range payload.ResponseActions {
		// Convert parameters to map[string]string
		parameters := make(map[string]string)
		for k, v := range action.Parameters {
			if str, ok := v.(string); ok {
				parameters[k] = str
			}
		}

		protoActions = append(protoActions, &managerspb.ManagerAction{
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

	return h.publisher.Publish(ctx, managerspb.ManagerAggregateChannel,
		ddd.NewEvent(managerspb.ManagerRequestProcessedEvent, &managerspb.ManagerRequestProcessed{
			Id:              managerID,
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

func (h domainHandlers[T]) onManagerLoggedOut(ctx context.Context, event ddd.Event) error {
	manager := event.Payload().(*domain.Manager)
	return h.publisher.Publish(ctx, managerspb.ManagerAggregateChannel,
		ddd.NewEvent(managerspb.ManagerLoggedOutEvent, &managerspb.ManagerLoggedOut{
			Id: manager.ID(),
		}),
	)
}

func (h domainHandlers[T]) onManagerAuthorized(ctx context.Context, event ddd.Event) error {
	manager := event.Payload().(*domain.Manager)
	return h.publisher.Publish(ctx, managerspb.ManagerAggregateChannel,
		ddd.NewEvent(managerspb.ManagerAuthorizedEvent, &managerspb.ManagerAuthorized{
			Id: manager.ID(),
		}),
	)
}

func (h domainHandlers[T]) onManagerParticipationEnabled(ctx context.Context, event ddd.Event) error {
	manager := event.Payload().(*domain.Manager)
	return h.publisher.Publish(ctx, managerspb.ManagerAggregateChannel,
		ddd.NewEvent(managerspb.ManagerParticipatingToggledEvent, &managerspb.ManagerParticipationToggled{
			Id:            manager.ID(),
			Participating: true,
		}),
	)
}

func (h domainHandlers[T]) onManagerParticipationDisabled(ctx context.Context, event ddd.Event) error {
	manager := event.Payload().(*domain.Manager)
	return h.publisher.Publish(ctx, managerspb.ManagerAggregateChannel,
		ddd.NewEvent(managerspb.ManagerParticipatingToggledEvent, &managerspb.ManagerParticipationToggled{
			Id:            manager.ID(),
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

	return h.publisher.Publish(ctx, managerspb.ConversationAggregateChannel,
		ddd.NewEvent(managerspb.ConversationCreatedEvent, &managerspb.ConversationCreated{
			ConversationId: payload.ConversationID,
			UserId:         payload.UserID,
			ManagerId:      payload.ManagerID,
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
	var pbActions []*managerspb.ManagerAction
	for _, act := range payload.ActionsTaken {
		params := make(map[string]string)
		for k, v := range act.Parameters {
			if str, ok := v.(string); ok {
				params[k] = str
			}
		}
		pbActions = append(pbActions, &managerspb.ManagerAction{
			Type:        act.Type,
			Endpoint:    act.Endpoint,
			Method:      act.Method,
			Parameters:  params,
			Description: act.Description,
		})
	}

	// Map role
	var role managerspb.MessageRole
	switch payload.Role {
	case domain.ManagerRole:
		role = managerspb.MessageRole_ASSISTANT
	case domain.SystemRole:
		role = managerspb.MessageRole_SYSTEM
	default:
		role = managerspb.MessageRole_USER
	}

	pbMsg := &managerspb.ConversationMessage{
		Id:           payload.ID,
		Role:         role,
		Content:      payload.Content,
		Timestamp:    payload.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		Metadata:     metaMap,
		ActionsTaken: pbActions,
	}

	return h.publisher.Publish(ctx, managerspb.ConversationAggregateChannel,
		ddd.NewEvent(managerspb.MessageAddedEvent, &managerspb.MessageAdded{
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

	return h.publisher.Publish(ctx, managerspb.ConversationAggregateChannel,
		ddd.NewEvent(managerspb.ConversationContextUpdatedEvent, &managerspb.ConversationContextUpdated{
			ConversationId: payload.ConversationID,
			Context:        ctxMap,
			UpdatedAt:      payload.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}),
	)
}

func (h domainHandlers[T]) onConversationArchived(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ConversationArchived)

	return h.publisher.Publish(ctx, managerspb.ConversationAggregateChannel,
		ddd.NewEvent(managerspb.ConversationArchivedEvent, &managerspb.ConversationArchived{
			ConversationId: payload.ConversationID,
			ArchivedAt:     payload.ArchivedAt.Format("2006-01-02T15:04:05Z07:00"),
		}),
	)
}

// domainToProtoCapabilities converts domain capabilities to protobuf capabilities with deduplication
func domainToProtoCapabilities(caps []domain.ManagerCapability) []managerspb.ManagerCapability {
	seen := make(map[managerspb.ManagerCapability]bool)
	result := make([]managerspb.ManagerCapability, 0, len(caps))

	for _, cap := range caps {
		var protoCap managerspb.ManagerCapability
		switch cap {
		case domain.CapabilityManagerManagement:
			protoCap = managerspb.ManagerCapability_ASSISTANT_MANAGEMENT
		case domain.CapabilityUserInteraction:
			protoCap = managerspb.ManagerCapability_USER_INTERACTION
		case domain.CapabilityDataAnalysis:
			protoCap = managerspb.ManagerCapability_DATA_ANALYSIS
		case domain.CapabilityLocationServices:
			protoCap = managerspb.ManagerCapability_LOCATION_SERVICES
		case domain.CapabilityAuthentication:
			protoCap = managerspb.ManagerCapability_AUTHENTICATION
		case domain.CapabilityPublicAPIAccess:
			protoCap = managerspb.ManagerCapability_PUBLIC_API_ACCESS
		case domain.CapabilityJailbreakResistant:
			protoCap = managerspb.ManagerCapability_JAILBREAK_RESISTANT
		case domain.CapabilityScopeEnforcement:
			protoCap = managerspb.ManagerCapability_SCOPE_ENFORCEMENT
		case domain.CapabilityDataRetrieval:
			protoCap = managerspb.ManagerCapability_DATA_RETRIEVAL
		case domain.CapabilitySearchAndFilter:
			protoCap = managerspb.ManagerCapability_SEARCH_AND_FILTER
		case domain.CapabilityPrivateAPIAccess:
			protoCap = managerspb.ManagerCapability_PRIVATE_API_ACCESS
		case domain.CapabilityUserDataAccess:
			protoCap = managerspb.ManagerCapability_USER_DATA_ACCESS
		case domain.CapabilityTokenManagement:
			protoCap = managerspb.ManagerCapability_TOKEN_MANAGEMENT
		case domain.CapabilityDataMasking:
			protoCap = managerspb.ManagerCapability_DATA_MASKING
		case domain.CapabilityAuditLogging:
			protoCap = managerspb.ManagerCapability_AUDIT_LOGGING
		case domain.CapabilityTextGeneration:
			protoCap = managerspb.ManagerCapability_TEXT_GENERATION
		case domain.CapabilityCodeGeneration:
			protoCap = managerspb.ManagerCapability_CODE_GENERATION
		case domain.CapabilityWebSearch:
			protoCap = managerspb.ManagerCapability_WEB_SEARCH
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
