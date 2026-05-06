package handlers

import (
	"context"
	"middleman/assistants/internal/constants"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/errorsotel"

	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type assistantCatalogHandlers[T ddd.Event] struct {
	catalog domain.CatalogRepository
}

var _ ddd.EventHandler[ddd.Event] = (*assistantCatalogHandlers[ddd.Event])(nil)

func NewAssistantCatalogHandlers(catalog domain.CatalogRepository) ddd.EventHandler[ddd.Event] {
	return assistantCatalogHandlers[ddd.Event]{
		catalog: catalog,
	}
}

func RegisterAssistantCatalogHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.AssistantCreatedEvent,
		domain.AssistantActivatedEvent,
		domain.AssistantDeactivatedEvent,
		domain.AssistantConfigurationUpdatedEvent,
	)
}

func RegisterAssistantCatalogHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		catalogHandlers := di.Get(ctx, constants.AssistantCatalogHandlersKey).(ddd.EventHandler[ddd.Event])
		return catalogHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])
	RegisterAssistantCatalogHandlers(subscriber, handlers)
}

func (h assistantCatalogHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling assistant catalog event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled assistant catalog event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling assistant catalog event", trace.WithAttributes(
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
	}
	return nil
}

func (h assistantCatalogHandlers[T]) onAssistantCreated(ctx context.Context, event ddd.Event) error {

	assistant := event.Payload().(*domain.Assistant)

	// Add to assistant catalog using the CatalogRepository interface with all required parameters
	return h.catalog.Add(
		ctx,
		assistant.ID(),
		assistant.Name,
		assistant.Description,
		assistant.UserID,
		assistant.Type,
		assistant.Capabilities,
		assistant.Active,
		assistant.Temperature,
		assistant.MaxTokens,
		assistant.SystemPrompt,
	)
}

func (h assistantCatalogHandlers[T]) onAssistantActivated(ctx context.Context, event ddd.Event) error {

	assistant := event.Payload().(*domain.Assistant)

	// Use partial update to only update the active status and timestamp
	return h.catalog.UpdateActiveStatus(ctx, assistant.ID(), assistant.Active, assistant.UpdatedAt)
}

func (h assistantCatalogHandlers[T]) onAssistantDeactivated(ctx context.Context, event ddd.Event) error {

	assistant := event.Payload().(*domain.Assistant)

	// Use partial update to only update the active status and timestamp
	return h.catalog.UpdateActiveStatus(ctx, assistant.ID(), assistant.Active, assistant.UpdatedAt)
}

func (h assistantCatalogHandlers[T]) onAssistantConfigurationUpdated(ctx context.Context, event ddd.Event) error {
	assistant := event.Payload().(*domain.Assistant)

	// Check if assistant has capabilities to determine which update method to use
	if len(assistant.Capabilities) > 0 {
		// Update with capabilities
		return h.catalog.UpdateConfigurationWithCapabilities(
			ctx,
			assistant.ID(),
			assistant.Temperature,
			assistant.MaxTokens,
			assistant.SystemPrompt,
			assistant.Capabilities,
			assistant.UpdatedAt,
		)
	}

	// Use partial update to only update configuration fields and timestamp
	return h.catalog.UpdateConfiguration(
		ctx,
		assistant.ID(),
		assistant.Temperature,
		assistant.MaxTokens,
		assistant.SystemPrompt,
		assistant.UpdatedAt,
	)
}
