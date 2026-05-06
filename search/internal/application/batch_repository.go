package application

import (
	"context"
	"middleman/search/internal/models"
)

// BatchFetchRepository defines batch fetch operations for repositories
type BatchFetchRepository interface {
	// FindBatch fetches multiple entities by their IDs
	// Returns a map of ID -> entity for efficient lookup
	// Missing IDs will not be present in the map
	FindBatch(ctx context.Context, ids []string) (map[string]interface{}, error)
}

// ProductBatchRepository extends ProductRepository with batch operations
type ProductBatchRepository interface {
	ProductRepository
	// FindBatch fetches multiple products by their IDs
	FindBatch(ctx context.Context, productIDs []string) (map[string]*models.Product, error)
}

// PostBatchRepository extends PostRepository with batch operations
type PostBatchRepository interface {
	PostRepository
	// FindBatch fetches multiple posts by their IDs
	FindBatch(ctx context.Context, postIDs []string) (map[string]*models.Post, error)
}

// UserBatchRepository extends UserRepository with batch operations
type UserBatchRepository interface {
	UserRepository
	// FindBatch fetches multiple users by their IDs
	FindBatch(ctx context.Context, userIDs []string) (map[string]*models.User, error)
}

// OrderBatchRepository extends OrderRepository with batch operations
type OrderBatchRepository interface {
	OrderRepository
	// FindBatch fetches multiple orders by their IDs
	FindBatch(ctx context.Context, orderIDs []string) (map[string]*models.Order, error)
}

// ServiceBatchRepository extends ServiceRepository with batch operations
type ServiceBatchRepository interface {
	ServiceRepository
	// FindBatch fetches multiple services by their IDs
	FindBatch(ctx context.Context, serviceIDs []string) (map[string]*models.Service, error)
}

// MetricBatchRepository already has batch operations via GetMultiple
// No need to extend it further

// BatchHelper provides utility functions for batch operations
type BatchHelper struct{}

// ChunkIDs splits a slice of IDs into chunks of specified size
func (BatchHelper) ChunkIDs(ids []string, chunkSize int) [][]string {
	if chunkSize <= 0 {
		chunkSize = 100 // Default chunk size
	}
	
	var chunks [][]string
	for i := 0; i < len(ids); i += chunkSize {
		end := i + chunkSize
		if end > len(ids) {
			end = len(ids)
		}
		chunks = append(chunks, ids[i:end])
	}
	return chunks
}

// MergeMaps merges multiple maps into one (string keys, interface{} values)
func (BatchHelper) MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// ExtractUniqueIDs extracts unique IDs from a slice
func (BatchHelper) ExtractUniqueIDs(ids []string) []string {
	seen := make(map[string]bool)
	unique := make([]string, 0, len(ids))
	
	for _, id := range ids {
		if id != "" && !seen[id] {
			seen[id] = true
			unique = append(unique, id)
		}
	}
	return unique
}