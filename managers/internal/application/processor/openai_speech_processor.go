package processor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"middleman/internal/ai"
	"middleman/managers/internal/application/services"
)

const (
	logPrefixSpeechProcessor = "[OpenAISpeechProcessor]"

	// Supported audio formats
	AudioFormatWAV  = "wav"
	AudioFormatMP3  = "mp3"
	AudioFormatM4A  = "m4a"
	AudioFormatFLAC = "flac"
	AudioFormatWEBM = "webm"
	AudioFormatMP4  = "mp4"
	AudioFormatOGG  = "ogg"

	// OpenAI Whisper models
	WhisperModel = "whisper-1"

	// Maximum audio file size (25MB as per OpenAI limits)
	MaxAudioSizeBytes = 25 * 1024 * 1024

	// Default response format
	DefaultResponseFormat = "verbose_json"
)

// OpenAISpeechProcessor implements the SpeechProcessor interface using OpenAI's Whisper API
type OpenAISpeechProcessor struct {
	aiClient ai.EnhancedAIService
	config   *SpeechProcessorConfig
}

// SpeechProcessorConfig holds configuration for speech processing
type SpeechProcessorConfig struct {
	Model            string        `json:"model"`             // Whisper model to use
	DefaultLanguage  string        `json:"default_language"`  // Default language code
	Temperature      float64       `json:"temperature"`       // Temperature for transcription
	ResponseFormat   string        `json:"response_format"`   // Response format
	EnableTimestamps bool          `json:"enable_timestamps"` // Whether to include timestamps
	EnableWordLevel  bool          `json:"enable_word_level"` // Whether to include word-level timestamps
	MaxRetries       int           `json:"max_retries"`       // Maximum retry attempts
	RetryDelay       time.Duration `json:"retry_delay"`       // Delay between retries
}

// NewOpenAISpeechProcessor creates a new OpenAI speech processor
func NewOpenAISpeechProcessor(aiClient ai.EnhancedAIService, config *SpeechProcessorConfig) services.SpeechProcessor {
	if aiClient == nil {
		log.Printf("%s ERROR: AI client cannot be nil, returning nil processor", logPrefixSpeechProcessor)
		return nil
	}

	if config == nil {
		config = &SpeechProcessorConfig{
			Model:            WhisperModel,
			DefaultLanguage:  "en",
			Temperature:      0.0,
			ResponseFormat:   DefaultResponseFormat,
			EnableTimestamps: true,
			EnableWordLevel:  true,
			MaxRetries:       3,
			RetryDelay:       1 * time.Second,
		}
	}

	return &OpenAISpeechProcessor{
		aiClient: aiClient,
		config:   config,
	}
}

// TranscribeAudio converts audio data to text using OpenAI's Whisper API
func (p *OpenAISpeechProcessor) TranscribeAudio(ctx context.Context, request services.SpeechTranscriptionRequest) (*services.SpeechTranscriptionResult, error) {
	startTime := time.Now()
	log.Printf("%s Starting audio transcription. Format: %s, Size: %d bytes",
		logPrefixSpeechProcessor, request.AudioFormat, len(request.AudioData))

	// Validate input
	if err := p.validateTranscriptionRequest(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Build OpenAI transcription request
	openaiRequest := p.buildOpenAIRequest(request)

	// Execute transcription with retries
	response, err := p.transcribeWithRetries(ctx, openaiRequest)
	if err != nil {
		log.Printf("%s Transcription failed: %v", logPrefixSpeechProcessor, err)
		return nil, fmt.Errorf("transcription failed: %w", err)
	}

	// Convert OpenAI response to our format
	result := p.convertToSpeechResult(response, startTime)

	log.Printf("%s Transcription completed successfully. Text: '%s', Confidence: %.2f, Duration: %v",
		logPrefixSpeechProcessor, result.Text, result.Confidence, result.Duration)

	return result, nil
}

// ValidateAudioFormat checks if the audio format is supported
func (p *OpenAISpeechProcessor) ValidateAudioFormat(format string) error {
	supportedFormats := p.GetSupportedFormats()

	normalizedFormat := strings.ToLower(strings.TrimSpace(format))
	for _, supported := range supportedFormats {
		if normalizedFormat == supported {
			return nil
		}
	}

	return fmt.Errorf("unsupported audio format: %s. Supported formats: %v", format, supportedFormats)
}

// GetSupportedLanguages returns list of supported language codes
func (p *OpenAISpeechProcessor) GetSupportedLanguages() []string {
	// OpenAI Whisper supports 99+ languages
	// Returning the most commonly used ones for performance
	return []string{
		"en", "es", "fr", "de", "it", "pt", "ru", "ja", "ko", "zh",
		"ar", "hi", "th", "vi", "nl", "sv", "da", "no", "fi", "pl",
		"cs", "sk", "hu", "ro", "bg", "hr", "sl", "et", "lv", "lt",
		"mt", "ga", "cy", "eu", "ca", "gl", "ast", "oc", "br", "co",
		"fo", "is", "gd", "gv", "kw", "lb", "rm", "sc", "vec", "wa",
		"af", "am", "as", "az", "ba", "be", "bn", "bo", "bs", "ce",
		"ck", "cv", "dv", "dz", "ee", "eo", "fa", "ff", "fj", "fy",
		"ha", "haw", "he", "hr", "ht", "hy", "ia", "id", "ig", "ik",
		"jw", "ka", "kk", "km", "kn", "ky", "la", "lb", "lg", "li",
		"ln", "lo", "lt", "lu", "lv", "mg", "mi", "mk", "ml", "mn",
		"mr", "ms", "my", "ne", "nn", "oc", "or", "pa", "ps", "qu",
		"rn", "sa", "sd", "si", "sk", "sm", "sn", "so", "sq", "sr",
		"su", "sw", "ta", "te", "tg", "tk", "tl", "tn", "to", "tr",
		"ts", "tt", "tw", "ty", "ug", "uk", "ur", "uz", "ve", "vo",
		"war", "wo", "xh", "yi", "yo", "zh-TW", "zu",
	}
}

// GetSupportedFormats returns list of supported audio formats
func (p *OpenAISpeechProcessor) GetSupportedFormats() []string {
	return []string{
		AudioFormatWAV,
		AudioFormatMP3,
		AudioFormatM4A,
		AudioFormatFLAC,
		AudioFormatWEBM,
		AudioFormatMP4,
		AudioFormatOGG,
	}
}

// validateTranscriptionRequest validates the transcription request
func (p *OpenAISpeechProcessor) validateTranscriptionRequest(request services.SpeechTranscriptionRequest) error {
	// Check audio data
	if len(request.AudioData) == 0 {
		return fmt.Errorf("audio data cannot be empty")
	}

	if len(request.AudioData) > MaxAudioSizeBytes {
		return fmt.Errorf("audio file too large: %d bytes (max %d bytes)",
			len(request.AudioData), MaxAudioSizeBytes)
	}

	// Validate audio format
	if request.AudioFormat != "" {
		if err := p.ValidateAudioFormat(request.AudioFormat); err != nil {
			return err
		}
	}

	// Validate language if provided
	if request.Language != "" {
		if err := p.validateLanguage(request.Language); err != nil {
			return err
		}
	}

	return nil
}

// validateLanguage checks if the language is supported
func (p *OpenAISpeechProcessor) validateLanguage(language string) error {
	supportedLanguages := p.GetSupportedLanguages()

	normalizedLang := strings.ToLower(strings.TrimSpace(language))
	for _, supported := range supportedLanguages {
		if normalizedLang == supported {
			return nil
		}
	}

	// Also check common language format variations
	if len(normalizedLang) > 2 {
		langCode := normalizedLang[:2]
		for _, supported := range supportedLanguages {
			if langCode == supported {
				return nil // Accept "en-US" for "en", etc.
			}
		}
	}

	return fmt.Errorf("unsupported language: %s", language)
}

// buildOpenAIRequest builds the OpenAI transcription request
func (p *OpenAISpeechProcessor) buildOpenAIRequest(request services.SpeechTranscriptionRequest) ai.AudioTranscriptionRequest {
	language := request.Language
	if language == "" {
		language = p.config.DefaultLanguage
	}

	format := request.AudioFormat
	if format == "" {
		format = AudioFormatWAV // Default
	}

	responseFormat := p.config.ResponseFormat
	if p.config.EnableTimestamps {
		responseFormat = "verbose_json"
	}

	return ai.AudioTranscriptionRequest{
		AudioData:      request.AudioData,
		Language:       language,
		ResponseFormat: responseFormat,
		Temperature:    p.config.Temperature,
		Model:          p.config.Model,
	}
}

// transcribeWithRetries executes transcription with retry logic
func (p *OpenAISpeechProcessor) transcribeWithRetries(ctx context.Context, request ai.AudioTranscriptionRequest) (*ai.AudioTranscriptionResponse, error) {
	var lastErr error

	for attempt := 1; attempt <= p.config.MaxRetries; attempt++ {
		// Type assert to access OpenAI client's TranscribeAudio method
		// We know this is an OpenAI client from the AIClientProvider architecture
		openAIClient, ok := p.aiClient.(*ai.OpenAIClient)
		if !ok {
			return nil, fmt.Errorf("speech processor requires OpenAI client, got %T", p.aiClient)
		}

		response, err := openAIClient.TranscribeAudio(ctx, request)
		if err == nil {
			return response, nil
		}

		lastErr = err

		// Don't retry on validation errors or context cancellation
		if strings.Contains(err.Error(), "validation") || ctx.Err() != nil {
			break
		}

		if attempt < p.config.MaxRetries {
			log.Printf("%s Transcription attempt %d failed, retrying in %v: %v",
				logPrefixSpeechProcessor, attempt, p.config.RetryDelay, err)

			select {
			case <-time.After(p.config.RetryDelay):
				continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	return nil, fmt.Errorf("transcription failed after %d attempts: %w", p.config.MaxRetries, lastErr)
}

// convertToSpeechResult converts OpenAI response to our SpeechTranscriptionResult format
func (p *OpenAISpeechProcessor) convertToSpeechResult(response *ai.AudioTranscriptionResponse, startTime time.Time) *services.SpeechTranscriptionResult {
	result := &services.SpeechTranscriptionResult{
		Text:       response.Text,
		Confidence: p.calculateConfidence(response),
		Language:   response.Language,
		Duration:   time.Duration(response.Duration * float64(time.Second)),
		Metadata: map[string]interface{}{
			"segments":        p.convertSegments(response.Segments),
			"processing_time": time.Since(startTime).Seconds(),
		},
	}

	return result
}

// calculateConfidence calculates confidence score from OpenAI response
func (p *OpenAISpeechProcessor) calculateConfidence(response *ai.AudioTranscriptionResponse) float64 {
	if len(response.Segments) == 0 {
		return 0.8 // Default confidence when no segments available
	}

	// Calculate average confidence from segments
	totalConfidence := 0.0
	validSegments := 0

	for _, segment := range response.Segments {
		// Use average log probability as confidence indicator
		if segment.AvgLogprob > -1.0 {
			// Convert log probability to confidence (0-1)
			confidence := 1.0 + (segment.AvgLogprob / 3.0) // Normalize roughly
			if confidence < 0 {
				confidence = 0
			}
			if confidence > 1 {
				confidence = 1
			}

			totalConfidence += confidence
			validSegments++
		}
	}

	if validSegments == 0 {
		return 0.8 // Default confidence
	}

	avgConfidence := totalConfidence / float64(validSegments)

	// Adjust confidence based on other factors
	if response.Text == "" {
		return 0.0
	}

	// Lower confidence for very short transcriptions
	if len(strings.TrimSpace(response.Text)) < 5 {
		avgConfidence *= 0.7
	}

	return avgConfidence
}

// convertSegments converts OpenAI segments to our TextSegment format
func (p *OpenAISpeechProcessor) convertSegments(openaiSegments []ai.TranscriptionSegment) []services.TextSegment {
	if !p.config.EnableTimestamps {
		return nil
	}

	segments := make([]services.TextSegment, 0, len(openaiSegments))

	for _, segment := range openaiSegments {
		textSegment := services.TextSegment{
			Text:      segment.Text,
			StartTime: segment.Start,
			EndTime:   segment.End,
		}

		segments = append(segments, textSegment)
	}

	return segments
}

// GetConfig returns the current speech processor configuration
func (p *OpenAISpeechProcessor) GetConfig() *SpeechProcessorConfig {
	return p.config
}

// UpdateConfig updates the speech processor configuration
func (p *OpenAISpeechProcessor) UpdateConfig(config *SpeechProcessorConfig) {
	if config != nil {
		p.config = config
	}
}
