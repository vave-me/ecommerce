package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"

	"middleman/baskets/internal/domain"
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

// Add inserts a user in the cache table.
// IMPORTANT: The column order MUST match the argument order you pass in ExecContext.
func (r UserCacheRepository) Add(ctx context.Context, userID, email, username, firstName, lastName, location string, enabled bool) error {
	const query = `
        INSERT INTO %s (
            id,
            email,
            username,
            first_name,
            last_name,
            location,
            enabled
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.db.ExecContext(
		ctx,
		r.table(query),
		userID,
		email,
		username,
		firstName,
		lastName,
		location,
		enabled,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Ignore "duplicate key" errors so we don't fail if it already exists
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil
			}
		}
		return err
	}
	return nil
}

// Rename updates only the first_name field for demonstration purposes.
func (r UserCacheRepository) Rename(ctx context.Context, userID, name string) error {
	const query = `
        UPDATE %s
        SET first_name = $2
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query), userID, name)
	return err
}

// Find retrieves a user from the cache; if not found, it falls back to the main user repository.
func (r UserCacheRepository) Find(ctx context.Context, userID string) (*domain.User, error) {
	const query = `
        SELECT
            email,
            username,
            first_name,
            last_name,
            location,
            enabled
        FROM %s
        WHERE id = $1
        LIMIT 1
    `

	user := &domain.User{
		ID: userID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), userID).Scan(
		&user.Email,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Location,
		&user.Enabled,
	)
	if err != nil {
		// If row not found in cache
		if errors.Is(err, sql.ErrNoRows) {
			// Use fallback
			fallbackUser, fbErr := r.fallback.Find(ctx, userID)
			if fbErr != nil {
				return nil, errors.Wrap(fbErr, "fallback user retrieval failed")
			}
			// Attempt to cache it
			cacheErr := r.Add(ctx,
				fallbackUser.ID,
				fallbackUser.Email,
				fallbackUser.Username,
				fallbackUser.FirstName,
				fallbackUser.LastName,
				fallbackUser.Location,
				fallbackUser.Enabled,
			)
			// Return user even if caching fails
			if cacheErr != nil {
				return fallbackUser, errors.Wrap(cacheErr, "failed to cache fallback user")
			}
			return fallbackUser, nil
		}
		// Some other error
		return nil, errors.Wrap(err, "scanning user")
	}

	// Found in cache
	return user, nil
}

func (r UserCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
