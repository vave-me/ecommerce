package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"middleman/internal/ddd"
	"middleman/managers/internal/application/services"
	"middleman/managers/internal/domain"
)

// ProcessImageInput represents the command to process image input for an manager.
type ProcessImageInput struct {
	ID           string                 `json:"id"`                      // Unique ID for this specific request/interaction
	ManagerID    string                 `json:"manager_id"`              // ID of the manager to use
	UserID       string                 `json:"user_id"`                 // ID of the user making the request
	ImageData    []byte                 `json:"image_data,omitempty"`    // Raw image data (JPEG, PNG, WebP, etc.)
	ImageURL     string                 `json:"image_url,omitempty"`     // URL to image (alternative to ImageData)
	ImageFormat  string                 `json:"image_format,omitempty"`  // Image format (jpeg, png, webp, gif, etc.)
	AnalysisType string                 `json:"analysis_type,omitempty"` // Type of analysis (describe, ocr, objects, faces, text, analyze, etc.)
	UserPrompt   string                 `json:"user_prompt,omitempty"`   // User's specific request about the image
	Context      map[string]interface{} `json:"context,omitempty"`       // Additional context
	Timestamp    time.Time              `json:"timestamp,omitempty"`     // Timestamp of the request
	RequestType  string                 `json:"request_type,omitempty"`  // Type of request (e.g., "image_analysis", "image_search", "vision_command")

	// Extended fields (used by gRPC server for advanced workflows)
	ProcessingMode string                 `json:"processing_mode,omitempty"` // Workflow mode detected (e.g., "analyze_image", "attach_image")
	ListingData    map[string]interface{} `json:"listing_data,omitempty"`    // Optional data extracted for listing creation / attachment
}

// ProcessImageInputResult holds the structured result of processing image input.
type ProcessImageInputResult struct {
	ResponseID         string
	AnalysisResult     string // The analysis result from vision AI
	ResponseMessage    string // AI manager response
	ResponseStatus     string
	ResponseConfidence float64 // Combined confidence (Vision + LLM)
	VisionConfidence   float64 // Vision analysis confidence
	LLMConfidence      float64 // LLM processing confidence
	ResponseTimestamp  time.Time
	ImageMetadata      map[string]interface{} // Image dimensions, format, etc.
	ExecutedActions    []domain.ManagerAction // Actions executed during processing
	ImageFormat        string                 // Image format used
	AnalysisType       string                 // Analysis type performed
	ProcessingTime     time.Duration          // Processing time
	InputSource        string                 // "binary_data" or "url"

	// Extended result fields expected by gRPC server
	ProcessingMode   string `json:"processing_mode,omitempty"`    // Echo of the processing mode that was executed
	ImageAttached    bool   `json:"image_attached,omitempty"`     // Indicates if the image was attached to a listing
	CreatedListingID string `json:"created_listing_id,omitempty"` // ID of the newly created listing (if any)
}

// ProcessImageInputHandler orchestrates image analysis and subsequent tool execution.
type ProcessImageInputHandler struct {
	managers        domain.ManagerRepository
	publisher       ddd.EventPublisher[ddd.Event]
	llmProcessor    services.LLMProcessor
	visionProcessor services.VisionProcessor // New interface for image analysis
}

// NewProcessImageInputHandler creates a new image input handler.
func NewProcessImageInputHandler(
	managers domain.ManagerRepository,
	visionProcessor services.VisionProcessor,
	llmProcessor services.LLMProcessor,
	publisher ddd.EventPublisher[ddd.Event],
) (ProcessImageInputHandler, error) {
	if managers == nil {
		return ProcessImageInputHandler{}, errors.New("managers repository cannot be nil")
	}
	if visionProcessor == nil {
		return ProcessImageInputHandler{}, errors.New("vision processor cannot be nil")
	}
	if llmProcessor == nil {
		return ProcessImageInputHandler{}, errors.New("LLM processor cannot be nil")
	}
	if publisher == nil {
		return ProcessImageInputHandler{}, errors.New("event publisher cannot be nil")
	}
	return ProcessImageInputHandler{
		managers:        managers,
		visionProcessor: visionProcessor,
		llmProcessor:    llmProcessor,
		publisher:       publisher,
	}, nil
}

// ProcessImageInput handles image analysis requests with AI vision processing
func (h *ProcessImageInputHandler) ProcessImageInput(ctx context.Context, cmd ProcessImageInput) (*ProcessImageInputResult, error) {
	startTime := time.Now()

	// Validate required fields
	if cmd.ManagerID == "" {
		return nil, fmt.Errorf("manager_id is required")
	}
	if cmd.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Validate image input - must have either data or URL
	if len(cmd.ImageData) == 0 && cmd.ImageURL == "" {
		return nil, fmt.Errorf("either image_data or image_url must be provided")
	}
	if len(cmd.ImageData) > 0 && cmd.ImageURL != "" {
		return nil, fmt.Errorf("provide either image_data or image_url, not both")
	}

	// Set default analysis type if not provided
	analysisType := cmd.AnalysisType
	if analysisType == "" {
		analysisType = "analyze"
	}

	// Load the manager
	manager, err := h.managers.Load(ctx, cmd.ManagerID)
	if err != nil {
		return nil, fmt.Errorf("failed to load manager: %w", err)
	}

	// Prepare vision analysis request
	visionRequest := services.VisionAnalysisRequest{
		ImageData:    cmd.ImageData, // May be empty if using URL
		ImageURL:     cmd.ImageURL,  // May be empty if using binary data
		ImageFormat:  cmd.ImageFormat,
		AnalysisType: analysisType,
		UserPrompt:   cmd.UserPrompt,
		Context:      cmd.Context,
		MaxTokens:    2000,   // Default for vision analysis
	}

	// Process image with vision AI
	visionResult, err := h.visionProcessor.AnalyzeImage(ctx, visionRequest)
	if err != nil {
		return &ProcessImageInputResult{
			ResponseID:         fmt.Sprintf("img_error_%d", time.Now().UnixNano()),
			ResponseMessage:    "Image analysis failed",
			ResponseStatus:     "error",
			ResponseTimestamp:  time.Now(),
			AnalysisResult:     "",
			ImageFormat:        cmd.ImageFormat,
			AnalysisType:       analysisType,
			VisionConfidence:   0.0,
			LLMConfidence:      0.0,
			ResponseConfidence: 0.0,
			ProcessingTime:     time.Since(startTime),
			InputSource:        getInputSource(cmd),
			ProcessingMode:     cmd.ProcessingMode,
			ImageAttached:      false,
			CreatedListingID:   "",
		}, nil
	}

	// Enhanced prompt for LLM with vision context
	enhancedPrompt := h.buildEnhancedPrompt(cmd.UserPrompt, visionResult, analysisType)

	// Process with LLM for tool execution and enhanced response
	response, actions, llmConfidence, err := h.llmProcessor.ProcessWithHistory(
		ctx,
		manager,
		enhancedPrompt,
		[]domain.ConversationMessage{}, // No history for simple image processing
		cmd.Context,
	)

	if err != nil {
		// Continue with vision-only result
		response = visionResult.Description
		actions = []domain.ManagerAction{}
		llmConfidence = 0.7
	}

	// Calculate combined confidence score
	combinedConfidence := (visionResult.Confidence + llmConfidence) / 2.0
	if combinedConfidence > 1.0 {
		combinedConfidence = 1.0
	}

	// Save manager state
	if err := h.managers.Save(ctx, manager); err != nil {
		// Non-critical error, continue
	}

	return &ProcessImageInputResult{
		ResponseID:         fmt.Sprintf("img_%d", time.Now().UnixNano()),
		ResponseMessage:    response,
		ResponseStatus:     "success",
		ResponseTimestamp:  time.Now(),
		ResponseConfidence: combinedConfidence,
		AnalysisResult:     visionResult.Description,
		ExecutedActions:    actions,
		ImageFormat:        cmd.ImageFormat,
		AnalysisType:       analysisType,
		VisionConfidence:   visionResult.Confidence,
		LLMConfidence:      llmConfidence,
		ProcessingTime:     time.Since(startTime),
		InputSource:        getInputSource(cmd),
		ProcessingMode:     cmd.ProcessingMode,
		ImageAttached:      false,
		CreatedListingID:   "",
	}, nil
}

// buildEnhancedPrompt creates an enhanced prompt combining user request with vision analysis
func (h *ProcessImageInputHandler) buildEnhancedPrompt(userPrompt string, visionResult *services.VisionAnalysisResult, analysisType string) string {
	visionPart := fmt.Sprintf("Image analysis result: %s", visionResult.Description)

	// If the detected processing mode implies listing creation, guide the LLM explicitly
	toolHint := ""
	if strings.Contains(strings.ToLower(analysisType), "attach") || strings.Contains(strings.ToLower(analysisType), "listing") {
		toolHint = `\n\nIf the user intends to sell the item, respond with a JSON object using the following schema *only* if all required fields are available: \n{\n  "name": "create_listing",\n  "arguments": {\n    "title": string,\n    "description": string,\n    "price": number,\n    "currency": "EUR",\n    "image_id": string\n  }\n}\nOtherwise answer normally.`
	}

	if userPrompt == "" {
		return fmt.Sprintf("I analyzed this image using %s analysis. %s%s", analysisType, visionPart, toolHint)
	}
	return fmt.Sprintf("%s\n\n%s%s", userPrompt, visionPart, toolHint)
}

// getInputSource determines whether binary data or URL was used
func getInputSource(cmd ProcessImageInput) string {
	if cmd.ImageURL != "" {
		return "url"
	}
	return "binary_data"
}
