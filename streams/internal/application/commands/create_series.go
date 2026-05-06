package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/streams/internal/domain"
)

type CreateSeries struct {
	SeriesID     string
	Title        string
	Description  string
	ThumbnailURL string
	Genre        []string
	Studio       string
}

type CreateSeriesHandler struct {
	series ddd.AggregateStore[*domain.Series]
}

func NewCreateSeriesHandler(series ddd.AggregateStore[*domain.Series]) CreateSeriesHandler {
	return CreateSeriesHandler{
		series: series,
	}
}

func (h CreateSeriesHandler) CreateSeries(ctx context.Context, cmd CreateSeries) error {
	series := domain.NewSeries(cmd.SeriesID)

	event, err := series.InitSeries(
		cmd.Title,
		cmd.Description,
		cmd.ThumbnailURL,
		cmd.Genre,
		cmd.Studio,
	)
	if err != nil {
		return err
	}
	series.AddEvent(event)

	return h.series.Save(ctx, series)
}