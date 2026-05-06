package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type OfferRepository interface {
	// Core Offer operations
	CreateOffer(ctx context.Context, userSellerID, productID string, price int64) (*models.CreateOfferResponse, error)
	ActivateOffer(ctx context.Context, offerID string) (*models.ActivateOfferResponse, error)
	CloseOffer(ctx context.Context, offerID, reason string) (*models.CloseOfferResponse, error)
	AcceptOffer(ctx context.Context, offerID, userCustomerID string) (*models.AcceptOfferResponse, error)
	GetOffer(ctx context.Context, offerID string) (*models.GetOfferResponse, error)
	ListOffers(ctx context.Context, userSellerID, userCustomerID, offerStatus string, page, limit int64) (*models.ListOffersResponse, error)

	// Offer Negotiation operations
	RequestOfferNegotiation(ctx context.Context, offerID string, proposedPrice int64, message string) (*models.NegotiationResponse, error)
	AcceptOfferNegotiation(ctx context.Context, offerID string, finalPrice int64) (*models.NegotiationResponse, error)
	DeclineOfferNegotiation(ctx context.Context, offerID, reason string) (*models.NegotiationResponse, error)

	// BuyNow operations
	CreateBuyNow(ctx context.Context, offerID string, finalPrice int64) (*models.CreateBuyNowResponse, error)
	ConfirmBuyNow(ctx context.Context, buyNowID string) (*models.ConfirmBuyNowResponse, error)
	CancelBuyNow(ctx context.Context, buyNowID string) (*models.CancelBuyNowResponse, error)

	// BuyNow Negotiation operations
	RequestBuyNowNegotiation(ctx context.Context, buyNowID string, newPrice int64, message string) (*models.NegotiationResponse, error)
	AcceptBuyNowNegotiation(ctx context.Context, buyNowID string, finalPrice int64) (*models.NegotiationResponse, error)
	DeclineBuyNowNegotiation(ctx context.Context, buyNowID, reason string) (*models.NegotiationResponse, error)

	// Lease operations
	CreateLease(ctx context.Context, offerID string, monthlyPrice, leaseTermMonths int64, hasBuyout bool, buyoutPrice int64) (*models.CreateLeaseResponse, error)
	StartLease(ctx context.Context, leaseID string) (*models.StartLeaseResponse, error)
	MakeLeasePayment(ctx context.Context, leaseID string, amount int64) (*models.MakeLeasePaymentResponse, error)
	ExecuteLeaseBuyout(ctx context.Context, leaseID string, buyoutAmount int64) (*models.ExecuteLeaseBuyoutResponse, error)
	EndLease(ctx context.Context, leaseID string) (*models.EndLeaseResponse, error)
	CancelLease(ctx context.Context, leaseID string) (*models.CancelLeaseResponse, error)
	DefaultLease(ctx context.Context, leaseID, reason string) (*models.DefaultLeaseResponse, error)

	// Lease Negotiation operations
	RequestLeaseNegotiation(ctx context.Context, leaseID string, proposedMonthlyPrice, proposedTermMonths int64, message string) (*models.NegotiationResponse, error)
	AcceptLeaseNegotiation(ctx context.Context, leaseID string, finalMonthlyPrice, finalTermMonths int64) (*models.NegotiationResponse, error)
	DeclineLeaseNegotiation(ctx context.Context, leaseID, reason string) (*models.NegotiationResponse, error)

	// BuyBack operations
	CreateBuyBack(ctx context.Context, offerID string, lockedPrice, redemptionFee int64, lockDurationDays int32, lockBuyerID string) (*models.CreateBuyBackResponse, error)
	RedeemBuyBack(ctx context.Context, buyBackID string) (*models.RedeemBuyBackResponse, error)
	ExpireBuyBack(ctx context.Context, buyBackID string) (*models.ExpireBuyBackResponse, error)
	CancelBuyBack(ctx context.Context, buyBackID string) (*models.CancelBuyBackResponse, error)

	// BuyBack Negotiation operations
	RequestBuyBackNegotiation(ctx context.Context, buyBackID string, newLockedPrice, newRedemptionFee int64, message string) (*models.NegotiationResponse, error)
	AcceptBuyBackNegotiation(ctx context.Context, buyBackID string, agreedLockedPrice, agreedRedemptionFee int64) (*models.NegotiationResponse, error)
	DeclineBuyBackNegotiation(ctx context.Context, buyBackID, reason string) (*models.NegotiationResponse, error)

	// Reservation operations
	CreateReservation(ctx context.Context, offerID string, lockedPrice, reservationFee int64, lockDurationDays int32, lockBuyerID string) (*models.CreateReservationResponse, error)
	RedeemReservation(ctx context.Context, reservationID string) (*models.RedeemReservationResponse, error)
	ExpireReservation(ctx context.Context, reservationID string) (*models.ExpireReservationResponse, error)
	CancelReservation(ctx context.Context, reservationID string) (*models.CancelReservationResponse, error)

	// Reservation Negotiation operations
	RequestReservationNegotiation(ctx context.Context, reservationID string, newLockedPrice, newReservationFee int64, message string) (*models.NegotiationResponse, error)
	AcceptReservationNegotiation(ctx context.Context, reservationID string, agreedLockedPrice, agreedReservationFee int64) (*models.NegotiationResponse, error)
	DeclineReservationNegotiation(ctx context.Context, reservationID, reason string) (*models.NegotiationResponse, error)

	// Additional query methods for AI tooling
	GetOffersByProduct(ctx context.Context, productID string, limit int64) ([]*models.Offer, error)
	GetOffersByUser(ctx context.Context, userID string, limit int64) ([]*models.Offer, error)
	SearchOffers(ctx context.Context, query string, limit int64) ([]*models.Offer, error)
	GetActiveLeases(ctx context.Context, userID string, limit int64) ([]*models.Lease, error)
	GetActiveBuyBacks(ctx context.Context, userID string, limit int64) ([]*models.BuyBack, error)
	GetActiveReservations(ctx context.Context, userID string, limit int64) ([]*models.Reservation, error)
}
