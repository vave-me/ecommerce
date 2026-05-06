package domain

import (
	"context"
)

// AssistantRepository defines the interface for assistant persistence (Event Sourcing)
type AssistantRepository interface {
	Load(ctx context.Context, id string) (*Assistant, error)
	Save(ctx context.Context, assistant *Assistant) error
}
