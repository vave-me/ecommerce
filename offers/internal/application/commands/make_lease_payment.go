package commands

import (
	"context"
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type MakeLeasePayment struct {
	ID          string // aggregator ID
	LeaseID     string
	Amount      int64
	PaymentDate time.Time
}

type MakeLeasePaymentHandler struct {
	leases    domain.LeaseRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewMakeLeasePaymentHandler(
	leases domain.LeaseRepository,
	publisher ddd.EventPublisher[ddd.Event],
) MakeLeasePaymentHandler {
	return MakeLeasePaymentHandler{
		leases:    leases,
		publisher: publisher,
	}
}

func (h MakeLeasePaymentHandler) MakeLeasePayment(ctx context.Context, cmd MakeLeasePayment) error {
	lease, err := h.leases.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading lease aggregator")
	}

	evt, err := lease.MakeLeasePayment(cmd.LeaseID, cmd.Amount, cmd.PaymentDate)
	if err != nil {
		return err
	}

	if err := h.leases.Save(ctx, lease); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, evt)
}
