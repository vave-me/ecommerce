package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"
	"middleman/activity/internal/domain"

	"middleman/internal/postgres"
)

type ProductCacheRepository struct {
	tableName string
	db        postgres.DB
	fallback  domain.ProductRepository
}

var _ domain.ProductCacheRepository = (*ProductCacheRepository)(nil)

func NewProductCacheRepository(tableName string, db postgres.DB, fallback domain.ProductRepository) ProductCacheRepository {
	return ProductCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r ProductCacheRepository) Add(ctx context.Context, productID, name, description string, price int64, userSellerID string, stock int64, sku, categoryID, active string) error {
	const query = `INSERT INTO %s (id, name, description, price, user_seller_id,stock, sku, category_id, active) VALUES ($1, $2, $3, $4,$5,$6, $7, $8,$9)`

	_, err := r.db.ExecContext(ctx, r.table(query), productID, name, description, price, userSellerID, stock, sku, categoryID, active)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil
			}
		}
	}

	return err
}

func (r ProductCacheRepository) Rebrand(ctx context.Context, productID, name string) error {
	const query = `UPDATE %s SET name = $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), productID, name)

	return err
}

func (r ProductCacheRepository) UpdatePrice(ctx context.Context, productID string, delta int64) error {
	const query = `UPDATE %s SET price = price + $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), productID, delta)

	return err
}

func (r ProductCacheRepository) Remove(ctx context.Context, productID string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), productID)

	return err
}

func (r ProductCacheRepository) Find(ctx context.Context, productID string) (*domain.Product, error) {
	const query = `SELECT name, description, price, user_seller_id,stock,sku,category_id, activeFROM %s WHERE id = $1 LIMIT 1`

	product := &domain.Product{
		ID: productID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), productID).Scan(&product.Name, &product.Description, &product.Price, &product.UserSellerID, &product.Stock, &product.SKU, &product.CategoryID, &product.Active)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning product")
		}
		product, err = r.fallback.Find(ctx, productID)
		if err != nil {
			return nil, errors.Wrap(err, "product fallback failed")
		}
		// attempt to add it to the cache
		return product, r.Add(ctx, product.ID, product.Name, product.Description, product.Price, product.UserSellerID, product.Stock, product.SKU, product.CategoryID, product.Active)
	}

	return product, nil
}

func (r ProductCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
