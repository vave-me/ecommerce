package domain

import (
	"fmt"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const FollowAggregate = "following.Follow"

var (
	ErrUserIdIsBlank         = errors.Wrap(errors.ErrBadRequest, "the sender id cannot be blank")
	ErrFollowedUserIdIsBlank = errors.Wrap(errors.ErrBadRequest, "the item id cannot be blank")
	ErrItemPriceIsNegative   = errors.Wrap(errors.ErrBadRequest, "the item price cannot be negative")
	ErrNotAPriceIncrease     = errors.Wrap(errors.ErrBadRequest, "the price change would be a decrease")
	ErrNotAPriceDecrease     = errors.Wrap(errors.ErrBadRequest, "the price change would be an increase")
)

type Follow struct {
	es.Aggregate
	UserID           string
	FollowedUserID   string
	FollowedUserType string
	CategoryID       string
	ParentID         string
	Content          string
	Approved         bool
	Flagged          bool
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Follow)(nil)

func NewFollow(id string) *Follow {
	return &Follow{
		Aggregate: es.NewAggregate(id, FollowAggregate),
	}
}
func (c *Follow) InitFollow(id, userID, followedUserID string, followedUserType FollowedUserType, content, categoryID, parentID string) (ddd.Event, error) {
	fmt.Printf("INIT COMMENT ")

	c.AddEvent(FollowAddedEvent, &FollowAdded{
		UserID:           userID,
		FollowedUserID:   followedUserID,
		FollowedUserType: followedUserType,
		Content:          content,
		CategoryID:       categoryID,
		ParentID:         parentID,
		Approved:         true,
		Flagged:          false,
	})

	return ddd.NewEvent(FollowAddedEvent, c), nil
}
func (c *Follow) Edit(content string) (ddd.Event, error) {
	c.AddEvent(FollowEditedEvent, &FollowEdited{
		Content: content,
	})

	return ddd.NewEvent(FollowEditedEvent, c), nil
}

func (c *Follow) Reject(rejected bool) (ddd.Event, error) {
	c.AddEvent(FollowRejectedEvent, &FollowRejected{
		Approve: rejected,
	})

	return ddd.NewEvent(FollowEditedEvent, c), nil
}
func (c *Follow) Approve(approval bool) (ddd.Event, error) {
	c.AddEvent(FollowApprovedEvent, &FollowApproved{
		Approved: approval,
	})

	return ddd.NewEvent(FollowApprovedEvent, c), nil
}

func (c *Follow) Flag(flagged bool) (ddd.Event, error) {
	c.AddEvent(FollowFlaggedEvent, &FollowFlagged{
		Flagged: flagged,
	})

	return ddd.NewEvent(FollowApprovedEvent, c), nil
}

func (c *Follow) Remove(followID string) (ddd.Event, error) {
	c.AddEvent(FollowRemovedEvent, &FollowRemoved{
		FollowID: followID,
	})

	return ddd.NewEvent(FollowRemovedEvent, c), nil
}

// ApplyEvent implements es.EventApplier
func (c *Follow) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *FollowAdded:
		c.UserID = payload.UserID
		c.FollowedUserID = payload.FollowedUserID
		c.FollowedUserType = payload.FollowedUserType.String()
		c.CategoryID = payload.CategoryID
		c.ParentID = payload.ParentID
		c.Content = payload.Content
		c.Flagged = false
		c.Approved = true
	case *FollowEdited:
		c.Content = payload.Content
	case FollowRejected:
		c.Approved = false
	case FollowFlagged:
		c.Flagged = true
	case FollowApproved:
		c.Approved = true
		c.Flagged = false
	case *FollowRemoved:

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", c, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (c *Follow) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *FollowV1:
		c.UserID = ss.UserID
		c.FollowedUserID = ss.FollowedUserID
		c.FollowedUserType = ss.FollowedUserType
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
func (s Follow) ToSnapshot() es.Snapshot {
	return FollowV1{
		ID:               s.ID(),
		UserID:           s.UserID,
		FollowedUserID:   s.FollowedUserID,
		FollowedUserType: s.FollowedUserType,
		Content:          s.Content,
		CategoryID:       s.CategoryID,
		ParentID:         s.ParentID,
		Approved:         s.Approved,
		Flagged:          s.Flagged,
	}
}
