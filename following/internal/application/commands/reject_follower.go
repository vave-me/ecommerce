package commands

import (
	"context"
	"middleman/following/internal/domain"
	"middleman/internal/ddd"
)

type RejectFollow struct {
	ID       string
	Rejected bool
}

type RejectFollowHandler struct {
	following domain.FollowRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRejectFollowHandler(following domain.FollowRepository, publisher ddd.EventPublisher[ddd.Event]) RejectFollowHandler {
	return RejectFollowHandler{
		following: following,
		publisher: publisher,
	}
}

func (h RejectFollowHandler) RejectFollow(ctx context.Context, cmd RejectFollow) error {
	follow, err := h.following.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	//TODO change to approve and implement
	event, err := follow.Reject(cmd.Rejected)
	if err != nil {
		return err
	}

	err = h.following.Save(ctx, follow)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
