package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/posts/internal/domain"
)

// ArchivePost marks a post as archived or inactive.
type ArchivePost struct {
	ID     string
	UserID string
}

type ArchivePostHandler struct {
	posts     domain.PostRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewArchivePostHandler(
	posts domain.PostRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ArchivePostHandler {
	return ArchivePostHandler{
		posts:     posts,
		publisher: publisher,
	}
}

func (h ArchivePostHandler) ArchivePost(ctx context.Context, cmd ArchivePost) error {
	post, err := h.posts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := post.Archive(cmd.UserID) // domain: post.Archive(userSellerID)
	if err != nil {
		return err
	}

	if err = h.posts.Save(ctx, post); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
