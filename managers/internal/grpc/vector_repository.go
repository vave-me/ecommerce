package grpc

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/vectors/vectorspb"
)

// VectorRepository provides vector database operations via gRPC with resilience patterns
type VectorRepository struct {
	endpoint       string
	auth           *auth.Auth
	circuitBreaker *CircuitBreaker
}

var _ domain.VectorRepository = (*VectorRepository)(nil)

// CircuitBreaker implements simple circuit breaker pattern for vector service
type CircuitBreaker struct {
	failureCount    int64
	lastFailureTime time.Time
	state           CircuitState
	maxFailures     int64
	timeout         time.Duration
}

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// NewVectorRepository creates a new vector repository with circuit breaker
func NewVectorRepository(endpoint string, authInstance *auth.Auth) *VectorRepository {
	return &VectorRepository{
		endpoint: endpoint,
		auth:     authInstance,
		circuitBreaker: &CircuitBreaker{
			maxFailures: 5,
			timeout:     30 * time.Second,
			state:       CircuitClosed,
		},
	}
}

// SearchSimilar finds entities similar to a given entity
func (r *VectorRepository) SearchSimilar(ctx context.Context, entityID string, entityType string, options domain.SearchSimilarOptions) (*domain.VectorSearchResults, error) {
	startTime := time.Now()

	// Check circuit breaker
	if !r.circuitBreaker.CanExecute() {
		log.Warn().Str("entity_id", entityID).Str("entity_type", entityType).Msg("[VECTOR_GRPC] Circuit breaker open, returning empty results")
		return &domain.VectorSearchResults{
			Results:      []*domain.VectorSearchResult{},
			TotalFound:   0,
			QueryTime:    time.Since(startTime),
			SourceFailed: true,
		}, nil
	}

	log.Info().
		Str("entity_id", entityID).
		Str("entity_type", entityType).
		Int64("top_k", options.TopK).
		Msg("[VECTOR_GRPC] SearchSimilar called")

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		r.circuitBreaker.RecordFailure()
		log.Error().Err(err).Msg("[VECTOR_GRPC] Failed to connect to vectors service")
		return &domain.VectorSearchResults{
			Results:      []*domain.VectorSearchResult{},
			TotalFound:   0,
			QueryTime:    time.Since(startTime),
			SourceFailed: true,
		}, nil
	}
	defer conn.Close()

	client := vectorspb.NewVectorServiceClient(conn)

	req := &vectorspb.SearchSimilarEntitiesRequest{
		EntityId:          entityID,
		EntityType:        entityType,
		TopK:              options.TopK,
		ExcludeSelf:       options.ExcludeSelf,
		TargetEntityTypes: options.TargetEntityTypes,
		Filters:           r.convertVectorFilters(options.Filters),
	}

	resp, err := client.SearchSimilarEntities(ctx, req)
	if err != nil {
		r.circuitBreaker.RecordFailure()
		log.Error().Err(err).
			Str("entity_id", entityID).
			Str("entity_type", entityType).
			Msg("[VECTOR_GRPC] SearchSimilarEntities RPC failed")

		// Return empty results instead of error for graceful degradation
		return &domain.VectorSearchResults{
			Results:      []*domain.VectorSearchResult{},
			TotalFound:   0,
			QueryTime:    time.Since(startTime),
			SourceFailed: true,
		}, nil
	}

	r.circuitBreaker.RecordSuccess()
	results := r.convertSearchResults(resp.Results)

	log.Info().
		Str("entity_id", entityID).
		Int("results_count", len(results)).
		Dur("query_time", time.Since(startTime)).
		Msg("[VECTOR_GRPC] SearchSimilar completed successfully")

	return &domain.VectorSearchResults{
		Results:      results,
		TotalFound:   int64(len(results)),
		MaxScore:     r.getMaxScore(results),
		MinScore:     r.getMinScore(results),
		QueryTime:    time.Since(startTime),
		SourceFailed: false,
	}, nil
}

// SearchByVector performs semantic search using raw vector
func (r *VectorRepository) SearchByVector(ctx context.Context, vector []float32, options domain.VectorSearchOptions) (*domain.VectorSearchResults, error) {
	startTime := time.Now()

	if !r.circuitBreaker.CanExecute() {
		log.Warn().Int("vector_dim", len(vector)).Msg("[VECTOR_GRPC] Circuit breaker open, returning empty results")
		return &domain.VectorSearchResults{
			Results:      []*domain.VectorSearchResult{},
			TotalFound:   0,
			QueryTime:    time.Since(startTime),
			SourceFailed: true,
		}, nil
	}

	log.Info().
		Int("vector_dim", len(vector)).
		Int64("top_k", options.TopK).
		Strs("entity_types", options.EntityTypes).
		Msg("[VECTOR_GRPC] SearchByVector called")

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		r.circuitBreaker.RecordFailure()
		log.Error().Err(err).Msg("[VECTOR_GRPC] Failed to connect to vectors service")
		return &domain.VectorSearchResults{
			Results:      []*domain.VectorSearchResult{},
			TotalFound:   0,
			QueryTime:    time.Since(startTime),
			SourceFailed: true,
		}, nil
	}
	defer conn.Close()

	client := vectorspb.NewVectorServiceClient(conn)

	req := &vectorspb.SearchByVectorRequest{
		Vector:         vector,
		EntityTypes:    options.EntityTypes,
		TopK:           options.TopK,
		ScoreThreshold: options.ScoreThreshold,
		WithVectors:    options.WithVectors,
		Filters:        r.convertVectorFilters(options.Filters),
	}

	resp, err := client.SearchByVector(ctx, req)
	if err != nil {
		r.circuitBreaker.RecordFailure()
		log.Error().Err(err).
			Int("vector_dim", len(vector)).
			Msg("[VECTOR_GRPC] SearchByVector RPC failed")

		return &domain.VectorSearchResults{
			Results:      []*domain.VectorSearchResult{},
			TotalFound:   0,
			QueryTime:    time.Since(startTime),
			SourceFailed: true,
		}, nil
	}

	r.circuitBreaker.RecordSuccess()
	results := r.convertSearchResults(resp.Results)

	log.Info().
		Int("vector_dim", len(vector)).
		Int("results_count", len(results)).
		Dur("query_time", time.Since(startTime)).
		Msg("[VECTOR_GRPC] SearchByVector completed successfully")

	return &domain.VectorSearchResults{
		Results:      results,
		TotalFound:   resp.TotalFound,
		MaxScore:     resp.MaxScore,
		MinScore:     resp.MinScore,
		QueryTime:    time.Since(startTime),
		SourceFailed: false,
	}, nil
}

// GetRecommendations provides personalized recommendations
func (r *VectorRepository) GetRecommendations(ctx context.Context, userVector []float32, options domain.RecommendationOptions) (*domain.VectorSearchResults, error) {
	startTime := time.Now()

	if !r.circuitBreaker.CanExecute() {
		log.Warn().Int("user_vector_dim", len(userVector)).Msg("[VECTOR_GRPC] Circuit breaker open, returning empty recommendations")
		return &domain.VectorSearchResults{
			Results:      []*domain.VectorSearchResult{},
			TotalFound:   0,
			QueryTime:    time.Since(startTime),
			SourceFailed: true,
		}, nil
	}

	log.Info().
		Int("user_vector_dim", len(userVector)).
		Int64("top_k", options.TopK).
		Float32("diversity_level", options.DiversityLevel).
		Msg("[VECTOR_GRPC] GetRecommendations called")

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		r.circuitBreaker.RecordFailure()
		log.Error().Err(err).Msg("[VECTOR_GRPC] Failed to connect to vectors service")
		return &domain.VectorSearchResults{
			Results:      []*domain.VectorSearchResult{},
			TotalFound:   0,
			QueryTime:    time.Since(startTime),
			SourceFailed: true,
		}, nil
	}
	defer conn.Close()

	client := vectorspb.NewVectorServiceClient(conn)

	req := &vectorspb.GetRecommendationsRequest{
		UserVector:     userVector,
		EntityTypes:    options.EntityTypes,
		TopK:           options.TopK,
		DiversityLevel: options.DiversityLevel,
		ExcludeIds:     options.ExcludeIDs,
		Filters:        r.convertVectorFilters(options.Filters),
	}

	resp, err := client.GetRecommendations(ctx, req)
	if err != nil {
		r.circuitBreaker.RecordFailure()
		log.Error().Err(err).
			Int("user_vector_dim", len(userVector)).
			Msg("[VECTOR_GRPC] GetRecommendations RPC failed")

		return &domain.VectorSearchResults{
			Results:      []*domain.VectorSearchResult{},
			TotalFound:   0,
			QueryTime:    time.Since(startTime),
			SourceFailed: true,
		}, nil
	}

	r.circuitBreaker.RecordSuccess()
	results := r.convertSearchResults(resp.Recommendations)

	log.Info().
		Int("user_vector_dim", len(userVector)).
		Int("recommendations_count", len(results)).
		Float32("diversity_score", resp.DiversityScore).
		Dur("query_time", time.Since(startTime)).
		Msg("[VECTOR_GRPC] GetRecommendations completed successfully")

	return &domain.VectorSearchResults{
		Results:      results,
		TotalFound:   int64(len(results)),
		MaxScore:     r.getMaxScore(results),
		MinScore:     r.getMinScore(results),
		QueryTime:    time.Since(startTime),
		SourceFailed: false,
	}, nil
}

// GetEntityContext retrieves contextual information about an entity
func (r *VectorRepository) GetEntityContext(ctx context.Context, entityID string, entityType string, options domain.ContextOptions) (*domain.EntityContext, error) {
	if !r.circuitBreaker.CanExecute() {
		log.Warn().Str("entity_id", entityID).Msg("[VECTOR_GRPC] Circuit breaker open, returning minimal context")
		return &domain.EntityContext{
			Entity: &domain.VectorEntity{
				ID:   entityID,
				Type: entityType,
			},
			Similar: []*domain.VectorSearchResult{},
			Failed:  true,
		}, nil
	}

	log.Info().
		Str("entity_id", entityID).
		Str("entity_type", entityType).
		Bool("include_similar", options.IncludeSimilar).
		Msg("[VECTOR_GRPC] GetEntityContext called")

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		r.circuitBreaker.RecordFailure()
		log.Error().Err(err).Msg("[VECTOR_GRPC] Failed to connect to vectors service")
		return &domain.EntityContext{
			Entity: &domain.VectorEntity{
				ID:   entityID,
				Type: entityType,
			},
			Similar: []*domain.VectorSearchResult{},
			Failed:  true,
		}, nil
	}
	defer conn.Close()

	client := vectorspb.NewVectorServiceClient(conn)

	req := &vectorspb.GetEntityContextRequest{
		EntityId:       entityID,
		EntityType:     entityType,
		IncludeSimilar: options.IncludeSimilar,
		ContextSize:    options.ContextSize,
	}

	resp, err := client.GetEntityContext(ctx, req)
	if err != nil {
		r.circuitBreaker.RecordFailure()
		log.Error().Err(err).
			Str("entity_id", entityID).
			Msg("[VECTOR_GRPC] GetEntityContext RPC failed")

		return &domain.EntityContext{
			Entity: &domain.VectorEntity{
				ID:   entityID,
				Type: entityType,
			},
			Similar: []*domain.VectorSearchResult{},
			Failed:  true,
		}, nil
	}

	r.circuitBreaker.RecordSuccess()

	context := &domain.EntityContext{
		Entity:  r.convertEntity(resp.Entity),
		Similar: r.convertSearchResults(resp.Similar),
		Metrics: r.convertEntityMetrics(resp.Metrics),
		Failed:  false,
	}

	log.Info().
		Str("entity_id", entityID).
		Int("similar_count", len(context.Similar)).
		Msg("[VECTOR_GRPC] GetEntityContext completed successfully")

	return context, nil
}

// Health checks if vector service is available
func (r *VectorRepository) Health(ctx context.Context) bool {
	if !r.circuitBreaker.CanExecute() {
		return false
	}

	// Create a short timeout context for health check
	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := r.dialWithAuth(healthCtx)
	if err != nil {
		r.circuitBreaker.RecordFailure()
		return false
	}
	defer conn.Close()

	// Try a simple operation to verify service health
	client := vectorspb.NewVectorServiceClient(conn)

	// Use a minimal search request as health check
	_, err = client.SearchByVector(healthCtx, &vectorspb.SearchByVectorRequest{
		Vector:  []float32{0.1, 0.2, 0.3}, // Minimal test vector
		TopK:    1,
		Filters: &vectorspb.VectorFilters{},
	})

	if err != nil {
		// Don't record failure for health checks to avoid affecting circuit breaker
		log.Debug().Err(err).Msg("[VECTOR_GRPC] Health check failed")
		return false
	}

	return true
}

// Circuit breaker methods
func (cb *CircuitBreaker) CanExecute() bool {
	now := time.Now()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		if now.Sub(cb.lastFailureTime) > cb.timeout {
			cb.state = CircuitHalfOpen
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	default:
		return false
	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.failureCount = 0
	cb.state = CircuitClosed
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.failureCount >= cb.maxFailures {
		cb.state = CircuitOpen
	}
}

// Helper methods for model conversion
func (r *VectorRepository) convertVectorFilters(filters *domain.VectorFilters) *vectorspb.VectorFilters {
	if filters == nil {
		return &vectorspb.VectorFilters{}
	}

	pbFilters := &vectorspb.VectorFilters{
		EntityTypes:     filters.EntityTypes,
		Statuses:        filters.Statuses,
		Categories:      filters.Categories,
		NegotiableOnly:  filters.NegotiableOnly,
		UserType:        filters.UserType,
		MetadataFilters: filters.MetadataFilters,
	}

	if filters.PriceRange != nil {
		pbFilters.PriceRange = &vectorspb.PriceRange{
			MinPrice: filters.PriceRange.MinPrice,
			MaxPrice: filters.PriceRange.MaxPrice,
		}
	}

	if filters.GeoFilter != nil {
		pbFilters.GeoFilter = &vectorspb.GeoFilter{
			Latitude:  filters.GeoFilter.Latitude,
			Longitude: filters.GeoFilter.Longitude,
			RadiusKm:  filters.GeoFilter.RadiusKm,
			Country:   filters.GeoFilter.Country,
			City:      filters.GeoFilter.City,
		}
	}

	if filters.TimeRange != nil {
		pbFilters.TimeRange = &vectorspb.TimeRange{
			StartTime: timestamppb.New(filters.TimeRange.StartTime),
			EndTime:   timestamppb.New(filters.TimeRange.EndTime),
		}
	}

	return pbFilters
}

func (r *VectorRepository) convertSearchResults(pbResults []*vectorspb.VectorSearchResult) []*domain.VectorSearchResult {
	results := make([]*domain.VectorSearchResult, len(pbResults))
	for i, pbResult := range pbResults {
		results[i] = &domain.VectorSearchResult{
			Entity:        r.convertEntity(pbResult.Entity),
			Score:         pbResult.Score,
			Vector:        pbResult.Vector,
			FeatureScores: pbResult.FeatureScores,
			Rank:          pbResult.Rank,
		}
	}
	return results
}

func (r *VectorRepository) convertEntity(pbEntity *vectorspb.Entity) *domain.VectorEntity {
	if pbEntity == nil {
		return nil
	}

	entity := &domain.VectorEntity{
		ID:          pbEntity.Id,
		Type:        pbEntity.Type,
		Name:        pbEntity.Name,
		Description: pbEntity.Description,
		Metadata:    pbEntity.Metadata,
		Tags:        pbEntity.Tags,
		Status:      pbEntity.Status,
		Thumbnail:   pbEntity.Thumbnail,
	}

	if pbEntity.CreatedAt != nil {
		entity.CreatedAt = pbEntity.CreatedAt.AsTime()
	}
	if pbEntity.UpdatedAt != nil {
		entity.UpdatedAt = pbEntity.UpdatedAt.AsTime()
	}

	if pbEntity.Location != nil {
		entity.Location = &domain.EntityLocation{
			Latitude:  pbEntity.Location.Latitude,
			Longitude: pbEntity.Location.Longitude,
			Address:   pbEntity.Location.Address,
			City:      pbEntity.Location.City,
			Country:   pbEntity.Location.Country,
		}
	}

	if pbEntity.Pricing != nil {
		entity.Pricing = &domain.EntityPricing{
			BasePrice:  pbEntity.Pricing.BasePrice,
			SalePrice:  pbEntity.Pricing.SalePrice,
			Currency:   pbEntity.Pricing.Currency,
			Negotiable: pbEntity.Pricing.Negotiable,
		}
	}

	return entity
}

func (r *VectorRepository) convertEntityMetrics(pbMetrics *vectorspb.EntityMetrics) *domain.EntityMetrics {
	if pbMetrics == nil {
		return nil
	}

	metrics := &domain.EntityMetrics{
		ViewCount:        pbMetrics.ViewCount,
		InteractionCount: pbMetrics.InteractionCount,
		PopularityScore:  pbMetrics.PopularityScore,
		CustomMetrics:    pbMetrics.CustomMetrics,
	}

	if pbMetrics.LastUpdated != nil {
		metrics.LastUpdated = pbMetrics.LastUpdated.AsTime()
	}

	return metrics
}

func (r *VectorRepository) getMaxScore(results []*domain.VectorSearchResult) float32 {
	if len(results) == 0 {
		return 0
	}
	max := results[0].Score
	for _, result := range results {
		if result.Score > max {
			max = result.Score
		}
	}
	return max
}

func (r *VectorRepository) getMinScore(results []*domain.VectorSearchResult) float32 {
	if len(results) == 0 {
		return 0
	}
	min := results[0].Score
	for _, result := range results {
		if result.Score < min {
			min = result.Score
		}
	}
	return min
}

// dial creates a gRPC connection to the vectors service
func (r *VectorRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

// dialWithAuth creates an authenticated gRPC connection
func (r *VectorRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	if r.auth == nil {
		return r.dial(ctx)
	}
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
