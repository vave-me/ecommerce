package queries

import (
	"context"
	"middleman/streams/internal/domain"
)

type SearchStreams struct {
	Query   string
	Filters domain.StreamFilters
}

type SearchStreamsHandler struct {
	streams domain.StreamRepository
}

func NewSearchStreamsHandler(streams domain.StreamRepository) SearchStreamsHandler {
	return SearchStreamsHandler{streams: streams}
}

func (h SearchStreamsHandler) SearchStreams(ctx context.Context, query SearchStreams) ([]*domain.Stream, error) {
	return h.streams.Search(ctx, query.Query, query.Filters)
}