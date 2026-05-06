package processor

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"middleman/assistants/internal/application/services"
	"middleman/internal/ai"
)

const (
	logPrefixVisionProcessor = "[OpenAIVisionProcessor]"

	// Supported image formats
	ImageFormatJPEG = "jpeg"
	ImageFormatJPG  = "jpg"
	ImageFormatPNG  = "png"
	ImageFormatWEBP = "webp"
	ImageFormatGIF  = "gif"

	// OpenAI Vision models - Updated for 2025
	VisionModel     = "gpt-4.1-mini"      // Fast, cost-effective GPT-4.1 mini
	VisionModelMini = "gpt-4o-mini"       // Previous generation mini
	VisionModelNew  = "gpt-4o-2024-11-20" // Latest GPT-4o version

	// Maximum image file size (20MB for OpenAI Vision)
	MaxImageSizeBytes = 20 * 1024 * 1024

	// Default parameters
	DefaultMaxTokens   = 500
	DefaultDetailLevel = "auto" // "low", "high", "auto"
)

// OpenAIVisionProcessor implements the VisionProcessor interface using OpenAI's Vision API
type OpenAIVisionProcessor struct {
	aiClient ai.EnhancedAIService // Use EnhancedAIService for consistency
	config   *VisionProcessorConfig
}

// VisionProcessorConfig holds configuration for vision processing
type VisionProcessorConfig struct {
	Model                 string        `json:"model"`                   // Vision model to use
	DefaultDetailLevel    string        `json:"default_detail_level"`    // Default detail level
	MaxTokens             int           `json:"max_tokens"`              // Maximum tokens for response
	Temperature           float64       `json:"temperature"`             // Temperature for generation
	EnableOCR             bool          `json:"enable_ocr"`              // Whether to enable OCR
	EnableObjectDetection bool          `json:"enable_object_detection"` // Whether to detect objects
	EnableFaceDetection   bool          `json:"enable_face_detection"`   // Whether to detect faces
	MaxRetries            int           `json:"max_retries"`             // Maximum retry attempts
	RetryDelay            time.Duration `json:"retry_delay"`             // Delay between retries
}

// NewOpenAIVisionProcessor creates a new OpenAI vision processor
func NewOpenAIVisionProcessor(aiClient ai.EnhancedAIService, config *VisionProcessorConfig) services.VisionProcessor {
	if aiClient == nil {
		log.Printf("%s ERROR: AI client cannot be nil, returning nil processor", logPrefixVisionProcessor)
		return nil
	}

	if config == nil {
		config = &VisionProcessorConfig{
			Model:                 VisionModel, // Use GPT-4.1 mini
			DefaultDetailLevel:    DefaultDetailLevel,
			MaxTokens:             DefaultMaxTokens,
			Temperature:           0.1,
			EnableOCR:             true,
			EnableObjectDetection: true,
			EnableFaceDetection:   false, // Privacy consideration
			MaxRetries:            3,
			RetryDelay:            1 * time.Second,
		}
	}

	return &OpenAIVisionProcessor{
		aiClient: aiClient,
		config:   config,
	}
}

// AnalyzeImage analyzes image content and extracts information using OpenAI's Vision API
func (p *OpenAIVisionProcessor) AnalyzeImage(ctx context.Context, request services.VisionAnalysisRequest) (*services.VisionAnalysisResult, error) {
	startTime := time.Now()
	log.Printf("%s Starting image analysis. Format: %s, Size: %d bytes, Type: %s",
		logPrefixVisionProcessor, request.ImageFormat, len(request.ImageData), request.AnalysisType)

	// Validate input
	if err := p.validateAnalysisRequest(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Build OpenAI vision request using proper 2025 format
	completionRequest := p.buildVisionRequest(request)

	// Execute analysis with retries
	response, err := p.analyzeWithRetries(ctx, completionRequest)
	if err != nil {
		log.Printf("%s Vision analysis failed: %v", logPrefixVisionProcessor, err)
		return nil, fmt.Errorf("vision analysis failed: %w", err)
	}

	// Convert OpenAI response to our format
	result := p.convertToVisionResult(response, request, startTime)

	log.Printf("%s Vision analysis completed successfully. Description: '%s', Confidence: %.2f, Duration: %v",
		logPrefixVisionProcessor, result.Description, result.Confidence, time.Since(startTime))

	return result, nil
}

// ValidateImageFormat checks if the image format is supported
func (p *OpenAIVisionProcessor) ValidateImageFormat(format string) error {
	format = strings.ToLower(format)
	supportedFormats := []string{ImageFormatJPEG, ImageFormatJPG, ImageFormatPNG, ImageFormatWEBP, ImageFormatGIF}

	for _, supported := range supportedFormats {
		if format == supported {
			return nil
		}
	}

	return fmt.Errorf("unsupported image format: %s. Supported formats: %v", format, supportedFormats)
}

// GetSupportedAnalysisTypes returns list of supported analysis types
func (p *OpenAIVisionProcessor) GetSupportedAnalysisTypes() []string {
	return []string{
		"describe", // General description
		"analyze",  // Detailed analysis
		"ocr",      // Text extraction
		"objects",  // Object detection
		"text",     // Text reading
		"classify", // Image classification
		"compare",  // Image comparison
		"explain",  // Explain what's happening
		"count",    // Count objects
		"identify", // Identify specific items
	}
}

// validateAnalysisRequest validates the vision analysis request
func (p *OpenAIVisionProcessor) validateAnalysisRequest(request services.VisionAnalysisRequest) error {
	// Must have either image data or URL, but not both
	if len(request.ImageData) == 0 && request.ImageURL == "" {
		return fmt.Errorf("either image data or image URL must be provided")
	}

	if len(request.ImageData) > 0 && request.ImageURL != "" {
		return fmt.Errorf("provide either image data or image URL, not both")
	}

	// Validate binary data if provided
	if len(request.ImageData) > 0 {
		if len(request.ImageData) > MaxImageSizeBytes {
			return fmt.Errorf("image size %d exceeds maximum allowed size %d", len(request.ImageData), MaxImageSizeBytes)
		}

		if err := p.ValidateImageFormat(request.ImageFormat); err != nil {
			return err
		}
	}

	// Validate URL if provided
	if request.ImageURL != "" {
		// Basic URL validation - more sophisticated validation could be added
		if !strings.HasPrefix(request.ImageURL, "http://") && !strings.HasPrefix(request.ImageURL, "https://") {
			return fmt.Errorf("image URL must be a valid HTTP or HTTPS URL")
		}
	}

	return nil
}

// buildVisionRequest creates an OpenAI vision request using 2025 API format
func (p *OpenAIVisionProcessor) buildVisionRequest(request services.VisionAnalysisRequest) ai.CompletionRequest {
	// Build the prompt based on analysis type and user prompt
	prompt := p.buildAnalysisPrompt(request.AnalysisType, request.UserPrompt)

	// Determine detail level - OpenAI supports "low", "high", "auto"
	detailLevel := p.config.DefaultDetailLevel

	// Determine max tokens
	maxTokens := request.MaxTokens
	if maxTokens <= 0 {
		maxTokens = p.config.MaxTokens
	}

	// Create image URL content - handle both binary data and URL inputs
	var imageContent map[string]interface{}

	if request.ImageURL != "" {
		// Use direct URL (public or signed URL)
		imageContent = map[string]interface{}{
			"type": "image_url",
			"image_url": map[string]interface{}{
				"url":    request.ImageURL,
				"detail": detailLevel,
			},
		}
	} else {
		// Use base64 encoded data URI for binary data
		imageContent = map[string]interface{}{
			"type": "image_url",
			"image_url": map[string]interface{}{
				"url":    fmt.Sprintf("data:image/%s;base64,%s", request.ImageFormat, p.encodeBase64(request.ImageData)),
				"detail": detailLevel,
			},
		}
	}

	// Create multimodal content using 2025 OpenAI format
	content := []interface{}{
		map[string]interface{}{
			"type": "text",
			"text": prompt,
		},
		imageContent,
	}

	// Create message with multimodal content
	message := ai.Message{
		Role:    "user",
		Content: content,
	}

	return ai.CompletionRequest{
		Model:       p.config.Model,
		Messages:    []ai.Message{message},
		MaxTokens:   &maxTokens,
		Temperature: &p.config.Temperature,
	}
}

// buildAnalysisPrompt creates a prompt based on analysis type and user request
func (p *OpenAIVisionProcessor) buildAnalysisPrompt(analysisType, userPrompt string) string {
	basePrompts := map[string]string{
		"describe": "Describe what you see in this image in detail.",
		"analyze":  "Provide a comprehensive analysis of this image, including objects, people, activities, setting, and any notable details.",
		"ocr":      "Extract and transcribe all text visible in this image. Maintain the original formatting and layout as much as possible.",
		"objects":  "Identify and list all objects visible in this image, including their approximate locations.",
		"text":     "Read and extract any text content from this image, including signs, labels, documents, or written material.",
		"classify": "Classify this image by category and provide relevant tags or labels.",
		"explain":  "Explain what is happening in this image, including any activities, interactions, or processes.",
		"count":    "Count and enumerate the different types of objects or people visible in this image.",
		"identify": "Identify and name specific items, brands, landmarks, or notable features in this image.",
	}

	basePrompt := basePrompts[analysisType]
	if basePrompt == "" {
		basePrompt = "Analyze and describe this image."
	}

	if userPrompt != "" {
		return fmt.Sprintf("%s\n\nSpecific request: %s", basePrompt, userPrompt)
	}

	return basePrompt
}

// analyzeWithRetries performs vision analysis with retry logic
func (p *OpenAIVisionProcessor) analyzeWithRetries(ctx context.Context, request ai.CompletionRequest) (*ai.CompletionResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("%s Retry attempt %d/%d", logPrefixVisionProcessor, attempt, p.config.MaxRetries)
			time.Sleep(p.config.RetryDelay * time.Duration(attempt))
		}

		response, err := p.aiClient.CreateCompletion(ctx, request)
		if err == nil {
			return response, nil
		}

		lastErr = err
		log.Printf("%s Attempt %d failed: %v", logPrefixVisionProcessor, attempt+1, err)
	}

	return nil, fmt.Errorf("all retry attempts failed: %w", lastErr)
}

// convertToVisionResult converts OpenAI response to our vision result format
func (p *OpenAIVisionProcessor) convertToVisionResult(response *ai.CompletionResponse, request services.VisionAnalysisRequest, startTime time.Time) *services.VisionAnalysisResult {
	// Extract main description
	description := ""
	confidence := 0.0

	if len(response.Choices) > 0 {
		description = response.Choices[0].Message.GetContentAsString()
		confidence = 0.9 // Default high confidence for successful response
	}

	// Basic metadata
	metadata := map[string]interface{}{
		"model":           p.config.Model,
		"analysis_type":   request.AnalysisType,
		"processing_time": time.Since(startTime).Seconds(),
		"image_format":    strings.ToUpper(request.ImageFormat),
		"file_size":       int64(len(request.ImageData)),
	}

	if response.Usage.TotalTokens > 0 {
		metadata["tokens_used"] = response.Usage.TotalTokens
		metadata["prompt_tokens"] = response.Usage.PromptTokens
		metadata["completion_tokens"] = response.Usage.CompletionTokens
	}

	// Parse detected objects, text, etc. from description if needed
	// This is a simplified implementation - in production, you might want
	// to use additional specialized models for object/text detection

	result := &services.VisionAnalysisResult{
		Description:     description,
		Confidence:      confidence,
		DetectedObjects: []services.DetectedObject{}, // Could be populated with specialized detection
		Metadata:        metadata,
	}

	return result
}

// extractCategories attempts to extract categories from the description
func (p *OpenAIVisionProcessor) extractCategories(description string) []string {
	// Simple keyword-based categorization - in production you might use more sophisticated NLP
	categories := []string{}
	description = strings.ToLower(description)

	categoryKeywords := map[string][]string{
		"people":     {"person", "people", "man", "woman", "child", "human"},
		"animals":    {"dog", "cat", "bird", "animal", "pet"},
		"vehicles":   {"car", "truck", "bike", "vehicle", "bus"},
		"buildings":  {"building", "house", "structure", "architecture"},
		"nature":     {"tree", "flower", "plant", "landscape", "outdoor"},
		"food":       {"food", "meal", "cooking", "kitchen", "restaurant"},
		"technology": {"computer", "phone", "device", "screen", "technology"},
		"indoor":     {"indoor", "inside", "room", "interior"},
		"outdoor":    {"outdoor", "outside", "exterior", "street"},
	}

	for category, keywords := range categoryKeywords {
		for _, keyword := range keywords {
			if strings.Contains(description, keyword) {
				categories = append(categories, category)
				break
			}
		}
	}

	return categories
}

// encodeBase64 encodes image data to base64 string
func (p *OpenAIVisionProcessor) encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
