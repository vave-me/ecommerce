package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type AcceptBuyNowNegotiation struct {
	ID              string // aggregator ID
	NegotiatedPrice int64
	UserCustomerID  string
}

type AcceptBuyNowNegotiationHandler struct {
	buyNowRepo domain.BuyNowRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewAcceptBuyNowNegotiationHandler(buyNowRepo domain.BuyNowRepository, publisher ddd.EventPublisher[ddd.Event]) AcceptBuyNowNegotiationHandler {
	return AcceptBuyNowNegotiationHandler{
		buyNowRepo: buyNowRepo,
		publisher:  publisher,
	}
}

func (h AcceptBuyNowNegotiationHandler) AcceptBuyNowNegotiation(ctx context.Context, cmd AcceptBuyNowNegotiation) error {
	agg, err := h.buyNowRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in AcceptBuyNowNegotiation")
	}

	evt, err := agg.AcceptNegotiation(cmd.NegotiatedPrice, cmd.UserCustomerID)
	if err != nil {
		return errors.Wrap(err, "aggregator AcceptNegotiation method")
	}

	if err := h.buyNowRepo.Save(ctx, agg); err != nil {
		return errors.Wrap(err, "saving aggregator after negotiation accept")
	}

	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing buy now negotiation accepted event")
}
