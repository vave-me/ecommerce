package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"middleman/internal/postgres"
	"middleman/metrics/internal/application"
	"middleman/metrics/internal/models"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"
)

type UserMetricCacheRepository struct {
	tableName string
	db        postgres.DB
}

var _ application.UserMetricCacheRepository = (*UserMetricCacheRepository)(nil)

func NewUserMetricCacheRepository(tableName string, db postgres.DB) UserMetricCacheRepository {
	return UserMetricCacheRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r UserMetricCacheRepository) AddUserMetric(ctx context.Context, userId, entityType string) error {
	const query = `INSERT INTO %s (
		id, entity_type, likes_count, dislikes_count, comments_count, 
		messages_count, shared_count, added_to_wishlist_count, added_to_basket_count, 
		visited_count, reported_count, follower_count, reviews_count, rating_count,
		videos_count, images_count, rating, review, category, category_id, category_slug,
		media_added_count, comment_added_count, liked_added_count, products_added_count, 
		videos_added_count, services_added_count, jobs_added_count, posts_added_count, 
		vehicles_added_count, properties_added_count, created_at, updated_at
	) VALUES ($1, $2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, '', '', '', 
	         0, 0, 0, 0, 0, 0, 0, 0, 0, 0, NOW(), NOW())`

	_, err := r.db.ExecContext(ctx, r.table(query), userId, entityType)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			// Ignore unique violation errors (user already exists)
			return nil
		}
		return errors.Wrap(err, "adding user metric")
	}

	return nil
}

func (r UserMetricCacheRepository) GetUserMetric(ctx context.Context, userId string) (*models.UserMetric, error) {
	const query = `SELECT 
		id, entity_type, likes_count, dislikes_count, comments_count, 
		messages_count, shared_count, added_to_wishlist_count, added_to_basket_count, 
		visited_count, reported_count, follower_count, reviews_count, rating_count,
		videos_count, images_count, rating, review, category, category_id, category_slug,
		media_added_count, comment_added_count, liked_added_count, products_added_count, 
		videos_added_count, services_added_count, jobs_added_count, posts_added_count, 
		vehicles_added_count, properties_added_count, created_at, updated_at
	FROM %s WHERE id = $1 LIMIT 1`

	metric := &models.UserMetric{
		ID: userId,
	}

	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, r.table(query), userId).Scan(
		&metric.ID, &metric.EntityType, &metric.LikesCount, &metric.DislikesCount, &metric.CommentsCount,
		&metric.MessagesCount, &metric.SharedCount, &metric.AddedToWishlistCount, &metric.AddedToBasketCount,
		&metric.VisitedCount, &metric.ReportedCount, &metric.FollowerCount, &metric.ReviewsCount, &metric.RatingCount,
		&metric.VideosCount, &metric.ImagesCount, &metric.Rating, &metric.Review, &metric.Category,
		&metric.CategoryID, &metric.CategorySlug, &metric.MediaAddedCount, &metric.CommentAddedCount,
		&metric.LikedAddedCount, &metric.ProductsAddedCount, &metric.VideosAddedCount,
		&metric.ServicesAddedCount, &metric.JobsAddedCount, &metric.PostsAddedCount,
		&metric.VehiclesAddedCount, &metric.PropertiesAddedCount, &createdAt, &updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msgf("user metric with ID %s not found", userId)
		}
		return nil, errors.Wrap(err, "getting user metric")
	}

	metric.CreatedAt = createdAt.Format(time.RFC3339)
	metric.UpdatedAt = updatedAt.Format(time.RFC3339)

	return metric, nil
}

func (r UserMetricCacheRepository) RemoveUserMetric(ctx context.Context, userId string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), userId)
	if err != nil {
		return errors.Wrap(err, "removing user metric")
	}

	return nil
}

func (r UserMetricCacheRepository) UpdateUserMetric(ctx context.Context, userId, metricType, metricTypeAction string) error {
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
	case models.MetricTypeUserProductAdd.String():
		query = `UPDATE %s SET products_added_count = products_added_count + 1, updated_at = NOW() WHERE id = $1`
	case models.MetricTypeUserPostAdd.String():
		query = `UPDATE %s SET posts_added_count = posts_added_count + 1, updated_at = NOW() WHERE id = $1`
	case models.MetricTypeUserServiceAdd.String():
		query = `UPDATE %s SET services_added_count = services_added_count + 1, updated_at = NOW() WHERE id = $1`
	case models.MetricTypeUserJobAdd.String():
		query = `UPDATE %s SET jobs_added_count = jobs_added_count + 1, updated_at = NOW() WHERE id = $1`
	case models.MetricTypeUserVehicleAdd.String():
		query = `UPDATE %s SET vehicles_added_count = vehicles_added_count + 1, updated_at = NOW() WHERE id = $1`
	case models.MetricTypeUserPropertyAdd.String():
		query = `UPDATE %s SET properties_added_count = properties_added_count + 1, updated_at = NOW() WHERE id = $1`
	default:
		return nil
	}

	_, err := r.db.ExecContext(ctx, r.table(query), userId)
	if err != nil {
		return errors.Wrap(err, "updating user metric")
	}

	return nil
}

func (r UserMetricCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
