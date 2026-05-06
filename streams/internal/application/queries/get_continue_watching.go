package queries

import (
	"context"
	"middleman/streams/internal/domain"
)

type GetContinueWatching struct {
	UserID string
}

type GetContinueWatchingHandler struct {
	streams domain.StreamRepository
}

func NewGetContinueWatchingHandler(streams domain.StreamRepository) GetContinueWatchingHandler {
	return GetContinueWatchingHandler{streams: streams}
}

func (h GetContinueWatchingHandler) GetContinueWatching(ctx context.Context, query GetContinueWatching) ([]*domain.Stream, error) {
	return h.streams.GetContinueWatching(ctx, query.UserID)
}