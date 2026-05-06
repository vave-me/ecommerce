package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type ReopenTicketHandler struct {
	tickets   domain.TicketRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewReopenTicketHandler(
	tickets domain.TicketRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ReopenTicketHandler {
	return ReopenTicketHandler{
		tickets:   tickets,
		publisher: publisher,
	}
}

func (h ReopenTicketHandler) ReopenTicket(ctx context.Context, cmd ReopenTicket) error {
	ticket, err := h.tickets.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := ticket.ReopenTicket(cmd.ReopenedBy, cmd.ReopenReason)
	if err != nil {
		return err
	}
	h.publisher.Publish(ctx, event)

	return h.tickets.Save(ctx, ticket)
}