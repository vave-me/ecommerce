package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type RebrandVariant struct {
	ID          string
	Name        string
	Description string
}

type RebrandVariantHandler struct {
	variants  domain.VariantRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRebrandVariantHandler(
	variants domain.VariantRepository,
	publisher ddd.EventPublisher[ddd.Event],
) RebrandVariantHandler {
	return RebrandVariantHandler{
		variants:  variants,
		publisher: publisher,
	}
}

func (h RebrandVariantHandler) RebrandVariant(ctx context.Context, cmd RebrandVariant) error {
	variant, err := h.variants.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := variant.Rebrand(cmd.Name, cmd.Description)
	// domain: variant.Rebrand name
	if err != nil {
		return err
	}

	if err = h.variants.Save(ctx, variant); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
