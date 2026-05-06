package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/streams/internal/domain"
)

type GrantUserAccess struct {
	StreamID   string
	UserID     string
	AccessType domain.AccessType
	Duration   int64 // in hours
}

type GrantUserAccessHandler struct {
	streams ddd.AggregateStore[*domain.Stream]
}

func NewGrantUserAccessHandler(streams ddd.AggregateStore[*domain.Stream]) GrantUserAccessHandler {
	return GrantUserAccessHandler{
		streams: streams,
	}
}

func (h GrantUserAccessHandler) GrantUserAccess(ctx context.Context, cmd GrantUserAccess) error {
	stream, err := h.streams.Load(ctx, cmd.StreamID)
	if err != nil {
		return err
	}

	event, err := stream.GrantUserAccess(cmd.UserID, cmd.AccessType, cmd.Duration)
	if err != nil {
		return err
	}
	stream.AddEvent(event)

	return h.streams.Save(ctx, stream)
}