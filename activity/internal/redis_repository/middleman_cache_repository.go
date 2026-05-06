package redis_repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/stackus/errors"

	"middleman/activity/internal/constants"
	"middleman/activity/internal/domain"
	"middleman/internal/di"
)

// MiddlemanCacheRepository caches data in Redis and writes to Postgres synchronously.
// It also falls back to Postgres on reads if Redis doesn’t have the data.
type MiddlemanCacheRepository struct {
	prefix   string
	fallback domain.MiddlemanRepository // Postgres fallback (source of truth)
}

// Compile-time check
var _ domain.MiddlemanCacheRepository = (*MiddlemanCacheRepository)(nil)

// NewMiddlemanCacheRepository constructs a new "write-through + read-through fallback" repository.
func NewMiddlemanCacheRepository(prefix string, fallback domain.MiddlemanRepository) MiddlemanCacheRepository {
	return MiddlemanCacheRepository{
		prefix:   prefix,
		fallback: fallback,
	}
}

// Add writes data to Redis, then immediately inserts into Postgres (no goroutines).
func (r MiddlemanCacheRepository) Add(
	ctx context.Context,
	activityID, userID string,
) error {
	// 1) Write to Redis
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	fmt.Printf("[Add] Received request: activityID=%s, userID=%s\n", activityID, userID)

	activity := &domain.MiddlemanActivity{
		ID:     activityID,
		UserID: userID,
	}

	data, err := json.Marshal(activity)
	if err != nil {
		fmt.Printf("[Add] Error marshalling activity: %v\n", err)
		return errors.Wrap(err, "marshalling activity data")
	}

	key := fmt.Sprintf("%s:activity:%s", r.prefix, activityID)
	if _, err := conn.Do("SET", key, data); err != nil {
		fmt.Printf("[Add] Error storing activity in Redis under key=%s: %v\n", key, err)
		return errors.Wrap(err, "storing activity in Redis")
	}

	fmt.Printf("[Add] Successfully stored activity in Redis: key=%s\n", key)

	// 2) Synchronize write to Postgres (no goroutine):
	fmt.Printf("[Add-Sync] Inserting activity %s into Postgres\n", activityID)
	if err := r.fallback.Add(context.Background(), activityID, userID); err != nil {
		fmt.Printf("[Add-Sync] Failed to insert activity %s into Postgres: %v\n", activityID, err)
		return errors.Wrap(err, "inserting activity in Postgres")
	}

	fmt.Printf("[Add-Sync] Successfully inserted activity %s into Postgres\n", activityID)
	return nil
}

// Remove deletes from Redis, then immediately removes from Postgres (no goroutines).
func (r MiddlemanCacheRepository) Remove(
	ctx context.Context,
	activityID string,
) error {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:activity:%s", r.prefix, activityID)
	if _, err := conn.Do("DEL", key); err != nil {
		return errors.Wrap(err, "deleting activity in Redis")
	}

	fmt.Printf("[Remove-Sync] Removing activity %s from Postgres\n", activityID)
	if err := r.fallback.Remove(context.Background(), activityID); err != nil {
		fmt.Printf("[Remove-Sync] Failed to remove activity %s in Postgres: %v\n", activityID, err)
		return errors.Wrap(err, "removing activity in Postgres")
	}

	fmt.Printf("[Remove-Sync] Successfully removed activity %s in Postgres\n", activityID)
	return nil
}

// Find locates the first matching userID in Redis. If not found, fallback to Postgres.
func (r MiddlemanCacheRepository) Find(
	ctx context.Context,
	userID string,
) (*domain.MiddlemanActivity, error) {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	// 1) Attempt to find in Redis
	pattern := fmt.Sprintf("%s:activity:*", r.prefix)
	keys, err := redis.Strings(conn.Do("KEYS", pattern))
	if err != nil {
		return nil, errors.Wrap(err, "searching activities in Redis")
	}
	for _, key := range keys {
		raw, gErr := redis.Bytes(conn.Do("GET", key))
		if gErr != nil {
			continue
		}
		activity := &domain.MiddlemanActivity{}
		if err := json.Unmarshal(raw, activity); err != nil {
			continue
		}
		if activity.UserID == userID {
			// Found it in Redis
			return activity, nil
		}
	}

	// 2) Fallback to Postgres if not in Redis
	dbActivity, err := r.fallback.Find(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	// Re-cache in Redis for next time
	recachedData, _ := json.Marshal(dbActivity)
	newKey := fmt.Sprintf("%s:activity:%s", r.prefix, dbActivity.ID)
	_, _ = conn.Do("SET", newKey, recachedData)
	fmt.Printf("[Find] Fallback found activity in Postgres, re-cached key=%s\n", newKey)
	return dbActivity, nil
}

// All tries to get all from Redis for the user. If empty, fallback to Postgres.
func (r MiddlemanCacheRepository) All(
	ctx context.Context,
	userID string,
) ([]*domain.MiddlemanActivity, error) {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	pattern := fmt.Sprintf("%s:activity:*", r.prefix)
	keys, err := redis.Strings(conn.Do("KEYS", pattern))
	if err != nil {
		return nil, errors.Wrap(err, "retrieving activity keys from Redis")
	}

	var activities []*domain.MiddlemanActivity
	for _, key := range keys {
		raw, gErr := redis.Bytes(conn.Do("GET", key))
		if gErr != nil || len(raw) == 0 {
			continue
		}
		activity := &domain.MiddlemanActivity{}
		if err := json.Unmarshal(raw, activity); err != nil {
			continue
		}
		if activity.UserID == userID {
			activities = append(activities, activity)
		}
	}

	if len(activities) > 0 {
		return activities, nil // Return from Redis
	}

	// Fallback to Postgres if Redis has none
	dbActivities, err := r.fallback.All(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	// Re-cache them
	for _, act := range dbActivities {
		recachedData, _ := json.Marshal(act)
		newKey := fmt.Sprintf("%s:activity:%s", r.prefix, act.ID)
		_, _ = conn.Do("SET", newKey, recachedData)
	}
	fmt.Printf("[All] Fallback found %d activities in Postgres, re-cached them in Redis\n", len(dbActivities))

	return dbActivities, nil
}
