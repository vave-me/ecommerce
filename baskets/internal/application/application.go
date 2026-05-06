package application

import (
	"context"
	"middleman/baskets/internal/domain"
	"middleman/internal/ddd"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stackus/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	StartBasket struct {
		ID             string
		UserCustomerID string
	}

	CancelBasket struct {
		ID string
	}

	CheckoutBasket struct {
		ID              string
		PaymentMethodID string
		PaymentIntentID string
		UserCustomerID  string
	}

	AddItem struct {
		ID        string
		ProductID string
		Quantity  int64
	}

	RemoveItem struct {
		ID        string
		ProductID string
		Quantity  int64
	}

	GetBasket struct {
		ID             string
		UserCustomerID string
	}
	GetTotalBasket struct {
		ID string
	}
	GetCurrentBasket struct {
		UserCustomerID string
	}
	ReopenBasket struct {
		ID     string
		Reason string
	}
	App interface {
		StartBasket(ctx context.Context, start StartBasket) error
		CancelBasket(ctx context.Context, cancel CancelBasket) error
		CheckoutBasket(ctx context.Context, checkout CheckoutBasket) error
		ReopenBasket(ctx context.Context, reopen ReopenBasket) error
		AddItem(ctx context.Context, add AddItem) error
		RemoveItem(ctx context.Context, remove RemoveItem) error
		GetBasket(ctx context.Context, get GetBasket) (*domain.Basket, error)
		GetCurrentBasket(ctx context.Context, get GetCurrentBasket) (*domain.CatalogBasket, error)
		GetTotalBasket(ctx context.Context, get GetTotalBasket) (int64, error)
	}

	Application struct {
		baskets   domain.BasketRepository
		users     domain.UserRepository
		products  domain.ProductRepository
		catalog   domain.CatalogRepository
		publisher ddd.EventPublisher[ddd.Event]
		logger    zerolog.Logger
	}
)

var _ App = (*Application)(nil)

func New(baskets domain.BasketRepository, users domain.UserRepository, products domain.ProductRepository,
	catalog domain.CatalogRepository, publisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		baskets:   baskets,
		users:     users,
		products:  products,
		catalog:   catalog,
		publisher: publisher,
		logger:    log.With().Str("service", "baskets").Logger(),
	}
}

func (a Application) StartBasket(ctx context.Context, start StartBasket) error {
	span := trace.SpanFromContext(ctx)
	logger := a.logger.With().
		Str("operation", "StartBasket").
		Str("basket_id", start.ID).
		Str("user_customer_id", start.UserCustomerID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("BASKETS_START_BASKET_BEGIN: Starting basket creation")

	basket, err := a.baskets.Load(ctx, start.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_START_BASKET_LOAD_FAILED: Failed to load basket")
		return err
	}

	logger.Debug().
		Str("basket_status", string(basket.Status)).
		Int("existing_items", len(basket.Items)).
		Msg("BASKETS_START_BASKET_LOADED: Basket loaded successfully")

	event, err := basket.Start(start.UserCustomerID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_START_BASKET_DOMAIN_FAILED: Domain start operation failed")
		return err
	}

	logger.Debug().
		Str("event_name", event.EventName()).
		Msg("BASKETS_START_BASKET_EVENT_CREATED: Domain event created")

	if err = a.baskets.Save(ctx, basket); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_START_BASKET_SAVE_FAILED: Failed to save basket")
		return err
	}

	logger.Debug().Msg("BASKETS_START_BASKET_SAVED: Basket saved to repository")

	if err = a.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_START_BASKET_PUBLISH_FAILED: Failed to publish event")
		return err
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Msg("BASKETS_START_BASKET_SUCCESS: Basket started successfully")

	span.AddEvent("basket_started", trace.WithAttributes(
		attribute.String("basket_id", start.ID),
		attribute.String("user_customer_id", start.UserCustomerID),
	))

	return nil
}

func (a Application) CancelBasket(ctx context.Context, cancel CancelBasket) error {
	span := trace.SpanFromContext(ctx)
	logger := a.logger.With().
		Str("operation", "CancelBasket").
		Str("basket_id", cancel.ID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("BASKETS_CANCEL_BASKET_BEGIN: Starting basket cancellation")

	basket, err := a.baskets.Load(ctx, cancel.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_CANCEL_BASKET_LOAD_FAILED: Failed to load basket")
		return err
	}

	logger.Debug().
		Str("basket_status", string(basket.Status)).
		Int("items_count", len(basket.Items)).
		Msg("BASKETS_CANCEL_BASKET_LOADED: Basket loaded for cancellation")

	event, err := basket.Cancel()
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_CANCEL_BASKET_DOMAIN_FAILED: Domain cancel operation failed")
		return err
	}

	if err = a.baskets.Save(ctx, basket); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_CANCEL_BASKET_SAVE_FAILED: Failed to save cancelled basket")
		return err
	}

	if err = a.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_CANCEL_BASKET_PUBLISH_FAILED: Failed to publish cancellation event")
		return err
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Msg("BASKETS_CANCEL_BASKET_SUCCESS: Basket cancelled successfully")

	return nil
}

func (a Application) CheckoutBasket(ctx context.Context, checkout CheckoutBasket) error {
	span := trace.SpanFromContext(ctx)
	logger := a.logger.With().
		Str("operation", "CheckoutBasket").
		Str("basket_id", checkout.ID).
		Str("user_customer_id", checkout.UserCustomerID).
		Str("payment_method_id", checkout.PaymentMethodID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("BASKETS_CHECKOUT_BEGIN: Starting basket checkout - CRITICAL WORKFLOW ENTRY POINT")

	basket, err := a.baskets.Load(ctx, checkout.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_CHECKOUT_LOAD_FAILED: Failed to load basket for checkout")
		return errors.Wrap(err, "baskets checkout")
	}

	logger.Info().
		Str("basket_status", string(basket.Status)).
		Int("items_count", len(basket.Items)).
		Int64("total_amount", basket.TotalAmount()).
		Msg("BASKETS_CHECKOUT_LOADED: Basket loaded for checkout")

	// Log each item for detailed tracking
	itemIndex := 0
	for productID, item := range basket.Items {
		logger.Debug().
			Int("item_index", itemIndex).
			Str("product_id", productID).
			Str("seller_id", item.UserSellerID).
			Int64("quantity", item.Quantity).
			Int64("price", item.ProductPrice).
			Msg("BASKETS_CHECKOUT_ITEM: Item details")
		itemIndex++
	}

	event, err := basket.Checkout(checkout.PaymentIntentID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_CHECKOUT_DOMAIN_FAILED: Domain checkout operation failed")
		return errors.Wrap(err, "baskets checkout")
	}

	logger.Info().
		Str("event_name", event.EventName()).
		Msg("BASKETS_CHECKOUT_EVENT_CREATED: BasketCheckedOut event created - TRIGGERING SAGA")

	if err = a.baskets.Save(ctx, basket); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_CHECKOUT_SAVE_FAILED: Failed to save checked out basket")
		return errors.Wrap(err, "basket checkout")
	}

	logger.Debug().
		Str("new_basket_status", string(basket.Status)).
		Msg("BASKETS_CHECKOUT_SAVED: Basket saved with checkout status")

	if err = a.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_CHECKOUT_PUBLISH_FAILED: CRITICAL - Failed to publish BasketCheckedOut event")
		return err
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Str("event_id", event.ID()).
		Msg("BASKETS_CHECKOUT_SUCCESS: BasketCheckedOut event published - SAGA SHOULD START")

	span.AddEvent("basket_checked_out", trace.WithAttributes(
		attribute.String("basket_id", checkout.ID),
		attribute.String("user_customer_id", checkout.UserCustomerID),
		attribute.Int64("total_amount", basket.TotalAmount()),
		attribute.Int("items_count", len(basket.Items)),
		attribute.String("event_id", event.ID()),
	))

	return nil
}

func (a Application) ReopenBasket(ctx context.Context, reopen ReopenBasket) error {
	span := trace.SpanFromContext(ctx)
	logger := a.logger.With().
		Str("operation", "ReopenBasket").
		Str("basket_id", reopen.ID).
		Str("reason", reopen.Reason).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("BASKETS_REOPEN_BEGIN: Reopening basket after failed checkout")

	basket, err := a.baskets.Load(ctx, reopen.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_REOPEN_LOAD_FAILED: Failed to load basket for reopening")
		return errors.Wrap(err, "baskets reopen")
	}

	logger.Info().
		Str("current_status", string(basket.Status)).
		Msg("BASKETS_REOPEN_LOADED: Basket loaded for reopening")

	event, err := basket.Reopen(reopen.Reason)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_REOPEN_DOMAIN_FAILED: Domain reopen operation failed")
		return errors.Wrap(err, "baskets reopen")
	}

	logger.Info().
		Str("event_name", event.EventName()).
		Msg("BASKETS_REOPEN_EVENT_CREATED: BasketReopened event created")

	if err = a.baskets.Save(ctx, basket); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_REOPEN_SAVE_FAILED: Failed to save reopened basket")
		return errors.Wrap(err, "basket reopen")
	}

	logger.Debug().
		Str("new_basket_status", string(basket.Status)).
		Msg("BASKETS_REOPEN_SAVED: Basket saved with open status")

	if err = a.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_REOPEN_PUBLISH_FAILED: Failed to publish BasketReopened event")
		return err
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Str("event_id", event.ID()).
		Msg("BASKETS_REOPEN_SUCCESS: Basket reopened successfully")

	span.AddEvent("basket_reopened", trace.WithAttributes(
		attribute.String("basket_id", reopen.ID),
		attribute.String("reason", reopen.Reason),
		attribute.String("event_id", event.ID()),
	))

	return nil
}

func (a Application) AddItem(ctx context.Context, add AddItem) error {
	span := trace.SpanFromContext(ctx)
	logger := a.logger.With().
		Str("operation", "AddItem").
		Str("basket_id", add.ID).
		Str("product_id", add.ProductID).
		Int64("quantity", add.Quantity).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("BASKETS_ADD_ITEM_BEGIN: Adding item to basket")

	basket, err := a.baskets.Load(ctx, add.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_ADD_ITEM_LOAD_FAILED: Failed to load basket")
		return err
	}

	logger.Debug().
		Str("basket_status", string(basket.Status)).
		Int("existing_items", len(basket.Items)).
		Msg("BASKETS_ADD_ITEM_BASKET_LOADED: Basket loaded for item addition")

	product, err := a.products.Find(ctx, add.ProductID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_ADD_ITEM_PRODUCT_NOT_FOUND: Product not found")
		return err
	}

	logger.Debug().
		Str("product_name", product.Name).
		Int64("product_price", product.BasePrice).
		Str("seller_id", product.UserSellerID).
		Msg("BASKETS_ADD_ITEM_PRODUCT_LOADED: Product details loaded")

	userSeller, err := a.users.Find(ctx, product.UserSellerID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_ADD_ITEM_SELLER_NOT_FOUND: Seller not found")
		return err
	}

	logger.Debug().
		Str("seller_name", userSeller.FirstName).
		Msg("BASKETS_ADD_ITEM_SELLER_LOADED: Seller details loaded")

	event, err := basket.AddItem(userSeller, product, add.Quantity)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_ADD_ITEM_DOMAIN_FAILED: Domain add item operation failed")
		return err
	}

	if err = a.baskets.Save(ctx, basket); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_ADD_ITEM_SAVE_FAILED: Failed to save basket with new item")
		return err
	}

	if err = a.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_ADD_ITEM_PUBLISH_FAILED: Failed to publish item added event")
		return err
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Int("new_items_count", len(basket.Items)).
		Int64("new_total", basket.TotalAmount()).
		Msg("BASKETS_ADD_ITEM_SUCCESS: Item added successfully")

	return nil
}

func (a Application) RemoveItem(ctx context.Context, remove RemoveItem) error {
	span := trace.SpanFromContext(ctx)
	logger := a.logger.With().
		Str("operation", "RemoveItem").
		Str("basket_id", remove.ID).
		Str("product_id", remove.ProductID).
		Int64("quantity", remove.Quantity).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("BASKETS_REMOVE_ITEM_BEGIN: Removing item from basket")

	product, err := a.products.Find(ctx, remove.ProductID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_REMOVE_ITEM_PRODUCT_NOT_FOUND: Product not found")
		return err
	}

	basket, err := a.baskets.Load(ctx, remove.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_REMOVE_ITEM_LOAD_FAILED: Failed to load basket")
		return err
	}

	logger.Debug().
		Int("items_before", len(basket.Items)).
		Int64("total_before", basket.TotalAmount()).
		Msg("BASKETS_REMOVE_ITEM_BEFORE_STATE: Basket state before removal")

	event, err := basket.RemoveItem(product, remove.Quantity)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_REMOVE_ITEM_DOMAIN_FAILED: Domain remove item operation failed")
		return err
	}

	if err = a.baskets.Save(ctx, basket); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_REMOVE_ITEM_SAVE_FAILED: Failed to save basket after item removal")
		return err
	}

	if err = a.publisher.Publish(ctx, event); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_REMOVE_ITEM_PUBLISH_FAILED: Failed to publish item removed event")
		return err
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Int("items_after", len(basket.Items)).
		Int64("total_after", basket.TotalAmount()).
		Msg("BASKETS_REMOVE_ITEM_SUCCESS: Item removed successfully")

	return nil
}

func (a Application) GetBasket(ctx context.Context, get GetBasket) (*domain.Basket, error) {
	span := trace.SpanFromContext(ctx)
	logger := a.logger.With().
		Str("operation", "GetBasket").
		Str("basket_id", get.ID).
		Str("user_customer_id", get.UserCustomerID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Debug().Msg("BASKETS_GET_BASKET_BEGIN: Retrieving basket")

	basket, err := a.baskets.Load(ctx, get.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_GET_BASKET_FAILED: Failed to load basket")
		return nil, err
	}

	logger.Debug().
		Str("basket_status", string(basket.Status)).
		Int("items_count", len(basket.Items)).
		Int64("total_amount", basket.TotalAmount()).
		Dur("duration_ms", time.Since(startTime)).
		Msg("BASKETS_GET_BASKET_SUCCESS: Basket retrieved successfully")

	return basket, nil
}

func (a Application) GetCurrentBasket(ctx context.Context, get GetCurrentBasket) (*domain.CatalogBasket, error) {
	span := trace.SpanFromContext(ctx)
	logger := a.logger.With().
		Str("operation", "GetCurrentBasket").
		Str("user_customer_id", get.UserCustomerID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Debug().Msg("BASKETS_GET_CURRENT_BASKET_BEGIN: Retrieving current basket")

	basket, err := a.catalog.Find(ctx, get.UserCustomerID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_GET_CURRENT_BASKET_FAILED: Failed to find current basket")
		return nil, err
	}

	logger.Debug().
		Dur("duration_ms", time.Since(startTime)).
		Msg("BASKETS_GET_CURRENT_BASKET_SUCCESS: Current basket retrieved")

	return basket, nil
}

func (a Application) GetTotalBasket(ctx context.Context, get GetTotalBasket) (int64, error) {
	span := trace.SpanFromContext(ctx)
	logger := a.logger.With().
		Str("operation", "GetTotalBasket").
		Str("basket_id", get.ID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Debug().Msg("BASKETS_GET_TOTAL_BEGIN: Calculating basket total")

	basket, err := a.baskets.Load(ctx, get.ID)
	if err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("BASKETS_GET_TOTAL_LOAD_FAILED: Failed to load basket for total calculation")
		return 0, err
	}

	total := basket.TotalAmount()
	logger.Debug().
		Int64("total_amount", total).
		Int("items_count", len(basket.Items)).
		Dur("duration_ms", time.Since(startTime)).
		Msg("BASKETS_GET_TOTAL_SUCCESS: Basket total calculated")

	return total, nil
}
