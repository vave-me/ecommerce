package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"middleman/internal/postgres"
	"middleman/metrics/internal/application"
	"middleman/metrics/internal/models"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"
)

type ItemMetricCacheRepository struct {
	tableName string
	db        postgres.DB
}

var _ application.ItemMetricCacheRepository = (*ItemMetricCacheRepository)(nil)

func NewItemMetricCacheRepository(tableName string, db postgres.DB) ItemMetricCacheRepository {
	return ItemMetricCacheRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r ItemMetricCacheRepository) AddMetric(ctx context.Context, itemId, entityType, categoryId string, price int64, lat, lng float64) error {
	const query = `INSERT INTO %s (
		id, entity_type, likes_count, dislikes_count, comments_count, 
		messages_count, shared_count, added_to_wishlist_count, added_to_basket_count, 
		visited_count, reported_count, follower_count, reviews_count, rating_count,
		videos_count, images_count, rating, review, category, category_id, category_slug,
        price, lat, lng, location, created_at, updated_at
	) VALUES ($1, $2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, '', $3, '', $4, $5, $6, ST_SetSRID(ST_MakePoint($6, $5), 4326), NOW(), NOW())`

	_, err := r.db.ExecContext(ctx, r.table(query), itemId, entityType, categoryId, price, lat, lng)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			// Ignore unique violation errors (metric already exists)
			return nil
		}
		return errors.Wrap(err, "adding item metric")
	}

	return nil
}

func (r ItemMetricCacheRepository) GetItemMetric(ctx context.Context, itemId string) (*models.ItemMetric, error) {
	const query = `SELECT 
		id, entity_type, likes_count, dislikes_count, comments_count, 
		messages_count, shared_count, added_to_wishlist_count, added_to_basket_count, 
		visited_count, reported_count, follower_count, reviews_count, rating_count,
		videos_count, images_count, rating, review, category, category_id, category_slug,
		price, lat, lng, created_at, updated_at
	FROM %s WHERE id = $1 LIMIT 1`

	metric := &models.ItemMetric{
		ID: itemId,
	}

	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, r.table(query), itemId).Scan(
		&metric.ID, &metric.EntityType, &metric.LikesCount, &metric.DislikesCount, &metric.CommentsCount,
		&metric.MessagesCount, &metric.SharedCount, &metric.AddedToWishlistCount, &metric.AddedToBasketCount,
		&metric.VisitedCount, &metric.ReportedCount, &metric.FollowerCount, &metric.ReviewsCount, &metric.RatingCount,
		&metric.VideosCount, &metric.ImagesCount, &metric.Rating, &metric.Review, &metric.Category,
		&metric.CategoryID, &metric.CategorySlug, &metric.Price, &metric.Lat, &metric.Lng, &createdAt, &updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msgf("item metric with ID %s not found", itemId)
		}
		return nil, errors.Wrap(err, "getting item metric")
	}

	metric.CreatedAt = createdAt.Format(time.RFC3339)
	metric.UpdatedAt = updatedAt.Format(time.RFC3339)

	return metric, nil
}

// GetItemsMetrics retrieves multiple metrics by IDs, with a limit of max 150 items
func (r ItemMetricCacheRepository) GetItemsMetrics(ctx context.Context, itemIds []string) ([]*models.ItemMetric, error) {
	if len(itemIds) == 0 {
		return []*models.ItemMetric{}, nil
	}

	// Limit to maximum 150 items
	if len(itemIds) > 150 {
		itemIds = itemIds[:150]
	}

	// Build the query with the correct number of placeholders
	baseQuery := `SELECT 
		id, entity_type, likes_count, dislikes_count, comments_count, 
		messages_count, shared_count, added_to_wishlist_count, added_to_basket_count, 
		visited_count, reported_count, follower_count, reviews_count, rating_count,
		videos_count, images_count, rating, review, category, category_id, category_slug,
		price, lat, lng, created_at, updated_at
	FROM %s WHERE id = ANY($1)`

	// Create params for the query
	idArray := "{" + strings.Join(itemIds, ",") + "}"

	// Execute the query
	rows, err := r.db.QueryContext(ctx, r.table(baseQuery), idArray)
	if err != nil {
		return nil, errors.Wrap(err, "querying for item metrics")
	}
	defer rows.Close()

	// Create a map to store results by ID for easier lookup
	results := make(map[string]*models.ItemMetric)

	// Process the results
	for rows.Next() {
		metric := &models.ItemMetric{}
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&metric.ID, &metric.EntityType, &metric.LikesCount, &metric.DislikesCount, &metric.CommentsCount,
			&metric.MessagesCount, &metric.SharedCount, &metric.AddedToWishlistCount, &metric.AddedToBasketCount,
			&metric.VisitedCount, &metric.ReportedCount, &metric.FollowerCount, &metric.ReviewsCount, &metric.RatingCount,
			&metric.VideosCount, &metric.ImagesCount, &metric.Rating, &metric.Review, &metric.Category,
			&metric.CategoryID, &metric.CategorySlug, &metric.Price, &metric.Lat, &metric.Lng, &createdAt, &updatedAt,
		)

		if err != nil {
			return nil, errors.Wrap(err, "scanning item metric row")
		}

		metric.CreatedAt = createdAt.Format(time.RFC3339)
		metric.UpdatedAt = updatedAt.Format(time.RFC3339)

		results[metric.ID] = metric
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "iterating item metric rows")
	}

	// Convert map to slice in the same order as requested
	metrics := make([]*models.ItemMetric, 0, len(results))
	for _, itemId := range itemIds {
		if metric, ok := results[itemId]; ok {
			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

func (r ItemMetricCacheRepository) RemoveItemMetric(ctx context.Context, itemId string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), itemId)
	if err != nil {
		return errors.Wrap(err, "removing item metric")
	}

	return nil
}

func (r ItemMetricCacheRepository) UpdateItemMetric(ctx context.Context, itemId, metricType, metricTypeAction string) error {
	var query string

	switch metricType {
	case models.MetricTypeCountLike.String():
		if metricTypeAction == models.MetricTypeActionAdd.String() {
			query = `UPDATE %s SET likes_count = likes_count + 1, updated_at = NOW() WHERE id = $1`
		} else {
			query = `UPDATE %s SET likes_count = GREATEST(likes_count - 1, 0), updated_at = NOW() WHERE id = $1`
		}
	case models.MetricTypeCountDislike.String():
		if metricTypeAction == models.MetricTypeActionAdd.String() {
			query = `UPDATE %s SET dislikes_count = dislikes_count + 1, updated_at = NOW() WHERE id = $1`
		} else {
			query = `UPDATE %s SET dislikes_count = GREATEST(dislikes_count - 1, 0), updated_at = NOW() WHERE id = $1`
		}
	case models.MetricTypeCountComment.String():
		if metricTypeAction == models.MetricTypeActionAdd.String() {
			query = `UPDATE %s SET comments_count = comments_count + 1, updated_at = NOW() WHERE id = $1`
		} else {
			query = `UPDATE %s SET comments_count = GREATEST(comments_count - 1, 0), updated_at = NOW() WHERE id = $1`
		}
	case models.MetricTypeCountMessage.String():
		if metricTypeAction == models.MetricTypeActionAdd.String() {
			query = `UPDATE %s SET messages_count = messages_count + 1, updated_at = NOW() WHERE id = $1`
		} else {
			query = `UPDATE %s SET messages_count = GREATEST(messages_count - 1, 0), updated_at = NOW() WHERE id = $1`
		}
	case models.MetricTypeCountShare.String():
		if metricTypeAction == models.MetricTypeActionAdd.String() {
			query = `UPDATE %s SET shared_count = shared_count + 1, updated_at = NOW() WHERE id = $1`
		} else {
			query = `UPDATE %s SET shared_count = GREATEST(shared_count - 1, 0), updated_at = NOW() WHERE id = $1`
		}
	case models.MetricTypeCountWishlist.String():
		if metricTypeAction == models.MetricTypeActionAdd.String() {
			query = `UPDATE %s SET added_to_wishlist_count = added_to_wishlist_count + 1, updated_at = NOW() WHERE id = $1`
		} else {
			query = `UPDATE %s SET added_to_wishlist_count = GREATEST(added_to_wishlist_count - 1, 0), updated_at = NOW() WHERE id = $1`
		}
	case models.MetricTypeCountBasket.String():
		if metricTypeAction == models.MetricTypeActionAdd.String() {
			query = `UPDATE %s SET added_to_basket_count = added_to_basket_count + 1, updated_at = NOW() WHERE id = $1`
		} else {
			query = `UPDATE %s SET added_to_basket_count = GREATEST(added_to_basket_count - 1, 0), updated_at = NOW() WHERE id = $1`
		}
	case models.MetricTypeCountVisit.String():
		if metricTypeAction == models.MetricTypeActionAdd.String() {
			query = `UPDATE %s SET visited_count = visited_count + 1, updated_at = NOW() WHERE id = $1`
		} else {
			query = `UPDATE %s SET visited_count = GREATEST(visited_count - 1, 0), updated_at = NOW() WHERE id = $1`
		}
	case models.MetricTypeCountReport.String():
		if metricTypeAction == models.MetricTypeActionAdd.String() {
			query = `UPDATE %s SET reported_count = reported_count + 1, updated_at = NOW() WHERE id = $1`
		} else {
			query = `UPDATE %s SET reported_count = GREATEST(reported_count - 1, 0), updated_at = NOW() WHERE id = $1`
		}
	case models.MetricTypeCountFollow.String():
		if metricTypeAction == models.MetricTypeActionAdd.String() {
			query = `UPDATE %s SET follower_count = follower_count + 1, updated_at = NOW() WHERE id = $1`
		} else {
			query = `UPDATE %s SET follower_count = GREATEST(follower_count - 1, 0), updated_at = NOW() WHERE id = $1`
		}
	case models.MetricTypeCountReview.String():
		if metricTypeAction == models.MetricTypeActionAdd.String() {
			query = `UPDATE %s SET reviews_count = reviews_count + 1, updated_at = NOW() WHERE id = $1`
		} else {
			query = `UPDATE %s SET reviews_count = GREATEST(reviews_count - 1, 0), updated_at = NOW() WHERE id = $1`
		}
	default:
		return nil
	}

	_, err := r.db.ExecContext(ctx, r.table(query), itemId)
	if err != nil {
		return errors.Wrap(err, "updating item metric")
	}

	return nil
}

func (r ItemMetricCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

func (r ItemMetricCacheRepository) GetHighestMetricsByType(ctx context.Context, metricType string, entityTypes []models.EntityType, categoryId string, lat, lng, radius float64, minPrice, maxPrice int64, createdFrom, CreatedTill string) ([]*models.ItemMetric, error) {
	return r.getMetricsByType(ctx, metricType, entityTypes, categoryId, lat, lng, radius, minPrice, maxPrice, createdFrom, CreatedTill, "DESC")
}

func (r ItemMetricCacheRepository) GetLowestMetricsByType(ctx context.Context, metricType string, entityTypes []models.EntityType, categoryId string, lat, lng, radius float64, minPrice, maxPrice int64, createdFrom, CreatedTill string) ([]*models.ItemMetric, error) {
	return r.getMetricsByType(ctx, metricType, entityTypes, categoryId, lat, lng, radius, minPrice, maxPrice, createdFrom, CreatedTill, "ASC")
}

func (r ItemMetricCacheRepository) getMetricsByType(ctx context.Context, metricType string, entityTypes []models.EntityType, categoryId string, lat, lng, radius float64, minPrice, maxPrice int64, createdFrom, createdTill, sortOrder string) ([]*models.ItemMetric, error) {
	var metricColumn string
	switch metricType {
	case models.MetricTypeCountLike.String():
		metricColumn = "likes_count"
	case models.MetricTypeCountDislike.String():
		metricColumn = "dislikes_count"
	case models.MetricTypeCountComment.String():
		metricColumn = "comments_count"
	case models.MetricTypeCountMessage.String():
		metricColumn = "messages_count"
	case models.MetricTypeCountShare.String():
		metricColumn = "shared_count"
	case models.MetricTypeCountWishlist.String():
		metricColumn = "added_to_wishlist_count"
	case models.MetricTypeCountBasket.String():
		metricColumn = "added_to_basket_count"
	case models.MetricTypeCountVisit.String():
		metricColumn = "visited_count"
	case models.MetricTypeCountReport.String():
		metricColumn = "reported_count"
	case models.MetricTypeCountFollow.String():
		metricColumn = "follower_count"
	case models.MetricTypeCountReview.String():
		metricColumn = "reviews_count"
	default:
		return []*models.ItemMetric{}, nil
	}

	var whereClause strings.Builder
	var args []interface{}
	argCount := 0

	// Build dynamic WHERE clause
	whereClause.WriteString("WHERE 1=1")

	// Category filter
	if categoryId != "" {
		argCount++
		whereClause.WriteString(fmt.Sprintf(" AND category_id = $%d", argCount))
		args = append(args, categoryId)
	}

	// Price range filter
	if minPrice > 0 {
		argCount++
		whereClause.WriteString(fmt.Sprintf(" AND price >= $%d", argCount))
		args = append(args, minPrice)
	}
	if maxPrice > 0 {
		argCount++
		whereClause.WriteString(fmt.Sprintf(" AND price <= $%d", argCount))
		args = append(args, maxPrice)
	}

	// Geospatial filter using PostGIS ST_DWithin
	if lat != 0 && lng != 0 && radius > 0 {
		argCount++
		whereClause.WriteString(fmt.Sprintf(" AND ST_DWithin(location, ST_SetSRID(ST_MakePoint($%d, $%d), 4326), $%d)", argCount, argCount+1, argCount+2))
		args = append(args, lng, lat, radius*1000) // Convert km to meters
		argCount += 2
	}

	// Date range filters
	if createdFrom != "" {
		argCount++
		whereClause.WriteString(fmt.Sprintf(" AND created_at >= $%d", argCount))
		args = append(args, createdFrom)
	}
	if createdTill != "" {
		argCount++
		whereClause.WriteString(fmt.Sprintf(" AND created_at <= $%d", argCount))
		args = append(args, createdTill)
	}

	// Entity type filter
	if len(entityTypes) > 0 {
		argCount++
		whereClause.WriteString(fmt.Sprintf(" AND entity_type = ANY($%d)", argCount))
		// Convert EntityType slice to string slice for PostgreSQL
		entityTypeStrings := make([]string, len(entityTypes))
		for i, entityType := range entityTypes {
			entityTypeStrings[i] = entityType.String()
		}
		args = append(args, entityTypeStrings)
	}

	query := fmt.Sprintf(`SELECT 
		id, entity_type, likes_count, dislikes_count, comments_count, 
		messages_count, shared_count, added_to_wishlist_count, added_to_basket_count, 
		visited_count, reported_count, follower_count, reviews_count, rating_count,
		videos_count, images_count, rating, review, category, category_id, category_slug,
		price, lat, lng, created_at, updated_at
	FROM %s %s 
	ORDER BY %s %s 
	LIMIT 100`, "%s", whereClause.String(), metricColumn, sortOrder)

	rows, err := r.db.QueryContext(ctx, r.table(query), args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying for metrics by type")
	}
	defer rows.Close()

	var metrics []*models.ItemMetric
	for rows.Next() {
		metric := &models.ItemMetric{}
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&metric.ID, &metric.EntityType, &metric.LikesCount, &metric.DislikesCount, &metric.CommentsCount,
			&metric.MessagesCount, &metric.SharedCount, &metric.AddedToWishlistCount, &metric.AddedToBasketCount,
			&metric.VisitedCount, &metric.ReportedCount, &metric.FollowerCount, &metric.ReviewsCount, &metric.RatingCount,
			&metric.VideosCount, &metric.ImagesCount, &metric.Rating, &metric.Review, &metric.Category,
			&metric.CategoryID, &metric.CategorySlug, &metric.Price, &metric.Lat, &metric.Lng, &createdAt, &updatedAt,
		)

		if err != nil {
			return nil, errors.Wrap(err, "scanning metric row")
		}

		metric.CreatedAt = createdAt.Format(time.RFC3339)
		metric.UpdatedAt = updatedAt.Format(time.RFC3339)

		metrics = append(metrics, metric)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "iterating metric rows")
	}

	return metrics, nil
}
