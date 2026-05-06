package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/ordering/internal/domain"
)

type UpdateOrderStatus struct {
	ID     string
	Status string
	Reason string
}

type UpdateOrderStatusHandler struct {
	orders    domain.OrderRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewUpdateOrderStatusHandler(orders domain.OrderRepository, publisher ddd.EventPublisher[ddd.Event]) UpdateOrderStatusHandler {
	return UpdateOrderStatusHandler{
		orders:    orders,
		publisher: publisher,
	}
}

func (h UpdateOrderStatusHandler) UpdateOrderStatus(ctx context.Context, cmd UpdateOrderStatus) error {
	order, err := h.orders.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Determine which transition to make based on target status
	targetStatus := domain.ToOrderStatus(cmd.Status)
	var event ddd.Event

	switch targetStatus {
	case domain.OrderIsApproved:
		// For approve, we need a shopping ID - for now, generate one
		shoppingID := "SHOP-" + cmd.ID
		event, err = order.Approve(shoppingID)
	case domain.OrderIsRejected:
		event, err = order.Reject()
	case domain.OrderIsCanceled:
		event, err = order.Cancel()
	case domain.OrderIsReady:
		event, err = order.Ready()
	case domain.OrderIsShipped:
		event, err = order.Ship()
	case domain.OrderIsDelivered:
		event, err = order.Deliver()
	case domain.OrderIsCompleted:
		// For complete, we need an invoice ID - for now, generate one
		invoiceID := "INV-" + cmd.ID
		event, err = order.Complete(invoiceID)
	default:
		return domain.ErrOrderCannotBeCancelled // Use as generic invalid transition error
	}

	if err != nil {
		return err
	}

	if err = h.orders.Save(ctx, order); err != nil {
		return err
	}

	// Publish the domain event if there is one
	if event != nil {
		return h.publisher.Publish(ctx, event)
	}

	return nil
}