package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// contextKey is a custom type to prevent context key collisions
type contextKey string

const claimsContextKey = contextKey("claims")
const userRoleContextKey = contextKey("userRole")

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type JwtUser struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Username  string  `json:"username"`
	Lat       float64 `json:"lat"`
	Long      float64 `json:"long"`
	Role      string  `json:"role"`
}

type TokenPairs struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Authenticator interface {
	GenerateTokenPair(user *JwtUser) (TokenPairs, error)
	ValidateRefreshToken(refreshToken string) (*JwtUser, error)
	// Add other necessary methods if needed
}

// Claims type for JWT
type Claims struct {
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Ensure Auth implements Authenticator
var _ Authenticator = (*Auth)(nil)

func (auth *Auth) GenerateTokenPair(user *JwtUser) (TokenPairs, error) {
	// Create access token
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	accessClaims["sub"] = user.ID
	accessClaims["aud"] = auth.Audience
	accessClaims["iss"] = auth.Issuer
	accessClaims["iat"] = time.Now().UTC().Unix()
	accessClaims["exp"] = time.Now().UTC().Add(auth.TokenExpiry).Unix()
	accessClaims["username"] = user.Username
	accessClaims["role"] = user.Role
	accessClaims["lat"] = user.Lat
	accessClaims["long"] = user.Long
	signedAccessToken, err := accessToken.SignedString([]byte(auth.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["sub"] = user.ID
	refreshClaims["iat"] = time.Now().UTC().Unix()
	refreshClaims["exp"] = time.Now().UTC().Add(auth.RefreshExpiry).Unix()
	refreshClaims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	refreshClaims["username"] = user.Username
	refreshClaims["role"] = user.Role
	refreshClaims["lat"] = user.Lat
	refreshClaims["long"] = user.Long
	signedRefreshToken, err := refreshToken.SignedString([]byte(auth.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	return TokenPairs{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (auth *Auth) ValidateRefreshToken(refreshToken string) (*JwtUser, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(auth.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &JwtUser{
			ID: claims["sub"].(string),
		}

		// Extract all other claims if they exist
		if username, ok := claims["username"].(string); ok {
			user.Username = username
		}
		if role, ok := claims["role"].(string); ok {
			user.Role = role
		}
		if lat, ok := claims["lat"].(float64); ok {
			user.Lat = lat
		}
		if long, ok := claims["long"].(float64); ok {
			user.Long = long
		}

		return user, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

// GetRefreshCookie returns a secure HTTP-only cookie for the refresh token
func (auth *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     auth.CookieName,
		Path:     auth.CookiePath,
		Value:    refreshToken,
		Expires:  time.Now().Add(auth.RefreshExpiry),
		MaxAge:   int(auth.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   auth.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

// GetExpiredRefreshCookie returns a cookie that expires immediately, used for clearing the refresh token
func (auth *Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     auth.CookieName,
		Path:     auth.CookiePath,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Domain:   auth.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

// GetTokenFromHeaderAndVerify extracts JWT from Authorization header and verifies it
func (auth *Auth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	w.Header().Add("Vary", "Authorization")

	// get auth header
	authHeader := r.Header.Get("Authorization")

	// sanity check
	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}

	// split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("invalid auth header")
	}

	// check to see if we have the word Bearer
	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("invalid auth header")
	}

	token := headerParts[1]

	// declare an empty claims
	claims := &Claims{}

	// parse the token
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(auth.Secret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}

	//if claims.Issuer != auth.Issuer {
	//	return "", nil, errors.New("invalid issuer")
	//}

	return token, claims, nil
}

func (auth *Auth) JwtInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing context metadata")
		}

		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(auth.Secret), nil
		})

		if err != nil || !token.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		//if claims.Issuer != auth.Issuer {
		//	return nil, status.Errorf(codes.Unauthenticated, "invalid issuer")
		//}

		// Add user info to the context
		ctx = context.WithValue(ctx, "userID", claims.Subject)
		ctx = context.WithValue(ctx, claimsContextKey, claims)
		if claims.Role != "" {
			ctx = context.WithValue(ctx, userRoleContextKey, claims.Role)
		}
		return handler(ctx, req)
	}
}

// UnaryServerInterceptor returns a new unary server interceptor that authenticates via JWT
func (auth *Auth) UnaryServerInterceptor(publicMethods map[string]bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Check if the method is public
		if strings.HasPrefix(info.FullMethod, "/grpc.reflection.v1alpha.ServerReflection") {
			return handler(ctx, req)
		}
		if isPublicMethod(info.FullMethod, publicMethods) {
			return handler(ctx, req)
		}
		var err error
		ctx, err = auth.authenticate(ctx)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptor that authenticates via JWT
func (auth *Auth) StreamServerInterceptor(publicMethods map[string]bool) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Check if the method is public
		if isPublicMethod(info.FullMethod, publicMethods) {
			return handler(srv, ss)
		}
		ctx, err := auth.authenticate(ss.Context())
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = ctx
		return handler(srv, wrapped)
	}
}

// authenticate authenticates a unary RPC using the Authorization header
func (auth *Auth) authenticate(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return ctx, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	tokenString := strings.TrimPrefix(authHeaders[0], "Bearer ")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}
		return []byte(auth.Secret), nil
	})
	if err != nil || !token.Valid {
		return ctx, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	// Additional claims validation
	err = claims.Valid()
	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	// Add claims to context
	ctx = context.WithValue(ctx, claimsContextKey, claims)
	
	// Log the parsed claims for debugging
	log.Printf("[AUTH] Parsed JWT claims - Subject: %s, Role: %s", claims.Subject, claims.Role)

	// Add role to context from claims
	if claims.Role != "" {
		ctx = context.WithValue(ctx, userRoleContextKey, claims.Role)
	}

	return ctx, nil
}

// ClaimsFromContext retrieves the JWT claims from the context
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	return claims, ok
}

// GetRoleFromContext retrieves the user role from the context
func GetRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(userRoleContextKey).(string)
	return role, ok
}

// GetRoleFromToken extracts the role from JWT MapClaims
func GetRoleFromToken(token *jwt.Token) (string, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if role, exists := claims["role"].(string); exists {
			return role, nil
		}
		return "", errors.New("role not found in token")
	}
	return "", errors.New("invalid token claims")
}

// Helper function to check if a method is public
func isPublicMethod(fullMethodName string, publicMethods map[string]bool) bool {
	_, ok := publicMethods[fullMethodName]
	return ok
}
func (auth *Auth) GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
func (auth *Auth) CookieAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetTokenFromCookie(r)
		if err != nil {
			// Handle missing or invalid token
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(auth.Secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), claimsContextKey, claims)
		if claims.Role != "" {
			ctx = context.WithValue(ctx, userRoleContextKey, claims.Role)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GeneratePasswordResetToken generates a JWT token specifically for password reset
func (auth *Auth) GeneratePasswordResetToken(userId, email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userId
	claims["iat"] = time.Now().UTC().Unix()
	claims["exp"] = time.Now().UTC().Add(auth.TokenExpiry).Unix()
	claims["purpose"] = "password_reset"
	claims["email"] = email

	return token.SignedString([]byte(auth.Secret))
}

// ValidatePasswordResetToken validates a password reset token
func (auth *Auth) ValidatePasswordResetToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(auth.Secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verify it's a password reset token
		if purpose, ok := claims["purpose"].(string); !ok || purpose != "password_reset" {
			return "", errors.New("invalid token purpose")
		}

		return claims["sub"].(string), nil
	}

	return "", errors.New("invalid token")
}

// UnaryClientInterceptor returns a new unary client interceptor that adds JWT authentication to outgoing calls
func (auth *Auth) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Get authentication token for outgoing call
		token, err := auth.getServiceToken(ctx)
		if err != nil {
			// If we can't get a token, proceed without authentication (for backwards compatibility)
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		// Add Authorization header to outgoing metadata
		md := metadata.Pairs("authorization", "Bearer "+token)
		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamClientInterceptor returns a new stream client interceptor that adds JWT authentication to outgoing calls
func (auth *Auth) StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// Get authentication token for outgoing call
		token, err := auth.getServiceToken(ctx)
		if err != nil {
			// If we can't get a token, proceed without authentication (for backwards compatibility)
			return streamer(ctx, desc, cc, method, opts...)
		}

		// Add Authorization header to outgoing metadata
		md := metadata.Pairs("authorization", "Bearer "+token)
		ctx = metadata.NewOutgoingContext(ctx, md)

		return streamer(ctx, desc, cc, method, opts...)
	}
}

// getServiceToken gets an appropriate JWT token for service-to-service communication
func (auth *Auth) getServiceToken(ctx context.Context) (string, error) {
	// Option 1: Try to propagate user token from incoming request
	if claims, ok := ClaimsFromContext(ctx); ok {
		return auth.generateTokenFromClaims(claims)
	}

	// Option 2: Generate service-to-service token
	return auth.generateServiceToken()
}

// generateTokenFromClaims creates a new token based on existing claims (token propagation)
func (auth *Auth) generateTokenFromClaims(claims *Claims) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)

	// Propagate user information
	tokenClaims["sub"] = claims.Subject
	tokenClaims["role"] = claims.Role
	tokenClaims["iss"] = auth.Issuer
	tokenClaims["aud"] = auth.Audience
	tokenClaims["iat"] = time.Now().UTC().Unix()
	tokenClaims["exp"] = time.Now().UTC().Add(time.Hour).Unix() // 1 hour expiry
	tokenClaims["propagated"] = true                            // Mark as propagated token

	return token.SignedString([]byte(auth.Secret))
}

// generateServiceToken creates a JWT token for service-to-service communication
func (auth *Auth) generateServiceToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// Service identity claims
	claims["sub"] = "assistants-service"
	claims["iss"] = auth.Issuer
	claims["aud"] = auth.Audience
	claims["iat"] = time.Now().UTC().Unix()
	claims["exp"] = time.Now().UTC().Add(time.Hour).Unix() // 1 hour expiry
	claims["service"] = true                               // Mark as service token
	claims["service_name"] = "assistants"

	return token.SignedString([]byte(auth.Secret))
}
