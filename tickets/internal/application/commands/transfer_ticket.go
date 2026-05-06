package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/tickets/internal/domain"
)

type TransferTicket struct {
	TicketID    string
	ToUserID    string
	ToUserName  string
	ToUserEmail string
	ToUserPhone string
	Reason      string
}

type TransferTicketHandler struct {
	tickets ddd.AggregateStore[*domain.Ticket]
}

func NewTransferTicketHandler(tickets ddd.AggregateStore[*domain.Ticket]) TransferTicketHandler {
	return TransferTicketHandler{
		tickets: tickets,
	}
}

func (h TransferTicketHandler) TransferTicket(ctx context.Context, cmd TransferTicket) error {
	ticket, err := h.tickets.Load(ctx, cmd.TicketID)
	if err != nil {
		return err
	}

	event, err := ticket.TransferTicket(
		cmd.ToUserID,
		cmd.ToUserName,
		cmd.ToUserEmail,
		cmd.ToUserPhone,
		cmd.Reason,
	)
	if err != nil {
		return err
	}
	ticket.AddEvent(event)

	return h.tickets.Save(ctx, ticket)
}