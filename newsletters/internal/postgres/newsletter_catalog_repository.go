package postgres

import (
	"context"
	"database/sql"
	"fmt"
	
	"middleman/internal/postgres"
	"middleman/newsletters/internal/domain"
	
	"github.com/stackus/errors"
)

type NewsletterCatalogRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.NewsletterCatalogRepository = (*NewsletterCatalogRepository)(nil)

func NewNewsletterCatalogRepository(tableName string, db postgres.DB) NewsletterCatalogRepository {
	return NewsletterCatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r NewsletterCatalogRepository) Find(ctx context.Context, id string) (*domain.CatalogNewsletter, error) {
	const query = `
		SELECT id, user_id, name, description, frequency, category, template_id, 
		       is_active, created_at, updated_at,
		       (SELECT COUNT(*) FROM newsletters.newsletter_subscriptions WHERE newsletter_id = n.id AND status = 'active') as subscriber_count
		FROM %s n
		WHERE id = $1`

	newsletter := &domain.CatalogNewsletter{}
	
	err := r.db.QueryRowContext(ctx, r.table(query), id).Scan(
		&newsletter.ID,
		&newsletter.UserID,
		&newsletter.Name,
		&newsletter.Description,
		&newsletter.Frequency,
		&newsletter.Category,
		&newsletter.TemplateID,
		&newsletter.IsActive,
		&newsletter.CreatedAt,
		&newsletter.UpdatedAt,
		&newsletter.SubscriberCount,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "newsletter not found")
		}
		return nil, errors.Wrap(err, "scanning newsletter")
	}
	
	return newsletter, nil
}

func (r NewsletterCatalogRepository) FindByUser(ctx context.Context, userID string, limit, offset int) ([]*domain.CatalogNewsletter, int, error) {
	const query = `
		SELECT id, user_id, name, description, frequency, category, template_id, 
		       is_active, created_at, updated_at,
		       (SELECT COUNT(*) FROM newsletters.newsletter_subscriptions WHERE newsletter_id = n.id AND status = 'active') as subscriber_count
		FROM %s n
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	
	const countQuery = `SELECT COUNT(*) FROM %s WHERE user_id = $1`

	var total int
	err := r.db.QueryRowContext(ctx, r.table(countQuery), userID).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(err, "counting newsletters")
	}

	rows, err := r.db.QueryContext(ctx, r.table(query), userID, limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying newsletters")
	}
	defer rows.Close()

	var newsletters []*domain.CatalogNewsletter
	for rows.Next() {
		newsletter := &domain.CatalogNewsletter{}
		err := rows.Scan(
			&newsletter.ID,
			&newsletter.UserID,
			&newsletter.Name,
			&newsletter.Description,
			&newsletter.Frequency,
			&newsletter.Category,
			&newsletter.TemplateID,
			&newsletter.IsActive,
			&newsletter.CreatedAt,
			&newsletter.UpdatedAt,
			&newsletter.SubscriberCount,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning newsletter")
		}
		newsletters = append(newsletters, newsletter)
	}

	return newsletters, total, nil
}

func (r NewsletterCatalogRepository) FindByCategory(ctx context.Context, category string, activeOnly bool, limit, offset int) ([]*domain.CatalogNewsletter, int, error) {
	query := `
		SELECT id, user_id, name, description, frequency, category, template_id, 
		       is_active, created_at, updated_at,
		       (SELECT COUNT(*) FROM newsletters.newsletter_subscriptions WHERE newsletter_id = n.id AND status = 'active') as subscriber_count
		FROM %s n
		WHERE category = $1`
	
	countQuery := `SELECT COUNT(*) FROM %s WHERE category = $1`
	
	args := []interface{}{category}
	
	if activeOnly {
		query += " AND is_active = true"
		countQuery += " AND is_active = true"
	}
	
	query += " ORDER BY created_at DESC LIMIT $2 OFFSET $3"

	var total int
	err := r.db.QueryRowContext(ctx, r.table(countQuery), args...).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(err, "counting newsletters")
	}

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, r.table(query), args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying newsletters")
	}
	defer rows.Close()

	var newsletters []*domain.CatalogNewsletter
	for rows.Next() {
		newsletter := &domain.CatalogNewsletter{}
		err := rows.Scan(
			&newsletter.ID,
			&newsletter.UserID,
			&newsletter.Name,
			&newsletter.Description,
			&newsletter.Frequency,
			&newsletter.Category,
			&newsletter.TemplateID,
			&newsletter.IsActive,
			&newsletter.CreatedAt,
			&newsletter.UpdatedAt,
			&newsletter.SubscriberCount,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning newsletter")
		}
		newsletters = append(newsletters, newsletter)
	}

	return newsletters, total, nil
}

func (r NewsletterCatalogRepository) FindAll(ctx context.Context, activeOnly bool, limit, offset int) ([]*domain.CatalogNewsletter, int, error) {
	query := `
		SELECT id, user_id, name, description, frequency, category, template_id, 
		       is_active, created_at, updated_at,
		       (SELECT COUNT(*) FROM newsletters.newsletter_subscriptions WHERE newsletter_id = n.id AND status = 'active') as subscriber_count
		FROM %s n`
	
	countQuery := `SELECT COUNT(*) FROM %s`
	
	if activeOnly {
		query += " WHERE is_active = true"
		countQuery += " WHERE is_active = true"
	}
	
	query += " ORDER BY created_at DESC LIMIT $1 OFFSET $2"

	var total int
	err := r.db.QueryRowContext(ctx, r.table(countQuery)).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(err, "counting newsletters")
	}

	rows, err := r.db.QueryContext(ctx, r.table(query), limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying newsletters")
	}
	defer rows.Close()

	var newsletters []*domain.CatalogNewsletter
	for rows.Next() {
		newsletter := &domain.CatalogNewsletter{}
		err := rows.Scan(
			&newsletter.ID,
			&newsletter.UserID,
			&newsletter.Name,
			&newsletter.Description,
			&newsletter.Frequency,
			&newsletter.Category,
			&newsletter.TemplateID,
			&newsletter.IsActive,
			&newsletter.CreatedAt,
			&newsletter.UpdatedAt,
			&newsletter.SubscriberCount,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning newsletter")
		}
		newsletters = append(newsletters, newsletter)
	}

	return newsletters, total, nil
}

func (r NewsletterCatalogRepository) Add(ctx context.Context, newsletter *domain.CatalogNewsletter) error {
	const query = `
		INSERT INTO %s (id, user_id, name, description, frequency, category, 
		                        template_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.ExecContext(ctx, r.table(query),
		newsletter.ID,
		newsletter.UserID,
		newsletter.Name,
		newsletter.Description,
		newsletter.Frequency,
		newsletter.Category,
		newsletter.TemplateID,
		newsletter.IsActive,
		newsletter.CreatedAt,
		newsletter.UpdatedAt,
	)

	return errors.Wrap(err, "inserting newsletter")
}

func (r NewsletterCatalogRepository) Update(ctx context.Context, newsletter *domain.CatalogNewsletter) error {
	const query = `
		UPDATE %s 
		SET name = $2, description = $3, frequency = $4, category = $5, 
		    template_id = $6, is_active = $7, updated_at = $8
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query),
		newsletter.ID,
		newsletter.Name,
		newsletter.Description,
		newsletter.Frequency,
		newsletter.Category,
		newsletter.TemplateID,
		newsletter.IsActive,
		newsletter.UpdatedAt,
	)

	return errors.Wrap(err, "updating newsletter")
}

func (r NewsletterCatalogRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), id)

	return errors.Wrap(err, "deleting newsletter")
}

func (r NewsletterCatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}