package queries

import (
	"context"
	"middleman/streams/internal/domain"
)

type GetLiveStream struct {
	StreamID string
}

type GetLiveStreamHandler struct {
	liveStreams domain.LiveStreamRepository
}

func NewGetLiveStreamHandler(liveStreams domain.LiveStreamRepository) GetLiveStreamHandler {
	return GetLiveStreamHandler{
		liveStreams: liveStreams,
	}
}

func (h GetLiveStreamHandler) GetLiveStream(ctx context.Context, query GetLiveStream) (*domain.LiveStream, error) {
	return h.liveStreams.Find(ctx, query.StreamID)
}