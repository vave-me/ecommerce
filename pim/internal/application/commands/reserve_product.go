package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/products/internal/domain"

	"github.com/stackus/errors"
)

//ITEM FREE RESERVEATION
//ITEM PAID RESERVATION - 10e,
//IF USER NOT BUYING HE IS LOOSING MONEY

// ReserveProduct reserves (holds) a specified quantity of a product's stock.
// If the product is not stock-managed (ManageStock = false) the command is a no-op.
// A reservation is treated as a temporary stock decrement.  When an order is
// finally shipped or cancelled a subsequent command should either commit or
// release the stock.
//
// NOTE: For now we reuse the existing ProductStockAdjustedEvent – a dedicated
// ProductReserved event can be introduced later if the workflow needs to
// distinguish between a reservation and a permanent deduction.
//
// Validation rules
//   • quantity must be > 0
//   • product.ManageStock must be true OR we silently ignore
//   • product.Stock must have sufficient units available
//
// On success: product.Stock is decreased by Quantity and the event is
// published via the supplied EventPublisher.
//
// On failure: returns an error (errors.ErrBadRequest when stock is
// insufficient).

type ReserveProduct struct {
	ProductID string
	Quantity  int64
}

type ReserveProductHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewReserveProductHandler(products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) ReserveProductHandler {
	return ReserveProductHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h ReserveProductHandler) ReserveProduct(ctx context.Context, cmd ReserveProduct) error {
	if cmd.Quantity <= 0 {
		return errors.Wrap(errors.ErrBadRequest, "reserve quantity must be positive")
	}

	product, err := h.products.Load(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	// If stock is not managed we treat reservation as NOOP (always success)
	if !product.ManageStock {
		return nil
	}

	if product.Stock < cmd.Quantity {
		return errors.Wrap(errors.ErrBadRequest, "insufficient stock to reserve")
	}

	newStock := product.Stock - cmd.Quantity

	event, err := product.AdjustStock(newStock)
	if err != nil {
		return err
	}

	if err = h.products.Save(ctx, product); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
