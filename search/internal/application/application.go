package application

import (
	"context"
	"fmt"
	"log"
	"math"
	"middleman/search/internal/constants"
	"middleman/search/internal/models"
	"middleman/search/internal/utils"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type (
	GetOrder struct {
		OrderID string
	}
	GetUser struct {
		UserID string
	}
	GetCatalog struct {
		UserID      string
		EntityTypes []string
	}
	SearchUsers struct {
		Name string
		Lat  float64
		Lon  float64
	}

	// Products
	SuggestProducts struct {
		Name string
	}
	GetProduct struct {
		ProductID string
	}
	SearchProductsWithFilters struct {
		Name             string
		CategoryID       string
		CategorySlug     string
		MinPrice         int64
		MaxPrice         int64
		Brand            string
		Condition        string
		Model            string
		Tags             []string
		ManageStock      bool
		MinStock         int64
		MaxStock         int64
		SKU              string
		Status           string
		Negotiable       bool
		UserType         string
		MiddlemanService bool
		HasVariants      bool
		ShippingCost     int64
		MinWeight        int64
		MaxWeight        int64
		MinHeight        int64
		MaxHeight        int64
		MinWidth         int64
		MaxWidth         int64
		MinDepth         int64
		MaxDepth         int64
		Offset           int64
		Limit            int64
		Lat              float64
		Lng              float64
		Radius           int64
		Page             int64
		PageSize         int64
		SortBy           string
		SortOrder        string
		Type             models.EntityType
	}
	SearchProductsWithCategorySlug struct {
		Name         string
		CategorySlug string
		Offset       int64
		Limit        int64
		Lat          float64
		Lng          float64
		Radius       int64
		Page         int64
		PageSize     int64
		SortBy       string
		SortOrder    string
		Type         models.EntityType
	}
	SearchProductsWithCategory struct {
		Name       string
		CategoryID string
		Offset     int64
		Limit      int64
		Lat        float64
		Lng        float64
		Radius     int64
		Page       int64
		PageSize   int64
		SortBy     string
		SortOrder  string
		Type       models.EntityType
	}

	SearchProductsWithTerm struct {
		Name string
	}

	//POSTS
	GetPost struct {
		PostID string
	}
	SuggestPosts struct {
		Name string
	}
	SearchPostsWithFilters struct {
		Name         string
		Description  string
		Tags         []string
		PostType     string
		UserType     string
		CategoryID   string
		CategorySlug string
		Status       string
		Thumbnail    string
		Offset       int64
		Limit        int64
		Lat          float64
		Lng          float64
		Radius       int64
		Page         int64
		PageSize     int64
		SortBy       string
		SortOrder    string
		Type         models.EntityType
	}
	SearchPostsWithTerm struct {
		Name string
	}

	SearchPostsWithCategorySlug struct {
		Name         string
		CategorySlug string
		Offset       int64
		Limit        int64
		Lat          float64
		Lng          float64
		Radius       int64
		Page         int64
		PageSize     int64
		SortBy       string
		SortOrder    string
		EntityType   models.EntityType
	}
	SearchPostsWithCategory struct {
		Name       string
		CategoryID string
		Offset     int64
		Limit      int64
		Lat        float64
		Lng        float64
		Radius     int64
		Page       int64
		PageSize   int64
		SortBy     string
		SortOrder  string
		EntityType models.EntityType
	}

	GetDeal struct {
		DealID string
	}
	SuggestDeals struct {
		Name string
	}
	SearchDealsWithFilters struct {
		Name         string
		Category     string
		MinPrice     int64
		MaxPrice     int64
		Brand        string
		Condition    string
		Model        string
		Tags         []string
		ManageStock  bool
		MinStock     int64
		MaxStock     int64
		SKU          string
		Status       string
		Negotiable   bool
		UserType     string
		HasVariants  bool
		ShippingCost int64
		MinWeight    int64
		MaxWeight    int64
		MinHeight    int64
		MaxHeight    int64
		MinWidth     int64
		MaxWidth     int64
		MinDepth     int64
		MaxDepth     int64
		Offset       int64
		Limit        int64
		Lat          float64
		Lng          float64
		Radius       int64
		Page         int64
		PageSize     int64
		SortBy       string
		SortOrder    string
		EntityType   models.EntityType
	}
	SearchDealsWithTerm struct {
		Name string
	}

	SearchDealsWithCategorySlug struct {
		Name         string
		CategorySlug string
		Offset       int64
		Limit        int64
		Lat          float64
		Lng          float64
		Radius       int64
		Page         int64
		PageSize     int64
		SortBy       string
		SortOrder    string
		EntityType   models.EntityType
	}
	SearchDealsWithCategory struct {
		Name       string
		CategoryID string
		Offset     int64
		Limit      int64
		Lat        float64
		Lng        float64
		Radius     int64
		Page       int64
		PageSize   int64
		SortBy     string
		SortOrder  string
		EntityType models.EntityType
	}
	GetJob struct {
		JobID string
	}
	SuggestJobs struct {
		Name string
	}
	SearchJobsWithFilters struct {
		Name              string
		SearchTerm        string
		Description       string
		MinSalary         int64
		MaxSalary         int64
		CategoryID        string
		CategorySlug      string
		SeniorityLevel    string
		RelocationSupport bool
		EmploymentType    string
		UserID            string
		UserType          string
		Tags              []string
		Offset            int64
		Limit             int64
		Lat               float64
		Lng               float64
		Radius            int64
		Page              int64
		PageSize          int64
		SortBy            string
		SortOrder         string
		EntityType        models.EntityType
	}
	GetService struct {
		ServiceID string
	}
	SuggestServices struct {
		Name string
	}
	SearchServicesWithFilters struct {
		CategoryID       string
		CategorySlug     string
		ServiceType      string
		UserID           string
		Status           models.Status
		SearchText       string
		MinPrice         int64
		MaxPrice         int64
		AvailableFrom    time.Time
		AvailableTo      time.Time
		HasVariants      bool
		Negotiable       bool
		MiddlemanService bool
		UserType         models.UserType
		Tags             []string
		Qualifications   []string
		Offset           int64
		Limit            int64
		Lat              float64
		Lng              float64
		Radius           int64
		Page             int64
		PageSize         int64
		SortBy           string
		SortOrder        string
		EntityType       models.EntityType
	}
	SearchServicesWithTerm struct {
		Name string
	}
	SearchServicesWithCategorySlug struct {
		Name         string
		CategorySlug string
		Offset       int64
		Limit        int64
		Lat          float64
		Lng          float64
		Radius       int64
		Page         int64
		PageSize     int64
		SortBy       string
		SortOrder    string
		EntityType   models.EntityType
	}
	SearchServicesWithCategory struct {
		Name       string
		CategoryID string
		Offset     int64
		Limit      int64
		Lat        float64
		Lng        float64
		Radius     int64
		Page       int64
		PageSize   int64
		SortBy     string
		SortOrder  string
		Type       models.EntityType
	}
	// Unified search types
	UnifiedSearchParams struct {
		SearchTerm       string
		EntityTypes      []string
		Page             int64
		PageSize         int64
		Lat              float64
		Lng              float64
		Radius           int64
		SortBy           string
		SortOrder        string
		MinPrice         int64
		MaxPrice         int64
		CategoryID       string
		CategorySlug     string
		UserType         string
		Negotiable       bool
		Brand            string
		Condition        string
		Model            string
		Tags             []string
		HasVariants      bool
		MiddlemanService bool
		Status           string
		ServiceType      string
	}

	UnifiedSearchResult struct {
		EntityType     string
		Product        *models.Product
		Post           *models.Post
		Service        *models.Service
		RelevanceScore float64
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	UnifiedSearchResults struct {
		Results      []UnifiedSearchResult
		TotalCount   int64
		Page         int64
		PageSize     int64
		CountsByType map[string]int64
	}

	UnifiedFeedParams struct {
		EntityTypes []string
		FeedType    string // "latest", "trending", "recommended"
		Page        int64
		PageSize    int64
		Lat         float64
		Lng         float64
		Radius      int64
		UserID      string // For personalized feeds
	}

	UnifiedFeedResults struct {
		Items         []UnifiedSearchResult
		TotalCount    int64
		NextPageToken string
	}

	Application interface {
		//USERS
		GetUser(ctx context.Context, get GetUser) (*models.User, error)

		//PRODUCTS
		GetProduct(ctx context.Context, get GetProduct) (*models.Product, error)
		SuggestProducts(ctx context.Context, get SuggestProducts) ([]*models.Product, error)
		SearchProductsWithFilters(ctx context.Context, search SearchProductsWithFilters) ([]*models.Product, error)
		SearchProductsWithCategorySlug(ctx context.Context, search SearchProductsWithCategorySlug) ([]*models.Product, error)
		SearchProductsWithCategory(ctx context.Context, search SearchProductsWithCategory) ([]*models.Product, error)
		SearchProductsWithTerm(ctx context.Context, search SearchProductsWithTerm) ([]*models.Product, error)

		//POSTS
		SuggestPosts(ctx context.Context, get SuggestPosts) ([]*models.Post, error)
		GetPost(ctx context.Context, get GetPost) (*models.Post, error)
		SearchPostsWithCategorySlug(ctx context.Context, search SearchPostsWithCategorySlug) ([]*models.Post, error)
		SearchPostsWithCategory(ctx context.Context, search SearchPostsWithCategory) ([]*models.Post, error)
		SearchPostsWithFilters(ctx context.Context, search SearchPostsWithFilters) ([]*models.Post, error)
		SearchPostsWithTerm(ctx context.Context, search SearchPostsWithTerm) ([]*models.Post, error)

		//DEALS
		//SERVICES
		SuggestServices(ctx context.Context, get SuggestServices) ([]*models.Service, error)
		GetService(ctx context.Context, get GetService) (*models.Service, error)
		SearchServicesWithCategorySlug(ctx context.Context, search SearchServicesWithCategorySlug) ([]*models.Service, error)
		SearchServicesWithCategory(ctx context.Context, search SearchServicesWithCategory) ([]*models.Service, error)
		SearchServicesWithFilters(ctx context.Context, search SearchServicesWithFilters) ([]*models.Service, error)
		SearchServicesWithTerm(ctx context.Context, search SearchServicesWithTerm) ([]*models.Service, error)
		// ORDERS
		GetOrder(ctx context.Context, get GetOrder) (*models.Order, error)

		// UNIFIED SEARCH
		UnifiedSearch(ctx context.Context, params UnifiedSearchParams) (*UnifiedSearchResults, error)
		UnifiedFeed(ctx context.Context, params UnifiedFeedParams) (*UnifiedFeedResults, error)
		GetCatalog(ctx context.Context, params UnifiedFeedParams) (*UnifiedFeedResults, error)
	}

	// app struct implements the Application interface
	app struct {
		orders      OrderRepository
		products    ProductCacheRepository
		variants    VariantCacheRepository
		posts       PostCacheRepository
		users       UserCacheRepository
		services    ServiceCacheRepository
		itemMetrics MetricRepository
	}
)

var _ Application = (*app)(nil)

func New(
	orders OrderRepository,
	products ProductCacheRepository,
	variants VariantCacheRepository,
	posts PostCacheRepository,
	users UserCacheRepository,
	services ServiceCacheRepository,
	metrics MetricRepository,

) *app {
	return &app{
		orders:      orders,
		products:    products,
		variants:    variants,
		posts:       posts,
		users:       users,
		services:    services,
		itemMetrics: metrics,
	}
}

func (a app) GetOrder(ctx context.Context, get GetOrder) (*models.Order, error) {
	return a.orders.Get(ctx, get.OrderID)
}
func (a app) GetUser(ctx context.Context, get GetUser) (*models.User, error) {
	return a.users.Find(ctx, get.UserID)
}

// getUnifiedResults is a shared helper that implements the common logic for UnifiedFeed and GetCatalog
func (a app) getUnifiedResults(ctx context.Context, params UnifiedFeedParams, isUserCatalog bool) (*UnifiedFeedResults, error) {
	// Limit concurrent entity types to prevent resource exhaustion
	const maxConcurrentTypes = 5
	// Return channel for results from each entity type
	type entityResult struct {
		entityType string
		items      []UnifiedSearchResult
		totalCount int64
		err        error
	}

	// Determine which entity types to include
	typesToInclude := params.EntityTypes
	if len(typesToInclude) == 0 {
		// Default to products and posts only
		typesToInclude = []string{
			models.ProductType.String(),
			models.PostType.String(),
			models.ServiceType.String(),
		}
	}
	log.Printf("[getUnifiedResults] Entity types to include: %v (from params: %v)", typesToInclude, params.EntityTypes)

	// Limit the number of entity types to prevent DoS
	if len(typesToInclude) > maxConcurrentTypes {
		typesToInclude = typesToInclude[:maxConcurrentTypes]
	}

	// Create a buffered result channel to avoid blocking goroutines
	resultsCh := make(chan entityResult, len(typesToInclude))

	// Create a context with reasonable timeout
	feedCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// WaitGroup for proper goroutine management
	var wg sync.WaitGroup

	// Launch parallel fetches with prioritization
	for _, entityType := range typesToInclude {
		wg.Add(1)

		// Create entity-specific timeout context
		entityCtx, entityCancel := context.WithTimeout(feedCtx, 1500*time.Millisecond)

		go func(et string, ectx context.Context, ecancel context.CancelFunc) {
			defer wg.Done()
			defer ecancel() // Ensure the context gets canceled when done

			log.Printf("[getUnifiedResults] Processing entity type: %s", et)

			var items []UnifiedSearchResult
			var err error

			// Dispatch to entity-specific fetch function with recoverable error handling
			panicHandler := utils.NewPanicHandler(nil) // Will add logger when available
			panicErr := panicHandler.RecoverWithError(ectx, fmt.Sprintf("%s feed handler", et), func() {
				if isUserCatalog {
					// User catalog mode - get user's items
					switch et {
					case models.ProductType.String():
						items, err = a.getUserProducts(ectx, params)
					case models.PostType.String():
						items, err = a.getUserPosts(ectx, params)
					case models.ServiceType.String():
						items, err = a.getUserServices(ectx, params)
					}
				} else {
					// General feed mode
					switch et {
					case models.ProductType.String():
						items, err = a.getFeedProducts(ectx, params)
					case models.PostType.String():
						items, err = a.getFeedPosts(ectx, params)
					case models.ServiceType.String():
						items, err = a.getFeedServices(ectx, params)
					default:
						log.Printf("[getUnifiedResults] Unknown entity type: %s", et)
					}
				}
			})

			log.Printf("[getUnifiedResults] Entity type %s returned %d items, err=%v", et, len(items), err)

			if panicErr != nil {
				err = panicErr
			}

			// Apply geo filtering here to avoid duplication in each helper
			if err == nil && params.Radius > 0 && params.Lat != 0 && params.Lng != 0 {
				items = filterByGeo(items, params.Lat, params.Lng, params.Radius)
			}

			// Estimate total count for this entity type
			// If we got a full page, there are likely more items
			totalCount := int64(len(items))
			if len(items) >= int(params.PageSize) {
				// Estimate there are at least (current page * page size) items
				totalCount = int64(params.Page)*int64(params.PageSize) + 1
			}

			// Send results through channel even if error (resilient processing)
			select {
			case resultsCh <- entityResult{entityType: et, items: items, totalCount: totalCount, err: err}:
				// Successfully sent
			case <-feedCtx.Done():
				// Parent context canceled, no need to send
			}
		}(entityType, entityCtx, entityCancel)
	}

	// Wait for all goroutines to complete in a separate goroutine
	go func() {
		wg.Wait()
		close(resultsCh) // Close channel when all goroutines are done
	}()

	// Collect results with timeout handling
	var resultsMutex sync.Mutex
	allResults := make([]UnifiedSearchResult, 0)
	errorsByType := &sync.Map{} // Thread-safe map for errors
	var successCount int32      // Use atomic for thread-safe counting
	var totalCount int64        // Track total count across all entity types

	// Create a timer channel for early return if we have sufficient results
	const minResultsForEarlyReturn = 20
	earlyReturnTimer := time.NewTimer(1500 * time.Millisecond) // increased timeout to ensure all entity types are processed

	// Wait for either all results or context deadline
	for i := 0; i < len(typesToInclude); i++ {
		select {
		case result := <-resultsCh:
			if result.err != nil {
				errorsByType.Store(result.entityType, result.err)
			} else {
				atomic.AddInt32(&successCount, 1)
				resultsMutex.Lock()
				allResults = append(allResults, result.items...)
				totalCount += result.totalCount
				resultsMutex.Unlock()
			}

			// Early return only if we have sufficient results from ALL entity types
			// This ensures services aren't skipped when products/posts return quickly
			if len(allResults) >= minResultsForEarlyReturn && atomic.LoadInt32(&successCount) >= int32(len(typesToInclude)) {
				// All entity types have returned results, we can safely return
				break
			}

		case <-earlyReturnTimer.C:
			// Log which entity types have completed
			log.Printf("[getUnifiedResults] Timer expired. Success count: %d/%d, Results: %d",
				atomic.LoadInt32(&successCount), len(typesToInclude), len(allResults))

			// If we have some results after initial timer, check if all entity types had a chance
			if len(allResults) > 0 && atomic.LoadInt32(&successCount) > 0 {
				// Only return early if we've processed most entity types
				if atomic.LoadInt32(&successCount) >= int32(len(typesToInclude)-1) {
					// Cancel context and drain
					cancel()
					go func() {
						drainCtx, drainCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
						defer drainCancel()
						for i := int(atomic.LoadInt32(&successCount)); i < len(typesToInclude); i++ {
							select {
							case <-resultsCh:
							case <-drainCtx.Done():
								return
							}
						}
					}()
					break
				}
			}
			// Otherwise continue waiting for all entity types

		case <-feedCtx.Done():
			if len(allResults) > 0 {
				// Return partial results if we have any
				break
			}
			return nil, ctx.Err()
		}
	}

	// Create the results structure
	results := &UnifiedFeedResults{
		Items:      allResults,
		TotalCount: totalCount,
	}

	// Sort results based on feed type - now optimized to avoid multiple sorts
	switch params.FeedType {
	case "latest":
		// Sort by creation date, newest first
		sort.Slice(results.Items, func(i, j int) bool {
			return results.Items[i].CreatedAt.After(results.Items[j].CreatedAt)
		})
	case "trending":
		// Use a combination of recency and relevance for trending
		sort.Slice(results.Items, func(i, j int) bool {
			scoreI := results.Items[i].RelevanceScore * (1.0 + float64(results.Items[i].CreatedAt.Unix())/86400.0)
			scoreJ := results.Items[j].RelevanceScore * (1.0 + float64(results.Items[j].CreatedAt.Unix())/86400.0)
			return scoreI > scoreJ
		})
	case "recommended":
		// For recommended feed, prioritize by relevance score
		sort.Slice(results.Items, func(i, j int) bool {
			return results.Items[i].RelevanceScore > results.Items[j].RelevanceScore
		})
	default:
		// Default to latest
		sort.Slice(results.Items, func(i, j int) bool {
			return results.Items[i].CreatedAt.After(results.Items[j].CreatedAt)
		})
	}

	// Apply pagination efficiently
	totalItems := len(results.Items)
	startIdx := 0
	endIdx := totalItems

	if params.PageSize > 0 {
		startIdx = int(params.Page-1) * int(params.PageSize)
		endIdx = startIdx + int(params.PageSize)

		// Handle pagination bounds - DO NOT RESET TO PAGE 1
		if startIdx >= totalItems {
			// Return empty results for out-of-bounds pages
			results.Items = []UnifiedSearchResult{}
			results.TotalCount = int64(totalItems)
			results.NextPageToken = "" // No more pages
			return results, nil
		}
		if endIdx > totalItems {
			endIdx = totalItems
		}
	}

	// Apply pagination
	if startIdx < endIdx {
		results.Items = results.Items[startIdx:endIdx]
	} else {
		results.Items = []UnifiedSearchResult{}
	}

	// Generate next page token only if more results exist
	if endIdx < totalItems {
		results.NextPageToken = fmt.Sprintf("%d", params.Page+1)
	}

	return results, nil
}

// Products
func (a app) GetProduct(ctx context.Context, get GetProduct) (*models.Product, error) {
	product, err := a.products.Find(ctx, get.ProductID)
	if err != nil {
		return nil, err
	}

	// Enrich with metrics
	entities := []interface{}{product}
	a.enrichEntitiesWithMetrics(ctx, entities)

	return product, nil
}
func (a app) SuggestProducts(ctx context.Context, get SuggestProducts) ([]*models.Product, error) {
	products, err := a.products.SuggestProducts(ctx, get.Name)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(products))
	for i, p := range products {
		entities[i] = p
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return products, nil
}
func (a app) SearchProductsWithFilters(ctx context.Context, search SearchProductsWithFilters) ([]*models.Product, error) {
	// Check if this is an advanced sort requiring metric service pre-filtering
	if isAdvancedSort(search.SortBy) {
		// ADVANCED SORT WORKFLOW: Metrics Service → Get sorted IDs → Fetch entities
		req := buildMetricSortRequest(
			[]string{"product"},
			search.CategoryID,
			search.MinPrice,
			search.MaxPrice,
			search.Lat,
			search.Lng,
			search.Radius,
			search.Limit,
		)

		// Get pre-sorted metrics from metrics service
		metrics, err := a.executeAdvancedSort(ctx, search.SortBy, req)
		if err != nil {
			return nil, err
		}

		// Fetch full entity details for the sorted IDs
		entities, err := a.fetchEntitiesByMetricResults(ctx, metrics)
		if err != nil {
			return nil, err
		}

		// Convert back to products (metrics already attached)
		products := make([]*models.Product, 0, len(entities))
		for _, entity := range entities {
			if product, ok := entity.(*models.Product); ok {
				products = append(products, product)
			}
		}

		return products, nil
	}

	// BASE SORT WORKFLOW: Search → Get items → Enrich with metrics
	products, err := a.products.SearchWithFilters(
		ctx,
		search.Name,
		search.CategoryID,
		search.CategorySlug,
		search.MinPrice,
		search.MaxPrice,
		search.Brand,
		search.Condition,
		search.Model,
		search.Tags,
		search.ManageStock,
		search.MinStock,
		search.MaxStock,
		search.SKU,
		search.Status,
		search.Negotiable,
		search.UserType,
		search.MiddlemanService,
		search.HasVariants,
		search.ShippingCost,
		search.MinWeight,
		search.MaxWeight,
		search.MinHeight,
		search.MaxHeight,
		search.MinWidth,
		search.MaxWidth,
		search.MinDepth,
		search.MaxDepth,
		search.Offset,
		search.Limit,
		search.Lat,
		search.Lng,
		search.Radius,
		search.Page,
		search.PageSize,
		search.SortBy,
		search.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(products))
	for i, p := range products {
		entities[i] = p
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return products, nil
}
func (a app) SearchProductsWithCategorySlug(ctx context.Context, search SearchProductsWithCategorySlug) ([]*models.Product, error) {
	products, err := a.products.SearchProductsWithCategorySlug(
		ctx,
		search.CategorySlug,
		search.Page,
		search.PageSize,
		search.SortBy,
		search.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(products))
	for i, p := range products {
		entities[i] = p
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return products, nil
}
func (a app) SearchProductsWithCategory(ctx context.Context, search SearchProductsWithCategory) ([]*models.Product, error) {
	products, err := a.products.SearchProductsWithCategory(
		ctx,
		search.CategoryID,
		search.Page,
		search.PageSize,
		search.SortBy,
		search.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(products))
	for i, p := range products {
		entities[i] = p
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return products, nil
}
func (a app) SearchProductsWithTerm(ctx context.Context, search SearchProductsWithTerm) ([]*models.Product, error) {
	products, err := a.products.SearchWithTerm(ctx, search.Name)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(products))
	for i, p := range products {
		entities[i] = p
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return products, nil
}

// POSTS
func (a app) GetPost(ctx context.Context, get GetPost) (*models.Post, error) {
	post, err := a.posts.Find(ctx, get.PostID)
	if err != nil {
		return nil, err
	}

	// Enrich with metrics
	entities := []interface{}{post}
	a.enrichEntitiesWithMetrics(ctx, entities)

	return post, nil
}
func (a app) SuggestPosts(ctx context.Context, get SuggestPosts) ([]*models.Post, error) {
	return a.posts.SuggestPosts(ctx, get.Name)
}
func (a app) SearchPostsWithFilters(ctx context.Context, search SearchPostsWithFilters) ([]*models.Post, error) {
	posts, err := a.posts.SearchPostsWithFilters(
		ctx,
		search.Name,
		search.Description,
		search.Tags,
		search.Status,
		search.UserType,
		search.Thumbnail,
		search.Offset,
		search.Limit,
		search.Lat,
		search.Lng,
		search.Radius,
		search.Page,
		search.PageSize,
		search.SortBy,
		search.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(posts))
	for i, p := range posts {
		entities[i] = p
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return posts, nil
}
func (a app) SearchPostsWithTerm(ctx context.Context, search SearchPostsWithTerm) ([]*models.Post, error) {
	posts, err := a.posts.SearchWithTerm(ctx, search.Name)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(posts))
	for i, p := range posts {
		entities[i] = p
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return posts, nil
}
func (a app) SearchPostsWithCategorySlug(ctx context.Context, search SearchPostsWithCategorySlug) ([]*models.Post, error) {
	posts, err := a.posts.SearchPostsWithCategorySlug(
		ctx,
		search.CategorySlug,
		search.Page,
		search.PageSize,
		search.SortBy,
		search.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(posts))
	for i, p := range posts {
		entities[i] = p
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return posts, nil
}

func (a app) SearchPostsWithCategory(ctx context.Context, search SearchPostsWithCategory) ([]*models.Post, error) {
	posts, err := a.posts.SearchPostsWithCategory(
		ctx,
		search.CategoryID,
		search.Page,
		search.PageSize,
		search.SortBy,
		search.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(posts))
	for i, p := range posts {
		entities[i] = p
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return posts, nil
}

func (a app) GetService(ctx context.Context, get GetService) (*models.Service, error) {
	service, err := a.services.Find(ctx, get.ServiceID)
	if err != nil {
		return nil, err
	}

	// Enrich with metrics
	entities := []interface{}{service}
	a.enrichEntitiesWithMetrics(ctx, entities)

	return service, nil
}
func (a app) SuggestServices(ctx context.Context, get SuggestServices) ([]*models.Service, error) {
	return a.services.SuggestServices(ctx, get.Name)
}
func (a app) SearchServicesWithFilters(ctx context.Context, search SearchServicesWithFilters) ([]*models.Service, error) {
	services, err := a.services.SearchServicesWithFilter(
		ctx,
		search.CategoryID,
		search.CategorySlug,
		search.ServiceType,
		search.UserID,
		search.Status,
		search.SearchText,
		search.MinPrice,
		search.MaxPrice,
		search.AvailableFrom,
		search.AvailableTo,
		search.HasVariants,
		search.Negotiable,
		search.MiddlemanService,
		search.UserType,
		search.Tags,
		search.Qualifications,
		search.Offset,
		search.Limit,
		search.Lat,
		search.Lng,
		search.Radius,
		search.Page,
		search.PageSize,
		search.SortBy,
		search.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(services))
	for i, s := range services {
		entities[i] = s
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return services, nil
}
func (a app) SearchServicesWithTerm(ctx context.Context, search SearchServicesWithTerm) ([]*models.Service, error) {
	services, err := a.services.SearchWithTerm(ctx, search.Name)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(services))
	for i, s := range services {
		entities[i] = s
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return services, nil
}
func (a app) SearchServicesWithCategorySlug(ctx context.Context, search SearchServicesWithCategorySlug) ([]*models.Service, error) {
	services, err := a.services.SearchServicesWithCategorySlug(
		ctx,
		search.CategorySlug,
		search.Page,
		search.PageSize,
		search.SortBy,
		search.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(services))
	for i, s := range services {
		entities[i] = s
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return services, nil
}
func (a app) SearchServicesWithCategory(ctx context.Context, search SearchServicesWithCategory) ([]*models.Service, error) {
	services, err := a.services.SearchServicesWithCategory(
		ctx,
		search.CategoryID,
		search.Page,
		search.PageSize,
		search.SortBy,
		search.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	entities := make([]interface{}, len(services))
	for i, s := range services {
		entities[i] = s
	}

	// Enrich with metrics
	a.enrichEntitiesWithMetrics(ctx, entities)

	return services, nil
}

// UNIFIED SEARCH
func (a app) UnifiedSearch(ctx context.Context, params UnifiedSearchParams) (*UnifiedSearchResults, error) {
	// Set default status if not provided
	if params.Status == "" {
		params.Status = models.StatusActive.String()
	}

	results := &UnifiedSearchResults{
		Results:      []UnifiedSearchResult{},
		CountsByType: make(map[string]int64),
		Page:         params.Page,
		PageSize:     params.PageSize,
	}

	// Determine which entity types to search
	typesToSearch := params.EntityTypes
	if len(typesToSearch) == 0 {
		// Default to products and posts only
		typesToSearch = []string{models.ProductType.String(), models.PostType.String()}
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Use a results channel for safe concurrent access
	type searchResult struct {
		entityType string
		results    []UnifiedSearchResult
		count      int64
		err        error
	}

	resultCh := make(chan searchResult, len(typesToSearch))

	// Launch parallel searches
	for _, entityType := range typesToSearch {
		go func(et string) {
			var typeResults []UnifiedSearchResult
			var typeCount int64
			var err error

			switch et {
			case models.ProductType.String():
				typeResults, typeCount, err = a.searchProducts(ctx, params)
			case models.PostType.String():
				typeResults, typeCount, err = a.searchPosts(ctx, params)
			case models.ServiceType.String():
				typeResults, typeCount, err = a.searchServices(ctx, params)
			}

			resultCh <- searchResult{
				entityType: et,
				results:    typeResults,
				count:      typeCount,
				err:        err,
			}
		}(entityType)
	}

	// Collect results
	for i := 0; i < len(typesToSearch); i++ {
		select {
		case result := <-resultCh:
			if result.err != nil {
				return nil, fmt.Errorf("searching %s: %w", result.entityType, result.err)
			}
			results.Results = append(results.Results, result.results...)
			results.CountsByType[result.entityType] = result.count
			results.TotalCount += result.count
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Sort results by relevance score with stable secondary sort by ID
	sort.Slice(results.Results, func(i, j int) bool {
		// First sort by relevance score (descending)
		if results.Results[i].RelevanceScore != results.Results[j].RelevanceScore {
			return results.Results[i].RelevanceScore > results.Results[j].RelevanceScore
		}

		// If relevance scores are equal, sort by creation date (descending) for stable order
		if !results.Results[i].CreatedAt.Equal(results.Results[j].CreatedAt) {
			return results.Results[i].CreatedAt.After(results.Results[j].CreatedAt)
		}

		// Finally, sort by entity ID for absolute stability
		iID := getEntityID(results.Results[i])
		jID := getEntityID(results.Results[j])
		return iID < jID
	})

	// Apply pagination to final results
	startIdx := 0
	endIdx := len(results.Results)
	if params.PageSize > 0 {
		startIdx = int(params.Page-1) * int(params.PageSize)
		endIdx = startIdx + int(params.PageSize)
		if startIdx >= len(results.Results) {
			startIdx = 0
			endIdx = 0
		}
		if endIdx > len(results.Results) {
			endIdx = len(results.Results)
		}
	}

	if startIdx < endIdx {
		results.Results = results.Results[startIdx:endIdx]
	} else {
		results.Results = []UnifiedSearchResult{}
	}

	return results, nil
}

// Helper methods for each entity type
func (a app) searchProducts(ctx context.Context, params UnifiedSearchParams) ([]UnifiedSearchResult, int64, error) {
	// Calculate proper pagination for this entity type
	// For unified search, we want to fetch more items per entity to ensure we have enough after sorting
	perEntityPageSize := params.PageSize
	if params.PageSize > 0 {
		// Fetch 2x the page size per entity to have enough for sorting and filtering
		perEntityPageSize = params.PageSize * 2
	}

	// Use SearchWithFilters which supports pagination
	products, err := a.products.SearchWithFilters(
		ctx,
		params.SearchTerm,                // name filter
		params.CategoryID,                // categoryID
		params.CategorySlug,              // categorySlug
		params.MinPrice, params.MaxPrice, // price range
		params.Brand,     // brand
		params.Condition, // condition
		params.Model,     // model
		params.Tags,      // tags
		false, 0, 0,      // stock management (no filter)
		"",                      // sku (no filter)
		params.Status,           // status (use from params or default)
		params.Negotiable,       // negotiable
		params.UserType,         // userType
		params.MiddlemanService, // middlemanService
		params.HasVariants,      // hasVariants
		0,                       // shippingCost (no filter)
		0, 0, 0, 0, 0, 0, 0, 0,  // weight, height, width, depth ranges (no filters)
		0, 0, // offset, limit (not used when using page/pageSize)
		params.Lat, params.Lng, params.Radius, // geo params
		params.Page, perEntityPageSize, // pagination
		params.SortBy, params.SortOrder, // use params sort order
	)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for this entity type
	// For now, estimate based on whether we got a full page
	totalCount := int64(len(products))
	if len(products) == int(perEntityPageSize) {
		// If we got a full page, there might be more
		totalCount = int64(params.Page)*perEntityPageSize + 1
	}

	// Convert to interface slice for enrichment
	if len(products) > 0 {
		entities := make([]interface{}, len(products))
		for i, p := range products {
			entities[i] = p
		}
		// Enrich with metrics
		a.enrichEntitiesWithMetrics(ctx, entities)
	}

	results := make([]UnifiedSearchResult, len(products))
	for i, product := range products {
		results[i] = UnifiedSearchResult{
			EntityType:     models.ProductType.String(),
			Product:        product,
			RelevanceScore: calculateRelevanceScore(params.SearchTerm, product.Name, product.Description),
			CreatedAt:      product.CreatedAt,
			UpdatedAt:      product.UpdatedAt,
		}
	}

	return results, totalCount, nil
}

func (a app) searchPosts(ctx context.Context, params UnifiedSearchParams) ([]UnifiedSearchResult, int64, error) {
	// Calculate proper pagination for this entity type
	perEntityPageSize := params.PageSize
	if params.PageSize > 0 {
		// Fetch 2x the page size per entity to have enough for sorting and filtering
		perEntityPageSize = params.PageSize * 2
	}

	// Use SearchPostsWithFilters which supports pagination
	posts, err := a.posts.SearchPostsWithFilters(
		ctx,
		params.SearchTerm, // name
		"",                // description
		params.Tags,       // tags
		params.Status,     // status (use from params or default)
		params.UserType,   // userType
		"",                // thumbnail
		0, 0,              // offset, limit (not used)
		params.Lat, params.Lng, params.Radius, // geo params
		params.Page, perEntityPageSize, // pagination
		params.SortBy, params.SortOrder, // use params sort order
	)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for this entity type
	totalCount := int64(len(posts))
	if len(posts) == int(perEntityPageSize) {
		// If we got a full page, there might be more
		totalCount = int64(params.Page)*perEntityPageSize + 1
	}

	// Convert to interface slice for enrichment
	if len(posts) > 0 {
		entities := make([]interface{}, len(posts))
		for i, p := range posts {
			entities[i] = p
		}
		// Enrich with metrics
		a.enrichEntitiesWithMetrics(ctx, entities)
	}

	results := make([]UnifiedSearchResult, len(posts))
	for i, post := range posts {
		results[i] = UnifiedSearchResult{
			EntityType:     models.PostType.String(),
			Post:           post,
			RelevanceScore: calculateRelevanceScore(params.SearchTerm, post.Name, post.Description),
			CreatedAt:      post.CreatedAt,
			UpdatedAt:      post.UpdatedAt,
		}
	}

	return results, totalCount, nil
}
func (a app) searchServices(ctx context.Context, params UnifiedSearchParams) ([]UnifiedSearchResult, int64, error) {
	// Calculate proper pagination for this entity type
	// For unified search, we want to fetch more items per entity to ensure we have enough after sorting
	perEntityPageSize := params.PageSize
	if params.PageSize > 0 {
		// Fetch 2x the page size per entity to have enough for sorting and filtering
		perEntityPageSize = params.PageSize * 2
	}

	// Use SearchServicesWithFilter which supports pagination and filters
	services, err := a.services.SearchServicesWithFilter(
		ctx,
		params.CategoryID,                  // categoryID
		params.CategorySlug,                // categorySlug
		params.ServiceType,                 // serviceType
		"",                                 // userID (no filter for unified search)
		models.ToStatus(params.Status),     // status (use from params or default)
		params.SearchTerm,                  // searchText
		params.MinPrice,                    // minPrice
		params.MaxPrice,                    // maxPrice
		time.Time{},                        // availableFrom (no filter)
		time.Time{},                        // availableTo (no filter)
		params.HasVariants,                 // hasVariants
		params.Negotiable,                  // negotiable
		params.MiddlemanService,            // middlemanService
		models.ToUserType(params.UserType), // userType
		params.Tags,                        // tags
		nil,                                // qualifications (no filter)
		0, 0,                               // offset, limit (not used when using page/pageSize)
		params.Lat, params.Lng, params.Radius, // geo params
		params.Page, perEntityPageSize, // pagination
		"created_at", "DESC", // stable sort order
	)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for this entity type
	// For now, estimate based on whether we got a full page
	totalCount := int64(len(services))
	if len(services) == int(perEntityPageSize) {
		// If we got a full page, there might be more
		totalCount = int64(params.Page)*perEntityPageSize + 1
	}

	// Convert to interface slice for enrichment
	if len(services) > 0 {
		entities := make([]interface{}, len(services))
		for i, s := range services {
			entities[i] = s
		}
		// Enrich with metrics
		a.enrichEntitiesWithMetrics(ctx, entities)
	}

	results := make([]UnifiedSearchResult, len(services))
	for i, service := range services {
		results[i] = UnifiedSearchResult{
			EntityType:     models.ServiceType.String(),
			Service:        service,
			RelevanceScore: calculateRelevanceScore(params.SearchTerm, service.Name, service.Description),
			CreatedAt:      service.CreatedAt,
			UpdatedAt:      service.UpdatedAt,
		}
	}

	return results, totalCount, nil
}

// calculateRelevanceScore computes a relevance score based on text matching
func calculateRelevanceScore(query, name, description string) float64 {
	query = strings.ToLower(query)
	name = strings.ToLower(name)
	description = strings.ToLower(description)

	score := 0.0

	// Exact name match gets highest score
	if name == query {
		score += 100.0
	} else if strings.Contains(name, query) {
		// Partial name match
		score += 50.0 * (float64(len(query)) / float64(len(name)))
	}

	// Description matches
	if strings.Contains(description, query) {
		score += 25.0 * (float64(len(query)) / float64(len(description)))
	}

	// Normalize score between 0 and 1
	if score > 100.0 {
		score = 100.0
	}
	score = score / 100.0

	return score
}

// getEntityID extracts the ID from a UnifiedSearchResult based on its entity type
func getEntityID(result UnifiedSearchResult) string {
	switch result.EntityType {
	case models.ProductType.String():
		if result.Product != nil {
			return result.Product.ProductID
		}
	case models.PostType.String():
		if result.Post != nil {
			return result.Post.PostID
		}
	}
	return ""
}

func (a app) UnifiedFeed(ctx context.Context, params UnifiedFeedParams) (*UnifiedFeedResults, error) {
	return a.getUnifiedResults(ctx, params, false) // false = general feed mode
}

func (a app) GetCatalog(ctx context.Context, params UnifiedFeedParams) (*UnifiedFeedResults, error) {
	return a.getUnifiedResults(ctx, params, true) // true = user catalog mode
}

// Helper function to filter results by geo location - centralized to avoid duplication
func filterByGeo(items []UnifiedSearchResult, lat, lng float64, radiusKm int64) []UnifiedSearchResult {
	if lat == 0 && lng == 0 || radiusKm <= 0 {
		return items
	}

	filtered := make([]UnifiedSearchResult, 0, len(items))
	for _, item := range items {
		var itemLat, itemLng float64

		// Extract lat/lng based on entity type
		switch {
		case item.Product != nil:
			itemLat, itemLng = item.Product.Lat, item.Product.Lng
		case item.Post != nil:
			itemLat, itemLng = item.Post.Lat, item.Post.Lng
		case item.Service != nil:
			itemLat, itemLng = item.Service.Lat, item.Service.Lng
		default:
			continue // Skip items with no geo data
		}

		if inRadius(lat, lng, itemLat, itemLng, float64(radiusKm)) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// Helper methods for feed generation
func (a app) getFeedProducts(ctx context.Context, params UnifiedFeedParams) ([]UnifiedSearchResult, error) {
	// For feeds, we want to fetch the full page size for each entity type
	// This ensures we have enough results even if one type has fewer items
	perEntityLimit := params.PageSize
	if perEntityLimit <= 0 {
		perEntityLimit = 20 // default
	}

	// Determine sort order based on feed type
	sortBy := "created_at"
	sortOrder := "DESC"
	if params.FeedType == "popular" || params.FeedType == "trending" {
		// For popular/trending, we'll re-sort after enrichment with metrics
		sortBy = "created_at"
		sortOrder = "DESC"
	}

	// Use existing search methods with proper pagination
	products, err := a.products.SearchWithFilters(ctx,
		"",   // name (no filter for feed)
		"",   // categoryID
		"",   // categorySlug
		0, 0, // price range (no filter)
		"", "", "", // brand, condition, model
		nil,         // tags
		false, 0, 0, // stock management
		"",                           // sku
		models.StatusActive.String(), // status
		false,                        // negotiable
		"",                           // userType
		false,                        // middlemanService
		false,                        // hasVariants
		0,                            // shippingCost
		0, 0, 0, 0, 0, 0, 0, 0,       // dimensions (no filter)
		0, 0, // offset, limit
		params.Lat, params.Lng, params.Radius,
		params.Page, perEntityLimit,
		sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	// Apply additional filtering and sorting based on feed type
	// Limit results to perEntityLimit to avoid unnecessary processing
	if len(products) > int(perEntityLimit)*2 {
		products = products[:int(perEntityLimit)*2]
	}

	// Convert to interface slice for enrichment
	if len(products) > 0 {
		entities := make([]interface{}, len(products))
		for i, p := range products {
			entities[i] = p
		}
		// Enrich with metrics
		a.enrichEntitiesWithMetrics(ctx, entities)
	}

	results := make([]UnifiedSearchResult, len(products))
	for i, product := range products {
		results[i] = UnifiedSearchResult{
			EntityType:     models.ProductType.String(),
			Product:        product,
			RelevanceScore: 0.5,               // Default score for feed items
			CreatedAt:      product.CreatedAt, // Use real timestamp
			UpdatedAt:      product.UpdatedAt, // Use real timestamp
		}
	}

	return results, nil
}

func (a app) getFeedPosts(ctx context.Context, params UnifiedFeedParams) ([]UnifiedSearchResult, error) {
	// For feeds, we want to fetch the full page size for each entity type
	perEntityLimit := params.PageSize
	if perEntityLimit <= 0 {
		perEntityLimit = 20 // default
	}

	// Determine sort order based on feed type
	sortBy := "created_at"
	sortOrder := "DESC"
	if params.FeedType == "popular" || params.FeedType == "trending" {
		// For popular/trending, we'll re-sort after enrichment with metrics
		sortBy = "created_at"
		sortOrder = "DESC"
	}

	// Use existing search methods with proper pagination
	posts, err := a.posts.SearchPostsWithFilters(ctx,
		"",                           // name
		"",                           // description
		nil,                          // tags
		models.StatusActive.String(), // status
		"",                           // userType
		"",                           // thumbnail
		0, 0,                         // offset, limit
		params.Lat, params.Lng, params.Radius,
		params.Page, perEntityLimit,
		sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	// Limit results to avoid unnecessary processing
	if len(posts) > int(perEntityLimit)*2 {
		posts = posts[:int(perEntityLimit)*2]
	}

	// Convert to interface slice for enrichment
	if len(posts) > 0 {
		entities := make([]interface{}, len(posts))
		for i, p := range posts {
			entities[i] = p
		}
		// Enrich with metrics
		a.enrichEntitiesWithMetrics(ctx, entities)
	}

	results := make([]UnifiedSearchResult, len(posts))
	for i, post := range posts {
		results[i] = UnifiedSearchResult{
			EntityType:     models.PostType.String(),
			Post:           post,
			RelevanceScore: 0.5,
			CreatedAt:      post.CreatedAt, // Use real timestamp
			UpdatedAt:      post.UpdatedAt, // Use real timestamp
		}
	}

	return results, nil
}
func (a app) getFeedServices(ctx context.Context, params UnifiedFeedParams) ([]UnifiedSearchResult, error) {
	log.Printf("[getFeedServices] Called with page=%d, pageSize=%d", params.Page, params.PageSize)

	// Fix pagination bug: ensure pageSize is used consistently
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20 // Default page size
	}

	// Use existing search methods with NO restrictive filters - use zeros for wide ranges
	services, err := a.services.SearchServicesWithFilter(ctx,
		"",                  // categoryID
		"",                  // categorySlug
		"",                  // serviceType
		"",                  // userID
		models.StatusActive, // status
		"",                  // searchText
		0, 0,                // price range
		time.Time{}, time.Time{}, // availability
		false,    // hasVariants
		false,    // negotiable
		false,    // middlemanService
		"",       // userType (empty = all)
		nil, nil, // tags, qualifications
		0, 0, // offset, limit
		params.Lat, params.Lng, params.Radius,
		params.Page, pageSize,
		"", "")
	if err != nil {
		log.Printf("[getFeedServices] Error calling SearchServicesWithFilter: %v", err)
		return nil, err
	}
	log.Printf("[getFeedServices] SearchServicesWithFilter returned %d services", len(services))

	// Limit results to pageSize to avoid unnecessary processing
	if len(services) > int(pageSize)*2 {
		services = services[:int(pageSize)*2]
	}

	// Convert to interface slice for enrichment
	if len(services) > 0 {
		entities := make([]interface{}, len(services))
		for i, s := range services {
			entities[i] = s
		}
		// Enrich with metrics
		a.enrichEntitiesWithMetrics(ctx, entities)
	}

	results := make([]UnifiedSearchResult, len(services))
	for i, service := range services {
		results[i] = UnifiedSearchResult{
			EntityType:     models.ServiceType.String(),
			Service:        service,
			RelevanceScore: 0.5,
			CreatedAt:      service.CreatedAt, // Use real timestamp
			UpdatedAt:      service.UpdatedAt, // Use real timestamp
		}
	}

	return results, nil
}

// inRadius checks if a point is within the given radius (km) of another point
func inRadius(lat1, lng1, lat2, lng2, radiusKm float64) bool {
	// Haversine formula
	const earthRadiusKm = 6371.0

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lng1Rad := lng1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lng2Rad := lng2 * math.Pi / 180

	// Differences
	dLat := lat2Rad - lat1Rad
	dLng := lng2Rad - lng1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadiusKm * c

	return distance <= radiusKm
}

// enrichEntitiesWithMetrics efficiently attaches metrics to a slice of entities
// with optimized batch processing and proper handling of array responses
func (a app) enrichEntitiesWithMetrics(ctx context.Context, entities []interface{}) {
	if len(entities) == 0 {
		return
	}

	// Check if parent context is already done
	if ctx.Err() != nil {
		return // Don't enrich if context is already canceled
	}

	// Create context with reasonable timeout for metrics enrichment
	enrichCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	// Group entities by type and extract IDs for batch processing
	entityGroups := make(map[string][]string)
	entityPtrs := make(map[string]map[string]interface{})

	// Initialize maps for each entity type
	entityGroups["product"] = make([]string, 0, len(entities))
	entityGroups["post"] = make([]string, 0, len(entities))
	entityGroups["service"] = make([]string, 0, len(entities))
	entityPtrs["product"] = make(map[string]interface{})
	entityPtrs["post"] = make(map[string]interface{})
	entityPtrs["service"] = make(map[string]interface{})
	// Count total entities to prioritize processing
	totalCount := 0

	// Group all entities by type and collect IDs
	for _, entity := range entities {
		switch v := entity.(type) {
		case *models.Product:
			if v.ProductID != "" {
				entityGroups["product"] = append(entityGroups["product"], v.ProductID)
				entityPtrs["product"][v.ProductID] = v
				totalCount++
			}
		case *models.Post:
			if v.PostID != "" {
				entityGroups["post"] = append(entityGroups["post"], v.PostID)
				entityPtrs["post"][v.PostID] = v
				totalCount++
			}
		case *models.Service:
			if v.ID != "" {
				entityGroups["service"] = append(entityGroups["service"], v.ID)
				entityPtrs["service"][v.ID] = v
				totalCount++
			}
		}
	}

	// Optimization: If there's only one entity, use direct GetItemMetric to reduce overhead
	if totalCount == 1 {
		var id string
		var entityType string
		var entityPtr interface{}

		// Find the single entity
		for eType, ids := range entityGroups {
			if len(ids) == 1 {
				id = ids[0]
				entityType = eType
				entityPtr = entityPtrs[eType][id]
				break
			}
		}

		// Use direct call for single entity
		metric, err := a.itemMetrics.GetItemMetric(ctx, id)
		if err != nil || metric == nil {
			return // silent fail
		}

		// Attach metric to entity
		switch entityType {
		case "product":
			if product, ok := entityPtr.(*models.Product); ok {
				product.Metrics = metric
			}
		case "post":
			if post, ok := entityPtr.(*models.Post); ok {
				post.Metrics = metric
			}
		case "service":
			if service, ok := entityPtr.(*models.Service); ok {
				service.Metrics = metric
			}
		}
		return
	}

	// For multiple entities, use batch processing in parallel
	var wg sync.WaitGroup

	// Adaptive batch size - use smaller batches for large result sets to ensure
	// we can process at least some data within timeout
	maxItems := 150 // Default maximum batch size
	if totalCount > 300 {
		maxItems = 100 // Reduce batch size for large result sets
	}

	// Define priorities for entity types (most important first)
	priorities := []string{"product", "post", "service"}

	// Process entity types in priority order
	for _, entityType := range priorities {
		ids := entityGroups[entityType]
		if len(ids) == 0 {
			continue
		}

		wg.Add(1)
		go func(eType string, entityIDs []string) {
			defer wg.Done()

			// Process in batches to respect the 150 item limit
			for i := 0; i < len(entityIDs); i += maxItems {
				// Check if context is canceled before starting a new batch
				select {
				case <-ctx.Done():
					return // Early exit if timeout occurs
				default:
					// Continue processing
				}

				end := i + maxItems
				if end > len(entityIDs) {
					end = len(entityIDs)
				}
				batchIDs := entityIDs[i:end]

				// Create batch context with reasonable timeout
				batchCtx, batchCancel := context.WithTimeout(enrichCtx, 300*time.Millisecond)

				// Fetch metrics for this batch (now returns array, not map)
				metrics, err := a.itemMetrics.GetItemsMetric(batchCtx, batchIDs)
				batchCancel() // Cancel immediately after use
				if err != nil {
					continue // Skip this batch on error, but process others
				}

				// Process each returned metric
				for _, metric := range metrics {
					if metric == nil {
						continue
					}

					id := metric.ID
					// Determine entity type and attach metric
					switch eType {
					case "product":
						if ptr, ok := entityPtrs[eType][id]; ok {
							if product, ok := ptr.(*models.Product); ok {
								product.Metrics = metric
							}
						}
					case "post":
						if ptr, ok := entityPtrs[eType][id]; ok {
							if post, ok := ptr.(*models.Post); ok {
								post.Metrics = metric
							}
						}
					case "service":
						if ptr, ok := entityPtrs[eType][id]; ok {
							if service, ok := ptr.(*models.Service); ok {
								service.Metrics = metric
							}
						}
					}
				}
			}
		}(entityType, ids)
	}

	// Use a timed wait to ensure we don't exceed our timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for all goroutines to complete or context to be canceled
	select {
	case <-done:
		// All processing completed
	case <-enrichCtx.Done():
		// Timeout or cancellation - return partial results
	case <-time.After(450 * time.Millisecond):
		// Safety timeout - ensure we return even if some goroutines are stuck
	}
}

// isAdvancedSort determines if the sort type requires metric service pre-filtering
func isAdvancedSort(sortBy string) bool {
	switch sortBy {
	case constants.SortByLikesHigh, constants.SortByLikesLow,
		constants.SortByCommentsHigh, constants.SortByCommentsLow,
		constants.SortByVisitedHigh, constants.SortByVisitedLow,
		constants.SortByRatingHigh, constants.SortByRatingLow,
		constants.SortByPopularityHigh, constants.SortByTrendingHigh:
		return true
	default:
		return false
	}
}

// getMetricTypeAndDirection extracts metric type and sorting direction from sort string
func getMetricTypeAndDirection(sortBy string) (string, bool) {
	switch sortBy {
	case constants.SortByLikesHigh:
		return constants.MetricTypeLikes, true
	case constants.SortByLikesLow:
		return constants.MetricTypeLikes, false
	case constants.SortByCommentsHigh:
		return constants.MetricTypeComments, true
	case constants.SortByCommentsLow:
		return constants.MetricTypeComments, false
	case constants.SortByVisitedHigh:
		return constants.MetricTypeVisited, true
	case constants.SortByVisitedLow:
		return constants.MetricTypeVisited, false
	case constants.SortByRatingHigh:
		return constants.MetricTypeRating, true
	case constants.SortByRatingLow:
		return constants.MetricTypeRating, false
	default:
		return constants.MetricTypeLikes, true // fallback
	}
}

// executeAdvancedSort implements the Advanced Sort Workflow
func (a app) executeAdvancedSort(ctx context.Context, sortBy string, req MetricSortRequest) ([]*models.ItemMetric, error) {
	metricType, isHighest := getMetricTypeAndDirection(sortBy)

	if isHighest {
		return a.itemMetrics.GetHighestMetricsByType(ctx, metricType, req)
	}
	return a.itemMetrics.GetLowestMetricsByType(ctx, metricType, req)
}

// buildMetricSortRequest creates MetricSortRequest from search parameters
func buildMetricSortRequest(entityTypes []string, categoryID string, minPrice, maxPrice int64,
	lat, lng float64, radius int64, limit int64) MetricSortRequest {
	return MetricSortRequest{
		EntityTypes: entityTypes,
		CategoryId:  categoryID,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		Limit:       int32(limit),
		Lat:         float32(lat),
		Lng:         float32(lng),
		Radius:      float32(radius),
	}
}

// fetchEntitiesByMetricResults retrieves full entity details from metric results
func (a app) fetchEntitiesByMetricResults(ctx context.Context, metrics []*models.ItemMetric) ([]interface{}, error) {
	if len(metrics) == 0 {
		return []interface{}{}, nil
	}

	// Group IDs by entity type
	entityGroups := make(map[string][]string)
	for _, metric := range metrics {
		if metric.EntityType != "" && metric.ID != "" {
			entityGroups[metric.EntityType] = append(entityGroups[metric.EntityType], metric.ID)
		}
	}

	entities := make([]interface{}, 0, len(metrics))

	// Fetch entities by type using batch methods to avoid N+1 queries
	for entityType, ids := range entityGroups {
		switch entityType {
		case "product":
			// Check if repository supports batch fetching
			if batchRepo, ok := a.products.(ProductBatchRepository); ok {
				// Use batch fetch for efficiency
				productMap, err := batchRepo.FindBatch(ctx, ids)
				if err != nil {
					log.Printf("[fetchEntitiesByMetricResults] Batch fetch failed for products: %v", err)
					// Fallback to individual fetches
					for _, id := range ids {
						if product, err := a.products.Find(ctx, id); err == nil && product != nil {
							entities = append(entities, product)
						}
					}
				} else {
					// Add products in original order
					for _, id := range ids {
						if product, found := productMap[id]; found && product != nil {
							entities = append(entities, product)
						}
					}
				}
			} else {
				// Repository doesn't support batch fetch, use individual queries
				for _, id := range ids {
					if product, err := a.products.Find(ctx, id); err == nil && product != nil {
						entities = append(entities, product)
					}
				}
			}
		case "post":
			// Check if repository supports batch fetching
			if batchRepo, ok := a.posts.(PostBatchRepository); ok {
				// Use batch fetch for efficiency
				postMap, err := batchRepo.FindBatch(ctx, ids)
				if err != nil {
					log.Printf("[fetchEntitiesByMetricResults] Batch fetch failed for posts: %v", err)
					// Fallback to individual fetches
					for _, id := range ids {
						if post, err := a.posts.Find(ctx, id); err == nil && post != nil {
							entities = append(entities, post)
						}
					}
				} else {
					// Add posts in original order
					for _, id := range ids {
						if post, found := postMap[id]; found && post != nil {
							entities = append(entities, post)
						}
					}
				}
			} else {
				// Repository doesn't support batch fetch, use individual queries
				for _, id := range ids {
					if post, err := a.posts.Find(ctx, id); err == nil && post != nil {
						entities = append(entities, post)
					}
				}
			}

		case "service":
			for _, id := range ids {
				if service, err := a.services.Find(ctx, id); err == nil && service != nil {
					entities = append(entities, service)
				}
			}
		}
	}

	// Attach metrics to entities
	metricMap := make(map[string]*models.ItemMetric)
	for _, metric := range metrics {
		metricMap[metric.ID] = metric
	}

	for _, entity := range entities {
		var id string
		switch v := entity.(type) {
		case *models.Product:
			id = v.ProductID
			v.Metrics = metricMap[id]
		case *models.Post:
			id = v.PostID
			v.Metrics = metricMap[id]
		case *models.Service:
			id = v.ID
			v.Metrics = metricMap[id]
		}
	}

	return entities, nil
}

func (a app) getUserProducts(ctx context.Context, params UnifiedFeedParams) ([]UnifiedSearchResult, error) {
	// Calculate fair share per entity type to avoid over-fetching
	totalEntityTypes := int64(2) // products and posts only
	if len(params.EntityTypes) > 0 {
		totalEntityTypes = int64(len(params.EntityTypes))
	}

	perEntityLimit := params.PageSize / totalEntityTypes
	if perEntityLimit < 1 {
		perEntityLimit = 1
	}

	// Use existing search methods with distributed pagination
	products, err := a.products.SearchWithFilters(ctx,
		"",   // name
		"",   // categoryID
		"",   // categorySlug
		0, 0, // price range
		"", "", "", // brand, condition, model
		nil,         // tags
		false, 0, 0, // stock management
		"",                           // sku
		models.StatusActive.String(), // status
		false,                        // negotiable
		"",                           // userType
		false,                        // middlemanService
		false,                        // hasVariants
		0,                            // shippingCost
		0, 0, 0, 0, 0, 0, 0, 0,       // dimensions
		0, perEntityLimit, // offset, limit
		params.Lat, params.Lng, params.Radius,
		params.Page, perEntityLimit,
		"", "")
	if err != nil {
		return nil, err
	}

	// Apply additional filtering and sorting based on feed type
	// Limit results to perEntityLimit to avoid unnecessary processing
	if len(products) > int(perEntityLimit)*2 {
		products = products[:int(perEntityLimit)*2]
	}

	// Convert to interface slice for enrichment
	if len(products) > 0 {
		entities := make([]interface{}, len(products))
		for i, p := range products {
			entities[i] = p
		}
		// Enrich with metrics
		a.enrichEntitiesWithMetrics(ctx, entities)
	}

	results := make([]UnifiedSearchResult, len(products))
	for i, product := range products {
		results[i] = UnifiedSearchResult{
			EntityType:     models.ProductType.String(),
			Product:        product,
			RelevanceScore: 0.5,               // Default score for feed items
			CreatedAt:      product.CreatedAt, // Use real timestamp
			UpdatedAt:      product.UpdatedAt, // Use real timestamp
		}
	}

	return results, nil
}

func (a app) getUserPosts(ctx context.Context, params UnifiedFeedParams) ([]UnifiedSearchResult, error) {
	// Fix pagination bug: ensure pageSize is used consistently
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20 // Default page size
	}

	// Use existing search methods
	posts, err := a.posts.GetCatalog(ctx, params.UserID)
	if err != nil {
		return nil, err
	}

	// Convert to interface slice for enrichment
	if len(posts) > 0 {
		entities := make([]interface{}, len(posts))
		for i, p := range posts {
			entities[i] = p
		}
		// Enrich with metrics
		a.enrichEntitiesWithMetrics(ctx, entities)
	}

	results := make([]UnifiedSearchResult, len(posts))
	for i, post := range posts {
		results[i] = UnifiedSearchResult{
			EntityType:     models.PostType.String(),
			Post:           post,
			RelevanceScore: 0.5,
			CreatedAt:      post.CreatedAt, // Use real timestamp
			UpdatedAt:      post.UpdatedAt, // Use real timestamp
		}
	}

	return results, nil
}
func (a app) getUserServices(ctx context.Context, params UnifiedFeedParams) ([]UnifiedSearchResult, error) {
	// Fix pagination bug: ensure pageSize is used consistently
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20 // Default page size
	}

	// Use existing search methods with NO restrictive filters - use zeros for wide ranges
	services, err := a.services.GetCatalog(ctx, params.UserID)
	if err != nil {
		return nil, err
	}

	// Limit results to pageSize to avoid unnecessary processing
	if len(services) > int(pageSize)*2 {
		services = services[:int(pageSize)*2]
	}

	// Convert to interface slice for enrichment
	if len(services) > 0 {
		entities := make([]interface{}, len(services))
		for i, s := range services {
			entities[i] = s
		}
		// Enrich with metrics
		a.enrichEntitiesWithMetrics(ctx, entities)
	}

	results := make([]UnifiedSearchResult, len(services))
	for i, service := range services {
		results[i] = UnifiedSearchResult{
			EntityType:     models.ServiceType.String(),
			Service:        service,
			RelevanceScore: 0.5,
			CreatedAt:      service.CreatedAt, // Use real timestamp
			UpdatedAt:      service.UpdatedAt, // Use real timestamp
		}
	}

	return results, nil
}
