package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stackus/errors"
	"middleman/activity/internal/domain"
	"middleman/internal/postgres"
)

type UserCacheRepository struct {
	tableName string
	db        postgres.DB
	fallback  domain.UserRepository
}

var _ domain.UserCacheRepository = (*UserCacheRepository)(nil)

func NewUserCacheRepository(tableName string, db postgres.DB, fallback domain.UserRepository) UserCacheRepository {
	return UserCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r UserCacheRepository) Add(ctx context.Context, userID, email, username, location string, enabled bool) error {
	const query = "INSERT INTO %s (id,email, username,location,enabled) VALUES ($1, $2, $3, $4, $5)"

	_, err := r.db.ExecContext(ctx, r.table(query), userID, email, username, location, enabled)

	return err
}

func (r UserCacheRepository) Rename(ctx context.Context, userID, username string) error {
	const query = `UPDATE %s SET username = $2 WHERE userID = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), userID, username)

	return err
}
func (r UserCacheRepository) Find(ctx context.Context, userID string) (*domain.User, error) {
	const query = `SELECT email, username, location, enabled FROM %s WHERE id = $1 LIMIT 1`

	user := &domain.User{
		ID: userID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), userID).Scan(&user.Email, &user.Username, &user.Location, &user.Enabled)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning user customer")
		}
		user, err = r.fallback.Find(ctx, userID)
		if err != nil {
			return nil, errors.Wrap(err, "customer fallback failed")
		}
		// attempt to add it to the cache
		return user, r.Add(ctx, user.ID, user.Email, user.Username, user.Location, user.Enabled)
	}

	return user, nil
}

func (r UserCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
