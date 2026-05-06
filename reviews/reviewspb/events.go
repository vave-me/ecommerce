package reviewspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	ReviewAggregateChannel = "middleman.reviews.events.Review"
	WebSocketChannel       = "reviews.>"
	AddReviewCommand       = "reviews.AddReview"
	ReplyAggregateChannel  = "middleman.reviews.events.Reply"
	ReplyAddedEvent        = "reviews.ReplyAdded"
	ReviewAddedEvent       = "reviewsapi.ReviewAdded"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Review events

	if err := serde.Register(&AddReview{}); err != nil {
		return err
	}

	if err := serde.Register(&ReviewAdded{}); err != nil {

	}
	return nil
}

func (*AddReview) Key() string   { return AddReviewCommand }
func (*ReviewAdded) Key() string { return ReviewAddedEvent }
