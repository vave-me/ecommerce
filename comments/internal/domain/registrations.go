package domain

import (
	"middleman/internal/es"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	CommentAddedEvent    = "comments.CommentAdded"
	CommentRemovedEvent  = "comments.CommentRemoved"
	CommentApprovedEvent = "comments.CommentApproved"
	CommentEditedEvent   = "comments.CommentEdited"
	CommentFlaggedEvent  = "comments.CommentFlagged"
	CommentRejectedEvent = "comments.CommentRejected"
)

func Registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Comment
	if err = serde.Register(Comment{}, func(v any) error {
		comment := v.(*Comment)
		comment.Aggregate = es.NewAggregate("", CommentAggregate)
		return nil
	}); err != nil {
		return
	}
	// comment events
	if err = serde.Register(CommentAdded{}); err != nil {
		return
	}
	if err = serde.Register(CommentEdited{}); err != nil {
		return
	}
	// comment snapshots
	if err = serde.RegisterKey(CommentV1{}.SnapshotName(), CommentV1{}); err != nil {
		return
	}

	return
}

func (Comment) Key() string { return CommentAggregate }

// Key implements registry.Registerable
func (CommentFlagged) Key() string  { return CommentFlaggedEvent }
func (CommentRejected) Key() string { return CommentRejectedEvent }

// Key implements registry.Registerable
func (CommentApproved) Key() string { return CommentApprovedEvent }

// Key implements registry.Registerable
func (CommentRemoved) Key() string { return CommentRemovedEvent }

// Key implements registry.Registerable
func (CommentEdited) Key() string { return CommentEditedEvent }
func (CommentAdded) Key() string  { return CommentAddedEvent }
