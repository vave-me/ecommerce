package domain

import "context"

type ActionRepository interface {
	Load(ctx context.Context, interactionID string) (*Action, error)
	Save(ctx context.Context, interaction *Action) error
}
