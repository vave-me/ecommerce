package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/tickets/internal/domain"
)

type InitializeSectorSeats struct {
	MatchID  string
	SectorID string
	Rows     []RowDefinition
}

type RowDefinition struct {
	RowID       string
	RowNumber   string
	StartSeat   int
	EndSeat     int
}

type InitializeSectorSeatsHandler struct {
	matches ddd.AggregateStore[*domain.Match]
}

func NewInitializeSectorSeatsHandler(matches ddd.AggregateStore[*domain.Match]) InitializeSectorSeatsHandler {
	return InitializeSectorSeatsHandler{
		matches: matches,
	}
}

func (h InitializeSectorSeatsHandler) InitializeSectorSeats(ctx context.Context, cmd InitializeSectorSeats) error {
	match, err := h.matches.Load(ctx, cmd.MatchID)
	if err != nil {
		return err
	}

	// Add rows and initialize seats
	for _, row := range cmd.Rows {
		// First add the row
		event, err := match.AddRowToSector(
			cmd.SectorID,
			row.RowID,
			row.RowNumber,
			row.EndSeat - row.StartSeat + 1,
		)
		if err != nil {
			return err
		}
		match.AddEvent(event)

		// Then initialize seats in the row
		event, err = match.InitializeSeatsInRow(
			cmd.SectorID,
			row.RowID,
			row.StartSeat,
			row.EndSeat,
		)
		if err != nil {
			return err
		}
		match.AddEvent(event)
	}

	return h.matches.Save(ctx, match)
}