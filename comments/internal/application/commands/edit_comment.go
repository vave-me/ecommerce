package commands

import (
	"context"
	"middleman/comments/internal/domain"
	"middleman/internal/ddd"
)

type EditComment struct {
	ID      string
	Content string
}

type EditCommentHandler struct {
	comments  domain.CommentRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewEditCommentHandler(comments domain.CommentRepository, publisher ddd.EventPublisher[ddd.Event]) EditCommentHandler {
	return EditCommentHandler{
		comments:  comments,
		publisher: publisher,
	}
}

func (h EditCommentHandler) EditComment(ctx context.Context, cmd EditComment) error {
	comment, err := h.comments.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := comment.Edit(cmd.Content)
	if err != nil {
		return err
	}

	err = h.comments.Save(ctx, comment)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
