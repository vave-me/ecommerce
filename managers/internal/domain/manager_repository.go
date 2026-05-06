package domain

import (
	"context"
)

// ManagerRepository defines the interface for manager persistence (Event Sourcing)
type ManagerRepository interface {
	Load(ctx context.Context, id string) (*Manager, error)
	Save(ctx context.Context, manager *Manager) error
}
