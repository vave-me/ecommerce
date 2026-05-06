package domain

import "time"

type LeaseV1 struct {
	OfferID         string
	MonthlyPrice    int64
	LeaseTermMonths int64
	HasBuyout       bool
	BuyoutPrice     int64
	LeaseStartDate  time.Time
	LeaseEndDate    time.Time
	LeaseStatus     LeaseStatus
}

func (LeaseV1) SnapshotName() string { return "offers.LeaseV1" }
