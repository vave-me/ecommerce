package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type PaymentRepository interface {
	// Payment operations
	AuthorizePaymentForCustomer(ctx context.Context, userCustomerID string, amount int64) (*models.AuthorizePaymentResponse, error)
	ConfirmPaymentWithMethod(ctx context.Context, paymentID, paymentMethodID string) (*models.ConfirmPaymentResponse, error)
	CaptureAuthorizedPaymentAmount(ctx context.Context, paymentID string, amountToCapture int64) (*models.CapturePaymentResponse, error)

	// Invoice operations
	CreateNewInvoiceForOrder(ctx context.Context, orderID, paymentID string, amount int64) (*models.CreateInvoiceResponse, error)
	AdjustInvoiceAmountWithReason(ctx context.Context, invoiceID string, amount int64, reason string) (*models.AdjustInvoiceResponse, error)
	PayInvoiceUsingPaymentMethod(ctx context.Context, invoiceID, paymentMethodID string) (*models.PayInvoiceResponse, error)
	CancelInvoiceWithReason(ctx context.Context, invoiceID, reason string) error

	// Webhook handling
	ProcessWebhookNotification(ctx context.Context, rawBody, signature string) (*models.HandleWebhookResponse, error)

	// Additional query methods for AI tooling
	GetPaymentDetailsByID(ctx context.Context, paymentID string) (*models.Payment, error)
	GetInvoiceDetailsByID(ctx context.Context, invoiceID string) (*models.Invoice, error)
	GetCustomerPaymentHistory(ctx context.Context, userCustomerID string) ([]*models.Payment, error)
	GetAllInvoicesForOrder(ctx context.Context, orderID string) ([]*models.Invoice, error)
	SearchPaymentsByStatus(ctx context.Context, status string, limit int64) ([]*models.Payment, error)
	SearchInvoicesByStatus(ctx context.Context, status string, limit int64) ([]*models.Invoice, error)
}
