package supportpb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

// Channels
const (
	SupportChannelAggregateChannel = "middleman.support.events.SupportChannel"
	TicketAggregateChannel         = "middleman.support.events.Ticket"
	CommunicationAggregateChannel  = "middleman.support.events.Communication"
)

// Support Channel Event Names
const (
	SupportChannelCreatedEvent          = "supportapi.SupportChannelCreated"
	SupportChannelSettingsUpdatedEvent  = "supportapi.SupportChannelSettingsUpdated"
	SupportChannelClosedEvent           = "supportapi.SupportChannelClosed"
	SupportChannelReactivatedEvent      = "supportapi.SupportChannelReactivated"
)

// Ticket Event Names
const (
	TicketCreatedEvent         = "supportapi.TicketCreated"
	TicketUpdatedEvent         = "supportapi.TicketUpdated"
	TicketAssignedEvent        = "supportapi.TicketAssigned"
	TicketPriorityUpdatedEvent = "supportapi.TicketPriorityUpdated"
	TicketEscalatedEvent       = "supportapi.TicketEscalated"
	TicketResolvedEvent        = "supportapi.TicketResolved"
	TicketReopenedEvent        = "supportapi.TicketReopened"
	TicketClosedEvent          = "supportapi.TicketClosed"
	TicketsMergedEvent         = "supportapi.TicketsMerged"
	TicketsLinkedEvent         = "supportapi.TicketsLinked"
)

// Communication Event Names
const (
	TicketReplyAddedEvent  = "supportapi.TicketReplyAdded"
	InternalNoteAddedEvent = "supportapi.InternalNoteAdded"
)

// AI Integration Event Names
const (
	AISupportEnabledEvent      = "supportapi.AISupportEnabled"
	AIConfigurationUpdatedEvent = "supportapi.AIConfigurationUpdated"
	HandoffToHumanOccurredEvent = "supportapi.HandoffToHumanOccurred"
	HandoffToAIOccurredEvent    = "supportapi.HandoffToAIOccurred"
	AIResponseGeneratedEvent    = "supportapi.AIResponseGenerated"
)

// Knowledge Base Event Names
const (
	KnowledgeArticleCreatedEvent = "supportapi.KnowledgeArticleCreated"
	ArticleLinkedToTicketEvent   = "supportapi.ArticleLinkedToTicket"
	ArticleRatedEvent            = "supportapi.ArticleRated"
)

// Analytics Event Names
const (
	SLABreachedEvent                    = "supportapi.SLABreached"
	CustomerSatisfactionRecordedEvent   = "supportapi.CustomerSatisfactionRecorded"
	ResolutionTimeRecordedEvent         = "supportapi.ResolutionTimeRecorded"
)

// Registrations and Serde
func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Support Channel Events
	if err := serde.Register(&SupportChannelCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&SupportChannelSettingsUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&SupportChannelClosed{}); err != nil {
		return err
	}
	if err := serde.Register(&SupportChannelReactivated{}); err != nil {
		return err
	}

	// Ticket Events
	if err := serde.Register(&TicketCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&TicketUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&TicketAssigned{}); err != nil {
		return err
	}
	if err := serde.Register(&TicketPriorityUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&TicketEscalated{}); err != nil {
		return err
	}
	if err := serde.Register(&TicketResolved{}); err != nil {
		return err
	}
	if err := serde.Register(&TicketReopened{}); err != nil {
		return err
	}
	if err := serde.Register(&TicketClosed{}); err != nil {
		return err
	}
	if err := serde.Register(&TicketsMerged{}); err != nil {
		return err
	}
	if err := serde.Register(&TicketsLinked{}); err != nil {
		return err
	}

	// Communication Events
	if err := serde.Register(&TicketReplyAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&InternalNoteAdded{}); err != nil {
		return err
	}

	// AI Integration Events
	if err := serde.Register(&AISupportEnabled{}); err != nil {
		return err
	}
	if err := serde.Register(&AIConfigurationUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&HandoffToHumanOccurred{}); err != nil {
		return err
	}
	if err := serde.Register(&HandoffToAIOccurred{}); err != nil {
		return err
	}
	if err := serde.Register(&AIResponseGenerated{}); err != nil {
		return err
	}

	// Knowledge Base Events
	if err := serde.Register(&KnowledgeArticleCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&ArticleLinkedToTicket{}); err != nil {
		return err
	}
	if err := serde.Register(&ArticleRated{}); err != nil {
		return err
	}

	// Analytics Events
	if err := serde.Register(&SLABreached{}); err != nil {
		return err
	}
	if err := serde.Register(&CustomerSatisfactionRecorded{}); err != nil {
		return err
	}
	if err := serde.Register(&ResolutionTimeRecorded{}); err != nil {
		return err
	}

	return nil
}

// Key() implementations for all events
// Support Channel Events
func (*SupportChannelCreated) Key() string         { return SupportChannelCreatedEvent }
func (*SupportChannelSettingsUpdated) Key() string { return SupportChannelSettingsUpdatedEvent }
func (*SupportChannelClosed) Key() string          { return SupportChannelClosedEvent }
func (*SupportChannelReactivated) Key() string     { return SupportChannelReactivatedEvent }

// Ticket Events
func (*TicketCreated) Key() string         { return TicketCreatedEvent }
func (*TicketUpdated) Key() string         { return TicketUpdatedEvent }
func (*TicketAssigned) Key() string        { return TicketAssignedEvent }
func (*TicketPriorityUpdated) Key() string { return TicketPriorityUpdatedEvent }
func (*TicketEscalated) Key() string       { return TicketEscalatedEvent }
func (*TicketResolved) Key() string        { return TicketResolvedEvent }
func (*TicketReopened) Key() string        { return TicketReopenedEvent }
func (*TicketClosed) Key() string          { return TicketClosedEvent }
func (*TicketsMerged) Key() string         { return TicketsMergedEvent }
func (*TicketsLinked) Key() string         { return TicketsLinkedEvent }

// Communication Events
func (*TicketReplyAdded) Key() string  { return TicketReplyAddedEvent }
func (*InternalNoteAdded) Key() string { return InternalNoteAddedEvent }

// AI Integration Events
func (*AISupportEnabled) Key() string      { return AISupportEnabledEvent }
func (*AIConfigurationUpdated) Key() string { return AIConfigurationUpdatedEvent }
func (*HandoffToHumanOccurred) Key() string { return HandoffToHumanOccurredEvent }
func (*HandoffToAIOccurred) Key() string    { return HandoffToAIOccurredEvent }
func (*AIResponseGenerated) Key() string    { return AIResponseGeneratedEvent }

// Knowledge Base Events
func (*KnowledgeArticleCreated) Key() string { return KnowledgeArticleCreatedEvent }
func (*ArticleLinkedToTicket) Key() string   { return ArticleLinkedToTicketEvent }
func (*ArticleRated) Key() string            { return ArticleRatedEvent }

// Analytics Events
func (*SLABreached) Key() string                  { return SLABreachedEvent }
func (*CustomerSatisfactionRecorded) Key() string { return CustomerSatisfactionRecordedEvent }
func (*ResolutionTimeRecorded) Key() string       { return ResolutionTimeRecordedEvent }