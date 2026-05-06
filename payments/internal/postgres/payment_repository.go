package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"middleman/internal/postgres"
	"middleman/payments/internal/application"
	"middleman/payments/internal/models"

	"github.com/rs/zerolog/log"
)

type PaymentRepository struct {
	tableName string
	db        postgres.DB
}

var _ application.PaymentRepository = (*PaymentRepository)(nil)

func NewPaymentRepository(tableName string, db postgres.DB) PaymentRepository {
	return PaymentRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r PaymentRepository) Save(ctx context.Context, payment *models.Payment) error {
	const query = `INSERT INTO %s (id, user_customer_id, amount, status, payment_intent_id, order_id, payment_method_id)
	              VALUES ($1, $2, $3, $4, $5, $6, $7)
	              ON CONFLICT (id) DO UPDATE
	              SET user_customer_id   = EXCLUDED.user_customer_id,
	                  amount             = EXCLUDED.amount,
	                  status             = EXCLUDED.status,
	                  payment_intent_id  = EXCLUDED.payment_intent_id,
	                  order_id           = EXCLUDED.order_id,
	                  payment_method_id  = EXCLUDED.payment_method_id`

	// Convert empty strings to NULL for nullable fields
	var orderID, paymentMethodID sql.NullString
	if payment.OrderID != "" {
		orderID = sql.NullString{String: payment.OrderID, Valid: true}
	}
	if payment.PaymentMethodID != "" {
		paymentMethodID = sql.NullString{String: payment.PaymentMethodID, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, r.table(query), 
		payment.ID, 
		payment.UserCustomerID, 
		payment.Amount, 
		payment.Status, 
		payment.PaymentIntentID,
		orderID,
		paymentMethodID)
	
	if err != nil {
		log.Error().Err(err).Str("payment_id", payment.ID).Msg("PaymentRepository.Save failed")
	} else {
		log.Debug().Str("payment_id", payment.ID).Msg("PaymentRepository.Save succeeded")
	}
	return err
}

func (r PaymentRepository) SavePaymentMethod(ctx context.Context, payment *models.Payment) error {
	const query = "UPDATE %s SET payment_method_id = $2, status = $3 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), payment.ID, payment.PaymentMethodID, payment.Status)
	if err != nil {
		log.Error().Err(err).Str("payment_id", payment.ID).Msg("PaymentRepository.SavePaymentMethod failed")
	} else {
		log.Debug().Str("payment_id", payment.ID).Msg("PaymentRepository.SavePaymentMethod succeeded")
	}
	return err
}

func (r PaymentRepository) Find(ctx context.Context, paymentID string) (*models.Payment, error) {
	const query = "SELECT user_customer_id, amount, payment_method_id, payment_intent_id, order_id, status FROM %s WHERE id = $1 LIMIT 1"

	payment := &models.Payment{
		ID: paymentID,
	}

	// Use sql.NullString for nullable fields
	var paymentMethodID, orderID sql.NullString

	err := r.db.QueryRowContext(ctx, r.table(query), paymentID).Scan(
		&payment.UserCustomerID,
		&payment.Amount,
		&paymentMethodID,
		&payment.PaymentIntentID,
		&orderID,
		&payment.Status,
	)
	
	if err != nil {
		log.Error().Err(err).Str("payment_id", paymentID).Msg("PaymentRepository.Find failed")
		return nil, err
	}
	
	// Convert NULL values to empty strings
	payment.PaymentMethodID = paymentMethodID.String
	payment.OrderID = orderID.String
	
	log.Debug().Str("payment_id", paymentID).Msg("PaymentRepository.Find succeeded")
	return payment, nil
}

func (r PaymentRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
