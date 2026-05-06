package grpc

import (
	"context"
	"database/sql"
	"middleman/assistants/assistantspb"
	"middleman/assistants/internal/application"
	"middleman/assistants/internal/constants"
	"middleman/assistants/internal/domain"
	"middleman/internal/di"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	assistantspb.UnimplementedAssistantsServiceServer
}

var _ assistantspb.AssistantsServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	assistantspb.RegisterAssistantsServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

// GetAssistant retrieves an assistant by ID (with transaction)
func (s serverTx) GetAssistant(ctx context.Context, request *assistantspb.GetAssistantRequest) (resp *assistantspb.GetAssistantResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetAssistant(ctx, request)
}

// ActivateAssistant activates an assistant (with transaction)
func (s serverTx) ActivateAssistant(ctx context.Context, request *assistantspb.ActivateAssistantRequest) (resp *assistantspb.ActivateAssistantResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ActivateAssistant(ctx, request)
}

// DeactivateAssistant deactivates an assistant (with transaction)
func (s serverTx) DeactivateAssistant(ctx context.Context, request *assistantspb.DeactivateAssistantRequest) (resp *assistantspb.DeactivateAssistantResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.DeactivateAssistant(ctx, request)
}

// UpdateAssistantConfiguration updates assistant configuration (with transaction)
func (s serverTx) UpdateAssistantConfiguration(ctx context.Context, request *assistantspb.UpdateAssistantConfigurationRequest) (resp *assistantspb.UpdateAssistantConfigurationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.UpdateAssistantConfiguration(ctx, request)
}

// GetAssistants retrieves all assistants (with transaction)
func (s serverTx) GetAssistants(ctx context.Context, request *assistantspb.GetAssistantsRequest) (resp *assistantspb.GetAssistantsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetAssistants(ctx, request)
}

// CreateConversation creates a new conversation (with transaction)
func (s serverTx) CreateConversation(ctx context.Context, request *assistantspb.CreateConversationRequest) (resp *assistantspb.CreateConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateConversation(ctx, request)
}

// GetConversation retrieves a conversation by ID (with transaction)
func (s serverTx) GetConversation(ctx context.Context, request *assistantspb.GetConversationRequest) (resp *assistantspb.GetConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetConversation(ctx, request)
}

// GetUserConversations retrieves conversations for a user (with transaction)
func (s serverTx) GetUserConversations(ctx context.Context, request *assistantspb.GetUserConversationsRequest) (resp *assistantspb.GetUserConversationsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetUserConversations(ctx, request)
}

// GetConversationMessages retrieves messages for a conversation (with transaction)
func (s serverTx) GetConversationMessages(ctx context.Context, request *assistantspb.GetConversationMessagesRequest) (resp *assistantspb.GetConversationMessagesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetConversationMessages(ctx, request)
}

// GetConversationStats retrieves conversation statistics for a user (with transaction)
func (s serverTx) GetConversationStats(ctx context.Context, request *assistantspb.GetConversationStatsRequest) (resp *assistantspb.GetConversationStatsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetConversationStats(ctx, request)
}

// AddMessageToConversation adds a message to an existing conversation (with transaction)
func (s serverTx) AddMessageToConversation(ctx context.Context, request *assistantspb.AddMessageToConversationRequest) (resp *assistantspb.AddMessageToConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddMessageToConversation(ctx, request)
}

// ChatWithConversation processes a message within an existing conversation context (with transaction)
func (s serverTx) ChatWithConversation(ctx context.Context, request *assistantspb.ChatWithConversationRequest) (resp *assistantspb.ChatWithConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ChatWithConversation(ctx, request)
}

// UpdateConversationContext updates the context of a conversation (with transaction)
func (s serverTx) UpdateConversationContext(ctx context.Context, request *assistantspb.UpdateConversationContextRequest) (resp *assistantspb.UpdateConversationContextResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.UpdateConversationContext(ctx, request)
}

// UpdateConversation updates conversation metadata (with transaction)
func (s serverTx) UpdateConversation(ctx context.Context, request *assistantspb.UpdateConversationRequest) (resp *assistantspb.UpdateConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.UpdateConversation(ctx, request)
}

// DeleteConversation deletes a conversation (with transaction)
func (s serverTx) DeleteConversation(ctx context.Context, request *assistantspb.DeleteConversationRequest) (resp *assistantspb.DeleteConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.DeleteConversation(ctx, request)
}

// ArchiveConversation archives a conversation (with transaction)
func (s serverTx) ArchiveConversation(ctx context.Context, request *assistantspb.ArchiveConversationRequest) (resp *assistantspb.ArchiveConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ArchiveConversation(ctx, request)
}

// Helper methods - delegate to the main server instance
// These don't need transactions since they're just data conversion

func (s serverTx) assistantFromDomain(assistant *domain.CatalogAssistant) *assistantspb.Assistant {
	// Create a temporary server instance to access helper methods
	next := server{app: di.Get(context.Background(), constants.ApplicationKey).(application.App)}
	return next.assistantFromDomain(assistant)
}

func (s serverTx) catalogAssistantFromDomain(assistant *domain.CatalogAssistant) *assistantspb.Assistant {
	// Create a temporary server instance to access helper methods
	next := server{app: di.Get(context.Background(), constants.ApplicationKey).(application.App)}
	return next.catalogAssistantFromDomain(assistant)
}

func (s serverTx) conversationViewToProto(conv *domain.ReadConversation) *assistantspb.Conversation {
	// Create a temporary server instance to access helper methods
	next := server{app: di.Get(context.Background(), constants.ApplicationKey).(application.App)}
	return next.conversationViewToProto(conv)
}

func (s serverTx) conversationMessageToProto(msg *domain.ReadMessage) *assistantspb.ConversationMessage {
	// Create a temporary server instance to access helper methods
	next := server{app: di.Get(context.Background(), constants.ApplicationKey).(application.App)}
	return next.conversationMessageToProto(msg)
}

func (s serverTx) convertParametersToStringMap(parameters map[string]interface{}) map[string]string {
	// Create a temporary server instance to access helper methods
	next := server{app: di.Get(context.Background(), constants.ApplicationKey).(application.App)}
	return next.convertParametersToStringMap(parameters)
}

// ProcessUserInput processes user input through an assistant (with transaction)
func (s serverTx) ProcessUserInput(ctx context.Context, request *assistantspb.ProcessUserInputRequest) (resp *assistantspb.ProcessUserInputResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ProcessUserInput(ctx, request)
}

// ProcessSpeechInput processes speech input through an assistant (with transaction)
func (s serverTx) ProcessSpeechInput(ctx context.Context, request *assistantspb.ProcessSpeechInputRequest) (resp *assistantspb.ProcessSpeechInputResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ProcessSpeechInput(ctx, request)
}

// ProcessImageInput processes image input through an assistant (with transaction)
func (s serverTx) ProcessImageInput(ctx context.Context, request *assistantspb.ProcessImageInputRequest) (resp *assistantspb.ProcessImageInputResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ProcessImageInput(ctx, request)
}

// ProcessDocumentInput processes document input through an assistant (with transaction)
func (s serverTx) ProcessDocumentInput(ctx context.Context, request *assistantspb.ProcessDocumentInputRequest) (resp *assistantspb.ProcessDocumentInputResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ProcessDocumentInput(ctx, request)
}

// ProcessAssistantRequest handles assistant request processing (with transaction)
func (s serverTx) ProcessAssistantRequest(ctx context.Context, request *assistantspb.ProcessAssistantRequestRequest) (resp *assistantspb.ProcessAssistantRequestResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ProcessAssistantRequest(ctx, request)
}

func (s serverTx) closeTx(tx *sql.Tx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}
