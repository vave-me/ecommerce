package domain

import (
	"context"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

type ImporterRepository interface {
	Load(ctx context.Context, sessionID string) (*Importer, error)
	Save(ctx context.Context, importer *Importer) error
}
