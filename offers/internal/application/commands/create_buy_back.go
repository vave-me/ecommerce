package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type CreateBuyBack struct {
	BuyBackID        string
	OfferID          string
	LockedPrice      int64
	RedemptionFee    int64
	LockBuyerID      string
	LockDurationDays int
}

type CreateBuyBackHandler struct {
	buyBacks  domain.BuyBackRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCreateBuyBackHandler(buyBacks domain.BuyBackRepository, publisher ddd.EventPublisher[ddd.Event]) CreateBuyBackHandler {
	return CreateBuyBackHandler{
		buyBacks:  buyBacks,
		publisher: publisher,
	}
}

func (h CreateBuyBackHandler) CreateBuyBack(ctx context.Context, cmd CreateBuyBack) error {
	buyBack, err := h.buyBacks.Load(ctx, cmd.BuyBackID)
	if err != nil {
		return errors.Wrap(err, "loading buyBack aggregator")
	}

	evt, err := buyBack.CreateBuyBack(cmd.OfferID, cmd.LockedPrice, cmd.RedemptionFee, cmd.LockBuyerID, cmd.LockDurationDays)
	if err != nil {
		return errors.Wrap(err, "CreateBuyBack aggregator method")
	}

	if saveErr := h.buyBacks.Save(ctx, buyBack); saveErr != nil {
		return errors.Wrap(saveErr, "saving buyBack aggregator")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing buyBack created event")
}
