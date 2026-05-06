package domain

import (
	"context"
)

type BuyBackRepository interface {
	Load(ctx context.Context, buyBackID string) (*BuyBack, error)
	Save(ctx context.Context, buyBack *BuyBack) error
}
