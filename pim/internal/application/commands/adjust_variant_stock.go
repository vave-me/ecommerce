package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type AdjustVariantStock struct {
	ID       string // variant ID
	NewStock int64
}

type AdjustVariantStockHandler struct {
	variants  domain.VariantRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAdjustVariantStockHandler(
	variants domain.VariantRepository,
	publisher ddd.EventPublisher[ddd.Event],
) AdjustVariantStockHandler {
	return AdjustVariantStockHandler{
		variants:  variants,
		publisher: publisher,
	}
}

func (h AdjustVariantStockHandler) AdjustVariantStock(ctx context.Context, cmd AdjustVariantStock) error {
	variant, err := h.variants.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := variant.AdjustStock(cmd.NewStock)
	if err != nil {
		return err
	}

	if err = h.variants.Save(ctx, variant); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
