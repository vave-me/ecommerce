package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type AcceptOffer struct {
	ID             string
	UserCustomerID string
}

type AcceptOfferHandler struct {
	offers    domain.OfferRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAcceptOfferHandler(
	offers domain.OfferRepository,
	publisher ddd.EventPublisher[ddd.Event],
) AcceptOfferHandler {
	return AcceptOfferHandler{offers, publisher}
}

func (h AcceptOfferHandler) AcceptOffer(ctx context.Context, cmd AcceptOffer) error {
	agg, err := h.offers.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	evt, err := agg.AcceptOffer(cmd.ID, cmd.UserCustomerID)
	if err != nil {
		return err
	}

	if err = h.offers.Save(ctx, agg); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, evt)
}
