package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type CreateReservation struct {
	ReservationID    string
	OfferID          string
	LockedPrice      int64
	RedemptionFee    int64
	LockBuyerID      string
	LockDurationDays int
}

type CreateReservationHandler struct {
	buyBacks  domain.ReservationRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCreateReservationHandler(buyBacks domain.ReservationRepository, publisher ddd.EventPublisher[ddd.Event]) CreateReservationHandler {
	return CreateReservationHandler{
		buyBacks:  buyBacks,
		publisher: publisher,
	}
}

func (h CreateReservationHandler) CreateReservation(ctx context.Context, cmd CreateReservation) error {
	buyBack, err := h.buyBacks.Load(ctx, cmd.ReservationID)
	if err != nil {
		return errors.Wrap(err, "loading buyBack aggregator")
	}

	evt, err := buyBack.CreateReservation(cmd.OfferID, cmd.LockedPrice, cmd.RedemptionFee, cmd.LockBuyerID, cmd.LockDurationDays)
	if err != nil {
		return errors.Wrap(err, "CreateReservation aggregator method")
	}

	if saveErr := h.buyBacks.Save(ctx, buyBack); saveErr != nil {
		return errors.Wrap(saveErr, "saving buyBack aggregator")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing buyBack created event")
}
