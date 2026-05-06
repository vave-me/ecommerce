package rpc

import (
	"context"
	"middleman/internal/auth"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	promgrpc "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/rs/zerolog/log"
	"github.com/stackus/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// ClientErrorUnaryInterceptor handles GRPC client-side errors
func clientErrorUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return errors.ReceiveGRPCError(invoker(ctx, method, req, reply, cc, opts...))
	}
}

func logStreamMetadataInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if ok {
			for key, values := range md {
				log.Printf("[MetadataInterceptor:STREAM] %s => %v", key, values)
			}
		} else {
			log.Printf("[MetadataInterceptor:STREAM] No metadata for method: %s", info.FullMethod)
		}
		return handler(srv, ss)
	}
}

// Dial establishes a gRPC connection using grpc.NewClient() with metrics and OpenTelemetry tracing
func Dial(ctx context.Context, endpoint string) (conn *grpc.ClientConn, err error) {
	clientMetrics := promgrpc.NewClientMetrics()

	opts := []grpc.DialOption{

		grpc.WithTransportCredentials(insecure.NewCredentials()),            // Use insecure credentials for local testing
		grpc.WithUnaryInterceptor(clientMetrics.UnaryClientInterceptor()),   // Prometheus unary interceptor
		grpc.WithStreamInterceptor(clientMetrics.StreamClientInterceptor()), // Prometheus stream interceptor
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),                  // OpenTelemetry client handler
	}

	return grpc.DialContext(ctx, endpoint, opts...)
}

// DialWithAuth establishes a gRPC connection with JWT authentication for service-to-service calls
// This is specifically designed for the assistants service to authenticate outgoing calls
func DialWithAuth(ctx context.Context, endpoint string, authInstance *auth.Auth) (conn *grpc.ClientConn, err error) {
	clientMetrics := promgrpc.NewClientMetrics()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			authInstance.UnaryClientInterceptor(),  // JWT Authentication for outgoing calls
			clientMetrics.UnaryClientInterceptor(), // Prometheus unary interceptor
			clientErrorUnaryInterceptor(),          // Error handling
		)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			authInstance.StreamClientInterceptor(),  // JWT Authentication for outgoing calls
			clientMetrics.StreamClientInterceptor(), // Prometheus stream interceptor
		)),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // OpenTelemetry client handler
	}

	return grpc.DialContext(ctx, endpoint, opts...)
}
