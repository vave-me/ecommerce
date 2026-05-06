package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type PaymentRepository interface {
	// Payment operations
	AuthorizePayment(ctx context.Context, userCustomerID string, amount int64) (*models.AuthorizePaymentResponse, error)
	ConfirmPayment(ctx context.Context, paymentID, paymentMethodID string) (*models.ConfirmPaymentResponse, error)
	CapturePayment(ctx context.Context, paymentID string, amountToCapture int64) (*models.CapturePaymentResponse, error)

	// Invoice operations
	CreateInvoice(ctx context.Context, orderID, paymentID string, amount int64) (*models.CreateInvoiceResponse, error)
	AdjustInvoice(ctx context.Context, invoiceID string, amount int64, reason string) (*models.AdjustInvoiceResponse, error)
	PayInvoice(ctx context.Context, invoiceID, paymentMethodID string) (*models.PayInvoiceResponse, error)
	CancelInvoice(ctx context.Context, invoiceID, reason string) error

	// Webhook handling
	HandleWebhook(ctx context.Context, rawBody, signature string) (*models.HandleWebhookResponse, error)

	// Additional query methods for AI tooling
	GetPayment(ctx context.Context, paymentID string) (*models.Payment, error)
	GetInvoice(ctx context.Context, invoiceID string) (*models.Invoice, error)
	GetPaymentsByCustomer(ctx context.Context, userCustomerID string) ([]*models.Payment, error)
	GetInvoicesByOrder(ctx context.Context, orderID string) ([]*models.Invoice, error)
	SearchPayments(ctx context.Context, status string, limit int64) ([]*models.Payment, error)
	SearchInvoices(ctx context.Context, status string, limit int64) ([]*models.Invoice, error)
}
