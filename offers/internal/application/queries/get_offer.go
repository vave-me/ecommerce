package queries

import (
	"context"
	"middleman/offers/internal/domain"
)

type GetOffer struct {
	OfferID string
}

type GetOfferHandler struct {
	offers domain.OfferRepository
}

func NewGetOfferHandler(offers domain.OfferRepository) GetOfferHandler {
	return GetOfferHandler{offers: offers}
}

func (h GetOfferHandler) GetOffer(ctx context.Context, query GetOffer) (*domain.Offer, error) {
	// Load the offer from the event store to get the actual state
	offer, err := h.offers.Load(ctx, query.OfferID)
	if err != nil {
		return nil, err
	}
	
	return offer, nil
}