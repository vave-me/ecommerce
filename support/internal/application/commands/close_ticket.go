package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type CloseTicketHandler struct {
	tickets   domain.TicketRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCloseTicketHandler(
	tickets domain.TicketRepository,
	publisher ddd.EventPublisher[ddd.Event],
) CloseTicketHandler {
	return CloseTicketHandler{
		tickets:   tickets,
		publisher: publisher,
	}
}

func (h CloseTicketHandler) CloseTicket(ctx context.Context, cmd CloseTicket) error {
	ticket, err := h.tickets.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := ticket.CloseTicket(cmd.ClosedBy, cmd.ClosureNotes, cmd.SatisfactionRating)
	if err != nil {
		return err
	}
	h.publisher.Publish(ctx, event)

	return h.tickets.Save(ctx, ticket)
}