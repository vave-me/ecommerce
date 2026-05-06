package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/media/internal/domain"
)

type (
	HasMedia struct {
		ID     string
		ItemID string
		UserID string
	}

	HasMediaHandler struct {
		medias    domain.MediaRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewHasMediaHandler(medias domain.MediaRepository, publisher ddd.EventPublisher[ddd.Event]) HasMediaHandler {
	return HasMediaHandler{
		medias:    medias,
		publisher: publisher,
	}
}

func (h HasMediaHandler) HasMedia(ctx context.Context, cmd HasMedia) error {

	media, err := h.medias.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := media.Delete(cmd.ID, cmd.UserID)
	if err != nil {
		return err
	}

	err = h.medias.Save(ctx, media)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
