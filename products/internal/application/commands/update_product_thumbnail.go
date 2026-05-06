package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type UpdateProductThumbnail struct {
	ID        string
	Thumbnail string
}

type UpdateProductThumbnailHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewUpdateProductThumbnailHandler(
	products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) UpdateProductThumbnailHandler {

	return UpdateProductThumbnailHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h UpdateProductThumbnailHandler) UpdateProductThumbnail(ctx context.Context, cmd UpdateProductThumbnail) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding product")
	}
	event, err := product.UpdateThumbnail(cmd.Thumbnail)
	if err != nil {
		return errors.Wrap(err, "initializing product")
	}

	err = h.products.Save(ctx, product)
	if err != nil {
		return errors.Wrap(err, "error adding product")
	}

	return h.publisher.Publish(ctx, event)
}
