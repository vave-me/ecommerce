package queries

import (
	"context"
	"middleman/tickets/internal/domain"
)

type GetMatch struct {
	MatchID string
}

type GetMatchHandler struct {
	matches domain.MatchRepository
}

func NewGetMatchHandler(matches domain.MatchRepository) GetMatchHandler {
	return GetMatchHandler{matches: matches}
}

func (h GetMatchHandler) GetMatch(ctx context.Context, query GetMatch) (*domain.Match, error) {
	return h.matches.Find(ctx, query.MatchID)
}