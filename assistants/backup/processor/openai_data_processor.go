package processor

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"middleman/assistants/internal/application/services"
	"middleman/assistants/internal/constants"
	"middleman/internal/ai"
)

const (
	logPrefixDataProcessor = "[OpenAIDataProcessor]"

	// Supported data formats
	DataFormatCSV  = "csv"
	DataFormatJSON = "json"
	DataFormatXLSX = "xlsx"
	DataFormatXLS  = "xls"
	DataFormatTSV  = "tsv"
	DataFormatXML  = "xml"
	DataFormatYAML = "yaml"
	DataFormatSQL  = "sql"

	// Analysis type constants for data (data-specific)
	AnalysisTypeValidateQuality = "validate_quality"
	AnalysisTypeProfile         = "profile"
	AnalysisTypeSchema          = "detect_schema"
	AnalysisTypeAnomalies       = "detect_anomalies"
	AnalysisTypeRelationships   = "find_relationships"

	// Size limits for data files
	MaxDataSizeBytes = 100 * 1024 * 1024 // 100MB limit for data files

	// Default model for data processing
	DefaultDataModel = "gpt-4.1-mini" // Fast, cost-effective GPT-4.1 mini
)

// OpenAIDataProcessor implements the DataProcessor interface using OpenAI's API
type OpenAIDataProcessor struct {
	aiClient ai.EnhancedAIService
	config   DataProcessorConfig
	analyzer *DataAnalyzer
}

// DataProcessorConfig holds configuration for the data processor
type DataProcessorConfig struct {
	Model           string
	MaxTokens       int
	Temperature     float32
	MaxRetries      int
	RequestTimeout  time.Duration
	SampleRows      int  // Number of sample rows to analyze
	EnableProfiling bool // Enable statistical profiling
}

// NewOpenAIDataProcessor creates a new OpenAI-based data processor
func NewOpenAIDataProcessor(aiClient ai.EnhancedAIService, config *DataProcessorConfig) services.DataProcessor {
	if aiClient == nil {
		log.Printf("%s ERROR: AI client cannot be nil, returning nil processor", logPrefixDataProcessor)
		return nil
	}

	if config == nil {
		config = &DataProcessorConfig{
			Model:           DefaultDataModel,
			MaxTokens:       4000,
			Temperature:     0.0, // Very low temperature for data analysis
			MaxRetries:      3,
			RequestTimeout:  60 * time.Second,
			SampleRows:      100,
			EnableProfiling: true,
		}
	}

	return &OpenAIDataProcessor{
		aiClient: aiClient,
		config:   *config,
		analyzer: NewDataAnalyzer(),
	}
}

// ProcessDataFile analyzes structured data files
func (p *OpenAIDataProcessor) ProcessDataFile(ctx context.Context, request services.DataProcessingRequest) (*services.DataProcessingResult, error) {
	startTime := time.Now()
	log.Printf("%s Starting data file processing. Format: %s, Size: %d bytes, Type: %s",
		logPrefixDataProcessor, request.DataFormat, len(request.DataContent), request.AnalysisType)

	// Add overall timeout for the entire operation
	processCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Validate input
	if err := p.validateDataRequest(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get data content
	dataContent, err := p.getDataContent(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get data content: %w", err)
	}

	// Use real data analyzer for supported formats
	var result *services.DataProcessingResult
	format := strings.ToLower(request.DataFormat)
	
	switch format {
	case DataFormatCSV, "text/csv":
		result, err = p.analyzer.AnalyzeCSV(dataContent)
	case DataFormatJSON, "application/json":
		result, err = p.analyzer.AnalyzeJSON(dataContent)
	case DataFormatXML, "application/xml", "text/xml":
		result, err = p.analyzer.AnalyzeXML(dataContent)
	case DataFormatYAML, "application/yaml", "text/yaml", "yml":
		result, err = p.analyzer.AnalyzeYAML(dataContent)
	default:
		// For unsupported formats or large files, use AI
		if len(dataContent) > 10000 {
			log.Printf("%s Large data file detected (%d bytes), using simplified processing", 
				logPrefixDataProcessor, len(dataContent))
			return p.processLargeDataFile(processCtx, request, startTime)
		}
		// Use AI for analysis
		result, err = p.processWithAI(processCtx, request, startTime)
	}

	if err != nil {
		return nil, fmt.Errorf("data analysis failed: %w", err)
	}

	// Add metadata
	if result.Metadata == nil {
		result.Metadata = make(map[string]interface{})
	}
	result.Metadata["processing_time"] = time.Since(startTime).Seconds()
	result.Metadata["analysis_type"] = request.AnalysisType
	result.Metadata["data_format"] = request.DataFormat

	// Handle specific analysis types
	if request.AnalysisType == AnalysisTypeValidateQuality {
		qualityResult, qErr := p.analyzer.ValidateDataQuality(dataContent, format)
		if qErr != nil {
			log.Printf("%s Warning: Quality validation failed: %v", logPrefixDataProcessor, qErr)
		} else {
			result.Metadata["quality_analysis"] = qualityResult
		}
	}

	// If user provided specific prompt, enhance with AI
	if request.UserPrompt != "" && len(dataContent) < 50000 { // Only for reasonable sizes
		aiEnhancement, aiErr := p.enhanceWithAI(processCtx, result, request)
		if aiErr != nil {
			log.Printf("%s Warning: AI enhancement failed: %v", logPrefixDataProcessor, aiErr)
		} else {
			result.Metadata["ai_insights"] = aiEnhancement
		}
	}

	log.Printf("%s Data processing completed successfully. Rows: %d, Columns: %d, Duration: %v",
		logPrefixDataProcessor, result.RowCount, result.ColumnCount, time.Since(startTime))

	return result, nil
}

// ValidateDataFormat checks if the data format is supported
func (p *OpenAIDataProcessor) ValidateDataFormat(format string) error {
	format = strings.ToLower(format)
	supportedFormats := p.GetSupportedDataFormats()

	for _, supported := range supportedFormats {
		if format == supported {
			return nil
		}
	}

	return fmt.Errorf("unsupported data format: %s. Supported formats: %v", format, supportedFormats)
}

// GetSupportedDataFormats returns list of supported data formats
func (p *OpenAIDataProcessor) GetSupportedDataFormats() []string {
	return []string{
		DataFormatCSV,
		DataFormatJSON,
		DataFormatXLSX,
		DataFormatXLS,
		DataFormatTSV,
		DataFormatXML,
		DataFormatYAML,
		DataFormatSQL,
	}
}

// AnalyzeDataStructure analyzes the structure of the data
func (p *OpenAIDataProcessor) AnalyzeDataStructure(ctx context.Context, request services.DataProcessingRequest) (*services.DataStructureResult, error) {
	startTime := time.Now()
	log.Printf("%s Starting data structure analysis. Format: %s", logPrefixDataProcessor, request.DataFormat)

	// Override analysis type for structure analysis
	structureRequest := request
	structureRequest.AnalysisType = constants.AnalysisTypeAnalyzeStructure

	// Build specialized structure analysis request
	completionRequest := p.buildStructureAnalysisRequest(structureRequest)

	// Execute processing
	response, err := p.processWithRetries(ctx, completionRequest)
	if err != nil {
		return nil, fmt.Errorf("structure analysis failed: %w", err)
	}

	// Convert to structure result
	result := p.convertToStructureResult(response, structureRequest, startTime)

	log.Printf("%s Data structure analysis completed. Schema detected: %v, Duration: %v",
		logPrefixDataProcessor, len(result.Schema.Columns) > 0, time.Since(startTime))

	return result, nil
}

// ValidateDataQuality checks data quality and integrity
func (p *OpenAIDataProcessor) ValidateDataQuality(ctx context.Context, request services.DataProcessingRequest) (*services.DataQualityResult, error) {
	startTime := time.Now()
	log.Printf("%s Starting data quality validation. Format: %s", logPrefixDataProcessor, request.DataFormat)

	// Get data content first
	data, err := p.getDataContent(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get data content: %w", err)
	}

	// Override analysis type for quality validation
	qualityRequest := request
	qualityRequest.AnalysisType = AnalysisTypeValidateQuality

	// Build specialized quality validation request
	completionRequest := p.buildQualityValidationRequest(qualityRequest)

	// Execute processing
	response, err := p.processWithRetries(ctx, completionRequest)
	if err != nil {
		return nil, fmt.Errorf("quality validation failed: %w", err)
	}

	// Convert to quality result
	result := p.convertToQualityResult(response, qualityRequest, data, startTime)

	log.Printf("%s Data quality validation completed. Quality score: %.2f, Issues: %d, Duration: %v",
		logPrefixDataProcessor, result.QualityScore, len(result.Issues), time.Since(startTime))

	return result, nil
}

// validateDataRequest validates the data processing request
func (p *OpenAIDataProcessor) validateDataRequest(request services.DataProcessingRequest) error {
	// Must have either data content or URL, but not both
	if len(request.DataContent) == 0 && request.DataURL == "" {
		return fmt.Errorf("either data content or data URL must be provided")
	}

	if len(request.DataContent) > 0 && request.DataURL != "" {
		return fmt.Errorf("provide either data content or data URL, not both")
	}

	// Validate binary data if provided
	if len(request.DataContent) > 0 {
		if len(request.DataContent) > MaxDataSizeBytes {
			return fmt.Errorf("data size %d exceeds maximum allowed size %d", len(request.DataContent), MaxDataSizeBytes)
		}

		if err := p.ValidateDataFormat(request.DataFormat); err != nil {
			return err
		}
	}

	// Validate URL if provided
	if request.DataURL != "" {
		if !strings.HasPrefix(request.DataURL, "http://") && !strings.HasPrefix(request.DataURL, "https://") {
			return fmt.Errorf("data URL must be a valid HTTP or HTTPS URL")
		}
	}

	return nil
}

// buildDataRequest creates an OpenAI completion request for data processing
func (p *OpenAIDataProcessor) buildDataRequest(request services.DataProcessingRequest) ai.CompletionRequest {
	// Build the prompt based on analysis type and user request
	prompt := p.buildDataPrompt(request.AnalysisType, request.UserPrompt, request.DataFormat)

	maxTokens := p.config.MaxTokens
	temperature := float64(p.config.Temperature)

	// Prepare the completion request
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

	// Add data content to the request
	if len(request.DataContent) > 0 {
		// For text-based data (CSV, JSON, etc.), include directly or as base64 for binary
		if p.isTextFormat(request.DataFormat) {
			dataContent := string(request.DataContent)
			// Truncate if too long for context
			if len(dataContent) > 50000 {
				dataContent = dataContent[:50000] + "\n... (truncated)"
			}
			completionRequest.Messages[1].Content = fmt.Sprintf("%s\n\nData to analyze:\n```%s\n%s\n```", prompt, request.DataFormat, dataContent)
		} else {
			// For binary formats like Excel, use base64
			base64Data := base64.StdEncoding.EncodeToString(request.DataContent)
			completionRequest.Messages[1].Content = fmt.Sprintf("%s\n\nData to analyze (base64):\n%s", prompt, base64Data)
		}
	} else if request.DataURL != "" {
		// For URLs, include the URL in the prompt
		completionRequest.Messages[1].Content = fmt.Sprintf("%s\n\nData URL: %s", prompt, request.DataURL)
	}

	return completionRequest
}

// buildStructureAnalysisRequest creates a specialized request for structure analysis
func (p *OpenAIDataProcessor) buildStructureAnalysisRequest(request services.DataProcessingRequest) ai.CompletionRequest {
	prompt := `Analyze the structure of this data file and provide:
- Detected schema with column names and data types
- Data relationships and patterns
- Anomalies and outliers
- Data quality assessment
- Suggested improvements

Return results in structured JSON format with confidence scores.`

	maxTokens := p.config.MaxTokens
	temperature := 0.0

	completionRequest := ai.CompletionRequest{
		Model:       p.config.Model,
		MaxTokens:   &maxTokens,
		Temperature: &temperature,
		Messages: []ai.Message{
			{
				Role:    "system",
				Content: "You are an expert data analyst specializing in data structure analysis. Analyze data files comprehensively and provide detailed insights about schema, patterns, and quality.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Add data content (similar to buildDataRequest)
	if len(request.DataContent) > 0 {
		if p.isTextFormat(request.DataFormat) {
			dataContent := string(request.DataContent)
			if len(dataContent) > 50000 {
				dataContent = dataContent[:50000] + "\n... (truncated)"
			}
			completionRequest.Messages[1].Content = fmt.Sprintf("%s\n\nData to analyze:\n```%s\n%s\n```", prompt, request.DataFormat, dataContent)
		} else {
			base64Data := base64.StdEncoding.EncodeToString(request.DataContent)
			completionRequest.Messages[1].Content = fmt.Sprintf("%s\n\nData to analyze (base64):\n%s", prompt, base64Data)
		}
	}

	return completionRequest
}

// buildQualityValidationRequest creates a specialized request for quality validation
func (p *OpenAIDataProcessor) buildQualityValidationRequest(request services.DataProcessingRequest) ai.CompletionRequest {
	prompt := `Validate the quality of this data file and provide:
- Overall quality score (0-1)
- Identified quality issues with severity levels
- Data completeness and consistency analysis
- Recommendations for improvement
- Automatically fixable issues

Return results in structured JSON format with detailed assessments.`

	maxTokens := p.config.MaxTokens
	temperature := 0.0

	completionRequest := ai.CompletionRequest{
		Model:       p.config.Model,
		MaxTokens:   &maxTokens,
		Temperature: &temperature,
		Messages: []ai.Message{
			{
				Role:    "system",
				Content: "You are an expert data quality analyst. Assess data files for quality issues, completeness, consistency, and provide actionable recommendations.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Add data content
	if len(request.DataContent) > 0 {
		if p.isTextFormat(request.DataFormat) {
			dataContent := string(request.DataContent)
			if len(dataContent) > 50000 {
				dataContent = dataContent[:50000] + "\n... (truncated)"
			}
			completionRequest.Messages[1].Content = fmt.Sprintf("%s\n\nData to analyze:\n```%s\n%s\n```", prompt, request.DataFormat, dataContent)
		} else {
			base64Data := base64.StdEncoding.EncodeToString(request.DataContent)
			completionRequest.Messages[1].Content = fmt.Sprintf("%s\n\nData to analyze (base64):\n%s", prompt, base64Data)
		}
	}

	return completionRequest
}

// buildSystemPrompt creates a system prompt based on analysis type
func (p *OpenAIDataProcessor) buildSystemPrompt(analysisType string) string {
	basePrompt := "You are an expert data analyst and scientist. Analyze data files accurately and provide detailed, structured responses in JSON format."

	specificPrompts := map[string]string{
		constants.AnalysisTypeAnalyzeStructure: "Focus on data schema detection, column types, relationships, and structural patterns.",
		AnalysisTypeValidateQuality:            "Assess data quality including completeness, consistency, accuracy, and validity.",
		constants.AnalysisTypeSummarize:        "Provide comprehensive data summaries with key statistics and insights.",
		AnalysisTypeProfile:                    "Generate detailed data profiles with statistical analysis and distributions.",
		AnalysisTypeSchema:                     "Detect and define comprehensive data schemas with constraints and relationships.",
		AnalysisTypeAnomalies:                  "Identify data anomalies, outliers, and unusual patterns with confidence scores.",
		AnalysisTypeRelationships:              "Find and analyze relationships between data elements and columns.",
		constants.AnalysisTypeClassify:         "Classify data types, categories, and content patterns.",
	}

	if specific, exists := specificPrompts[analysisType]; exists {
		return fmt.Sprintf("%s %s", basePrompt, specific)
	}

	return basePrompt
}

// buildDataPrompt creates a prompt based on analysis type and user input
func (p *OpenAIDataProcessor) buildDataPrompt(analysisType, userPrompt, dataFormat string) string {
	basePrompts := map[string]string{
		constants.AnalysisTypeAnalyzeStructure: "Analyze the structure of this data file. Identify columns, data types, relationships, patterns, and any structural anomalies.",
		AnalysisTypeValidateQuality:  "Validate the quality of this data file. Check for completeness, consistency, accuracy, duplicates, and provide quality metrics.",
		constants.AnalysisTypeSummarize:        "Provide a comprehensive summary of this data file including row/column counts, data types, key statistics, and notable insights.",
		AnalysisTypeProfile:          "Generate a detailed data profile with statistical analysis, distributions, patterns, and data characteristics.",
		AnalysisTypeSchema:           "Detect and define the schema of this data file including column definitions, constraints, and relationships.",
		AnalysisTypeAnomalies:        "Identify anomalies, outliers, and unusual patterns in this data file with confidence scores and explanations.",
		AnalysisTypeRelationships:    "Find and analyze relationships between columns and data elements in this data file.",
		constants.AnalysisTypeClassify:         "Classify the data types, content categories, and patterns found in this data file.",
	}

	basePrompt := basePrompts[analysisType]
	if basePrompt == "" {
		basePrompt = "Analyze this data file and provide detailed insights about its content, structure, and quality."
	}

	// Add data format context
	formatContext := fmt.Sprintf("Data format: %s", strings.ToUpper(dataFormat))

	// Combine with user prompt if provided
	if userPrompt != "" {
		return fmt.Sprintf("%s\n\n%s\n\nSpecific request: %s", formatContext, basePrompt, userPrompt)
	}

	return fmt.Sprintf("%s\n\n%s", formatContext, basePrompt)
}

// processWithRetries performs data processing with retry logic
func (p *OpenAIDataProcessor) processWithRetries(ctx context.Context, request ai.CompletionRequest) (*ai.CompletionResponse, error) {
	var lastErr error

	for attempt := 1; attempt <= p.config.MaxRetries; attempt++ {
		log.Printf("%s Processing attempt %d/%d", logPrefixDataProcessor, attempt, p.config.MaxRetries)

		// Create a timeout context for each attempt
		attemptCtx, cancel := context.WithTimeout(ctx, p.config.RequestTimeout)
		
		response, err := p.aiClient.CreateCompletion(attemptCtx, request)
		cancel() // Cancel immediately after use to prevent goroutine leak
		
		if err == nil {
			return response, nil
		}

		lastErr = err
		log.Printf("%s Attempt %d failed: %v", logPrefixDataProcessor, attempt, err)

		// Check if context was cancelled
		if ctx.Err() != nil {
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		}

		if attempt < p.config.MaxRetries {
			backoff := time.Duration(attempt) * time.Second
			log.Printf("%s Retrying in %v...", logPrefixDataProcessor, backoff)
			
			select {
			case <-time.After(backoff):
				// Continue to next attempt
			case <-ctx.Done():
				return nil, fmt.Errorf("context cancelled during backoff: %w", ctx.Err())
			}
		}
	}

	return nil, fmt.Errorf("all retry attempts failed. Last error: %w", lastErr)
}

// isTextFormat checks if the data format is text-based
func (p *OpenAIDataProcessor) isTextFormat(format string) bool {
	textFormats := []string{DataFormatCSV, DataFormatJSON, DataFormatTSV, DataFormatXML, DataFormatYAML, DataFormatSQL}
	format = strings.ToLower(format)

	for _, textFormat := range textFormats {
		if format == textFormat {
			return true
		}
	}
	return false
}

// convertToDataResult converts OpenAI response to DataProcessingResult
func (p *OpenAIDataProcessor) convertToDataResult(response *ai.CompletionResponse, request services.DataProcessingRequest, startTime time.Time) *services.DataProcessingResult {
	// Extract main content
	content := ""

	if len(response.Choices) > 0 {
		content = response.Choices[0].Message.GetContentAsString()
	}

	// Basic metadata
	metadata := map[string]interface{}{
		"model":           p.config.Model,
		"analysis_type":   request.AnalysisType,
		"processing_time": time.Since(startTime).Seconds(),
		"data_format":     request.DataFormat,
	}

	if response.Usage.TotalTokens > 0 {
		metadata["tokens_used"] = response.Usage.TotalTokens
		metadata["prompt_tokens"] = response.Usage.PromptTokens
		metadata["completion_tokens"] = response.Usage.CompletionTokens
	}

	// Estimate data properties (would be enhanced with actual parsing)
	result := &services.DataProcessingResult{
		Summary:         p.extractSummary(content, request.AnalysisType),
		RowCount:        p.estimateRowCount(request),
		ColumnCount:     p.estimateColumnCount(request),
		SampleData:      []map[string]interface{}{},    // Would be populated with samples
		Metadata:        metadata,
	}

	return result
}

// convertToStructureResult converts response to DataStructureResult
func (p *OpenAIDataProcessor) convertToStructureResult(response *ai.CompletionResponse, request services.DataProcessingRequest, startTime time.Time) *services.DataStructureResult {
	// Basic metadata
	metadata := map[string]interface{}{
		"model":           p.config.Model,
		"processing_time": time.Since(startTime).Seconds(),
		"data_format":     request.DataFormat,
	}

	result := &services.DataStructureResult{
		Schema:         services.DataSchema{},         // Would be populated from JSON response
		Relationships:  []services.DataRelationship{}, // Would be populated from analysis
		Patterns:       []services.DataPattern{},      // Would be populated from analysis
		Anomalies:      []services.DataAnomaly{},      // Would be populated from analysis
		Completeness:   0.9,                           // Would be calculated from analysis
		Consistency:    0.9,                           // Would be calculated from analysis
		Validity:       0.9,                           // Would be calculated from analysis
		Suggestions:    []string{},                    // Would be populated from analysis
		ProcessingTime: time.Since(startTime),
		Metadata:       metadata,
	}

	return result
}

// convertToQualityResult converts response to DataQualityResult
func (p *OpenAIDataProcessor) convertToQualityResult(response *ai.CompletionResponse, qualityRequest services.DataProcessingRequest, data []byte, startTime time.Time) *services.DataQualityResult {
	overallScore := 0.8 // Would be extracted from analysis

	if len(response.Choices) > 0 {
		// Try to extract score from response
		// This is where we'd parse the AI response for quality metrics
	}

	// Basic metadata
	metadata := map[string]interface{}{
		"model":           p.config.Model,
		"processing_time": time.Since(startTime).Seconds(),
		"data_format":     qualityRequest.DataFormat,
	}

	// Use real analyzer for quality validation
	qualityResult, err := p.analyzer.ValidateDataQuality(data, qualityRequest.DataFormat)
	if err != nil {
		// Fallback to basic quality result
		qualityResult = &services.DataQualityResult{
			QualityScore:    overallScore,
			Issues:          []services.DataQualityIssue{},
			MissingValues:   map[string]int{},
			DuplicateRows:   0,
			DataTypes:       map[string]string{},
			ValidityChecks:  []services.ValidityCheck{},
			Recommendations: []string{"Unable to perform full quality analysis"},
			Metadata:        metadata,
		}
	}
	
	// Add metadata
	qualityResult.Metadata = metadata

	return qualityResult
}

// Helper methods for estimation and extraction

func (p *OpenAIDataProcessor) extractSummary(content, analysisType string) string {
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

func (p *OpenAIDataProcessor) estimateRowCount(request services.DataProcessingRequest) int {
	// Basic row count estimation based on data format
	if len(request.DataContent) == 0 {
		return 0
	}

	switch strings.ToLower(request.DataFormat) {
	case DataFormatCSV, DataFormatTSV:
		// Count newlines in CSV/TSV
		content := string(request.DataContent)
		return strings.Count(content, "\n")
	case DataFormatJSON:
		// Estimate based on JSON structure
		content := string(request.DataContent)
		return strings.Count(content, "{") // Rough estimate
	default:
		return 1000 // Default estimate
	}
}

func (p *OpenAIDataProcessor) estimateColumnCount(request services.DataProcessingRequest) int {
	// Basic column count estimation
	if len(request.DataContent) == 0 {
		return 0
	}

	switch strings.ToLower(request.DataFormat) {
	case DataFormatCSV, DataFormatTSV:
		// Count delimiters in first line
		content := string(request.DataContent)
		lines := strings.Split(content, "\n")
		if len(lines) > 0 {
			delimiter := ","
			if request.DataFormat == DataFormatTSV {
				delimiter = "\t"
			}
			return strings.Count(lines[0], delimiter) + 1
		}
	default:
		return 10 // Default estimate
	}

	return 0
}

// processLargeDataFile handles large data files with a simplified approach
func (p *OpenAIDataProcessor) processLargeDataFile(ctx context.Context, request services.DataProcessingRequest, startTime time.Time) (*services.DataProcessingResult, error) {
	log.Printf("%s Processing large data file with simplified approach", logPrefixDataProcessor)
	
	// Extract basic information without sending full content to AI
	dataContent := string(request.DataContent)
	
	// For JSON, try to parse structure
	if strings.ToLower(request.DataFormat) == DataFormatJSON {
		return p.processJSONDataLocally(ctx, request, dataContent, startTime)
	}
	
	// For other formats, use sampling
	sample := p.extractDataSample(dataContent, request.DataFormat)
	
	// Create a modified request with just the sample
	sampleRequest := request
	sampleRequest.DataContent = []byte(sample)
	
	// Build and process with sample
	completionRequest := p.buildDataRequest(sampleRequest)
	
	response, err := p.processWithRetries(ctx, completionRequest)
	if err != nil {
		log.Printf("%s Sample processing failed: %v", logPrefixDataProcessor, err)
		// Fall back to basic analysis
		return p.createBasicDataResult(request, dataContent, startTime), nil
	}
	
	// Convert response and adjust for full data
	result := p.convertToDataResult(response, request, startTime)
	
	// Update counts based on full data
	result.RowCount = p.estimateRowCount(request)
	result.ColumnCount = p.estimateColumnCount(request)
	
	return result, nil
}

// processJSONDataLocally processes JSON data without AI for large files
func (p *OpenAIDataProcessor) processJSONDataLocally(ctx context.Context, request services.DataProcessingRequest, dataContent string, startTime time.Time) (*services.DataProcessingResult, error) {
	log.Printf("%s Processing JSON data locally", logPrefixDataProcessor)
	
	// Basic JSON analysis
	var jsonData interface{}
	if err := json.Unmarshal([]byte(dataContent), &jsonData); err != nil {
		return nil, fmt.Errorf("invalid JSON data: %w", err)
	}
	
	// Analyze structure
	var rowCount, columnCount int
	columns := make(map[string]bool)
	
	switch data := jsonData.(type) {
	case []interface{}:
		rowCount = len(data)
		// Analyze first few items for columns
		for i, item := range data {
			if i > 10 { // Sample first 10 items
				break
			}
			if obj, ok := item.(map[string]interface{}); ok {
				for key := range obj {
					columns[key] = true
				}
			}
		}
		columnCount = len(columns)
	case map[string]interface{}:
		rowCount = 1
		columnCount = len(data)
		for key := range data {
			columns[key] = true
		}
	}
	
	// Create result
	result := &services.DataProcessingResult{
		Summary:     fmt.Sprintf("JSON data with %d records and %d fields", rowCount, columnCount),
		RowCount:    rowCount,
		ColumnCount: columnCount,
		SampleData:  []map[string]interface{}{},
		Metadata: map[string]interface{}{
			"processing_type": "local",
			"data_format":     "json",
			"file_size":       len(dataContent),
			"processing_time": time.Since(startTime).Seconds(),
			"columns":         p.createColumnInfo(columns),
		},
	}
	
	// Add user prompt analysis if needed
	if request.UserPrompt != "" {
		// Use AI for just the user's question with data summary
		summaryRequest := ai.CompletionRequest{
			Model:       p.config.Model,
			MaxTokens:   intPtr(500),
			Temperature: float64Ptr(0.1),
			Messages: []ai.Message{
				{
					Role:    "system",
					Content: "You are a data analyst. Answer questions about data based on the provided summary.",
				},
				{
					Role:    "user",
					Content: fmt.Sprintf("Data summary: %s\n\nUser question: %s", result.Summary, request.UserPrompt),
				},
			},
		}
		
		if response, err := p.aiClient.CreateCompletion(ctx, summaryRequest); err == nil && len(response.Choices) > 0 {
			result.Summary = response.Choices[0].Message.GetContentAsString()
		}
	}
	
	return result, nil
}

// extractDataSample extracts a representative sample from data
func (p *OpenAIDataProcessor) extractDataSample(dataContent string, format string) string {
	const maxSampleSize = 5000 // 5KB sample
	
	if len(dataContent) <= maxSampleSize {
		return dataContent
	}
	
	// For line-based formats, take first N lines
	if format == DataFormatCSV || format == DataFormatTSV {
		lines := strings.Split(dataContent, "\n")
		sampleLines := 50
		if len(lines) < sampleLines {
			sampleLines = len(lines)
		}
		return strings.Join(lines[:sampleLines], "\n") + "\n... (truncated)"
	}
	
	// For other formats, take beginning and end
	return dataContent[:maxSampleSize/2] + "\n... (truncated) ...\n" + dataContent[len(dataContent)-maxSampleSize/2:]
}

// createBasicDataResult creates a basic result without AI processing
func (p *OpenAIDataProcessor) createBasicDataResult(request services.DataProcessingRequest, dataContent string, startTime time.Time) *services.DataProcessingResult {
	rowCount := p.estimateRowCount(request)
	columnCount := p.estimateColumnCount(request)
	
	return &services.DataProcessingResult{
		Summary:     fmt.Sprintf("Data file analysis: %s format with approximately %d rows and %d columns", request.DataFormat, rowCount, columnCount),
		RowCount:    rowCount,
		ColumnCount: columnCount,
		SampleData:  []map[string]interface{}{},
		Metadata: map[string]interface{}{
			"processing_type": "basic",
			"data_format":     request.DataFormat,
			"file_size":       len(dataContent),
			"processing_time": time.Since(startTime).Seconds(),
			"columns":         []services.ColumnInfo{},
			"recommendations": []string{"Consider using smaller data samples for detailed AI analysis"},
		},
	}
}

// createColumnInfo creates column info from column names
func (p *OpenAIDataProcessor) createColumnInfo(columns map[string]bool) []services.ColumnInfo {
	result := make([]services.ColumnInfo, 0, len(columns))
	for col := range columns {
		result = append(result, services.ColumnInfo{
			Name:         col,
			DataType:     "unknown",
			SampleValues: []interface{}{},
			NullCount:    0,
			UniqueCount:  0,
			Statistics:   map[string]interface{}{},
		})
	}
	return result
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

// processWithAI processes data using AI when analyzer doesn't support the format
func (p *OpenAIDataProcessor) processWithAI(ctx context.Context, request services.DataProcessingRequest, startTime time.Time) (*services.DataProcessingResult, error) {
	// Build processing request
	completionRequest := p.buildDataRequest(request)

	// Execute processing with retries
	response, err := p.processWithRetries(ctx, completionRequest)
	if err != nil {
		log.Printf("%s Data processing failed: %v", logPrefixDataProcessor, err)
		return nil, fmt.Errorf("data processing failed: %w", err)
	}

	// Convert response to data result
	result := p.convertToDataResult(response, request, startTime)
	return result, nil
}

// enhanceWithAI enhances the result with AI insights based on user prompt
func (p *OpenAIDataProcessor) enhanceWithAI(ctx context.Context, result *services.DataProcessingResult, request services.DataProcessingRequest) (map[string]interface{}, error) {
	// Build context from result
	contextData := map[string]interface{}{
		"summary":      result.Summary,
		"row_count":    result.RowCount,
		"column_count": result.ColumnCount,
		"metadata":     result.Metadata,
	}

	contextJSON, _ := json.MarshalIndent(contextData, "", "  ")

	// Prepare messages
	messages := []ai.Message{
		{
			Role:    "system",
			Content: "You are analyzing data based on user requirements. Provide insights and answers to their specific questions.",
		},
		{
			Role: "user",
			Content: fmt.Sprintf("Data Analysis Context:\n%s\n\nUser Question: %s\n\nProvide detailed analysis and insights.",
				string(contextJSON), request.UserPrompt),
		},
	}

	// Create completion request
	temp := float64(p.config.Temperature)
	maxTokens := p.config.MaxTokens
	completionReq := ai.CompletionRequest{
		Messages:    messages,
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	}

	// Get AI response
	response, err := p.aiClient.CreateCompletion(ctx, completionReq)
	if err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	content := response.Choices[0].Message.GetContentAsString()

	// Try to parse as JSON
	var aiResult map[string]interface{}
	if err := json.Unmarshal([]byte(content), &aiResult); err != nil {
		aiResult = map[string]interface{}{
			"analysis": content,
		}
	}

	return aiResult, nil
}

// getDataContent retrieves data content from request
func (p *OpenAIDataProcessor) getDataContent(request services.DataProcessingRequest) ([]byte, error) {
	if request.DataContent != nil {
		return request.DataContent, nil
	}

	if request.DataURL != "" {
		// TODO: Implement URL fetching if needed
		return nil, fmt.Errorf("URL fetching not implemented")
	}

	return nil, fmt.Errorf("no data content provided")
}
