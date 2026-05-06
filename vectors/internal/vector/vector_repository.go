package vector

import (
	"context"
	"fmt"
	"time"

	"middleman/vectors/internal/ports"
)

// VectorRepository implements the domain repository for vector operations
type VectorRepository struct {
	embeddingProvider ports.EmbeddingClientProvider
	vectorProvider    ports.VectorClientProvider
	vectorService     *VectorService
	embeddingService  *EmbeddingService
}

// NewVectorRepository creates a new vector repository
func NewVectorRepository(
	embeddingProvider ports.EmbeddingClientProvider,
	vectorProvider ports.VectorClientProvider,
	vectorService *VectorService,
	embeddingService *EmbeddingService,
) *VectorRepository {
	return &VectorRepository{
		embeddingProvider: embeddingProvider,
		vectorProvider:    vectorProvider,
		vectorService:     vectorService,
		embeddingService:  embeddingService,
	}
}

// ===============================
// ENTITY INDEXING OPERATIONS
// ===============================

// IndexEntity indexes an entity with vector embedding
func (r *VectorRepository) IndexEntity(ctx context.Context, entityID string, entityType string, entityData interface{}) error {
	// Convert entity data to map
	entityMap, ok := entityData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("entity data must be a map[string]interface{}")
	}

	// Get embedding client
	embeddingClient, _, err := r.embeddingProvider.GetOptimalProvider(ctx, entityType)
	if err != nil {
		return fmt.Errorf("failed to get embedding provider: %w", err)
	}

	// Generate embedding
	embedding, err := embeddingClient.GenerateEntityEmbedding(ctx, entityMap)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Get vector client
	vectorClient, _, err := r.vectorProvider.GetOptimalProvider(ctx, "index")
	if err != nil {
		return fmt.Errorf("failed to get vector provider: %w", err)
	}

	// Create vector point
	point := ports.VectorPoint{
		ID:     fmt.Sprintf("%s:%s", entityType, entityID),
		Vector: embedding,
		Metadata: map[string]interface{}{
			"entity_id":   entityID,
			"entity_type": entityType,
			"indexed_at":  time.Now().Unix(),
		},
	}

	// Index the vector
	return vectorClient.IndexVector(ctx, point)
}

// BatchIndexEntities indexes multiple entities
func (r *VectorRepository) BatchIndexEntities(ctx context.Context, entities []ports.EntityIndexRequest) error {
	if len(entities) == 0 {
		return nil
	}

	// Get embedding client
	embeddingClient, _, err := r.embeddingProvider.GetOptimalProvider(ctx, "batch")
	if err != nil {
		return fmt.Errorf("failed to get embedding provider: %w", err)
	}

	// Get vector client
	vectorClient, _, err := r.vectorProvider.GetOptimalProvider(ctx, "batch_index")
	if err != nil {
		return fmt.Errorf("failed to get vector provider: %w", err)
	}

	// Process entities in batches
	points := make([]ports.VectorPoint, 0, len(entities))

	for _, entity := range entities {
		// Entity data is already a map
		entityMap := entity.EntityData

		// Generate embedding
		embedding, err := embeddingClient.GenerateEntityEmbedding(ctx, entityMap)
		if err != nil {
			continue // Skip entities that fail embedding generation
		}

		// Create vector point
		point := ports.VectorPoint{
			ID:     fmt.Sprintf("%s:%s", entity.EntityType, entity.EntityID),
			Vector: embedding,
			Metadata: map[string]interface{}{
				"entity_id":   entity.EntityID,
				"entity_type": entity.EntityType,
				"indexed_at":  time.Now().Unix(),
				"strategy":    entity.Strategy,
			},
		}

		points = append(points, point)
	}

	// Batch index all points
	return vectorClient.BatchIndexVectors(ctx, points)
}

// ReindexEntity reindexes an entity with new data
func (r *VectorRepository) ReindexEntity(ctx context.Context, entityID string, entityType string, entityData interface{}) error {
	// Remove existing vector first
	if err := r.RemoveEntity(ctx, entityID); err != nil {
		// Continue even if removal fails (entity might not exist)
	}

	// Index with new data
	return r.IndexEntity(ctx, entityID, entityType, entityData)
}

// ===============================
// ENTITY REMOVAL OPERATIONS
// ===============================

// RemoveEntity removes an entity's vector
func (r *VectorRepository) RemoveEntity(ctx context.Context, entityID string) error {
	// Get vector client
	vectorClient, err := r.vectorProvider.GetDefaultClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to get vector provider: %w", err)
	}

	// Delete the vector (using entityID as the key)
	return vectorClient.DeleteVector(ctx, entityID)
}

// BatchRemoveEntities removes multiple entities' vectors
func (r *VectorRepository) BatchRemoveEntities(ctx context.Context, entityIDs []string) error {
	if len(entityIDs) == 0 {
		return nil
	}

	// Get vector client
	vectorClient, _, err := r.vectorProvider.GetOptimalProvider(ctx, "batch_delete")
	if err != nil {
		return fmt.Errorf("failed to get vector provider: %w", err)
	}

	// Batch delete vectors
	return vectorClient.BatchDeleteVectors(ctx, entityIDs)
}

// RemoveEntitiesByType removes all entities of a specific type
func (r *VectorRepository) RemoveEntitiesByType(ctx context.Context, entityType string) error {
	// This would require a more sophisticated implementation that can query by metadata
	// For now, return an error indicating it's not implemented
	return fmt.Errorf("removing entities by type is not currently implemented")
}

// ===============================
// SEARCH OPERATIONS
// ===============================

// SearchSimilarEntities performs vector similarity search
func (r *VectorRepository) SearchSimilarEntities(ctx context.Context, params ports.VectorSearchParams) (*ports.VectorSearchResults, error) {
	// This requires a vector to search with - return error if not provided
	return nil, fmt.Errorf("SearchSimilarEntities requires SearchByVector to be called with an actual vector")
}

// SearchByVector performs vector similarity search with a provided vector
func (r *VectorRepository) SearchByVector(ctx context.Context, vector []float32, params ports.VectorSearchParams) (*ports.VectorSearchResults, error) {
	if len(vector) == 0 {
		return nil, fmt.Errorf("vector cannot be empty")
	}

	// Get vector client
	vectorClient, _, err := r.vectorProvider.GetOptimalProvider(ctx, "search")
	if err != nil {
		return nil, fmt.Errorf("failed to get vector provider: %w", err)
	}

	// Build search request
	searchReq := ports.SearchRequest{
		Vector:         vector,
		TopK:           params.TopK,
		ScoreThreshold: float32(params.ScoreThreshold),
		WithVector:     params.IncludeVector,
		EntityTypes:    params.EntityTypes,
		Filter:         params.Filters,
	}

	// Perform search
	startTime := time.Now()
	response, err := vectorClient.SearchSimilar(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}
	searchTime := time.Since(startTime)

	// Convert results
	results := make([]ports.VectorSearchResult, len(response.Points))
	for i, point := range response.Points {
		entityID := ""
		entityType := ""

		if point.Metadata != nil {
			if id, ok := point.Metadata["entity_id"].(string); ok {
				entityID = id
			}
			if typ, ok := point.Metadata["entity_type"].(string); ok {
				entityType = typ
			}
		}

		results[i] = ports.VectorSearchResult{
			EntityID:   entityID,
			EntityType: entityType,
			Score:      float64(point.Score),
			Distance:   1.0 - float64(point.Score), // Convert score to distance
			Vector:     point.Vector,
			Metadata:   point.Metadata,
		}
	}

	return &ports.VectorSearchResults{
		Results:    results,
		TotalCount: response.Total,
		SearchTime: searchTime,
		Query:      params,
	}, nil
}

// FindSimilarToEntity finds entities similar to a specific entity
func (r *VectorRepository) FindSimilarToEntity(ctx context.Context, entityID string, entityType string, params ports.VectorSearchParams) (*ports.VectorSearchResults, error) {
	// Get the entity's vector first
	vector, err := r.GetEntityVector(ctx, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity vector: %w", err)
	}

	// Use the entity's vector for search
	return r.SearchByVector(ctx, vector, params)
}

// ===============================
// ENTITY VECTOR OPERATIONS
// ===============================

// GetEntityVector retrieves an entity's vector
func (r *VectorRepository) GetEntityVector(ctx context.Context, entityID string) ([]float32, error) {
	// Get vector client
	vectorClient, err := r.vectorProvider.GetDefaultClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vector provider: %w", err)
	}

	// Get the vector point
	point, err := vectorClient.GetVector(ctx, entityID, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get vector: %w", err)
	}

	return point.Vector, nil
}

// HasEntityVector checks if an entity has a vector
func (r *VectorRepository) HasEntityVector(ctx context.Context, entityID string) (bool, error) {
	// Get vector client
	vectorClient, err := r.vectorProvider.GetDefaultClient(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get vector provider: %w", err)
	}

	// Try to get the vector without the actual vector data
	_, err = vectorClient.GetVector(ctx, entityID, false)
	if err != nil {
		return false, nil // Assume it doesn't exist
	}

	return true, nil
}

// GetEntityVectors retrieves vectors for multiple entities
func (r *VectorRepository) GetEntityVectors(ctx context.Context, entityIDs []string) (map[string][]float32, error) {
	if len(entityIDs) == 0 {
		return make(map[string][]float32), nil
	}

	// Get vector client
	vectorClient, err := r.vectorProvider.GetDefaultClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vector provider: %w", err)
	}

	vectors := make(map[string][]float32)

	// Get vectors one by one (could be optimized with batch operations)
	for _, entityID := range entityIDs {
		point, err := vectorClient.GetVector(ctx, entityID, true)
		if err != nil {
			continue // Skip entities that don't have vectors
		}
		vectors[entityID] = point.Vector
	}

	return vectors, nil
}

// ===============================
// HEALTH AND DIAGNOSTICS
// ===============================

// HealthCheck performs a health check on the vector repository
func (r *VectorRepository) HealthCheck(ctx context.Context) error {
	// Check embedding provider health
	embeddingClient, err := r.embeddingProvider.GetDefaultClient(ctx)
	if err != nil {
		return fmt.Errorf("embedding provider unhealthy: %w", err)
	}

	if err := embeddingClient.HealthCheck(ctx); err != nil {
		return fmt.Errorf("embedding client health check failed: %w", err)
	}

	// Check vector provider health
	vectorClient, err := r.vectorProvider.GetDefaultClient(ctx)
	if err != nil {
		return fmt.Errorf("vector provider unhealthy: %w", err)
	}

	if err := vectorClient.HealthCheck(ctx); err != nil {
		return fmt.Errorf("vector client health check failed: %w", err)
	}

	return nil
}

// GetVectorStats returns vector repository statistics
func (r *VectorRepository) GetVectorStats(ctx context.Context) (*ports.VectorRepositoryStats, error) {
	// Get vector client
	vectorClient, err := r.vectorProvider.GetDefaultClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vector provider: %w", err)
	}

	// Get stats from vector database
	dbStats, err := vectorClient.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vector database stats: %w", err)
	}

	// Convert to repository stats
	return &ports.VectorRepositoryStats{
		TotalVectors:      dbStats.TotalVectors,
		VectorsByType:     make(map[string]int64), // Would need to query metadata for accurate counts
		IndexSize:         dbStats.IndexSize,
		MemoryUsage:       dbStats.MemoryUsage,
		AverageQueryTime:  dbStats.QueryLatency,
		IndexingRate:      dbStats.IndexingRate,
		EmbeddingProvider: "openai",
		VectorProvider:    "qdrant",
		HealthStatus:      "healthy",
		LastUpdated:       time.Now(),
		ProviderStats:     dbStats.Collections,
	}, nil
}

// ValidateVectorIndex validates the vector index integrity
func (r *VectorRepository) ValidateVectorIndex(ctx context.Context) error {
	// Perform basic health check
	return r.HealthCheck(ctx)
}
