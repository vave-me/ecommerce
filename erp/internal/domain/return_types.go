package domain

import "time"

// ReturnStatus represents the status of a return
type ReturnStatus string

const (
	ReturnStatusPending    ReturnStatus = "pending"
	ReturnStatusApproved   ReturnStatus = "approved"
	ReturnStatusProcessing ReturnStatus = "processing"
	ReturnStatusCompleted  ReturnStatus = "completed"
	ReturnStatusRejected   ReturnStatus = "rejected"
	ReturnStatusCanceled   ReturnStatus = "canceled"
)

// ReturnReason represents the reason for return
type ReturnReason string

const (
	ReturnReasonDefective    ReturnReason = "defective"
	ReturnReasonNotAsDesc    ReturnReason = "not_as_described"
	ReturnReasonWrongItem    ReturnReason = "wrong_item"
	ReturnReasonDamaged      ReturnReason = "damaged"
	ReturnReasonNotNeeded    ReturnReason = "not_needed"
	ReturnReasonLateDelivery ReturnReason = "late_delivery"
	ReturnReasonOther        ReturnReason = "other"
)

// RefundMethod represents how the refund will be processed
type RefundMethod string

const (
	RefundMethodOriginal    RefundMethod = "original_payment"
	RefundMethodStoreCredit RefundMethod = "store_credit"
	RefundMethodExchange    RefundMethod = "exchange"
	RefundMethodPartial     RefundMethod = "partial_refund"
)

// ItemCondition represents the condition of returned item
type ItemCondition string

const (
	ItemConditionNew       ItemCondition = "new"
	ItemConditionOpened    ItemCondition = "opened"
	ItemConditionUsed      ItemCondition = "used"
	ItemConditionDefective ItemCondition = "defective"
	ItemConditionDamaged   ItemCondition = "damaged"
)

// ReturnItem represents an item being returned
type ReturnItem struct {
	SKU             string
	ProductName     string
	Quantity        int
	Condition       ItemCondition
	RestockingFee   float64
	RefundAmount    float64
	ExchangeForSKU  string // If exchanging for different item
	SerialNumbers   []string
	InspectionNotes string
}

// RefundDetails contains refund information
type RefundDetails struct {
	Method        RefundMethod
	Amount        float64
	ProcessedAt   time.Time
	TransactionID string
}