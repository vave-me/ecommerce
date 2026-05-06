package redisrepo

import (
	"context"
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

type CatalogCacheRepository struct {
	prefix   string
	fallback domain.CatalogRepository
}

var _ domain.CatalogCacheRepository = (*CatalogCacheRepository)(nil)

func NewCatalogCacheRepository(prefix string, fallback domain.CatalogRepository) (*CatalogCacheRepository, error) {
	return &CatalogCacheRepository{
		prefix:   prefix,
		fallback: fallback,
	}, nil
}

// -----------------------------------------------------------
// 1) AddCategory
// -----------------------------------------------------------
func (r *CatalogCacheRepository) AddCategory(
	ctx context.Context,
	id, description, parentID, googleCategoryID, slug string,
) error {
	// Step A: Write to the fallback DB
	err := r.fallback.AddCategory(ctx, id, description, parentID, googleCategoryID, slug)
	if err != nil {
		return errors.Wrap(err, "fallback AddCategory failed")
	}

	// Step B: Cache in Redis
	if cacheErr := r.cacheSingleCategory(ctx, &domain.CatalogCategory{
		ID:               id,
		Description:      description,
		ParentID:         parentID,
		GoogleCategoryID: googleCategoryID,
		Slug:             slug,
		IsActive:         true, // assume default
	}); cacheErr != nil {
		log.Printf("Warning: caching new category %s in Redis failed: %v", id, cacheErr)
		// not returning the error to keep AddCategory atomic
	}
	return nil
}

// -----------------------------------------------------------
// 2) UpdateCategory
// -----------------------------------------------------------
func (r *CatalogCacheRepository) UpdateCategory(
	ctx context.Context,
	categoryID, description, parentID, googleCategoryID, slug string,
) error {
	// Step A: fallback DB
	err := r.fallback.UpdateCategory(ctx, categoryID, description, parentID, googleCategoryID, slug)
	if err != nil {
		return err
	}

	// Step B: Re-sync Redis
	cp, findErr := r.fallback.Find(ctx, categoryID)
	if findErr == nil && cp != nil {
		if cacheErr := r.cacheSingleCategory(ctx, cp); cacheErr != nil {
			log.Printf("Warning: re-caching category %s in Redis: %v", categoryID, cacheErr)
		}
	}
	return nil
}

// -----------------------------------------------------------
// 3) RemoveCategory
// -----------------------------------------------------------
func (r *CatalogCacheRepository) RemoveCategory(
	ctx context.Context,
	categoryID string,
	userID string, // or userSellerID
) error {
	// Step A: fallback
	if err := r.fallback.RemoveCategory(ctx, categoryID, userID); err != nil {
		return errors.Wrap(err, "fallback RemoveCategory failed")
	}

	// Step B: remove from Redis
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	mainKey := fmt.Sprintf("%s:category:%s", r.prefix, categoryID)
	if _, delErr := conn.Do("DEL", mainKey); delErr != nil {
		log.Printf("Warning: removing category %s from Redis: %v", categoryID, delErr)
	}
	return nil
}

// -----------------------------------------------------------
// 4) ArchiveCategory
// -----------------------------------------------------------
func (r *CatalogCacheRepository) ArchiveCategory(
	ctx context.Context,
	categoryID, userID string,
) error {
	if err := r.fallback.ArchiveCategory(ctx, categoryID, userID); err != nil {
		return err
	}
	// Optionally mark "IsActive=false" in Redis
	return r.updateRedisHashFields(ctx, categoryID, map[string]interface{}{
		"IsActive":  0,
		"updatedAt": time.Now().Unix(),
	})
}

// -----------------------------------------------------------
// 5) RebrandCategory
// -----------------------------------------------------------
func (r *CatalogCacheRepository) RebrandCategory(
	ctx context.Context,
	categoryID, newSlug, newDescription string,
) error {
	// fallback
	if err := r.fallback.RebrandCategory(ctx, categoryID, newSlug, newDescription); err != nil {
		return err
	}
	// update redis
	return r.updateRedisHashFields(ctx, categoryID, map[string]interface{}{
		"Slug":        newSlug,
		"Description": newDescription,
		"updatedAt":   time.Now().Unix(),
	})
}

// -----------------------------------------------------------
// 6) Find
// -----------------------------------------------------------
func (r *CatalogCacheRepository) Find(ctx context.Context, categoryID string) (*domain.CatalogCategory, error) {
	// Step A: check Redis
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:category:%s", r.prefix, categoryID)
	values, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return nil, errors.Wrap(err, "redis HGETALL error in Find")
	}
	if len(values) == 0 {
		// Step B: fallback
		cat, fErr := r.fallback.Find(ctx, categoryID)
		if fErr != nil {
			return nil, errors.Wrap(fErr, "fallback Find failed")
		}
		// Step C: if found, cache
		if cat != nil {
			if err2 := r.cacheSingleCategory(ctx, cat); err2 != nil {
				log.Printf("Warning: caching fallback category %s in Redis: %v", categoryID, err2)
			}
		}
		return cat, nil
	}

	// parse from Redis
	cat, mapErr := r.mapRedisFields(values)
	if mapErr != nil {
		return nil, mapErr
	}
	return cat, nil
}

// -----------------------------------------------------------
// 7) GetCategories
// -----------------------------------------------------------
func (r *CatalogCacheRepository) GetCategories(
	ctx context.Context,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogCategory, int64, error) {
	// We rely on fallback
	categories, totalCount, err := r.fallback.GetCategories(ctx, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, 0, errors.Wrap(err, "fallback GetCategories failed")
	}
	// (Optional) caching
	return categories, totalCount, nil
}

// -----------------------------------------------------------
// 7) GetCategories
// -----------------------------------------------------------
func (r *CatalogCacheRepository) GetMainCategories(
	ctx context.Context,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogCategory, int64, error) {
	// We rely on fallback
	categories, totalCount, err := r.fallback.GetMainCategories(ctx, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, 0, errors.Wrap(err, "fallback GetCategories failed")
	}
	// (Optional) caching
	return categories, totalCount, nil
}

// -----------------------------------------------------------
// 8) GetCatalog
// -----------------------------------------------------------
func (r *CatalogCacheRepository) GetCatalog(
	ctx context.Context,
	userID string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogCategory, int64, error) {
	return r.fallback.GetCatalog(ctx, userID, page, pageSize, sortBy, sortOrder)
}

// -----------------------------------------------------------
// 9) GetSubCategories (NEW METHOD)
// -----------------------------------------------------------
// If you prefer caching subcategories in Redis, adapt similarly
// to how you're handling Find. Here we simply call fallback.
func (r *CatalogCacheRepository) GetSubCategories(
	ctx context.Context,
	parentCategoryID string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.CatalogCategory, int64, error) {
	cats, totalCount, err := r.fallback.GetSubCategories(ctx, parentCategoryID, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, 0, errors.Wrap(err, "fallback GetSubCategories failed")
	}
	// (Optionally) store or retrieve from Redis if desired
	return cats, totalCount, nil
}

// -----------------------------------------------------------
// Internal caching logic
// -----------------------------------------------------------
func (r *CatalogCacheRepository) cacheSingleCategory(ctx context.Context, cat *domain.CatalogCategory) error {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:category:%s", r.prefix, cat.ID)
	// Convert IsActive => int(0/1) or store as string
	isActiveVal := 0
	if cat.IsActive {
		isActiveVal = 1
	}

	// Build the HSET fields
	categoryHash := map[string]interface{}{
		"ID":               cat.ID,
		"Description":      cat.Description,
		"ParentID":         cat.ParentID,
		"GoogleCategoryID": cat.GoogleCategoryID,
		"Slug":             cat.Slug,
		"IsActive":         isActiveVal,
		"updatedAt":        time.Now().Unix(),
	}

	_, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(categoryHash)...)
	return errors.Wrap(err, "HMSET for category cache")
}

func (r *CatalogCacheRepository) mapRedisFields(values map[string]string) (*domain.CatalogCategory, error) {
	isActiveInt, _ := strconv.Atoi(values["IsActive"])
	cat := &domain.CatalogCategory{
		ID:               values["ID"],
		Description:      values["Description"],
		ParentID:         values["ParentID"],
		GoogleCategoryID: values["GoogleCategoryID"],
		Slug:             values["Slug"],
		IsActive:         (isActiveInt == 1),
	}
	return cat, nil
}

func (r *CatalogCacheRepository) updateRedisHashFields(
	ctx context.Context,
	categoryID string,
	updates map[string]interface{},
) error {
	redisPool := di.Get(ctx, constants.RedisPoolKey).(*redis.Pool)
	conn := redisPool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s:category:%s", r.prefix, categoryID)
	exists, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		log.Printf("Warning: EXISTS check for category %s failed: %v", categoryID, err)
		return nil
	}
	if exists == 0 {
		return nil // not cached
	}
	args := redis.Args{}.Add(key)
	for field, value := range updates {
		args = args.Add(field).Add(value)
	}
	_, err = conn.Do("HMSET", args...)
	if err != nil {
		log.Printf("Warning: HMSET partial update in Redis for category %s: %v", categoryID, err)
	}
	return nil
}
