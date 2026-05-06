package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	
	"middleman/internal/postgres"
	"middleman/newsletters/internal/domain"
	
	"github.com/stackus/errors"
)

type TemplateCatalogRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.TemplateCatalogRepository = (*TemplateCatalogRepository)(nil)

func NewTemplateCatalogRepository(tableName string, db postgres.DB) TemplateCatalogRepository {
	return TemplateCatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r TemplateCatalogRepository) Find(ctx context.Context, id string) (*domain.CatalogTemplate, error) {
	const query = `
		SELECT id, user_id, name, description, html_template, text_template, 
		       variables, preview_data, is_public, created_at, updated_at
		FROM %s
		WHERE id = $1`

	template := &domain.CatalogTemplate{}
	var variables, previewData []byte
	
	err := r.db.QueryRowContext(ctx, r.table(query), id).Scan(
		&template.ID,
		&template.UserID,
		&template.Name,
		&template.Description,
		&template.HTMLTemplate,
		&template.TextTemplate,
		&variables,
		&previewData,
		&template.IsPublic,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "template not found")
		}
		return nil, errors.Wrap(err, "scanning template")
	}

	// Parse JSON fields
	if len(variables) > 0 {
		if err := json.Unmarshal(variables, &template.Variables); err != nil {
			return nil, errors.Wrap(err, "unmarshaling variables")
		}
	} else {
		template.Variables = make(map[string]string)
	}

	if len(previewData) > 0 {
		if err := json.Unmarshal(previewData, &template.PreviewData); err != nil {
			return nil, errors.Wrap(err, "unmarshaling preview data")
		}
	} else {
		template.PreviewData = make(map[string]string)
	}
	
	return template, nil
}

func (r TemplateCatalogRepository) FindByUser(ctx context.Context, userID string, limit, offset int) ([]*domain.CatalogTemplate, int, error) {
	const query = `
		SELECT id, user_id, name, description, html_template, text_template, 
		       variables, preview_data, is_public, created_at, updated_at
		FROM %s
		WHERE user_id = $1 OR is_public = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	
	const countQuery = `SELECT COUNT(*) FROM %s WHERE user_id = $1 OR is_public = true`

	var total int
	err := r.db.QueryRowContext(ctx, r.table(countQuery), userID).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(err, "counting templates")
	}

	rows, err := r.db.QueryContext(ctx, r.table(query), userID, limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying templates")
	}
	defer rows.Close()

	var templates []*domain.CatalogTemplate
	for rows.Next() {
		template := &domain.CatalogTemplate{}
		var variables, previewData []byte
		
		err := rows.Scan(
			&template.ID,
			&template.UserID,
			&template.Name,
			&template.Description,
			&template.HTMLTemplate,
			&template.TextTemplate,
			&variables,
			&previewData,
			&template.IsPublic,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning template")
		}

		// Parse JSON fields
		if len(variables) > 0 {
			if err := json.Unmarshal(variables, &template.Variables); err != nil {
				template.Variables = make(map[string]string)
			}
		} else {
			template.Variables = make(map[string]string)
		}

		if len(previewData) > 0 {
			if err := json.Unmarshal(previewData, &template.PreviewData); err != nil {
				template.PreviewData = make(map[string]string)
			}
		} else {
			template.PreviewData = make(map[string]string)
		}
		
		templates = append(templates, template)
	}

	return templates, total, nil
}

func (r TemplateCatalogRepository) FindPublic(ctx context.Context, limit, offset int) ([]*domain.CatalogTemplate, int, error) {
	const query = `
		SELECT id, user_id, name, description, html_template, text_template, 
		       variables, preview_data, is_public, created_at, updated_at
		FROM %s
		WHERE is_public = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	
	const countQuery = `SELECT COUNT(*) FROM %s WHERE is_public = true`

	var total int
	err := r.db.QueryRowContext(ctx, r.table(countQuery)).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(err, "counting templates")
	}

	rows, err := r.db.QueryContext(ctx, r.table(query), limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying templates")
	}
	defer rows.Close()

	var templates []*domain.CatalogTemplate
	for rows.Next() {
		template := &domain.CatalogTemplate{}
		var variables, previewData []byte
		
		err := rows.Scan(
			&template.ID,
			&template.UserID,
			&template.Name,
			&template.Description,
			&template.HTMLTemplate,
			&template.TextTemplate,
			&variables,
			&previewData,
			&template.IsPublic,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning template")
		}

		// Parse JSON fields
		if len(variables) > 0 {
			if err := json.Unmarshal(variables, &template.Variables); err != nil {
				template.Variables = make(map[string]string)
			}
		} else {
			template.Variables = make(map[string]string)
		}

		if len(previewData) > 0 {
			if err := json.Unmarshal(previewData, &template.PreviewData); err != nil {
				template.PreviewData = make(map[string]string)
			}
		} else {
			template.PreviewData = make(map[string]string)
		}
		
		templates = append(templates, template)
	}

	return templates, total, nil
}

func (r TemplateCatalogRepository) Add(ctx context.Context, template *domain.CatalogTemplate) error {
	variables, err := json.Marshal(template.Variables)
	if err != nil {
		return errors.Wrap(err, "marshaling variables")
	}

	previewData, err := json.Marshal(template.PreviewData)
	if err != nil {
		return errors.Wrap(err, "marshaling preview data")
	}

	const query = `
		INSERT INTO %s (id, user_id, name, description, html_template, text_template, 
		                                 variables, preview_data, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = r.db.ExecContext(ctx, r.table(query),
		template.ID,
		template.UserID,
		template.Name,
		template.Description,
		template.HTMLTemplate,
		template.TextTemplate,
		variables,
		previewData,
		template.IsPublic,
		template.CreatedAt,
		template.UpdatedAt,
	)

	return errors.Wrap(err, "inserting template")
}

func (r TemplateCatalogRepository) Update(ctx context.Context, template *domain.CatalogTemplate) error {
	variables, err := json.Marshal(template.Variables)
	if err != nil {
		return errors.Wrap(err, "marshaling variables")
	}

	previewData, err := json.Marshal(template.PreviewData)
	if err != nil {
		return errors.Wrap(err, "marshaling preview data")
	}

	const query = `
		UPDATE %s 
		SET name = $2, description = $3, html_template = $4, text_template = $5,
		    variables = $6, preview_data = $7, is_public = $8, updated_at = $9
		WHERE id = $1`

	_, err = r.db.ExecContext(ctx, r.table(query),
		template.ID,
		template.Name,
		template.Description,
		template.HTMLTemplate,
		template.TextTemplate,
		variables,
		previewData,
		template.IsPublic,
		template.UpdatedAt,
	)

	return errors.Wrap(err, "updating template")
}

func (r TemplateCatalogRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), id)

	return errors.Wrap(err, "deleting template")
}

func (r TemplateCatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}