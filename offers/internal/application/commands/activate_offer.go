package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

// ActivateOffer is the command payload
type ActivateOffer struct {
	ID        string // OfferID
	ProductID string
	UserID    string
}

// ActivateOfferHandler handles activating an offer
type ActivateOfferHandler struct {
	offers    domain.OfferRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewActivateOfferHandler(
	offers domain.OfferRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ActivateOfferHandler {
	return ActivateOfferHandler{
		offers:    offers,
		publisher: publisher,
	}
}

func (h ActivateOfferHandler) ActivateOffer(ctx context.Context, cmd ActivateOffer) error {
	agg, err := h.offers.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	evt, err := agg.ActivateOffer(cmd.ProductID, cmd.UserID) // e.g. agg.Activate() in your domain.
	if err != nil {
		return err
	}

	if err = h.offers.Save(ctx, agg); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, evt)
}
