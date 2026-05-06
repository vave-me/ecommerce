package application

import (
	"context"
	"fmt"
	"middleman/vectors/internal/models"
	"middleman/vectors/internal/ports"

	"sync"

	"time"
)

// ===============================
// EMBEDDING INTERFACE ABSTRACTION
// ===============================

// TransformationStrategy defines strategies for text transformation
type TransformationStrategy string

const (
	StrategyEnhanced   TransformationStrategy = "enhanced"
	StrategyOptimized  TransformationStrategy = "optimized"
	StrategyMinimal    TransformationStrategy = "minimal"
	StrategyContextual TransformationStrategy = "contextual"
	StrategyMultimodal TransformationStrategy = "multimodal"
)

// EmbeddingInterface defines the interface for embedding services
type EmbeddingInterface interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
	GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
	GenerateEntityEmbedding(ctx context.Context, entityData map[string]interface{}) ([]float32, error)
	GenerateEntityEmbeddingWithPrompt(ctx context.Context, entityType string, entityData map[string]interface{}, strategy TransformationStrategy) ([]float32, error)
	GenerateOptimizedEmbedding(ctx context.Context, entityType string, entityData map[string]interface{}, optimization string) ([]float32, error)
	GetDimensions() int
	GetModel() string
	IsPromptEnabled() bool
}

// ===============================
// APPLICATION LAYER TYPES
// ===============================

type (
	VectorSearchParams struct {
		Vector          []float32              // Raw embedding vector from LLM
		EntityTypes     []string               // Filter by entity types
		TopK            int64                  // Number of results to return
		ScoreThreshold  float64                // Minimum similarity score
		Filters         map[string]interface{} // Additional filters
		IncludeMetadata bool                   // Whether to include full entity data
	}

	SimilarEntitySearchParams struct {
		EntityID       string                 // Find entities similar to this one
		EntityType     string                 // Type of the reference entity
		Vector         []float32              // Alternative: provide vector directly
		TopK           int64                  // Number of similar entities to return
		ScoreThreshold float64                // Minimum similarity score
		ExcludeSelf    bool                   // Exclude the reference entity from results
		Filters        map[string]interface{} // Additional filters
	}

	EntityContextParams struct {
		EntityIDs   []string // Get context for these entities
		EntityTypes []string // Filter by types
		ContextType string   // "related", "similar", "recommendations"
		MaxItems    int64    // Max items to return
	}

	RecommendationParams struct {
		UserID         string                 // User for personalized recommendations
		UserVector     []float32              // User preference vector
		EntityTypes    []string               // Types to recommend
		TopK           int64                  // Number of recommendations
		DiversityLevel float64                // 0.0-1.0, higher = more diverse
		Filters        map[string]interface{} // Additional filters
	}

	GetEntityByIdParams struct {
		EntityID      string // Entity ID to retrieve
		EntityType    string // Entity type
		IncludeVector bool   // Whether to return the vector
	}

	// Vector indexing parameters for integration events
	EntityIndexingParams struct {
		EntityID     string                 // Entity ID to index
		EntityType   string                 // Type of entity
		EntityData   interface{}            // Full entity data
		Strategy     TransformationStrategy // Embedding strategy
		ForceReindex bool                   // Force reindexing even if exists
		Async        bool                   // Whether to index asynchronously
	}

	BatchEntityIndexingParams struct {
		Entities  []EntityIndexingParams // Entities to index
		Strategy  TransformationStrategy // Default strategy for all
		BatchSize int                    // Batch processing size
		Async     bool                   // Whether to process asynchronously
	}

	VectorSearchResults struct {
		Results    []VectorSearchResult
		TotalCount int64
		SearchTime time.Duration
		ScoreStats VectorScoreStats
	}

	VectorSearchResult struct {
		EntityID    string                 // Entity identifier
		EntityType  string                 // Type of entity
		Score       float64                // Similarity score
		Vector      []float32              // Entity vector (if requested)
		Entity      interface{}            // Full entity data (Product, Post, etc.)
		Metadata    map[string]interface{} // Additional metadata
		Distance    float64                // Vector distance
		Explanation string                 // Why this was returned
	}

	EntityContextResults struct {
		EntityContext  []EntityContext
		TotalCount     int64
		ContextType    string
		GenerationTime time.Duration
	}

	EntityContext struct {
		EntityID     string                 // Entity ID
		EntityType   string                 // Entity type
		Entity       interface{}            // Full entity data
		Relationship string                 // Relationship to query context
		Score        float64                // Relevance score
		Metadata     map[string]interface{} // Context metadata
	}

	EntityResult struct {
		EntityID   string                 // Entity ID
		EntityType string                 // Entity type
		Entity     interface{}            // Full entity data
		Vector     []float32              // Entity vector (if requested)
		Metadata   map[string]interface{} // Additional metadata
		Found      bool                   // Whether entity was found
	}

	VectorScoreStats struct {
		MinScore    float64 // Minimum score in results
		MaxScore    float64 // Maximum score in results
		AvgScore    float64 // Average score
		MedianScore float64 // Median score
	}

	// Vector indexing result
	IndexingResult struct {
		EntityID   string        // Entity that was indexed
		EntityType string        // Type of entity
		Success    bool          // Whether indexing succeeded
		Error      error         // Error if any
		Duration   time.Duration // Time taken to index
		VectorID   string        // Vector database ID
		Dimensions int           // Vector dimensions
	}

	BatchIndexingResult struct {
		Results      []IndexingResult       // Individual results
		TotalCount   int                    // Total entities processed
		SuccessCount int                    // Successfully indexed
		FailureCount int                    // Failed to index
		Duration     time.Duration          // Total processing time
		Strategy     TransformationStrategy // Strategy used
	}

	Application interface {
		// VECTOR INDEXING METHODS - Called from integration event handlers
		IndexEntity(ctx context.Context, params EntityIndexingParams) (*IndexingResult, error)
		BatchIndexEntities(ctx context.Context, params BatchEntityIndexingParams) (*BatchIndexingResult, error)
		RemoveEntityVector(ctx context.Context, entityID string, entityType string) error
		BatchRemoveEntityVectors(ctx context.Context, entityIDs []string) error
		ReindexEntity(ctx context.Context, entityID string, entityType string, entityData interface{}) (*IndexingResult, error)

		// VECTOR SEARCH METHODS - Accept raw vectors from LLMs
		SearchByVector(ctx context.Context, params VectorSearchParams) (*VectorSearchResults, error)
		SearchSimilarEntities(ctx context.Context, params SimilarEntitySearchParams) (*VectorSearchResults, error)
		GetEntityContext(ctx context.Context, params EntityContextParams) (*EntityContextResults, error)
		GetRecommendations(ctx context.Context, params RecommendationParams) (*VectorSearchResults, error)
		GetEntityById(ctx context.Context, params GetEntityByIdParams) (*EntityResult, error)

		// VECTOR HEALTH AND DIAGNOSTICS
		GetVectorHealth(ctx context.Context) (*VectorHealthStatus, error)
		ValidateVectorIndex(ctx context.Context) (*IndexValidationResult, error)
	}

	VectorHealthStatus struct {
		Status            string                     // Overall health status
		EmbeddingProvider string                     // Current embedding provider
		VectorProvider    string                     // Current vector database provider
		TotalVectors      int64                      // Total indexed vectors
		VectorsByType     map[string]int64           // Vectors by entity type
		ProviderHealth    map[string]*ProviderHealth // Health of each provider
		LastIndexed       time.Time                  // Last successful indexing
		IndexingRate      float64                    // Current indexing rate
		Errors            []string                   // Recent errors
	}

	ProviderHealth struct {
		Provider    string        // Provider name
		Status      string        // Health status
		Latency     time.Duration // Average latency
		ErrorRate   float64       // Error rate percentage
		LastChecked time.Time     // Last health check
	}

	IndexValidationResult struct {
		Valid            bool                   // Whether index is valid
		TotalEntities    int64                  // Total entities checked
		IndexedEntities  int64                  // Entities with vectors
		MissingVectors   []string               // Entity IDs missing vectors
		CorruptedVectors []string               // Entity IDs with corrupted vectors
		Issues           []IndexValidationIssue // Detailed issues
		Recommendations  []string               // Recommended actions
	}

	IndexValidationIssue struct {
		EntityID    string // Entity with issue
		EntityType  string // Type of entity
		IssueType   string // Type of issue
		Description string // Issue description
		Severity    string // Issue severity
	}

	// app struct implements the Application interface
	app struct {
		// Repository dependencies
		vectorRepo      ports.VectorRepositoryPort
		repos           ports.RepositoryCollection
		embedding       EmbeddingInterface
		mutex           sync.RWMutex
		healthCache     *VectorHealthStatus
		lastHealthCheck time.Time
	}
)

var _ Application = (*app)(nil)

func New(
	vectorRepo ports.VectorRepositoryPort,
	repos ports.RepositoryCollection,
	embedding EmbeddingInterface,
) *app {
	return &app{
		vectorRepo: vectorRepo,
		repos:      repos,
		embedding:  embedding,
	}
}

// ===============================
// VECTOR INDEXING METHODS
// ===============================

// IndexEntity indexes a single entity - called from integration event handlers
func (a *app) IndexEntity(ctx context.Context, params EntityIndexingParams) (*IndexingResult, error) {
	start := time.Now()

	result := &IndexingResult{
		EntityID:   params.EntityID,
		EntityType: params.EntityType,
		Duration:   time.Since(start),
	}

	// Convert entity data to map for embedding generation
	entityData, err := a.convertEntityToMap(params.EntityData)
	if err != nil {
		result.Error = fmt.Errorf("failed to convert entity data: %w", err)
		return result, err
	}

	// Generate embedding using the specified strategy
	var embedding []float32
	if params.Strategy != "" && a.embedding.IsPromptEnabled() {
		embedding, err = a.embedding.GenerateEntityEmbeddingWithPrompt(ctx, params.EntityType, entityData, params.Strategy)
	} else {
		embedding, err = a.embedding.GenerateEntityEmbedding(ctx, entityData)
	}

	if err != nil {
		result.Error = fmt.Errorf("failed to generate embedding: %w", err)
		return result, err
	}

	result.Dimensions = len(embedding)

	// Index in vector repository
	if params.ForceReindex {
		err = a.vectorRepo.ReindexEntity(ctx, params.EntityID, params.EntityType, entityData)
	} else {
		err = a.vectorRepo.IndexEntity(ctx, params.EntityID, params.EntityType, entityData)
	}

	if err != nil {
		result.Error = fmt.Errorf("failed to index entity: %w", err)
		return result, err
	}

	result.Success = true
	result.VectorID = params.EntityID
	result.Duration = time.Since(start)

	return result, nil
}

// BatchIndexEntities indexes multiple entities efficiently
func (a *app) BatchIndexEntities(ctx context.Context, params BatchEntityIndexingParams) (*BatchIndexingResult, error) {
	start := time.Now()

	result := &BatchIndexingResult{
		TotalCount: len(params.Entities),
		Strategy:   params.Strategy,
		Results:    make([]IndexingResult, 0, len(params.Entities)),
	}

	// Process in batches
	batchSize := params.BatchSize
	if batchSize <= 0 {
		batchSize = 50 // Default batch size
	}

	for i := 0; i < len(params.Entities); i += batchSize {
		end := i + batchSize
		if end > len(params.Entities) {
			end = len(params.Entities)
		}

		batch := params.Entities[i:end]
		batchResults := a.processBatch(ctx, batch, params.Strategy)
		result.Results = append(result.Results, batchResults...)

		// Count successes and failures
		for _, r := range batchResults {
			if r.Success {
				result.SuccessCount++
			} else {
				result.FailureCount++
			}
		}
	}

	result.Duration = time.Since(start)
	return result, nil
}

// processBatch processes a batch of entities for indexing
func (a *app) processBatch(ctx context.Context, entities []EntityIndexingParams, strategy TransformationStrategy) []IndexingResult {
	results := make([]IndexingResult, 0, len(entities))

	// Convert entities to index requests
	indexRequests := make([]ports.EntityIndexRequest, 0, len(entities))
	for _, entity := range entities {
		entityData, err := a.convertEntityToMap(entity.EntityData)
		if err != nil {
			results = append(results, IndexingResult{
				EntityID:   entity.EntityID,
				EntityType: entity.EntityType,
				Success:    false,
				Error:      fmt.Errorf("failed to convert entity data: %w", err),
			})
			continue
		}

		indexRequests = append(indexRequests, ports.EntityIndexRequest{
			EntityID:   entity.EntityID,
			EntityType: entity.EntityType,
			EntityData: entityData,
			Strategy:   string(strategy),
		})
	}

	// Batch index entities
	if len(indexRequests) > 0 {
		err := a.vectorRepo.BatchIndexEntities(ctx, indexRequests)
		if err != nil {
			// If batch fails, mark all as failed
			for _, entity := range entities {
				results = append(results, IndexingResult{
					EntityID:   entity.EntityID,
					EntityType: entity.EntityType,
					Success:    false,
					Error:      fmt.Errorf("batch indexing failed: %w", err),
				})
			}
		} else {
			// Mark all as successful
			for _, entity := range entities {
				results = append(results, IndexingResult{
					EntityID:   entity.EntityID,
					EntityType: entity.EntityType,
					Success:    true,
					VectorID:   entity.EntityID,
					Dimensions: a.embedding.GetDimensions(),
				})
			}
		}
	}

	return results
}

// RemoveEntityVector removes a vector for an entity
func (a *app) RemoveEntityVector(ctx context.Context, entityID string, entityType string) error {
	return a.vectorRepo.RemoveEntity(ctx, entityID)
}

// BatchRemoveEntityVectors removes vectors for multiple entities
func (a *app) BatchRemoveEntityVectors(ctx context.Context, entityIDs []string) error {
	return a.vectorRepo.BatchRemoveEntities(ctx, entityIDs)
}

// ReindexEntity forces reindexing of an entity
func (a *app) ReindexEntity(ctx context.Context, entityID string, entityType string, entityData interface{}) (*IndexingResult, error) {
	params := EntityIndexingParams{
		EntityID:     entityID,
		EntityType:   entityType,
		EntityData:   entityData,
		Strategy:     StrategyOptimized,
		ForceReindex: true,
	}
	return a.IndexEntity(ctx, params)
}

// ===============================
// HEALTH AND DIAGNOSTICS
// ===============================

// GetVectorHealth returns the health status of the vector system
func (a *app) GetVectorHealth(ctx context.Context) (*VectorHealthStatus, error) {
	a.mutex.RLock()
	if a.healthCache != nil && time.Since(a.lastHealthCheck) < 30*time.Second {
		defer a.mutex.RUnlock()
		return a.healthCache, nil
	}
	a.mutex.RUnlock()

	// Perform health check
	a.mutex.Lock()
	defer a.mutex.Unlock()

	status := &VectorHealthStatus{
		EmbeddingProvider: a.embedding.GetModel(),
		VectorsByType:     make(map[string]int64),
		ProviderHealth:    make(map[string]*ProviderHealth),
		LastIndexed:       time.Now(),
	}

	// Check vector repository health
	if err := a.vectorRepo.HealthCheck(ctx); err != nil {
		status.Status = "unhealthy"
		status.Errors = append(status.Errors, fmt.Sprintf("Vector repository error: %v", err))
	} else {
		status.Status = "healthy"
	}

	// Get vector statistics
	if stats, err := a.vectorRepo.GetVectorStats(ctx); err == nil {
		status.TotalVectors = stats.TotalVectors
		status.VectorsByType = stats.VectorsByType
		status.IndexingRate = stats.IndexingRate
		status.VectorProvider = stats.VectorProvider
	}

	a.healthCache = status
	a.lastHealthCheck = time.Now()

	return status, nil
}

// ValidateVectorIndex validates the vector index integrity
func (a *app) ValidateVectorIndex(ctx context.Context) (*IndexValidationResult, error) {
	result := &IndexValidationResult{
		MissingVectors:   make([]string, 0),
		CorruptedVectors: make([]string, 0),
		Issues:           make([]IndexValidationIssue, 0),
		Recommendations:  make([]string, 0),
	}

	// Validate the vector repository
	if err := a.vectorRepo.ValidateVectorIndex(ctx); err != nil {
		result.Valid = false
		result.Issues = append(result.Issues, IndexValidationIssue{
			IssueType:   "validation_error",
			Description: fmt.Sprintf("Vector index validation failed: %v", err),
			Severity:    "high",
		})
		result.Recommendations = append(result.Recommendations, "Run vector index rebuild")
	} else {
		result.Valid = true
	}

	return result, nil
}

// ===============================
// HELPER METHODS
// ===============================

// convertEntityToMap converts entity data to map format for embedding generation
func (a *app) convertEntityToMap(entityData interface{}) (map[string]interface{}, error) {
	switch v := entityData.(type) {
	case map[string]interface{}:
		return v, nil
	case *models.Product:
		return map[string]interface{}{
			"id":            v.ProductID,
			"name":          v.Name,
			"description":   v.Description,
			"brand":         v.Brand,
			"model":         v.Model,
			"category_id":   v.CategoryID,
			"category_slug": v.CategorySlug,
			"tags":          v.Tags,
			"condition":     string(v.Condition),
			"price":         v.BasePrice,
			"user_id":       v.UserSellerID,
			"status":        string(v.Status),
			"negotiable":    v.Negotiable,
			"user_type":     string(v.UserType),
			"lat":           v.Lat,
			"lng":           v.Lng,
		}, nil
	case *models.Post:
		return map[string]interface{}{
			"id":          v.PostID,
			"name":        v.Name,
			"description": v.Description,
			"tags":        v.Tags,
			"user_id":     v.UserID,
			"status":      string(v.Status),
			"lat":         v.Lat,
			"lng":         v.Lng,
		}, nil
	case *models.Service:
		return map[string]interface{}{
			"id":            v.ID,
			"name":          v.Name,
			"description":   v.Description,
			"category_id":   v.CategoryID,
			"category_slug": v.CategorySlug,
			"tags":          v.Tags,
			"price":         v.BasePrice,
			"user_id":       v.UserID,
			"status":        v.Status,
			"negotiable":    v.Negotiable,
			"user_type":     v.UserType,
			"lat":           v.Lat,
			"lng":           v.Lng,
		}, nil
	case *models.User:
		fullName := v.FirstName + " " + v.LastName
		if fullName == " " {
			fullName = v.Username
		}
		return map[string]interface{}{
			"id":       v.ID,
			"name":     fullName,
			"email":    v.Email,
			"username": v.Username,
			"location": v.Location,
			"enabled":  v.Enabled,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported entity type: %T", entityData)
	}
}

// ===============================
// VECTOR SEARCH IMPLEMENTATIONS
// ===============================

// ASSISTANT/LLM-FOCUSED METHODS - These methods are designed for LLM consumption

func (a *app) SearchByVector(ctx context.Context, params VectorSearchParams) (*VectorSearchResults, error) {
	// Convert application params to repository params
	repoParams := ports.VectorSearchParams{
		EntityTypes:    params.EntityTypes,
		TopK:           params.TopK,
		ScoreThreshold: params.ScoreThreshold,
		Filters:        params.Filters,
		IncludeVector:  true,
		IncludeEntity:  params.IncludeMetadata,
	}

	// Search using the vector repository
	repoResults, err := a.vectorRepo.SearchByVector(ctx, params.Vector, repoParams)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// Convert repository results to application results
	results := make([]VectorSearchResult, len(repoResults.Results))
	for i, result := range repoResults.Results {
		results[i] = VectorSearchResult{
			EntityID:    result.EntityID,
			EntityType:  result.EntityType,
			Score:       result.Score,
			Distance:    result.Distance,
			Vector:      result.Vector,
			Entity:      result.Entity,
			Metadata:    result.Metadata,
			Explanation: result.Explanation,
		}
	}

	return &VectorSearchResults{
		Results:    results,
		TotalCount: repoResults.TotalCount,
		SearchTime: repoResults.SearchTime,
		ScoreStats: VectorScoreStats{}, // TODO: Calculate stats from results
	}, nil
}

func (a *app) SearchSimilarEntities(ctx context.Context, params SimilarEntitySearchParams) (*VectorSearchResults, error) {
	// Convert application params to repository params
	repoParams := ports.VectorSearchParams{
		EntityTypes:    []string{params.EntityType},
		TopK:           params.TopK,
		ScoreThreshold: params.ScoreThreshold,
		Filters:        params.Filters,
		IncludeVector:  true,
		IncludeEntity:  true,
	}

	var repoResults *ports.VectorSearchResults
	var err error

	if len(params.Vector) > 0 {
		// Use provided vector
		repoResults, err = a.vectorRepo.SearchByVector(ctx, params.Vector, repoParams)
	} else {
		// Find similar to entity
		repoResults, err = a.vectorRepo.FindSimilarToEntity(ctx, params.EntityID, params.EntityType, repoParams)
	}

	if err != nil {
		return nil, fmt.Errorf("similar entity search failed: %w", err)
	}

	// Convert repository results to application results
	results := make([]VectorSearchResult, 0, len(repoResults.Results))
	for _, result := range repoResults.Results {
		// Exclude self if requested
		if params.ExcludeSelf && result.EntityID == params.EntityID {
			continue
		}

		results = append(results, VectorSearchResult{
			EntityID:    result.EntityID,
			EntityType:  result.EntityType,
			Score:       result.Score,
			Distance:    result.Distance,
			Vector:      result.Vector,
			Entity:      result.Entity,
			Metadata:    result.Metadata,
			Explanation: result.Explanation,
		})
	}

	return &VectorSearchResults{
		Results:    results,
		TotalCount: int64(len(results)),
		SearchTime: repoResults.SearchTime,
		ScoreStats: VectorScoreStats{}, // TODO: Calculate stats from results
	}, nil
}

func (a *app) GetEntityContext(ctx context.Context, params EntityContextParams) (*EntityContextResults, error) {
	// Get contextual information for specified entities
	// This could include related entities, categories, etc.
	return &EntityContextResults{
		EntityContext:  []EntityContext{},
		TotalCount:     0,
		ContextType:    params.ContextType,
		GenerationTime: time.Millisecond * 50,
	}, nil
}

func (a *app) GetRecommendations(ctx context.Context, params RecommendationParams) (*VectorSearchResults, error) {
	// Generate recommendations based on user preferences or vector
	var searchVector []float32
	var err error

	if len(params.UserVector) > 0 {
		searchVector = params.UserVector
	} else {
		// TODO: Generate user preference vector from user history
		// For now, return empty results
		return &VectorSearchResults{
			Results:    []VectorSearchResult{},
			TotalCount: 0,
			SearchTime: time.Millisecond * 150,
			ScoreStats: VectorScoreStats{},
		}, nil
	}

	// Convert application params to repository params
	repoParams := ports.VectorSearchParams{
		EntityTypes:    params.EntityTypes,
		TopK:           params.TopK,
		ScoreThreshold: 0.0, // More lenient for recommendations
		Filters:        params.Filters,
		IncludeVector:  false,
		IncludeEntity:  true,
	}

	repoResults, err := a.vectorRepo.SearchByVector(ctx, searchVector, repoParams)
	if err != nil {
		return nil, fmt.Errorf("recommendation search failed: %w", err)
	}

	// Apply diversity filtering if requested
	results := a.applyDiversityFiltering(repoResults.Results, params.DiversityLevel)

	// Convert repository results to application results
	appResults := make([]VectorSearchResult, len(results))
	for i, result := range results {
		appResults[i] = VectorSearchResult{
			EntityID:    result.EntityID,
			EntityType:  result.EntityType,
			Score:       result.Score,
			Distance:    result.Distance,
			Vector:      result.Vector,
			Entity:      result.Entity,
			Metadata:    result.Metadata,
			Explanation: "Recommended based on user preferences",
		}
	}

	return &VectorSearchResults{
		Results:    appResults,
		TotalCount: int64(len(appResults)),
		SearchTime: repoResults.SearchTime,
		ScoreStats: VectorScoreStats{}, // TODO: Calculate stats from results
	}, nil
}

func (a *app) GetEntityById(ctx context.Context, params GetEntityByIdParams) (*EntityResult, error) {
	// Retrieve a specific entity by ID with optional vector
	var entity interface{}
	var err error

	switch params.EntityType {
	case "product":
		entity, err = a.repos.Products().Find(ctx, params.EntityID)
	case "post":
		entity, err = a.repos.Posts().Find(ctx, params.EntityID)
	case "service":
		entity, err = a.repos.Services().Find(ctx, params.EntityID)
	case "user":
		entity, err = a.repos.Users().Find(ctx, params.EntityID)
	default:
		return &EntityResult{Found: false}, fmt.Errorf("unsupported entity type: %s", params.EntityType)
	}

	if err != nil {
		return &EntityResult{Found: false}, err
	}

	result := &EntityResult{
		EntityID:   params.EntityID,
		EntityType: params.EntityType,
		Entity:     entity,
		Metadata:   map[string]interface{}{},
		Found:      true,
	}

	// Get vector if requested
	if params.IncludeVector {
		vector, err := a.vectorRepo.GetEntityVector(ctx, params.EntityID)
		if err == nil {
			result.Vector = vector
		}
		// Don't fail if vector is missing
	}

	return result, nil
}

// applyDiversityFiltering applies diversity filtering to search results
func (a *app) applyDiversityFiltering(results []ports.VectorSearchResult, diversityLevel float64) []ports.VectorSearchResult {
	if diversityLevel <= 0.0 || len(results) <= 1 {
		return results
	}

	// Simple diversity filtering: spread results across different entity types and categories
	seenCategories := make(map[string]int)
	diverseResults := make([]ports.VectorSearchResult, 0, len(results))

	for _, result := range results {
		// Extract category from metadata
		category := "unknown"
		if categoryID, exists := result.Metadata["category_id"]; exists {
			if catStr, ok := categoryID.(string); ok {
				category = catStr
			}
		}

		// Apply diversity filtering based on category frequency
		categoryCount := seenCategories[category]
		maxCategoryCount := int(float64(len(results))*(1.0-diversityLevel) + 1)

		if categoryCount < maxCategoryCount {
			diverseResults = append(diverseResults, result)
			seenCategories[category] = categoryCount + 1
		}
	}

	return diverseResults
}
