package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/search/internal/application"
	"middleman/search/internal/models"
)

type UserCacheRepository struct {
	tableName string
	db        postgres.DB
	fallback  application.UserRepository
}

var _ application.UserCacheRepository = (*UserCacheRepository)(nil)

func NewUserCacheRepository(tableName string, db postgres.DB, fallback application.UserRepository) UserCacheRepository {
	return UserCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r UserCacheRepository) Add(ctx context.Context, userID, email, username, firstName, lastName, location string, enabled bool) error {
	const query = "INSERT INTO %s (id, email, username, first_name, last_name, location, enabled) VALUES ($1, $2, $3, $4,$5, $6, $7)"

	_, err := r.db.ExecContext(ctx, r.table(query), userID, email, username, firstName, lastName, location, enabled)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			// Ignore unique violation errors (user already exists)
			return nil
		}
		return errors.Wrap(err, "adding user")
	}

	return nil
}

func (r UserCacheRepository) Rename(ctx context.Context, userID string, name string) error {
	const query = "UPDATE %s SET username = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), userID, name)
	if err != nil {
		return errors.Wrap(err, "renaming user")
	}

	return nil
}

func (r UserCacheRepository) Update(ctx context.Context, user *models.User) error {
	const query = "UPDATE %s SET first_name = $2, last_name = $3, email = $4 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), user.ID, user.FirstName, user.LastName, user.Email)
	if err != nil {
		return errors.Wrap(err, "updating user")
	}

	return nil
}

func (r UserCacheRepository) Get(ctx context.Context, userID string) (*models.User, error) {
	const query = "SELECT  email, username,first_name, last_name, email, location, enabled FROM %s WHERE id = $1 LIMIT 1"

	user := &models.User{
		ID: userID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), userID).Scan(&user.Email, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Location, &user.Enabled)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msgf("user with ID %s not found", userID)
		}
		return nil, errors.Wrap(err, "getting user")
	}

	return user, nil
}

func (r UserCacheRepository) Find(ctx context.Context, userID string) (*models.User, error) {
	user, err := r.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			// Try fallback repository
			user, err = r.fallback.Find(ctx, userID)
			if err != nil {
				return nil, errors.Wrap(err, "user fallback failed")
			}
			// Attempt to add it to the cache
			err = r.Add(ctx, user.ID, user.Email, user.Username, user.FirstName, user.LastName, user.Location, user.Enabled)
			if err != nil {
				return nil, errors.Wrap(err, "adding user to cache")
			}
			return user, nil
		}
		return nil, err
	}
	return user, nil
}

//func (r UserCacheRepository) Search(ctx context.Context, search application.SearchUsers) ([]*models.User, error) {
//	const baseQuery = "SELECT id, first_name, last_name, email FROM %s"
//
//	var (
//		conditions []string
//		args       []interface{}
//		argID      = 1
//	)
//
//	// Build conditions based on filters
//	if search.Filters.Name != "" {
//		conditions = append(conditions, fmt.Sprintf("(first_name || ' ' || last_name) ILIKE $%d", argID))
//		args = append(args, "%"+search.Filters.Name+"%")
//		argID++
//	}
//
//	if search.Filters.Email != "" {
//		conditions = append(conditions, fmt.Sprintf("email ILIKE $%d", argID))
//		args = append(args, "%"+search.Filters.Email+"%")
//		argID++
//	}
//
//	// Build the final query
//	query := r.table(baseQuery)
//
//	if len(conditions) > 0 {
//		query += " WHERE " + strings.Join(conditions, " AND ")
//	}
//
//	// Add ORDER BY, LIMIT, and OFFSET
//	query += " ORDER BY first_name ASC, last_name ASC"
//
//	if search.Limit > 0 {
//		query += fmt.Sprintf(" LIMIT $%d", argID)
//		args = append(args, search.Limit)
//		argID++
//	}
//
//	if search.Offset > 0 {
//		query += fmt.Sprintf(" OFFSET $%d", argID)
//		args = append(args, search.Offset)
//		argID++
//	}
//
//	// Execute the query
//	rows, err := r.db.QueryContext(ctx, query, args...)
//	if err != nil {
//		return nil, errors.Wrap(err, "querying users")
//	}
//	defer rows.Close()
//
//	var users []*models.User
//
//	for rows.Next() {
//		user := &models.User{}
//		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
//		if err != nil {
//			return nil, errors.Wrap(err, "scanning user")
//		}
//
//		users = append(users, user)
//	}
//
//	if err := rows.Err(); err != nil {
//		return nil, errors.Wrap(err, "iterating user rows")
//	}
//
//	return users, nil
//}

func (r UserCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
