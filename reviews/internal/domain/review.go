package domain

import (
	"fmt"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const ReviewAggregate = "reviews.Review"

var (
	ErrSenderIdIsBlank     = errors.Wrap(errors.ErrBadRequest, "the sender id cannot be blank")
	ErrItemIdIsBlank       = errors.Wrap(errors.ErrBadRequest, "the item id cannot be blank")
	ErrItemPriceIsNegative = errors.Wrap(errors.ErrBadRequest, "the item price cannot be negative")
	ErrNotAPriceIncrease   = errors.Wrap(errors.ErrBadRequest, "the price change would be a decrease")
	ErrNotAPriceDecrease   = errors.Wrap(errors.ErrBadRequest, "the price change would be an increase")
)

type Review struct {
	es.Aggregate
	SenderID   string
	ItemID     string
	ItemType   string
	CategoryID string
	ParentID   string
	Content    string
	Approved   bool
	Flagged    bool
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Review)(nil)

func NewReview(id string) *Review {
	return &Review{
		Aggregate: es.NewAggregate(id, ReviewAggregate),
	}
}
func (c *Review) InitReview(id, senderID, itemID string, itemType ItemType, content, categoryID, parentID string) (ddd.Event, error) {
	fmt.Printf("INIT REVIEW ")

	c.AddEvent(ReviewAddedEvent, &ReviewAdded{
		SenderID:   senderID,
		ItemID:     itemID,
		ItemType:   itemType,
		Content:    content,
		CategoryID: categoryID,
		ParentID:   parentID,
		Approved:   true,
		Flagged:    false,
	})

	return ddd.NewEvent(ReviewAddedEvent, c), nil
}
func (c *Review) Edit(content string) (ddd.Event, error) {
	c.AddEvent(ReviewEditedEvent, &ReviewEdited{
		Content: content,
	})

	return ddd.NewEvent(ReviewEditedEvent, c), nil
}

func (c *Review) Reject(rejected bool) (ddd.Event, error) {
	c.AddEvent(ReviewRejectedEvent, &ReviewRejected{
		Approve: rejected,
	})

	return ddd.NewEvent(ReviewEditedEvent, c), nil
}
func (c *Review) Approve(approval bool) (ddd.Event, error) {
	c.AddEvent(ReviewApprovedEvent, &ReviewApproved{
		Approved: approval,
	})

	return ddd.NewEvent(ReviewApprovedEvent, c), nil
}

func (c *Review) Flag(flagged bool) (ddd.Event, error) {
	c.AddEvent(ReviewFlaggedEvent, &ReviewFlagged{
		Flagged: flagged,
	})

	return ddd.NewEvent(ReviewApprovedEvent, c), nil
}

func (c *Review) Remove(reviewID string) (ddd.Event, error) {
	c.AddEvent(ReviewRemovedEvent, &ReviewRemoved{
		ReviewID: reviewID,
	})

	return ddd.NewEvent(ReviewRemovedEvent, c), nil
}

// ApplyEvent implements es.EventApplier
func (c *Review) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *ReviewAdded:
		c.SenderID = payload.SenderID
		c.ItemID = payload.ItemID
		c.ItemType = payload.ItemType.String()
		c.CategoryID = payload.CategoryID
		c.ParentID = payload.ParentID
		c.Content = payload.Content
		c.Flagged = false
		c.Approved = true
	case *ReviewEdited:
		c.Content = payload.Content
	case ReviewRejected:
		c.Approved = false
	case ReviewFlagged:
		c.Flagged = true
	case ReviewApproved:
		c.Approved = true
		c.Flagged = false
	case *ReviewRemoved:

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", c, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (c *Review) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ReviewV1:
		c.SenderID = ss.SenderID
		c.ItemID = ss.ItemID
		c.ItemType = ss.ItemType
		c.Content = ss.Content
		c.CategoryID = ss.CategoryID
		c.ParentID = ss.ParentID
		c.Approved = ss.Approved

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", c, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (s Review) ToSnapshot() es.Snapshot {
	return ReviewV1{
		ID:         s.ID(),
		SenderID:   s.SenderID,
		ItemID:     s.ItemID,
		ItemType:   s.ItemType,
		Content:    s.Content,
		CategoryID: s.CategoryID,
		ParentID:   s.ParentID,
		Approved:   s.Approved,
		Flagged:    s.Flagged,
	}
}
