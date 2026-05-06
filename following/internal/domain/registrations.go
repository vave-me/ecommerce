package domain

import (
	"middleman/internal/es"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	FollowAddedEvent    = "following.FollowAdded"
	FollowRemovedEvent  = "following.FollowRemoved"
	FollowApprovedEvent = "following.FollowApproved"
	FollowEditedEvent   = "following.FollowEdited"
	FollowFlaggedEvent  = "following.FollowFlagged"
	FollowRejectedEvent = "following.FollowRejected"
)

func Registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Follow
	if err = serde.Register(Follow{}, func(v any) error {
		follow := v.(*Follow)
		follow.Aggregate = es.NewAggregate("", FollowAggregate)
		return nil
	}); err != nil {
		return
	}
	// follow events
	if err = serde.Register(FollowAdded{}); err != nil {
		return
	}
	if err = serde.Register(FollowEdited{}); err != nil {
		return
	}
	// follow snapshots
	if err = serde.RegisterKey(FollowV1{}.SnapshotName(), FollowV1{}); err != nil {
		return
	}

	return
}

func (Follow) Key() string { return FollowAggregate }

// Key implements registry.Registerable
func (FollowFlagged) Key() string  { return FollowFlaggedEvent }
func (FollowRejected) Key() string { return FollowRejectedEvent }

// Key implements registry.Registerable
func (FollowApproved) Key() string { return FollowApprovedEvent }

// Key implements registry.Registerable
func (FollowRemoved) Key() string { return FollowRemovedEvent }

// Key implements registry.Registerable
func (FollowEdited) Key() string { return FollowEditedEvent }
func (FollowAdded) Key() string  { return FollowAddedEvent }
