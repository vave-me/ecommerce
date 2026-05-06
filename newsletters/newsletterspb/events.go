package newsletterspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	// Aggregate channels
	NewsletterAggregateChannel    = "middleman.newsletters.events.Newsletter"
	SubscriptionAggregateChannel  = "middleman.newsletters.events.Subscription"
	EditionAggregateChannel       = "middleman.newsletters.events.Edition"
	TemplateAggregateChannel      = "middleman.newsletters.events.Template"
	
	// Newsletter events
	NewsletterCreatedEvent      = "newslettersapi.NewsletterCreated"
	NewsletterUpdatedEvent      = "newslettersapi.NewsletterUpdated"
	NewsletterActivatedEvent    = "newslettersapi.NewsletterActivated"
	NewsletterDeactivatedEvent  = "newslettersapi.NewsletterDeactivated"
	NewsletterDeletedEvent      = "newslettersapi.NewsletterDeleted"
	
	// Subscription events
	SubscribedEvent             = "newslettersapi.Subscribed"
	UnsubscribedEvent           = "newslettersapi.Unsubscribed"
	PreferencesUpdatedEvent     = "newslettersapi.PreferencesUpdated"
	SubscriptionPausedEvent     = "newslettersapi.SubscriptionPaused"
	SubscriptionResumedEvent    = "newslettersapi.SubscriptionResumed"
	
	// Edition events
	EditionCreatedEvent         = "newslettersapi.EditionCreated"
	EditionUpdatedEvent         = "newslettersapi.EditionUpdated"
	EditionScheduledEvent       = "newslettersapi.EditionScheduled"
	EditionSendingEvent         = "newslettersapi.EditionSending"
	EditionSentEvent            = "newslettersapi.EditionSent"
	
	// Template events
	TemplateCreatedEvent        = "newslettersapi.TemplateCreated"
	TemplateUpdatedEvent        = "newslettersapi.TemplateUpdated"
	TemplateDeletedEvent        = "newslettersapi.TemplateDeleted"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Newsletter events
	if err := serde.Register(&NewsletterCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&NewsletterUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&NewsletterActivated{}); err != nil {
		return err
	}
	if err := serde.Register(&NewsletterDeactivated{}); err != nil {
		return err
	}
	if err := serde.Register(&NewsletterDeleted{}); err != nil {
		return err
	}
	
	// Subscription events
	if err := serde.Register(&UserSubscribed{}); err != nil {
		return err
	}
	if err := serde.Register(&UserUnsubscribed{}); err != nil {
		return err
	}
	if err := serde.Register(&SubscriptionPreferencesUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&SubscriptionPaused{}); err != nil {
		return err
	}
	if err := serde.Register(&SubscriptionResumed{}); err != nil {
		return err
	}
	
	// Edition events
	if err := serde.Register(&EditionCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&EditionUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&EditionScheduled{}); err != nil {
		return err
	}
	if err := serde.Register(&EditionSending{}); err != nil {
		return err
	}
	if err := serde.Register(&EditionSent{}); err != nil {
		return err
	}
	
	// Template events
	if err := serde.Register(&TemplateCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&TemplateUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&TemplateDeleted{}); err != nil {
		return err
	}

	return nil
}

// Newsletter event keys
func (*NewsletterCreated) Key() string     { return NewsletterCreatedEvent }
func (*NewsletterUpdated) Key() string     { return NewsletterUpdatedEvent }
func (*NewsletterActivated) Key() string   { return NewsletterActivatedEvent }
func (*NewsletterDeactivated) Key() string { return NewsletterDeactivatedEvent }
func (*NewsletterDeleted) Key() string     { return NewsletterDeletedEvent }

// Subscription event keys
func (*UserSubscribed) Key() string           { return SubscribedEvent }
func (*UserUnsubscribed) Key() string         { return UnsubscribedEvent }
func (*SubscriptionPreferencesUpdated) Key() string   { return PreferencesUpdatedEvent }
func (*SubscriptionPaused) Key() string   { return SubscriptionPausedEvent }
func (*SubscriptionResumed) Key() string  { return SubscriptionResumedEvent }

// Edition event keys
func (*EditionCreated) Key() string   { return EditionCreatedEvent }
func (*EditionUpdated) Key() string   { return EditionUpdatedEvent }
func (*EditionScheduled) Key() string { return EditionScheduledEvent }
func (*EditionSending) Key() string   { return EditionSendingEvent }
func (*EditionSent) Key() string      { return EditionSentEvent }

// Template event keys
func (*TemplateCreated) Key() string { return TemplateCreatedEvent }
func (*TemplateUpdated) Key() string { return TemplateUpdatedEvent }
func (*TemplateDeleted) Key() string { return TemplateDeletedEvent }