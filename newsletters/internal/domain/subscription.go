package domain

import (
	"time"

	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const SubscriptionAggregate = "newsletters.Subscription"

type Subscription struct {
	es.Aggregate
	UserID         string
	NewsletterID   string
	Status         SubscriptionStatus
	Preferences    SubscriptionPreferences
	SubscribedAt   time.Time
	UnsubscribedAt *time.Time
}

type SubscriptionPreferences struct {
	FrequencyOverride *NewsletterFrequency
	Topics            []string
	Format            ContentFormat
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Subscription)(nil)

func NewSubscription(id string) *Subscription {
	return &Subscription{
		Aggregate: es.NewAggregate(id, SubscriptionAggregate),
	}
}

func (s *Subscription) Subscribe(userID, newsletterID string, preferences SubscriptionPreferences) (ddd.Event, error) {
	if userID == "" {
		return nil, errors.ErrBadRequest.Msg("user ID is required")
	}
	if newsletterID == "" {
		return nil, errors.ErrBadRequest.Msg("newsletter ID is required")
	}

	s.AddEvent(SubscribedEvent, &Subscribed{
		UserID:            userID,
		NewsletterID:      newsletterID,
		FrequencyOverride: preferences.FrequencyOverride.String(),
		Topics:            preferences.Topics,
		Format:            preferences.Format.String(),
	})

	return ddd.NewEvent(SubscribedEvent, s), nil
}

func (s *Subscription) UpdatePreferences(preferences SubscriptionPreferences) (ddd.Event, error) {
	if s.Status != ActiveStatus {
		return nil, errors.ErrBadRequest.Msg("cannot update preferences for inactive subscription")
	}

	s.AddEvent(PreferencesUpdatedEvent, &PreferencesUpdated{
		FrequencyOverride: preferences.FrequencyOverride.String(),
		Topics:            preferences.Topics,
		Format:            preferences.Format.String(),
	})

	return ddd.NewEvent(PreferencesUpdatedEvent, s), nil
}

func (s *Subscription) Pause() (ddd.Event, error) {
	if s.Status != ActiveStatus {
		return nil, errors.ErrBadRequest.Msg("can only pause active subscriptions")
	}

	s.AddEvent(SubscriptionPausedEvent, &SubscriptionPaused{})
	return ddd.NewEvent(SubscriptionPausedEvent, s), nil
}

func (s *Subscription) Resume() (ddd.Event, error) {
	if s.Status != PausedStatus {
		return nil, errors.ErrBadRequest.Msg("can only resume paused subscriptions")
	}

	s.AddEvent(SubscriptionResumedEvent, &SubscriptionResumed{})
	return ddd.NewEvent(SubscriptionResumedEvent, s), nil
}

func (s *Subscription) Unsubscribe(reason string) (ddd.Event, error) {
	if s.Status == UnsubscribedStatus {
		return nil, errors.ErrBadRequest.Msg("already unsubscribed")
	}

	s.AddEvent(UnsubscribedEvent, &Unsubscribed{
		Reason: reason,
	})

	return ddd.NewEvent(UnsubscribedEvent, s), nil
}

// Key implements registry.Registerable
func (Subscription) Key() string { return SubscriptionAggregate }

// ApplyEvent implements es.EventApplier
func (s *Subscription) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *Subscribed:
		s.UserID = payload.UserID
		s.NewsletterID = payload.NewsletterID
		s.Status = ActiveStatus
		s.Preferences = SubscriptionPreferences{
			FrequencyOverride: ToNewsletterFrequencyPtr(payload.FrequencyOverride),
			Topics:            payload.Topics,
			Format:            ToContentFormat(payload.Format),
		}
		s.SubscribedAt = event.OccurredAt()

	case *PreferencesUpdated:
		s.Preferences = SubscriptionPreferences{
			FrequencyOverride: ToNewsletterFrequencyPtr(payload.FrequencyOverride),
			Topics:            payload.Topics,
			Format:            ToContentFormat(payload.Format),
		}

	case *SubscriptionPaused:
		s.Status = PausedStatus

	case *SubscriptionResumed:
		s.Status = ActiveStatus

	case *Unsubscribed:
		s.Status = UnsubscribedStatus
		unsubAt := event.OccurredAt()
		s.UnsubscribedAt = &unsubAt

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", s, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (s *Subscription) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *SubscriptionV1:
		s.UserID = ss.UserID
		s.NewsletterID = ss.NewsletterID
		s.Status = ToSubscriptionStatus(ss.Status)
		s.Preferences = SubscriptionPreferences{
			FrequencyOverride: ToNewsletterFrequencyPtr(ss.FrequencyOverride),
			Topics:            ss.Topics,
			Format:            ToContentFormat(ss.Format),
		}
		s.SubscribedAt = ss.SubscribedAt
		s.UnsubscribedAt = ss.UnsubscribedAt

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", s, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (s Subscription) ToSnapshot() es.Snapshot {
	var freqOverride string
	if s.Preferences.FrequencyOverride != nil {
		freqOverride = s.Preferences.FrequencyOverride.String()
	}

	return SubscriptionV1{
		UserID:            s.UserID,
		NewsletterID:      s.NewsletterID,
		Status:            s.Status.String(),
		FrequencyOverride: freqOverride,
		Topics:            s.Preferences.Topics,
		Format:            s.Preferences.Format.String(),
		SubscribedAt:      s.SubscribedAt,
		UnsubscribedAt:    s.UnsubscribedAt,
	}
}