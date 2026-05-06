// File: search/internal/postgres/product_cache_repository.go

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"middleman/internal/di"
	"middleman/search/internal/application"
	"middleman/search/internal/constants"
	"middleman/search/internal/models"
	"middleman/search/internal/utils"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/stackus/errors"
)

type ProductCacheRepository struct {
	fallback       application.ProductRepository
	circuitBreaker *utils.CircuitBreaker
}

// Ensure we implement the ProductCacheRepository interface.
var _ application.ProductCacheRepository = (*ProductCacheRepository)(nil)

// NewProductCacheRepository constructs a RediSearch-based repo plus a fallback.
func NewProductCacheRepository(fallback application.ProductRepository) *ProductCacheRepository {
	return &ProductCacheRepository{
		fallback:       fallback,
		circuitBreaker: utils.NewCircuitBreaker(5, 30*time.Second), // Open after 5 failures, reset after 30s
	}
}

// -----------------------------------------------------------------------------
// Implementation of ProductCacheRepository interface
// -----------------------------------------------------------------------------

// Add indexes a new product in RediSearch. Optionally also add to fallback DB.
func (r *ProductCacheRepository) Add(
	ctx context.Context,
	productID string,
	name string,
	description string,
	basePrice int64,
	userSellerID string,
	categoryID string,
	categorySlug string,
	brand string,
	condition string,
	model string,
	tags []string,
	manageStock bool,
	stock int64,
	sku string,
	attributes []models.Attribute,
	weight int64,
	height int64,
	width int64,
	depth int64,
	status string,
	negotiable bool,
	userType string,
	middlemanService bool,
	shippingCost int64,
	hasVariants bool,
	options []models.Option,
	lat float64,
	lng float64,
	thumbnail string,
	entityType models.EntityType,
) error {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	tagsJSON, _ := json.Marshal(tags)
	optsJSON, _ := json.Marshal(options)
	attrsJSON, _ := json.Marshal(attributes)

	locationString := fmt.Sprintf("%.6f,%.6f", lng, lat)

	doc := redisearch.NewDocument(productID, 1.0).
		Set("product_id", productID).
		Set("name", name).
		Set("description", description).
		Set("base_price", basePrice).
		Set("user_seller_id", userSellerID).
		Set("category_id", categoryID).
		Set("category_slug", categorySlug).
		Set("brand", brand).
		Set("condition", condition).
		Set("model", model).
		Set("tags", string(tagsJSON)).
		Set("manage_stock", boolToInt(manageStock)).
		Set("stock", stock).
		Set("sku", sku).
		Set("attributes", string(attrsJSON)).
		Set("weight", weight).
		Set("height", height).
		Set("width", width).
		Set("depth", depth).
		Set("status", status).
		Set("negotiable", boolToInt(negotiable)).
		Set("user_type", userType).
		Set("middleman_service", boolToInt(middlemanService)).
		Set("shipping_cost", shippingCost).
		Set("has_variants", boolToInt(hasVariants)).
		Set("options", string(optsJSON)).
		Set("entity_type", entityType.String()).
		Set("thumbnail", thumbnail).
		Set("location", locationString)

	// Use replace option to prevent "Document already exists" errors
	if err := client.IndexOptions(redisearch.IndexingOptions{Replace: true}, doc); err != nil {
		return errors.Wrapf(err, "indexing product %s in RediSearch", productID)
	}
	return nil
}
func (r *ProductCacheRepository) GetCatalog(ctx context.Context, userID string) ([]*models.Product, error) {
	// Call SearchDealsWithFilter with term as name and defaults for other filters/pagination
	return r.fallback.GetCatalog(ctx, userID)
}

// Rebrand is an example partial update. We'll call Update for re-indexing.
func (r *ProductCacheRepository) Rebrand(
	ctx context.Context,
	productID string,
	name string,
	description string,
	basePrice int64,
	stock int64,
	sku string,
	categoryID string,
	status string,
) error {
	// Call Update with all parameters, using empty values for fields we don't want to update
	return r.Update(ctx,
		productID, name, description,
		basePrice,
		categoryID, "", "", "", "",
		nil,   // tags
		false, // manageStock
		stock,
		sku,
		nil,        // attributes ([]models.Attribute)
		0, 0, 0, 0, // weight, height, width, depth
		status,
		false, false, // negotiable, middlemanService
		"",    // userType
		0,     // shippingCost
		false, // hasVariants
		nil,   // option ([]models.Option)
		"",    // thumbnail
		0, 0)  // lat, lng
}

// Update properly handles all product fields like other cache repositories
func (r *ProductCacheRepository) Update(ctx context.Context,
	productID, name, description string,
	basePrice int64,
	categoryID, categorySlug, brand, condition, model string,
	tags []string,
	manageStock bool,
	stock int64,
	sku string,
	attributes []models.Attribute,
	weight, height, width, depth int64,
	status string,
	negotiable, middlemanService bool,
	userType string,
	shippingCost int64,
	hasVariants bool,
	option []models.Option,
	thumbnail string, lat, lng float64) error {

	// CRITICAL FIX: Validate productID before proceeding
	if productID == "" {
		log.Printf("⚠️ Update: productID is empty - cannot update product")
		return errors.ErrInvalidArgument.Msg("productID cannot be empty")
	}

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// 1. Retrieve existing product from fallback to preserve fields not in update
	existingProduct, err := r.fallback.Find(ctx, productID)
	if err != nil {
		return errors.Wrap(err, "finding product in fallback for update")
	}
	if existingProduct == nil {
		return errors.Wrapf(errors.ErrNotFound, "no fallback product found for ID=%s", productID)
	}

	// 2. Apply updates to existing product data, preserving non-empty values
	if name != "" {
		existingProduct.Name = name
	}
	if description != "" {
		existingProduct.Description = description
	}
	if basePrice > 0 {
		existingProduct.BasePrice = basePrice
	}
	if categoryID != "" {
		existingProduct.CategoryID = categoryID
	}
	if categorySlug != "" {
		existingProduct.CategorySlug = categorySlug
	}
	if brand != "" {
		existingProduct.Brand = brand
	}
	if condition != "" {
		existingProduct.Condition = condition
	}
	if model != "" {
		existingProduct.Model = model
	}
	if len(tags) > 0 {
		existingProduct.Tags = tags
	}
	// Always update boolean and numeric fields
	existingProduct.ManageStock = manageStock
	if stock >= 0 {
		existingProduct.Stock = stock
	}
	if sku != "" {
		existingProduct.SKU = sku
	}
	if len(attributes) > 0 {
		// Use the provided attributes directly since they're already []models.Attribute
		existingProduct.Attributes = attributes
	}
	if weight > 0 {
		existingProduct.Weight = weight
	}
	if height > 0 {
		existingProduct.Height = height
	}
	if width > 0 {
		existingProduct.Width = width
	}
	if depth > 0 {
		existingProduct.Depth = depth
	}
	if status != "" {
		existingProduct.Status = status
	}
	existingProduct.Negotiable = negotiable
	existingProduct.MiddlemanService = middlemanService
	if userType != "" {
		existingProduct.UserType = userType
	}
	if shippingCost >= 0 {
		existingProduct.ShippingCost = shippingCost
	}
	existingProduct.HasVariants = hasVariants
	if len(option) > 0 {
		existingProduct.Options = option
	}
	if thumbnail != "" {
		existingProduct.Thumbnail = thumbnail
	}
	if lat != 0 && lng != 0 {
		existingProduct.Lat = lat
		existingProduct.Lng = lng
	}

	// 3. CRITICAL FIX: Ensure EntityType is always set to ProductType
	if existingProduct.EntityType == "" || existingProduct.EntityType.String() == "" {
		existingProduct.EntityType = models.ProductType
		log.Printf("🔧 [Update] Fixed EntityType for product %s - set to ProductType", productID)
	}

	// 4. Delete old document and re-index with updated data
	client.DeleteDocument(productID)

	// 5. Re-index the updated product
	log.Printf("✅ [Update] Successfully updating product %s with all fields (EntityType: %s)", productID, existingProduct.EntityType.String())
	return r.addOrUpdateDoc(ctx, client, existingProduct)
}

// UpdateThumbnail updates only the thumbnail in Redis (and optionally fallback).
func (r *ProductCacheRepository) UpdateThumbnail(ctx context.Context, productID string, thumbnail string) error {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	client.DeleteDocument(productID)

	// Retrieve the product from fallback
	updated, err := r.fallback.Find(ctx, productID)
	if err != nil {
		return errors.Wrap(err, "finding product in fallback for thumbnail update")
	}
	if updated == nil {
		return errors.Wrapf(errors.ErrNotFound, "no fallback product found for ID=%s", productID)
	}

	// If your fallback DB also needs the updated thumbnail, do so here...

	// Temporarily override the thumbnail for the reindex
	updated.Thumbnail = thumbnail

	// CRITICAL FIX: Ensure EntityType is always set to ProductType
	if updated.EntityType == "" || updated.EntityType.String() == "" {
		updated.EntityType = models.ProductType
		log.Printf("🔧 [UpdateThumbnail] Fixed EntityType for product %s - set to ProductType", productID)
	}

	// Re-index with the updated thumbnail
	return r.addOrUpdateDoc(ctx, client, updated)
}

// Remove deletes from fallback plus from RediSearch.
func (r *ProductCacheRepository) Remove(ctx context.Context, productID string) error {

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)
	if err := client.DeleteDocument(productID); err != nil {
		return errors.Wrapf(err, "removing product %s from RediSearch", productID)
	}
	return nil
}

// Get is a trivial pass-through to fallback for a single product by ID.
func (r *ProductCacheRepository) Get(ctx context.Context, productID string) (*models.Product, error) {
	return r.fallback.Find(ctx, productID)
}

// Find tries RediSearch. If doc is missing => fallback => reindex if found.
func (r *ProductCacheRepository) Find(ctx context.Context, productID string) (*models.Product, error) {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// CRITICAL FIX: Escape product ID for TAG field to handle special characters like hyphens
	escapedProductID := redisearch.EscapeTextFileString(productID)
	q := redisearch.NewQuery(fmt.Sprintf("@entity_type:{%s} @product_id:{%s}", models.ProductType.String(), escapedProductID)).
		SetReturnFields(
			"name", "description", "base_price", "user_seller_id", "category_id", "category_slug",
			"brand", "condition", "model", "tags", "manage_stock", "stock", "sku",
			"attributes", "weight", "height", "width", "depth", "status", "negotiable",
			"user_type", "middleman_service", "shipping_cost", "has_variants", "options",
			"entity_type", "thumbnail",
			"location", "created_at", "updated_at",
		).
		Limit(0, 1)

	var docs []redisearch.Document
	var searchErr error
	
	// Use circuit breaker for Redis operations
	err := r.circuitBreaker.Call(ctx, func() error {
		docs, _, searchErr = client.Search(q)
		return searchErr
	})
	
	if err != nil {
		if errors.Is(err, errors.ErrUnavailable) {
			// Circuit breaker is open, go directly to fallback
			log.Printf("Find: Circuit breaker open, using fallback for ID=%s", productID)
			return r.fetchFromFallbackAndMaybeReindex(ctx, client, productID)
		}
		log.Printf("Find: RediSearch query error for ID=%s: %v", productID, err)
		return nil, errors.Wrap(err, "redisearch search error")
	}

	if len(docs) == 0 {
		log.Printf("Find: doc not in RediSearch => fallback for ID=%s", productID)
		return r.fetchFromFallbackAndMaybeReindex(ctx, client, productID)
	}

	doc := docs[0]
	p, parseErr := r.parseDocToProduct(doc)
	if parseErr != nil {
		return nil, parseErr
	}

	// If entityType != ProductType => fallback logic (optional)
	if p.EntityType != models.ProductType {
		log.Printf("[Find] entityType mismatch => fallback for productID=%s (type=%s)",
			productID, p.EntityType)
		return r.fallbackForWrongType(ctx, client, doc.Id)
	}
	return p, nil
}

// SuggestProducts uses a prefix search on the "name" field in RediSearch.
func (r *ProductCacheRepository) SuggestProducts(ctx context.Context, searchName string) ([]*models.Product, error) {
	if len(searchName) == 0 {
		return []*models.Product{}, nil
	}
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Use QueryBuilder for consistent query construction
	qb := NewQueryBuilder(models.ProductType)

	// Add prefix search for auto-completion
	escapedName := redisearch.EscapeTextFileString(searchName)
	qb.WithCustomFilter(fmt.Sprintf("@name:%s*", escapedName))

	// Limit to 10 results for suggestions
	qb.WithPagination(0, 10)

	// Set return fields
	qb.WithReturnFields(
		"product_id", "name", "description", "base_price", "user_seller_id",
		"stock", "sku", "attributes", "category_id", "category_slug", "brand", "condition",
		"model", "tags", "manage_stock", "weight", "height", "width", "depth",
		"status", "negotiable", "user_type", "middleman_service",
		"shipping_cost", "has_variants", "options", "entity_type", "thumbnail",
		"location", "created_at", "updated_at",
	)

	// Build the final query
	_, query := qb.Build()

	// Execute the search
	docs, _, err := client.Search(query)
	if err != nil {
		return nil, errors.Wrap(err, "prefix search error in RediSearch")
	}

	// Parse results
	products := make([]*models.Product, 0, len(docs))
	for _, doc := range docs {
		product, err := r.parseDocToProduct(doc)
		if err != nil {
			log.Printf("[WARNING] SuggestProducts: skipping product ID=%s due to parse error: %v",
				doc.Id, err)
			continue
		}
		products = append(products, product)
	}

	return products, nil
}

// SearchWithTerm => minimal pass-through to SearchWithFilters with minimal filters.
func (r *ProductCacheRepository) SearchWithTerm(ctx context.Context, name string) ([]*models.Product, error) {
	return r.SearchWithFilters(
		ctx,
		name, // name
		"",   // categoryID
		"",   // categorySlug
		0,    // minPrice
		0,    // maxPrice
		"",   // brand
		"",   // condition
		"",   // model
		nil,  // tags
		false,
		0,
		0,
		"",
		"",
		false,
		"",
		false,
		false,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		"",
		"",
	)
}

// SearchWithFilters tries RediSearch. If no docs => fallback => reindex.
func (r *ProductCacheRepository) SearchWithFilters(
	ctx context.Context,
	name string,
	categoryID string,
	categorySlug string,
	minPrice int64,
	maxPrice int64,
	brand string,
	condition string,
	model string,
	tags []string,
	manageStock bool,
	minStock int64,
	maxStock int64,
	sku string,
	status string,
	negotiable bool,
	userType string,
	middlemanService bool,
	hasVariants bool,
	shippingCost int64,
	minWeight int64,
	maxWeight int64,
	minHeight int64,
	maxHeight int64,
	minWidth int64,
	maxWidth int64,
	minDepth int64,
	maxDepth int64,
	offset int64,
	limit int64,
	lat float64,
	lng float64,
	radius int64,
	page int64,
	pageSize int64,
	sortBy string,
	sortOrder string,
) ([]*models.Product, error) {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Build the search query with entity type filter
	builder := NewQueryBuilder(models.ProductType)

	// Add name search if provided
	if name != "" {
		builder.WithNameFilter(name)
	}

	// Category ID filter
	if categoryID != "" {
		builder.WithCustomFilter(fmt.Sprintf("@category_id:{%s}", redisearch.EscapeTextFileString(categoryID)))
	}

	// Category slug filter
	if categorySlug != "" {
		builder.WithCustomFilter(fmt.Sprintf("@category_slug:{%s}", redisearch.EscapeTextFileString(categorySlug)))
	}

	// Price range filter
	builder.WithPriceRange(minPrice, maxPrice)

	// Brand filter
	if brand != "" {
		builder.WithCustomFilter(fmt.Sprintf("@brand:{%s}", redisearch.EscapeTextFileString(brand)))
	}

	// Condition filter
	if condition != "" {
		builder.WithCustomFilter(fmt.Sprintf("@condition:{%s}", redisearch.EscapeTextFileString(condition)))
	}

	// Model filter
	if model != "" {
		builder.WithCustomFilter(fmt.Sprintf("@model:{%s}", redisearch.EscapeTextFileString(model)))
	}

	// Tags filter
	if len(tags) > 0 {
		tagParts := make([]string, len(tags))
		for i, t := range tags {
			tagParts[i] = fmt.Sprintf("\"%s\"", redisearch.EscapeTextFileString(t))
		}
		builder.WithCustomFilter(fmt.Sprintf("@tags:{%s}", strings.Join(tagParts, "|")))
	}

	// Stock management filters
	if manageStock {
		builder.WithCustomFilter("@manage_stock:[1 1]")
	}

	// Only add stock filter for meaningful constraints
	if minStock > 0 || (maxStock > 0 && maxStock < 9999999) {
		low, high := "-inf", "+inf"
		if minStock > 0 {
			low = fmt.Sprintf("%d", minStock)
		}
		if maxStock > 0 && maxStock < 9999999 {
			high = fmt.Sprintf("%d", maxStock)
		}
		builder.WithCustomFilter(fmt.Sprintf("@stock:[%s %s]", low, high))
	}

	// SKU filter
	if sku != "" {
		builder.WithCustomFilter(fmt.Sprintf("@sku:{%s}", redisearch.EscapeTextFileString(sku)))
	}

	// Status filter
	if status != "" {
		builder.WithStatus(status)
	}

	// Flag filters
	if negotiable {
		builder.WithCustomFilter("@negotiable:[1 1]")
	}

	if userType != "" {
		builder.WithCustomFilter(fmt.Sprintf("@user_type:{%s}", redisearch.EscapeTextFileString(userType)))
	}

	if middlemanService {
		builder.WithCustomFilter("@middleman_service:[1 1]")
	}

	if hasVariants {
		builder.WithCustomFilter("@has_variants:[1 1]")
	}

	// Dimension filters - only add for meaningful constraints
	if minWeight > 0 || (maxWeight > 0 && maxWeight < 999999) {
		low, high := "-inf", "+inf"
		if minWeight > 0 {
			low = fmt.Sprintf("%d", minWeight)
		}
		if maxWeight > 0 && maxWeight < 999999 {
			high = fmt.Sprintf("%d", maxWeight)
		}
		builder.WithCustomFilter(fmt.Sprintf("@weight:[%s %s]", low, high))
	}

	if minHeight > 0 || (maxHeight > 0 && maxHeight < 9999999) {
		low, high := "-inf", "+inf"
		if minHeight > 0 {
			low = fmt.Sprintf("%d", minHeight)
		}
		if maxHeight > 0 && maxHeight < 9999999 {
			high = fmt.Sprintf("%d", maxHeight)
		}
		builder.WithCustomFilter(fmt.Sprintf("@height:[%s %s]", low, high))
	}

	if minWidth > 0 || (maxWidth > 0 && maxWidth < 999999) {
		low, high := "-inf", "+inf"
		if minWidth > 0 {
			low = fmt.Sprintf("%d", minWidth)
		}
		if maxWidth > 0 && maxWidth < 999999 {
			high = fmt.Sprintf("%d", maxWidth)
		}
		builder.WithCustomFilter(fmt.Sprintf("@width:[%s %s]", low, high))
	}

	if minDepth > 0 || (maxDepth > 0 && maxDepth < 9999999) {
		low, high := "-inf", "+inf"
		if minDepth > 0 {
			low = fmt.Sprintf("%d", minDepth)
		}
		if maxDepth > 0 && maxDepth < 9999999 {
			high = fmt.Sprintf("%d", maxDepth)
		}
		builder.WithCustomFilter(fmt.Sprintf("@depth:[%s %s]", low, high))
	}

	// Geo filter
	if lat != 0 && lng != 0 && radius > 0 {
		builder.WithGeoFilter(lat, lng, radius)
	}

	// Pagination setup
	finalOffset := offset
	finalLimit := limit
	if pageSize > 0 {
		if page < 1 {
			page = 1
		}
		finalOffset = (page - 1) * pageSize
		finalLimit = pageSize
	}
	if finalLimit <= 0 {
		finalLimit = 50
	}

	// Set pagination
	builder.WithPagination(int(finalOffset), int(finalLimit))

	// Set sorting if requested
	if sortBy != "" {
		sortDesc := strings.ToLower(sortOrder) == "desc"
		builder.WithSorting(sortBy, sortDesc)
	}

	// Set return fields
	builder.WithReturnFields(
		"product_id", "name", "description", "base_price", "user_seller_id",
		"category_id", "category_slug", "brand", "condition", "model", "tags", "manage_stock",
		"stock", "sku", "attributes", "weight", "height", "width", "depth",
		"status", "negotiable", "user_type", "middleman_service",
		"shipping_cost", "has_variants", "options", "entity_type", "thumbnail",
		"location",
	)

	// Build the final query
	_, query := builder.Build()

	// Execute the search
	docs, total, err := client.Search(query)

	if err != nil {
		return nil, err
	}

	// If no docs found => fallback => reindex
	if len(docs) == 0 {
		log.Printf("[SearchWithFilters] No docs => fallback. name=%q lat=%.6f lng=%.6f radius=%d",
			name, lat, lng, radius)

		fallbackProds, fallbackErr := r.fallback.SearchWithFilters(
			ctx,
			name,
			categoryID,
			categorySlug,
			minPrice,
			maxPrice,
			brand,
			condition,
			model,
			tags,
			manageStock,
			minStock,
			maxStock,
			sku,
			status,
			negotiable,
			userType,
			middlemanService,
			hasVariants,
			shippingCost,
			minWeight,
			maxWeight,
			minHeight,
			maxHeight,
			minWidth,
			maxWidth,
			minDepth,
			maxDepth,
			offset,
			limit,
			lat,
			lng,
			radius,
			page,
			pageSize,
			sortBy,
			sortOrder,
		)
		if fallbackErr != nil {
			return nil, errors.Wrap(fallbackErr, "fallback DB error in SearchWithFilters")
		}

		// Reindex asynchronously with proper context and rate limiting
		if len(fallbackProds) > 0 && len(fallbackProds) <= 100 { // Only reindex if reasonable number
			panicHandler := utils.NewPanicHandler(func(ctx context.Context, format string, args ...interface{}) {
				log.Printf(format, args...)
			})
			panicHandler.SafeGo(ctx, "product reindexing", func() {
				reindexCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				// Rate limit reindexing to prevent overwhelming Redis
				ticker := time.NewTicker(10 * time.Millisecond) // 100 ops/sec max
				defer ticker.Stop()

				for _, p := range fallbackProds {
					select {
					case <-reindexCtx.Done():
						return
					case <-ticker.C:
						if err := r.addOrUpdateDoc(reindexCtx, client, p); err != nil {
							log.Printf("[WARNING] Failed to reindex product %s: %v", p.ProductID, err)
						}
					}
				}
			})
		}

		return fallbackProds, nil
	}

	products := make([]*models.Product, 0, len(docs))
	for _, doc := range docs {
		product, err := r.parseDocToProduct(doc)
		if err != nil {
			log.Printf("[WARNING] SearchWithFilters: skipping product ID=%s due to parse error: %v",
				doc.Id, err)
			continue
		}
		products = append(products, product)
	}

	log.Printf("[SearchWithFilters] Found %d/%d products after parsing", len(products), total)
	return products, nil
}

func (r *ProductCacheRepository) SearchProductsWithCategorySlug(
	ctx context.Context,
	categorySlug string,
	page int64,
	pageSize int64,
	sortBy string,
	sortOrder string,
) ([]*models.Product, error) {

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Use the QueryBuilder pattern for consistent query construction
	qb := NewQueryBuilder(models.ProductType)

	// Add category slug filter if provided
	if categorySlug != "" {
		qb.WithCustomFilter(fmt.Sprintf("@category_slug:{%s}", redisearch.EscapeTextFileString(categorySlug)))
	}

	// Configure pagination
	finalOffset := int((page - 1) * pageSize)
	if finalOffset < 0 {
		finalOffset = 0
	}
	finalLimit := int(pageSize)
	if finalLimit <= 0 {
		finalLimit = 50
	}
	qb.WithPagination(finalOffset, finalLimit)

	// Set sorting if provided
	if sortBy != "" {
		sortDesc := strings.ToLower(sortOrder) == "desc"
		qb.WithSorting(sortBy, sortDesc)
	}

	// Set fields to return
	qb.WithReturnFields(
		"product_id", "name", "description", "base_price", "user_seller_id",
		"category_id", "category_slug", "brand", "condition", "model", "tags", "manage_stock",
		"stock", "sku", "attributes", "weight", "height", "width", "depth",
		"status", "negotiable", "user_type", "middleman_service",
		"shipping_cost", "has_variants", "options", "entity_type", "thumbnail",
		"location",
	)

	// Build the final query
	_, query := qb.Build()

	// Execute the search
	docs, total, err := client.Search(query)
	if err != nil {
		return nil, errors.Wrap(err, "RediSearch query error in SearchProductsWithCategorySlug")
	}

	// If nothing found in Redis => fallback => reindex
	if len(docs) == 0 {
		log.Printf("[SearchProductsWithCategorySlug] no docs => fallback (categorySlug=%q)", categorySlug)

		fallbackProds, fallbackErr := r.fallback.SearchProductsWithCategorySlug(
			ctx,
			categorySlug,
			page,
			pageSize,
			sortBy,
			sortOrder,
		)
		if fallbackErr != nil {
			return nil, errors.Wrap(fallbackErr, "fallback DB error in SearchProductsWithCategorySlug")
		}

		// Reindex asynchronously with timeout to prevent goroutine leaks
		panicHandler := utils.NewPanicHandler(func(ctx context.Context, format string, args ...interface{}) {
			log.Printf(format, args...)
		})
		panicHandler.SafeGo(ctx, "product category reindexing", func() {
			reindexCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			for _, p := range fallbackProds {
				select {
				case <-reindexCtx.Done():
					return
				default:
					if err := r.addOrUpdateDoc(reindexCtx, client, p); err != nil {
						continue
					}
				}
			}
		})

		return fallbackProds, nil
	}

	// Otherwise parse all the found docs
	products := make([]*models.Product, 0, len(docs))
	for _, doc := range docs {
		product, err := r.parseDocToProduct(doc)
		if err != nil {
			log.Printf("[WARNING] SearchProductsWithCategorySlug: skipping product ID=%s due to parse error: %v",
				doc.Id, err)
			continue
		}
		products = append(products, product)
	}

	log.Printf("[SearchProductsWithCategorySlug] Found %d/%d products after parsing", len(products), total)
	return products, nil
}

func (r *ProductCacheRepository) SearchProductsWithCategory(
	ctx context.Context,
	categoryId string,
	page int64,
	pageSize int64,
	sortBy string,
	sortOrder string,
) ([]*models.Product, error) {

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Use the QueryBuilder pattern for consistent query construction
	qb := NewQueryBuilder(models.ProductType)

	// Add category ID filter if provided
	if categoryId != "" {
		qb.WithCustomFilter(fmt.Sprintf("@category_id:{%s}", redisearch.EscapeTextFileString(categoryId)))
	}

	// Configure pagination
	finalOffset := int((page - 1) * pageSize)
	if finalOffset < 0 {
		finalOffset = 0
	}
	finalLimit := int(pageSize)
	if finalLimit <= 0 {
		finalLimit = 50
	}
	qb.WithPagination(finalOffset, finalLimit)

	// Set sorting if provided
	if sortBy != "" {
		sortDesc := strings.ToLower(sortOrder) == "desc"
		qb.WithSorting(sortBy, sortDesc)
	}

	// Set fields to return
	qb.WithReturnFields(
		"product_id", "name", "description", "base_price", "user_seller_id",
		"category_id", "category_slug", "brand", "condition", "model", "tags", "manage_stock",
		"stock", "sku", "attributes", "weight", "height", "width", "depth",
		"status", "negotiable", "user_type", "middleman_service",
		"shipping_cost", "has_variants", "options", "entity_type", "thumbnail",
		"location",
	)

	// Build the final query
	_, query := qb.Build()

	// Execute the search
	docs, total, err := client.Search(query)
	if err != nil {
		return nil, errors.Wrap(err, "RediSearch query error in SearchProductsWithCategory")
	}

	// If nothing found in Redis => fallback => reindex
	if len(docs) == 0 {
		log.Printf("[SearchProductsWithCategory] no docs => fallback (categoryId=%q)", categoryId)

		fallbackProds, fallbackErr := r.fallback.SearchProductsWithCategory(
			ctx,
			categoryId,
			page,
			pageSize,
			sortBy,
			sortOrder,
		)
		if fallbackErr != nil {
			return nil, errors.Wrap(fallbackErr, "fallback DB error in SearchProductsWithCategory")
		}

		// Reindex asynchronously with timeout to prevent goroutine leaks
		panicHandler := utils.NewPanicHandler(func(ctx context.Context, format string, args ...interface{}) {
			log.Printf(format, args...)
		})
		panicHandler.SafeGo(ctx, "product category reindexing", func() {
			reindexCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			for _, p := range fallbackProds {
				select {
				case <-reindexCtx.Done():
					return
				default:
					if err := r.addOrUpdateDoc(reindexCtx, client, p); err != nil {
						continue
					}
				}
			}
		})

		return fallbackProds, nil
	}

	// Otherwise parse all the found docs
	products := make([]*models.Product, 0, len(docs))
	for _, doc := range docs {
		product, err := r.parseDocToProduct(doc)
		if err != nil {
			log.Printf("[WARNING] SearchProductsWithCategory: skipping product ID=%s due to parse error: %v",
				doc.Id, err)
			continue
		}
		products = append(products, product)
	}

	log.Printf("[SearchProductsWithCategory] Found %d/%d products after parsing", len(products), total)
	return products, nil
}

// FindBatch retrieves multiple products by their IDs using Redis MGET for efficiency
func (r *ProductCacheRepository) FindBatch(ctx context.Context, productIDs []string) (map[string]*models.Product, error) {
	if len(productIDs) == 0 {
		return make(map[string]*models.Product), nil
	}

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)
	result := make(map[string]*models.Product, len(productIDs))
	
	// First, try to get from Redis using individual document fetches
	// RediSearch doesn't have a native batch get, so we'll fetch in parallel
	type fetchResult struct {
		id      string
		product *models.Product
		err     error
	}
	
	ch := make(chan fetchResult, len(productIDs))
	
	// Use worker pool to limit concurrent fetches
	const maxWorkers = 10
	sem := make(chan struct{}, maxWorkers)
	
	for _, id := range productIDs {
		sem <- struct{}{} // Acquire semaphore
		go func(productID string) {
			defer func() { <-sem }() // Release semaphore
			
			doc, err := client.Get(productID)
			if err != nil {
				ch <- fetchResult{id: productID, err: err}
				return
			}
			
			if doc != nil {
				product, err := r.parseDocToProduct(*doc)
				ch <- fetchResult{id: productID, product: product, err: err}
			} else {
				ch <- fetchResult{id: productID, err: fmt.Errorf("document not found")}
			}
		}(id)
	}
	
	// Collect results and track missing IDs
	var missingIDs []string
	for i := 0; i < len(productIDs); i++ {
		res := <-ch
		if res.err != nil {
			missingIDs = append(missingIDs, res.id)
		} else if res.product != nil {
			result[res.id] = res.product
		}
	}
	
	// If any IDs are missing, try fallback
	if len(missingIDs) > 0 && r.fallback != nil {
		log.Printf("[FindBatch] %d products not in cache, trying fallback", len(missingIDs))
		
		// Check if fallback supports batch fetch
		if batchFallback, ok := r.fallback.(application.ProductBatchRepository); ok {
			fallbackProducts, err := batchFallback.FindBatch(ctx, missingIDs)
			if err != nil {
				log.Printf("[FindBatch] Fallback batch fetch failed: %v", err)
			} else {
				// Add fallback results to result map and reindex them
				for id, product := range fallbackProducts {
					result[id] = product
					// Reindex asynchronously
					go func(p *models.Product) {
						if err := r.addOrUpdateDoc(context.Background(), client, p); err != nil {
							log.Printf("[FindBatch] Failed to reindex product %s: %v", p.ProductID, err)
						}
					}(product)
				}
			}
		} else {
			// Fallback doesn't support batch, fetch individually
			for _, id := range missingIDs {
				product, err := r.fallback.Find(ctx, id)
				if err == nil && product != nil {
					result[id] = product
					// Reindex asynchronously
					go func(p *models.Product) {
						if err := r.addOrUpdateDoc(context.Background(), client, p); err != nil {
							log.Printf("[FindBatch] Failed to reindex product %s: %v", p.ProductID, err)
						}
					}(product)
				}
			}
		}
	}
	
	return result, nil
}

// -----------------------------------------------------------------------------
// Internal helpers
// -----------------------------------------------------------------------------

func (r *ProductCacheRepository) addOrUpdateDoc(ctx context.Context, client redisearch.Client, p *models.Product) error {

	tagsJSON, _ := json.Marshal(p.Tags)
	optsJSON, _ := json.Marshal(p.Options)
	attrsJSON, _ := json.Marshal(p.Attributes)
	locStr := fmt.Sprintf("%.6f,%.6f", p.Lng, p.Lat)

	// Handle timestamps - set current time if not provided
	now := time.Now()
	createdAt := p.CreatedAt
	updatedAt := p.UpdatedAt

	if createdAt.IsZero() {
		createdAt = now
		p.CreatedAt = createdAt // Update the model
	}
	if updatedAt.IsZero() {
		updatedAt = now
		p.UpdatedAt = updatedAt // Update the model
	}

	doc := redisearch.NewDocument(p.ProductID, 1.0).
		Set("product_id", p.ProductID).
		Set("name", safeString(p.Name)).
		Set("description", safeString(p.Description)).
		Set("base_price", p.BasePrice).
		Set("user_seller_id", safeString(p.UserSellerID)).
		Set("category_id", safeString(p.CategoryID)).
		Set("category_slug", safeString(p.CategorySlug)).
		Set("brand", safeString(p.Brand)).
		Set("condition", safeString(p.Condition)).
		Set("model", safeString(p.Model)).
		Set("tags", string(tagsJSON)).
		Set("manage_stock", boolToInt(p.ManageStock)).
		Set("stock", p.Stock).
		Set("sku", safeString(p.SKU)).
		Set("attributes", string(attrsJSON)).
		Set("weight", p.Weight).
		Set("height", p.Height).
		Set("width", p.Width).
		Set("depth", p.Depth).
		Set("status", safeString(p.Status)).
		Set("negotiable", boolToInt(p.Negotiable)).
		Set("user_type", safeString(p.UserType)).
		Set("middleman_service", boolToInt(p.MiddlemanService)).
		Set("shipping_cost", p.ShippingCost).
		Set("has_variants", boolToInt(p.HasVariants)).
		Set("options", string(optsJSON)).
		Set("entity_type", p.EntityType.String()).
		Set("thumbnail", p.Thumbnail).
		Set("location", locStr).
		Set("created_at", createdAt.Unix()).
		Set("updated_at", updatedAt.Unix())

	// Use replace option to prevent "Document already exists" errors
	return client.IndexOptions(redisearch.IndexingOptions{Replace: true}, doc)
}

func (r *ProductCacheRepository) parseDocToProduct(doc redisearch.Document) (*models.Product, error) {
	p := &models.Product{ProductID: doc.Id}

	p.Name = strVal(doc.Properties["name"])
	p.Description = strVal(doc.Properties["description"])
	baseVal, err := parseInt64(doc.Properties["base_price"], "base_price", doc.Id)
	if err != nil {
		return nil, err
	}
	p.BasePrice = baseVal
	p.UserSellerID = strVal(doc.Properties["user_seller_id"])
	p.CategoryID = strVal(doc.Properties["category_id"])
	p.CategorySlug = strVal(doc.Properties["category_slug"])
	p.Brand = strVal(doc.Properties["brand"])
	p.Condition = strVal(doc.Properties["condition"])
	p.Model = strVal(doc.Properties["model"])

	// tags => JSON array
	if rawT, ok := doc.Properties["tags"].(string); ok && rawT != "" {
		var t []string
		if e := json.Unmarshal([]byte(rawT), &t); e == nil {
			p.Tags = t
		}
	}

	// manage_stock => bool
	mVal, _ := parseInt64(doc.Properties["manage_stock"], "manage_stock", doc.Id)
	p.ManageStock = (mVal == 1)

	// stock => int64
	stVal, stErr := parseInt64(doc.Properties["stock"], "stock", doc.Id)
	if stErr != nil {
		return nil, stErr
	}
	p.Stock = stVal
	p.SKU = strVal(doc.Properties["sku"])

	// attributes => JSON
	if rawA, ok := doc.Properties["attributes"].(string); ok && rawA != "" {
		var a []models.Attribute
		if e := json.Unmarshal([]byte(rawA), &a); e == nil {
			p.Attributes = a
		}
	}

	// weight, height, width, depth => int64
	p.Weight, _ = parseInt64(doc.Properties["weight"], "weight", doc.Id)
	p.Height, _ = parseInt64(doc.Properties["height"], "height", doc.Id)
	p.Width, _ = parseInt64(doc.Properties["width"], "width", doc.Id)
	p.Depth, _ = parseInt64(doc.Properties["depth"], "depth", doc.Id)

	// status => string
	p.Status = strVal(doc.Properties["status"])

	// negotiable => bool
	negVal, _ := parseInt64(doc.Properties["negotiable"], "negotiable", doc.Id)
	p.Negotiable = (negVal == 1)

	// user_type => string
	p.UserType = strVal(doc.Properties["user_type"])

	// middleman_service => bool
	mmVal, _ := parseInt64(doc.Properties["middleman_service"], "middleman_service", doc.Id)
	p.MiddlemanService = (mmVal == 1)

	// shipping_cost => int64
	shipVal, _ := parseInt64(doc.Properties["shipping_cost"], "shipping_cost", doc.Id)
	p.ShippingCost = shipVal

	// has_variants => bool
	hvVal, _ := parseInt64(doc.Properties["has_variants"], "has_variants", doc.Id)
	p.HasVariants = (hvVal == 1)

	// options => JSON array
	if rawOpts, ok := doc.Properties["options"].(string); ok && rawOpts != "" {
		var opts []models.Option
		if e := json.Unmarshal([]byte(rawOpts), &opts); e == nil {
			p.Options = opts
		}
	}

	// entity_type => domain model
	p.EntityType = models.ToEntityType(strVal(doc.Properties["entity_type"]))

	// thumbnail => string
	p.Thumbnail = strVal(doc.Properties["thumbnail"])

	// location => "lon,lat"
	if rawLoc, ok := doc.Properties["location"].(string); ok && rawLoc != "" {
		parts := strings.Split(rawLoc, ",")
		if len(parts) == 2 {
			lonF, _ := strconv.ParseFloat(parts[0], 64)
			latF, _ := strconv.ParseFloat(parts[1], 64)
			p.Lng = lonF
			p.Lat = latF
		}
	}

	// timestamps => Unix timestamps
	if createdAtUnix, err := parseInt64(doc.Properties["created_at"], "created_at", doc.Id); err == nil && createdAtUnix > 0 {
		p.CreatedAt = time.Unix(createdAtUnix, 0)
	}
	if updatedAtUnix, err := parseInt64(doc.Properties["updated_at"], "updated_at", doc.Id); err == nil && updatedAtUnix > 0 {
		p.UpdatedAt = time.Unix(updatedAtUnix, 0)
	}

	return p, nil
}

func (r *ProductCacheRepository) fallbackForWrongType(
	ctx context.Context,
	client redisearch.Client,
	docID string,
) (*models.Product, error) {
	// First check if we can determine entity type from the document
	entityType := ""

	// Try to get the entity_type field from the document
	q := redisearch.NewQuery(fmt.Sprintf("@$id:%s", docID)).
		SetReturnFields("entity_type").
		Limit(0, 1)

	docs, _, err := client.Search(q)
	if err == nil && len(docs) > 0 {
		if typeVal, ok := docs[0].Properties["entity_type"].(string); ok && typeVal != "" {
			entityType = typeVal
			log.Printf("[fallbackForWrongType] docID=%s entityType=%s => fallback", docID, entityType)
		}
	}

	// Remove the doc from RediSearch first
	if err := client.DeleteDocument(docID); err != nil {
		log.Printf("[fallbackForWrongType] could not delete docID=%s: %v", docID, err)
	}

	// Use the appropriate repository based on entity type
	var result interface{}
	var fetchErr error

	switch entityType {

	case models.PostType.String():
		// Try post repository if available
		if postRepo, ok := di.Get(ctx, constants.PostsRepoKey).(application.PostRepository); ok {
			result, fetchErr = postRepo.Find(ctx, docID)
			if fetchErr == nil && result != nil {
				log.Printf("[fallbackForWrongType] Found post with ID=%s", docID)
				return nil, errors.Wrap(errors.ErrInvalidArgument, "document is a post, not a product")
			}
		}

	}

	// Default: try product repository as fallback
	log.Printf("[fallbackForWrongType] Find: retrieving product with ID=%s", docID)
	fbProd, err := r.fallback.Find(ctx, docID)
	if err != nil {
		return nil, errors.Wrap(err, "fallbackForWrongType => fallback Find error")
	}
	if fbProd == nil {
		// nothing to reindex
		return nil, nil
	}
	return fbProd, nil
}

func (r *ProductCacheRepository) fetchFromFallbackAndMaybeReindex(
	ctx context.Context,
	client redisearch.Client,
	productID string,
) (*models.Product, error) {
	fbProd, fbErr := r.fallback.Find(ctx, productID)
	if fbErr != nil {
		return nil, errors.Wrap(fbErr, "fallback find error in fetchFromFallbackAndMaybeReindex")
	}
	if fbProd == nil {
		return nil, nil
	}

	// CRITICAL FIX: Ensure EntityType is always set to ProductType before reindexing
	if fbProd.EntityType == "" || fbProd.EntityType.String() == "" {
		fbProd.EntityType = models.ProductType
		log.Printf("🔧 [fetchFromFallbackAndMaybeReindex] Fixed EntityType for product %s - set to ProductType", productID)
	}

	// else reindex it
	if err := r.addOrUpdateDoc(ctx, client, fbProd); err != nil {
		log.Printf("fetchFromFallbackAndMaybeReindex: reindex product ID=%s error: %v", productID, err)
	}
	return fbProd, nil
}

// -----------------------------------------------------------------------------
// Simple type-conversion helpers
// -----------------------------------------------------------------------------
// Utility functions are now in shared_utils.go
// ============================================================================
// DEDICATED EVENT-SPECIFIC METHODS (avoiding bottlenecks)
// ============================================================================

func (r *ProductCacheRepository) IncreasePrice(ctx context.Context, productID string, newPrice int64) error {
	// Validate productID before proceeding
	if productID == "" {
		log.Printf("⚠️ IncreasePrice: productID is empty - cannot update product")
		return errors.ErrInvalidArgument.Msg("productID cannot be empty")
	}

	// Use Update method to ensure all fields are preserved
	return r.Update(ctx, productID, "", "", newPrice, "", "", "", "", "", nil, false, 0, "", nil, 0, 0, 0, 0, "", false, false, "", 0, false, nil, "", 0, 0)
}

func (r *ProductCacheRepository) DecreasePrice(ctx context.Context, productID string, newPrice int64) error {
	// Validate productID before proceeding
	if productID == "" {
		log.Printf("⚠️ DecreasePrice: productID is empty - cannot update product")
		return errors.ErrInvalidArgument.Msg("productID cannot be empty")
	}

	// Use Update method to ensure all fields are preserved
	return r.Update(ctx, productID, "", "", newPrice, "", "", "", "", "", nil, false, 0, "", nil, 0, 0, 0, 0, "", false, false, "", 0, false, nil, "", 0, 0)
}

func (r *ProductCacheRepository) MarkAsLeased(ctx context.Context, productID string) error {
	// Validate productID before proceeding
	if productID == "" {
		log.Printf("⚠️ MarkAsLeased: productID is empty - cannot update product")
		return errors.ErrInvalidArgument.Msg("productID cannot be empty")
	}

	// Use Update method to ensure all fields are preserved
	return r.Update(ctx, productID, "", "", 0, "", "", "", "", "", nil, false, 0, "", nil, 0, 0, 0, 0, "leased", false, false, "", 0, false, nil, "", 0, 0)
}

func (r *ProductCacheRepository) MarkAsSold(ctx context.Context, productID string) error {
	// Validate productID before proceeding
	if productID == "" {
		log.Printf("⚠️ MarkAsSold: productID is empty - cannot update product")
		return errors.ErrInvalidArgument.Msg("productID cannot be empty")
	}

	// Use Update method to ensure all fields are preserved
	return r.Update(ctx, productID, "", "", 0, "", "", "", "", "", nil, false, 0, "", nil, 0, 0, 0, 0, "sold", false, false, "", 0, false, nil, "", 0, 0)
}

func (r *ProductCacheRepository) MarkAsPawned(ctx context.Context, productID string) error {
	// Validate productID before proceeding
	if productID == "" {
		log.Printf("⚠️ MarkAsPawned: productID is empty - cannot update product")
		return errors.ErrInvalidArgument.Msg("productID cannot be empty")
	}

	// Use Update method to ensure all fields are preserved
	return r.Update(ctx, productID, "", "", 0, "", "", "", "", "", nil, false, 0, "", nil, 0, 0, 0, 0, "pawned", false, false, "", 0, false, nil, "", 0, 0)
}

func (r *ProductCacheRepository) AdjustStock(ctx context.Context, productID string, newStock int64) error {
	// Validate productID before proceeding
	if productID == "" {
		log.Printf("⚠️ AdjustStock: productID is empty - cannot update product")
		return errors.ErrInvalidArgument.Msg("productID cannot be empty")
	}

	// Use Update method to ensure all fields are preserved
	return r.Update(ctx, productID, "", "", 0, "", "", "", "", "", nil, false, newStock, "", nil, 0, 0, 0, 0, "", false, false, "", 0, false, nil, "", 0, 0)
}

func (r *ProductCacheRepository) ToggleNegotiable(ctx context.Context, productID string, negotiable bool) error {
	// Validate productID before proceeding
	if productID == "" {
		log.Printf("⚠️ ToggleNegotiable: productID is empty - cannot update product")
		return errors.ErrInvalidArgument.Msg("productID cannot be empty")
	}

	// Use Update method to ensure all fields are preserved
	return r.Update(ctx, productID, "", "", 0, "", "", "", "", "", nil, false, 0, "", nil, 0, 0, 0, 0, "", negotiable, false, "", 0, false, nil, "", 0, 0)
}

func (r *ProductCacheRepository) ArchiveProduct(ctx context.Context, productID string) error {
	// Validate productID before proceeding
	if productID == "" {
		log.Printf("⚠️ ArchiveProduct: productID is empty - cannot update product")
		return errors.ErrInvalidArgument.Msg("productID cannot be empty")
	}

	// Use Update method to ensure all fields are preserved
	return r.Update(ctx, productID, "", "", 0, "", "", "", "", "", nil, false, 0, "", nil, 0, 0, 0, 0, "archived", false, false, "", 0, false, nil, "", 0, 0)
}
