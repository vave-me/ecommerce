package domain

import "context"

type InteractionRepository interface {
	Load(ctx context.Context, interactionID string) (*Interaction, error)
	Save(ctx context.Context, interaction *Interaction) error
}
