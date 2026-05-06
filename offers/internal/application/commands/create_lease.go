package commands

import (
	"context"
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type CreateLease struct {
	ID              string
	OfferID         string
	MonthlyPrice    int64
	LeaseTermMonths int64
	HasBuyout       bool
	BuyoutPrice     int64
	LeaseStartDate  time.Time
	LeaseEndDate    time.Time
	InitialStatus   domain.LeaseStatus
}

type CreateLeaseHandler struct {
	leases    domain.LeaseRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCreateLeaseHandler(
	leases domain.LeaseRepository,
	publisher ddd.EventPublisher[ddd.Event],
) CreateLeaseHandler {
	return CreateLeaseHandler{
		leases:    leases,
		publisher: publisher,
	}
}

func (h CreateLeaseHandler) CreateLease(ctx context.Context, cmd CreateLease) error {
	// Load or create aggregator
	lease, err := h.leases.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading lease aggregator")
	}

	evt, err := lease.CreateLease(
		cmd.OfferID,
		cmd.MonthlyPrice,
		cmd.LeaseTermMonths,
		cmd.HasBuyout,
		cmd.BuyoutPrice,
		cmd.LeaseStartDate,
		cmd.LeaseEndDate,
		cmd.InitialStatus,
	)
	if err != nil {
		return err
	}

	if err := h.leases.Save(ctx, lease); err != nil {
		return errors.Wrap(err, "saving lease aggregator")
	}
	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing domain event")
}
