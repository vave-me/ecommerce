package processor

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"middleman/managers/internal/application/services"
	"middleman/managers/internal/constants"
	"middleman/internal/ai"
)

const (
	logPrefixDocumentProcessor = "[OpenAIDocumentProcessor]"

	// Document format constants
	DocumentFormatPDF  = "pdf"
	DocumentFormatDOCX = "docx"
	DocumentFormatTXT  = "txt"
	DocumentFormatRTF  = "rtf"
	DocumentFormatHTML = "html"
	DocumentFormatMD   = "markdown"
	DocumentFormatODT  = "odt"

	// Analysis type constants for documents (document-specific)
	AnalysisTypeExtractText     = "extract_text"
	AnalysisTypeExtractTables   = "extract_tables"
	AnalysisTypeExtractForms    = "extract_forms"
	AnalysisTypeExtractEntities = "extract_entities"
	AnalysisTypeAnalyze         = "analyze"

	// Size limits
	MaxDocumentSizeBytes = 50 * 1024 * 1024 // 50MB limit for documents

	// OpenAI model for document processing
	DefaultDocumentModel = "gpt-4.1-mini" // Fast, cost-effective GPT-4.1 mini
)

// OpenAIDocumentProcessor implements the DocumentProcessor interface using OpenAI's API
type OpenAIDocumentProcessor struct {
	aiClient ai.EnhancedAIService // Use the existing ai service
	config   DocumentProcessorConfig
}

// DocumentProcessorConfig holds configuration for the document processor
type DocumentProcessorConfig struct {
	Model            string
	MaxTokens        int
	Temperature      float32
	MaxRetries       int
	RequestTimeout   time.Duration
	EnableStructured bool // Enable structured data extraction
}

// NewOpenAIDocumentProcessor creates a new OpenAI-based document processor
func NewOpenAIDocumentProcessor(aiClient ai.EnhancedAIService, config *DocumentProcessorConfig) services.DocumentProcessor {
	if aiClient == nil {
		log.Printf("%s ERROR: AI client cannot be nil, returning nil processor", logPrefixDocumentProcessor)
		return nil
	}

	if config == nil {
		config = &DocumentProcessorConfig{
			Model:            DefaultDocumentModel,
			MaxTokens:        4000,
			Temperature:      0.1, // Low temperature for factual document analysis
			MaxRetries:       3,
			RequestTimeout:   60 * time.Second,
			EnableStructured: true,
		}
	}

	return &OpenAIDocumentProcessor{
		aiClient: aiClient,
		config:   *config,
	}
}

// ProcessDocument analyzes and extracts information from documents
func (p *OpenAIDocumentProcessor) ProcessDocument(ctx context.Context, request services.DocumentProcessingRequest) (*services.DocumentProcessingResult, error) {
	startTime := time.Now()
	log.Printf("%s Starting document processing. Format: %s, Size: %d bytes, Type: %s",
		logPrefixDocumentProcessor, request.DocumentFormat, len(request.DocumentData), request.AnalysisType)

	// Validate input
	if err := p.validateDocumentRequest(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Build processing request using existing ai.CompletionRequest
	completionRequest := p.buildDocumentRequest(request)

	// Execute processing with retries
	response, err := p.processWithRetries(ctx, completionRequest)
	if err != nil {
		log.Printf("%s Document processing failed: %v", logPrefixDocumentProcessor, err)
		return nil, fmt.Errorf("document processing failed: %w", err)
	}

	// Convert response to document result
	result := p.convertToDocumentResult(response, request, startTime)

	log.Printf("%s Document processing completed successfully. Type: %s, Pages: %d, Duration: %v",
		logPrefixDocumentProcessor, result.DocumentType, result.PageCount, time.Since(startTime))

	return result, nil
}

// ValidateDocumentFormat checks if the document format is supported
func (p *OpenAIDocumentProcessor) ValidateDocumentFormat(format string) error {
	format = strings.ToLower(format)
	supportedFormats := p.GetSupportedDocumentTypes()

	for _, supported := range supportedFormats {
		if format == supported {
			return nil
		}
	}

	return fmt.Errorf("unsupported document format: %s. Supported formats: %v", format, supportedFormats)
}

// GetSupportedDocumentTypes returns list of supported document types
func (p *OpenAIDocumentProcessor) GetSupportedDocumentTypes() []string {
	return []string{
		DocumentFormatPDF,
		DocumentFormatDOCX,
		DocumentFormatTXT,
		DocumentFormatRTF,
		DocumentFormatHTML,
		DocumentFormatMD,
		DocumentFormatODT,
	}
}

// GetSupportedAnalysisTypes returns list of supported analysis types for documents
func (p *OpenAIDocumentProcessor) GetSupportedAnalysisTypes() []string {
	return []string{
		AnalysisTypeExtractText,
		constants.AnalysisTypeAnalyzeStructure,
		constants.AnalysisTypeSummarize,
		AnalysisTypeExtractTables,
		AnalysisTypeExtractForms,
		AnalysisTypeExtractEntities,
		constants.AnalysisTypeClassify,
		AnalysisTypeAnalyze,
	}
}

// ExtractStructuredData extracts structured data from documents (tables, forms, etc.)
func (p *OpenAIDocumentProcessor) ExtractStructuredData(ctx context.Context, request services.DocumentProcessingRequest) (*services.StructuredDataResult, error) {
	startTime := time.Now()
	log.Printf("%s Starting structured data extraction. Format: %s", logPrefixDocumentProcessor, request.DocumentFormat)

	// Override analysis type for structured extraction
	structuredRequest := request
	structuredRequest.AnalysisType = "extract_structured"

	// Build specialized structured extraction request
	completionRequest := p.buildStructuredExtractionRequest(structuredRequest)

	// Execute processing
	response, err := p.processWithRetries(ctx, completionRequest)
	if err != nil {
		return nil, fmt.Errorf("structured extraction failed: %w", err)
	}

	// Convert to structured result
	result := p.convertToStructuredResult(response, structuredRequest, startTime)

	log.Printf("%s Structured data extraction completed. Tables: %d, Key-Value Pairs: %d, Duration: %v",
		logPrefixDocumentProcessor, len(result.Tables), len(result.KeyValuePairs), time.Since(startTime))

	return result, nil
}

// validateDocumentRequest validates the document processing request
func (p *OpenAIDocumentProcessor) validateDocumentRequest(request services.DocumentProcessingRequest) error {
	// Must have either document data or URL, but not both
	if len(request.DocumentData) == 0 && request.DocumentURL == "" {
		return fmt.Errorf("either document data or document URL must be provided")
	}

	if len(request.DocumentData) > 0 && request.DocumentURL != "" {
		return fmt.Errorf("provide either document data or document URL, not both")
	}

	// Validate binary data if provided
	if len(request.DocumentData) > 0 {
		if len(request.DocumentData) > MaxDocumentSizeBytes {
			return fmt.Errorf("document size %d exceeds maximum allowed size %d", len(request.DocumentData), MaxDocumentSizeBytes)
		}

		if err := p.ValidateDocumentFormat(request.DocumentFormat); err != nil {
			return err
		}
	}

	// Validate URL if provided
	if request.DocumentURL != "" {
		if !strings.HasPrefix(request.DocumentURL, "http://") && !strings.HasPrefix(request.DocumentURL, "https://") {
			return fmt.Errorf("document URL must be a valid HTTP or HTTPS URL")
		}
	}

	return nil
}

// buildDocumentRequest creates an OpenAI completion request for document processing
func (p *OpenAIDocumentProcessor) buildDocumentRequest(request services.DocumentProcessingRequest) ai.CompletionRequest {
	// Build the prompt based on analysis type and user request
	prompt := p.buildDocumentPrompt(request.AnalysisType, request.UserPrompt, request.DocumentFormat)

	// Prepare the completion request using existing ai.CompletionRequest
	maxTokens := p.config.MaxTokens
	temperature := float64(p.config.Temperature)

	completionRequest := ai.CompletionRequest{
		Model:       p.config.Model,
		MaxTokens:   &maxTokens,
		Temperature: &temperature,
		Messages: []ai.Message{
			{
				Role:    "system",
				Content: p.buildSystemPrompt(request.AnalysisType),
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Add document data to the request
	if len(request.DocumentData) > 0 {
		// For binary documents, we need to convert to base64 for the API
		base64Data := base64.StdEncoding.EncodeToString(request.DocumentData)
		documentContent := fmt.Sprintf("Document data (base64): %s", base64Data)

		// Add document content to user message
		completionRequest.Messages[1].Content = fmt.Sprintf("%s\n\nDocument to analyze:\n%s", prompt, documentContent)
	} else if request.DocumentURL != "" {
		// For URLs, include the URL in the prompt
		completionRequest.Messages[1].Content = fmt.Sprintf("%s\n\nDocument URL: %s", prompt, request.DocumentURL)
	}

	return completionRequest
}

// buildStructuredExtractionRequest creates a specialized request for structured data extraction
func (p *OpenAIDocumentProcessor) buildStructuredExtractionRequest(request services.DocumentProcessingRequest) ai.CompletionRequest {
	prompt := `Extract all structured data from this document including:
- Tables with headers and data
- Forms with field names and values  
- Lists (ordered and unordered)
- Key-value pairs
- Any other structured information

Return the results in a structured JSON format with detailed extraction confidence scores.`

	maxTokens := p.config.MaxTokens
	temperature := 0.0 // Very low temperature for precise extraction

	completionRequest := ai.CompletionRequest{
		Model:       p.config.Model,
		MaxTokens:   &maxTokens,
		Temperature: &temperature,
		Messages: []ai.Message{
			{
				Role:    "system",
				Content: "You are an expert document analysis system specializing in structured data extraction. Extract all structured elements with high precision and provide confidence scores.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Add document data
	if len(request.DocumentData) > 0 {
		base64Data := base64.StdEncoding.EncodeToString(request.DocumentData)
		documentContent := fmt.Sprintf("Document data (base64): %s", base64Data)
		completionRequest.Messages[1].Content = fmt.Sprintf("%s\n\nDocument to analyze:\n%s", prompt, documentContent)
	}

	return completionRequest
}

// buildSystemPrompt creates a system prompt based on analysis type
func (p *OpenAIDocumentProcessor) buildSystemPrompt(analysisType string) string {
	basePrompt := "You are an expert document analysis system. Analyze documents accurately and provide detailed, structured responses in JSON format."

	specificPrompts := map[string]string{
		AnalysisTypeExtractText:     "Focus on extracting clean, well-formatted text while preserving document structure.",
		AnalysisTypeExtractTables:   "Identify and extract all tables with proper headers, data, and relationships.",
		AnalysisTypeExtractForms:    "Locate and extract form fields, labels, values, and form structure.",
		AnalysisTypeExtractEntities: "Identify named entities (people, organizations, locations, dates, etc.) with confidence scores.",
		AnalysisTypeAnalyze:         "Perform comprehensive analysis including content, structure, and metadata.",
	}

	if specific, exists := specificPrompts[analysisType]; exists {
		return fmt.Sprintf("%s %s", basePrompt, specific)
	}

	return basePrompt
}

// buildDocumentPrompt creates a prompt based on analysis type and user input
func (p *OpenAIDocumentProcessor) buildDocumentPrompt(analysisType, userPrompt, documentFormat string) string {
	basePrompts := map[string]string{
		AnalysisTypeExtractText:     "Extract all text content from this document. Preserve formatting, structure, and hierarchy. Include headings, paragraphs, lists, and any other textual elements.",
		AnalysisTypeExtractTables:   "Identify and extract all tables from this document. For each table, provide headers, all data rows, captions, and table structure information.",
		AnalysisTypeExtractForms:    "Locate and extract all form elements including field names, values, input types, and form structure. Include any form instructions or labels.",
		AnalysisTypeExtractEntities: "Extract all named entities from this document including people, organizations, locations, dates, phone numbers, emails, and other important entities.",
		AnalysisTypeAnalyze:         "Perform a comprehensive analysis of this document including content analysis, structural analysis, and metadata extraction.",
	}

	basePrompt := basePrompts[analysisType]
	if basePrompt == "" {
		basePrompt = "Analyze this document and provide detailed information about its content and structure."
	}

	// Add document format context
	formatContext := fmt.Sprintf("Document format: %s", strings.ToUpper(documentFormat))

	// Combine with user prompt if provided
	if userPrompt != "" {
		return fmt.Sprintf("%s\n\n%s\n\nSpecific request: %s", formatContext, basePrompt, userPrompt)
	}

	return fmt.Sprintf("%s\n\n%s", formatContext, basePrompt)
}

// processWithRetries performs document processing with retry logic using existing ai client
func (p *OpenAIDocumentProcessor) processWithRetries(ctx context.Context, request ai.CompletionRequest) (*ai.CompletionResponse, error) {
	var lastErr error

	for attempt := 1; attempt <= p.config.MaxRetries; attempt++ {
		log.Printf("%s Processing attempt %d/%d", logPrefixDocumentProcessor, attempt, p.config.MaxRetries)

		response, err := p.aiClient.CreateCompletion(ctx, request)
		if err == nil {
			return response, nil
		}

		lastErr = err
		log.Printf("%s Attempt %d failed: %v", logPrefixDocumentProcessor, attempt, err)

		if attempt < p.config.MaxRetries {
			backoff := time.Duration(attempt) * time.Second
			log.Printf("%s Retrying in %v...", logPrefixDocumentProcessor, backoff)
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("all retry attempts failed. Last error: %w", lastErr)
}

// convertToDocumentResult converts OpenAI response to DocumentProcessingResult
func (p *OpenAIDocumentProcessor) convertToDocumentResult(response *ai.CompletionResponse, request services.DocumentProcessingRequest, startTime time.Time) *services.DocumentProcessingResult {
	// Extract main content
	content := ""
	confidence := 0.0

	if len(response.Choices) > 0 {
		content = response.Choices[0].Message.GetContentAsString()
		confidence = 0.9 // Default high confidence for successful response
	}

	// Basic metadata
	metadata := map[string]interface{}{
		"model":           p.config.Model,
		"analysis_type":   request.AnalysisType,
		"processing_time": time.Since(startTime).Seconds(),
		"document_format": request.DocumentFormat,
	}

	if response.Usage.TotalTokens > 0 {
		metadata["tokens_used"] = response.Usage.TotalTokens
		metadata["prompt_tokens"] = response.Usage.PromptTokens
		metadata["completion_tokens"] = response.Usage.CompletionTokens
	}

	// Estimate document properties
	wordCount := len(strings.Fields(content))

	result := &services.DocumentProcessingResult{
		ExtractedText: content,
		Summary:       p.extractSummary(content, request.AnalysisType),
		DocumentType:  p.detectDocumentType(request.DocumentFormat, content),
		PageCount:     p.estimatePageCount(content, request.DocumentFormat),
		Metadata: map[string]interface{}{
			"confidence":       confidence,
			"language":         "en", // Could be enhanced with language detection
			"word_count":       wordCount,
			"extracted_tables": p.extractTables(content),
			"sections":         p.extractSections(content),
			"keywords":         p.extractKeywords(content),
			"entities":         p.extractEntities(content),
			"model":            p.config.Model,
			"processing_time":  time.Since(startTime).Seconds(),
			"document_format":  request.DocumentFormat,
		},
	}

	return result
}

// convertToStructuredResult converts response to StructuredDataResult
func (p *OpenAIDocumentProcessor) convertToStructuredResult(response *ai.CompletionResponse, request services.DocumentProcessingRequest, startTime time.Time) *services.StructuredDataResult {
	confidence := 0.9
	if len(response.Choices) > 0 {
		confidence = 0.95 // High confidence for structured extraction
	}

	result := &services.StructuredDataResult{
		Tables:        []services.ExtractedTable{}, // Would be populated from JSON response
		Lists:         []services.ExtractedList{},  // Would be populated from JSON response
		KeyValuePairs: map[string]interface{}{},    // Would be populated from JSON response
		Metadata: map[string]interface{}{
			"confidence":      confidence,
			"structure_type":  "mixed",
			"processing_time": time.Since(startTime).Seconds(),
			"model":           p.config.Model,
			"document_format": request.DocumentFormat,
		},
	}

	return result
}

// Helper methods for content extraction and analysis

func (p *OpenAIDocumentProcessor) extractSummary(content, analysisType string) string {
	if analysisType == constants.AnalysisTypeSummarize {
		return content // For summary analysis, the content is the summary
	}

	// For other analysis types, create a brief summary
	words := strings.Fields(content)
	if len(words) > 50 {
		return strings.Join(words[:50], " ") + "..."
	}
	return content
}

func (p *OpenAIDocumentProcessor) detectDocumentType(format, content string) string {
	// Basic document type detection based on format and content patterns
	switch strings.ToLower(format) {
	case "pdf":
		return "PDF Document"
	case "docx":
		return "Word Document"
	case "txt":
		return "Text Document"
	case "html":
		return "HTML Document"
	case "md", "markdown":
		return "Markdown Document"
	default:
		return "Unknown Document"
	}
}

func (p *OpenAIDocumentProcessor) estimatePageCount(content, format string) int {
	// Basic page count estimation
	wordCount := len(strings.Fields(content))

	switch strings.ToLower(format) {
	case "pdf", "docx":
		// Assume ~250 words per page for formatted documents
		return (wordCount / 250) + 1
	default:
		// For text documents, estimate based on characters
		return (len(content) / 2000) + 1
	}
}

func (p *OpenAIDocumentProcessor) extractTables(content string) []services.ExtractedTable {
	// This would be enhanced to actually parse JSON response for table data
	return []services.ExtractedTable{}
}

func (p *OpenAIDocumentProcessor) extractSections(content string) []services.DocumentSection {
	// This would be enhanced to actually parse document sections
	return []services.DocumentSection{}
}

func (p *OpenAIDocumentProcessor) extractKeywords(content string) []string {
	// Basic keyword extraction - would be enhanced with NLP
	words := strings.Fields(content)
	keywords := []string{}

	for _, word := range words {
		if len(word) > 4 && strings.ToUpper(word) == word {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

func (p *OpenAIDocumentProcessor) extractEntities(content string) []services.NamedEntity {
	// This would be enhanced with actual named entity recognition
	return []services.NamedEntity{}
}
