package domain

import (
	"fmt"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const CommentAggregate = "comments.Comment"

var (
	ErrSenderIdIsBlank     = errors.Wrap(errors.ErrBadRequest, "the sender id cannot be blank")
	ErrItemIdIsBlank       = errors.Wrap(errors.ErrBadRequest, "the item id cannot be blank")
	ErrItemPriceIsNegative = errors.Wrap(errors.ErrBadRequest, "the item price cannot be negative")
	ErrNotAPriceIncrease   = errors.Wrap(errors.ErrBadRequest, "the price change would be a decrease")
	ErrNotAPriceDecrease   = errors.Wrap(errors.ErrBadRequest, "the price change would be an increase")
)

type Comment struct {
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
} = (*Comment)(nil)

func NewComment(id string) *Comment {
	return &Comment{
		Aggregate: es.NewAggregate(id, CommentAggregate),
	}
}
func (c *Comment) InitComment(id, senderID, itemID string, itemType ItemType, content, categoryID, parentID string) (ddd.Event, error) {
	fmt.Printf("INIT COMMENT ")

	c.AddEvent(CommentAddedEvent, &CommentAdded{
		SenderID:   senderID,
		ItemID:     itemID,
		ItemType:   itemType,
		Content:    content,
		CategoryID: categoryID,
		ParentID:   parentID,
		Approved:   true,
		Flagged:    false,
	})

	return ddd.NewEvent(CommentAddedEvent, c), nil
}
func (c *Comment) Edit(content string) (ddd.Event, error) {
	c.AddEvent(CommentEditedEvent, &CommentEdited{
		Content: content,
	})

	return ddd.NewEvent(CommentEditedEvent, c), nil
}

func (c *Comment) Reject(rejected bool) (ddd.Event, error) {
	c.AddEvent(CommentRejectedEvent, &CommentRejected{
		Approve: rejected,
	})

	return ddd.NewEvent(CommentEditedEvent, c), nil
}
func (c *Comment) Approve(approval bool) (ddd.Event, error) {
	c.AddEvent(CommentApprovedEvent, &CommentApproved{
		Approved: approval,
	})

	return ddd.NewEvent(CommentApprovedEvent, c), nil
}

func (c *Comment) Flag(flagged bool) (ddd.Event, error) {
	c.AddEvent(CommentFlaggedEvent, &CommentFlagged{
		Flagged: flagged,
	})

	return ddd.NewEvent(CommentApprovedEvent, c), nil
}

func (c *Comment) Remove(commentID string) (ddd.Event, error) {
	c.AddEvent(CommentRemovedEvent, &CommentRemoved{
		CommentID: commentID,
	})

	return ddd.NewEvent(CommentRemovedEvent, c), nil
}

// ApplyEvent implements es.EventApplier
func (c *Comment) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *CommentAdded:
		c.SenderID = payload.SenderID
		c.ItemID = payload.ItemID
		c.ItemType = payload.ItemType.String()
		c.CategoryID = payload.CategoryID
		c.ParentID = payload.ParentID
		c.Content = payload.Content
		c.Flagged = false
		c.Approved = true
	case *CommentEdited:
		c.Content = payload.Content
	case CommentRejected:
		c.Approved = false
	case CommentFlagged:
		c.Flagged = true
	case CommentApproved:
		c.Approved = true
		c.Flagged = false
	case *CommentRemoved:

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", c, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (c *Comment) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *CommentV1:
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
func (s Comment) ToSnapshot() es.Snapshot {
	return CommentV1{
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
