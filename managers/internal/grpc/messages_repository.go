package grpc

import (
	"context"
	"fmt"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
	"middleman/messages/messagespb"
	"time"

	"google.golang.org/grpc"
)

// MessagesRepository calls the remote messages service (gRPC) as a fallback.
type MessagesRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.MessagesRepository = (*MessagesRepository)(nil)

// NewMessagesRepositoryWithAuth creates a new MessagesRepository with JWT authentication support
func NewMessagesRepository(endpoint string, authInstance *auth.Auth) MessagesRepository {
	return MessagesRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// Conversation Management Methods

// StartConversation starts a new conversation between a sender and a recipient regarding a specific item
func (r MessagesRepository) StartConversation(ctx context.Context, senderID, recipientID, itemID string) (*models.ConversationResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := messagespb.NewMessagesServiceClient(conn)
	resp, err := client.StartConversation(ctx, &messagespb.StartConversationRequest{
		RecipientId: recipientID,
		ItemId:      itemID,
	})
	if err != nil {
		return nil, fmt.Errorf("StartConversation RPC failed: %w", err)
	}

	return &models.ConversationResponse{
		ID: resp.GetId(),
	}, nil
}

// RestoreConversation restores an archived conversation by its unique identifier
func (r MessagesRepository) RestoreConversation(ctx context.Context, conversationID string) (*models.ConversationStatusResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := messagespb.NewMessagesServiceClient(conn)
	resp, err := client.RestoreConversation(ctx, &messagespb.RestoreConversationRequest{
		Id: conversationID,
	})
	if err != nil {
		return nil, fmt.Errorf("RestoreConversation RPC failed: %w", err)
	}

	return &models.ConversationStatusResponse{
		ID:                 resp.GetId(),
		ConversationStatus: resp.GetConversationStatus(),
	}, nil
}

// ArchiveConversation archives an active conversation by its unique identifier
func (r MessagesRepository) ArchiveConversation(ctx context.Context, conversationID string) (*models.ConversationStatusResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := messagespb.NewMessagesServiceClient(conn)
	resp, err := client.ArchiveConversation(ctx, &messagespb.ArchiveConversationRequest{
		Id: conversationID,
	})
	if err != nil {
		return nil, fmt.Errorf("ArchiveConversation RPC failed: %w", err)
	}

	return &models.ConversationStatusResponse{
		ID:                 resp.GetId(),
		ConversationStatus: resp.GetConversationStatus(),
	}, nil
}

// GetConversation retrieves details of a specific conversation by its unique identifier
func (r MessagesRepository) GetConversation(ctx context.Context, conversationID string) (*models.Conversation, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := messagespb.NewMessagesServiceClient(conn)
	resp, err := client.GetConversation(ctx, &messagespb.GetConversationRequest{
		Id: conversationID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetConversation RPC failed: %w", err)
	}

	return r.convertConversationFromPb(resp.GetConversation()), nil
}

// GetConversationByRecipientAndItem finds a conversation between users for a specific item
func (r MessagesRepository) GetConversationByRecipientAndItem(ctx context.Context, recipientID, itemID string) (*models.Conversation, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := messagespb.NewMessagesServiceClient(conn)
	resp, err := client.GetConversationByRecipientAndItem(ctx, &messagespb.GetConversationByRecipientAndItemRequest{
		RecipientId: recipientID,
		ItemId:      itemID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetConversationByRecipientAndItem RPC failed: %w", err)
	}

	return r.convertConversationFromPb(resp.GetConversation()), nil
}

// GetConversations retrieves a list of all conversations associated with a specific user
func (r MessagesRepository) GetConversations(ctx context.Context, userID string, page, limit int64) (*models.ConversationsResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := messagespb.NewMessagesServiceClient(conn)
	resp, err := client.GetConversations(ctx, &messagespb.GetConversationsRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return nil, fmt.Errorf("GetConversations RPC failed: %w", err)
	}

	conversations := make([]*models.Conversation, len(resp.GetConversations()))
	for i, pbConv := range resp.GetConversations() {
		conversations[i] = r.convertConversationFromPb(pbConv)
	}

	return &models.ConversationsResponse{
		Conversations: conversations,
		Total:         resp.GetTotal(),
		Page:          resp.GetPage(),
		Limit:         resp.GetLimit(),
	}, nil
}

// GetActiveConversations retrieves a list of active conversations for a specific user
func (r MessagesRepository) GetActiveConversations(ctx context.Context, userID string, page, limit int64) (*models.ConversationsResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := messagespb.NewMessagesServiceClient(conn)
	resp, err := client.GetActiveConversations(ctx, &messagespb.GetActiveConversationsRequest{
		UserId: userID,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("GetActiveConversations RPC failed: %w", err)
	}

	conversations := make([]*models.Conversation, len(resp.GetConversations()))
	for i, pbConv := range resp.GetConversations() {
		conversations[i] = r.convertConversationFromPb(pbConv)
	}

	return &models.ConversationsResponse{
		Conversations: conversations,
		Total:         resp.GetTotal(),
		Page:          resp.GetPage(),
		Limit:         resp.GetLimit(),
	}, nil
}

// Message Management Methods

// SendMessage sends a message via the messages microservice
func (r MessagesRepository) SendMessage(ctx context.Context, messageID, conversationID, senderID, recipientID, itemID, body string, isRead bool) (*models.MessageSentResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := messagespb.NewMessagesServiceClient(conn)
	resp, err := client.SendMessage(ctx, &messagespb.SendMessageRequest{
		Id:             messageID,
		ConversationId: conversationID,
		RecipientId:    recipientID,
		ItemId:         itemID,
		Body:           body,
		IsRead:         isRead,
	})
	if err != nil {
		return nil, fmt.Errorf("SendMessage RPC failed: %w", err)
	}

	sentAt := time.Now()
	if resp.GetSentAt() != nil {
		sentAt = resp.GetSentAt().AsTime()
	}

	return &models.MessageSentResponse{
		ID:     resp.GetId(),
		SentAt: sentAt,
	}, nil
}

// DeleteMessage deletes a message by ID
func (r MessagesRepository) DeleteMessage(ctx context.Context, messageID string) error {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := messagespb.NewMessagesServiceClient(conn)
	_, err = client.DeleteMessage(ctx, &messagespb.DeleteMessageRequest{
		Id: messageID,
	})
	if err != nil {
		return fmt.Errorf("DeleteMessage RPC failed: %w", err)
	}

	return nil
}

// GetMessage retrieves a message by ID from the messages microservice
func (r MessagesRepository) GetMessage(ctx context.Context, messageID string) (*models.Message, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := messagespb.NewMessagesServiceClient(conn)
	resp, err := client.GetMessage(ctx, &messagespb.GetMessageRequest{
		Id: messageID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetMessage RPC failed: %w", err)
	}

	return r.convertMessageFromPb(resp.GetMessage()), nil
}

// GetMessages retrieves a list of messages within a specific conversation
func (r MessagesRepository) GetMessages(ctx context.Context, conversationID string, page, limit int64) (*models.MessagesResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := messagespb.NewMessagesServiceClient(conn)
	resp, err := client.GetMessages(ctx, &messagespb.GetMessagesRequest{
		ConversationId: conversationID,
		Page:           page,
		Limit:          limit,
	})
	if err != nil {
		return nil, fmt.Errorf("GetMessages RPC failed: %w", err)
	}

	messages := make([]*models.Message, len(resp.GetMessages()))
	for i, pbMsg := range resp.GetMessages() {
		messages[i] = r.convertMessageFromPb(pbMsg)
	}

	return &models.MessagesResponse{
		Messages: messages,
		Total:    resp.GetTotal(),
		Page:     resp.GetPage(),
		Limit:    resp.GetLimit(),
	}, nil
}

// Legacy methods simplified - use standard methods with appropriate pagination

// Helper conversion methods

// convertConversationFromPb converts protobuf Conversation to domain model
func (r MessagesRepository) convertConversationFromPb(pbConv *messagespb.Conversation) *models.Conversation {
	if pbConv == nil {
		return nil
	}

	return &models.Conversation{
		ID:                 pbConv.GetId(),
		SenderID:           pbConv.GetSenderId(),
		RecipientID:        pbConv.GetRecipientId(),
		ItemID:             pbConv.GetItemId(),
		ConversationStatus: pbConv.GetConversationStatus(),
		Active:             pbConv.GetConversationStatus() == "active", // Backward compatibility
		CreatedAt:          time.Now(),                                 // Would be set from protobuf timestamp if available
		UpdatedAt:          time.Now(),
	}
}

// convertMessageFromPb converts protobuf Message to domain model
func (r MessagesRepository) convertMessageFromPb(pbMsg *messagespb.Message) *models.Message {
	if pbMsg == nil {
		return nil
	}

	return &models.Message{
		ID:             pbMsg.GetId(),
		ConversationID: pbMsg.GetConversationId(),
		SenderID:       pbMsg.GetSenderId(),
		RecipientID:    pbMsg.GetRecipientId(),
		ItemID:         pbMsg.GetItemId(),
		Body:           pbMsg.GetBody(),
		IsRead:         pbMsg.GetIsRead(),
		CreatedAt:      time.Now(), // Would be set from protobuf timestamp if available
		UpdatedAt:      time.Now(),
	}
}

// dial establishes a gRPC connection to the messages service
// dial sets up a gRPC connection with the microservice endpoint
func (r MessagesRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r MessagesRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}

// Additional methods needed by tool service

// MarkMessageAsRead marks a message as read for a specific user
func (r MessagesRepository) MarkMessageAsRead(ctx context.Context, messageID, userID string) error {
	// For now, return nil as this would need to be implemented in the messages service
	// This is a placeholder implementation
	return nil
}

// GetUnreadMessagesCount returns the count of unread messages for a user
func (r MessagesRepository) GetUnreadMessagesCount(ctx context.Context, userID string) (int64, error) {
	// For now, return 0 as this would need to be implemented in the messages service
	// This is a placeholder implementation
	return 0, nil
}
