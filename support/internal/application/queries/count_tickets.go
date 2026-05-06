package queries

import (
	"context"
	"middleman/support/internal/domain"
)

type CountTickets struct {
	Query   string
	Filters map[string]interface{}
}

type CountTicketsHandler struct {
	ticketCatalog domain.TicketCatalogRepository
}

func NewCountTicketsHandler(ticketCatalog domain.TicketCatalogRepository) CountTicketsHandler {
	return CountTicketsHandler{
		ticketCatalog: ticketCatalog,
	}
}

func (h CountTicketsHandler) CountTickets(ctx context.Context, query CountTickets) (int, error) {
	// Use search with limit 0 to get total count
	tickets, err := h.ticketCatalog.Search(ctx, query.Query, query.Filters, 0, 0)
	if err != nil {
		return 0, err
	}
	// For now, return the length of results
	// In a real implementation, this should be a dedicated count query
	return len(tickets), nil
}