package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type ItemMetricRepository interface {
	GetItemMetric(ctx context.Context, itemId string) (*models.ItemMetric, error)
	GetItemsMetric(ctx context.Context, itemIds []string) ([]*models.ItemMetric, error)
	GetHighestMetricsByType(ctx context.Context, metricType string, req MetricSortRequest) ([]*models.ItemMetric, error)
	GetLowestMetricsByType(ctx context.Context, metricType string, req MetricSortRequest) ([]*models.ItemMetric, error)
}

