package application

import (
	"context"
	"middleman/metrics/internal/models"
)

type (
	UpdateItemMetric struct {
		ItemID           string
		MetricType       models.MetricTypeCount
		MetricTypeAction models.MetricTypeAction
	}
	ShareItem struct {
		ItemID string
	}
	VisitItem struct {
		ItemID string
	}
	UpdateUserMetric struct {
		UserID           string
		MetricType       models.MetricTypeCount
		MetricTypeAction models.MetricTypeAction
	}
	GetUserMetric struct {
		UserID string
	}
	GetItemMetric struct {
		ItemID string
	}

	GetHighestMetricsByType struct {
		MetricType  models.MetricTypeCount
		EntityTypes []string
		CategoryID  string
		Lat         float64
		Lng         float64
		Radius      float64
		MinPrice    int64
		MaxPrice    int64
		CreatedFrom string
		CreatedTill string
	}
	GetLowestMetricsByType struct {
		MetricType  models.MetricTypeCount
		EntityTypes []string
		CategoryID  string
		Lat         float64
		Lng         float64
		Radius      float64
		MinPrice    int64
		MaxPrice    int64
		CreatedFrom string
		CreatedTill string
	}

	GetItemMetrics struct {
		ItemIDs []string
	}
	Application interface {
		UpdateItemMetric(ctx context.Context, update UpdateItemMetric) error
		UpdateUserMetric(ctx context.Context, update UpdateUserMetric) error
		GetUserMetric(ctx context.Context, getUserMetric GetUserMetric) (*models.UserMetric, error)
		GetItemMetric(ctx context.Context, getItemMetric GetItemMetric) (*models.ItemMetric, error)
		GetItemMetrics(ctx context.Context, getItemMetric GetItemMetrics) ([]*models.ItemMetric, error)
		GetHighestMetricsByType(ctx context.Context, getLowestMetricsByType GetHighestMetricsByType) ([]*models.ItemMetric, error)
		GetLowestMetricsByType(ctx context.Context, getLowestMetricsByType GetLowestMetricsByType) ([]*models.ItemMetric, error)
	}

	app struct {
		itemMetrics ItemMetricRepository
		userMetrics UserMetricRepository
	}
)

var _ Application = (*app)(nil)

func New(
	itemMetrics ItemMetricRepository,
	userMetrics UserMetricRepository,

) *app {
	return &app{
		itemMetrics: itemMetrics,
		userMetrics: userMetrics,
	}
}
func (a app) UpdateItemMetric(ctx context.Context, update UpdateItemMetric) error {
	return a.itemMetrics.UpdateItemMetric(ctx, update.ItemID, update.MetricType.String(), update.MetricTypeAction.String())
}

func (a app) UpdateUserMetric(ctx context.Context, update UpdateUserMetric) error {
	return a.userMetrics.UpdateUserMetric(ctx, update.UserID, update.MetricType.String(), update.MetricTypeAction.String())
}

func (a app) GetUserMetric(ctx context.Context, getUserMetric GetUserMetric) (*models.UserMetric, error) {
	return a.userMetrics.GetUserMetric(ctx, getUserMetric.UserID)
}
func (a app) GetItemMetric(ctx context.Context, getItemMetric GetItemMetric) (*models.ItemMetric, error) {
	return a.itemMetrics.GetItemMetric(ctx, getItemMetric.ItemID)
}
func (a app) GetItemMetrics(ctx context.Context, getItemMetrics GetItemMetrics) ([]*models.ItemMetric, error) {
	return a.itemMetrics.GetItemsMetrics(ctx, getItemMetrics.ItemIDs)
}
func (a app) GetHighestMetricsByType(ctx context.Context, getHighestItemMetricsByType GetHighestMetricsByType) ([]*models.ItemMetric, error) {
	// Convert string entity types to models.EntityType
	entityTypes := make([]models.EntityType, len(getHighestItemMetricsByType.EntityTypes))
	for i, entityTypeStr := range getHighestItemMetricsByType.EntityTypes {
		entityTypes[i] = models.ToEntityType(entityTypeStr)
	}

	return a.itemMetrics.GetHighestMetricsByType(ctx, getHighestItemMetricsByType.MetricType.String(), entityTypes, getHighestItemMetricsByType.CategoryID, getHighestItemMetricsByType.Lat, getHighestItemMetricsByType.Lng, getHighestItemMetricsByType.Radius, getHighestItemMetricsByType.MinPrice, getHighestItemMetricsByType.MaxPrice, getHighestItemMetricsByType.CreatedFrom, getHighestItemMetricsByType.CreatedTill)
}

func (a app) GetLowestMetricsByType(ctx context.Context, getLowestItemMetricsByType GetLowestMetricsByType) ([]*models.ItemMetric, error) {
	// Convert string entity types to models.EntityType
	entityTypes := make([]models.EntityType, len(getLowestItemMetricsByType.EntityTypes))
	for i, entityTypeStr := range getLowestItemMetricsByType.EntityTypes {
		entityTypes[i] = models.ToEntityType(entityTypeStr)
	}

	return a.itemMetrics.GetLowestMetricsByType(ctx, getLowestItemMetricsByType.MetricType.String(), entityTypes, getLowestItemMetricsByType.CategoryID, getLowestItemMetricsByType.Lat, getLowestItemMetricsByType.Lng, getLowestItemMetricsByType.Radius, getLowestItemMetricsByType.MinPrice, getLowestItemMetricsByType.MaxPrice, getLowestItemMetricsByType.CreatedFrom, getLowestItemMetricsByType.CreatedTill)
}
