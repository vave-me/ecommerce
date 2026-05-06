// File: search/internal/redis/variant_cache_repository.go
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/stackus/errors"

	"middleman/internal/di"
	"middleman/search/internal/application"
	"middleman/search/internal/constants"
	"middleman/search/internal/models"
)

// VariantCacheRepository implements application.VariantCacheRepository,
// using RediSearch as a cache/index and delegating to a fallback
// for actual DB persistence.
type VariantCacheRepository struct {
	fallback application.VariantRepository // The fallback repository for DB operations
}

// Compile-time check that VariantCacheRepository implements the interface
var _ application.VariantCacheRepository = (*VariantCacheRepository)(nil)

// NewVariantCacheRepository constructs a new repository with the fallback.
func NewVariantCacheRepository(fallback application.VariantRepository) *VariantCacheRepository {
	return &VariantCacheRepository{
		fallback: fallback,
	}
}

// createIndex verifies that the unified RediSearch index exists and is ready.
// The actual unified index is created by SearchSystem.initRedisearch().
func (r *VariantCacheRepository) createIndex(ctx context.Context) error {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	log.Printf("VariantCacheRepository: Verifying unified index exists and is ready")

	// Check if the unified index exists and get its info
	info, err := client.Info()
	if err != nil {
		log.Printf("VariantCacheRepository: Index verification failed: %v", err)
		return fmt.Errorf("unified index not available: %w", err)
	}

	log.Printf("VariantCacheRepository: Unified index verified successfully")
	log.Printf("VariantCacheRepository: Index contains %v documents", info.DocCount)

	return nil
}

// -----------------------------------------------------------------------------
// Add
// -----------------------------------------------------------------------------
func (r *VariantCacheRepository) Add(
	ctx context.Context,
	variantID, productID, name, sku, barcode string,
	variantPrice int64,
	currencyCode string,
	stock, weight, height, width, depth int64,
	attributes []models.Attribute,
	isAvailable bool,
	hasOptions bool,
	options []models.Option,
) error {
	// (Optional) A) Persist in fallback DB:
	// err := r.fallback.Add(ctx, variantID, productID, name, sku, barcode, ...)
	// if err != nil { return errors.Wrap(err, "adding variant in fallback DB") }

	// B) Index in RediSearch
	return r.indexVariantInRedis(ctx, &models.Variant{
		VariantID:    variantID,
		ProductID:    productID,
		Name:         name,
		SKU:          sku,
		Barcode:      barcode,
		VariantPrice: variantPrice,
		CurrencyCode: currencyCode,
		Stock:        stock,
		Weight:       weight,
		Height:       height,
		Width:        width,
		Depth:        depth,
		Attributes:   attributes,
		IsAvailable:  isAvailable,
		HasOptions:   hasOptions,
		Options:      options,
	})
}

// -----------------------------------------------------------------------------
// Find
// -----------------------------------------------------------------------------
func (r *VariantCacheRepository) Find(ctx context.Context, variantID string) (*models.Variant, error) {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// CRITICAL FIX: Escape variant ID for TAG field to handle special characters like hyphens
	escapedVariantID := redisearch.EscapeTextFileString(variantID)
	query := redisearch.NewQuery(fmt.Sprintf("@variant_id:{%s}", escapedVariantID)).
		SetReturnFields(
			"variant_id", "product_id", "name", "sku", "barcode", "variant_price",
			"currency_code", "stock", "weight", "height", "width", "depth",
			"attributes", "is_available", "has_options", "options", "entity_type",
		).
		Limit(0, 1)

	docs, _, err := client.Search(query)
	if err != nil {
		log.Printf("Find(variant): RediSearch error for ID=%s: %v", variantID, err)
		return nil, errors.Wrap(err, "searching variant in RediSearch")
	}

	if len(docs) == 0 {
		// fallback DB
		log.Printf("Find(variant): Not in Redis, fallback to DB for ID=%s", variantID)
		v, fe := r.fallback.Find(ctx, variantID)
		if fe != nil {
			return nil, fe
		}
		// Optionally cache for next time
		if v != nil {
			if e := r.indexVariantInRedis(ctx, v); e != nil {
				log.Printf("Warning: caching variant ID %s in RediSearch: %v", variantID, e)
			}
		}
		return v, nil
	}

	return r.parseDocToVariant(docs[0])
}

// -----------------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------------
func (r *VariantCacheRepository) Update(
	ctx context.Context,
	variantID string,
	newVariantPrice int64,
	newStock int64,
	newName string,
	newAttributes []models.Attribute,
) error {
	// A) fallback DB update
	err := r.fallback.Update(ctx, variantID, newVariantPrice, newStock, newName, newAttributes)
	if err != nil {
		return errors.Wrap(err, "fallback updating variant")
	}

	// B) re-index in Redis
	updated, findErr := r.fallback.Find(ctx, variantID)
	if findErr != nil {
		return errors.Wrap(findErr, "finding updated variant in fallback for re-index")
	}

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)
	// Remove old doc
	if delErr := client.DeleteDocument(variantID); delErr != nil {
		log.Printf("Update(variant): Could not delete doc ID=%s: %v", variantID, delErr)
	}

	return r.indexVariantInRedis(ctx, updated)
}

// -----------------------------------------------------------------------------
// Remove
// -----------------------------------------------------------------------------
func (r *VariantCacheRepository) Remove(ctx context.Context, variantID string) error {
	// A) fallback
	if err := r.fallback.Remove(ctx, variantID); err != nil {
		return errors.Wrap(err, "removing variant from fallback DB")
	}

	// B) remove from RediSearch
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)
	if err := client.DeleteDocument(variantID); err != nil {
		return errors.Wrap(err, "removing variant doc in RediSearch")
	}
	return nil
}

// -----------------------------------------------------------------------------
// Rebrand
// -----------------------------------------------------------------------------
func (r *VariantCacheRepository) Rebrand(
	ctx context.Context,
	variantID string,
	name string,
	variantPrice int64,
	stock int64,
	attributes []models.Attribute,
	isAvailable bool,
) error {
	// 1) fetch from fallback
	v, err := r.fallback.Find(ctx, variantID)
	if err != nil {
		return errors.Wrap(err, "finding variant for rebrand fallback")
	}

	// 2) modify fields
	v.Name = name
	v.VariantPrice = variantPrice
	v.Stock = stock
	v.Attributes = attributes
	v.IsAvailable = isAvailable

	// 3) update fallback
	updateErr := r.fallback.Update(ctx, variantID, variantPrice, stock, name, attributes)
	if updateErr != nil {
		return errors.Wrap(updateErr, "rebranding variant in fallback DB")
	}

	// 4) re-index in Redis
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)
	if delErr := client.DeleteDocument(variantID); delErr != nil {
		log.Printf("Rebrand(variant): Could not delete doc ID=%s: %v", variantID, delErr)
	}
	return r.indexVariantInRedis(ctx, v)
}

// -----------------------------------------------------------------------------
// SearchWithFilters
// -----------------------------------------------------------------------------
func (r *VariantCacheRepository) SearchWithFilters(
	ctx context.Context,
	name string,
	minPrice int64,
	maxPrice int64,
	offset int64,
	limit int64,
) ([]*models.Variant, error) {
	log.Println("VariantCacheRepository: SearchWithFilters started")

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)
	var qb strings.Builder

	// partial match on name
	if name != "" {
		escaped := redisearch.EscapeTextFileString(name)
		fmt.Fprintf(&qb, "@name:(%s*) ", escaped)
	}
	// numeric range for variant_price
	if minPrice > 0 || maxPrice > 0 {
		if minPrice < 0 || maxPrice < 0 {
			return nil, errors.ErrInvalidArgument.Msg("negative prices invalid for variant search")
		}
		if maxPrice == 0 {
			maxPrice = 999999999
		}
		fmt.Fprintf(&qb, "@variant_price:[%d %d] ", minPrice, maxPrice)
	}

	rawQuery := strings.TrimSpace(qb.String())
	if rawQuery == "" {
		rawQuery = "*"
	}

	query := redisearch.NewQuery(rawQuery).
		SetReturnFields(
			"variant_id", "product_id", "name", "sku", "barcode", "variant_price",
			"currency_code", "stock", "weight", "height", "width", "depth",
			"attributes", "is_available", "has_options", "options", "entity_type",
		).
		Limit(int(offset), int(limit))

	docs, _, err := client.Search(query)
	if err != nil {
		return nil, errors.Wrap(err, "RediSearch variant query error")
	}

	results := make([]*models.Variant, 0, len(docs))
	for _, doc := range docs {
		v, parseErr := r.parseDocToVariant(doc)
		if parseErr != nil {
			log.Printf("SearchWithFilters(variant): skipping doc ID=%s parse error: %v", doc.Id, parseErr)
			continue
		}
		results = append(results, v)
	}
	return results, nil
}

// -----------------------------------------------------------------------------
// SearchWithTerm
// -----------------------------------------------------------------------------
func (r *VariantCacheRepository) SearchWithTerm(ctx context.Context, name string) ([]*models.Variant, error) {
	// Reuse SearchWithFilters with an unbounded price range, offset=0, limit=100
	return r.SearchWithFilters(ctx, name, 0, 0, 0, 100)
}

// -----------------------------------------------------------------------------
// SuggestVariants
// -----------------------------------------------------------------------------
func (r *VariantCacheRepository) SuggestVariants(ctx context.Context, partialName string) ([]*models.Variant, error) {
	if partialName == "" {
		return []*models.Variant{}, nil
	}
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	escapedName := redisearch.EscapeTextFileString(partialName)
	query := redisearch.NewQuery(fmt.Sprintf("@name:%s*", escapedName)).
		SetReturnFields(
			"variant_id", "product_id", "name", "sku", "barcode", "variant_price",
			"currency_code", "stock", "weight", "height", "width", "depth",
			"attributes", "is_available", "has_options", "options", "entity_type",
		).
		Limit(0, 10)

	docs, _, err := client.Search(query)
	if err != nil {
		return nil, errors.Wrap(err, "SuggestVariants: RediSearch prefix search error")
	}

	suggestions := make([]*models.Variant, 0, len(docs))
	for _, doc := range docs {
		v, parseErr := r.parseDocToVariant(doc)
		if parseErr != nil {
			log.Printf("SuggestVariants: skipping docID=%s parse error: %v", doc.Id, parseErr)
			continue
		}
		suggestions = append(suggestions, v)
	}
	return suggestions, nil
}

// -----------------------------------------------------------------------------
// Internals
// -----------------------------------------------------------------------------

// indexVariantInRedis indexes a variant as a Redis doc with doc.Id = variantID
func (r *VariantCacheRepository) indexVariantInRedis(ctx context.Context, v *models.Variant) error {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Marshal []Attribute and []Option to JSON for storage
	attrJSON, attrErr := json.Marshal(v.Attributes)
	if attrErr != nil {
		log.Printf("indexVariantInRedis: Error marshaling attributes for %s: %v", v.VariantID, attrErr)
		attrJSON = []byte("[]")
	}
	optJSON, optErr := json.Marshal(v.Options)
	if optErr != nil {
		log.Printf("indexVariantInRedis: Error marshaling options for %s: %v", v.VariantID, optErr)
		optJSON = []byte("[]")
	}

	doc := redisearch.NewDocument(v.VariantID, 1.0).
		Set("variant_id", v.VariantID).
		Set("product_id", v.ProductID).
		Set("name", v.Name).
		Set("sku", v.SKU).
		Set("barcode", v.Barcode).
		Set("variant_price", v.VariantPrice).
		Set("currency_code", v.CurrencyCode).
		Set("stock", v.Stock).
		Set("weight", v.Weight).
		Set("height", v.Height).
		Set("width", v.Width).
		Set("depth", v.Depth).
		Set("attributes", string(attrJSON)).
		Set("is_available", stringVal(v.IsAvailable)).
		Set("has_options", stringVal(v.HasOptions)).
		Set("options", string(optJSON)).
		Set("entity_type", "variant") // CRITICAL FIX: Add entity_type for filtering

	return client.IndexOptions(redisearch.DefaultIndexingOptions, doc)
}

// parseDocToVariant constructs a models.Variant from a redisearch.Document
func (r *VariantCacheRepository) parseDocToVariant(doc redisearch.Document) (*models.Variant, error) {
	v := &models.Variant{VariantID: doc.Id}

	v.ProductID = stringVal(doc.Properties["product_id"])
	v.Name = stringVal(doc.Properties["name"])
	v.SKU = stringVal(doc.Properties["sku"])
	v.Barcode = stringVal(doc.Properties["barcode"])

	vp, err := parseInt64(doc.Properties["variant_price"], "variant_price", doc.Id)
	if err != nil {
		return nil, err
	}
	v.VariantPrice = vp

	v.CurrencyCode = stringVal(doc.Properties["currency_code"])

	st, err := parseInt64(doc.Properties["stock"], "stock", doc.Id)
	if err != nil {
		return nil, err
	}
	v.Stock = st

	w, _ := parseInt64(doc.Properties["weight"], "weight", doc.Id)
	v.Weight = w
	h, _ := parseInt64(doc.Properties["height"], "height", doc.Id)
	v.Height = h
	wd, _ := parseInt64(doc.Properties["width"], "width", doc.Id)
	v.Width = wd
	d, _ := parseInt64(doc.Properties["depth"], "depth", doc.Id)
	v.Depth = d

	// attributes => []models.Attribute stored as JSON
	if rawAttr, ok := doc.Properties["attributes"].(string); ok && rawAttr != "" {
		var at []models.Attribute
		if e := json.Unmarshal([]byte(rawAttr), &at); e == nil {
			v.Attributes = at
		}
	}

	ia, err := parseInt64(doc.Properties["is_available"], "is_available", doc.Id)
	if err == nil {
		v.IsAvailable = (ia == 1)
	}

	ho, err := parseInt64(doc.Properties["has_options"], "has_options", doc.Id)
	if err == nil {
		v.HasOptions = (ho == 1)
	}

	// options => []models.Option stored as JSON
	if rawOpts, ok := doc.Properties["options"].(string); ok && rawOpts != "" {
		var opts []models.Option
		if e := json.Unmarshal([]byte(rawOpts), &opts); e == nil {
			v.Options = opts
		}
	}

	return v, nil
}

// stringVal tries a string cast or returns ""
func stringVal(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
