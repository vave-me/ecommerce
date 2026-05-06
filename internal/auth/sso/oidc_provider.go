package sso

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/stackus/errors"
	"golang.org/x/oauth2"
)

// OIDCProvider implements the Provider interface for OpenID Connect providers
type OIDCProvider struct {
	name         string
	config       *Config
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	provider     *oidc.Provider
}

// NewOIDCProvider creates a new OIDC provider
func NewOIDCProvider(ctx context.Context, name string, config *Config) (*OIDCProvider, error) {
	if config.DiscoveryURL == "" && config.Issuer == "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "either discovery URL or issuer must be provided")
	}

	issuer := config.Issuer
	if config.DiscoveryURL != "" {
		issuer = config.DiscoveryURL
	}

	// Initialize OIDC provider
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create OIDC provider")
	}

	// Configure OAuth2
	scopes := config.Scopes
	if len(scopes) == 0 {
		scopes = []string{oidc.ScopeOpenID, "profile", "email"}
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  config.CallbackURL,
		Scopes:       scopes,
	}

	// Create ID token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

	return &OIDCProvider{
		name:         name,
		config:       config,
		oauth2Config: oauth2Config,
		verifier:     verifier,
		provider:     provider,
	}, nil
}

// GetProviderName returns the name of the provider
func (p *OIDCProvider) GetProviderName() string {
	return p.name
}

// GetProviderType returns "oidc"
func (p *OIDCProvider) GetProviderType() string {
	return "oidc"
}

// GetAuthorizationURL generates the authorization URL
func (p *OIDCProvider) GetAuthorizationURL(state string, redirectURI string) string {
	// Override redirect URI if provided
	config := *p.oauth2Config
	if redirectURI != "" {
		config.RedirectURL = redirectURI
	}

	return config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode exchanges an authorization code for tokens
func (p *OIDCProvider) ExchangeCode(ctx context.Context, code string, redirectURI string) (*TokenResponse, error) {
	// Override redirect URI if provided
	config := *p.oauth2Config
	if redirectURI != "" {
		config.RedirectURL = redirectURI
	}

	oauth2Token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange authorization code")
	}

	// Extract raw ID token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, errors.Wrap(errors.ErrInternal, "no id_token in token response")
	}

	return &TokenResponse{
		AccessToken:  oauth2Token.AccessToken,
		RefreshToken: oauth2Token.RefreshToken,
		IDToken:      rawIDToken,
		TokenType:    oauth2Token.TokenType,
		ExpiresIn:    int(oauth2Token.Expiry.Unix() - oauth2Token.Expiry.Unix()),
	}, nil
}

// GetUserInfo retrieves user information using the access token
func (p *OIDCProvider) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	userInfo, err := p.provider.UserInfo(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user info")
	}

	var claims map[string]interface{}
	if err := userInfo.Claims(&claims); err != nil {
		return nil, errors.Wrap(err, "failed to parse user info claims")
	}

	// Extract standard claims
	info := &UserInfo{
		Provider: p.name,
		Raw:      claims,
	}

	// Map standard OIDC claims
	if sub, ok := claims["sub"].(string); ok {
		info.ID = sub
	}
	if email, ok := claims["email"].(string); ok {
		info.Email = email
	}
	if emailVerified, ok := claims["email_verified"].(bool); ok {
		info.EmailVerified = emailVerified
	}
	if name, ok := claims["name"].(string); ok {
		info.Name = name
	}
	if givenName, ok := claims["given_name"].(string); ok {
		info.FirstName = givenName
	}
	if familyName, ok := claims["family_name"].(string); ok {
		info.LastName = familyName
	}
	if picture, ok := claims["picture"].(string); ok {
		info.Picture = picture
	}
	if locale, ok := claims["locale"].(string); ok {
		info.Locale = locale
	}

	return info, nil
}

// ValidateIDToken validates and parses an ID token
func (p *OIDCProvider) ValidateIDToken(ctx context.Context, idToken string) (*Claims, error) {
	token, err := p.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify ID token")
	}

	var claims map[string]interface{}
	if err := token.Claims(&claims); err != nil {
		return nil, errors.Wrap(err, "failed to parse ID token claims")
	}

	// Extract standard claims
	result := &Claims{
		Issuer:         token.Issuer,
		Subject:        token.Subject,
		Audience:       token.Audience,
		ExpirationTime: token.Expiry,
		IssuedAt:       token.IssuedAt,
		Raw:            claims,
	}

	// Map additional claims
	if email, ok := claims["email"].(string); ok {
		result.Email = email
	}
	if emailVerified, ok := claims["email_verified"].(bool); ok {
		result.EmailVerified = emailVerified
	}
	if name, ok := claims["name"].(string); ok {
		result.Name = name
	}
	if picture, ok := claims["picture"].(string); ok {
		result.Picture = picture
	}
	if nonce, ok := claims["nonce"].(string); ok {
		result.Nonce = nonce
	}

	return result, nil
}

// RefreshToken refreshes an access token using a refresh token
func (p *OIDCProvider) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	tokenSource := p.oauth2Config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	})

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, errors.Wrap(err, "failed to refresh token")
	}

	// Extract raw ID token if present
	rawIDToken, _ := newToken.Extra("id_token").(string)

	return &TokenResponse{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		IDToken:      rawIDToken,
		TokenType:    newToken.TokenType,
		ExpiresIn:    int(newToken.Expiry.Unix() - newToken.Expiry.Unix()),
	}, nil
}

// RevokeToken revokes a token (if supported by the provider)
func (p *OIDCProvider) RevokeToken(ctx context.Context, token string) error {
	// Check if provider supports revocation endpoint
	var revocationEndpoint string
	if p.provider != nil {
		// Try to get revocation endpoint from provider metadata
		// This is provider-specific and may not be standardized
		endpoint := strings.TrimSuffix(p.config.Issuer, "/") + "/oauth2/revoke"
		revocationEndpoint = endpoint
	}

	if revocationEndpoint == "" {
		// Revocation not supported
		return nil
	}

	// Implement token revocation
	// This would typically involve making an HTTP POST request to the revocation endpoint
	// For now, we'll return nil as revocation is optional
	return nil
}

// ProviderConfig contains pre-configured settings for common OIDC providers
var ProviderConfigs = map[string]func(clientID, clientSecret, callbackURL string) *Config{
	"google": func(clientID, clientSecret, callbackURL string) *Config {
		return &Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Issuer:       "https://accounts.google.com",
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
			CallbackURL:  callbackURL,
		}
	},
	"azure-ad": func(clientID, clientSecret, callbackURL string) *Config {
		return &Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			// Note: For Azure AD, the issuer includes the tenant ID
			// This should be configured per deployment
			DiscoveryURL: fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", "common"),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
			CallbackURL:  callbackURL,
		}
	},
	"okta": func(clientID, clientSecret, callbackURL string) *Config {
		return &Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			// Note: For Okta, the issuer includes the Okta domain
			// This should be configured per deployment
			Scopes:      []string{oidc.ScopeOpenID, "profile", "email"},
			CallbackURL: callbackURL,
		}
	},
}