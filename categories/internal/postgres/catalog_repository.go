package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackus/errors"
	"middleman/categories/internal/domain"
	"middleman/internal/postgres"
)

// CatalogRepository implements domain.CatalogRepository using a Postgres table.
type CatalogRepository struct {
	tableName string
	db        postgres.DB
}

// Compile-time check that CatalogRepository satisfies domain.CatalogRepository.
var _ domain.CatalogRepository = (*CatalogRepository)(nil)

// NewCatalogRepository constructs a Postgres-based CatalogRepository.
func NewCatalogRepository(tableName string, db postgres.DB) CatalogRepository {
	return CatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

// -----------------------------------------------------------------------------
// 1) AddCategory
// -----------------------------------------------------------------------------
func (r CatalogRepository) AddCategory(
	ctx context.Context,
	categoryID string,
	description string,
	parentID string,
	googleCategoryID string,
	tags []string,
	isActive bool,
	slug string,
	seoTitle string,
	seoKeywords []string,
	seoDesc string,
	categoryType string,
	lang string,
) error {
	const q = `
		INSERT INTO %s (
			category_id,
			description,
			parent_id,
			google_category_id,
			tags,
			slug,
			is_active,
			seo_title,
			seo_keywords,
			seo_desc,
			category_type,
			lang
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`
	_, err := r.db.ExecContext(ctx, r.table(q),
		categoryID,
		description,
		parentID,
		googleCategoryID,
		sliceToString(tags),        // $5
		slug,                       // $6
		isActive,                   // $7
		seoTitle,                   // $8
		sliceToString(seoKeywords), // $9
		seoDesc,                    // $10
		categoryType,               // $11
		lang,                       // $12
	)
	return errors.Wrap(err, "inserting new category")
}

// -----------------------------------------------------------------------------
// 2) UpdateCategory
// -----------------------------------------------------------------------------
func (r CatalogRepository) UpdateCategory(
	ctx context.Context,
	categoryID string,
	description string,
	parentID string,
	googleCategoryID string,
	tags []string,
	isActive bool,
	slug string,
	seoTitle string,
	seoKeywords []string,
	seoDesc string,
	categoryType string,
	lang string,
) error {
	const q = `
		UPDATE %s
		SET
			description        = $2,
			parent_id          = $3,
			google_category_id = $4,
			tags               = $5,
			is_active          = $6,
			slug               = $7,
			seo_title          = $8,
			seo_keywords       = $9,
			seo_desc           = $10,
			category_type      = $11,
			lang               = $12,
			updated_at         = NOW()
		WHERE category_id = $1 AND lang = $12
	`
	_, err := r.db.ExecContext(ctx, r.table(q),
		categoryID,
		description,
		parentID,
		googleCategoryID,
		sliceToString(tags),        // $5
		isActive,                   // $6
		slug,                       // $7
		seoTitle,                   // $8
		sliceToString(seoKeywords), // $9
		seoDesc,                    // $10
		categoryType,               // $11
		lang,                       // $12
	)
	return errors.Wrap(err, "updating category")
}

// -----------------------------------------------------------------------------
// 3) RemoveCategory
// -----------------------------------------------------------------------------
func (r CatalogRepository) RemoveCategory(
	ctx context.Context,
	categoryID string,
	userID string, // not used yet
) error {
	const q = `
		DELETE FROM %s
		WHERE category_id = $1
	`
	_, err := r.db.ExecContext(ctx, r.table(q), categoryID)
	return errors.Wrap(err, "removing category")
}

// -----------------------------------------------------------------------------
// 4) ArchiveCategory
// -----------------------------------------------------------------------------
func (r CatalogRepository) ArchiveCategory(
	ctx context.Context,
	categoryID string,
	userID string, // not used yet
) error {
	const q = `
		UPDATE %s
		SET is_active = FALSE,
			updated_at = NOW()
		WHERE category_id = $1
	`
	_, err := r.db.ExecContext(ctx, r.table(q), categoryID)
	return errors.Wrap(err, "archiving category")
}

// -----------------------------------------------------------------------------
// 5) RebrandCategory
// -----------------------------------------------------------------------------
func (r CatalogRepository) RebrandCategory(
	ctx context.Context,
	categoryID string,
	newSlug string,
	newDescription string,
) error {
	const q = `
		UPDATE %s
		SET
			slug        = $2,
			description = $3,
			updated_at  = NOW()
		WHERE category_id = $1
	`
	_, err := r.db.ExecContext(ctx, r.table(q), categoryID, newSlug, newDescription)
	return errors.Wrap(err, "rebranding category")
}

// -----------------------------------------------------------------------------
// 6) Find
// -----------------------------------------------------------------------------
func (r CatalogRepository) Find(
	ctx context.Context,
	categoryID string,
) (*domain.CatalogCategory, error) {
	const q = `
		SELECT
			category_id,
			description,
			parent_id,
			google_category_id,
			is_active,
			slug,
			seo_title,
			seo_keywords,
			seo_desc,
			category_type,
			lang
		FROM %s
		WHERE category_id = $1
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, r.table(q), categoryID)

	cc := &domain.CatalogCategory{}
	var keywordsStr string
	if err := row.Scan(
		&cc.ID,
		&cc.Description,
		&cc.ParentID,
		&cc.GoogleCategoryID,
		&cc.IsActive,
		&cc.Slug,
		&cc.SeoTitle,
		&keywordsStr,
		&cc.SeoDesc,
		&cc.CategoryType,
		&cc.Lang,
	); err != nil {
		return nil, errors.Wrap(err, "scanning category from DB")
	}
	cc.SeoKeywords = stringToSlice(keywordsStr)
	return cc, nil
}

// -----------------------------------------------------------------------------
// 7a) GetCategories
// -----------------------------------------------------------------------------
func (r CatalogRepository) GetCategories(
	ctx context.Context,
	categoryType string,
	lang string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogCategory, int64, error) {
	validSorts := map[string]bool{"slug": true, "updated_at": true, "category_id": true}
	if !validSorts[sortBy] {
		sortBy = "slug"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	offset := (page - 1) * pageSize

	// total count
	countQ := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE lang = $1 AND category_type = $2
	`, r.tableName)

	var total int64
	if err := r.db.QueryRowContext(ctx, countQ, lang, categoryType).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, "counting categories")
	}

	// page query
	q := fmt.Sprintf(`
		SELECT
			category_id,
			description,
			google_category_id,
			is_active,
			slug,
			seo_title,
			seo_keywords,
			seo_desc,
			category_type,
			lang
		FROM %s
		WHERE lang = $1 AND category_type = $2
		ORDER BY %s %s
		LIMIT $3
		OFFSET $4
	`, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, q, lang, categoryType, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying categories")
	}
	defer rows.Close()

	var cats []*domain.CatalogCategory
	for rows.Next() {
		c := &domain.CatalogCategory{}
		var keywordsStr string
		if scanErr := rows.Scan(
			&c.ID,
			&c.Description,
			&c.GoogleCategoryID,
			&c.IsActive,
			&c.Slug,
			&c.SeoTitle,
			&keywordsStr,
			&c.SeoDesc,
			&c.CategoryType,
			&c.Lang,
		); scanErr != nil {
			return nil, 0, errors.Wrap(scanErr, "scanning category row")
		}
		c.SeoKeywords = stringToSlice(keywordsStr)
		cats = append(cats, c)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading category rows final")
	}
	return cats, total, nil
}

// -----------------------------------------------------------------------------
// 7b) GetMainCategories
// -----------------------------------------------------------------------------
func (r CatalogRepository) GetMainCategories(
	ctx context.Context,
	categoryType string,
	lang string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogCategory, int64, error) {
	validSorts := map[string]bool{"slug": true, "updated_at": true, "category_id": true}
	if !validSorts[sortBy] {
		sortBy = "slug"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	offset := (page - 1) * pageSize

	countQ := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE lang = $1
		  AND category_type = $2
		  AND (parent_id IS NULL OR parent_id = '')
	`, r.tableName)

	var total int64
	if err := r.db.QueryRowContext(ctx, countQ, lang, categoryType).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, "counting main categories")
	}

	q := fmt.Sprintf(`
		SELECT
			category_id,
			description,
			google_category_id,
			is_active,
			slug,
			seo_title,
			seo_keywords,
			seo_desc,
			category_type,
			lang
		FROM %s
		WHERE lang = $1
		  AND category_type = $2
		  AND (parent_id IS NULL OR parent_id = '')
		ORDER BY %s %s
		LIMIT $3
		OFFSET $4
	`, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, q, lang, categoryType, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying main categories")
	}
	defer rows.Close()

	var cats []*domain.CatalogCategory
	for rows.Next() {
		c := &domain.CatalogCategory{}
		var keywordsStr string
		if scanErr := rows.Scan(
			&c.ID,
			&c.Description,
			&c.GoogleCategoryID,
			&c.IsActive,
			&c.Slug,
			&c.SeoTitle,
			&keywordsStr,
			&c.SeoDesc,
			&c.CategoryType,
			&c.Lang,
		); scanErr != nil {
			return nil, 0, errors.Wrap(scanErr, "scanning main category row")
		}
		c.SeoKeywords = stringToSlice(keywordsStr)
		cats = append(cats, c)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading main category rows final")
	}
	return cats, total, nil
}

// -----------------------------------------------------------------------------
// 7b) GetMainCategories
// -----------------------------------------------------------------------------
func (r CatalogRepository) GetAllMainCategories(
	ctx context.Context,
	lang string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogCategory, int64, error) {
	validSorts := map[string]bool{"slug": true, "updated_at": true, "category_id": true}
	if !validSorts[sortBy] {
		sortBy = "slug"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	offset := (page - 1) * pageSize

	countQ := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE lang = $1
		  AND (parent_id IS NULL OR parent_id = '')
	`, r.tableName)

	var total int64
	if err := r.db.QueryRowContext(ctx, countQ, lang).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, "counting main categories")
	}

	q := fmt.Sprintf(`
		SELECT
			category_id,
			description,
			google_category_id,
			is_active,
			slug,
			seo_title,
			seo_keywords,
			seo_desc,
			category_type,
			lang
		FROM %s
		WHERE lang = $1
		  AND (parent_id IS NULL OR parent_id = '')
		ORDER BY %s %s
		LIMIT $2
		OFFSET $3
	`, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, q, lang, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying main categories")
	}
	defer rows.Close()

	var cats []*domain.CatalogCategory
	for rows.Next() {
		c := &domain.CatalogCategory{}
		var keywordsStr string
		if scanErr := rows.Scan(
			&c.ID,
			&c.Description,
			&c.GoogleCategoryID,
			&c.IsActive,
			&c.Slug,
			&c.SeoTitle,
			&keywordsStr,
			&c.SeoDesc,
			&c.CategoryType,
			&c.Lang,
		); scanErr != nil {
			return nil, 0, errors.Wrap(scanErr, "scanning main category row")
		}
		c.SeoKeywords = stringToSlice(keywordsStr)
		cats = append(cats, c)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading main category rows final")
	}
	return cats, total, nil
}

// -----------------------------------------------------------------------------
// 7c) GetSubCategories
// -----------------------------------------------------------------------------
func (r CatalogRepository) GetSubCategories(
	ctx context.Context,
	lang string,
	parentCategoryID string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogCategory, int64, error) {
	validSorts := map[string]bool{"slug": true, "updated_at": true, "category_id": true}
	if !validSorts[sortBy] {
		sortBy = "slug"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	offset := (page - 1) * pageSize

	countQ := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE lang = $1 AND parent_id = $2
	`, r.tableName)

	var total int64
	if err := r.db.QueryRowContext(ctx, countQ, lang, parentCategoryID).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, "counting subcategories")
	}

	q := fmt.Sprintf(`
		SELECT
			category_id,
			description,
			parent_id,
			google_category_id,
			is_active,
			slug,
			seo_title,
			seo_keywords,
			seo_desc,
			category_type,
			lang
		FROM %s
		WHERE lang = $1 AND parent_id = $2
		ORDER BY %s %s
		LIMIT $3
		OFFSET $4
	`, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, q, lang, parentCategoryID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying subcategories")
	}
	defer rows.Close()

	var cats []*domain.CatalogCategory
	for rows.Next() {
		c := &domain.CatalogCategory{}
		var keywordsStr string
		if scanErr := rows.Scan(
			&c.ID,
			&c.Description,
			&c.ParentID,
			&c.GoogleCategoryID,
			&c.IsActive,
			&c.Slug,
			&c.SeoTitle,
			&keywordsStr,
			&c.SeoDesc,
			&c.CategoryType,
			&c.Lang,
		); scanErr != nil {
			return nil, 0, errors.Wrap(scanErr, "scanning subcategory row")
		}
		c.SeoKeywords = stringToSlice(keywordsStr)
		cats = append(cats, c)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading subcategory rows final")
	}
	return cats, total, nil
}

// -----------------------------------------------------------------------------
// 8) GetCatalog
// -----------------------------------------------------------------------------
func (r CatalogRepository) GetCatalog(
	ctx context.Context,
	lang string,
	userID string, // reserved for future ownership filtering
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogCategory, int64, error) {
	validSorts := map[string]bool{"slug": true, "updated_at": true, "category_id": true}
	if !validSorts[sortBy] {
		sortBy = "slug"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	offset := (page - 1) * pageSize

	countQ := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE lang = $1 AND is_active = TRUE
	`, r.tableName)

	var total int64
	if err := r.db.QueryRowContext(ctx, countQ, lang).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, "counting catalog categories")
	}

	q := fmt.Sprintf(`
		SELECT
			category_id,
			description,
			parent_id,
			google_category_id,
			is_active,
			slug,
			seo_title,
			seo_keywords,
			seo_desc,
			category_type,
			lang
		FROM %s
		WHERE lang = $1 AND is_active = TRUE
		ORDER BY %s %s
		LIMIT $2
		OFFSET $3
	`, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, q, lang, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying catalog")
	}
	defer rows.Close()

	var cats []*domain.CatalogCategory
	for rows.Next() {
		c := &domain.CatalogCategory{}
		var keywordsStr string
		if scanErr := rows.Scan(
			&c.ID,
			&c.Description,
			&c.ParentID,
			&c.GoogleCategoryID,
			&c.IsActive,
			&c.Slug,
			&c.SeoTitle,
			&keywordsStr,
			&c.SeoDesc,
			&c.CategoryType,
			&c.Lang,
		); scanErr != nil {
			return nil, 0, errors.Wrap(scanErr, "scanning catalog row")
		}
		c.SeoKeywords = stringToSlice(keywordsStr)
		cats = append(cats, c)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "reading catalog rows final")
	}
	return cats, total, nil
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

// table is a tiny helper that substitutes the table name.
func (r CatalogRepository) table(q string) string {
	return fmt.Sprintf(q, r.tableName)
}

// Converts a slice to a single comma-separated string ("red,blue,yellow").
func sliceToString(sl []string) string {
	return strings.Join(sl, ",")
}

// Parses the comma-separated representation back into a slice.
func stringToSlice(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
