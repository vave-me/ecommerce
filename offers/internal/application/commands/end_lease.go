package commands

import (
	"context"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type EndLease struct {
	ID      string // aggregator ID
	LeaseID string
}

type EndLeaseHandler struct {
	leases    domain.LeaseRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewEndLeaseHandler(
	leases domain.LeaseRepository,
	publisher ddd.EventPublisher[ddd.Event],
) EndLeaseHandler {
	return EndLeaseHandler{
		leases:    leases,
		publisher: publisher,
	}
}

func (h EndLeaseHandler) EndLease(ctx context.Context, cmd EndLease) error {
	lease, err := h.leases.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading lease aggregator")
	}

	evt, err := lease.EndLease(cmd.LeaseID)
	if err != nil {
		return err
	}

	if err := h.leases.Save(ctx, lease); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, evt)
}
