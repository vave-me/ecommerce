package grpc

import (
	"context"
	"database/sql"
	"middleman/internal/di"
	"middleman/support/internal/application"
	"middleman/support/internal/constants"
	"middleman/support/supportpb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	supportpb.UnimplementedSupportServiceServer
}

var _ supportpb.SupportServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	supportpb.RegisterSupportServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateSupportChannel(ctx context.Context, request *supportpb.CreateSupportChannelRequest) (resp *supportpb.CreateSupportChannelResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateSupportChannel(ctx, request)
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
