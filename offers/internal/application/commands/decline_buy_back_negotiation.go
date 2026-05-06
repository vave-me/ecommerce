package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type DeclineBuyBackNegotiation struct {
	BuyBackID string
	Reason    string
}

type DeclineBuyBackNegotiationHandler struct {
	buyBacks  domain.BuyBackRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDeclineBuyBackNegotiationHandler(buyBacks domain.BuyBackRepository, publisher ddd.EventPublisher[ddd.Event]) DeclineBuyBackNegotiationHandler {
	return DeclineBuyBackNegotiationHandler{
		buyBacks:  buyBacks,
		publisher: publisher,
	}
}

func (h DeclineBuyBackNegotiationHandler) DeclineBuyBackNegotiation(ctx context.Context, cmd DeclineBuyBackNegotiation) error {
	buyBack, err := h.buyBacks.Load(ctx, cmd.BuyBackID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in DeclineBuyBackNegotiation")
	}

	evt, err := buyBack.DeclineNegotiation(cmd.Reason)
	if err != nil {
		return errors.Wrap(err, "aggregator DeclineNegotiation method")
	}

	if err := h.buyBacks.Save(ctx, buyBack); err != nil {
		return errors.Wrap(err, "saving aggregator after negotiation decline")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing negotiation declined event")
}
