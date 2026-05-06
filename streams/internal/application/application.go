package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/streams/internal/application/commands"
	"middleman/streams/internal/application/queries"
	"middleman/streams/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		// Stream commands
		CreateStream(ctx context.Context, cmd commands.CreateStream) error
		PublishStream(ctx context.Context, cmd commands.PublishStream) error
		SetStreamPricing(ctx context.Context, cmd commands.SetStreamPricing) error
		GrantUserAccess(ctx context.Context, cmd commands.GrantUserAccess) error
		UpdateWatchProgress(ctx context.Context, cmd commands.UpdateWatchProgress) error
		
		// Series commands
		CreateSeries(ctx context.Context, cmd commands.CreateSeries) error
		AddSeason(ctx context.Context, cmd commands.AddSeason) error
		AddEpisode(ctx context.Context, cmd commands.AddEpisode) error
	}

	Queries interface {
		// Stream queries
		GetStream(ctx context.Context, query queries.GetStream) (*domain.Stream, error)
		GetUserStreams(ctx context.Context, query queries.GetUserStreams) ([]*domain.Stream, error)
		SearchStreams(ctx context.Context, query queries.SearchStreams) ([]*domain.Stream, error)
		GetContinueWatching(ctx context.Context, query queries.GetContinueWatching) ([]*domain.Stream, error)
		
		// Catalog queries
		GetCatalog(ctx context.Context, query queries.GetCatalog) (*domain.StreamCatalog, error)
		
		// Series queries
		GetSeries(ctx context.Context, query queries.GetSeries) (*domain.Series, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		// Stream command handlers
		commands.CreateStreamHandler
		commands.PublishStreamHandler
		commands.SetStreamPricingHandler
		commands.GrantUserAccessHandler
		commands.UpdateWatchProgressHandler
		
		// Series command handlers
		commands.CreateSeriesHandler
		commands.AddSeasonHandler
		commands.AddEpisodeHandler
	}

	appQueries struct {
		// Stream query handlers
		queries.GetStreamHandler
		queries.GetUserStreamsHandler
		queries.SearchStreamsHandler
		queries.GetContinueWatchingHandler
		
		// Catalog query handlers
		queries.GetCatalogHandler
		
		// Series query handlers
		queries.GetSeriesHandler
	}
)

var _ App = (*Application)(nil)

func New(
	streams ddd.AggregateStore[*domain.Stream],
	series ddd.AggregateStore[*domain.Series],
	streamRepo domain.StreamRepository,
	seriesRepo domain.SeriesRepository,
	catalogRepo domain.StreamCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		appCommands: appCommands{
			// Stream commands
			CreateStreamHandler:        commands.NewCreateStreamHandler(streams),
			PublishStreamHandler:       commands.NewPublishStreamHandler(streams),
			SetStreamPricingHandler:    commands.NewSetStreamPricingHandler(streams),
			GrantUserAccessHandler:     commands.NewGrantUserAccessHandler(streams),
			UpdateWatchProgressHandler: commands.NewUpdateWatchProgressHandler(streams),
			
			// Series commands
			CreateSeriesHandler: commands.NewCreateSeriesHandler(series),
			AddSeasonHandler:    commands.NewAddSeasonHandler(series),
			AddEpisodeHandler:   commands.NewAddEpisodeHandler(series),
		},
		appQueries: appQueries{
			// Stream queries
			GetStreamHandler:           queries.NewGetStreamHandler(streamRepo),
			GetUserStreamsHandler:      queries.NewGetUserStreamsHandler(streamRepo),
			SearchStreamsHandler:       queries.NewSearchStreamsHandler(streamRepo),
			GetContinueWatchingHandler: queries.NewGetContinueWatchingHandler(streamRepo),
			
			// Catalog queries
			GetCatalogHandler: queries.NewGetCatalogHandler(catalogRepo),
			
			// Series queries
			GetSeriesHandler: queries.NewGetSeriesHandler(seriesRepo),
		},
	}
}