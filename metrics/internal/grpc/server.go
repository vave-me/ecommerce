// File: search/internal/grpc/server.go
package grpc

import (
	"context"
	"fmt"
	"middleman/internal/errorsotel"
	"middleman/metrics/internal/application"
	"middleman/metrics/internal/models"
	"middleman/metrics/metricspb"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	app application.Application
	metricspb.UnimplementedMetricsServiceServer
}

func RegisterServer(
	ctx context.Context,
	app application.Application,
	registrar grpc.ServiceRegistrar,
) error {
	metricspb.RegisterMetricsServiceServer(registrar, server{app: app})
	return nil
}
func (s server) GetUserMetric(ctx context.Context, req *metricspb.GetUserMetricRequest) (*metricspb.GetUserMetricResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "Get user metric")

	span.SetAttributes(attribute.String("User ID", req.GetUserId()))

	userMetric, err := s.app.GetUserMetric(ctx, application.GetUserMetric{
		UserID: req.GetUserId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &metricspb.GetUserMetricResponse{
		Metric: s.userMetricFromDomain(userMetric),
	}, nil
}
func (s server) ShareItem(ctx context.Context, req *metricspb.ShareItemRequest) (*metricspb.ShareItemResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "Get user metric")

	span.SetAttributes(attribute.String("Item ID", req.GetItemId()))
	err := s.app.UpdateItemMetric(ctx, application.UpdateItemMetric{
		ItemID:           req.GetItemId(),
		MetricType:       models.MetricTypeCountShare,
		MetricTypeAction: models.MetricTypeActionAdd,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &metricspb.ShareItemResponse{}, nil
}
func (s server) VisitItem(ctx context.Context, req *metricspb.VisitItemRequest) (*metricspb.VisitItemResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "Get user metric")

	span.SetAttributes(attribute.String("Item ID", req.GetItemId()))
	err := s.app.UpdateItemMetric(ctx, application.UpdateItemMetric{
		ItemID:           req.GetItemId(),
		MetricType:       models.MetricTypeCountVisit,
		MetricTypeAction: models.MetricTypeActionAdd,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &metricspb.VisitItemResponse{}, nil
}

func (s server) GetItemMetric(ctx context.Context, req *metricspb.GetItemMetricRequest) (*metricspb.GetItemMetricResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "Get item metric")

	span.SetAttributes(attribute.String("Item ID", req.GetItemId()))

	metric, err := s.app.GetItemMetric(ctx, application.GetItemMetric{
		ItemID: req.GetItemId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &metricspb.GetItemMetricResponse{
		Metric: s.itemMetricFromDomain(metric),
	}, nil
}

func (s server) GetItemsMetric(ctx context.Context, req *metricspb.GetItemsMetricRequest) (*metricspb.GetItemsMetricResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "Get items metrics batch")

	span.SetAttributes(attribute.Int("Item Count", len(req.GetItemIds())))

	// Apply limit if specified, default max 150
	itemIDs := req.GetItemIds()
	limit := req.GetLimit()
	if limit <= 0 || limit > 150 {
		limit = 150
	}
	if int(limit) < len(itemIDs) {
		itemIDs = itemIDs[:limit]
	}

	metrics, err := s.app.GetItemMetrics(ctx, application.GetItemMetrics{
		ItemIDs: itemIDs,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Convert domain models to protobuf models
	protoMetrics := make([]*metricspb.ItemMetric, 0, len(metrics))
	for _, metric := range metrics {
		protoMetrics = append(protoMetrics, s.itemMetricFromDomain(metric))
	}

	return &metricspb.GetItemsMetricResponse{
		Metrics: protoMetrics,
	}, nil
}

func (s server) GetHighestMetricsByType(ctx context.Context, req *metricspb.GetHighestMetricsByTypeRequest) (*metricspb.GetItemsMetricResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "Get highest metrics by type")

	span.SetAttributes(
		attribute.String("Metric Type", req.GetMetricType()),
		attribute.String("Category ID", req.GetCategoryId()),
		attribute.Float64("Lat", float64(req.GetLat())),
		attribute.Float64("Lng", float64(req.GetLng())),
		attribute.Float64("Radius", float64(req.GetRadius())),
		attribute.Int64("Min Price", req.GetMinPrice()),
		attribute.Int64("Max Price", req.GetMaxPrice()),
	)

	metrics, err := s.app.GetHighestMetricsByType(ctx, application.GetHighestMetricsByType{
		MetricType:  models.ToMetricTypeCount(req.GetMetricType()),
		EntityTypes: req.GetEntityTypes(),
		CategoryID:  req.GetCategoryId(),
		Lat:         float64(req.GetLat()),
		Lng:         float64(req.GetLng()),
		Radius:      float64(req.GetRadius()),
		MinPrice:    req.GetMinPrice(),
		MaxPrice:    req.GetMaxPrice(),
		CreatedFrom: req.GetCreatedFrom(),
		CreatedTill: req.GetCreatedTo(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Apply limit if specified, default max 100
	limit := req.GetLimit()
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if int(limit) < len(metrics) {
		metrics = metrics[:limit]
	}

	// Convert domain models to protobuf models
	protoMetrics := make([]*metricspb.ItemMetric, 0, len(metrics))
	for _, metric := range metrics {
		protoMetrics = append(protoMetrics, s.itemMetricFromDomain(metric))
	}

	return &metricspb.GetItemsMetricResponse{
		Metrics: protoMetrics,
	}, nil
}

func (s server) GetLowestMetricsByType(ctx context.Context, req *metricspb.GetLowestMetricsByTypeRequest) (*metricspb.GetItemsMetricResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "Get lowest metrics by type")

	span.SetAttributes(
		attribute.String("Metric Type", req.GetMetricType()),
		attribute.String("Category ID", req.GetCategoryId()),
		attribute.Float64("Lat", float64(req.GetLat())),
		attribute.Float64("Lng", float64(req.GetLng())),
		attribute.Float64("Radius", float64(req.GetRadius())),
		attribute.Int64("Min Price", req.GetMinPrice()),
		attribute.Int64("Max Price", req.GetMaxPrice()),
	)

	metrics, err := s.app.GetLowestMetricsByType(ctx, application.GetLowestMetricsByType{
		MetricType:  models.ToMetricTypeCount(req.GetMetricType()),
		EntityTypes: req.GetEntityTypes(),
		CategoryID:  req.GetCategoryId(),
		Lat:         float64(req.GetLat()),
		Lng:         float64(req.GetLng()),
		Radius:      float64(req.GetRadius()),
		MinPrice:    req.GetMinPrice(),
		MaxPrice:    req.GetMaxPrice(),
		CreatedFrom: req.GetCreatedFrom(),
		CreatedTill: req.GetCreatedTo(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Apply limit if specified, default max 100
	limit := req.GetLimit()
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	if int(limit) < len(metrics) {
		metrics = metrics[:limit]
	}

	// Convert domain models to protobuf models
	protoMetrics := make([]*metricspb.ItemMetric, 0, len(metrics))
	for _, metric := range metrics {
		protoMetrics = append(protoMetrics, s.itemMetricFromDomain(metric))
	}

	return &metricspb.GetItemsMetricResponse{
		Metrics: protoMetrics,
	}, nil
}

func (s server) UpdateItemMetric(ctx context.Context, req *metricspb.UpdateItemMetricRequest) (*metricspb.UpdateItemMetricResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "Update item metric")

	span.SetAttributes(attribute.String("Item ID", req.GetItemId()))

	err := s.app.UpdateItemMetric(ctx, application.UpdateItemMetric{
		ItemID:           req.GetItemId(),
		MetricType:       models.MetricTypeCount(req.GetMetricType()),
		MetricTypeAction: models.MetricTypeAction(req.GetMetricTypeAction()),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &metricspb.UpdateItemMetricResponse{
		ItemId: req.GetItemId(),
	}, nil
}

func (s server) UpdateUserMetric(ctx context.Context, req *metricspb.UpdateUserMetricRequest) (*metricspb.UpdateUserMetricResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "Update user metric")

	span.SetAttributes(attribute.String("User ID", req.GetUserId()))

	err := s.app.UpdateUserMetric(ctx, application.UpdateUserMetric{
		UserID:           req.GetUserId(),
		MetricType:       models.MetricTypeCount(req.GetMetricType()),
		MetricTypeAction: models.MetricTypeAction(req.GetMetricTypeAction()),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &metricspb.UpdateUserMetricResponse{
		UserId: req.GetUserId(),
	}, nil
}

func handlePanic(span trace.Span, methodName string) {
	if r := recover(); r != nil {
		err := fmt.Errorf("panic recovered in %s: %v", methodName, r)
		if span != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, "panic")
		}
		log.Printf("Panic recovered in %s: %v", methodName, r)
	}
}

func (s server) userMetricFromDomain(userMetric *models.UserMetric) *metricspb.UserMetric {
	if userMetric == nil {
		return nil
	}

	return &metricspb.UserMetric{
		UserId:               userMetric.ID,
		EntityType:           userMetric.EntityType,
		LikesCount:           userMetric.LikesCount,
		DislikesCount:        userMetric.DislikesCount,
		CommentsCount:        userMetric.CommentsCount,
		MessagesCount:        userMetric.MessagesCount,
		SharedCount:          userMetric.SharedCount,
		AddedToWishlistCount: userMetric.AddedToWishlistCount,
		AddedToBasketCount:   userMetric.AddedToBasketCount,
		VisitedCount:         userMetric.VisitedCount,
		ReportedCount:        userMetric.ReportedCount,
		FollowerCount:        userMetric.FollowerCount,
		ReviewCount:          userMetric.ReviewsCount,
		RatingCount:          userMetric.RatingCount,
		VideosCount:          userMetric.VideosCount,
		ImagesCount:          userMetric.ImagesCount,
		Rating:               userMetric.Rating,
		Category:             userMetric.Category,
		CategoryId:           userMetric.CategoryID,
		CategorySlug:         userMetric.CategorySlug,
		MediaAddedCount:      userMetric.MediaAddedCount,
		CommentAddedCount:    userMetric.CommentAddedCount,
		LikedCount:           userMetric.LikedAddedCount,
		DislikedCount:        userMetric.DislikesCount,
		ProductsAddedCount:   userMetric.ProductsAddedCount,
		VideosAddedCount:     userMetric.VideosAddedCount,
		ImagesAddedCount:     userMetric.ImagesCount,
		SeriesAddedCount:     userMetric.ServicesAddedCount,
		JobsAddedCount:       userMetric.JobsAddedCount,
		PostsAddedCount:      userMetric.PostsAddedCount,
		VehiclesAddedCount:   userMetric.VehiclesAddedCount,
		PropertiesAddedCount: userMetric.PropertiesAddedCount,
	}
}

func (s server) itemMetricFromDomain(itemMetric *models.ItemMetric) *metricspb.ItemMetric {
	if itemMetric == nil {
		return nil
	}

	return &metricspb.ItemMetric{
		ItemId:               itemMetric.ID,
		EntityType:           itemMetric.EntityType,
		LikesCount:           itemMetric.LikesCount,
		DislikesCount:        itemMetric.DislikesCount,
		CommentsCount:        itemMetric.CommentsCount,
		MessagesCount:        itemMetric.MessagesCount,
		SharedCount:          itemMetric.SharedCount,
		AddedToWishlistCount: itemMetric.AddedToWishlistCount,
		AddedToBasketCount:   itemMetric.AddedToBasketCount,
		VisitedCount:         itemMetric.VisitedCount,
		ReportedCount:        itemMetric.ReportedCount,
		FollowerCount:        itemMetric.FollowerCount,
		ReviewCount:          itemMetric.ReviewsCount,
		RatingCount:          itemMetric.RatingCount,
		VideosCount:          itemMetric.VideosCount,
		ImagesCount:          itemMetric.ImagesCount,
		Rating:               itemMetric.Rating,
		Category:             itemMetric.Category,
		CategoryId:           itemMetric.CategoryID,
		CategorySlug:         itemMetric.CategorySlug,
		MediaCount:           itemMetric.VideosCount + itemMetric.ImagesCount,
		Price:                itemMetric.Price,
		Lat:                  float32(itemMetric.Lat),
		Lng:                  float32(itemMetric.Lng),
	}
}
