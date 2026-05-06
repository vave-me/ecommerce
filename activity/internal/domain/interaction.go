package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const InteractionAggregate = "activity.Interaction"

var (
	ErrInteractionAlreadyCreated = errors.Wrap(errors.ErrBadRequest, "the interaction cannot be recreated")
	ErrInteractionHasNoPayload   = errors.Wrap(errors.ErrBadRequest, "the interaction has no payload")
	ErrItemIDCannotBeBlank       = errors.Wrap(errors.ErrBadRequest, "the item id cannot be blank")
)

type Interaction struct {
	es.Aggregate
	ActivityID string
	ItemID     string // ID of the item being liked or disliked (e.g., listing, comment, product)
	ItemType   string
	ActionType string // The action, either "like" or "dislike"
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Interaction)(nil)

func NewInteraction(id string) *Interaction {
	return &Interaction{
		Aggregate: es.NewAggregate(id, InteractionAggregate),
	}
}

func (i *Interaction) AddInteraction(activityID string, itemID string, itemType string, action string) (ddd.Event, error) {

	i.AddEvent(InteractionAddedEvent, &InteractionAdded{
		ActivityID: activityID,
		ItemID:     itemID,
		ItemType:   itemType,
		ActionType: action,
	})

	return ddd.NewEvent(InteractionAddedEvent, i), nil
}

func (i *Interaction) Remove(activityID string, itemID string) (ddd.Event, error) {

	if ItemType(i.ItemType) != UnknownType {
		return nil, ErrInteractionAlreadyCreated
	}
	if itemID == "" {
		return nil, ErrItemIDCannotBeBlank
	}

	i.AddEvent(InteractionRemovedEvent, &InteractionRemoved{
		ActivityID: activityID,
		ItemID:     itemID,
	})

	return ddd.NewEvent(InteractionAddedEvent, i), nil
}

func (i *Interaction) Update(activityID string, itemID string, actionType string) (ddd.Event, error) {

	if ItemType(i.ItemType) != UnknownType {
		return nil, ErrInteractionAlreadyCreated
	}
	if itemID == "" {
		return nil, ErrItemIDCannotBeBlank
	}

	i.AddEvent(InteractionUpdatedEvent, &InteractionUpdated{
		ActivityID: activityID,
		ItemID:     itemID,
		ActionType: actionType,
	})

	return ddd.NewEvent(InteractionAddedEvent, i), nil
}
func (b *Interaction) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *InteractionAdded:
		b.ActivityID = payload.ActivityID
		b.ItemID = payload.ItemID
		b.ItemType = payload.ItemType
		b.ActionType = payload.ActionType

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", b, event.EventName(), payload)
	}

	return nil
}

func (b *Interaction) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *InteractionVi:
		b.ActivityID = ss.ActivityID
		b.ItemID = ss.ItemID
		b.ActionType = ss.ActionType
		b.ItemType = ss.ItemType

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", b, snapshot)
	}

	return nil
}

func (b *Interaction) ToSnapshot() es.Snapshot {
	return &InteractionVi{
		ActivityID: b.ActivityID,
		ItemID:     b.ItemID,
		ActionType: b.ActionType,
		ItemType:   b.ItemType,
	}
}
