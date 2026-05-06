package application

import (
	"context"
	"middleman/activity/internal/application/commands"
	"middleman/activity/internal/application/queries"
	"middleman/activity/internal/domain"
	"middleman/internal/ddd"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		CreateActivity(ctx context.Context, cmd commands.CreateActivity) error
		AddInteraction(ctx context.Context, cmd commands.AddInteraction) error
		RemoveInteraction(ctx context.Context, cmd commands.RemoveInteraction) error
		UpdateInteraction(ctx context.Context, cmd commands.UpdateInteraction) error
	}
	Queries interface {
		GetInteraction(ctx context.Context, query queries.GetInteraction) (*domain.MiddlemanInteraction, error)
		GetInteractions(ctx context.Context, query queries.GetInteractions) ([]*domain.MiddlemanInteraction, error)
		GetMostLiked(ctx context.Context, query queries.GetMostLiked) ([]*domain.MostReactionResult, error)
		GetMostDisliked(ctx context.Context, query queries.GetMostDisliked) ([]*domain.MostReactionResult, error)
		GetActivities(ctx context.Context, query queries.GetActivities) ([]*domain.MiddlemanActivity, error)
		GetActivity(ctx context.Context, query queries.GetActivity) (*domain.MiddlemanActivity, error)
	}

	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.CreateActivityHandler
		commands.AddInteractionHandler
		commands.RemoveInteractionHandler
		commands.UpdateInteractionHandler
	}
	appQueries struct {
		queries.GetInteractionHandler
		queries.GetInteractionsHandler
		queries.GetActivityHandler
		queries.GetActivitiesHandler
		queries.GetMostLikedHandler
		queries.GetMostDislikedHandler
	}
)

var _ App = (*Application)(nil)

func New(activity domain.ActivityRepository, interactions domain.InteractionRepository,
	middleman domain.MiddlemanRepository, itemInteractions domain.MiddlemanInteractionRepository,
	publisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateActivityHandler:    commands.NewCreateActivityHandler(activity, publisher),
			AddInteractionHandler:    commands.NewAddInteractionHandler(interactions, publisher),
			RemoveInteractionHandler: commands.NewRemoveInteractionHandler(interactions, publisher),
			UpdateInteractionHandler: commands.NewUpdateInteractionHandler(interactions, publisher),
		},
		appQueries: appQueries{
			GetActivitiesHandler:   queries.NewGetActivitiesHandler(middleman),
			GetActivityHandler:     queries.NewGetActivityHandler(middleman),
			GetInteractionHandler:  queries.NewGetInteractionHandler(itemInteractions),
			GetInteractionsHandler: queries.NewGetInteractionsHandler(itemInteractions),
			GetMostLikedHandler:    queries.NewGetMostLikedHandler(itemInteractions),
			GetMostDislikedHandler: queries.NewGetMostDislikedHandler(itemInteractions),
		},
	}
}
