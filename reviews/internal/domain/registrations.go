package domain

import (
	"middleman/internal/es"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	ReviewAddedEvent    = "reviews.ReviewAdded"
	ReviewRemovedEvent  = "reviews.ReviewRemoved"
	ReviewApprovedEvent = "reviews.ReviewApproved"
	ReviewEditedEvent   = "reviews.ReviewEdited"
	ReviewFlaggedEvent  = "reviews.ReviewFlagged"
	ReviewRejectedEvent = "reviews.ReviewRejected"
)

func Registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Review
	if err = serde.Register(Review{}, func(v any) error {
		review := v.(*Review)
		review.Aggregate = es.NewAggregate("", ReviewAggregate)
		return nil
	}); err != nil {
		return
	}
	// review events
	if err = serde.Register(ReviewAdded{}); err != nil {
		return
	}
	if err = serde.Register(ReviewEdited{}); err != nil {
		return
	}
	// review snapshots
	if err = serde.RegisterKey(ReviewV1{}.SnapshotName(), ReviewV1{}); err != nil {
		return
	}

	return
}

func (Review) Key() string { return ReviewAggregate }

// Key implements registry.Registerable
func (ReviewFlagged) Key() string  { return ReviewFlaggedEvent }
func (ReviewRejected) Key() string { return ReviewRejectedEvent }

// Key implements registry.Registerable
func (ReviewApproved) Key() string { return ReviewApprovedEvent }

// Key implements registry.Registerable
func (ReviewRemoved) Key() string { return ReviewRemovedEvent }

// Key implements registry.Registerable
func (ReviewEdited) Key() string { return ReviewEditedEvent }
func (ReviewAdded) Key() string  { return ReviewAddedEvent }
