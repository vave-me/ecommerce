package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
	"time"
)

type (
	MakeBuyNowPaid struct {
		ID               string
		OfferID          string
		MonthlyPrice     int64
		BuyNowTermMonths int64
		HasBuyout        bool
		BuyoutPrice      int64
		BuyNowStartDate  time.Time
		BuyNowEndDate    time.Time
		BuyNowStatus     domain.BuyNowStatus
		UserCustomerID   string
	}

	MakeBuyNowPaidHandler struct {
		leasing   domain.BuyNowRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewMakeBuyNowPaidHandler(leasing domain.BuyNowRepository, publisher ddd.EventPublisher[ddd.Event]) MakeBuyNowPaidHandler {
	return MakeBuyNowPaidHandler{
		leasing:   leasing,
		publisher: publisher,
	}
}

func (h MakeBuyNowPaidHandler) MakeBuyNowPaid(ctx context.Context, cmd MakeBuyNowPaid) error {
	lease, err := h.leasing.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := lease.ConfirmBuyNow()
	if err != nil {
		return err
	}

	err = h.leasing.Save(ctx, lease)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
