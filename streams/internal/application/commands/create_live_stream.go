package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/streams/internal/domain"
	"time"
)

type CreateLiveStream struct {
	ID                 string
	Title              string
	Description        string
	EventType          string
	HomeTeam           string
	AwayTeam           string
	Competition        string
	Season             string
	MatchDay           int
	Stadium            string
	ScheduledStartTime time.Time
	ScheduledEndTime   time.Time
}

type CreateLiveStreamHandler struct {
	liveStreams domain.LiveStreamRepository
	publisher   ddd.EventPublisher
}

func NewCreateLiveStreamHandler(liveStreams domain.LiveStreamRepository, publisher ddd.EventPublisher) CreateLiveStreamHandler {
	return CreateLiveStreamHandler{
		liveStreams: liveStreams,
		publisher:   publisher,
	}
}

func (h CreateLiveStreamHandler) CreateLiveStream(ctx context.Context, cmd CreateLiveStream) error {
	stream := domain.NewLiveStream(cmd.ID)

	if err := stream.InitLiveStream(
		cmd.Title,
		cmd.Description,
		cmd.EventType,
		cmd.HomeTeam,
		cmd.AwayTeam,
		cmd.Competition,
		cmd.Season,
		cmd.MatchDay,
		cmd.Stadium,
		cmd.ScheduledStartTime,
		cmd.ScheduledEndTime,
	); err != nil {
		return err
	}

	if err := h.liveStreams.Save(ctx, stream); err != nil {
		return err
	}

	// Publish domain events
	if err := h.publisher.Publish(ctx, stream.Events()...); err != nil {
		return err
	}

	return nil
}