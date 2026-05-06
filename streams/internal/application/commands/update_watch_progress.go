package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/streams/internal/domain"
)

type UpdateWatchProgress struct {
	StreamID  string
	UserID    string
	Progress  int64 // in seconds
	Completed bool
}

type UpdateWatchProgressHandler struct {
	streams ddd.AggregateStore[*domain.Stream]
}

func NewUpdateWatchProgressHandler(streams ddd.AggregateStore[*domain.Stream]) UpdateWatchProgressHandler {
	return UpdateWatchProgressHandler{
		streams: streams,
	}
}

func (h UpdateWatchProgressHandler) UpdateWatchProgress(ctx context.Context, cmd UpdateWatchProgress) error {
	stream, err := h.streams.Load(ctx, cmd.StreamID)
	if err != nil {
		return err
	}

	event, err := stream.UpdateWatchProgress(cmd.UserID, cmd.Progress, cmd.Completed)
	if err != nil {
		return err
	}
	stream.AddEvent(event)

	return h.streams.Save(ctx, stream)
}