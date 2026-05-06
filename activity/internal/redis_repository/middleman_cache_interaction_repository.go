package redis_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/stackus/errors"

	"middleman/activity/internal/constants"
	"middleman/activity/internal/domain"
	"middleman/internal/di"
)

type MiddlemanCacheInteractionRepository struct {
	prefix   string
	fallback domain.MiddlemanInteractionRepository
}

// Ensure we satisfy the domain interface at compile time:
var _ domain.MiddlemanCacheInteractionRepository = (*MiddlemanCacheInteractionRepository)(nil)

func NewMiddlemanCacheInteractionRepository(
	prefix string,
	fallback domain.MiddlemanInteractionRepository,
) MiddlemanCacheInteractionRepository {
	return MiddlemanCacheInteractionRepository{
		prefix:   prefix,
		fallback: fallback,
	}
}

// -------------------
// WRITE-THROUGH (Add)
// -------------------
// 1) Write to Redis
// 2) Immediately write to Postgres (no goroutine)
func (r MiddlemanCacheInteractionRepository) Add(
	ctx context.Context,
	interactionID, activityID, itemID, itemType, actionType string,
) error {
	// 1) Write to Redis
	pool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := pool.Get()
	defer conn.Close()

	interaction := domain.MiddlemanInteraction{
		ID:         interactionID,
		ActivityID: activityID,
		ItemID:     itemID,
		ItemType:   itemType,
		ActionType: actionType,
	}

	raw, err := json.Marshal(interaction)
	if err != nil {
		return errors.Wrap(err, "marshalling interaction")
	}

	key := fmt.Sprintf("%s:interaction:%s", r.prefix, interactionID)
	if _, err := conn.Do("SET", key, raw); err != nil {
		return errors.Wrap(err, "storing interaction in Redis")
	}

	// If like/dislike, increment counters
	if actionType == "like" || actionType == "dislike" {
		statsKey := fmt.Sprintf("%s:stats:%s:%s", r.prefix, actionType, itemType)
		if _, err := conn.Do("ZINCRBY", statsKey, 1, itemID); err != nil {
			return errors.Wrap(err, "incrementing reaction counter in Redis")
		}
	}

	// 2) Immediately write to Postgres
	if err := r.fallback.Add(ctx, interactionID, activityID, itemID, itemType, actionType); err != nil {
		return errors.Wrap(err, "inserting interaction into Postgres")
	}

	return nil
}

// ----------------------
// WRITE-THROUGH (Update)
// ----------------------
// 1) Retrieve & modify in Redis
// 2) Immediately update Postgres (no goroutine)
func (r MiddlemanCacheInteractionRepository) Update(
	ctx context.Context,
	interactionID, newActionType string,
) error {
	// 1) Retrieve from Redis
	pool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := pool.Get()
	defer conn.Close()

	existing, err := r.Find(ctx, interactionID)
	if err != nil {
		return err // not found or other error
	}
	oldAction := existing.ActionType
	itemID := existing.ItemID
	itemType := existing.ItemType

	// Adjust counters if action changed
	if oldAction != newActionType {
		if oldAction == "like" || oldAction == "dislike" {
			oldKey := fmt.Sprintf("%s:stats:%s:%s", r.prefix, oldAction, itemType)
			if _, decrErr := conn.Do("ZINCRBY", oldKey, -1, itemID); decrErr != nil {
				return errors.Wrap(decrErr, "decrementing old reaction counter")
			}
		}
		if newActionType == "like" || newActionType == "dislike" {
			newKey := fmt.Sprintf("%s:stats:%s:%s", r.prefix, newActionType, itemType)
			if _, incrErr := conn.Do("ZINCRBY", newKey, 1, itemID); incrErr != nil {
				return errors.Wrap(incrErr, "incrementing new reaction counter")
			}
		}
		existing.ActionType = newActionType
	}

	// Store updated record in Redis
	updatedData, err := json.Marshal(existing)
	if err != nil {
		return errors.Wrap(err, "marshalling updated interaction")
	}
	mainKey := fmt.Sprintf("%s:interaction:%s", r.prefix, interactionID)
	if _, err := conn.Do("SET", mainKey, updatedData); err != nil {
		return errors.Wrap(err, "storing updated interaction in Redis")
	}

	// 2) Synchronously update Postgres
	if err := r.fallback.Update(ctx, interactionID, newActionType); err != nil {
		return errors.Wrap(err, "updating interaction in Postgres")
	}

	return nil
}

// ----------------------
// WRITE-THROUGH (Remove)
// ----------------------
// 1) Remove from Redis
// 2) Remove from Postgres
func (r MiddlemanCacheInteractionRepository) Remove(
	ctx context.Context,
	interactionID string,
) error {
	pool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := pool.Get()
	defer conn.Close()

	// Retrieve from Redis so we know old action
	existing, err := r.Find(ctx, interactionID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			// nothing to remove
			return nil
		}
		return err
	}

	oldAction := existing.ActionType
	itemID := existing.ItemID
	itemType := existing.ItemType

	// 1) Remove from Redis
	mainKey := fmt.Sprintf("%s:interaction:%s", r.prefix, interactionID)
	if _, err := conn.Do("DEL", mainKey); err != nil {
		return errors.Wrap(err, "deleting interaction from Redis")
	}

	// If it was like/dislike, decrement counters
	if oldAction == "like" || oldAction == "dislike" {
		oldKey := fmt.Sprintf("%s:stats:%s:%s", r.prefix, oldAction, itemType)
		if _, decrErr := conn.Do("ZINCRBY", oldKey, -1, itemID); decrErr != nil {
			return errors.Wrap(decrErr, "decrementing reaction counter in Redis")
		}
	}

	// 2) Immediately remove from Postgres
	if err := r.fallback.Remove(ctx, interactionID); err != nil {
		return errors.Wrap(err, "removing interaction in Postgres")
	}

	return nil
}

// ----------------------------
// FIND (Read only from Redis?)
// ----------------------------
// If you want fallback, do so. We'll show the optional fallback logic.
func (r MiddlemanCacheInteractionRepository) Find(
	ctx context.Context,
	interactionID string,
) (*domain.MiddlemanInteraction, error) {
	pool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:interaction:%s", r.prefix, interactionID)
	raw, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		// Optionally fallback to Postgres if not found:
		// fallbackInteraction, fbErr := r.fallback.Find(ctx, interactionID)
		// if fbErr != nil {
		//     return nil, fbErr
		// }
		// // Re-cache if found
		// recached, _ := json.Marshal(fallbackInteraction)
		// _, _ = conn.Do("SET", key, recached)
		// return fallbackInteraction, nil

		return nil, errors.ErrNotFound.Msg("interaction not found in Redis")
	} else if err != nil {
		return nil, errors.Wrap(err, "retrieving interaction from Redis")
	}

	var interaction domain.MiddlemanInteraction
	if unmarshalErr := json.Unmarshal(raw, &interaction); unmarshalErr != nil {
		return nil, errors.Wrap(unmarshalErr, "unmarshalling interaction")
	}
	return &interaction, nil
}

// All scans Redis for all keys matching the activityID. (No fallback).
// If you want fallback, you could do a similar pattern as in Add or an "All" fallback.
func (r MiddlemanCacheInteractionRepository) All(
	ctx context.Context,
	activityID string,
) ([]*domain.MiddlemanInteraction, error) {
	pool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := pool.Get()
	defer conn.Close()

	pattern := fmt.Sprintf("%s:interaction:*", r.prefix)
	keys, err := redis.Strings(conn.Do("KEYS", pattern))
	if err != nil {
		return nil, errors.Wrap(err, "retrieving interaction keys from Redis")
	}

	var results []*domain.MiddlemanInteraction
	for _, key := range keys {
		raw, gErr := redis.Bytes(conn.Do("GET", key))
		if gErr != nil || len(raw) == 0 {
			continue
		}
		var interaction domain.MiddlemanInteraction
		if err := json.Unmarshal(raw, &interaction); err != nil {
			continue
		}
		if interaction.ActivityID == activityID {
			results = append(results, &interaction)
		}
	}

	// Optional fallback if you want:
	// if len(results) == 0 {
	//     dbInteractions, dbErr := r.fallback.All(ctx, activityID)
	//     if dbErr != nil {
	//         return nil, dbErr
	//     }
	//     // Re-cache in Redis
	//     for _, it := range dbInteractions {
	//         recached, _ := json.Marshal(it)
	//         key := fmt.Sprintf("%s:interaction:%s", r.prefix, it.ID)
	//         conn.Do("SET", key, recached)
	//     }
	//     return dbInteractions, nil
	// }

	return results, nil
}

// GetMostLiked reads from Redis ZSET "stats:like:<itemType>" in descending order. (No fallback).
func (r MiddlemanCacheInteractionRepository) GetMostLiked(
	ctx context.Context,
	itemType string,
	limit int64,
) ([]*domain.MostReactionResult, error) {
	return r.getMostByAction(ctx, itemType, "like", limit)
}

// GetMostDisliked reads from Redis ZSET "stats:dislike:<itemType>" in descending order. (No fallback).
func (r MiddlemanCacheInteractionRepository) GetMostDisliked(
	ctx context.Context,
	itemType string,
	limit int64,
) ([]*domain.MostReactionResult, error) {
	return r.getMostByAction(ctx, itemType, "dislike", limit)
}

// getMostByAction does a synchronous ZREVRANGE on Redis, no fallback.
func (r MiddlemanCacheInteractionRepository) getMostByAction(
	ctx context.Context,
	itemType, action string,
	limit int64,
) ([]*domain.MostReactionResult, error) {
	pool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:stats:%s:%s", r.prefix, action, itemType)
	if limit <= 0 {
		limit = -1 // no limit
	}
	end := limit - 1

	values, err := redis.Strings(conn.Do("ZREVRANGE", key, 0, end, "WITHSCORES"))
	if err != nil {
		return nil, errors.Wrap(err, "retrieving top items from Redis")
	}
	if len(values) == 0 {
		return nil, nil
	}

	results := make([]*domain.MostReactionResult, 0, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		itemID := values[i]
		scoreStr := values[i+1]

		score, convErr := strconv.ParseInt(scoreStr, 10, 64)
		if convErr != nil {
			// skip malformed
			continue
		}
		results = append(results, &domain.MostReactionResult{
			ItemID:   itemID,
			ItemType: itemType,
			Action:   action,
			Count:    score,
		})
	}
	return results, nil
}
