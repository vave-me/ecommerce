package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"middleman/internal/postgres"
	"middleman/wishlists/internal/domain"

	"github.com/stackus/errors"
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

func (r CatalogRepository) AddWishlistItem(ctx context.Context, wishlistItemID, wishlistID, itemID, entityType string) error {
	const query = `INSERT INTO %s (id, wishlist_id,item_id, entity_type) VALUES ($1, $2, $3,$4)`

	_, err := r.db.ExecContext(ctx, r.table(query), wishlistItemID, wishlistID, itemID)

	return err
}

func (r CatalogRepository) RemoveWishlistItem(ctx context.Context, wishlistItemID string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), wishlistItemID)

	return err
}

func (r CatalogRepository) Find(ctx context.Context, wishlistItemID string) (*domain.CatalogWishlistItem, error) {
	const query = `SELECT wishlist_id, item_id,entity_type  FROM %s WHERE id = $1 LIMIT 1`

	wishlistItem := &domain.CatalogWishlistItem{
		ID: wishlistItemID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), wishlistItemID).Scan(&wishlistItem.WishlistID, &wishlistItem.ItemID, &wishlistItem.EntityType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msg("product with that ID does not exist")
		}
		return nil, errors.Wrap(err, "scanning wishlist products")
	}

	return wishlistItem, nil
}

func (r CatalogRepository) GetWishlistItems(ctx context.Context, wishlistID string) (wishlistItems []*domain.CatalogWishlistItem, err error) {
	const query = `SELECT id,item_id, entity_type FROM %s WHERE wishlist_id = $1`

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query), wishlistID)
	if err != nil {
		return nil, errors.Wrap(err, "querying wishlists products")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing wishlist rows")
		}
	}(rows)

	for rows.Next() {
		product := &domain.CatalogWishlistItem{
			WishlistID: wishlistID,
		}
		err := rows.Scan(&product.ID, &product.ItemID, &product.EntityType)
		if err != nil {
			return nil, errors.Wrap(err, "scanning wishlist product")
		}

		wishlistItems = append(wishlistItems, product)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing wishlist products rows")
	}

	return wishlistItems, nil
}

func (r CatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
