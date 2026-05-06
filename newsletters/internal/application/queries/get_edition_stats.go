package queries

import (
	"context"
	"middleman/newsletters/internal/domain"
)

type GetEditionStats struct {
	EditionID string
}

type EditionStats struct {
	Recipients   int
	Delivered    int
	Opened       int
	Clicked      int
	Bounced      int
	Complaints   int
	OpenRate     float32
	ClickRate    float32
}

type GetEditionStatsHandler struct {
	editionCatalog domain.EditionCatalogRepository
}

func NewGetEditionStatsHandler(editionCatalog domain.EditionCatalogRepository) GetEditionStatsHandler {
	return GetEditionStatsHandler{
		editionCatalog: editionCatalog,
	}
}

func (h GetEditionStatsHandler) GetEditionStats(ctx context.Context, query GetEditionStats) (*EditionStats, error) {
	edition, err := h.editionCatalog.Find(ctx, query.EditionID)
	if err != nil {
		return nil, err
	}

	// TODO: Get actual stats from send logs

	return &EditionStats{
		Recipients: edition.RecipientCount,
		Delivered:  edition.RecipientCount, // TODO: Get from send logs
		Opened:     0,                      // TODO: Get from send logs
		Clicked:    0,                      // TODO: Get from send logs
		Bounced:    0,                      // TODO: Get from send logs
		Complaints: 0,                      // TODO: Get from send logs
		OpenRate:   0,                      // TODO: Calculate
		ClickRate:  0,                      // TODO: Calculate
	}, nil
}