package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/posts/internal/domain"
)

type RemovePost struct {
	ID     string
	UserID string
}

type RemovePostHandler struct {
	posts     domain.PostRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRemovePostHandler(posts domain.PostRepository, publisher ddd.EventPublisher[ddd.Event]) RemovePostHandler {
	return RemovePostHandler{
		posts:     posts,
		publisher: publisher,
	}
}

func (h RemovePostHandler) RemovePost(ctx context.Context, cmd RemovePost) error {
	post, err := h.posts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := post.Remove(cmd.UserID)
	if err != nil {
		return err
	}

	err = h.posts.Save(ctx, post)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
