package queries

import (
	"context"
	"middleman/support/internal/domain"
)

type SearchTickets struct {
	Query      string
	Filters    map[string]interface{}
	Limit      int
	Offset     int
	SortBy     string
	Descending bool
}

type SearchTicketsHandler struct {
	ticketCatalog domain.TicketCatalogRepository
}

func NewSearchTicketsHandler(ticketCatalog domain.TicketCatalogRepository) SearchTicketsHandler {
	return SearchTicketsHandler{
		ticketCatalog: ticketCatalog,
	}
}

func (h SearchTicketsHandler) SearchTickets(ctx context.Context, query SearchTickets) ([]*domain.TicketCatalog, error) {
	return h.ticketCatalog.Search(ctx, query.Query, query.Filters, query.Limit, query.Offset)
}