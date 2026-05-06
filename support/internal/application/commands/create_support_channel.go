package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type CreateSupportChannel struct {
	ID          string
	UserID      string
	BusinessID  string
	ChannelType domain.SupportChannelType
	Settings    domain.SupportChannelSettings
}

type CreateSupportChannelHandler struct {
	channels  domain.SupportChannelRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCreateSupportChannelHandler(
	channels domain.SupportChannelRepository,
	publisher ddd.EventPublisher[ddd.Event],
) CreateSupportChannelHandler {
	return CreateSupportChannelHandler{
		channels:  channels,
		publisher: publisher,
	}
}

func (h CreateSupportChannelHandler) CreateSupportChannel(ctx context.Context, cmd CreateSupportChannel) error {
	channel, err := h.channels.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := channel.InitChannel(cmd.UserID, cmd.BusinessID, cmd.ChannelType, cmd.Settings)
	if err != nil {
		return err
	}

	err = h.channels.Save(ctx, channel)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}