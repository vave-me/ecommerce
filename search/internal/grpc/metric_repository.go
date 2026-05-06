// File: search/internal/grpc/metric_repository.go
package grpc

import (
	"context"
	"sync"

	"google.golang.org/grpc"

	"middleman/internal/rpc"
	"middleman/metrics/metricspb"
	"middleman/search/internal/application"
	"middleman/search/internal/models"
)

// MetricRepository calls the remote metrics service (gRPC).
type MetricRepository struct {
	endpoint string
	conn     *grpc.ClientConn
	client   metricspb.MetricsServiceClient
	mu       sync.RWMutex
	cache    map[string]*models.ItemMetric
	connOnce sync.Once
}

var _ application.ItemMetricRepository = (*MetricRepository)(nil)

// NewMetricRepository instantiates the gRPC-based metrics repo.
func NewMetricRepository(endpoint string) *MetricRepository {
	return &MetricRepository{
		endpoint: endpoint,
		cache:    make(map[string]*models.ItemMetric, 1000),
	}
}

// GetItemMetric retrieves metrics by item ID from the metrics microservice (via gRPC).
func (r *MetricRepository) GetItemMetric(ctx context.Context, itemID string) (*models.ItemMetric, error) {
	// Check cache first
	r.mu.RLock()
	if metric, found := r.cache[itemID]; found {
		r.mu.RUnlock()
		return metric, nil
	}
	r.mu.RUnlock()

	// Set up connection if not already established
	err := r.ensureConnection(ctx)
	if err != nil {
		return r.createDefaultMetric(itemID), nil
	}

	// Make the RPC call
	resp, err := r.client.GetItemMetric(ctx, &metricspb.GetItemMetricRequest{ItemId: itemID})
	if err != nil {
		return r.createDefaultMetric(itemID), nil
	}

	// Convert from protobuf format to internal model
	itemMetric := &models.ItemMetric{
		ID:                   resp.Metric.GetItemId(),
		EntityType:           resp.Metric.GetEntityType(),
		LikesCount:           resp.Metric.GetLikesCount(),
		DislikesCount:        resp.Metric.GetDislikesCount(),
		CommentsCount:        resp.Metric.GetCommentsCount(),
		SharedCount:          resp.Metric.GetSharedCount(),
		AddedToWishlistCount: resp.Metric.GetAddedToWishlistCount(),
		AddedToBasketCount:   resp.Metric.GetAddedToBasketCount(),
		VisitedCount:         resp.Metric.GetVisitedCount(),
		ReportedCount:        resp.Metric.GetReportedCount(),
		FollowerCount:        resp.Metric.GetFollowerCount(),
		ReviewsCount:         resp.Metric.GetReviewCount(),
		RatingCount:          resp.Metric.GetRatingCount(),
		VideosCount:          resp.Metric.GetVideosCount(),
		ImagesCount:          resp.Metric.GetImagesCount(),
		Rating:               resp.Metric.GetRating(),
		Category:             resp.Metric.GetCategory(),
		CategoryID:           resp.Metric.GetCategoryId(),
		CategorySlug:         resp.Metric.GetCategorySlug(),
	}

	// Cache the result
	r.mu.Lock()
	r.cache[itemID] = itemMetric
	r.mu.Unlock()

	return itemMetric, nil
}

// GetItemsMetric retrieves metrics for multiple items at once (batch operation)
// Limited to 150 items per request
func (r *MetricRepository) GetItemsMetric(ctx context.Context, itemIDs []string) ([]*models.ItemMetric, error) {
	if len(itemIDs) == 0 {
		return make([]*models.ItemMetric, 0), nil
	}

	// Limit to 150 items max per request
	if len(itemIDs) > 150 {
		itemIDs = itemIDs[:150]
	}

	// Check cache first for all requested items
	result := make([]*models.ItemMetric, 0, len(itemIDs))
	uncachedIDs := make([]string, 0, len(itemIDs))
	cachedMetricsByID := make(map[string]*models.ItemMetric)

	r.mu.RLock()
	for _, id := range itemIDs {
		if metric, found := r.cache[id]; found {
			cachedMetricsByID[id] = metric
		} else {
			uncachedIDs = append(uncachedIDs, id)
		}
	}
	r.mu.RUnlock()

	// If all items were in cache, return immediately
	if len(uncachedIDs) == 0 {
		// Convert map to ordered array matching request order
		for _, id := range itemIDs {
			if metric, found := cachedMetricsByID[id]; found {
				result = append(result, metric)
			}
		}
		return result, nil
	}

	// Set up connection if not already established
	err := r.ensureConnection(ctx)
	if err != nil {
		// Return default metrics for uncached items
		for _, id := range itemIDs {
			if metric, found := cachedMetricsByID[id]; found {
				result = append(result, metric)
			} else {
				result = append(result, r.createDefaultMetric(id))
			}
		}
		return result, nil
	}

	// Make the RPC call for batch request
	resp, err := r.client.GetItemsMetric(ctx, &metricspb.GetItemsMetricRequest{ItemIds: uncachedIDs})
	if err != nil {
		// Return default metrics for uncached items if request failed
		for _, id := range itemIDs {
			if metric, found := cachedMetricsByID[id]; found {
				result = append(result, metric)
			} else {
				result = append(result, r.createDefaultMetric(id))
			}
		}
		return result, nil
	}

	// Process the response - create map for fast lookup
	fetchedMetrics := make(map[string]*models.ItemMetric)

	r.mu.Lock()
	// Convert from protobuf format to internal models and update cache
	for _, pbMetric := range resp.Metrics {
		if pbMetric == nil {
			continue
		}

		itemMetric := &models.ItemMetric{
			ID:                   pbMetric.GetItemId(),
			EntityType:           pbMetric.GetEntityType(),
			LikesCount:           pbMetric.GetLikesCount(),
			DislikesCount:        pbMetric.GetDislikesCount(),
			CommentsCount:        pbMetric.GetCommentsCount(),
			SharedCount:          pbMetric.GetSharedCount(),
			AddedToWishlistCount: pbMetric.GetAddedToWishlistCount(),
			AddedToBasketCount:   pbMetric.GetAddedToBasketCount(),
			VisitedCount:         pbMetric.GetVisitedCount(),
			ReportedCount:        pbMetric.GetReportedCount(),
			FollowerCount:        pbMetric.GetFollowerCount(),
			ReviewsCount:         pbMetric.GetReviewCount(),
			RatingCount:          pbMetric.GetRatingCount(),
			VideosCount:          pbMetric.GetVideosCount(),
			ImagesCount:          pbMetric.GetImagesCount(),
			Rating:               pbMetric.GetRating(),
			Category:             pbMetric.GetCategory(),
			CategoryID:           pbMetric.GetCategoryId(),
			CategorySlug:         pbMetric.GetCategorySlug(),
		}

		// Add to temp map and cache
		fetchedMetrics[itemMetric.ID] = itemMetric
		r.cache[itemMetric.ID] = itemMetric
	}
	r.mu.Unlock()

	// Create final result array in original item order
	for _, id := range itemIDs {
		if metric, found := cachedMetricsByID[id]; found {
			// Use cached metric
			result = append(result, metric)
		} else if metric, found := fetchedMetrics[id]; found {
			// Use newly fetched metric
			result = append(result, metric)
		} else {
			// Create default metric for missing items
			result = append(result, r.createDefaultMetric(id))
		}
	}

	return result, nil
}

// ensureConnection makes sure we have an active connection to the metrics service
func (r *MetricRepository) ensureConnection(ctx context.Context) error {
	var err error
	r.connOnce.Do(func() {
		r.conn, err = rpc.Dial(ctx, r.endpoint)
		if err == nil && r.conn != nil {
			r.client = metricspb.NewMetricsServiceClient(r.conn)
		}
	})
	return err
}

// createDefaultMetric generates a default metric object when the service fails
func (r *MetricRepository) createDefaultMetric(itemID string) *models.ItemMetric {
	defaultMetric := &models.ItemMetric{
		ID:                   itemID,
		CommentsCount:        0,
		MessagesCount:        0,
		DislikesCount:        0,
		LikesCount:           0,
		AddedToBasketCount:   0,
		AddedToWishlistCount: 0,
		SharedCount:          0,
		VisitedCount:         0,
		Rating:               4, // Reasonable default rating
		RatingCount:          0,
	}

	// Cache the default metric
	r.mu.Lock()
	r.cache[itemID] = defaultMetric
	r.mu.Unlock()

	return defaultMetric
}

// GetHighestMetricsByType retrieves top items sorted by specific metric type
func (r *MetricRepository) GetHighestMetricsByType(ctx context.Context, metricType string, req application.MetricSortRequest) ([]*models.ItemMetric, error) {
	err := r.ensureConnection(ctx)
	if err != nil {
		return []*models.ItemMetric{}, nil
	}

	pbReq := &metricspb.GetHighestMetricsByTypeRequest{
		MetricType:  metricType,
		EntityTypes: req.EntityTypes,
		CategoryId:  req.CategoryId,
		MinPrice:    req.MinPrice,
		MaxPrice:    req.MaxPrice,
		Limit:       req.Limit,
		Lat:         req.Lat,
		Lng:         req.Lng,
		Radius:      req.Radius,
		CreatedFrom: req.CreatedFrom,
		CreatedTo:   req.CreatedTo,
	}

	resp, err := r.client.GetHighestMetricsByType(ctx, pbReq)
	if err != nil {
		return []*models.ItemMetric{}, nil
	}

	return r.convertMetricsResponse(resp.Metrics), nil
}

// GetLowestMetricsByType retrieves bottom items sorted by specific metric type
func (r *MetricRepository) GetLowestMetricsByType(ctx context.Context, metricType string, req application.MetricSortRequest) ([]*models.ItemMetric, error) {
	err := r.ensureConnection(ctx)
	if err != nil {
		return []*models.ItemMetric{}, nil
	}

	pbReq := &metricspb.GetLowestMetricsByTypeRequest{
		MetricType:  metricType,
		EntityTypes: req.EntityTypes,
		CategoryId:  req.CategoryId,
		MinPrice:    req.MinPrice,
		MaxPrice:    req.MaxPrice,
		Limit:       req.Limit,
		Lat:         req.Lat,
		Lng:         req.Lng,
		Radius:      req.Radius,
		CreatedFrom: req.CreatedFrom,
		CreatedTo:   req.CreatedTo,
	}

	resp, err := r.client.GetLowestMetricsByType(ctx, pbReq)
	if err != nil {
		return []*models.ItemMetric{}, nil
	}

	return r.convertMetricsResponse(resp.Metrics), nil
}

// convertMetricsResponse converts protobuf metrics to internal models
func (r *MetricRepository) convertMetricsResponse(pbMetrics []*metricspb.ItemMetric) []*models.ItemMetric {
	metrics := make([]*models.ItemMetric, 0, len(pbMetrics))

	r.mu.Lock()
	for _, pbMetric := range pbMetrics {
		if pbMetric == nil {
			continue
		}

		itemMetric := &models.ItemMetric{
			ID:                   pbMetric.GetItemId(),
			EntityType:           pbMetric.GetEntityType(),
			LikesCount:           pbMetric.GetLikesCount(),
			DislikesCount:        pbMetric.GetDislikesCount(),
			CommentsCount:        pbMetric.GetCommentsCount(),
			SharedCount:          pbMetric.GetSharedCount(),
			AddedToWishlistCount: pbMetric.GetAddedToWishlistCount(),
			AddedToBasketCount:   pbMetric.GetAddedToBasketCount(),
			VisitedCount:         pbMetric.GetVisitedCount(),
			ReportedCount:        pbMetric.GetReportedCount(),
			FollowerCount:        pbMetric.GetFollowerCount(),
			ReviewsCount:         pbMetric.GetReviewCount(),
			RatingCount:          pbMetric.GetRatingCount(),
			VideosCount:          pbMetric.GetVideosCount(),
			ImagesCount:          pbMetric.GetImagesCount(),
			Rating:               pbMetric.GetRating(),
			Category:             pbMetric.GetCategory(),
			CategoryID:           pbMetric.GetCategoryId(),
			CategorySlug:         pbMetric.GetCategorySlug(),
		}

		metrics = append(metrics, itemMetric)
		r.cache[itemMetric.ID] = itemMetric
	}
	r.mu.Unlock()

	return metrics
}
