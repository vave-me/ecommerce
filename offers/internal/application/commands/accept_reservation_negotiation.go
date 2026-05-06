package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type AcceptReservationNegotiation struct {
	ReservationID   string
	NegotiatedPrice int64
	UserCustomerID  string
}

type AcceptReservationNegotiationHandler struct {
	buyBacks  domain.ReservationRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAcceptReservationNegotiationHandler(buyBacks domain.ReservationRepository, publisher ddd.EventPublisher[ddd.Event]) AcceptReservationNegotiationHandler {
	return AcceptReservationNegotiationHandler{
		buyBacks:  buyBacks,
		publisher: publisher,
	}
}

func (h AcceptReservationNegotiationHandler) AcceptReservationNegotiation(ctx context.Context, cmd AcceptReservationNegotiation) error {
	buyBack, err := h.buyBacks.Load(ctx, cmd.ReservationID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in AcceptReservationNegotiation")
	}

	evt, err := buyBack.AcceptNegotiation(cmd.NegotiatedPrice, cmd.UserCustomerID)
	if err != nil {
		return errors.Wrap(err, "aggregator AcceptNegotiation method")
	}

	if err := h.buyBacks.Save(ctx, buyBack); err != nil {
		return errors.Wrap(err, "saving aggregator after negotiation acceptance")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing negotiation accepted event")
}
