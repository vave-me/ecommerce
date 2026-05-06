package application

import (
	"context"
	"middleman/metrics/internal/models"
)

type ItemMetricRepository interface {
	AddMetric(
		ctx context.Context,
		itemId, entityType string,
		categoryId string,
		price int64, lat, lng float64,
	) error
	GetItemMetric(ctx context.Context, itemIds string) (*models.ItemMetric, error)
	GetItemsMetrics(ctx context.Context, itemId []string) ([]*models.ItemMetric, error)
	GetHighestMetricsByType(ctx context.Context, metricType string, entityTypes []models.EntityType, categoryId string, lat, lng, radius float64, minPrice, maxPrice int64, createdFrom, CreatedTill string) ([]*models.ItemMetric, error)
	GetLowestMetricsByType(ctx context.Context, metricType string, entityTypes []models.EntityType, categoryId string, lat, lng, radius float64, minPrice, maxPrice int64, createdFrom, CreatedTill string) ([]*models.ItemMetric, error)
	RemoveItemMetric(ctx context.Context, itemId string) error
	UpdateItemMetric(ctx context.Context, itemId, metricType, metricTypeAction string) error
}

type ItemMetricCacheRepository interface {
	ItemMetricRepository
}
