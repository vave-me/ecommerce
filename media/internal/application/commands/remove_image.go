package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/media/internal/domain"
)

type RemoveImage struct {
	ID      string
	MediaID string
}

type RemoveImageHandler struct {
	images    domain.ImageRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRemoveImageHandler(images domain.ImageRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveImageHandler {
	return RemoveImageHandler{
		images:    images,
		publisher: publisher,
	}
}

func (h RemoveImageHandler) RemoveImage(ctx context.Context, cmd RemoveImage) error {
	image, err := h.images.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding image")
	}

	event, err := image.Remove(cmd.ID, cmd.MediaID)
	if err != nil {
		return errors.Wrap(err, "remove image")
	}

	return errors.Wrap(h.publisher.Publish(ctx, event), "publishing domain event")
}
