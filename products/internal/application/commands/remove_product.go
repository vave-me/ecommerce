package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type RemoveProduct struct {
	ID     string
	UserID string
}

type RemoveProductHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRemoveProductHandler(products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveProductHandler {
	return RemoveProductHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h RemoveProductHandler) RemoveProduct(ctx context.Context, cmd RemoveProduct) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := product.Remove(cmd.ID, cmd.UserID)
	if err != nil {
		return err
	}

	err = h.products.Save(ctx, product)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
