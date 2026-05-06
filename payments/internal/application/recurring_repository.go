package application

import (
	"context"
	"middleman/payments/internal/models"
)

type RecurringRepository interface {
	Save(ctx context.Context, payment *models.RecurringPaymentPlan) error
	Find(ctx context.Context, paymentID string) (*models.RecurringPaymentPlan, error)
	SavePaymentMethod(ctx context.Context, payment *models.RecurringPaymentPlan) error
}
