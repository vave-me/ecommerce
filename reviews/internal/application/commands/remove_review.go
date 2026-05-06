package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/reviews/internal/domain"
)

type RemoveReview struct {
	ID string
}

type RemoveReviewHandler struct {
	reviews   domain.ReviewRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRemoveReviewHandler(reviews domain.ReviewRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveReviewHandler {
	return RemoveReviewHandler{
		reviews:   reviews,
		publisher: publisher,
	}
}

func (h RemoveReviewHandler) RemoveReview(ctx context.Context, cmd RemoveReview) error {
	review, err := h.reviews.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := review.Remove(cmd.ID)
	if err != nil {
		return err
	}

	err = h.reviews.Save(ctx, review)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
