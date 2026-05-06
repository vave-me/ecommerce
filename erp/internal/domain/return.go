package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"
	"strconv"
	"time"

	"github.com/stackus/errors"
)

const ReturnAggregate = "erp.Return"

var (
	ErrReturnAlreadyCreated      = errors.Wrap(errors.ErrBadRequest, "return already created")
	ErrReturnNumberIsBlank       = errors.Wrap(errors.ErrBadRequest, "return number cannot be blank")
	ErrOriginalOrderIsBlank      = errors.Wrap(errors.ErrBadRequest, "original order ID cannot be blank")
	ErrReturnNoItems             = errors.Wrap(errors.ErrBadRequest, "return must have at least one item")
	ErrReturnMustBePending       = errors.Wrap(errors.ErrBadRequest, "return must be pending to approve")
	ErrReturnMustBeApproved      = errors.Wrap(errors.ErrBadRequest, "return must be approved before processing")
	ErrReturnMustBeProcessing    = errors.Wrap(errors.ErrBadRequest, "return must be processing to complete")
	ErrOnlyPendingCanBeRejected = errors.Wrap(errors.ErrBadRequest, "only pending returns can be rejected")
	ErrInvalidRestockStatus      = errors.Wrap(errors.ErrBadRequest, "return must be processing or completed to restock items")
	ErrERPReturnIDIsBlank        = errors.Wrap(errors.ErrBadRequest, "ERP return ID cannot be blank")
	ErrSyncErrorIsNil            = errors.Wrap(errors.ErrBadRequest, "error is required to record sync failure")
)

// Return represents a return merchandise authorization aggregate
type Return struct {
	es.Aggregate
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

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Return)(nil)

func NewReturn(id string) *Return {
	return &Return{
		Aggregate: es.NewAggregate(id, ReturnAggregate),
	}
}

// Key implements registry.Registerable
func (Return) Key() string { return ReturnAggregate }

// CreateReturn creates a new return
func (r *Return) CreateReturn(
	returnNumber string,
	originalOrderID string,
	customerID string,
	customerName string,
	customerEmail string,
	reason ReturnReason,
	items []ReturnItem,
	refundMethod RefundMethod,
	refundAmount float64,
	warehouseID string,
	notes string,
	connectorID string,
) (ddd.Event, error) {
	if r.ReturnNumber != "" {
		return nil, ErrReturnAlreadyCreated
	}

	if returnNumber == "" {
		return nil, ErrReturnNumberIsBlank
	}

	if originalOrderID == "" {
		return nil, ErrOriginalOrderIsBlank
	}

	if len(items) == 0 {
		return nil, ErrReturnNoItems
	}

	// Convert items for event
	eventItems := make([]ReturnItemEvent, len(items))
	for i, item := range items {
		eventItems[i] = ReturnItemEvent{
			SKU:             item.SKU,
			ProductName:     item.ProductName,
			Quantity:        item.Quantity,
			Condition:       string(item.Condition),
			RestockingFee:   item.RestockingFee,
			RefundAmount:    item.RefundAmount,
			ExchangeForSKU:  item.ExchangeForSKU,
			SerialNumbers:   item.SerialNumbers,
			InspectionNotes: item.InspectionNotes,
		}
	}

	r.AddEvent(ReturnCreatedEvent, &ReturnCreated{
		ReturnNumber:    returnNumber,
		OriginalOrderID: originalOrderID,
		CustomerID:      customerID,
		CustomerName:    customerName,
		CustomerEmail:   customerEmail,
		Reason:          string(reason),
		Items:           eventItems,
		RefundMethod:    string(refundMethod),
		RefundAmount:    refundAmount,
		WarehouseID:     warehouseID,
		Notes:           notes,
		CreatedAt:       time.Now(),
		ConnectorID:     connectorID,
	})

	return ddd.NewEvent(ReturnCreatedEvent, r), nil
}

// ApproveReturn approves the return
func (r *Return) ApproveReturn(approvedBy string) (ddd.Event, error) {
	if r.Status != ReturnStatusPending {
		return nil, ErrReturnMustBePending
	}

	r.AddEvent(ReturnApprovedEvent, &ReturnApproved{
		ApprovedAt: time.Now(),
		ApprovedBy: approvedBy,
	})

	return ddd.NewEvent(ReturnApprovedEvent, r), nil
}

// ProcessReturn marks the return as being processed
func (r *Return) ProcessReturn(erpReturnID string) (ddd.Event, error) {
	if r.Status != ReturnStatusApproved {
		return nil, ErrReturnMustBeApproved
	}

	r.AddEvent(ReturnProcessedEvent, &ReturnProcessed{
		ProcessedAt: time.Now(),
		ERPReturnID: erpReturnID,
	})

	return ddd.NewEvent(ReturnProcessedEvent, r), nil
}

// CompleteReturn marks the return as completed
func (r *Return) CompleteReturn(refundTransactionID string) (ddd.Event, error) {
	if r.Status != ReturnStatusProcessing {
		return nil, ErrReturnMustBeProcessing
	}

	r.AddEvent(ReturnCompletedEvent, &ReturnCompleted{
		CompletedAt:         time.Now(),
		RefundProcessedAt:   time.Now(),
		RefundTransactionID: refundTransactionID,
	})

	return ddd.NewEvent(ReturnCompletedEvent, r), nil
}

// RejectReturn rejects the return
func (r *Return) RejectReturn(reason string) (ddd.Event, error) {
	if r.Status != ReturnStatusPending {
		return nil, ErrOnlyPendingCanBeRejected
	}

	r.AddEvent(ReturnRejectedEvent, &ReturnRejected{
		RejectedAt: time.Now(),
		Reason:     reason,
	})

	return ddd.NewEvent(ReturnRejectedEvent, r), nil
}

// RestockItems records that items have been restocked
func (r *Return) RestockItems(items []RestockedItem) (ddd.Event, error) {
	if r.Status != ReturnStatusProcessing && r.Status != ReturnStatusCompleted {
		return nil, ErrInvalidRestockStatus
	}

	r.AddEvent(ReturnItemsRestockedEvent, &ReturnItemsRestocked{
		Items:       items,
		RestockedAt: time.Now(),
	})

	return ddd.NewEvent(ReturnItemsRestockedEvent, r), nil
}

// LinkToERP links the return to an ERP return ID
func (r *Return) LinkToERP(erpReturnID string) (ddd.Event, error) {
	if erpReturnID == "" {
		return nil, ErrERPReturnIDIsBlank
	}

	r.AddEvent(ReturnLinkedToERPEvent, &ReturnLinkedToERP{
		ERPReturnID: erpReturnID,
		LinkedAt:    time.Now(),
	})

	return ddd.NewEvent(ReturnLinkedToERPEvent, r), nil
}

// RecordSyncFailure records that ERP sync failed
func (r *Return) RecordSyncFailure(err error) (ddd.Event, error) {
	if err == nil {
		return nil, ErrSyncErrorIsNil
	}

	r.AddEvent(ReturnSyncFailedEvent, &ReturnSyncFailed{
		Error:    err.Error(),
		FailedAt: time.Now(),
	})

	return ddd.NewEvent(ReturnSyncFailedEvent, r), nil
}

// ApplyEvent implements es.EventApplier
func (r *Return) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {
	case *ReturnCreated:
		r.ReturnNumber = e.ReturnNumber
		r.OriginalOrderID = e.OriginalOrderID
		r.CustomerID = e.CustomerID
		r.CustomerName = e.CustomerName
		r.CustomerEmail = e.CustomerEmail
		r.Status = ReturnStatusPending
		r.Reason = ReturnReason(e.Reason)
		r.WarehouseID = e.WarehouseID
		r.Notes = e.Notes
		r.CreatedAt = e.CreatedAt
		r.ConnectorID = e.ConnectorID

		// Convert items
		r.Items = make([]ReturnItem, len(e.Items))
		for i, item := range e.Items {
			r.Items[i] = ReturnItem{
				SKU:             item.SKU,
				ProductName:     item.ProductName,
				Quantity:        item.Quantity,
				Condition:       ItemCondition(item.Condition),
				RestockingFee:   item.RestockingFee,
				RefundAmount:    item.RefundAmount,
				ExchangeForSKU:  item.ExchangeForSKU,
				SerialNumbers:   item.SerialNumbers,
				InspectionNotes: item.InspectionNotes,
			}
		}

		r.RefundDetails = RefundDetails{
			Method: RefundMethod(e.RefundMethod),
			Amount: e.RefundAmount,
		}

	case *ReturnApproved:
		r.Status = ReturnStatusApproved
		r.ApprovedAt = &e.ApprovedAt
		r.ApprovedBy = e.ApprovedBy

	case *ReturnProcessed:
		r.Status = ReturnStatusProcessing
		r.ProcessedAt = &e.ProcessedAt
		r.ERPReturnID = e.ERPReturnID

	case *ReturnCompleted:
		r.Status = ReturnStatusCompleted
		r.CompletedAt = &e.CompletedAt
		r.RefundDetails.ProcessedAt = e.RefundProcessedAt
		r.RefundDetails.TransactionID = e.RefundTransactionID

	case *ReturnRejected:
		r.Status = ReturnStatusRejected
		r.RejectedAt = &e.RejectedAt
		r.RejectedReason = e.Reason

	case *ReturnItemsRestocked:
		// Update item conditions if needed
		for i, item := range r.Items {
			for _, restocked := range e.Items {
				if item.SKU == restocked.SKU {
					note := item.InspectionNotes
					if note != "" {
						note += "; "
					}
					note += "Restocked " + strconv.Itoa(restocked.Quantity) + " units at " + e.RestockedAt.Format(time.RFC3339)
					r.Items[i].InspectionNotes = note
				}
			}
		}

	case *ReturnLinkedToERP:
		r.ERPReturnID = e.ERPReturnID
		r.LastSyncedAt = &e.LinkedAt
		r.SyncError = ""

	case *ReturnSyncFailed:
		r.SyncError = e.Error
		r.LastSyncedAt = &e.FailedAt

	default:
		return errors.ErrInternal.Msgf(
			"%T received the event %s with unexpected payload %T",
			r, event.EventName(), e)
	}
	return nil
}

// ApplySnapshot implements es.Snapshotter
func (r *Return) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ReturnV1:
		r.ReturnNumber = ss.ReturnNumber
		r.OriginalOrderID = ss.OriginalOrderID
		r.CustomerID = ss.CustomerID
		r.CustomerName = ss.CustomerName
		r.CustomerEmail = ss.CustomerEmail
		r.Status = ss.Status
		r.Reason = ss.Reason
		r.Items = ss.Items
		r.RefundDetails = ss.RefundDetails
		r.WarehouseID = ss.WarehouseID
		r.Notes = ss.Notes
		r.CreatedAt = ss.CreatedAt
		r.ApprovedAt = ss.ApprovedAt
		r.ApprovedBy = ss.ApprovedBy
		r.ProcessedAt = ss.ProcessedAt
		r.CompletedAt = ss.CompletedAt
		r.RejectedAt = ss.RejectedAt
		r.RejectedReason = ss.RejectedReason
		r.ERPReturnID = ss.ERPReturnID
		r.ConnectorID = ss.ConnectorID
		r.LastSyncedAt = ss.LastSyncedAt
		r.SyncError = ss.SyncError

	default:
		return errors.ErrInternal.Msgf(
			"%T received the unexpected snapshot %T", r, snapshot)
	}
	return nil
}

// ToSnapshot implements es.Snapshotter
func (r Return) ToSnapshot() es.Snapshot {
	return ReturnV1{
		ReturnNumber:    r.ReturnNumber,
		OriginalOrderID: r.OriginalOrderID,
		CustomerID:      r.CustomerID,
		CustomerName:    r.CustomerName,
		CustomerEmail:   r.CustomerEmail,
		Status:          r.Status,
		Reason:          r.Reason,
		Items:           r.Items,
		RefundDetails:   r.RefundDetails,
		WarehouseID:     r.WarehouseID,
		Notes:           r.Notes,
		CreatedAt:       r.CreatedAt,
		ApprovedAt:      r.ApprovedAt,
		ApprovedBy:      r.ApprovedBy,
		ProcessedAt:     r.ProcessedAt,
		CompletedAt:     r.CompletedAt,
		RejectedAt:      r.RejectedAt,
		RejectedReason:  r.RejectedReason,
		ERPReturnID:     r.ERPReturnID,
		ConnectorID:     r.ConnectorID,
		LastSyncedAt:    r.LastSyncedAt,
		SyncError:       r.SyncError,
	}
}

// ShouldRestock determines if items should be restocked based on their condition
func (r *Return) ShouldRestock(condition ItemCondition) bool {
	switch condition {
	case ItemConditionNew, ItemConditionOpened:
		return true
	case ItemConditionUsed:
		return true // May require inspection
	case ItemConditionDefective, ItemConditionDamaged:
		return false // Don't restock damaged items
	default:
		return false
	}
}
