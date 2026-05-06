package commands

import (
	"context"
	"github.com/google/uuid"
	"github.com/stackus/errors"
	"middleman/assistants/internal/application/services"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"
	"time"
)

type AddMessageToConversation struct {
	ConversationID     string                   `json:"conversation_id"`
	AssistantID        string                   `json:"assistant_id"`
	MessageID          string                   `json:"message_id"`
	Role               domain.MessageRole       `json:"role"`
	Content            string                   `json:"content"`
	Metadata           map[string]interface{}   `json:"metadata,omitempty"`
	ActionsTaken       []domain.AssistantAction `json:"actions_taken,omitempty"`
	UserID             string                   `json:"user_id,omitempty"`
	ProcessWithLLM     bool                     `json:"process_with_llm,omitempty"`
	MaxHistoryMessages int                      `json:"max_history_messages,omitempty"`
}

type AddMessageToConversationResult struct {
	MessageID          string                   `json:"message_id"`
	AssistantMessageID string                   `json:"assistant_message_id,omitempty"`
	AssistantResponse  string                   `json:"assistant_response,omitempty"`
	Actions            []domain.AssistantAction `json:"actions,omitempty"`
	Confidence         float64                  `json:"confidence,omitempty"`
}

type AddMessageToConversationHandler struct {
	conversations domain.ConversationRepository
	readMessages  domain.ReadMessagesRepository
	assistants    domain.AssistantRepository
	llmProcessor  services.LLMProcessor
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewAddMessageToConversationHandler(
	conversations domain.ConversationRepository,
	readMessages domain.ReadMessagesRepository,
	assistants domain.AssistantRepository,
	llmProcessor services.LLMProcessor,
	publisher ddd.EventPublisher[ddd.Event],
) AddMessageToConversationHandler {
	return AddMessageToConversationHandler{
		conversations: conversations,
		readMessages:  readMessages,
		assistants:    assistants,
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
		cmd.AssistantID,
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
		// Load assistant
		if conversation.AssistantID == "" {
			return nil, errors.Wrap(errors.ErrNotFound, "conversation has no assistant")
		}

		assistant, err := h.assistants.Load(ctx, conversation.AssistantID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load assistant")
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
			assistant,
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

		// Add assistant response message
		assistantMessageID := uuid.New().String()
		assistantEvent, err := conversation.AddMessage(
			cmd.ConversationID,
			cmd.AssistantID,
			assistantMessageID,
			domain.AssistantRole,
			response,
			map[string]interface{}{
				"timestamp":  time.Now().Unix(),
				"confidence": confidence,
				"actions":    len(actions),
			},
			actions,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to add assistant message")
		}

		// SAVE assistant message
		if err = h.conversations.Save(ctx, conversation); err != nil {
			return nil, errors.Wrap(err, "error saving assistant message")
		}

		// PUBLISH assistant message event immediately
		if err = h.publisher.Publish(ctx, assistantEvent); err != nil {
			return nil, errors.Wrap(err, "error publishing assistant message event")
		}

		// RELOAD conversation after assistant message
		conversation, err = h.conversations.Load(ctx, cmd.ConversationID)
		if err != nil {
			return nil, errors.Wrap(err, "error reloading conversation after assistant message")
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
		result.AssistantMessageID = assistantMessageID
		result.AssistantResponse = response
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
	context["assistant_id"] = conversation.AssistantID
	context["message_count"] = len(history)
	context["timestamp"] = time.Now().Unix()

	return context
}

func (h AddMessageToConversationHandler) updateContext(baseContext map[string]interface{}, response string, actions []domain.AssistantAction, confidence float64) map[string]interface{} {
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
