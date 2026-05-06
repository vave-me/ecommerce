package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"middleman/internal/postgres"
	"time"

	"github.com/stackus/errors"

	"middleman/users/internal/domain"
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

// 1) AddUser with lat/lng -> location geography(Point,4326)
func (r MiddlemanRepository) AddUser(
	ctx context.Context,
	id string,
	email string,
	username string,
	firstname string,
	lastname string,
	googleID string,
	enabled bool,
	lat float64,
	lng float64,
	thumbnail string, // if you also want to store user thumbnail
	role string,
) error {

	// We omit the old 'location string' param
	// We'll insert lat/lng into `location` using ST_SetSRID(ST_MakePoint($lng, $lat),4326)
	const query = `
        INSERT INTO %s (
            id,
            email,
            username,
            firstname,
            lastname,
            google_id,
            enabled,
            thumbnail,
            lat,lng,            
            location,
            role
        )
        VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
            ST_SetSRID(ST_MakePoint($10, $9), 4326),
            $11
        )
    `
	// Ensuring googleID is passed as a string, which matches the domain model type
	_, err := r.db.ExecContext(ctx, r.table(query),
		id,
		email,
		username,
		firstname,
		lastname,
		googleID, // This should be string type in the database schema
		enabled,
		thumbnail,
		lat,
		lng,
		role,
	)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return err
	}
	return nil
}

// 1) AddUser with lat/lng -> location geography(Point,4326)
func (r MiddlemanRepository) UpdateUser(
	ctx context.Context,
	id string,
	username string,
	firstname string,
	lastname string,
	bio string,
	privacy string,
	background string,
	lat float64,
	lng float64,
	thumbnail string, // if you also want to store user thumbnail
	role string,
) error {

	const query = `
        UPDATE %s SET
            username = $2,
            firstname = $3,
            lastname = $4,
            bio = $5,
            privacy = $6,
            background = $7,
            thumbnail = $8,
            lat = $9, 
            lng = $10,
            location = ST_SetSRID(ST_MakePoint($10, $9), 4326),
            role = $11
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		id,
		username,
		firstname,
		lastname,
		bio,
		privacy,
		background,
		thumbnail,
		lat,
		lng,
		role,
	)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}
	return nil
}

func (r MiddlemanRepository) EnableUser(ctx context.Context, userID string, participating bool) error {
	const query = "UPDATE %s SET enabled = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), userID, participating)

	return err
}

func (r MiddlemanRepository) RenameUser(ctx context.Context, userID, name string) error {
	const query = "UPDATE %s SET username = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), userID, name)

	return err
}

func (r MiddlemanRepository) Find(ctx context.Context, userID string) (*domain.MiddlemanUser, error) {
	const query = "SELECT email, username, firstname, lastname, bio, background, privacy, enabled, google_id, lat, lng, thumbnail, role FROM %s WHERE id = $1 LIMIT 1"
	user := &domain.MiddlemanUser{
		ID: userID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), userID).Scan(
		&user.Email,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Bio,
		&user.Background,
		&user.Privacy,
		&user.Enabled,
		&user.GoogleID,
		&user.Lat,
		&user.Lng,
		&user.Thumbnail,
		&user.Role,
	)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Find: error finding user %s: %v", userID, err))
	}

	return user, err
}

func (r MiddlemanRepository) FindByEmail(ctx context.Context, email string) (*domain.MiddlemanUser, error) {
	const query = "SELECT id, email, username, lat, lng, firstname, lastname, enabled, google_id, role FROM %s WHERE email = $1 LIMIT 1"
	user := &domain.MiddlemanUser{}

	// Using a sql.NullString to handle potentially null google_id values
	var googleID sql.NullString

	err := r.db.QueryRowContext(ctx, r.table(query), email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Lat,
		&user.Lng,
		&user.FirstName,
		&user.LastName,
		&user.Enabled,
		&googleID,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, fmt.Sprintf("user not found with email: %s", email))
		}
		return nil, err
	}

	// Only assign the GoogleID if it's valid/non-null
	if googleID.Valid {
		user.GoogleID = googleID.String
	}

	return user, nil
}

func (r MiddlemanRepository) FindByGoogleID(ctx context.Context, googleID string) (*domain.MiddlemanUser, error) {
	const query = "SELECT id, email, username, lat, lng, firstname, lastname, enabled, role FROM %s WHERE google_id = $1 LIMIT 1"
	user := &domain.MiddlemanUser{}

	// When scanning from the database, google_id is a string type column
	err := r.db.QueryRowContext(ctx, r.table(query), googleID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Lat,
		&user.Lng,
		&user.FirstName,
		&user.LastName,
		&user.Enabled,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, fmt.Sprintf("user not found with googleID: %s", googleID))
		}
		return nil, err
	}

	// Set the GoogleID field since we queried by it
	user.GoogleID = googleID
	return user, nil
}

func (r MiddlemanRepository) All(ctx context.Context) (users []*domain.MiddlemanUser, err error) {
	const query = "SELECT id, email, username, firstname, lastname, enabled, google_id, lat, lng, thumbnail, role FROM %s"
	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, errors.Wrap(err, "querying users")
	}
	defer rows.Close()

	for rows.Next() {
		user := new(domain.MiddlemanUser)
		var googleID sql.NullString

		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.Enabled, &googleID, &user.Lat, &user.Lng, &user.Thumbnail, &user.Role)
		if err != nil {
			return nil, errors.Wrap(err, "scanning user")
		}

		if googleID.Valid {
			user.GoogleID = googleID.String
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing user rows")
	}

	return users, nil
}

func (r MiddlemanRepository) AllEnabled(ctx context.Context) (users []*domain.MiddlemanUser, err error) {
	const query = "SELECT id, email, username, firstname, lastname, enabled, google_id, lat, lng, thumbnail, role FROM %s WHERE enabled is true"

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, errors.Wrap(err, "querying enabled users")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing enabled user rows")
		}
	}(rows)

	for rows.Next() {
		user := new(domain.MiddlemanUser)
		var googleID sql.NullString

		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.Enabled, &googleID, &user.Lat, &user.Lng, &user.Thumbnail, &user.Role)
		if err != nil {
			return nil, errors.Wrap(err, "scanning user")
		}

		if googleID.Valid {
			user.GoogleID = googleID.String
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing enabled users rows")
	}

	return users, nil
}

// LogUserIn logs the user in by storing the necessary session/token data
func (r MiddlemanRepository) LogUserIn(ctx context.Context, userID string) error {
	const query = "UPDATE %s SET last_login = $2 WHERE id = $1"
	_, err := r.db.ExecContext(ctx, r.table(query), userID, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// ID        string
// Username  string
// Lat       float64
// Lng       float64
// Thumbnail string
func (r MiddlemanRepository) FindSimple(ctx context.Context, userID string) (*domain.MiddlemanViewUser, error) {
	const query = "SELECT username, lat, lng, location, thumbnail, bio, background, privacy FROM %s WHERE id = $1 LIMIT 1"
	user := &domain.MiddlemanViewUser{
		ID: userID,
	}
	err := r.db.QueryRowContext(ctx, r.table(query), userID).Scan(
		&user.Username,
		&user.Lat,
		&user.Lng,
		&user.Location,
		&user.Thumbnail,
		&user.Bio,
		&user.Background,
		&user.Privacy,
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Find: error finding user %s: %v", userID, err))
	}
	return user, err
}

// LogUserOut logs the user out by clearing the necessary session/token data
func (r MiddlemanRepository) LogUserOut(ctx context.Context, userID string) error {
	const query = "UPDATE %s SET last_login = '' WHERE id = $1"
	_, err := r.db.ExecContext(ctx, r.table(query), userID, time.Now())
	if err != nil {
		return err
	}
	return nil
}
func (r MiddlemanRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
