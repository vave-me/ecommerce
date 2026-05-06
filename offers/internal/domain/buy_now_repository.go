package domain

import (
	"context"
)

type BuyNowRepository interface {
	Load(ctx context.Context, leaseID string) (*BuyNow, error)
	Save(ctx context.Context, lease *BuyNow) error
}
