package commands

import (
	"context"
	"github.com/google/uuid"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/managers/internal/application/services"
	"middleman/managers/internal/domain"
	"time"
)

type AddMessageToConversation struct {
	ConversationID     string                 `json:"conversation_id"`
	ManagerID          string                 `json:"manager_id"`
	MessageID          string                 `json:"message_id"`
	Role               domain.MessageRole     `json:"role"`
	Content            string                 `json:"content"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
	ActionsTaken       []domain.ManagerAction `json:"actions_taken,omitempty"`
	UserID             string                 `json:"user_id,omitempty"`
	ProcessWithLLM     bool                   `json:"process_with_llm,omitempty"`
	MaxHistoryMessages int                    `json:"max_history_messages,omitempty"`
}

type AddMessageToConversationResult struct {
	MessageID        string                 `json:"message_id"`
	ManagerMessageID string                 `json:"manager_message_id,omitempty"`
	ManagerResponse  string                 `json:"manager_response,omitempty"`
	Actions          []domain.ManagerAction `json:"actions,omitempty"`
	Confidence       float64                `json:"confidence,omitempty"`
}

type AddMessageToConversationHandler struct {
	conversations domain.ConversationRepository
	readMessages  domain.ReadMessagesRepository
	managers      domain.ManagerRepository
	llmProcessor  services.LLMProcessor
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewAddMessageToConversationHandler(
	conversations domain.ConversationRepository,
	readMessages domain.ReadMessagesRepository,
	managers domain.ManagerRepository,
	llmProcessor services.LLMProcessor,
	publisher ddd.EventPublisher[ddd.Event],
) AddMessageToConversationHandler {
	return AddMessageToConversationHandler{
		conversations: conversations,
		readMessages:  readMessages,
		managers:      managers,
		llmProcessor:  llmProcessor,
		publisher:     publisher,
	}
}

func (h AddMessageToConversationHandler) AddMessageToConversation(ctx context.Context, cmd AddMessageToConversation) (*AddMessageToConversationResult, error) {
	conversation, err := h.conversations.Load(ctx, cmd.ConversationID)
	if err != nil {
		return nil, errors.Wrap(err, "error loading conversation")
	}

	// Add the user message
	event, err := conversation.AddMessage(
		cmd.ConversationID,
		cmd.ManagerID,
		cmd.MessageID,
		cmd.Role,
		cmd.Content,
		cmd.Metadata,
		cmd.ActionsTaken,
	)
	if err != nil {
		return nil, err
	}

	// workflow is corrumpeted
	result := &AddMessageToConversationResult{
		MessageID: cmd.MessageID,
	}

	// If LLM processing is requested and this is a user message
	if cmd.ProcessWithLLM && cmd.Role == domain.UserRole && cmd.UserID != "" {
		// Load manager
		if conversation.ManagerID == "" {
			return nil, errors.Wrap(errors.ErrNotFound, "conversation has no manager")
		}

		manager, err := h.managers.Load(ctx, conversation.ManagerID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load manager")
		}

		// Load conversation history
		history, err := h.loadConversationHistory(ctx, cmd.ConversationID, cmd.UserID, cmd.MaxHistoryMessages)
		if err != nil {
			return nil, err
		}

		// Build context
		enhancedContext := h.buildContext(conversation, cmd, history)

		// Process with LLM
		response, actions, confidence, err := h.llmProcessor.ProcessWithHistory(
			ctx,
			manager,
			cmd.Content,
			history,
			enhancedContext,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to process with LLM")
		}

		// SAVE conversation with user message
		if err = h.conversations.Save(ctx, conversation); err != nil {
			return nil, errors.Wrap(err, "error saving user message")
		}

		// PUBLISH user message event immediately
		if err = h.publisher.Publish(ctx, event); err != nil {
			return nil, errors.Wrap(err, "error publishing user message event")
		}

		// RELOAD conversation after user message
		conversation, err = h.conversations.Load(ctx, cmd.ConversationID)
		if err != nil {
			return nil, errors.Wrap(err, "error reloading conversation after user message")
		}

		// Add manager response message
		managerMessageID := uuid.New().String()
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
		updatedContext := h.updateContext(enhancedContext, response, actions, confidence)
		contextEvent, err := conversation.UpdateContext(updatedContext)
		if err != nil {
			return nil, errors.Wrap(err, "failed to update context")
		}

		// SAVE context update
		if err = h.conversations.Save(ctx, conversation); err != nil {
			return nil, errors.Wrap(err, "error saving context update")
		}

		// PUBLISH context update event immediately
		if err = h.publisher.Publish(ctx, contextEvent); err != nil {
			return nil, errors.Wrap(err, "error publishing context update event")
		}

		// Update result
		result.ManagerMessageID = managerMessageID
		result.ManagerResponse = response
		result.Actions = actions
		result.Confidence = confidence

		return result, nil
	}

	// Save conversation
	if err = h.conversations.Save(ctx, conversation); err != nil {
		return nil, errors.Wrap(err, "error saving conversation")
	}

	// Publish event (even single events use the same pattern)
	if err = h.publisher.Publish(ctx, event); err != nil {
		return nil, errors.Wrap(err, "failed to publish event")
	}

	return result, nil
}

func (h AddMessageToConversationHandler) loadConversationHistory(ctx context.Context, conversationID, userID string, maxMessages int) ([]domain.ConversationMessage, error) {
	limit := maxMessages
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
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

func (h AddMessageToConversationHandler) buildContext(conversation *domain.Conversation, cmd AddMessageToConversation, history []domain.ConversationMessage) map[string]interface{} {
	context := make(map[string]interface{})

	for k, v := range conversation.Context {
		context[k] = v
	}

	for k, v := range cmd.Metadata {
		context[k] = v
	}

	context["conversation_id"] = cmd.ConversationID
	context["user_id"] = cmd.UserID
	context["manager_id"] = conversation.ManagerID
	context["message_count"] = len(history)
	context["timestamp"] = time.Now().Unix()

	return context
}

func (h AddMessageToConversationHandler) updateContext(baseContext map[string]interface{}, response string, actions []domain.ManagerAction, confidence float64) map[string]interface{} {
	updatedContext := make(map[string]interface{})
	for k, v := range baseContext {
		updatedContext[k] = v
	}

	updatedContext["last_response_confidence"] = confidence
	updatedContext["last_tools_used_count"] = len(actions)
	updatedContext["last_interaction_timestamp"] = time.Now().Unix()

	if confidence > 0.8 {
		updatedContext["last_interaction_quality"] = "high"
	} else if confidence > 0.6 {
		updatedContext["last_interaction_quality"] = "medium"
	} else {
		updatedContext["last_interaction_quality"] = "low"
	}

	if len(actions) > 0 {
		toolTypes := make([]string, 0, len(actions))
		for _, action := range actions {
			if action.Type != "" {
				toolTypes = append(toolTypes, action.Type)
			}
		}
		updatedContext["last_tools_used"] = toolTypes
	}

	return updatedContext
}
