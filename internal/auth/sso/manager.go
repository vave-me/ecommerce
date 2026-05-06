package sso

import (
	"context"
	"sync"

	"github.com/stackus/errors"
)

// Manager manages multiple SSO providers
type Manager struct {
	providers map[string]Provider
	mu        sync.RWMutex
}

// NewManager creates a new SSO manager
func NewManager() *Manager {
	return &Manager{
		providers: make(map[string]Provider),
	}
}

// RegisterProvider registers a new SSO provider
func (m *Manager) RegisterProvider(name string, provider Provider) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.providers[name]; exists {
		return errors.Wrap(errors.ErrAlreadyExists, "provider already registered")
	}

	m.providers[name] = provider
	return nil
}

// GetProvider retrieves a provider by name
func (m *Manager) GetProvider(name string) (Provider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, exists := m.providers[name]
	if !exists {
		return nil, errors.Wrap(errors.ErrNotFound, "provider not found")
	}

	return provider, nil
}

// ListProviders returns a list of all registered provider names
func (m *Manager) ListProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.providers))
	for name := range m.providers {
		names = append(names, name)
	}

	return names
}

// RemoveProvider removes a provider from the manager
func (m *Manager) RemoveProvider(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.providers[name]; !exists {
		return errors.Wrap(errors.ErrNotFound, "provider not found")
	}

	delete(m.providers, name)
	return nil
}

// InitializeFromConfig initializes providers from configuration
func (m *Manager) InitializeFromConfig(ctx context.Context, configs map[string]*Config) error {
	for name, config := range configs {
		var provider Provider
		var err error

		// Determine provider type based on configuration
		if config.DiscoveryURL != "" || config.Issuer != "" {
			// OIDC provider
			provider, err = NewOIDCProvider(ctx, name, config)
		} else if config.MetadataURL != "" {
			// SAML provider (to be implemented)
			return errors.Wrap(errors.ErrUnimplemented, "SAML providers not yet implemented")
		} else {
			return errors.Wrap(errors.ErrBadRequest, "unable to determine provider type from configuration")
		}

		if err != nil {
			return errors.Wrapf(err, "failed to initialize provider %s", name)
		}

		if err := m.RegisterProvider(name, provider); err != nil {
			return errors.Wrapf(err, "failed to register provider %s", name)
		}
	}

	return nil
}