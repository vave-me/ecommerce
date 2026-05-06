package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type MarkProductLeased struct {
	ID              string
	UserSellerID    string
	MonthlyPrice    int64
	LeaseTermMonths int64
}

type MarkProductLeasedHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewMarkProductLeasedHandler(
	products domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
) MarkProductLeasedHandler {
	return MarkProductLeasedHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h MarkProductLeasedHandler) MarkProductLeased(ctx context.Context, cmd MarkProductLeased) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := product.MarkLeased(cmd.UserSellerID, cmd.MonthlyPrice, cmd.LeaseTermMonths)
	// domain: product.MarkLeased
	if err != nil {
		return err
	}

	if err = h.products.Save(ctx, product); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
