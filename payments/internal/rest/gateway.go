package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"middleman/payments/paymentspb"
)

func RegisterGateway(ctx context.Context, mux *chi.Mux, grpcAddr string) error {
	const apiRoot = "/api/payments"

	// Create a new gRPC-Gateway mux that includes our custom header matcher
	gwMux := runtime.NewServeMux()

	// Register the PaymentsService handler to route incoming HTTP -> gRPC
	err := paymentspb.RegisterPaymentsServiceHandlerFromEndpoint(
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
