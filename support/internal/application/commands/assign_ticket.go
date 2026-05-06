package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type AssignTicket struct {
	ID               string
	AssigneeID       string
	AssigneeType     domain.AssigneeType
	AssignedBy       string
	AssignmentReason string
}

type AssignTicketHandler struct {
	tickets   domain.TicketRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAssignTicketHandler(
	tickets domain.TicketRepository,
	publisher ddd.EventPublisher[ddd.Event],
) AssignTicketHandler {
	return AssignTicketHandler{
		tickets:   tickets,
		publisher: publisher,
	}
}

func (h AssignTicketHandler) AssignTicket(ctx context.Context, cmd AssignTicket) error {
	ticket, err := h.tickets.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := ticket.AssignTicket(cmd.AssigneeID, cmd.AssigneeType, cmd.AssignedBy, cmd.AssignmentReason)
	if err != nil {
		return err
	}

	err = h.tickets.Save(ctx, ticket)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}