package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"middleman/assistants/assistantspb"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterGateway(ctx context.Context, mux *chi.Mux, grpcAddr string) error {
	const apiRoot = "/api/assistants"

	// Create gateway with custom header matcher to forward Authorization header
	gateway := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(header string) (string, bool) {
			// Forward Authorization header to gRPC metadata
			if strings.ToLower(header) == "authorization" {
				return header, true
			}
			// Default behavior for other headers
			return runtime.DefaultHeaderMatcher(header)
		}),
	)
	
	err := assistantspb.RegisterAssistantsServiceHandlerFromEndpoint(ctx, gateway, grpcAddr, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		return err
	}

	// mount the GRPC gateway
	mux.Mount(apiRoot, gateway)

	return nil
}
