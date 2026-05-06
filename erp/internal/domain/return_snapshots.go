package domain

import "time"

type ReturnV1 struct {
	ReturnNumber    string
	OriginalOrderID string
	CustomerID      string
	CustomerName    string
	CustomerEmail   string
	Status          ReturnStatus
	Reason          ReturnReason
	Items           []ReturnItem
	RefundDetails   RefundDetails
	WarehouseID     string
	Notes           string
	CreatedAt       time.Time
	ApprovedAt      *time.Time
	ApprovedBy      string
	ProcessedAt     *time.Time
	CompletedAt     *time.Time
	RejectedAt      *time.Time
	RejectedReason  string
	ERPReturnID     string
	ConnectorID     string
	LastSyncedAt    *time.Time
	SyncError       string
}

func (ReturnV1) SnapshotName() string { return "erp.ReturnV1" }