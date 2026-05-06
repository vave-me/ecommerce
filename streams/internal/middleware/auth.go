package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"
	
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"middleman/streams/internal/domain"
)

// Claims represents JWT claims
type Claims struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	Role     string   `json:"role"`
	Scopes   []string `json:"scopes"`
	TenantID string   `json:"tenant_id"`
	jwt.RegisteredClaims
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	JWTSecret    string
	APIKeyHeader string
	Logger       *zap.Logger
}

// JWTAuth creates JWT authentication middleware
func JWTAuth(config *AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondUnauthorized(w, "missing authorization header")
				return
			}
			
			// Check Bearer token format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondUnauthorized(w, "invalid authorization format")
				return
			}
			
			tokenString := parts[1]
			
			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(config.JWTSecret), nil
			})
			
			if err != nil {
				config.Logger.Debug("JWT validation failed", zap.Error(err))
				respondUnauthorized(w, "invalid token")
				return
			}
			
			claims, ok := token.Claims.(*Claims)
			if !ok || !token.Valid {
				respondUnauthorized(w, "invalid token claims")
				return
			}
			
			// Check token expiration
			if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
				respondUnauthorized(w, "token expired")
				return
			}
			
			// Add claims to context
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "user_email", claims.Email)
			ctx = context.WithValue(ctx, "user_role", claims.Role)
			ctx = context.WithValue(ctx, "user_scopes", claims.Scopes)
			ctx = context.WithValue(ctx, "tenant_id", claims.TenantID)
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// APIKeyAuth creates API key authentication middleware
func APIKeyAuth(config *AuthConfig, validateKey func(string) (*Claims, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract API key from header
			apiKey := r.Header.Get(config.APIKeyHeader)
			if apiKey == "" {
				respondUnauthorized(w, "missing API key")
				return
			}
			
			// Validate API key
			claims, err := validateKey(apiKey)
			if err != nil {
				config.Logger.Debug("API key validation failed", 
					zap.Error(err),
					zap.String("key_prefix", apiKey[:min(8, len(apiKey))]+"..."))
				respondUnauthorized(w, "invalid API key")
				return
			}
			
			// Add claims to context
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "user_email", claims.Email)
			ctx = context.WithValue(ctx, "user_role", claims.Role)
			ctx = context.WithValue(ctx, "user_scopes", claims.Scopes)
			ctx = context.WithValue(ctx, "tenant_id", claims.TenantID)
			ctx = context.WithValue(ctx, "api_key", true)
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireScopes ensures the user has required scopes
func RequireScopes(scopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userScopes, ok := r.Context().Value("user_scopes").([]string)
			if !ok {
				respondForbidden(w, "no scopes found")
				return
			}
			
			// Check if user has all required scopes
			for _, required := range scopes {
				found := false
				for _, userScope := range userScopes {
					if userScope == required {
						found = true
						break
					}
				}
				if !found {
					respondForbidden(w, "insufficient scopes")
					return
				}
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole ensures the user has the required role
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("user_role").(string)
			if !ok {
				respondForbidden(w, "no role found")
				return
			}
			
			// Check if user has required role
			found := false
			for _, role := range roles {
				if userRole == role {
					found = true
					break
				}
			}
			
			if !found {
				respondForbidden(w, "insufficient role")
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// OptionalAuth attempts authentication but doesn't require it
func OptionalAuth(config *AuthConfig) func(http.Handler) http.Handler {
	jwtAuth := JWTAuth(config)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If Authorization header exists, validate it
			if r.Header.Get("Authorization") != "" {
				// Create a custom response writer to capture auth errors
				captureWriter := &responseCapture{ResponseWriter: w}
				jwtAuth(next).ServeHTTP(captureWriter, r)
				
				// If auth failed, continue anyway
				if captureWriter.statusCode == http.StatusUnauthorized {
					next.ServeHTTP(w, r)
				}
			} else {
				// No auth header, continue
				next.ServeHTTP(w, r)
			}
		})
	}
}

// Helper functions
func respondUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"` + message + `"}`))
}

func respondForbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"error":"` + message + `"}`))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// responseCapture captures response status for optional auth
type responseCapture struct {
	http.ResponseWriter
	statusCode int
}

func (rc *responseCapture) WriteHeader(code int) {
	rc.statusCode = code
	rc.ResponseWriter.WriteHeader(code)
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", domain.ErrUnauthorized
	}
	return userID, nil
}

// GetUserRole extracts user role from context
func GetUserRole(ctx context.Context) (string, error) {
	role, ok := ctx.Value("user_role").(string)
	if !ok {
		return "", domain.ErrUnauthorized
	}
	return role, nil
}

// IsAdmin checks if user has admin role
func IsAdmin(ctx context.Context) bool {
	role, _ := GetUserRole(ctx)
	return role == "admin" || role == "super_admin"
}