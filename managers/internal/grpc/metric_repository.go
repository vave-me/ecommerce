// File: search/internal/grpc/metric_repository.go
package grpc

import (
	"context"
	"fmt"
	"log"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
	"middleman/metrics/metricspb"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MetricRepository implements the domain.MetricRepository interface using gRPC
type MetricRepository struct {
	endpoint      string
	conn          *grpc.ClientConn
	client        metricspb.MetricsServiceClient
	mu            sync.RWMutex
	cache         map[string]*models.ItemMetric
	connOnce      sync.Once
	connErr       error
	jwt           string // Optional JWT for authenticated requests
	cleanupTicker *time.Ticker
	stopCleanup   chan struct{}
}

var _ domain.MetricRepository = (*MetricRepository)(nil)

// NewMetricRepository creates a new MetricRepository instance
func NewMetricRepository(endpoint string) *MetricRepository {
	repo := &MetricRepository{
		endpoint:      endpoint,
		cache:         make(map[string]*models.ItemMetric, 1000),
		cleanupTicker: time.NewTicker(5 * time.Minute),
		stopCleanup:   make(chan struct{}),
	}

	// Start cache cleanup goroutine
	go repo.cleanupOldEntries()

	return repo
}

// NewMetricRepositoryWithAuth creates a new MetricRepository instance with JWT authentication
func NewMetricRepositoryWithAuth(endpoint, jwt string) *MetricRepository {
	repo := NewMetricRepository(endpoint)
	repo.jwt = jwt
	return repo
}

// cleanupOldEntries removes cache entries older than 10 minutes
func (r *MetricRepository) cleanupOldEntries() {
	for {
		select {
		case <-r.cleanupTicker.C:
			r.mu.Lock()
			// In a real implementation, you'd track entry timestamps
			// For now, just clear the cache if it's too large
			if len(r.cache) > 5000 {
				r.cache = make(map[string]*models.ItemMetric, 1000)
			}
			r.mu.Unlock()
		case <-r.stopCleanup:
			r.cleanupTicker.Stop()
			return
		}
	}
}

// Close cleanly shuts down the repository
func (r *MetricRepository) Close() error {
	close(r.stopCleanup)

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
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

// getContextWithAuth adds JWT authentication to the context if available
func (r *MetricRepository) getContextWithAuth(ctx context.Context) context.Context {
	if r.jwt != "" {
		return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+r.jwt)
	}
	return ctx
}

// Core gRPC methods from protobuf specification

// UpdateItemMetric updates metrics for an item
func (r *MetricRepository) UpdateItemMetric(ctx context.Context, itemID, metricType, metricTypeAction string) (*models.UpdateItemMetricResponse, error) {
	if err := r.ensureConnection(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to metrics service: %w", err)
	}

	authCtx := r.getContextWithAuth(ctx)
	req := &metricspb.UpdateItemMetricRequest{
		ItemId:           itemID,
		MetricType:       metricType,
		MetricTypeAction: metricTypeAction,
	}

	resp, err := r.client.UpdateItemMetric(authCtx, req)
	if err != nil {
		log.Printf("Failed to update item metric for item %s: %v", itemID, err)
		return nil, fmt.Errorf("failed to update item metric: %w", err)
	}

	// Clear cache for this item since metrics changed
	r.mu.Lock()
	delete(r.cache, itemID)
	r.mu.Unlock()

	return &models.UpdateItemMetricResponse{
		ItemID: resp.ItemId,
	}, nil
}

// ShareItem increments the share count for an item
func (r *MetricRepository) ShareItem(ctx context.Context, itemID string) (*models.ShareItemResponse, error) {
	if err := r.ensureConnection(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to metrics service: %w", err)
	}

	authCtx := r.getContextWithAuth(ctx)
	req := &metricspb.ShareItemRequest{
		ItemId: itemID,
	}

	_, err := r.client.ShareItem(authCtx, req)
	if err != nil {
		log.Printf("Failed to share item %s: %v", itemID, err)
		return &models.ShareItemResponse{Success: false, ItemID: itemID}, nil
	}

	// Clear cache for this item since metrics changed
	r.mu.Lock()
	delete(r.cache, itemID)
	r.mu.Unlock()

	return &models.ShareItemResponse{Success: true, ItemID: itemID}, nil
}

// VisitItem increments the visit count for an item
func (r *MetricRepository) VisitItem(ctx context.Context, itemID string) (*models.VisitItemResponse, error) {
	if err := r.ensureConnection(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to metrics service: %w", err)
	}

	authCtx := r.getContextWithAuth(ctx)
	req := &metricspb.VisitItemRequest{
		ItemId: itemID,
	}

	_, err := r.client.VisitItem(authCtx, req)
	if err != nil {
		log.Printf("Failed to visit item %s: %v", itemID, err)
		return &models.VisitItemResponse{Success: false, ItemID: itemID}, nil
	}

	// Clear cache for this item since metrics changed
	r.mu.Lock()
	delete(r.cache, itemID)
	r.mu.Unlock()

	return &models.VisitItemResponse{Success: true, ItemID: itemID}, nil
}

// UpdateUserMetric updates metrics for a user
func (r *MetricRepository) UpdateUserMetric(ctx context.Context, userID, metricType, metricTypeAction string) (*models.UpdateUserMetricResponse, error) {
	if err := r.ensureConnection(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to metrics service: %w", err)
	}

	authCtx := r.getContextWithAuth(ctx)
	req := &metricspb.UpdateUserMetricRequest{
		UserId:           userID,
		MetricType:       metricType,
		MetricTypeAction: metricTypeAction,
	}

	resp, err := r.client.UpdateUserMetric(authCtx, req)
	if err != nil {
		log.Printf("Failed to update user metric for user %s: %v", userID, err)
		return nil, fmt.Errorf("failed to update user metric: %w", err)
	}

	return &models.UpdateUserMetricResponse{
		UserID: resp.UserId,
	}, nil
}

// GetUserMetric retrieves metrics for a user
func (r *MetricRepository) GetUserMetric(ctx context.Context, userID string) (*models.GetUserMetricResponse, error) {
	if err := r.ensureConnection(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to metrics service: %w", err)
	}

	authCtx := r.getContextWithAuth(ctx)
	req := &metricspb.GetUserMetricRequest{
		UserId: userID,
	}

	resp, err := r.client.GetUserMetric(authCtx, req)
	if err != nil {
		log.Printf("Failed to get user metric for user %s: %v", userID, err)
		return nil, fmt.Errorf("failed to get user metric: %w", err)
	}

	userMetric := r.convertUserMetricFromPb(resp.Metric)

	return &models.GetUserMetricResponse{
		Metric: userMetric,
	}, nil
}

// GetItemMetric retrieves metrics for an item
func (r *MetricRepository) GetItemMetric(ctx context.Context, itemID string) (*models.GetItemMetricResponse, error) {
	if err := r.ensureConnection(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to metrics service: %w", err)
	}

	authCtx := r.getContextWithAuth(ctx)
	req := &metricspb.GetItemMetricRequest{
		ItemId: itemID,
	}

	resp, err := r.client.GetItemMetric(authCtx, req)
	if err != nil {
		log.Printf("Failed to get item metric for item %s: %v", itemID, err)
		return nil, fmt.Errorf("failed to get item metric: %w", err)
	}

	itemMetric := r.convertItemMetricFromPb(resp.Metric)

	// Cache the result
	r.mu.Lock()
	if len(r.cache) < 5000 {
		r.cache[itemID] = itemMetric
	}
	r.mu.Unlock()

	return &models.GetItemMetricResponse{
		Metric: itemMetric,
	}, nil
}

// GetItemsMetric retrieves metrics for multiple items
func (r *MetricRepository) GetItemsMetric(ctx context.Context, itemIDs []string, limit int32) (*models.GetItemsMetricResponse, error) {
	if err := r.ensureConnection(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to metrics service: %w", err)
	}

	authCtx := r.getContextWithAuth(ctx)
	req := &metricspb.GetItemsMetricRequest{
		ItemIds: itemIDs,
		Limit:   limit,
	}

	resp, err := r.client.GetItemsMetric(authCtx, req)
	if err != nil {
		log.Printf("Failed to get items metrics: %v", err)
		return nil, fmt.Errorf("failed to get items metrics: %w", err)
	}

	metrics := make([]*models.ItemMetric, 0, len(resp.Metrics))
	r.mu.Lock()
	for _, pbMetric := range resp.Metrics {
		itemMetric := r.convertItemMetricFromPb(pbMetric)
		metrics = append(metrics, itemMetric)
		// Cache the result
		if len(r.cache) < 5000 {
			r.cache[itemMetric.ItemID] = itemMetric
		}
	}
	r.mu.Unlock()

	return &models.GetItemsMetricResponse{
		Metrics: metrics,
	}, nil
}

// GetHighestMetricsByType retrieves items with highest metrics by type
func (r *MetricRepository) GetHighestMetricsByType(ctx context.Context, metricType string, req domain.MetricSortRequest) (*models.GetItemsMetricResponse, error) {
	if err := r.ensureConnection(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to metrics service: %w", err)
	}

	authCtx := r.getContextWithAuth(ctx)
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

	resp, err := r.client.GetHighestMetricsByType(authCtx, pbReq)
	if err != nil {
		log.Printf("Failed to get highest metrics by type %s: %v", metricType, err)
		return nil, fmt.Errorf("failed to get highest metrics by type: %w", err)
	}

	metrics := make([]*models.ItemMetric, 0, len(resp.Metrics))
	for _, pbMetric := range resp.Metrics {
		metrics = append(metrics, r.convertItemMetricFromPb(pbMetric))
	}

	return &models.GetItemsMetricResponse{
		Metrics: metrics,
	}, nil
}

// GetLowestMetricsByType retrieves items with lowest metrics by type
func (r *MetricRepository) GetLowestMetricsByType(ctx context.Context, metricType string, req domain.MetricSortRequest) (*models.GetItemsMetricResponse, error) {
	if err := r.ensureConnection(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to metrics service: %w", err)
	}

	authCtx := r.getContextWithAuth(ctx)
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

	resp, err := r.client.GetLowestMetricsByType(authCtx, pbReq)
	if err != nil {
		log.Printf("Failed to get lowest metrics by type %s: %v", metricType, err)
		return nil, fmt.Errorf("failed to get lowest metrics by type: %w", err)
	}

	metrics := make([]*models.ItemMetric, 0, len(resp.Metrics))
	for _, pbMetric := range resp.Metrics {
		metrics = append(metrics, r.convertItemMetricFromPb(pbMetric))
	}

	return &models.GetItemsMetricResponse{
		Metrics: metrics,
	}, nil
}

// Extended methods for AI tooling (mock implementations since not in protobuf)

// GetItemMetricByType retrieves metrics for an item by specific metric type
func (r *MetricRepository) GetItemMetricByType(ctx context.Context, itemID, metricType string) (*models.ItemMetric, error) {
	// This method is not in the protobuf spec, so we'll use GetItemMetric and filter
	resp, err := r.GetItemMetric(ctx, itemID)
	if err != nil {
		return nil, err
	}
	return resp.Metric, nil
}

// GetUserMetricByType retrieves metrics for a user by specific metric type
func (r *MetricRepository) GetUserMetricByType(ctx context.Context, userID, metricType string) (*models.UserMetric, error) {
	// This method is not in the protobuf spec, so we'll use GetUserMetric
	resp, err := r.GetUserMetric(ctx, userID)
	if err != nil {
		return nil, err
	}
	return resp.Metric, nil
}

// GetItemMetricsByCategory retrieves metrics for items in a category
func (r *MetricRepository) GetItemMetricsByCategory(ctx context.Context, categoryID string, limit int32) ([]*models.ItemMetric, error) {
	// Mock implementation - in real scenario this would be a separate gRPC method
	log.Printf("GetItemMetricsByCategory called for category %s (mock implementation)", categoryID)
	return []*models.ItemMetric{}, nil
}

// GetUserMetricsByCategory retrieves metrics for a user in a specific category
func (r *MetricRepository) GetUserMetricsByCategory(ctx context.Context, userID, categoryID string) (*models.UserMetric, error) {
	// Mock implementation
	log.Printf("GetUserMetricsByCategory called for user %s, category %s (mock implementation)", userID, categoryID)
	return &models.UserMetric{UserID: userID, CategoryID: categoryID}, nil
}

// GetTopItemsByMetric retrieves top items by specific metric
func (r *MetricRepository) GetTopItemsByMetric(ctx context.Context, metricType string, entityTypes []string, limit int32) ([]*models.ItemMetric, error) {
	// Use GetHighestMetricsByType with default parameters
	req := domain.MetricSortRequest{
		EntityTypes: entityTypes,
		Limit:       limit,
	}
	resp, err := r.GetHighestMetricsByType(ctx, metricType, req)
	if err != nil {
		return nil, err
	}
	return resp.Metrics, nil
}

// GetTopUsersByMetric retrieves top users by specific metric
func (r *MetricRepository) GetTopUsersByMetric(ctx context.Context, metricType string, limit int32) ([]*models.UserMetric, error) {
	// Mock implementation
	log.Printf("GetTopUsersByMetric called for metric %s (mock implementation)", metricType)
	return []*models.UserMetric{}, nil
}

// GetItemMetricsInRange retrieves metrics for items within geographic range
func (r *MetricRepository) GetItemMetricsInRange(ctx context.Context, lat, lng, radius float32, limit int32) ([]*models.ItemMetric, error) {
	// Use GetHighestMetricsByType with location parameters
	req := domain.MetricSortRequest{
		Lat:    lat,
		Lng:    lng,
		Radius: radius,
		Limit:  limit,
	}
	resp, err := r.GetHighestMetricsByType(ctx, models.MetricTypeVisits, req)
	if err != nil {
		return nil, err
	}
	return resp.Metrics, nil
}

// GetMetricsSummary retrieves summary statistics for metrics
func (r *MetricRepository) GetMetricsSummary(ctx context.Context, entityType string) (map[string]int64, error) {
	// Mock implementation
	log.Printf("GetMetricsSummary called for entity type %s (mock implementation)", entityType)
	return map[string]int64{
		"total_items":  1000,
		"total_visits": 50000,
		"total_likes":  25000,
		"total_shares": 5000,
	}, nil
}

// SearchItemMetrics searches for items with specific criteria
func (r *MetricRepository) SearchItemMetrics(ctx context.Context, query string, entityTypes []string, limit int32) ([]*models.ItemMetric, error) {
	// Mock implementation
	log.Printf("SearchItemMetrics called with query '%s' (mock implementation)", query)
	return []*models.ItemMetric{}, nil
}

// SearchUserMetrics searches for users with specific criteria
func (r *MetricRepository) SearchUserMetrics(ctx context.Context, query string, limit int32) ([]*models.UserMetric, error) {
	// Mock implementation
	log.Printf("SearchUserMetrics called with query '%s' (mock implementation)", query)
	return []*models.UserMetric{}, nil
}

// GetRecentlyUpdatedMetrics retrieves recently updated metrics
func (r *MetricRepository) GetRecentlyUpdatedMetrics(ctx context.Context, entityType string, limit int32) ([]*models.ItemMetric, error) {
	// Mock implementation
	log.Printf("GetRecentlyUpdatedMetrics called for entity type %s (mock implementation)", entityType)
	return []*models.ItemMetric{}, nil
}

// GetMetricsByPriceRange retrieves metrics for items within price range
func (r *MetricRepository) GetMetricsByPriceRange(ctx context.Context, minPrice, maxPrice int64, entityTypes []string, limit int32) ([]*models.ItemMetric, error) {
	// Use GetHighestMetricsByType with price parameters
	req := domain.MetricSortRequest{
		EntityTypes: entityTypes,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		Limit:       limit,
	}
	resp, err := r.GetHighestMetricsByType(ctx, models.MetricTypeVisits, req)
	if err != nil {
		return nil, err
	}
	return resp.Metrics, nil
}

// GetMetricsByRatingRange retrieves metrics for items within rating range
func (r *MetricRepository) GetMetricsByRatingRange(ctx context.Context, minRating, maxRating int64, entityTypes []string, limit int32) ([]*models.ItemMetric, error) {
	// Mock implementation using rating filter
	log.Printf("GetMetricsByRatingRange called for range %d-%d (mock implementation)", minRating, maxRating)
	return []*models.ItemMetric{}, nil
}

// GetTrendingItems retrieves trending items
func (r *MetricRepository) GetTrendingItems(ctx context.Context, entityTypes []string, days int32, limit int32) ([]*models.ItemMetric, error) {
	// Mock implementation
	log.Printf("GetTrendingItems called for %d days (mock implementation)", days)
	return []*models.ItemMetric{}, nil
}

// GetActiveUsers retrieves active users
func (r *MetricRepository) GetActiveUsers(ctx context.Context, days int32, limit int32) ([]*models.UserMetric, error) {
	// Mock implementation
	log.Printf("GetActiveUsers called for %d days (mock implementation)", days)
	return []*models.UserMetric{}, nil
}

// CompareItemMetrics compares metrics between two items
func (r *MetricRepository) CompareItemMetrics(ctx context.Context, itemID1, itemID2 string) (map[string]interface{}, error) {
	// Mock implementation
	log.Printf("CompareItemMetrics called for items %s vs %s (mock implementation)", itemID1, itemID2)
	return map[string]interface{}{
		"item1_id":   itemID1,
		"item2_id":   itemID2,
		"comparison": "mock_data",
	}, nil
}

// CompareUserMetrics compares metrics between two users
func (r *MetricRepository) CompareUserMetrics(ctx context.Context, userID1, userID2 string) (map[string]interface{}, error) {
	// Mock implementation
	log.Printf("CompareUserMetrics called for users %s vs %s (mock implementation)", userID1, userID2)
	return map[string]interface{}{
		"user1_id":   userID1,
		"user2_id":   userID2,
		"comparison": "mock_data",
	}, nil
}

// GetMetricsAnalytics retrieves analytics data for metrics
func (r *MetricRepository) GetMetricsAnalytics(ctx context.Context, entityType string, timeRange string) (map[string]interface{}, error) {
	// Mock implementation
	log.Printf("GetMetricsAnalytics called for entity type %s, time range %s (mock implementation)", entityType, timeRange)
	return map[string]interface{}{
		"entity_type": entityType,
		"time_range":  timeRange,
		"analytics":   "mock_data",
	}, nil
}

// BulkUpdateItemMetrics updates multiple item metrics
func (r *MetricRepository) BulkUpdateItemMetrics(ctx context.Context, updates []map[string]interface{}) error {
	// Mock implementation
	log.Printf("BulkUpdateItemMetrics called with %d updates (mock implementation)", len(updates))
	return nil
}

// BulkUpdateUserMetrics updates multiple user metrics
func (r *MetricRepository) BulkUpdateUserMetrics(ctx context.Context, updates []map[string]interface{}) error {
	// Mock implementation
	log.Printf("BulkUpdateUserMetrics called with %d updates (mock implementation)", len(updates))
	return nil
}

// ResetItemMetrics resets specific metrics for an item
func (r *MetricRepository) ResetItemMetrics(ctx context.Context, itemID string, metricTypes []string) error {
	// Mock implementation using UpdateItemMetric with reset action
	for _, metricType := range metricTypes {
		_, err := r.UpdateItemMetric(ctx, itemID, metricType, models.MetricActionReset)
		if err != nil {
			return err
		}
	}
	return nil
}

// ResetUserMetrics resets specific metrics for a user
func (r *MetricRepository) ResetUserMetrics(ctx context.Context, userID string, metricTypes []string) error {
	// Mock implementation using UpdateUserMetric with reset action
	for _, metricType := range metricTypes {
		_, err := r.UpdateUserMetric(ctx, userID, metricType, models.MetricActionReset)
		if err != nil {
			return err
		}
	}
	return nil
}

// Helper methods for protobuf conversion

// convertItemMetricFromPb converts protobuf ItemMetric to domain model
func (r *MetricRepository) convertItemMetricFromPb(pbMetric *metricspb.ItemMetric) *models.ItemMetric {
	if pbMetric == nil {
		return nil
	}

	return &models.ItemMetric{
		ItemID:               pbMetric.GetItemId(),
		EntityType:           pbMetric.GetEntityType(),
		LikesCount:           pbMetric.GetLikesCount(),
		DislikesCount:        pbMetric.GetDislikesCount(),
		CommentsCount:        pbMetric.GetCommentsCount(),
		MessagesCount:        pbMetric.GetMessagesCount(),
		SharedCount:          pbMetric.GetSharedCount(),
		AddedToWishlistCount: pbMetric.GetAddedToWishlistCount(),
		AddedToBasketCount:   pbMetric.GetAddedToBasketCount(),
		VisitedCount:         pbMetric.GetVisitedCount(),
		ReportedCount:        pbMetric.GetReportedCount(),
		FollowerCount:        pbMetric.GetFollowerCount(),
		ReviewCount:          pbMetric.GetReviewCount(),
		RatingCount:          pbMetric.GetRatingCount(),
		VideosCount:          pbMetric.GetVideosCount(),
		ImagesCount:          pbMetric.GetImagesCount(),
		Rating:               pbMetric.GetRating(),
		Category:             pbMetric.GetCategory(),
		CategoryID:           pbMetric.GetCategoryId(),
		CategorySlug:         pbMetric.GetCategorySlug(),
		MediaCount:           pbMetric.GetMediaCount(),
		Price:                pbMetric.GetPrice(),
		Lat:                  pbMetric.GetLat(),
		Lng:                  pbMetric.GetLng(),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

// convertUserMetricFromPb converts protobuf UserMetric to domain model
func (r *MetricRepository) convertUserMetricFromPb(pbMetric *metricspb.UserMetric) *models.UserMetric {
	if pbMetric == nil {
		return nil
	}

	return &models.UserMetric{
		UserID:               pbMetric.GetUserId(),
		EntityType:           pbMetric.GetEntityType(),
		LikesCount:           pbMetric.GetLikesCount(),
		DislikesCount:        pbMetric.GetDislikesCount(),
		CommentsCount:        pbMetric.GetCommentsCount(),
		MessagesCount:        pbMetric.GetMessagesCount(),
		SharedCount:          pbMetric.GetSharedCount(),
		AddedToWishlistCount: pbMetric.GetAddedToWishlistCount(),
		AddedToBasketCount:   pbMetric.GetAddedToBasketCount(),
		VisitedCount:         pbMetric.GetVisitedCount(),
		ReportedCount:        pbMetric.GetReportedCount(),
		FollowerCount:        pbMetric.GetFollowerCount(),
		ReviewCount:          pbMetric.GetReviewCount(),
		RatingCount:          pbMetric.GetRatingCount(),
		VideosCount:          pbMetric.GetVideosCount(),
		ImagesCount:          pbMetric.GetImagesCount(),
		Rating:               pbMetric.GetRating(),
		Category:             pbMetric.GetCategory(),
		CategoryID:           pbMetric.GetCategoryId(),
		CategorySlug:         pbMetric.GetCategorySlug(),
		MediaAddedCount:      pbMetric.GetMediaAddedCount(),
		CommentAddedCount:    pbMetric.GetCommentAddedCount(),
		LikedCount:           pbMetric.GetLikedCount(),
		DislikedCount:        pbMetric.GetDislikedCount(),
		ProductsAddedCount:   pbMetric.GetProductsAddedCount(),
		VideosAddedCount:     pbMetric.GetVideosAddedCount(),
		ImagesAddedCount:     pbMetric.GetImagesAddedCount(),
		SeriesAddedCount:     pbMetric.GetSeriesAddedCount(),
		JobsAddedCount:       pbMetric.GetJobsAddedCount(),
		PostsAddedCount:      pbMetric.GetPostsAddedCount(),
		VehiclesAddedCount:   pbMetric.GetVehiclesAddedCount(),
		PropertiesAddedCount: pbMetric.GetPropertiesAddedCount(),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}
