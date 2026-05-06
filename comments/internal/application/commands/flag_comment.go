package commands

import (
	"context"
	"middleman/comments/internal/domain"
	"middleman/internal/ddd"
)

type FlagComment struct {
	ID      string
	Flagged bool
}

type FlagCommentHandler struct {
	comments  domain.CommentRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewFlagCommentHandler(comments domain.CommentRepository, publisher ddd.EventPublisher[ddd.Event]) FlagCommentHandler {
	return FlagCommentHandler{
		comments:  comments,
		publisher: publisher,
	}
}

func (h FlagCommentHandler) FlagComment(ctx context.Context, cmd FlagComment) error {
	comment, err := h.comments.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	//TODO change to flag and implement
	event, err := comment.Flag(cmd.Flagged)
	if err != nil {
		return err
	}

	err = h.comments.Save(ctx, comment)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
