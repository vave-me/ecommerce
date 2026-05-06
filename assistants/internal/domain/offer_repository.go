package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type OfferRepository interface {
	// Core Offer operations
	CreateNewSellerOffer(ctx context.Context, userSellerID, productID string, price int64) (*models.CreateOfferResponse, error)
	ActivateExistingOffer(ctx context.Context, offerID string) (*models.ActivateOfferResponse, error)
	CloseOfferWithReason(ctx context.Context, offerID, reason string) (*models.CloseOfferResponse, error)
	AcceptOfferByCustomer(ctx context.Context, offerID, userCustomerID string) (*models.AcceptOfferResponse, error)
	GetOfferDetailsByID(ctx context.Context, offerID string) (*models.GetOfferResponse, error)
	ListOffersWithFilters(ctx context.Context, userSellerID, userCustomerID, offerStatus string, page, limit int64) (*models.ListOffersResponse, error)

	// Offer Negotiation operations
	RequestOfferPriceNegotiation(ctx context.Context, offerID string, proposedPrice int64, message string) (*models.NegotiationResponse, error)
	AcceptNegotiatedOfferPrice(ctx context.Context, offerID string, finalPrice int64) (*models.NegotiationResponse, error)
	DeclineOfferNegotiationRequest(ctx context.Context, offerID, reason string) (*models.NegotiationResponse, error)

	// BuyNow operations
	CreateBuyNowTransaction(ctx context.Context, offerID string, finalPrice int64) (*models.CreateBuyNowResponse, error)
	ConfirmBuyNowPurchase(ctx context.Context, buyNowID string) (*models.ConfirmBuyNowResponse, error)
	CancelBuyNowTransaction(ctx context.Context, buyNowID string) (*models.CancelBuyNowResponse, error)

	// BuyNow Negotiation operations
	RequestBuyNowPriceAdjustment(ctx context.Context, buyNowID string, newPrice int64, message string) (*models.NegotiationResponse, error)
	AcceptBuyNowPriceChange(ctx context.Context, buyNowID string, finalPrice int64) (*models.NegotiationResponse, error)
	DeclineBuyNowPriceNegotiation(ctx context.Context, buyNowID, reason string) (*models.NegotiationResponse, error)

	// Lease operations
	CreateLeaseAgreementForOffer(ctx context.Context, offerID string, monthlyPrice, leaseTermMonths int64, hasBuyout bool, buyoutPrice int64) (*models.CreateLeaseResponse, error)
	StartActiveLeaseAgreement(ctx context.Context, leaseID string) (*models.StartLeaseResponse, error)
	RecordMonthlyLeasePayment(ctx context.Context, leaseID string, amount int64) (*models.MakeLeasePaymentResponse, error)
	ExecuteLeaseBuyoutOption(ctx context.Context, leaseID string, buyoutAmount int64) (*models.ExecuteLeaseBuyoutResponse, error)
	CompleteLeaseAgreement(ctx context.Context, leaseID string) (*models.EndLeaseResponse, error)
	CancelActiveLeaseAgreement(ctx context.Context, leaseID string) (*models.CancelLeaseResponse, error)
	MarkLeaseAsDefaulted(ctx context.Context, leaseID, reason string) (*models.DefaultLeaseResponse, error)

	// Lease Negotiation operations
	RequestLeaseTermsNegotiation(ctx context.Context, leaseID string, proposedMonthlyPrice, proposedTermMonths int64, message string) (*models.NegotiationResponse, error)
	AcceptNegotiatedLeaseTerms(ctx context.Context, leaseID string, finalMonthlyPrice, finalTermMonths int64) (*models.NegotiationResponse, error)
	DeclineLeaseTermsNegotiation(ctx context.Context, leaseID, reason string) (*models.NegotiationResponse, error)

	// BuyBack operations
	CreateBuyBackAgreement(ctx context.Context, offerID string, lockedPrice, redemptionFee int64, lockDurationDays int32, lockBuyerID string) (*models.CreateBuyBackResponse, error)
	RedeemBuyBackOption(ctx context.Context, buyBackID string) (*models.RedeemBuyBackResponse, error)
	ExpireBuyBackAgreement(ctx context.Context, buyBackID string) (*models.ExpireBuyBackResponse, error)
	CancelBuyBackAgreement(ctx context.Context, buyBackID string) (*models.CancelBuyBackResponse, error)

	// BuyBack Negotiation operations
	RequestBuyBackTermsNegotiation(ctx context.Context, buyBackID string, newLockedPrice, newRedemptionFee int64, message string) (*models.NegotiationResponse, error)
	AcceptNegotiatedBuyBackTerms(ctx context.Context, buyBackID string, agreedLockedPrice, agreedRedemptionFee int64) (*models.NegotiationResponse, error)
	DeclineBuyBackTermsNegotiation(ctx context.Context, buyBackID, reason string) (*models.NegotiationResponse, error)

	// Reservation operations
	CreateOfferReservation(ctx context.Context, offerID string, lockedPrice, reservationFee int64, lockDurationDays int32, lockBuyerID string) (*models.CreateReservationResponse, error)
	RedeemOfferReservation(ctx context.Context, reservationID string) (*models.RedeemReservationResponse, error)
	ExpireOfferReservation(ctx context.Context, reservationID string) (*models.ExpireReservationResponse, error)
	CancelOfferReservation(ctx context.Context, reservationID string) (*models.CancelReservationResponse, error)

	// Reservation Negotiation operations
	RequestReservationTermsNegotiation(ctx context.Context, reservationID string, newLockedPrice, newReservationFee int64, message string) (*models.NegotiationResponse, error)
	AcceptNegotiatedReservationTerms(ctx context.Context, reservationID string, agreedLockedPrice, agreedReservationFee int64) (*models.NegotiationResponse, error)
	DeclineReservationTermsNegotiation(ctx context.Context, reservationID, reason string) (*models.NegotiationResponse, error)

	// Additional query methods for AI tooling
	GetAllOffersForProduct(ctx context.Context, productID string, limit int64) ([]*models.Offer, error)
	GetAllOffersFromUser(ctx context.Context, userID string, limit int64) ([]*models.Offer, error)
	SearchOffersByKeyword(ctx context.Context, query string, limit int64) ([]*models.Offer, error)
	GetActiveLeaseAgreements(ctx context.Context, userID string, limit int64) ([]*models.Lease, error)
	GetActiveBuyBackAgreements(ctx context.Context, userID string, limit int64) ([]*models.BuyBack, error)
	GetActiveOfferReservations(ctx context.Context, userID string, limit int64) ([]*models.Reservation, error)
}
