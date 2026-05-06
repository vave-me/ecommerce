package queries

import (
	"context"
	"middleman/streams/internal/domain"
)

type GetUserStreams struct {
	UserID string
}

type GetUserStreamsHandler struct {
	streams domain.StreamRepository
}

func NewGetUserStreamsHandler(streams domain.StreamRepository) GetUserStreamsHandler {
	return GetUserStreamsHandler{streams: streams}
}

func (h GetUserStreamsHandler) GetUserStreams(ctx context.Context, query GetUserStreams) ([]*domain.Stream, error) {
	return h.streams.FindByUserAccess(ctx, query.UserID)
}