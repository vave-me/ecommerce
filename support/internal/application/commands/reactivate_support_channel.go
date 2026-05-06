package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type ReactivateSupportChannelHandler struct {
	channels  domain.SupportChannelRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewReactivateSupportChannelHandler(
	channels domain.SupportChannelRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ReactivateSupportChannelHandler {
	return ReactivateSupportChannelHandler{
		channels:  channels,
		publisher: publisher,
	}
}

func (h ReactivateSupportChannelHandler) ReactivateSupportChannel(ctx context.Context, cmd ReactivateSupportChannel) error {
	channel, err := h.channels.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := channel.ReactivateChannel(cmd.ReactivatedBy)
	if err != nil {
		return err
	}
	h.publisher.Publish(ctx, event)

	return h.channels.Save(ctx, channel)
}