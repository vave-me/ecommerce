package commands

import (
	"context"
	"middleman/comments/internal/domain"
	"middleman/internal/ddd"
)

type ApproveComment struct {
	ID       string
	Approval bool
}

type ApproveCommentHandler struct {
	comments  domain.CommentRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewApproveCommentHandler(comments domain.CommentRepository, publisher ddd.EventPublisher[ddd.Event]) ApproveCommentHandler {
	return ApproveCommentHandler{
		comments:  comments,
		publisher: publisher,
	}
}

func (h ApproveCommentHandler) ApproveComment(ctx context.Context, cmd ApproveComment) error {
	comment, err := h.comments.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	//TODO change to approve and implement
	event, err := comment.Approve(cmd.Approval)
	if err != nil {
		return err
	}

	err = h.comments.Save(ctx, comment)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
