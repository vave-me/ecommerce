package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const WishlistItemAggregate = "wishlists.WishlistItem"

type WishlistItem struct {
	es.Aggregate
	WishlistID string
	ItemID     string
	EntityType string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*WishlistItem)(nil)

func NewWishlistItem(id string) *WishlistItem {
	return &WishlistItem{
		Aggregate: es.NewAggregate(id, WishlistItemAggregate),
	}
}

func (w *WishlistItem) InitWishlistItem(id, wishlistID, itemID, entityType string) (ddd.Event, error) {

	w.AddEvent(WishlistItemAddedEvent, &WishlistItemAdded{
		WishlistID: wishlistID,
		ItemID:     itemID,
		EntityType: entityType,
	})

	return ddd.NewEvent(WishlistItemAddedEvent, w), nil
}
func (w *WishlistItem) Remove(id string) (ddd.Event, error) {
	w.AddEvent(WishlistItemRemovedEvent, &WishlistItemRemoved{
		ID: id,
	})

	return ddd.NewEvent(WishlistItemRemovedEvent, w), nil
}

func (w *WishlistItem) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *WishlistItemAdded:
		w.WishlistID = payload.WishlistID
		w.ItemID = payload.ItemID
		w.EntityType = payload.EntityType
	case *WishlistItemRemoved:

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", w, event.EventName(), payload)
	}

	return nil
}

// Key implements registry.Registerable
func (WishlistItem) Key() string { return WishlistItemAggregate }

func (w *WishlistItem) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *WishlistItemV1:
		w.WishlistID = ss.WishlistID
		w.ItemID = ss.ItemID

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", w, snapshot)
	}

	return nil
}

func (w WishlistItem) ToSnapshot() es.Snapshot {
	return WishlistItemV1{
		WishlistID: w.WishlistID,
		ItemID:     w.ItemID,
	}
}
