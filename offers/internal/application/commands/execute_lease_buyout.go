package commands

import (
	"context"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type ExecuteLeaseBuyout struct {
	ID          string // aggregator ID
	LeaseID     string
	BuyoutPrice int64
}

type ExecuteLeaseBuyoutHandler struct {
	leases    domain.LeaseRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewExecuteLeaseBuyoutHandler(
	leases domain.LeaseRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ExecuteLeaseBuyoutHandler {
	return ExecuteLeaseBuyoutHandler{
		leases:    leases,
		publisher: publisher,
	}
}

func (h ExecuteLeaseBuyoutHandler) ExecuteLeaseBuyout(ctx context.Context, cmd ExecuteLeaseBuyout) error {
	lease, err := h.leases.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading lease aggregator")
	}

	evt, err := lease.ExecuteLeaseBuyout(cmd.LeaseID, cmd.BuyoutPrice)
	if err != nil {
		return err
	}

	if err := h.leases.Save(ctx, lease); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, evt)
}
