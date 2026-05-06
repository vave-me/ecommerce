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

// ProcessDocumentInput represents the command to process document input for an manager.
type ProcessDocumentInput struct {
	ID             string                 `json:"id"`
	ManagerID      string                 `json:"manager_id"`
	UserID         string                 `json:"user_id"`
	DocumentData   []byte                 `json:"document_data,omitempty"`
	DocumentURL    string                 `json:"document_url,omitempty"`
	DocumentFormat string                 `json:"document_format,omitempty"`
	AnalysisType   string                 `json:"analysis_type,omitempty"`
	UserPrompt     string                 `json:"user_prompt,omitempty"`
	Context        map[string]interface{} `json:"context,omitempty"`
	Timestamp      time.Time              `json:"timestamp,omitempty"`
	RequestType    string                 `json:"request_type,omitempty"`
}

// ProcessDocumentInputResult holds the structured result of processing document input.
type ProcessDocumentInputResult struct {
	ResponseID           string                 `json:"response_id"`
	ExtractedContent     string                 `json:"extracted_content"`
	AnalysisResult       string                 `json:"analysis_result"`
	ResponseMessage      string                 `json:"response_message"`
	ResponseStatus       string                 `json:"response_status"`
	ResponseConfidence   float64                `json:"response_confidence"`
	ProcessingConfidence float64                `json:"processing_confidence"`
	LLMConfidence        float64                `json:"llm_confidence"`
	ResponseTimestamp    time.Time              `json:"response_timestamp"`
	DocumentMetadata     map[string]interface{} `json:"document_metadata,omitempty"`
	ExecutedActions      []domain.ManagerAction `json:"executed_actions,omitempty"`
	DocumentFormat       string                 `json:"document_format,omitempty"`
	AnalysisType         string                 `json:"analysis_type,omitempty"`
	ProcessingTime       time.Duration          `json:"processing_time"`
	InputSource          string                 `json:"input_source"`
	WordCount            int                    `json:"word_count"`
	PageCount            int                    `json:"page_count"`
}

// ProcessDocumentInputHandler orchestrates document processing and subsequent tool execution.
type ProcessDocumentInputHandler struct {
	managers          domain.ManagerRepository
	publisher         ddd.EventPublisher[ddd.Event]
	llmProcessor      services.LLMProcessor
	documentProcessor services.DocumentProcessor // Interface for document processing
	dataProcessor     services.DataProcessor     // Interface for structured data processing
}

// NewProcessDocumentInputHandler creates a new document input handler.
func NewProcessDocumentInputHandler(
	managers domain.ManagerRepository,
	documentProcessor services.DocumentProcessor,
	dataProcessor services.DataProcessor,
	llmProcessor services.LLMProcessor,
	publisher ddd.EventPublisher[ddd.Event],
) (ProcessDocumentInputHandler, error) {
	if managers == nil {
		return ProcessDocumentInputHandler{}, errors.New("managers repository cannot be nil")
	}
	if documentProcessor == nil {
		return ProcessDocumentInputHandler{}, errors.New("document processor cannot be nil")
	}
	if dataProcessor == nil {
		return ProcessDocumentInputHandler{}, errors.New("data processor cannot be nil")
	}
	if llmProcessor == nil {
		return ProcessDocumentInputHandler{}, errors.New("LLM processor cannot be nil")
	}
	if publisher == nil {
		return ProcessDocumentInputHandler{}, errors.New("event publisher cannot be nil")
	}
	return ProcessDocumentInputHandler{
		managers:          managers,
		documentProcessor: documentProcessor,
		dataProcessor:     dataProcessor,
		llmProcessor:      llmProcessor,
		publisher:         publisher,
	}, nil
}

// ProcessDocumentInput handles document processing requests with AI document analysis
func (h *ProcessDocumentInputHandler) ProcessDocumentInput(ctx context.Context, cmd ProcessDocumentInput) (*ProcessDocumentInputResult, error) {
	startTime := time.Now()

	// Validate required fields
	if cmd.ManagerID == "" {
		return nil, fmt.Errorf("manager_id is required")
	}
	if cmd.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Validate document input - must have either data or URL
	if len(cmd.DocumentData) == 0 && cmd.DocumentURL == "" {
		return nil, fmt.Errorf("either document_data or document_url must be provided")
	}
	if len(cmd.DocumentData) > 0 && cmd.DocumentURL != "" {
		return nil, fmt.Errorf("provide either document_data or document_url, not both")
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

	// Determine if this is structured data (CSV, JSON, etc.) or document (PDF, DOCX, etc.)
	isStructuredData := h.isStructuredDataFormat(cmd.DocumentFormat)

	var extractedContent string
	var analysisResult string
	var processingConfidence float64
	var documentMetadata map[string]interface{}
	var wordCount, pageCount int

	if isStructuredData {
		// Process as structured data

		dataRequest := services.DataProcessingRequest{
			DataContent:  cmd.DocumentData,
			DataURL:      cmd.DocumentURL,
			DataFormat:   cmd.DocumentFormat,
			AnalysisType: analysisType,
			UserPrompt:   cmd.UserPrompt,
			Context:      cmd.Context,
		}

		dataResult, err := h.dataProcessor.ProcessDataFile(ctx, dataRequest)
		if err != nil {
			return nil, fmt.Errorf("data processing failed: %w", err)
		}

		extractedContent = dataResult.Summary
		analysisResult = fmt.Sprintf("Data analysis completed. Rows: %d, Columns: %d", dataResult.RowCount, dataResult.ColumnCount)
		processingConfidence = 1.0 // Default confidence for data processing
		documentMetadata = dataResult.Metadata
		wordCount = len(extractedContent)
		pageCount = 1 // Data files are typically single "page"

	} else {
		// Process as document

		documentRequest := services.DocumentProcessingRequest{
			DocumentData:   cmd.DocumentData,
			DocumentURL:    cmd.DocumentURL,
			DocumentFormat: cmd.DocumentFormat,
			AnalysisType:   analysisType,
			UserPrompt:     cmd.UserPrompt,
			Context:        cmd.Context,
			MaxTokens:      2000,
			// Extract content in mixed mode by default
		}

		documentResult, err := h.documentProcessor.ProcessDocument(ctx, documentRequest)
		if err != nil {
			return nil, fmt.Errorf("document processing failed: %w", err)
		}

		extractedContent = documentResult.ExtractedText
		analysisResult = documentResult.Summary
		processingConfidence = 1.0 // Default confidence for document processing
		documentMetadata = documentResult.Metadata
		wordCount = len(documentResult.ExtractedText) // Approximate word count
		pageCount = documentResult.PageCount
	}

	// STEP 2: LLM Processing with Tool Execution

	// Enhance context with document metadata
	enhancedContext := make(map[string]interface{})
	for k, v := range cmd.Context {
		enhancedContext[k] = v
	}
	enhancedContext["input_type"] = "document"
	enhancedContext["document_format"] = cmd.DocumentFormat
	enhancedContext["analysis_type"] = analysisType
	enhancedContext["processing_confidence"] = processingConfidence
	enhancedContext["word_count"] = wordCount
	enhancedContext["page_count"] = pageCount
	enhancedContext["is_structured_data"] = isStructuredData

	// Build enhanced prompt combining user request with extracted content
	enhancedPrompt := h.buildEnhancedPrompt(cmd.UserPrompt, extractedContent, analysisResult, cmd.DocumentFormat, analysisType)

	response, actions, llmConfidence, err := h.llmProcessor.ProcessWithHistory(
		ctx,
		manager,
		enhancedPrompt,
		[]domain.ConversationMessage{}, // No history for simple document processing
		enhancedContext,
	)
	if err != nil {
		return nil, fmt.Errorf("LLM processing failed: %w", err)
	}

	// STEP 3: Calculate Combined Confidence
	combinedConfidence := (processingConfidence + llmConfidence) / 2.0

	// Adjust for low processing confidence
	if processingConfidence < 0.7 {
		combinedConfidence *= 0.8 // Penalize low processing confidence
	}

	// Save manager state
	if err := h.managers.Save(ctx, manager); err != nil {
		// Non-critical error, continue
	}

	// Return comprehensive result
	result := &ProcessDocumentInputResult{
		ResponseID:           fmt.Sprintf("doc_%d", time.Now().UnixNano()),
		ExtractedContent:     extractedContent,
		AnalysisResult:       analysisResult,
		ResponseMessage:      response,
		ResponseStatus:       "success",
		ResponseConfidence:   combinedConfidence,
		ProcessingConfidence: processingConfidence,
		LLMConfidence:        llmConfidence,
		ResponseTimestamp:    time.Now(),
		DocumentMetadata:     documentMetadata,
		ExecutedActions:      actions,
		DocumentFormat:       cmd.DocumentFormat,
		AnalysisType:         analysisType,
		ProcessingTime:       time.Since(startTime),
		InputSource:          getDocumentInputSource(cmd),
		WordCount:            wordCount,
		PageCount:            pageCount,
	}

	return result, nil
}

// isStructuredDataFormat determines if the format is structured data vs document
func (h *ProcessDocumentInputHandler) isStructuredDataFormat(format string) bool {
	structuredFormats := []string{"csv", "json", "xlsx", "xls", "tsv", "xml", "yaml", "sql"}
	format = strings.ToLower(format)

	for _, structuredFormat := range structuredFormats {
		if format == structuredFormat {
			return true
		}
	}
	return false
}

// buildEnhancedPrompt creates an enhanced prompt combining user request with extracted content
func (h *ProcessDocumentInputHandler) buildEnhancedPrompt(userPrompt, extractedContent, analysisResult, documentFormat, analysisType string) string {
	// Base prompt with extracted content
	basePrompt := fmt.Sprintf("I have processed a %s document using %s analysis. Here's what I found:\n\n%s",
		strings.ToUpper(documentFormat), analysisType, analysisResult)

	// Add extracted content if not too long
	if len(extractedContent) > 0 {
		contentToInclude := extractedContent
		if len(contentToInclude) > 3000 {
			contentToInclude = contentToInclude[:3000] + "... (content truncated)"
		}
		basePrompt += fmt.Sprintf("\n\nExtracted content:\n%s", contentToInclude)
	}

	// Add user prompt if provided
	if userPrompt != "" {
		return fmt.Sprintf("%s\n\nUser request: %s", basePrompt, userPrompt)
	}

	return basePrompt
}

// getDocumentInputSource determines whether binary data or URL was used
func getDocumentInputSource(cmd ProcessDocumentInput) string {
	if cmd.DocumentURL != "" {
		return "url"
	}
	return "binary_data"
}
