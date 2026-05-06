package queries

import (
	"context"
	"middleman/streams/internal/domain"
)

type GetSeries struct {
	SeriesID string
}

type GetSeriesHandler struct {
	series domain.SeriesRepository
}

func NewGetSeriesHandler(series domain.SeriesRepository) GetSeriesHandler {
	return GetSeriesHandler{series: series}
}

func (h GetSeriesHandler) GetSeries(ctx context.Context, query GetSeries) (*domain.Series, error) {
	return h.series.Find(ctx, query.SeriesID)
}