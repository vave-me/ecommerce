package domain

import "time"

type ReservationV1 struct {
	OfferID             string
	LockedPrice         int64
	ReservationFee      int64
	LockBuyerID         string
	LockStartAt         time.Time
	LockExpiresAt       time.Time
	ReservationStatus   ReservationStatus
	IsFree              bool
	FreeReservationUsed bool
	NegotiationComments string
	NegotiatedPrice     int64
}

func (ReservationV1) SnapshotName() string { return "offers.ReservationV1" }
