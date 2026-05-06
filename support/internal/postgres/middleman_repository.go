package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"middleman/internal/postgres"
	"middleman/support/internal/domain"

	"github.com/stackus/errors"
)

type MiddlemanRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.MiddlemanRepository = (*MiddlemanRepository)(nil)

func NewMiddlemanRepository(tableName string, db postgres.DB) MiddlemanRepository {
	return MiddlemanRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r MiddlemanRepository) Add(ctx context.Context, supportID, userSellerID, userCustomerID, productID string, price int64) error {
	const query = "INSERT INTO %s (id, user_seller_id, user_customer_id, product_id, price) VALUES ($1, $2, $3, $4,$5)"

	_, err := r.db.ExecContext(ctx, r.table(query), supportID, userSellerID, userCustomerID, productID, price)

	return err
}

func (r MiddlemanRepository) Find(ctx context.Context, supportID string) (*domain.MiddlemanSupport, error) {
	const query = "SELECT user_seller_id, user_customer_id, product_id, price FROM %s WHERE id = $1 LIMIT 1"

	support := &domain.MiddlemanSupport{
		ID: supportID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), supportID).Scan(&support.UserSellerID, &support.UserCustomerID, &support.ProductID, &support.Price)
	if err != nil {
		return nil, errors.Wrap(err, "scanning support")
	}

	return support, nil
}

func (r MiddlemanRepository) All(ctx context.Context, userID string) (supports []*domain.MiddlemanSupport, err error) {
	const query = "SELECT id,user_seller_id, user_customer_id, product_id, price FROM %s WHERE user_customer_id = $1 OR user_seller_id = $1"

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, errors.Wrap(err, "querying support")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing support rows")
		}
	}(rows)

	for rows.Next() {
		support := new(domain.MiddlemanSupport)
		err := rows.Scan(&support.ID, &support.UserSellerID, &support.UserCustomerID, &support.ProductID, &support.Price)
		if err != nil {
			return nil, errors.Wrap(err, "scanning support")
		}

		supports = append(supports, support)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing support rows")
	}

	return supports, nil
}

func (r MiddlemanRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
