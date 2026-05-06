package handlers

import (
	"context"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"middleman/managers/internal/constants"
	"middleman/managers/internal/domain"

	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type managerCatalogHandlers[T ddd.Event] struct {
	catalog domain.CatalogRepository
}

var _ ddd.EventHandler[ddd.Event] = (*managerCatalogHandlers[ddd.Event])(nil)

func NewManagerCatalogHandlers(catalog domain.CatalogRepository) ddd.EventHandler[ddd.Event] {
	return managerCatalogHandlers[ddd.Event]{
		catalog: catalog,
	}
}

func RegisterManagerCatalogHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.ManagerCreatedEvent,
		domain.ManagerActivatedEvent,
		domain.ManagerDeactivatedEvent,
		domain.ManagerConfigurationUpdatedEvent,
	)
}

func RegisterManagerCatalogHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		catalogHandlers := di.Get(ctx, constants.ManagerCatalogHandlersKey).(ddd.EventHandler[ddd.Event])
		return catalogHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])
	RegisterManagerCatalogHandlers(subscriber, handlers)
}

func (h managerCatalogHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling manager catalog event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled manager catalog event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling manager catalog event", trace.WithAttributes(
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
	}
	return nil
}

func (h managerCatalogHandlers[T]) onManagerCreated(ctx context.Context, event ddd.Event) error {

	manager := event.Payload().(*domain.Manager)

	// Add to manager catalog using the CatalogRepository interface with all required parameters
	return h.catalog.Add(
		ctx,
		manager.ID(),
		manager.Name,
		manager.Description,
		manager.UserID,
		manager.Type,
		manager.Capabilities,
		manager.Active,
		manager.Temperature,
		manager.MaxTokens,
		manager.SystemPrompt,
	)
}

func (h managerCatalogHandlers[T]) onManagerActivated(ctx context.Context, event ddd.Event) error {

	manager := event.Payload().(*domain.Manager)

	// Use partial update to only update the active status and timestamp
	return h.catalog.UpdateActiveStatus(ctx, manager.ID(), manager.Active, manager.UpdatedAt)
}

func (h managerCatalogHandlers[T]) onManagerDeactivated(ctx context.Context, event ddd.Event) error {

	manager := event.Payload().(*domain.Manager)

	// Use partial update to only update the active status and timestamp
	return h.catalog.UpdateActiveStatus(ctx, manager.ID(), manager.Active, manager.UpdatedAt)
}

func (h managerCatalogHandlers[T]) onManagerConfigurationUpdated(ctx context.Context, event ddd.Event) error {
	manager := event.Payload().(*domain.Manager)

	// Check if manager has capabilities to determine which update method to use
	if len(manager.Capabilities) > 0 {
		// Update with capabilities
		return h.catalog.UpdateConfigurationWithCapabilities(
			ctx,
			manager.ID(),
			manager.Temperature,
			manager.MaxTokens,
			manager.SystemPrompt,
			manager.Capabilities,
			manager.UpdatedAt,
		)
	}

	// Use partial update to only update configuration fields and timestamp
	return h.catalog.UpdateConfiguration(
		ctx,
		manager.ID(),
		manager.Temperature,
		manager.MaxTokens,
		manager.SystemPrompt,
		manager.UpdatedAt,
	)
}
