package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/products/internal/domain"
)

type AddProduct struct {
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
	SKU              string
	Attributes       []domain.Attribute
	Weight           int64
	Height           int64
	Width            int64
	Depth            int64
	Status           domain.ProductStatus
	Negotiable       bool
	UserType         domain.UserType
	MiddlemanService bool
	ShippingCost     int64
	HasVariants      bool
	Options          []domain.Option
	Thumbnail        string
	Lat              float64
	Lng              float64
}

type AddProductHandler struct {
	products  domain.ProductRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddProductHandler(
	products domain.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) AddProductHandler {

	return AddProductHandler{
		products:  products,
		publisher: publisher,
	}
}

func (h AddProductHandler) AddProduct(ctx context.Context, cmd AddProduct) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading product")
	}

	event, err := product.InitProduct(cmd.ID, cmd.Name, cmd.Description, cmd.BasePrice, cmd.UserSellerID, cmd.CategoryID, cmd.CategorySlug, cmd.Brand, cmd.Condition, cmd.Model, cmd.Tags, cmd.ManageStock, cmd.Stock, cmd.SKU, cmd.Attributes, cmd.Weight, cmd.Height, cmd.Width, cmd.Depth, cmd.Status, cmd.Negotiable, cmd.UserType, cmd.MiddlemanService, cmd.ShippingCost, cmd.HasVariants, cmd.Options, cmd.Thumbnail, cmd.Lat, cmd.Lng)
	if err != nil {
		return errors.Wrap(err, "initializing product")
	}

	err = h.products.Save(ctx, product)
	if err != nil {
		return errors.Wrap(err, "error adding product")
	}

	return errors.Wrap(h.publisher.Publish(ctx, event), "publishing domain event")
}
