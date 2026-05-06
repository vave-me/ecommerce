package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/offers/internal/domain"
)

type CreateBuyNow struct {
	ID         string // aggregator ID
	OfferID    string
	FinalPrice int64
}

type CreateBuyNowHandler struct {
	buyNowRepo domain.BuyNowRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewCreateBuyNowHandler(buyNowRepo domain.BuyNowRepository, publisher ddd.EventPublisher[ddd.Event]) CreateBuyNowHandler {
	return CreateBuyNowHandler{
		buyNowRepo: buyNowRepo,
		publisher:  publisher,
	}
}

func (h CreateBuyNowHandler) CreateBuyNow(ctx context.Context, cmd CreateBuyNow) error {
	buyNowAgg, err := h.buyNowRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "loading buyNow aggregator")
	}

	evt, err := buyNowAgg.CreateBuyNow(cmd.OfferID, cmd.FinalPrice)
	if err != nil {
		return errors.Wrap(err, "CreateBuyNow aggregator method")
	}

	if err := h.buyNowRepo.Save(ctx, buyNowAgg); err != nil {
		return errors.Wrap(err, "saving buyNow aggregator")
	}

	return errors.Wrap(h.publisher.Publish(ctx, evt), "publishing buy now created event")
}
