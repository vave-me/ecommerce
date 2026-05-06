package handlers

import (
	"context"
	"time"

	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/tickets/internal/application"
	"middleman/tickets/internal/application/commands"
	"middleman/tickets/ticketspb"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// commandHandlers routes incoming NATS/JetStream commands to the application.
// Only ReserveTicket is handled for now.

type commandHandlers struct {
	app application.App
}

var _ ddd.CommandHandler[ddd.Command] = (*commandHandlers)(nil)

// NewCommandHandlers wires registry-based serialization, a reply publisher and middlewares.
func NewCommandHandlers(reg registry.Registry, app application.App, replyPublisher am.ReplyPublisher, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewCommandHandler(reg, replyPublisher, commandHandlers{app: app}, mws...)
}

// RegisterCommandHandlers subscribes the given subscriber to the tickets command channel.
func RegisterCommandHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(ticketspb.CommandChannel, handlers, am.MessageFilter{
		ticketspb.ReserveTicketCommand,
		ticketspb.ReleaseTicketCommand,
		ticketspb.ReserveTicketsCommand,
		ticketspb.ReleaseTicketsCommand,
	}, am.GroupName("tickets-commands"))
	return err
}

// HandleCommand dispatches based on CommandName.
func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (reply ddd.Reply, err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent("error", trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		}
		span.AddEvent("done", trace.WithAttributes(attribute.Int64("ms", time.Since(started).Milliseconds())))
	}(time.Now())

	switch cmd.CommandName() {
	case ticketspb.ReserveTicketCommand:
		return h.doReserveTicket(ctx, cmd)
	case ticketspb.ReleaseTicketCommand:
		return h.doReleaseTicket(ctx, cmd)
	case ticketspb.ReserveTicketsCommand:
		return h.doReserveTickets(ctx, cmd)
	case ticketspb.ReleaseTicketsCommand:
		return h.doReleaseTickets(ctx, cmd)
	}
	return nil, nil
}

func (h commandHandlers) doReserveTicket(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*ticketspb.ReserveTicket)

	if err := h.app.ReserveTicket(ctx, commands.ReserveTicket{
		TicketID: payload.GetTicketId(),
		Quantity: payload.GetQuantity(),
	}); err != nil {
		return nil, err
	}

	// publish success reply so saga can continue
	return ddd.NewReply(am.SuccessReply, nil), nil
}

func (h commandHandlers) doReleaseTicket(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*ticketspb.ReleaseTicket)

	if err := h.app.ReleaseTicket(ctx, commands.ReleaseTicket{
		TicketID: payload.GetTicketId(),
		Quantity: payload.GetQuantity(),
	}); err != nil {
		return nil, err
	}

	return ddd.NewReply(am.SuccessReply, nil), nil
}

func (h commandHandlers) doReserveTickets(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*ticketspb.ReserveTickets)

	// Process each item in the batch
	for _, item := range payload.GetItems() {
		if err := h.app.ReserveTicket(ctx, commands.ReserveTicket{
			TicketID: item.GetTicketId(),
			Quantity: item.GetQuantity(),
		}); err != nil {
			// If any reservation fails, return error
			// The saga will handle compensation
			return nil, err
		}
	}

	// All reservations successful
	return ddd.NewReply(am.SuccessReply, nil), nil
}

func (h commandHandlers) doReleaseTickets(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*ticketspb.ReleaseTickets)

	// Process each item in the batch
	for _, item := range payload.GetItems() {
		if err := h.app.ReleaseTicket(ctx, commands.ReleaseTicket{
			TicketID: item.GetTicketId(),
			Quantity: item.GetQuantity(),
		}); err != nil {
			// Continue with other releases even if one fails
			// This is a compensation action, so we try to release as much as possible
			continue
		}
	}

	return ddd.NewReply(am.SuccessReply, nil), nil
}
