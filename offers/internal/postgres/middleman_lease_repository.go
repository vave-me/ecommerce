package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"middleman/internal/postgres"
	"middleman/offers/internal/domain"

	"github.com/stackus/errors"
)

type MiddlemanLeaseRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.MiddlemanLeaseRepository = (*MiddlemanLeaseRepository)(nil)

func NewMiddlemanLeaseRepository(tableName string, db postgres.DB) MiddlemanLeaseRepository {
	return MiddlemanLeaseRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r MiddlemanLeaseRepository) Add(ctx context.Context, offerID, userSellerID, userCustomerID, productID string, price int64) error {
	const query = "INSERT INTO %s (id, user_seller_id, user_customer_id, product_id, price) VALUES ($1, $2, $3, $4,$5)"

	_, err := r.db.ExecContext(ctx, r.table(query), offerID, userSellerID, userCustomerID, productID, price)

	return err
}

func (r MiddlemanLeaseRepository) Find(ctx context.Context, offerID string) (*domain.MiddlemanOffer, error) {
	const query = "SELECT user_seller_id, user_customer_id, product_id, price FROM %s WHERE id = $1 LIMIT 1"

	offer := &domain.MiddlemanOffer{
		ID: offerID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), offerID).Scan(&offer.UserSellerID, &offer.UserCustomerID, &offer.ProductID, &offer.Price)
	if err != nil {
		return nil, errors.Wrap(err, "scanning offer")
	}

	return offer, nil
}

func (r MiddlemanLeaseRepository) All(ctx context.Context, userID string) (offers []*domain.MiddlemanOffer, err error) {
	const query = "SELECT id,user_seller_id, user_customer_id, product_id, price FROM %s WHERE user_customer_id = $1 OR user_seller_id = $1"

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, errors.Wrap(err, "querying offers")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing offers rows")
		}
	}(rows)

	for rows.Next() {
		offer := new(domain.MiddlemanOffer)
		err := rows.Scan(&offer.ID, &offer.UserSellerID, &offer.UserCustomerID, &offer.ProductID, &offer.Price)
		if err != nil {
			return nil, errors.Wrap(err, "scanning offer")
		}

		offers = append(offers, offer)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing offer rows")
	}

	return offers, nil
}

func (r MiddlemanLeaseRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
