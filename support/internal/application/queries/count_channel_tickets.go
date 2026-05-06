package queries

import (
	"context"
	"middleman/support/internal/domain"
)

type CountChannelTickets struct {
	ChannelID string
	Status    *string
}

type CountChannelTicketsHandler struct {
	ticketCatalog domain.TicketCatalogRepository
}

func NewCountChannelTicketsHandler(ticketCatalog domain.TicketCatalogRepository) CountChannelTicketsHandler {
	return CountChannelTicketsHandler{
		ticketCatalog: ticketCatalog,
	}
}

func (h CountChannelTicketsHandler) CountChannelTickets(ctx context.Context, query CountChannelTickets) (int, error) {
	return h.ticketCatalog.Count(ctx, query.ChannelID, query.Status)
}