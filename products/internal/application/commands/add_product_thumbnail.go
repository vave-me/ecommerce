package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type AddProductThumbnail struct {
	ID        string
	Thumbnail string
}

type AddProductThumbnailHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddProductThumbnailHandler(
	products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) AddProductThumbnailHandler {

	return AddProductThumbnailHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h AddProductThumbnailHandler) AddProductThumbnail(ctx context.Context, cmd AddProductThumbnail) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding product")
	}

	event, err := product.AddThumbnail(cmd.Thumbnail)
	if err != nil {
		return errors.Wrap(err, "initializing product")
	}

	err = h.products.Save(ctx, product)
	if err != nil {
		return errors.Wrap(err, "error adding product")
	}

	return h.publisher.Publish(ctx, event)
}
