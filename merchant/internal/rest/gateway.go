package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"middleman/merchant/merchantpb"
)

func RegisterGateway(ctx context.Context, mux *chi.Mux, grpcAddr string) error {
	const apiRoot = "/api/merchant"

	// Create a new gRPC-Gateway mux
	gwMux := runtime.NewServeMux()

	// Register the MerchantService handler to route incoming HTTP -> gRPC
	err := merchantpb.RegisterMerchantServiceHandlerFromEndpoint(
		ctx,
		gwMux,
		grpcAddr,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	)
	if err != nil {
		return err
	}

	mux.Mount(apiRoot, gwMux)

	return nil
}