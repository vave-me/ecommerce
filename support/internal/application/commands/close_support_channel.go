package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type CloseSupportChannelHandler struct {
	channels  domain.SupportChannelRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCloseSupportChannelHandler(
	channels domain.SupportChannelRepository,
	publisher ddd.EventPublisher[ddd.Event],
) CloseSupportChannelHandler {
	return CloseSupportChannelHandler{
		channels:  channels,
		publisher: publisher,
	}
}

func (h CloseSupportChannelHandler) CloseSupportChannel(ctx context.Context, cmd CloseSupportChannel) error {
	channel, err := h.channels.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := channel.CloseChannel(cmd.ClosedBy, cmd.Reason)
	if err != nil {
		return err
	}
	h.publisher.Publish(ctx, event)

	return h.channels.Save(ctx, channel)
}