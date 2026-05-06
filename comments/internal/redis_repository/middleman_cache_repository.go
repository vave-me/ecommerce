package redis_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stackus/errors"
	"middleman/comments/internal/constants"
	"middleman/comments/internal/domain"
	"middleman/internal/di"
)

type MiddlemanCacheRepository struct {
	prefix   string
	fallback domain.MiddlemanRepository
}

func NewMiddlemanCacheRepository(prefix string, fallback domain.MiddlemanRepository) MiddlemanCacheRepository {
	return MiddlemanCacheRepository{
		prefix:   prefix,
		fallback: fallback,
	}
}

func (r MiddlemanCacheRepository) Add(
	ctx context.Context,
	commentID, senderID, itemID string,
	itemType domain.ItemType,
	content, categoryID, parentID string,
) error {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// 1) Create the comment struct
	comment := domain.MiddlemanComment{
		ID:         commentID,
		SenderID:   senderID,
		ItemID:     itemID,
		ItemType:   itemType.String(),
		Content:    content,
		CategoryID: categoryID,
		ParentID:   parentID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Approved:   true,  // or default
		Flagged:    false, // or default
	}

	// 2) Marshal it to JSON
	data, err := json.Marshal(comment)
	if err != nil {
		return errors.Wrap(err, "marshalling comment")
	}

	// 3) Store comment JSON in "prefix:comment:{commentID}"
	commentKey := fmt.Sprintf("%s:comment:%s", r.prefix, commentID)
	if _, err = db.Do("SET", commentKey, data); err != nil {
		return errors.Wrap(err, "storing comment in Redis")
	}

	// 4) Add commentID to item set: "prefix:item:{itemID}:comments"
	itemKey := fmt.Sprintf("%s:item:%s:comments", r.prefix, itemID)
	if _, err = db.Do("SADD", itemKey, commentID); err != nil {
		return errors.Wrap(err, "adding comment ID to item set")
	}

	// 5) Add commentID to sender set: "prefix:sender:{senderID}:comments"
	senderKey := fmt.Sprintf("%s:sender:%s:comments", r.prefix, senderID)
	if _, err = db.Do("SADD", senderKey, commentID); err != nil {
		return errors.Wrap(err, "adding comment ID to sender set")
	}

	// 6) Add commentID to item_type set: "prefix:item_type:{itemType}:comments"
	typeKey := fmt.Sprintf("%s:item_type:%s:comments", r.prefix, itemType.String())
	if _, err = db.Do("SADD", typeKey, commentID); err != nil {
		return errors.Wrap(err, "adding comment ID to item_type set")
	}

	// 7) Update aggregator stats
	//    a) global item stats
	statsKey := fmt.Sprintf("%s:stats:comment_count", r.prefix)
	if _, err = db.Do("ZINCRBY", statsKey, 1, itemID); err != nil {
		return errors.Wrap(err, "incrementing item comment count")
	}

	//    b) category-level stats
	if categoryID != "" {
		catStatsKey := fmt.Sprintf("%s:stats:comment_count:category:%s", r.prefix, categoryID)
		if _, err = db.Do("ZINCRBY", catStatsKey, 1, itemID); err != nil {
			return errors.Wrap(err, "incrementing category-level comment count")
		}
	}

	//    c) item_type-level stats
	if itemType.String() != "" {
		typeStatsKey := fmt.Sprintf("%s:stats:comment_count:item_type:%s", r.prefix, itemType.String())
		if _, err = db.Do("ZINCRBY", typeStatsKey, 1, itemID); err != nil {
			return errors.Wrap(err, "incrementing itemType comment count")
		}
	}

	return nil
}

func (r MiddlemanCacheRepository) Find(ctx context.Context, commentID, itemID string) (*domain.MiddlemanComment, error) {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// Directly GET "prefix:comment:{commentID}"
	commentKey := fmt.Sprintf("%s:comment:%s", r.prefix, commentID)
	data, err := redis.Bytes(db.Do("GET", commentKey))
	if err != nil {
		if err == redis.ErrNil {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, "retrieving comment from Redis")
	}

	var comment domain.MiddlemanComment
	if err := json.Unmarshal(data, &comment); err != nil {
		return nil, errors.Wrap(err, "unmarshalling comment JSON")
	}

	// Optionally verify itemID
	if comment.ItemID != itemID {
		return nil, errors.ErrNotFound
	}
	return &comment, nil
}

func (r MiddlemanCacheRepository) All(ctx context.Context, itemID string) ([]*domain.MiddlemanComment, error) {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// S = "prefix:item:{itemID}:comments"
	itemKey := fmt.Sprintf("%s:item:%s:comments", r.prefix, itemID)
	ids, err := redis.Strings(db.Do("SMEMBERS", itemKey))
	if err != nil {
		return nil, errors.Wrap(err, "retrieving comment IDs by item")
	}
	if len(ids) == 0 {
		return nil, nil
	}

	// MGET all keys "prefix:comment:{commentID}"
	keys := make([]interface{}, 0, len(ids))
	for _, cid := range ids {
		keys = append(keys, fmt.Sprintf("%s:comment:%s", r.prefix, cid))
	}

	values, err := redis.ByteSlices(db.Do("MGET", keys...))
	if err != nil {
		return nil, errors.Wrap(err, "retrieving comments with MGET")
	}

	comments := make([]*domain.MiddlemanComment, 0, len(values))
	for _, val := range values {
		if val == nil {
			continue
		}
		c := new(domain.MiddlemanComment)
		if err := json.Unmarshal(val, c); err != nil {
			return nil, errors.Wrap(err, "unmarshalling comment JSON")
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r MiddlemanCacheRepository) FindBySenderID(ctx context.Context, senderID string) ([]*domain.MiddlemanComment, error) {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// "prefix:sender:{senderID}:comments"
	senderKey := fmt.Sprintf("%s:sender:%s:comments", r.prefix, senderID)
	ids, err := redis.Strings(db.Do("SMEMBERS", senderKey))
	if err != nil {
		return nil, errors.Wrap(err, "retrieving comment IDs by sender")
	}
	if len(ids) == 0 {
		return nil, nil
	}

	keys := make([]interface{}, 0, len(ids))
	for _, cid := range ids {
		keys = append(keys, fmt.Sprintf("%s:comment:%s", r.prefix, cid))
	}

	values, err := redis.ByteSlices(db.Do("MGET", keys...))
	if err != nil {
		return nil, errors.Wrap(err, "retrieving comments with MGET")
	}

	var comments []*domain.MiddlemanComment
	for _, val := range values {
		if val == nil {
			continue
		}
		c := new(domain.MiddlemanComment)
		if err := json.Unmarshal(val, c); err != nil {
			return nil, errors.Wrap(err, "unmarshalling comment JSON")
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// New method: retrieve all comments of a certain item_type
func (r MiddlemanCacheRepository) AllByItemType(ctx context.Context, itemType domain.ItemType) ([]*domain.MiddlemanComment, error) {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// "prefix:item_type:{itemType}:comments"
	typeKey := fmt.Sprintf("%s:item_type:%s:comments", r.prefix, itemType.String())
	ids, err := redis.Strings(db.Do("SMEMBERS", typeKey))
	if err != nil {
		return nil, errors.Wrap(err, "retrieving comment IDs by item_type")
	}
	if len(ids) == 0 {
		return nil, nil
	}

	// MGET "prefix:comment:{commentID}"
	keys := make([]interface{}, 0, len(ids))
	for _, cid := range ids {
		keys = append(keys, fmt.Sprintf("%s:comment:%s", r.prefix, cid))
	}

	values, err := redis.ByteSlices(db.Do("MGET", keys...))
	if err != nil {
		return nil, errors.Wrap(err, "MGET by item_type comments")
	}

	var comments []*domain.MiddlemanComment
	for _, val := range values {
		if val == nil {
			continue
		}
		c := new(domain.MiddlemanComment)
		if err := json.Unmarshal(val, c); err != nil {
			return nil, errors.Wrap(err, "unmarshalling comment JSON")
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r MiddlemanCacheRepository) Remove(ctx context.Context, commentID string) error {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// 1) GET the comment
	commentKey := fmt.Sprintf("%s:comment:%s", r.prefix, commentID)
	data, err := redis.Bytes(db.Do("GET", commentKey))
	if err != nil {
		if err == redis.ErrNil {
			return nil // or errors.ErrNotFound
		}
		return errors.Wrap(err, "getting comment for removal")
	}

	var comment domain.MiddlemanComment
	if err := json.Unmarshal(data, &comment); err != nil {
		return errors.Wrap(err, "unmarshalling comment prior to removal")
	}

	// 2) Remove from item set
	itemKey := fmt.Sprintf("%s:item:%s:comments", r.prefix, comment.ItemID)
	if _, err = db.Do("SREM", itemKey, commentID); err != nil {
		return errors.Wrap(err, "removing comment ID from item set")
	}

	// 3) Remove from sender set
	senderKey := fmt.Sprintf("%s:sender:%s:comments", r.prefix, comment.SenderID)
	if _, err = db.Do("SREM", senderKey, commentID); err != nil {
		return errors.Wrap(err, "removing comment ID from sender set")
	}

	// 4) Remove from item_type set
	if comment.ItemType != "" {
		typeKey := fmt.Sprintf("%s:item_type:%s:comments", r.prefix, comment.ItemType)
		if _, err = db.Do("SREM", typeKey, commentID); err != nil {
			return errors.Wrap(err, "removing comment ID from itemType set")
		}
	}

	// 5) Delete the comment key
	if _, err = db.Do("DEL", commentKey); err != nil {
		return errors.Wrap(err, "deleting comment key")
	}

	// 6) Decrement aggregator stats
	statsKey := fmt.Sprintf("%s:stats:comment_count", r.prefix)
	if _, err = db.Do("ZINCRBY", statsKey, -1, comment.ItemID); err != nil {
		return errors.Wrap(err, "decrementing aggregator stats")
	}

	// If category is set, decrement category aggregator
	if comment.CategoryID != "" {
		catStatsKey := fmt.Sprintf("%s:stats:comment_count:category:%s", r.prefix, comment.CategoryID)
		if _, err = db.Do("ZINCRBY", catStatsKey, -1, comment.ItemID); err != nil {
			return errors.Wrap(err, "decrementing category aggregator stats")
		}
	}

	// Decrement item_type aggregator
	if comment.ItemType != "" {
		typeStatsKey := fmt.Sprintf("%s:stats:comment_count:item_type:%s", r.prefix, comment.ItemType)
		if _, err = db.Do("ZINCRBY", typeStatsKey, -1, comment.ItemID); err != nil {
			return errors.Wrap(err, "decrementing itemType aggregator stats")
		}
	}

	return nil
}

// For aggregator queries (e.g., MostCommentedItemsByItemType), we do something similar to category:
func (r MiddlemanCacheRepository) MostCommentedItemsByItemType(
	ctx context.Context,
	itemType domain.ItemType,
	limit, offset int,
) ([]*domain.ItemCommentCount, error) {

	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	// e.g. "prefix:stats:comment_count:item_type:books"
	typeStatsKey := fmt.Sprintf("%s:stats:comment_count:item_type:%s", r.prefix, itemType.String())

	end := offset + limit - 1
	if limit <= 0 {
		end = -1
	}

	// ZREVRANGE ... WITHSCORES
	values, err := redis.Strings(db.Do("ZREVRANGE", typeStatsKey, offset, end, "WITHSCORES"))
	if err != nil {
		return nil, errors.Wrap(err, "retrieving top commented items from Redis by item_type")
	}

	var results []*domain.ItemCommentCount
	for i := 0; i < len(values); i += 2 {
		itemID := values[i]
		scoreStr := values[i+1]

		score, convErr := redis.Int(scoreStr, nil)
		if convErr != nil {
			return nil, errors.Wrap(convErr, "converting score to int")
		}
		results = append(results, &domain.ItemCommentCount{
			ItemID:        itemID,
			ItemType:      itemType.String(),
			CategoryID:    "", // not tracked in this aggregator
			CommentsCount: int64(score),
		})
	}
	return results, nil
}

// The aggregator "most commented items" logic is the same approach (using sorted sets).
func (r MiddlemanCacheRepository) MostCommentedItems(
	ctx context.Context,
	limit, offset int,
) ([]*domain.ItemCommentCount, error) {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	statsKey := fmt.Sprintf("%s:stats:comment_count", r.prefix)

	end := offset + limit - 1
	if limit <= 0 {
		end = -1
	}
	values, err := redis.Strings(db.Do("ZREVRANGE", statsKey, offset, end, "WITHSCORES"))
	if err != nil {
		return nil, errors.Wrap(err, "retrieving top commented items from Redis")
	}

	var results []*domain.ItemCommentCount
	for i := 0; i < len(values); i += 2 {
		itemID := values[i]
		scoreStr := values[i+1]

		score, convErr := redis.Int(scoreStr, nil)
		if convErr != nil {
			return nil, errors.Wrap(convErr, "converting score to int")
		}
		results = append(results, &domain.ItemCommentCount{
			ItemID:        itemID,
			ItemType:      "", // Not tracked in this global ZSET
			CategoryID:    "",
			CommentsCount: int64(score),
		})
	}

	return results, nil
}

func (r MiddlemanCacheRepository) MostCommentedItemsByCategory(
	ctx context.Context,
	itemType domain.ItemType,
	categoryID string,
	limit, offset int,
) ([]*domain.ItemCommentCount, error) {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	db := redisPool.Get()
	defer db.Close()

	catStatsKey := fmt.Sprintf("%s:stats:comment_count:category:%s", r.prefix, categoryID)

	end := offset + limit - 1
	if limit <= 0 {
		end = -1
	}
	values, err := redis.Strings(db.Do("ZREVRANGE", catStatsKey, offset, end, "WITHSCORES"))
	if err != nil {
		return nil, errors.Wrap(err, "retrieving top commented items from Redis (category)")
	}

	var results []*domain.ItemCommentCount
	for i := 0; i < len(values); i += 2 {
		itemID := values[i]
		scoreStr := values[i+1]

		score, convErr := redis.Int(scoreStr, nil)
		if convErr != nil {
			return nil, errors.Wrap(convErr, "converting score to int")
		}
		results = append(results, &domain.ItemCommentCount{
			ItemID:        itemID,
			ItemType:      itemType.String(),
			CategoryID:    categoryID,
			CommentsCount: int64(score),
		})
	}
	return results, nil
}
