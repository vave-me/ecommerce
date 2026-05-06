// File: search/internal/redis/post_cache_repository.go
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

// PostCacheRepository is an implementation of application.PostCacheRepository
// that uses RediSearch as a "cache"/search layer, with a fallback PostRepository
// for deeper DB or gRPC calls.
type PostCacheRepository struct {
	fallback       application.PostRepository
	circuitBreaker *utils.CircuitBreaker
}

// Ensure we implement the PostCacheRepository interface.
var _ application.PostCacheRepository = (*PostCacheRepository)(nil)

// NewPostCacheRepository constructs a RediSearch-based repo plus a fallback.
func NewPostCacheRepository(fallback application.PostRepository) *PostCacheRepository {
	return &PostCacheRepository{
		fallback:       fallback,
		circuitBreaker: utils.NewCircuitBreaker(5, 30*time.Second), // Open after 5 failures, reset after 30s
	}
}

// getPanicHandler creates a panic handler for safe goroutine execution
func (r *PostCacheRepository) getPanicHandler() *utils.PanicHandler {
	return utils.NewPanicHandler(func(ctx context.Context, format string, args ...interface{}) {
		log.Printf(format, args...)
	})
}

// -----------------------------------------------------------------------------
// Implementation of PostCacheRepository interface
// -----------------------------------------------------------------------------

// Add indexes a new post in RediSearch. Optionally also add to fallback DB.
func (r *PostCacheRepository) Add(
	ctx context.Context,
	postID string,
	name string,
	description string,
	userID string,
	tags []string,
	status string,
	lat float64,
	lng float64,
	thumbnail string,
	entityType models.EntityType,
) error {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	tagsJSON, _ := json.Marshal(tags)
	locationString := fmt.Sprintf("%.6f,%.6f", lng, lat)

	doc := redisearch.NewDocument(postID, 1.0).
		Set("post_id", postID).
		Set("name", name).
		Set("description", description).
		Set("user_id", userID).
		Set("thumbnail", thumbnail).
		Set("status", status).
		Set("tags", string(tagsJSON)).
		Set("entity_type", entityType.String()). // Explicitly setting to "post" - PostType value
		Set("location", locationString)

	// Use replace option to prevent "Document already exists" errors
	if err := client.IndexOptions(redisearch.IndexingOptions{Replace: true}, doc); err != nil {
		return errors.Wrapf(err, "indexing post %s in RediSearch", postID)
	}
	return nil
}
func (r *PostCacheRepository) GetCatalog(ctx context.Context, userID string) ([]*models.Post, error) {
	// Call SearchDealsWithFilter with term as name and defaults for other filters/pagination
	return r.fallback.GetCatalog(ctx, userID)
}

// Rebrand is an example partial update. We'll call Update for re-indexing.
func (r *PostCacheRepository) Rebrand(
	ctx context.Context,
	postID string,
	name string,
	description string,
) error {
	// Example: you might update the fallback DB first, then reindex
	return r.Update(ctx, postID, 0.0)
}

// Update removes the old doc from RediSearch, fetches from fallback, and re-indexes.
func (r *PostCacheRepository) Update(ctx context.Context, postID string, _ float64) error {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	if err := client.DeleteDocument(postID); err != nil {
		log.Printf("[WARNING] Failed to delete document %s from RediSearch: %v", postID, err)
	}

	// Retrieve updated post from fallback
	updated, err := r.fallback.Find(ctx, postID)
	if err != nil {
		return errors.Wrap(err, "finding updated post in fallback for reindex")
	}
	if updated == nil {
		return errors.Wrapf(errors.ErrNotFound, "no fallback post found for ID=%s", postID)
	}

	// Re-index the updated post
	return r.addOrUpdateDoc(ctx, client, updated)
}
func (r *PostCacheRepository) UpdatePost(
	ctx context.Context,
	postID, name, description string,
	tags []string,
	status, thumbnail string,
) error {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// 1. Retrieve the existing post from the fallback (database).
	existingPost, err := r.fallback.Find(ctx, postID)
	if err != nil {
		return errors.Wrap(err, "finding post in fallback for update")
	}
	if existingPost == nil {
		return errors.Wrapf(errors.ErrNotFound, "no fallback post found for ID=%s", postID)
	}

	// 2. Apply the updates to the existing post data.
	existingPost.Name = name
	existingPost.Description = description
	existingPost.Tags = tags
	existingPost.Status = status
	existingPost.Thumbnail = thumbnail

	if err := client.DeleteDocument(postID); err != nil {
		log.Printf("[WARNING] Failed to delete document %s from RediSearch: %v", postID, err)
	}

	// 4. Re-index the updated post in RediSearch.
	return r.addOrUpdateDoc(ctx, client, existingPost)
}

// UpdateThumbnail updates only the thumbnail in Redis (and optionally fallback, if your
// fallback supports that). It then reindexes the post in RediSearch.
func (r *PostCacheRepository) UpdateThumbnail(ctx context.Context, postID string, thumbnail string) error {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	if err := client.DeleteDocument(postID); err != nil {
		log.Printf("[WARNING] Failed to delete document %s from RediSearch: %v", postID, err)
	}

	// Retrieve the post from fallback
	updated, err := r.fallback.Find(ctx, postID)
	if err != nil {
		return errors.Wrap(err, "finding post in fallback for thumbnail update")
	}
	if updated == nil {
		return errors.Wrapf(errors.ErrNotFound, "no fallback post found for ID=%s", postID)
	}

	// If your fallback DB also needs the updated thumbnail, you would do so here
	// (e.g., r.fallback.UpdateThumbnail(...)) if that method existed.

	// Temporarily override the thumbnail for the reindex
	updated.Thumbnail = thumbnail

	// Re-index with the updated thumbnail
	return r.addOrUpdateDoc(ctx, client, updated)
}

// Remove deletes from fallback plus from RediSearch.
func (r *PostCacheRepository) Remove(ctx context.Context, postID string) error {
	// Delete from fallback if needed:
	if err := r.fallback.Remove(ctx, postID); err != nil {
		return errors.Wrap(err, "removing from fallback DB")
	}

	// Also remove from RediSearch
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)
	if err := client.DeleteDocument(postID); err != nil {
		return errors.Wrapf(err, "removing post %s from RediSearch", postID)
	}
	return nil
}

// Get is a trivial pass-through to fallback for a single post by ID.
func (r *PostCacheRepository) Get(ctx context.Context, postID string) (*models.Post, error) {
	return r.fallback.Find(ctx, postID)
}

// Find tries RediSearch. If doc missing => fallback => reindex if found.
func (r *PostCacheRepository) Find(ctx context.Context, postID string) (*models.Post, error) {
	// Input validation
	if postID == "" {
		return nil, errors.ErrInvalidArgument.Msg("postID cannot be empty")
	}

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	var post *models.Post
	err := r.circuitBreaker.Call(ctx, func() error {
		// CRITICAL FIX: Escape post ID for TAG field to handle special characters like hyphens
		escapedPostID := redisearch.EscapeTextFileString(postID)
		q := redisearch.NewQuery(fmt.Sprintf("@entity_type:{%s} @post_id:{%s}", models.PostType.String(), escapedPostID)).
			SetReturnFields(
				"post_id", "name", "description", "user_id",
				"thumbnail", "status", "tags",
				"entity_type",
				"location", "created_at", "updated_at",
			).
			Limit(0, 1)

		docs, _, searchErr := client.Search(q)
		if searchErr != nil {
			return searchErr
		}
		
		if len(docs) == 0 {
			return errors.ErrNotFound.Msgf("post %s not found", postID)
		}
		
		// Parse the document
		var parseErr error
		post, parseErr = r.parseDocToPost(docs[0])
		return parseErr
	})

	if err != nil {
		log.Printf("[Find] RediSearch query error for postID=%s: %v. Trying fallback.", postID, err)
		
		// Check if circuit breaker is open
		if errors.Is(err, errors.ErrUnavailable) {
			log.Printf("[Find] Circuit breaker is open, going directly to fallback")
		}
		
		// Try fallback on any error (including not found)
		if r.fallback != nil {
			return r.fetchFromFallbackAndMaybeReindex(ctx, client, postID)
		}
		return nil, err
	}

	// If we got here, the circuit breaker call succeeded and we have a post
	return post, nil
}

// SuggestPosts uses a prefix search on "name" field in RediSearch.
func (r *PostCacheRepository) SuggestPosts(ctx context.Context, searchName string) ([]*models.Post, error) {
	if len(searchName) == 0 {
		return []*models.Post{}, nil
	}

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Create a context with timeout
	queryCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// Use QueryBuilder pattern for consistent query construction
	qb := NewQueryBuilder(models.PostType)

	// Add name filter with prefix matching
	escapedName := redisearch.EscapeTextFileString(searchName)
	qb.WithCustomFilter(fmt.Sprintf("@name:%s*", escapedName))

	// Set pagination to limit results
	qb.WithPagination(0, 10)

	// Set fields to return
	qb.WithReturnFields(
		"post_id", "name", "description", "user_id",
		"thumbnail", "status", "tags", "entity_type", "location",
		"created_at", "updated_at",
	)

	// Build the final query
	_, query := qb.Build()

	// Execute search with timeout context
	docs, _, err := client.Search(query)
	if err != nil {
		if queryCtx.Err() == context.DeadlineExceeded {
			log.Printf("[SuggestPosts] query timed out for searchName=%s", searchName)
		}
		return nil, errors.Wrap(err, "RediSearch suggest posts error")
	}

	var suggestions []*models.Post
	for _, doc := range docs {
		parsed, parseErr := r.parseDocToPost(doc)
		if parseErr != nil {
			log.Printf("[SuggestPosts] skip docID=%s parse error: %v", doc.Id, parseErr)
			continue
		}
		// Exclude or fallback if entityType != PostType
		if parsed.EntityType != models.PostType {
			log.Printf("[SuggestPosts] docID=%s entityType=%s => skip or fallback",
				doc.Id, parsed.EntityType)
			// You can skip, or do an inline fallback if you wish:
			// _, _ = r.fallbackForWrongType(ctx, client, doc.Id)
			continue
		}
		suggestions = append(suggestions, parsed)
	}
	return suggestions, nil
}

// SearchWithTerm delegates to SearchWithFilters with minimal filters.
func (r *PostCacheRepository) SearchWithTerm(ctx context.Context, name string) ([]*models.Post, error) {
	return r.SearchPostsWithFilters(
		ctx,
		name, // description
		"",   // description
		nil,  // tags
		"",   // status
		"",
		"", // thumbnail
		0,  // offset
		0,  // limit
		0,  // lat
		0,  // lng
		0,  // radius
		0,  // page
		0,  // pageSize
		"", // sortBy
		"", // sortOrder
	)
}

func (r *PostCacheRepository) SearchPostsWithCategorySlug(ctx context.Context, categorySlug string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error) {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Create a context with timeout to prevent long-running queries
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Use QueryBuilder pattern for consistent query construction
	qb := NewQueryBuilder(models.PostType)

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
		"post_id", "name", "description", "user_id",
		"category_id", "category_slug", "tags", "status", "user_type", "entity_type", "thumbnail",
		"location", "created_at", "updated_at",
	)

	// Build the final query
	_, query := qb.Build()

	// Execute search with timeout context
	docs, total, err := client.Search(query)
	if err != nil {
		if queryCtx.Err() == context.DeadlineExceeded {
			log.Printf("[SearchPostsWithCategorySlug] query timed out for categorySlug=%s", categorySlug)
		}
		return nil, errors.Wrap(err, "RediSearch query error in SearchPostsWithCategorySlug")
	}

	// If nothing found in Redis => fallback => reindex
	if len(docs) == 0 {
		log.Printf("[SearchPostsWithCategorySlug] no docs => fallback (categorySlug=%q)", categorySlug)

		fallbackPosts, fallbackErr := r.fallback.SearchPostsWithCategorySlug(
			ctx,
			categorySlug,
			page,
			pageSize,
			sortBy,
			sortOrder,
		)
		if fallbackErr != nil {
			return nil, fallbackErr
		}

		// Reindex asynchronously with proper context, rate limiting, and panic protection
		if len(fallbackPosts) > 0 && len(fallbackPosts) <= 100 { // Only reindex if reasonable number
			panicHandler := r.getPanicHandler()
			panicHandler.SafeGo(ctx, "post reindexing for category slug", func() {
				reindexCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				
				// Rate limit reindexing to prevent overwhelming Redis
				ticker := time.NewTicker(10 * time.Millisecond) // 100 ops/sec max
				defer ticker.Stop()
				
				for _, p := range fallbackPosts {
					select {
					case <-reindexCtx.Done():
						return
					case <-ticker.C:
						if err := r.addOrUpdateDoc(reindexCtx, client, p); err != nil {
							log.Printf("[WARNING] Failed to reindex post %s: %v", p.PostID, err)
						}
					}
				}
			})
		}

		return fallbackPosts, nil
	}

	// Otherwise parse all the found docs
	var results []*models.Post
	for _, doc := range docs {
		p, parseErr := r.parseDocToPost(doc)
		if parseErr != nil {
			log.Printf("[SearchPostsWithCategorySlug] skipping docID=%s parse err: %v", doc.Id, parseErr)
			continue
		}
		// If entityType != PostType => fallback for a "correct" doc
		if p.EntityType != models.PostType {
			log.Printf("[SearchPostsWithCategorySlug] docID=%s entityType=%s => fallback", doc.Id, p.EntityType)
			correctPost, fbErr := r.fallbackForWrongType(ctx, client, doc.Id)
			if fbErr != nil {
				log.Printf("[SearchPostsWithCategorySlug] fallback error docID=%s: %v", doc.Id, fbErr)
				continue
			}
			if correctPost != nil {
				results = append(results, correctPost)
			}
			continue
		}
		results = append(results, p)
	}

	log.Printf("[SearchPostsWithCategorySlug] returning %d docs, total=%d", len(results), total)
	return results, nil
}

func (r *PostCacheRepository) SearchPostsWithCategory(
	ctx context.Context,
	categoryId string,
	page int64,
	pageSize int64,
	sortBy string,
	sortOrder string,
) ([]*models.Post, error) {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Create a context with timeout to prevent long-running queries
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Use QueryBuilder pattern for consistent query construction
	qb := NewQueryBuilder(models.PostType)

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
		"post_id", "name", "description", "user_id",
		"category_id", "category_slug", "tags", "status", "user_type",
		"options", "entity_type", "thumbnail",
		"location",
	)

	// Build the final query
	_, query := qb.Build()

	// Execute search with timeout context
	docs, total, err := client.Search(query)
	if err != nil {
		if queryCtx.Err() == context.DeadlineExceeded {
			log.Printf("[SearchPostsWithCategory] query timed out for categoryId=%s", categoryId)
		}
		return nil, errors.Wrap(err, "RediSearch query error in SearchPostsWithCategory")
	}

	// If nothing found in Redis => fallback => reindex
	if len(docs) == 0 {
		log.Printf("[SearchPostsWithCategory] no docs => fallback (categoryId=%q)", categoryId)

		fallbackPosts, fallbackErr := r.fallback.SearchPostsWithCategory(
			ctx,
			categoryId,
			page,
			pageSize,
			sortBy,
			sortOrder,
		)
		if fallbackErr != nil {
			return nil, fallbackErr
		}

		// Reindex asynchronously with proper context, rate limiting, and panic protection
		if len(fallbackPosts) > 0 && len(fallbackPosts) <= 100 { // Only reindex if reasonable number
			panicHandler := r.getPanicHandler()
			panicHandler.SafeGo(ctx, "post reindexing for category slug", func() {
				reindexCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				
				// Rate limit reindexing to prevent overwhelming Redis
				ticker := time.NewTicker(10 * time.Millisecond) // 100 ops/sec max
				defer ticker.Stop()
				
				for _, p := range fallbackPosts {
					select {
					case <-reindexCtx.Done():
						return
					case <-ticker.C:
						if err := r.addOrUpdateDoc(reindexCtx, client, p); err != nil {
							log.Printf("[WARNING] Failed to reindex post %s: %v", p.PostID, err)
						}
					}
				}
			})
		}

		return fallbackPosts, nil
	}

	// Otherwise parse all the found docs
	var results []*models.Post
	for _, doc := range docs {
		p, parseErr := r.parseDocToPost(doc)
		if parseErr != nil {
			log.Printf("[SearchPostsWithCategory] skipping docID=%s parse err: %v", doc.Id, parseErr)
			continue
		}
		// If entityType != PostType => fallback for a "correct" doc
		if p.EntityType != models.PostType {
			log.Printf("[SearchPostsWithCategory] docID=%s entityType=%s => fallback", doc.Id, p.EntityType)
			correctPost, fbErr := r.fallbackForWrongType(ctx, client, doc.Id)
			if fbErr != nil {
				log.Printf("[SearchPostsWithCategory] fallback error docID=%s: %v", doc.Id, fbErr)
				continue
			}
			if correctPost != nil {
				results = append(results, correctPost)
			}
			continue
		}
		results = append(results, p)
	}

	log.Printf("[SearchPostsWithCategory] returning %d docs, total=%d", len(results), total)
	return results, nil
}

// SearchWithFilters tries RediSearch. If no docs => fallback => reindex.
func (r *PostCacheRepository) SearchPostsWithFilters(
	ctx context.Context,
	name string,
	description string,
	tags []string,
	status string,
	userType string,
	thumbnail string,
	offset int64,
	limit int64,
	lat, lng float64,
	radius int64,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*models.Post, error) {
	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Create a context with timeout to prevent long-running queries
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Use QueryBuilder pattern for consistent query construction
	qb := NewQueryBuilder(models.PostType)

	// Add name filter if provided
	if name != "" {
		qb.WithNameFilter(name)
	}

	// Add description filter if provided
	if description != "" {
		qb.WithCustomFilter(fmt.Sprintf("@description:(%s)",
			redisearch.EscapeTextFileString(description)))
	}

	// Add thumbnail filter if provided
	if thumbnail != "" {
		qb.WithCustomFilter(fmt.Sprintf("@thumbnail:(%s)",
			redisearch.EscapeTextFileString(thumbnail)))
	}

	// Add status filter if provided
	if status != "" {
		qb.WithStatus(status)
	}

	// Add user type filter if provided
	if userType != "" {
		qb.WithCustomFilter(fmt.Sprintf("@user_type:{%s}",
			redisearch.EscapeTextFileString(userType)))
	}

	// Add tags filter if provided
	if len(tags) > 0 {
		tagParts := make([]string, len(tags))
		for i, t := range tags {
			tagParts[i] = redisearch.EscapeTextFileString(t)
		}
		qb.WithCustomFilter(fmt.Sprintf("@tags:{%s}", strings.Join(tagParts, "|")))
	}

	// Add geo filter if provided - this properly optimizes spatial queries
	if lat != 0 && lng != 0 && radius > 0 {
		qb.WithGeoFilter(lat, lng, radius)
	}

	// Configure pagination
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
	qb.WithPagination(int(finalOffset), int(finalLimit))

	// Set sorting if provided
	if sortBy != "" {
		sortDesc := strings.ToLower(sortOrder) == "desc"
		qb.WithSorting(sortBy, sortDesc)
	}

	// Set fields to return
	qb.WithReturnFields(
		"post_id", "name", "description", "user_id", "thumbnail", "status",
		"tags", "entity_type", "location",
	)

	// Build the final query
	_, query := qb.Build()

	// Execute search with timeout context
	docs, total, err := client.Search(query)
	if err != nil {
		if queryCtx.Err() == context.DeadlineExceeded {
			log.Printf("[SearchWithFilters] query timed out for name=%q lat=%.6f lng=%.6f radius=%d",
				name, lat, lng, radius)
			// Fall back to a simpler query without geo filtering
			if lat != 0 && lng != 0 && radius > 0 {
				return r.fallbackToSimpleQuery(ctx, name, description, tags, status, userType,
					thumbnail, offset, limit, page, pageSize, sortBy, sortOrder)
			}
		}
		return nil, errors.Wrap(err, "RediSearch query error in SearchWithFilters")
	}

	// If no docs found => fallback => reindex
	if len(docs) == 0 {
		log.Printf("[SearchWithFilters] No docs => fallback. name=%q lat=%.6f lng=%.6f radius=%d",
			name, lat, lng, radius)

		fallbackPosts, fallbackErr := r.fallback.SearchPostsWithFilters(
			ctx,
			name,
			description,
			tags,
			status,
			userType,
			thumbnail,
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
			return nil, fallbackErr
		}

		// Reindex asynchronously with proper context, rate limiting, and panic protection
		if len(fallbackPosts) > 0 && len(fallbackPosts) <= 100 { // Only reindex if reasonable number
			panicHandler := r.getPanicHandler()
			panicHandler.SafeGo(ctx, "post reindexing for filters", func() {
				reindexCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				
				// Rate limit reindexing to prevent overwhelming Redis
				ticker := time.NewTicker(10 * time.Millisecond) // 100 ops/sec max
				defer ticker.Stop()
				
				for _, p := range fallbackPosts {
					select {
					case <-reindexCtx.Done():
						return
					case <-ticker.C:
						if err := r.addOrUpdateDoc(reindexCtx, client, p); err != nil {
							log.Printf("[WARNING] Failed to reindex post %s: %v", p.PostID, err)
						}
					}
				}
			})
		}

		return fallbackPosts, nil
	}

	var results []*models.Post
	for _, doc := range docs {
		p, parseErr := r.parseDocToPost(doc)
		if parseErr != nil {
			log.Printf("[SearchWithFilters] skipping docID=%s parse err: %v", doc.Id, parseErr)
			continue
		}
		// If entityType is not PostType => fallback
		if p.EntityType != models.PostType {
			continue
		}
		results = append(results, p)
	}
	log.Printf("[SearchWithFilters] found %d docs, total=%d", len(results), total)
	return results, nil
}

// fallbackToSimpleQuery is a circuit breaker fallback when geo queries time out
func (r *PostCacheRepository) fallbackToSimpleQuery(
	ctx context.Context,
	name string,
	description string,
	tags []string,
	status string,
	userType string,
	thumbnail string,
	offset int64,
	limit int64,
	page int64,
	pageSize int64,
	sortBy string,
	sortOrder string,
) ([]*models.Post, error) {
	log.Println("[PostCacheRepository] Falling back to simple query without geo filtering")

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)

	// Create a simpler query without geo filtering
	qb := NewQueryBuilder(models.PostType)

	if name != "" {
		qb.WithNameFilter(name)
	}

	if description != "" {
		qb.WithCustomFilter(fmt.Sprintf("@description:(%s)",
			redisearch.EscapeTextFileString(description)))
	}

	if status != "" {
		qb.WithStatus(status)
	}

	// Configure pagination
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
	qb.WithPagination(int(finalOffset), int(finalLimit))

	// Set return fields
	qb.WithReturnFields(
		"post_id", "name", "description", "user_id", "thumbnail", "status",
		"tags", "entity_type", "location",
	)

	// Build and execute the simpler query
	_, query := qb.Build()
	docs, _, err := client.Search(query)
	if err != nil {
		return nil, errors.Wrap(err, "RediSearch simple fallback query error")
	}

	var results []*models.Post
	for _, doc := range docs {
		p, parseErr := r.parseDocToPost(doc)
		if parseErr != nil {
			continue
		}
		if p.EntityType != models.PostType {
			continue
		}
		results = append(results, p)
	}

	return results, nil
}

// addOrUpdateDoc re-indexes a post doc from the domain struct.
func (r *PostCacheRepository) addOrUpdateDoc(ctx context.Context, client redisearch.Client, p *models.Post) error {
	tagsJSON, _ := json.Marshal(p.Tags)
	locStr := fmt.Sprintf("%.6f,%.6f", p.Lng, p.Lat)

	doc := redisearch.NewDocument(p.PostID, 1.0).
		Set("post_id", p.PostID).
		Set("name", safeString(p.Name)).
		Set("description", safeString(p.Description)).
		Set("user_id", safeString(p.UserID)).
		Set("thumbnail", safeString(p.Thumbnail)).
		Set("status", safeString(p.Status)).
		Set("tags", string(tagsJSON)).
		Set("entity_type", models.PostType.String()). // Use proper entity type constant
		Set("location", locStr)

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

	doc.Set("created_at", createdAt.Unix())
	doc.Set("updated_at", updatedAt.Unix())

	// Use replace option to prevent "Document already exists" errors
	return client.IndexOptions(redisearch.IndexingOptions{Replace: true}, doc)
}

// parseDocToPost maps a RediSearch Document => *models.Post.
func (r *PostCacheRepository) parseDocToPost(doc redisearch.Document) (*models.Post, error) {
	p := &models.Post{PostID: doc.Id}

	p.Name = strVal(doc.Properties["name"])
	p.Description = strVal(doc.Properties["description"])
	p.UserID = strVal(doc.Properties["user_id"])
	p.Thumbnail = strVal(doc.Properties["thumbnail"])
	p.Status = strVal(doc.Properties["status"])

	// entity_type => domain model
	p.EntityType = models.ToEntityType(strVal(doc.Properties["entity_type"]))

	// tags => JSON array
	if rawT, ok := doc.Properties["tags"].(string); ok && rawT != "" {
		var t []string
		if e := json.Unmarshal([]byte(rawT), &t); e == nil {
			p.Tags = t
		}
	}

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

	// Parse timestamps => Unix timestamps
	if createdAtUnix, err := parseInt64(doc.Properties["created_at"], "created_at", p.PostID); err == nil && createdAtUnix > 0 {
		p.CreatedAt = time.Unix(createdAtUnix, 0)
	}
	if updatedAtUnix, err := parseInt64(doc.Properties["updated_at"], "updated_at", p.PostID); err == nil && updatedAtUnix > 0 {
		p.UpdatedAt = time.Unix(updatedAtUnix, 0)
	}

	return p, nil
}
func (r *PostCacheRepository) fallbackForWrongType(
	ctx context.Context,
	client redisearch.Client,
	docID string,
) (*models.Post, error) {
	// Remove the doc from RediSearch first
	if err := client.DeleteDocument(docID); err != nil {
		log.Printf("[fallbackForWrongType] could not delete docID=%s: %v", docID, err)
	}
	// fetch from fallback
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

func (r *PostCacheRepository) fetchFromFallbackAndMaybeReindex(
	ctx context.Context,
	client redisearch.Client,
	postID string,
) (*models.Post, error) {
	fbPost, fbErr := r.fallback.Find(ctx, postID)
	if fbErr != nil {
		return nil, errors.Wrap(fbErr, "fallback find error in fetchFromFallbackAndMaybeReindex")
	}
	if fbPost == nil {
		return nil, nil
	}
	
	// Reindex asynchronously with panic protection
	panicHandler := r.getPanicHandler()
	panicHandler.SafeGo(ctx, "post reindexing from fallback", func() {
		reindexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		if err := r.addOrUpdateDoc(reindexCtx, client, fbPost); err != nil {
			log.Printf("[WARNING] Failed to reindex post %s from fallback: %v", postID, err)
		}
	})
	
	return fbPost, nil
}

// FindBatch retrieves multiple posts by their IDs using parallel fetches for efficiency
func (r *PostCacheRepository) FindBatch(ctx context.Context, postIDs []string) (map[string]*models.Post, error) {
	if len(postIDs) == 0 {
		return make(map[string]*models.Post), nil
	}

	client := di.Get(ctx, constants.RedisearchClientKey).(redisearch.Client)
	result := make(map[string]*models.Post, len(postIDs))
	
	// First, try to get from Redis using individual document fetches
	// RediSearch doesn't have a native batch get, so we'll fetch in parallel
	type fetchResult struct {
		id   string
		post *models.Post
		err  error
	}
	
	ch := make(chan fetchResult, len(postIDs))
	
	// Use worker pool to limit concurrent fetches
	const maxWorkers = 10
	sem := make(chan struct{}, maxWorkers)
	
	panicHandler := r.getPanicHandler()
	
	for _, id := range postIDs {
		sem <- struct{}{} // Acquire semaphore
		panicHandler.SafeGo(ctx, "post batch fetch", func() {
			defer func() { <-sem }() // Release semaphore
			
			postID := id // Capture loop variable
			
			// Use circuit breaker for each fetch
			var post *models.Post
			err := r.circuitBreaker.Call(ctx, func() error {
				escapedPostID := redisearch.EscapeTextFileString(postID)
				q := redisearch.NewQuery(fmt.Sprintf("@entity_type:{%s} @post_id:{%s}", models.PostType.String(), escapedPostID)).
					SetReturnFields(
						"post_id", "name", "description", "user_id",
						"thumbnail", "status", "tags",
						"entity_type",
						"location", "created_at", "updated_at",
					).
					Limit(0, 1)

				docs, _, searchErr := client.Search(q)
				if searchErr != nil {
					return searchErr
				}
				
				if len(docs) == 0 {
					return errors.ErrNotFound.Msgf("post %s not found", postID)
				}
				
				var parseErr error
				post, parseErr = r.parseDocToPost(docs[0])
				return parseErr
			})
			
			ch <- fetchResult{id: postID, post: post, err: err}
		})
	}
	
	// Collect results and track missing IDs
	var missingIDs []string
	for i := 0; i < len(postIDs); i++ {
		res := <-ch
		if res.err != nil {
			missingIDs = append(missingIDs, res.id)
		} else if res.post != nil {
			result[res.id] = res.post
		}
	}
	
	// If any IDs are missing, try fallback
	if len(missingIDs) > 0 && r.fallback != nil {
		log.Printf("[FindBatch] %d posts not in cache, trying fallback", len(missingIDs))
		
		// Check if fallback supports batch fetch
		if batchFallback, ok := r.fallback.(application.PostBatchRepository); ok {
			fallbackPosts, err := batchFallback.FindBatch(ctx, missingIDs)
			if err != nil {
				log.Printf("[FindBatch] Fallback batch fetch failed: %v", err)
			} else {
				// Add fallback results to result map and reindex them
				if len(fallbackPosts) > 0 && len(fallbackPosts) <= 100 { // Only reindex reasonable number
					for id, post := range fallbackPosts {
						result[id] = post
						// Reindex asynchronously with rate limiting
						panicHandler.SafeGo(ctx, "post batch reindexing", func() {
							p := post // Capture loop variable
							reindexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
							defer cancel()
							
							if err := r.addOrUpdateDoc(reindexCtx, client, p); err != nil {
								log.Printf("[FindBatch] Failed to reindex post %s: %v", p.PostID, err)
							}
						})
					}
				}
			}
		} else {
			// Fallback doesn't support batch, fetch individually
			for _, id := range missingIDs {
				post, err := r.fallback.Find(ctx, id)
				if err == nil && post != nil {
					result[id] = post
					// Reindex asynchronously
					panicHandler.SafeGo(ctx, "post individual reindexing", func() {
						p := post // Capture variable
						reindexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()
						
						if err := r.addOrUpdateDoc(reindexCtx, client, p); err != nil {
							log.Printf("[FindBatch] Failed to reindex post %s: %v", p.PostID, err)
						}
					})
				}
			}
		}
	}
	
	return result, nil
}

// -----------------------------------------------------------------------------
// Simple type-conversion helpers
// -----------------------------------------------------------------------------
