package stripe

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/paymentmethod"
	"github.com/stripe/stripe-go/v81/webhook"
)

// StripeClient encapsulates all Stripe API interactions.
type StripeClient struct {
	APIKey        string
	WebhookSecret string
	PaymentIntent *paymentintent.Client
}

// PaymentIntentData is a simple struct if you want to store local data about a PaymentIntent.
type PaymentIntentData struct {
	ID       string
	Amount   int64
	Currency string
	// Add more fields if needed
}

// CaptureParams allows partial or full captures of a PaymentIntent.
type CaptureParams struct {
	AmountToCapture *int64 // if nil => full capture
}

// NewStripeClient creates a new client with the provided API key and webhook secret.
func NewStripeClient(apiKey, webhookSecret string) *StripeClient {
	return &StripeClient{
		APIKey:        apiKey,
		WebhookSecret: webhookSecret,
		PaymentIntent: &paymentintent.Client{
			B:   stripe.GetBackend(stripe.APIBackend),
			Key: apiKey,
		},
	}
}

// CreatePaymentIntent creates a PaymentIntent in Stripe, returning the object with ID + ClientSecret.
// Adjust the parameters as needed for your application.
func (c *StripeClient) CreatePaymentIntent(
	amount int64,
	currency string,
	paymentMethodTypes []string, // e.g. []string{"card"}
	confirm bool, // if you want to confirm immediately
	captureMethod stripe.PaymentIntentCaptureMethod, // e.g. stripe.PaymentIntentCaptureMethodManual
) (*stripe.PaymentIntent, error) {

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
	}

	// Optionally set payment method types
	if len(paymentMethodTypes) > 0 {
		var pmt []*string
		for _, t := range paymentMethodTypes {
			pmt = append(pmt, stripe.String(t))
		}
		params.PaymentMethodTypes = pmt
	}

	// If you want to confirm right away (or handle "requires_action")
	params.Confirm = stripe.Bool(confirm)

	// If you want manual capture (for partial captures later)
	if captureMethod != "" {
		params.CaptureMethod = (*string)(&captureMethod)
	}

	intent, err := c.PaymentIntent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}
	return intent, nil
}

// RetrievePaymentIntent fetches a PaymentIntent by ID.
func (c *StripeClient) RetrievePaymentIntent(id string) (*stripe.PaymentIntent, error) {
	intent, err := c.PaymentIntent.Get(id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment intent '%s': %w", id, err)
	}
	return intent, nil
}

// ConfirmPaymentIntent finalizes the payment (capture) with a PaymentMethod.
func (c *StripeClient) ConfirmPaymentIntent(id, paymentMethodID string) (*stripe.PaymentIntent, error) {
	log.Debug().
		Str("payment_intent_id", id).
		Str("payment_method_id", paymentMethodID).
		Msg("[StripeClient] ConfirmPaymentIntent called")

	params := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String(paymentMethodID),
	}
	intent, err := c.PaymentIntent.Confirm(id, params)
	if err != nil {
		log.Error().Err(err).
			Str("payment_intent_id", id).
			Msg("[StripeClient] ERROR confirming PaymentIntent")
		return nil, fmt.Errorf("failed to confirm payment intent '%s': %w", id, err)
	}

	log.Debug().
		Str("payment_intent_id", id).
		Str("status", string(intent.Status)).
		Msg("[StripeClient] PaymentIntent confirmed successfully")
	return intent, nil
}

// CancelPaymentIntent cancels a PaymentIntent that is not yet captured or in a final state.
func (c *StripeClient) CancelPaymentIntent(id string, params *stripe.PaymentIntentCancelParams) (*stripe.PaymentIntent, error) {
	intent, err := c.PaymentIntent.Cancel(id, params)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel payment intent '%s': %w", id, err)
	}
	return intent, nil
}

// CapturePaymentIntent captures the funds of an existing PaymentIntent (partial or full).
func (c *StripeClient) CapturePaymentIntent(id string, captureParams *CaptureParams) (*stripe.PaymentIntent, error) {
	stripeParams := &stripe.PaymentIntentCaptureParams{}
	if captureParams != nil && captureParams.AmountToCapture != nil {
		stripeParams.AmountToCapture = stripe.Int64(*captureParams.AmountToCapture)
	}

	intent, err := c.PaymentIntent.Capture(id, stripeParams)
	if err != nil {
		return nil, fmt.Errorf("failed to capture payment intent '%s': %w", id, err)
	}
	return intent, nil
}

// UpdatePaymentIntent updates a PaymentIntent with new parameters (e.g., amount, metadata).
func (c *StripeClient) UpdatePaymentIntent(id string, params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	intent, err := c.PaymentIntent.Update(id, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment intent '%s': %w", id, err)
	}
	return intent, nil
}

// ------------------ PaymentMethod / Customer methods ------------------ //

// GetPaymentMethod retrieves a PaymentMethod by ID.
func (c *StripeClient) GetPaymentMethod(id string) (*stripe.PaymentMethod, error) {
	pm, err := paymentmethod.Get(id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment method '%s': %w", id, err)
	}
	return pm, nil
}

// CreatePaymentMethod creates a new PaymentMethod with card and billing details.
func (c *StripeClient) CreatePaymentMethod(
	cardParams *stripe.PaymentMethodCardParams,
	billingDetails *stripe.PaymentMethodBillingDetailsParams,
) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodParams{
		Type:           stripe.String(string(stripe.PaymentMethodTypeCard)),
		Card:           cardParams,
		BillingDetails: billingDetails,
	}
	pm, err := paymentmethod.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment method: %w", err)
	}
	return pm, nil
}

// AttachPaymentMethodToCustomer attaches a payment method to a specific customer.
func (c *StripeClient) AttachPaymentMethodToCustomer(paymentMethodID, customerID string) (*stripe.PaymentMethod, error) {
	attachParams := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}
	pm, err := paymentmethod.Attach(paymentMethodID, attachParams)
	if err != nil {
		return nil, fmt.Errorf("failed to attach payment method '%s' to customer '%s': %w", paymentMethodID, customerID, err)
	}
	return pm, nil
}

// DetachPaymentMethod detaches a payment method from its customer.
func (c *StripeClient) DetachPaymentMethod(paymentMethodID string) (*stripe.PaymentMethod, error) {
	pm, err := paymentmethod.Detach(paymentMethodID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to detach payment method '%s': %w", paymentMethodID, err)
	}
	return pm, nil
}

// ListPaymentMethods returns all the card payment methods attached to a given customer.
func (c *StripeClient) ListPaymentMethods(customerID string) ([]*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerID),
		Type:     stripe.String(string(stripe.PaymentMethodTypeCard)),
	}
	i := paymentmethod.List(params)

	var paymentMethods []*stripe.PaymentMethod
	for i.Next() {
		paymentMethods = append(paymentMethods, i.PaymentMethod())
	}
	if err := i.Err(); err != nil {
		return nil, fmt.Errorf("failed to list payment methods for customer '%s': %w", customerID, err)
	}
	return paymentMethods, nil
}

// UpdatePaymentMethod updates the billing details of a PaymentMethod.
func (c *StripeClient) UpdatePaymentMethod(
	paymentMethodID string,
	billingDetails *stripe.PaymentMethodBillingDetailsParams,
) (*stripe.PaymentMethod, error) {
	pmParams := &stripe.PaymentMethodParams{
		BillingDetails: billingDetails,
	}
	pm, err := paymentmethod.Update(paymentMethodID, pmParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment method '%s': %w", paymentMethodID, err)
	}
	return pm, nil
}

// ------------------ Webhook Parsing ------------------ //

// ParseWebhookEvent parses the Stripe event from the raw body.
// Usually called after ValidateSign if you want to do multi-step verification.
func (c *StripeClient) ParseWebhookEvent(rawBody []byte, signature, secret string) (*stripe.Event, error) {
	event, err := webhook.ConstructEventWithOptions(rawBody, signature, secret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true})
	if err != nil {
		return nil, fmt.Errorf("failed to verify webhook signature: %w", err)
	}
	return &event, nil
}

// ValidateSign demonstrates multiple validation checks. If you prefer a single standard check,
// you can omit the extra calls. The main calls used by Stripe docs are either `ConstructEvent` or `ValidatePayload`.
func (c *StripeClient) ValidateSign(rawBody []byte, signature, secret string) error {
	// 1) Validate ignoring tolerance
	err := webhook.ValidatePayloadIgnoringTolerance(rawBody, signature, secret)
	if err != nil {
		return fmt.Errorf("failed ValidatePayloadIgnoringTolerance: %w", err)
	}

	// 2) Validate with an extremely large tolerance
	err = webhook.ValidatePayloadWithTolerance(rawBody, signature, secret, 100000000000000000)
	if err != nil {
		return fmt.Errorf("failed ValidatePayloadWithTolerance: %w", err)
	}

	// 3) Validate with the default tolerance (currently 5 min).
	err = webhook.ValidatePayload(rawBody, signature, secret)
	if err != nil {
		return fmt.Errorf("failed ValidatePayload: %w", err)
	}

	return nil
}

// ExtractPaymentIntent attempts to parse a PaymentIntent from the event’s Data.Raw field.
// Useful for "payment_intent.succeeded", "payment_intent.created", "payment_intent.requires_action", etc.
func (c *StripeClient) ExtractPaymentIntent(event *stripe.Event) (*stripe.PaymentIntent, error) {
	log.Debug().
		Str("event_id", event.ID).
		Str("event_type", string(event.Type)).
		Msg("[StripeClient] ExtractPaymentIntent")

	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		log.Error().Err(err).
			Str("event_id", event.ID).
			Msg("[StripeClient] ERROR extracting PaymentIntent")
		return nil, fmt.Errorf("failed to parse PaymentIntent from event: %w", err)
	}
	return &pi, nil
}

// ExtractPaymentMethod attempts to parse a PaymentMethod from the event’s Data.Raw field.
// Useful for events like "payment_method.attached".
func (c *StripeClient) ExtractPaymentMethod(event *stripe.Event) (*stripe.PaymentMethod, error) {
	log.Debug().
		Str("event_id", event.ID).
		Str("event_type", string(event.Type)).
		Msg("[StripeClient] ExtractPaymentMethod")

	var pm stripe.PaymentMethod
	if err := json.Unmarshal(event.Data.Raw, &pm); err != nil {
		log.Error().Err(err).
			Str("event_id", event.ID).
			Msg("[StripeClient] ERROR extracting PaymentMethod")
		return nil, fmt.Errorf("failed to parse PaymentMethod from event: %w", err)
	}
	return &pm, nil
}

// ExtractCharge attempts to parse a Charge from the event’s Data.Raw field.
// Useful for events like "charge.succeeded", "charge.updated", etc.
func (c *StripeClient) ExtractCharge(event *stripe.Event) (*stripe.Charge, error) {
	log.Debug().
		Str("event_id", event.ID).
		Str("event_type", string(event.Type)).
		Msg("[StripeClient] ExtractCharge")

	var ch stripe.Charge
	if err := json.Unmarshal(event.Data.Raw, &ch); err != nil {
		log.Error().Err(err).
			Str("event_id", event.ID).
			Msg("[StripeClient] ERROR extracting Charge")
		return nil, fmt.Errorf("failed to parse Charge from event: %w", err)
	}
	return &ch, nil
}
