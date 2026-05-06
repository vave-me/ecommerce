package domain

import (
	"context"
)

type PaymentRepository interface {
	// Payment operations
	AuthorizePayment(ctx context.Context, userCustomerID string, amount int64) (*AuthorizePaymentResponse, error)
	ConfirmPayment(ctx context.Context, paymentID, paymentMethodID string) (*ConfirmPaymentResponse, error)
	CapturePayment(ctx context.Context, paymentID string, amountToCapture int64) (*CapturePaymentResponse, error)

	// Invoice operations
	CreateInvoice(ctx context.Context, orderID, paymentID string, amount int64) (*CreateInvoiceResponse, error)
	AdjustInvoice(ctx context.Context, invoiceID string, amount int64, reason string) (*AdjustInvoiceResponse, error)
	PayInvoice(ctx context.Context, invoiceID, paymentMethodID string) (*PayInvoiceResponse, error)
	CancelInvoice(ctx context.Context, invoiceID, reason string) error

	// Webhook handling
	HandleWebhook(ctx context.Context, rawBody, signature string) (*HandleWebhookResponse, error)

	// Additional query methods for AI tooling
	GetPayment(ctx context.Context, paymentID string) (*Payment, error)
	GetInvoice(ctx context.Context, invoiceID string) (*Invoice, error)
	GetPaymentsByCustomer(ctx context.Context, userCustomerID string) ([]*Payment, error)
	GetInvoicesByOrder(ctx context.Context, orderID string) ([]*Invoice, error)
	SearchPayments(ctx context.Context, status string, limit int64) ([]*Payment, error)
	SearchInvoices(ctx context.Context, status string, limit int64) ([]*Invoice, error)
}
