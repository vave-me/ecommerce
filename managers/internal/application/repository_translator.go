package application

import (
	"fmt"
	"strconv"

	"middleman/managers/internal/application/services"
	"middleman/managers/internal/models"
)

// RepositoryTranslator handles repository query translation
type RepositoryTranslator struct{}

// NewRepositoryTranslator creates a new repository translator
func NewRepositoryTranslator() *RepositoryTranslator {
	return &RepositoryTranslator{}
}

// ValidateQuery validates a repository query
func (r *RepositoryTranslator) ValidateQuery(query services.RepositoryQuery) error {
	if query.EntityType == "" {
		return fmt.Errorf("entity type is required")
	}
	if query.Operation == "" {
		return fmt.Errorf("operation is required")
	}
	return nil
}

// TranslateAIRequest translates an AI request to a repository query
func (r *RepositoryTranslator) TranslateAIRequest(aiRequest map[string]interface{}) (*services.RepositoryQuery, error) {
	entityType, ok := aiRequest["entity_type"].(string)
	if !ok {
		return nil, fmt.Errorf("entity_type is required")
	}

	operation, ok := aiRequest["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation is required")
	}

	parameters, _ := aiRequest["parameters"].(map[string]interface{})
	if parameters == nil {
		parameters = make(map[string]interface{})
	}

	// Extract standard query parameters
	queryParams := r.extractQueryParameters(parameters)

	// Convert to simple parameter map
	paramsMap := r.convertToMap(queryParams)

	return &services.RepositoryQuery{
		EntityType: models.EntityType(entityType),
		Operation:  services.OperationType(operation),
		Parameters: paramsMap,
	}, nil
}

// extractQueryParameters extracts standard query parameters from a map
func (r *RepositoryTranslator) extractQueryParameters(params map[string]interface{}) services.QueryParameters {
	return services.QueryParameters{
		ID:          r.extractString(params, "id"),
		UserID:      r.extractString(params, "user_id"),
		ItemID:      r.extractString(params, "item_id"),
		SearchTerm:  r.extractString(params, "search_term"),
		Name:        r.extractString(params, "name"),
		Description: r.extractString(params, "description"),
		CategoryID:  r.extractString(params, "category_id"),
		MinPrice:    r.extractInt64(params, "min_price"),
		MaxPrice:    r.extractInt64(params, "max_price"),
		Brand:       r.extractString(params, "brand"),
		Condition:   r.extractString(params, "condition"),
		Model:       r.extractString(params, "model"),
		Status:      r.extractString(params, "status"),
		Page:        int(r.extractInt64(params, "page")),
		PageSize:    int(r.extractInt64(params, "page_size")),
		SortBy:      r.extractString(params, "sort_by"),
		SortOrder:   r.extractString(params, "sort_order"),
	}
}

// convertToMap converts QueryParameters to a simple map
func (r *RepositoryTranslator) convertToMap(params services.QueryParameters) map[string]interface{} {
	result := make(map[string]interface{})
	
	if params.ID != "" {
		result["id"] = params.ID
	}
	if params.UserID != "" {
		result["user_id"] = params.UserID
	}
	if params.ItemID != "" {
		result["item_id"] = params.ItemID
	}
	if params.SearchTerm != "" {
		result["search_term"] = params.SearchTerm
	}
	if params.Name != "" {
		result["name"] = params.Name
	}
	if params.Description != "" {
		result["description"] = params.Description
	}
	if params.CategoryID != "" {
		result["category_id"] = params.CategoryID
	}
	if params.MinPrice > 0 {
		result["min_price"] = params.MinPrice
	}
	if params.MaxPrice > 0 {
		result["max_price"] = params.MaxPrice
	}
	if params.Brand != "" {
		result["brand"] = params.Brand
	}
	if params.Condition != "" {
		result["condition"] = params.Condition
	}
	if params.Model != "" {
		result["model"] = params.Model
	}
	if params.Status != "" {
		result["status"] = params.Status
	}
	if params.Page > 0 {
		result["page"] = params.Page
	}
	if params.PageSize > 0 {
		result["page_size"] = params.PageSize
	}
	if params.SortBy != "" {
		result["sort_by"] = params.SortBy
	}
	if params.SortOrder != "" {
		result["sort_order"] = params.SortOrder
	}
	
	return result
}

// Helper methods

func (r *RepositoryTranslator) extractString(params map[string]interface{}, key string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (r *RepositoryTranslator) extractInt64(params map[string]interface{}, key string) int64 {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case int64:
			return v
		case int:
			return int64(v)
		case float64:
			return int64(v)
		case string:
			if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
				return parsed
			}
		}
	}
	return 0
}