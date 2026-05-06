package grpc

import (
	"context"
	"database/sql"
	"middleman/internal/di"
	"middleman/notifications/internal/application"
	"middleman/notifications/internal/constants"
	"middleman/notifications/notificationspb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	notificationspb.UnimplementedNotificationsServiceServer
}

var _ notificationspb.NotificationsServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	notificationspb.RegisterNotificationsServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) ListAlerts(ctx context.Context, request *notificationspb.ListAlertsRequest) (resp *notificationspb.ListAlertsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ListAlerts(ctx, request)
}
func (s serverTx) GetAlertsByType(ctx context.Context, request *notificationspb.GetAlertsByTypeRequest) (resp *notificationspb.GetAlertsByTypeResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetAlertsByType(ctx, request)
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
