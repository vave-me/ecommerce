package vector

import (
	"context"
	"fmt"
	"middleman/vectors/internal/constants"
	"strings"
)

// EmbeddingService generates vector embeddings for text content using LLM-guided transformation
type EmbeddingService struct {
	dimensions    int
	model         string
	promptEnabled bool
	llmClient     LLMClient // Interface for LLM interactions
}

// LLMClient interface for interacting with language models
type LLMClient interface {
	Transform(ctx context.Context, text string, prompt string) (string, error)
	IsAvailable() bool
}

// EmbeddingConfig holds embedding service configuration
type EmbeddingConfig struct {
	Model         string    // e.g., "text-embedding-3-small", "all-MiniLM-L6-v2"
	Dimensions    int       // Vector dimensions (default: 384)
	APIKey        string    // For external embedding services
	PromptEnabled bool      // Enable LLM-guided transformation
	LLMClient     LLMClient // Optional LLM client for text transformation
}

// TransformationStrategy defines how to transform entity data
type TransformationStrategy struct {
	BasePrompt        string
	EnhancementPrompt string
	ContextType       string
	PerformanceLevel  string
	QualityTarget     string
}

// NewEmbeddingService creates a new embedding service with prompt capabilities
func NewEmbeddingService(config EmbeddingConfig) *EmbeddingService {
	dimensions := config.Dimensions
	if dimensions <= 0 {
		dimensions = 384 // Default for many open-source models
	}

	model := config.Model
	if model == "" {
		model = "all-MiniLM-L6-v2"
	}

	return &EmbeddingService{
		dimensions:    dimensions,
		model:         model,
		promptEnabled: config.PromptEnabled,
		llmClient:     config.LLMClient,
	}
}

// GenerateEmbedding converts text to vector embedding
func (e *EmbeddingService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if text = strings.TrimSpace(text); text == "" {
		return nil, fmt.Errorf("empty text provided")
	}

	// Generate deterministic but varied embedding based on text content
	// In production, this would call an actual embedding model/API
	embedding := make([]float32, e.dimensions)
	hash := e.textToHash(text)

	for i := 0; i < e.dimensions; i++ {
		// Generate pseudo-random float between -1 and 1 with better distribution
		hash = hash*1103515245 + 12345
		normalized := float32((float64(hash%2000000) - 1000000) / 1000000.0)
		embedding[i] = normalized * 0.8 // Scale down to avoid extreme values
	}

	// Normalize the vector
	return e.normalizeVector(embedding), nil
}

// GenerateEntityEmbeddingWithPrompt creates embeddings using LLM-guided transformation
func (e *EmbeddingService) GenerateEntityEmbeddingWithPrompt(
	ctx context.Context,
	entityType string,
	entityData map[string]interface{},
	strategy TransformationStrategy,
) ([]float32, error) {

	// If prompts are disabled or no LLM client, fall back to basic transformation
	if !e.promptEnabled || e.llmClient == nil || !e.llmClient.IsAvailable() {
		return e.GenerateEntityEmbedding(ctx, entityData)
	}

	// Build comprehensive prompt
	prompt := e.buildTransformationPrompt(entityType, strategy)

	// Convert entity data to initial text representation
	entityText := e.entityToText(entityData)

	// Use LLM to transform the text according to the prompt
	transformedText, err := e.llmClient.Transform(ctx, entityText, prompt)
	if err != nil {
		// Fallback to basic transformation on LLM error
		return e.GenerateEntityEmbedding(ctx, entityData)
	}

	// Generate embedding from the LLM-transformed text
	return e.GenerateEmbedding(ctx, transformedText)
}

// GenerateOptimizedEmbedding creates embeddings optimized for specific use cases
func (e *EmbeddingService) GenerateOptimizedEmbedding(
	ctx context.Context,
	entityType string,
	entityData map[string]interface{},
	optimization string, // "search", "recommendation", "similarity", "real-time"
) ([]float32, error) {

	strategy := e.buildOptimizationStrategy(entityType, optimization)
	return e.GenerateEntityEmbeddingWithPrompt(ctx, entityType, entityData, strategy)
}

// GenerateBatchEmbeddingsWithStrategy generates embeddings for multiple entities using consistent strategy
func (e *EmbeddingService) GenerateBatchEmbeddingsWithStrategy(
	ctx context.Context,
	entities []EntityData,
	strategy TransformationStrategy,
) ([][]float32, error) {

	if len(entities) == 0 {
		return nil, fmt.Errorf("no entities provided")
	}

	embeddings := make([][]float32, 0, len(entities))
	for _, entity := range entities {
		embedding, err := e.GenerateEntityEmbeddingWithPrompt(
			ctx,
			entity.Type,
			entity.Data,
			strategy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for entity %s: %w", entity.ID, err)
		}
		embeddings = append(embeddings, embedding)
	}

	return embeddings, nil
}

// EntityData represents an entity for batch processing
type EntityData struct {
	ID   string
	Type string
	Data map[string]interface{}
}

// buildTransformationPrompt constructs a comprehensive prompt for entity transformation
func (e *EmbeddingService) buildTransformationPrompt(entityType string, strategy TransformationStrategy) string {
	// Get base prompt for entity type
	basePrompt := e.getBasePromptForEntity(entityType)

	// If strategy specifies a base prompt, use it
	if strategy.BasePrompt != "" {
		basePrompt = strategy.BasePrompt
	}

	// Build final prompt
	prompt := basePrompt

	// Add enhancement prompt if specified
	if strategy.EnhancementPrompt != "" {
		enhancedTemplate := strings.Replace(
			constants.BaseEnhancedPrompt,
			"{BASE_PROMPT}",
			basePrompt,
			1,
		)
		prompt = strings.Replace(
			enhancedTemplate,
			"{ENHANCEMENT_PROMPT}",
			strategy.EnhancementPrompt,
			1,
		)
	}

	// Add context adaptation if specified
	if strategy.ContextType != "" {
		contextTemplate := strings.Replace(
			constants.ContextAdaptationPrompt,
			"{BASE_PROMPT}",
			prompt,
			1,
		)
		contextTemplate = strings.Replace(
			contextTemplate,
			"{CONTEXT_TYPE}",
			strategy.ContextType,
			1,
		)
		prompt = contextTemplate
	}

	return prompt
}

// getBasePromptForEntity returns the appropriate base transformation prompt for an entity type
func (e *EmbeddingService) getBasePromptForEntity(entityType string) string {
	switch strings.ToLower(entityType) {
	case "product", "products":
		return constants.ProductVectorPrompt
	case "job", "jobs":
		return constants.JobVectorPrompt
	case "vehicle", "vehicles":
		return constants.VehicleVectorPrompt
	case "property", "properties":
		return constants.PropertyVectorPrompt
	case "post", "posts":
		return constants.PostVectorPrompt
	case "deal", "deals":
		return constants.DealVectorPrompt
	case "service", "services":
		return constants.ServiceVectorPrompt
	case "user", "users":
		return constants.UserVectorPrompt
	default:
		// Generic prompt for unknown entity types
		return `Transform this entity into a comprehensive text representation optimized for vector similarity search.
Focus on key attributes, unique characteristics, and searchable qualities.
Emphasize features that enable effective discovery and matching.
Output format: A rich, descriptive text highlighting the entity's essential qualities.`
	}
}

// buildOptimizationStrategy creates a transformation strategy for specific optimization goals
func (e *EmbeddingService) buildOptimizationStrategy(entityType, optimization string) TransformationStrategy {
	strategy := TransformationStrategy{}

	switch strings.ToLower(optimization) {
	case "search", "discovery":
		strategy.EnhancementPrompt = constants.SearchOptimizedPrompt
		strategy.PerformanceLevel = "balanced"

	case "recommendation", "suggestions":
		strategy.EnhancementPrompt = constants.RecommendationPrompt
		strategy.QualityTarget = "high_relevance"

	case "similarity", "matching":
		strategy.EnhancementPrompt = constants.SimilarityMatchingPrompt
		strategy.QualityTarget = "high_precision"

	case "real-time", "fast":
		strategy.EnhancementPrompt = constants.RealTimeOptimizationPrompt
		strategy.PerformanceLevel = "speed_optimized"

	case "quality", "premium":
		strategy.EnhancementPrompt = constants.EmbeddingQualityPrompt
		strategy.QualityTarget = "maximum_quality"

	case "ecommerce", "shopping":
		strategy.EnhancementPrompt = constants.EcommercePrompt
		strategy.ContextType = "ECOMMERCE"

	case "social", "community":
		strategy.EnhancementPrompt = constants.SocialPlatformPrompt
		strategy.ContextType = "SOCIAL"

	case "professional", "business":
		strategy.EnhancementPrompt = constants.ProfessionalNetworkingPrompt
		strategy.ContextType = "PROFESSIONAL"

	case "educational", "learning":
		strategy.EnhancementPrompt = constants.EducationalPrompt
		strategy.ContextType = "EDUCATIONAL"

	default:
		// Balanced approach for unknown optimizations
		strategy.EnhancementPrompt = constants.BalancedPrompt
		strategy.PerformanceLevel = "balanced"
	}

	return strategy
}

// GenerateBatchEmbeddings generates embeddings for multiple texts efficiently
func (e *EmbeddingService) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("no texts provided")
	}

	embeddings := make([][]float32, 0, len(texts))
	for _, text := range texts {
		embedding, err := e.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for text: %w", err)
		}
		embeddings = append(embeddings, embedding)
	}

	return embeddings, nil
}

// GenerateEntityEmbedding creates embeddings for structured entity data
func (e *EmbeddingService) GenerateEntityEmbedding(ctx context.Context, entityData map[string]interface{}) ([]float32, error) {
	text := e.entityToText(entityData)
	return e.GenerateEmbedding(ctx, text)
}

// GetDimensions returns the embedding dimensions
func (e *EmbeddingService) GetDimensions() int {
	return e.dimensions
}

// GetModel returns the embedding model name
func (e *EmbeddingService) GetModel() string {
	return e.model
}

// IsPromptEnabled returns whether LLM-guided transformation is enabled
func (e *EmbeddingService) IsPromptEnabled() bool {
	return e.promptEnabled && e.llmClient != nil && e.llmClient.IsAvailable()
}

// Helper methods

func (e *EmbeddingService) textToHash(text string) uint32 {
	var hash uint32 = 5381
	for _, char := range text {
		hash = ((hash << 5) + hash) + uint32(char)
	}
	return hash
}

func (e *EmbeddingService) normalizeVector(vector []float32) []float32 {
	var magnitude float32
	for _, value := range vector {
		magnitude += value * value
	}

	if magnitude == 0 {
		return vector
	}

	magnitude = float32(1.0 / (magnitude * magnitude)) // Fast inverse square root approximation

	normalized := make([]float32, len(vector))
	for i, value := range vector {
		normalized[i] = value * magnitude
	}

	return normalized
}

func (e *EmbeddingService) entityToText(entityData map[string]interface{}) string {
	var parts []string

	// Extract important fields in priority order
	priority := []string{"name", "title", "description", "content", "brand", "model", "category", "tags"}

	for _, field := range priority {
		if value, exists := entityData[field]; exists && value != nil {
			if str := e.valueToString(value); str != "" {
				parts = append(parts, str)
			}
		}
	}

	// Add remaining fields
	for key, value := range entityData {
		if !e.isPriorityField(key, priority) && value != nil {
			if str := e.valueToString(value); str != "" {
				parts = append(parts, fmt.Sprintf("%s: %s", key, str))
			}
		}
	}

	return strings.Join(parts, " ")
}

func (e *EmbeddingService) valueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case []string:
		return strings.Join(v, " ")
	case int, int64, float32, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func (e *EmbeddingService) isPriorityField(field string, priority []string) bool {
	for _, p := range priority {
		if field == p {
			return true
		}
	}
	return false
}

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

// Ensure service implements the interface
var _ EmbeddingInterface = (*EmbeddingService)(nil)
