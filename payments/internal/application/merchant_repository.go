package application

import (
	"context"
	"middleman/payments/internal/models"
)

type MerchantRepository interface {
	Find(ctx context.Context, invoiceID string) (*models.Invoice, error)
	Save(ctx context.Context, invoice *models.Invoice) error
	Update(ctx context.Context, invoice *models.Invoice) error
}
