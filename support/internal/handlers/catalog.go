package handlers

import (
	"context"

	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/support/internal/domain"
)

// catalogHandlers listens to domain events and updates the Catalog DB
type catalogHandlers[T ddd.Event] struct {
	channelCatalog domain.SupportChannelCatalogRepository
	ticketCatalog  domain.TicketCatalogRepository
}

var _ ddd.EventHandler[ddd.Event] = (*catalogHandlers[ddd.Event])(nil)

func NewCatalogHandlers(
	channelCatalog domain.SupportChannelCatalogRepository,
	ticketCatalog domain.TicketCatalogRepository,
) ddd.EventHandler[ddd.Event] {
	return catalogHandlers[ddd.Event]{
		channelCatalog: channelCatalog,
		ticketCatalog:  ticketCatalog,
	}
}

func RegisterCatalogHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		// Support Channel events
		domain.SupportChannelCreatedEvent,
		domain.SupportChannelSettingsUpdatedEvent,
		domain.SupportChannelClosedEvent,
		domain.SupportChannelReactivatedEvent,

		// Ticket events
		domain.TicketCreatedEvent,
		domain.TicketUpdatedEvent,
		domain.TicketAssignedEvent,
		domain.TicketPriorityUpdatedEvent,
		domain.TicketEscalatedEvent,
		domain.TicketResolvedEvent,
		domain.TicketReopenedEvent,
		domain.TicketClosedEvent,
	)
}

func RegisterCatalogHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		catalogHandlers := di.Get(ctx, "catalogHandlers").(ddd.EventHandler[ddd.Event])
		return catalogHandlers.HandleEvent(ctx, event)
	})
	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])

	RegisterCatalogHandlers(subscriber, handlers)
}

func (h catalogHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	// Support Channel events
	case domain.SupportChannelCreatedEvent:
		return h.onSupportChannelCreated(ctx, event)
	case domain.SupportChannelSettingsUpdatedEvent:
		return h.onSupportChannelSettingsUpdated(ctx, event)
	case domain.SupportChannelClosedEvent:
		return h.onSupportChannelClosed(ctx, event)
	case domain.SupportChannelReactivatedEvent:
		return h.onSupportChannelReactivated(ctx, event)

	// Ticket events
	case domain.TicketCreatedEvent:
		return h.onTicketCreated(ctx, event)
	case domain.TicketUpdatedEvent:
		return h.onTicketUpdated(ctx, event)
	case domain.TicketAssignedEvent:
		return h.onTicketAssigned(ctx, event)
	case domain.TicketPriorityUpdatedEvent:
		return h.onTicketPriorityUpdated(ctx, event)
	case domain.TicketEscalatedEvent:
		return h.onTicketEscalated(ctx, event)
	case domain.TicketResolvedEvent:
		return h.onTicketResolved(ctx, event)
	case domain.TicketReopenedEvent:
		return h.onTicketReopened(ctx, event)
	case domain.TicketClosedEvent:
		return h.onTicketClosed(ctx, event)
	}
	return nil
}

// Support Channel event handlers
func (h catalogHandlers[T]) onSupportChannelCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.SupportChannelCreated)
	channelID := event.Metadata().Get(ddd.AggregateIDKey).(string)

	channel := &domain.SupportChannelCatalog{
		ID:           channelID,
		UserID:       payload.UserID,
		BusinessID:   payload.BusinessID,
		ChannelType:  string(payload.ChannelType),
		Active:       true,
		OpenTickets:  0,
		TotalTickets: 0,
		CreatedAt:    event.OccurredAt(),
		UpdatedAt:    event.OccurredAt(),
	}

	return h.channelCatalog.Add(ctx, channel)
}

func (h catalogHandlers[T]) onSupportChannelSettingsUpdated(ctx context.Context, event ddd.Event) error {
	channelID := event.Metadata().Get(ddd.AggregateIDKey).(string)

	channel, err := h.channelCatalog.Find(ctx, channelID)
	if err != nil {
		return err
	}

	channel.UpdatedAt = event.OccurredAt()
	return h.channelCatalog.Update(ctx, channel)
}

func (h catalogHandlers[T]) onSupportChannelClosed(ctx context.Context, event ddd.Event) error {
	channelID := event.Metadata().Get(ddd.AggregateIDKey).(string)

	channel, err := h.channelCatalog.Find(ctx, channelID)
	if err != nil {
		return err
	}

	channel.Active = false
	channel.UpdatedAt = event.OccurredAt()
	return h.channelCatalog.Update(ctx, channel)
}

func (h catalogHandlers[T]) onSupportChannelReactivated(ctx context.Context, event ddd.Event) error {
	channelID := event.Metadata().Get(ddd.AggregateIDKey).(string)

	channel, err := h.channelCatalog.Find(ctx, channelID)
	if err != nil {
		return err
	}

	channel.Active = true
	channel.UpdatedAt = event.OccurredAt()
	return h.channelCatalog.Update(ctx, channel)
}

// Ticket event handlers
func (h catalogHandlers[T]) onTicketCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketCreated)
	ticketID := event.Metadata().Get(ddd.AggregateIDKey).(string)

	ticket := &domain.TicketCatalog{
		ID:        ticketID,
		ChannelID: payload.ChannelID,
		Title:     payload.Title,
		Status:    string(domain.TicketStatusSubmitted),
		Priority:  string(payload.Priority),
		Category:  string(payload.Category),
		CreatedBy: payload.CreatedBy,
		CreatedAt: event.OccurredAt(),
		UpdatedAt: event.OccurredAt(),
	}

	if err := h.ticketCatalog.Add(ctx, ticket); err != nil {
		return err
	}

	// Update channel ticket count
	channel, err := h.channelCatalog.Find(ctx, payload.ChannelID)
	if err != nil {
		return err
	}

	channel.OpenTickets++
	channel.TotalTickets++
	return h.channelCatalog.Update(ctx, channel)
}

func (h catalogHandlers[T]) onTicketUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketUpdated)
	ticketID := event.Metadata().Get(ddd.AggregateIDKey).(string)

	ticket, err := h.ticketCatalog.Find(ctx, ticketID)
	if err != nil {
		return err
	}

	if payload.Title != nil {
		ticket.Title = *payload.Title
	}
	if payload.Category != nil {
		ticket.Category = string(*payload.Category)
	}

	ticket.UpdatedAt = event.OccurredAt()
	return h.ticketCatalog.Update(ctx, ticket)
}

func (h catalogHandlers[T]) onTicketAssigned(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketAssigned)
	ticketID := event.Metadata().Get(ddd.AggregateIDKey).(string)

	ticket, err := h.ticketCatalog.Find(ctx, ticketID)
	if err != nil {
		return err
	}

	ticket.AssigneeID = &payload.AssigneeID
	assigneeType := string(payload.AssigneeType)
	ticket.AssigneeType = &assigneeType
	ticket.Status = string(domain.TicketStatusAssigned)
	ticket.UpdatedAt = event.OccurredAt()

	return h.ticketCatalog.Update(ctx, ticket)
}

func (h catalogHandlers[T]) onTicketPriorityUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TicketPriorityUpdated)
	ticketID := event.Metadata().Get(ddd.AggregateIDKey).(string)

	ticket, err := h.ticketCatalog.Find(ctx, ticketID)
	if err != nil {
		return err
	}

	ticket.Priority = string(payload.NewPriority)
	ticket.UpdatedAt = event.OccurredAt()

	return h.ticketCatalog.Update(ctx, ticket)
}

func (h catalogHandlers[T]) onTicketEscalated(ctx context.Context, event ddd.Event) error {
	ticketID := event.Metadata().Get(ddd.AggregateIDKey).(string)

	ticket, err := h.ticketCatalog.Find(ctx, ticketID)
	if err != nil {
		return err
	}

	ticket.Status = string(domain.TicketStatusEscalated)
	ticket.UpdatedAt = event.OccurredAt()

	return h.ticketCatalog.Update(ctx, ticket)
}

func (h catalogHandlers[T]) onTicketResolved(ctx context.Context, event ddd.Event) error {
	ticketID := event.Metadata().Get(ddd.AggregateIDKey).(string)

	ticket, err := h.ticketCatalog.Find(ctx, ticketID)
	if err != nil {
		return err
	}

	oldStatus := ticket.Status
	ticket.Status = string(domain.TicketStatusResolved)
	ticket.UpdatedAt = event.OccurredAt()

	if err := h.ticketCatalog.Update(ctx, ticket); err != nil {
		return err
	}

	// Update channel open tickets count if status changed from open to resolved
	if oldStatus != string(domain.TicketStatusResolved) && oldStatus != string(domain.TicketStatusClosed) {
		channel, err := h.channelCatalog.Find(ctx, ticket.ChannelID)
		if err != nil {
			return err
		}

		if channel.OpenTickets > 0 {
			channel.OpenTickets--
		}
		return h.channelCatalog.Update(ctx, channel)
	}

	return nil
}

func (h catalogHandlers[T]) onTicketReopened(ctx context.Context, event ddd.Event) error {
	ticketID := event.Metadata().Get(ddd.AggregateIDKey).(string)

	ticket, err := h.ticketCatalog.Find(ctx, ticketID)
	if err != nil {
		return err
	}

	oldStatus := ticket.Status
	ticket.Status = string(domain.TicketStatusReopened)
	ticket.UpdatedAt = event.OccurredAt()

	if err := h.ticketCatalog.Update(ctx, ticket); err != nil {
		return err
	}

	// Update channel open tickets count if status changed from closed/resolved to reopened
	if oldStatus == string(domain.TicketStatusResolved) || oldStatus == string(domain.TicketStatusClosed) {
		channel, err := h.channelCatalog.Find(ctx, ticket.ChannelID)
		if err != nil {
			return err
		}

		channel.OpenTickets++
		return h.channelCatalog.Update(ctx, channel)
	}

	return nil
}

func (h catalogHandlers[T]) onTicketClosed(ctx context.Context, event ddd.Event) error {
	ticketID := event.Metadata().Get(ddd.AggregateIDKey).(string)

	ticket, err := h.ticketCatalog.Find(ctx, ticketID)
	if err != nil {
		return err
	}

	oldStatus := ticket.Status
	ticket.Status = string(domain.TicketStatusClosed)
	ticket.UpdatedAt = event.OccurredAt()

	if err := h.ticketCatalog.Update(ctx, ticket); err != nil {
		return err
	}

	// Update channel open tickets count if status changed from open to closed
	if oldStatus != string(domain.TicketStatusResolved) && oldStatus != string(domain.TicketStatusClosed) {
		channel, err := h.channelCatalog.Find(ctx, ticket.ChannelID)
		if err != nil {
			return err
		}

		if channel.OpenTickets > 0 {
			channel.OpenTickets--
		}
		return h.channelCatalog.Update(ctx, channel)
	}

	return nil
}
