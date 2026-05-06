package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type RemoveVariant struct {
	ID     string
	UserID string // optional if needed
}

type RemoveVariantHandler struct {
	variants  domain.VariantRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRemoveVariantHandler(
	variants domain.VariantRepository,
	publisher ddd.EventPublisher[ddd.Event],
) RemoveVariantHandler {
	return RemoveVariantHandler{
		variants:  variants,
		publisher: publisher,
	}
}

func (h RemoveVariantHandler) RemoveVariant(ctx context.Context, cmd RemoveVariant) error {
	variant, err := h.variants.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := variant.Remove(cmd.ID) // domain: variant.Remove
	if err != nil {
		return err
	}

	if err = h.variants.Save(ctx, variant); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
