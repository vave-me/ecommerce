package adapters

import (
	"context"
	"middleman/payments/internal/application"
)

// paymentDomainWrapper implements the `PaymentDomain` interface
// by forwarding calls to your full application.App.
type paymentDomainWrapper struct {
	app application.App
}

// NewPaymentDomainWrapper constructs a wrapper that satisfies
// the minimal `PaymentDomain` interface by delegating to the full `application.App`.
func NewPaymentDomainWrapper(app application.App) application.PaymentDomain {
	return &paymentDomainWrapper{app: app}
}

// ConfirmPayment calls `ConfirmPayment` on the underlying app.
func (w *paymentDomainWrapper) ConfirmPayment(ctx context.Context, cmd application.ConfirmPaymentCommand) error {
	// Delegates to the real domain logic in app
	return w.app.ConfirmPayment(ctx, cmd)
}

// AuthorizePayment calls `AuthorizePayment` on the underlying app.
func (w *paymentDomainWrapper) AuthorizePayment(ctx context.Context, cmd application.AuthorizePaymentCommand) error {
	return w.app.AuthorizePayment(ctx, cmd)
}
