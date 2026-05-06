package domain

import (
	"context"
)

type ProductCacheRepository interface {
	Add(ctx context.Context, productID, userSellerD, name, description string, price int64) error
	Rebrand(ctx context.Context, productID, name string) error
	UpdatePrice(ctx context.Context, productID string, delta int64) error
	Remove(ctx context.Context, productID string) error
	ProductRepository
}
