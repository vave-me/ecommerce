package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stackus/errors"
	"middleman/baskets/internal/domain"

	"middleman/internal/postgres"
)

type CatalogRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.CatalogRepository = (*CatalogRepository)(nil)

func NewCatalogRepository(tableName string, db postgres.DB) CatalogRepository {
	return CatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r CatalogRepository) Add(ctx context.Context, basketID, userID string, status domain.BasketStatus) error {
	const query = `INSERT INTO %s (id, user_id, basket_status) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, r.table(query), basketID, userID, status)

	return err
}

func (r CatalogRepository) Update(ctx context.Context, basketID string, status domain.BasketStatus) error {
	const query = `UPDATE %s SET status = $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), basketID, status)

	return err
}

func (r CatalogRepository) Remove(ctx context.Context, basketID string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), basketID)

	return err
}

func (r CatalogRepository) Find(ctx context.Context, userID string) (*domain.CatalogBasket, error) {
	const query = `SELECT id,basket_status FROM %s WHERE user_id = $1 LIMIT 1`

	basket := &domain.CatalogBasket{
		UserID: userID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), userID).Scan(&basket.ID, &basket.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msg("basket with that ID does not exist")
		}
		return nil, errors.Wrap(err, "scanning basket")
	}

	return basket, nil
}

func (r CatalogRepository) All(ctx context.Context, userID string) (baskets []*domain.CatalogBasket, err error) {
	const query = `SELECT id,user_id, basket_status FROM %s WHERE id = $1 LIMIT 1`

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query), userID)
	if err != nil {
		return nil, errors.Wrap(err, "querying baskets")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing basket rows")
		}
	}(rows)

	for rows.Next() {
		basket := &domain.CatalogBasket{
			UserID: userID,
		}
		err := rows.Scan(&basket.ID, &basket.UserID, basket.Status)
		if err != nil {
			return nil, errors.Wrap(err, "scanning basket")
		}

		baskets = append(baskets, basket)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing basket rows")
	}

	return baskets, nil
}

func (r CatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
