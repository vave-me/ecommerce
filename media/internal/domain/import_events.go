package domain

const (
	BulkImportStartedEvent   = "media.BulkImportStarted"
	ImportBatchAddedEvent    = "media.ImportBatchAdded"
	ImportItemProcessedEvent = "media.ImportItemProcessed"
	ImportItemFailedEvent    = "media.ImportItemFailed"
	BulkImportCompletedEvent = "media.BulkImportCompleted"
	BulkImportCancelledEvent = "media.BulkImportCancelled"
)

type BulkImportStarted struct {
	SessionID          string
	ExternalSystemID   string
	ExternalSystemType string
	EstimatedCount     int32
	UserID             string
}

func (BulkImportStarted) Key() string { return BulkImportStartedEvent }

type ImportBatchAdded struct {
	SessionID   string
	BatchSize   int32
	BatchNumber int32
}

func (ImportBatchAdded) Key() string { return ImportBatchAddedEvent }

type ImportItemProcessed struct {
	SessionID  string
	ExternalID string
	MediaID    string
	ImageID    string
	Status     string
}

func (ImportItemProcessed) Key() string { return ImportItemProcessedEvent }

type ImportItemFailed struct {
	SessionID    string
	ExternalID   string
	ErrorCode    string
	ErrorMessage string
	RetryCount   int32
}

func (ImportItemFailed) Key() string { return ImportItemFailedEvent }

type BulkImportCompleted struct {
	SessionID         string
	TotalProcessed    int32
	SuccessfulImports int32
	FailedImports     int32
	DurationMs        int64
}

func (BulkImportCompleted) Key() string { return BulkImportCompletedEvent }

type BulkImportCancelled struct {
	SessionID               string
	Reason                  string
	ProcessedAtCancellation int32
}

func (BulkImportCancelled) Key() string { return BulkImportCancelledEvent }
