package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type RequestLeaseNegotiation struct {
	ID      string
	Message string
}

type RequestLeaseNegotiationHandler struct {
	leases    domain.LeaseRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRequestLeaseNegotiationHandler(
	leases domain.LeaseRepository,
	publisher ddd.EventPublisher[ddd.Event],
) RequestLeaseNegotiationHandler {
	return RequestLeaseNegotiationHandler{leases, publisher}
}

func (h RequestLeaseNegotiationHandler) RequestLeaseNegotiation(ctx context.Context, cmd RequestLeaseNegotiation) error {
	lease, err := h.leases.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	evt, err := lease.RequestNegotiation(cmd.ID, cmd.ID, 1, 1)
	if err != nil {
		return err
	}

	if err = h.leases.Save(ctx, lease); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, evt)
}
