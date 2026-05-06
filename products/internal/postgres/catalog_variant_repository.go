package postgres

import (
	"context"
	"fmt"

	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/products/internal/domain"
)

type CatalogVariantRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.CatalogVariantRepository = (*CatalogVariantRepository)(nil)

func NewCatalogVariantRepository(tableName string, db postgres.DB) CatalogVariantRepository {
	return CatalogVariantRepository{
		tableName: tableName,
		db:        db,
	}
}

// -----------------------------------------------------------------------------
// 1) AddVariant
// -----------------------------------------------------------------------------
func (r CatalogVariantRepository) AddVariant(
	ctx context.Context,
	variantID, productID string,
	status domain.ProductStatus,
	sku, barcode string,
	condition domain.ProductCondition,
	variantPrice int64,
	currencyCode string,
	stock, weight, height, width, depth int64,
	attributes []domain.Attribute,
	isAvailable bool,
	hasOptions bool,
	options []domain.Option,
) error {
	// Convert slices to strings for storage
	attrStr := attributesToString(attributes)
	optsStr := optionsToString(options)

	const query = `
        INSERT INTO %s (
            id,
            product_id,
            status,
            sku,
            barcode,
            condition,
            variant_price,
            currency_code,
            stock,
            weight,
            height,
            width,
            depth,
            attributes,
            is_available,
            has_options,
            options
        ) VALUES (
            $1,$2,$3,$4,$5,$6,
            $7,$8,$9,$10,$11,
            $12,$13,$14,$15,$16,$17
        )
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		variantID,
		productID,
		status, // <-- store status as text or enum
		sku,
		barcode,
		condition, // <-- store condition as text or enum
		variantPrice,
		currencyCode,
		stock,
		weight,
		height,
		width,
		depth,
		attrStr,
		isAvailable,
		hasOptions,
		optsStr,
	)
	return err
}

// -----------------------------------------------------------------------------
// 2) UpdateVariantPrice
// -----------------------------------------------------------------------------
func (r CatalogVariantRepository) UpdateVariantPrice(ctx context.Context, variantID string, oldPrice, newPrice int64) error {
	const query = `
        UPDATE %s
        SET variant_price = $3,
            updated_at    = NOW()
        WHERE id = $1
          AND variant_price = $2
    `
	res, err := r.db.ExecContext(ctx, r.table(query), variantID, oldPrice, newPrice)
	if err != nil {
		return errors.Wrap(err, "updating variant price")
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.Wrap(err, "no variant updated; old_price mismatch or variant not found")
	}
	return nil
}

// -----------------------------------------------------------------------------
// 3) RemoveVariant
// -----------------------------------------------------------------------------
func (r CatalogVariantRepository) RemoveVariant(ctx context.Context, variantID, userSellerID string) error {
	// If your schema doesn't store userSellerID in variants, remove that part.
	// Otherwise, add "AND user_seller_id = $2" (plus an argument).
	const query = `
        DELETE FROM %s
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query), variantID /*, userSellerID */)
	return err
}

// -----------------------------------------------------------------------------
// 4) FindVariant
// -----------------------------------------------------------------------------
func (r CatalogVariantRepository) FindVariant(ctx context.Context, variantID string) (*domain.CatalogVariant, error) {
	const query = `
        SELECT
            product_id,
            status,
            sku,
            barcode,
            condition,
            variant_price,
            currency_code,
            stock,
            weight,
            height,
            width,
            depth,
            attributes,
            is_available,
            has_options,
            options
        FROM %s
        WHERE id = $1
        LIMIT 1
    `
	row := r.db.QueryRowContext(ctx, r.table(query), variantID)

	cv := &domain.CatalogVariant{ID: variantID}
	var (
		statusStr string
		condStr   string
		attrStr   string
		optsStr   string
	)

	err := row.Scan(
		&cv.ProductID,
		&statusStr,
		&cv.SKU,
		&cv.Barcode,
		&condStr,
		&cv.VariantPrice,
		&cv.CurrencyCode,
		&cv.Stock,
		&cv.Weight,
		&cv.Height,
		&cv.Width,
		&cv.Depth,
		&attrStr,
		&cv.IsAvailable,
		&cv.HasOptions,
		&optsStr,
	)
	if err != nil {
		return nil, errors.Wrap(err, "finding variant")
	}

	// Convert textual status/condition to domain enums
	cv.Status = domain.ToProductStatus(statusStr)
	cv.Condition = domain.ToProductCondition(condStr)

	// Convert stored strings -> slices
	cv.Attributes = stringToAttributes(attrStr)
	cv.Options = stringToOptions(optsStr)

	return cv, nil
}

// -----------------------------------------------------------------------------
// 5) GetVariant - get a single variant by ID
// -----------------------------------------------------------------------------
func (r CatalogVariantRepository) GetVariant(ctx context.Context, variantID string) (*domain.CatalogVariant, error) {
	// First, get the catalog variant
	catalogVariant, err := r.FindVariant(ctx, variantID)
	if err != nil {
		return nil, err
	}
	if catalogVariant == nil {
		return nil, nil
	}

	// Convert CatalogVariant to Variant
	return &domain.CatalogVariant{
		ID:           catalogVariant.ID,
		ProductID:    catalogVariant.ProductID,
		Status:       catalogVariant.Status,
		SKU:          catalogVariant.SKU,
		Barcode:      catalogVariant.Barcode,
		Condition:    catalogVariant.Condition,
		VariantPrice: catalogVariant.VariantPrice,
		CurrencyCode: catalogVariant.CurrencyCode,
		Stock:        catalogVariant.Stock,
		Weight:       catalogVariant.Weight,
		Height:       catalogVariant.Height,
		Width:        catalogVariant.Width,
		Depth:        catalogVariant.Depth,
		Attributes:   catalogVariant.Attributes,
		IsAvailable:  catalogVariant.IsAvailable,
		HasOptions:   catalogVariant.HasOptions,
		Options:      catalogVariant.Options,
	}, nil
}

// 6) GetVariants
// -----------------------------------------------------------------------------
func (r CatalogVariantRepository) GetVariants(
	ctx context.Context,
	productID string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogVariant, int64, error) {

	offset := (page - 1) * pageSize
	if sortBy == "" {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE product_id = $1`, r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQ, productID).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting variants")
	}

	q := fmt.Sprintf(`
        SELECT
            id,
            product_id,
            status,
            sku,
            barcode,
            condition,
            variant_price,
            currency_code,
            stock,
            weight,
            height,
            width,
            depth,
            attributes,
            is_available,
            has_options,
            options
        FROM %s
        WHERE product_id = $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, q, productID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "retrieving variants")
	}
	defer rows.Close()

	var results []*domain.CatalogVariant
	for rows.Next() {
		cv := &domain.CatalogVariant{}
		var (
			statusStr string
			condStr   string
			attrStr   string
			optsStr   string
		)
		if err := rows.Scan(
			&cv.ID,
			&cv.ProductID,
			&statusStr,
			&cv.SKU,
			&cv.Barcode,
			&condStr,
			&cv.VariantPrice,
			&cv.CurrencyCode,
			&cv.Stock,
			&cv.Weight,
			&cv.Height,
			&cv.Width,
			&cv.Depth,
			&attrStr,
			&cv.IsAvailable,
			&cv.HasOptions,
			&optsStr,
		); err != nil {
			return nil, 0, errors.Wrap(err, "scanning variant row")
		}
		cv.Status = domain.ToProductStatus(statusStr)
		cv.Condition = domain.ToProductCondition(condStr)
		cv.Attributes = stringToAttributes(attrStr)
		cv.Options = stringToOptions(optsStr)

		results = append(results, cv)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "finishing variant rows")
	}

	return results, totalCount, nil
}

// -----------------------------------------------------------------------------
// 6) GetVariantsByProduct
// -----------------------------------------------------------------------------
func (r CatalogVariantRepository) GetVariantsByProduct(
	ctx context.Context,
	productID string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogVariant, int64, error) {

	offset := (page - 1) * pageSize
	if sortBy == "" {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE product_id = $1`, r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQ, productID).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting variants by product")
	}

	q := fmt.Sprintf(`
        SELECT
            id,
            product_id,
            status,
            sku,
            barcode,
            condition,
            variant_price,
            currency_code,
            stock,
            weight,
            height,
            width,
            depth,
            attributes,
            is_available,
            has_options,
            options
        FROM %s
        WHERE product_id = $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, q, productID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "retrieving variants by product")
	}
	defer rows.Close()

	var results []*domain.CatalogVariant
	for rows.Next() {
		cv := &domain.CatalogVariant{}
		var (
			statusStr string
			condStr   string
			attrStr   string
			optsStr   string
		)
		if err := rows.Scan(
			&cv.ID,
			&cv.ProductID,
			&statusStr,
			&cv.SKU,
			&cv.Barcode,
			&condStr,
			&cv.VariantPrice,
			&cv.CurrencyCode,
			&cv.Stock,
			&cv.Weight,
			&cv.Height,
			&cv.Width,
			&cv.Depth,
			&attrStr,
			&cv.IsAvailable,
			&cv.HasOptions,
			&optsStr,
		); err != nil {
			return nil, 0, errors.Wrap(err, "scanning variant by product row")
		}
		cv.Status = domain.ToProductStatus(statusStr)
		cv.Condition = domain.ToProductCondition(condStr)
		cv.Attributes = stringToAttributes(attrStr)
		cv.Options = stringToOptions(optsStr)

		results = append(results, cv)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "finishing variant rows by product")
	}
	return results, totalCount, nil
}

// -----------------------------------------------------------------------------
// 7) GetVariantsByCategory
// -----------------------------------------------------------------------------
func (r CatalogVariantRepository) GetVariantsByCategory(
	ctx context.Context,
	categoryID string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogVariant, int64, error) {
	// Implementation would require either storing category_id in the variants table
	// or doing a JOIN with products. We'll just return an error for now.
	return nil, 0, errors.Wrap(errors.ErrBadRequest,
		"Not Implemented: GetVariantsByCategory requires a join on the products table or storing category_id in variants")
}

// -----------------------------------------------------------------------------
// 8) RebrandVariant
// -----------------------------------------------------------------------------
func (r CatalogVariantRepository) RebrandVariant(
	ctx context.Context,
	variantID string,
	newName string,
	newDescription string,
) error {
	// If your table doesn't have columns `name` / `description`, remove them or add them.
	const query = `
        UPDATE %s
        SET name = $2,
            description = $3,
            updated_at = NOW()
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query), variantID, newName, newDescription)
	return err
}

// -----------------------------------------------------------------------------
// 9) AdjustVariantStock
// -----------------------------------------------------------------------------
func (r CatalogVariantRepository) AdjustVariantStock(
	ctx context.Context,
	variantID string,
	oldStock, newStock int64,
) error {
	const query = `
        UPDATE %s
        SET stock = $3,
            updated_at = NOW()
        WHERE id = $1
          AND stock = $2
    `
	res, err := r.db.ExecContext(ctx, r.table(query), variantID, oldStock, newStock)
	if err != nil {
		return errors.Wrap(err, "adjusting variant stock")
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.Wrap(err, "no variant updated; stock mismatch or not found")
	}
	return nil
}

// -----------------------------------------------------------------------------
// 10) ArchiveVariant
// -----------------------------------------------------------------------------
func (r CatalogVariantRepository) ArchiveVariant(ctx context.Context, variantID string) error {
	// Mark as unavailable or set a "status=archived" if you have that column
	const query = `
        UPDATE %s
        SET is_available = FALSE,
            updated_at = NOW()
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query), variantID)
	return err
}

// -----------------------------------------------------------------------------
// Helper to format table name
// -----------------------------------------------------------------------------
func (r CatalogVariantRepository) table(q string) string {
	return fmt.Sprintf(q, r.tableName)
}
