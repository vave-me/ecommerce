package oidcclient

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Config holds the options required to create a GoogleOIDCClient.
type Config struct {
	ClientID      string
	Issuer        string
	CacheDuration time.Duration
}
type GoogleUserClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Picture       string `json:"picture"` // Thumbnail/profile picture URL
	Subject       string `json:"sub"`     // Unique user ID from Google
}

// GoogleOIDCClient encapsulates OIDC Provider + IDTokenVerifier.
// It is concurrency-safe. Provider metadata (incl. Google JWK keys) is
// periodically refreshed, so long-running processes stay up-to-date.
type GoogleOIDCClient struct {
	cfg       Config
	mu        sync.RWMutex
	verifier  *oidc.IDTokenVerifier
	provider  *oidc.Provider
	expiresAt time.Time
}

// ---------------------------
// Constructor
// ---------------------------

// NewGoogleOIDCClient builds the client and fetches Google’s JWK keys once.
func NewGoogleOIDCClient(ctx context.Context, clientID, issuer string) *GoogleOIDCClient {
	// Basic validation
	if clientID == "" {
		return nil
	}
	if issuer == "" {
		issuer = "https://accounts.google.com"
	}

	// Compose a minimal Config with sane defaults
	cfg := Config{
		ClientID:      clientID,
		Issuer:        issuer,
		CacheDuration: 6 * time.Hour,
	}

	c := &GoogleOIDCClient{cfg: cfg}
	if err := c.reload(ctx); err != nil {
		return nil
	}
	return c
}

// ---------------------------
// Public Methods
// ---------------------------

// VerifyIDToken validates signature, issuer, audience & expiry.
// It transparently refreshes JWKs / metadata when cache is stale.
func (c *GoogleOIDCClient) VerifyIDToken(ctx context.Context, raw string) (*oidc.IDToken, error) {
	if raw == "" {
		return nil, errors.New("oidc client: idToken is empty")
	}

	// Refresh provider if needed
	if err := c.refreshIfNeeded(ctx); err != nil {
		return nil, err
	}

	c.mu.RLock()
	ver := c.verifier
	c.mu.RUnlock()

	idTok, err := ver.Verify(ctx, raw)
	if err != nil {
		return nil, fmt.Errorf("oidc client: verify failed: %w", err)
	}
	return idTok, nil
}

// ParseClaims unmarshals the JWT payload into dst (struct/map).
// Example:
//
//	var c struct{ Email string `json:"email"` }
//	err := client.ParseClaims(idTok, &c)
func (c *GoogleOIDCClient) ParseClaims(idTok *oidc.IDToken, dst interface{}) error {
	if idTok == nil {
		return errors.New("oidc client: idToken is nil")
	}
	if err := idTok.Claims(dst); err != nil {
		return fmt.Errorf("oidc client: parsing claims: %w", err)
	}
	return nil
}

// Provider returns a *oidc.Provider (read-only) in case callers need
// advanced OIDC flows (e.g., code-exchange).
func (c *GoogleOIDCClient) Provider() *oidc.Provider {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.provider
}

// ---------------------------
// Internal helpers
// ---------------------------

func (c *GoogleOIDCClient) refreshIfNeeded(ctx context.Context) error {
	c.mu.RLock()
	need := time.Now().After(c.expiresAt)
	c.mu.RUnlock()
	if !need {
		return nil
	}
	// Only one goroutine should reload
	return c.reload(ctx)
}

func (c *GoogleOIDCClient) reload(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If another goroutine refreshed already, no need to redo.
	if time.Now().Before(c.expiresAt) && c.verifier != nil {
		return nil
	}

	provider, err := oidc.NewProvider(ctx, c.cfg.Issuer)
	if err != nil {
		return fmt.Errorf("oidc client: fetching provider metadata: %w", err)
	}

	// Build verifier
	ver := provider.Verifier(&oidc.Config{ClientID: c.cfg.ClientID})

	// Store
	c.provider = provider
	c.verifier = ver
	c.expiresAt = time.Now().Add(c.cfg.CacheDuration)
	return nil
}

// ---------------------------
// (Optional) OAuth2 Config Forwarder
// ---------------------------

// OAuth2Config returns a Google OAuth2 config pre-filled with clientID and scopes.
// Handy if later you need full OAuth (auth-code flow) instead of only ID-tokens.
func (c *GoogleOIDCClient) OAuth2Config(redirectURL string, scopes ...string) oauth2.Config {
	if len(scopes) == 0 {
		scopes = []string{oidc.ScopeOpenID, "profile", "email"}
	}
	return oauth2.Config{
		ClientID:     c.cfg.ClientID,
		Endpoint:     c.provider.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		ClientSecret: "", // public client, no secret needed in browser/mobile
	}
}
