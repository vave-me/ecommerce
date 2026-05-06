package followingpb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	FollowAggregateChannel = "middleman.following.events.Follow"
	WebSocketChannel       = "following.>"
	AddFollowCommand       = "following.AddFollow"
	ReplyAggregateChannel  = "middleman.following.events.Reply"
	ReplyAddedEvent        = "following.ReplyAdded"
	FollowAddedEvent       = "followingapi.FollowAdded"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Follow events

	if err := serde.Register(&AddFollow{}); err != nil {
		return err
	}

	if err := serde.Register(&FollowAdded{}); err != nil {

	}
	return nil
}

func (*AddFollow) Key() string   { return AddFollowCommand }
func (*FollowAdded) Key() string { return FollowAddedEvent }
