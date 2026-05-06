package postgres

import (
	"context"
	"fmt"
	"middleman/internal/postgres"
	"middleman/payments/internal/application"
	"middleman/payments/internal/models"

	"github.com/rs/zerolog/log"
	"github.com/stackus/errors"
)

type InvoiceRepository struct {
	tableName string
	db        postgres.DB
}

var _ application.InvoiceRepository = (*InvoiceRepository)(nil)

func NewInvoiceRepository(tableName string, db postgres.DB) InvoiceRepository {
	return InvoiceRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r InvoiceRepository) Find(ctx context.Context, invoiceID string) (*models.Invoice, error) {
	const query = "SELECT order_id, amount, paid_amt, status FROM %s WHERE id = $1 LIMIT 1"

	invoice := &models.Invoice{
		ID: invoiceID,
	}
	var status string
	err := r.db.QueryRowContext(ctx, r.table(query), invoiceID).Scan(&invoice.OrderID, &invoice.Amount, &invoice.PaidAmt, &status)
	if err != nil {
		log.Error().Err(err).Str("invoice_id", invoiceID).Msg("InvoiceRepository.Find query failed")
		return nil, errors.Wrap(err, "scanning invoice")
	}

	invoice.Status, err = r.statusToDomain(status)
	if err != nil {
		log.Error().Err(err).Str("invoice_id", invoiceID).Msg("InvoiceRepository.statusToDomain failed")
		return nil, err
	}
	log.Debug().Str("invoice_id", invoiceID).Msg("InvoiceRepository.Find succeeded")

	return invoice, nil
}

func (r InvoiceRepository) Save(ctx context.Context, invoice *models.Invoice) error {
	const query = "INSERT INTO %s (id, order_id, amount, paid_amt, status) VALUES ($1, $2, $3, $4, $5)"

	_, err := r.db.ExecContext(ctx, r.table(query), invoice.ID, invoice.OrderID, invoice.Amount, invoice.PaidAmt, invoice.Status.String())
	if err != nil {
		log.Error().Err(err).Str("invoice_id", invoice.ID).Msg("InvoiceRepository.Save failed")
	} else {
		log.Debug().Str("invoice_id", invoice.ID).Msg("InvoiceRepository.Save succeeded")
	}
	return err
}

func (r InvoiceRepository) Update(ctx context.Context, invoice *models.Invoice) error {
	const query = "UPDATE %s SET amount = $2, paid_amt = $3, status = $4 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), invoice.ID, invoice.Amount, invoice.PaidAmt, invoice.Status.String())
	if err != nil {
		log.Error().Err(err).Str("invoice_id", invoice.ID).Msg("InvoiceRepository.Update failed")
	} else {
		log.Debug().Str("invoice_id", invoice.ID).Msg("InvoiceRepository.Update succeeded")
	}
	return err
}

func (r InvoiceRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

func (r InvoiceRepository) statusToDomain(status string) (models.InvoiceStatus, error) {
	switch status {
	case models.InvoiceIsPending.String():
		return models.InvoiceIsPending, nil
	case models.InvoiceIsPaid.String():
		return models.InvoiceIsPaid, nil
	case models.InvoiceIsPartially.String():
		return models.InvoiceIsPartially, nil
	case models.InvoiceIsCanceled.String():
		return models.InvoiceIsCanceled, nil
	default:
		return models.InvoiceIsUnknown, fmt.Errorf("unknown invoice status: %s", status)
	}
}
