package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/posts/internal/application/commands"
	"middleman/posts/internal/application/queries"
	"middleman/posts/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		AddPost(ctx context.Context, cmd commands.AddPost) error
		UpdatePost(ctx context.Context, cmd commands.UpdatePost) error
		RemovePost(ctx context.Context, cmd commands.RemovePost) error
		ArchivePost(ctx context.Context, cmd commands.ArchivePost) error
		AddPostThumbnail(ctx context.Context, cmd commands.AddPostThumbnail) error
		UpdatePostThumbnail(ctx context.Context, cmd commands.UpdatePostThumbnail) error
	}

	Queries interface {
		GetPost(ctx context.Context, query queries.GetPost) (*domain.CatalogPost, error)
		GetPosts(ctx context.Context, query queries.GetPosts) ([]*domain.CatalogPost, int64, error)
		GetUserPosts(ctx context.Context, query queries.GetUserPosts) ([]*domain.CatalogPost, int64, error)
		GetPublicCatalog(ctx context.Context, query queries.GetPublicCatalog) ([]*domain.CatalogPost, int64, error)
		GetPostsWithFilters(ctx context.Context, query queries.GetPostsWithFilters) ([]*domain.CatalogPost, int64, error)
		GetPostsByCategory(ctx context.Context, query queries.GetPostsByCategory) ([]*domain.CatalogPost, int64, error)
		GetPostsByCategorySlug(ctx context.Context, query queries.GetPostsByCategorySlug) ([]*domain.CatalogPost, int64, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.AddPostHandler
		commands.UpdatePostHandler
		commands.RemovePostHandler
		commands.ArchivePostHandler
		commands.AddPostThumbnailHandler
		commands.UpdatePostThumbnailHandler
	}

	appQueries struct {
		queries.GetPostHandler
		queries.GetPostsHandler
		queries.GetUserPostsHandler
		queries.GetPublicCatalogHandler
		queries.GetPostsWithFiltersHandler
		queries.GetPostsByCategoryHandler
		queries.GetPostsByCategorySlugHandler
	}
)

var _ App = (*Application)(nil)

func New(posts domain.PostRepository,
	catalog domain.CatalogRepository,
	publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		appCommands: appCommands{
			AddPostHandler:             commands.NewAddPostHandler(posts, publisher),
			UpdatePostHandler:          commands.NewUpdatePostHandler(posts, publisher),
			RemovePostHandler:          commands.NewRemovePostHandler(posts, publisher),
			ArchivePostHandler:         commands.NewArchivePostHandler(posts, publisher),
			AddPostThumbnailHandler:    commands.NewAddPostThumbnailHandler(posts, publisher),
			UpdatePostThumbnailHandler: commands.NewUpdatePostThumbnailHandler(posts, publisher),
		},
		appQueries: appQueries{
			GetPostHandler:                queries.NewGetPostHandler(catalog),
			GetPostsHandler:               queries.NewGetPostsHandler(catalog),
			GetUserPostsHandler:           queries.NewGetUserPostsHandler(catalog),
			GetPublicCatalogHandler:       queries.NewGetPublicCatalogHandler(catalog),
			GetPostsWithFiltersHandler:    queries.NewGetPostsWithFiltersHandler(catalog),
			GetPostsByCategoryHandler:     queries.NewGetPostsByCategoryHandler(catalog),
			GetPostsByCategorySlugHandler: queries.NewGetPostsByCategorySlugHandler(catalog),
		},
	}
}
