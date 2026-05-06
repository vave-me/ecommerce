package grpc

import (
	"context"
	"log" // Added for logging

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"middleman/internal/auth"
	"middleman/internal/errorsotel"
	"middleman/messages/internal/application"
	"middleman/messages/internal/application/commands"
	"middleman/messages/internal/application/queries"
	"middleman/messages/internal/domain"
	"middleman/messages/messagespb"
)

type server struct {
	app application.App
	messagespb.UnimplementedMessagesServiceServer
}

var _ messagespb.MessagesServiceServer = (*server)(nil)

func RegisterServer(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	log.Println("Registering MessagesServiceServer...")

	messagespb.RegisterMessagesServiceServer(registrar, server{app: app})
	log.Println("MessagesServiceServer registered successfully.")
	return nil
}

func (s server) StartConversation(ctx context.Context, request *messagespb.StartConversationRequest) (*messagespb.StartConversationResponse, error) {
	log.Println("StartConversation called. Request:", request)

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)
	conversationID := uuid.New().String()

	span.SetAttributes(
		attribute.String("conversationID", conversationID),
	)

	err := s.app.StartConversation(ctx, commands.StartConversation{
		ID:          conversationID,
		SenderID:    userID,
		RecipientID: request.GetRecipientId(),
		ItemID:      request.GetItemId(),
	})
	if err != nil {
		log.Println("Error in StartConversation:", err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Println("StartConversation completed successfully. Returning response with conversation ID:", conversationID)
	return &messagespb.StartConversationResponse{
		Id: conversationID,
	}, nil
}

func (s server) GetConversation(ctx context.Context, request *messagespb.GetConversationRequest) (*messagespb.GetConversationResponse, error) {
	log.Println("GetConversation called. Request:", request)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConversationID", request.GetId()),
	)

	conversation, err := s.app.GetConversation(ctx, queries.GetConversation{ID: request.GetId()})
	if err != nil {
		log.Println("Error in GetConversation:", err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Println("GetConversation completed successfully. Returning conversation:", conversation.ID)
	return &messagespb.GetConversationResponse{Conversation: s.conversationFromDomain(conversation)}, nil
}

func (s server) GetConversationByRecipientAndItem(ctx context.Context, request *messagespb.GetConversationByRecipientAndItemRequest) (*messagespb.GetConversationByRecipientAndItemResponse, error) {
	log.Println("GetConversationByRecipientAndItem called. Request:", request)

	span := trace.SpanFromContext(ctx)
	claims, ok := auth.ClaimsFromContext(ctx)

	senderID := claims.Subject
	log.Println("Extracted senderID from context:", senderID)

	if !ok {
		log.Println("Unauthenticated request for GetConversationByRecipientAndItem.")
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	span.SetAttributes(
		attribute.String("RecipientID", request.GetRecipientId()),
	)

	conversation, err := s.app.GetConversationByRecipientAndItem(ctx, queries.GetConversationByRecipientAndItem{
		SenderID:    senderID,
		RecipientID: request.GetRecipientId(),
		ItemID:      request.GetItemId(),
	})
	if err != nil {
		log.Println("Error in GetConversationByRecipientAndItem:", err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Println("GetConversationByRecipientAndItem completed. Returning conversation:", conversation.ID)
	return &messagespb.GetConversationByRecipientAndItemResponse{Conversation: s.conversationFromDomain(conversation)}, nil
}

func (s server) GetConversations(ctx context.Context, request *messagespb.GetConversationsRequest) (*messagespb.GetConversationsResponse, error) {
	log.Println("GetConversations called. Request:", request)
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	senderID := claims.Subject
	span := trace.SpanFromContext(ctx)

	conversations, err := s.app.GetConversations(ctx, queries.GetConversations{
		UserID: senderID,
	})
	if err != nil {
		log.Println("Error in GetConversations:", err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Println("Fetched conversations. Count:", len(conversations))

	protoConversations := []*messagespb.Conversation{}
	for _, conversation := range conversations {
		protoConversations = append(protoConversations, s.conversationFromDomain(conversation))
	}

	log.Println("GetConversations completed successfully. Returning", len(protoConversations), "conversations.")
	return &messagespb.GetConversationsResponse{
		Conversations: protoConversations,
	}, nil
}

func (s server) GetActiveConversations(ctx context.Context, request *messagespb.GetActiveConversationsRequest) (*messagespb.GetActiveConversationsResponse, error) {
	log.Println("GetActiveConversations called. Request:", request)

	span := trace.SpanFromContext(ctx)

	conversations, err := s.app.GetConversations(ctx, queries.GetConversations{})
	if err != nil {
		log.Println("Error in GetActiveConversations:", err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Println("Fetched active conversations. Count:", len(conversations))

	protoConversations := []*messagespb.Conversation{}
	for _, conversation := range conversations {
		protoConversations = append(protoConversations, s.conversationFromDomain(conversation))
	}

	log.Println("GetActiveConversations completed successfully. Returning", len(protoConversations), "conversations.")
	return &messagespb.GetActiveConversationsResponse{
		Conversations: protoConversations,
	}, nil
}

func (s server) SendMessage(ctx context.Context, request *messagespb.SendMessageRequest) (*messagespb.SendMessageResponse, error) {
	log.Println("SendMessage called. Request:", request)
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("MessageID", request.GetId()),
	)

	err := s.app.SendMessage(ctx, commands.SendMessage{
		ID:             request.GetId(),
		ConversationID: request.GetConversationId(),
		SenderID:       userID,
		RecipientID:    request.GetRecipientId(),
		ItemID:         request.GetItemId(),
		Body:           request.GetBody(),
		IsRead:         request.GetIsRead(),
	})
	if err != nil {
		log.Println("Error in SendMessage:", err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Println("SendMessage completed successfully. MessageID:", request.GetId())
	return &messagespb.SendMessageResponse{Id: request.GetId()}, nil
}

func (s server) DeleteMessage(ctx context.Context, request *messagespb.DeleteMessageRequest) (*messagespb.DeleteMessageResponse, error) {
	log.Println("DeleteMessage called. Request:", request)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("MessageID", request.GetId()),
	)

	err := s.app.DeleteMessage(ctx, commands.DeleteMessage{
		ID: request.GetId(),
	})
	if err != nil {
		log.Println("Error in DeleteMessage:", err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Println("DeleteMessage completed successfully. MessageID:", request.GetId())
	return &messagespb.DeleteMessageResponse{}, nil
}

func (s server) GetMessage(ctx context.Context, request *messagespb.GetMessageRequest) (*messagespb.GetMessageResponse, error) {
	log.Println("GetMessage called. Request:", request)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("MessageID", request.GetId()),
	)

	message, err := s.app.GetMessage(ctx, queries.GetMessage{
		MessageID: request.GetId(),
	})
	if err != nil {
		log.Println("Error in GetMessage:", err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Println("GetMessage completed successfully. MessageID:", message.ID)
	return &messagespb.GetMessageResponse{Message: s.messageFromDomain(message)}, nil
}

func (s server) GetMessages(ctx context.Context, request *messagespb.GetMessagesRequest) (*messagespb.GetMessagesResponse, error) {
	log.Println("GetMessages called. Request:", request)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConversationID", request.GetConversationId()),
	)

	messages, err := s.app.GetMessages(ctx, queries.GetMessages{ConversationID: request.GetConversationId()})
	if err != nil {
		log.Println("Error in GetMessages:", err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	log.Println("Fetched messages. Count:", len(messages))

	protoMessages := make([]*messagespb.Message, len(messages))
	for i, message := range messages {
		protoMessages[i] = s.messageFromDomain(message)
	}

	log.Println("GetMessages completed successfully. Returning", len(protoMessages), "messages.")
	return &messagespb.GetMessagesResponse{
		Messages: protoMessages,
	}, nil
}

func (s server) conversationFromDomain(conversation *domain.MiddlemanConversation) *messagespb.Conversation {
	return &messagespb.Conversation{
		Id:          conversation.ID,
		SenderId:    conversation.SenderID,
		RecipientId: conversation.RecipientID,
		ItemId:      conversation.ItemID,
	}
}

func (s server) messageFromDomain(message *domain.MiddlemanMessage) *messagespb.Message {
	return &messagespb.Message{
		Id:             message.ID,
		ConversationId: message.ConversationID,
		SenderId:       message.SenderID,
		RecipientId:    message.RecipientID,
		ItemId:         message.ItemID,
		Body:           message.Body,
	}
}
