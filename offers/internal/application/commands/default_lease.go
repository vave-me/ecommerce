package commands

import (
	"context"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type DefaultLease struct {
	ID      string // aggregator ID
	LeaseID string
	Reason  string
}

type DefaultLeaseHandler struct {
	leases    domain.LeaseRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDefaultLeaseHandler(
	leases domain.LeaseRepository,
	publisher ddd.EventPublisher[ddd.Event],
) DefaultLeaseHandler {
	return DefaultLeaseHandler{
		leases:    leases,
		publisher: publisher,
	}
}

func (h DefaultLeaseHandler) DefaultLease(ctx context.Context, cmd DefaultLease) error {
	lease, err := h.leases.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading lease aggregator")
	}

	evt, err := lease.DefaultLease(cmd.LeaseID, cmd.Reason)
	if err != nil {
		return err
	}

	if err := h.leases.Save(ctx, lease); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, evt)
}
