package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/wishlists/internal/application/commands"
	"middleman/wishlists/internal/application/queries"
	"middleman/wishlists/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		CreateWishlist(ctx context.Context, cmd commands.CreateWishlist) error
		RemoveWishlist(ctx context.Context, cmd commands.RemoveWishlist) error
		AddWishlistItem(ctx context.Context, cmd commands.AddWishlistItem) error
		RemoveWishlistItem(ctx context.Context, cmd commands.RemoveWishlistItem) error
	}
	Queries interface {
		GetWishlist(ctx context.Context, query queries.GetWishlist) (*domain.MiddlemanWishlist, error)
		GetWishlists(ctx context.Context, query queries.GetWishlists) ([]*domain.MiddlemanWishlist, error)
		GetWishlistItem(ctx context.Context, query queries.GetWishlistItem) (*domain.CatalogWishlistItem, error)
		GetWishlistItems(ctx context.Context, query queries.GetWishlistItems) ([]*domain.CatalogWishlistItem, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.CreateWishlistHandler
		commands.RemoveWishlistHandler
		commands.AddWishlistItemHandler
		commands.RemoveWishlistItemHandler
	}
	appQueries struct {
		queries.GetWishlistHandler
		queries.GetWishlistsHandler
		queries.GetWishlistItemHandler
		queries.GetWishlistItemsHandler
	}
)

var _ App = (*Application)(nil)

func New(wishlists domain.WishlistRepository, wishlistItems domain.WishlistItemRepository, catalog domain.CatalogRepository, middleman domain.MiddlemanRepository, publisher ddd.EventPublisher[ddd.Event]) *Application {

	return &Application{
		appCommands: appCommands{
			CreateWishlistHandler:     commands.NewCreateWishlistHandler(wishlists, publisher),
			AddWishlistItemHandler:    commands.NewAddWishlistItemsHandler(wishlistItems, publisher),
			RemoveWishlistHandler:     commands.NewRemoveWishlistHandler(wishlists, publisher),
			RemoveWishlistItemHandler: commands.NewRemoveWishlistItemHandler(wishlistItems, publisher),
		},
		appQueries: appQueries{
			GetWishlistHandler:      queries.NewGetWishlistHandler(middleman),
			GetWishlistsHandler:     queries.NewGetWishlistsHandler(middleman),
			GetWishlistItemHandler:  queries.NewGetWishlistItemHandler(catalog),
			GetWishlistItemsHandler: queries.NewGetWishlistItemsHandler(catalog),
		},
	}
}
