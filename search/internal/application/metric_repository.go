package application

import (
	"context"
	"middleman/search/internal/models"
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
	GetItemMetric(ctx context.Context, itemID string) (*models.ItemMetric, error)
	GetItemsMetric(ctx context.Context, itemIDs []string) ([]*models.ItemMetric, error)
	// Advanced sorting methods
	GetHighestMetricsByType(ctx context.Context, metricType string, req MetricSortRequest) ([]*models.ItemMetric, error)
	GetLowestMetricsByType(ctx context.Context, metricType string, req MetricSortRequest) ([]*models.ItemMetric, error)
}

type MetricCacheRepository interface {
	MetricRepository
}
