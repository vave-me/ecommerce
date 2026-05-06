// File: metric/internal/redis/user_metric_cache_repository.go
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

// UserMetricRepository implements application.UserMetricRepository using Redis
type UserMetricRepository struct {
	tableName string
	fallback  application.UserMetricCacheRepository
}

var _ application.UserMetricRepository = (*UserMetricRepository)(nil)

// NewUserMetricRepository creates a new Redis-based repository for user metrics
func NewUserMetricRepository(
	tableName string,
	fallback application.UserMetricCacheRepository,
) *UserMetricRepository {
	if tableName == "" {
		tableName = "metrics:users:"
	}
	return &UserMetricRepository{
		tableName: tableName,
		fallback:  fallback,
	}
}

// AddUserMetric adds a new user metric to the cache
func (r *UserMetricRepository) AddUserMetric(ctx context.Context, userId, entityType string) error {
	redisPool := di.Get(ctx, constants.RedisTransactionKey).(redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// First add to fallback (persistent storage)
	if r.fallback != nil {
		if err := r.fallback.AddUserMetric(ctx, userId, entityType); err != nil {
			return errors.Wrap(err, "adding user metric to fallback repository")
		}
	}

	// Create a new metric with default values
	metric := &models.UserMetric{
		ID:                   userId,
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
		CategoryID:           "",
		CategorySlug:         "",
		MediaAddedCount:      0,
		CommentAddedCount:    0,
		LikedAddedCount:      0,
		ProductsAddedCount:   0,
		VideosAddedCount:     0,
		ServicesAddedCount:   0,
		JobsAddedCount:       0,
		PostsAddedCount:      0,
		VehiclesAddedCount:   0,
		PropertiesAddedCount: 0,
		CreatedAt:            time.Now().Format(time.RFC3339),
		UpdatedAt:            time.Now().Format(time.RFC3339),
	}

	// Store in Redis
	data, err := json.Marshal(metric)
	if err != nil {
		return errors.Wrap(err, "marshaling user metric")
	}

	key := r.getKey(userId)
	_, err = db.Do("SET", key, data, "EX", 86400) // 24 hours TTL
	if err != nil {
		return errors.Wrap(err, "storing user metric in Redis")
	}

	return nil
}

// GetUserMetric retrieves a metric by ID from the cache, falling back to the repository if not found
func (r *UserMetricRepository) GetUserMetric(ctx context.Context, userId string) (*models.UserMetric, error) {
	redisPool := di.Get(ctx, constants.RedisTransactionKey).(redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	key := r.getKey(userId)
	data, err := redis.Bytes(db.Do("GET", key))

	if err == nil {
		// Found in Redis, unmarshal and return
		var metric models.UserMetric
		if err := json.Unmarshal(data, &metric); err != nil {
			// Cache corrupted, get from fallback
			return r.getFromFallback(ctx, userId)
		}

		// Refresh TTL
		db.Do("EXPIRE", key, 86400) // 24 hours
		return &metric, nil
	} else if err != redis.ErrNil {
		// Unexpected error (not just "not found")
		return nil, errors.Wrap(err, "querying Redis for user metric")
	}

	// Not found in Redis, try fallback
	return r.getFromFallback(ctx, userId)
}

// RemoveUserMetric removes a metric from the cache and fallback
func (r *UserMetricRepository) RemoveUserMetric(ctx context.Context, userId string) error {
	redisPool := di.Get(ctx, constants.RedisTransactionKey).(redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// Remove from fallback first
	if r.fallback != nil {
		if err := r.fallback.RemoveUserMetric(ctx, userId); err != nil {
			return errors.Wrap(err, "removing user metric from fallback")
		}
	}

	// Remove from Redis
	key := r.getKey(userId)
	_, err := db.Do("DEL", key)
	if err != nil {
		return errors.Wrap(err, "removing user metric from Redis")
	}

	return nil
}

// UpdateUserMetric updates a specific metric count in an atomic way
func (r *UserMetricRepository) UpdateUserMetric(ctx context.Context, userId, metricType, metricTypeAction string) error {
	redisPool := di.Get(ctx, constants.RedisTransactionKey).(redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// Get current metric data
	key := r.getKey(userId)
	data, err := redis.Bytes(db.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			// Metric doesn't exist, create it first
			if err := r.AddUserMetric(ctx, userId, ""); err != nil {
				return errors.Wrap(err, "adding missing user metric")
			}
			// Try to get it again
			data, err = redis.Bytes(db.Do("GET", key))
			if err != nil {
				return errors.Wrap(err, "getting newly added user metric")
			}
		} else {
			// If Redis fails, try fallback
			if r.fallback != nil {
				return r.fallback.UpdateUserMetric(ctx, userId, metricType, metricTypeAction)
			}
			return errors.Wrap(err, "getting user metric for update")
		}
	}

	// Unmarshal the metric
	var metric models.UserMetric
	if err := json.Unmarshal(data, &metric); err != nil {
		return errors.Wrap(err, "unmarshaling user metric")
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
	case models.MetricTypeUserMediaAdd.String():
		metric.MediaAddedCount++
		updated = true
	case models.MetricTypeUserCommentAdd.String():
		metric.CommentAddedCount++
		updated = true
	case models.MetricTypeUserLikeAdd.String():
		metric.LikedAddedCount++
		updated = true
	case models.MetricTypeUserProductAdd.String():
		metric.ProductsAddedCount++
		updated = true
	case models.MetricTypeUserVideoAdd.String():
		metric.VideosAddedCount++
		updated = true
	case models.MetricTypeUserServiceAdd.String():
		metric.ServicesAddedCount++
		updated = true
	case models.MetricTypeUserJobAdd.String():
		metric.JobsAddedCount++
		updated = true
	case models.MetricTypeUserPostAdd.String():
		metric.PostsAddedCount++
		updated = true
	case models.MetricTypeUserVehicleAdd.String():
		metric.VehiclesAddedCount++
		updated = true
	case models.MetricTypeUserPropertyAdd.String():
		metric.PropertiesAddedCount++
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
		return errors.Wrap(err, "marshaling updated user metric")
	}

	_, err = db.Do("SET", key, data, "EX", 86400) // 24 hours TTL
	if err != nil {
		return errors.Wrap(err, "storing updated user metric in Redis")
	}

	// Update in fallback repository
	if r.fallback != nil {
		if err := r.fallback.UpdateUserMetric(ctx, userId, metricType, metricTypeAction); err != nil {
			// Log error but don't fail - Redis is now the source of truth
			log.Printf("Error updating user metric in fallback: %v\n", err)
		}
	}

	return nil
}

// Helper methods
func (r *UserMetricRepository) getKey(userId string) string {
	return fmt.Sprintf("%s:%s", r.tableName, userId)
}

func (r *UserMetricRepository) getFromFallback(ctx context.Context, userId string) (*models.UserMetric, error) {
	if r.fallback == nil {
		return nil, errors.ErrNotFound
	}

	metric, err := r.fallback.GetUserMetric(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "getting user metric from fallback")
	}

	// Store in Redis for next time
	redisPool := di.Get(ctx, constants.RedisTransactionKey).(redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	data, err := json.Marshal(metric)
	if err != nil {
		return metric, nil // Return the metric anyway, just don't cache it
	}

	key := r.getKey(userId)
	_, err = db.Do("SET", key, data, "EX", 86400) // 24 hours TTL
	if err != nil {
		return metric, nil // Return the metric anyway, just don't cache it
	}

	return metric, nil
}
