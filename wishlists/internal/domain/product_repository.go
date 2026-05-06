package domain

import (
	"context"
)

type ProductRepository interface {
	Find(ctx context.Context, itemID string) (*Product, error)
}
