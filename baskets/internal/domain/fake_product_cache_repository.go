package domain

import (
	"context"

	"github.com/stackus/errors"
)

type FakeProductCacheRepository struct {
	products map[string]*Product
}

// Confirm that FakeProductCacheRepository implements ProductCacheRepository
var _ ProductCacheRepository = (*FakeProductCacheRepository)(nil)

// NewFakeProductCacheRepository returns a new fake repo with an in-memory map
func NewFakeProductCacheRepository() *FakeProductCacheRepository {
	return &FakeProductCacheRepository{
		products: map[string]*Product{},
	}
}

// Add implements ProductCacheRepository.Add
func (r *FakeProductCacheRepository) Add(
	ctx context.Context,
	id, name, description string,
	basePrice int64,
	userSellerID, categoryID, brand string,
	condition ProductCondition,
	model string,
	tags []string,
	manageStock bool,
	stock int64,
	sku string,
	attributes []Attribute,
	weight, height, width, depth int64,
	status ProductStatus,
	negotiable bool,
	userType UserType,
	middlemanService bool,
	shippingCost int64,
	hasVariants bool,
	options []Option,
	thumbnail string,
	lat, long float64,
) error {
	r.products[id] = &Product{
		ID:               id,
		Name:             name,
		Description:      description,
		BasePrice:        basePrice,
		UserSellerID:     userSellerID,
		CategoryID:       categoryID,
		Brand:            brand,
		Condition:        condition,
		Model:            model,
		Tags:             tags,
		ManageStock:      manageStock,
		Stock:            stock,
		SKU:              sku,
		Attributes:       attributes,
		Weight:           weight,
		Height:           height,
		Width:            width,
		Depth:            depth,
		Status:           status,
		Negotiable:       negotiable,
		UserType:         userType,
		MiddlemanService: middlemanService,
		ShippingCost:     shippingCost,
		HasVariants:      hasVariants,
		Options:          options,
		Thumbnail:        thumbnail,
		Lat:              lat,
		Lng:              long,
	}
	return nil
}

// Rebrand implements ProductCacheRepository.Rebrand
func (r *FakeProductCacheRepository) Rebrand(
	ctx context.Context,
	productID, name, description string,
	price int64,
	stock int64,
	sku, categoryID string,
) error {
	if product, exists := r.products[productID]; exists {
		product.Name = name
		product.Description = description
		product.BasePrice = price
		product.Stock = stock
		product.SKU = sku
		product.CategoryID = categoryID
	}
	return nil
}

// UpdatePrice implements ProductCacheRepository.UpdatePrice
func (r *FakeProductCacheRepository) UpdatePrice(
	ctx context.Context,
	productID string,
	oldPrice, newPrice int64,
) error {
	if product, exists := r.products[productID]; exists {
		// One option is to do: product.BasePrice = product.BasePrice - oldPrice + newPrice
		// but typically you might just set product.BasePrice = newPrice
		// depending on your actual business logic.
		product.BasePrice = newPrice
	}
	return nil
}

// Remove implements ProductCacheRepository.Remove
func (r *FakeProductCacheRepository) Remove(ctx context.Context, productID string) error {
	delete(r.products, productID)
	return nil
}

// Find implements ProductCacheRepository.Find
func (r *FakeProductCacheRepository) Find(ctx context.Context, productID string) (*Product, error) {
	if product, exists := r.products[productID]; exists {
		return product, nil
	}
	return nil, errors.ErrNotFound.Msgf("product with id: `%s` does not exist", productID)
}

// Reset is a helper to clear out and re-seed the map with given products
func (r *FakeProductCacheRepository) Reset(products ...*Product) {
	r.products = make(map[string]*Product)
	for _, product := range products {
		r.products[product.ID] = product
	}
}
