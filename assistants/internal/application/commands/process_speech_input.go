package commands

import (
	"context"
	"fmt"
	"time"

	"middleman/assistants/internal/application/services"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"
)


// ProcessSpeechInput represents the command to process speech input for an assistant.
type ProcessSpeechInput struct {
	ID          string                 `json:"id"`                     // Unique ID for this specific request/interaction
	AssistantID string                 `json:"assistant_id"`           // ID of the assistant to use
	UserID      string                 `json:"user_id"`                // ID of the user making the request
	AudioData   []byte                 `json:"audio_data,omitempty"`   // Raw audio data (WAV, MP3, M4A, etc.)
	AudioFormat string                 `json:"audio_format,omitempty"` // Audio format (wav, mp3, m4a, flac, etc.)
	Language    string                 `json:"language,omitempty"`     // Language code (en-US, es-ES, etc.) - optional
	Context     map[string]interface{} `json:"context,omitempty"`      // Additional context
	Timestamp   time.Time              `json:"timestamp,omitempty"`    // Timestamp of the request
	RequestType string                 `json:"request_type,omitempty"` // Type of request (e.g., "voice_chat", "voice_command", "voice_query")
}

// ProcessSpeechInputResult holds the structured result of processing speech input.
type ProcessSpeechInputResult struct {
	ResponseID         string                   `json:"response_id"`
	TranscribedText    string                   `json:"transcribed_text"` // The transcribed text from speech
	ResponseMessage    string                   `json:"response_message"` // AI assistant response
	ResponseStatus     string                   `json:"response_status"`
	ResponseConfidence float64                  `json:"response_confidence"` // Combined confidence (STT + LLM)
	STTConfidence      float64                  `json:"stt_confidence"`      // Speech-to-text confidence
	LLMConfidence      float64                  `json:"llm_confidence"`      // LLM processing confidence
	ResponseTimestamp  time.Time                `json:"response_timestamp"`
	AudioDuration      time.Duration            `json:"audio_duration"`             // Duration of processed audio
	ExecutedActions    []domain.AssistantAction `json:"executed_actions,omitempty"` // Actions executed during processing
}

// ProcessSpeechInputHandler orchestrates speech-to-text and subsequent tool execution.
type ProcessSpeechInputHandler struct {
	assistants      domain.AssistantRepository
	publisher       ddd.EventPublisher[ddd.Event]
	llmProcessor    services.LLMProcessor
	speechProcessor services.SpeechProcessor // New interface for STT
}

// NewProcessSpeechInputHandler creates a new speech input handler.
func NewProcessSpeechInputHandler(
	assistants domain.AssistantRepository,
	speechProcessor services.SpeechProcessor,
	llmProcessor services.LLMProcessor,
	publisher ddd.EventPublisher[ddd.Event],
) ProcessSpeechInputHandler {
	if assistants == nil || speechProcessor == nil || llmProcessor == nil || publisher == nil {
		panic("Critical dependencies cannot be nil")
	}
	return ProcessSpeechInputHandler{
		assistants:      assistants,
		speechProcessor: speechProcessor,
		llmProcessor:    llmProcessor,
		publisher:       publisher,
	}
}

// ProcessSpeechInput handles the speech input workflow: STT → LLM → Tools → Response.
func (h ProcessSpeechInputHandler) ProcessSpeechInput(ctx context.Context, cmd ProcessSpeechInput) (*ProcessSpeechInputResult, error) {

	// Validate audio data
	if len(cmd.AudioData) == 0 {
		return nil, fmt.Errorf("audio data cannot be empty")
	}
	if cmd.AudioFormat == "" {
		cmd.AudioFormat = "wav" // Default format
	}

	// Load the assistant
	assistant, err := h.assistants.Load(ctx, cmd.AssistantID)
	if err != nil {
		return nil, fmt.Errorf("failed to load assistant: %w", err)
	}

	// STEP 1: Speech-to-Text Transcription
	transcriptionRequest := services.SpeechTranscriptionRequest{
		AudioData:   cmd.AudioData,
		AudioFormat: cmd.AudioFormat,
		Language:    cmd.Language,
		Context:     cmd.Context,
	}

	transcriptionResult, err := h.speechProcessor.TranscribeAudio(ctx, transcriptionRequest)
	if err != nil {
		return nil, fmt.Errorf("speech transcription failed: %w", err)
	}

	if transcriptionResult.Text == "" {
		return &ProcessSpeechInputResult{
			ResponseID:         cmd.ID,
			TranscribedText:    "",
			ResponseMessage:    "I couldn't understand the audio. Could you please try again?",
			ResponseStatus:     "transcription_failed",
			ResponseConfidence: 0.0,
			STTConfidence:      transcriptionResult.Confidence,
			LLMConfidence:      0.0,
			ResponseTimestamp:  time.Now(),
			AudioDuration:      transcriptionResult.Duration,
			ExecutedActions:    []domain.AssistantAction{},
		}, nil
	}

	// STEP 2: LLM Processing with Tool Execution

	// Enhance context with speech metadata
	enhancedContext := make(map[string]interface{})
	for k, v := range cmd.Context {
		enhancedContext[k] = v
	}
	enhancedContext["input_type"] = "speech"
	enhancedContext["audio_format"] = cmd.AudioFormat
	enhancedContext["audio_duration"] = transcriptionResult.Duration.Seconds()
	enhancedContext["stt_confidence"] = transcriptionResult.Confidence
	enhancedContext["language"] = transcriptionResult.Language

	response, actions, llmConfidence, err := h.llmProcessor.ProcessWithHistory(
		ctx,
		assistant,
		transcriptionResult.Text,
		[]domain.ConversationMessage{}, // No history for simple speech processing
		enhancedContext,
	)
	if err != nil {
		return nil, fmt.Errorf("LLM processing failed: %w", err)
	}

	// STEP 3: Calculate Combined Confidence
	combinedConfidence := (transcriptionResult.Confidence + llmConfidence) / 2.0

	// Adjust for low STT confidence
	if transcriptionResult.Confidence < 0.7 {
		combinedConfidence *= 0.8 // Penalize low STT confidence
	}

	// Save assistant state
	if err := h.assistants.Save(ctx, assistant); err != nil {
		// Non-critical error, continue
	}

	// Return comprehensive result
	result := &ProcessSpeechInputResult{
		ResponseID:         cmd.ID,
		TranscribedText:    transcriptionResult.Text,
		ResponseMessage:    response,
		ResponseStatus:     "completed",
		ResponseConfidence: combinedConfidence,
		STTConfidence:      transcriptionResult.Confidence,
		LLMConfidence:      llmConfidence,
		ResponseTimestamp:  time.Now(),
		AudioDuration:      transcriptionResult.Duration,
		ExecutedActions:    actions,
	}

	return result, nil
}
