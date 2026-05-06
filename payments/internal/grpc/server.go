package grpc

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"middleman/internal/di"
	"middleman/internal/errorsotel"
	stripeMiddleman "middleman/internal/stripe"
	"middleman/payments/internal/application"
	"middleman/payments/internal/constants"
	"middleman/payments/paymentspb"
)

// server implements the PaymentsServiceServer gRPC interface.
type server struct {
	app application.App
	paymentspb.UnimplementedPaymentsServiceServer
}

var _ paymentspb.PaymentsServiceServer = (*server)(nil)

// RegisterServer registers this gRPC server with the provided registrar.
func RegisterServer(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	paymentspb.RegisterPaymentsServiceServer(registrar, server{app: app})
	return nil
}

func (s server) AuthorizePayment(ctx context.Context, request *paymentspb.AuthorizePaymentRequest) (*paymentspb.AuthorizePaymentResponse, error) {
	span := trace.SpanFromContext(ctx)
	client := di.Get(ctx, constants.Stripe).(*stripeMiddleman.StripeClient)

	userID := request.GetUserCustomerId()
	amount := request.GetAmount()

	// Generate a local Payment ID (or rely on Stripe’s PaymentIntent ID)
	localPaymentID := uuid.New().String()
	span.SetAttributes(
		attribute.String("PaymentID", localPaymentID),
		attribute.String("UserCustomerID", userID),
	)

	// 1) Create PaymentIntent in Stripe
	paymentIntent, err := client.CreatePaymentIntent(amount, "eur", nil, false, "")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// 2) Create internal Payment record (PaymentStatus=AUTHORIZED)
	err = s.app.AuthorizePayment(ctx, application.AuthorizePaymentCommand{
		PaymentID:      paymentIntent.ID,
		UserCustomerID: userID,
		Amount:         amount,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &paymentspb.AuthorizePaymentResponse{
		Id:           paymentIntent.ID,
		ClientSecret: paymentIntent.ClientSecret,
	}, nil
}

func (s server) ConfirmPayment(ctx context.Context, request *paymentspb.ConfirmPaymentRequest) (*paymentspb.ConfirmPaymentResponse, error) {
	span := trace.SpanFromContext(ctx)
	client := di.Get(ctx, constants.Stripe).(*stripeMiddleman.StripeClient)

	paymentIntentID := request.GetId()
	paymentMethodID := request.GetPaymentMethodId()

	span.SetAttributes(attribute.String("PaymentID", paymentIntentID))

	// 1) Confirm the PaymentIntent in Stripe
	intent, err := client.ConfirmPaymentIntent(paymentIntentID, paymentMethodID)
	if err != nil {
		log.Error().Msgf("[ConfirmPayment] ERROR confirming PaymentIntent='%s': %v", paymentIntentID, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// 2) Update local Payment record (PaymentStatus=CONFIRMED)
	err = s.app.ConfirmPayment(ctx, application.ConfirmPaymentCommand{
		PaymentID: intent.ID,
	})
	if err != nil {
		log.Error().Msgf("[ConfirmPayment] ERROR updating local payment record for PaymentID='%s': %v", intent.ID, err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &paymentspb.ConfirmPaymentResponse{
		Id:            intent.ID,
		PaymentStatus: string(intent.Status),
	}, nil
}

func (s server) CreateInvoice(ctx context.Context, request *paymentspb.CreateInvoiceRequest) (*paymentspb.CreateInvoiceResponse, error) {
	span := trace.SpanFromContext(ctx)
	invoiceID := uuid.New().String()

	span.SetAttributes(
		attribute.String("InvoiceID", invoiceID),
		attribute.String("OrderID", request.GetOrderId()),
	)

	err := s.app.CreateInvoice(ctx, application.CreateInvoiceCommand{
		InvoiceID: invoiceID,
		OrderID:   request.GetOrderId(),
		Amount:    request.GetAmount(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &paymentspb.CreateInvoiceResponse{
		Id: invoiceID,
		// We could set CreatedAt if your app returns a timestamp
	}, nil
}

func (s server) AdjustInvoice(ctx context.Context, request *paymentspb.AdjustInvoiceRequest) (*paymentspb.AdjustInvoiceResponse, error) {
	span := trace.SpanFromContext(ctx)
	invoiceID := request.GetId()

	span.SetAttributes(attribute.String("InvoiceID", invoiceID))

	err := s.app.AdjustInvoice(ctx, application.AdjustInvoiceCommand{
		InvoiceID: invoiceID,
		Amount:    request.GetAmount(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &paymentspb.AdjustInvoiceResponse{}, nil
}

func (s server) CancelInvoice(ctx context.Context, request *paymentspb.CancelInvoiceRequest) (*paymentspb.CancelInvoiceResponse, error) {
	span := trace.SpanFromContext(ctx)
	invoiceID := request.GetId()

	span.SetAttributes(attribute.String("InvoiceID", invoiceID))

	err := s.app.CancelInvoice(ctx, application.CancelInvoiceCommand{
		InvoiceID: invoiceID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &paymentspb.CancelInvoiceResponse{}, nil
}

func (s server) CapturePayment(ctx context.Context, request *paymentspb.CapturePaymentRequest) (*paymentspb.CapturePaymentResponse, error) {
	span := trace.SpanFromContext(ctx)
	client := di.Get(ctx, constants.Stripe).(*stripeMiddleman.StripeClient)

	paymentIntentID := request.GetPaymentId()
	amountToCapture := request.GetAmountToCapture() // optional partial capture

	span.SetAttributes(attribute.String("PaymentIntentID", paymentIntentID))

	// 1) Capture PaymentIntent in Stripe
	captureParams := &stripeMiddleman.CaptureParams{}
	if amountToCapture > 0 {
		captureParams.AmountToCapture = &amountToCapture
	}
	captured, err := client.CapturePaymentIntent(paymentIntentID, captureParams)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// (Optional) Update local Payment record to PaymentStatusConfirmed or partial
	return &paymentspb.CapturePaymentResponse{
		PaymentId:     captured.ID,
		PaymentStatus: string(captured.Status),
	}, nil
}
