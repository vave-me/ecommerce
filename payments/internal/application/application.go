package application

import (
	"context"
	"fmt"
	"middleman/internal/ddd"
	"middleman/payments/internal/models"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stackus/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	// Payment-level commands (UNCHANGED)
	AuthorizePaymentCommand struct {
		PaymentID      string
		BasketID       string
		UserCustomerID string
		OrderID        string
		Amount         int64
	}
	ConfirmPaymentCommand struct {
		ID        string
		PaymentID string
	}
	HandleWebhookCommand struct {
		PaymentID string
	}

	// Invoice-level commands (UNCHANGED)
	CreateInvoiceCommand struct {
		ID        string
		InvoiceID string
		OrderID   string
		PaymentID string // optional association
		Amount    int64
	}
	AdjustInvoiceCommand struct {
		InvoiceID string
		Amount    int64
	}
	PayInvoiceCommand struct {
		InvoiceID string
		PayAmount int64 // how much the user is paying this time
	}
	CancelInvoiceCommand struct {
		InvoiceID string
	}

	// Recurring Payment commands (UNCHANGED)
	SetupRecurringPaymentCommand struct {
		PlanID         string
		UserCustomerID string
		Amount         int64
		Frequency      string // e.g. "monthly", "weekly", etc.
		StartDate      time.Time
	}
	ChargeRecurringPaymentCommand struct {
		PlanID string
	}
)

// PaymentDomain is the minimal interface that the webhook route needs.
// For example, only ConfirmPayment is required for "payment_intent.succeeded".
//
// ADDED: This interface is separate from the full App interface.
type PaymentDomain interface {
	ConfirmPayment(ctx context.Context, cmd ConfirmPaymentCommand) error
	AuthorizePayment(ctx context.Context, cmd AuthorizePaymentCommand) error
	// add other domain methods if needed, e.g. CancelPayment
}

// App interface is the full feature set for your Payment service
// (UNCHANGED, except we note it extends the same ConfirmPayment signature)
type App interface {
	// Payment commands
	AuthorizePayment(ctx context.Context, cmd AuthorizePaymentCommand) error
	ConfirmPayment(ctx context.Context, cmd ConfirmPaymentCommand) error
	HandleWebhook(ctx context.Context, cmd HandleWebhookCommand) error

	// Invoice commands
	CreateInvoice(ctx context.Context, cmd CreateInvoiceCommand) error
	AdjustInvoice(ctx context.Context, cmd AdjustInvoiceCommand) error
	PayInvoice(ctx context.Context, cmd PayInvoiceCommand) error
	CancelInvoice(ctx context.Context, cmd CancelInvoiceCommand) error

	// Recurring Payment commands
	SetupRecurringPayment(ctx context.Context, cmd SetupRecurringPaymentCommand) error
	ChargeRecurringPayment(ctx context.Context, cmd ChargeRecurringPaymentCommand) error
}

// Implementation of both interfaces:

// Application satisfies both App and PaymentDomain, because it has ConfirmPayment.
type Application struct {
	invoices      InvoiceRepository
	payments      PaymentRepository
	recurringRepo RecurringRepository
	publisher     ddd.EventPublisher[ddd.Event]
}

// Ensure it implements App fully
var _ App = (*Application)(nil)

// ALSO ensure it implements PaymentDomain for minimal "webhook" usage
var _ PaymentDomain = (*Application)(nil)

// New constructor (UNCHANGED)
func New(
	invoices InvoiceRepository,
	payments PaymentRepository,
	recurring RecurringRepository,
	publisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		invoices:      invoices,
		payments:      payments,
		recurringRepo: recurring,
		publisher:     publisher,
	}
}

// AuthorizePayment (UPDATED with comprehensive logging)
func (a *Application) AuthorizePayment(ctx context.Context, cmd AuthorizePaymentCommand) error {
	span := trace.SpanFromContext(ctx)
	logger := log.With().
		Str("service", "payments").
		Str("operation", "AuthorizePayment").
		Str("payment_id", cmd.PaymentID).
		Str("basket_id", cmd.BasketID).
		Str("user_customer_id", cmd.UserCustomerID).
		Str("order_id", cmd.OrderID).
		Int64("amount", cmd.Amount).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("PAYMENTS_AUTHORIZE_BEGIN: Starting payment authorization - SAGA STEP 3")

	payment := &models.Payment{
		ID:              cmd.PaymentID,
		UserCustomerID:  cmd.UserCustomerID,
		Amount:          cmd.Amount,
		Status:          models.PaymentStatusAuthorized,
		PaymentIntentID: cmd.PaymentID,
		OrderID:         cmd.OrderID,
	}

	logger.Debug().
		Str("payment_status", string(payment.Status)).
		Str("payment_intent_id", payment.PaymentIntentID).
		Msg("PAYMENTS_AUTHORIZE_PAYMENT_CREATED: Payment object created")

	// Save the payment record
	if err := a.payments.Save(ctx, payment); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("PAYMENTS_AUTHORIZE_SAVE_FAILED: Failed to save payment record")
		return err
	}

	logger.Debug().Msg("PAYMENTS_AUTHORIZE_SAVED: Payment record saved to repository")

	// Publish PaymentAuthorizedEvent
	evt := ddd.NewEvent(models.PaymentAuthorizedEvent, &models.PaymentAuthorized{
		PaymentID: cmd.PaymentID,
		Amount:    cmd.Amount,
		UserID:    cmd.UserCustomerID,
		BasketID:  cmd.BasketID,
	})

	logger.Info().
		Str("event_name", evt.EventName()).
		Str("event_id", evt.ID()).
		Msg("PAYMENTS_AUTHORIZE_EVENT_CREATED: PaymentAuthorized event created")

	if err := a.publisher.Publish(ctx, evt); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Str("event_name", evt.EventName()).
			Msg("PAYMENTS_AUTHORIZE_PUBLISH_FAILED: Failed to publish PaymentAuthorized event")
		return err
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Str("event_id", evt.ID()).
		Msg("PAYMENTS_AUTHORIZE_SUCCESS: Payment authorized and event published - SAGA STEP 3 COMPLETE")

	span.AddEvent("payment_authorized", trace.WithAttributes(
		attribute.String("payment_id", cmd.PaymentID),
		attribute.String("user_customer_id", cmd.UserCustomerID),
		attribute.String("order_id", cmd.OrderID),
		attribute.Int64("amount", cmd.Amount),
		attribute.String("event_id", evt.ID()),
	))

	return nil
}

// ConfirmPayment (UPDATED with comprehensive logging)
func (a *Application) ConfirmPayment(ctx context.Context, cmd ConfirmPaymentCommand) error {
	span := trace.SpanFromContext(ctx)
	logger := log.With().
		Str("service", "payments").
		Str("operation", "ConfirmPayment").
		Str("payment_id", cmd.PaymentID).
		Str("correlation_id", span.SpanContext().TraceID().String()).
		Logger()

	startTime := time.Now()
	logger.Info().Msg("PAYMENTS_CONFIRM_BEGIN: Starting payment confirmation - SAGA STEP 4")

	payment, err := a.payments.Find(ctx, cmd.PaymentID)
	if err != nil || payment == nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("PAYMENTS_CONFIRM_NOT_FOUND: Payment not found for confirmation")
		return errors.Wrap(errors.ErrNotFound, "payment not found for confirm")
	}

	logger.Debug().
		Str("current_status", string(payment.Status)).
		Int64("amount", payment.Amount).
		Str("user_customer_id", payment.UserCustomerID).
		Msg("PAYMENTS_CONFIRM_LOADED: Payment loaded for confirmation")

	previousStatus := payment.Status
	payment.Status = models.PaymentStatusConfirmed

	logger.Info().
		Str("previous_status", string(previousStatus)).
		Str("new_status", string(payment.Status)).
		Msg("PAYMENTS_CONFIRM_STATUS_UPDATED: Payment status updated to confirmed")

	if err = a.payments.Save(ctx, payment); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Msg("PAYMENTS_CONFIRM_SAVE_FAILED: Failed to save confirmed payment")
		fmt.Printf("ERROR IN THE SAVE %s", err.Error())
		return err
	}

	logger.Debug().Msg("PAYMENTS_CONFIRM_SAVED: Confirmed payment saved to repository")

	// Publish PaymentConfirmedEvent
	evt := ddd.NewEvent(models.PaymentConfirmedEvent, &models.PaymentConfirmed{
		PaymentID: cmd.PaymentID,
		Amount:    payment.Amount,
	})

	logger.Info().
		Str("event_name", evt.EventName()).
		Str("event_id", evt.ID()).
		Msg("PAYMENTS_CONFIRM_EVENT_CREATED: PaymentConfirmed event created")

	if err := a.publisher.Publish(ctx, evt); err != nil {
		logger.Error().Err(err).
			Dur("duration_ms", time.Since(startTime)).
			Str("event_name", evt.EventName()).
			Msg("PAYMENTS_CONFIRM_PUBLISH_FAILED: Failed to publish PaymentConfirmed event")
		return err
	}

	logger.Info().
		Dur("duration_ms", time.Since(startTime)).
		Str("event_id", evt.ID()).
		Msg("PAYMENTS_CONFIRM_SUCCESS: Payment confirmed and event published - SAGA STEP 4 COMPLETE")

	span.AddEvent("payment_confirmed", trace.WithAttributes(
		attribute.String("payment_id", cmd.PaymentID),
		attribute.Int64("amount", payment.Amount),
		attribute.String("event_id", evt.ID()),
		attribute.String("previous_status", string(previousStatus)),
		attribute.String("new_status", string(payment.Status)),
	))

	return nil
}

// HandleWebhook (UNCHANGED)
func (a *Application) HandleWebhook(ctx context.Context, cmd HandleWebhookCommand) error {
	payment, err := a.payments.Find(ctx, cmd.PaymentID)
	if err != nil || payment == nil {
		return errors.Wrap(errors.ErrNotFound, "payment not found in webhook")
	}

	// Possibly update PaymentMethod or handle gateway callbacks
	if err := a.payments.SavePaymentMethod(ctx, payment); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save payment method from webhook")
	}
	return nil
}

// --------------- Invoice logic (UNCHANGED) ---------------
func (a *Application) CreateInvoice(ctx context.Context, cmd CreateInvoiceCommand) error {
	invoice := &models.Invoice{
		ID:      cmd.InvoiceID,
		OrderID: cmd.OrderID,
		Amount:  cmd.Amount,
		PaidAmt: 0,
		Status:  models.InvoiceIsPending,
	}
	return a.invoices.Save(ctx, invoice)
}

func (a *Application) AdjustInvoice(ctx context.Context, cmd AdjustInvoiceCommand) error {
	invoice, err := a.invoices.Find(ctx, cmd.InvoiceID)
	if err != nil {
		return err
	}
	invoice.Amount = cmd.Amount
	// Could publish "InvoiceAdjustedEvent" if needed
	return a.invoices.Update(ctx, invoice)
}

func (a *Application) PayInvoice(ctx context.Context, cmd PayInvoiceCommand) error {
	invoice, err := a.invoices.Find(ctx, cmd.InvoiceID)
	if err != nil {
		return err
	}
	if invoice.Status == models.InvoiceIsCanceled || invoice.Status == models.InvoiceIsPaid {
		return errors.Wrap(errors.ErrBadRequest, "invoice cannot be paid in current status")
	}

	originalPaid := invoice.PaidAmt
	invoice.PaidAmt += cmd.PayAmount

	if invoice.PaidAmt < invoice.Amount {
		invoice.Status = models.InvoiceIsPartially
		partialEvt := ddd.NewEvent(models.InvoicePartialPaidEvent, &models.InvoicePartialPaid{
			InvoiceID:  invoice.ID,
			OrderID:    invoice.OrderID,
			PaidAmount: cmd.PayAmount,
			Remaining:  invoice.Amount - invoice.PaidAmt,
		})
		log.Debug().Str("event", partialEvt.EventName()).Str("invoice_id", invoice.ID).Msg("Publishing domain event")
		if errPub := a.publisher.Publish(ctx, partialEvt); errPub != nil {
			return errPub
		}
	} else {
		invoice.Status = models.InvoiceIsPaid
		paidEvt := ddd.NewEvent(models.InvoicePaidEvent, &models.InvoicePaid{
			InvoiceID: invoice.ID,
			OrderID:   invoice.OrderID,
		})
		log.Debug().Str("event", paidEvt.EventName()).Str("invoice_id", invoice.ID).Msg("Publishing domain event")
		if errPub := a.publisher.Publish(ctx, paidEvt); errPub != nil {
			invoice.PaidAmt = originalPaid
			return errPub
		}
	}

	return a.invoices.Update(ctx, invoice)
}

func (a *Application) CancelInvoice(ctx context.Context, cmd CancelInvoiceCommand) error {
	invoice, err := a.invoices.Find(ctx, cmd.InvoiceID)
	if err != nil {
		return err
	}

	if invoice.Status == models.InvoiceIsPaid {
		return errors.Wrap(errors.ErrBadRequest, "cannot cancel an already paid invoice")
	}
	if invoice.Status == models.InvoiceIsCanceled {
		return errors.Wrap(errors.ErrBadRequest, "invoice already canceled")
	}

	invoice.Status = models.InvoiceIsCanceled

	canceledEvt := ddd.NewEvent(models.InvoiceCanceledEvent, &models.InvoiceCanceled{
		InvoiceID: invoice.ID,
		OrderID:   invoice.OrderID,
	})
	log.Debug().Str("event", canceledEvt.EventName()).Str("invoice_id", invoice.ID).Msg("Publishing domain event")
	if errPub := a.publisher.Publish(ctx, canceledEvt); errPub != nil {
		return errPub
	}

	return a.invoices.Update(ctx, invoice)
}

// --------------- Recurring logic (UNCHANGED) ---------------
func (a *Application) SetupRecurringPayment(ctx context.Context, cmd SetupRecurringPaymentCommand) error {
	recurringPlan := &models.RecurringPaymentPlan{
		PlanID:         cmd.PlanID,
		UserCustomerID: cmd.UserCustomerID,
		Amount:         cmd.Amount,
		Frequency:      cmd.Frequency,
		StartDate:      cmd.StartDate,
		Status:         models.RecurringPlanStatusActive,
	}
	if err := a.recurringRepo.Save(ctx, recurringPlan); err != nil {
		return err
	}

	evt := ddd.NewEvent(models.RecurringPaymentSetupEvent, &models.RecurringPaymentSetup{
		PlanID:         cmd.PlanID,
		UserCustomerID: cmd.UserCustomerID,
		Amount:         cmd.Amount,
		Frequency:      cmd.Frequency,
		StartDate:      cmd.StartDate,
	})
	log.Debug().Str("event", evt.EventName()).Str("plan_id", cmd.PlanID).Msg("Publishing domain event")
	return a.publisher.Publish(ctx, evt)
}

// ChargeRecurringPayment
func (a *Application) ChargeRecurringPayment(ctx context.Context, cmd ChargeRecurringPaymentCommand) error {
	plan, err := a.recurringRepo.Find(ctx, cmd.PlanID)
	if err != nil || plan == nil {
		return errors.Wrap(errors.ErrNotFound, "recurring plan not found")
	}
	if plan.Status != models.RecurringPlanStatusActive {
		return errors.Wrap(errors.ErrBadRequest, "recurring plan is not active")
	}

	payment := &models.Payment{
		ID:             "auto-" + cmd.PlanID + "-" + time.Now().Format("20060102_150405"),
		UserCustomerID: plan.UserCustomerID,
		Amount:         plan.Amount,
		Status:         models.PaymentStatusAuthorized,
	}
	if err := a.payments.Save(ctx, payment); err != nil {
		return errors.Wrap(err, "failed to create payment for recurring plan")
	}

	now := time.Now()
	plan.LastChargedAt = now
	plan.NextDueDate = computeNextDueDate(plan.Frequency, now)
	if err := a.recurringRepo.Save(ctx, plan); err != nil {
		return err
	}

	evt := ddd.NewEvent(models.RecurringPaymentChargedEvent, &models.RecurringPaymentCharged{
		PlanID:    plan.PlanID,
		PaymentID: payment.ID,
		Amount:    plan.Amount,
	})
	log.Debug().Str("event", evt.EventName()).Str("plan_id", plan.PlanID).Msg("Publishing domain event")
	return a.publisher.Publish(ctx, evt)
}

func computeNextDueDate(freq string, lastCharged time.Time) time.Time {
	switch freq {
	case "monthly":
		return lastCharged.AddDate(0, 1, 0)
	case "weekly":
		return lastCharged.AddDate(0, 0, 7)
	default:
		return lastCharged.AddDate(0, 1, 0)
	}
}
