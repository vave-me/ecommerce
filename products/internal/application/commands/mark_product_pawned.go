package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type MarkProductPawned struct {
	ID            string
	UserSellerID  string
	LockedPrice   int64
	RedemptionFee int64
}

type MarkProductPawnedHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewMarkProductPawnedHandler(
	products domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
) MarkProductPawnedHandler {
	return MarkProductPawnedHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h MarkProductPawnedHandler) MarkProductPawned(ctx context.Context, cmd MarkProductPawned) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := product.MarkPawned(cmd.UserSellerID, cmd.LockedPrice, cmd.RedemptionFee)
	// domain: product.MarkPawned
	if err != nil {
		return err
	}

	if err = h.products.Save(ctx, product); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
