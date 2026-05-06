package grpc

import (
	"context"
	"database/sql"
	"middleman/internal/di"
	"middleman/messages/internal/application"
	"middleman/messages/internal/constants"
	"middleman/messages/messagespb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	messagespb.UnimplementedMessagesServiceServer
}

var _ messagespb.MessagesServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	messagespb.RegisterMessagesServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) StartConversation(ctx context.Context, request *messagespb.StartConversationRequest) (resp *messagespb.StartConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.StartConversation(ctx, request)
}

func (s serverTx) RestoreConversation(ctx context.Context, request *messagespb.RestoreConversationRequest) (resp *messagespb.RestoreConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RestoreConversation(ctx, request)
}

func (s serverTx) ArchiveConversation(ctx context.Context, request *messagespb.ArchiveConversationRequest) (resp *messagespb.ArchiveConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ArchiveConversation(ctx, request)
}

func (s serverTx) GetConversation(ctx context.Context, request *messagespb.GetConversationRequest) (resp *messagespb.GetConversationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetConversation(ctx, request)
}
func (s serverTx) GetConversationByRecipientAndItem(ctx context.Context, request *messagespb.GetConversationByRecipientAndItemRequest) (resp *messagespb.GetConversationByRecipientAndItemResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetConversationByRecipientAndItem(ctx, request)
}
func (s serverTx) GetConversations(ctx context.Context, request *messagespb.GetConversationsRequest) (resp *messagespb.GetConversationsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetConversations(ctx, request)
}

func (s serverTx) GetActiveConversations(ctx context.Context, request *messagespb.GetActiveConversationsRequest) (resp *messagespb.GetActiveConversationsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetActiveConversations(ctx, request)
}

func (s serverTx) SendMessage(ctx context.Context, request *messagespb.SendMessageRequest) (resp *messagespb.SendMessageResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.SendMessage(ctx, request)
}

func (s serverTx) DeleteMessage(ctx context.Context, request *messagespb.DeleteMessageRequest) (resp *messagespb.DeleteMessageResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.DeleteMessage(ctx, request)
}

func (s serverTx) GetMessage(ctx context.Context, request *messagespb.GetMessageRequest) (resp *messagespb.GetMessageResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetMessage(ctx, request)
}

func (s serverTx) GetMessages(ctx context.Context, request *messagespb.GetMessagesRequest) (resp *messagespb.GetMessagesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetMessages(ctx, request)
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
