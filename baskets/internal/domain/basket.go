package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"

	"github.com/stackus/errors"
)

const BasketAggregate = "baskets.Basket"

var (
	ErrBasketHasNoItems         = errors.Wrap(errors.ErrBadRequest, "the basket has no items")
	ErrBasketCannotBeModified   = errors.Wrap(errors.ErrBadRequest, "the basket cannot be modified")
	ErrBasketCannotBeCancelled  = errors.Wrap(errors.ErrBadRequest, "the basket cannot be cancelled")
	ErrQuantityCannotBeNegative = errors.Wrap(errors.ErrBadRequest, "the item quantity cannot be negative")
	ErrUserIDCannotBeBlank      = errors.Wrap(errors.ErrBadRequest, "the user id cannot be blank")
)

type Basket struct {
	es.Aggregate
	UserCustomerID  string
	Items           map[string]Item
	Status          BasketStatus
	OfferType       string
	PaymentIntentID string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Basket)(nil)

// NewBasket creates a fresh basket with an ID
func NewBasket(id string) *Basket {
	return &Basket{
		Aggregate: es.NewAggregate(id, BasketAggregate),
		Items:     make(map[string]Item),
	}
}

func (b *Basket) IsOpen() bool {
	return b.Status == BasketIsOpen
}

func (b *Basket) IsCancellable() bool {
	return b.Status == BasketIsOpen
}

// Start initializes the basket if it's not already open/checked out/canceled
func (b *Basket) Start(userCustomerID string) (ddd.Event, error) {
	if b.Status == BasketIsOpen || b.Status == BasketIsCheckedOut || b.Status == BasketIsCanceled {
		return nil, ErrBasketCannotBeModified
	}
	if userCustomerID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	b.AddEvent(BasketStartedEvent, &BasketStarted{
		UserCustomerID: userCustomerID,
		Status:         BasketIsOpen,
	})
	return ddd.NewEvent(BasketStartedEvent, b), nil
}

func (b *Basket) Cancel() (ddd.Event, error) {
	if !b.IsCancellable() {
		return nil, ErrBasketCannotBeCancelled
	}
	b.AddEvent(BasketCanceledEvent, &BasketCanceled{
		Status: BasketIsCanceled,
	})
	return ddd.NewEvent(BasketCanceledEvent, b), nil
}

func (b *Basket) Checkout(paymentIntentID string) (ddd.Event, error) {
	if !b.IsOpen() {
		return nil, ErrBasketCannotBeModified
	}
	if len(b.Items) == 0 {
		return nil, ErrBasketHasNoItems
	}

	b.AddEvent(BasketCheckedOutEvent, &BasketCheckedOut{
		Status:          BasketIsCheckedOut,
		PaymentIntentID: paymentIntentID,
	})
	return ddd.NewEvent(BasketCheckedOutEvent, b), nil
}

// Reopen allows a checked out basket to be reopened if checkout failed
func (b *Basket) Reopen(reason string) (ddd.Event, error) {
	if b.Status != BasketIsCheckedOut {
		return nil, ErrBasketCannotBeModified
	}

	b.AddEvent(BasketReopenedEvent, &BasketReopened{
		Status: BasketIsOpen,
		Reason: reason,
	})
	return ddd.NewEvent(BasketReopenedEvent, b), nil
}

// AddItem requires that the basket is open, and quantity >= 0
func (b *Basket) AddItem(user *User, product *Product, quantity int64) (ddd.Event, error) {
	if !b.IsOpen() {
		return nil, ErrBasketCannotBeModified
	}
	if quantity < 0 {
		return nil, ErrQuantityCannotBeNegative
	}

	b.AddEvent(BasketItemAddedEvent, &BasketItemAdded{
		Item: Item{
			UserSellerID:   user.ID,
			ProductID:      product.ID,
			UserSellerName: user.FirstName,
			ProductName:    product.Name,
			ProductPrice:   product.BasePrice,
			Quantity:       quantity,
		},
	})
	return ddd.NewEvent(BasketItemAddedEvent, b), nil
}

// RemoveItem requires the basket be open, quantity >=0. If the item doesn't exist, do nothing
func (b *Basket) RemoveItem(product *Product, quantity int64) (ddd.Event, error) {
	if !b.IsOpen() {
		return nil, ErrBasketCannotBeModified
	}
	if quantity < 0 {
		return nil, ErrQuantityCannotBeNegative
	}

	if _, exists := b.Items[product.ID]; exists {
		b.AddEvent(BasketItemRemovedEvent, &BasketItemRemoved{
			ProductID: product.ID,
			Quantity:  quantity,
		})
	}
	return ddd.NewEvent(BasketItemRemovedEvent, b), nil
}

// ApplyEvent mutates the Basket to incorporate the event payload
func (b *Basket) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *BasketStarted:
		b.UserCustomerID = payload.UserCustomerID
		b.Status = payload.Status

	case *BasketItemAdded:
		if item, exists := b.Items[payload.Item.ProductID]; exists {
			item.Quantity += payload.Item.Quantity
			b.Items[payload.Item.ProductID] = item
		} else {
			b.Items[payload.Item.ProductID] = payload.Item
		}

	case *BasketItemRemoved:
		if item, exists := b.Items[payload.ProductID]; exists {
			newQty := item.Quantity - payload.Quantity
			if newQty <= 0 {
				delete(b.Items, payload.ProductID)
			} else {
				item.Quantity = newQty
				b.Items[payload.ProductID] = item
			}
		}

	case *BasketCanceled:
		b.Items = make(map[string]Item)
		b.Status = payload.Status

	case *BasketCheckedOut:
		b.Status = payload.Status
		b.PaymentIntentID = payload.PaymentIntentID

	case *BasketReopened:
		b.Status = payload.Status
		b.PaymentIntentID = "" // Clear payment intent when reopening

	default:
		return errors.ErrInternal.
			Msgf("%T received the event %s with unexpected payload %T", b, event.EventName(), payload)
	}
	return nil
}

// TotalAmount sums productPrice * quantity for all items
func (b *Basket) TotalAmount() int64 {
	var total int64
	for _, itm := range b.Items {
		total += itm.ProductPrice * itm.Quantity
	}
	return total
}

// Snapshot / Restore
func (b *Basket) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *BasketV1:
		b.UserCustomerID = ss.UserCustomerID
		b.Items = ss.Items
		b.Status = ss.Status
		b.PaymentIntentID = ss.PaymentIntentID
	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", b, snapshot)
	}
	return nil
}

func (b *Basket) ToSnapshot() es.Snapshot {
	return &BasketV1{
		UserCustomerID:  b.UserCustomerID,
		Items:           b.Items,
		Status:          b.Status,
		PaymentIntentID: b.PaymentIntentID,
	}
}
