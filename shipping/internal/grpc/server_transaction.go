package grpc

import (
	"context"
	"database/sql"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"middleman/internal/di"
	"middleman/shipping/internal/application"
	"middleman/shipping/internal/constants"
	"middleman/shipping/shippingpb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	shippingpb.UnimplementedShippingServiceServer
	s3Client *s3.Client
	bucket   string
}

var _ shippingpb.ShippingServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar,
) error {
	shippingpb.RegisterShippingServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateShipping(ctx context.Context, request *shippingpb.CreateShippingRequest) (resp *shippingpb.CreateShippingResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.CreateShipping(ctx, request)
}

func (s serverTx) CancelShipment(ctx context.Context, request *shippingpb.CancelShipmentRequest) (resp *shippingpb.CancelShipmentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.CancelShipment(ctx, request)
}

func (s serverTx) SchedulePickup(ctx context.Context, request *shippingpb.SchedulePickupRequest) (resp *shippingpb.SchedulePickupResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.SchedulePickup(ctx, request)
}

func (s serverTx) StartShipment(ctx context.Context, request *shippingpb.StartShipmentRequest) (resp *shippingpb.StartShipmentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.StartShipment(ctx, request)
}

func (s serverTx) UpdateShipmentStatus(ctx context.Context, request *shippingpb.UpdateShipmentStatusRequest) (resp *shippingpb.UpdateShipmentStatusResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateShipmentStatus(ctx, request)
}

func (s serverTx) AssignCarrier(ctx context.Context, request *shippingpb.AssignCarrierRequest) (resp *shippingpb.AssignCarrierResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AssignCarrier(ctx, request)
}

func (s serverTx) MarkShipmentAsDelivered(ctx context.Context, request *shippingpb.MarkShipmentAsDeliveredRequest) (resp *shippingpb.MarkShipmentAsDeliveredResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.MarkShipmentAsDelivered(ctx, request)
}

func (s serverTx) ReturnShipment(ctx context.Context, request *shippingpb.ReturnShipmentRequest) (resp *shippingpb.ReturnShipmentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ReturnShipment(ctx, request)
}

func (s serverTx) GetShipment(ctx context.Context, request *shippingpb.GetShipmentRequest) (*shippingpb.GetShipmentResponse, error) {
	// Read-only operation, no transaction needed
	next := server{app: di.Get(s.c.Scoped(ctx), constants.ApplicationKey).(application.App)}
	return next.GetShipment(ctx, request)
}

func (s serverTx) ListShipments(ctx context.Context, request *shippingpb.ListShipmentsRequest) (*shippingpb.ListShipmentsResponse, error) {
	// Read-only operation, no transaction needed
	next := server{app: di.Get(s.c.Scoped(ctx), constants.ApplicationKey).(application.App)}
	return next.ListShipments(ctx, request)
}

func (s serverTx) TrackShipping(ctx context.Context, request *shippingpb.TrackShippingRequest) (*shippingpb.TrackShippingResponse, error) {
	// Read-only operation, no transaction needed
	next := server{app: di.Get(s.c.Scoped(ctx), constants.ApplicationKey).(application.App)}
	return next.TrackShipping(ctx, request)
}

func (s serverTx) GetShipmentHistory(ctx context.Context, request *shippingpb.GetShipmentHistoryRequest) (*shippingpb.GetShipmentHistoryResponse, error) {
	// Read-only operation, no transaction needed
	next := server{app: di.Get(s.c.Scoped(ctx), constants.ApplicationKey).(application.App)}
	return next.GetShipmentHistory(ctx, request)
}

func (s serverTx) GetLabel(ctx context.Context, request *shippingpb.GetLabelRequest) (*shippingpb.GetLabelResponse, error) {
	// Read-only operation, no transaction needed
	next := server{app: di.Get(s.c.Scoped(ctx), constants.ApplicationKey).(application.App)}
	return next.GetLabel(ctx, request)
}

func (s serverTx) GetRates(ctx context.Context, request *shippingpb.GetRatesRequest) (*shippingpb.GetRatesResponse, error) {
	// Read-only operation, no transaction needed
	next := server{app: di.Get(s.c.Scoped(ctx), constants.ApplicationKey).(application.App)}
	return next.GetRates(ctx, request)
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
