package internal

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"middleman/cosec/internal/models"
	"middleman/internal/ddd"
	"middleman/internal/sec"
	"middleman/ordering/orderingpb"
	"middleman/payments/paymentspb"
	"middleman/products/productspb"
	"middleman/shipping/shippingpb"

	"github.com/google/uuid"
	"github.com/stackus/errors"
)

// CheckoutSaga orchestrates: BasketCheckedOut ➜ CreateOrder ➜ ReserveStock ➜ AuthorizePayment ➜ ConfirmPayment ➜ ApproveOrder ➜ CreateShipment

const CheckoutSagaName = "cosec.Checkout"
const CheckoutSagaReplyChannel = "middleman.cosec.replies.Checkout"

type checkoutSaga struct {
	sec.Saga[*models.CheckoutData]
	logger zerolog.Logger
}

func NewCheckoutSaga() sec.Saga[*models.CheckoutData] {
	s := checkoutSaga{
		Saga:   sec.NewSaga[*models.CheckoutData](CheckoutSagaName, CheckoutSagaReplyChannel),
		logger: log.With().Str("service", "cosec").Str("saga", "checkout").Logger(),
	}

	// Compensation for the entire saga - reopen basket first, then reject order
	s.AddStep().Compensation(s.reopenBasket)
	s.AddStep().Compensation(s.rejectOrder)

	// 1. CreateOrder - with proper ID generation
	s.AddStep().Action(s.createOrder).OnActionReply(orderingpb.OrderCreatedEvent, s.onOrderCreated)

	// 2. ReserveStock - for ALL items with compensation
	s.AddStep().Action(s.reserveStock).Compensation(s.releaseStock)

	// 3. AuthorizePayment - pre-auth with compensation
	s.AddStep().Action(s.authorizePayment).OnActionReply(paymentspb.PaymentAuthorizedEvent, s.onPaymentAuthorized).Compensation(s.cancelPayment)

	// 4. ConfirmPayment - capture funds
	s.AddStep().Action(s.confirmPayment).Compensation(s.refundPayment)

	// 5. ApproveOrder - finalize order
	s.AddStep().Action(s.approveOrder)

	// 6. CreateShipment - initiate shipping
	s.AddStep().Action(s.createShipment)

	return s
}

// --- Compensations --------------------------------------------------------

func (s checkoutSaga) rejectOrder(ctx context.Context, d *models.CheckoutData) (string, ddd.Command, error) {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "rejectOrder").
		Str("basket_id", d.BasketID).
		Str("order_id", d.OrderID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	logger.Info().Msg("COSEC_SAGA_REJECT_ORDER_BEGIN: Starting order rejection compensation")

	if d.OrderID == "" {
		logger.Debug().Msg("COSEC_SAGA_REJECT_ORDER_SKIP: No order ID, skipping rejection")
		return "", nil, nil // nothing created yet
	}

	cmd := ddd.NewCommand(orderingpb.RejectOrderCommand, &orderingpb.RejectOrder{Id: d.OrderID})

	logger.Info().
		Str("command", orderingpb.RejectOrderCommand).
		Str("order_id", d.OrderID).
		Msg("COSEC_SAGA_REJECT_ORDER_COMMAND: Sending reject order command")

	return orderingpb.CommandChannel, cmd, nil
}

func (s checkoutSaga) releaseStock(ctx context.Context, d *models.CheckoutData) (string, ddd.Command, error) {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "releaseStock").
		Str("basket_id", d.BasketID).
		Str("order_id", d.OrderID).
		Int("items_count", len(d.Items)).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	logger.Info().Msg("COSEC_SAGA_RELEASE_STOCK_BEGIN: Starting stock release compensation")

	if len(d.Items) == 0 {
		logger.Debug().Msg("COSEC_SAGA_RELEASE_STOCK_SKIP: No items to release")
		return "", nil, nil
	}

	// Build batch release request for ALL items
	releaseItems := make([]*productspb.ReleaseProductItem, 0, len(d.Items))
	for i, item := range d.Items {
		logger.Debug().
			Int("item_index", i).
			Str("product_id", item.ProductID).
			Int64("quantity", item.Quantity).
			Msg("COSEC_SAGA_RELEASE_STOCK_ITEM: Adding item to batch release")

		releaseItems = append(releaseItems, &productspb.ReleaseProductItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	// Send batch release command
	cmd := ddd.NewCommand(productspb.ReleaseProductsCommand, &productspb.ReleaseProducts{
		OrderId: d.OrderID,
		Items:   releaseItems,
	})

	logger.Info().
		Str("command", productspb.ReleaseProductsCommand).
		Int("items_count", len(releaseItems)).
		Msg("COSEC_SAGA_RELEASE_STOCK_COMMAND: Sending batch release stock command")

	return productspb.CommandChannel, cmd, nil
}

func (s checkoutSaga) cancelPayment(ctx context.Context, d *models.CheckoutData) (string, ddd.Command, error) {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "cancelPayment").
		Str("basket_id", d.BasketID).
		Str("payment_id", d.PaymentID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	logger.Info().Msg("COSEC_SAGA_CANCEL_PAYMENT_BEGIN: Starting payment cancellation compensation")

	if d.PaymentID == "" {
		logger.Debug().Msg("COSEC_SAGA_CANCEL_PAYMENT_SKIP: No payment ID, skipping cancellation")
		return "", nil, nil
	}

	// Use existing command - payment service should handle cancellation
	cmd := ddd.NewCommand(paymentspb.AuthorizePaymentCommand, &paymentspb.AuthorizePaymentRequest{
		UserCustomerId: d.UserID,
		Amount:         -d.Total, // Negative amount for cancellation
	})

	logger.Info().
		Str("command", paymentspb.AuthorizePaymentCommand).
		Int64("amount", -d.Total).
		Msg("COSEC_SAGA_CANCEL_PAYMENT_COMMAND: Sending payment cancellation command")

	return paymentspb.CommandChannel, cmd, nil
}

func (s checkoutSaga) refundPayment(ctx context.Context, d *models.CheckoutData) (string, ddd.Command, error) {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "refundPayment").
		Str("basket_id", d.BasketID).
		Str("payment_id", d.PaymentID).
		Int64("amount", d.Total).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	logger.Info().Msg("COSEC_SAGA_REFUND_PAYMENT_BEGIN: Starting payment refund compensation")

	if d.PaymentID == "" {
		logger.Debug().Msg("COSEC_SAGA_REFUND_PAYMENT_SKIP: No payment ID, skipping refund")
		return "", nil, nil
	}

	// Use confirm payment with negative amount for refund
	cmd := ddd.NewCommand(paymentspb.ConfirmPaymentCommand, &paymentspb.ConfirmPayment{
		Id:      d.PaymentID,
		Amount:  -d.Total, // Negative amount for refund
		OrderId: d.OrderID,
	})

	logger.Info().
		Str("command", paymentspb.ConfirmPaymentCommand).
		Int64("refund_amount", -d.Total).
		Msg("COSEC_SAGA_REFUND_PAYMENT_COMMAND: Sending payment refund command")

	return paymentspb.CommandChannel, cmd, nil
}

func (s checkoutSaga) reopenBasket(ctx context.Context, d *models.CheckoutData) (string, ddd.Command, error) {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "reopenBasket").
		Str("basket_id", d.BasketID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	logger.Info().Msg("COSEC_SAGA_REOPEN_BASKET_BEGIN: Reopening basket after failed checkout")

	// Since baskets service doesn't use commands, we'll need to call the gRPC method directly
	// For now, log that this needs to be done
	logger.Warn().
		Str("basket_id", d.BasketID).
		Msg("COSEC_SAGA_REOPEN_BASKET_TODO: Basket reopening needs to be implemented via direct gRPC call")

	// Return nil to indicate no command to send
	return "", nil, nil
}

// --- Actions --------------------------------------------------------------

// Step 1: CreateOrder - with proper ID generation
func (s checkoutSaga) createOrder(ctx context.Context, d *models.CheckoutData) (string, ddd.Command, error) {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "createOrder").
		Str("basket_id", d.BasketID).
		Str("user_id", d.UserID).
		Int("items_count", len(d.Items)).
		Int64("total", d.Total).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("COSEC_SAGA_CREATE_ORDER_BEGIN: Starting order creation - SAGA STEP 1")

	// Generate OrderID if not present
	if d.OrderID == "" {
		d.OrderID = uuid.New().String()
		logger.Info().
			Str("generated_order_id", d.OrderID).
			Msg("COSEC_SAGA_CREATE_ORDER_ID_GENERATED: Generated new order ID")
	}

	// Convert CheckoutItems to protobuf Items for the command
	orderItems := make([]*orderingpb.Item, len(d.Items))
	for i, item := range d.Items {
		orderItems[i] = &orderingpb.Item{
			UserSellerId:   item.SellerID,
			ProductId:      item.ProductID,
			UserSellerName: item.SellerName,
			ProductName:    item.ProductName,
			Price:          item.Price,
			Quantity:       item.Quantity,
		}
		logger.Debug().
			Int("item_index", i).
			Str("product_id", item.ProductID).
			Str("seller_id", item.SellerID).
			Int64("quantity", item.Quantity).
			Int64("price", item.Price).
			Msg("COSEC_SAGA_CREATE_ORDER_ITEM: Converting item to order item")
	}

	cmd := ddd.NewCommand(orderingpb.CreateOrderCommand, &orderingpb.CreateOrderRequest{
		Id:             d.OrderID,
		UserCustomerId: d.UserID,
		Items:          orderItems,
		BasketId:       d.BasketID,
		PaymentIntent:  d.PaymentID,
	})

	logger.Info().
		Str("command", orderingpb.CreateOrderCommand).
		Str("order_id", d.OrderID).
		Str("user_customer_id", d.UserID).
		Dur("duration_ms", time.Since(startTime)).
		Msg("COSEC_SAGA_CREATE_ORDER_COMMAND: Sending create order command - SAGA STEP 1")

	span.AddEvent("saga_create_order", trace.WithAttributes(
		attribute.String("order_id", d.OrderID),
		attribute.String("basket_id", d.BasketID),
		attribute.String("user_id", d.UserID),
		attribute.Int("items_count", len(d.Items)),
		attribute.Int64("total", d.Total),
	))

	return orderingpb.CommandChannel, cmd, nil
}

func (s checkoutSaga) onOrderCreated(ctx context.Context, d *models.CheckoutData, reply ddd.Reply) error {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "onOrderCreated").
		Str("basket_id", d.BasketID).
		Str("expected_order_id", d.OrderID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	logger.Info().Msg("COSEC_SAGA_ORDER_CREATED_REPLY: Processing OrderCreated reply")

	payload := reply.Payload().(*orderingpb.OrderCreated)

	logger.Debug().
		Str("reply_order_id", payload.GetId()).
		Str("expected_order_id", d.OrderID).
		Msg("COSEC_SAGA_ORDER_CREATED_VALIDATION: Validating order ID match")

	// Confirm OrderID is set correctly
	if d.OrderID != payload.GetId() {
		logger.Error().
			Str("expected_order_id", d.OrderID).
			Str("received_order_id", payload.GetId()).
			Msg("COSEC_SAGA_ORDER_CREATED_MISMATCH: Order ID mismatch detected")
		return errors.ErrInternal.Msgf("OrderID mismatch: expected %s, got %s", d.OrderID, payload.GetId())
	}

	logger.Info().
		Str("order_id", d.OrderID).
		Msg("COSEC_SAGA_ORDER_CREATED_SUCCESS: Order creation confirmed - PROCEEDING TO STEP 2")

	return nil
}

// Step 2: ReserveStock - for ALL items
func (s checkoutSaga) reserveStock(ctx context.Context, d *models.CheckoutData) (string, ddd.Command, error) {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "reserveStock").
		Str("basket_id", d.BasketID).
		Str("order_id", d.OrderID).
		Int("items_count", len(d.Items)).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("COSEC_SAGA_RESERVE_STOCK_BEGIN: Starting stock reservation - SAGA STEP 2")

	if len(d.Items) == 0 {
		logger.Error().Msg("COSEC_SAGA_RESERVE_STOCK_NO_ITEMS: No items to reserve")
		return "", nil, errors.ErrBadRequest.Msg("no items to reserve")
	}

	// Build batch reservation request for ALL items
	reserveItems := make([]*productspb.ReserveProductItem, 0, len(d.Items))
	for i, item := range d.Items {
		logger.Debug().
			Int("item_index", i).
			Str("product_id", item.ProductID).
			Int64("quantity", item.Quantity).
			Msg("COSEC_SAGA_RESERVE_STOCK_ITEM: Adding item to batch reservation")

		reserveItems = append(reserveItems, &productspb.ReserveProductItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	// Send batch reservation command
	cmd := ddd.NewCommand(productspb.ReserveProductsCommand, &productspb.ReserveProducts{
		OrderId: d.OrderID,
		Items:   reserveItems,
	})

	logger.Info().
		Str("command", productspb.ReserveProductsCommand).
		Int("items_count", len(reserveItems)).
		Dur("duration_ms", time.Since(startTime)).
		Msg("COSEC_SAGA_RESERVE_STOCK_COMMAND: Sending batch reserve stock command - SAGA STEP 2")

	return productspb.CommandChannel, cmd, nil
}

// Step 3: AuthorizePayment - pre-authorization (SKIP if payment intent already exists)
func (s checkoutSaga) authorizePayment(ctx context.Context, d *models.CheckoutData) (string, ddd.Command, error) {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "authorizePayment").
		Str("basket_id", d.BasketID).
		Str("order_id", d.OrderID).
		Str("user_id", d.UserID).
		Str("payment_id", d.PaymentID).
		Int64("amount", d.Total).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("COSEC_SAGA_AUTHORIZE_PAYMENT_BEGIN: Starting payment authorization - SAGA STEP 3")

	if d.OrderID == "" {
		logger.Error().Msg("COSEC_SAGA_AUTHORIZE_PAYMENT_NO_ORDER: Order ID required for payment authorization")
		return "", nil, errors.ErrBadRequest.Msg("orderID required for payment authorization")
	}

	// SKIP AUTHORIZATION if payment intent already exists (from frontend)
	if d.PaymentID != "" {
		logger.Info().
			Str("existing_payment_id", d.PaymentID).
			Dur("duration_ms", time.Since(startTime)).
			Msg("COSEC_SAGA_AUTHORIZE_PAYMENT_SKIP: Payment intent already exists, skipping authorization - PROCEEDING TO STEP 4")

		// Return nil to skip this step and proceed to next step
		return "", nil, nil
	}

	// Only create new payment intent if none exists
	cmd := ddd.NewCommand(paymentspb.AuthorizePaymentCommand, &paymentspb.AuthorizePaymentRequest{
		UserCustomerId: d.UserID,
		Amount:         d.Total,
		OrderId:        d.OrderID,
		BasketId:       d.BasketID,
	})

	logger.Info().
		Str("command", paymentspb.AuthorizePaymentCommand).
		Str("user_customer_id", d.UserID).
		Int64("amount", d.Total).
		Dur("duration_ms", time.Since(startTime)).
		Msg("COSEC_SAGA_AUTHORIZE_PAYMENT_COMMAND: Sending authorize payment command - SAGA STEP 3")

	span.AddEvent("saga_authorize_payment", trace.WithAttributes(
		attribute.String("order_id", d.OrderID),
		attribute.String("user_id", d.UserID),
		attribute.Int64("amount", d.Total),
	))

	return paymentspb.CommandChannel, cmd, nil
}

func (s checkoutSaga) onPaymentAuthorized(ctx context.Context, d *models.CheckoutData, reply ddd.Reply) error {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "onPaymentAuthorized").
		Str("basket_id", d.BasketID).
		Str("order_id", d.OrderID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	logger.Info().Msg("COSEC_SAGA_PAYMENT_AUTHORIZED_REPLY: Processing PaymentAuthorized reply")

	p := reply.Payload().(*paymentspb.PaymentAuthorized)
	d.PaymentID = p.GetPaymentIntentId()

	logger.Info().
		Str("payment_id", d.PaymentID).
		Int64("amount", p.GetAmount()).
		Msg("COSEC_SAGA_PAYMENT_AUTHORIZED_SUCCESS: Payment authorization confirmed - PROCEEDING TO STEP 4")

	span.AddEvent("saga_payment_authorized", trace.WithAttributes(
		attribute.String("payment_id", d.PaymentID),
		attribute.Int64("amount", p.GetAmount()),
	))

	return nil
}

// Step 4: ConfirmPayment - capture funds
func (s checkoutSaga) confirmPayment(ctx context.Context, d *models.CheckoutData) (string, ddd.Command, error) {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "confirmPayment").
		Str("basket_id", d.BasketID).
		Str("order_id", d.OrderID).
		Str("payment_id", d.PaymentID).
		Int64("amount", d.Total).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("COSEC_SAGA_CONFIRM_PAYMENT_BEGIN: Starting payment confirmation - SAGA STEP 4")

	if d.PaymentID == "" {
		logger.Error().Msg("COSEC_SAGA_CONFIRM_PAYMENT_NO_PAYMENT: Payment ID required for payment confirmation")
		return "", nil, errors.ErrBadRequest.Msg("paymentID required for payment confirmation")
	}

	cmd := ddd.NewCommand(paymentspb.ConfirmPaymentCommand, &paymentspb.ConfirmPayment{
		Id:      d.PaymentID,
		Amount:  d.Total,
		OrderId: d.OrderID,
	})

	logger.Info().
		Str("command", paymentspb.ConfirmPaymentCommand).
		Str("payment_id", d.PaymentID).
		Int64("amount", d.Total).
		Dur("duration_ms", time.Since(startTime)).
		Msg("COSEC_SAGA_CONFIRM_PAYMENT_COMMAND: Sending confirm payment command - SAGA STEP 4")

	return paymentspb.CommandChannel, cmd, nil
}

// Step 5: ApproveOrder - finalize order
func (s checkoutSaga) approveOrder(ctx context.Context, d *models.CheckoutData) (string, ddd.Command, error) {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "approveOrder").
		Str("basket_id", d.BasketID).
		Str("order_id", d.OrderID).
		Str("payment_id", d.PaymentID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("COSEC_SAGA_APPROVE_ORDER_BEGIN: Starting order approval - SAGA STEP 5")

	if d.OrderID == "" {
		logger.Error().Msg("COSEC_SAGA_APPROVE_ORDER_NO_ORDER: Order ID missing in saga data")
		return "", nil, errors.ErrBadRequest.Msg("orderID missing in saga data")
	}

	cmd := ddd.NewCommand(orderingpb.ApproveOrderCommand, &orderingpb.ApproveOrder{
		Id: d.OrderID,
	})

	logger.Info().
		Str("command", orderingpb.ApproveOrderCommand).
		Str("order_id", d.OrderID).
		Dur("duration_ms", time.Since(startTime)).
		Msg("COSEC_SAGA_APPROVE_ORDER_COMMAND: Sending approve order command - SAGA STEP 5")

	return orderingpb.CommandChannel, cmd, nil
}

// Step 6: CreateShipment - initiate shipping
func (s checkoutSaga) createShipment(ctx context.Context, d *models.CheckoutData) (string, ddd.Command, error) {
	span := trace.SpanFromContext(ctx)
	logger := s.logger.With().
		Str("operation", "createShipment").
		Str("basket_id", d.BasketID).
		Str("order_id", d.OrderID).
		Str("payment_id", d.PaymentID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("COSEC_SAGA_CREATE_SHIPMENT_BEGIN: Starting shipment creation - SAGA STEP 6 (FINAL)")

	if d.OrderID == "" {
		logger.Error().Msg("COSEC_SAGA_CREATE_SHIPMENT_NO_ORDER: Order ID missing for shipment")
		return "", nil, errors.ErrBadRequest.Msg("orderID missing for shipment")
	}

	cmd := ddd.NewCommand(shippingpb.CreateShipmentCommand, &shippingpb.CreateShipment{
		OrderId:  d.OrderID,
		BasketId: d.BasketID,
	})

	logger.Info().
		Str("command", shippingpb.CreateShipmentCommand).
		Str("order_id", d.OrderID).
		Str("basket_id", d.BasketID).
		Dur("duration_ms", time.Since(startTime)).
		Msg("COSEC_SAGA_CREATE_SHIPMENT_COMMAND: Sending create shipment command - SAGA STEP 6 (FINAL)")

	span.AddEvent("saga_create_shipment", trace.WithAttributes(
		attribute.String("order_id", d.OrderID),
		attribute.String("basket_id", d.BasketID),
	))

	return shippingpb.CommandChannel, cmd, nil
}
