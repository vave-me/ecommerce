package redis_repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/stackus/errors"

	"middleman/internal/di"
	"middleman/scheduler/internal/constants"
	"middleman/scheduler/internal/domain"
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
	schedulerID, userID string,
) error {
	// 1) Write to Redis
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	fmt.Printf("[Add] Received request: schedulerID=%s, userID=%s\n", schedulerID, userID)

	scheduler := &domain.MiddlemanScheduler{
		ID:     schedulerID,
		UserID: userID,
	}

	data, err := json.Marshal(scheduler)
	if err != nil {
		fmt.Printf("[Add] Error marshalling scheduler: %v\n", err)
		return errors.Wrap(err, "marshalling scheduler data")
	}

	key := fmt.Sprintf("%s:scheduler:%s", r.prefix, schedulerID)
	if _, err := conn.Do("SET", key, data); err != nil {
		fmt.Printf("[Add] Error storing scheduler in Redis under key=%s: %v\n", key, err)
		return errors.Wrap(err, "storing scheduler in Redis")
	}

	fmt.Printf("[Add] Successfully stored scheduler in Redis: key=%s\n", key)

	// 2) Synchronize write to Postgres (no goroutine):
	fmt.Printf("[Add-Sync] Inserting scheduler %s into Postgres\n", schedulerID)
	if err := r.fallback.Add(context.Background(), schedulerID, userID); err != nil {
		fmt.Printf("[Add-Sync] Failed to insert scheduler %s into Postgres: %v\n", schedulerID, err)
		return errors.Wrap(err, "inserting scheduler in Postgres")
	}

	fmt.Printf("[Add-Sync] Successfully inserted scheduler %s into Postgres\n", schedulerID)
	return nil
}

// Remove deletes from Redis, then immediately removes from Postgres (no goroutines).
func (r MiddlemanCacheRepository) Remove(
	ctx context.Context,
	schedulerID string,
) error {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:scheduler:%s", r.prefix, schedulerID)
	if _, err := conn.Do("DEL", key); err != nil {
		return errors.Wrap(err, "deleting scheduler in Redis")
	}

	fmt.Printf("[Remove-Sync] Removing scheduler %s from Postgres\n", schedulerID)
	if err := r.fallback.Remove(context.Background(), schedulerID); err != nil {
		fmt.Printf("[Remove-Sync] Failed to remove scheduler %s in Postgres: %v\n", schedulerID, err)
		return errors.Wrap(err, "removing scheduler in Postgres")
	}

	fmt.Printf("[Remove-Sync] Successfully removed scheduler %s in Postgres\n", schedulerID)
	return nil
}

// Find locates the first matching userID in Redis. If not found, fallback to Postgres.
func (r MiddlemanCacheRepository) Find(
	ctx context.Context,
	userID string,
) (*domain.MiddlemanScheduler, error) {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	// 1) Attempt to find in Redis
	pattern := fmt.Sprintf("%s:scheduler:*", r.prefix)
	keys, err := redis.Strings(conn.Do("KEYS", pattern))
	if err != nil {
		return nil, errors.Wrap(err, "searching activities in Redis")
	}
	for _, key := range keys {
		raw, gErr := redis.Bytes(conn.Do("GET", key))
		if gErr != nil {
			continue
		}
		scheduler := &domain.MiddlemanScheduler{}
		if err := json.Unmarshal(raw, scheduler); err != nil {
			continue
		}
		if scheduler.UserID == userID {
			// Found it in Redis
			return scheduler, nil
		}
	}

	// 2) Fallback to Postgres if not in Redis
	dbScheduler, err := r.fallback.Find(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	// Re-cache in Redis for next time
	recachedData, _ := json.Marshal(dbScheduler)
	newKey := fmt.Sprintf("%s:scheduler:%s", r.prefix, dbScheduler.ID)
	_, _ = conn.Do("SET", newKey, recachedData)
	fmt.Printf("[Find] Fallback found scheduler in Postgres, re-cached key=%s\n", newKey)
	return dbScheduler, nil
}

// All tries to get all from Redis for the user. If empty, fallback to Postgres.
func (r MiddlemanCacheRepository) All(
	ctx context.Context,
	userID string,
) ([]*domain.MiddlemanScheduler, error) {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	pattern := fmt.Sprintf("%s:scheduler:*", r.prefix)
	keys, err := redis.Strings(conn.Do("KEYS", pattern))
	if err != nil {
		return nil, errors.Wrap(err, "retrieving scheduler keys from Redis")
	}

	var activities []*domain.MiddlemanScheduler
	for _, key := range keys {
		raw, gErr := redis.Bytes(conn.Do("GET", key))
		if gErr != nil || len(raw) == 0 {
			continue
		}
		scheduler := &domain.MiddlemanScheduler{}
		if err := json.Unmarshal(raw, scheduler); err != nil {
			continue
		}
		if scheduler.UserID == userID {
			activities = append(activities, scheduler)
		}
	}

	if len(activities) > 0 {
		return activities, nil // Return from Redis
	}

	// Fallback to Postgres if Redis has none
	dbSchedulers, err := r.fallback.All(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	// Re-cache them
	for _, act := range dbSchedulers {
		recachedData, _ := json.Marshal(act)
		newKey := fmt.Sprintf("%s:scheduler:%s", r.prefix, act.ID)
		_, _ = conn.Do("SET", newKey, recachedData)
	}
	fmt.Printf("[All] Fallback found %d activities in Postgres, re-cached them in Redis\n", len(dbSchedulers))

	return dbSchedulers, nil
}
