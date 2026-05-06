package commentspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	CommentAggregateChannel = "middleman.comments.events.Comment"
	WebSocketChannel        = "comments.>"
	AddCommentCommand       = "comments.AddComment"
	ReplyAggregateChannel   = "middleman.comments.events.Reply"
	ReplyAddedEvent         = "comments.ReplyAdded"
	CommentAddedEvent       = "commentsapi.CommentAdded"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Comment events

	if err := serde.Register(&AddComment{}); err != nil {
		return err
	}

	if err := serde.Register(&CommentAdded{}); err != nil {

	}
	return nil
}

func (*AddComment) Key() string   { return AddCommentCommand }
func (*CommentAdded) Key() string { return CommentAddedEvent }
