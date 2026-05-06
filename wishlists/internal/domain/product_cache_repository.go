package domain

import (
	"context"
)

type ProductCacheRepository interface {
	Add(ctx context.Context, itemID, name, description string, price int64, userSellerID string, stock int64, sku string, categoryID string) error
	Rebrand(ctx context.Context, itemID, name string) error
	UpdatePrice(ctx context.Context, itemID string, delta int64) error
	Remove(ctx context.Context, itemID string) error
	ProductRepository
}
