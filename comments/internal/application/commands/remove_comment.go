package commands

import (
	"context"
	"middleman/comments/internal/domain"
	"middleman/internal/ddd"
)

type RemoveComment struct {
	ID string
}

type RemoveCommentHandler struct {
	comments  domain.CommentRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRemoveCommentHandler(comments domain.CommentRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveCommentHandler {
	return RemoveCommentHandler{
		comments:  comments,
		publisher: publisher,
	}
}

func (h RemoveCommentHandler) RemoveComment(ctx context.Context, cmd RemoveComment) error {
	comment, err := h.comments.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := comment.Remove(cmd.ID)
	if err != nil {
		return err
	}

	err = h.comments.Save(ctx, comment)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
