package queries

import (
	"context"
	"middleman/streams/internal/domain"
)

type GetStream struct {
	StreamID string
	UserID   string // Optional - to check access
}

type GetStreamHandler struct {
	streams domain.StreamRepository
}

func NewGetStreamHandler(streams domain.StreamRepository) GetStreamHandler {
	return GetStreamHandler{streams: streams}
}

func (h GetStreamHandler) GetStream(ctx context.Context, query GetStream) (*domain.Stream, error) {
	stream, err := h.streams.Find(ctx, query.StreamID)
	if err != nil {
		return nil, err
	}

	// If UserID is provided, check if user has access
	if query.UserID != "" {
		if access, exists := stream.UserAccess[query.UserID]; exists {
			// Check if access is still valid
			if access.ExpiresAt.Before(stream.CreatedAt) {
				// Access expired, remove from response
				delete(stream.UserAccess, query.UserID)
			}
		}
	}

	return stream, nil
}