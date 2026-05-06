package grpc

import (
	"context"
	"middleman/internal/di"
	"middleman/vectors/internal/application"
	"middleman/vectors/internal/constants"
	"middleman/vectors/vectorspb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	vectorspb.UnimplementedVectorServiceServer
}

var _ vectorspb.VectorServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	vectorspb.RegisterVectorServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) SearchByVector(ctx context.Context, request *vectorspb.SearchByVectorRequest) (resp *vectorspb.SearchByVectorResponse, err error) {
	ctx = s.c.Scoped(ctx)

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchByVector(ctx, request)
}

func (s serverTx) SearchSimilarEntities(ctx context.Context, request *vectorspb.SearchSimilarEntitiesRequest) (resp *vectorspb.SearchSimilarEntitiesResponse, err error) {
	ctx = s.c.Scoped(ctx)

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchSimilarEntities(ctx, request)
}

func (s serverTx) GetEntityContext(ctx context.Context, request *vectorspb.GetEntityContextRequest) (resp *vectorspb.GetEntityContextResponse, err error) {
	ctx = s.c.Scoped(ctx)

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetEntityContext(ctx, request)
}

func (s serverTx) GetRecommendations(ctx context.Context, request *vectorspb.GetRecommendationsRequest) (resp *vectorspb.GetRecommendationsResponse, err error) {
	ctx = s.c.Scoped(ctx)

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetRecommendations(ctx, request)
}

func (s serverTx) GetEntityById(ctx context.Context, request *vectorspb.GetEntityByIdRequest) (resp *vectorspb.GetEntityByIdResponse, err error) {
	ctx = s.c.Scoped(ctx)

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetEntityById(ctx, request)
}
