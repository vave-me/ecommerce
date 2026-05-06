package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/posts/internal/domain"
)

type UpdatePostThumbnail struct {
	ID        string
	Thumbnail string
}

type UpdatePostThumbnailHandler struct {
	posts     domain.PostRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewUpdatePostThumbnailHandler(
	posts domain.PostRepository, publisher ddd.EventPublisher[ddd.Event]) UpdatePostThumbnailHandler {

	return UpdatePostThumbnailHandler{
		posts:     posts,
		publisher: publisher,
	}
}

func (h UpdatePostThumbnailHandler) UpdatePostThumbnail(ctx context.Context, cmd UpdatePostThumbnail) error {
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
