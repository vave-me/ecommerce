package services

import (
	"context"
	"middleman/managers/internal/domain"
	"time"
)

// Core AI and Processing Interfaces

// AIClientProvider is already defined in ai_interfaces.go

// LLMProcessor handles LLM integration for managers
type LLMProcessor interface {
	ProcessWithHistory(ctx context.Context, manager *domain.Manager, currentMessage string, history []domain.ConversationMessage, contextData map[string]interface{}) (string, []domain.ManagerAction, float64, error)
	ShouldUseAI(manager interface{}) bool
}

// SpeechProcessor handles speech-to-text operations
type SpeechProcessor interface {
	TranscribeAudio(ctx context.Context, request SpeechTranscriptionRequest) (*SpeechTranscriptionResult, error)
	ValidateAudioFormat(format string) error
	GetSupportedLanguages() []string
}

// VisionProcessor handles image analysis operations
type VisionProcessor interface {
	AnalyzeImage(ctx context.Context, request VisionAnalysisRequest) (*VisionAnalysisResult, error)
	ValidateImageFormat(format string) error
	GetSupportedAnalysisTypes() []string
}

// DocumentProcessor handles document processing
type DocumentProcessor interface {
	ProcessDocument(ctx context.Context, request DocumentProcessingRequest) (*DocumentProcessingResult, error)
	ValidateDocumentFormat(format string) error
	GetSupportedDocumentTypes() []string
}

// DataProcessor handles structured data files
type DataProcessor interface {
	ProcessDataFile(ctx context.Context, request DataProcessingRequest) (*DataProcessingResult, error)
	ValidateDataFormat(format string) error
	GetSupportedDataFormats() []string
}

// Request/Response Types for Speech Processing

type SpeechTranscriptionRequest struct {
	AudioData   []byte                 `json:"audio_data"`
	AudioFormat string                 `json:"audio_format"`
	Language    string                 `json:"language,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

type SpeechTranscriptionResult struct {
	Text             string                 `json:"text"`
	Confidence       float64                `json:"confidence"`
	Language         string                 `json:"language"`
	Duration         time.Duration          `json:"duration"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// Request/Response Types for Vision Processing

type VisionAnalysisRequest struct {
	ImageData   []byte                 `json:"image_data,omitempty"`
	ImageURL    string                 `json:"image_url,omitempty"`
	ImageFormat string                 `json:"image_format"`
	AnalysisType string                `json:"analysis_type"`
	UserPrompt  string                 `json:"user_prompt"`
	Context     map[string]interface{} `json:"context,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
}

type VisionAnalysisResult struct {
	Description     string                 `json:"description"`
	Confidence      float64                `json:"confidence"`
	DetectedObjects []DetectedObject       `json:"detected_objects,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

type DetectedObject struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Request/Response Types for Document Processing

type DocumentProcessingRequest struct {
	DocumentData   []byte                 `json:"document_data,omitempty"`
	DocumentURL    string                 `json:"document_url,omitempty"`
	DocumentFormat string                 `json:"document_format"`
	AnalysisType   string                 `json:"analysis_type"`
	UserPrompt     string                 `json:"user_prompt"`
	Context        map[string]interface{} `json:"context,omitempty"`
	MaxTokens      int                    `json:"max_tokens,omitempty"`
}

type DocumentProcessingResult struct {
	ExtractedText   string                 `json:"extracted_text"`
	Summary         string                 `json:"summary,omitempty"`
	DocumentType    string                 `json:"document_type"`
	PageCount       int                    `json:"page_count"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// Request/Response Types for Data Processing

type DataProcessingRequest struct {
	DataContent  []byte                 `json:"data_content,omitempty"`
	DataURL      string                 `json:"data_url,omitempty"`
	DataFormat   string                 `json:"data_format"`
	AnalysisType string                 `json:"analysis_type"`
	UserPrompt   string                 `json:"user_prompt"`
	Context      map[string]interface{} `json:"context,omitempty"`
}

type DataProcessingResult struct {
	Summary      string                   `json:"summary"`
	RowCount     int                      `json:"row_count"`
	ColumnCount  int                      `json:"column_count"`
	SampleData   []map[string]interface{} `json:"sample_data,omitempty"`
	Metadata     map[string]interface{}   `json:"metadata,omitempty"`
}

// Additional types for data processing

type DataStructureResult struct {
	Schema         DataSchema         `json:"schema"`
	Relationships  []DataRelationship `json:"relationships"`
	Patterns       []DataPattern      `json:"patterns"`
	Anomalies      []DataAnomaly      `json:"anomalies"`
	Completeness   float64            `json:"completeness"`
	Consistency    float64            `json:"consistency"`
	Validity       float64            `json:"validity"`
	Suggestions    []string           `json:"suggestions"`
	ProcessingTime time.Duration      `json:"processing_time"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type DataQualityResult struct {
	QualityScore      float64                `json:"quality_score"`
	Issues            []DataQualityIssue     `json:"issues"`
	MissingValues     map[string]int         `json:"missing_values"`
	DuplicateRows     int                    `json:"duplicate_rows"`
	DataTypes         map[string]string      `json:"data_types"`
	ValidityChecks    []ValidityCheck        `json:"validity_checks"`
	Recommendations   []string               `json:"recommendations"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

type DataSchema struct {
	Columns       []ColumnInfo           `json:"columns"`
	RowCount      int                    `json:"row_count"`
	FileFormat    string                 `json:"file_format"`
	Encoding      string                 `json:"encoding"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

type ColumnInfo struct {
	Name         string                 `json:"name"`
	DataType     string                 `json:"data_type"`
	Description  string                 `json:"description,omitempty"`
	SampleValues []interface{}          `json:"sample_values,omitempty"`
	NullCount    int                    `json:"null_count"`
	UniqueCount  int                    `json:"unique_count"`
	Statistics   map[string]interface{} `json:"statistics,omitempty"`
}

type DataRelationship struct {
	SourceColumn string `json:"source_column"`
	TargetColumn string `json:"target_column"`
	Type         string `json:"type"`
	Strength     float64 `json:"strength"`
	Description  string `json:"description,omitempty"`
}

type DataPattern struct {
	Column      string  `json:"column"`
	Pattern     string  `json:"pattern"`
	Frequency   int     `json:"frequency"`
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description,omitempty"`
}

type DataAnomaly struct {
	Column      string                 `json:"column"`
	Row         int                    `json:"row,omitempty"`
	Type        string                 `json:"type"`
	Value       interface{}            `json:"value,omitempty"`
	Expected    interface{}            `json:"expected,omitempty"`
	Severity    string                 `json:"severity"`
	Description string                 `json:"description"`
}

type DataQualityIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Column      string `json:"column,omitempty"`
	Count       int    `json:"count"`
	Description string `json:"description"`
	Examples    []interface{} `json:"examples,omitempty"`
}

type ValidityCheck struct {
	Name        string `json:"name"`
	Passed      bool   `json:"passed"`
	Description string `json:"description"`
	Details     string `json:"details,omitempty"`
}

// Additional types for document processing

type StructuredDataResult struct {
	Tables          []ExtractedTable       `json:"tables,omitempty"`
	KeyValuePairs   map[string]interface{} `json:"key_value_pairs,omitempty"`
	Lists           []ExtractedList        `json:"lists,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

type ExtractedTable struct {
	Title    string                   `json:"title,omitempty"`
	Headers  []string                 `json:"headers"`
	Rows     [][]string               `json:"rows"`
	Metadata map[string]interface{}   `json:"metadata,omitempty"`
}

type ExtractedList struct {
	Title    string   `json:"title,omitempty"`
	Items    []string `json:"items"`
	Ordered  bool     `json:"ordered"`
}

type DocumentSection struct {
	Title    string                 `json:"title"`
	Content  string                 `json:"content"`
	Level    int                    `json:"level"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type NamedEntity struct {
	Text     string                 `json:"text"`
	Type     string                 `json:"type"`
	Category string                 `json:"category,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Additional types for speech processing

type TextSegment struct {
	Text      string  `json:"text"`
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
	Speaker   string  `json:"speaker,omitempty"`
}