package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/posts/internal/domain"
)

type AddPostThumbnail struct {
	ID        string
	Thumbnail string

}

type AddPostThumbnailHandler struct {
	posts     domain.PostRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddPostThumbnailHandler(
	posts domain.PostRepository, publisher ddd.EventPublisher[ddd.Event]) AddPostThumbnailHandler {

	return AddPostThumbnailHandler{
		posts:     posts,
		publisher: publisher,
	}
}

func (h AddPostThumbnailHandler) AddPostThumbnail(ctx context.Context, cmd AddPostThumbnail) error {
	post, err := h.posts.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding post")
	}

	event, err := post.AddThumbnail(cmd.Thumbnail)
	if err != nil {
		return errors.Wrap(err, "initializing post")
	}

	err = h.posts.Save(ctx, post)
	if err != nil {
		return errors.Wrap(err, "error adding post")
	}

	return h.publisher.Publish(ctx, event)
}
