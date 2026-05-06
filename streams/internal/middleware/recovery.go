package middleware

import (
	"net/http"
	"runtime/debug"
	
	"go.uber.org/zap"
)

// Recovery middleware recovers from panics and logs them
func Recovery(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic with stack trace
					logger.Error("Panic recovered",
						zap.Any("error", err),
						zap.String("path", r.URL.Path),
						zap.String("method", r.Method),
						zap.String("remote_addr", r.RemoteAddr),
						zap.String("user_agent", r.UserAgent()),
						zap.String("stack", string(debug.Stack())))
					
					// Send error response
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"Internal Server Error"}`))
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryWithMetrics adds metrics to panic recovery
func RecoveryWithMetrics(logger *zap.Logger, panicCounter func()) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Increment panic counter
					if panicCounter != nil {
						panicCounter()
					}
					
					// Log the panic with stack trace
					logger.Error("Panic recovered",
						zap.Any("error", err),
						zap.String("path", r.URL.Path),
						zap.String("method", r.Method),
						zap.String("remote_addr", r.RemoteAddr),
						zap.String("user_agent", r.UserAgent()),
						zap.String("stack", string(debug.Stack())))
					
					// Send error response
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"Internal Server Error"}`))
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}