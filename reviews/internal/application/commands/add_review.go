package commands

import (
	"context"
	"fmt"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/reviews/internal/domain"
)

type AddReview struct {
	ID         string
	SenderID   string
	ItemID     string
	ItemType   domain.ItemType
	Content    string
	CategoryID string
	ParentID   string
}

type AddReviewHandler struct {
	reviews   domain.ReviewRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddReviewHandler(
	reviews domain.ReviewRepository, publisher ddd.EventPublisher[ddd.Event]) AddReviewHandler {

	return AddReviewHandler{
		reviews:   reviews,
		publisher: publisher,
	}
}

func (h AddReviewHandler) AddReview(ctx context.Context, cmd AddReview) error {

	fmt.Printf("ItemID %s, SenderID: %s, Content: %s, ParentID: %s", cmd.ItemID, cmd.SenderID, cmd.Content, cmd.ItemID)
	review, err := h.reviews.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding review")
	}
	fmt.Println("start initialization review ")
	event, err := review.InitReview(cmd.ID, cmd.SenderID, cmd.ItemID, cmd.ItemType, cmd.Content, cmd.CategoryID, cmd.ParentID)
	if err != nil {
		return errors.Wrap(err, "initializing review")
	}
	err = h.reviews.Save(ctx, review)
	if err != nil {
		return errors.Wrap(err, "error adding review")
	}

	return h.publisher.Publish(ctx, event)
}
