package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

// buyer or potential user requests a negotiation
type RequestBuyNowNegotiation struct {
	ID       string // aggregator ID
	Comments string // what the buyer is requesting or reason
}

type RequestBuyNowNegotiationHandler struct {
	buyNowRepo domain.BuyNowRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewRequestBuyNowNegotiationHandler(buyNowRepo domain.BuyNowRepository, publisher ddd.EventPublisher[ddd.Event]) RequestBuyNowNegotiationHandler {
	return RequestBuyNowNegotiationHandler{
		buyNowRepo: buyNowRepo,
		publisher:  publisher,
	}
}

func (h RequestBuyNowNegotiationHandler) RequestBuyNowNegotiation(ctx context.Context, cmd RequestBuyNowNegotiation) error {
	agg, err := h.buyNowRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in RequestBuyNowNegotiation")
	}

	evt, err := agg.RequestNegotiation(cmd.Comments)
	if err != nil {
		return errors.Wrap(err, "aggregator RequestNegotiation method")
	}

	if err := h.buyNowRepo.Save(ctx, agg); err != nil {
		return errors.Wrap(err, "saving aggregator after negotiation request")
	}

	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing buy now negotiation requested event")
}
