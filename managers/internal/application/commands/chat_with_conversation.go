package commands

import (
	"context"
	"fmt"
	"middleman/internal/ddd"
	"middleman/managers/internal/application/services"
	"middleman/managers/internal/domain"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/stackus/errors"
)

type ChatWithConversation struct {
	ManagerID          string                 `json:"manager_id"`
	ConversationID     string                 `json:"conversation_id"`
	UserID             string                 `json:"user_id"`
	Message            string                 `json:"message"`
	Context            map[string]interface{} `json:"context,omitempty"`
	MaxHistoryMessages int                    `json:"max_history_messages,omitempty"`
}

type ChatWithConversationResult struct {
	Response   string                 `json:"response"`
	Actions    []domain.ManagerAction `json:"actions,omitempty"`
	Confidence float64                `json:"confidence"`
	Status     string                 `json:"status"`
	Data       map[string]interface{} `json:"data,omitempty"`
	MessageID  string                 `json:"message_id"`
}

type ChatWithConversationHandler struct {
	conversations     domain.ConversationRepository
	readConversations domain.ReadConversationRepository
	readMessages      domain.ReadMessagesRepository
	managers          domain.ManagerRepository
	llmProcessor      services.LLMProcessor
	publisher         ddd.EventPublisher[ddd.Event]
}

func NewChatWithConversationHandler(
	conversations domain.ConversationRepository,
	readConversations domain.ReadConversationRepository,
	readMessages domain.ReadMessagesRepository,
	managers domain.ManagerRepository,
	llmProcessor services.LLMProcessor,
	publisher ddd.EventPublisher[ddd.Event],
) ChatWithConversationHandler {
	return ChatWithConversationHandler{
		conversations:     conversations,
		readConversations: readConversations,
		readMessages:      readMessages,
		managers:          managers,
		llmProcessor:      llmProcessor,
		publisher:         publisher,
	}
}

func (h ChatWithConversationHandler) ChatWithConversation(ctx context.Context, cmd ChatWithConversation) (*ChatWithConversationResult, error) {
	// Load and validate conversation
	conversation, err := h.conversations.Load(ctx, cmd.ConversationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load conversation")
	}

	// Get manager ID - prefer from command, fallback to conversation
	managerID := cmd.ManagerID
	if managerID == "" {
		managerID = conversation.ManagerID
	}
	if managerID == "" {
		return nil, errors.Wrap(errors.ErrNotFound, "no manager ID provided and conversation has no manager ID")
	}

	// Load manager configuration
	manager, err := h.managers.Load(ctx, managerID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load manager")
	}

	// Prepare conversation history
	history, err := h.loadConversationHistory(ctx, cmd.ConversationID, cmd.UserID, cmd.MaxHistoryMessages)
	if err != nil {
		return nil, err
	}

	// Build enhanced context
	enhancedContext := h.buildEnhancedContext(conversation, cmd, history)

	// Process message with LLM
	response, actions, confidence, err := h.llmProcessor.ProcessWithHistory(
		ctx,
		manager,
		cmd.Message,
		history,
		enhancedContext,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process message with LLM")
	}

	// Generate message IDs
	userMessageID := uuid.New().String()
	managerMessageID := uuid.New().String()

	// Add user message
	userEvent, err := conversation.AddMessage(
		cmd.ConversationID,
		cmd.ManagerID,
		userMessageID,
		domain.UserRole,
		cmd.Message,
		map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"user_id":   cmd.UserID,
		},
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add user message")
	}

	// SAVE user message
	if err = h.conversations.Save(ctx, conversation); err != nil {
		return nil, errors.Wrap(err, "error saving user message")
	}

	// PUBLISH user message event immediately
	if err = h.publisher.Publish(ctx, userEvent); err != nil {
		return nil, errors.Wrap(err, "error publishing user message event")
	}

	// RELOAD conversation after user message
	conversation, err = h.conversations.Load(ctx, cmd.ConversationID)
	if err != nil {
		return nil, errors.Wrap(err, "error reloading conversation after user message")
	}

	// Add manager message
	managerEvent, err := conversation.AddMessage(
		cmd.ConversationID,
		cmd.ManagerID,
		managerMessageID,
		domain.ManagerRole,
		response,
		map[string]interface{}{
			"timestamp":  time.Now().Unix(),
			"confidence": confidence,
			"actions":    len(actions),
		},
		actions,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add manager message")
	}

	// SAVE manager message
	if err = h.conversations.Save(ctx, conversation); err != nil {
		return nil, errors.Wrap(err, "error saving manager message")
	}

	// PUBLISH manager message event immediately
	if err = h.publisher.Publish(ctx, managerEvent); err != nil {
		return nil, errors.Wrap(err, "error publishing manager message event")
	}

	// RELOAD conversation after manager message
	conversation, err = h.conversations.Load(ctx, cmd.ConversationID)
	if err != nil {
		return nil, errors.Wrap(err, "error reloading conversation after manager message")
	}

	// Update context
	updatedContext := h.buildUpdatedContext(enhancedContext, history, response, actions, confidence)
	contextEvent, err := conversation.UpdateContext(updatedContext)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update conversation context")
	}

	// SAVE context update
	if err = h.conversations.Save(ctx, conversation); err != nil {
		return nil, errors.Wrap(err, "error saving context update")
	}

	// PUBLISH context update event immediately
	if err = h.publisher.Publish(ctx, contextEvent); err != nil {
		return nil, errors.Wrap(err, "error publishing context update event")
	}

	// Process structured data for frontend
	structuredData, naturalLanguageResponse := h.processActionsForDisplay(actions, response, cmd.Message, managerID)

	return &ChatWithConversationResult{
		Response:   naturalLanguageResponse,
		Actions:    actions,
		Confidence: confidence,
		Status:     "success",
		Data:       structuredData,
		MessageID:  managerMessageID,
	}, nil
}

// loadConversationHistory retrieves message history from read model
func (h ChatWithConversationHandler) loadConversationHistory(ctx context.Context, conversationID, userID string, maxMessages int) ([]domain.ConversationMessage, error) {
	limit := maxMessages
	if limit <= 0 {
		limit = 50
	}

	messageViews, err := h.readMessages.GetConversationMessages(ctx, conversationID, userID, limit, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get conversation history")
	}

	history := make([]domain.ConversationMessage, len(messageViews))
	for i, msgView := range messageViews {
		history[i] = domain.ConversationMessage{
			ID:           msgView.ID,
			Role:         msgView.Role,
			Content:      msgView.Content,
			Timestamp:    msgView.Timestamp,
			Metadata:     msgView.Metadata,
			ActionsTaken: msgView.ActionsTaken,
		}
	}

	return history, nil
}

// buildEnhancedContext merges conversation and request contexts
func (h ChatWithConversationHandler) buildEnhancedContext(conversation *domain.Conversation, cmd ChatWithConversation, history []domain.ConversationMessage) map[string]interface{} {
	enhancedContext := make(map[string]interface{})

	// Copy existing conversation context
	for k, v := range conversation.Context {
		enhancedContext[k] = v
	}

	// Merge request context
	for k, v := range cmd.Context {
		enhancedContext[k] = v
	}

	// Add metadata
	enhancedContext["conversation_id"] = cmd.ConversationID
	enhancedContext["user_id"] = cmd.UserID
	enhancedContext["manager_id"] = conversation.ManagerID
	enhancedContext["message_count"] = len(history)
	enhancedContext["last_interaction"] = time.Now().Unix()

	return enhancedContext
}

// buildUpdatedContext creates context update with conversation insights
func (h ChatWithConversationHandler) buildUpdatedContext(baseContext map[string]interface{}, history []domain.ConversationMessage, response string, actions []domain.ManagerAction, confidence float64) map[string]interface{} {
	updatedContext := make(map[string]interface{})
	for k, v := range baseContext {
		updatedContext[k] = v
	}

	// Add AI-derived insights
	updatedContext["last_response_confidence"] = confidence
	updatedContext["tools_used_count"] = len(actions)
	updatedContext["total_message_count"] = len(history) + 2
	updatedContext["last_interaction_timestamp"] = time.Now().Unix()
	updatedContext["conversation_active"] = true
	updatedContext["last_user_message_length"] = len(baseContext["user_id"].(string))
	updatedContext["last_manager_response_length"] = len(response)

	// Quality metrics
	switch {
	case confidence > 0.8:
		updatedContext["last_interaction_quality"] = "high"
	case confidence > 0.6:
		updatedContext["last_interaction_quality"] = "medium"
	default:
		updatedContext["last_interaction_quality"] = "low"
	}

	// Tool usage tracking
	if len(actions) > 0 {
		toolNames := make([]string, 0, len(actions))
		toolTypes := make(map[string]bool)

		for _, action := range actions {
			toolNames = append(toolNames, action.Type)
			if action.Type != "" {
				toolTypes[action.Type] = true
			}
		}

		updatedContext["last_tools_used"] = toolNames
		updatedContext["tool_types_used"] = mapKeys(toolTypes)
	}

	return updatedContext
}

// processActionsForDisplay extracts structured data from actions for frontend display
func (h ChatWithConversationHandler) processActionsForDisplay(actions []domain.ManagerAction, response, userMessage, managerID string) (map[string]interface{}, string) {
	structuredData := make(map[string]interface{})
	structuredData["metadata"] = map[string]interface{}{
		"queryInterpretation": userMessage,
		"managerId":           managerID,
	}

	var items []interface{}

	for _, action := range actions {
		if action.Type != "tool_execution" || action.Parameters == nil {
			continue
		}

		result, ok := action.Parameters["result"]
		if !ok || result == nil {
			continue
		}

		resultMap, ok := result.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract entity type and process results
		entityType := extractEntityType(action.Endpoint)

		if results, exists := resultMap["results"]; exists {
			if resultArray, ok := results.([]interface{}); ok {
				items = h.convertResultsToCards(resultArray, entityType, items)
			}

			// Store raw results for compatibility
			structuredData[entityType] = results
		}

		// Update metadata
		metadata := structuredData["metadata"].(map[string]interface{})
		if count, exists := resultMap["count"]; exists {
			metadata[entityType+"_count"] = count
			metadata["totalCount"] = count
		}
	}

	if len(items) > 0 {
		structuredData["items"] = items
	}

	naturalLanguageResponse := cleanNaturalLanguageResponse(response, len(items) > 0)

	return structuredData, naturalLanguageResponse
}

// convertResultsToCards converts raw results to frontend-ready card format
func (h ChatWithConversationHandler) convertResultsToCards(results []interface{}, entityType string, existingItems []interface{}) []interface{} {
	timestamp := time.Now()
	items := existingItems

	for i, item := range results {
		cardItem := map[string]interface{}{
			"id":             fmt.Sprintf("%s_%d_%d", entityType, timestamp.Unix(), i),
			"entityType":     mapEntityTypeForFrontend(entityType),
			"source":         "manager",
			"timestamp":      timestamp.Format(time.RFC3339),
			"relevanceScore": calculateRelevanceScore(i, len(results)),
		}

		// Add entity data under appropriate key
		frontendType := mapEntityTypeForFrontend(entityType)
		cardItem[frontendType] = item

		items = append(items, cardItem)
	}

	return items
}

// Utility functions

func mapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// extractEntityType extracts the entity type from endpoint name
func extractEntityType(endpoint string) string {
	entityTypeMap := map[string]string{
		"product": "products",
		"service": "services",
		"post":    "posts",
		"user":    "users",
		"order":   "orders",
	}

	lowerEndpoint := strings.ToLower(endpoint)
	for key, value := range entityTypeMap {
		if strings.Contains(lowerEndpoint, key) {
			return value
		}
	}

	// Fallback: extract first part before underscore
	if parts := strings.Split(endpoint, "_"); len(parts) > 0 {
		return parts[0]
	}

	return "unknown"
}

// mapEntityTypeForFrontend maps backend entity types to frontend card types
func mapEntityTypeForFrontend(entityType string) string {
	mapping := map[string]string{
		"products": "product",
		"product":  "product",
		"services": "service",
		"service":  "service",
		"posts":    "post",
		"post":     "post",
		"users":    "user",
		"user":     "user",
		"orders":   "order",
		"order":    "order",
	}

	if mapped, exists := mapping[entityType]; exists {
		return mapped
	}
	return entityType
}

// calculateRelevanceScore calculates relevance based on position in results
func calculateRelevanceScore(position, total int) float64 {
	if total == 0 {
		return 1.0
	}
	// Higher score for items appearing earlier
	return 1.0 - (float64(position) / float64(total) * 0.5)
}

// cleanNaturalLanguageResponse removes data dumps from response
func cleanNaturalLanguageResponse(response string, hasStructuredData bool) string {
	if !hasStructuredData {
		return response
	}

	lines := strings.Split(response, "\n")
	var cleanedLines []string
	var inDataSection bool

	dataPatterns := []string{"Item ", "Product ", "Service ", "ID:", "Name:"}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if entering data section
		if !inDataSection {
			for _, pattern := range dataPatterns {
				if strings.HasPrefix(trimmed, pattern) ||
					(pattern == "Name:" && strings.Contains(trimmed, pattern)) {
					inDataSection = true
					break
				}
			}
		}

		// Skip data lines
		if inDataSection {
			continue
		}

		// Keep non-empty summary lines
		if trimmed != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	result := strings.TrimSpace(strings.Join(cleanedLines, "\n"))

	// Provide default message if everything was removed
	if result == "" {
		result = "I found the results you were looking for. Here they are:"
	}

	return result
}
