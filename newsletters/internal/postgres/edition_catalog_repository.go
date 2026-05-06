package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	
	"middleman/internal/postgres"
	"middleman/newsletters/internal/domain"
	
	"github.com/stackus/errors"
)

type EditionCatalogRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.EditionCatalogRepository = (*EditionCatalogRepository)(nil)

func NewEditionCatalogRepository(tableName string, db postgres.DB) EditionCatalogRepository {
	return EditionCatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r EditionCatalogRepository) Find(ctx context.Context, id string) (*domain.CatalogEdition, error) {
	const query = `
		SELECT id, newsletter_id, subject, content_html, content_text, template_data,
		       scheduled_at, sent_at, status, created_by, created_at, updated_at
		FROM %s
		WHERE id = $1`

	edition := &domain.CatalogEdition{}
	var templateData []byte
	
	err := r.db.QueryRowContext(ctx, r.table(query), id).Scan(
		&edition.ID,
		&edition.NewsletterID,
		&edition.Subject,
		&edition.ContentHTML,
		&edition.ContentText,
		&templateData,
		&edition.ScheduledAt,
		&edition.SentAt,
		&edition.Status,
		&edition.CreatedBy,
		&edition.CreatedAt,
		&edition.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "edition not found")
		}
		return nil, errors.Wrap(err, "scanning edition")
	}

	// Parse template data JSON
	if len(templateData) > 0 {
		if err := json.Unmarshal(templateData, &edition.TemplateData); err != nil {
			return nil, errors.Wrap(err, "unmarshaling template data")
		}
	} else {
		edition.TemplateData = make(map[string]string)
	}

	// Get recipient count from send logs
	const countQuery = `SELECT COUNT(DISTINCT user_id) FROM newsletters.newsletter_send_logs WHERE edition_id = $1`
	_ = r.db.QueryRowContext(ctx, countQuery, id).Scan(&edition.RecipientCount)
	
	return edition, nil
}

func (r EditionCatalogRepository) FindByNewsletter(ctx context.Context, newsletterID string, status string, limit, offset int) ([]*domain.CatalogEdition, int, error) {
	query := `
		SELECT id, newsletter_id, subject, content_html, content_text, template_data,
		       scheduled_at, sent_at, status, created_by, created_at, updated_at
		FROM %s
		WHERE newsletter_id = $1`
	
	countQuery := `SELECT COUNT(*) FROM %s WHERE newsletter_id = $1`
	
	args := []interface{}{newsletterID}
	
	if status != "" {
		query += " AND status = $2"
		countQuery += " AND status = $2"
		args = append(args, status)
	}
	
	// Add placeholders for limit and offset
	placeholderOffset := len(args) + 1
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", placeholderOffset, placeholderOffset+1)

	var total int
	err := r.db.QueryRowContext(ctx, r.table(countQuery), args...).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(err, "counting editions")
	}

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, r.table(query), args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying editions")
	}
	defer rows.Close()

	var editions []*domain.CatalogEdition
	for rows.Next() {
		edition := &domain.CatalogEdition{}
		var templateData []byte
		
		err := rows.Scan(
			&edition.ID,
			&edition.NewsletterID,
			&edition.Subject,
			&edition.ContentHTML,
			&edition.ContentText,
			&templateData,
			&edition.ScheduledAt,
			&edition.SentAt,
			&edition.Status,
			&edition.CreatedBy,
			&edition.CreatedAt,
			&edition.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning edition")
		}

		// Parse template data JSON
		if len(templateData) > 0 {
			if err := json.Unmarshal(templateData, &edition.TemplateData); err != nil {
				edition.TemplateData = make(map[string]string)
			}
		} else {
			edition.TemplateData = make(map[string]string)
		}

		// Get recipient count from send logs
		const countQuery = `SELECT COUNT(DISTINCT user_id) FROM newsletters.newsletter_send_logs WHERE edition_id = $1`
		_ = r.db.QueryRowContext(ctx, countQuery, edition.ID).Scan(&edition.RecipientCount)
		
		editions = append(editions, edition)
	}

	return editions, total, nil
}

func (r EditionCatalogRepository) FindScheduled(ctx context.Context, before time.Time) ([]*domain.CatalogEdition, error) {
	const query = `
		SELECT id, newsletter_id, subject, content_html, content_text, template_data,
		       scheduled_at, sent_at, status, created_by, created_at, updated_at
		FROM %s
		WHERE status = 'scheduled' AND scheduled_at <= $1
		ORDER BY scheduled_at ASC`

	rows, err := r.db.QueryContext(ctx, r.table(query), before)
	if err != nil {
		return nil, errors.Wrap(err, "querying scheduled editions")
	}
	defer rows.Close()

	var editions []*domain.CatalogEdition
	for rows.Next() {
		edition := &domain.CatalogEdition{}
		var templateData []byte
		
		err := rows.Scan(
			&edition.ID,
			&edition.NewsletterID,
			&edition.Subject,
			&edition.ContentHTML,
			&edition.ContentText,
			&templateData,
			&edition.ScheduledAt,
			&edition.SentAt,
			&edition.Status,
			&edition.CreatedBy,
			&edition.CreatedAt,
			&edition.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "scanning edition")
		}

		// Parse template data JSON
		if len(templateData) > 0 {
			if err := json.Unmarshal(templateData, &edition.TemplateData); err != nil {
				edition.TemplateData = make(map[string]string)
			}
		} else {
			edition.TemplateData = make(map[string]string)
		}
		
		editions = append(editions, edition)
	}

	return editions, nil
}

func (r EditionCatalogRepository) Add(ctx context.Context, edition *domain.CatalogEdition) error {
	templateData, err := json.Marshal(edition.TemplateData)
	if err != nil {
		return errors.Wrap(err, "marshaling template data")
	}

	const query = `
		INSERT INTO %s (id, newsletter_id, subject, content_html, content_text, 
		                                template_data, scheduled_at, sent_at, status, created_by, 
		                                created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err = r.db.ExecContext(ctx, r.table(query),
		edition.ID,
		edition.NewsletterID,
		edition.Subject,
		edition.ContentHTML,
		edition.ContentText,
		templateData,
		edition.ScheduledAt,
		edition.SentAt,
		edition.Status,
		edition.CreatedBy,
		edition.CreatedAt,
		edition.UpdatedAt,
	)

	return errors.Wrap(err, "inserting edition")
}

func (r EditionCatalogRepository) Update(ctx context.Context, edition *domain.CatalogEdition) error {
	templateData, err := json.Marshal(edition.TemplateData)
	if err != nil {
		return errors.Wrap(err, "marshaling template data")
	}

	const query = `
		UPDATE %s 
		SET subject = $2, content_html = $3, content_text = $4, template_data = $5,
		    scheduled_at = $6, sent_at = $7, status = $8, updated_at = $9
		WHERE id = $1`

	_, err = r.db.ExecContext(ctx, r.table(query),
		edition.ID,
		edition.Subject,
		edition.ContentHTML,
		edition.ContentText,
		templateData,
		edition.ScheduledAt,
		edition.SentAt,
		edition.Status,
		edition.UpdatedAt,
	)

	return errors.Wrap(err, "updating edition")
}

func (r EditionCatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}