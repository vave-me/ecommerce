package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/reviews/internal/domain"
)

type FlagReview struct {
	ID      string
	Flagged bool
}

type FlagReviewHandler struct {
	reviews   domain.ReviewRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewFlagReviewHandler(reviews domain.ReviewRepository, publisher ddd.EventPublisher[ddd.Event]) FlagReviewHandler {
	return FlagReviewHandler{
		reviews:   reviews,
		publisher: publisher,
	}
}

func (h FlagReviewHandler) FlagReview(ctx context.Context, cmd FlagReview) error {
	review, err := h.reviews.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	//TODO change to flag and implement
	event, err := review.Flag(cmd.Flagged)
	if err != nil {
		return err
	}

	err = h.reviews.Save(ctx, review)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
