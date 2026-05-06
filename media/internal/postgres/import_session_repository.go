package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/media/internal/domain"
)

// ImportSessionRepository satisfies domain.ImportSessionRepository
type ImportSessionRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.ImportSessionRepository = (*ImportSessionRepository)(nil)

// Constructor
func NewImportSessionRepository(tableName string, db postgres.DB) ImportSessionRepository {
	return ImportSessionRepository{
		tableName: tableName,
		db:        db,
	}
}

// Create adds a new import session
func (r ImportSessionRepository) Create(ctx context.Context, session *domain.ImportSession) error {
	const query = `
      INSERT INTO %s (
        id, 
        external_system_id, 
        external_system_type, 
        total_images, 
        processed_images, 
        failed_images, 
        status, 
        started_at, 
        metadata
      )
      VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9
      )
    `

	metadata, err := json.Marshal(session.Metadata)
	if err != nil {
		return errors.Wrap(err, "marshaling metadata")
	}

	_, err = r.db.ExecContext(ctx, r.table(query),
		session.ID,
		session.ExternalSystemID,
		session.ExternalSystemType,
		session.TotalImages,
		session.ProcessedImages,
		session.FailedImages,
		session.Status,
		session.StartedAt,
		metadata,
	)
	if err != nil {
		return errors.Wrap(err, "inserting import session")
	}
	return nil
}

// Get retrieves an import session by ID
func (r ImportSessionRepository) Get(ctx context.Context, sessionID string) (*domain.ImportSession, error) {
	const query = `
      SELECT 
        id, 
        external_system_id, 
        external_system_type, 
        total_images, 
        processed_images, 
        failed_images, 
        status, 
        started_at, 
        completed_at, 
        metadata
      FROM %s
      WHERE id = $1
      LIMIT 1
    `

	row := r.db.QueryRowContext(ctx, r.table(query), sessionID)

	session := &domain.ImportSession{}
	var metadata []byte
	var completedAt sql.NullTime

	if err := row.Scan(
		&session.ID,
		&session.ExternalSystemID,
		&session.ExternalSystemType,
		&session.TotalImages,
		&session.ProcessedImages,
		&session.FailedImages,
		&session.Status,
		&session.StartedAt,
		&completedAt,
		&metadata,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrImportSessionNotFound
		}
		return nil, errors.Wrap(err, "scanning import session")
	}

	if completedAt.Valid {
		session.CompletedAt = completedAt.Time
	}

	// Parse metadata JSON
	if len(metadata) > 0 {
		if err := json.Unmarshal(metadata, &session.Metadata); err != nil {
			return nil, errors.Wrap(err, "unmarshaling metadata")
		}
	} else {
		session.Metadata = make(map[string]string)
	}

	return session, nil
}

// Update modifies an existing import session
func (r ImportSessionRepository) Update(ctx context.Context, session *domain.ImportSession) error {
	const query = `
      UPDATE %s 
      SET 
        total_images = $2, 
        processed_images = $3, 
        failed_images = $4, 
        status = $5, 
        completed_at = $6, 
        metadata = $7
      WHERE id = $1
    `

	metadata, err := json.Marshal(session.Metadata)
	if err != nil {
		return errors.Wrap(err, "marshaling metadata")
	}

	var completedAt sql.NullTime
	if !session.CompletedAt.IsZero() {
		completedAt = sql.NullTime{Time: session.CompletedAt, Valid: true}
	}

	result, err := r.db.ExecContext(ctx, r.table(query),
		session.ID,
		session.TotalImages,
		session.ProcessedImages,
		session.FailedImages,
		session.Status,
		completedAt,
		metadata,
	)
	if err != nil {
		return errors.Wrap(err, "updating import session")
	}

	// Check if any row was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "checking rows affected")
	}

	if rowsAffected == 0 {
		return domain.ErrImportSessionNotFound
	}

	return nil
}

// Helper to format queries
func (r ImportSessionRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
