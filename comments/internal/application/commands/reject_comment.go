package commands

import (
	"context"
	"middleman/comments/internal/domain"
	"middleman/internal/ddd"
)

type RejectComment struct {
	ID       string
	Rejected bool
}

type RejectCommentHandler struct {
	comments  domain.CommentRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRejectCommentHandler(comments domain.CommentRepository, publisher ddd.EventPublisher[ddd.Event]) RejectCommentHandler {
	return RejectCommentHandler{
		comments:  comments,
		publisher: publisher,
	}
}

func (h RejectCommentHandler) RejectComment(ctx context.Context, cmd RejectComment) error {
	comment, err := h.comments.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	//TODO change to approve and implement
	event, err := comment.Reject(cmd.Rejected)
	if err != nil {
		return err
	}

	err = h.comments.Save(ctx, comment)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
