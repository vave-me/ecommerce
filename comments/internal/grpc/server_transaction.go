package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/comments/commentspb"
	"middleman/comments/internal/application"
	"middleman/comments/internal/constants"
	"middleman/internal/di"
)

type serverTx struct {
	c di.Container
	commentspb.UnimplementedCommentsServiceServer
}

var _ commentspb.CommentsServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	commentspb.RegisterCommentsServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) AddComment(ctx context.Context, request *commentspb.AddCommentRequest) (resp *commentspb.AddCommentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddComment(ctx, request)
}
func (s serverTx) EditComment(ctx context.Context, request *commentspb.EditCommentRequest) (resp *commentspb.EditCommentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.EditComment(ctx, request)
}
func (s serverTx) ApproveComment(ctx context.Context, request *commentspb.ApproveCommentRequest) (resp *commentspb.ApproveCommentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ApproveComment(ctx, request)
}
func (s serverTx) FlagComment(ctx context.Context, request *commentspb.FlagCommentRequest) (resp *commentspb.FlagCommentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.FlagComment(ctx, request)
}
func (s serverTx) RemoveComment(ctx context.Context, request *commentspb.RemoveCommentRequest) (resp *commentspb.RemoveCommentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RemoveComment(ctx, request)
}
func (s serverTx) RejectComment(ctx context.Context, request *commentspb.RejectCommentRequest) (resp *commentspb.RejectCommentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RejectComment(ctx, request)
}
func (s serverTx) GetComment(ctx context.Context, request *commentspb.GetCommentRequest) (resp *commentspb.GetCommentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetComment(ctx, request)
}
func (s serverTx) GetComments(ctx context.Context, request *commentspb.GetCommentsRequest) (resp *commentspb.GetCommentsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetComments(ctx, request)
}

func (s serverTx) GetMostCommentedItems(ctx context.Context, request *commentspb.GetMostCommentedRequest) (resp *commentspb.GetMostCommentedResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetMostCommented(ctx, request)
}

func (s serverTx) GetMostCommentedItemsByCategory(ctx context.Context, request *commentspb.GetMostCommentedByCategoryRequest) (resp *commentspb.GetMostCommentedByCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetMostCommentedByCategory(ctx, request)
}

func (s serverTx) GetCommentsBySender(ctx context.Context, request *commentspb.GetCommentsBySenderRequest) (resp *commentspb.GetCommentsBySenderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetCommentsBySender(ctx, request)
}
func (s serverTx) GetApprovedComments(ctx context.Context, request *commentspb.GetApprovedCommentsRequest) (resp *commentspb.GetApprovedCommentsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetApprovedComments(ctx, request)
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
