package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type IncreaseVariantPrice struct {
	ID    string
	Price int64
}

type IncreaseVariantPriceHandler struct {
	variants  domain.VariantRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewIncreaseVariantPriceHandler(variants domain.VariantRepository, publisher ddd.EventPublisher[ddd.Event]) IncreaseVariantPriceHandler {
	return IncreaseVariantPriceHandler{
		variants:  variants,
		publisher: publisher,
	}
}

func (h IncreaseVariantPriceHandler) IncreaseVariantPrice(ctx context.Context, cmd IncreaseVariantPrice) error {
	variant, err := h.variants.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := variant.IncreasePrice(cmd.ID, cmd.Price)
	if err != nil {
		return err
	}

	err = h.variants.Save(ctx, variant)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
