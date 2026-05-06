package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/tickets/internal/domain"
	"time"
)

type CreateMatch struct {
	MatchID        string
	HomeTeam       domain.Team
	AwayTeam       domain.Team
	Competition    domain.Competition
	MatchDate      time.Time
	Stadium        domain.Stadium
	SalesStartDate time.Time
	SalesEndDate   time.Time
}

type CreateMatchHandler struct {
	matches ddd.AggregateStore[*domain.Match]
}

func NewCreateMatchHandler(matches ddd.AggregateStore[*domain.Match]) CreateMatchHandler {
	return CreateMatchHandler{
		matches: matches,
	}
}

func (h CreateMatchHandler) CreateMatch(ctx context.Context, cmd CreateMatch) error {
	match := domain.NewMatch(cmd.MatchID)

	event, err := match.InitMatch(
		cmd.HomeTeam,
		cmd.AwayTeam,
		cmd.Competition,
		cmd.MatchDate,
		cmd.Stadium,
		cmd.SalesStartDate,
		cmd.SalesEndDate,
	)
	if err != nil {
		return err
	}
	match.AddEvent(event)

	return h.matches.Save(ctx, match)
}