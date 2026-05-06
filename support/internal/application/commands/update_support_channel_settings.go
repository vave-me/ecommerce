package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/support/internal/domain"
)

type UpdateSupportChannelSettingsHandler struct {
	channels  domain.SupportChannelRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewUpdateSupportChannelSettingsHandler(
	channels domain.SupportChannelRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UpdateSupportChannelSettingsHandler {
	return UpdateSupportChannelSettingsHandler{
		channels:  channels,
		publisher: publisher,
	}
}

func (h UpdateSupportChannelSettingsHandler) UpdateSupportChannelSettings(ctx context.Context, cmd UpdateSupportChannelSettings) error {
	channel, err := h.channels.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := channel.UpdateSettings(cmd.Settings)
	if err != nil {
		return err
	}
	h.publisher.Publish(ctx, event)

	return h.channels.Save(ctx, channel)
}