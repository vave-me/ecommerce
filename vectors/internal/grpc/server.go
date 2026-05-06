// File: vector/internal/grpc/server.go
package grpc

import (
	"context"
	"fmt"
	"log"

	"middleman/internal/errorsotel"
	"middleman/vectors/internal/application"
	"middleman/vectors/vectorspb"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	app application.Application
	vectorspb.UnimplementedVectorServiceServer
}

func RegisterServer(
	ctx context.Context,
	app application.Application,
	registrar grpc.ServiceRegistrar,
) error {
	vectorspb.RegisterVectorServiceServer(registrar, server{app: app})
	return nil
}

func handlePanic(span trace.Span, methodName string) {
	if r := recover(); r != nil {
		err := fmt.Errorf("panic recovered in %s: %v", methodName, r)
		if span != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, "panic")
		}
		log.Printf("Panic recovered in %s: %v", methodName, r)
	}
}

// ------------------------------
// Pure Vector Operations - LLM-Focused
// ------------------------------

func (s server) SearchByVector(ctx context.Context, req *vectorspb.SearchByVectorRequest) (*vectorspb.SearchByVectorResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchByVector")

	span.SetAttributes(
		attribute.Int("vector_dimensions", len(req.GetVector())),
		attribute.StringSlice("entity_types", req.GetEntityTypes()),
		attribute.Int64("top_k", req.GetTopK()),
	)

	log.Printf("[VectorService] SearchByVector: vector_dim=%d entity_types=%v top_k=%d",
		len(req.GetVector()), req.GetEntityTypes(), req.GetTopK())

	params := application.VectorSearchParams{
		Vector:          req.GetVector(),
		EntityTypes:     req.GetEntityTypes(),
		TopK:            req.GetTopK(),
		ScoreThreshold:  float64(req.GetScoreThreshold()),
		Filters:         s.convertVectorFilters(req.GetFilters()),
		IncludeMetadata: req.GetWithVectors(),
	}

	results, err := s.app.SearchByVector(ctx, params)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoResults := make([]*vectorspb.VectorSearchResult, len(results.Results))
	for i, result := range results.Results {
		protoResults[i] = s.vectorSearchResultFromDomain(result)
	}

	return &vectorspb.SearchByVectorResponse{
		Results:    protoResults,
		Stats:      s.vectorScoreStatsFromDomain(results.ScoreStats),
		TotalFound: results.TotalCount,
		MaxScore:   float32(results.ScoreStats.MaxScore),
		MinScore:   float32(results.ScoreStats.MinScore),
	}, nil
}

func (s server) SearchSimilarEntities(ctx context.Context, req *vectorspb.SearchSimilarEntitiesRequest) (*vectorspb.SearchSimilarEntitiesResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchSimilarEntities")

	span.SetAttributes(
		attribute.String("entity_id", req.GetEntityId()),
		attribute.String("entity_type", req.GetEntityType()),
		attribute.Int64("top_k", req.GetTopK()),
	)

	log.Printf("[VectorService] SearchSimilarEntities: entity_id=%s entity_type=%s top_k=%d",
		req.GetEntityId(), req.GetEntityType(), req.GetTopK())

	params := application.SimilarEntitySearchParams{
		EntityID:    req.GetEntityId(),
		EntityType:  req.GetEntityType(),
		TopK:        req.GetTopK(),
		ExcludeSelf: req.GetExcludeSelf(),
		Filters:     s.convertVectorFilters(req.GetFilters()),
	}

	results, err := s.app.SearchSimilarEntities(ctx, params)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoResults := make([]*vectorspb.VectorSearchResult, len(results.Results))
	for i, result := range results.Results {
		protoResults[i] = s.vectorSearchResultFromDomain(result)
	}

	return &vectorspb.SearchSimilarEntitiesResponse{
		Results:        protoResults,
		SourceEntityId: req.GetEntityId(),
		Stats:          s.vectorScoreStatsFromDomain(results.ScoreStats),
	}, nil
}

func (s server) GetEntityContext(ctx context.Context, req *vectorspb.GetEntityContextRequest) (*vectorspb.GetEntityContextResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "GetEntityContext")

	span.SetAttributes(
		attribute.String("entity_id", req.GetEntityId()),
		attribute.String("entity_type", req.GetEntityType()),
		attribute.Bool("include_similar", req.GetIncludeSimilar()),
	)

	log.Printf("[VectorService] GetEntityContext: entity_id=%s entity_type=%s include_similar=%t",
		req.GetEntityId(), req.GetEntityType(), req.GetIncludeSimilar())

	params := application.EntityContextParams{
		EntityIDs:   []string{req.GetEntityId()},
		EntityTypes: []string{req.GetEntityType()},
		ContextType: "similar",
		MaxItems:    req.GetContextSize(),
	}

	_, err := s.app.GetEntityContext(ctx, params)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var similar []*vectorspb.VectorSearchResult
	// For now, return empty similar - would need to implement in application layer

	return &vectorspb.GetEntityContextResponse{
		Entity:  &vectorspb.Entity{}, // Placeholder - would need proper conversion
		Similar: similar,
		Metrics: &vectorspb.EntityMetrics{},
	}, nil
}

func (s server) GetRecommendations(ctx context.Context, req *vectorspb.GetRecommendationsRequest) (*vectorspb.GetRecommendationsResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "GetRecommendations")

	span.SetAttributes(
		attribute.Int("user_vector_dimensions", len(req.GetUserVector())),
		attribute.StringSlice("entity_types", req.GetEntityTypes()),
		attribute.Int64("top_k", req.GetTopK()),
		attribute.Float64("diversity_level", float64(req.GetDiversityLevel())),
	)

	log.Printf("[VectorService] GetRecommendations: vector_dim=%d entity_types=%v top_k=%d diversity=%.2f",
		len(req.GetUserVector()), req.GetEntityTypes(), req.GetTopK(), req.GetDiversityLevel())

	params := application.RecommendationParams{
		UserVector:     req.GetUserVector(),
		EntityTypes:    req.GetEntityTypes(),
		TopK:           req.GetTopK(),
		DiversityLevel: float64(req.GetDiversityLevel()),
		Filters:        s.convertVectorFilters(req.GetFilters()),
	}

	results, err := s.app.GetRecommendations(ctx, params)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoResults := make([]*vectorspb.VectorSearchResult, len(results.Results))
	for i, result := range results.Results {
		protoResults[i] = s.vectorSearchResultFromDomain(result)
	}

	return &vectorspb.GetRecommendationsResponse{
		Recommendations: protoResults,
		DiversityScore:  0.0, // Placeholder - would need to implement in application layer
		Stats:           s.vectorScoreStatsFromDomain(results.ScoreStats),
	}, nil
}

func (s server) GetEntityById(ctx context.Context, req *vectorspb.GetEntityByIdRequest) (*vectorspb.GetEntityByIdResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "GetEntityById")

	span.SetAttributes(
		attribute.String("entity_id", req.GetEntityId()),
		attribute.String("entity_type", req.GetEntityType()),
		attribute.Bool("with_vector", req.GetWithVector()),
	)

	log.Printf("[VectorService] GetEntityById: entity_id=%s entity_type=%s with_vector=%t",
		req.GetEntityId(), req.GetEntityType(), req.GetWithVector())

	params := application.GetEntityByIdParams{
		EntityID:      req.GetEntityId(),
		EntityType:    req.GetEntityType(),
		IncludeVector: req.GetWithVector(),
	}

	result, err := s.app.GetEntityById(ctx, params)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &vectorspb.GetEntityByIdResponse{
		Entity:  s.entityFromResult(result),
		Vector:  result.Vector,
		Metrics: &vectorspb.EntityMetrics{},
	}, nil
}

// ------------------------------
// Helper Methods - Pure Vector Domain Conversion
// ------------------------------

func (s server) vectorSearchResultFromDomain(result application.VectorSearchResult) *vectorspb.VectorSearchResult {
	featureScores := make(map[string]float32)
	// Convert from float64 to float32 for protobuf
	// Note: feature_scores would need to be added to application.VectorSearchResult

	return &vectorspb.VectorSearchResult{
		Entity:        s.entityFromInterface(result.Entity, result.EntityType, result.EntityID),
		Score:         float32(result.Score),
		Vector:        result.Vector,
		FeatureScores: featureScores,
		Rank:          0, // Placeholder - would need to be added to application layer
	}
}

func (s server) entityFromInterface(entity interface{}, entityType, entityID string) *vectorspb.Entity {
	if entity == nil {
		return &vectorspb.Entity{
			Id:   entityID,
			Type: entityType,
		}
	}

	// For now, return basic entity - would need proper conversion based on entity type
	return &vectorspb.Entity{
		Id:   entityID,
		Type: entityType,
	}
}

func (s server) entityFromResult(result *application.EntityResult) *vectorspb.Entity {
	if result == nil || !result.Found {
		return nil
	}

	return &vectorspb.Entity{
		Id:   result.EntityID,
		Type: result.EntityType,
	}
}

func (s server) vectorScoreStatsFromDomain(stats application.VectorScoreStats) *vectorspb.VectorScoreStats {
	return &vectorspb.VectorScoreStats{
		MaxScore:         float32(stats.MaxScore),
		MinScore:         float32(stats.MinScore),
		AvgScore:         float32(stats.AvgScore),
		MedianScore:      float32(stats.MedianScore),
		TotalComparisons: 0,        // Would need to be added to application layer
		VectorDimensions: 0,        // Would need to be added to application layer
		SimilarityMetric: "cosine", // Default metric
	}
}

func (s server) convertVectorFilters(pbFilters *vectorspb.VectorFilters) map[string]interface{} {
	if pbFilters == nil {
		return make(map[string]interface{})
	}

	filters := make(map[string]interface{})

	if len(pbFilters.GetEntityTypes()) > 0 {
		filters["entity_types"] = pbFilters.GetEntityTypes()
	}
	if len(pbFilters.GetStatuses()) > 0 {
		filters["statuses"] = pbFilters.GetStatuses()
	}
	if len(pbFilters.GetCategories()) > 0 {
		filters["categories"] = pbFilters.GetCategories()
	}
	if pbFilters.GetNegotiableOnly() {
		filters["negotiable_only"] = true
	}
	if pbFilters.GetUserType() != "" {
		filters["user_type"] = pbFilters.GetUserType()
	}

	// Add price range filter
	if pbFilters.GetPriceRange() != nil {
		filters["min_price"] = pbFilters.GetPriceRange().GetMinPrice()
		filters["max_price"] = pbFilters.GetPriceRange().GetMaxPrice()
	}

	// Add geo filter
	if pbFilters.GetGeoFilter() != nil {
		filters["latitude"] = pbFilters.GetGeoFilter().GetLatitude()
		filters["longitude"] = pbFilters.GetGeoFilter().GetLongitude()
		filters["radius_km"] = pbFilters.GetGeoFilter().GetRadiusKm()
		filters["country"] = pbFilters.GetGeoFilter().GetCountry()
		filters["city"] = pbFilters.GetGeoFilter().GetCity()
	}

	// Add metadata filters
	for k, v := range pbFilters.GetMetadataFilters() {
		filters[k] = v
	}

	return filters
}
