package handlers

import (
	"context"
	"log"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"

	"middleman/posts/postspb"
	"middleman/products/productspb"

	"middleman/users/userspb"
	"middleman/vectors/internal/application"
	"middleman/vectors/internal/models"

	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Pure vector-oriented integration handlers - Uses application layer for vector operations
type vectorIntegrationHandlers[T ddd.Event] struct {
	app application.Application
}

var _ ddd.EventHandler[ddd.Event] = (*vectorIntegrationHandlers[ddd.Event])(nil)

func NewVectorIntegrationEventHandlers(
	reg registry.Registry,
	app application.Application,
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewEventHandler(reg, vectorIntegrationHandlers[ddd.Event]{
		app: app,
	}, zerolog.Logger{}, mws...)
}

func RegisterVectorIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	// Product events - comprehensive vector indexing
	if _, err = subscriber.Subscribe(productspb.ProductAggregateChannel, handlers, am.MessageFilter{
		productspb.ProductAddedEvent,
		productspb.ProductRebrandedEvent,
		productspb.ProductUpdatedEvent,
		productspb.ProductRemovedEvent,
		productspb.ProductPriceIncreasedEvent,
		productspb.ProductPriceDecreasedEvent,
		productspb.ProductLeasedEvent,
		productspb.ProductSoldEvent,
		productspb.ProductPawnedEvent,
		productspb.ProductStockAdjustedEvent,
		productspb.ProductNegotiableToggledEvent,
		productspb.ProductArchivedEvent,
		productspb.ProductThumbnailAddedEvent,
		productspb.ProductThumbnailUpdatedEvent,
	}, am.GroupName("vector-products")); err != nil {
		return
	}

	// Post events
	if _, err = subscriber.Subscribe(postspb.PostAggregateChannel, handlers, am.MessageFilter{
		postspb.PostAddedEvent,
		postspb.PostRemovedEvent,
		postspb.PostUpdatedEvent,
		postspb.PostArchivedEvent,
		postspb.PostThumbnailAddedEvent,
		postspb.PostThumbnailUpdatedEvent,
	}, am.GroupName("vector-posts")); err != nil {
		return
	}

	// Deal events

	// User events (only for basic user indexing)
	if _, err = subscriber.Subscribe(userspb.UserAggregateChannel, handlers, am.MessageFilter{
		userspb.UserCreatedEvent,
		userspb.UserRenamedEvent,
	}, am.GroupName("vector-users")); err != nil {
		return
	}

	return
}

func (h vectorIntegrationHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling vector integration event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled vector integration event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling vector integration event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	// User events (basic indexing)
	case userspb.UserCreatedEvent:
		return h.onUserCreated(ctx, event)
	case userspb.UserRenamedEvent:
		return h.onUserUpdated(ctx, event)

	// Product events - all trigger vector operations
	case productspb.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case productspb.ProductUpdatedEvent,
		productspb.ProductRebrandedEvent,
		productspb.ProductPriceIncreasedEvent,
		productspb.ProductPriceDecreasedEvent,
		productspb.ProductLeasedEvent,
		productspb.ProductSoldEvent,
		productspb.ProductPawnedEvent,
		productspb.ProductStockAdjustedEvent,
		productspb.ProductNegotiableToggledEvent,
		productspb.ProductThumbnailAddedEvent,
		productspb.ProductThumbnailUpdatedEvent:
		return h.onProductUpdated(ctx, event)
	case productspb.ProductRemovedEvent,
		productspb.ProductArchivedEvent:
		return h.onProductRemoved(ctx, event)

	// Post events
	case postspb.PostAddedEvent:
		return h.onPostAdded(ctx, event)
	case postspb.PostUpdatedEvent,
		postspb.PostThumbnailAddedEvent,
		postspb.PostThumbnailUpdatedEvent:
		return h.onPostUpdated(ctx, event)
	case postspb.PostRemovedEvent,
		postspb.PostArchivedEvent:
		return h.onPostRemoved(ctx, event)

	default:
		log.Printf("Unhandled vector integration event: %s", event.EventName())
	}

	return nil
}

// ------------------------------
// User Events - Basic Vector Operations
// ------------------------------

func (h vectorIntegrationHandlers[T]) onUserCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*userspb.UserCreated)

	user := &models.User{
		ID:       payload.GetId(),
		Email:    payload.GetEmail(),
		Username: payload.GetUserName(),
		// Basic user indexing - users are not primary searchable entities
	}

	return h.indexEntityAsync(ctx, "user_created", user.ID, func() error {
		params := application.EntityIndexingParams{
			EntityID:   user.ID,
			EntityType: "user",
			EntityData: user,
			Strategy:   application.StrategyMinimal, // Users are not primary searchable entities
		}
		_, err := h.app.IndexEntity(ctx, params)
		return err
	})
}

func (h vectorIntegrationHandlers[T]) onUserUpdated(ctx context.Context, event T) error {
	var userID string

	switch event.EventName() {
	case userspb.UserRenamedEvent:
		payload := event.Payload().(*userspb.UserRenamed)
		userID = payload.GetId()
	}

	return h.removeAndReindexAsync(ctx, "user_updated", userID, "user")
}

// ------------------------------
// Product Events - Vector Operations Only
// ------------------------------

func (h vectorIntegrationHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductAdded)

	product := &models.Product{
		ProductID:    payload.GetId(),
		Name:         payload.GetName(),
		Description:  payload.GetDescription(),
		BasePrice:    payload.GetBasePrice(),
		UserSellerID: payload.GetUserSellerId(),
		CategoryID:   payload.GetCategoryId(),
		CategorySlug: payload.GetCategorySlug(),
		Brand:        payload.GetBrand(),
		Condition:    payload.GetCondition(),
		Model:        payload.GetModel(),
		Tags:         payload.GetTags(),
		Status:       payload.GetStatus(),
		Negotiable:   payload.GetNegotiable(),
		UserType:     payload.GetUserType(),
		Thumbnail:    payload.GetThumbnail(),
		Lat:          float64(payload.GetLat()),
		Lng:          float64(payload.GetLng()),
		EntityType:   models.ProductType,
	}

	return h.indexEntityAsync(ctx, "product_added", product.ProductID, func() error {
		params := application.EntityIndexingParams{
			EntityID:   product.ProductID,
			EntityType: "product",
			EntityData: product,
			Strategy:   application.StrategyOptimized,
		}
		_, err := h.app.IndexEntity(ctx, params)
		return err
	})
}

func (h vectorIntegrationHandlers[T]) onProductUpdated(ctx context.Context, event ddd.Event) error {
	var productID string

	switch event.EventName() {
	case productspb.ProductUpdatedEvent:
		payload := event.Payload().(*productspb.ProductUpdated)
		productID = payload.GetId()
	case productspb.ProductRebrandedEvent:
		payload := event.Payload().(*productspb.ProductRebranded)
		productID = payload.GetId()
	default:
		// Handle other product update events by extracting ID
		return nil
	}

	return h.removeAndReindexAsync(ctx, "product_updated", productID, string(models.ProductType))
}

func (h vectorIntegrationHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	var productID string

	switch event.EventName() {
	case productspb.ProductRemovedEvent:
		payload := event.Payload().(*productspb.ProductRemoved)
		productID = payload.GetId()
	case productspb.ProductArchivedEvent:
		payload := event.Payload().(*productspb.ProductArchived)
		productID = payload.GetId()
	}

	return h.removeEntityAsync(ctx, "product_removed", productID)
}

// ------------------------------
// Post Events - Vector Operations Only
// ------------------------------

func (h vectorIntegrationHandlers[T]) onPostAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*postspb.PostAdded)

	post := &models.Post{
		PostID:      payload.GetId(),
		Name:        payload.GetName(),
		Description: payload.GetDescription(),
		UserID:      payload.GetUserId(),
		Tags:        payload.GetTags(),
		Status:      payload.GetStatus(),
		Thumbnail:   payload.GetThumbnail(),
		Lat:         float64(payload.GetLat()),
		Lng:         float64(payload.GetLng()),
		EntityType:  models.PostType,
	}

	return h.indexEntityAsync(ctx, "post_added", post.PostID, func() error {
		params := application.EntityIndexingParams{
			EntityID:   post.PostID,
			EntityType: "post",
			EntityData: post,
			Strategy:   application.StrategyOptimized,
		}
		_, err := h.app.IndexEntity(ctx, params)
		return err
	})
}

func (h vectorIntegrationHandlers[T]) onPostUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*postspb.PostUpdated)

	return h.removeAndReindexAsync(ctx, "post_updated", payload.GetId(), string(models.PostType))
}

func (h vectorIntegrationHandlers[T]) onPostRemoved(ctx context.Context, event ddd.Event) error {
	var postID string

	switch event.EventName() {
	case postspb.PostRemovedEvent:
		payload := event.Payload().(*postspb.PostRemoved)
		postID = payload.GetId()
	case postspb.PostArchivedEvent:
		payload := event.Payload().(*postspb.PostArchived)
		postID = payload.GetId()
	}

	return h.removeEntityAsync(ctx, "post_removed", postID)
}

// ------------------------------
// Helper Methods for Async Vector Operations
// ------------------------------

func (h vectorIntegrationHandlers[T]) indexEntityAsync(ctx context.Context, operation, entityID string, indexFunc func() error) error {
	if h.app == nil {
		log.Printf("Application service not available, skipping %s for entity %s", operation, entityID)
		return nil
	}

	go func() {
		if err := indexFunc(); err != nil {
			log.Printf("Failed %s for entity %s: %v", operation, entityID, err)
		} else {
			log.Printf("Successfully completed %s for entity %s", operation, entityID)
		}
	}()

	return nil
}

func (h vectorIntegrationHandlers[T]) removeEntityAsync(ctx context.Context, operation, entityID string) error {
	if h.app == nil {
		log.Printf("Application service not available, skipping %s for entity %s", operation, entityID)
		return nil
	}

	go func() {
		if err := h.app.RemoveEntityVector(ctx, entityID, "unknown"); err != nil {
			log.Printf("Failed %s for entity %s: %v", operation, entityID, err)
		} else {
			log.Printf("Successfully completed %s for entity %s", operation, entityID)
		}
	}()

	return nil
}

func (h vectorIntegrationHandlers[T]) removeAndReindexAsync(ctx context.Context, operation, entityID, entityType string) error {
	if h.app == nil {
		log.Printf("Application service not available, skipping %s for entity %s", operation, entityID)
		return nil
	}

	go func() {
		// First remove the old vector
		if err := h.app.RemoveEntityVector(ctx, entityID, entityType); err != nil {
			log.Printf("Failed to remove entity %s during %s: %v", entityID, operation, err)
		}

		// For reindexing, we would need to fetch fresh data from the source
		// This would typically be done by triggering a reindex event or fetching from repository
		// For now, log that reindexing is needed
		log.Printf("Entity %s (%s) requires reindexing after %s - would need to fetch fresh data", entityID, entityType, operation)
	}()

	return nil
}
