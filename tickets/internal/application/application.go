package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/tickets/internal/application/commands"
	"middleman/tickets/internal/application/queries"
	"middleman/tickets/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		// Match commands
		CreateMatch(ctx context.Context, cmd commands.CreateMatch) error
		AddSector(ctx context.Context, cmd commands.AddSector) error
		InitializeSectorSeats(ctx context.Context, cmd commands.InitializeSectorSeats) error
		PublishMatch(ctx context.Context, cmd commands.PublishMatch) error
		
		// Ticket commands
		PurchaseTicket(ctx context.Context, cmd commands.PurchaseTicket) error
		TransferTicket(ctx context.Context, cmd commands.TransferTicket) error
		ValidateTicket(ctx context.Context, cmd commands.ValidateTicket) error
	}

	Queries interface {
		// Match queries
		GetMatch(ctx context.Context, query queries.GetMatch) (*domain.Match, error)
		GetMatchCatalog(ctx context.Context, query queries.GetMatchCatalog) (*domain.MatchCatalog, error)
		GetSeatMap(ctx context.Context, query queries.GetSeatMap) (*domain.SectorSeatMap, error)
		SearchMatches(ctx context.Context, query queries.SearchMatches) ([]*domain.Match, error)
		
		// Ticket queries
		GetUserTickets(ctx context.Context, query queries.GetUserTickets) ([]*domain.Ticket, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		// Match command handlers
		commands.CreateMatchHandler
		commands.AddSectorHandler
		commands.InitializeSectorSeatsHandler
		commands.PublishMatchHandler
		
		// Ticket command handlers
		commands.PurchaseTicketHandler
		commands.TransferTicketHandler
		commands.ValidateTicketHandler
	}

	appQueries struct {
		// Match query handlers
		queries.GetMatchHandler
		queries.GetMatchCatalogHandler
		queries.GetSeatMapHandler
		queries.SearchMatchesHandler
		
		// Ticket query handlers
		queries.GetUserTicketsHandler
	}
)

var _ App = (*Application)(nil)

func New(
	matches ddd.AggregateStore[*domain.Match],
	tickets ddd.AggregateStore[*domain.Ticket],
	matchRepo domain.MatchRepository,
	ticketRepo domain.TicketRepository,
	matchCatalogRepo domain.MatchCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		appCommands: appCommands{
			// Match commands
			CreateMatchHandler:           commands.NewCreateMatchHandler(matches),
			AddSectorHandler:             commands.NewAddSectorHandler(matches),
			InitializeSectorSeatsHandler: commands.NewInitializeSectorSeatsHandler(matches),
			PublishMatchHandler:          commands.NewPublishMatchHandler(matches),
			
			// Ticket commands
			PurchaseTicketHandler: commands.NewPurchaseTicketHandler(matches, tickets),
			TransferTicketHandler: commands.NewTransferTicketHandler(tickets),
			ValidateTicketHandler: commands.NewValidateTicketHandler(tickets),
		},
		appQueries: appQueries{
			// Match queries
			GetMatchHandler:        queries.NewGetMatchHandler(matchRepo),
			GetMatchCatalogHandler: queries.NewGetMatchCatalogHandler(matchCatalogRepo),
			GetSeatMapHandler:      queries.NewGetSeatMapHandler(matchCatalogRepo),
			SearchMatchesHandler:   queries.NewSearchMatchesHandler(matchRepo),
			
			// Ticket queries
			GetUserTicketsHandler: queries.NewGetUserTicketsHandler(ticketRepo),
		},
	}
}