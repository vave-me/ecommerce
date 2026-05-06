package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/shipping/internal/domain"
)

type (
	CreateShipping struct {
		ID              string
		OrderID         string
		BasketID        string
		ProductID       string
		TrackingNumber  string
		LabelUrl        string
		SenderName      string
		SenderAddress   string
		ReceiverName    string
		ReceiverAddress string
		Weight          string
		Dimensions      string
		ServiceType     string
	}

	CreateShippingHandler struct {
		shippings domain.ShippingRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewCreateShippingHandler(shippings domain.ShippingRepository, publisher ddd.EventPublisher[ddd.Event]) CreateShippingHandler {
	return CreateShippingHandler{
		shippings: shippings,
		publisher: publisher,
	}
}

func (h CreateShippingHandler) CreateShipping(ctx context.Context, cmd CreateShipping) error {

	shipping, err := h.shippings.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := shipping.InitShipment(cmd.ProductID, cmd.OrderID, cmd.BasketID, cmd.TrackingNumber, cmd.LabelUrl, cmd.SenderName, cmd.SenderAddress, cmd.ReceiverName, cmd.ReceiverAddress, cmd.Weight, cmd.Dimensions, cmd.ServiceType)
	if err != nil {
		return err
	}

	err = h.shippings.Save(ctx, shipping)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
