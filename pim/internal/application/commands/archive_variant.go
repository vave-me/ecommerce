package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type ArchiveVariant struct {
	ID string // variant ID
}

type ArchiveVariantHandler struct {
	variants  domain.VariantRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewArchiveVariantHandler(
	variants domain.VariantRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ArchiveVariantHandler {
	return ArchiveVariantHandler{
		variants:  variants,
		publisher: publisher,
	}
}

func (h ArchiveVariantHandler) ArchiveVariant(ctx context.Context, cmd ArchiveVariant) error {
	variant, err := h.variants.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := variant.Archive() // domain: variant.Archive
	if err != nil {
		return err
	}

	if err = h.variants.Save(ctx, variant); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
