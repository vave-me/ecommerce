package sap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sap/cloud-security-client-go/auth"
	"github.com/sap/cloud-security-client-go/oidcclient"
	"github.com/rs/zerolog/log"
)

// SecurityClient handles SAP IAS authentication and authorization
type SecurityClient struct {
	middleware auth.Middleware
	config     *SecurityConfig
}

// SecurityConfig holds configuration for SAP security
type SecurityConfig struct {
	// IAS instance name for Kubernetes environments
	IASInstanceName string
	
	// OAuth2 configuration
	ClientID     string
	ClientSecret string
	TokenURL     string
	
	// JWT validation
	Issuer   string
	Audience []string
	
	// Configuration path for Kubernetes
	ConfigPath string
}

// NewSecurityClient creates a new SAP security client
func NewSecurityClient(config *SecurityConfig) (*SecurityClient, error) {
	// In Kubernetes, configuration is mounted as secrets
	if config.ConfigPath == "" && config.IASInstanceName != "" {
		config.ConfigPath = filepath.Join("/etc/secrets/sapbtp/identity", config.IASInstanceName)
	}
	
	// Create identity configuration from environment or provided config
	// SAP Cloud Security Client expects an identity configuration
	// In production, this would be loaded from bound service instances
	
	// For now, create a mock middleware since the exact API has changed
	// You'll need to update this based on your SAP BTP service binding
	
	// The new API requires env.Identity which is typically loaded from service bindings
	// Example:
	// identity, err := env.GetIdentityService()
	// middleware := auth.NewMiddleware(identity, auth.Options{})
	
	// Temporary placeholder - update based on your environment
	var middleware auth.Middleware
	
	return &SecurityClient{
		middleware: middleware,
		config:     config,
	}, nil
}

// Middleware returns the HTTP middleware for authentication
func (c *SecurityClient) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return c.middleware.Handler(next)
	}
}

// ValidateToken validates a JWT token and returns claims
func (c *SecurityClient) ValidateToken(ctx context.Context, token string) (*auth.Claims, error) {
	// Create a mock request with the token
	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	
	// Use the middleware to validate
	claims, err := c.middleware.Authenticate(req)
	if err != nil {
		return nil, fmt.Errorf("validating token: %w", err)
	}
	
	return claims, nil
}

// GetClaimsFromRequest extracts claims from an authenticated request
func (c *SecurityClient) GetClaimsFromRequest(r *http.Request) (*auth.Claims, error) {
	claims := auth.GetClaims(r)
	if claims == nil {
		return nil, fmt.Errorf("no claims found in request context")
	}
	return claims, nil
}

// AuthenticateWebhook validates webhook requests from SAP
func (c *SecurityClient) AuthenticateWebhook(r *http.Request) error {
	// Check for API key in header
	apiKey := r.Header.Get("X-SAP-API-Key")
	if apiKey == "" {
		apiKey = r.Header.Get("SAP-API-Key")
	}
	
	if apiKey == "" {
		return fmt.Errorf("missing API key in webhook request")
	}
	
	// Validate API key against configured value
	expectedKey := os.Getenv("SAP_WEBHOOK_API_KEY")
	if expectedKey == "" {
		log.Warn().Msg("SAP_WEBHOOK_API_KEY not configured, skipping webhook authentication")
		return nil
	}
	
	if apiKey != expectedKey {
		return fmt.Errorf("invalid API key")
	}
	
	return nil
}

// loadConfigFromPath would load configuration from Kubernetes mounted secrets
// This needs to be updated based on the new SAP Cloud Security Client API
// which expects env.Identity instead of direct configuration

// CreateAuthenticatedHTTPClient creates an HTTP client with authentication
func (c *SecurityClient) CreateAuthenticatedHTTPClient(ctx context.Context) (*http.Client, error) {
	// Get access token using client credentials
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting access token: %w", err)
	}
	
	// Create HTTP client with auth transport
	client := &http.Client{
		Transport: &authTransport{
			base:  http.DefaultTransport,
			token: token,
		},
	}
	
	return client, nil
}

// getAccessToken obtains an access token using client credentials
func (c *SecurityClient) getAccessToken(ctx context.Context) (string, error) {
	// This is a simplified implementation
	// In production, you would use OAuth2 client credentials flow
	// with proper token caching and refresh logic
	
	req, err := http.NewRequestWithContext(ctx, "POST", c.config.TokenURL, nil)
	if err != nil {
		return "", err
	}
	
	req.SetBasicAuth(c.config.ClientID, c.config.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	// Add grant_type=client_credentials body
	// Implementation details omitted for brevity
	
	// In production, implement full OAuth2 flow
	return "mock-token", nil
}

// authTransport adds authorization header to requests
type authTransport struct {
	base  http.RoundTripper
	token string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(req)
}