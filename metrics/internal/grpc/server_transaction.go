package grpc

import (
	"context"
	"database/sql"
	"middleman/internal/di"
	"middleman/metrics/internal/application"
	"middleman/metrics/internal/constants"
	"middleman/metrics/metricspb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	metricspb.UnimplementedMetricsServiceServer
}

var _ metricspb.MetricsServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	metricspb.RegisterMetricsServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) GetItemMetric(ctx context.Context, request *metricspb.GetItemMetricRequest) (resp *metricspb.GetItemMetricResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetItemMetric(ctx, request)
}

func (s serverTx) GetItemsMetric(ctx context.Context, request *metricspb.GetItemsMetricRequest) (resp *metricspb.GetItemsMetricResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetItemsMetric(ctx, request)
}

func (s serverTx) GetUserMetric(ctx context.Context, request *metricspb.GetUserMetricRequest) (resp *metricspb.GetUserMetricResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetUserMetric(ctx, request)
}

func (s serverTx) GetHighestMetricsByType(ctx context.Context, request *metricspb.GetHighestMetricsByTypeRequest) (resp *metricspb.GetItemsMetricResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetHighestMetricsByType(ctx, request)
}

func (s serverTx) GetLowestMetricsByType(ctx context.Context, request *metricspb.GetLowestMetricsByTypeRequest) (resp *metricspb.GetItemsMetricResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetLowestMetricsByType(ctx, request)
}
func (s serverTx) ShareItem(ctx context.Context, request *metricspb.ShareItemRequest) (resp *metricspb.ShareItemResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.ShareItem(ctx, request)
}
func (s serverTx) VisitItem(ctx context.Context, request *metricspb.VisitItemRequest) (resp *metricspb.VisitItemResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.VisitItem(ctx, request)
}

func (s serverTx) UpdateItemMetric(ctx context.Context, request *metricspb.UpdateItemMetricRequest) (resp *metricspb.UpdateItemMetricResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.UpdateItemMetric(ctx, request)
}

func (s serverTx) UpdateUserMetric(ctx context.Context, request *metricspb.UpdateUserMetricRequest) (resp *metricspb.UpdateUserMetricResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.UpdateUserMetric(ctx, request)
}

func (s serverTx) closeTx(tx *sql.Tx, err error) error {
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
