package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type CancelLease struct {
	LeaseID string
}

type CancelLeaseHandler struct {
	leases    domain.LeaseRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCancelLeaseHandler(leases domain.LeaseRepository, publisher ddd.EventPublisher[ddd.Event]) CancelLeaseHandler {
	return CancelLeaseHandler{
		leases:    leases,
		publisher: publisher,
	}
}

func (h CancelLeaseHandler) CancelLease(ctx context.Context, cmd CancelLease) error {
	lease, err := h.leases.Load(ctx, cmd.LeaseID)
	if err != nil {
		return errors.Wrap(err, "loading aggregator in CancelLease")
	}
	evt, err := lease.CancelLease()
	if err != nil {
		return errors.Wrap(err, "CancelLease aggregator method")
	}

	if err := h.leases.Save(ctx, lease); err != nil {
		return errors.Wrap(err, "saving aggregator after cancel")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing cancel event")
}
