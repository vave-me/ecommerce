package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/streams/internal/domain"
)

type AddSeason struct {
	SeriesID     string
	SeasonNumber int
	Title        string
	Description  string
	ThumbnailURL string
}

type AddSeasonHandler struct {
	series ddd.AggregateStore[*domain.Series]
}

func NewAddSeasonHandler(series ddd.AggregateStore[*domain.Series]) AddSeasonHandler {
	return AddSeasonHandler{
		series: series,
	}
}

func (h AddSeasonHandler) AddSeason(ctx context.Context, cmd AddSeason) error {
	series, err := h.series.Load(ctx, cmd.SeriesID)
	if err != nil {
		return err
	}

	event, err := series.AddSeason(
		cmd.SeasonNumber,
		cmd.Title,
		cmd.Description,
		cmd.ThumbnailURL,
	)
	if err != nil {
		return err
	}
	series.AddEvent(event)

	return h.series.Save(ctx, series)
}