package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type CancelReservation struct {
	ReservationID string
}

type CancelReservationHandler struct {
	buyBacks  domain.ReservationRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCancelReservationHandler(buyBacks domain.ReservationRepository, publisher ddd.EventPublisher[ddd.Event]) CancelReservationHandler {
	return CancelReservationHandler{
		buyBacks:  buyBacks,
		publisher: publisher,
	}
}

func (h CancelReservationHandler) CancelReservation(ctx context.Context, cmd CancelReservation) error {
	buyBack, err := h.buyBacks.Load(ctx, cmd.ReservationID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in CancelReservation")
	}
	evt, err := buyBack.CancelReservation()
	if err != nil {
		return errors.Wrap(err, "CancelReservation aggregator method")
	}

	if err := h.buyBacks.Save(ctx, buyBack); err != nil {
		return errors.Wrap(err, "saving aggregator after cancel")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing cancel event")
}
