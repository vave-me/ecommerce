package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type MarkProductSold struct {
	ID           string
	UserSellerID string
	FinalPrice   int64 // optional
}

type MarkProductSoldHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewMarkProductSoldHandler(
	products domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
) MarkProductSoldHandler {
	return MarkProductSoldHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h MarkProductSoldHandler) MarkProductSold(ctx context.Context, cmd MarkProductSold) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := product.MarkSold(cmd.UserSellerID, cmd.FinalPrice)
	// domain: product.MarkSold(...)
	if err != nil {
		return err
	}

	if err = h.products.Save(ctx, product); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
