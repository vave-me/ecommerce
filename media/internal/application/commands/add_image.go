package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/media/internal/domain"
)

type AddImage struct {
	ID           string
	MediaID      string
	DisplayOrder int
	IsMain       bool
	Url          string
	Metadata     string
	FileType     string
	Thumbnail    string
	UserID       string
}

type AddImageHandler struct {
	images    domain.ImageRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddImageHandler(images domain.ImageRepository, publisher ddd.EventPublisher[ddd.Event]) AddImageHandler {
	return AddImageHandler{
		images:    images,
		publisher: publisher,
	}
}

func (h AddImageHandler) AddImage(ctx context.Context, cmd AddImage) error {
	image, err := h.images.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding image")
	}

	event, err := image.InitImage(cmd.ID, cmd.MediaID, cmd.DisplayOrder, cmd.IsMain, cmd.Url, cmd.Metadata, cmd.FileType, cmd.Thumbnail, cmd.UserID)
	if err != nil {
		return errors.Wrap(err, "initializing image")
	}

	err = h.images.Save(ctx, image)
	if err != nil {
		return errors.Wrap(err, "error adding image")
	}

	return errors.Wrap(h.publisher.Publish(ctx, event), "publishing domain event")
}
