package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type AcceptBuyBackNegotiation struct {
	BuyBackID       string
	NegotiatedPrice int64
	UserCustomerID  string
}

type AcceptBuyBackNegotiationHandler struct {
	buyBacks  domain.BuyBackRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAcceptBuyBackNegotiationHandler(buyBacks domain.BuyBackRepository, publisher ddd.EventPublisher[ddd.Event]) AcceptBuyBackNegotiationHandler {
	return AcceptBuyBackNegotiationHandler{
		buyBacks:  buyBacks,
		publisher: publisher,
	}
}

func (h AcceptBuyBackNegotiationHandler) AcceptBuyBackNegotiation(ctx context.Context, cmd AcceptBuyBackNegotiation) error {
	buyBack, err := h.buyBacks.Load(ctx, cmd.BuyBackID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in AcceptBuyBackNegotiation")
	}

	evt, err := buyBack.AcceptNegotiation(cmd.NegotiatedPrice, cmd.UserCustomerID)
	if err != nil {
		return errors.Wrap(err, "aggregator AcceptNegotiation method")
	}

	if err := h.buyBacks.Save(ctx, buyBack); err != nil {
		return errors.Wrap(err, "saving aggregator after negotiation acceptance")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing negotiation accepted event")
}
