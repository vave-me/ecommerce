package commands

import (
	"context"
	"fmt"
	"github.com/stackus/errors"
	"middleman/following/internal/domain"
	"middleman/internal/ddd"
)

type AddFollow struct {
	ID               string
	UserID           string
	FollowedUserID   string
	FollowedUserType domain.FollowedUserType
	Content          string
	CategoryID       string
	ParentID         string
}

type AddFollowHandler struct {
	following domain.FollowRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddFollowHandler(
	following domain.FollowRepository, publisher ddd.EventPublisher[ddd.Event]) AddFollowHandler {

	return AddFollowHandler{
		following: following,
		publisher: publisher,
	}
}

func (h AddFollowHandler) AddFollow(ctx context.Context, cmd AddFollow) error {

	fmt.Printf("FollowedUserID %s, UserID: %s, Content: %s, ParentID: %s", cmd.FollowedUserID, cmd.UserID, cmd.Content, cmd.FollowedUserID)
	follow, err := h.following.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding follow")
	}
	fmt.Println("start initialization follow ")
	event, err := follow.InitFollow(cmd.ID, cmd.UserID, cmd.FollowedUserID, cmd.FollowedUserType, cmd.Content, cmd.CategoryID, cmd.ParentID)
	if err != nil {
		return errors.Wrap(err, "initializing follow")
	}
	err = h.following.Save(ctx, follow)
	if err != nil {
		return errors.Wrap(err, "error adding follow")
	}

	return h.publisher.Publish(ctx, event)
}
