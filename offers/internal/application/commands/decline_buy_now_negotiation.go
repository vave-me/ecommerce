package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type DeclineBuyNowNegotiation struct {
	ID     string // aggregator ID
	Reason string
}

type DeclineBuyNowNegotiationHandler struct {
	buyNowRepo domain.BuyNowRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewDeclineBuyNowNegotiationHandler(buyNowRepo domain.BuyNowRepository, publisher ddd.EventPublisher[ddd.Event]) DeclineBuyNowNegotiationHandler {
	return DeclineBuyNowNegotiationHandler{
		buyNowRepo: buyNowRepo,
		publisher:  publisher,
	}
}

func (h DeclineBuyNowNegotiationHandler) DeclineBuyNowNegotiation(ctx context.Context, cmd DeclineBuyNowNegotiation) error {
	agg, err := h.buyNowRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in DeclineBuyNowNegotiation")
	}

	evt, err := agg.DeclineNegotiation(cmd.Reason)
	if err != nil {
		return errors.Wrap(err, "aggregator DeclineNegotiation method")
	}

	if err := h.buyNowRepo.Save(ctx, agg); err != nil {
		return errors.Wrap(err, "saving aggregator after negotiation decline")
	}

	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing buy now negotiation declined event")
}
