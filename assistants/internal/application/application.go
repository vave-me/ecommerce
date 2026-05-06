package application

import (
	"context"
	"fmt"

	"middleman/assistants/internal/application/commands"
	"middleman/assistants/internal/application/queries"
	"middleman/assistants/internal/application/services"
	"middleman/assistants/internal/application/tools"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
	externalai "middleman/internal/ai"
	"middleman/internal/ddd"
)

// App interface defines the complete application contract
type App interface {
	Commands
	Queries
	Services
	Tools
}

// Commands interface defines all command operations
type Commands interface {
	// Assistant management
	CreateAdminAssistant(ctx context.Context, cmd commands.CreateAdminAssistant) error
	CreateBusinessAssistant(ctx context.Context, cmd commands.CreateBusinessAssistant) error
	CreateSupportAssistant(ctx context.Context, cmd commands.CreateSupportAssistant) error
	CreateSchedulerAssistant(ctx context.Context, cmd commands.CreateSchedulerAssistant) error
	CreateUserAssistant(ctx context.Context, cmd commands.CreateUserAssistant) error
	ActivateAssistant(ctx context.Context, cmd commands.ActivateAssistant) error
	DeactivateAssistant(ctx context.Context, cmd commands.DeactivateAssistant) error
	UpdateAssistantConfiguration(ctx context.Context, cmd commands.UpdateAssistantConfiguration) error

	// Input processing
	ProcessUserInput(ctx context.Context, cmd commands.ProcessUserInput) (*commands.ProcessUserInputResult, error)
	ProcessSpeechInput(ctx context.Context, cmd commands.ProcessSpeechInput) (*commands.ProcessSpeechInputResult, error)
	ProcessImageInput(ctx context.Context, cmd commands.ProcessImageInput) (*commands.ProcessImageInputResult, error)
	ProcessDocumentInput(ctx context.Context, cmd commands.ProcessDocumentInput) (*commands.ProcessDocumentInputResult, error)

	// Conversation management
	CreateConversation(ctx context.Context, cmd commands.CreateConversation) error
	AddMessageToConversation(ctx context.Context, cmd commands.AddMessageToConversation) (*commands.AddMessageToConversationResult, error)
	ChatWithConversation(ctx context.Context, cmd commands.ChatWithConversation) (*commands.ChatWithConversationResult, error)
	UpdateConversationContext(ctx context.Context, cmd commands.UpdateConversationContext) error
}

// Queries interface defines all query operations
type Queries interface {
	GetAssistant(ctx context.Context, query queries.GetAssistant) (*domain.CatalogAssistant, error)
	GetAssistants(ctx context.Context, query queries.GetAssistants) ([]*domain.CatalogAssistant, error)
	GetConversation(ctx context.Context, query queries.GetConversation) (*domain.ReadConversation, error)
	GetUserConversations(ctx context.Context, query queries.GetUserConversations) (*queries.GetUserConversationsResult, error)
	GetConversationMessages(ctx context.Context, query queries.GetConversationMessages) ([]*domain.ReadMessage, error)
	GetConversationStats(ctx context.Context, query queries.GetConversationStats) (*domain.ConversationStats, error)
	GetOrCreateUserAssistant(ctx context.Context, query queries.GetOrCreateUserAssistant) (*domain.Assistant, error)
}

// Services interface defines repository and AI service operations
type Services interface {
	Execute(ctx context.Context, query services.RepositoryQuery) (*services.RepositoryResponse, error)
	ValidateQuery(query services.RepositoryQuery) error
	GetSupportedOperations(entityType models.EntityType) []services.OperationType
	TranslateAIRequest(aiRequest map[string]interface{}) (*services.RepositoryQuery, error)
}

// Tools interface defines tool execution operations
type Tools interface {
	ExecuteTools(ctx context.Context, toolCalls []externalai.ToolCall, execCtx *tools.ToolExecutionContext) ([]tools.ToolExecutionResult, error)
	ExecuteToolsStream(ctx context.Context, toolCalls []externalai.ToolCall, execCtx *tools.ToolExecutionContext) (<-chan tools.ToolExecutionStream, error)
}

// Application is the main application service that orchestrates all operations
type Application struct {
	// Embedded structs for clean organization
	appCommands
	appQueries
	appProcessors
	appServices
	appTools
	assistants        domain.AssistantRepository
	conversations     domain.ConversationRepository
	readConversations domain.ReadConversationRepository
	readMessages      domain.ReadMessagesRepository
	catalogReader     domain.CatalogRepository
	vectorRepository  domain.VectorRepository
	publisher         ddd.EventPublisher[ddd.Event]
}

// appCommands contains all command handlers
type appCommands struct {
	// Assistant handlers
	commands.CreateAdminAssistantHandler
	commands.CreateBusinessAssistantHandler
	commands.CreateSupportAssistantHandler
	commands.CreateSchedulerAssistantHandler
	commands.CreateUserAssistantHandler
	commands.ActivateAssistantHandler
	commands.DeactivateAssistantHandler
	commands.UpdateAssistantConfigurationHandler

	// Input processing handlers
	commands.ProcessUserInputHandler
	commands.ProcessSpeechInputHandler
	commands.ProcessImageInputHandler
	commands.ProcessDocumentInputHandler

	// Conversation handlers
	commands.CreateConversationHandler
	commands.AddMessageToConversationHandler
	commands.ChatWithConversationHandler
	commands.UpdateConversationContextHandler
}

// appQueries contains all query handlers
type appQueries struct {
	queries.GetAssistantHandler
	queries.GetAssistantsHandler
	queries.GetConversationHandler
	queries.GetUserConversationsHandler
	queries.GetConversationMessagesHandler
	queries.GetConversationStatsHandler
	queries.GetOrCreateUserAssistantHandler
}

// appProcessors contains all AI processors
type appProcessors struct {
	llmProcessor      services.LLMProcessor
	speechProcessor   services.SpeechProcessor
	visionProcessor   services.VisionProcessor
	documentProcessor services.DocumentProcessor
	dataProcessor     services.DataProcessor
	promptProvider    domain.SystemPromptProvider
}

// appServices contains service-level operations
type appServices struct {
	toolRegistry         *tools.ToolRegistry
	repositoryTranslator *RepositoryTranslator
}

// appTools contains tool executors
type appTools struct {
	toolExecutor   *SimplifiedToolExecutor
	streamExecutor *SimplifiedStreamingExecutor
}

// Ensure Application implements App interface
var _ App = (*Application)(nil)

// New creates a new Application instance with all required dependencies
func New(
	assistants domain.AssistantRepository,
	conversations domain.ConversationRepository,
	readConversations domain.ReadConversationRepository,
	readMessages domain.ReadMessagesRepository,
	catalogReader domain.CatalogRepository,
	vectorRepository domain.VectorRepository,
	toolRegistry *tools.ToolRegistry,
	config *Config,
	publisher ddd.EventPublisher[ddd.Event],
) *Application {
	// Validate and set defaults
	config = ValidateConfig(config)
	app := &Application{
		assistants:        assistants,
		conversations:     conversations,
		readConversations: readConversations,
		readMessages:      readMessages,
		catalogReader:     catalogReader,
		vectorRepository:  vectorRepository,
		publisher:         publisher,
	}

	// Initialize processors
	app.appProcessors = appProcessors{
		llmProcessor:      config.LLMProcessor,
		speechProcessor:   config.SpeechProcessor,
		visionProcessor:   config.VisionProcessor,
		documentProcessor: config.DocumentProcessor,
		dataProcessor:     config.DataProcessor,
		promptProvider:    config.PromptProvider,
	}

	// Initialize command handlers
	app.appCommands = appCommands{
		CreateAdminAssistantHandler:         commands.NewCreateAdminAssistantHandler(assistants, publisher, config.PromptProvider),
		CreateBusinessAssistantHandler:      commands.NewCreateBusinessAssistantHandler(assistants, publisher, config.PromptProvider),
		CreateSupportAssistantHandler:       commands.NewCreateSupportAssistantHandler(assistants, publisher, config.PromptProvider),
		CreateSchedulerAssistantHandler:     commands.NewCreateSchedulerAssistantHandler(assistants, publisher, config.PromptProvider),
		CreateUserAssistantHandler:          commands.NewCreateUserAssistantHandler(assistants, publisher, config.PromptProvider),
		ActivateAssistantHandler:            commands.NewActivateAssistantHandler(assistants, publisher),
		DeactivateAssistantHandler:          commands.NewDeactivateAssistantHandler(assistants, publisher),
		UpdateAssistantConfigurationHandler: commands.NewUpdateAssistantConfigurationHandler(assistants, publisher),
		ProcessUserInputHandler:             commands.NewProcessUserInputHandler(assistants, config.LLMProcessor, publisher, config.PromptProvider),
		ProcessSpeechInputHandler:           commands.NewProcessSpeechInputHandler(assistants, config.SpeechProcessor, config.LLMProcessor, publisher),
		ProcessImageInputHandler:            commands.NewProcessImageInputHandler(assistants, config.VisionProcessor, config.LLMProcessor, publisher),
		ProcessDocumentInputHandler:         commands.NewProcessDocumentInputHandler(assistants, config.DocumentProcessor, config.DataProcessor, config.LLMProcessor, publisher),
		CreateConversationHandler:           commands.NewCreateConversationHandler(conversations, publisher),
		AddMessageToConversationHandler:     commands.NewAddMessageToConversationHandler(conversations, readMessages, assistants, config.LLMProcessor, publisher),
		ChatWithConversationHandler:         commands.NewChatWithConversationHandler(conversations, readConversations, readMessages, assistants, config.LLMProcessor, publisher),
		UpdateConversationContextHandler:    commands.NewUpdateConversationContextHandler(conversations, readConversations, publisher),
	}

	// Initialize query handlers
	app.appQueries = appQueries{
		GetAssistantHandler:             queries.NewGetAssistantHandler(catalogReader),
		GetAssistantsHandler:            queries.NewGetAssistantsHandler(catalogReader),
		GetConversationHandler:          queries.NewGetConversationHandler(readConversations),
		GetUserConversationsHandler:     queries.NewGetUserConversationsHandler(readConversations),
		GetConversationMessagesHandler:  queries.NewGetConversationMessagesHandler(readMessages),
		GetConversationStatsHandler:     queries.NewGetConversationStatsHandler(conversations),
		GetOrCreateUserAssistantHandler: queries.NewGetOrCreateUserAssistantHandler(assistants, catalogReader, publisher),
	}

	// Initialize services
	app.appServices = appServices{
		toolRegistry:         toolRegistry,
		repositoryTranslator: NewRepositoryTranslator(),
	}

	// Initialize tool executors
	app.appTools = appTools{
		toolExecutor:   NewSimplifiedToolExecutor(toolRegistry, config.ToolConfig),
		streamExecutor: NewSimplifiedStreamingExecutor(toolRegistry, config.StreamingConfig),
	}

	return app
}

// Service implementations

func (a *Application) Execute(ctx context.Context, query services.RepositoryQuery) (*services.RepositoryResponse, error) {
	if err := a.ValidateQuery(query); err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}
	
	// Execute through tool registry
	result, err := a.toolRegistry.ExecuteTool(ctx, string(query.Operation), query.Parameters)
	if err != nil {
		return nil, err
	}
	
	return &services.RepositoryResponse{
		Data: result,
		Metadata: services.ResponseMetadata{
			EntityType: query.EntityType,
			Operation:  query.Operation,
		},
	}, nil
}

func (a *Application) ValidateQuery(query services.RepositoryQuery) error {
	return a.repositoryTranslator.ValidateQuery(query)
}

func (a *Application) GetSupportedOperations(entityType models.EntityType) []services.OperationType {
	// Get operations from tool definitions
	toolDefs := a.toolRegistry.GetToolDefinitions()
	operations := make([]services.OperationType, 0)
	
	entityPrefix := string(entityType) + "_"
	for _, toolDef := range toolDefs {
		if toolDef.Function.Name[:len(entityPrefix)] == entityPrefix {
			opName := toolDef.Function.Name[len(entityPrefix):]
			operations = append(operations, services.OperationType(opName))
		}
	}
	
	return operations
}

func (a *Application) TranslateAIRequest(aiRequest map[string]interface{}) (*services.RepositoryQuery, error) {
	return a.repositoryTranslator.TranslateAIRequest(aiRequest)
}

// Tool implementations

func (a *Application) ExecuteTools(ctx context.Context, toolCalls []externalai.ToolCall, execCtx *tools.ToolExecutionContext) ([]tools.ToolExecutionResult, error) {
	return a.toolExecutor.ExecuteTools(ctx, toolCalls, execCtx)
}

func (a *Application) ExecuteToolsStream(ctx context.Context, toolCalls []externalai.ToolCall, execCtx *tools.ToolExecutionContext) (<-chan tools.ToolExecutionStream, error) {
	return a.streamExecutor.ExecuteToolsStream(ctx, toolCalls, execCtx)
}