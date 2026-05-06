package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type RequestBuyBackNegotiation struct {
	BuyBackID string
	Comments  string
}

type RequestBuyBackNegotiationHandler struct {
	buyBacks  domain.BuyBackRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRequestBuyBackNegotiationHandler(buyBacks domain.BuyBackRepository, publisher ddd.EventPublisher[ddd.Event]) RequestBuyBackNegotiationHandler {
	return RequestBuyBackNegotiationHandler{
		buyBacks:  buyBacks,
		publisher: publisher,
	}
}

func (h RequestBuyBackNegotiationHandler) RequestBuyBackNegotiation(ctx context.Context, cmd RequestBuyBackNegotiation) error {
	buyBack, err := h.buyBacks.Load(ctx, cmd.BuyBackID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in RequestBuyBackNegotiation")
	}

	evt, err := buyBack.RequestNegotiation(cmd.Comments)
	if err != nil {
		return errors.Wrap(err, "aggregator RequestNegotiation method")
	}

	if err := h.buyBacks.Save(ctx, buyBack); err != nil {
		return errors.Wrap(err, "saving aggregator after negotiation request")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing negotiation requested event")
}
