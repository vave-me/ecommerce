package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/posts/internal/domain"
)

type UpdatePost struct {
	ID           string
	Name         string
	TypeOfPost   domain.TypeOfPost
	CategoryID   string
	CategorySlug string
	Description  string
	Tags         []string
	Status       domain.PostStatus
	Thumbnail    string
}

type UpdatePostHandler struct {
	posts     domain.PostRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewUpdatePostHandler(
	posts domain.PostRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UpdatePostHandler {
	return UpdatePostHandler{
		posts:     posts,
		publisher: publisher,
	}
}

func (h UpdatePostHandler) UpdatePost(ctx context.Context, cmd UpdatePost) error {
	post, err := h.posts.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// domain method: post.Update(...) or partial setters
	event, err := post.Update(cmd.Name, cmd.Description, cmd.TypeOfPost, cmd.CategoryID, cmd.CategorySlug, cmd.Tags, cmd.Status, cmd.Thumbnail)
	if err != nil {
		return err
	}

	if err = h.posts.Save(ctx, post); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
