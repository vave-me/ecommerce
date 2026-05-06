package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"middleman/internal/di"
	"middleman/streams/internal/application"
	"middleman/streams/streamspb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterGateway(ctx context.Context, mux *chi.Mux, grpcAddr string) error {
	const apiRoot = "/api/streams"

	gateway := runtime.NewServeMux()
	err := streamspb.RegisterStreamsServiceHandlerFromEndpoint(ctx, gateway, grpcAddr, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		return err
	}

	// mount the GRPC gateway
	mux.Mount(apiRoot, gateway)

	return nil
}

// RegisterGatewayWithWebhooks registers the gateway with webhook routes
func RegisterGatewayWithWebhooks(ctx context.Context, container di.Container, mux *chi.Mux, grpcAddr string, app application.StreamingApp) error {
	const apiRoot = "/api/streams"

	gateway := runtime.NewServeMux()
	err := streamspb.RegisterStreamsServiceHandlerFromEndpoint(ctx, gateway, grpcAddr, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		return err
	}

	// Create a sub-router for the API
	apiRouter := chi.NewRouter()
	
	// Mount the GRPC gateway
	apiRouter.Mount("/", gateway)
	
	// Register webhook routes
	RegisterWebhookRoutes(container, apiRouter, app)
	
	// Mount the entire API router
	mux.Mount(apiRoot, apiRouter)

	return nil
}
