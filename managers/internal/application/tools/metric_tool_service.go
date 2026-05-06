package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// MetricToolService handles analytics and performance tracking operations
type MetricToolService struct {
	metrics domain.MetricRepository
}

// NewMetricToolService creates a new metric tool service
func NewMetricToolService(metricRepo domain.MetricRepository) *MetricToolService {
	return &MetricToolService{
		metrics: metricRepo,
	}
}

// ExecuteOperation executes a metric operation with streaming progress
func (s *MetricToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "progress",
		Progress: 25.0,
		Metadata: map[string]interface{}{
			"step":      "extracting_parameters",
			"operation": operation,
		},
		Timestamp: time.Now().Unix(),
	}

	// Extract common parameters
	userID := getStringParam(parameters, "user_id", "")
	itemID := getStringParam(parameters, "item_id", "")
	if itemID == "" {
		itemID = getStringParam(parameters, "id", "")
	}
	metricType := getStringParam(parameters, "metric_type", "")
	metricTypeAction := getStringParam(parameters, "metric_type_action", "")
	categoryID := getStringParam(parameters, "category_id", "")
	limit := int32(getInt64Param(parameters, "limit", 20))

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "progress",
		Progress: 50.0,
		Metadata: map[string]interface{}{
			"step":      "executing_operation",
			"operation": operation,
		},
		Timestamp: time.Now().Unix(),
	}

	switch operation {
	case "update_item_metric":
		return s.updateItemMetric(ctx, itemID, metricType, metricTypeAction, streamChan, toolID)
	case "share_item":
		return s.shareItem(ctx, itemID, streamChan, toolID)
	case "visit_item":
		return s.visitItem(ctx, itemID, streamChan, toolID)
	case "update_user_metric":
		return s.updateUserMetric(ctx, userID, metricType, metricTypeAction, streamChan, toolID)
	case "get_user_metric":
		return s.getUserMetric(ctx, userID, streamChan, toolID)
	case "get_item_metric", "find":
		return s.getItemMetric(ctx, itemID, streamChan, toolID)
	case "get_item_metric_by_type":
		return s.getItemMetricByType(ctx, itemID, metricType, streamChan, toolID)
	case "get_user_metric_by_type":
		return s.getUserMetricByType(ctx, userID, metricType, streamChan, toolID)
	case "get_item_metrics_by_category":
		return s.getItemMetricsByCategory(ctx, categoryID, limit, streamChan, toolID)
	case "get_user_metrics_by_category":
		return s.getUserMetricsByCategory(ctx, userID, categoryID, streamChan, toolID)
	case "get_top_items_by_metric":
		return s.getTopItemsByMetric(ctx, metricType, parameters, limit, streamChan, toolID)
	case "get_top_users_by_metric":
		return s.getTopUsersByMetric(ctx, metricType, limit, streamChan, toolID)
	case "get_metrics_summary":
		entityType := getStringParam(parameters, "entity_type", "item")
		return s.getMetricsSummary(ctx, entityType, streamChan, toolID)
	default:
		return s.handleUnsupportedOperation(operation, streamChan, toolID)
	}
}

func (s *MetricToolService) updateItemMetric(
	ctx context.Context,
	itemID, metricType, metricTypeAction string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	if itemID == "" || metricType == "" || metricTypeAction == "" {
		return nil, fmt.Errorf("item_id, metric_type, and metric_type_action parameters required")
	}

	log.Printf("MetricToolService: Updating item metric for item %s, type: %s, action: %s", itemID, metricType, metricTypeAction)
	result, err := s.metrics.UpdateItemMetric(ctx, itemID, metricType, metricTypeAction)
	if err != nil {
		return nil, fmt.Errorf("update item metric failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"metric":      result,
			"item_id":     itemID,
			"metric_type": metricType,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "update_item_metric",
		"result":      result,
		"item_id":     itemID,
	}, nil
}

func (s *MetricToolService) shareItem(
	ctx context.Context,
	itemID string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	if itemID == "" {
		return nil, fmt.Errorf("item_id parameter required")
	}

	log.Printf("MetricToolService: Sharing item %s", itemID)
	result, err := s.metrics.ShareItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("share item failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"share_result": result,
			"item_id":      itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "share_item",
		"result":      result,
		"item_id":     itemID,
	}, nil
}

func (s *MetricToolService) visitItem(
	ctx context.Context,
	itemID string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	if itemID == "" {
		return nil, fmt.Errorf("item_id parameter required")
	}

	log.Printf("MetricToolService: Visiting item %s", itemID)
	result, err := s.metrics.VisitItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("visit item failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"visit_result": result,
			"item_id":      itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "visit_item",
		"result":      result,
		"item_id":     itemID,
	}, nil
}

func (s *MetricToolService) updateUserMetric(
	ctx context.Context,
	userID, metricType, metricTypeAction string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	if userID == "" || metricType == "" || metricTypeAction == "" {
		return nil, fmt.Errorf("user_id, metric_type, and metric_type_action parameters required")
	}

	log.Printf("MetricToolService: Updating user metric for user %s, type: %s, action: %s", userID, metricType, metricTypeAction)
	result, err := s.metrics.UpdateUserMetric(ctx, userID, metricType, metricTypeAction)
	if err != nil {
		return nil, fmt.Errorf("update user metric failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"metric":      result,
			"user_id":     userID,
			"metric_type": metricType,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "update_user_metric",
		"result":      result,
		"user_id":     userID,
	}, nil
}

func (s *MetricToolService) getUserMetric(
	ctx context.Context,
	userID string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	log.Printf("MetricToolService: Getting user metric for user %s", userID)
	result, err := s.metrics.GetUserMetric(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user metric failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"metric":  result,
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "get_user_metric",
		"result":      result,
		"user_id":     userID,
	}, nil
}

func (s *MetricToolService) getItemMetric(
	ctx context.Context,
	itemID string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	if itemID == "" {
		return nil, fmt.Errorf("item_id parameter required")
	}

	log.Printf("MetricToolService: Getting item metric for item %s", itemID)
	result, err := s.metrics.GetItemMetric(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("get item metric failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"metric":  result,
			"item_id": itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "get_item_metric",
		"result":      result,
		"item_id":     itemID,
	}, nil
}

func (s *MetricToolService) getItemMetricByType(
	ctx context.Context,
	itemID, metricType string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	if itemID == "" || metricType == "" {
		return nil, fmt.Errorf("item_id and metric_type parameters required")
	}

	log.Printf("MetricToolService: Getting item metric by type for item %s, type: %s", itemID, metricType)
	result, err := s.metrics.GetItemMetricByType(ctx, itemID, metricType)
	if err != nil {
		return nil, fmt.Errorf("get item metric by type failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"metric":      result,
			"item_id":     itemID,
			"metric_type": metricType,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "get_item_metric_by_type",
		"result":      result,
		"item_id":     itemID,
		"metric_type": metricType,
	}, nil
}

func (s *MetricToolService) getUserMetricByType(
	ctx context.Context,
	userID, metricType string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	if userID == "" || metricType == "" {
		return nil, fmt.Errorf("user_id and metric_type parameters required")
	}

	log.Printf("MetricToolService: Getting user metric by type for user %s, type: %s", userID, metricType)
	result, err := s.metrics.GetUserMetricByType(ctx, userID, metricType)
	if err != nil {
		return nil, fmt.Errorf("get user metric by type failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"metric":      result,
			"user_id":     userID,
			"metric_type": metricType,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "get_user_metric_by_type",
		"result":      result,
		"user_id":     userID,
		"metric_type": metricType,
	}, nil
}

func (s *MetricToolService) getItemMetricsByCategory(
	ctx context.Context,
	categoryID string,
	limit int32,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	if categoryID == "" {
		return nil, fmt.Errorf("category_id parameter required")
	}

	log.Printf("MetricToolService: Getting item metrics by category %s", categoryID)
	result, err := s.metrics.GetItemMetricsByCategory(ctx, categoryID, limit)
	if err != nil {
		return nil, fmt.Errorf("get item metrics by category failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"metrics":     result,
			"category_id": categoryID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "get_item_metrics_by_category",
		"result":      result,
		"category_id": categoryID,
	}, nil
}

func (s *MetricToolService) getUserMetricsByCategory(
	ctx context.Context,
	userID, categoryID string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	if userID == "" || categoryID == "" {
		return nil, fmt.Errorf("user_id and category_id parameters required")
	}

	log.Printf("MetricToolService: Getting user metrics by category for user %s, category: %s", userID, categoryID)
	result, err := s.metrics.GetUserMetricsByCategory(ctx, userID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("get user metrics by category failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"metrics":     result,
			"user_id":     userID,
			"category_id": categoryID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "get_user_metrics_by_category",
		"result":      result,
		"user_id":     userID,
		"category_id": categoryID,
	}, nil
}

func (s *MetricToolService) getTopItemsByMetric(
	ctx context.Context,
	metricType string,
	parameters map[string]interface{},
	limit int32,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	if metricType == "" {
		return nil, fmt.Errorf("metric_type parameter required")
	}

	// Extract entity types array
	var entityTypes []string
	if entityTypesRaw, ok := parameters["entity_types"].([]interface{}); ok {
		for _, et := range entityTypesRaw {
			if etStr, ok := et.(string); ok {
				entityTypes = append(entityTypes, etStr)
			}
		}
	}

	log.Printf("MetricToolService: Getting top items by metric %s", metricType)
	result, err := s.metrics.GetTopItemsByMetric(ctx, metricType, entityTypes, limit)
	if err != nil {
		return nil, fmt.Errorf("get top items by metric failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"top_items":   result,
			"metric_type": metricType,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "get_top_items_by_metric",
		"result":      result,
		"metric_type": metricType,
	}, nil
}

func (s *MetricToolService) getTopUsersByMetric(
	ctx context.Context,
	metricType string,
	limit int32,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	if metricType == "" {
		return nil, fmt.Errorf("metric_type parameter required")
	}

	log.Printf("MetricToolService: Getting top users by metric %s", metricType)
	result, err := s.metrics.GetTopUsersByMetric(ctx, metricType, limit)
	if err != nil {
		return nil, fmt.Errorf("get top users by metric failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"top_users":   result,
			"metric_type": metricType,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "get_top_users_by_metric",
		"result":      result,
		"metric_type": metricType,
	}, nil
}

func (s *MetricToolService) getMetricsSummary(
	ctx context.Context,
	entityType string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	log.Printf("MetricToolService: Getting metrics summary for entity type %s", entityType)
	result, err := s.metrics.GetMetricsSummary(ctx, entityType)
	if err != nil {
		return nil, fmt.Errorf("get metrics summary failed: %w", err)
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "metric_operation",
		Status:   "completed",
		Progress: 100.0,
		Result: map[string]interface{}{
			"summary":     result,
			"entity_type": entityType,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   "get_metrics_summary",
		"result":      result,
		"for_entity":  entityType,
	}, nil
}

func (s *MetricToolService) handleUnsupportedOperation(
	operation string,
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	streamChan <- ToolExecutionStream{
		ID:        toolID,
		ToolName:  "metric_operation",
		Status:    "error",
		Progress:  100.0,
		Error:     fmt.Sprintf("Metric operation '%s' not implemented", operation),
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "metric",
		"operation":   operation,
		"message":     fmt.Sprintf("Metric operation '%s' not implemented yet", operation),
	}, nil
}
