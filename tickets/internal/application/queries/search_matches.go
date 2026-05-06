package queries

import (
	"context"
	"middleman/tickets/internal/domain"
)

type SearchMatches struct {
	Criteria domain.MatchSearchCriteria
}

type SearchMatchesHandler struct {
	matches domain.MatchRepository
}

func NewSearchMatchesHandler(matches domain.MatchRepository) SearchMatchesHandler {
	return SearchMatchesHandler{matches: matches}
}

func (h SearchMatchesHandler) SearchMatches(ctx context.Context, query SearchMatches) ([]*domain.Match, error) {
	return h.matches.SearchMatches(ctx, query.Criteria)
}