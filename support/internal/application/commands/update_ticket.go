package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type UpdateTicketHandler struct {
	tickets   domain.TicketRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewUpdateTicketHandler(
	tickets domain.TicketRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UpdateTicketHandler {
	return UpdateTicketHandler{
		tickets:   tickets,
		publisher: publisher,
	}
}

func (h UpdateTicketHandler) UpdateTicket(ctx context.Context, cmd UpdateTicket) error {
	ticket, err := h.tickets.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := ticket.UpdateTicket(cmd.Title, cmd.Description, cmd.Category, cmd.Tags, cmd.Metadata, cmd.UpdatedBy)
	if err != nil {
		return err
	}
	h.publisher.Publish(ctx, event)

	return h.tickets.Save(ctx, ticket)
}