package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"middleman/internal/auth"
	"middleman/managers/internal/application"
	"middleman/managers/internal/application/commands"
	"middleman/managers/internal/application/processor"
	"middleman/managers/internal/application/queries"
	"middleman/managers/internal/domain"
	"middleman/managers/managerspb"

	"google.golang.org/grpc"
)

type server struct {
	app application.App

	managerspb.ManagersServiceServer
}

var _ managerspb.ManagersServiceServer = (*server)(nil)

func RegisterServer(app application.App, sr, registrar grpc.ServiceRegistrar) error {
	managerspb.RegisterManagersServiceServer(registrar, server{
		app: app,
	})
	return nil
}

// CreateManager creates a new manager
func (s server) CreateManager(ctx context.Context, request *managerspb.CreateManagerRequest) (*managerspb.CreateManagerResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Only admin users can create managers
	if claims.Role != "admin" {
		return nil, status.Error(grpc_code.PermissionDenied, "admin role required to create managers")
	}

	// Create the manager - always use defaults from domain
	managerID := uuid.New().String()
	err := s.app.CreateManager(ctx, commands.CreateManager{
		ID:           managerID,
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		UserID:       request.GetUserId(),
		Type:         "", // Empty string will use default in domain
		Capabilities: protoToDomainCapabilities(request.GetCapabilities()),
		Temperature:  request.GetTemperature(),
		MaxTokens:    int(request.GetMaxTokens()),
		SystemPrompt: request.GetSystemPrompt(),
	})
	if err != nil {
		return nil, err
	}

	return &managerspb.CreateManagerResponse{
		ManagerId: managerID,
	}, nil
}

// GetManager retrieves an manager by ID
func (s server) GetManager(ctx context.Context, request *managerspb.GetManagerRequest) (*managerspb.GetManagerResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Log claims for debugging admin role population
	log.Printf("[GetManager] Claims extracted - UserID: %s, Role: %s", claims.Subject, claims.Role)

	manager, err := s.app.GetManager(ctx, queries.GetManager{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &managerspb.GetManagerResponse{
		Manager: s.managerFromDomain(manager),
	}, nil
}

// ActivateManager activates an manager
func (s server) ActivateManager(ctx context.Context, request *managerspb.ActivateManagerRequest) (*managerspb.ActivateManagerResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Log claims for debugging admin role population
	log.Printf("[ActivateManager] Claims extracted - UserID: %s, Role: %s", claims.Subject, claims.Role)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ManagerID", request.GetId()))

	err := s.app.ActivateManager(ctx, commands.ActivateManager{
		ID: request.GetId(),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &managerspb.ActivateManagerResponse{}, nil
}

// DeactivateManager deactivates an manager
func (s server) DeactivateManager(ctx context.Context, request *managerspb.DeactivateManagerRequest) (*managerspb.DeactivateManagerResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Log claims for debugging admin role population
	log.Printf("[DeactivateManager] Claims extracted - UserID: %s, Role: %s", claims.Subject, claims.Role)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ManagerID", request.GetId()))

	err := s.app.DeactivateManager(ctx, commands.DeactivateManager{
		ID: request.GetId(),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &managerspb.DeactivateManagerResponse{}, nil
}

// UpdateManagerConfiguration updates manager configuration
func (s server) UpdateManagerConfiguration(ctx context.Context, request *managerspb.UpdateManagerConfigurationRequest) (*managerspb.UpdateManagerConfigurationResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[UpdateManagerConfiguration] Claims extracted - UserID: %s, Role: %s", userID, claims.Role)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ManagerID", request.GetId()))

	// Build configuration from structured fields
	config := commands.UpdateManagerConfiguration{
		ID:     request.GetId(),
		UserID: userID,
	}

	// Only update fields that are provided (non-zero values)
	if request.GetName() != "" {
		config.Name = request.GetName()
	}

	if request.GetDescription() != "" {
		config.Description = request.GetDescription()
	}

	if request.GetTemperature() != 0 {
		config.Temperature = request.GetTemperature()
	}

	if request.GetMaxTokens() != 0 {
		config.MaxTokens = int(request.GetMaxTokens())
	}

	if request.GetSystemPrompt() != "" {
		config.SystemPrompt = request.GetSystemPrompt()
	} else if request.GetSystemPrompt() == "" && len(request.GetCapabilities()) > 0 {
		// If capabilities are being updated but no prompt provided, use consolidated prompt
		config.SystemPrompt = processor.NewEnhancedLLMInterface().GetCompleteSystemPrompt()
		log.Printf("[UpdateManagerConfiguration] Using consolidated prompt for manager %s due to capability update", request.GetId())
	}

	if len(request.GetCapabilities()) > 0 {
		// Convert protobuf capabilities to domain capabilities with deduplication
		config.Capabilities = protoToDomainCapabilities(request.GetCapabilities())
	}

	err := s.app.UpdateManagerConfiguration(ctx, config)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &managerspb.UpdateManagerConfigurationResponse{}, nil
}

// GetManagers retrieves all managers
func (s server) GetManagers(ctx context.Context, request *managerspb.GetManagersRequest) (*managerspb.GetManagersResponse, error) {
	log.Printf("[GetManagers] START - Request received")

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		log.Printf("[GetManagers] ERROR - No authentication claims found")
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[GetManagers] Claims extracted - UserID: %s, Role: %s, Limit: %d, Page: %d",
		userID, claims.Role, request.GetLimit(), request.GetPage())

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("UserID", userID),
		attribute.Int("RequestLimit", int(request.GetLimit())),
	)

	// Use application layer to get managers with proper consistency
	log.Printf("[GetManagers] Calling app.GetManagers with UserID=%s, Limit=%d, Offset=%d",
		userID, int(request.GetLimit()), int(request.GetPage())*int(request.GetLimit()))

	managers, err := s.app.GetManagers(ctx, queries.GetManagers{
		UserID: userID,
		Limit:  int(request.GetLimit()),
		Offset: int(request.GetPage()) * int(request.GetLimit()),
	})

	if err != nil {
		log.Printf("[GetManagers] ERROR from app.GetManagers: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Printf("[GetManagers] app.GetManagers returned %d managers", len(managers))

	// Convert domain managers to protobuf
	protoManagers := make([]*managerspb.Manager, len(managers))
	for i, manager := range managers {
		protoManagers[i] = s.catalogManagerFromDomain(manager)
	}

	// If no managers found, return empty list
	// The consciousness manager is created separately on startup
	if len(protoManagers) == 0 {
		log.Printf("[GetManagers] No managers found for user %s, returning empty list", userID)
		return &managerspb.GetManagersResponse{
			Managers:   []*managerspb.Manager{},
			TotalCount: 0,
		}, nil
	}

	log.Printf("[GetManagers] SUCCESS - Returning %d managers", len(protoManagers))
	return &managerspb.GetManagersResponse{
		Managers:   protoManagers,
		TotalCount: int32(len(protoManagers)),
	}, nil
}

// convertStringMapToInterface converts string map to interface map
func convertStringMapToInterface(context map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range context {
		result[k] = v
	}
	return result
}

func (s server) managerFromDomain(manager *domain.CatalogManager) *managerspb.Manager {
	// Convert domain capabilities to protobuf capabilities with deduplication
	capabilities := domainToProtoCapabilities(manager.Capabilities)

	// Convert domain type to proto type
	var protoType managerspb.ManagerType
	switch manager.Type {
	case domain.ManagerTypeAdmin:
		protoType = managerspb.ManagerType_ADMIN
	case domain.ManagerTypeBusiness:
		protoType = managerspb.ManagerType_BUSINESS
	case domain.ManagerTypeSupport:
		protoType = managerspb.ManagerType_SUPPORT
	case domain.ManagerTypeScheduler:
		protoType = managerspb.ManagerType_SCHEDULER
	default:
		protoType = managerspb.ManagerType_STANDARD
	}

	// Create manager with simple capability flags for frontend
	pbManager := &managerspb.Manager{
		Id:           manager.ID,
		Name:         manager.Name,
		Description:  manager.Description,
		Type:         protoType,
		Capabilities: capabilities,
		Active:       manager.Active,
		Temperature:  manager.Temperature,
		MaxTokens:    int32(manager.MaxTokens),
		SystemPrompt: manager.SystemPrompt,
		CreatedAt:    manager.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    manager.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return pbManager
}

func (s server) catalogManagerFromDomain(manager *domain.CatalogManager) *managerspb.Manager {
	// Convert domain capabilities to protobuf capabilities with deduplication
	capabilities := domainToProtoCapabilities(manager.Capabilities)

	return &managerspb.Manager{
		Id:           manager.ID,
		Name:         manager.Name,
		Description:  manager.Description,
		Capabilities: capabilities,
		Active:       manager.Active,
		Temperature:  manager.Temperature,
		MaxTokens:    int32(manager.MaxTokens),
		SystemPrompt: manager.SystemPrompt,
		CreatedAt:    manager.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    manager.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Conversation Management Endpoints

// CreateConversation creates a new conversation
func (s server) CreateConversation(ctx context.Context, request *managerspb.CreateConversationRequest) (*managerspb.CreateConversationResponse, error) {
	span := trace.SpanFromContext(ctx)
	conversationID := uuid.New().String()

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		log.Printf("[CREATE_CONVERSATION] ERROR: Authentication failed - no JWT claims found")
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[CreateConversation] Claims extracted - UserID: %s, Role: %s", userID, claims.Role)

	// Creating new conversation
	log.Printf("[CREATE_CONVERSATION] Request ManagerID: %s", request.GetManagerId())

	span.SetAttributes(
		attribute.String("ConversationID", conversationID),
		attribute.String("UserID", userID),
		attribute.String("ManagerID", request.GetManagerId()),
	)

	// Convert protobuf context to map[string]interface{}
	initialContext := make(map[string]interface{})
	for k, v := range request.GetInitialContext() {
		initialContext[k] = v
	}
	// Context converted

	// If manager ID is not provided, generate a new one
	managerID := request.GetManagerId()
	if managerID == "" {
		managerID = uuid.New().String()
		log.Printf("[CREATE_CONVERSATION] No ManagerID provided, generated new one: %s", managerID)
	}

	// Create conversation via application layer
	log.Printf("[CREATE_CONVERSATION] Creating conversation with: ID=%s, UserID=%s, ManagerID=%s",
		conversationID, userID, managerID)

	err := s.app.CreateConversation(ctx, commands.CreateConversation{
		ID:             conversationID,
		UserID:         userID,
		ManagerID:      managerID,
		InitialContext: initialContext,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Conversation created successfully

	return &managerspb.CreateConversationResponse{
		ConversationId: conversationID,
	}, nil
}

// GetConversation retrieves a conversation by ID
func (s server) GetConversation(ctx context.Context, request *managerspb.GetConversationRequest) (*managerspb.GetConversationResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[GetConversation] Claims extracted - UserID: %s, Role: %s", userID, claims.Role)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConversationID", request.GetConversationId()),
		attribute.String("UserID", userID),
	)

	conversation, err := s.app.GetConversation(ctx, queries.GetConversation{
		ConversationID: request.GetConversationId(),
		UserID:         userID,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &managerspb.GetConversationResponse{
		Conversation: s.conversationViewToProto(conversation),
	}, nil
}

// GetUserConversations retrieves conversations for a user
func (s server) GetUserConversations(ctx context.Context, request *managerspb.GetUserConversationsRequest) (*managerspb.GetUserConversationsResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[GetUserConversations] Claims extracted - UserID: %s, Role: %s", userID, claims.Role)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("UserID", userID),
		attribute.Bool("ActiveOnly", request.GetActiveOnly()),
		attribute.Int("Page", int(request.GetPage())),
		attribute.Int("Limit", int(request.GetLimit())),
	)

	// Extract pagination parameters
	page := int(request.GetPage())
	limit := int(request.GetLimit())

	// Set defaults if not provided
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	// Calculate offset from page
	offset := (page - 1) * limit

	result, err := s.app.GetUserConversations(ctx, queries.GetUserConversations{
		UserID:     userID,
		ActiveOnly: request.GetActiveOnly(),
		Limit:      limit,
		Offset:     offset,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	conversations := make([]*managerspb.Conversation, len(result.Conversations))
	for i, conv := range result.Conversations {
		conversations[i] = s.conversationViewToProto(conv)
	}

	return &managerspb.GetUserConversationsResponse{
		Conversations: conversations,
		TotalCount:    int32(result.TotalCount),
	}, nil
}

// GetConversationMessages retrieves messages for a conversation
func (s server) GetConversationMessages(ctx context.Context, request *managerspb.GetConversationMessagesRequest) (*managerspb.GetConversationMessagesResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[GetConversationMessages] Claims extracted - UserID: %s, Role: %s", userID, claims.Role)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConversationID", request.GetConversationId()),
		attribute.String("UserID", userID),
		attribute.Int("Page", int(request.GetPage())),
		attribute.Int("Limit", int(request.GetLimit())),
	)

	// Extract pagination parameters
	page := int(request.GetPage())
	limit := int(request.GetLimit())

	// Set defaults if not provided
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50 // Default higher for messages
	}

	// Calculate offset from page
	offset := (page - 1) * limit

	messages, err := s.app.GetConversationMessages(ctx, queries.GetConversationMessages{
		ConversationID: request.GetConversationId(),
		UserID:         userID,
		Limit:          limit,
		Offset:         offset,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoMessages := make([]*managerspb.ConversationMessage, len(messages))
	for i, msg := range messages {
		protoMessages[i] = s.conversationMessageToProto(msg)
	}

	return &managerspb.GetConversationMessagesResponse{
		Messages:   protoMessages,
		TotalCount: int32(len(messages)), // Note: This should be the actual total count from pagination
	}, nil
}

// AddMessageToConversation adds a message to an existing conversation
func (s server) AddMessageToConversation(ctx context.Context, request *managerspb.AddMessageToConversationRequest) (*managerspb.AddMessageToConversationResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	messageID := uuid.New().String()
	span.SetAttributes(
		attribute.String("ConversationID", request.GetConversationId()),
		attribute.String("UserID", userID),
		attribute.String("MessageID", messageID),
		attribute.String("Role", request.GetRole().String()),
	)

	// Convert protobuf metadata to map[string]interface{}
	metadata := make(map[string]interface{})
	for k, v := range request.GetMetadata() {
		metadata[k] = v
	}

	// Convert protobuf role to domain role
	var role domain.MessageRole
	switch request.GetRole() {
	case managerspb.MessageRole_USER:
		role = domain.UserRole
	case managerspb.MessageRole_ASSISTANT:
		role = domain.ManagerRole
	case managerspb.MessageRole_SYSTEM:
		role = domain.SystemRole
	default:
		role = domain.UserRole
	}

	// For user messages, default to processing with LLM
	processWithLLM := role == domain.UserRole

	result, err := s.app.AddMessageToConversation(ctx, commands.AddMessageToConversation{
		ConversationID:     request.GetConversationId(),
		MessageID:          messageID,
		ManagerID:          request.GetManagerId(),
		Role:               role,
		Content:            request.GetContent(),
		Metadata:           metadata,
		UserID:             userID,
		ProcessWithLLM:     processWithLLM,
		MaxHistoryMessages: 20, // Default to last 20 messages
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Build response with LLM results if processed
	response := &managerspb.AddMessageToConversationResponse{
		MessageId: result.MessageID,
	}

	// If LLM processing occurred, include the manager's response
	if result.ManagerMessageID != "" {
		// For backward compatibility, still return the user's message ID
		// But also log that manager responded
		log.Printf("[AddMessageToConversation] LLM processed - User msg: %s, Manager msg: %s",
			result.MessageID, result.ManagerMessageID)
	}

	return response, nil
}

// ChatWithConversation processes a message within an existing conversation context
func (s server) ChatWithConversation(ctx context.Context, request *managerspb.ChatWithConversationRequest) (*managerspb.ChatWithConversationResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[ChatWithConversation] Claims extracted - UserID: %s, Role: %s", userID, claims.Role)

	log.Printf("[GRPC_CHAT] ========== ChatWithConversation START ==========")
	log.Printf("[GRPC_CHAT] Request received - ConversationID: %s, ManagerID: %s, UserID: %s, Message: %s",
		request.GetConversationId(), request.GetManagerId(), userID, request.GetMessage())
	log.Printf("[GRPC_CHAT] MaxHistoryMessages: %d, Context keys: %d",
		request.GetMaxHistoryMessages(), len(request.GetContext()))

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConversationID", request.GetConversationId()),
		attribute.String("ManagerID", request.GetManagerId()),
		attribute.String("UserID", userID),
		attribute.Int("MaxHistoryMessages", int(request.GetMaxHistoryMessages())),
	)

	// Get manager ID - required field
	managerID := request.GetManagerId()
	if managerID == "" {
		// Use the consciousness manager ID
		managerID = "consciousness-manager-001"
		log.Printf("[GRPC_CHAT] No manager ID provided, using consciousness manager: %s", managerID)
	}

	// Handle conversation creation if needed
	conversationID := request.GetConversationId()
	if conversationID == "" {
		// Create new conversation if conversation_id is empty
		log.Printf("[GRPC_CHAT] No conversation ID provided, creating new conversation")
		newConversationID := uuid.New().String()

		// Convert protobuf context to map[string]interface{} for initial context
		initialContext := make(map[string]interface{})
		for k, v := range request.GetContext() {
			initialContext[k] = v
		}

		// Create conversation
		err := s.app.CreateConversation(ctx, commands.CreateConversation{
			ID:             newConversationID,
			UserID:         userID,
			ManagerID:      managerID,
			InitialContext: initialContext,
		})
		if err != nil {
			log.Printf("[GRPC_CHAT] ERROR: Failed to create conversation: %v", err)
			return nil, status.Errorf(grpc_code.Internal, "failed to create conversation: %v", err)
		}
		conversationID = newConversationID
		log.Printf("[GRPC_CHAT] Created new conversation: %s", conversationID)
	}

	// Convert protobuf context to map[string]interface{}
	context := make(map[string]interface{})
	for k, v := range request.GetContext() {
		context[k] = v
	}
	log.Printf("[GRPC_CHAT] Context conversion completed, %d keys converted", len(context))

	log.Printf("[GRPC_CHAT] Calling application layer ChatWithConversation with userID: %s, conversationID: %s, managerID: %s...",
		userID, conversationID, managerID)
	result, err := s.app.ChatWithConversation(ctx, commands.ChatWithConversation{
		ConversationID:     conversationID,
		ManagerID:          managerID,
		UserID:             userID,
		Message:            request.GetMessage(),
		Context:            context,
		MaxHistoryMessages: int(request.GetMaxHistoryMessages()),
	})

	if err != nil {
		log.Printf("[GRPC_CHAT] ERROR: Application layer returned error: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Printf("[GRPC_CHAT] SUCCESS: Application layer returned result!")
	log.Printf("[GRPC_CHAT] Result details - Response length: %d, Actions: %d, Confidence: %.2f, Status: %s, MessageID: %s",
		len(result.Response), len(result.Actions), result.Confidence, result.Status, result.MessageID)

	if result.Data != nil {
		log.Printf("[GRPC_CHAT] Result data keys: %d", len(result.Data))
		for k, v := range result.Data {
			log.Printf("[GRPC_CHAT] Data[%s] = %v", k, v)
		}
	}

	log.Printf("[GRPC_CHAT] Converting %d domain actions to protobuf actions...", len(result.Actions))
	// Convert domain actions to protobuf actions
	actions := make([]*managerspb.ManagerAction, len(result.Actions))
	for i, action := range result.Actions {
		log.Printf("[GRPC_CHAT] Converting action %d: Type=%s, Endpoint=%s", i, action.Type, action.Endpoint)
		actions[i] = &managerspb.ManagerAction{
			Type:        action.Type,
			Endpoint:    action.Endpoint,
			Method:      action.Method,
			Parameters:  s.convertParametersToStringMap(action.Parameters),
			Description: action.Description,
		}
	}
	log.Printf("[GRPC_CHAT] Actions conversion completed successfully")

	log.Printf("[GRPC_CHAT] Converting result data to string map...")
	// Convert result data to string map
	data := make(map[string]string)
	if result.Data != nil {
		for k, v := range result.Data {
			var strValue string
			// Check if value needs JSON serialization
			switch val := v.(type) {
			case string:
				strValue = val
			case []byte:
				strValue = string(val)
			case int, int32, int64, float32, float64, bool:
				strValue = fmt.Sprintf("%v", val)
			default:
				// For complex types (structs, slices, maps), use JSON encoding
				jsonBytes, err := json.Marshal(val)
				if err != nil {
					log.Printf("[GRPC_CHAT] WARNING: Failed to marshal %s to JSON: %v", k, err)
					strValue = fmt.Sprintf("%v", val)
				} else {
					strValue = string(jsonBytes)
				}
			}
			data[k] = strValue
			log.Printf("[GRPC_CHAT] Data conversion: %s = %s", k, strValue)
		}
	}
	log.Printf("[GRPC_CHAT] Data conversion completed, %d keys converted", len(data))

	log.Printf("[GRPC_CHAT] Constructing final gRPC response...")
	response := &managerspb.ChatWithConversationResponse{
		Response:   result.Response,
		Actions:    actions,
		Confidence: result.Confidence,
		Status:     result.Status,
		Data:       data,
		MessageId:  result.MessageID,
	}

	log.Printf("[GRPC_CHAT] Final response created - Response length: %d, Actions: %d, Status: %s",
		len(response.Response), len(response.Actions), response.Status)
	log.Printf("[GRPC_CHAT] ========== ChatWithConversation SUCCESS ==========")

	return response, nil
}

// UpdateConversationContext updates the context of a conversation
func (s server) UpdateConversationContext(ctx context.Context, request *managerspb.UpdateConversationContextRequest) (*managerspb.UpdateConversationContextResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[UpdateConversationContext] Claims extracted - UserID: %s, Role: %s", userID, claims.Role)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConversationID", request.GetConversationId()),
		attribute.String("UserID", userID),
	)

	// Convert protobuf context to map[string]interface{}
	context := make(map[string]interface{})
	for k, v := range request.GetContext() {
		context[k] = v
	}

	err := s.app.UpdateConversationContext(ctx, commands.UpdateConversationContext{
		ConversationID: request.GetConversationId(),
		UserID:         userID,
		Context:        context,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &managerspb.UpdateConversationContextResponse{}, nil
}

// UpdateConversation updates conversation metadata (Frontend-expected endpoint)
func (s server) UpdateConversation(ctx context.Context, request *managerspb.UpdateConversationRequest) (*managerspb.UpdateConversationResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[UpdateConversation] Claims extracted - UserID: %s, Role: %s", userID, claims.Role)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConversationID", request.GetConversationId()),
		attribute.String("UserID", userID),
	)

	// Convert protobuf metadata to map[string]interface{}
	metadata := make(map[string]interface{})
	for k, v := range request.GetMetadata() {
		metadata[k] = v
	}

	err := s.app.UpdateConversation(ctx, commands.UpdateConversation{
		ConversationID: request.GetConversationId(),
		UserID:         userID,
		Metadata:       metadata,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Get updated conversation to return
	conversation, err := s.app.GetConversation(ctx, queries.GetConversation{
		ConversationID: request.GetConversationId(),
		UserID:         userID,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &managerspb.UpdateConversationResponse{
		Conversation: s.conversationViewToProto(conversation),
	}, nil
}

// DeleteConversation deletes a conversation (Frontend-expected endpoint)
func (s server) DeleteConversation(ctx context.Context, request *managerspb.DeleteConversationRequest) (*managerspb.DeleteConversationResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[DeleteConversation] Claims extracted - UserID: %s, Role: %s", userID, claims.Role)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConversationID", request.GetConversationId()),
		attribute.String("UserID", userID),
	)

	err := s.app.DeleteConversation(ctx, commands.DeleteConversation{
		ConversationID: request.GetConversationId(),
		UserID:         userID,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &managerspb.DeleteConversationResponse{}, nil
}

// ArchiveConversation archives a conversation
func (s server) ArchiveConversation(ctx context.Context, request *managerspb.ArchiveConversationRequest) (*managerspb.ArchiveConversationResponse, error) {
	// This endpoint is not implemented in the application layer yet
	// It would require a new command handler
	return &managerspb.ArchiveConversationResponse{}, nil
}

func (s server) conversationViewToProto(conv *domain.ReadConversation) *managerspb.Conversation {
	// Convert context to string map
	context := make(map[string]string)
	if conv.Context != nil {
		for k, v := range conv.Context {
			context[k] = fmt.Sprintf("%v", v)
		}
	}

	return &managerspb.Conversation{
		Id:        conv.ID,
		UserId:    conv.UserID,
		ManagerId: conv.ManagerID,
		Messages:  []*managerspb.ConversationMessage{}, // Messages loaded separately
		CreatedAt: conv.CreatedAt.Format(time.RFC3339),
		UpdatedAt: conv.UpdatedAt.Format(time.RFC3339),
		Active:    conv.Active,
		Context:   context,
	}
}

func (s server) conversationMessageToProto(msg *domain.ReadMessage) *managerspb.ConversationMessage {

	// Convert metadata to string map
	metadata := make(map[string]string)
	if msg.Metadata != nil {
		for k, v := range msg.Metadata {
			metadata[k] = fmt.Sprintf("%v", v)
		}
	}

	// Convert domain role to protobuf role
	var role managerspb.MessageRole
	switch msg.Role {
	case domain.UserRole:
		role = managerspb.MessageRole_USER
	case domain.ManagerRole:
		role = managerspb.MessageRole_ASSISTANT
	case domain.SystemRole:
		role = managerspb.MessageRole_SYSTEM
	default:
		role = managerspb.MessageRole_USER
	}

	// Convert actions
	actions := make([]*managerspb.ManagerAction, len(msg.ActionsTaken))
	for i, action := range msg.ActionsTaken {
		actions[i] = &managerspb.ManagerAction{
			Type:        action.Type,
			Endpoint:    action.Endpoint,
			Method:      action.Method,
			Parameters:  s.convertParametersToStringMap(action.Parameters),
			Description: action.Description,
		}
	}

	return &managerspb.ConversationMessage{
		Id:           msg.ID,
		Role:         role,
		Content:      msg.Content,
		Timestamp:    msg.Timestamp.Format(time.RFC3339),
		Metadata:     metadata,
		ActionsTaken: actions,
	}
}

func (s server) convertParametersToStringMap(parameters map[string]interface{}) map[string]string {
	result := make(map[string]string)
	if parameters != nil {
		for k, v := range parameters {
			result[k] = fmt.Sprintf("%v", v)
		}
	}
	return result
}

// ProcessUserInput processes user input through an manager
func (s server) ProcessUserInput(ctx context.Context, request *managerspb.ProcessUserInputRequest) (*managerspb.ProcessUserInputResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ManagerID", request.GetManagerId()),
		attribute.String("UserID", userID),
		attribute.String("RequestType", request.GetRequestType()),
	)

	// Generate a unique request ID since the request doesn't have one
	requestID := fmt.Sprintf("req_%s_%d", userID, time.Now().UnixNano())

	// Convert protobuf context to interface map
	requestContext := convertStringMapToInterface(request.GetContext())

	// Add timeout context for AI processing to prevent 504 timeouts
	// Set a reasonable timeout for AI processing (5 minutes)
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	result, err := s.app.ProcessUserInput(timeoutCtx, commands.ProcessUserInput{
		ID:          requestID,
		ManagerID:   request.GetManagerId(),
		UserID:      userID,
		Message:     request.GetMessage(),
		Context:     requestContext,
		RequestType: request.GetRequestType(),
		Timestamp:   time.Now(),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		// Check if it's a timeout error and provide specific message
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return &managerspb.ProcessUserInputResponse{
				Id:         uuid.New().String(),
				RequestId:  requestID,
				Message:    "Your request is taking longer than expected to process. Please try a simpler request or try again later.",
				Data:       map[string]string{"error": "processing_timeout"},
				Actions:    []*managerspb.ManagerAction{},
				Timestamp:  time.Now().Format("2006-01-02T15:04:05Z07:00"),
				Status:     "timeout",
				Confidence: 0.1,
			}, nil
		}

		return nil, err
	}

	// Convert domain actions to protobuf actions
	var protoActions []*managerspb.ManagerAction
	for _, action := range result.ExecutedActions {
		protoActions = append(protoActions, &managerspb.ManagerAction{
			Type:        action.Type,
			Endpoint:    action.Endpoint,
			Method:      action.Method,
			Parameters:  s.convertParametersToStringMap(action.Parameters),
			Description: action.Description,
		})
	}

	return &managerspb.ProcessUserInputResponse{
		Id:         result.ResponseID,
		RequestId:  requestID,
		Message:    result.ResponseMessage,
		Data:       map[string]string{"status": result.ResponseStatus}, // Simple status data
		Actions:    protoActions,
		Timestamp:  result.ResponseTimestamp.Format("2006-01-02T15:04:05Z07:00"),
		Status:     result.ResponseStatus,
		Confidence: result.ResponseConfidence,
	}, nil
}

// ProcessImageInput handles image analysis requests with specialized model selection
func (s server) ProcessImageInput(ctx context.Context, req *managerspb.ProcessImageInputRequest) (*managerspb.ProcessImageInputResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	log.Printf("[ProcessImageInput] Starting image processing with model selection. ManagerID: %s, ImageFormat: %s, Analysis: %s",
		req.ManagerId, req.ImageFormat, req.AnalysisType)

	// Validate request
	if req.ManagerId == "" {
		return nil, status.Errorf(grpc_code.InvalidArgument, "manager_id is required")
	}
	if len(req.ImageData) == 0 && req.ImageUrl == "" {
		return nil, status.Errorf(grpc_code.InvalidArgument, "either image_data or image_url must be provided")
	}

	// Detect processing mode based on analysis type and user prompt
	processingMode := s.detectImageProcessingMode(req)

	log.Printf("[ProcessImageInput] Detected processing mode: %s for prompt: '%s'", processingMode, req.UserPrompt)

	// Create command with intelligent workflow detection
	cmd := commands.ProcessImageInput{
		ID:             uuid.New().String(),
		ManagerID:      req.ManagerId,
		UserID:         userID,
		ImageData:      req.ImageData,
		ImageURL:       req.ImageUrl,
		ImageFormat:    req.ImageFormat,
		AnalysisType:   req.AnalysisType,
		UserPrompt:     req.UserPrompt,
		Context:        convertStringMapToInterface(req.Context),
		Timestamp:      time.Now(),
		RequestType:    "image_processing",
		ProcessingMode: processingMode,
		ListingData:    s.extractListingDataFromContext(req.Context),
	}

	// Execute with enhanced model selection and workflow detection
	result, err := s.app.ProcessImageInput(ctx, cmd)
	if err != nil {
		log.Printf("[ProcessImageInput] Error processing image: %v", err)
		return nil, status.Errorf(grpc_code.Internal, "failed to process image: %v", err)
	}

	log.Printf("[ProcessImageInput] Image processing completed. Mode: %s, Confidence: %.2f, Processing time: %v, Listing ID: %s",
		result.ProcessingMode, result.ResponseConfidence, result.ProcessingTime, result.CreatedListingID)

	// Convert result data to string map with enhanced metrics
	data := make(map[string]string)
	if result.ImageMetadata != nil {
		for k, v := range result.ImageMetadata {
			data[k] = fmt.Sprintf("%v", v)
		}
	}
	data["processing_time_ms"] = fmt.Sprintf("%.0f", float64(result.ProcessingTime.Milliseconds()))
	data["input_source"] = result.InputSource
	data["processing_mode"] = result.ProcessingMode
	data["image_attached"] = fmt.Sprintf("%t", result.ImageAttached)
	if result.CreatedListingID != "" {
		data["created_listing_id"] = result.CreatedListingID
	}

	return &managerspb.ProcessImageInputResponse{
		Id:                 result.ResponseID,
		RequestId:          cmd.ID,
		AnalysisResult:     result.AnalysisResult,
		Message:            result.ResponseMessage,
		Data:               data,
		Actions:            s.convertDomainActionsToProto(result.ExecutedActions),
		Timestamp:          result.ResponseTimestamp.Format(time.RFC3339),
		Status:             result.ResponseStatus,
		Confidence:         result.ResponseConfidence,
		AnalysisConfidence: result.VisionConfidence,
	}, nil
}

// detectImageProcessingMode determines the appropriate processing workflow
func (s server) detectImageProcessingMode(req *managerspb.ProcessImageInputRequest) string {

	// Check analysis type for explicit intent
	if req.AnalysisType == "attach" || req.AnalysisType == "attachment" {
		return "attach_image"
	}

	// Check user prompt for intent keywords
	if req.UserPrompt != "" {
		prompt := strings.ToLower(req.UserPrompt)

		// Check for listing creation keywords (extended)
		listingKeywords := []string{
			"create listing", "create product", "make listing", "list this",
			"list item", "list", "sell", "sell this", "sell item", "add product",
			"post for sale", "post item", "put for sale", "put item for sale",
			"list for sale", "sell for", "€", "eur", "euro", "price",
			"marketplace listing", "product listing",
		}

		for _, keyword := range listingKeywords {
			if strings.Contains(prompt, keyword) {
				return "attach_image"
			}
		}

		// Check for analysis keywords (extended)
		analysisKeywords := []string{
			"analyze", "analysis", "describe", "what is", "what's in", "identify",
			"recognize", "detect", "examine", "ocr", "read text", "extract text",
			"count", "objects", "details",
		}

		for _, keyword := range analysisKeywords {
			if strings.Contains(prompt, keyword) {
				return "analyze_image"
			}
		}
	}

	// Check context for listing data
	if req.Context != nil {
		for key := range req.Context {
			if strings.Contains(strings.ToLower(key), "listing") ||
				strings.Contains(strings.ToLower(key), "product") ||
				strings.Contains(strings.ToLower(key), "price") {
				return "attach_image"
			}
		}
	}

	// Default to analysis mode
	return "analyze_image"
}

// extractListingDataFromContext extracts listing-related data from request context
func (s server) extractListingDataFromContext(context map[string]string) map[string]interface{} {
	listingData := make(map[string]interface{})

	if context == nil {
		return listingData
	}

	// Extract relevant listing fields
	listingFields := []string{
		"name", "title", "description", "price", "category", "brand",
		"condition", "location", "specifications", "features", "tags",
	}

	for _, field := range listingFields {
		if value, exists := context[field]; exists {
			listingData[field] = value
		}
		// Also check with "listing_" prefix
		if value, exists := context["listing_"+field]; exists {
			listingData[field] = value
		}
	}

	return listingData
}

// ProcessSpeechInput handles speech processing requests with specialized model selection
func (s server) ProcessSpeechInput(ctx context.Context, req *managerspb.ProcessSpeechInputRequest) (*managerspb.ProcessSpeechInputResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	log.Printf("[ProcessSpeechInput] Starting speech processing with model selection. ManagerID: %s, AudioFormat: %s",
		req.ManagerId, req.AudioFormat)

	// Validate request
	if req.ManagerId == "" {
		return nil, status.Errorf(grpc_code.InvalidArgument, "manager_id is required")
	}
	if len(req.AudioData) == 0 {
		return nil, status.Errorf(grpc_code.InvalidArgument, "audio_data is required")
	}

	// Create command with proper input type detection
	cmd := commands.ProcessSpeechInput{
		ID:          uuid.New().String(),
		ManagerID:   req.ManagerId,
		UserID:      userID,
		AudioData:   req.AudioData,
		AudioFormat: req.AudioFormat,
		Language:    req.Language,
		Context:     convertStringMapToInterface(req.Context),
		Timestamp:   time.Now(),
		RequestType: "speech_processing",
	}

	// Execute with enhanced model selection
	result, err := s.app.ProcessSpeechInput(ctx, cmd)
	if err != nil {
		log.Printf("[ProcessSpeechInput] Error processing speech: %v", err)
		return nil, status.Errorf(grpc_code.Internal, "failed to process speech: %v", err)
	}

	log.Printf("[ProcessSpeechInput] Speech processing completed successfully. Transcribed: '%s', Confidence: %.2f",
		result.TranscribedText, result.ResponseConfidence)

	// Convert result data to string map
	data := map[string]string{
		"status":         result.ResponseStatus,
		"audio_duration": fmt.Sprintf("%.2fs", result.AudioDuration.Seconds()),
		"stt_confidence": fmt.Sprintf("%.2f", result.STTConfidence),
		"llm_confidence": fmt.Sprintf("%.2f", result.LLMConfidence),
	}

	return &managerspb.ProcessSpeechInputResponse{
		Id:              result.ResponseID,
		RequestId:       cmd.ID,
		TranscribedText: result.TranscribedText,
		Message:         result.ResponseMessage,
		Data:            data,
		Actions:         s.convertDomainActionsToProto(result.ExecutedActions),
		Timestamp:       result.ResponseTimestamp.Format(time.RFC3339),
		Status:          result.ResponseStatus,
		Confidence:      result.ResponseConfidence,

		TranscriptionConfidence: result.STTConfidence,
	}, nil
}

// ProcessDocumentInput processes document input for analysis and tool execution
func (s server) ProcessDocumentInput(ctx context.Context, request *managerspb.ProcessDocumentInputRequest) (*managerspb.ProcessDocumentInputResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ManagerID", request.GetManagerId()),
		attribute.String("UserID", userID),
		attribute.String("DocumentFormat", request.GetDocumentFormat()),
		attribute.String("AnalysisType", request.GetAnalysisType()),
		attribute.String("RequestType", request.GetRequestType()),
	)

	// Generate a unique request ID
	requestID := fmt.Sprintf("doc_req_%s_%d", userID, time.Now().UnixNano())

	// Convert protobuf context to interface map
	context := convertStringMapToInterface(request.GetContext())

	result, err := s.app.ProcessDocumentInput(ctx, commands.ProcessDocumentInput{
		ID:             requestID,
		ManagerID:      request.GetManagerId(),
		UserID:         userID,
		DocumentData:   request.GetDocumentData(),
		DocumentURL:    request.GetDocumentUrl(),
		DocumentFormat: request.GetDocumentFormat(),
		AnalysisType:   request.GetAnalysisType(),
		UserPrompt:     request.GetUserPrompt(),
		Context:        context,
		RequestType:    request.GetRequestType(),
		Timestamp:      time.Now(),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Convert domain actions to protobuf actions
	var protoActions []*managerspb.ManagerAction
	for _, action := range result.ExecutedActions {
		protoActions = append(protoActions, &managerspb.ManagerAction{
			Type:        action.Type,
			Endpoint:    action.Endpoint,
			Method:      action.Method,
			Parameters:  s.convertParametersToStringMap(action.Parameters),
			Description: action.Description,
		})
	}

	// Convert result data to string map with document-specific metrics
	data := map[string]string{
		"status":                result.ResponseStatus,
		"processing_confidence": fmt.Sprintf("%.2f", result.ProcessingConfidence),
		"llm_confidence":        fmt.Sprintf("%.2f", result.LLMConfidence),
		"document_format":       result.DocumentFormat,
		"analysis_type":         result.AnalysisType,
		"input_source":          result.InputSource,
		"processing_time":       fmt.Sprintf("%.2fs", result.ProcessingTime.Seconds()),
	}

	// Add document metadata if available
	if result.DocumentMetadata != nil {
		for k, v := range result.DocumentMetadata {
			data[k] = fmt.Sprintf("%v", v)
		}
	}

	return &managerspb.ProcessDocumentInputResponse{
		Id:                   result.ResponseID,
		RequestId:            requestID,
		ExtractedContent:     result.ExtractedContent,
		AnalysisResult:       result.AnalysisResult,
		Message:              result.ResponseMessage,
		Data:                 data,
		Actions:              protoActions,
		Timestamp:            result.ResponseTimestamp.Format("2006-01-02T15:04:05Z07:00"),
		Status:               result.ResponseStatus,
		Confidence:           result.ResponseConfidence,
		ProcessingConfidence: result.ProcessingConfidence,
		WordCount:            int32(result.WordCount),
		PageCount:            int32(result.PageCount),
	}, nil
}

// GetConversationStats retrieves conversation statistics for a user
func (s server) GetConversationStats(ctx context.Context, request *managerspb.GetConversationStatsRequest) (*managerspb.GetConversationStatsResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)

	// Extract user ID from JWT claims if not provided in request

	if userID == "" {
		claims, ok := auth.ClaimsFromContext(ctx)
		if !ok {
			return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
		}
		userID = claims.Subject
	}

	span.SetAttributes(
		attribute.String("UserID", userID),
	)

	stats, err := s.app.GetConversationStats(ctx, queries.GetConversationStats{
		UserID: userID,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		// Return default stats instead of error for better UX
		return &managerspb.GetConversationStatsResponse{
			Stats: &managerspb.ConversationStats{
				TotalConversations:         0,
				ActiveConversations:        0,
				TotalMessages:              0,
				MessagesToday:              0,
				MessagesThisWeek:           0,
				MessagesThisMonth:          0,
				FirstConversationAt:        "",
				LastConversationAt:         "",
				AvgMessagesPerConversation: 0.0,
				MostUsedManagerId:          "",
			},
		}, nil
	}

	return &managerspb.GetConversationStatsResponse{
		Stats: s.conversationStatsToProto(stats),
	}, nil
}

func (s server) conversationStatsToProto(stats *domain.ConversationStats) *managerspb.ConversationStats {
	if stats == nil {
		return &managerspb.ConversationStats{}
	}

	var firstConversationAt, lastConversationAt string
	if !stats.FirstConversationAt.IsZero() {
		firstConversationAt = stats.FirstConversationAt.Format(time.RFC3339)
	}
	if !stats.LastConversationAt.IsZero() {
		lastConversationAt = stats.LastConversationAt.Format(time.RFC3339)
	}

	return &managerspb.ConversationStats{
		TotalConversations:         stats.TotalConversations,
		ActiveConversations:        stats.ActiveConversations,
		TotalMessages:              stats.TotalMessages,
		MessagesToday:              stats.MessagesToday,
		MessagesThisWeek:           stats.MessagesThisWeek,
		MessagesThisMonth:          stats.MessagesThisMonth,
		FirstConversationAt:        firstConversationAt,
		LastConversationAt:         lastConversationAt,
		AvgMessagesPerConversation: stats.AvgMessagesPerConversation,
		MostUsedManagerId:          stats.MostUsedManagerID,
	}
}

// convertDomainActionsToProto converts domain actions to protobuf actions
func (s server) convertDomainActionsToProto(actions []domain.ManagerAction) []*managerspb.ManagerAction {
	protoActions := make([]*managerspb.ManagerAction, len(actions))
	for i, action := range actions {
		protoActions[i] = &managerspb.ManagerAction{
			Type:        action.Type,
			Endpoint:    action.Endpoint,
			Method:      action.Method,
			Parameters:  s.convertParametersToStringMap(action.Parameters),
			Description: action.Description,
		}
	}
	return protoActions
}

// ProcessManagerRequest handles manager request processing
func (s server) ProcessManagerRequest(ctx context.Context, request *managerspb.ProcessManagerRequestRequest) (*managerspb.ProcessManagerRequestResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ManagerID", request.GetManagerId()),
		attribute.String("RequestID", request.GetRequestId()),
		attribute.String("UserID", userID),
		attribute.String("Message", request.GetMessage()),
	)

	// Get the manager
	manager, err := s.app.GetManager(ctx, queries.GetManager{
		ID: request.GetManagerId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.Wrap(err, "failed to get manager")
	}

	// Use ProcessUserInput command for processing
	result, err := s.app.ProcessUserInput(ctx, commands.ProcessUserInput{
		ID:          request.GetRequestId(),
		ManagerID:   manager.ID,
		UserID:      userID,
		Message:     request.GetMessage(),
		Context:     convertStringMapToInterface(request.GetContext()),
		RequestType: "general",
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.Wrap(err, "failed to process manager request")
	}

	// Convert actions
	actions := s.convertDomainActionsToProto(result.ExecutedActions)

	// Create empty data map (ProcessUserInputResult doesn't have Data field)
	data := make(map[string]string)

	return &managerspb.ProcessManagerRequestResponse{
		Id:         result.ResponseID,
		RequestId:  request.GetRequestId(),
		Message:    result.ResponseMessage,
		Data:       data,
		Actions:    actions,
		Timestamp:  result.ResponseTimestamp.Format(time.RFC3339),
		Status:     result.ResponseStatus,
		Confidence: result.ResponseConfidence,
		ManagerId:  manager.ID,
		UserId:     userID,
	}, nil
}

// managerHasLatestConfiguration checks if manager has the latest consolidated configuration
func (s server) managerHasLatestConfiguration(manager *domain.CatalogManager) bool {
	expectedPrompt := processor.NewEnhancedLLMInterface().GetCompleteSystemPrompt()
	expectedCaps := s.getStandardCapabilities()

	// Check if system prompt matches
	if manager.SystemPrompt != expectedPrompt {
		return false
	}

	// Check if capabilities match
	if len(manager.Capabilities) != len(expectedCaps) {
		return false
	}

	capMap := make(map[domain.ManagerCapability]bool)
	for _, cap := range manager.Capabilities {
		capMap[cap] = true
	}

	for _, expectedCap := range expectedCaps {
		if !capMap[expectedCap] {
			return false
		}
	}

	return true
}

// getStandardCapabilities returns the standard set of capabilities for default managers
func (s server) getStandardCapabilities() []domain.ManagerCapability {
	return []domain.ManagerCapability{
		domain.CapabilityUserInteraction,
		domain.CapabilityDataAnalysis,
		domain.CapabilitySearchAndFilter,
		domain.CapabilityDataRetrieval,
		domain.CapabilityPublicAPIAccess,
		domain.CapabilityJailbreakResistant,
		domain.CapabilityScopeEnforcement,
		domain.CapabilityTextGeneration,
		domain.CapabilityCodeGeneration,
	}
}

// domainToProtoCapability converts a domain capability to protobuf capability
func domainToProtoCapability(cap domain.ManagerCapability) (managerspb.ManagerCapability, error) {
	switch cap {
	case domain.CapabilityManagerManagement:
		return managerspb.ManagerCapability_ASSISTANT_MANAGEMENT, nil
	case domain.CapabilityUserInteraction:
		return managerspb.ManagerCapability_USER_INTERACTION, nil
	case domain.CapabilityDataAnalysis:
		return managerspb.ManagerCapability_DATA_ANALYSIS, nil
	case domain.CapabilityLocationServices:
		return managerspb.ManagerCapability_LOCATION_SERVICES, nil
	case domain.CapabilityAuthentication:
		return managerspb.ManagerCapability_AUTHENTICATION, nil
	case domain.CapabilityPublicAPIAccess:
		return managerspb.ManagerCapability_PUBLIC_API_ACCESS, nil
	case domain.CapabilityJailbreakResistant:
		return managerspb.ManagerCapability_JAILBREAK_RESISTANT, nil
	case domain.CapabilityScopeEnforcement:
		return managerspb.ManagerCapability_SCOPE_ENFORCEMENT, nil
	case domain.CapabilityDataRetrieval:
		return managerspb.ManagerCapability_DATA_RETRIEVAL, nil
	case domain.CapabilitySearchAndFilter:
		return managerspb.ManagerCapability_SEARCH_AND_FILTER, nil
	case domain.CapabilityPrivateAPIAccess:
		return managerspb.ManagerCapability_PRIVATE_API_ACCESS, nil
	case domain.CapabilityUserDataAccess:
		return managerspb.ManagerCapability_USER_DATA_ACCESS, nil
	case domain.CapabilityTokenManagement:
		return managerspb.ManagerCapability_TOKEN_MANAGEMENT, nil
	case domain.CapabilityDataMasking:
		return managerspb.ManagerCapability_DATA_MASKING, nil
	case domain.CapabilityAuditLogging:
		return managerspb.ManagerCapability_AUDIT_LOGGING, nil
	case domain.CapabilityTextGeneration:
		return managerspb.ManagerCapability_TEXT_GENERATION, nil
	case domain.CapabilityCodeGeneration:
		return managerspb.ManagerCapability_CODE_GENERATION, nil
	case domain.CapabilityWebSearch:
		return managerspb.ManagerCapability_WEB_SEARCH, nil
	default:
		return 0, fmt.Errorf("unknown domain capability: %s", cap)
	}
}

// protoToDomainCapability converts a protobuf capability to domain capability
func protoToDomainCapability(cap managerspb.ManagerCapability) (domain.ManagerCapability, error) {
	switch cap {
	case managerspb.ManagerCapability_ASSISTANT_MANAGEMENT:
		return domain.CapabilityManagerManagement, nil
	case managerspb.ManagerCapability_USER_INTERACTION:
		return domain.CapabilityUserInteraction, nil
	case managerspb.ManagerCapability_DATA_ANALYSIS:
		return domain.CapabilityDataAnalysis, nil
	case managerspb.ManagerCapability_LOCATION_SERVICES:
		return domain.CapabilityLocationServices, nil
	case managerspb.ManagerCapability_AUTHENTICATION:
		return domain.CapabilityAuthentication, nil
	case managerspb.ManagerCapability_PUBLIC_API_ACCESS:
		return domain.CapabilityPublicAPIAccess, nil
	case managerspb.ManagerCapability_JAILBREAK_RESISTANT:
		return domain.CapabilityJailbreakResistant, nil
	case managerspb.ManagerCapability_SCOPE_ENFORCEMENT:
		return domain.CapabilityScopeEnforcement, nil
	case managerspb.ManagerCapability_DATA_RETRIEVAL:
		return domain.CapabilityDataRetrieval, nil
	case managerspb.ManagerCapability_SEARCH_AND_FILTER:
		return domain.CapabilitySearchAndFilter, nil
	case managerspb.ManagerCapability_PRIVATE_API_ACCESS:
		return domain.CapabilityPrivateAPIAccess, nil
	case managerspb.ManagerCapability_USER_DATA_ACCESS:
		return domain.CapabilityUserDataAccess, nil
	case managerspb.ManagerCapability_TOKEN_MANAGEMENT:
		return domain.CapabilityTokenManagement, nil
	case managerspb.ManagerCapability_DATA_MASKING:
		return domain.CapabilityDataMasking, nil
	case managerspb.ManagerCapability_AUDIT_LOGGING:
		return domain.CapabilityAuditLogging, nil
	case managerspb.ManagerCapability_TEXT_GENERATION:
		return domain.CapabilityTextGeneration, nil
	case managerspb.ManagerCapability_CODE_GENERATION:
		return domain.CapabilityCodeGeneration, nil
	case managerspb.ManagerCapability_WEB_SEARCH:
		return domain.CapabilityWebSearch, nil
	default:
		return "", fmt.Errorf("unknown protobuf capability: %v", cap)
	}
}

// domainToProtoCapabilities converts domain capabilities to protobuf capabilities with deduplication
func domainToProtoCapabilities(caps []domain.ManagerCapability) []managerspb.ManagerCapability {
	seen := make(map[managerspb.ManagerCapability]bool)
	result := make([]managerspb.ManagerCapability, 0, len(caps))

	for _, cap := range caps {
		protoCap, err := domainToProtoCapability(cap)
		if err != nil {
			log.Printf("Warning: failed to convert capability %s: %v", cap, err)
			continue
		}
		if !seen[protoCap] {
			seen[protoCap] = true
			result = append(result, protoCap)
		}
	}

	return result
}

// protoToDomainCapabilities converts protobuf capabilities to domain capabilities with deduplication
func protoToDomainCapabilities(caps []managerspb.ManagerCapability) []domain.ManagerCapability {
	seen := make(map[domain.ManagerCapability]bool)
	result := make([]domain.ManagerCapability, 0, len(caps))

	for _, cap := range caps {
		domainCap, err := protoToDomainCapability(cap)
		if err != nil {
			log.Printf("Warning: failed to convert capability %v: %v", cap, err)
			continue
		}
		if !seen[domainCap] {
			seen[domainCap] = true
			result = append(result, domainCap)
		}
	}

	return result
}
