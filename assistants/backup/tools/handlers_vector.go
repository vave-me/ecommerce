package tools

import (
	"context"
	"fmt"
	"middleman/assistants/internal/domain"
)

// ==================== VECTOR HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeVectorHandlers() {
	r.handlers["vector_search_similar"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		entityID := getStringParam(params, "entity_id")
		entityType := getStringParam(params, "entity_type")
		topK := getInt64Param(params, "top_k", 10)

		// Validate required parameters
		if err := ValidateIDParam("entity_id", entityID); err != nil {
			return nil, fmt.Errorf("invalid entity_id: %w", err)
		}
		if entityType == "" {
			return nil, fmt.Errorf("entity_type is required")
		}
		if topK <= 0 || topK > 100 {
			return nil, fmt.Errorf("top_k must be between 1 and 100")
		}

		options := domain.SearchSimilarOptions{
			TopK:              topK,
			ExcludeSelf:       getBoolParam(params, "exclude_self", true),
			TargetEntityTypes: getStringArrayParam(params, "target_entity_types"),
			Filters:           createVectorFilters(params),
		}
		return reg.vectorRepo.FindSimilarEntitiesByVectorMatch(ctx, entityID, entityType, options)
	}

	r.handlers["vector_search_by_vector"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		vector := getFloat32ArrayParam(params, "vector")
		topK := getInt64Param(params, "top_k", 10)
		scoreThreshold := getFloat64Param(params, "score_threshold", 0.0)

		// Validate required parameters
		if len(vector) == 0 {
			return nil, fmt.Errorf("vector array cannot be empty")
		}
		if topK <= 0 || topK > 100 {
			return nil, fmt.Errorf("top_k must be between 1 and 100")
		}
		if scoreThreshold < 0.0 || scoreThreshold > 1.0 {
			return nil, fmt.Errorf("score_threshold must be between 0.0 and 1.0")
		}

		options := domain.VectorSearchOptions{
			EntityTypes:    getStringArrayParam(params, "entity_types"),
			TopK:           topK,
			ScoreThreshold: float32(scoreThreshold),
			WithVectors:    getBoolParam(params, "with_vectors", false),
			Filters:        createVectorFilters(params),
		}
		return reg.vectorRepo.PerformSemanticSearchWithVector(ctx, vector, options)
	}

	r.handlers["vector_get_recommendations"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userVector := getFloat32ArrayParam(params, "user_vector")
		topK := getInt64Param(params, "top_k", 10)
		diversityLevel := getFloat64Param(params, "diversity_level", 0.5)

		// Validate required parameters
		if len(userVector) == 0 {
			return nil, fmt.Errorf("user_vector array cannot be empty")
		}
		if topK <= 0 || topK > 100 {
			return nil, fmt.Errorf("top_k must be between 1 and 100")
		}
		if diversityLevel < 0.0 || diversityLevel > 1.0 {
			return nil, fmt.Errorf("diversity_level must be between 0.0 and 1.0")
		}

		options := domain.RecommendationOptions{
			EntityTypes:    getStringArrayParam(params, "entity_types"),
			TopK:           topK,
			DiversityLevel: float32(diversityLevel),
			ExcludeIDs:     getStringArrayParam(params, "exclude_ids"),
			Filters:        createVectorFilters(params),
		}
		return reg.vectorRepo.GetPersonalizedRecommendationsForUser(ctx, userVector, options)
	}

	r.handlers["vector_get_entity_context"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		entityID := getStringParam(params, "entity_id")
		entityType := getStringParam(params, "entity_type")
		contextSize := getInt64Param(params, "context_size", 5)

		// Validate required parameters
		if err := ValidateIDParam("entity_id", entityID); err != nil {
			return nil, fmt.Errorf("invalid entity_id: %w", err)
		}
		if entityType == "" {
			return nil, fmt.Errorf("entity_type is required")
		}
		if contextSize <= 0 || contextSize > 50 {
			return nil, fmt.Errorf("context_size must be between 1 and 50")
		}

		options := domain.ContextOptions{
			IncludeSimilar: getBoolParam(params, "include_similar", true),
			ContextSize:    contextSize,
		}
		return reg.vectorRepo.RetrieveEntityContextualInformation(ctx, entityID, entityType, options)
	}

	r.handlers["vector_health"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.vectorRepo.CheckVectorServiceHealthStatus(ctx), nil
	}
}