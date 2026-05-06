package grpc

import (
	"context"
	"database/sql"
	"middleman/internal/di"
	"middleman/managers/internal/application"
	"middleman/managers/internal/constants"
	"middleman/managers/internal/domain"
	"middleman/managers/managerspb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	managerspb.UnimplementedManagersServiceServer
}

var _ managerspb.ManagersServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	managerspb.RegisterManagersServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

// CreateManager creates a new manager (with transaction)
func (s serverTx) CreateManager(ctx context.Context, request *managerspb.CreateManagerRequest) (resp *managerspb.CreateManagerResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateManager(ctx, request)
}

// GetManager retrieves an manager by ID (with transaction)
func (s serverTx) GetManager(ctx context.Context, request *managerspb.GetManagerRequest) (resp *managerspb.GetManagerResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetManager(ctx, request)
}

// ActivateManager activates an manager (with transaction)
func (s serverTx) ActivateManager(ctx context.Context, request *managerspb.ActivateManagerRequest) (resp *managerspb.ActivateManagerResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ActivateManager(ctx, request)
}

// DeactivateManager deactivates an manager (with transaction)
func (s serverTx) DeactivateManager(ctx context.Context, request *managerspb.DeactivateManagerRequest) (resp *managerspb.DeactivateManagerResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.DeactivateManager(ctx, request)
}

// UpdateManagerConfiguration updates manager configuration (with transaction)
func (s serverTx) UpdateManagerConfiguration(ctx context.Context, request *managerspb.UpdateManagerConfigurationRequest) (resp *managerspb.UpdateManagerConfigurationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.UpdateManagerConfiguration(ctx, request)
}

// GetManagers retrieves all managers (with transaction)
func (s serverTx) GetManagers(ctx context.Context, request *managerspb.GetManagersRequest) (resp *managerspb.GetManagersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetManagers(ctx, request)
}

// CreateConversation creates a new conversation (with transaction)
func (s serverTx) CreateConversation(ctx context.Context, request *managerspb.CreateConversationRequest) (resp *managerspb.CreateConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateConversation(ctx, request)
}

// GetConversation retrieves a conversation by ID (with transaction)
func (s serverTx) GetConversation(ctx context.Context, request *managerspb.GetConversationRequest) (resp *managerspb.GetConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetConversation(ctx, request)
}

// GetUserConversations retrieves conversations for a user (with transaction)
func (s serverTx) GetUserConversations(ctx context.Context, request *managerspb.GetUserConversationsRequest) (resp *managerspb.GetUserConversationsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetUserConversations(ctx, request)
}

// GetConversationMessages retrieves messages for a conversation (with transaction)
func (s serverTx) GetConversationMessages(ctx context.Context, request *managerspb.GetConversationMessagesRequest) (resp *managerspb.GetConversationMessagesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetConversationMessages(ctx, request)
}

// GetConversationStats retrieves conversation statistics for a user (with transaction)
func (s serverTx) GetConversationStats(ctx context.Context, request *managerspb.GetConversationStatsRequest) (resp *managerspb.GetConversationStatsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetConversationStats(ctx, request)
}

// AddMessageToConversation adds a message to an existing conversation (with transaction)
func (s serverTx) AddMessageToConversation(ctx context.Context, request *managerspb.AddMessageToConversationRequest) (resp *managerspb.AddMessageToConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddMessageToConversation(ctx, request)
}

// ChatWithConversation processes a message within an existing conversation context (with transaction)
func (s serverTx) ChatWithConversation(ctx context.Context, request *managerspb.ChatWithConversationRequest) (resp *managerspb.ChatWithConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ChatWithConversation(ctx, request)
}

// UpdateConversationContext updates the context of a conversation (with transaction)
func (s serverTx) UpdateConversationContext(ctx context.Context, request *managerspb.UpdateConversationContextRequest) (resp *managerspb.UpdateConversationContextResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.UpdateConversationContext(ctx, request)
}

// UpdateConversation updates conversation metadata (with transaction)
func (s serverTx) UpdateConversation(ctx context.Context, request *managerspb.UpdateConversationRequest) (resp *managerspb.UpdateConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.UpdateConversation(ctx, request)
}

// DeleteConversation deletes a conversation (with transaction)
func (s serverTx) DeleteConversation(ctx context.Context, request *managerspb.DeleteConversationRequest) (resp *managerspb.DeleteConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.DeleteConversation(ctx, request)
}

// ArchiveConversation archives a conversation (with transaction)
func (s serverTx) ArchiveConversation(ctx context.Context, request *managerspb.ArchiveConversationRequest) (resp *managerspb.ArchiveConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ArchiveConversation(ctx, request)
}

// Helper methods - delegate to the main server instance
// These don't need transactions since they're just data conversion

func (s serverTx) managerFromDomain(manager *domain.CatalogManager) *managerspb.Manager {
	// Create a temporary server instance to access helper methods
	next := server{app: di.Get(context.Background(), constants.ApplicationKey).(application.App)}
	return next.managerFromDomain(manager)
}

func (s serverTx) catalogManagerFromDomain(manager *domain.CatalogManager) *managerspb.Manager {
	// Create a temporary server instance to access helper methods
	next := server{app: di.Get(context.Background(), constants.ApplicationKey).(application.App)}
	return next.catalogManagerFromDomain(manager)
}

func (s serverTx) conversationViewToProto(conv *domain.ReadConversation) *managerspb.Conversation {
	// Create a temporary server instance to access helper methods
	next := server{app: di.Get(context.Background(), constants.ApplicationKey).(application.App)}
	return next.conversationViewToProto(conv)
}

func (s serverTx) conversationMessageToProto(msg *domain.ReadMessage) *managerspb.ConversationMessage {
	// Create a temporary server instance to access helper methods
	next := server{app: di.Get(context.Background(), constants.ApplicationKey).(application.App)}
	return next.conversationMessageToProto(msg)
}

func (s serverTx) convertParametersToStringMap(parameters map[string]interface{}) map[string]string {
	// Create a temporary server instance to access helper methods
	next := server{app: di.Get(context.Background(), constants.ApplicationKey).(application.App)}
	return next.convertParametersToStringMap(parameters)
}

// ProcessUserInput processes user input through an manager (with transaction)
func (s serverTx) ProcessUserInput(ctx context.Context, request *managerspb.ProcessUserInputRequest) (resp *managerspb.ProcessUserInputResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ProcessUserInput(ctx, request)
}

// ProcessSpeechInput processes speech input through an manager (with transaction)
func (s serverTx) ProcessSpeechInput(ctx context.Context, request *managerspb.ProcessSpeechInputRequest) (resp *managerspb.ProcessSpeechInputResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ProcessSpeechInput(ctx, request)
}

// ProcessImageInput processes image input through an manager (with transaction)
func (s serverTx) ProcessImageInput(ctx context.Context, request *managerspb.ProcessImageInputRequest) (resp *managerspb.ProcessImageInputResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ProcessImageInput(ctx, request)
}

// ProcessDocumentInput processes document input through an manager (with transaction)
func (s serverTx) ProcessDocumentInput(ctx context.Context, request *managerspb.ProcessDocumentInputRequest) (resp *managerspb.ProcessDocumentInputResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ProcessDocumentInput(ctx, request)
}

// ProcessManagerRequest handles manager request processing (with transaction)
func (s serverTx) ProcessManagerRequest(ctx context.Context, request *managerspb.ProcessManagerRequestRequest) (resp *managerspb.ProcessManagerRequestResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ProcessManagerRequest(ctx, request)
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
