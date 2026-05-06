package offerspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	// Channels
	OfferAggregateChannel       = "middleman.offers.events.Offer"
	LeaseAggregateChannel       = "middleman.offers.events.Lease"
	BuyBackAggregateChannel     = "middleman.offers.events.BuyBack"
	BuyNowAggregateChannel      = "middleman.offers.events.BuyNow"
	ReservationAggregateChannel = "middleman.offers.events.Reservation"

	// Offer events
	OfferCreatedEvent   = "offersapi.OfferCreated"
	OfferActivatedEvent = "offersapi.OfferActivated"
	OfferClosedEvent    = "offersapi.OfferClosed"
	OfferAcceptedEvent  = "offersapi.OfferAccepted"

	// BuyNow events
	BuyNowCreatedEvent   = "offersapi.BuyNowCreated"
	BuyNowConfirmedEvent = "offersapi.BuyNowConfirmed"

	// Lease events
	LeaseAddedEvent     = "offersapi.LeaseAdded"
	LeaseCreatedEvent   = "offersapi.LeaseCreated"
	LeaseDefaultedEvent = "offersapi.LeaseDefaulted"
	LeaseEndedEvent     = "offersapi.LeaseEnded"

	// BuyBack events
	BuyBackCreatedEvent  = "offersapi.BuyBackCreated"
	BuyBackCanceledEvent = "offersapi.BuyBackCanceled"
	BuyBackExpiredEvent  = "offersapi.BuyBackExpired"
	BuyBackRedeemedEvent = "offersapi.BuyBackRedeemed"

	// Reservation events
	ReservationCreatedEvent  = "offersapi.ReservationCreated"
	ReservationCanceledEvent = "offersapi.ReservationCanceled"
	ReservationExpiredEvent  = "offersapi.ReservationExpired"
	ReservationRedeemedEvent = "offersapi.ReservationRedeemed"
)

func (*OfferCreated) Key() string   { return OfferCreatedEvent }
func (*OfferActivated) Key() string { return OfferActivatedEvent }
func (*OfferClosed) Key() string    { return OfferClosedEvent }
func (*OfferAccepted) Key() string  { return OfferAcceptedEvent }

func (*BuyNowCreated) Key() string   { return BuyNowCreatedEvent }
func (*BuyNowConfirmed) Key() string { return BuyNowConfirmedEvent }

// func (*LeaseAdded) Key() string     { return LeaseAddedEvent }
func (*LeaseCreated) Key() string   { return LeaseCreatedEvent }
func (*LeaseDefaulted) Key() string { return LeaseDefaultedEvent }
func (*LeaseEnded) Key() string     { return LeaseEndedEvent }

func (*BuyBackCreated) Key() string  { return BuyBackCreatedEvent }
func (*BuyBackCanceled) Key() string { return BuyBackCanceledEvent }
func (*BuyBackExpired) Key() string  { return BuyBackExpiredEvent }
func (*BuyBackRedeemed) Key() string { return BuyBackRedeemedEvent }

func (*ReservationCreated) Key() string  { return ReservationCreatedEvent }
func (*ReservationCanceled) Key() string { return ReservationCanceledEvent }
func (*ReservationExpired) Key() string  { return ReservationExpiredEvent }
func (*ReservationRedeemed) Key() string { return ReservationRedeemedEvent }

// -----------------------------------------------------------------------------
// Registration logic
// -----------------------------------------------------------------------------

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Offer events
	if err := serde.Register(&OfferCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&OfferActivated{}); err != nil {
		return err
	}
	if err := serde.Register(&OfferClosed{}); err != nil {
		return err
	}
	if err := serde.Register(&OfferAccepted{}); err != nil {
		return err
	}

	// BuyNow events
	if err := serde.Register(&BuyNowCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&BuyNowConfirmed{}); err != nil {
		return err
	}

	//// Lease events
	//if err := serde.Register(&LeaseAdded{}); err != nil {
	//	return err
	//}
	if err := serde.Register(&LeaseCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&LeaseDefaulted{}); err != nil {
		return err
	}
	if err := serde.Register(&LeaseEnded{}); err != nil {
		return err
	}

	// BuyBack events
	if err := serde.Register(&BuyBackCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&BuyBackCanceled{}); err != nil {
		return err
	}
	if err := serde.Register(&BuyBackExpired{}); err != nil {
		return err
	}
	if err := serde.Register(&BuyBackRedeemed{}); err != nil {
		return err
	}

	// Reservation events
	if err := serde.Register(&ReservationCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&ReservationCanceled{}); err != nil {
		return err
	}
	if err := serde.Register(&ReservationExpired{}); err != nil {
		return err
	}
	if err := serde.Register(&ReservationRedeemed{}); err != nil {
		return err
	}

	return nil
}
