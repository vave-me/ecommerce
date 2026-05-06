package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/wishlists/internal/domain"
)

type RemoveWishlist struct {
	ID string
}

type RemoveWishlistHandler struct {
	wishlists domain.WishlistRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRemoveWishlistHandler(wishlists domain.WishlistRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveWishlistHandler {
	return RemoveWishlistHandler{
		wishlists: wishlists,
		publisher: publisher,
	}
}

func (h RemoveWishlistHandler) RemoveWishlist(ctx context.Context, cmd RemoveWishlist) error {
	wishlist, err := h.wishlists.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := wishlist.Remove()
	if err != nil {
		return err
	}

	err = h.wishlists.Save(ctx, wishlist)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
