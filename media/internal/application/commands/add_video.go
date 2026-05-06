package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/media/internal/domain"
)

type AddVideo struct {
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

type AddVideoHandler struct {
	videos    domain.VideoRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddVideoHandler(videos domain.VideoRepository, publisher ddd.EventPublisher[ddd.Event]) AddVideoHandler {
	return AddVideoHandler{
		videos:    videos,
		publisher: publisher,
	}
}

func (h AddVideoHandler) AddVideo(ctx context.Context, cmd AddVideo) error {
	video, err := h.videos.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding video")
	}

	event, err := video.InitVideo(cmd.ID, cmd.MediaID, cmd.DisplayOrder, cmd.IsMain, cmd.Url, cmd.Metadata, cmd.FileType, cmd.UserID)
	if err != nil {
		return errors.Wrap(err, "initializing video")
	}

	err = h.videos.Save(ctx, video)
	if err != nil {
		return errors.Wrap(err, "error adding video")
	}

	return errors.Wrap(h.publisher.Publish(ctx, event), "publishing domain event")
}
