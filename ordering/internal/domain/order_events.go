package domain

const (
	OrderCreatedEvent   = "ordering.OrderCreated"
	OrderRejectedEvent  = "ordering.OrderRejected"
	OrderApprovedEvent  = "ordering.OrderApproved"
	OrderCanceledEvent  = "ordering.OrderCanceled"
	OrderReadiedEvent   = "ordering.OrderReadied"
	OrderCompletedEvent = "ordering.OrderCompleted"
	OrderShippedEvent   = "ordering.OrderShipped"
	OrderDeliveredEvent = "ordering.OrderDelivered"
)

type OrderCreated struct {
	UserCustomerID string
	BasketID       string
	ShoppingID     string
	Items          []Item
	PaymentIntent  string
}

func (OrderCreated) Key() string { return OrderCreatedEvent }

type OrderRejected struct{}

func (OrderRejected) Key() string { return OrderRejectedEvent }

type OrderApproved struct {
	ShoppingID string
}

func (OrderApproved) Key() string { return OrderApprovedEvent }

type OrderCanceled struct {
	UserCustomerID  string
	PaymentMethodID string
}

func (OrderCanceled) Key() string { return OrderCanceledEvent }

type OrderReadied struct {
	UserCustomerID  string
	PaymentMethodID string
	Total           int64
}

func (OrderReadied) Key() string { return OrderReadiedEvent }

type OrderCompleted struct {
	UserCustomerID string
	InvoiceID      string
}

func (OrderCompleted) Key() string { return OrderCompletedEvent }

// NEW events for shipping & delivery
type OrderShipped struct{}

func (OrderShipped) Key() string { return OrderShippedEvent }

type OrderDelivered struct{}

func (OrderDelivered) Key() string { return OrderDeliveredEvent }
