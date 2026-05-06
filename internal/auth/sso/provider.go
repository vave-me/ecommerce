package sso

import (
	"context"
	"time"
)

// Provider represents an SSO provider interface
type Provider interface {
	// GetProviderName returns the name of the provider (e.g., "google", "azure-ad", "okta")
	GetProviderName() string

	// GetProviderType returns the type of the provider (e.g., "oidc", "saml")
	GetProviderType() string

	// GetAuthorizationURL generates the authorization URL for the OAuth/SAML flow
	GetAuthorizationURL(state string, redirectURI string) string

	// ExchangeCode exchanges an authorization code for tokens (OAuth flow)
	ExchangeCode(ctx context.Context, code string, redirectURI string) (*TokenResponse, error)

	// GetUserInfo retrieves user information from the provider
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)

	// ValidateIDToken validates and parses an ID token (OIDC)
	ValidateIDToken(ctx context.Context, idToken string) (*Claims, error)

	// RefreshToken refreshes an access token using a refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)

	// RevokeToken revokes an access or refresh token
	RevokeToken(ctx context.Context, token string) error
}

// TokenResponse represents the response from a token exchange
type TokenResponse struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
	TokenType    string
	ExpiresIn    int
	Scope        string
}

// UserInfo represents user information from an SSO provider
type UserInfo struct {
	ID            string
	Email         string
	EmailVerified bool
	Name          string
	FirstName     string
	LastName      string
	Picture       string
	Locale        string
	Provider      string
	Raw           map[string]interface{} // Raw claims for provider-specific data
}

// Claims represents the claims in an ID token
type Claims struct {
	Issuer            string
	Subject           string
	Audience          []string
	ExpirationTime    time.Time
	NotBefore         time.Time
	IssuedAt          time.Time
	Email             string
	EmailVerified     bool
	Name              string
	Picture           string
	Nonce             string
	AuthorizationTime time.Time
	Raw               map[string]interface{}
}

// Config represents the configuration for an SSO provider
type Config struct {
	ClientID     string
	ClientSecret string
	Issuer       string
	Scopes       []string
	// OIDC specific
	DiscoveryURL string
	// SAML specific
	MetadataURL     string
	EntityID        string
	AssertionConsumerServiceURL string
	// Common
	CallbackURL string
}