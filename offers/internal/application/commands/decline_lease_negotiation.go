package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type DeclineLeaseNegotiation struct {
	ID     string
	Reason string
}

type DeclineLeaseNegotiationHandler struct {
	leases    domain.LeaseRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDeclineLeaseNegotiationHandler(
	leases domain.LeaseRepository,
	publisher ddd.EventPublisher[ddd.Event],
) DeclineLeaseNegotiationHandler {
	return DeclineLeaseNegotiationHandler{leases, publisher}
}

func (h DeclineLeaseNegotiationHandler) DeclineLeaseNegotiation(ctx context.Context, cmd DeclineLeaseNegotiation) error {
	lease, err := h.leases.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	evt, err := lease.DeclineLease(cmd.ID, cmd.Reason)
	if err != nil {
		return err
	}

	if err = h.leases.Save(ctx, lease); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, evt)
}
