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

type ProductCacheRepository struct {
	tableName string
	db        postgres.DB
	fallback  domain.ProductRepository
}

var _ domain.ProductCacheRepository = (*ProductCacheRepository)(nil)

// NewProductCacheRepository constructor
func NewProductCacheRepository(tableName string, db postgres.DB, fallback domain.ProductRepository) ProductCacheRepository {
	return ProductCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

// -----------------------------------------------------------------------------
// Add - Similar to CatalogRepository.AddProduct, but for the 'products_cache' table
// -----------------------------------------------------------------------------
func (r ProductCacheRepository) Add(
	ctx context.Context,
	id, name, description string,
	basePrice int64,
	userSellerID, categoryID, brand string,
	condition domain.ProductCondition,
	model string,
	tags []string,
	manageStock bool,
	stock int64,
	sku string,
	attributes []domain.Attribute,
	weight, height, width, depth int64,
	status domain.ProductStatus,
	negotiable bool,
	userType domain.UserType,
	middlemanService bool,
	shippingCost int64,
	hasVariants bool,
	options []domain.Option,
	thumbnail string,
	lat, lng float64,
) error {
	// 1. Convert slices to strings so we can store in TEXT columns
	tagsStr := sliceToString(tags)
	attrsStr := attributesToString(attributes)
	optsStr := optionsToString(options)

	// 2. Insert (with location as ST_SetSRID(ST_MakePoint($lng, $lat),4326))
	const query = `
      INSERT INTO %s (
        id,
        name,
        description,
        base_price,
        user_seller_id,
        category_id,
        brand,
        condition,
        model,
        tags,
        attributes,
        manage_stock,
        stock,
        sku,
        weight,
        height,
        width,
        depth,
        status,
        negotiable,
        user_type,
        middleman_service,
        shipping_cost,
        has_variants,
        options,
        thumbnail,
        lat, lng
      )
      VALUES (
        $1, $2, $3,
        $4, $5, $6,
        $7, $8, $9,
        $10, $11,
        $12, $13, $14,
        $15, $16, $17, $18,
        $19, $20, $21, $22,
        $23, $24, $25, $26,$27,$28
        
      )
    `
	_, err := r.db.ExecContext(
		ctx,
		r.table(query),
		id,
		name,
		description,
		basePrice,
		userSellerID,
		categoryID,
		brand,
		condition,
		model,
		tagsStr,
		attrsStr,
		manageStock,
		stock,
		sku,
		weight,
		height,
		width,
		depth,
		status,
		negotiable,
		userType,
		middlemanService,
		shippingCost,
		hasVariants,
		optsStr,
		thumbnail,
		lng,
		lat,
	)
	if err != nil {
		// If there's a unique violation, maybe do nothing (already cached)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil
			}
		}
		return err
	}
	return nil
}

// -----------------------------------------------------------------------------
// Rebrand - Similar to CatalogRepository.RebrandProduct
// -----------------------------------------------------------------------------
func (r ProductCacheRepository) Rebrand(
	ctx context.Context,
	productID, name, description string,
	price, stock int64,
	sku, categoryID string,
) error {
	const query = `
      UPDATE %s
      SET
        name         = $2,
        description  = $3,
        base_price   = $4,
        stock        = $5,
        sku          = $6,
        category_id  = $7,
        updated_at   = NOW()
      WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		productID,
		name,
		description,
		price,
		stock,
		sku,
		categoryID,
	)
	return err
}

// -----------------------------------------------------------------------------
// UpdatePrice - like CatalogRepository.UpdatePrice
// -----------------------------------------------------------------------------
func (r ProductCacheRepository) UpdatePrice(
	ctx context.Context,
	productID string,
	oldPrice, newPrice int64,
) error {
	// If your domain logic requires matching oldPrice, do that:
	const query = `
      UPDATE %s
      SET base_price = $3,
          updated_at = NOW()
      WHERE id = $1
        AND base_price = $2
    `
	res, err := r.db.ExecContext(ctx, r.table(query), productID, oldPrice, newPrice)
	if err != nil {
		return errors.Wrap(err, "updating product price")
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.Wrap(err, "no product updated; possibly base_price mismatch or product not found")
	}
	return nil
}

// -----------------------------------------------------------------------------
// Remove - like CatalogRepository.RemoveProduct
// -----------------------------------------------------------------------------
func (r ProductCacheRepository) Remove(ctx context.Context, productID string) error {
	const query = `
      DELETE FROM %s
      WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query), productID)
	return err
}

// -----------------------------------------------------------------------------
// Find - checks cache, if missing => fallback repository
// -----------------------------------------------------------------------------
func (r ProductCacheRepository) Find(ctx context.Context, productID string) (*domain.Product, error) {
	const query = `
      SELECT
        name,
        description,
        base_price,
        user_seller_id,
        category_id,
        brand,
        condition,
        model,
        tags,
        attributes,
        manage_stock,
        stock,
        sku,
        weight,
        height,
        width,
        depth,
        status,
        negotiable,
        user_type,
        middleman_service,
        shipping_cost,
        has_variants,
        options,
        thumbnail,
      lat,
    lng
      FROM %s
      WHERE id = $1
      LIMIT 1
    `
	product := &domain.Product{ID: productID}
	var (
		tagsStr  string
		attrsStr string
		optsStr  string
		latVal   sql.NullFloat64
		lngVal   sql.NullFloat64
		mmsBool  bool
	)
	err := r.db.QueryRowContext(ctx, r.table(query), productID).Scan(
		&product.Name,
		&product.Description,
		&product.BasePrice,
		&product.UserSellerID,
		&product.CategoryID,
		&product.Brand,
		&product.Condition,
		&product.Model,
		&tagsStr,
		&attrsStr,
		&product.ManageStock,
		&product.Stock,
		&product.SKU,
		&product.Weight,
		&product.Height,
		&product.Width,
		&product.Depth,
		&product.Status,
		&product.Negotiable,
		&product.UserType,
		&mmsBool,
		&product.ShippingCost,
		&product.HasVariants,
		&optsStr,
		&product.Thumbnail,
		&latVal,
		&lngVal,
	)
	if err != nil {
		// If not found, fallback
		if errors.Is(err, sql.ErrNoRows) {
			fallbackProd, fbErr := r.fallback.Find(ctx, productID)
			if fbErr != nil {
				return nil, errors.Wrap(fbErr, "product fallback failed")
			}
			// Attempt to add it to the cache
			cacheErr := r.Add(
				ctx,
				fallbackProd.ID,
				fallbackProd.Name,
				fallbackProd.Description,
				fallbackProd.BasePrice,
				fallbackProd.UserSellerID,
				fallbackProd.CategoryID,
				fallbackProd.Brand,
				fallbackProd.Condition,
				fallbackProd.Model,
				fallbackProd.Tags,
				fallbackProd.ManageStock,
				fallbackProd.Stock,
				fallbackProd.SKU,
				fallbackProd.Attributes,
				fallbackProd.Weight,
				fallbackProd.Height,
				fallbackProd.Width,
				fallbackProd.Depth,
				fallbackProd.Status,
				fallbackProd.Negotiable,
				fallbackProd.UserType,
				fallbackProd.MiddlemanService,
				fallbackProd.ShippingCost,
				fallbackProd.HasVariants,
				fallbackProd.Options,
				fallbackProd.Thumbnail,
				fallbackProd.Lat,
				fallbackProd.Lng,
			)
			if cacheErr != nil {
				// Return the product anyway, even if caching failed
				return fallbackProd, errors.Wrap(cacheErr, "failed to cache fallback product")
			}
			return fallbackProd, nil
		}
		// Some other SQL error
		return nil, errors.Wrap(err, "scanning product from cache")
	}

	// Convert the stored strings to slices
	product.Tags = stringToSlice(tagsStr)
	product.Attributes = stringToAttributes(attrsStr)
	product.MiddlemanService = mmsBool
	product.Options = stringToOptions(optsStr)

	// If location is present, set lat/lng
	if latVal.Valid {
		product.Lat = latVal.Float64
	}
	if lngVal.Valid {
		product.Lng = lngVal.Float64
	}

	return product, nil
}

// table helps embed the table name in queries
func (r ProductCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

// -----------------------------------------------------------------------------
// Helper functions for converting slices <-> string
// (Same approach as in the provided CatalogRepository code.)
// -----------------------------------------------------------------------------

func sliceToString(sl []string) string {
	// e.g. convert ["red","blue"] => "[\"red\",\"blue\"]"
	return fmt.Sprintf("%q", sl)
}

func stringToSlice(s string) []string {
	// stub: parse string => []string
	return []string{}
}

func attributesToString(attrs []domain.Attribute) string {
	// Possibly store as JSON. For now, just do a naive representation.
	return fmt.Sprintf("%q", attrs)
}

func stringToAttributes(s string) []domain.Attribute {
	// stub: parse the string => []domain.Attribute
	return []domain.Attribute{}
}

func optionsToString(opts []domain.Option) string {
	return fmt.Sprintf("%q", opts)
}

func stringToOptions(s string) []domain.Option {
	return []domain.Option{}
}
