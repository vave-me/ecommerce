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

	"middleman/assistants/assistantspb"
	"middleman/assistants/internal/application"
	"middleman/assistants/internal/application/commands"

	"middleman/assistants/internal/application/queries"
	"middleman/assistants/internal/constants"
	"middleman/assistants/internal/domain"
	"middleman/internal/auth"

	"google.golang.org/grpc"
)

type server struct {
	app application.App

	assistantspb.AssistantsServiceServer
}

var _ assistantspb.AssistantsServiceServer = (*server)(nil)

func RegisterServer(app application.App, sr, registrar grpc.ServiceRegistrar) error {
	assistantspb.RegisterAssistantsServiceServer(registrar, server{
		app: app,
	})
	return nil
}

// GetAssistant retrieves an assistant by ID
func (s server) GetAssistant(ctx context.Context, request *assistantspb.GetAssistantRequest) (*assistantspb.GetAssistantResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Log claims for debugging admin role population
	log.Printf("[GetAssistant] Claims extracted - UserID: %s, Role: %s", claims.Subject, claims.Role)

	assistant, err := s.app.GetAssistant(ctx, queries.GetAssistant{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &assistantspb.GetAssistantResponse{
		Assistant: s.assistantFromDomain(assistant),
	}, nil
}

// ActivateAssistant activates an assistant
func (s server) ActivateAssistant(ctx context.Context, request *assistantspb.ActivateAssistantRequest) (*assistantspb.ActivateAssistantResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Log claims for debugging admin role population
	log.Printf("[ActivateAssistant] Claims extracted - UserID: %s, Role: %s", claims.Subject, claims.Role)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("AssistantID", request.GetId()))

	err := s.app.ActivateAssistant(ctx, commands.ActivateAssistant{
		ID: request.GetId(),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &assistantspb.ActivateAssistantResponse{}, nil
}

// DeactivateAssistant deactivates an assistant
func (s server) DeactivateAssistant(ctx context.Context, request *assistantspb.DeactivateAssistantRequest) (*assistantspb.DeactivateAssistantResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Log claims for debugging admin role population
	log.Printf("[DeactivateAssistant] Claims extracted - UserID: %s, Role: %s", claims.Subject, claims.Role)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("AssistantID", request.GetId()))

	err := s.app.DeactivateAssistant(ctx, commands.DeactivateAssistant{
		ID: request.GetId(),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &assistantspb.DeactivateAssistantResponse{}, nil
}

// UpdateAssistantConfiguration updates assistant configuration
func (s server) UpdateAssistantConfiguration(ctx context.Context, request *assistantspb.UpdateAssistantConfigurationRequest) (*assistantspb.UpdateAssistantConfigurationResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[UpdateAssistantConfiguration] Claims extracted - UserID: %s, Role: %s", userID, claims.Role)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("AssistantID", request.GetId()))

	// Build configuration from structured fields
	config := commands.UpdateAssistantConfiguration{
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
		config.SystemPrompt = constants.SystemPrompt // Using production system prompt
		log.Printf("[UpdateAssistantConfiguration] Using consolidated prompt for assistant %s due to capability update", request.GetId())
	}

	if len(request.GetCapabilities()) > 0 {
		// Convert protobuf capabilities to domain capabilities with deduplication
		config.Capabilities = protoToDomainCapabilities(request.GetCapabilities())
	}

	err := s.app.UpdateAssistantConfiguration(ctx, config)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &assistantspb.UpdateAssistantConfigurationResponse{}, nil
}

// GetAssistants retrieves all assistants
func (s server) GetAssistants(ctx context.Context, request *assistantspb.GetAssistantsRequest) (*assistantspb.GetAssistantsResponse, error) {
	log.Printf("[GetAssistants] START - Request received")

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		log.Printf("[GetAssistants] ERROR - No authentication claims found")
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[GetAssistants] Claims extracted - UserID: %s, Role: %s, Limit: %d, Page: %d",
		userID, claims.Role, request.GetLimit(), request.GetPage())

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("UserID", userID),
		attribute.Int("RequestLimit", int(request.GetLimit())),
	)

	// Use application layer to get assistants with proper consistency
	log.Printf("[GetAssistants] Calling app.GetAssistants with UserID=%s, Limit=%d, Offset=%d",
		userID, int(request.GetLimit()), int(request.GetPage())*int(request.GetLimit()))

	assistants, err := s.app.GetAssistants(ctx, queries.GetAssistants{
		UserID: userID,
		Limit:  int(request.GetLimit()),
		Offset: int(request.GetPage()) * int(request.GetLimit()),
	})

	if err != nil {
		log.Printf("[GetAssistants] ERROR from app.GetAssistants: %v", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Printf("[GetAssistants] app.GetAssistants returned %d assistants", len(assistants))

	// Convert domain assistants to protobuf
	protoAssistants := make([]*assistantspb.Assistant, len(assistants))
	for i, assistant := range assistants {
		protoAssistants[i] = s.catalogAssistantFromDomain(assistant)
	}

	// If no assistants found, create and return a default assistant
	if len(protoAssistants) == 0 {
		log.Printf("[GetAssistants] No assistants found, attempting to create default assistant")
		defaultAssistant, err := s.createStandardDefaultAssistant(ctx, userID)
		if err != nil {
			log.Printf("[GetAssistants] ERROR - Failed to create default assistant: %v", err)
			// Return empty list rather than error to maintain compatibility
			return &assistantspb.GetAssistantsResponse{
				Assistants: []*assistantspb.Assistant{},
				TotalCount: 0,
			}, nil
		}
		log.Printf("[GetAssistants] Default assistant created successfully")
		protoAssistants = []*assistantspb.Assistant{s.assistantFromDomain(defaultAssistant)}
	}

	log.Printf("[GetAssistants] SUCCESS - Returning %d assistants", len(protoAssistants))
	return &assistantspb.GetAssistantsResponse{
		Assistants: protoAssistants,
		TotalCount: int32(len(protoAssistants)),
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

func (s server) assistantFromDomain(assistant *domain.CatalogAssistant) *assistantspb.Assistant {
	// Convert domain capabilities to protobuf capabilities with deduplication
	capabilities := domainToProtoCapabilities(assistant.Capabilities)

	// Convert domain type to proto type
	var protoType assistantspb.AssistantType
	switch assistant.Type {
	case domain.AssistantTypeAdmin:
		protoType = assistantspb.AssistantType_ADMIN
	case domain.AssistantTypeBusiness:
		protoType = assistantspb.AssistantType_BUSINESS
	case domain.AssistantTypeSupport:
		protoType = assistantspb.AssistantType_SUPPORT
	case domain.AssistantTypeScheduler:
		protoType = assistantspb.AssistantType_SCHEDULER
	default:
		protoType = assistantspb.AssistantType_STANDARD
	}

	// Create assistant with simple capability flags for frontend
	pbAssistant := &assistantspb.Assistant{
		Id:           assistant.ID,
		Name:         assistant.Name,
		Description:  assistant.Description,
		Type:         protoType,
		Capabilities: capabilities,
		Active:       assistant.Active,
		Temperature:  assistant.Temperature,
		MaxTokens:    int32(assistant.MaxTokens),
		SystemPrompt: assistant.SystemPrompt,
		CreatedAt:    assistant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    assistant.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return pbAssistant
}

func (s server) catalogAssistantFromDomain(assistant *domain.CatalogAssistant) *assistantspb.Assistant {
	// Convert domain capabilities to protobuf capabilities with deduplication
	capabilities := domainToProtoCapabilities(assistant.Capabilities)

	return &assistantspb.Assistant{
		Id:           assistant.ID,
		Name:         assistant.Name,
		Description:  assistant.Description,
		Capabilities: capabilities,
		Active:       assistant.Active,
		Temperature:  assistant.Temperature,
		MaxTokens:    int32(assistant.MaxTokens),
		SystemPrompt: assistant.SystemPrompt,
		CreatedAt:    assistant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    assistant.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Conversation Management Endpoints

// CreateConversation creates a new conversation
func (s server) CreateConversation(ctx context.Context, request *assistantspb.CreateConversationRequest) (*assistantspb.CreateConversationResponse, error) {
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
	log.Printf("[CREATE_CONVERSATION] Request AssistantID: %s", request.GetAssistantId())

	span.SetAttributes(
		attribute.String("ConversationID", conversationID),
		attribute.String("UserID", userID),
		attribute.String("AssistantID", request.GetAssistantId()),
	)

	// Convert protobuf context to map[string]interface{}
	initialContext := make(map[string]interface{})
	for k, v := range request.GetInitialContext() {
		initialContext[k] = v
	}
	// Context converted

	// If assistant ID is not provided, generate a new one
	assistantID := request.GetAssistantId()
	if assistantID == "" {
		assistantID = uuid.New().String()
		log.Printf("[CREATE_CONVERSATION] No AssistantID provided, generated new one: %s", assistantID)
	}

	// Create conversation via application layer
	log.Printf("[CREATE_CONVERSATION] Creating conversation with: ID=%s, UserID=%s, AssistantID=%s",
		conversationID, userID, assistantID)

	err := s.app.CreateConversation(ctx, commands.CreateConversation{
		ID:             conversationID,
		UserID:         userID,
		AssistantID:    assistantID,
		InitialContext: initialContext,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Conversation created successfully

	return &assistantspb.CreateConversationResponse{
		ConversationId: conversationID,
	}, nil
}

// GetConversation retrieves a conversation by ID
func (s server) GetConversation(ctx context.Context, request *assistantspb.GetConversationRequest) (*assistantspb.GetConversationResponse, error) {
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

	return &assistantspb.GetConversationResponse{
		Conversation: s.conversationViewToProto(conversation),
	}, nil
}

// GetUserConversations retrieves conversations for a user
func (s server) GetUserConversations(ctx context.Context, request *assistantspb.GetUserConversationsRequest) (*assistantspb.GetUserConversationsResponse, error) {
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

	conversations := make([]*assistantspb.Conversation, len(result.Conversations))
	for i, conv := range result.Conversations {
		conversations[i] = s.conversationViewToProto(conv)
	}

	return &assistantspb.GetUserConversationsResponse{
		Conversations: conversations,
		TotalCount:    int32(result.TotalCount),
	}, nil
}

// GetConversationMessages retrieves messages for a conversation
func (s server) GetConversationMessages(ctx context.Context, request *assistantspb.GetConversationMessagesRequest) (*assistantspb.GetConversationMessagesResponse, error) {
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

	protoMessages := make([]*assistantspb.ConversationMessage, len(messages))
	for i, msg := range messages {
		protoMessages[i] = s.conversationMessageToProto(msg)
	}

	return &assistantspb.GetConversationMessagesResponse{
		Messages:   protoMessages,
		TotalCount: int32(len(messages)), // Note: This should be the actual total count from pagination
	}, nil
}

// AddMessageToConversation adds a message to an existing conversation
func (s server) AddMessageToConversation(ctx context.Context, request *assistantspb.AddMessageToConversationRequest) (*assistantspb.AddMessageToConversationResponse, error) {
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
	case assistantspb.MessageRole_USER:
		role = domain.UserRole
	case assistantspb.MessageRole_ASSISTANT:
		role = domain.AssistantRole
	case assistantspb.MessageRole_SYSTEM:
		role = domain.SystemRole
	default:
		role = domain.UserRole
	}

	// For user messages, default to processing with LLM
	processWithLLM := role == domain.UserRole

	result, err := s.app.AddMessageToConversation(ctx, commands.AddMessageToConversation{
		ConversationID:     request.GetConversationId(),
		MessageID:          messageID,
		AssistantID:        request.GetAssistantId(),
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
	response := &assistantspb.AddMessageToConversationResponse{
		MessageId: result.MessageID,
	}

	// If LLM processing occurred, include the assistant's response
	if result.AssistantMessageID != "" {
		// For backward compatibility, still return the user's message ID
		// But also log that assistant responded
		log.Printf("[AddMessageToConversation] LLM processed - User msg: %s, Assistant msg: %s",
			result.MessageID, result.AssistantMessageID)
	}

	return response, nil
}

// ChatWithConversation processes a message within an existing conversation context
func (s server) ChatWithConversation(ctx context.Context, request *assistantspb.ChatWithConversationRequest) (*assistantspb.ChatWithConversationResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	// Log claims for debugging admin role population
	log.Printf("[ChatWithConversation] Claims extracted - UserID: %s, Role: %s", userID, claims.Role)

	log.Printf("[GRPC_CHAT] ========== ChatWithConversation START ==========")
	log.Printf("[GRPC_CHAT] Request received - ConversationID: %s, AssistantID: %s, UserID: %s, Message: %s",
		request.GetConversationId(), request.GetAssistantId(), userID, request.GetMessage())
	log.Printf("[GRPC_CHAT] MaxHistoryMessages: %d, Context keys: %d",
		request.GetMaxHistoryMessages(), len(request.GetContext()))

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConversationID", request.GetConversationId()),
		attribute.String("AssistantID", request.GetAssistantId()),
		attribute.String("UserID", userID),
		attribute.Int("MaxHistoryMessages", int(request.GetMaxHistoryMessages())),
	)

	// Get assistant ID - required field
	assistantID := request.GetAssistantId()
	if assistantID == "" {
		// If not provided, try to get user's existing assistant
		log.Printf("[GRPC_CHAT] No assistant ID provided, looking for user's existing assistant")
		assistants, err := s.app.GetAssistants(ctx, queries.GetAssistants{
			UserID: userID,
			Limit:  10,
		})
		if err == nil && len(assistants) > 0 {
			// Find first active assistant
			for _, assistant := range assistants {
				if assistant.Active {
					assistantID = assistant.ID
					log.Printf("[GRPC_CHAT] Found user's existing active assistant: %s", assistantID)
					break
				}
			}
		}

		if assistantID == "" {
			log.Printf("[GRPC_CHAT] ERROR: No assistant ID provided and no existing assistant found for user")
			return nil, status.Error(grpc_code.InvalidArgument, "assistant_id is required")
		}
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
			AssistantID:    assistantID,
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

	log.Printf("[GRPC_CHAT] Calling application layer ChatWithConversation with userID: %s, conversationID: %s, assistantID: %s...",
		userID, conversationID, assistantID)
	result, err := s.app.ChatWithConversation(ctx, commands.ChatWithConversation{
		ConversationID:     conversationID,
		AssistantID:        assistantID,
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
	actions := make([]*assistantspb.AssistantAction, len(result.Actions))
	for i, action := range result.Actions {
		log.Printf("[GRPC_CHAT] Converting action %d: Type=%s, Endpoint=%s", i, action.Type, action.Endpoint)
		actions[i] = &assistantspb.AssistantAction{
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
	response := &assistantspb.ChatWithConversationResponse{
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
func (s server) UpdateConversationContext(ctx context.Context, request *assistantspb.UpdateConversationContextRequest) (*assistantspb.UpdateConversationContextResponse, error) {
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

	return &assistantspb.UpdateConversationContextResponse{}, nil
}

// UpdateConversation updates conversation metadata (Frontend-expected endpoint)
func (s server) UpdateConversation(ctx context.Context, request *assistantspb.UpdateConversationRequest) (*assistantspb.UpdateConversationResponse, error) {
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

	// TODO: Implement UpdateConversation command in application layer
	// For now, return the same conversation structure that frontend expects
	conversation := &assistantspb.Conversation{
		Id:      request.GetConversationId(),
		UserId:  userID,
		Active:  true,
		Context: request.GetMetadata(),
	}

	return &assistantspb.UpdateConversationResponse{
		Conversation: conversation,
	}, nil
}

// DeleteConversation deletes a conversation (Frontend-expected endpoint)
func (s server) DeleteConversation(ctx context.Context, request *assistantspb.DeleteConversationRequest) (*assistantspb.DeleteConversationResponse, error) {
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

	// TODO: Implement DeleteConversation command in application layer
	// For now, return empty response as frontend expects
	return &assistantspb.DeleteConversationResponse{}, nil
}

// ArchiveConversation archives a conversation
func (s server) ArchiveConversation(ctx context.Context, request *assistantspb.ArchiveConversationRequest) (*assistantspb.ArchiveConversationResponse, error) {
	// This endpoint is not implemented in the application layer yet
	// It would require a new command handler
	return &assistantspb.ArchiveConversationResponse{}, nil
}

func (s server) conversationViewToProto(conv *domain.ReadConversation) *assistantspb.Conversation {
	// Convert context to string map
	context := make(map[string]string)
	if conv.Context != nil {
		for k, v := range conv.Context {
			context[k] = fmt.Sprintf("%v", v)
		}
	}

	return &assistantspb.Conversation{
		Id:          conv.ID,
		UserId:      conv.UserID,
		AssistantId: conv.AssistantID,
		Messages:    []*assistantspb.ConversationMessage{}, // Messages loaded separately
		CreatedAt:   conv.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   conv.UpdatedAt.Format(time.RFC3339),
		Active:      conv.Active,
		Context:     context,
	}
}

func (s server) conversationMessageToProto(msg *domain.ReadMessage) *assistantspb.ConversationMessage {

	// Convert metadata to string map
	metadata := make(map[string]string)
	if msg.Metadata != nil {
		for k, v := range msg.Metadata {
			metadata[k] = fmt.Sprintf("%v", v)
		}
	}

	// Convert domain role to protobuf role
	var role assistantspb.MessageRole
	switch msg.Role {
	case domain.UserRole:
		role = assistantspb.MessageRole_USER
	case domain.AssistantRole:
		role = assistantspb.MessageRole_ASSISTANT
	case domain.SystemRole:
		role = assistantspb.MessageRole_SYSTEM
	default:
		role = assistantspb.MessageRole_USER
	}

	// Convert actions
	actions := make([]*assistantspb.AssistantAction, len(msg.ActionsTaken))
	for i, action := range msg.ActionsTaken {
		actions[i] = &assistantspb.AssistantAction{
			Type:        action.Type,
			Endpoint:    action.Endpoint,
			Method:      action.Method,
			Parameters:  s.convertParametersToStringMap(action.Parameters),
			Description: action.Description,
		}
	}

	return &assistantspb.ConversationMessage{
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

// ProcessUserInput processes user input through an assistant
func (s server) ProcessUserInput(ctx context.Context, request *assistantspb.ProcessUserInputRequest) (*assistantspb.ProcessUserInputResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("AssistantID", request.GetAssistantId()),
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
		AssistantID: request.GetAssistantId(),
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
			return &assistantspb.ProcessUserInputResponse{
				Id:         uuid.New().String(),
				RequestId:  requestID,
				Message:    "Your request is taking longer than expected to process. Please try a simpler request or try again later.",
				Data:       map[string]string{"error": "processing_timeout"},
				Actions:    []*assistantspb.AssistantAction{},
				Timestamp:  time.Now().Format("2006-01-02T15:04:05Z07:00"),
				Status:     "timeout",
				Confidence: 0.1,
			}, nil
		}

		return nil, err
	}

	// Convert domain actions to protobuf actions
	var protoActions []*assistantspb.AssistantAction
	for _, action := range result.ExecutedActions {
		protoActions = append(protoActions, &assistantspb.AssistantAction{
			Type:        action.Type,
			Endpoint:    action.Endpoint,
			Method:      action.Method,
			Parameters:  s.convertParametersToStringMap(action.Parameters),
			Description: action.Description,
		})
	}

	return &assistantspb.ProcessUserInputResponse{
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
func (s server) ProcessImageInput(ctx context.Context, req *assistantspb.ProcessImageInputRequest) (*assistantspb.ProcessImageInputResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	log.Printf("[ProcessImageInput] Starting image processing with model selection. AssistantID: %s, ImageFormat: %s, Analysis: %s",
		req.AssistantId, req.ImageFormat, req.AnalysisType)

	// Validate request
	if req.AssistantId == "" {
		return nil, status.Errorf(grpc_code.InvalidArgument, "assistant_id is required")
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
		AssistantID:    req.AssistantId,
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

	return &assistantspb.ProcessImageInputResponse{
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
func (s server) detectImageProcessingMode(req *assistantspb.ProcessImageInputRequest) string {

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
func (s server) ProcessSpeechInput(ctx context.Context, req *assistantspb.ProcessSpeechInputRequest) (*assistantspb.ProcessSpeechInputResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	log.Printf("[ProcessSpeechInput] Starting speech processing with model selection. AssistantID: %s, AudioFormat: %s",
		req.AssistantId, req.AudioFormat)

	// Validate request
	if req.AssistantId == "" {
		return nil, status.Errorf(grpc_code.InvalidArgument, "assistant_id is required")
	}
	if len(req.AudioData) == 0 {
		return nil, status.Errorf(grpc_code.InvalidArgument, "audio_data is required")
	}

	// Create command with proper input type detection
	cmd := commands.ProcessSpeechInput{
		ID:          uuid.New().String(),
		AssistantID: req.AssistantId,
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

	return &assistantspb.ProcessSpeechInputResponse{
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
func (s server) ProcessDocumentInput(ctx context.Context, request *assistantspb.ProcessDocumentInputRequest) (*assistantspb.ProcessDocumentInputResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("AssistantID", request.GetAssistantId()),
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
		AssistantID:    request.GetAssistantId(),
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
	var protoActions []*assistantspb.AssistantAction
	for _, action := range result.ExecutedActions {
		protoActions = append(protoActions, &assistantspb.AssistantAction{
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

	return &assistantspb.ProcessDocumentInputResponse{
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
func (s server) GetConversationStats(ctx context.Context, request *assistantspb.GetConversationStatsRequest) (*assistantspb.GetConversationStatsResponse, error) {

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
		return &assistantspb.GetConversationStatsResponse{
			Stats: &assistantspb.ConversationStats{
				TotalConversations:         0,
				ActiveConversations:        0,
				TotalMessages:              0,
				MessagesToday:              0,
				MessagesThisWeek:           0,
				MessagesThisMonth:          0,
				FirstConversationAt:        "",
				LastConversationAt:         "",
				AvgMessagesPerConversation: 0.0,
				MostUsedAssistantId:        "",
			},
		}, nil
	}

	return &assistantspb.GetConversationStatsResponse{
		Stats: s.conversationStatsToProto(stats),
	}, nil
}

func (s server) conversationStatsToProto(stats *domain.ConversationStats) *assistantspb.ConversationStats {
	if stats == nil {
		return &assistantspb.ConversationStats{}
	}

	var firstConversationAt, lastConversationAt string
	if !stats.FirstConversationAt.IsZero() {
		firstConversationAt = stats.FirstConversationAt.Format(time.RFC3339)
	}
	if !stats.LastConversationAt.IsZero() {
		lastConversationAt = stats.LastConversationAt.Format(time.RFC3339)
	}

	return &assistantspb.ConversationStats{
		TotalConversations:         stats.TotalConversations,
		ActiveConversations:        stats.ActiveConversations,
		TotalMessages:              stats.TotalMessages,
		MessagesToday:              stats.MessagesToday,
		MessagesThisWeek:           stats.MessagesThisWeek,
		MessagesThisMonth:          stats.MessagesThisMonth,
		FirstConversationAt:        firstConversationAt,
		LastConversationAt:         lastConversationAt,
		AvgMessagesPerConversation: stats.AvgMessagesPerConversation,
		MostUsedAssistantId:        stats.MostUsedAssistantID,
	}
}

// convertDomainActionsToProto converts domain actions to protobuf actions
func (s server) convertDomainActionsToProto(actions []domain.AssistantAction) []*assistantspb.AssistantAction {
	protoActions := make([]*assistantspb.AssistantAction, len(actions))
	for i, action := range actions {
		protoActions[i] = &assistantspb.AssistantAction{
			Type:        action.Type,
			Endpoint:    action.Endpoint,
			Method:      action.Method,
			Parameters:  s.convertParametersToStringMap(action.Parameters),
			Description: action.Description,
		}
	}
	return protoActions
}

// ProcessAssistantRequest handles assistant request processing
func (s server) ProcessAssistantRequest(ctx context.Context, request *assistantspb.ProcessAssistantRequestRequest) (*assistantspb.ProcessAssistantRequestResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("AssistantID", request.GetAssistantId()),
		attribute.String("RequestID", request.GetRequestId()),
		attribute.String("UserID", userID),
		attribute.String("Message", request.GetMessage()),
	)

	// Get the assistant
	assistant, err := s.app.GetAssistant(ctx, queries.GetAssistant{
		ID: request.GetAssistantId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.Wrap(err, "failed to get assistant")
	}

	// Use ProcessUserInput command for processing
	result, err := s.app.ProcessUserInput(ctx, commands.ProcessUserInput{
		ID:          request.GetRequestId(),
		AssistantID: assistant.ID,
		UserID:      userID,
		Message:     request.GetMessage(),
		Context:     convertStringMapToInterface(request.GetContext()),
		RequestType: "general",
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.Wrap(err, "failed to process assistant request")
	}

	// Convert actions
	actions := s.convertDomainActionsToProto(result.ExecutedActions)

	// Create empty data map (ProcessUserInputResult doesn't have Data field)
	data := make(map[string]string)

	return &assistantspb.ProcessAssistantRequestResponse{
		Id:          result.ResponseID,
		RequestId:   request.GetRequestId(),
		Message:     result.ResponseMessage,
		Data:        data,
		Actions:     actions,
		Timestamp:   result.ResponseTimestamp.Format(time.RFC3339),
		Status:      result.ResponseStatus,
		Confidence:  result.ResponseConfidence,
		AssistantId: assistant.ID,
		UserId:      userID,
	}, nil
}

// createStandardDefaultAssistant creates a default assistant with consolidated prompt and capabilities
func (s server) createStandardDefaultAssistant(ctx context.Context, userID string) (*domain.CatalogAssistant, error) {

	// Generate a new assistant ID
	assistantID := uuid.New().String()

	err := s.app.CreateUserAssistant(ctx, commands.CreateUserAssistant{
		ID:     assistantID,
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create default assistant: %w", err)
	}

	// Retrieve the created assistant
	assistant, err := s.app.GetAssistant(ctx, queries.GetAssistant{ID: assistantID})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created assistant: %w", err)
	}

	return assistant, nil
}

// assistantHasLatestConfiguration checks if assistant has the latest consolidated configuration
func (s server) assistantHasLatestConfiguration(assistant *domain.CatalogAssistant) bool {
	expectedPrompt := constants.SystemPrompt // Using production system prompt
	expectedCaps := s.getStandardCapabilities()

	// Check if system prompt matches
	if assistant.SystemPrompt != expectedPrompt {
		return false
	}

	// Check if capabilities match
	if len(assistant.Capabilities) != len(expectedCaps) {
		return false
	}

	capMap := make(map[domain.AssistantCapability]bool)
	for _, cap := range assistant.Capabilities {
		capMap[cap] = true
	}

	for _, expectedCap := range expectedCaps {
		if !capMap[expectedCap] {
			return false
		}
	}

	return true
}

// getStandardCapabilities returns the standard set of capabilities for default assistants
func (s server) getStandardCapabilities() []domain.AssistantCapability {
	return []domain.AssistantCapability{
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
func domainToProtoCapability(cap domain.AssistantCapability) (assistantspb.AssistantCapability, error) {
	switch cap {
	case domain.CapabilityAssistantManagement:
		return assistantspb.AssistantCapability_ASSISTANT_MANAGEMENT, nil
	case domain.CapabilityUserInteraction:
		return assistantspb.AssistantCapability_USER_INTERACTION, nil
	case domain.CapabilityDataAnalysis:
		return assistantspb.AssistantCapability_DATA_ANALYSIS, nil
	case domain.CapabilityLocationServices:
		return assistantspb.AssistantCapability_LOCATION_SERVICES, nil
	case domain.CapabilityAuthentication:
		return assistantspb.AssistantCapability_AUTHENTICATION, nil
	case domain.CapabilityPublicAPIAccess:
		return assistantspb.AssistantCapability_PUBLIC_API_ACCESS, nil
	case domain.CapabilityJailbreakResistant:
		return assistantspb.AssistantCapability_JAILBREAK_RESISTANT, nil
	case domain.CapabilityScopeEnforcement:
		return assistantspb.AssistantCapability_SCOPE_ENFORCEMENT, nil
	case domain.CapabilityDataRetrieval:
		return assistantspb.AssistantCapability_DATA_RETRIEVAL, nil
	case domain.CapabilitySearchAndFilter:
		return assistantspb.AssistantCapability_SEARCH_AND_FILTER, nil
	case domain.CapabilityPrivateAPIAccess:
		return assistantspb.AssistantCapability_PRIVATE_API_ACCESS, nil
	case domain.CapabilityUserDataAccess:
		return assistantspb.AssistantCapability_USER_DATA_ACCESS, nil
	case domain.CapabilityTokenManagement:
		return assistantspb.AssistantCapability_TOKEN_MANAGEMENT, nil
	case domain.CapabilityDataMasking:
		return assistantspb.AssistantCapability_DATA_MASKING, nil
	case domain.CapabilityAuditLogging:
		return assistantspb.AssistantCapability_AUDIT_LOGGING, nil
	case domain.CapabilityTextGeneration:
		return assistantspb.AssistantCapability_TEXT_GENERATION, nil
	case domain.CapabilityCodeGeneration:
		return assistantspb.AssistantCapability_CODE_GENERATION, nil
	case domain.CapabilityWebSearch:
		return assistantspb.AssistantCapability_WEB_SEARCH, nil
	default:
		return 0, fmt.Errorf("unknown domain capability: %s", cap)
	}
}

// protoToDomainCapability converts a protobuf capability to domain capability
func protoToDomainCapability(cap assistantspb.AssistantCapability) (domain.AssistantCapability, error) {
	switch cap {
	case assistantspb.AssistantCapability_ASSISTANT_MANAGEMENT:
		return domain.CapabilityAssistantManagement, nil
	case assistantspb.AssistantCapability_USER_INTERACTION:
		return domain.CapabilityUserInteraction, nil
	case assistantspb.AssistantCapability_DATA_ANALYSIS:
		return domain.CapabilityDataAnalysis, nil
	case assistantspb.AssistantCapability_LOCATION_SERVICES:
		return domain.CapabilityLocationServices, nil
	case assistantspb.AssistantCapability_AUTHENTICATION:
		return domain.CapabilityAuthentication, nil
	case assistantspb.AssistantCapability_PUBLIC_API_ACCESS:
		return domain.CapabilityPublicAPIAccess, nil
	case assistantspb.AssistantCapability_JAILBREAK_RESISTANT:
		return domain.CapabilityJailbreakResistant, nil
	case assistantspb.AssistantCapability_SCOPE_ENFORCEMENT:
		return domain.CapabilityScopeEnforcement, nil
	case assistantspb.AssistantCapability_DATA_RETRIEVAL:
		return domain.CapabilityDataRetrieval, nil
	case assistantspb.AssistantCapability_SEARCH_AND_FILTER:
		return domain.CapabilitySearchAndFilter, nil
	case assistantspb.AssistantCapability_PRIVATE_API_ACCESS:
		return domain.CapabilityPrivateAPIAccess, nil
	case assistantspb.AssistantCapability_USER_DATA_ACCESS:
		return domain.CapabilityUserDataAccess, nil
	case assistantspb.AssistantCapability_TOKEN_MANAGEMENT:
		return domain.CapabilityTokenManagement, nil
	case assistantspb.AssistantCapability_DATA_MASKING:
		return domain.CapabilityDataMasking, nil
	case assistantspb.AssistantCapability_AUDIT_LOGGING:
		return domain.CapabilityAuditLogging, nil
	case assistantspb.AssistantCapability_TEXT_GENERATION:
		return domain.CapabilityTextGeneration, nil
	case assistantspb.AssistantCapability_CODE_GENERATION:
		return domain.CapabilityCodeGeneration, nil
	case assistantspb.AssistantCapability_WEB_SEARCH:
		return domain.CapabilityWebSearch, nil
	default:
		return "", fmt.Errorf("unknown protobuf capability: %v", cap)
	}
}

// domainToProtoCapabilities converts domain capabilities to protobuf capabilities with deduplication
func domainToProtoCapabilities(caps []domain.AssistantCapability) []assistantspb.AssistantCapability {
	seen := make(map[assistantspb.AssistantCapability]bool)
	result := make([]assistantspb.AssistantCapability, 0, len(caps))

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
func protoToDomainCapabilities(caps []assistantspb.AssistantCapability) []domain.AssistantCapability {
	seen := make(map[domain.AssistantCapability]bool)
	result := make([]domain.AssistantCapability, 0, len(caps))

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
