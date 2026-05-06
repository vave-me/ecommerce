package domain

import "time"

// Return Event Names
const (
	ReturnCreatedEvent        = "erp.ReturnCreated"
	ReturnApprovedEvent       = "erp.ReturnApproved"
	ReturnProcessedEvent      = "erp.ReturnProcessed"
	ReturnCompletedEvent      = "erp.ReturnCompleted"
	ReturnRejectedEvent       = "erp.ReturnRejected"
	ReturnItemsRestockedEvent = "erp.ReturnItemsRestocked"
	ReturnLinkedToERPEvent    = "erp.ReturnLinkedToERP"
	ReturnSyncFailedEvent     = "erp.ReturnSyncFailed"
)

// ReturnItemEvent represents a return item in events
type ReturnItemEvent struct {
	SKU             string
	ProductName     string
	Quantity        int
	Condition       string
	RestockingFee   float64
	RefundAmount    float64
	ExchangeForSKU  string
	SerialNumbers   []string
	InspectionNotes string
}

// RestockedItem represents an item that was restocked
type RestockedItem struct {
	SKU        string
	Quantity   int
	LocationID string
}

// ReturnCreated event - fired when a return is created
type ReturnCreated struct {
	ReturnNumber    string
	OriginalOrderID string
	CustomerID      string
	CustomerName    string
	CustomerEmail   string
	Reason          string
	Items           []ReturnItemEvent
	RefundMethod    string
	RefundAmount    float64
	WarehouseID     string
	Notes           string
	CreatedAt       time.Time
	ConnectorID     string
}

// Key implements event registry
func (ReturnCreated) Key() string { return ReturnCreatedEvent }

// ReturnApproved event - fired when a return is approved
type ReturnApproved struct {
	ApprovedAt time.Time
	ApprovedBy string
}

// Key implements event registry
func (ReturnApproved) Key() string { return ReturnApprovedEvent }

// ReturnProcessed event - fired when return processing starts
type ReturnProcessed struct {
	ProcessedAt time.Time
	ERPReturnID string
}

// Key implements event registry
func (ReturnProcessed) Key() string { return ReturnProcessedEvent }

// ReturnCompleted event - fired when a return is completed
type ReturnCompleted struct {
	CompletedAt         time.Time
	RefundProcessedAt   time.Time
	RefundTransactionID string
}

// Key implements event registry
func (ReturnCompleted) Key() string { return ReturnCompletedEvent }

// ReturnRejected event - fired when a return is rejected
type ReturnRejected struct {
	RejectedAt time.Time
	Reason     string
}

// Key implements event registry
func (ReturnRejected) Key() string { return ReturnRejectedEvent }

// ReturnItemsRestocked event - fired when return items are restocked
type ReturnItemsRestocked struct {
	Items       []RestockedItem
	RestockedAt time.Time
}

// Key implements event registry
func (ReturnItemsRestocked) Key() string { return ReturnItemsRestockedEvent }

// ReturnLinkedToERP event - fired when return is linked to ERP
type ReturnLinkedToERP struct {
	ERPReturnID string
	LinkedAt    time.Time
}

// Key implements event registry
func (ReturnLinkedToERP) Key() string { return ReturnLinkedToERPEvent }

// ReturnSyncFailed event - fired when ERP sync fails
type ReturnSyncFailed struct {
	Error    string
	FailedAt time.Time
}

// Key implements event registry
func (ReturnSyncFailed) Key() string { return ReturnSyncFailedEvent }