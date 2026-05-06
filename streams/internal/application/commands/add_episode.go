package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/streams/internal/domain"
	"time"
)

type AddEpisode struct {
	SeriesID      string
	SeasonNumber  int
	EpisodeNumber int
	StreamID      string
	Title         string
	Duration      int64
	AirDate       time.Time
}

type AddEpisodeHandler struct {
	series ddd.AggregateStore[*domain.Series]
}

func NewAddEpisodeHandler(series ddd.AggregateStore[*domain.Series]) AddEpisodeHandler {
	return AddEpisodeHandler{
		series: series,
	}
}

func (h AddEpisodeHandler) AddEpisode(ctx context.Context, cmd AddEpisode) error {
	series, err := h.series.Load(ctx, cmd.SeriesID)
	if err != nil {
		return err
	}

	event, err := series.AddEpisodeToSeason(
		cmd.SeasonNumber,
		cmd.EpisodeNumber,
		cmd.StreamID,
		cmd.Title,
		cmd.Duration,
		cmd.AirDate,
	)
	if err != nil {
		return err
	}
	series.AddEvent(event)

	return h.series.Save(ctx, series)
}