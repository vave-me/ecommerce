package domain

type OfferType string

const (
	OfferTypeDirect OfferType = "DIRECT"
	OfferTypeLock   OfferType = "LOCK"
	OfferTypeLease  OfferType = "LEASE"
)
