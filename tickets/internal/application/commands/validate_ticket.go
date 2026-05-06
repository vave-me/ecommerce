package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/tickets/internal/domain"
)

type ValidateTicket struct {
	TicketID string
	Gate     string
	QRCode   string
}

type ValidateTicketHandler struct {
	tickets ddd.AggregateStore[*domain.Ticket]
}

func NewValidateTicketHandler(tickets ddd.AggregateStore[*domain.Ticket]) ValidateTicketHandler {
	return ValidateTicketHandler{
		tickets: tickets,
	}
}

func (h ValidateTicketHandler) ValidateTicket(ctx context.Context, cmd ValidateTicket) error {
	ticket, err := h.tickets.Load(ctx, cmd.TicketID)
	if err != nil {
		return err
	}

	// Validate the ticket
	event, err := ticket.ValidateTicket(cmd.Gate, cmd.QRCode)
	if err != nil {
		return err
	}
	ticket.AddEvent(event)

	// If validation successful, mark as used
	if err == nil {
		useEvent, err := ticket.UseTicket(cmd.Gate)
		if err != nil {
			return err
		}
		ticket.AddEvent(useEvent)
	}

	return h.tickets.Save(ctx, ticket)
}