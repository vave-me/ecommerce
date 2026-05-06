package postgres

import (
	"context"
	"fmt"
	"middleman/internal/postgres"
	"middleman/wishlists/internal/domain"

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

func (r MiddlemanRepository) AddWishlist(ctx context.Context, wishlistID, userID, name string) error {
	const query = "INSERT INTO %s (id, user_id, name) VALUES ($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, r.table(query), wishlistID, userID, name)

	return err
}

func (r MiddlemanRepository) Find(ctx context.Context, userID, name string) (*domain.MiddlemanWishlist, error) {
	const query = "SELECT id FROM %s WHERE user_id = $1 AND name = $2 LIMIT 1"

	wishlist := &domain.MiddlemanWishlist{
		UserID: userID,
		Name:   name,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), userID, name).Scan(&wishlist.ID)
	if err != nil {
		return nil, errors.Wrap(err, "scanning wishlist")
	}

	return wishlist, nil
}
func (r MiddlemanRepository) GetWishlists(ctx context.Context, userID string) ([]*domain.MiddlemanWishlist, error) {

	const query = "SELECT id, name FROM %s WHERE user_id = $1"

	rows, err := r.db.QueryContext(ctx, r.table(query), userID)

	if err != nil {
		return nil, errors.Wrap(err, "querying all wishlists")
	}
	defer rows.Close()

	var wishlists []*domain.MiddlemanWishlist
	for rows.Next() {
		wishlist := &domain.MiddlemanWishlist{
			UserID: userID,
		}
		err := rows.Scan(
			&wishlist.ID,
			&wishlist.Name,
		)
		if err != nil {
			return nil, err
		}
		wishlists = append(wishlists, wishlist)
	}

	return wishlists, nil
}

func (r MiddlemanRepository) Remove(ctx context.Context, wishlistID string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), wishlistID)

	return err
}
func (r MiddlemanRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
