package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

// ArchiveProduct marks a product as archived or inactive.
type ArchiveProduct struct {
	ID           string
	UserSellerID string
}

type ArchiveProductHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewArchiveProductHandler(
	products domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ArchiveProductHandler {
	return ArchiveProductHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h ArchiveProductHandler) ArchiveProduct(ctx context.Context, cmd ArchiveProduct) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := product.Archive(cmd.UserSellerID) // domain: product.Archive(userSellerID)
	if err != nil {
		return err
	}

	if err = h.products.Save(ctx, product); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
