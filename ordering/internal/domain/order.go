package domain

import (
	"fmt"
	"middleman/internal/ddd"
	"middleman/internal/es"

	"github.com/stackus/errors"
)

// The Aggregate name.
const OrderAggregate = "ordering.Order"

// Potential domain errors.
var (
	ErrOrderAlreadyCreated          = errors.Wrap(errors.ErrBadRequest, "the order cannot be recreated")
	ErrOrderHasNoItems              = errors.Wrap(errors.ErrBadRequest, "the order has no items")
	ErrOrderCannotBeCancelled       = errors.Wrap(errors.ErrBadRequest, "the order cannot be cancelled in this status")
	ErrCustomerIDCannotBeBlank      = errors.Wrap(errors.ErrBadRequest, "the customer id cannot be blank")
	ErrPaymentMethodIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the payment id cannot be blank")
)

type Order struct {
	es.Aggregate
	UserCustomerID  string
	PaymentMethodID string
	InvoiceID       string
	ShoppingID      string
	BasketID        string
	PaymentIntent   string
	Items           []Item
	Status          OrderStatus
}

// Ensure Order implements ES interfaces.
var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Order)(nil)

// NewOrder is a constructor if you want a partially constructed order in memory (e.g. for creation).
func NewOrder(id string) *Order {
	return &Order{
		Aggregate: es.NewAggregate(id, OrderAggregate),
		Status:    OrderUnknown,
	}
}

// Key returns the aggregate name (for ES frameworks).
func (Order) Key() string { return OrderAggregate }

// CreateOrder triggers an OrderCreatedEvent if the order is in unknown status.
func (o *Order) CreateOrder(id, userCustomerID, basketID, paymentIntent string, items []Item) (ddd.Event, error) {
	if o.Status != OrderUnknown {
		return nil, ErrOrderAlreadyCreated
	}
	if len(items) == 0 {
		return nil, ErrOrderHasNoItems
	}
	if userCustomerID == "" {
		fmt.Println("user customer id :", userCustomerID)
		return nil, ErrCustomerIDCannotBeBlank
	}

	// Fire the event
	o.AddEvent(OrderCreatedEvent, &OrderCreated{
		UserCustomerID: userCustomerID,
		BasketID:       basketID,
		Items:          items,
		PaymentIntent:  paymentIntent,
	})
	return ddd.NewEvent(OrderCreatedEvent, o), nil
}

// Reject sets the order to REJECTED from PENDING, presumably if payment or other checks failed.
func (o *Order) Reject() (ddd.Event, error) {
	// Possibly check if status == PENDING
	if o.Status != OrderIsPending {
		return nil, errors.Wrap(errors.ErrBadRequest, "only a pending order can be rejected")
	}
	o.AddEvent(OrderRejectedEvent, &OrderRejected{})
	return ddd.NewEvent(OrderRejectedEvent, o), nil
}

// Approve sets order to APPROVED from PENDING, presumably after payment or stock checks.
func (o *Order) Approve(shoppingID string) (ddd.Event, error) {
	// Check if status == PENDING
	if o.Status != OrderIsPending {
		return nil, errors.Wrap(errors.ErrBadRequest, "only a pending order can be approved")
	}
	o.AddEvent(OrderApprovedEvent, &OrderApproved{
		ShoppingID: shoppingID,
	})
	return ddd.NewEvent(OrderApprovedEvent, o), nil
}

// Cancel sets order to CANCELED if it's in a valid state to do so.
func (o *Order) Cancel() (ddd.Event, error) {
	// If you allow cancellation from PENDING, APPROVED, etc.
	// If it's REJECTED or COMPLETED, you might disallow
	if o.Status == OrderIsRejected || o.Status == OrderIsCompleted {
		return nil, ErrOrderCannotBeCancelled
	}

	o.AddEvent(OrderCanceledEvent, &OrderCanceled{
		UserCustomerID:  o.UserCustomerID,
		PaymentMethodID: o.PaymentMethodID,
	})
	return ddd.NewEvent(OrderCanceledEvent, o), nil
}

// Ready sets the order to READY from APPROVED (meaning items are picked/packed).
func (o *Order) Ready() (ddd.Event, error) {
	if o.Status != OrderIsApproved {
		return nil, errors.Wrap(errors.ErrBadRequest, "order must be approved before ready")
	}
	o.AddEvent(OrderReadiedEvent, &OrderReadied{
		UserCustomerID:  o.UserCustomerID,
		PaymentMethodID: o.PaymentMethodID,
		Total:           o.GetTotal(),
	})
	return ddd.NewEvent(OrderReadiedEvent, o), nil
}

// Ship is a NEW method to handle shipping transition from READY -> SHIPPED
func (o *Order) Ship() (ddd.Event, error) {
	if o.Status != OrderIsReady {
		return nil, errors.Wrap(errors.ErrBadRequest, "order must be READY to be shipped")
	}
	o.AddEvent(OrderShippedEvent, &OrderShipped{})
	return ddd.NewEvent(OrderShippedEvent, o), nil
}

// Deliver is a NEW method to handle final physical delivery from SHIPPED -> DELIVERED
func (o *Order) Deliver() (ddd.Event, error) {
	if o.Status != OrderIsShipped {
		return nil, errors.Wrap(errors.ErrBadRequest, "order must be SHIPPED to be delivered")
	}
	o.AddEvent(OrderDeliveredEvent, &OrderDelivered{})
	return ddd.NewEvent(OrderDeliveredEvent, o), nil
}

// Complete sets the order to COMPLETED from DELIVERED, or possibly from other states if business allows
func (o *Order) Complete(invoiceID string) (ddd.Event, error) {
	// e.g., check status == DELIVERED
	if o.Status != OrderIsDelivered {
		return nil, errors.Wrap(errors.ErrBadRequest, "order must be DELIVERED to be completed")
	}
	if invoiceID == "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "invoiceID cannot be blank when completing order")
	}
	o.AddEvent(OrderCompletedEvent, &OrderCompleted{
		UserCustomerID: o.UserCustomerID,
		InvoiceID:      invoiceID,
	})
	return ddd.NewEvent(OrderCompletedEvent, o), nil
}

// Helper to calculate total cost
func (o Order) GetTotal() int64 {
	var total int64
	for _, item := range o.Items {
		total += item.Price * int64(item.Quantity)
	}
	return total
}

// ---- ES: ApplyEvent ----

// ApplyEvent updates the aggregate's state based on the event's payload.
func (o *Order) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *OrderCreated:
		o.UserCustomerID = payload.UserCustomerID
		o.BasketID = payload.BasketID
		o.Items = payload.Items
		o.PaymentIntent = payload.PaymentIntent
		o.Status = OrderIsPending

	case *OrderRejected:
		o.Status = OrderIsRejected

	case *OrderApproved:
		o.ShoppingID = payload.ShoppingID
		o.Status = OrderIsApproved

	case *OrderCanceled:
		o.Status = OrderIsCanceled

	case *OrderReadied:
		o.Status = OrderIsReady

	case *OrderShipped:
		o.Status = OrderIsShipped

	case *OrderDelivered:
		o.Status = OrderIsDelivered

	case *OrderCompleted:
		o.InvoiceID = payload.InvoiceID
		o.Status = OrderIsCompleted

	default:
		return errors.ErrInternal.
			Msgf("%T received an unknown payload type %T for event %s",
				o, payload, event.EventName())
	}
	return nil
}

// ---- ES: Snapshot ----

func (o *Order) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *OrderV1:
		o.UserCustomerID = ss.UserCustomerID
		o.PaymentMethodID = ss.PaymentMethodID
		o.InvoiceID = ss.InvoiceID
		o.ShoppingID = ss.ShoppingID
		o.Items = ss.Items
		o.Status = ss.Status
		o.PaymentIntent = ss.PaymentIntent
		o.BasketID = ss.BasketID

	default:
		return errors.ErrInternal.
			Msgf("%T received an unexpected snapshot %T", o, snapshot)
	}
	return nil
}

// ToSnapshot returns the current state as a snapshot.
func (o *Order) ToSnapshot() es.Snapshot {
	return &OrderV1{
		UserCustomerID:  o.UserCustomerID,
		PaymentMethodID: o.PaymentMethodID,
		InvoiceID:       o.InvoiceID,
		ShoppingID:      o.ShoppingID,
		Items:           o.Items,
		Status:          o.Status,
		PaymentIntent:   o.PaymentIntent,
		BasketID:        o.BasketID,
	}
}
