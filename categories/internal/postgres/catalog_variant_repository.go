package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/stackus/errors"
	"middleman/categories/internal/domain"
	"middleman/internal/postgres"
)

type CatalogFilterRepository struct {
	tableName string
	db        postgres.DB
}

// Ensure it implements domain.CatalogFilterRepository
var _ domain.CatalogFilterRepository = (*CatalogFilterRepository)(nil)

// NewCatalogFilterRepository constructs a Postgres repository for the Filter domain entity.
func NewCatalogFilterRepository(tableName string, db postgres.DB) CatalogFilterRepository {
	return CatalogFilterRepository{
		tableName: tableName,
		db:        db,
	}
}

// table is a small helper function to format `tableName` into queries.
func (r CatalogFilterRepository) table(q string) string {
	return fmt.Sprintf(q, r.tableName)
}

// -----------------------------------------------------------------------------
// 1) AddFilter
// -----------------------------------------------------------------------------
func (r CatalogFilterRepository) AddFilter(
	ctx context.Context,
	filterID string,
	categoryID string,
	name string,
	filterType domain.FilterType,
	values []string,
	isActive bool,
) error {
	// We store `values` as JSON text. Adjust as needed.
	valuesData, _ := json.Marshal(values)

	// Possibly store created_at, updated_at, etc.
	const query = `
        INSERT INTO %s (
            id,
            category_id,
            name,
            filter_type,
            values,
            is_active
        )
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		filterID,
		categoryID,
		name,
		string(filterType), // convert the enum to string
		string(valuesData), // store JSON string
		isActive,
	)
	return errors.Wrap(err, "inserting new filter")
}

// -----------------------------------------------------------------------------
// 2) UpdateFilter (Rename, change filterType, change values, etc.)
// -----------------------------------------------------------------------------
func (r CatalogFilterRepository) UpdateFilter(
	ctx context.Context,
	filterID string,
	name string,
	filterType domain.FilterType,
	values []string,
) error {
	valuesData, _ := json.Marshal(values)
	const query = `
        UPDATE %s
        SET 
            name         = $2,
            filter_type  = $3,
            values       = $4,
            updated_at   = NOW()
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		filterID,
		name,
		string(filterType),
		string(valuesData),
	)
	return errors.Wrap(err, "updating filter")
}

// -----------------------------------------------------------------------------
// 3) RemoveFilter
// -----------------------------------------------------------------------------
func (r CatalogFilterRepository) RemoveFilter(
	ctx context.Context,
	filterID string,
	userID string, // if you do not store user ownership, just ignore
) error {
	const query = `
        DELETE FROM %s
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query), filterID)
	return errors.Wrap(err, "removing filter")
}

// -----------------------------------------------------------------------------
// 4) ArchiveFilter (set isActive = false)
// -----------------------------------------------------------------------------
func (r CatalogFilterRepository) ArchiveFilter(ctx context.Context, filterID string) error {
	const query = `
        UPDATE %s
        SET is_active = FALSE,
            updated_at = NOW()
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query), filterID)
	return errors.Wrap(err, "archiving filter")
}

// -----------------------------------------------------------------------------
// 5) RebrandFilter (rename only, ignoring domain details if you want a 2nd approach)
// -----------------------------------------------------------------------------
func (r CatalogFilterRepository) RebrandFilter(ctx context.Context, filterID string, newName string) error {
	const query = `
        UPDATE %s
        SET name = $2,
            updated_at = NOW()
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query), filterID, newName)
	return errors.Wrap(err, "rebranding filter")
}

// -----------------------------------------------------------------------------
// 6) FindFilter
// -----------------------------------------------------------------------------
func (r CatalogFilterRepository) FindFilter(ctx context.Context, filterID string) (*domain.CatalogFilter, error) {
	const query = `
        SELECT
            category_id,
            name,
            filter_type,
            values,
            is_active
        FROM %s
        WHERE id = $1
        LIMIT 1
    `
	row := r.db.QueryRowContext(ctx, r.table(query), filterID)

	cv := &domain.CatalogFilter{ID: filterID}
	var (
		filterTypeStr string
		valuesStr     string
		active        bool
	)
	err := row.Scan(
		&cv.CategoryID,
		&cv.Name,
		&filterTypeStr,
		&valuesStr,
		&active,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // or return a domain-specific error
		}
		return nil, errors.Wrap(err, "finding filter in DB")
	}

	cv.FilterType = domain.FilterType(filterTypeStr)
	cv.IsActive = active

	// parse JSON array of strings
	var vals []string
	_ = json.Unmarshal([]byte(valuesStr), &vals)
	cv.Values = vals

	return cv, nil
}

// -----------------------------------------------------------------------------
// 7) GetFilters
// -----------------------------------------------------------------------------
func (r CatalogFilterRepository) GetFilters(
	ctx context.Context,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogFilter, int64, error) {
	if sortBy == "" {
		sortBy = "name"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	offset := (page - 1) * pageSize

	// Count
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, r.tableName)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQ).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, "counting filters")
	}

	// Query
	q := fmt.Sprintf(`
        SELECT
            id,
            category_id,
            name,
            filter_type,
            values,
            is_active
        FROM %s
        ORDER BY %s %s
        LIMIT $1 OFFSET $2
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, q, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying filters")
	}
	defer rows.Close()

	var results []*domain.CatalogFilter
	for rows.Next() {
		cf := &domain.CatalogFilter{}
		var (
			filterTypeStr string
			valuesStr     string
			active        bool
		)
		scanErr := rows.Scan(
			&cf.ID,
			&cf.CategoryID,
			&cf.Name,
			&filterTypeStr,
			&valuesStr,
			&active,
		)
		if scanErr != nil {
			return nil, 0, errors.Wrap(scanErr, "scanning filter row")
		}

		cf.FilterType = domain.FilterType(filterTypeStr)
		cf.IsActive = active

		var vals []string
		_ = json.Unmarshal([]byte(valuesStr), &vals)
		cf.Values = vals

		results = append(results, cf)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "finalizing filter row results")
	}
	return results, total, nil
}

// -----------------------------------------------------------------------------
// 8) GetFiltersByCategory
// -----------------------------------------------------------------------------
func (r CatalogFilterRepository) GetFiltersByCategory(
	ctx context.Context,
	categoryID string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogFilter, int64, error) {
	if sortBy == "" {
		sortBy = "name"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	offset := (page - 1) * pageSize

	// Count
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE category_id = $1`, r.tableName)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQ, categoryID).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, "counting filters by category")
	}

	// Query
	q := fmt.Sprintf(`
        SELECT
            id,
            category_id,
            name,
            filter_type,
            values,
            is_active
        FROM %s
        WHERE category_id = $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, q, categoryID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying filters by category")
	}
	defer rows.Close()

	var results []*domain.CatalogFilter
	for rows.Next() {
		cf := &domain.CatalogFilter{}
		var (
			filterTypeStr string
			valuesStr     string
			active        bool
		)
		scanErr := rows.Scan(
			&cf.ID,
			&cf.CategoryID,
			&cf.Name,
			&filterTypeStr,
			&valuesStr,
			&active,
		)
		if scanErr != nil {
			return nil, 0, errors.Wrap(scanErr, "scanning filter row by category")
		}

		cf.FilterType = domain.FilterType(filterTypeStr)
		cf.IsActive = active

		var vals []string
		_ = json.Unmarshal([]byte(valuesStr), &vals)
		cf.Values = vals

		results = append(results, cf)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "finalizing filter results by category")
	}
	return results, total, nil
}
