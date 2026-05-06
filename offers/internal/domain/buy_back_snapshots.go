package domain

import "time"

type BuyBackV1 struct {
	OfferID             string
	LockedPrice         int64
	RedemptionFee       int64
	LockBuyerID         string
	LockStartAt         time.Time
	LockExpiresAt       time.Time
	BuyBackStatus       BuyBackStatus
	NegotiatedPrice     int64
	NegotiationComments string
}

func (BuyBackV1) SnapshotName() string { return "offers.BuyBackV1" }
