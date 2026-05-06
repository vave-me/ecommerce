package commands

import (
	"context"
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type StartLease struct {
	ID        string
	LeaseID   string
	StartDate time.Time
}

type StartLeaseHandler struct {
	leases    domain.LeaseRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewStartLeaseHandler(
	leases domain.LeaseRepository,
	publisher ddd.EventPublisher[ddd.Event],
) StartLeaseHandler {
	return StartLeaseHandler{
		leases:    leases,
		publisher: publisher,
	}
}

func (h StartLeaseHandler) StartLease(ctx context.Context, cmd StartLease) error {
	lease, err := h.leases.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading lease aggregator")
	}

	evt, err := lease.StartLease(cmd.LeaseID, cmd.StartDate)
	if err != nil {
		return err
	}

	if err := h.leases.Save(ctx, lease); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, evt)
}
