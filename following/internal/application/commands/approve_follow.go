package commands

import (
	"context"
	"middleman/following/internal/domain"
	"middleman/internal/ddd"
)

type ApproveFollow struct {
	ID       string
	Approval bool
}

type ApproveFollowHandler struct {
	following domain.FollowRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewApproveFollowHandler(following domain.FollowRepository, publisher ddd.EventPublisher[ddd.Event]) ApproveFollowHandler {
	return ApproveFollowHandler{
		following: following,
		publisher: publisher,
	}
}

func (h ApproveFollowHandler) ApproveFollow(ctx context.Context, cmd ApproveFollow) error {
	follow, err := h.following.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	//TODO change to approve and implement
	event, err := follow.Approve(cmd.Approval)
	if err != nil {
		return err
	}

	err = h.following.Save(ctx, follow)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
