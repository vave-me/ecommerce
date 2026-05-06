package redis_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stackus/errors"

	"middleman/internal/di"
	"middleman/scheduler/internal/constants"
	"middleman/scheduler/internal/domain"
)

type MiddlemanCacheActionRepository struct {
	prefix   string
	fallback domain.MiddlemanActionRepository
}

var _ domain.MiddlemanCacheActionRepository = (*MiddlemanCacheActionRepository)(nil)

func NewMiddlemanCacheActionRepository(
	prefix string,
	fallback domain.MiddlemanActionRepository,
) MiddlemanCacheActionRepository {
	return MiddlemanCacheActionRepository{
		prefix:   prefix,
		fallback: fallback,
	}
}

func (r MiddlemanCacheActionRepository) Add(
	ctx context.Context,
	actionID, schedulerID, task string,
	executionTime time.Time,
) error {
	// Write to Redis
	pool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := pool.Get()
	defer conn.Close()

	action := domain.MiddlemanAction{
		ID:                  actionID,
		SchedulerID:         schedulerID,
		NaturalLanguageTask: task,
		ExecutionTime:       executionTime,
		Status:              "pending",
		CreatedAt:           time.Now(),
	}

	raw, err := json.Marshal(action)
	if err != nil {
		return errors.Wrap(err, "marshalling action")
	}

	key := fmt.Sprintf("%s:action:%s", r.prefix, actionID)
	if _, err := conn.Do("SET", key, raw); err != nil {
		return errors.Wrap(err, "storing action in Redis")
	}

	// Add to pending actions sorted set for efficient lookup by execution time
	pendingKey := fmt.Sprintf("%s:pending_actions", r.prefix)
	score := float64(executionTime.Unix())
	if _, err := conn.Do("ZADD", pendingKey, score, actionID); err != nil {
		return errors.Wrap(err, "adding to pending actions set")
	}

	// Write through to PostgreSQL
	if err := r.fallback.Add(ctx, actionID, schedulerID, task, executionTime); err != nil {
		return errors.Wrap(err, "inserting action into PostgreSQL")
	}

	return nil
}

func (r MiddlemanCacheActionRepository) UpdateStatus(
	ctx context.Context,
	actionID, status, result, errorMessage string,
) error {
	pool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := pool.Get()
	defer conn.Close()

	// Get existing action
	action, err := r.Find(ctx, actionID)
	if err != nil {
		return err
	}

	// Update fields
	action.Status = status
	now := time.Now()
	action.ExecutedAt = &now
	action.Result = result
	action.ErrorMessage = errorMessage

	// Save updated action
	raw, err := json.Marshal(action)
	if err != nil {
		return errors.Wrap(err, "marshalling updated action")
	}

	key := fmt.Sprintf("%s:action:%s", r.prefix, actionID)
	if _, err := conn.Do("SET", key, raw); err != nil {
		return errors.Wrap(err, "storing updated action in Redis")
	}

	// Remove from pending set if no longer pending
	if status != "pending" {
		pendingKey := fmt.Sprintf("%s:pending_actions", r.prefix)
		if _, err := conn.Do("ZREM", pendingKey, actionID); err != nil {
			return errors.Wrap(err, "removing from pending actions set")
		}
	}

	// Write through to PostgreSQL
	if err := r.fallback.UpdateStatus(ctx, actionID, status, result, errorMessage); err != nil {
		return errors.Wrap(err, "updating action in PostgreSQL")
	}

	return nil
}

func (r MiddlemanCacheActionRepository) Remove(
	ctx context.Context,
	actionID string,
) error {
	pool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := pool.Get()
	defer conn.Close()

	// Remove from Redis
	key := fmt.Sprintf("%s:action:%s", r.prefix, actionID)
	if _, err := conn.Do("DEL", key); err != nil {
		return errors.Wrap(err, "deleting action from Redis")
	}

	// Remove from pending set
	pendingKey := fmt.Sprintf("%s:pending_actions", r.prefix)
	if _, err := conn.Do("ZREM", pendingKey, actionID); err != nil {
		return errors.Wrap(err, "removing from pending actions set")
	}

	// Remove from PostgreSQL
	if err := r.fallback.Remove(ctx, actionID); err != nil {
		return errors.Wrap(err, "removing action from PostgreSQL")
	}

	return nil
}

func (r MiddlemanCacheActionRepository) Find(
	ctx context.Context,
	actionID string,
) (*domain.MiddlemanAction, error) {
	pool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:action:%s", r.prefix, actionID)
	raw, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		// Fallback to PostgreSQL
		action, fbErr := r.fallback.Find(ctx, actionID)
		if fbErr != nil {
			return nil, fbErr
		}
		// Re-cache
		recached, _ := json.Marshal(action)
		_, _ = conn.Do("SET", key, recached)
		return action, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "retrieving action from Redis")
	}

	var action domain.MiddlemanAction
	if unmarshalErr := json.Unmarshal(raw, &action); unmarshalErr != nil {
		return nil, errors.Wrap(unmarshalErr, "unmarshalling action")
	}
	return &action, nil
}

func (r MiddlemanCacheActionRepository) All(
	ctx context.Context,
	schedulerID string,
) ([]*domain.MiddlemanAction, error) {
	// For All operations, we'll use PostgreSQL as the source of truth
	// since Redis KEYS command is expensive
	return r.fallback.All(ctx, schedulerID)
}

func (r MiddlemanCacheActionRepository) GetPendingActions(
	ctx context.Context,
	beforeTime time.Time,
) ([]*domain.MiddlemanAction, error) {
	pool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := pool.Get()
	defer conn.Close()

	pendingKey := fmt.Sprintf("%s:pending_actions", r.prefix)
	maxScore := float64(beforeTime.Unix())
	
	// Get action IDs with execution time <= beforeTime
	actionIDs, err := redis.Strings(conn.Do("ZRANGEBYSCORE", pendingKey, "-inf", maxScore))
	if err != nil {
		return nil, errors.Wrap(err, "getting pending actions from Redis")
	}

	if len(actionIDs) == 0 {
		return []*domain.MiddlemanAction{}, nil
	}

	// Fetch full action details
	var actions []*domain.MiddlemanAction
	for _, actionID := range actionIDs {
		action, err := r.Find(ctx, actionID)
		if err != nil {
			continue // Skip if not found
		}
		if action.Status == "pending" && action.ExecutionTime.Before(beforeTime) {
			actions = append(actions, action)
		}
	}

	return actions, nil
}