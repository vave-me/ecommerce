package tools

import (
	"context"
	"fmt"
	"middleman/assistants/internal/domain"
)

// ==================== METRIC HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeMetricHandlers() {
	// Core gRPC methods
	r.handlers["metric_update_item"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		metricType := getStringParam(params, "metric_type")
		metricTypeAction := getStringParam(params, "metric_type_action")
		if itemID == "" || metricType == "" || metricTypeAction == "" {
			return nil, fmt.Errorf("item_id, metric_type, and metric_type_action are required")
		}
		return reg.metricRepo.UpdateItemMetric(ctx, itemID, metricType, metricTypeAction)
	}

	r.handlers["metric_share_item"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		if itemID == "" {
			return nil, fmt.Errorf("item_id is required")
		}
		return reg.metricRepo.ShareItem(ctx, itemID)
	}

	r.handlers["metric_visit_item"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		if itemID == "" {
			return nil, fmt.Errorf("item_id is required")
		}
		return reg.metricRepo.VisitItem(ctx, itemID)
	}

	r.handlers["metric_update_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		metricType := getStringParam(params, "metric_type")
		metricTypeAction := getStringParam(params, "metric_type_action")
		if userID == "" || metricType == "" || metricTypeAction == "" {
			return nil, fmt.Errorf("user_id, metric_type, and metric_type_action are required")
		}
		return reg.metricRepo.UpdateUserMetric(ctx, userID, metricType, metricTypeAction)
	}

	r.handlers["metric_get_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return reg.metricRepo.GetUserMetric(ctx, userID)
	}

	r.handlers["metric_get_item"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		if itemID == "" {
			return nil, fmt.Errorf("item_id is required")
		}
		return reg.metricRepo.GetItemMetric(ctx, itemID)
	}

	r.handlers["metric_get_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemIDs := getStringArrayParam(params, "item_ids")
		limit := int32(getInt64Param(params, "limit", 50))
		if len(itemIDs) == 0 {
			return nil, fmt.Errorf("item_ids is required")
		}
		return reg.metricRepo.GetItemsMetric(ctx, itemIDs, limit)
	}

	r.handlers["metric_get_highest_by_type"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		metricType := getStringParam(params, "metric_type")
		if metricType == "" {
			return nil, fmt.Errorf("metric_type is required")
		}
		
		req := domain.MetricSortRequest{
			EntityTypes: getStringArrayParam(params, "entity_types"),
			CategoryId:  getStringParam(params, "category_id"),
			MinPrice:    getInt64Param(params, "min_price", 0),
			MaxPrice:    getInt64Param(params, "max_price", 0),
			Limit:       int32(getInt64Param(params, "limit", 10)),
			Lat:         float32(getFloat64Param(params, "lat", 0)),
			Lng:         float32(getFloat64Param(params, "lng", 0)),
			Radius:      float32(getFloat64Param(params, "radius", 0)),
			CreatedFrom: getStringParam(params, "created_from"),
			CreatedTo:   getStringParam(params, "created_to"),
		}
		return reg.metricRepo.GetHighestMetricsByType(ctx, metricType, req)
	}

	r.handlers["metric_get_lowest_by_type"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		metricType := getStringParam(params, "metric_type")
		if metricType == "" {
			return nil, fmt.Errorf("metric_type is required")
		}
		
		req := domain.MetricSortRequest{
			EntityTypes: getStringArrayParam(params, "entity_types"),
			CategoryId:  getStringParam(params, "category_id"),
			MinPrice:    getInt64Param(params, "min_price", 0),
			MaxPrice:    getInt64Param(params, "max_price", 0),
			Limit:       int32(getInt64Param(params, "limit", 10)),
			Lat:         float32(getFloat64Param(params, "lat", 0)),
			Lng:         float32(getFloat64Param(params, "lng", 0)),
			Radius:      float32(getFloat64Param(params, "radius", 0)),
			CreatedFrom: getStringParam(params, "created_from"),
			CreatedTo:   getStringParam(params, "created_to"),
		}
		return reg.metricRepo.GetLowestMetricsByType(ctx, metricType, req)
	}

	// Extended methods for AI tooling
	r.handlers["metric_get_item_by_type"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		metricType := getStringParam(params, "metric_type")
		if itemID == "" || metricType == "" {
			return nil, fmt.Errorf("item_id and metric_type are required")
		}
		return reg.metricRepo.GetItemMetricByType(ctx, itemID, metricType)
	}

	r.handlers["metric_get_user_by_type"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		metricType := getStringParam(params, "metric_type")
		if userID == "" || metricType == "" {
			return nil, fmt.Errorf("user_id and metric_type are required")
		}
		return reg.metricRepo.GetUserMetricByType(ctx, userID, metricType)
	}

	r.handlers["metric_get_items_by_category"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categoryID := getStringParam(params, "category_id")
		limit := int32(getInt64Param(params, "limit", 50))
		if categoryID == "" {
			return nil, fmt.Errorf("category_id is required")
		}
		return reg.metricRepo.GetItemMetricsByCategory(ctx, categoryID, limit)
	}

	r.handlers["metric_get_user_by_category"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		categoryID := getStringParam(params, "category_id")
		if userID == "" || categoryID == "" {
			return nil, fmt.Errorf("user_id and category_id are required")
		}
		return reg.metricRepo.GetUserMetricsByCategory(ctx, userID, categoryID)
	}

	r.handlers["metric_get_top_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		metricType := getStringParam(params, "metric_type")
		entityTypes := getStringArrayParam(params, "entity_types")
		limit := int32(getInt64Param(params, "limit", 10))
		if metricType == "" {
			return nil, fmt.Errorf("metric_type is required")
		}
		return reg.metricRepo.GetTopItemsByMetric(ctx, metricType, entityTypes, limit)
	}

	r.handlers["metric_get_top_users"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		metricType := getStringParam(params, "metric_type")
		limit := int32(getInt64Param(params, "limit", 10))
		if metricType == "" {
			return nil, fmt.Errorf("metric_type is required")
		}
		return reg.metricRepo.GetTopUsersByMetric(ctx, metricType, limit)
	}

	r.handlers["metric_get_items_in_range"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		lat := float32(getFloat64Param(params, "lat", 0))
		lng := float32(getFloat64Param(params, "lng", 0))
		radius := float32(getFloat64Param(params, "radius", 0))
		limit := int32(getInt64Param(params, "limit", 50))
		if lat == 0 || lng == 0 || radius == 0 {
			return nil, fmt.Errorf("lat, lng, and radius are required")
		}
		return reg.metricRepo.GetItemMetricsInRange(ctx, lat, lng, radius, limit)
	}

	r.handlers["metric_get_summary"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		entityType := getStringParam(params, "entity_type")
		if entityType == "" {
			return nil, fmt.Errorf("entity_type is required")
		}
		return reg.metricRepo.GetMetricsSummary(ctx, entityType)
	}

	r.handlers["metric_search_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		query := getStringParam(params, "query")
		entityTypes := getStringArrayParam(params, "entity_types")
		limit := int32(getInt64Param(params, "limit", 50))
		if query == "" {
			return nil, fmt.Errorf("query is required")
		}
		return reg.metricRepo.SearchItemMetrics(ctx, query, entityTypes, limit)
	}

	r.handlers["metric_search_users"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		query := getStringParam(params, "query")
		limit := int32(getInt64Param(params, "limit", 50))
		if query == "" {
			return nil, fmt.Errorf("query is required")
		}
		return reg.metricRepo.SearchUserMetrics(ctx, query, limit)
	}

	r.handlers["metric_get_trending_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		entityTypes := getStringArrayParam(params, "entity_types")
		days := int32(getInt64Param(params, "days", 7))
		limit := int32(getInt64Param(params, "limit", 20))
		return reg.metricRepo.GetTrendingItems(ctx, entityTypes, days, limit)
	}

	r.handlers["metric_get_active_users"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		days := int32(getInt64Param(params, "days", 7))
		limit := int32(getInt64Param(params, "limit", 50))
		return reg.metricRepo.GetActiveUsers(ctx, days, limit)
	}

	r.handlers["metric_compare_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID1 := getStringParam(params, "item_id_1")
		itemID2 := getStringParam(params, "item_id_2")
		if itemID1 == "" || itemID2 == "" {
			return nil, fmt.Errorf("item_id_1 and item_id_2 are required")
		}
		return reg.metricRepo.CompareItemMetrics(ctx, itemID1, itemID2)
	}

	r.handlers["metric_compare_users"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID1 := getStringParam(params, "user_id_1")
		userID2 := getStringParam(params, "user_id_2")
		if userID1 == "" || userID2 == "" {
			return nil, fmt.Errorf("user_id_1 and user_id_2 are required")
		}
		return reg.metricRepo.CompareUserMetrics(ctx, userID1, userID2)
	}

	r.handlers["metric_get_analytics"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		entityType := getStringParam(params, "entity_type")
		timeRange := getStringParam(params, "time_range", "7d")
		if entityType == "" {
			return nil, fmt.Errorf("entity_type is required")
		}
		return reg.metricRepo.GetMetricsAnalytics(ctx, entityType, timeRange)
	}

	r.handlers["metric_reset_item"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		metricTypes := getStringArrayParam(params, "metric_types")
		if itemID == "" {
			return nil, fmt.Errorf("item_id is required")
		}
		return nil, reg.metricRepo.ResetItemMetrics(ctx, itemID, metricTypes)
	}

	r.handlers["metric_reset_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		metricTypes := getStringArrayParam(params, "metric_types")
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return nil, reg.metricRepo.ResetUserMetrics(ctx, userID, metricTypes)
	}
}