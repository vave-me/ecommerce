package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"middleman/internal/postgres"
	"middleman/wishlists/internal/domain"

	"github.com/stackus/errors"
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

func (r ProductCacheRepository) Add(ctx context.Context, itemID, name, description string, price int64, userID string, stock int64, sku string, categoryID string) error {
	const query = `INSERT INTO %s (id, NAME,description,price, user_seller_id,stock,sku,category_id) VALUES ($1, $2, $3, $4, $5,$6,$7,$8) ON CONFLICT DO NOTHING`

	_, err := r.db.ExecContext(ctx, r.table(query), itemID, name, description, price, userID, stock, sku, categoryID)

	return err
}

func (r ProductCacheRepository) Rebrand(ctx context.Context, itemID, name string) error {
	const query = `UPDATE %s SET NAME = $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), itemID, name)

	return err
}

func (r ProductCacheRepository) UpdatePrice(ctx context.Context, itemID string, delta int64) error {
	const query = `UPDATE %s SET price = price + $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), itemID, delta)

	return err
}

func (r ProductCacheRepository) Remove(ctx context.Context, itemID string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), itemID)

	return err
}

func (r ProductCacheRepository) Find(ctx context.Context, itemID string) (*domain.Product, error) {
	const query = `SELECT name, description, price, user_seller_id, stock,sku,category_id,active,  FROM %s WHERE id = $1 LIMIT 1`

	product := &domain.Product{
		ID: itemID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), itemID).Scan(&product.Name, &product.Description, &product.Price, &product.UserSellerID, &product.Stock, &product.SKU, &product.CategoryID, &product.Active)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning product")
		}
		product, err = r.fallback.Find(ctx, itemID)
		if err != nil {
			return nil, errors.Wrap(err, "product fallback failed")
		}
		// attempt to add it to the cache
		return product, r.Add(ctx, product.ID, product.Name, product.Description, product.Price, product.UserSellerID, product.Stock, product.SKU, product.CategoryID)
	}

	return product, nil
}

func (r ProductCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
