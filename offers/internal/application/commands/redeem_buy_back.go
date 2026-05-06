package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type RedeemBuyBack struct {
	BuyBackID string
}

type RedeemBuyBackHandler struct {
	buyBacks  domain.BuyBackRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRedeemBuyBackHandler(buyBacks domain.BuyBackRepository, publisher ddd.EventPublisher[ddd.Event]) RedeemBuyBackHandler {
	return RedeemBuyBackHandler{
		buyBacks:  buyBacks,
		publisher: publisher,
	}
}

func (h RedeemBuyBackHandler) RedeemBuyBack(ctx context.Context, cmd RedeemBuyBack) error {
	buyBack, err := h.buyBacks.Load(ctx, cmd.BuyBackID)
	if err != nil {
		return errors.Wrap(err, "loading buyBack aggregator")
	}

	evt, err := buyBack.RedeemBuyBack()
	if err != nil {
		return errors.Wrap(err, "RedeemBuyBack aggregator method")
	}

	if err := h.buyBacks.Save(ctx, buyBack); err != nil {
		return errors.Wrap(err, "saving aggregator after redemption")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing redeem event")
}
