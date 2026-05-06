package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type RedeemReservation struct {
	ReservationID string
}

type RedeemReservationHandler struct {
	reservations domain.ReservationRepository
	publisher    ddd.EventPublisher[ddd.Event]
}

func NewRedeemReservationHandler(reservations domain.ReservationRepository, publisher ddd.EventPublisher[ddd.Event]) RedeemReservationHandler {
	return RedeemReservationHandler{
		reservations: reservations,
		publisher:    publisher,
	}
}

func (h RedeemReservationHandler) RedeemReservation(ctx context.Context, cmd RedeemReservation) error {
	buyBack, err := h.reservations.Load(ctx, cmd.ReservationID)
	if err != nil {
		return errors.Wrap(err, "loading buyBack aggregator")
	}

	evt, err := buyBack.RedeemReservation()
	if err != nil {
		return errors.Wrap(err, "RedeemReservation aggregator method")
	}

	if err := h.reservations.Save(ctx, buyBack); err != nil {
		return errors.Wrap(err, "saving aggregator after redemption")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing redeem event")
}
