package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/media/internal/domain"
)

type (
	RemoveMedia struct {
		ID     string
		ItemID string
		UserID string
	}

	RemoveMediaHandler struct {
		medias    domain.MediaRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewRemoveMediaHandler(medias domain.MediaRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveMediaHandler {
	return RemoveMediaHandler{
		medias:    medias,
		publisher: publisher,
	}
}

func (h RemoveMediaHandler) RemoveMedia(ctx context.Context, cmd RemoveMedia) error {

	media, err := h.medias.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := media.Delete(cmd.ID, cmd.UserID)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
