package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryInterceptor provides panic recovery for gRPC
func RecoveryInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				logger.Error().
					Interface("panic", r).
					Str("method", info.FullMethod).
					Bytes("stack", stack).
					Msg("gRPC handler panic recovered")
				
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// RecoveryMiddleware provides panic recovery for HTTP
func RecoveryMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := debug.Stack()
					logger.Error().
						Interface("panic", rec).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Bytes("stack", stack).
						Msg("HTTP handler panic recovered")

					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// StreamRecoveryInterceptor provides panic recovery for gRPC streams
func StreamRecoveryInterceptor(logger zerolog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				logger.Error().
					Interface("panic", r).
					Str("method", info.FullMethod).
					Bytes("stack", stack).
					Msg("gRPC stream handler panic recovered")
				
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(srv, stream)
	}
}

// LoggingInterceptor logs all gRPC requests
func LoggingInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Debug().
			Str("method", info.FullMethod).
			Interface("request", req).
			Msg("gRPC request")

		resp, err := handler(ctx, req)

		if err != nil {
			logger.Error().
				Str("method", info.FullMethod).
				Err(err).
				Msg("gRPC request failed")
		} else {
			logger.Debug().
				Str("method", info.FullMethod).
				Interface("response", resp).
				Msg("gRPC request succeeded")
		}

		return resp, err
	}
}

// RequestIDMiddleware adds a request ID to the context
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("req-%d", r.Context().Value("request_start").(int64))
		}

		ctx := context.WithValue(r.Context(), "request_id", requestID)
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}