package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/geocoding/geocodingpb"
	"middleman/geocoding/internal/application"
	"middleman/geocoding/internal/constants"
	"middleman/internal/di"
)

type serverTx struct {
	c di.Container
	geocodingpb.UnimplementedGeocodingServiceServer
}

var _ geocodingpb.GeocodingServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	geocodingpb.RegisterGeocodingServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) GeocodeAddress(ctx context.Context, request *geocodingpb.GeocodeAddressRequest) (resp *geocodingpb.GeocodeAddressResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GeocodeAddress(ctx, request)
}
func (s serverTx) SuggestAddress(ctx context.Context, request *geocodingpb.SuggestAddressRequest) (resp *geocodingpb.SuggestAddressResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.SuggestAddress(ctx, request)
}
func (s serverTx) SuggestCity(ctx context.Context, request *geocodingpb.SuggestCityRequest) (resp *geocodingpb.SuggestCityResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.SuggestCity(ctx, request)
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
