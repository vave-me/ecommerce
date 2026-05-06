package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/media/internal/domain"
)

type (
	UpdateMedia struct {
		ID       string
		ItemID   string
		ItemType domain.ItemType
		UserID   string
		Status   domain.MediaStatus
	}

	UpdateMediaHandler struct {
		medias    domain.MediaRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewUpdateMediaHandler(medias domain.MediaRepository, publisher ddd.EventPublisher[ddd.Event]) UpdateMediaHandler {
	return UpdateMediaHandler{
		medias:    medias,
		publisher: publisher,
	}
}

func (h UpdateMediaHandler) UpdateMedia(ctx context.Context, cmd UpdateMedia) error {

	media, err := h.medias.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := media.UpdateMedia(cmd.ItemID, cmd.ItemType, cmd.UserID, cmd.Status)
	if err != nil {
		return err
	}

	err = h.medias.Save(ctx, media)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
