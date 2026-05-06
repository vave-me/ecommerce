package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type CloseOffer struct {
	ID     string // OfferID
	Reason string
}

type CloseOfferHandler struct {
	offers    domain.OfferRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCloseOfferHandler(
	offers domain.OfferRepository,
	publisher ddd.EventPublisher[ddd.Event],
) CloseOfferHandler {
	return CloseOfferHandler{offers, publisher}
}

func (h CloseOfferHandler) CloseOffer(ctx context.Context, cmd CloseOffer) error {
	agg, err := h.offers.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	evt, err := agg.CloseOffer(cmd.ID, cmd.Reason)
	if err != nil {
		return err
	}

	if err = h.offers.Save(ctx, agg); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, evt)
}
