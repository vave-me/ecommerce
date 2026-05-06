package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

// AdjustProductStock is the command payload to change a product's stock.
type AdjustProductStock struct {
	ID       string // Product ID
	NewStock int64
}

type AdjustProductStockHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAdjustProductStockHandler(
	products domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
) AdjustProductStockHandler {
	return AdjustProductStockHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h AdjustProductStockHandler) AdjustProductStock(ctx context.Context, cmd AdjustProductStock) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := product.AdjustStock(cmd.NewStock) // domain method (to be implemented)
	if err != nil {
		return err
	}

	if err = h.products.Save(ctx, product); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
