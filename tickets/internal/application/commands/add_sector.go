package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/tickets/internal/domain"
)

type AddSector struct {
	MatchID       string
	SectorID      string
	Name          string
	Level         int
	Category      string
	TotalSeats    int64
	BasePrice     int64
	Amenities     []string
	EntranceGates []string
}

type AddSectorHandler struct {
	matches ddd.AggregateStore[*domain.Match]
}

func NewAddSectorHandler(matches ddd.AggregateStore[*domain.Match]) AddSectorHandler {
	return AddSectorHandler{
		matches: matches,
	}
}

func (h AddSectorHandler) AddSector(ctx context.Context, cmd AddSector) error {
	match, err := h.matches.Load(ctx, cmd.MatchID)
	if err != nil {
		return err
	}

	event, err := match.AddSector(
		cmd.SectorID,
		cmd.Name,
		cmd.Level,
		cmd.Category,
		cmd.TotalSeats,
		cmd.BasePrice,
		cmd.Amenities,
		cmd.EntranceGates,
	)
	if err != nil {
		return err
	}
	match.AddEvent(event)

	return h.matches.Save(ctx, match)
}