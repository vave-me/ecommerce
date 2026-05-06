package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/notifications/internal/domain"
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

func (r UserCacheRepository) Add(ctx context.Context, userID, firstName, lastName, email string) error {
	const query = "INSERT INTO %s (id, first_name, last_name, email) VALUES ($1, $2, $3, $4)"

	_, err := r.db.ExecContext(ctx, r.table(query), userID, firstName, lastName, email)

	return err
}

func (r UserCacheRepository) Rename(ctx context.Context, userID, firstName string) error {
	const query = `UPDATE %s SET first_name = $2 WHERE userID = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), userID, firstName)

	return err
}
func (r UserCacheRepository) Find(ctx context.Context, userID string) (*domain.User, error) {
	const query = `SELECT first_name, last_name, email FROM %s WHERE id = $1 LIMIT 1`

	user := &domain.User{
		ID: userID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), userID).Scan(&user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning user customer")
		}
		user, err = r.fallback.Find(ctx, userID)
		if err != nil {
			return nil, errors.Wrap(err, "customer fallback failed")
		}
		// attempt to add it to the cache
		return user, r.Add(ctx, user.ID, user.FirstName, user.LastName, user.Email)
	}

	return user, nil
}

func (r UserCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
