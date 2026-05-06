package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/streams/internal/domain"
)

type PublishStream struct {
	StreamID string
}

type PublishStreamHandler struct {
	streams ddd.AggregateStore[*domain.Stream]
}

func NewPublishStreamHandler(streams ddd.AggregateStore[*domain.Stream]) PublishStreamHandler {
	return PublishStreamHandler{
		streams: streams,
	}
}

func (h PublishStreamHandler) PublishStream(ctx context.Context, cmd PublishStream) error {
	stream, err := h.streams.Load(ctx, cmd.StreamID)
	if err != nil {
		return err
	}

	event, err := stream.PublishStream()
	if err != nil {
		return err
	}
	stream.AddEvent(event)

	return h.streams.Save(ctx, stream)
}