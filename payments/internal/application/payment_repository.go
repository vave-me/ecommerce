package application

import (
	"context"
	"middleman/payments/internal/models"
)

type PaymentRepository interface {
	Save(ctx context.Context, payment *models.Payment) error
	Find(ctx context.Context, paymentID string) (*models.Payment, error)
	SavePaymentMethod(ctx context.Context, payment *models.Payment) error
}
