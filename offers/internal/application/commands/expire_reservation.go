package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type ExpireReservation struct {
	ReservationID string
}

type ExpireReservationHandler struct {
	buyBacks  domain.ReservationRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewExpireReservationHandler(buyBacks domain.ReservationRepository, publisher ddd.EventPublisher[ddd.Event]) ExpireReservationHandler {
	return ExpireReservationHandler{
		buyBacks:  buyBacks,
		publisher: publisher,
	}
}

func (h ExpireReservationHandler) ExpireReservation(ctx context.Context, cmd ExpireReservation) error {
	buyBack, err := h.buyBacks.Load(ctx, cmd.ReservationID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in ExpireReservation")
	}
	evt, err := buyBack.ExpireReservation()
	if err != nil {
		return errors.Wrap(err, "ExpireReservation aggregator method")
	}

	if err := h.buyBacks.Save(ctx, buyBack); err != nil {
		return errors.Wrap(err, "saving aggregator after expire")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing expire event")
}
