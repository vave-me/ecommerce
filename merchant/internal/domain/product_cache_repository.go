package domain

import (
	"context"
)

type ProductCacheRepository interface {
	Add(ctx context.Context, productID, name, description string, price int64, userSellerID string, stock int64, sku string, categoryID string) error
	Rebrand(ctx context.Context, productID, name, description string, price int64, stock int64, sku string, categoryID string) error
	UpdatePrice(ctx context.Context, productID string, delta int64) error
	Remove(ctx context.Context, productID string) error
	ProductRepository
}
