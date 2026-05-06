package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type AcceptLeaseNegotiation struct {
	ID             string
	NewPrice       int64
	UserCustomerID string
}

type AcceptLeaseNegotiationHandler struct {
	leases    domain.LeaseRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAcceptLeaseNegotiationHandler(
	leases domain.LeaseRepository,
	publisher ddd.EventPublisher[ddd.Event],
) AcceptLeaseNegotiationHandler {
	return AcceptLeaseNegotiationHandler{leases, publisher}
}

func (h AcceptLeaseNegotiationHandler) AcceptLeaseNegotiation(ctx context.Context, cmd AcceptLeaseNegotiation) error {
	lease, err := h.leases.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	evt, err := lease.AcceptNegotiation(cmd.NewPrice, cmd.UserCustomerID)
	if err != nil {
		return err
	}

	if err = h.leases.Save(ctx, lease); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, evt)
}
