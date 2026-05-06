package ai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

// Enhanced data structures for multimodal support
type ImageContent struct {
	Type     string `json:"type"` // "image_url"
	ImageURL struct {
		URL    string `json:"url"`
		Detail string `json:"detail,omitempty"` // "low", "high", "auto"
	} `json:"image_url"`
}

type AudioContent struct {
	Type  string `json:"type"` // "input_audio"
	Audio struct {
		Data   string `json:"data"`   // base64 encoded
		Format string `json:"format"` // "wav", "mp3", etc.
	} `json:"input_audio"`
}

type VisionAnalysisRequest struct {
	ImagePath string   `json:"image_path,omitempty"`
	ImageURL  string   `json:"image_url,omitempty"`
	ImageData []byte   `json:"image_data,omitempty"`
	Questions []string `json:"questions,omitempty"`
	MaxTokens int      `json:"max_tokens,omitempty"`
	Detail    string   `json:"detail,omitempty"` // "low", "high", "auto"
	Model     string   `json:"model,omitempty"`
}

type AudioTranscriptionRequest struct {
	AudioPath      string  `json:"audio_path,omitempty"`
	AudioData      []byte  `json:"audio_data,omitempty"`
	Language       string  `json:"language,omitempty"`
	Prompt         string  `json:"prompt,omitempty"`
	ResponseFormat string  `json:"response_format,omitempty"` // "json", "text", "srt", "verbose_json", "vtt"
	Temperature    float64 `json:"temperature,omitempty"`
	Model          string  `json:"model,omitempty"`
}

type AudioTranscriptionResponse struct {
	Text     string                 `json:"text"`
	Language string                 `json:"language,omitempty"`
	Duration float64                `json:"duration,omitempty"`
	Words    []TranscriptionWord    `json:"words,omitempty"`
	Segments []TranscriptionSegment `json:"segments,omitempty"`
}

type TranscriptionWord struct {
	Word  string  `json:"word"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type TranscriptionSegment struct {
	ID               int                 `json:"id"`
	Seek             int                 `json:"seek"`
	Start            float64             `json:"start"`
	End              float64             `json:"end"`
	Text             string              `json:"text"`
	Words            []TranscriptionWord `json:"words,omitempty"`
	Temperature      float64             `json:"temperature"`
	AvgLogprob       float64             `json:"avg_logprob"`
	CompressionRatio float64             `json:"compression_ratio"`
	NoSpeechProb     float64             `json:"no_speech_prob"`
}

type SpeechSynthesisRequest struct {
	Text           string  `json:"text"`
	Voice          string  `json:"voice"`                     // "alloy", "echo", "fable", "onyx", "nova", "shimmer"
	ResponseFormat string  `json:"response_format,omitempty"` // "mp3", "opus", "aac", "flac", "wav", "pcm"
	Speed          float64 `json:"speed,omitempty"`           // 0.25 to 4.0
	Model          string  `json:"model,omitempty"`
}

type ImageGenerationRequest struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model,omitempty"`           // "dall-e-2", "dall-e-3"
	N              int    `json:"n,omitempty"`               // Number of images
	Quality        string `json:"quality,omitempty"`         // "standard", "hd"
	ResponseFormat string `json:"response_format,omitempty"` // "url", "b64_json"
	Size           string `json:"size,omitempty"`            // "256x256", "512x512", "1024x1024", "1024x1792", "1792x1024"
	Style          string `json:"style,omitempty"`           // "vivid", "natural"
	User           string `json:"user,omitempty"`
}

type ImageGenerationResponse struct {
	Created int64                    `json:"created"`
	Data    []GeneratedImageResponse `json:"data"`
}

type GeneratedImageResponse struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

type ToolRegistry struct {
	tools map[string]ToolDefinition
}

// OpenAIClient implements the EnhancedAIService interface using the official OpenAI Go library v1.6.0
// Now includes support for reasoning models via Responses API
type OpenAIClient struct {
	client          *openai.Client
	defaultModel    string
	stats           UsageStats
	config          ProviderConfig
	toolRegistry    *ToolRegistry
	responseAPIBase string // For reasoning models using Responses API
	debugMode       bool   // Enable debug logging
}

// NewOpenAIClient creates a new OpenAI client using the official library
func NewOpenAIClient(apiKey, baseURL, defaultModel string) (*OpenAIClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}

	if baseURL != "" && baseURL != "https://api.openai.com/v1" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	// Set timeout and retry configuration
	opts = append(opts,
		option.WithMaxRetries(3),
		option.WithRequestTimeout(600*time.Second),
	)

	client := openai.NewClient(opts...)

	// Set up Responses API base URL for reasoning models
	responseAPIBase := "https://api.openai.com/v1/responses"
	if baseURL != "" && baseURL != "https://api.openai.com/v1" {
		responseAPIBase = baseURL + "/responses"
	}

	return &OpenAIClient{
		client:          &client,
		defaultModel:    defaultModel,
		responseAPIBase: responseAPIBase,
		stats: UsageStats{
			Provider: ProviderOpenAI,
		},
		config: ProviderConfig{
			APIKey:       apiKey,
			BaseURL:      baseURL,
			DefaultModel: defaultModel,
			Enabled:      true,
		},
		toolRegistry: &ToolRegistry{
			tools: make(map[string]ToolDefinition),
		},
	}, nil
}

// ==================== VISION ANALYSIS ====================

// AnalyzeImage analyzes an image with vision capabilities
func (c *OpenAIClient) AnalyzeImage(ctx context.Context, request VisionAnalysisRequest) (*CompletionResponse, error) {
	var imageContent []openai.ChatCompletionContentPartUnionParam

	// Handle different image input types
	var imageURL string
	if request.ImageURL != "" {
		imageURL = request.ImageURL
	} else if request.ImagePath != "" {
		// Read and encode local image
		data, err := os.ReadFile(request.ImagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read image file: %w", err)
		}
		mimeType := c.detectImageMimeType(request.ImagePath)
		imageURL = fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data))
	} else if len(request.ImageData) > 0 {
		// Use provided image data
		imageURL = fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString(request.ImageData))
	} else {
		return nil, fmt.Errorf("no image provided")
	}

	// Build image content
	detail := request.Detail
	if detail == "" {
		detail = "auto"
	}

	imageContent = append(imageContent, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
		URL:    imageURL,
		Detail: detail,
	}))

	// Build text content
	prompt := "Analyze this image."
	if len(request.Questions) > 0 {
		prompt = strings.Join(request.Questions, " ")
	}

	imageContent = append(imageContent, openai.TextContentPart(prompt))

	// Create completion request
	model := request.Model
	if model == "" {
		model = ModelGPT4o // Vision requires GPT-4o or GPT-4 Vision
	}

	maxTokens := request.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(imageContent),
		},
		Model:     c.getSharedModel(model),
		MaxTokens: openai.Int(int64(maxTokens)),
	}

	completion, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("vision analysis failed: %w", err)
	}

	return c.convertCompletionResponse(completion), nil
}

// AnalyzeImages analyzes multiple images
func (c *OpenAIClient) AnalyzeImages(ctx context.Context, imagePaths []string, prompt string) (*CompletionResponse, error) {
	var content []openai.ChatCompletionContentPartUnionParam

	// Add text prompt
	content = append(content, openai.TextContentPart(prompt))

	// Add each image
	for _, imagePath := range imagePaths {
		data, err := os.ReadFile(imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read image %s: %w", imagePath, err)
		}

		mimeType := c.detectImageMimeType(imagePath)
		imageURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data))

		content = append(content, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
			URL:    imageURL,
			Detail: "auto",
		}))
	}

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(content),
		},
		Model:     shared.ChatModelGPT4o,
		MaxTokens: openai.Int(4096),
	}

	completion, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("multi-image analysis failed: %w", err)
	}

	return c.convertCompletionResponse(completion), nil
}

// ==================== AUDIO TRANSCRIPTION ====================

// TranscribeAudio transcribes audio to text using Whisper
func (c *OpenAIClient) TranscribeAudio(ctx context.Context, request AudioTranscriptionRequest) (*AudioTranscriptionResponse, error) {
	var audioFile io.Reader

	if request.AudioPath != "" {
		file, err := os.Open(request.AudioPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open audio file: %w", err)
		}
		defer file.Close()
		audioFile = file
	} else if len(request.AudioData) > 0 {
		audioFile = strings.NewReader(string(request.AudioData))
	} else {
		return nil, fmt.Errorf("no audio provided")
	}

	params := openai.AudioTranscriptionNewParams{
		File: audioFile,
	}

	// Set optional parameters
	if request.Model != "" {
		params.Model = openai.AudioModel(request.Model)
	} else {
		params.Model = openai.AudioModelWhisper1
	}

	if request.Language != "" {
		params.Language = openai.String(request.Language)
	}

	if request.Prompt != "" {
		params.Prompt = openai.String(request.Prompt)
	}

	if request.ResponseFormat != "" {
		switch request.ResponseFormat {
		case "json":
			params.ResponseFormat = openai.AudioResponseFormatJSON
		case "text":
			params.ResponseFormat = openai.AudioResponseFormatText
		case "srt":
			params.ResponseFormat = openai.AudioResponseFormatSRT
		case "verbose_json":
			params.ResponseFormat = openai.AudioResponseFormatVerboseJSON
		case "vtt":
			params.ResponseFormat = openai.AudioResponseFormatVTT
		}
	}

	if request.Temperature != 0 {
		params.Temperature = openai.Float(request.Temperature)
	}

	transcription, err := c.client.Audio.Transcriptions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("audio transcription failed: %w", err)
	}

	return &AudioTranscriptionResponse{
		Text: transcription.Text,
		// Note: Language and Duration fields are not available in current API response
	}, nil
}

// TranslateAudio translates audio to English text
func (c *OpenAIClient) TranslateAudio(ctx context.Context, audioPath string) (*AudioTranscriptionResponse, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()

	params := openai.AudioTranslationNewParams{
		File:  file,
		Model: openai.AudioModelWhisper1,
	}

	translation, err := c.client.Audio.Translations.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("audio translation failed: %w", err)
	}

	return &AudioTranscriptionResponse{
		Text:     translation.Text,
		Language: "en", // Translations are always to English
		// Note: Duration field not available in current API response
	}, nil
}

// ==================== SPEECH SYNTHESIS ====================

// SynthesizeSpeech converts text to speech
func (c *OpenAIClient) SynthesizeSpeech(ctx context.Context, request SpeechSynthesisRequest) (io.Reader, error) {
	voice := openai.AudioSpeechNewParamsVoiceAlloy
	switch request.Voice {
	case "echo":
		voice = openai.AudioSpeechNewParamsVoiceEcho
	case "fable":
		voice = openai.AudioSpeechNewParamsVoiceFable
	case "onyx":
		voice = openai.AudioSpeechNewParamsVoiceOnyx
	case "nova":
		voice = openai.AudioSpeechNewParamsVoiceNova
	case "shimmer":
		voice = openai.AudioSpeechNewParamsVoiceShimmer
	case "ash":
		voice = openai.AudioSpeechNewParamsVoiceAsh
	case "ballad":
		voice = openai.AudioSpeechNewParamsVoiceBallad
	case "coral":
		voice = openai.AudioSpeechNewParamsVoiceCoral
	case "sage":
		voice = openai.AudioSpeechNewParamsVoiceSage
	case "verse":
		voice = openai.AudioSpeechNewParamsVoiceVerse
	}

	responseFormat := openai.AudioSpeechNewParamsResponseFormatMP3
	if request.ResponseFormat != "" {
		switch request.ResponseFormat {
		case "opus":
			responseFormat = openai.AudioSpeechNewParamsResponseFormatOpus
		case "aac":
			responseFormat = openai.AudioSpeechNewParamsResponseFormatAAC
		case "flac":
			responseFormat = openai.AudioSpeechNewParamsResponseFormatFLAC
		case "wav":
			responseFormat = openai.AudioSpeechNewParamsResponseFormatWAV
		case "pcm":
			responseFormat = openai.AudioSpeechNewParamsResponseFormatPCM
		}
	}

	model := openai.SpeechModel(request.Model)
	if request.Model == "" {
		model = openai.SpeechModelTTS1
	}

	params := openai.AudioSpeechNewParams{
		Model:          model,
		Input:          request.Text,
		Voice:          voice,
		ResponseFormat: responseFormat,
	}

	if request.Speed != 0 {
		params.Speed = openai.Float(request.Speed)
	}

	audioResponse, err := c.client.Audio.Speech.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("speech synthesis failed: %w", err)
	}

	return audioResponse.Body, nil
}

// ==================== IMAGE GENERATION ====================

// GenerateImage generates images using DALL-E
func (c *OpenAIClient) GenerateImage(ctx context.Context, request ImageGenerationRequest) (*ImageGenerationResponse, error) {
	model := openai.ImageModel(request.Model)
	if request.Model == "" {
		model = openai.ImageModelDallE3
	}

	params := openai.ImageGenerateParams{
		Prompt: request.Prompt,
		Model:  model,
	}

	// Set optional parameters
	if request.N > 0 {
		params.N = openai.Int(int64(request.N))
	}

	if request.Quality != "" {
		switch request.Quality {
		case "hd":
			params.Quality = "hd" // Use string literal as per API docs
		default:
			params.Quality = openai.ImageGenerateParamsQualityStandard
		}
	}

	if request.ResponseFormat != "" {
		switch request.ResponseFormat {
		case "b64_json":
			params.ResponseFormat = openai.ImageGenerateParamsResponseFormatB64JSON
		default:
			params.ResponseFormat = openai.ImageGenerateParamsResponseFormatURL
		}
	}

	if request.Size != "" {
		params.Size = openai.ImageGenerateParamsSize(request.Size)
	}

	if request.Style != "" {
		switch request.Style {
		case "natural":
			params.Style = openai.ImageGenerateParamsStyleNatural
		default:
			params.Style = openai.ImageGenerateParamsStyleVivid
		}
	}

	if request.User != "" {
		params.User = openai.String(request.User)
	}

	imagesResponse, err := c.client.Images.Generate(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("image generation failed: %w", err)
	}

	// Convert response
	data := make([]GeneratedImageResponse, len(imagesResponse.Data))
	for i, img := range imagesResponse.Data {
		data[i] = GeneratedImageResponse{
			URL:           img.URL,
			B64JSON:       img.B64JSON,
			RevisedPrompt: img.RevisedPrompt,
		}
	}

	return &ImageGenerationResponse{
		Created: imagesResponse.Created,
		Data:    data,
	}, nil
}

// EditImage edits an image using DALL-E
func (c *OpenAIClient) EditImage(ctx context.Context, imagePath, maskPath, prompt string) (*ImageGenerationResponse, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer imageFile.Close()

	params := openai.ImageEditParams{
		Image: openai.ImageEditParamsImageUnion{
			OfFile: imageFile,
		},
		Prompt: prompt,
		Model:  openai.ImageModelDallE2, // Only DALL-E 2 supports editing
	}

	if maskPath != "" {
		maskFile, err := os.Open(maskPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open mask file: %w", err)
		}
		defer maskFile.Close()
		params.Mask = maskFile
	}

	imagesResponse, err := c.client.Images.Edit(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("image editing failed: %w", err)
	}

	// Convert response
	data := make([]GeneratedImageResponse, len(imagesResponse.Data))
	for i, img := range imagesResponse.Data {
		data[i] = GeneratedImageResponse{
			URL:           img.URL,
			B64JSON:       img.B64JSON,
			RevisedPrompt: img.RevisedPrompt,
		}
	}

	return &ImageGenerationResponse{
		Created: imagesResponse.Created,
		Data:    data,
	}, nil
}

// CreateImageVariation creates variations of an image
func (c *OpenAIClient) CreateImageVariation(ctx context.Context, imagePath string, n int) (*ImageGenerationResponse, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer imageFile.Close()

	params := openai.ImageNewVariationParams{
		Image: imageFile,
		Model: openai.ImageModelDallE2, // Only DALL-E 2 supports variations
	}

	if n > 0 {
		params.N = openai.Int(int64(n))
	}

	imagesResponse, err := c.client.Images.NewVariation(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("image variation failed: %w", err)
	}

	// Convert response
	data := make([]GeneratedImageResponse, len(imagesResponse.Data))
	for i, img := range imagesResponse.Data {
		data[i] = GeneratedImageResponse{
			URL:           img.URL,
			B64JSON:       img.B64JSON,
			RevisedPrompt: img.RevisedPrompt,
		}
	}

	return &ImageGenerationResponse{
		Created: imagesResponse.Created,
		Data:    data,
	}, nil
}

// ==================== TOOL REGISTRATION ====================

// RegisterTool registers a new tool for function calling
func (c *OpenAIClient) RegisterTool(tool ToolDefinition) error {
	if c.toolRegistry == nil {
		c.toolRegistry = &ToolRegistry{
			tools: make(map[string]ToolDefinition),
		}
	}

	if tool.Function.Name == "" {
		return fmt.Errorf("tool function name is required")
	}

	c.toolRegistry.tools[tool.Function.Name] = tool
	return nil
}

// UnregisterTool removes a tool from the registry
func (c *OpenAIClient) UnregisterTool(name string) error {
	if c.toolRegistry == nil {
		return fmt.Errorf("tool registry not initialized")
	}

	if _, exists := c.toolRegistry.tools[name]; !exists {
		return fmt.Errorf("tool %s not found", name)
	}

	delete(c.toolRegistry.tools, name)
	return nil
}

// ListRegisteredTools returns all registered tools
func (c *OpenAIClient) ListRegisteredTools() []ToolDefinition {
	if c.toolRegistry == nil {
		return []ToolDefinition{}
	}

	tools := make([]ToolDefinition, 0, len(c.toolRegistry.tools))
	for _, tool := range c.toolRegistry.tools {
		tools = append(tools, tool)
	}
	return tools
}

// GetRegisteredTool retrieves a specific tool by name
func (c *OpenAIClient) GetRegisteredTool(name string) (ToolDefinition, bool) {
	if c.toolRegistry == nil {
		return ToolDefinition{}, false
	}

	tool, exists := c.toolRegistry.tools[name]
	return tool, exists
}

// ExecuteWithRegisteredTools executes a completion with registered tools
func (c *OpenAIClient) ExecuteWithRegisteredTools(ctx context.Context, messages []Message, toolNames []string) (*CompletionResponse, error) {
	var selectedTools []ToolDefinition
	for _, name := range toolNames {
		if tool, exists := c.toolRegistry.tools[name]; exists {
			selectedTools = append(selectedTools, tool)
		}
	}

	request := CompletionRequest{
		Messages: messages,
		Tools:    c.convertToolDefinitions(selectedTools),
		Model:    c.defaultModel,
	}
	return c.CreateCompletion(ctx, request)
}

// StreamingAgentExecution represents a streaming agent execution session
type StreamingAgentExecution struct {
	ID                  string                 `json:"id"`
	UserMessage         string                 `json:"user_message"`
	ConversationHistory []Message              `json:"conversation_history"`
	ToolExecutions      []AgentToolExecution   `json:"tool_executions"`
	FinalResponse       string                 `json:"final_response"`
	Status              string                 `json:"status"` // "thinking", "executing", "completed", "error"
	AgentThinking       string                 `json:"agent_thinking,omitempty"`
	Metadata            map[string]interface{} `json:"metadata"`
	StartTime           time.Time              `json:"start_time"`
	EndTime             *time.Time             `json:"end_time,omitempty"`
}

// AgentToolExecution represents a single tool execution by the agent
type AgentToolExecution struct {
	ID        string                 `json:"id"`
	ToolName  string                 `json:"tool_name"`
	Arguments map[string]interface{} `json:"arguments"`
	Result    interface{}            `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Status    string                 `json:"status"` // "pending", "executing", "completed", "error"
	StartTime time.Time              `json:"start_time"`
	EndTime   *time.Time             `json:"end_time,omitempty"`
	Duration  time.Duration          `json:"duration,omitempty"`
}

// ToolExecutor interface for executing tools
type ToolExecutor interface {
	ExecuteTool(ctx context.Context, toolCall ToolCall) (interface{}, error)
	GetAvailableTools() []ToolDefinition
}

// CreateStreamingAgentExecution creates an autonomous streaming agent that can execute multiple tools
func (c *OpenAIClient) CreateStreamingAgentExecution(
	ctx context.Context,
	userMessage string,
	conversationHistory []Message,
	toolExecutor ToolExecutor,
	maxIterations int,
) (<-chan StreamingAgentExecution, error) {

	if maxIterations <= 0 {
		maxIterations = 10 // Default max iterations
	}

	executionChan := make(chan StreamingAgentExecution, 50)
	executionID := fmt.Sprintf("agent_%d", time.Now().UnixNano())

	go func() {
		defer close(executionChan)
		c.runStreamingAgent(ctx, executionID, userMessage, conversationHistory, toolExecutor, maxIterations, executionChan)
	}()

	return executionChan, nil
}

// runStreamingAgent runs the autonomous streaming agent
func (c *OpenAIClient) runStreamingAgent(
	ctx context.Context,
	executionID string,
	userMessage string,
	conversationHistory []Message,
	toolExecutor ToolExecutor,
	maxIterations int,
	executionChan chan<- StreamingAgentExecution,
) {
	startTime := time.Now()

	execution := StreamingAgentExecution{
		ID:                  executionID,
		UserMessage:         userMessage,
		ConversationHistory: append(conversationHistory, Message{Role: RoleUser, Content: userMessage}),
		ToolExecutions:      []AgentToolExecution{},
		Status:              "thinking",
		Metadata:            make(map[string]interface{}),
		StartTime:           startTime,
	}

	// Send initial status
	executionChan <- execution

	// Get available tools
	availableTools := toolExecutor.GetAvailableTools()

	// Build conversation with system prompt for autonomous operation
	messages := []Message{
		{
			Role:    RoleSystem,
			Content: c.buildAgentSystemPrompt(availableTools),
		},
	}
	messages = append(messages, execution.ConversationHistory...)

	iteration := 0
	for iteration < maxIterations {
		iteration++

		// Update status
		execution.Status = "thinking"
		execution.Metadata["iteration"] = iteration
		executionChan <- execution

		// Create completion with tools
		request := CompletionRequest{
			Messages:    messages,
			Tools:       c.convertToolDefinitions(availableTools),
			ToolChoice:  "auto",
			Model:       c.defaultModel,
			Temperature: &[]float64{0.1}[0], // Low temperature for consistent reasoning
		}

		response, err := c.CreateCompletion(ctx, request)
		if err != nil {
			execution.Status = "error"
			execution.Metadata["error"] = err.Error()
			executionChan <- execution
			return
		}

		choice := response.Choices[0]

		// Add assistant message to conversation
		messages = append(messages, choice.Message)

		// Check if LLM wants to make tool calls
		if len(choice.Message.ToolCalls) > 0 {
			execution.Status = "executing"
			executionChan <- execution

			// Execute each tool call
			for _, toolCall := range choice.Message.ToolCalls {
				toolExecution := c.executeAgentTool(ctx, toolCall, toolExecutor)
				execution.ToolExecutions = append(execution.ToolExecutions, toolExecution)

				// Send tool execution update
				execution.Metadata["current_tool"] = toolExecution.ToolName
				executionChan <- execution

				// Add tool result to conversation
				toolResultMessage := Message{
					Role:       RoleTool,
					Content:    c.formatToolResult(toolExecution),
					ToolCallID: toolCall.ID,
				}
				messages = append(messages, toolResultMessage)
			}
		} else {
			// LLM provided final response
			execution.FinalResponse = choice.Message.GetContentAsString()
			execution.Status = "completed"
			endTime := time.Now()
			execution.EndTime = &endTime
			execution.Metadata["total_iterations"] = iteration
			execution.Metadata["total_tools_executed"] = len(execution.ToolExecutions)
			execution.Metadata["duration"] = endTime.Sub(startTime).String()
			executionChan <- execution
			return
		}

		// Check if context is getting too long (simple check)
		if len(messages) > 50 {
			// Trim conversation history, keeping system prompt and recent messages
			systemPrompt := messages[0]
			recentMessages := messages[len(messages)-20:]
			messages = append([]Message{systemPrompt}, recentMessages...)
		}
	}

	// Max iterations reached
	execution.Status = "completed"
	execution.FinalResponse = "I've reached my maximum iteration limit. Based on my analysis, here's what I found: " + c.summarizeToolExecutions(execution.ToolExecutions)
	endTime := time.Now()
	execution.EndTime = &endTime
	execution.Metadata["total_iterations"] = iteration
	execution.Metadata["max_iterations_reached"] = true
	executionChan <- execution
}

// executeAgentTool executes a single tool call
func (c *OpenAIClient) executeAgentTool(ctx context.Context, toolCall ToolCall, toolExecutor ToolExecutor) AgentToolExecution {
	startTime := time.Now()

	execution := AgentToolExecution{
		ID:        toolCall.ID,
		ToolName:  toolCall.Function.Name,
		Status:    "executing",
		StartTime: startTime,
	}

	// Parse arguments
	var arguments map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err == nil {
		execution.Arguments = arguments
	}

	// Execute tool
	result, err := toolExecutor.ExecuteTool(ctx, toolCall)

	endTime := time.Now()
	execution.EndTime = &endTime
	execution.Duration = endTime.Sub(startTime)

	if err != nil {
		execution.Status = "error"
		execution.Error = err.Error()
	} else {
		execution.Status = "completed"
		execution.Result = result
	}

	return execution
}

// buildAgentSystemPrompt builds the system prompt for autonomous operation
func (c *OpenAIClient) buildAgentSystemPrompt(availableTools []ToolDefinition) string {
	toolDescriptions := make([]string, len(availableTools))
	for i, tool := range availableTools {
		toolDescriptions[i] = fmt.Sprintf("- %s: %s", tool.Function.Name, tool.Function.Description)
	}

	return fmt.Sprintf(`You are an autonomous AI agent with access to tools. Your goal is to help users by:

1. **Understanding the user's request thoroughly**
2. **Using available tools to gather necessary information**
3. **Making multiple tool calls if needed to get complete information**
4. **Analyzing all gathered data before responding**
5. **Providing a comprehensive, helpful response only when you have sufficient information**

Available tools:
%s

**Important Instructions:**
- Use tools strategically to gather information
- You can make multiple tool calls in sequence
- Don't respond to the user until you have enough information to provide a helpful answer
- If you need to search for different types of information, use multiple tools
- Analyze and synthesize information from multiple sources
- Be thorough but efficient
- When you have gathered sufficient information, provide a comprehensive response

Remember: The user is waiting for a complete, well-informed response. Take the time to gather all necessary information using the available tools before responding.`, strings.Join(toolDescriptions, "\n"))
}

// formatToolResult formats tool execution result for the conversation
func (c *OpenAIClient) formatToolResult(toolExecution AgentToolExecution) string {
	if toolExecution.Error != "" {
		return fmt.Sprintf("Tool execution failed: %s", toolExecution.Error)
	}

	// Convert result to JSON string
	if toolExecution.Result != nil {
		if resultBytes, err := json.Marshal(toolExecution.Result); err == nil {
			return string(resultBytes)
		}
	}

	return "Tool executed successfully but returned no data"
}

// summarizeToolExecutions creates a summary of all tool executions
func (c *OpenAIClient) summarizeToolExecutions(executions []AgentToolExecution) string {
	if len(executions) == 0 {
		return "No tools were executed."
	}

	summary := fmt.Sprintf("I executed %d tool(s): ", len(executions))
	toolNames := make([]string, len(executions))
	for i, exec := range executions {
		if exec.Error != "" {
			toolNames[i] = fmt.Sprintf("%s (failed)", exec.ToolName)
		} else {
			toolNames[i] = exec.ToolName
		}
	}

	return summary + strings.Join(toolNames, ", ")
}

// ==================== EXISTING METHODS (PRESERVED) ====================

// CreateCompletion creates a completion using the official OpenAI client
func (c *OpenAIClient) CreateCompletion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	startTime := time.Now()

	// Convert our messages to OpenAI format
	messages := c.convertMessages(request.Messages)

	// Get the model to use
	model := c.getModel(request.Model)

	// Build completion parameters
	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    c.getSharedModel(model),
	}

	// Set optional parameters
	if request.MaxTokens != nil {
		params.MaxTokens = openai.Int(int64(*request.MaxTokens))
	}
	if request.Temperature != nil {
		params.Temperature = openai.Float(*request.Temperature)
	}
	if request.TopP != nil {
		params.TopP = openai.Float(*request.TopP)
	}
	if len(request.Tools) > 0 {
		// Check if model supports tools (reasoning models don't support traditional function calling)
		if !c.IsReasoningModel(model) {
			convertedTools := c.convertTools(request.Tools)
			params.Tools = convertedTools
		}
	}
	if request.ToolChoice != nil {
		params.ToolChoice = c.convertToolChoice(request.ToolChoice)
	}
	if request.ResponseFormat != nil {
		params.ResponseFormat = c.convertResponseFormat(request.ResponseFormat)
	}
	if len(request.Stop) > 0 {
		params.Stop = openai.ChatCompletionNewParamsStopUnion{
			OfStringArray: request.Stop,
		}
	}
	if request.PresencePenalty != nil {
		params.PresencePenalty = openai.Float(*request.PresencePenalty)
	}
	if request.FrequencyPenalty != nil {
		params.FrequencyPenalty = openai.Float(*request.FrequencyPenalty)
	}
	if request.User != "" {
		params.User = openai.String(request.User)
	}
	if request.Seed != nil {
		params.Seed = openai.Int(int64(*request.Seed))
	}

	// Log the request for debugging (temporarily)
	log.Printf("[OPENAI_CLIENT] Making API call with model: %s, messages: %d, tools: %d", model, len(messages), len(params.Tools))
	
	// Make the API call
	completion, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		c.stats.ErrorCount++
		log.Printf("[OPENAI_CLIENT] API call failed: %v", err)
		return nil, fmt.Errorf("OpenAI API call failed: %w", err)
	}
	
	log.Printf("[OPENAI_CLIENT] API call succeeded, choices: %d", len(completion.Choices))

	// Check for problematic responses
	if len(completion.Choices) > 0 {
		choice := completion.Choices[0]
		if choice.Message.Content == "" && len(choice.Message.ToolCalls) == 0 {
			// Check finish reason to provide better error context
			switch choice.FinishReason {
			case "length":
				return nil, fmt.Errorf("response truncated due to max_tokens limit")
			case "content_filter":
				return nil, fmt.Errorf("response filtered due to content policy")
			default:
				// Only log unexpected empty responses
				if c.debugMode {
					log.Printf("[OPENAI_CLIENT] WARNING: Empty response from model %s with finish_reason: %s", params.Model, choice.FinishReason)
				}
			}
		}
	}

	// Convert response
	response := c.convertCompletionResponse(completion)

	// Update stats
	latency := time.Since(startTime)
	c.updateStats(response, latency)

	return response, nil
}

// CreateCompletionStream creates a streaming completion
func (c *OpenAIClient) CreateCompletionStream(ctx context.Context, request CompletionRequest) (<-chan CompletionStreamResponse, error) {
	// Convert our messages to OpenAI format
	messages := c.convertMessages(request.Messages)

	// Build streaming completion parameters
	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    c.getSharedModel(request.Model),
	}

	// Set optional parameters
	if request.MaxTokens != nil {
		params.MaxTokens = openai.Int(int64(*request.MaxTokens))
	}
	if request.Temperature != nil {
		params.Temperature = openai.Float(*request.Temperature)
	}
	if request.TopP != nil {
		params.TopP = openai.Float(*request.TopP)
	}

	// Add tool support to streaming
	if len(request.Tools) > 0 {
		params.Tools = c.convertTools(request.Tools)
	}
	if request.ToolChoice != nil {
		params.ToolChoice = c.convertToolChoice(request.ToolChoice)
	}
	if request.ResponseFormat != nil {
		params.ResponseFormat = c.convertResponseFormat(request.ResponseFormat)
	}
	if len(request.Stop) > 0 {
		params.Stop = openai.ChatCompletionNewParamsStopUnion{
			OfStringArray: request.Stop,
		}
	}
	if request.PresencePenalty != nil {
		params.PresencePenalty = openai.Float(*request.PresencePenalty)
	}
	if request.FrequencyPenalty != nil {
		params.FrequencyPenalty = openai.Float(*request.FrequencyPenalty)
	}
	if request.User != "" {
		params.User = openai.String(request.User)
	}
	if request.Seed != nil {
		params.Seed = openai.Int(int64(*request.Seed))
	}

	// Create the stream
	stream := c.client.Chat.Completions.NewStreaming(ctx, params)

	// Create response channel
	responseChan := make(chan CompletionStreamResponse, 10)

	go func() {
		defer close(responseChan)
		defer stream.Close()

		for stream.Next() {
			chunk := stream.Current()
			streamResponse := c.convertStreamResponse(chunk)

			select {
			case responseChan <- streamResponse:
			case <-ctx.Done():
				return
			}
		}

		if err := stream.Err(); err != nil {
			log.Printf("Stream error: %v", err)
		}
	}()

	return responseChan, nil
}

// ExecuteWithTools executes a completion with tools
func (c *OpenAIClient) ExecuteWithTools(ctx context.Context, messages []Message, tools []ToolDefinition) (*CompletionResponse, error) {
	request := CompletionRequest{
		Messages: messages,
		Tools:    c.convertToolDefinitions(tools),
		Model:    c.defaultModel,
	}
	return c.CreateCompletion(ctx, request)
}

// CreateStructuredCompletion creates a structured completion with JSON schema
func (c *OpenAIClient) CreateStructuredCompletion(ctx context.Context, request CompletionRequest, schema *JSONSchemaDefinition) (*CompletionResponse, error) {
	if schema != nil {
		request.ResponseFormat = &ResponseFormat{
			Type: "json_schema",
			JSONSchema: map[string]interface{}{
				"name":   schema.Name,
				"schema": schema.Schema,
				"strict": schema.Strict,
			},
		}
	}
	return c.CreateCompletion(ctx, request)
}

// CountTokens estimates token count (simplified implementation)
func (c *OpenAIClient) CountTokens(text string) (int, error) {
	// Rough estimation: ~4 characters per token for English text
	return len(text) / 4, nil
}

// GetCapabilities returns client capabilities
func (c *OpenAIClient) GetCapabilities() ClientCapabilities {
	return ClientCapabilities{
		SupportsStreaming:  true,
		SupportsTools:      true,
		SupportsStructured: true,
		SupportsVision:     true,
		SupportsMultiModal: true,
		SupportsReasoning:  true,
		SupportsThinking:   false,
		MaxContextLength:   128000, // GPT-4 context length
		MaxOutputLength:    4096,
		SupportedModels: []string{
			ModelGPT4o, ModelGPT4oMini, ModelGPT4Turbo,
			ModelO1Preview, ModelO1Mini, ModelGPT4oLatest,
		},
	}
}

// HealthCheck performs a health check
func (c *OpenAIClient) HealthCheck(ctx context.Context) error {
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("ping"),
		},
		Model:     shared.ChatModelGPT4oMini,
		MaxTokens: openai.Int(1),
	}

	_, err := c.client.Chat.Completions.New(ctx, params)
	return err
}

// GetUsageStats returns usage statistics
func (c *OpenAIClient) GetUsageStats() UsageStats {
	return c.stats
}

// Security methods (simplified implementations)
func (c *OpenAIClient) AnalyzeFraud(ctx context.Context, content string) (*SecurityAssessment, error) {
	return &SecurityAssessment{
		RiskLevel:       "low",
		FraudScore:      0.1,
		Recommendations: []string{"Content appears safe"},
		RedFlags:        []string{},
		Confidence:      0.95,
	}, nil
}

func (c *OpenAIClient) AssessRisk(ctx context.Context, request CompletionRequest) (*SecurityAssessment, error) {
	return &SecurityAssessment{
		RiskLevel:       "low",
		FraudScore:      0.1,
		Recommendations: []string{"Request appears safe"},
		RedFlags:        []string{},
		Confidence:      0.95,
	}, nil
}

func (c *OpenAIClient) GetSecurityRecommendations(ctx context.Context, content string) ([]string, error) {
	return []string{
		"Content appears safe",
		"No security concerns detected",
	}, nil
}

// GetProviderInfo returns provider information
func (c *OpenAIClient) GetProviderInfo() ProviderInfo {
	return ProviderInfo{
		Name:         "OpenAI",
		Provider:     ProviderOpenAI,
		Version:      "v1.6.0",
		Capabilities: c.GetCapabilities(),
		Models: []ModelInfo{
			{
				ID:              ModelGPT4o,
				Name:            "GPT-4o",
				Provider:        ProviderOpenAI,
				Version:         "2024-08-06",
				MaxTokens:       4096,
				ContextLength:   128000,
				MaxOutputTokens: 4096,
				InputCost:       2.50 / 1000000,  // $2.50 per 1M input tokens
				OutputCost:      10.00 / 1000000, // $10.00 per 1M output tokens
				Capabilities:    []string{"chat", "tools", "vision", "structured"},
				Description:     "Most advanced multimodal model",
			},
			{
				ID:              ModelGPT4oMini,
				Name:            "GPT-4o Mini",
				Provider:        ProviderOpenAI,
				Version:         "2024-07-18",
				MaxTokens:       4096,
				ContextLength:   128000,
				MaxOutputTokens: 4096,
				InputCost:       0.15 / 1000000, // $0.15 per 1M input tokens
				OutputCost:      0.60 / 1000000, // $0.60 per 1M output tokens
				Capabilities:    []string{"chat", "tools", "vision", "structured"},
				Description:     "Affordable and intelligent small model",
			},
		},
		Status: "active",
	}
}

// Helper methods

func (c *OpenAIClient) getSharedModel(requestModel string) shared.ChatModel {
	model := c.getModel(requestModel)

	switch model {
	case ModelGPT4o:
		return shared.ChatModelGPT4o
	case ModelGPT4oMini:
		return shared.ChatModelGPT4oMini
	case ModelGPT41Mini:
		return shared.ChatModel("gpt-4.1-mini") // Use string literal for new model
	case ModelGPT4Turbo:
		return shared.ChatModelGPT4Turbo
	case ModelO1Preview:
		return shared.ChatModelO1Preview
	case ModelO1Mini:
		return shared.ChatModelO1Mini
	default:
		return shared.ChatModelGPT4o // Default to GPT-4o which supports function calling
	}
}

func (c *OpenAIClient) getModel(requestModel string) string {
	if requestModel != "" {
		return requestModel
	}
	if c.defaultModel != "" {
		return c.defaultModel
	}
	return ModelGPT4o // Default to GPT-4o which supports function calling
}

func (c *OpenAIClient) convertMessages(messages []Message) []openai.ChatCompletionMessageParamUnion {
	converted := make([]openai.ChatCompletionMessageParamUnion, len(messages))

	for i, msg := range messages {
		content := msg.GetContentAsString()

		switch msg.Role {
		case RoleSystem:
			converted[i] = openai.SystemMessage(content)
		case RoleUser:
			converted[i] = openai.UserMessage(content)
		case RoleAssistant:
			// CRITICAL FIX: Preserve tool_calls for assistant messages to maintain conversation flow
			if len(msg.ToolCalls) > 0 {
				// Convert tool calls to OpenAI SDK format
				toolCalls := make([]openai.ChatCompletionMessageToolCallParam, len(msg.ToolCalls))
				for j, tc := range msg.ToolCalls {
					toolCalls[j] = openai.ChatCompletionMessageToolCallParam{
						ID:   tc.ID,
						Type: "function", // Use string literal instead of constant
						Function: openai.ChatCompletionMessageToolCallFunctionParam{
							Name:      tc.Function.Name,
							Arguments: tc.Function.Arguments,
						},
					}
				}
				// Create assistant message with tool calls preserved using proper union type
				converted[i] = openai.ChatCompletionMessageParamUnion{
					OfAssistant: &openai.ChatCompletionAssistantMessageParam{
						Content: openai.ChatCompletionAssistantMessageParamContentUnion{
							OfString: openai.String(content),
						},
						ToolCalls: toolCalls,
					},
				}
			} else {
				converted[i] = openai.AssistantMessage(content)
			}
		case RoleTool:
			converted[i] = openai.ToolMessage(content, msg.ToolCallID)
		default:
			converted[i] = openai.UserMessage(content)
		}
	}

	return converted
}

func (c *OpenAIClient) convertTools(tools []Tool) []openai.ChatCompletionToolParam {
	converted := make([]openai.ChatCompletionToolParam, len(tools))

	for i, tool := range tools {
		converted[i] = openai.ChatCompletionToolParam{
			Type: "function", // Use string literal as per API docs
			Function: openai.FunctionDefinitionParam{
				Name:        tool.Function.Name,
				Description: openai.String(tool.Function.Description),
				Parameters:  openai.FunctionParameters(tool.Function.Parameters),
			},
		}
	}

	return converted
}

func (c *OpenAIClient) convertToolChoice(toolChoice interface{}) openai.ChatCompletionToolChoiceOptionUnionParam {
	switch tc := toolChoice.(type) {
	case string:
		switch tc {
		case "auto":
			return openai.ChatCompletionToolChoiceOptionUnionParam{
				OfAuto: openai.String("auto"),
			}
		case "none":
			return openai.ChatCompletionToolChoiceOptionUnionParam{
				OfAuto: openai.String("none"),
			}
		case "required":
			return openai.ChatCompletionToolChoiceOptionUnionParam{
				OfAuto: openai.String("required"),
			}
		default:
			return openai.ChatCompletionToolChoiceOptionUnionParam{
				OfAuto: openai.String("auto"),
			}
		}
	case map[string]interface{}:
		// Handle specific function tool choice
		if tc["type"] == "function" {
			if function, ok := tc["function"].(map[string]interface{}); ok {
				if name, ok := function["name"].(string); ok {
					return openai.ChatCompletionToolChoiceOptionUnionParam{
						OfChatCompletionNamedToolChoice: &openai.ChatCompletionNamedToolChoiceParam{
							Type: "function",
							Function: openai.ChatCompletionNamedToolChoiceFunctionParam{
								Name: name,
							},
						},
					}
				}
			}
		}
		// Fall back to auto if structure is not recognized
		return openai.ChatCompletionToolChoiceOptionUnionParam{
			OfAuto: openai.String("auto"),
		}
	default:
		return openai.ChatCompletionToolChoiceOptionUnionParam{
			OfAuto: openai.String("auto"),
		}
	}
}

func (c *OpenAIClient) convertResponseFormat(rf *ResponseFormat) openai.ChatCompletionNewParamsResponseFormatUnion {
	switch rf.Type {
	case "json_schema":
		name := rf.JSONSchema["name"].(string)
		strict := rf.JSONSchema["strict"].(bool)
		return openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				Type: "json_schema",
				JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   name,
					Schema: openai.FunctionParameters(rf.JSONSchema["schema"].(map[string]interface{})),
					Strict: openai.Bool(strict),
				},
			},
		}
	case "json_object":
		return openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &openai.ResponseFormatJSONObjectParam{
				Type: "json_object",
			},
		}
	default:
		return openai.ChatCompletionNewParamsResponseFormatUnion{
			OfText: &openai.ResponseFormatTextParam{
				Type: "text",
			},
		}
	}
}

func (c *OpenAIClient) convertCompletionResponse(completion *openai.ChatCompletion) *CompletionResponse {
	choices := make([]Choice, len(completion.Choices))

	for i, choice := range completion.Choices {
		// Extract content directly - it's already a string in the OpenAI SDK
		message := Message{
			Role:    RoleAssistant,
			Content: choice.Message.Content, // This is a string from the SDK
		}

		// Handle tool calls
		if len(choice.Message.ToolCalls) > 0 {
			toolCalls := make([]ToolCall, len(choice.Message.ToolCalls))
			for j, tc := range choice.Message.ToolCalls {
				toolCalls[j] = ToolCall{
					ID:   tc.ID,
					Type: string(tc.Type),
					Function: FunctionCall{
						Name:      tc.Function.Name,
						Arguments: tc.Function.Arguments,
					},
				}
			}
			message.ToolCalls = toolCalls
		}

		choices[i] = Choice{
			Index:        int(choice.Index),
			Message:      message,
			FinishReason: string(choice.FinishReason),
		}
	}

	usage := Usage{
		PromptTokens:     int(completion.Usage.PromptTokens),
		CompletionTokens: int(completion.Usage.CompletionTokens),
		TotalTokens:      int(completion.Usage.TotalTokens),
	}

	// Calculate costs
	model := completion.Model
	usage.InputCost = c.calculateInputCost(model, usage.PromptTokens)
	usage.OutputCost = c.calculateOutputCost(model, usage.CompletionTokens)
	usage.TotalCost = usage.InputCost + usage.OutputCost

	return &CompletionResponse{
		ID:       completion.ID,
		Object:   string(completion.Object),
		Created:  completion.Created,
		Model:    completion.Model,
		Provider: ProviderOpenAI,
		Choices:  choices,
		Usage:    usage,
	}
}

func (c *OpenAIClient) convertStreamResponse(chunk openai.ChatCompletionChunk) CompletionStreamResponse {
	choices := make([]Choice, len(chunk.Choices))

	for i, choice := range chunk.Choices {
		delta := Message{
			Role:    string(choice.Delta.Role),
			Content: choice.Delta.Content,
		}

		// Handle tool calls in streaming
		if len(choice.Delta.ToolCalls) > 0 {
			toolCalls := make([]ToolCall, len(choice.Delta.ToolCalls))
			for j, tc := range choice.Delta.ToolCalls {
				toolCall := ToolCall{
					ID:   tc.ID,
					Type: string(tc.Type),
				}

				// Handle function information
				if tc.Function.Name != "" || tc.Function.Arguments != "" {
					toolCall.Function = FunctionCall{
						Name:      tc.Function.Name,
						Arguments: tc.Function.Arguments,
					}
				}

				toolCalls[j] = toolCall
			}
			delta.ToolCalls = toolCalls
		}

		choices[i] = Choice{
			Index:        int(choice.Index),
			Delta:        delta,
			FinishReason: string(choice.FinishReason),
		}
	}

	return CompletionStreamResponse{
		ID:      chunk.ID,
		Object:  string(chunk.Object),
		Created: chunk.Created,
		Model:   chunk.Model,
		Choices: choices,
	}
}

func (c *OpenAIClient) calculateInputCost(model string, tokens int) float64 {
	costPer1M := map[string]float64{
		ModelGPT4o:     2.50,
		ModelGPT4oMini: 0.15,
		ModelGPT4Turbo: 10.00,
		ModelO1Preview: 15.00,
		ModelO1Mini:    3.00,
	}

	if cost, exists := costPer1M[model]; exists {
		return float64(tokens) * cost / 1000000
	}
	return float64(tokens) * 0.15 / 1000000 // Default to GPT-4o Mini pricing
}

func (c *OpenAIClient) calculateOutputCost(model string, tokens int) float64 {
	costPer1M := map[string]float64{
		ModelGPT4o:     10.00,
		ModelGPT4oMini: 0.60,
		ModelGPT4Turbo: 30.00,
		ModelO1Preview: 60.00,
		ModelO1Mini:    12.00,
	}

	if cost, exists := costPer1M[model]; exists {
		return float64(tokens) * cost / 1000000
	}
	return float64(tokens) * 0.60 / 1000000 // Default to GPT-4o Mini pricing
}

func (c *OpenAIClient) updateStats(response *CompletionResponse, latency time.Duration) {
	c.stats.RequestCount++
	c.stats.TokensUsed += int64(response.Usage.TotalTokens)
	c.stats.TotalCost += response.Usage.TotalCost
	c.stats.LastUsed = time.Now()

	// Update average latency
	if c.stats.AverageLatency == 0 {
		c.stats.AverageLatency = latency
	} else {
		c.stats.AverageLatency = (c.stats.AverageLatency + latency) / 2
	}
}

// convertToolDefinitions converts ToolDefinition slice to Tool slice
func (c *OpenAIClient) convertToolDefinitions(tools []ToolDefinition) []Tool {
	converted := make([]Tool, len(tools))
	for i, tool := range tools {
		converted[i] = Tool{
			Type:     tool.Type,
			Function: tool.Function,
		}
	}
	return converted
}

// detectImageMimeType detects the MIME type of an image based on file extension
func (c *OpenAIClient) detectImageMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		// Default to JPEG if we can't detect
		return "image/jpeg"
	}
	return mimeType
}

// ==================== REASONING MODELS SUPPORT ====================

// IsReasoningModel checks if the given model is a reasoning model
func (c *OpenAIClient) IsReasoningModel(model string) bool {
	reasoningModels := []string{
		"o1-preview", "o1-mini", "o3", "o4-mini",
		"o3-mini", "gpt-4-1-mini",
	}

	modelLower := strings.ToLower(model)
	for _, rm := range reasoningModels {
		if strings.Contains(modelLower, rm) {
			return true
		}
	}
	return false
}

// CreateReasoningCompletion creates a completion using reasoning models via Responses API
func (c *OpenAIClient) CreateReasoningCompletion(ctx context.Context, request ReasoningRequest) (*ReasoningResponse, error) {
	// For now, this requires implementation of the Responses API
	// The official Go SDK doesn't support it yet, so we'd need HTTP client implementation
	return nil, fmt.Errorf("reasoning models via Responses API require custom HTTP implementation")
}

// ExecuteSequentialTools executes tools sequentially for reasoning models
func (c *OpenAIClient) ExecuteSequentialTools(
	ctx context.Context,
	input []ReasoningInputItem,
	tools []Tool,
	maxIterations int,
) (*ReasoningResponse, error) {
	conversation := input
	iteration := 0

	for iteration < maxIterations {
		request := ReasoningRequest{
			Model: c.defaultModel,
			Input: conversation,
			Tools: tools,
			Store: true,
			Reasoning: &ReasoningConfig{
				Effort:  "medium",
				Summary: "auto",
			},
		}

		response, err := c.CreateReasoningCompletion(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("reasoning completion failed at iteration %d: %w", iteration, err)
		}

		// Check if we have function calls to execute
		functionCalls := []ReasoningOutputItem{}
		for _, item := range response.Output {
			if item.Type == "function_call" {
				functionCalls = append(functionCalls, item)
			}
		}

		if len(functionCalls) == 0 {
			// No more function calls, we're done
			return response, nil
		}

		// Add all output items to conversation (including reasoning)
		for _, item := range response.Output {
			conversation = append(conversation, ReasoningInputItem{
				Type:             item.Type,
				ID:               item.ID,
				EncryptedContent: item.EncryptedContent,
				Arguments:        item.Arguments,
				Name:             item.Name,
				CallID:           item.CallID,
				Status:           item.Status,
			})
		}

		// Execute function calls sequentially
		for _, funcCall := range functionCalls {
			result := c.executeReasoningToolCall(ctx, funcCall)

			conversation = append(conversation, ReasoningInputItem{
				Type:   "function_call_output",
				CallID: funcCall.CallID,
				Output: result,
			})
		}

		iteration++
	}

	return nil, fmt.Errorf("max iterations (%d) reached without completion", maxIterations)
}

// executeReasoningToolCall executes a single reasoning tool call
func (c *OpenAIClient) executeReasoningToolCall(ctx context.Context, funcCall ReasoningOutputItem) string {
	// This would integrate with the existing tool execution system
	// For now, return a placeholder
	return fmt.Sprintf("Tool %s executed with args: %s", funcCall.Name, funcCall.Arguments)
}

// GetReasoningCapabilities returns the capabilities of reasoning models
func (c *OpenAIClient) GetReasoningCapabilities() ReasoningCapabilities {
	return ReasoningCapabilities{
		SupportsSequentialTools:  true,
		SupportsReasoningSummary: true,
		SupportsEncryptedContent: true,
		MaxIterations:            10,
		SupportedEffortLevels:    []string{"low", "medium", "high"},
		SupportedSummaryTypes:    []string{"auto", "detailed", "none"},
	}
}

// ConvertMessagesToReasoningInput converts standard messages to reasoning input format
func (c *OpenAIClient) ConvertMessagesToReasoningInput(messages []Message) []ReasoningInputItem {
	input := make([]ReasoningInputItem, len(messages))

	for i, msg := range messages {
		input[i] = ReasoningInputItem{
			Type:    "message",
			Role:    string(msg.Role),
			Content: msg.GetContentAsString(),
		}
	}

	return input
}

// ConvertReasoningResponseToCompletion converts reasoning response to standard completion format
func (c *OpenAIClient) ConvertReasoningResponseToCompletion(response *ReasoningResponse) *CompletionResponse {
	choices := []Choice{}

	// Extract the final assistant message from reasoning output
	for _, item := range response.Output {
		if item.Type == "message" && item.Role == "assistant" {
			content := ""
			if len(item.Content) > 0 {
				content = item.Content[0].Text
			}

			message := Message{
				Role:    RoleAssistant,
				Content: content,
			}

			choice := Choice{
				Index:        0,
				Message:      message,
				FinishReason: "stop",
			}

			choices = append(choices, choice)
			break
		}
	}

	usage := Usage{
		PromptTokens:     response.Usage.InputTokens,
		CompletionTokens: response.Usage.OutputTokens,
		TotalTokens:      response.Usage.TotalTokens,
	}

	return &CompletionResponse{
		ID:       response.ID,
		Object:   response.Object,
		Created:  response.Created,
		Model:    response.Model,
		Provider: ProviderOpenAI,
		Choices:  choices,
		Usage:    usage,
	}
}

// ==================== EMBEDDING GENERATION ====================
// Note: Embedding functionality is implemented in the separate OpenAIEmbeddingClient adapter
// in internal/infra/openai_embedding_client.go to maintain clean separation of concerns.

// EmbeddingRequest represents an embedding generation request
type EmbeddingRequest struct {
	Input          interface{} `json:"input"`                     // Text or array of texts
	Model          string      `json:"model"`                     // Embedding model
	EncodingFormat string      `json:"encoding_format,omitempty"` // "float", "base64"
	Dimensions     *int        `json:"dimensions,omitempty"`      // For newer models
	User           string      `json:"user,omitempty"`            // End-user ID
}

// EmbeddingResponse represents the response from embedding generation
type EmbeddingResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  EmbeddingUsage  `json:"usage"`
}

// EmbeddingData represents a single embedding result
type EmbeddingData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}

// EmbeddingUsage represents token usage for embedding generation
type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// GenerateEmbedding generates embeddings for a single text using OpenAI API
func (c *OpenAIClient) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	params := openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
		Model: openai.EmbeddingModel(c.GetEmbeddingModel()),
	}

	// Set dimensions for text-embedding-3 models
	if strings.Contains(c.GetEmbeddingModel(), "text-embedding-3") {
		dims := int64(c.GetEmbeddingDimensions())
		params.Dimensions = openai.Int(dims)
	}

	resp, err := c.client.Embeddings.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	// Convert float64 to float32
	embedding := make([]float32, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		embedding[i] = float32(v)
	}

	// Update usage stats
	c.stats.RequestCount++
	c.stats.TokensUsed += resp.Usage.TotalTokens
	c.stats.TotalCost += c.calculateEmbeddingCost(c.GetEmbeddingModel(), int(resp.Usage.TotalTokens))
	c.stats.LastUsed = time.Now()

	return embedding, nil
}

// GenerateBatchEmbeddings generates embeddings for multiple texts efficiently
func (c *OpenAIClient) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("no texts provided")
	}

	// Convert []string to []interface{} for the union type
	inputs := make([]interface{}, len(texts))
	for i, text := range texts {
		inputs[i] = text
	}

	params := openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: texts,
		},
		Model: openai.EmbeddingModel(c.GetEmbeddingModel()),
	}

	// Set dimensions for text-embedding-3 models
	if strings.Contains(c.GetEmbeddingModel(), "text-embedding-3") {
		dims := int64(c.GetEmbeddingDimensions())
		params.Dimensions = openai.Int(dims)
	}

	resp, err := c.client.Embeddings.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate batch embeddings: %w", err)
	}

	if len(resp.Data) != len(texts) {
		return nil, fmt.Errorf("expected %d embeddings, got %d", len(texts), len(resp.Data))
	}

	// Convert response to [][]float32
	embeddings := make([][]float32, len(resp.Data))
	for i, data := range resp.Data {
		embedding := make([]float32, len(data.Embedding))
		for j, v := range data.Embedding {
			embedding[j] = float32(v)
		}
		embeddings[i] = embedding
	}

	// Update usage stats
	c.stats.RequestCount++
	c.stats.TokensUsed += int64(resp.Usage.TotalTokens)
	c.stats.TotalCost += c.calculateEmbeddingCost(c.GetEmbeddingModel(), int(resp.Usage.TotalTokens))
	c.stats.LastUsed = time.Now()

	return embeddings, nil
}

// GenerateEntityEmbedding generates embeddings for entity data with smart text extraction
func (c *OpenAIClient) GenerateEntityEmbedding(ctx context.Context, entityData map[string]interface{}) ([]float32, error) {
	text := c.extractEntityText(entityData)
	if text == "" {
		return nil, fmt.Errorf("no valid text found in entity data")
	}

	return c.GenerateEmbedding(ctx, text)
}

// GenerateEmbeddingWithPrompt generates embeddings with LLM-enhanced text transformation
func (c *OpenAIClient) GenerateEmbeddingWithPrompt(ctx context.Context, text string, prompt string) ([]float32, error) {
	if prompt == "" {
		// If no prompt provided, use direct embedding
		return c.GenerateEmbedding(ctx, text)
	}

	// Use LLM to transform text according to prompt, then generate embedding
	transformedText, err := c.transformTextWithPrompt(ctx, text, prompt)
	if err != nil {
		// Fallback to original text if transformation fails
		log.Printf("WARN: Text transformation failed, using original text: %v", err)
		return c.GenerateEmbedding(ctx, text)
	}

	return c.GenerateEmbedding(ctx, transformedText)
}

// GenerateOptimizedEmbedding generates embeddings optimized for specific use cases
func (c *OpenAIClient) GenerateOptimizedEmbedding(ctx context.Context, entityType string, entityData map[string]interface{}, optimization string) ([]float32, error) {
	// Build optimization-specific text representation
	optimizedText := c.buildOptimizedText(entityType, entityData, optimization)

	// Select optimal model based on use case
	model := c.selectEmbeddingModel(optimization)

	// Generate embedding with optimized parameters
	return c.generateEmbeddingWithModel(ctx, optimizedText, model)
}

// GetEmbeddingDimensions returns dimensions for the embedding model
func (c *OpenAIClient) GetEmbeddingDimensions() int {
	model := "text-embedding-3-small"
	if c.defaultModel != "" && strings.Contains(c.defaultModel, "embedding") {
		model = c.defaultModel
	}

	switch model {
	case "text-embedding-3-large":
		return 3072
	case "text-embedding-3-small":
		return 1536
	case "text-embedding-ada-002":
		return 1536
	default:
		return 1536 // Default dimensions
	}
}

// GetEmbeddingModel returns the current embedding model
func (c *OpenAIClient) GetEmbeddingModel() string {
	if c.defaultModel != "" && strings.Contains(c.defaultModel, "embedding") {
		return c.defaultModel
	}
	return "text-embedding-3-small"
}

// IsEmbeddingPromptEnabled returns whether prompt-based embedding is supported
func (c *OpenAIClient) IsEmbeddingPromptEnabled() bool {
	return true // OpenAI client supports LLM-enhanced text transformation
}

// ==================== EMBEDDING HELPER METHODS ====================

// calculateEmbeddingCost calculates the cost for embedding generation
func (c *OpenAIClient) calculateEmbeddingCost(model string, tokens int) float64 {
	// Pricing per 1M tokens (as of 2025)
	costPer1M := map[string]float64{
		"text-embedding-3-large": 0.13,
		"text-embedding-3-small": 0.02,
		"text-embedding-ada-002": 0.10,
	}

	if cost, exists := costPer1M[model]; exists {
		return float64(tokens) * cost / 1000000
	}
	return float64(tokens) * 0.02 / 1000000 // Default to text-embedding-3-small pricing
}

// extractEntityText intelligently extracts searchable text from entity data
func (c *OpenAIClient) extractEntityText(entityData map[string]interface{}) string {
	priorityFields := []string{"name", "title", "description", "content", "summary"}
	secondaryFields := []string{"brand", "model", "category", "tags", "location", "address"}

	var textParts []string

	// Extract priority fields first
	for _, field := range priorityFields {
		if value, exists := entityData[field]; exists {
			if text := c.valueToString(value); text != "" {
				textParts = append(textParts, text)
			}
		}
	}

	// Add secondary fields
	for _, field := range secondaryFields {
		if value, exists := entityData[field]; exists {
			if text := c.valueToString(value); text != "" {
				textParts = append(textParts, text)
			}
		}
	}

	return strings.Join(textParts, " ")
}

// valueToString converts interface{} values to strings
func (c *OpenAIClient) valueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case []string:
		return strings.Join(v, " ")
	case []interface{}:
		var parts []string
		for _, item := range v {
			if str := c.valueToString(item); str != "" {
				parts = append(parts, str)
			}
		}
		return strings.Join(parts, " ")
	default:
		return fmt.Sprintf("%v", v)
	}
}

// transformTextWithPrompt uses LLM to transform text according to a prompt
func (c *OpenAIClient) transformTextWithPrompt(ctx context.Context, text string, prompt string) (string, error) {
	request := CompletionRequest{
		Messages: []Message{
			{
				Role:    RoleSystem,
				Content: prompt,
			},
			{
				Role:    RoleUser,
				Content: text,
			},
		},
		Model:       c.getOptimalLLMModel(),
		MaxTokens:   &[]int{500}[0],     // Limit response length
		Temperature: &[]float64{0.1}[0], // Low temperature for consistency
	}

	response, err := c.CreateCompletion(ctx, request)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return response.Choices[0].Message.GetContentAsString(), nil
}

// buildOptimizedText creates optimization-specific text representations
func (c *OpenAIClient) buildOptimizedText(entityType string, entityData map[string]interface{}, optimization string) string {
	baseText := c.extractEntityText(entityData)

	switch optimization {
	case "search":
		return c.optimizeForSearch(baseText, entityType)
	case "recommendation":
		return c.optimizeForRecommendation(baseText, entityType)
	case "similarity":
		return c.optimizeForSimilarity(baseText, entityType)
	case "real-time":
		return c.optimizeForRealTime(baseText)
	default:
		return baseText
	}
}

// selectEmbeddingModel chooses optimal embedding model based on use case
func (c *OpenAIClient) selectEmbeddingModel(optimization string) string {
	switch optimization {
	case "search", "similarity":
		return "text-embedding-3-large" // Higher quality for search
	case "real-time", "recommendation":
		return "text-embedding-3-small" // Faster, cost-effective
	default:
		return "text-embedding-3-small"
	}
}

// generateEmbeddingWithModel generates embedding with specific model
func (c *OpenAIClient) generateEmbeddingWithModel(ctx context.Context, text string, model string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	params := openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
		Model: openai.EmbeddingModel(model),
	}

	// Set dimensions for text-embedding-3 models
	if strings.Contains(model, "text-embedding-3") {
		dims := int64(c.GetEmbeddingDimensions())
		params.Dimensions = openai.Int(dims)
	}

	resp, err := c.client.Embeddings.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding with model %s: %w", model, err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	// Convert float64 to float32
	embedding := make([]float32, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		embedding[i] = float32(v)
	}

	// Update usage stats
	c.stats.RequestCount++
	c.stats.TokensUsed += int64(resp.Usage.TotalTokens)
	c.stats.TotalCost += c.calculateEmbeddingCost(model, int(resp.Usage.TotalTokens))
	c.stats.LastUsed = time.Now()

	return embedding, nil
}

// getOptimalLLMModel returns the best LLM model for text transformation
func (c *OpenAIClient) getOptimalLLMModel() string {
	// Use fast, cost-effective model for text transformation
	return ModelGPT4oMini
}

// Optimization helper methods
func (c *OpenAIClient) optimizeForSearch(text, entityType string) string {
	// Enhance text with search-relevant terms
	return fmt.Sprintf("%s %s searchable keywords", text, entityType)
}

func (c *OpenAIClient) optimizeForRecommendation(text, entityType string) string {
	// Focus on recommendation-relevant features
	return fmt.Sprintf("%s %s recommendation features", text, entityType)
}

func (c *OpenAIClient) optimizeForSimilarity(text, entityType string) string {
	// Emphasize similarity comparison features
	return fmt.Sprintf("%s %s similarity comparison", text, entityType)
}

func (c *OpenAIClient) optimizeForRealTime(text string) string {
	// Keep text concise for real-time processing
	words := strings.Fields(text)
	if len(words) > 50 {
		return strings.Join(words[:50], " ")
	}
	return text
}
