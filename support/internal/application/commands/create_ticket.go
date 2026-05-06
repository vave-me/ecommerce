package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type CreateTicket struct {
	ID          string
	ChannelID   string
	Title       string
	Description string
	Category    domain.TicketCategory
	Priority    domain.TicketPriority
	Tags        []string
	Metadata    map[string]string
	CreatedBy   string
	Attachments []domain.Attachment
}

type CreateTicketHandler struct {
	tickets       domain.TicketRepository
	channels      domain.SupportChannelRepository
	channelCatalog domain.SupportChannelCatalogRepository
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewCreateTicketHandler(
	tickets domain.TicketRepository,
	channels domain.SupportChannelRepository,
	channelCatalog domain.SupportChannelCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) CreateTicketHandler {
	return CreateTicketHandler{
		tickets:       tickets,
		channels:      channels,
		channelCatalog: channelCatalog,
		publisher:     publisher,
	}
}

func (h CreateTicketHandler) CreateTicket(ctx context.Context, cmd CreateTicket) error {
	// Verify channel exists and is active
	channelCatalog, err := h.channelCatalog.Find(ctx, cmd.ChannelID)
	if err != nil {
		return err
	}
	if !channelCatalog.Active {
		return domain.ErrChannelNotActive
	}

	// Create the ticket
	ticket, err := h.tickets.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := ticket.InitTicket(
		cmd.ChannelID,
		cmd.Title,
		cmd.Description,
		cmd.Category,
		cmd.Priority,
		cmd.Tags,
		cmd.Metadata,
		cmd.CreatedBy,
		cmd.Attachments,
	)
	if err != nil {
		return err
	}

	err = h.tickets.Save(ctx, ticket)
	if err != nil {
		return err
	}

	// Update channel ticket count
	channel, err := h.channels.Load(ctx, cmd.ChannelID)
	if err != nil {
		return err
	}

	channel.IncrementTicketCount()
	err = h.channels.Save(ctx, channel)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}