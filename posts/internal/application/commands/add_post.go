package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/posts/internal/domain"
)

type AddPost struct {
	ID           string
	Name         string
	Description  string
	TypeOfPost   domain.TypeOfPost
	CategoryID   string
	CategorySlug string
	UserID       string
	UserType     domain.UserType
	Tags         []string
	Status       domain.PostStatus
	Thumbnail    string
	Lat          float64
	Lng          float64
}

type AddPostHandler struct {
	posts     domain.PostRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddPostHandler(
	posts domain.PostRepository, publisher ddd.EventPublisher[ddd.Event]) AddPostHandler {

	return AddPostHandler{
		posts:     posts,
		publisher: publisher,
	}
}

func (h AddPostHandler) AddPost(ctx context.Context, cmd AddPost) error {
	post, err := h.posts.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding post")
	}

	event, err := post.InitPost(cmd.ID, cmd.Name, cmd.Description, cmd.TypeOfPost, cmd.CategoryID, cmd.CategorySlug, cmd.UserID, cmd.UserType, cmd.Tags, cmd.Status, cmd.Thumbnail, cmd.Lat, cmd.Lng)
	if err != nil {
		return errors.Wrap(err, "initializing post")
	}

	err = h.posts.Save(ctx, post)
	if err != nil {
		return errors.Wrap(err, "error adding post")
	}

	return h.publisher.Publish(ctx, event)
}
