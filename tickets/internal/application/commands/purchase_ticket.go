package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/tickets/internal/domain"
)

type PurchaseTicket struct {
	TicketID     string
	MatchID      string
	SectorID     string
	RowID        string
	SeatNumber   int
	TicketType   domain.TicketType
	UserID       string
	UserName     string
	UserEmail    string
	UserPhone    string
	PaymentID    string
	OrderID      string
	Price        int64
	Transferable bool
}

type PurchaseTicketHandler struct {
	matches ddd.AggregateStore[*domain.Match]
	tickets ddd.AggregateStore[*domain.Ticket]
}

func NewPurchaseTicketHandler(
	matches ddd.AggregateStore[*domain.Match],
	tickets ddd.AggregateStore[*domain.Ticket],
) PurchaseTicketHandler {
	return PurchaseTicketHandler{
		matches: matches,
		tickets: tickets,
	}
}

func (h PurchaseTicketHandler) PurchaseTicket(ctx context.Context, cmd PurchaseTicket) error {
	// Load the match to get match details
	match, err := h.matches.Load(ctx, cmd.MatchID)
	if err != nil {
		return err
	}

	// Get seat information
	sector := match.Sectors[cmd.SectorID]
	if sector == nil {
		return domain.ErrSectorNotFound
	}

	// Create the ticket
	ticket := domain.NewTicket(cmd.TicketID)
	
	event, err := ticket.InitTicket(
		cmd.MatchID,
		match.MatchDate,
		match.HomeTeam.Name,
		match.AwayTeam.Name,
		match.Competition.Name,
		match.Stadium.Name,
		cmd.SectorID,
		sector.Name,
		cmd.RowID,
		sector.Rows[cmd.RowID].Number,
		cmd.SeatNumber,
		sector.EntranceGates[0], // Use first entrance gate
		cmd.TicketType,
		sector.Category,
		cmd.Price,
		cmd.UserID,
		cmd.UserName,
		cmd.UserEmail,
		cmd.UserPhone,
		cmd.UserID, // Purchaser is same as owner initially
		cmd.PaymentID,
		cmd.OrderID,
		cmd.Transferable,
	)
	if err != nil {
		return err
	}
	ticket.AddEvent(event)

	// Save the ticket
	if err := h.tickets.Save(ctx, ticket); err != nil {
		return err
	}

	// Update the match to record the ticket purchase
	match.AddEvent(domain.TicketPurchasedEvent, &domain.TicketPurchased{
		MatchID:     cmd.MatchID,
		TicketID:    cmd.TicketID,
		SectorID:    cmd.SectorID,
		RowID:       cmd.RowID,
		SeatNumber:  cmd.SeatNumber,
		UserID:      cmd.UserID,
		Price:       cmd.Price,
		Category:    sector.Category,
		PurchasedAt: ticket.CreatedAt,
	})

	return h.matches.Save(ctx, match)
}