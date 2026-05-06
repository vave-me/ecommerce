// File: metric/internal/redis/item_metric_cache_repository.go
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stackus/errors"

	"middleman/internal/di"
	"middleman/metrics/internal/application"
	"middleman/metrics/internal/constants"
	"middleman/metrics/internal/models"
)

// ItemMetricRepository implements application.ItemMetricRepository using Redis
type ItemMetricRepository struct {
	tableName string
	fallback  application.ItemMetricCacheRepository
}

var _ application.ItemMetricRepository = (*ItemMetricRepository)(nil)

// NewItemMetricRepository creates a new Redis-based repository for item metrics
func NewItemMetricRepository(
	tableName string,
	fallback application.ItemMetricCacheRepository,
) *ItemMetricRepository {
	if tableName == "" {
		tableName = "metrics:items:"
	}
	return &ItemMetricRepository{
		tableName: tableName,
		fallback:  fallback,
	}
}

// GetItemsMetrics retrieves multiple metrics by IDs, with a limit of max 150 items
func (r *ItemMetricRepository) GetItemsMetrics(ctx context.Context, itemIds []string) ([]*models.ItemMetric, error) {
	if len(itemIds) == 0 {
		return []*models.ItemMetric{}, nil
	}

	// Limit to maximum 150 items
	if len(itemIds) > 150 {
		itemIds = itemIds[:150]
	}

	redisPool := di.Get(ctx, constants.RedisTransactionKey).(redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// Track which items are found in Redis
	foundItems := make(map[string]*models.ItemMetric)
	var missingIds []string

	// Pipeline Redis requests for efficiency
	for _, itemId := range itemIds {
		key := r.getKey(itemId)
		db.Send("GET", key)
	}
	db.Flush()

	// Process responses
	for _, itemId := range itemIds {
		data, err := redis.Bytes(db.Receive())
		if err == nil {
			var metric models.ItemMetric
			if err := json.Unmarshal(data, &metric); err != nil {
				// On unmarshal error, mark as missing to get from fallback
				missingIds = append(missingIds, itemId)
				continue
			}

			// Refresh TTL in a fire-and-forget manner (don't wait for response)
			key := r.getKey(itemId)
			db.Send("EXPIRE", key, 86400) // 24 hours

			foundItems[itemId] = &metric
		} else if err == redis.ErrNil {
			// Not found in Redis, will get from fallback
			missingIds = append(missingIds, itemId)
		} else {
			// Unexpected error, try fallback for all remaining items
			if r.fallback != nil {
				fallbackMetrics, err := r.fallback.GetItemsMetrics(ctx, itemIds)
				if err != nil {
					return nil, errors.Wrap(err, "getting item metrics from fallback")
				}
				return fallbackMetrics, nil
			}
			return nil, errors.Wrap(err, "querying Redis for item metrics")
		}
	}
	db.Flush() // Flush the EXPIRE commands

	// If we have missing items and a fallback repository, get them from there
	if len(missingIds) > 0 && r.fallback != nil {
		fallbackMetrics, err := r.fallback.GetItemsMetrics(ctx, missingIds)
		if err != nil {
			return nil, errors.Wrap(err, "getting missing item metrics from fallback")
		}

		// Store fallback results in Redis for next time (in parallel)
		for _, metric := range fallbackMetrics {
			data, err := json.Marshal(metric)
			if err != nil {
				continue // Skip caching this one
			}

			key := r.getKey(metric.ID)
			db.Send("SET", key, data, "EX", 86400)

			foundItems[metric.ID] = metric
		}
		db.Flush() // Flush the SET commands
	}

	// Convert map to slice in the same order as requested
	result := make([]*models.ItemMetric, 0, len(foundItems))
	for _, itemId := range itemIds {
		if metric, ok := foundItems[itemId]; ok {
			result = append(result, metric)
		}
	}

	return result, nil
}

// AddMetric adds a new item metric to the cache
func (r *ItemMetricRepository) AddMetric(ctx context.Context, itemId, entityType, categoryId string, price int64, lat, lng float64) error {
	redisPool := di.Get(ctx, constants.RedisTransactionKey).(redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// First add to fallback (persistent storage)
	if r.fallback != nil {
		if err := r.fallback.AddMetric(ctx, itemId, entityType, categoryId, price, lat, lng); err != nil {
			return errors.Wrap(err, "adding item metric to fallback repository")
		}
	}

	// Create a new metric with default values
	metric := &models.ItemMetric{
		ID:                   itemId,
		EntityType:           entityType,
		LikesCount:           0,
		DislikesCount:        0,
		CommentsCount:        0,
		MessagesCount:        0,
		SharedCount:          0,
		AddedToWishlistCount: 0,
		AddedToBasketCount:   0,
		VisitedCount:         0,
		ReportedCount:        0,
		FollowerCount:        0,
		ReviewsCount:         0,
		RatingCount:          0,
		VideosCount:          0,
		ImagesCount:          0,
		Rating:               0,
		Review:               0,
		Category:             "",
		CategoryID:           categoryId,
		CategorySlug:         "",
		Price:                price,
		Lat:                  lat,
		Lng:                  lng,
		CreatedAt:            time.Now().Format(time.RFC3339),
		UpdatedAt:            time.Now().Format(time.RFC3339),
	}

	// Store in Redis
	data, err := json.Marshal(metric)
	if err != nil {
		return errors.Wrap(err, "marshaling item metric")
	}

	key := r.getKey(itemId)
	_, err = db.Do("SET", key, data, "EX", 86400) // 24 hours TTL
	if err != nil {
		return errors.Wrap(err, "storing item metric in Redis")
	}

	return nil
}

// GetItemMetric retrieves a metric by ID from the cache, falling back to the repository if not found
func (r *ItemMetricRepository) GetItemMetric(ctx context.Context, itemId string) (*models.ItemMetric, error) {
	redisPool := di.Get(ctx, constants.RedisTransactionKey).(redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	key := r.getKey(itemId)
	data, err := redis.Bytes(db.Do("GET", key))

	if err == nil {
		// Found in Redis, unmarshal and return
		var metric models.ItemMetric
		if err := json.Unmarshal(data, &metric); err != nil {
			// Cache corrupted, get from fallback
			return r.getFromFallback(ctx, itemId)
		}

		// Refresh TTL
		db.Do("EXPIRE", key, 86400) // 24 hours
		return &metric, nil
	} else if err != redis.ErrNil {
		// Unexpected error (not just "not found")
		if r.fallback != nil {
			return r.getFromFallback(ctx, itemId)
		}
		return nil, errors.Wrap(err, "querying Redis for item metric")
	}

	// Not found in Redis, try fallback
	return r.getFromFallback(ctx, itemId)
}

// RemoveItemMetric removes a metric from the cache
func (r *ItemMetricRepository) RemoveItemMetric(ctx context.Context, itemId string) error {
	redisPool := di.Get(ctx, constants.RedisTransactionKey).(redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// Remove from fallback first
	if r.fallback != nil {
		if err := r.fallback.RemoveItemMetric(ctx, itemId); err != nil {
			return errors.Wrap(err, "removing item metric from fallback")
		}
	}

	// Remove from Redis
	key := r.getKey(itemId)
	_, err := db.Do("DEL", key)
	if err != nil {
		return errors.Wrap(err, "removing item metric from Redis")
	}

	return nil
}

// UpdateItemMetric updates a specific metric count
func (r *ItemMetricRepository) UpdateItemMetric(ctx context.Context, itemId, metricType, metricTypeAction string) error {
	redisPool := di.Get(ctx, constants.RedisTransactionKey).(redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// Get current metric data
	key := r.getKey(itemId)
	data, err := redis.Bytes(db.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			// Metric doesn't exist, create it first
			if err := r.AddMetric(ctx, itemId, "", "", 0, 0.0, 0.0); err != nil {
				return errors.Wrap(err, "adding missing item metric")
			}
			// Try to get it again
			data, err = redis.Bytes(db.Do("GET", key))
			if err != nil {
				return errors.Wrap(err, "getting newly added item metric")
			}
		} else {
			// If Redis fails, try fallback
			if r.fallback != nil {
				return r.fallback.UpdateItemMetric(ctx, itemId, metricType, metricTypeAction)
			}
			return errors.Wrap(err, "getting item metric for update")
		}
	}

	// Unmarshal the metric
	var metric models.ItemMetric
	if err := json.Unmarshal(data, &metric); err != nil {
		return errors.Wrap(err, "unmarshaling item metric")
	}

	// Update the appropriate field based on metric type
	var updated bool
	isIncrement := metricTypeAction == models.MetricTypeActionAdd.String()

	switch metricType {
	case models.MetricTypeCountLike.String():
		if isIncrement {
			metric.LikesCount++
		} else if metric.LikesCount > 0 {
			metric.LikesCount--
		}
		updated = true
	case models.MetricTypeCountDislike.String():
		if isIncrement {
			metric.DislikesCount++
		} else if metric.DislikesCount > 0 {
			metric.DislikesCount--
		}
		updated = true
	case models.MetricTypeCountComment.String():
		if isIncrement {
			metric.CommentsCount++
		} else if metric.CommentsCount > 0 {
			metric.CommentsCount--
		}
		updated = true
	case models.MetricTypeCountMessage.String():
		if isIncrement {
			metric.MessagesCount++
		} else if metric.MessagesCount > 0 {
			metric.MessagesCount--
		}
		updated = true
	case models.MetricTypeCountShare.String():
		if isIncrement {
			metric.SharedCount++
		} else if metric.SharedCount > 0 {
			metric.SharedCount--
		}
		updated = true
	case models.MetricTypeCountWishlist.String():
		if isIncrement {
			metric.AddedToWishlistCount++
		} else if metric.AddedToWishlistCount > 0 {
			metric.AddedToWishlistCount--
		}
		updated = true
	case models.MetricTypeCountBasket.String():
		if isIncrement {
			metric.AddedToBasketCount++
		} else if metric.AddedToBasketCount > 0 {
			metric.AddedToBasketCount--
		}
		updated = true
	case models.MetricTypeCountVisit.String():
		if isIncrement {
			metric.VisitedCount++
		} else if metric.VisitedCount > 0 {
			metric.VisitedCount--
		}
		updated = true
	case models.MetricTypeCountReport.String():
		if isIncrement {
			metric.ReportedCount++
		} else if metric.ReportedCount > 0 {
			metric.ReportedCount--
		}
		updated = true
	case models.MetricTypeCountFollow.String():
		if isIncrement {
			metric.FollowerCount++
		} else if metric.FollowerCount > 0 {
			metric.FollowerCount--
		}
		updated = true
	case models.MetricTypeCountReview.String():
		if isIncrement {
			metric.ReviewsCount++
		} else if metric.ReviewsCount > 0 {
			metric.ReviewsCount--
		}
		updated = true
	}

	if !updated {
		return nil // No field was updated, probably an unknown metric type
	}

	// Update the timestamp
	metric.UpdatedAt = time.Now().Format(time.RFC3339)

	// Save to Redis
	data, err = json.Marshal(metric)
	if err != nil {
		return errors.Wrap(err, "marshaling updated item metric")
	}

	_, err = db.Do("SET", key, data, "EX", 86400) // 24 hours TTL
	if err != nil {
		return errors.Wrap(err, "storing updated item metric in Redis")
	}

	// Update in fallback repository
	if r.fallback != nil {
		if err := r.fallback.UpdateItemMetric(ctx, itemId, metricType, metricTypeAction); err != nil {
			// Log error but don't fail - Redis is now the source of truth
			log.Printf("Error updating item metric in fallback: %v\n", err)
		}
	}

	return nil
}

// GetHighestMetricsByType delegates to fallback repository for geospatial queries
// Redis is not optimal for complex geospatial queries, so we use PostgreSQL directly
func (r *ItemMetricRepository) GetHighestMetricsByType(ctx context.Context, metricType string, entityTypes []models.EntityType, categoryId string, lat, lng, radius float64, minPrice, maxPrice int64, createdFrom, createdTill string) ([]*models.ItemMetric, error) {
	if r.fallback == nil {
		return []*models.ItemMetric{}, errors.ErrNotFound.Msg("fallback repository required for geospatial queries")
	}

	// For geospatial queries with complex filtering, use PostgreSQL directly
	return r.fallback.GetHighestMetricsByType(ctx, metricType, entityTypes, categoryId, lat, lng, radius, minPrice, maxPrice, createdFrom, createdTill)
}

// GetLowestMetricsByType delegates to fallback repository for geospatial queries
// Redis is not optimal for complex geospatial queries, so we use PostgreSQL directly
func (r *ItemMetricRepository) GetLowestMetricsByType(ctx context.Context, metricType string, entityTypes []models.EntityType, categoryId string, lat, lng, radius float64, minPrice, maxPrice int64, createdFrom, createdTill string) ([]*models.ItemMetric, error) {
	if r.fallback == nil {
		return []*models.ItemMetric{}, errors.ErrNotFound.Msg("fallback repository required for geospatial queries")
	}

	// For geospatial queries with complex filtering, use PostgreSQL directly
	return r.fallback.GetLowestMetricsByType(ctx, metricType, entityTypes, categoryId, lat, lng, radius, minPrice, maxPrice, createdFrom, createdTill)
}

// Helper methods
func (r *ItemMetricRepository) getKey(itemId string) string {
	return fmt.Sprintf("%s:%s", r.tableName, itemId)
}

func (r *ItemMetricRepository) getFromFallback(ctx context.Context, itemId string) (*models.ItemMetric, error) {
	if r.fallback == nil {
		return nil, errors.ErrNotFound
	}

	metric, err := r.fallback.GetItemMetric(ctx, itemId)
	if err != nil {
		return nil, errors.Wrap(err, "getting item metric from fallback")
	}

	// Store in Redis for next time
	redisPool := di.Get(ctx, constants.RedisTransactionKey).(redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	data, err := json.Marshal(metric)
	if err != nil {
		return metric, nil // Return the metric anyway, just don't cache it
	}

	key := r.getKey(itemId)
	_, err = db.Do("SET", key, data, "EX", 86400) // 24 hours TTL
	if err != nil {
		return metric, nil // Return the metric anyway, just don't cache it
	}

	return metric, nil
}
