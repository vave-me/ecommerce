package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/streams/internal/domain"
)

type CreateStream struct {
	StreamID      string
	Title         string
	Description   string
	Synopsis      string
	StreamType    domain.StreamType
	StreamURL     string
	ThumbnailURL  string
	TrailerURL    string
	Duration      int64
	ContentRating domain.ContentRating
	AccessType    domain.AccessType
	Genre         []string
	Language      string
	Country       string
	Studio        string
}

type CreateStreamHandler struct {
	streams ddd.AggregateStore[*domain.Stream]
}

func NewCreateStreamHandler(streams ddd.AggregateStore[*domain.Stream]) CreateStreamHandler {
	return CreateStreamHandler{
		streams: streams,
	}
}

func (h CreateStreamHandler) CreateStream(ctx context.Context, cmd CreateStream) error {
	stream := domain.NewStream(cmd.StreamID)

	event, err := stream.InitStream(
		cmd.Title,
		cmd.Description,
		cmd.Synopsis,
		cmd.StreamType,
		cmd.StreamURL,
		cmd.ThumbnailURL,
		cmd.TrailerURL,
		cmd.Duration,
		cmd.ContentRating,
		cmd.AccessType,
		cmd.Genre,
		cmd.Language,
		cmd.Country,
		cmd.Studio,
	)
	if err != nil {
		return err
	}
	stream.AddEvent(event)

	return h.streams.Save(ctx, stream)
}