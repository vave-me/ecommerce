package application

import (
	"context"
	"middleman/following/internal/application/commands"
	"middleman/following/internal/application/queries"
	"middleman/following/internal/domain"
	"middleman/internal/ddd"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		AddFollow(ctx context.Context, cmd commands.AddFollow) error
		ApproveFollow(ctx context.Context, cmd commands.ApproveFollow) error
		FlagFollow(ctx context.Context, cmd commands.FlagFollow) error
		RejectFollow(ctx context.Context, cmd commands.RejectFollow) error
		RemoveFollow(ctx context.Context, cmd commands.RemoveFollow) error
	}
	Queries interface {
		GetFollow(ctx context.Context, query queries.GetFollow) (*domain.MiddlemanFollow, error)
		GetFollowing(ctx context.Context, query queries.GetFollowing) ([]*domain.MiddlemanFollow, error)
		GetFollowingBySender(ctx context.Context, query queries.GetFollowingBySender) ([]*domain.MiddlemanFollow, error)
		GetMostFollowed(ctx context.Context, query queries.GetMostFollowed) ([]*domain.ItemFollowCount, error)
		GetMostFollowedByCategory(ctx context.Context, query queries.GetMostFollowedByCategory) ([]*domain.ItemFollowCount, error)
		GetApprovedFollowing(ctx context.Context, query queries.GetApprovedFollowing) ([]*domain.MiddlemanFollow, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.AddFollowHandler
		commands.ApproveFollowHandler
		commands.FlagFollowHandler
		commands.RejectFollowHandler
		commands.RemoveFollowHandler
	}

	appQueries struct {
		queries.GetFollowHandler
		queries.GetFollowingHandler
		queries.GetFollowingBySenderHandler
		queries.GetMostFollowedHandler
		queries.GetMostFollowedByCategoryHandler
		queries.GetApprovedFollowingHandler
	}
)

var _ App = (*Application)(nil)

func New(following domain.FollowRepository, middlemanFollowing domain.MiddlemanCacheRepository,
	publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		appCommands: appCommands{
			AddFollowHandler:     commands.NewAddFollowHandler(following, publisher),
			ApproveFollowHandler: commands.NewApproveFollowHandler(following, publisher),
			FlagFollowHandler:    commands.NewFlagFollowHandler(following, publisher),
			RejectFollowHandler:  commands.NewRejectFollowHandler(following, publisher),
			RemoveFollowHandler:  commands.NewRemoveFollowHandler(following, publisher),
		},
		appQueries: appQueries{
			GetFollowHandler:                 queries.NewGetFollowHandler(middlemanFollowing),
			GetFollowingHandler:              queries.NewGetFollowingHandler(middlemanFollowing),
			GetFollowingBySenderHandler:      queries.NewGetFollowingBySenderHandler(middlemanFollowing),
			GetMostFollowedHandler:           queries.NewGetMostFollowedHandler(middlemanFollowing),
			GetMostFollowedByCategoryHandler: queries.NewGetMostFollowedByCategoryHandler(middlemanFollowing),
			GetApprovedFollowingHandler:      queries.NewGetApprovedFollowingHandler(middlemanFollowing),
		},
	}
}
