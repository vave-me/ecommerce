package queries

import (
	"context"
	"middleman/support/internal/domain"
)

type GetTicket struct {
	ID string
}

type GetTicketHandler struct {
	catalog domain.TicketCatalogRepository
}

func NewGetTicketHandler(catalog domain.TicketCatalogRepository) GetTicketHandler {
	return GetTicketHandler{
		catalog: catalog,
	}
}

func (h GetTicketHandler) GetTicket(ctx context.Context, query GetTicket) (*domain.TicketCatalog, error) {
	return h.catalog.Find(ctx, query.ID)
}