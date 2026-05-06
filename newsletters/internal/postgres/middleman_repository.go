package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"middleman/internal/postgres"
	"middleman/newsletters/internal/domain"

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

func (r MiddlemanRepository) Add(ctx context.Context, newsletterID, userSellerID, userCustomerID, productID string, price int64) error {
	const query = "INSERT INTO %s (id, user_seller_id, user_customer_id, product_id, price) VALUES ($1, $2, $3, $4,$5)"

	_, err := r.db.ExecContext(ctx, r.table(query), newsletterID, userSellerID, userCustomerID, productID, price)

	return err
}

func (r MiddlemanRepository) Find(ctx context.Context, newsletterID string) (*domain.MiddlemanNewsletter, error) {
	const query = "SELECT user_seller_id, user_customer_id, product_id, price FROM %s WHERE id = $1 LIMIT 1"

	newsletter := &domain.MiddlemanNewsletter{
		ID: newsletterID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), newsletterID).Scan(&newsletter.UserSellerID, &newsletter.UserCustomerID, &newsletter.ProductID, &newsletter.Price)
	if err != nil {
		return nil, errors.Wrap(err, "scanning newsletter")
	}

	return newsletter, nil
}

func (r MiddlemanRepository) All(ctx context.Context, userID string) (newsletters []*domain.MiddlemanNewsletter, err error) {
	const query = "SELECT id,user_seller_id, user_customer_id, product_id, price FROM %s WHERE user_customer_id = $1 OR user_seller_id = $1"

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, errors.Wrap(err, "querying newsletters")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing newsletters rows")
		}
	}(rows)

	for rows.Next() {
		newsletter := new(domain.MiddlemanNewsletter)
		err := rows.Scan(&newsletter.ID, &newsletter.UserSellerID, &newsletter.UserCustomerID, &newsletter.ProductID, &newsletter.Price)
		if err != nil {
			return nil, errors.Wrap(err, "scanning newsletter")
		}

		newsletters = append(newsletters, newsletter)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing newsletter rows")
	}

	return newsletters, nil
}

func (r MiddlemanRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
