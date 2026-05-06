package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type RequestReservationNegotiation struct {
	ReservationID string
	Comments      string
}

type RequestReservationNegotiationHandler struct {
	buyBacks  domain.ReservationRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRequestReservationNegotiationHandler(buyBacks domain.ReservationRepository, publisher ddd.EventPublisher[ddd.Event]) RequestReservationNegotiationHandler {
	return RequestReservationNegotiationHandler{
		buyBacks:  buyBacks,
		publisher: publisher,
	}
}

func (h RequestReservationNegotiationHandler) RequestReservationNegotiation(ctx context.Context, cmd RequestReservationNegotiation) error {
	buyBack, err := h.buyBacks.Load(ctx, cmd.ReservationID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in RequestReservationNegotiation")
	}

	evt, err := buyBack.RequestNegotiation(cmd.Comments)
	if err != nil {
		return errors.Wrap(err, "aggregator RequestNegotiation method")
	}

	if err := h.buyBacks.Save(ctx, buyBack); err != nil {
		return errors.Wrap(err, "saving aggregator after negotiation request")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing negotiation requested event")
}
