package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type CancelBuyNow struct {
	ID string // aggregator ID
}

type CancelBuyNowHandler struct {
	buyNowRepo domain.BuyNowRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewCancelBuyNowHandler(buyNowRepo domain.BuyNowRepository, publisher ddd.EventPublisher[ddd.Event]) CancelBuyNowHandler {
	return CancelBuyNowHandler{
		buyNowRepo: buyNowRepo,
		publisher:  publisher,
	}
}

func (h CancelBuyNowHandler) CancelBuyNow(ctx context.Context, cmd CancelBuyNow) error {
	buyNowAgg, err := h.buyNowRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in CancelBuyNow")
	}

	evt, err := buyNowAgg.CancelBuyNow()
	if err != nil {
		return errors.Wrap(err, "CancelBuyNow aggregator method")
	}

	if err := h.buyNowRepo.Save(ctx, buyNowAgg); err != nil {
		return errors.Wrap(err, "saving aggregator after cancel")
	}

	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing buy now canceled event")
}
