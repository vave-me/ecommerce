package application

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	
	"middleman/managers/internal/application/commands"
	"middleman/managers/internal/application/queries"
	"middleman/managers/internal/application/services"
	"middleman/managers/internal/application/tools"
	"middleman/managers/internal/application/consciousness"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
	externalai "middleman/internal/ai"
	"middleman/internal/ddd"
)

// App interface defines the complete application contract
type App interface {
	Commands
	Queries
	Services
	Tools
	EventProcessor
}

// EventProcessor interface for handling platform events
type EventProcessor interface {
	ProcessPlatformEvent(ctx context.Context, event ddd.Event) error
}

// Commands interface defines all command operations
type Commands interface {
	// Manager management
	CreateManager(ctx context.Context, cmd commands.CreateManager) error
	CreateAdminManager(ctx context.Context, cmd commands.CreateAdminManager) error
	CreateBusinessManager(ctx context.Context, cmd commands.CreateBusinessManager) error
	CreateSupportManager(ctx context.Context, cmd commands.CreateSupportManager) error
	CreateUserManager(ctx context.Context, cmd commands.CreateUserManager) error
	CreateSchedulerManager(ctx context.Context, cmd commands.CreateSchedulerManager) error
	ActivateManager(ctx context.Context, cmd commands.ActivateManager) error
	DeactivateManager(ctx context.Context, cmd commands.DeactivateManager) error
	UpdateManagerConfiguration(ctx context.Context, cmd commands.UpdateManagerConfiguration) error
	EnsureConsciousnessManager(ctx context.Context, cmd commands.EnsureConsciousnessManager) error
	EnsureSchedulerManager(ctx context.Context, cmd commands.EnsureSchedulerManager) error
	
	// Input processing
	ProcessUserInput(ctx context.Context, cmd commands.ProcessUserInput) (*commands.ProcessUserInputResult, error)
	ProcessSpeechInput(ctx context.Context, cmd commands.ProcessSpeechInput) (*commands.ProcessSpeechInputResult, error)
	ProcessImageInput(ctx context.Context, cmd commands.ProcessImageInput) (*commands.ProcessImageInputResult, error)
	ProcessDocumentInput(ctx context.Context, cmd commands.ProcessDocumentInput) (*commands.ProcessDocumentInputResult, error)
	
	// Conversation management
	CreateConversation(ctx context.Context, cmd commands.CreateConversation) error
	AddManagerToConversation(ctx context.Context, cmd commands.AddManagerToConversation) (*commands.AddManagerToConversationResult, error)
	ChatWithConversation(ctx context.Context, cmd commands.ChatWithConversation) (*commands.ChatWithConversationResult, error)
	UpdateConversationContext(ctx context.Context, cmd commands.UpdateConversationContext) error
	UpdateConversation(ctx context.Context, cmd commands.UpdateConversation) error
	DeleteConversation(ctx context.Context, cmd commands.DeleteConversation) error
}

// Queries interface defines all query operations
type Queries interface {
	GetManager(ctx context.Context, query queries.GetManager) (*domain.CatalogManager, error)
	GetManagers(ctx context.Context, query queries.GetManagers) ([]*domain.CatalogManager, error)
	GetOrCreateUserManager(ctx context.Context, query queries.GetOrCreateUserManager) (*domain.Manager, error)
	GetConversation(ctx context.Context, query queries.GetConversation) (*domain.ReadConversation, error)
	GetUserConversations(ctx context.Context, query queries.GetUserConversations) (*queries.GetUserConversationsResult, error)
	GetConversationMessages(ctx context.Context, query queries.GetConversationMessages) ([]*domain.ReadMessage, error)
	GetConversationStats(ctx context.Context, query queries.GetConversationStats) (*domain.ConversationStats, error)
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
	
	// Repositories
	managers          domain.ManagerRepository
	conversations     domain.ConversationRepository
	readConversations domain.ReadConversationRepository
	readMessages      domain.ReadMessagesRepository
	catalogReader     domain.CatalogRepository
	vectorRepository  domain.VectorRepository
	publisher         ddd.EventPublisher[ddd.Event]
	
	// Consciousness modules
	memoryCore           *consciousness.MemoryCore
	patternDetector      *consciousness.PatternDetector
	learningProcessor    *consciousness.LearningProcessor
	decisionOrchestrator *consciousness.DecisionOrchestrator
	actionExecutor       *consciousness.ActionExecutor
	
	// AI infrastructure
	aiManager            externalai.AIClientManager
}

// appCommands contains all command handlers
type appCommands struct {
	// Manager handlers
	commands.CreateManagerHandler
	commands.CreateAdminManagerHandler
	commands.CreateBusinessManagerHandler
	commands.CreateSupportManagerHandler
	commands.CreateUserManagerHandler
	commands.CreateSchedulerManagerHandler
	commands.ActivateManagerHandler
	commands.DeactivateManagerHandler
	commands.UpdateManagerConfigurationHandler
	commands.EnsureConsciousnessManagerHandler
	commands.EnsureSchedulerManagerHandler
	
	// Input processing handlers
	commands.ProcessUserInputHandler
	commands.ProcessSpeechInputHandler
	commands.ProcessImageInputHandler
	commands.ProcessDocumentInputHandler
	
	// Conversation handlers
	commands.CreateConversationHandler
	commands.AddManagerToConversationHandler
	commands.ChatWithConversationHandler
	commands.UpdateConversationContextHandler
	commands.UpdateConversationHandler
	commands.DeleteConversationHandler
}

// appQueries contains all query handlers
type appQueries struct {
	queries.GetManagerHandler
	queries.GetManagersHandler
	queries.GetOrCreateUserManagerHandler
	queries.GetConversationHandler
	queries.GetUserConversationsHandler
	queries.GetConversationMessagesHandler
	queries.GetConversationStatsHandler
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
	toolServiceRegistry  *tools.ToolServiceRegistry
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
	managers domain.ManagerRepository,
	conversations domain.ConversationRepository,
	readConversations domain.ReadConversationRepository,
	readMessages domain.ReadMessagesRepository,
	catalogReader domain.CatalogRepository,
	vectorRepository domain.VectorRepository,
	toolServiceRegistry *tools.ToolServiceRegistry,
	config *Config,
	publisher ddd.EventPublisher[ddd.Event],
) (*Application, error) {
	// Validate and set defaults
	config = ValidateConfig(config)
	
	app := &Application{
		managers:          managers,
		conversations:     conversations,
		readConversations: readConversations,
		readMessages:      readMessages,
		catalogReader:     catalogReader,
		vectorRepository:  vectorRepository,
		publisher:         publisher,
	}
	
	// Initialize consciousness modules
	logger := zerolog.Logger{}
	app.memoryCore = consciousness.NewMemoryCore(vectorRepository, logger)
	app.patternDetector = consciousness.NewPatternDetector(
		consciousness.NewActivitySurgeDetector(),
		consciousness.NewAbandonmentRiskDetector(),
		consciousness.NewFraudRiskDetector(),
		consciousness.NewSupportCrisisDetector(),
		consciousness.NewUserSurgeDetector(),
		consciousness.NewCancellationSpikeDetector(),
	)
	app.patternDetector.SetMemoryCore(app.memoryCore)
	app.learningProcessor = consciousness.NewLearningProcessor(vectorRepository)
	app.decisionOrchestrator = consciousness.NewDecisionOrchestrator()
	app.actionExecutor = consciousness.NewActionExecutor(toolServiceRegistry)
	
	// Initialize AI infrastructure
	app.aiManager = config.AIClientManager
	if app.aiManager == nil {
		// Create default AI manager if not provided
		factory := externalai.NewClientFactory()
		app.aiManager = externalai.NewClientManager(factory)
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
	processUserInputHandler, err := commands.NewProcessUserInputHandler(managers, config.LLMProcessor, publisher, config.PromptProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create ProcessUserInputHandler: %w", err)
	}
	
	processSpeechInputHandler, err := commands.NewProcessSpeechInputHandler(managers, config.SpeechProcessor, config.LLMProcessor, publisher)
	if err != nil {
		return nil, fmt.Errorf("failed to create ProcessSpeechInputHandler: %w", err)
	}
	
	processImageInputHandler, err := commands.NewProcessImageInputHandler(managers, config.VisionProcessor, config.LLMProcessor, publisher)
	if err != nil {
		return nil, fmt.Errorf("failed to create ProcessImageInputHandler: %w", err)
	}
	
	processDocumentInputHandler, err := commands.NewProcessDocumentInputHandler(managers, config.DocumentProcessor, config.DataProcessor, config.LLMProcessor, publisher)
	if err != nil {
		return nil, fmt.Errorf("failed to create ProcessDocumentInputHandler: %w", err)
	}
	
	app.appCommands = appCommands{
		CreateManagerHandler:                commands.NewCreateManagerHandler(managers, publisher, config.PromptProvider),
		CreateAdminManagerHandler:           commands.NewCreateAdminManagerHandler(managers, publisher, config.PromptProvider),
		CreateBusinessManagerHandler:        commands.NewCreateBusinessManagerHandler(managers, publisher, config.PromptProvider),
		CreateSupportManagerHandler:         commands.NewCreateSupportManagerHandler(managers, publisher, config.PromptProvider),
		CreateUserManagerHandler:            commands.NewCreateUserManagerHandler(managers, publisher, config.PromptProvider),
		CreateSchedulerManagerHandler:       commands.NewCreateSchedulerManagerHandler(managers, publisher, config.PromptProvider),
		ActivateManagerHandler:              commands.NewActivateManagerHandler(managers, publisher),
		DeactivateManagerHandler:            commands.NewDeactivateManagerHandler(managers, publisher),
		UpdateManagerConfigurationHandler:   commands.NewUpdateManagerConfigurationHandler(managers, publisher),
		EnsureConsciousnessManagerHandler:   commands.NewEnsureConsciousnessManagerHandler(managers, publisher, config.PromptProvider),
		EnsureSchedulerManagerHandler:       commands.NewEnsureSchedulerManagerHandler(managers, publisher, config.PromptProvider),
		ProcessUserInputHandler:             processUserInputHandler,
		ProcessSpeechInputHandler:           processSpeechInputHandler,
		ProcessImageInputHandler:            processImageInputHandler,
		ProcessDocumentInputHandler:         processDocumentInputHandler,
		CreateConversationHandler:           commands.NewCreateConversationHandler(conversations, publisher),
		AddManagerToConversationHandler:     commands.NewAddManagerToConversationHandler(conversations, readMessages, managers, config.LLMProcessor, publisher),
		ChatWithConversationHandler:         commands.NewChatWithConversationHandler(conversations, readConversations, readMessages, managers, config.LLMProcessor, publisher),
		UpdateConversationContextHandler:    commands.NewUpdateConversationContextHandler(conversations, readConversations, publisher),
		UpdateConversationHandler:           commands.NewUpdateConversationHandler(conversations, publisher),
		DeleteConversationHandler:           commands.NewDeleteConversationHandler(conversations, publisher),
	}
	
	// Initialize query handlers
	app.appQueries = appQueries{
		GetManagerHandler:              queries.NewGetManagerHandler(catalogReader),
		GetManagersHandler:             queries.NewGetManagersHandler(catalogReader),
		GetOrCreateUserManagerHandler:  queries.NewGetOrCreateUserManagerHandler(managers, catalogReader, publisher),
		GetConversationHandler:         queries.NewGetConversationHandler(readConversations),
		GetUserConversationsHandler:    queries.NewGetUserConversationsHandler(readConversations),
		GetConversationMessagesHandler: queries.NewGetConversationMessagesHandler(readMessages),
		GetConversationStatsHandler:    queries.NewGetConversationStatsHandler(conversations),
	}
	
	// Initialize services
	app.appServices = appServices{
		toolServiceRegistry:  toolServiceRegistry,
		repositoryTranslator: NewRepositoryTranslator(),
	}
	
	// Initialize tool executors
	app.appTools = appTools{
		toolExecutor:   NewSimplifiedToolExecutor(toolServiceRegistry, config.ToolConfig),
		streamExecutor: NewSimplifiedStreamingExecutor(toolServiceRegistry, config.StreamingConfig),
	}
	
	return app, nil
}

// GetMemoryCore returns the memory core component
func (a *Application) GetMemoryCore() *consciousness.MemoryCore {
	return a.memoryCore
}

// GetPatternDetector returns the pattern detector component
func (a *Application) GetPatternDetector() *consciousness.PatternDetector {
	return a.patternDetector
}

// GetLearningProcessor returns the learning processor component
func (a *Application) GetLearningProcessor() *consciousness.LearningProcessor {
	return a.learningProcessor
}

// GetDecisionOrchestrator returns the decision orchestrator component
func (a *Application) GetDecisionOrchestrator() *consciousness.DecisionOrchestrator {
	return a.decisionOrchestrator
}

// GetActionExecutor returns the action executor component
func (a *Application) GetActionExecutor() *consciousness.ActionExecutor {
	return a.actionExecutor
}

// GetAIManager returns the AI client manager
func (a *Application) GetAIManager() externalai.AIClientManager {
	return a.aiManager
}

// Service implementations

func (a *Application) Execute(ctx context.Context, query services.RepositoryQuery) (*services.RepositoryResponse, error) {
	if err := a.ValidateQuery(query); err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}
	
	// Execute through tool service registry
	result, err := a.toolServiceRegistry.ExecuteTool(ctx, string(query.Operation), query.Parameters)
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
	toolDefs := a.toolServiceRegistry.GetToolDefinitions()
	operations := make([]services.OperationType, 0)
	
	entityPrefix := string(entityType) + "_"
	for _, toolDef := range toolDefs {
		if len(toolDef.Function.Name) > len(entityPrefix) && toolDef.Function.Name[:len(entityPrefix)] == entityPrefix {
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

// ProcessPlatformEvent processes incoming platform events with production-ready features
func (a *Application) ProcessPlatformEvent(ctx context.Context, event ddd.Event) error {
	logger := zerolog.Ctx(ctx)
	span := trace.SpanFromContext(ctx)
	
	// Add tracing
	span.SetAttributes(
		attribute.String("event.type", event.EventName()),
		attribute.String("event.id", event.ID()),
	)
	
	// Measure processing time
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		span.SetAttributes(attribute.Int64("processing.duration_ms", duration.Milliseconds()))
		
		// Log slow processing
		if duration > 100*time.Millisecond {
			logger.Warn().
				Dur("duration", duration).
				Str("event_type", event.EventName()).
				Msg("Slow event processing detected")
		}
	}()
	
	// Store event in memory with retry
	var storeErr error
	for i := 0; i < 3; i++ {
		if storeErr = a.memoryCore.StoreEvent(ctx, event); storeErr == nil {
			break
		}
		if i < 2 {
			time.Sleep(time.Duration(i+1) * 10 * time.Millisecond)
		}
	}
	if storeErr != nil {
		span.RecordError(storeErr)
		logger.Warn().
			Err(storeErr).
			Int("retry_count", 3).
			Msg("Failed to store event in memory after retries")
	}
	
	// Detect patterns with timeout
	patternCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	
	pattern := a.patternDetector.DetectPattern(patternCtx, event)
	if pattern == nil {
		span.AddEvent("No pattern detected")
		return nil
	}
	
	// Record pattern detection
	span.SetAttributes(
		attribute.String("pattern.type", pattern.Type),
		attribute.Float64("pattern.confidence", pattern.Confidence),
	)
	
	logger.Info().
		Str("pattern_type", pattern.Type).
		Float64("confidence", pattern.Confidence).
		Interface("pattern_data", pattern.Data).
		Msg("Pattern detected from event")
	
	// Process learning insights with error handling
	insightCtx, insightCancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer insightCancel()
	
	insight, err := a.learningProcessor.ProcessPattern(insightCtx, pattern)
	if err != nil {
		span.RecordError(err)
		logger.Warn().
			Err(err).
			Str("pattern_type", pattern.Type).
			Msg("Failed to process learning from pattern")
	}
	
	// Make decision with validation
	decision := a.decisionOrchestrator.MakeDecision(ctx, pattern, insight)
	if decision == nil {
		span.AddEvent("No decision made")
		return nil
	}
	
	span.SetAttributes(
		attribute.String("decision.type", decision.Type),
		attribute.String("decision.priority", decision.Priority),
		attribute.Int("decision.actions_count", len(decision.Actions)),
	)
	
	// Only process high priority decisions
	if decision.Priority != "high" && decision.Priority != "critical" {
		logger.Debug().
			Str("priority", decision.Priority).
			Str("decision_type", decision.Type).
			Msg("Decision priority too low for immediate action")
		return nil
	}
	
	logger.Info().
		Str("decision_type", decision.Type).
		Str("priority", decision.Priority).
		Int("actions_count", len(decision.Actions)).
		Msg("High priority decision made")
	
	// Ensure consciousness manager exists with proper error handling
	if err := a.EnsureConsciousnessManager(ctx, commands.EnsureConsciousnessManager{
		ManagerID: "consciousness-manager",
	}); err != nil {
		span.RecordError(err)
		logger.Error().
			Err(err).
			Msg("Failed to ensure consciousness manager")
		return fmt.Errorf("failed to ensure consciousness manager: %w", err)
	}
	
	// Execute actions with proper error tracking
	successCount := 0
	failureCount := 0
	
	for i, action := range decision.Actions {
		actionSpan := trace.SpanFromContext(ctx)
		actionSpan.AddEvent("Executing action", trace.WithAttributes(
			attribute.Int("action.index", i),
			attribute.String("action.type", action.Type),
		))
		
		// Execute with timeout
		actionCtx, actionCancel := context.WithTimeout(ctx, 30*time.Second)
		err := a.actionExecutor.ExecuteAction(actionCtx, action)
		actionCancel()
		
		if err != nil {
			failureCount++
			actionSpan.RecordError(err)
			logger.Error().
				Err(err).
				Str("action_type", action.Type).
				Int("action_index", i).
				Interface("action_params", action.Parameters).
				Msg("Failed to execute consciousness action")
		} else {
			successCount++
			logger.Info().
				Str("action_type", action.Type).
				Int("action_index", i).
				Msg("Successfully executed consciousness action")
		}
	}
	
	// Record execution summary
	span.SetAttributes(
		attribute.Int("actions.success_count", successCount),
		attribute.Int("actions.failure_count", failureCount),
	)
	
	// Log execution summary
	logger.Info().
		Int("total_actions", len(decision.Actions)).
		Int("successful_actions", successCount).
		Int("failed_actions", failureCount).
		Str("decision_type", decision.Type).
		Msg("Completed decision execution")
	
	return nil
}