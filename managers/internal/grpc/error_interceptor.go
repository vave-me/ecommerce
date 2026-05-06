package grpc

import (
	"context"

	models "middleman/managers/internal/models"

	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorInterceptor converts returned *models.ServiceError into gRPC status errors
// so clients (including LLM agents) receive structured, predictable feedback.
func ErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		// Unwrap ServiceError if present
		if se, ok := err.(*models.ServiceError); ok {
			// Map validation-error to InvalidArgument, others to Internal by default
			code := grpc_code.Internal
			if se.Code == "validation-error" {
				code = grpc_code.InvalidArgument
			}
			return nil, status.Errorf(code, "%s", se.Message)
		}

		return nil, err // passthrough raw error
	}
}
