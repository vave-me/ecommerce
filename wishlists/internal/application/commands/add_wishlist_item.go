package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/wishlists/internal/domain"
)

type AddWishlistItem struct {
	ID         string
	WishlistID string
	ItemID     string
	EntityType string
}

type AddWishlistItemHandler struct {
	wishlistItems domain.WishlistItemRepository
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewAddWishlistItemsHandler(wishlistItems domain.WishlistItemRepository, publisher ddd.EventPublisher[ddd.Event]) AddWishlistItemHandler {
	return AddWishlistItemHandler{
		wishlistItems: wishlistItems,
		publisher:     publisher,
	}
}

func (h *AddWishlistItemHandler) AddWishlistItem(ctx context.Context, cmd AddWishlistItem) error {
	wishlistItem, err := h.wishlistItems.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding item to the wishlist")
	}

	event, err := wishlistItem.InitWishlistItem(cmd.ID, cmd.WishlistID, cmd.ItemID, cmd.EntityType)
	if err != nil {
		return errors.Wrap(err, "initializing wishlist item")
	}

	err = h.wishlistItems.Save(ctx, wishlistItem)
	if err != nil {
		return errors.Wrap(err, "error adding wishlist item")
	}

	return errors.Wrap(h.publisher.Publish(ctx, event), "publishing domain event")
}
