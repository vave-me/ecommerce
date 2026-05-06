package commands

import (
	"context"
	"middleman/following/internal/domain"
	"middleman/internal/ddd"
)

type RemoveFollow struct {
	ID string
}

type RemoveFollowHandler struct {
	following domain.FollowRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRemoveFollowHandler(following domain.FollowRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveFollowHandler {
	return RemoveFollowHandler{
		following: following,
		publisher: publisher,
	}
}

func (h RemoveFollowHandler) RemoveFollow(ctx context.Context, cmd RemoveFollow) error {
	follow, err := h.following.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := follow.Remove(cmd.ID)
	if err != nil {
		return err
	}

	err = h.following.Save(ctx, follow)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
