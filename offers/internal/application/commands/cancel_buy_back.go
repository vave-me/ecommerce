package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type CancelBuyBack struct {
	BuyBackID string
}

type CancelBuyBackHandler struct {
	buyBacks  domain.BuyBackRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCancelBuyBackHandler(buyBacks domain.BuyBackRepository, publisher ddd.EventPublisher[ddd.Event]) CancelBuyBackHandler {
	return CancelBuyBackHandler{
		buyBacks:  buyBacks,
		publisher: publisher,
	}
}

func (h CancelBuyBackHandler) CancelBuyBack(ctx context.Context, cmd CancelBuyBack) error {
	buyBack, err := h.buyBacks.Load(ctx, cmd.BuyBackID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in CancelBuyBack")
	}
	evt, err := buyBack.CancelBuyBack()
	if err != nil {
		return errors.Wrap(err, "CancelBuyBack aggregator method")
	}

	if err := h.buyBacks.Save(ctx, buyBack); err != nil {
		return errors.Wrap(err, "saving aggregator after cancel")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing cancel event")
}
