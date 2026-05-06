package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type ConfirmBuyNow struct {
	ID string // aggregator ID
}

type ConfirmBuyNowHandler struct {
	buyNowRepo domain.BuyNowRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewConfirmBuyNowHandler(buyNowRepo domain.BuyNowRepository, publisher ddd.EventPublisher[ddd.Event]) ConfirmBuyNowHandler {
	return ConfirmBuyNowHandler{
		buyNowRepo: buyNowRepo,
		publisher:  publisher,
	}
}

func (h ConfirmBuyNowHandler) ConfirmBuyNow(ctx context.Context, cmd ConfirmBuyNow) error {
	buyNowAgg, err := h.buyNowRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading buyNow aggregator")
	}

	evt, err := buyNowAgg.ConfirmBuyNow()
	if err != nil {
		return errors.Wrap(err, "ConfirmBuyNow aggregator method")
	}

	if err := h.buyNowRepo.Save(ctx, buyNowAgg); err != nil {
		return errors.Wrap(err, "saving aggregator after confirm")
	}

	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing buy now confirmed event")
}
