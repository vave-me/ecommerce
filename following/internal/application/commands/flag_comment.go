package commands

import (
	"context"
	"middleman/following/internal/domain"
	"middleman/internal/ddd"
)

type FlagFollow struct {
	ID      string
	Flagged bool
}

type FlagFollowHandler struct {
	following domain.FollowRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewFlagFollowHandler(following domain.FollowRepository, publisher ddd.EventPublisher[ddd.Event]) FlagFollowHandler {
	return FlagFollowHandler{
		following: following,
		publisher: publisher,
	}
}

func (h FlagFollowHandler) FlagFollow(ctx context.Context, cmd FlagFollow) error {
	follow, err := h.following.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	//TODO change to flag and implement
	event, err := follow.Flag(cmd.Flagged)
	if err != nil {
		return err
	}

	err = h.following.Save(ctx, follow)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
