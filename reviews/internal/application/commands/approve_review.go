package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/reviews/internal/domain"
)

type ApproveReview struct {
	ID       string
	Approval bool
}

type ApproveReviewHandler struct {
	reviews   domain.ReviewRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewApproveReviewHandler(reviews domain.ReviewRepository, publisher ddd.EventPublisher[ddd.Event]) ApproveReviewHandler {
	return ApproveReviewHandler{
		reviews:   reviews,
		publisher: publisher,
	}
}

func (h ApproveReviewHandler) ApproveReview(ctx context.Context, cmd ApproveReview) error {
	review, err := h.reviews.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	//TODO change to approve and implement
	event, err := review.Approve(cmd.Approval)
	if err != nil {
		return err
	}

	err = h.reviews.Save(ctx, review)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
