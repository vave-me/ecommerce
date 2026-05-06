package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/media/internal/domain"
)

type (
	CreateMedia struct {
		ID       string
		ItemID   string
		ItemType domain.ItemType
		UserID   string
		Status   domain.MediaStatus
	}

	CreateMediaHandler struct {
		medias    domain.MediaRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewCreateMediaHandler(medias domain.MediaRepository, publisher ddd.EventPublisher[ddd.Event]) CreateMediaHandler {
	return CreateMediaHandler{
		medias:    medias,
		publisher: publisher,
	}
}

func (h CreateMediaHandler) CreateMedia(ctx context.Context, cmd CreateMedia) error {

	media, err := h.medias.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := media.InitMedia(cmd.ItemID, cmd.ItemType, cmd.UserID, cmd.Status)
	if err != nil {
		return err
	}

	err = h.medias.Save(ctx, media)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
