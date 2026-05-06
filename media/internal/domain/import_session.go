package domain

import (
	"context"
	"errors"
	"time"
)

type ImportStatus string

const (
	ImportStatusPending    ImportStatus = "pending"
	ImportStatusProcessing ImportStatus = "processing"
	ImportStatusCompleted  ImportStatus = "completed"
	ImportStatusFailed     ImportStatus = "failed"
	ImportStatusCancelled  ImportStatus = "cancelled"
)

type ImportItemStatus string

const (
	ImportItemStatusPending    ImportItemStatus = "pending"
	ImportItemStatusFetching   ImportItemStatus = "fetching"
	ImportItemStatusProcessing ImportItemStatus = "processing"
	ImportItemStatusCompleted  ImportItemStatus = "completed"
	ImportItemStatusFailed     ImportItemStatus = "failed"
)

var (
	ErrImportSessionNotFound  = errors.New("import session not found")
	ErrImportSessionNotActive = errors.New("import session is not active")
	ErrImportAlreadyFinished  = errors.New("import session already finished")
)

type ImportSession struct {
	ID                 string
	ExternalSystemID   string
	ExternalSystemType string
	TotalImages        int32
	ProcessedImages    int32
	FailedImages       int32
	Status             ImportStatus
	StartedAt          time.Time
	CompletedAt        time.Time
	Metadata           map[string]string
	BatchCount         int32
}

type ImportItem struct {
	ID           string
	SessionID    string
	ExternalID   string
	SKU          string
	ProductID    string // Product ID resolved from SKU
	ImageURL     string
	Status       ImportItemStatus
	ErrorMessage string
	RetryCount   int32
	MediaID      string
	ImageID      string
	ProcessedAt  time.Time
	Metadata     map[string]string
	DisplayOrder int32
}

type ImportSessionRepository interface {
	Create(ctx context.Context, session *ImportSession) error
	Get(ctx context.Context, sessionID string) (*ImportSession, error)
	Update(ctx context.Context, session *ImportSession) error
}

type ImportItemRepository interface {
	CreateBatch(ctx context.Context, items []*ImportItem) error
	GetBySession(ctx context.Context, sessionID string, status ImportItemStatus) ([]*ImportItem, error)
	Update(ctx context.Context, item *ImportItem) error
	GetPendingItems(ctx context.Context, limit int) ([]*ImportItem, error)
}