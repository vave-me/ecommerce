package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stackus/errors"
	"middleman/categories/internal/constants"
	"middleman/categories/internal/domain"
	"middleman/internal/di"
)

type CatalogFilterCacheRepository struct {
	prefix   string
	fallback domain.CatalogFilterRepository
}

// Satisfy the embedded interface
var _ domain.CatalogFilterCacheRepository = (*CatalogFilterCacheRepository)(nil)

func NewCatalogFilterCacheRepository(prefix string, fallback domain.CatalogFilterRepository) *CatalogFilterCacheRepository {
	return &CatalogFilterCacheRepository{
		prefix:   prefix,
		fallback: fallback,
	}
}

func (r *CatalogFilterCacheRepository) AddFilter(
	ctx context.Context,
	filterID string,
	categoryID string,
	name string,
	filterType domain.FilterType,
	values []string,
	isActive bool,
) error {
	// 1) DB first
	err := r.fallback.AddFilter(ctx, filterID, categoryID, name, filterType, values, isActive)
	if err != nil {
		return errors.Wrap(err, "fallback AddFilter failed")
	}

	// 2) Cache in Redis
	if cacheErr := r.cacheSingleFilter(ctx, &domain.CatalogFilter{
		ID:         filterID,
		CategoryID: categoryID,
		Name:       name,
		FilterType: filterType,
		Values:     values,
		IsActive:   isActive,
	}); cacheErr != nil {
		log.Printf("Warning: caching new filter %s in Redis failed: %v", filterID, cacheErr)
	}
	return nil
}

func (r *CatalogFilterCacheRepository) UpdateFilter(
	ctx context.Context,
	filterID string,
	name string,
	filterType domain.FilterType,
	values []string,
) error {
	// 1) Fallback
	err := r.fallback.UpdateFilter(ctx, filterID, name, filterType, values)
	if err != nil {
		return errors.Wrap(err, "fallback UpdateFilter failed")
	}

	// 2) Re-sync from DB
	cv, findErr := r.fallback.FindFilter(ctx, filterID)
	if findErr == nil && cv != nil {
		if cacheErr := r.cacheSingleFilter(ctx, cv); cacheErr != nil {
			log.Printf("Warning: re-caching updated filter %s: %v", filterID, cacheErr)
		}
	}
	return nil
}

func (r *CatalogFilterCacheRepository) RemoveFilter(
	ctx context.Context,
	filterID string,
	userID string,
) error {
	// DB
	if err := r.fallback.RemoveFilter(ctx, filterID, userID); err != nil {
		return errors.Wrap(err, "fallback RemoveFilter failed")
	}
	// Redis
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:filter:%s", r.prefix, filterID)
	if _, delErr := conn.Do("DEL", key); delErr != nil {
		log.Printf("Warning: removing filter %s from Redis: %v", filterID, delErr)
	}
	return nil
}

func (r *CatalogFilterCacheRepository) ArchiveFilter(ctx context.Context, filterID string) error {
	// DB
	if err := r.fallback.ArchiveFilter(ctx, filterID); err != nil {
		return err
	}
	// Redis update
	updates := map[string]interface{}{
		"IsActive":  0, // or store as string "false"
		"updatedAt": time.Now().Unix(),
	}
	return r.updateRedisHashFields(ctx, filterID, updates)
}

func (r *CatalogFilterCacheRepository) RebrandFilter(
	ctx context.Context,
	filterID string,
	newName string,
) error {
	// DB
	if err := r.fallback.RebrandFilter(ctx, filterID, newName); err != nil {
		return err
	}
	// Redis
	updates := map[string]interface{}{
		"Name":      newName,
		"updatedAt": time.Now().Unix(),
	}
	return r.updateRedisHashFields(ctx, filterID, updates)
}

func (r *CatalogFilterCacheRepository) FindFilter(
	ctx context.Context,
	filterID string,
) (*domain.CatalogFilter, error) {
	// Redis
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:filter:%s", r.prefix, filterID)
	values, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return nil, errors.Wrap(err, "redis HGETALL for filter")
	}
	if len(values) == 0 {
		// fallback
		cv, fErr := r.fallback.FindFilter(ctx, filterID)
		if fErr != nil {
			return nil, fErr
		}
		// cache if found
		if cv != nil {
			if cErr := r.cacheSingleFilter(ctx, cv); cErr != nil {
				log.Printf("Warning: caching fallback filter %s in Redis: %v", filterID, cErr)
			}
		}
		return cv, nil
	}

	// parse from Redis
	return r.mapRedisFields(values)
}

func (r *CatalogFilterCacheRepository) GetFilters(
	ctx context.Context,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogFilter, int64, error) {
	return r.fallback.GetFilters(ctx, page, pageSize, sortBy, sortOrder)
}

func (r *CatalogFilterCacheRepository) GetFiltersByCategory(
	ctx context.Context,
	categoryID string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogFilter, int64, error) {
	return r.fallback.GetFiltersByCategory(ctx, categoryID, page, pageSize, sortBy, sortOrder)
}

func (r *CatalogFilterCacheRepository) cacheSingleFilter(ctx context.Context, cv *domain.CatalogFilter) error {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:filter:%s", r.prefix, cv.ID)
	// Convert IsActive => int
	isActiveVal := 0
	if cv.IsActive {
		isActiveVal = 1
	}
	// Convert Values => JSON
	valuesJSON, _ := json.Marshal(cv.Values)

	filterHash := map[string]interface{}{
		"ID":         cv.ID,
		"CategoryID": cv.CategoryID,
		"Name":       cv.Name,
		"FilterType": cv.FilterType.String(), // if FilterType is an enum
		"Values":     string(valuesJSON),
		"IsActive":   isActiveVal,
		"updatedAt":  time.Now().Unix(),
	}

	_, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(filterHash)...)
	return errors.Wrap(err, "HMSET for filter in Redis")
}

func (r *CatalogFilterCacheRepository) mapRedisFields(values map[string]string) (*domain.CatalogFilter, error) {
	isActiveInt, _ := strconv.Atoi(values["IsActive"])
	var filterValues []string
	_ = json.Unmarshal([]byte(values["Values"]), &filterValues)

	cv := &domain.CatalogFilter{
		ID:         values["ID"],
		CategoryID: values["CategoryID"],
		Name:       values["Name"],
		// parse FilterType from string if you want
		FilterType: domain.FilterType(values["FilterType"]),
		Values:     filterValues,
		IsActive:   (isActiveInt == 1),
	}
	return cv, nil
}

func (r *CatalogFilterCacheRepository) updateRedisHashFields(ctx context.Context, filterID string, updates map[string]interface{}) error {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:filter:%s", r.prefix, filterID)
	exists, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		log.Printf("Warning: EXISTS check for filter %s failed: %v", filterID, err)
		return nil
	}
	if exists == 0 {
		return nil
	}
	args := redis.Args{}.Add(key)
	for f, v := range updates {
		args = args.Add(f).Add(v)
	}
	if _, err := conn.Do("HMSET", args...); err != nil {
		log.Printf("Warning: HMSET partial update for filter %s: %v", filterID, err)
	}
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
