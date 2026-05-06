package queries

import (
	"context"
	"middleman/tickets/internal/domain"
)

type GetUserTickets struct {
	UserID       string
	IncludePast  bool
	IncludeFuture bool
}

type GetUserTicketsHandler struct {
	tickets domain.TicketRepository
}

func NewGetUserTicketsHandler(tickets domain.TicketRepository) GetUserTicketsHandler {
	return GetUserTicketsHandler{tickets: tickets}
}

func (h GetUserTicketsHandler) GetUserTickets(ctx context.Context, query GetUserTickets) ([]*domain.Ticket, error) {
	var allTickets []*domain.Ticket
	
	if query.IncludeFuture {
		upcoming, err := h.tickets.GetUserUpcomingTickets(ctx, query.UserID)
		if err != nil {
			return nil, err
		}
		allTickets = append(allTickets, upcoming...)
	}
	
	if query.IncludePast {
		past, err := h.tickets.GetUserPastTickets(ctx, query.UserID, 50) // Limit to 50 past tickets
		if err != nil {
			return nil, err
		}
		allTickets = append(allTickets, past...)
	}
	
	return allTickets, nil
}