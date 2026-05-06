package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const WishlistAggregate = "wishlists.Wishlist"

type Wishlist struct {
	es.Aggregate
	UserID      string
	Name        string
	Description string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Wishlist)(nil)

func NewWishlist(id string) *Wishlist {
	return &Wishlist{
		Aggregate: es.NewAggregate(id, WishlistAggregate),
	}
}
func (w *Wishlist) InitWishlist(userID, name string) (ddd.Event, error) {

	w.AddEvent(WishlistCreatedEvent, &WishlistCreated{
		UserID: userID,
		Name:   name,
	})

	return ddd.NewEvent(WishlistCreatedEvent, w), nil
}

func (w *Wishlist) Remove() (ddd.Event, error) {
	w.AddEvent(WishlistRemovedEvent, &WishlistRemoved{
		WishlistID: w.ID(),
		UserID:     w.UserID,
	})

	return ddd.NewEvent(WishlistRemovedEvent, w), nil
}

// Key implements registry.Registerable
func (Wishlist) Key() string { return WishlistAggregate }

// ApplyEvent implements es.EventApplier
func (w *Wishlist) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *WishlistCreated:
		w.UserID = payload.UserID
		w.Name = payload.Name

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", w, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (w *Wishlist) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *WishlistV1:

		w.UserID = ss.UserID
		w.Name = ss.Name
	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", w, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (w Wishlist) ToSnapshot() es.Snapshot {
	return WishlistV1{}
}
