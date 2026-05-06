package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type DeclineReservationNegotiation struct {
	ReservationID string
	Reason        string
}

type DeclineReservationNegotiationHandler struct {
	reservations domain.ReservationRepository
	publisher    ddd.EventPublisher[ddd.Event]
}

func NewDeclineReservationNegotiationHandler(reservations domain.ReservationRepository, publisher ddd.EventPublisher[ddd.Event]) DeclineReservationNegotiationHandler {
	return DeclineReservationNegotiationHandler{
		reservations: reservations,
		publisher:    publisher,
	}
}

func (h DeclineReservationNegotiationHandler) DeclineReservationNegotiation(ctx context.Context, cmd DeclineReservationNegotiation) error {
	reservation, err := h.reservations.Load(ctx, cmd.ReservationID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in DeclineReservationNegotiation")
	}

	evt, err := reservation.DeclineNegotiation(cmd.Reason)
	if err != nil {
		return errors.Wrap(err, "aggregator DeclineNegotiation method")
	}

	if err := h.reservations.Save(ctx, reservation); err != nil {
		return errors.Wrap(err, "saving aggregator after negotiation decline")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing negotiation declined event")
}
