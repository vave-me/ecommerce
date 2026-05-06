package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type ExpireBuyBack struct {
	BuyBackID string
}

type ExpireBuyBackHandler struct {
	buyBacks  domain.BuyBackRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewExpireBuyBackHandler(buyBacks domain.BuyBackRepository, publisher ddd.EventPublisher[ddd.Event]) ExpireBuyBackHandler {
	return ExpireBuyBackHandler{
		buyBacks:  buyBacks,
		publisher: publisher,
	}
}

func (h ExpireBuyBackHandler) ExpireBuyBack(ctx context.Context, cmd ExpireBuyBack) error {
	buyBack, err := h.buyBacks.Load(ctx, cmd.BuyBackID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in ExpireBuyBack")
	}
	evt, err := buyBack.ExpireBuyBack()
	if err != nil {
		return errors.Wrap(err, "ExpireBuyBack aggregator method")
	}

	if err := h.buyBacks.Save(ctx, buyBack); err != nil {
		return errors.Wrap(err, "saving aggregator after expire")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing expire event")
}
