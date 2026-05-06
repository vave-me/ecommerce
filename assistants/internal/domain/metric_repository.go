package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type MetricSortRequest struct {
	EntityTypes []string
	CategoryId  string
	MinPrice    int64
	MaxPrice    int64
	Limit       int32
	Lat         float32
	Lng         float32
	Radius      float32
	CreatedFrom string
	CreatedTo   string
}

type MetricRepository interface {
	// Core gRPC methods from protobuf specification
	UpdateItemMetric(ctx context.Context, itemID, metricType, metricTypeAction string) (*models.UpdateItemMetricResponse, error)
	ShareItem(ctx context.Context, itemID string) (*models.ShareItemResponse, error)
	VisitItem(ctx context.Context, itemID string) (*models.VisitItemResponse, error)
	UpdateUserMetric(ctx context.Context, userID, metricType, metricTypeAction string) (*models.UpdateUserMetricResponse, error)
	GetUserMetric(ctx context.Context, userID string) (*models.GetUserMetricResponse, error)
	GetItemMetric(ctx context.Context, itemID string) (*models.GetItemMetricResponse, error)
	GetItemsMetric(ctx context.Context, itemIDs []string, limit int32) (*models.GetItemsMetricResponse, error)
	GetHighestMetricsByType(ctx context.Context, metricType string, req MetricSortRequest) (*models.GetItemsMetricResponse, error)
	GetLowestMetricsByType(ctx context.Context, metricType string, req MetricSortRequest) (*models.GetItemsMetricResponse, error)

	// Extended methods for AI tooling and comprehensive metrics management
	GetItemMetricByType(ctx context.Context, itemID, metricType string) (*models.ItemMetric, error)
	GetUserMetricByType(ctx context.Context, userID, metricType string) (*models.UserMetric, error)
	GetItemMetricsByCategory(ctx context.Context, categoryID string, limit int32) ([]*models.ItemMetric, error)
	GetUserMetricsByCategory(ctx context.Context, userID, categoryID string) (*models.UserMetric, error)
	GetTopItemsByMetric(ctx context.Context, metricType string, entityTypes []string, limit int32) ([]*models.ItemMetric, error)
	GetTopUsersByMetric(ctx context.Context, metricType string, limit int32) ([]*models.UserMetric, error)
	GetItemMetricsInRange(ctx context.Context, lat, lng, radius float32, limit int32) ([]*models.ItemMetric, error)
	GetMetricsSummary(ctx context.Context, entityType string) (map[string]int64, error)
	SearchItemMetrics(ctx context.Context, query string, entityTypes []string, limit int32) ([]*models.ItemMetric, error)
	SearchUserMetrics(ctx context.Context, query string, limit int32) ([]*models.UserMetric, error)
	GetRecentlyUpdatedMetrics(ctx context.Context, entityType string, limit int32) ([]*models.ItemMetric, error)
	GetMetricsByPriceRange(ctx context.Context, minPrice, maxPrice int64, entityTypes []string, limit int32) ([]*models.ItemMetric, error)
	GetMetricsByRatingRange(ctx context.Context, minRating, maxRating int64, entityTypes []string, limit int32) ([]*models.ItemMetric, error)
	GetTrendingItems(ctx context.Context, entityTypes []string, days int32, limit int32) ([]*models.ItemMetric, error)
	GetActiveUsers(ctx context.Context, days int32, limit int32) ([]*models.UserMetric, error)
	CompareItemMetrics(ctx context.Context, itemID1, itemID2 string) (map[string]interface{}, error)
	CompareUserMetrics(ctx context.Context, userID1, userID2 string) (map[string]interface{}, error)
	GetMetricsAnalytics(ctx context.Context, entityType string, timeRange string) (map[string]interface{}, error)
	BulkUpdateItemMetrics(ctx context.Context, updates []map[string]interface{}) error
	BulkUpdateUserMetrics(ctx context.Context, updates []map[string]interface{}) error
	ResetItemMetrics(ctx context.Context, itemID string, metricTypes []string) error
	ResetUserMetrics(ctx context.Context, userID string, metricTypes []string) error
}

type MetricCacheRepository interface {
	MetricRepository
}
