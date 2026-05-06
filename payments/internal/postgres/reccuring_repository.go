package postgres

import (
	"context"
	"fmt"
	"middleman/internal/postgres"
	"middleman/payments/internal/application"
	"middleman/payments/internal/models"

	"github.com/rs/zerolog/log"
)

type RecurringRepository struct {
	tableName string
	db        postgres.DB
}

// Ensure we implement the interface:
var _ application.RecurringRepository = (*RecurringRepository)(nil)

func NewRecurringRepository(tableName string, db postgres.DB) RecurringRepository {
	return RecurringRepository{
		tableName: tableName,
		db:        db,
	}
}

// Save inserts a new recurring payment plan.
func (r RecurringRepository) Save(ctx context.Context, plan *models.RecurringPaymentPlan) error {
	const query = `
	  INSERT INTO %s (
		plan_id,
		user_customer_id,
		amount,
		frequency,
		start_date,
		last_charged_at,
		next_due_date,
		status,
		payment_method_id
	  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(
		ctx,
		r.table(query),
		plan.PlanID,
		plan.UserCustomerID,
		plan.Amount,
		plan.Frequency,
		plan.StartDate,
		plan.LastChargedAt,
		plan.NextDueDate,
		plan.Status,
		plan.PaymentMethodID,
	)
	if err != nil {
		log.Error().Err(err).Str("plan_id", plan.PlanID).Msg("RecurringRepository.Save failed")
	} else {
		log.Debug().Str("plan_id", plan.PlanID).Msg("RecurringRepository.Save succeeded")
	}
	return err
}

// SavePaymentMethod updates the plan's payment method and status
func (r RecurringRepository) SavePaymentMethod(ctx context.Context, plan *models.RecurringPaymentPlan) error {
	const query = `
	  UPDATE %s
	     SET payment_method_id = $2,
	         status            = $3
	   WHERE plan_id          = $1
	`

	_, err := r.db.ExecContext(
		ctx,
		r.table(query),
		plan.PlanID,
		plan.PaymentMethodID,
		plan.Status,
	)
	if err != nil {
		log.Error().Err(err).Str("plan_id", plan.PlanID).Msg("RecurringRepository.SavePaymentMethod failed")
	} else {
		log.Debug().Str("plan_id", plan.PlanID).Msg("RecurringRepository.SavePaymentMethod succeeded")
	}
	return err
}

// Find retrieves a single recurring plan by its PlanID.
func (r RecurringRepository) Find(ctx context.Context, planID string) (*models.RecurringPaymentPlan, error) {
	const query = `
	  SELECT
	    user_customer_id,
	    amount,
	    frequency,
	    start_date,
	    last_charged_at,
	    next_due_date,
	    status,
	    payment_method_id
	  FROM %s
	  WHERE plan_id = $1
	  LIMIT 1
	`

	plan := &models.RecurringPaymentPlan{
		PlanID: planID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), planID).Scan(
		&plan.UserCustomerID,
		&plan.Amount,
		&plan.Frequency,
		&plan.StartDate,
		&plan.LastChargedAt,
		&plan.NextDueDate,
		&plan.Status,
		&plan.PaymentMethodID,
	)
	if err != nil {
		log.Error().Err(err).Str("plan_id", planID).Msg("RecurringRepository.Find failed")
	} else {
		log.Debug().Str("plan_id", planID).Msg("RecurringRepository.Find succeeded")
	}
	return plan, err
}

// table is a small helper to inject the table name into queries
func (r RecurringRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
