package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type EscalateTicketHandler struct {
	tickets   domain.TicketRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewEscalateTicketHandler(
	tickets domain.TicketRepository,
	publisher ddd.EventPublisher[ddd.Event],
) EscalateTicketHandler {
	return EscalateTicketHandler{
		tickets:   tickets,
		publisher: publisher,
	}
}

func (h EscalateTicketHandler) EscalateTicket(ctx context.Context, cmd EscalateTicket) error {
	ticket, err := h.tickets.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := ticket.EscalateTicket(cmd.EscalationTier, cmd.EscalatedBy, cmd.EscalationReason, cmd.EscalationNotes)
	if err != nil {
		return err
	}
	h.publisher.Publish(ctx, event)

	return h.tickets.Save(ctx, ticket)
}