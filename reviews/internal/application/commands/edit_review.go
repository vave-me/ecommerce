package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/reviews/internal/domain"
)

type EditReview struct {
	ID      string
	Content string
}

type EditReviewHandler struct {
	reviews   domain.ReviewRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewEditReviewHandler(reviews domain.ReviewRepository, publisher ddd.EventPublisher[ddd.Event]) EditReviewHandler {
	return EditReviewHandler{
		reviews:   reviews,
		publisher: publisher,
	}
}

func (h EditReviewHandler) EditReview(ctx context.Context, cmd EditReview) error {
	review, err := h.reviews.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := review.Edit(cmd.Content)
	if err != nil {
		return err
	}

	err = h.reviews.Save(ctx, review)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
