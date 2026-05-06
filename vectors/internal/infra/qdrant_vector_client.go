package infra

import (
	"context"
	"middleman/vectors/internal/ports"
	"middleman/vectors/internal/vector"
	"time"
)

// QdrantVectorClient adapts the vector service to the VectorDatabaseClient interface
type QdrantVectorClient struct {
	service *vector.VectorService
}

// NewQdrantVectorClient creates a new Qdrant vector client adapter
func NewQdrantVectorClient(service *vector.VectorService) *QdrantVectorClient {
	return &QdrantVectorClient{
		service: service,
	}
}

// IndexVector indexes a single vector point
func (c *QdrantVectorClient) IndexVector(ctx context.Context, point ports.VectorPoint) error {
	vectorPoint := vector.VectorPoint{
		ID:       point.ID,
		Vector:   point.Vector,
		Metadata: point.Metadata,
	}

	return c.service.IndexVector(ctx, vectorPoint)
}

// BatchIndexVectors indexes multiple vector points
func (c *QdrantVectorClient) BatchIndexVectors(ctx context.Context, points []ports.VectorPoint) error {
	vectorPoints := make([]vector.VectorPoint, len(points))
	for i, point := range points {
		vectorPoints[i] = vector.VectorPoint{
			ID:       point.ID,
			Vector:   point.Vector,
			Metadata: point.Metadata,
		}
	}

	return c.service.BatchIndexVectors(ctx, vectorPoints)
}

// SearchSimilar performs vector similarity search
func (c *QdrantVectorClient) SearchSimilar(ctx context.Context, req ports.SearchRequest) (*ports.SearchResponse, error) {
	searchReq := vector.SearchRequest{
		Vector:         req.Vector,
		TopK:           req.TopK,
		ScoreThreshold: req.ScoreThreshold,
		Filter:         req.Filter,
		WithVector:     req.WithVector,
	}

	response, err := c.service.SearchSimilar(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	points := make([]ports.VectorPoint, len(response.Points))
	for i, point := range response.Points {
		points[i] = ports.VectorPoint{
			ID:       point.ID,
			Vector:   point.Vector,
			Metadata: point.Metadata,
			Score:    point.Score,
		}
	}

	return &ports.SearchResponse{
		Points: points,
		Total:  response.Total,
		Took:   time.Millisecond * 100, // Approximate response time
	}, nil
}

// GetVector retrieves a specific vector by ID
func (c *QdrantVectorClient) GetVector(ctx context.Context, id string, withVector bool) (*ports.VectorPoint, error) {
	point, err := c.service.GetVector(ctx, id, withVector)
	if err != nil {
		return nil, err
	}

	return &ports.VectorPoint{
		ID:       point.ID,
		Vector:   point.Vector,
		Metadata: point.Metadata,
		Score:    point.Score,
	}, nil
}

// DeleteVector removes a vector by ID
func (c *QdrantVectorClient) DeleteVector(ctx context.Context, id string) error {
	return c.service.DeleteVector(ctx, id)
}

// BatchDeleteVectors removes multiple vectors by IDs
func (c *QdrantVectorClient) BatchDeleteVectors(ctx context.Context, ids []string) error {
	return c.service.BatchDeleteVectors(ctx, ids)
}

// HealthCheck performs a health check on the vector database
func (c *QdrantVectorClient) HealthCheck(ctx context.Context) error {
	// Use a valid UUID for health check
	healthCheckID := "00000000-0000-0000-0000-000000000000"
	_, err := c.service.GetVector(ctx, healthCheckID, false)
	if err != nil && err.Error() != "vector not found" && err.Error() != "point not found" {
		return err
	}
	return nil
}

// GetStats returns vector database statistics
func (c *QdrantVectorClient) GetStats(ctx context.Context) (*ports.VectorDBStats, error) {
	// Return basic stats - in a real implementation, this would query Qdrant for actual statistics
	return &ports.VectorDBStats{
		TotalVectors:   0, // Would need to query Qdrant for actual count
		IndexSize:      0,
		MemoryUsage:    0,
		QueryLatency:   time.Millisecond * 50,
		IndexingRate:   100.0,
		Collections:    map[string]interface{}{"vectors": "active"},
		ConnectionPool: nil,
		LastUpdated:    time.Now(),
	}, nil
}

// GetProviderInfo returns information about the vector database provider
func (c *QdrantVectorClient) GetProviderInfo() ports.VectorDBProviderInfo {
	return ports.VectorDBProviderInfo{
		Provider:    "qdrant",
		Version:     "1.0.0",
		Endpoint:    "qdrant:6334",
		Collection:  "vectors",
		Dimensions:  1536,
		IndexType:   "hnsw",
		ConnectedAt: time.Now(),
	}
}

// Close closes the connection to the vector database
func (c *QdrantVectorClient) Close() error {
	return c.service.Close()
}

// IsConnected returns whether the client is connected to the vector database
func (c *QdrantVectorClient) IsConnected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.HealthCheck(ctx)
	return err == nil
}
