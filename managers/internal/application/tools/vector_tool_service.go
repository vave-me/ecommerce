package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"

	"middleman/managers/internal/domain"
)

// VectorToolService provides vector database capabilities for AI managers
type VectorToolService struct {
	vectorRepo domain.VectorRepository
}

// NewVectorToolService creates a new vector tool service
func NewVectorToolService(vectorRepo domain.VectorRepository) *VectorToolService {
	return &VectorToolService{
		vectorRepo: vectorRepo,
	}
}

// GetAvailableTools returns vector-related tools for AI managers
func (s *VectorToolService) GetAvailableTools() []Tool {
	return []Tool{
		{
			Name:        "search_similar_entities",
			Description: "Find entities similar to a given entity using vector similarity",
			Parameters: map[string]ParameterInfo{
				"entity_id": {
					Type:        "string",
					Required:    true,
					Description: "ID of the entity to find similar items for",
				},
				"entity_type": {
					Type:        "string",
					Required:    true,
					Description: "Type of entity (product, post, deal, etc.)",
				},
				"top_k": {
					Type:        "integer",
					Required:    false,
					Description: "Number of similar entities to return (default: 5)",
				},
				"target_types": {
					Type:        "array",
					Required:    false,
					Description: "Limit search to specific entity types",
				},
				"exclude_self": {
					Type:        "boolean",
					Required:    false,
					Description: "Exclude the source entity from results (default: true)",
				},
			},
		},
		{
			Name:        "get_entity_context",
			Description: "Get contextual information and similar entities for a specific entity",
			Parameters: map[string]ParameterInfo{
				"entity_id": {
					Type:        "string",
					Required:    true,
					Description: "ID of the entity to get context for",
				},
				"entity_type": {
					Type:        "string",
					Required:    true,
					Description: "Type of entity (product, post, deal, etc.)",
				},
				"include_similar": {
					Type:        "boolean",
					Required:    false,
					Description: "Include similar entities in context (default: true)",
				},
				"context_size": {
					Type:        "integer",
					Required:    false,
					Description: "Number of similar items to include (default: 3)",
				},
			},
		},
		{
			Name:        "get_recommendations",
			Description: "Get personalized recommendations based on user preferences (requires user vector)",
			Parameters: map[string]ParameterInfo{
				"entity_types": {
					Type:        "array",
					Required:    false,
					Description: "Types of entities to recommend (product, post, etc.)",
				},
				"top_k": {
					Type:        "integer",
					Required:    false,
					Description: "Number of recommendations to return (default: 10)",
				},
				"diversity_level": {
					Type:        "number",
					Required:    false,
					Description: "Recommendation diversity level (0.0-1.0, default: 0.5)",
				},
				"exclude_ids": {
					Type:        "array",
					Required:    false,
					Description: "Entity IDs to exclude from recommendations",
				},
				"filters": {
					Type:        "object",
					Required:    false,
					Description: "Additional filters for recommendations",
				},
			},
		},
		{
			Name:        "check_vector_service_health",
			Description: "Check if vector database service is available",
			Parameters:  map[string]ParameterInfo{},
		},
	}
}

// ExecuteTool executes a vector-related tool operation
func (s *VectorToolService) ExecuteTool(ctx context.Context, toolName string, parameters map[string]interface{}) (interface{}, error) {
	log.Info().
		Str("tool_name", toolName).
		Interface("parameters", parameters).
		Msg("[VECTOR_TOOL] Executing vector tool")

	switch toolName {
	case "search_similar_entities":
		return s.searchSimilarEntities(ctx, parameters)
	case "get_entity_context":
		return s.getEntityContext(ctx, parameters)
	case "get_recommendations":
		return s.getRecommendations(ctx, parameters)
	case "check_vector_service_health":
		return s.checkVectorServiceHealth(ctx)
	default:
		return nil, fmt.Errorf("unknown vector tool: %s", toolName)
	}
}

// searchSimilarEntities finds entities similar to a given entity
func (s *VectorToolService) searchSimilarEntities(ctx context.Context, parameters map[string]interface{}) (interface{}, error) {
	// Extract required parameters
	entityID, ok := parameters["entity_id"].(string)
	if !ok || entityID == "" {
		return nil, fmt.Errorf("entity_id is required and must be a string")
	}

	entityType, ok := parameters["entity_type"].(string)
	if !ok || entityType == "" {
		return nil, fmt.Errorf("entity_type is required and must be a string")
	}

	// Extract optional parameters with defaults
	topK := int64(5)
	if topKParam, ok := parameters["top_k"]; ok {
		if topKFloat, ok := topKParam.(float64); ok {
			topK = int64(topKFloat)
		} else if topKStr, ok := topKParam.(string); ok {
			if parsed, err := strconv.ParseInt(topKStr, 10, 64); err == nil {
				topK = parsed
			}
		}
	}

	excludeSelf := true
	if excludeSelfParam, ok := parameters["exclude_self"]; ok {
		if excludeSelfBool, ok := excludeSelfParam.(bool); ok {
			excludeSelf = excludeSelfBool
		}
	}

	var targetTypes []string
	if targetTypesParam, ok := parameters["target_types"]; ok {
		if targetTypesArray, ok := targetTypesParam.([]interface{}); ok {
			for _, t := range targetTypesArray {
				if typeStr, ok := t.(string); ok {
					targetTypes = append(targetTypes, typeStr)
				}
			}
		}
	}

	options := domain.SearchSimilarOptions{
		TopK:              topK,
		TargetEntityTypes: targetTypes,
		ExcludeSelf:       excludeSelf,
		ScoreThreshold:    0.0, // No minimum threshold by default
	}

	results, err := s.vectorRepo.SearchSimilar(ctx, entityID, entityType, options)
	if err != nil {
		log.Error().Err(err).
			Str("entity_id", entityID).
			Str("entity_type", entityType).
			Msg("[VECTOR_TOOL] Failed to search similar entities")
		return nil, fmt.Errorf("failed to search similar entities: %w", err)
	}

	// Format results for AI manager consumption
	response := map[string]interface{}{
		"entity_id":        entityID,
		"entity_type":      entityType,
		"total_found":      results.TotalFound,
		"source_failed":    results.SourceFailed,
		"query_time_ms":    results.QueryTime.Milliseconds(),
		"similar_entities": s.formatSearchResults(results.Results),
	}

	if results.SourceFailed {
		response["message"] = "Vector service unavailable, returning empty results"
	}

	log.Info().
		Str("entity_id", entityID).
		Int("results_count", len(results.Results)).
		Bool("source_failed", results.SourceFailed).
		Msg("[VECTOR_TOOL] Similar entities search completed")

	return response, nil
}

// getEntityContext retrieves contextual information about an entity
func (s *VectorToolService) getEntityContext(ctx context.Context, parameters map[string]interface{}) (interface{}, error) {
	entityID, ok := parameters["entity_id"].(string)
	if !ok || entityID == "" {
		return nil, fmt.Errorf("entity_id is required and must be a string")
	}

	entityType, ok := parameters["entity_type"].(string)
	if !ok || entityType == "" {
		return nil, fmt.Errorf("entity_type is required and must be a string")
	}

	includeSimilar := true
	if includeSimilarParam, ok := parameters["include_similar"]; ok {
		if includeSimilarBool, ok := includeSimilarParam.(bool); ok {
			includeSimilar = includeSimilarBool
		}
	}

	contextSize := int64(3)
	if contextSizeParam, ok := parameters["context_size"]; ok {
		if contextSizeFloat, ok := contextSizeParam.(float64); ok {
			contextSize = int64(contextSizeFloat)
		}
	}

	options := domain.ContextOptions{
		IncludeSimilar: includeSimilar,
		ContextSize:    contextSize,
	}

	context, err := s.vectorRepo.GetEntityContext(ctx, entityID, entityType, options)
	if err != nil {
		log.Error().Err(err).
			Str("entity_id", entityID).
			Str("entity_type", entityType).
			Msg("[VECTOR_TOOL] Failed to get entity context")
		return nil, fmt.Errorf("failed to get entity context: %w", err)
	}

	response := map[string]interface{}{
		"entity_id":   entityID,
		"entity_type": entityType,
		"failed":      context.Failed,
		"entity":      s.formatVectorEntity(context.Entity),
		"similar":     s.formatSearchResults(context.Similar),
		"metrics":     s.formatEntityMetrics(context.Metrics),
	}

	if context.Failed {
		response["message"] = "Vector service unavailable, returning minimal context"
	}

	log.Info().
		Str("entity_id", entityID).
		Int("similar_count", len(context.Similar)).
		Bool("failed", context.Failed).
		Msg("[VECTOR_TOOL] Entity context retrieved")

	return response, nil
}

// getRecommendations provides personalized recommendations
func (s *VectorToolService) getRecommendations(ctx context.Context, parameters map[string]interface{}) (interface{}, error) {
	// Note: This requires user vector which typically comes from user behavior analysis
	// For now, we'll return a placeholder indicating this feature needs user vector input

	topK := int64(10)
	if topKParam, ok := parameters["top_k"]; ok {
		if topKFloat, ok := topKParam.(float64); ok {
			topK = int64(topKFloat)
		}
	}

	diversityLevel := float32(0.5)
	if diversityParam, ok := parameters["diversity_level"]; ok {
		if diversityFloat, ok := diversityParam.(float64); ok {
			diversityLevel = float32(diversityFloat)
		}
	}

	var entityTypes []string
	if entityTypesParam, ok := parameters["entity_types"]; ok {
		if entityTypesArray, ok := entityTypesParam.([]interface{}); ok {
			for _, t := range entityTypesArray {
				if typeStr, ok := t.(string); ok {
					entityTypes = append(entityTypes, typeStr)
				}
			}
		}
	}

	var excludeIDs []string
	if excludeIDsParam, ok := parameters["exclude_ids"]; ok {
		if excludeIDsArray, ok := excludeIDsParam.([]interface{}); ok {
			for _, id := range excludeIDsArray {
				if idStr, ok := id.(string); ok {
					excludeIDs = append(excludeIDs, idStr)
				}
			}
		}
	}

	// For demonstration, return a message indicating this needs user vector
	response := map[string]interface{}{
		"message": "Personalized recommendations require user preference vector which is not yet available. This feature would analyze user behavior to generate a preference vector and then use vector similarity to recommend relevant items.",
		"requested_params": map[string]interface{}{
			"top_k":           topK,
			"diversity_level": diversityLevel,
			"entity_types":    entityTypes,
			"exclude_ids":     excludeIDs,
		},
		"feature_status": "requires_user_vector_implementation",
	}

	log.Info().
		Int64("top_k", topK).
		Float32("diversity_level", diversityLevel).
		Strs("entity_types", entityTypes).
		Msg("[VECTOR_TOOL] Recommendations requested but requires user vector")

	return response, nil
}

// checkVectorServiceHealth checks if vector service is available
func (s *VectorToolService) checkVectorServiceHealth(ctx context.Context) (interface{}, error) {
	isHealthy := s.vectorRepo.Health(ctx)

	response := map[string]interface{}{
		"healthy": isHealthy,
		"status":  "ok",
	}

	if !isHealthy {
		response["status"] = "unavailable"
		response["message"] = "Vector database service is currently unavailable"
	}

	log.Info().
		Bool("healthy", isHealthy).
		Msg("[VECTOR_TOOL] Vector service health check completed")

	return response, nil
}

// Helper methods for formatting results

func (s *VectorToolService) formatSearchResults(results []*domain.VectorSearchResult) []map[string]interface{} {
	formatted := make([]map[string]interface{}, len(results))
	for i, result := range results {
		formatted[i] = map[string]interface{}{
			"entity":           s.formatVectorEntity(result.Entity),
			"similarity_score": result.Score,
			"rank":             result.Rank,
			"feature_scores":   result.FeatureScores,
		}

		if len(result.Vector) > 0 {
			formatted[i]["vector_included"] = true
			formatted[i]["vector_dimension"] = len(result.Vector)
		}
	}
	return formatted
}

func (s *VectorToolService) formatVectorEntity(entity *domain.VectorEntity) map[string]interface{} {
	if entity == nil {
		return nil
	}

	formatted := map[string]interface{}{
		"id":          entity.ID,
		"type":        entity.Type,
		"name":        entity.Name,
		"description": entity.Description,
		"status":      entity.Status,
		"tags":        entity.Tags,
		"metadata":    entity.Metadata,
	}

	if !entity.CreatedAt.IsZero() {
		formatted["created_at"] = entity.CreatedAt.Format("2006-01-02T15:04:05Z")
	}
	if !entity.UpdatedAt.IsZero() {
		formatted["updated_at"] = entity.UpdatedAt.Format("2006-01-02T15:04:05Z")
	}

	if entity.Location != nil {
		formatted["location"] = map[string]interface{}{
			"latitude":  entity.Location.Latitude,
			"longitude": entity.Location.Longitude,
			"address":   entity.Location.Address,
			"city":      entity.Location.City,
			"country":   entity.Location.Country,
		}
	}

	if entity.Pricing != nil {
		formatted["pricing"] = map[string]interface{}{
			"base_price": entity.Pricing.BasePrice,
			"sale_price": entity.Pricing.SalePrice,
			"currency":   entity.Pricing.Currency,
			"negotiable": entity.Pricing.Negotiable,
		}
	}

	if entity.Thumbnail != "" {
		formatted["thumbnail"] = entity.Thumbnail
	}

	return formatted
}

func (s *VectorToolService) formatEntityMetrics(metrics *domain.EntityMetrics) map[string]interface{} {
	if metrics == nil {
		return nil
	}

	formatted := map[string]interface{}{
		"view_count":        metrics.ViewCount,
		"interaction_count": metrics.InteractionCount,
		"popularity_score":  metrics.PopularityScore,
		"custom_metrics":    metrics.CustomMetrics,
	}

	if !metrics.LastUpdated.IsZero() {
		formatted["last_updated"] = metrics.LastUpdated.Format("2006-01-02T15:04:05Z")
	}

	return formatted
}

// Tool and ParameterInfo types (should be defined in a common tools package)
type Tool struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Parameters  map[string]ParameterInfo `json:"parameters"`
}

type ParameterInfo struct {
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}
