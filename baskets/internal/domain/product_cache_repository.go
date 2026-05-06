package domain

import (
	"context"
)

type ProductCacheRepository interface {
	Add(ctx context.Context,
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
		lat, long float64) error
	Rebrand(ctx context.Context, productID, name, description string, price int64, stock int64, sku, categoryID string) error
	UpdatePrice(ctx context.Context, productID string, oldPrice, newPrice int64) error
	Remove(ctx context.Context, productID string) error
	ProductRepository
}
