package vector

import (
	"context"
	"fmt"
	"log"
	"strings"

	"middleman/vectors/internal/models"
)

// IntegrationService handles entity indexing into the vector database
type IntegrationService struct {
	vectorService    *VectorService
	embeddingService *EmbeddingService
}

// NewIntegrationService creates a new integration service
func NewIntegrationService(vectorService *VectorService, embeddingService *EmbeddingService) *IntegrationService {
	return &IntegrationService{
		vectorService:    vectorService,
		embeddingService: embeddingService,
	}
}

// IndexEntity indexes any entity into the vector database
func (is *IntegrationService) IndexEntity(ctx context.Context, entityID string, entityType string, entityData map[string]interface{}) error {
	// Generate embedding from entity data
	embedding, err := is.embeddingService.GenerateEntityEmbedding(ctx, entityData)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Create vector point
	point := VectorPoint{
		ID:       entityID,
		Vector:   embedding,
		Metadata: is.enrichMetadata(entityType, entityData),
	}

	// Index in vector database
	return is.vectorService.IndexVector(ctx, point)
}

// BatchIndexEntities indexes multiple entities efficiently
func (is *IntegrationService) BatchIndexEntities(ctx context.Context, entities []EntityIndexRequest) error {
	if len(entities) == 0 {
		return nil
	}

	points := make([]VectorPoint, 0, len(entities))

	for _, entity := range entities {
		embedding, err := is.embeddingService.GenerateEntityEmbedding(ctx, entity.Data)
		if err != nil {
			log.Printf("Failed to generate embedding for entity %s: %v", entity.ID, err)
			continue
		}

		points = append(points, VectorPoint{
			ID:       entity.ID,
			Vector:   embedding,
			Metadata: is.enrichMetadata(entity.Type, entity.Data),
		})
	}

	return is.vectorService.BatchIndexVectors(ctx, points)
}

// IndexProduct indexes a product entity
func (is *IntegrationService) IndexProduct(ctx context.Context, product *models.Product) error {
	data := map[string]interface{}{
		"name":          product.Name,
		"description":   product.Description,
		"brand":         product.Brand,
		"model":         product.Model,
		"category_id":   product.CategoryID,
		"category_slug": product.CategorySlug,
		"tags":          product.Tags,
		"condition":     string(product.Condition),
		"price":         product.BasePrice,
		"user_id":       product.UserSellerID,
		"status":        string(product.Status),
		"negotiable":    product.Negotiable,
		"user_type":     string(product.UserType),
		"lat":           product.Lat,
		"lng":           product.Lng,
	}

	return is.IndexEntity(ctx, product.ProductID, "product", data)
}

// IndexPost indexes a post entity
func (is *IntegrationService) IndexPost(ctx context.Context, post *models.Post) error {
	data := map[string]interface{}{
		"name":        post.Name,
		"description": post.Description,
		"tags":        post.Tags,
		"user_id":     post.UserID,
		"status":      string(post.Status),
		"lat":         post.Lat,
		"lng":         post.Lng,
	}

	return is.IndexEntity(ctx, post.PostID, "post", data)
}

// IndexService indexes a service entity
func (is *IntegrationService) IndexService(ctx context.Context, service *models.Service) error {
	data := map[string]interface{}{
		"name":           service.Name,
		"description":    service.Description,
		"provider_name":  service.ProviderName,
		"category_id":    service.CategoryID,
		"category_slug":  service.CategorySlug,
		"tags":           service.Tags,
		"service_type":   service.ServiceType,
		"qualifications": service.Qualifications,
		"price":          service.BasePrice,
		"user_id":        service.UserID,
		"status":         string(service.Status),
		"negotiable":     service.Negotiable,
		"user_type":      string(service.UserType),
		"lat":            service.Lat,
		"lng":            service.Lng,
	}

	return is.IndexEntity(ctx, service.ID, "service", data)
}

// RemoveEntity removes an entity from the vector database
func (is *IntegrationService) RemoveEntity(ctx context.Context, entityID string) error {
	return is.vectorService.DeleteVector(ctx, entityID)
}

// BatchRemoveEntities removes multiple entities from the vector database
func (is *IntegrationService) BatchRemoveEntities(ctx context.Context, entityIDs []string) error {
	return is.vectorService.BatchDeleteVectors(ctx, entityIDs)
}

// Helper types and methods

// EntityIndexRequest represents a request to index an entity
type EntityIndexRequest struct {
	ID   string                 `json:"id"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// enrichMetadata adds additional metadata to entity data
func (is *IntegrationService) enrichMetadata(entityType string, data map[string]interface{}) map[string]interface{} {
	metadata := make(map[string]interface{})

	// Copy original data
	for k, v := range data {
		metadata[k] = v
	}

	// Add entity type
	metadata["entity_type"] = entityType

	// Add searchable text field
	metadata["searchable_text"] = is.buildSearchableText(data)

	return metadata
}

// buildSearchableText creates a searchable text representation of entity data
func (is *IntegrationService) buildSearchableText(data map[string]interface{}) string {
	var parts []string

	// Priority fields for search
	searchFields := []string{"name", "title", "description", "brand", "model", "tags", "category_slug"}

	for _, field := range searchFields {
		if value, exists := data[field]; exists && value != nil {
			if str := is.valueToSearchString(value); str != "" {
				parts = append(parts, str)
			}
		}
	}

	return strings.Join(parts, " ")
}

func (is *IntegrationService) valueToSearchString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case []string:
		return strings.Join(v, " ")
	default:
		return ""
	}
}
