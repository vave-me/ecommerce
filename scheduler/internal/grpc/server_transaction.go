package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/internal/di"
	"middleman/scheduler/internal/application"
	"middleman/scheduler/internal/constants"
	"middleman/scheduler/schedulerspb"
)

type serverTx struct {
	c di.Container
	schedulerspb.UnimplementedSchedulerServiceServer
}

var _ schedulerspb.SchedulerServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	schedulerspb.RegisterSchedulerServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}
func (s serverTx) CreateScheduler(ctx context.Context, request *schedulerspb.CreateSchedulerRequest) (resp *schedulerspb.CreateSchedulerResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateScheduler(ctx, request)
}

func (s serverTx) GetScheduler(ctx context.Context, request *schedulerspb.GetSchedulerRequest) (resp *schedulerspb.GetSchedulerResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetScheduler(ctx, request)
}

func (s serverTx) GetSchedulers(ctx context.Context, request *schedulerspb.GetSchedulersRequest) (resp *schedulerspb.GetSchedulersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetSchedulers(ctx, request)
}

func (s serverTx) GetActions(ctx context.Context, request *schedulerspb.GetActionsRequest) (resp *schedulerspb.GetActionsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetActions(ctx, request)
}


func (s serverTx) GetAction(ctx context.Context, request *schedulerspb.GetActionRequest) (resp *schedulerspb.GetActionResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetAction(ctx, request)
}
func (s serverTx) RemoveAction(ctx context.Context, request *schedulerspb.RemoveActionRequest) (resp *schedulerspb.RemoveActionResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RemoveAction(ctx, request)
}
func (s serverTx) AddAction(ctx context.Context, request *schedulerspb.AddActionRequest) (resp *schedulerspb.AddActionResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddAction(ctx, request)
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
