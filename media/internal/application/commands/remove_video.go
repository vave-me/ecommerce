package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/media/internal/domain"
)

type RemoveVideo struct {
	ID      string
	MediaID string
}

type RemoveVideoHandler struct {
	videos    domain.VideoRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRemoveVideoHandler(videos domain.VideoRepository, publisher ddd.EventPublisher[ddd.Event]) RemoveVideoHandler {
	return RemoveVideoHandler{
		videos:    videos,
		publisher: publisher,
	}
}

func (h RemoveVideoHandler) RemoveVideo(ctx context.Context, cmd RemoveVideo) error {
	video, err := h.videos.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading video")
	}

	event, err := video.Remove(cmd.ID, cmd.MediaID)

	return h.publisher.Publish(ctx, event)
}
