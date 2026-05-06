package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/products/internal/domain"

	"github.com/stackus/errors"
)

// ReleaseProduct adds back previously reserved quantity (compensation).
// It is essentially AdjustStock(oldStock + quantity).

type ReleaseProduct struct {
	ProductID string
	Quantity  int64
}

type ReleaseProductHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewReleaseProductHandler(products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) ReleaseProductHandler {
	return ReleaseProductHandler{products: products, publisher: publisher}
}

func (h ReleaseProductHandler) ReleaseProduct(ctx context.Context, cmd ReleaseProduct) error {
	if cmd.Quantity <= 0 {
		return errors.Wrap(errors.ErrBadRequest, "release quantity must be positive")
	}

	product, err := h.products.Load(ctx, cmd.ProductID)
	if err != nil {
		return err
	}
	if !product.ManageStock {
		return nil // noop
	}

	newStock := product.Stock + cmd.Quantity

	evt, err := product.AdjustStock(newStock)
	if err != nil {
		return err
	}
	if err = h.products.Save(ctx, product); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, evt)
}
