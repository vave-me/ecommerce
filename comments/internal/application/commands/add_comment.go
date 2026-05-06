package commands

import (
	"context"
	"fmt"
	"github.com/stackus/errors"
	"middleman/comments/internal/domain"
	"middleman/internal/ddd"
)

type AddComment struct {
	ID         string
	SenderID   string
	ItemID     string
	ItemType   domain.ItemType
	Content    string
	CategoryID string
	ParentID   string
}

type AddCommentHandler struct {
	comments  domain.CommentRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddCommentHandler(
	comments domain.CommentRepository, publisher ddd.EventPublisher[ddd.Event]) AddCommentHandler {

	return AddCommentHandler{
		comments:  comments,
		publisher: publisher,
	}
}

func (h AddCommentHandler) AddComment(ctx context.Context, cmd AddComment) error {

	fmt.Printf("ItemID %s, SenderID: %s, Content: %s, ParentID: %s", cmd.ItemID, cmd.SenderID, cmd.Content, cmd.ItemID)
	comment, err := h.comments.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding comment")
	}
	fmt.Println("start initialization comment ")
	event, err := comment.InitComment(cmd.ID, cmd.SenderID, cmd.ItemID, cmd.ItemType, cmd.Content, cmd.CategoryID, cmd.ParentID)
	if err != nil {
		return errors.Wrap(err, "initializing comment")
	}
	err = h.comments.Save(ctx, comment)
	if err != nil {
		return errors.Wrap(err, "error adding comment")
	}

	return h.publisher.Publish(ctx, event)
}
