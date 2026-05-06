package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/activity/activitypb"
	"middleman/activity/internal/application"
	"middleman/activity/internal/constants"
	"middleman/internal/di"
)

type serverTx struct {
	c di.Container
	activitypb.UnimplementedActivityServiceServer
}

var _ activitypb.ActivityServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	activitypb.RegisterActivityServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}
func (s serverTx) CreateActivity(ctx context.Context, request *activitypb.CreateActivityRequest) (resp *activitypb.CreateActivityResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateActivity(ctx, request)
}

func (s serverTx) GetActivity(ctx context.Context, request *activitypb.GetActivityRequest) (resp *activitypb.GetActivityResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetActivity(ctx, request)
}

func (s serverTx) GetActivities(ctx context.Context, request *activitypb.GetActivitiesRequest) (resp *activitypb.GetActivitiesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetActivities(ctx, request)
}

func (s serverTx) GetInteractions(ctx context.Context, request *activitypb.GetInteractionsRequest) (resp *activitypb.GetInteractionsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetInteractions(ctx, request)
}

func (s serverTx) GetMostLiked(ctx context.Context, request *activitypb.GetMostLikedRequest) (resp *activitypb.GetMostLikedResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetMostLiked(ctx, request)
}

func (s serverTx) GetMostDisliked(ctx context.Context, request *activitypb.GetMostDislikedRequest) (resp *activitypb.GetMostDislikedResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetMostDisliked(ctx, request)
}

func (s serverTx) GetInteraction(ctx context.Context, request *activitypb.GetInteractionRequest) (resp *activitypb.GetInteractionResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetInteraction(ctx, request)
}
func (s serverTx) RemoveInteraction(ctx context.Context, request *activitypb.RemoveInteractionRequest) (resp *activitypb.RemoveInteractionResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RemoveInteraction(ctx, request)
}
func (s serverTx) AddInteraction(ctx context.Context, request *activitypb.AddInteractionRequest) (resp *activitypb.AddInteractionResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddInteraction(ctx, request)
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
