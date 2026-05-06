package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type DecreaseVariantPrice struct {
	ID    string
	Price int64
}

type DecreaseVariantPriceHandler struct {
	variants  domain.VariantRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDecreaseVariantPriceHandler(variants domain.VariantRepository, publisher ddd.EventPublisher[ddd.Event]) DecreaseVariantPriceHandler {
	return DecreaseVariantPriceHandler{
		variants:  variants,
		publisher: publisher,
	}
}

func (h DecreaseVariantPriceHandler) DecreaseVariantPrice(ctx context.Context, cmd DecreaseVariantPrice) error {
	variant, err := h.variants.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := variant.DecreasePrice(cmd.ID, cmd.Price)
	if err != nil {
		return err
	}

	err = h.variants.Save(ctx, variant)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
