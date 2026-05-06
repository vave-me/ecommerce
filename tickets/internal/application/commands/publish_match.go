package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/tickets/internal/domain"
)

type PublishMatch struct {
	MatchID string
}

type PublishMatchHandler struct {
	matches ddd.AggregateStore[*domain.Match]
}

func NewPublishMatchHandler(matches ddd.AggregateStore[*domain.Match]) PublishMatchHandler {
	return PublishMatchHandler{
		matches: matches,
	}
}

func (h PublishMatchHandler) PublishMatch(ctx context.Context, cmd PublishMatch) error {
	match, err := h.matches.Load(ctx, cmd.MatchID)
	if err != nil {
		return err
	}

	event, err := match.PublishMatch()
	if err != nil {
		return err
	}
	match.AddEvent(event)

	return h.matches.Save(ctx, match)
}