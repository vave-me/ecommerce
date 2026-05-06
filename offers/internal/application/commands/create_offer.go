package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type (
	CreateOffer struct {
		ID           string
		UserSellerID string // ID of the seller creating the offer
		ProductID    string // ID of the item being offered
		Price        int64  // Discount or price offer

	}

	CreateOfferHandler struct {
		offers    domain.OfferRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewCreateOfferHandler(offers domain.OfferRepository, publisher ddd.EventPublisher[ddd.Event]) CreateOfferHandler {
	return CreateOfferHandler{
		offers:    offers,
		publisher: publisher,
	}
}

func (h CreateOfferHandler) CreateOffer(ctx context.Context, cmd CreateOffer) error {
	offer, err := h.offers.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := offer.CreateOffer(cmd.UserSellerID, cmd.ProductID, cmd.Price)
	if err != nil {
		return err
	}

	err = h.offers.Save(ctx, offer)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
