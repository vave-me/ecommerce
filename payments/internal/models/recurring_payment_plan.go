package models

import "time"

// Recurring Payment Plan
const (
	RecurringPlanStatusActive   = "ACTIVE"
	RecurringPlanStatusCanceled = "CANCELED"
)

type RecurringPaymentPlan struct {
	PlanID          string
	UserCustomerID  string
	Amount          int64
	Frequency       string
	StartDate       time.Time
	LastChargedAt   time.Time
	NextDueDate     time.Time
	Status          string
	PaymentMethodID string
}
