package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/reviews/internal/domain"
)

type RejectReview struct {
	ID       string
	Rejected bool
}

type RejectReviewHandler struct {
	reviews   domain.ReviewRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRejectReviewHandler(reviews domain.ReviewRepository, publisher ddd.EventPublisher[ddd.Event]) RejectReviewHandler {
	return RejectReviewHandler{
		reviews:   reviews,
		publisher: publisher,
	}
}

func (h RejectReviewHandler) RejectReview(ctx context.Context, cmd RejectReview) error {
	review, err := h.reviews.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	//TODO change to approve and implement
	event, err := review.Reject(cmd.Rejected)
	if err != nil {
		return err
	}

	err = h.reviews.Save(ctx, review)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
