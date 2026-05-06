package models

import "time"

const RecurringPaymentSetupEvent = "payments.RecurringPaymentSetup"
const RecurringPaymentChargedEvent = "payments.RecurringPaymentCharged"

type RecurringPaymentSetup struct {
	PlanID         string
	UserCustomerID string
	Amount         int64
	Frequency      string
	StartDate      time.Time
}
type RecurringPaymentCharged struct {
	PlanID    string
	PaymentID string
	Amount    int64
}
