package domain

import (
	"context"
)

type OfferRepository interface {
	Load(ctx context.Context, offerID string) (*Offer, error)
	Save(ctx context.Context, offer *Offer) error
}
