package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type UpdateProduct struct {
	ID               string
	Name             string
	Description      string
	BasePrice        int64
	UserSellerID     string
	CategoryID       string
	CategorySlug     string
	Brand            string
	Condition        domain.ProductCondition
	Model            string
	Tags             []string
	ManageStock      bool
	Stock            int64
	Sku              string
	Attributes       []domain.Attribute
	Weight           int64
	Height           int64
	Width            int64
	Depth            int64
	Status           domain.ProductStatus
	Thumbnail        string
	Negotiable       bool
	UserType         domain.UserType
	MiddlemanService bool
	ShippingCost     int64
	HasVariants      bool
	Options          []domain.Option
	Lat              float64
	Lng              float64
}

type UpdateProductHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewUpdateProductHandler(
	products domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UpdateProductHandler {
	return UpdateProductHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h UpdateProductHandler) UpdateProduct(ctx context.Context, cmd UpdateProduct) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// domain method: product.Update(...) or partial setters
	event, err := product.Update(
		cmd.Name,
		cmd.Description,
		cmd.BasePrice,
		cmd.UserSellerID,
		cmd.CategoryID,
		cmd.CategorySlug,
		cmd.Brand,
		cmd.Condition,
		cmd.Model,
		cmd.Tags,
		cmd.ManageStock,
		cmd.Stock,
		cmd.Sku,
		cmd.Attributes,
		cmd.Weight,
		cmd.Height,
		cmd.Width,
		cmd.Depth,
		cmd.Status,
		cmd.Thumbnail,
		cmd.Negotiable,
		cmd.UserType,
		cmd.MiddlemanService,
		cmd.ShippingCost,
		cmd.HasVariants,
		cmd.Options,
		cmd.Lat,
		cmd.Lng,
	)
	if err != nil {
		return err
	}

	if err = h.products.Save(ctx, product); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
