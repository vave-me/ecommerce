package commands

import (
	"context"
	"fmt"
	"middleman/internal/ddd"
	"middleman/wishlists/internal/domain"
)

type (
	CreateWishlist struct {
		ID     string
		Name   string
		UserID string
	}

	CreateWishlistHandler struct {
		wishlists domain.WishlistRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewCreateWishlistHandler(wishlists domain.WishlistRepository, publisher ddd.EventPublisher[ddd.Event]) CreateWishlistHandler {
	return CreateWishlistHandler{
		wishlists: wishlists,
		publisher: publisher,
	}
}

func (h CreateWishlistHandler) CreateWishlist(ctx context.Context, cmd CreateWishlist) error {
	wishlist, err := h.wishlists.Load(ctx, cmd.ID)
	if err != nil {
		fmt.Println("Wishlist cannot be loaded")
		return err
	}

	event, err := wishlist.InitWishlist(cmd.UserID, cmd.Name)
	if err != nil {
		return err
	}

	err = h.wishlists.Save(ctx, wishlist)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
