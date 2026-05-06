// File: search/internal/redis/service_cache_repository.go
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/stackus/errors"

	"middleman/internal/di"
	"middleman/search/internal/application"
	"middleman/search/internal/constants"
	"middleman/search/internal/models"
	"middleman/search/internal/utils"
)

type ServiceCacheRepository struct {
	tableName      string
	fallback       application.ServiceRepository
	circuitBreaker *utils.CircuitBreaker
}

var _ application.ServiceCacheRepository = (*ServiceCacheRepository)(nil)

func NewServiceCacheRepository(fallback application.ServiceRepository) *ServiceCacheRepository {
	return &ServiceCacheRepository{
		fallback:       fallback,
		circuitBreaker: utils.NewCircuitBreaker(5, 30*time.Second), // Open after 5 failures, reset after 30s
	}
}

// createIndex verifies that the unified RediSearch index exists and is ready.
// The actual unified index is created by SearchSystem.initRedisearch().
// getPanicHandler creates a panic handler for safe goroutine execution
func (r *ServiceCacheRepository) getPanicHandler() *utils.PanicHandler {
	return utils.NewPanicHandler(func(ctx context.Context, format string, args ...interface{}) {
		log.Printf(format, args...)
	})
}

func (r *ServiceCacheRepository) createIndex(ctx context.Context) error {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	log.Printf("[createIndex] Verifying unified index exists and is ready")

	// Check if the unified index exists and get its info
	info, err := client.Info()
	if err != nil {
		log.Printf("[createIndex] Index verification failed: %v", err)
		return fmt.Errorf("unified index not available: %w", err)
	}

	log.Printf("[createIndex] Unified index verified successfully")
	log.Printf("[createIndex] Index contains %v documents", info.DocCount)

	return nil
}

func (r *ServiceCacheRepository) Add(ctx context.Context,
	id, name, description, serviceType string,
	basePrice int64, pricing []string, availability string,
	providerName, userID, categoryID, categorySlug string,
	descriptionShort, descriptionLong string,
	qualifications []string,
	contact, faq string,
	tags []string,
	status models.Status,
	userType models.UserType,
	shippingCost int64,
	negotiable, hasVariants, middlemanService bool,
	attributes []string,
	options []string,
	thumbnail string,
	lat, long float64, // Using 'long' as per interface
	entityType models.EntityType) error {

	log.Printf("[Add] Adding service %s to search index with status: %s", id, status.String())
	
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	pricingJSON, pErr := json.Marshal(pricing)
	if pErr != nil {
		log.Printf("[Add] Error marshaling pricing for %s: %v", id, pErr)
		pricingJSON = []byte("[]")
	}
	qualJSON, qErr := json.Marshal(qualifications)
	if qErr != nil {
		log.Printf("[Add] Error marshaling qualifications for %s: %v", id, qErr)
		qualJSON = []byte("[]")
	}
	attrJSON, aErr := json.Marshal(attributes)
	if aErr != nil {
		log.Printf("[Add] Error marshaling attributes []string for %s: %v", id, aErr)
		attrJSON = []byte("[]")
	}
	optJSON, oErr := json.Marshal(options)
	if oErr != nil {
		log.Printf("[Add] Error marshaling options []string for %s: %v", id, oErr)
		optJSON = []byte("[]")
	}

	tagsForIndex := ""
	if len(tags) > 0 {
		tagsForIndex = strings.Join(tags, ";")
	}
	// FIXED: Always store location data like product repository does
	locationString := fmt.Sprintf("%.6f,%.6f", long, lat) // Use 'long' from param

	doc := redisearch.NewDocument(id, 1.0).
		Set("service_id", id).
		Set("name", safeString(name)).
		Set("description", safeString(description)).
		Set("service_type", safeString(serviceType)).
		Set("base_price", basePrice).
		Set("pricing", string(pricingJSON)).
		Set("availability", safeString(availability)).
		Set("provider_name", safeString(providerName)).
		Set("user_id", safeString(userID)).
		Set("category_id", safeString(categoryID)).
		Set("category_slug", safeString(categorySlug)).
		Set("qualifications", string(qualJSON)).
		Set("contact", safeString(contact)).
		Set("faq", safeString(faq)).
		Set("tags", tagsForIndex).
		Set("status", status.String()).
		Set("user_type", userType.String()).
		Set("shipping_cost", shippingCost).
		Set("negotiable", boolToInt(negotiable)).
		Set("has_variants", boolToInt(hasVariants)).
		Set("middleman_service", boolToInt(middlemanService)).
		Set("attributes", string(attrJSON)). // Store marshaled []string
		Set("options", string(optJSON)).     // Store marshaled []string
		Set("thumbnail", safeString(thumbnail)).
		Set("entity_type", entityType.String()).
		Set("location", locationString)

	// Use replace option to prevent "Document already exists" errors
	if err := client.IndexOptions(redisearch.IndexingOptions{Replace: true}, doc); err != nil {
		log.Printf("[Add] ERROR indexing service %s: %v", id, err)
		return errors.Wrapf(err, "indexing service %s in RediSearch", id)
	}
	log.Printf("[Add] Successfully indexed service %s with status: %s", id, status.String())
	return nil
}

func (r *ServiceCacheRepository) Remove(ctx context.Context, serviceID string) error {
	if r.fallback != nil {
		log.Printf("[Remove] Fallback removal for service %s not performed (ServiceRepository has no Remove method).", serviceID)
	}
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)
	if err := client.DeleteDocument(serviceID); err != nil {
		if !strings.Contains(strings.ToLower(err.Error()), "unknown document") {
			return errors.Wrapf(err, "removing service %s from RediSearch", serviceID)
		}
		log.Printf("[Remove] Service %s not found in RediSearch.", serviceID)
	}
	return nil
}

func (r *ServiceCacheRepository) UpdateService(ctx context.Context,
	id, name, description, serviceType string,
	basePrice int64, pricing []string, availability string,
	providerName, userID, categoryID, categorySlug string,
	descriptionShort, descriptionLong string,
	qualifications []string,
	contact, faq string,
	tags []string,
	status models.Status,
	userType models.UserType,
	shippingCost int64,
	negotiable, hasVariants, middlemanService bool,
	attributes []string,
	options []string,
	thumbnail string,
	lat, long float64,
	entityType models.EntityType) error {

	log.Printf("[UpdateService] Updating service %s in search index", id)
	
	// Use the same Add method logic to update the document
	// RediSearch will replace the existing document with the same ID
	return r.Add(ctx, id, name, description, serviceType,
		basePrice, pricing, availability,
		providerName, userID, categoryID, categorySlug,
		descriptionShort, descriptionLong,
		qualifications,
		contact, faq,
		tags,
		status,
		userType,
		shippingCost,
		negotiable, hasVariants, middlemanService,
		attributes,
		options,
		thumbnail,
		lat, long,
		entityType)
}

func (r *ServiceCacheRepository) GetCatalog(ctx context.Context, userID string) ([]*models.Service, error) {
	// Call SearchDealsWithFilter with term as name and defaults for other filters/pagination
	return r.fallback.GetCatalog(ctx, userID)
}
func (r *ServiceCacheRepository) Find(ctx context.Context, serviceID string) (*models.Service, error) {
	// Input validation
	if serviceID == "" {
		return nil, errors.ErrInvalidArgument.Msg("serviceID cannot be empty")
	}

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	var service *models.Service
	err := r.circuitBreaker.Call(ctx, func() error {
		// CRITICAL FIX: Escape service ID for TAG field to handle special characters like hyphens
		escapedServiceID := redisearch.EscapeTextFileString(serviceID)
		q := redisearch.NewQuery(fmt.Sprintf("@entity_type:{%s} @service_id:{%s}", models.ServiceType.String(), escapedServiceID)).
			SetReturnFields(r.getReturnFields()...).
			Limit(0, 1)

		docs, _, searchErr := client.Search(q)
		if searchErr != nil {
			return searchErr
		}

		if len(docs) == 0 {
			return errors.ErrNotFound.Msgf("service %s not found", serviceID)
		}

		// Parse the document
		retrievedServiceID := strVal(docs[0].Properties["service_id"])
		if retrievedServiceID == "" {
			retrievedServiceID = docs[0].Id
		}

		var parseErr error
		service, parseErr = r.parseDocToService(docs[0], retrievedServiceID)
		if parseErr != nil {
			return parseErr
		}
		service.ID = retrievedServiceID
		return nil
	})

	if err != nil {
		log.Printf("[Find] RediSearch query error for serviceID=%s: %v. Trying fallback.", serviceID, err)

		// Check if circuit breaker is open
		if errors.Is(err, errors.ErrUnavailable) {
			log.Printf("[Find] Circuit breaker is open, going directly to fallback")
		}

		// Try fallback on any error (including not found)
		if r.fallback != nil {
			return r.fetchFromFallbackAndMaybeReindex(ctx, client, serviceID, serviceID)
		}
		return nil, err
	}

	// If we got here, the circuit breaker call succeeded and we have a service
	return service, nil
}

func (r *ServiceCacheRepository) SuggestServices(ctx context.Context, namePrefix string) ([]*models.Service, error) {
	if namePrefix == "" {
		return []*models.Service{}, nil
	}
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Create a context with timeout to prevent long-running queries
	queryCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// Use QueryBuilder pattern for consistent query construction
	qb := NewQueryBuilder(models.ServiceType)

	// Add name filter with prefix matching
	escapedName := redisearch.EscapeTextFileString(namePrefix)
	qb.WithCustomFilter(fmt.Sprintf("@name:%s*", escapedName))

	// Set pagination to limit results
	qb.WithPagination(0, 10)

	// Set fields to return
	qb.WithReturnFields(r.getReturnFields()...)

	// Build the final query
	_, query := qb.Build()

	// Execute search with timeout context
	docs, _, err := client.Search(query)
	if err != nil {
		if queryCtx.Err() == context.DeadlineExceeded {
			log.Printf("[SuggestServices] query timed out for namePrefix=%s", namePrefix)
		}
		return nil, errors.Wrap(err, "RediSearch suggest services error")
	}

	var suggestions []*models.Service
	for _, doc := range docs {
		serviceKeyInRedis := doc.Id
		retrievedServiceID := strVal(doc.Properties["service_id"])
		if retrievedServiceID == "" {
			retrievedServiceID = serviceKeyInRedis
		}
		service, parseErr := r.parseDocToService(doc, retrievedServiceID)
		if parseErr != nil {
			log.Printf("[SuggestServices] Skipping docID=%s (serviceID %s) parse error: %v", doc.Id, retrievedServiceID, parseErr)
			continue
		}
		service.ID = retrievedServiceID
		suggestions = append(suggestions, service)
	}
	return suggestions, nil
}

func (r *ServiceCacheRepository) SearchWithTerm(ctx context.Context, term string) ([]*models.Service, error) {
	return r.SearchServicesWithFilter(
		ctx,
		"", "", "", "",
		"",
		term,
		0, 0,
		time.Time{}, time.Time{},
		false, false, false,
		"",
		nil, nil,
		0, 20,
		0, 0, 0,
		1, 20,
		"", "",
	)
}

func (r *ServiceCacheRepository) SearchServicesWithCategorySlug(ctx context.Context, categorySlug string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Service, error) {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Create a context with timeout to prevent long-running queries
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Use QueryBuilder pattern for consistent query construction
	qb := NewQueryBuilder(models.ServiceType)

	// Add category slug filter if provided
	if categorySlug != "" {
		qb.WithCustomFilter(fmt.Sprintf("@category_slug:{%s}", redisearch.EscapeTextFileString(categorySlug)))
	}

	// Configure pagination
	finalOffset, finalLimit := calculatePagination(page, pageSize)
	if finalLimit <= 0 {
		finalLimit = 50
	}
	qb.WithPagination(int(finalOffset), int(finalLimit))

	// Set sorting if provided
	if sortBy != "" {
		sortDesc := strings.ToLower(sortOrder) == "desc"
		qb.WithSorting(sortBy, sortDesc)
	}

	// Set fields to return
	qb.WithReturnFields(r.getReturnFields()...)

	// Build the final query
	_, query := qb.Build()

	// Execute search with timeout context
	docs, total, err := client.Search(query)
	if err != nil {
		if queryCtx.Err() == context.DeadlineExceeded {
			log.Printf("[SearchServicesWithCategorySlug] query timed out for categorySlug=%s", categorySlug)
		}
		return nil, errors.Wrapf(err, "RediSearch error in SearchServicesWithCategorySlug")
	}

	if len(docs) == 0 && total == 0 {
		log.Printf("[SearchServicesWithCategorySlug] cache miss for slug '%s' => fallback", categorySlug)
		if r.fallback != nil {
			fallbackServices, fallbackErr := r.fallback.SearchServicesWithCategorySlug(ctx, categorySlug, page, pageSize, sortBy, sortOrder)
			if fallbackErr != nil {
				return nil, errors.Wrap(fallbackErr, "fallback SearchServicesWithCategorySlug error")
			}

			// Reindex asynchronously with proper context, rate limiting, and panic protection
			if len(fallbackServices) > 0 && len(fallbackServices) <= 100 { // Only reindex if reasonable number
				panicHandler := r.getPanicHandler()
				panicHandler.SafeGo(ctx, "service reindexing", func() {
					reindexCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()

					// Rate limit reindexing to prevent overwhelming Redis
					ticker := time.NewTicker(10 * time.Millisecond) // 100 ops/sec max
					defer ticker.Stop()

					for _, s := range fallbackServices {
						if s == nil || s.ID == "" {
							continue
						}

						// CRITICAL FIX: Ensure EntityType is set correctly for reindexing
						if s.EntityType == "" || s.EntityType == models.UnknownType {
							s.EntityType = models.ServiceType
						}

						select {
						case <-reindexCtx.Done():
							return
						case <-ticker.C:
							if err := r.addOrUpdateDoc(reindexCtx, client, s.ID, s); err != nil {
								log.Printf("[WARNING] Failed to reindex service %s: %v", s.ID, err)
							}
						}
					}
				})
			}

			return fallbackServices, nil
		}
		return []*models.Service{}, nil
	}

	var results []*models.Service
	for _, doc := range docs {
		serviceKeyInRedis := doc.Id
		retrievedServiceID := strVal(doc.Properties["service_id"])
		if retrievedServiceID == "" {
			retrievedServiceID = serviceKeyInRedis
		}
		if svc, err := r.parseDocToService(doc, retrievedServiceID); err == nil {
			svc.ID = retrievedServiceID
			results = append(results, svc)
		} else {
			log.Printf("[SearchServicesWithCategorySlug] Parse error: %v", err)
		}
	}
	log.Printf("[SearchServicesWithCategorySlug] returning %d docs, total=%d", len(results), total)
	return results, nil
}

func (r *ServiceCacheRepository) SearchServicesWithCategory(ctx context.Context, categoryId string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Service, error) {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Create a context with timeout to prevent long-running queries
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Use QueryBuilder pattern for consistent query construction
	qb := NewQueryBuilder(models.ServiceType)

	// Add category ID filter if provided
	if categoryId != "" {
		qb.WithCustomFilter(fmt.Sprintf("@category_id:{%s}", redisearch.EscapeTextFileString(categoryId)))
	}

	// Configure pagination
	finalOffset, finalLimit := calculatePagination(page, pageSize)
	if finalLimit <= 0 {
		finalLimit = 50
	}
	qb.WithPagination(int(finalOffset), int(finalLimit))

	// Set sorting if provided
	if sortBy != "" {
		sortDesc := strings.ToLower(sortOrder) == "desc"
		qb.WithSorting(sortBy, sortDesc)
	}

	// Set fields to return
	qb.WithReturnFields(r.getReturnFields()...)

	// Build the final query
	_, query := qb.Build()

	// Execute search with timeout context
	docs, total, err := client.Search(query)
	if err != nil {
		if queryCtx.Err() == context.DeadlineExceeded {
			log.Printf("[SearchServicesWithCategory] query timed out for categoryId=%s", categoryId)
		}
		return nil, errors.Wrapf(err, "RediSearch error in SearchServicesWithCategory")
	}

	if len(docs) == 0 && total == 0 {
		log.Printf("[SearchServicesWithCategory] cache miss for categoryID '%s' => fallback", categoryId)
		if r.fallback != nil {
			fallbackServices, fallbackErr := r.fallback.SearchServicesWithCategory(ctx, categoryId, page, pageSize, sortBy, sortOrder)
			if fallbackErr != nil {
				return nil, errors.Wrap(fallbackErr, "fallback SearchServicesWithCategory error")
			}

			// Reindex asynchronously with proper context, rate limiting, and panic protection
			if len(fallbackServices) > 0 && len(fallbackServices) <= 100 { // Only reindex if reasonable number
				panicHandler := r.getPanicHandler()
				panicHandler.SafeGo(ctx, "service reindexing", func() {
					reindexCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()

					// Rate limit reindexing to prevent overwhelming Redis
					ticker := time.NewTicker(10 * time.Millisecond) // 100 ops/sec max
					defer ticker.Stop()

					for _, s := range fallbackServices {
						if s == nil || s.ID == "" {
							continue
						}

						// CRITICAL FIX: Ensure EntityType is set correctly for reindexing
						if s.EntityType == "" || s.EntityType == models.UnknownType {
							s.EntityType = models.ServiceType
						}

						select {
						case <-reindexCtx.Done():
							return
						case <-ticker.C:
							if err := r.addOrUpdateDoc(reindexCtx, client, s.ID, s); err != nil {
								log.Printf("[WARNING] Failed to reindex service %s: %v", s.ID, err)
							}
						}
					}
				})
			}

			return fallbackServices, nil
		}
		return []*models.Service{}, nil
	}

	var results []*models.Service
	for _, doc := range docs {
		serviceKeyInRedis := doc.Id
		retrievedServiceID := strVal(doc.Properties["service_id"])
		if retrievedServiceID == "" {
			retrievedServiceID = serviceKeyInRedis
		}
		if svc, err := r.parseDocToService(doc, retrievedServiceID); err == nil {
			svc.ID = retrievedServiceID
			results = append(results, svc)
		} else {
			log.Printf("[SearchServicesWithCategory] Parse error: %v", err)
		}
	}
	log.Printf("[SearchServicesWithCategory] returning %d docs, total=%d", len(results), total)
	return results, nil
}

func (r *ServiceCacheRepository) SearchServicesWithFilter(
	ctx context.Context,
	categoryID string,
	categorySlug string,
	serviceType string,
	userID string,
	status models.Status,
	searchText string,
	minPrice int64,
	maxPrice int64,
	availableFrom time.Time,
	availableTo time.Time,
	hasVariants bool,
	negotiable bool,
	middlemanService bool,
	userType models.UserType,
	tags []string,
	qualifications []string,
	offsetParam int64,
	limitParam int64,
	lat float64,
	lng float64,
	radius int64,
	page int64,
	pageSize int64,
	sortBy string,
	sortOrder string,
) ([]*models.Service, error) {
	log.Printf("[SearchServicesWithFilter] Searching with searchText=%q status=%s categoryID=%s categorySlug=%s", 
		searchText, status.String(), categoryID, categorySlug)
	
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Create a context with timeout to prevent long-running queries
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Use QueryBuilder pattern for consistent query construction
	qb := NewQueryBuilder(models.ServiceType)

	// Add search text across multiple fields if provided
	if searchText != "" {
		escTerm := redisearch.EscapeTextFileString(searchText)
		qb.WithCustomFilter(fmt.Sprintf("@name|description|provider_name|tags|service_type:(%s)", escTerm))
	}

	// Add category filters if provided
	if categoryID != "" {
		qb.WithCustomFilter(fmt.Sprintf("@category_id:{%s}", redisearch.EscapeTextFileString(categoryID)))
	}

	if categorySlug != "" {
		qb.WithCustomFilter(fmt.Sprintf("@category_slug:{%s}", redisearch.EscapeTextFileString(categorySlug)))
	}

	// Add service type filter if provided
	if serviceType != "" {
		qb.WithCustomFilter(fmt.Sprintf("@service_type:{%s}", redisearch.EscapeTextFileString(serviceType)))
	}

	// Add user ID filter if provided
	if userID != "" {
		qb.WithCustomFilter(fmt.Sprintf("@user_id:{%s}", redisearch.EscapeTextFileString(userID)))
	}

	// Add status filter if provided
	if statusStr := status.String(); statusStr != "" && statusStr != "unknown" {
		qb.WithStatus(statusStr)
	}

	// Add user type filter if provided
	if userTypeStr := userType.String(); userTypeStr != "" && userTypeStr != "unknown" {
		qb.WithCustomFilter(fmt.Sprintf("@user_type:{%s}", redisearch.EscapeTextFileString(userTypeStr)))
	}

	// Add tags filter if provided
	if len(tags) > 0 {
		tagParts := make([]string, len(tags))
		for i, t := range tags {
			tagParts[i] = redisearch.EscapeTextFileString(t)
		}
		qb.WithCustomFilter(fmt.Sprintf("@tags:{%s}", strings.Join(tagParts, "|")))
	}

	// Add price range filter
	qb.WithPriceRange(minPrice, maxPrice)

	// Add boolean filters
	if hasVariants {
		qb.WithCustomFilter("@has_variants:[1 1]")
	}

	if negotiable {
		qb.WithCustomFilter("@negotiable:[1 1]")
	}

	if middlemanService {
		qb.WithCustomFilter("@middleman_service:[1 1]")
	}

	// Add geo filter if provided
	if lat != 0 && lng != 0 && radius > 0 {
		qb.WithGeoFilter(lat, lng, radius)
	}

	// Configure pagination
	finalOffset, finalLimit := calculatePaginationWithParams(page, pageSize, offsetParam, limitParam)
	if finalLimit <= 0 {
		finalLimit = 50
	}
	qb.WithPagination(int(finalOffset), int(finalLimit))

	// Set sorting if provided
	if sortBy != "" {
		sortDesc := strings.ToLower(sortOrder) == "desc" // Use DESC unless explicitly "ASC"
		qb.WithSorting(sortBy, sortDesc)
	}

	// Set fields to return
	qb.WithReturnFields(r.getReturnFields()...)

	// Build the final query
	queryStr, query := qb.Build()
	log.Printf("[SearchServicesWithFilter] Query: %s", queryStr)

	// Execute search with timeout context
	docs, total, err := client.Search(query)
	log.Printf("[SearchServicesWithFilter] Results: %d docs, total=%d, err=%v", len(docs), total, err)
	if err != nil {
		if queryCtx.Err() == context.DeadlineExceeded {
			log.Printf("[SearchServicesWithFilter] query timed out for lat=%.6f lng=%.6f radius=%d",
				lat, lng, radius)
			// Fall back to a simpler query without geo filtering if that was the issue
			if lat != 0 && lng != 0 && radius > 0 {
				// Create a simpler query without geo filtering
				simpleQb := NewQueryBuilder(models.ServiceType)

				// Add all the same filters except geo
				if searchText != "" {
					escTerm := redisearch.EscapeTextFileString(searchText)
					simpleQb.WithCustomFilter(fmt.Sprintf("@name|description|provider_name|tags|service_type:(%s)", escTerm))
				}
				if categoryID != "" {
					simpleQb.WithCustomFilter(fmt.Sprintf("@category_id:{%s}", redisearch.EscapeTextFileString(categoryID)))
				}
				if categorySlug != "" {
					simpleQb.WithCustomFilter(fmt.Sprintf("@category_slug:{%s}", redisearch.EscapeTextFileString(categorySlug)))
				}
				if serviceType != "" {
					simpleQb.WithCustomFilter(fmt.Sprintf("@service_type:{%s}", redisearch.EscapeTextFileString(serviceType)))
				}
				if userID != "" {
					simpleQb.WithCustomFilter(fmt.Sprintf("@user_id:{%s}", redisearch.EscapeTextFileString(userID)))
				}
				if statusStr := status.String(); statusStr != "" && statusStr != "unknown" {
					simpleQb.WithStatus(statusStr)
				}

				// Configure pagination and return fields
				simpleQb.WithPagination(int(finalOffset), int(finalLimit))
				simpleQb.WithReturnFields(r.getReturnFields()...)

				// Build and execute the simpler query
				_, simpleQuery := simpleQb.Build()
				simpleDocs, _, simpleErr := client.Search(simpleQuery)
				if simpleErr != nil {
					log.Printf("[SearchServicesWithFilter] Simple fallback query also failed: %v", simpleErr)
				} else {
					// Process results from simple query
					var simpleResults []*models.Service
					for _, doc := range simpleDocs {
						serviceKeyInRedis := doc.Id
						retrievedServiceID := strVal(doc.Properties["service_id"])
						if retrievedServiceID == "" {
							retrievedServiceID = serviceKeyInRedis
						}
						p, parseErr := r.parseDocToService(doc, retrievedServiceID)
						if parseErr != nil {
							continue
						}
						p.ID = retrievedServiceID
						simpleResults = append(simpleResults, p)
					}
					if len(simpleResults) > 0 {
						log.Printf("[SearchServicesWithFilter] Simple fallback query returned %d results", len(simpleResults))
						return simpleResults, nil
					}
				}
			}
		}
		return nil, errors.Wrapf(err, "RediSearch query error in SearchServicesWithFilter")
	}

	if len(docs) == 0 {
		log.Printf("[SearchServicesWithFilter] No docs in cache (total=%d) => fallback", total)
		if r.fallback != nil {
			fallbackServices, fallbackErr := r.fallback.SearchServicesWithFilter(ctx, categoryID, categorySlug, serviceType, userID, status, searchText, minPrice, maxPrice, availableFrom, availableTo, hasVariants, negotiable, middlemanService, userType, tags, qualifications, offsetParam, limitParam, lat, lng, radius, page, pageSize, sortBy, sortOrder)
			if fallbackErr != nil {
				return nil, errors.Wrap(fallbackErr, "fallback SearchServicesWithFilter error")
			}

			// Reindex asynchronously with proper context, rate limiting, and panic protection
			if len(fallbackServices) > 0 && len(fallbackServices) <= 100 { // Only reindex if reasonable number
				panicHandler := r.getPanicHandler()
				panicHandler.SafeGo(ctx, "service reindexing", func() {
					reindexCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()

					// Rate limit reindexing to prevent overwhelming Redis
					ticker := time.NewTicker(10 * time.Millisecond) // 100 ops/sec max
					defer ticker.Stop()

					for _, s := range fallbackServices {
						if s == nil || s.ID == "" {
							continue
						}

						// CRITICAL FIX: Ensure EntityType is set correctly for reindexing
						if s.EntityType == "" || s.EntityType == models.UnknownType {
							s.EntityType = models.ServiceType
						}

						select {
						case <-reindexCtx.Done():
							return
						case <-ticker.C:
							if err := r.addOrUpdateDoc(reindexCtx, client, s.ID, s); err != nil {
								log.Printf("[WARNING] Failed to reindex service %s: %v", s.ID, err)
							}
						}
					}
				})
			}

			return fallbackServices, nil
		}
		return []*models.Service{}, nil
	}

	var results []*models.Service
	for _, doc := range docs {
		serviceKeyInRedis := doc.Id
		retrievedServiceID := strVal(doc.Properties["service_id"])
		if retrievedServiceID == "" {
			retrievedServiceID = serviceKeyInRedis
		}
		p, parseErr := r.parseDocToService(doc, retrievedServiceID)
		if parseErr != nil {
			log.Printf("[SearchServicesWithFilter] Skipping docID=%s (serviceID %s) parse error: %v", doc.Id, retrievedServiceID, parseErr)
			continue
		}
		p.ID = retrievedServiceID
		results = append(results, p)
	}
	log.Printf("[SearchServicesWithFilter] returning %d docs from RediSearch, total found=%d", len(results), total)
	return results, nil
}

func (r *ServiceCacheRepository) addOrUpdateDoc(ctx context.Context, client redisearch.Client, serviceID string, s *models.Service) error {
	if s == nil {
		return errors.ErrInvalidArgument.Msg("cannot index nil service")
	}
	if serviceID == "" {
		return errors.ErrInvalidArgument.Msg("serviceID is missing for indexing")
	}

	doc := redisearch.NewDocument(serviceID, 1.0)
	doc.Set("service_id", serviceID)
	doc.Set("name", safeString(s.Name))
	doc.Set("description", safeString(s.Description))
	doc.Set("service_type", safeString(s.ServiceType))
	doc.Set("base_price", s.BasePrice)
	if len(s.Pricing) > 0 {
		pJSON, err := json.Marshal(s.Pricing)
		if err == nil {
			doc.Set("pricing", string(pJSON))
		} else {
			doc.Set("pricing", "[]")
			log.Printf("[addOrUpdateDoc] Error marshaling pricing for %s: %v", serviceID, err)
		}
	} else {
		doc.Set("pricing", "[]")
	}
	doc.Set("availability", safeString(s.Availability))
	doc.Set("provider_name", safeString(s.ProviderName))
	doc.Set("user_id", safeString(s.UserID))
	doc.Set("category_id", safeString(s.CategoryID))
	doc.Set("category_slug", safeString(s.CategorySlug))
	if len(s.Qualifications) > 0 {
		qualJSON, err := json.Marshal(s.Qualifications)
		if err == nil {
			doc.Set("qualifications", string(qualJSON))
		} else {
			doc.Set("qualifications", "[]")
			log.Printf("[addOrUpdateDoc] Error marshaling qualifications for %s: %v", serviceID, err)
		}
	} else {
		doc.Set("qualifications", "[]")
	}
	doc.Set("contact", safeString(s.Contact))
	doc.Set("faq", safeString(s.Faq))
	if len(s.Tags) > 0 {
		doc.Set("tags", strings.Join(s.Tags, ";"))
	} else {
		doc.Set("tags", "")
	}
	doc.Set("status", safeString(s.Status))
	doc.Set("user_type", safeString(s.UserType))
	doc.Set("shipping_cost", s.ShippingCost)
	doc.Set("has_variants", boolToInt(s.HasVariants))
	doc.Set("middleman_service", boolToInt(s.MiddlemanService))
	doc.Set("negotiable", boolToInt(s.Negotiable))
	if len(s.Attributes) > 0 {
		attrJSON, err := json.Marshal(s.Attributes)
		if err == nil {
			doc.Set("attributes", string(attrJSON))
		} else {
			doc.Set("attributes", "[]")
			log.Printf("[addOrUpdateDoc] Error marshaling attributes for %s: %v", serviceID, err)
		}
	} else {
		doc.Set("attributes", "[]")
	}
	if len(s.Options) > 0 {
		optJSON, err := json.Marshal(s.Options)
		if err == nil {
			doc.Set("options", string(optJSON))
		} else {
			doc.Set("options", "[]")
			log.Printf("[addOrUpdateDoc] Error marshaling options for %s: %v", serviceID, err)
		}
	} else {
		doc.Set("options", "[]")
	}
	doc.Set("thumbnail", safeString(s.Thumbnail))
	doc.Set("entity_type", s.EntityType.String())
	// FIXED: Always store location data like product repository does
	doc.Set("location", fmt.Sprintf("%.6f,%.6f", s.Lng, s.Lat))

	// Handle timestamps - set current time if not provided
	now := time.Now()
	createdAt := s.CreatedAt
	updatedAt := s.UpdatedAt

	if createdAt.IsZero() {
		createdAt = now
		s.CreatedAt = createdAt // Update the model
	}
	if updatedAt.IsZero() {
		updatedAt = now
		s.UpdatedAt = updatedAt // Update the model
	}

	doc.Set("created_at", createdAt.Unix())
	doc.Set("updated_at", updatedAt.Unix())

	return client.IndexOptions(redisearch.IndexingOptions{Replace: true}, doc)
}

func (r *ServiceCacheRepository) parseDocToService(doc redisearch.Document, retrievedServiceID string) (*models.Service, error) {
	p := &models.Service{ID: retrievedServiceID}
	var err error

	p.Name = strVal(doc.Properties["name"])
	p.Description = strVal(doc.Properties["description"])
	p.ServiceType = strVal(doc.Properties["service_type"])
	p.BasePrice, err = parseInt64(doc.Properties["base_price"], "base_price", retrievedServiceID)
	if err != nil {
		log.Printf("[parseDocToService] Parse error for service %s: %v", retrievedServiceID, err)
	}
	if rawPricing := strVal(doc.Properties["pricing"]); rawPricing != "" && rawPricing != "[]" {
		if jErr := json.Unmarshal([]byte(rawPricing), &p.Pricing); jErr != nil {
			log.Printf("[parseDocToService] Parse error for service %s: pricing: %v", retrievedServiceID, jErr)
		}
	}
	p.Availability = strVal(doc.Properties["availability"])
	p.ProviderName = strVal(doc.Properties["provider_name"])
	p.UserID = strVal(doc.Properties["user_id"])
	p.CategoryID = strVal(doc.Properties["category_id"])
	p.CategorySlug = strVal(doc.Properties["category_slug"])
	if rawQual := strVal(doc.Properties["qualifications"]); rawQual != "" && rawQual != "[]" {
		if jErr := json.Unmarshal([]byte(rawQual), &p.Qualifications); jErr != nil {
			log.Printf("[parseDocToService] Parse error for service %s: qualifications: %v", retrievedServiceID, jErr)
		}
	}
	p.Contact = strVal(doc.Properties["contact"])
	p.Faq = strVal(doc.Properties["faq"])
	if rawTags := strVal(doc.Properties["tags"]); rawTags != "" {
		p.Tags = strings.Split(rawTags, ";")
	}
	p.Status = strVal(doc.Properties["status"])
	p.UserType = strVal(doc.Properties["user_type"])
	p.ShippingCost, err = parseInt64(doc.Properties["shipping_cost"], "shipping_cost", retrievedServiceID)
	if err != nil {
		log.Printf("[parseDocToService] Parse error for service %s: %v", retrievedServiceID, err)
	}
	p.HasVariants = parseBoolVal(doc.Properties["has_variants"], "has_variants", retrievedServiceID)
	p.MiddlemanService = parseBoolVal(doc.Properties["middleman_service"], "middleman_service", retrievedServiceID)
	p.Negotiable = parseBoolVal(doc.Properties["negotiable"], "negotiable", retrievedServiceID)
	if rawAttrs := strVal(doc.Properties["attributes"]); rawAttrs != "" && rawAttrs != "[]" {
		if jErr := json.Unmarshal([]byte(rawAttrs), &p.Attributes); jErr != nil {
			log.Printf("[parseDocToService] Parse error for service %s: attributes: %v", retrievedServiceID, jErr)
		}
	}
	if rawOpts := strVal(doc.Properties["options"]); rawOpts != "" && rawOpts != "[]" {
		if jErr := json.Unmarshal([]byte(rawOpts), &p.Options); jErr != nil {
			log.Printf("[parseDocToService] Parse error for service %s: options: %v", retrievedServiceID, jErr)
		}
	}
	p.Thumbnail = strVal(doc.Properties["thumbnail"])
	ltStr := strVal(doc.Properties["entity_type"])
	lt := models.ToEntityType(ltStr)
	p.EntityType = lt

	if rawLoc := strVal(doc.Properties["location"]); rawLoc != "" {
		parts := strings.Split(rawLoc, ",")
		if len(parts) == 2 {
			if lngF, lErr := strconv.ParseFloat(parts[0], 64); lErr == nil {
				p.Lng = lngF
			} else {
				log.Printf("[parseDocToService] Parse error for service %s: lng: %v", retrievedServiceID, lErr)
			}
			if latF, lErr := strconv.ParseFloat(parts[1], 64); lErr == nil {
				p.Lat = latF
			} else {
				log.Printf("[parseDocToService] Parse error for service %s: lat: %v", retrievedServiceID, lErr)
			}
		} else {
			log.Printf("[parseDocToService] Parse warning for service %s: location malformed: '%s'", retrievedServiceID, rawLoc)
		}
	}

	// Parse timestamps => Unix timestamps
	if createdAtUnix, err := parseInt64(doc.Properties["created_at"], "created_at", retrievedServiceID); err == nil && createdAtUnix > 0 {
		p.CreatedAt = time.Unix(createdAtUnix, 0)
	}
	if updatedAtUnix, err := parseInt64(doc.Properties["updated_at"], "updated_at", retrievedServiceID); err == nil && updatedAtUnix > 0 {
		p.UpdatedAt = time.Unix(updatedAtUnix, 0)
	}

	return p, nil
}

func (r *ServiceCacheRepository) fallbackForWrongType(ctx context.Context, client redisearch.Client, docKey string, serviceID string) (*models.Service, error) {
	log.Printf("[fallbackForWrongType] Removing mismatched service doc key %s from cache (logical serviceID %s)", docKey, serviceID)
	if err := client.DeleteDocument(docKey); err != nil {
		log.Printf("[fallbackForWrongType] could not delete mismatched docKey=%s: %v", docKey, err)
	}
	if r.fallback != nil {
		fbService, err := r.fallback.Find(ctx, serviceID)
		if err != nil {
			return nil, errors.Wrapf(err, "fallbackForWrongType => fallback Find error for serviceID %s", serviceID)
		}
		if fbService == nil {
			log.Printf("[fallbackForWrongType] serviceID %s not found in fallback either.", serviceID)
			return nil, nil
		}
		if fbService.ID == "" {
			fbService.ID = serviceID
		}
		log.Printf("[fallbackForWrongType] Re-indexing correct service %s from fallback.", serviceID)
		panicHandler := r.getPanicHandler()
		panicHandler.SafeGo(ctx, "service reindexing wrong type fix", func() {
			reindexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := r.addOrUpdateDoc(reindexCtx, client, serviceID, fbService); err != nil {
				log.Printf("[WARNING] Failed to reindex service %s after wrong type fix: %v", serviceID, err)
			}
		})
		return fbService, nil
	}
	return nil, errors.ErrInvalidArgument.Msg("fallback repository not configured")
}

func (r *ServiceCacheRepository) fetchFromFallbackAndMaybeReindex(ctx context.Context, client redisearch.Client, docKey string, serviceID string) (*models.Service, error) {
	if r.fallback == nil {
		return nil, errors.ErrNotFound
	}
	fbService, fbErr := r.fallback.Find(ctx, serviceID)
	if fbErr != nil {
		return nil, errors.Wrapf(fbErr, "fallback find error for service %s", serviceID)
	}
	if fbService == nil {
		return nil, errors.ErrNotFound
	}
	if fbService.ID == "" {
		fbService.ID = serviceID
	}

	panicHandler := r.getPanicHandler()
	panicHandler.SafeGo(ctx, "service reindexing from fallback", func() {
		reindexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := r.addOrUpdateDoc(reindexCtx, client, serviceID, fbService); err != nil {
			log.Printf("[WARNING] Failed to reindex service %s from fallback: %v", serviceID, err)
		}
	})

	return fbService, nil
}

// FindBatch retrieves multiple services by their IDs using parallel fetches for efficiency
func (r *ServiceCacheRepository) FindBatch(ctx context.Context, serviceIDs []string) (map[string]*models.Service, error) {
	if len(serviceIDs) == 0 {
		return make(map[string]*models.Service), nil
	}

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)
	result := make(map[string]*models.Service, len(serviceIDs))

	// First, try to get from Redis using individual document fetches
	// RediSearch doesn't have a native batch get, so we'll fetch in parallel
	type fetchResult struct {
		id      string
		service *models.Service
		err     error
	}

	ch := make(chan fetchResult, len(serviceIDs))

	// Use worker pool to limit concurrent fetches
	const maxWorkers = 10
	sem := make(chan struct{}, maxWorkers)

	panicHandler := r.getPanicHandler()

	for _, id := range serviceIDs {
		sem <- struct{}{} // Acquire semaphore
		panicHandler.SafeGo(ctx, "service batch fetch", func() {
			defer func() { <-sem }() // Release semaphore

			serviceID := id // Capture loop variable

			// Use circuit breaker for each fetch
			var service *models.Service
			err := r.circuitBreaker.Call(ctx, func() error {
				escapedServiceID := redisearch.EscapeTextFileString(serviceID)
				q := redisearch.NewQuery(fmt.Sprintf("@entity_type:{%s} @service_id:{%s}", models.ServiceType.String(), escapedServiceID)).
					SetReturnFields(r.getReturnFields()...).
					Limit(0, 1)

				docs, _, searchErr := client.Search(q)
				if searchErr != nil {
					return searchErr
				}

				if len(docs) == 0 {
					return errors.ErrNotFound.Msgf("service %s not found", serviceID)
				}

				retrievedServiceID := strVal(docs[0].Properties["service_id"])
				if retrievedServiceID == "" {
					retrievedServiceID = docs[0].Id
				}

				var parseErr error
				service, parseErr = r.parseDocToService(docs[0], retrievedServiceID)
				if parseErr != nil {
					return parseErr
				}
				service.ID = retrievedServiceID
				return nil
			})

			ch <- fetchResult{id: serviceID, service: service, err: err}
		})
	}

	// Collect results and track missing IDs
	var missingIDs []string
	for i := 0; i < len(serviceIDs); i++ {
		res := <-ch
		if res.err != nil {
			missingIDs = append(missingIDs, res.id)
		} else if res.service != nil {
			result[res.id] = res.service
		}
	}

	// If any IDs are missing, try fallback
	if len(missingIDs) > 0 && r.fallback != nil {
		log.Printf("[FindBatch] %d services not in cache, trying fallback", len(missingIDs))

		// Check if fallback supports batch fetch
		if batchFallback, ok := r.fallback.(application.ServiceBatchRepository); ok {
			fallbackServices, err := batchFallback.FindBatch(ctx, missingIDs)
			if err != nil {
				log.Printf("[FindBatch] Fallback batch fetch failed: %v", err)
			} else {
				// Add fallback results to result map and reindex them
				if len(fallbackServices) > 0 && len(fallbackServices) <= 100 { // Only reindex reasonable number
					for id, service := range fallbackServices {
						result[id] = service
						// Reindex asynchronously with rate limiting
						panicHandler.SafeGo(ctx, "service batch reindexing", func() {
							s := service // Capture loop variable
							reindexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
							defer cancel()

							if err := r.addOrUpdateDoc(reindexCtx, client, s.ID, s); err != nil {
								log.Printf("[FindBatch] Failed to reindex service %s: %v", s.ID, err)
							}
						})
					}
				}
			}
		} else {
			// Fallback doesn't support batch, fetch individually
			for _, id := range missingIDs {
				service, err := r.fallback.Find(ctx, id)
				if err == nil && service != nil {
					result[id] = service
					// Reindex asynchronously
					panicHandler.SafeGo(ctx, "service individual reindexing", func() {
						s := service // Capture variable
						reindexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()

						if err := r.addOrUpdateDoc(reindexCtx, client, s.ID, s); err != nil {
							log.Printf("[FindBatch] Failed to reindex service %s: %v", s.ID, err)
						}
					})
				}
			}
		}
	}

	return result, nil
}

func (r *ServiceCacheRepository) getReturnFields() []string {
	return []string{
		"service_id", "name", "description", "service_type", "base_price", "pricing",
		"availability", "description_short", "description_long",
		"provider_name", "user_id", "category_id", "category_slug", "qualifications",
		"contact", "faq", "tags", "status", "user_type", "shipping_cost",
		"negotiable", "has_variants", "middleman_service", "attributes", "options",
		"thumbnail", "entity_type", "location", "created_at", "updated_at",
	}
}
