package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/wishlists/internal/domain"
)

type RemoveWishlistItem struct {
	ID string
}

type RemoveWishlistItemHandler struct {
	wishlistItems domain.WishlistItemRepository
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewRemoveWishlistItemHandler(wishlistItems domain.WishlistItemRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveWishlistItemHandler {
	return RemoveWishlistItemHandler{
		wishlistItems: wishlistItems,
		publisher:     publisher,
	}
}

func (h RemoveWishlistItemHandler) RemoveWishlistItem(ctx context.Context, cmd RemoveWishlistItem) error {
	wishlistItem, err := h.wishlistItems.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := wishlistItem.Remove(cmd.ID)
	if err != nil {
		return err
	}

	err = h.wishlistItems.Save(ctx, wishlistItem)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
