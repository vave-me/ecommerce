package grpc

import (
	"context"
	"fmt"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/offers/offerspb"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// OfferRepository calls the remote offers service (gRPC).
type OfferRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.OfferRepository = (*OfferRepository)(nil)

// NewOfferRepositoryWithAuth creates a new OfferRepository with JWT authentication support
func NewOfferRepository(endpoint string, authInstance *auth.Auth) OfferRepository {
	return OfferRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// Core Offer operations

// CreateNewSellerOffer creates a new offer
func (r OfferRepository) CreateNewSellerOffer(ctx context.Context, userSellerID, productID string, price int64) (*models.CreateOfferResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.CreateOffer(ctx, &offerspb.CreateOfferRequest{
		UserSellerId: userSellerID,
		ProductId:    productID,
		Price:        price,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateOffer RPC failed: %w", err)
	}

	var createdAt time.Time
	if resp.GetCreatedAt() != nil {
		createdAt = resp.GetCreatedAt().AsTime()
	}

	return &models.CreateOfferResponse{
		ID:        resp.GetId(),
		CreatedAt: createdAt,
	}, nil
}

// ActivateExistingOffer activates an offer
func (r OfferRepository) ActivateExistingOffer(ctx context.Context, offerID string) (*models.ActivateOfferResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.ActivateOffer(ctx, &offerspb.ActivateOfferRequest{
		OfferId: offerID,
	})
	if err != nil {
		return nil, fmt.Errorf("ActivateOffer RPC failed: %w", err)
	}

	return &models.ActivateOfferResponse{
		OfferID:     resp.GetOfferId(),
		OfferStatus: resp.GetOfferStatus(),
	}, nil
}

// CloseOfferWithReason closes an offer
func (r OfferRepository) CloseOfferWithReason(ctx context.Context, offerID, reason string) (*models.CloseOfferResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.CloseOffer(ctx, &offerspb.CloseOfferRequest{
		OfferId: offerID,
		Reason:  reason,
	})
	if err != nil {
		return nil, fmt.Errorf("CloseOffer RPC failed: %w", err)
	}

	return &models.CloseOfferResponse{
		OfferID:     resp.GetOfferId(),
		OfferStatus: resp.GetOfferStatus(),
	}, nil
}

// AcceptOfferByCustomer accepts an offer
func (r OfferRepository) AcceptOfferByCustomer(ctx context.Context, offerID, userCustomerID string) (*models.AcceptOfferResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.AcceptOffer(ctx, &offerspb.AcceptOfferRequest{
		OfferId:        offerID,
		UserCustomerId: userCustomerID,
	})
	if err != nil {
		return nil, fmt.Errorf("AcceptOffer RPC failed: %w", err)
	}

	return &models.AcceptOfferResponse{
		OfferID:     resp.GetOfferId(),
		OfferStatus: resp.GetOfferStatus(),
	}, nil
}

// GetOfferDetailsByID retrieves an offer by ID
func (r OfferRepository) GetOfferDetailsByID(ctx context.Context, offerID string) (*models.GetOfferResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.GetOffer(ctx, &offerspb.GetOfferRequest{
		OfferId: offerID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetOffer RPC failed: %w", err)
	}

	offer := r.convertOfferFromPb(resp.GetOffer())

	return &models.GetOfferResponse{
		Offer: *offer,
	}, nil
}

// ListOffersWithFilters lists offers with filters
func (r OfferRepository) ListOffersWithFilters(ctx context.Context, userSellerID, userCustomerID, offerStatus string, page, limit int64) (*models.ListOffersResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.ListOffers(ctx, &offerspb.ListOffersRequest{
		UserSellerId:   userSellerID,
		UserCustomerId: userCustomerID,
		OfferStatus:    offerStatus,
		Page:           page,
		Limit:          limit,
	})
	if err != nil {
		return nil, fmt.Errorf("ListOffers RPC failed: %w", err)
	}

	offers := make([]models.Offer, len(resp.GetOffers()))
	for i, pbOffer := range resp.GetOffers() {
		offer := r.convertOfferFromPb(pbOffer)
		offers[i] = *offer
	}

	return &models.ListOffersResponse{
		Offers: offers,
		Total:  resp.GetTotal(),
		Page:   resp.GetPage(),
		Limit:  resp.GetLimit(),
	}, nil
}

// convertOfferFromPb converts protobuf Offer to domain Offer
func (r OfferRepository) convertOfferFromPb(pbOffer *offerspb.Offer) *models.Offer {
	if pbOffer == nil {
		return nil
	}

	return &models.Offer{
		ID:             pbOffer.GetId(),
		UserSellerID:   pbOffer.GetUserSellerId(),
		UserCustomerID: pbOffer.GetUserCustomerId(),
		ProductID:      pbOffer.GetProductId(),
		Price:          pbOffer.GetPrice(),
		OfferStatus:    pbOffer.GetOfferStatus(),
		CreatedAt:      time.Now(), // Would be set from protobuf timestamp if available
		UpdatedAt:      time.Now(),
	}
}

// dial establishes a gRPC connection to the offers service
// dial sets up a gRPC connection with the microservice endpoint
func (r OfferRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r OfferRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}

// Offer Negotiation operations - These methods are placeholders as the protobuf service doesn't define them

// RequestOfferPriceNegotiation requests negotiation for an offer
func (r OfferRepository) RequestOfferPriceNegotiation(ctx context.Context, offerID string, proposedPrice int64, message string) (*models.NegotiationResponse, error) {
	// This method is not implemented in the protobuf service
	return &models.NegotiationResponse{
		ID:                offerID,
		NegotiationStatus: models.NegotiationStatusRequested,
	}, nil
}

// AcceptNegotiatedOfferPrice accepts an offer negotiation
func (r OfferRepository) AcceptNegotiatedOfferPrice(ctx context.Context, offerID string, finalPrice int64) (*models.NegotiationResponse, error) {
	// This method is not implemented in the protobuf service
	return &models.NegotiationResponse{
		ID:                offerID,
		NegotiationStatus: models.NegotiationStatusAccepted,
	}, nil
}

// DeclineOfferNegotiationRequest declines an offer negotiation
func (r OfferRepository) DeclineOfferNegotiationRequest(ctx context.Context, offerID, reason string) (*models.NegotiationResponse, error) {
	// This method is not implemented in the protobuf service
	return &models.NegotiationResponse{
		ID:                offerID,
		NegotiationStatus: models.NegotiationStatusDeclined,
	}, nil
}

// BuyNow operations

// CreateBuyNowTransaction creates a buy-now transaction
func (r OfferRepository) CreateBuyNowTransaction(ctx context.Context, offerID string, finalPrice int64) (*models.CreateBuyNowResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.CreateBuyNow(ctx, &offerspb.CreateBuyNowRequest{
		OfferId:    offerID,
		FinalPrice: finalPrice,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateBuyNow RPC failed: %w", err)
	}

	var createdAt time.Time
	if resp.GetCreatedAt() != nil {
		createdAt = resp.GetCreatedAt().AsTime()
	}

	return &models.CreateBuyNowResponse{
		BuyNowID:  resp.GetBuyNowId(),
		CreatedAt: createdAt,
	}, nil
}

// ConfirmBuyNowPurchase confirms a buy-now transaction
func (r OfferRepository) ConfirmBuyNowPurchase(ctx context.Context, buyNowID string) (*models.ConfirmBuyNowResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.ConfirmBuyNow(ctx, &offerspb.ConfirmBuyNowRequest{
		BuyNowId: buyNowID,
	})
	if err != nil {
		return nil, fmt.Errorf("ConfirmBuyNow RPC failed: %w", err)
	}

	var confirmedAt time.Time
	if resp.GetConfirmedAt() != nil {
		confirmedAt = resp.GetConfirmedAt().AsTime()
	}

	return &models.ConfirmBuyNowResponse{
		BuyNowID:    resp.GetBuyNowId(),
		Status:      resp.GetStatus(),
		ConfirmedAt: confirmedAt,
	}, nil
}

// CancelBuyNowTransaction cancels a buy-now transaction
func (r OfferRepository) CancelBuyNowTransaction(ctx context.Context, buyNowID string) (*models.CancelBuyNowResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.CancelBuyNow(ctx, &offerspb.CancelBuyNowRequest{
		BuyNowId: buyNowID,
	})
	if err != nil {
		return nil, fmt.Errorf("CancelBuyNow RPC failed: %w", err)
	}

	var canceledAt time.Time
	if resp.GetCanceledAt() != nil {
		canceledAt = resp.GetCanceledAt().AsTime()
	}

	return &models.CancelBuyNowResponse{
		BuyNowID:     resp.GetBuyNowId(),
		BuyNowStatus: resp.GetBuyNowStatus(),
		CanceledAt:   canceledAt,
	}, nil
}

// BuyNow Negotiation operations

// RequestBuyNowPriceAdjustment requests negotiation for a buy-now transaction
func (r OfferRepository) RequestBuyNowPriceAdjustment(ctx context.Context, buyNowID string, newPrice int64, message string) (*models.NegotiationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.RequestBuyNowNegotiation(ctx, &offerspb.RequestBuyNowNegotiationRequest{
		BuyNowId: buyNowID,
		NewPrice: newPrice,
		Message:  message,
	})
	if err != nil {
		return nil, fmt.Errorf("RequestBuyNowNegotiation RPC failed: %w", err)
	}

	return &models.NegotiationResponse{
		ID:                resp.GetBuyNowId(),
		NegotiationStatus: resp.GetNegotiationStatus(),
	}, nil
}

// AcceptBuyNowPriceChange accepts a buy-now negotiation
func (r OfferRepository) AcceptBuyNowPriceChange(ctx context.Context, buyNowID string, finalPrice int64) (*models.NegotiationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.AcceptBuyNowNegotiation(ctx, &offerspb.AcceptBuyNowNegotiationRequest{
		BuyNowId:   buyNowID,
		FinalPrice: finalPrice,
	})
	if err != nil {
		return nil, fmt.Errorf("AcceptBuyNowNegotiation RPC failed: %w", err)
	}

	return &models.NegotiationResponse{
		ID:                resp.GetBuyNowId(),
		NegotiationStatus: resp.GetNegotiationStatus(),
	}, nil
}

// DeclineBuyNowPriceNegotiation declines a buy-now negotiation
func (r OfferRepository) DeclineBuyNowPriceNegotiation(ctx context.Context, buyNowID, reason string) (*models.NegotiationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.DeclineBuyNowNegotiation(ctx, &offerspb.DeclineBuyNowNegotiationRequest{
		BuyNowId: buyNowID,
		Reason:   reason,
	})
	if err != nil {
		return nil, fmt.Errorf("DeclineBuyNowNegotiation RPC failed: %w", err)
	}

	return &models.NegotiationResponse{
		ID:                resp.GetBuyNowId(),
		NegotiationStatus: resp.GetNegotiationStatus(),
	}, nil
}

// Lease operations

// CreateLeaseAgreementForOffer creates a new lease
func (r OfferRepository) CreateLeaseAgreementForOffer(ctx context.Context, offerID string, monthlyPrice, leaseTermMonths int64, hasBuyout bool, buyoutPrice int64) (*models.CreateLeaseResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.CreateLease(ctx, &offerspb.CreateLeaseRequest{
		OfferId:         offerID,
		MonthlyPrice:    monthlyPrice,
		LeaseTermMonths: leaseTermMonths,
		HasBuyout:       hasBuyout,
		BuyoutPrice:     buyoutPrice,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateLease RPC failed: %w", err)
	}

	var createdAt time.Time
	if resp.GetCreatedAt() != nil {
		createdAt = resp.GetCreatedAt().AsTime()
	}

	return &models.CreateLeaseResponse{
		LeaseID:     resp.GetLeaseId(),
		LeaseStatus: resp.GetLeaseStatus(),
		CreatedAt:   createdAt,
	}, nil
}

// StartActiveLeaseAgreement starts a lease
func (r OfferRepository) StartActiveLeaseAgreement(ctx context.Context, leaseID string) (*models.StartLeaseResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.StartLease(ctx, &offerspb.StartLeaseRequest{
		LeaseId: leaseID,
	})
	if err != nil {
		return nil, fmt.Errorf("StartLease RPC failed: %w", err)
	}

	var startedAt, endDate time.Time
	if resp.GetStartedAt() != nil {
		startedAt = resp.GetStartedAt().AsTime()
	}
	if resp.GetEndDate() != nil {
		endDate = resp.GetEndDate().AsTime()
	}

	return &models.StartLeaseResponse{
		LeaseID:     resp.GetLeaseId(),
		LeaseStatus: resp.GetLeaseStatus(),
		StartedAt:   startedAt,
		EndDate:     endDate,
	}, nil
}

// RecordMonthlyLeasePayment makes a lease payment
func (r OfferRepository) RecordMonthlyLeasePayment(ctx context.Context, leaseID string, amount int64) (*models.MakeLeasePaymentResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.MakeLeasePayment(ctx, &offerspb.MakeLeasePaymentRequest{
		LeaseId: leaseID,
		Amount:  amount,
	})
	if err != nil {
		return nil, fmt.Errorf("MakeLeasePayment RPC failed: %w", err)
	}

	var paymentDate time.Time
	if resp.GetPaymentDate() != nil {
		paymentDate = resp.GetPaymentDate().AsTime()
	}

	return &models.MakeLeasePaymentResponse{
		LeaseID:     resp.GetLeaseId(),
		PaymentDate: paymentDate,
	}, nil
}

// ExecuteLeaseBuyoutOption executes a lease buyout
func (r OfferRepository) ExecuteLeaseBuyoutOption(ctx context.Context, leaseID string, buyoutAmount int64) (*models.ExecuteLeaseBuyoutResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.ExecuteLeaseBuyout(ctx, &offerspb.ExecuteLeaseBuyoutRequest{
		LeaseId:      leaseID,
		BuyoutAmount: buyoutAmount,
	})
	if err != nil {
		return nil, fmt.Errorf("ExecuteLeaseBuyout RPC failed: %w", err)
	}

	var executedAt time.Time
	if resp.GetExecutedAt() != nil {
		executedAt = resp.GetExecutedAt().AsTime()
	}

	return &models.ExecuteLeaseBuyoutResponse{
		LeaseID:     resp.GetLeaseId(),
		LeaseStatus: resp.GetLeaseStatus(),
		ExecutedAt:  executedAt,
	}, nil
}

// CompleteLeaseAgreement ends a lease
func (r OfferRepository) CompleteLeaseAgreement(ctx context.Context, leaseID string) (*models.EndLeaseResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.EndLease(ctx, &offerspb.EndLeaseRequest{
		LeaseId: leaseID,
	})
	if err != nil {
		return nil, fmt.Errorf("EndLease RPC failed: %w", err)
	}

	var endedAt time.Time
	if resp.GetEndedAt() != nil {
		endedAt = resp.GetEndedAt().AsTime()
	}

	return &models.EndLeaseResponse{
		LeaseID:     resp.GetLeaseId(),
		LeaseStatus: resp.GetLeaseStatus(),
		EndedAt:     endedAt,
	}, nil
}

// CancelActiveLeaseAgreement cancels a lease
func (r OfferRepository) CancelActiveLeaseAgreement(ctx context.Context, leaseID string) (*models.CancelLeaseResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.CancelLease(ctx, &offerspb.CancelLeaseRequest{
		LeaseId: leaseID,
	})
	if err != nil {
		return nil, fmt.Errorf("CancelLease RPC failed: %w", err)
	}

	var canceledAt time.Time
	if resp.GetCanceledAt() != nil {
		canceledAt = resp.GetCanceledAt().AsTime()
	}

	return &models.CancelLeaseResponse{
		LeaseID:     resp.GetLeaseId(),
		LeaseStatus: resp.GetLeaseStatus(),
		CanceledAt:  canceledAt,
	}, nil
}

// MarkLeaseAsDefaulted defaults a lease
func (r OfferRepository) MarkLeaseAsDefaulted(ctx context.Context, leaseID, reason string) (*models.DefaultLeaseResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.DefaultLease(ctx, &offerspb.DefaultLeaseRequest{
		LeaseId: leaseID,
		Reason:  reason,
	})
	if err != nil {
		return nil, fmt.Errorf("DefaultLease RPC failed: %w", err)
	}

	var defaultedAt time.Time
	if resp.GetDefaultedAt() != nil {
		defaultedAt = resp.GetDefaultedAt().AsTime()
	}

	return &models.DefaultLeaseResponse{
		LeaseID:     resp.GetLeaseId(),
		LeaseStatus: resp.GetLeaseStatus(),
		DefaultedAt: defaultedAt,
	}, nil
}

// Lease Negotiation operations

// RequestLeaseTermsNegotiation requests negotiation for a lease
func (r OfferRepository) RequestLeaseTermsNegotiation(ctx context.Context, leaseID string, proposedMonthlyPrice, proposedTermMonths int64, message string) (*models.NegotiationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.RequestLeaseNegotiation(ctx, &offerspb.RequestLeaseNegotiationRequest{
		LeaseId:              leaseID,
		ProposedMonthlyPrice: proposedMonthlyPrice,
		ProposedTermMonths:   proposedTermMonths,
		Message:              message,
	})
	if err != nil {
		return nil, fmt.Errorf("RequestLeaseNegotiation RPC failed: %w", err)
	}

	return &models.NegotiationResponse{
		ID:                resp.GetLeaseId(),
		NegotiationStatus: resp.GetNegotiationStatus(),
	}, nil
}

// AcceptNegotiatedLeaseTerms accepts a lease negotiation
func (r OfferRepository) AcceptNegotiatedLeaseTerms(ctx context.Context, leaseID string, finalMonthlyPrice, finalTermMonths int64) (*models.NegotiationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.AcceptLeaseNegotiation(ctx, &offerspb.AcceptLeaseNegotiationRequest{
		LeaseId:           leaseID,
		FinalMonthlyPrice: finalMonthlyPrice,
		FinalTermMonths:   finalTermMonths,
	})
	if err != nil {
		return nil, fmt.Errorf("AcceptLeaseNegotiation RPC failed: %w", err)
	}

	return &models.NegotiationResponse{
		ID:                resp.GetLeaseId(),
		NegotiationStatus: resp.GetNegotiationStatus(),
	}, nil
}

// DeclineLeaseTermsNegotiation declines a lease negotiation
func (r OfferRepository) DeclineLeaseTermsNegotiation(ctx context.Context, leaseID, reason string) (*models.NegotiationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.DeclineLeaseNegotiation(ctx, &offerspb.DeclineLeaseNegotiationRequest{
		LeaseId: leaseID,
		Reason:  reason,
	})
	if err != nil {
		return nil, fmt.Errorf("DeclineLeaseNegotiation RPC failed: %w", err)
	}

	return &models.NegotiationResponse{
		ID:                resp.GetLeaseId(),
		NegotiationStatus: resp.GetNegotiationStatus(),
	}, nil
}

// BuyBack operations

// CreateBuyBackAgreement creates a new buy-back agreement
func (r OfferRepository) CreateBuyBackAgreement(ctx context.Context, offerID string, lockedPrice, redemptionFee int64, lockDurationDays int32, lockBuyerID string) (*models.CreateBuyBackResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.CreateBuyBack(ctx, &offerspb.CreateBuyBackRequest{
		OfferId:          offerID,
		LockedPrice:      lockedPrice,
		RedemptionFee:    redemptionFee,
		LockDurationDays: lockDurationDays,
		LockBuyerId:      lockBuyerID,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateBuyBack RPC failed: %w", err)
	}

	var createdAt time.Time
	if resp.GetCreatedAt() != nil {
		createdAt = resp.GetCreatedAt().AsTime()
	}

	return &models.CreateBuyBackResponse{
		BuyBackID:     resp.GetBuyBackId(),
		BuyBackStatus: resp.GetBuyBackStatus(),
		CreatedAt:     createdAt,
	}, nil
}

// RedeemBuyBackOption redeems a buy-back agreement
func (r OfferRepository) RedeemBuyBackOption(ctx context.Context, buyBackID string) (*models.RedeemBuyBackResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.RedeemBuyBack(ctx, &offerspb.RedeemBuyBackRequest{
		BuyBackId: buyBackID,
	})
	if err != nil {
		return nil, fmt.Errorf("RedeemBuyBack RPC failed: %w", err)
	}

	var redeemedAt time.Time
	if resp.GetRedeemedAt() != nil {
		redeemedAt = resp.GetRedeemedAt().AsTime()
	}

	return &models.RedeemBuyBackResponse{
		BuyBackID:     resp.GetBuyBackId(),
		BuyBackStatus: resp.GetBuyBackStatus(),
		RedeemedAt:    redeemedAt,
	}, nil
}

// ExpireBuyBackAgreement expires a buy-back agreement
func (r OfferRepository) ExpireBuyBackAgreement(ctx context.Context, buyBackID string) (*models.ExpireBuyBackResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.ExpireBuyBack(ctx, &offerspb.ExpireBuyBackRequest{
		BuyBackId: buyBackID,
	})
	if err != nil {
		return nil, fmt.Errorf("ExpireBuyBack RPC failed: %w", err)
	}

	var expiredAt time.Time
	if resp.GetExpiredAt() != nil {
		expiredAt = resp.GetExpiredAt().AsTime()
	}

	return &models.ExpireBuyBackResponse{
		BuyBackID:     resp.GetBuyBackId(),
		BuyBackStatus: resp.GetBuyBackStatus(),
		ExpiredAt:     expiredAt,
	}, nil
}

// CancelBuyBackAgreement cancels a buy-back agreement
func (r OfferRepository) CancelBuyBackAgreement(ctx context.Context, buyBackID string) (*models.CancelBuyBackResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.CancelBuyBack(ctx, &offerspb.CancelBuyBackRequest{
		BuyBackId: buyBackID,
	})
	if err != nil {
		return nil, fmt.Errorf("CancelBuyBack RPC failed: %w", err)
	}

	var canceledAt time.Time
	if resp.GetCanceledAt() != nil {
		canceledAt = resp.GetCanceledAt().AsTime()
	}

	return &models.CancelBuyBackResponse{
		BuyBackID:     resp.GetBuyBackId(),
		BuyBackStatus: resp.GetBuyBackStatus(),
		CanceledAt:    canceledAt,
	}, nil
}

// BuyBack Negotiation operations

// RequestBuyBackTermsNegotiation requests negotiation for a buy-back
func (r OfferRepository) RequestBuyBackTermsNegotiation(ctx context.Context, buyBackID string, newLockedPrice, newRedemptionFee int64, message string) (*models.NegotiationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.RequestBuyBackNegotiation(ctx, &offerspb.RequestBuyBackNegotiationRequest{
		BuyBackId:        buyBackID,
		NewLockedPrice:   newLockedPrice,
		NewRedemptionFee: newRedemptionFee,
		Message:          message,
	})
	if err != nil {
		return nil, fmt.Errorf("RequestBuyBackNegotiation RPC failed: %w", err)
	}

	return &models.NegotiationResponse{
		ID:                resp.GetBuyBackId(),
		NegotiationStatus: resp.GetNegotiationStatus(),
	}, nil
}

// AcceptNegotiatedBuyBackTerms accepts a buy-back negotiation
func (r OfferRepository) AcceptNegotiatedBuyBackTerms(ctx context.Context, buyBackID string, agreedLockedPrice, agreedRedemptionFee int64) (*models.NegotiationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.AcceptBuyBackNegotiation(ctx, &offerspb.AcceptBuyBackNegotiationRequest{
		BuyBackId:           buyBackID,
		AgreedLockedPrice:   agreedLockedPrice,
		AgreedRedemptionFee: agreedRedemptionFee,
	})
	if err != nil {
		return nil, fmt.Errorf("AcceptBuyBackNegotiation RPC failed: %w", err)
	}

	return &models.NegotiationResponse{
		ID:                resp.GetBuyBackId(),
		NegotiationStatus: resp.GetNegotiationStatus(),
	}, nil
}

// DeclineBuyBackTermsNegotiation declines a buy-back negotiation
func (r OfferRepository) DeclineBuyBackTermsNegotiation(ctx context.Context, buyBackID, reason string) (*models.NegotiationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.DeclineBuyBackNegotiation(ctx, &offerspb.DeclineBuyBackNegotiationRequest{
		BuyBackId: buyBackID,
		Reason:    reason,
	})
	if err != nil {
		return nil, fmt.Errorf("DeclineBuyBackNegotiation RPC failed: %w", err)
	}

	return &models.NegotiationResponse{
		ID:                resp.GetBuyBackId(),
		NegotiationStatus: resp.GetNegotiationStatus(),
	}, nil
}

// Reservation operations

// CreateOfferReservation creates a new reservation
func (r OfferRepository) CreateOfferReservation(ctx context.Context, offerID string, lockedPrice, reservationFee int64, lockDurationDays int32, lockBuyerID string) (*models.CreateReservationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.CreateReservation(ctx, &offerspb.CreateReservationRequest{
		OfferId:          offerID,
		LockedPrice:      lockedPrice,
		ReservationFee:   reservationFee,
		LockDurationDays: lockDurationDays,
		LockBuyerId:      lockBuyerID,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateReservation RPC failed: %w", err)
	}

	var createdAt time.Time
	if resp.GetCreatedAt() != nil {
		createdAt = resp.GetCreatedAt().AsTime()
	}

	return &models.CreateReservationResponse{
		ReservationID:     resp.GetReservationId(),
		ReservationStatus: resp.GetReservationStatus(),
		CreatedAt:         createdAt,
		IsFree:            resp.GetIsFree(),
	}, nil
}

// RedeemOfferReservation redeems a reservation
func (r OfferRepository) RedeemOfferReservation(ctx context.Context, reservationID string) (*models.RedeemReservationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.RedeemReservation(ctx, &offerspb.RedeemReservationRequest{
		ReservationId: reservationID,
	})
	if err != nil {
		return nil, fmt.Errorf("RedeemReservation RPC failed: %w", err)
	}

	var redeemedAt time.Time
	if resp.GetRedeemedAt() != nil {
		redeemedAt = resp.GetRedeemedAt().AsTime()
	}

	return &models.RedeemReservationResponse{
		ReservationID:     resp.GetReservationId(),
		ReservationStatus: resp.GetReservationStatus(),
		RedeemedAt:        redeemedAt,
	}, nil
}

// ExpireOfferReservation expires a reservation
func (r OfferRepository) ExpireOfferReservation(ctx context.Context, reservationID string) (*models.ExpireReservationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.ExpireReservation(ctx, &offerspb.ExpireReservationRequest{
		ReservationId: reservationID,
	})
	if err != nil {
		return nil, fmt.Errorf("ExpireReservation RPC failed: %w", err)
	}

	var expiredAt time.Time
	if resp.GetExpiredAt() != nil {
		expiredAt = resp.GetExpiredAt().AsTime()
	}

	return &models.ExpireReservationResponse{
		ReservationID:     resp.GetReservationId(),
		ReservationStatus: resp.GetReservationStatus(),
		ExpiredAt:         expiredAt,
	}, nil
}

// CancelOfferReservation cancels a reservation
func (r OfferRepository) CancelOfferReservation(ctx context.Context, reservationID string) (*models.CancelReservationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.CancelReservation(ctx, &offerspb.CancelReservationRequest{
		ReservationId: reservationID,
	})
	if err != nil {
		return nil, fmt.Errorf("CancelReservation RPC failed: %w", err)
	}

	var canceledAt time.Time
	if resp.GetCanceledAt() != nil {
		canceledAt = resp.GetCanceledAt().AsTime()
	}

	return &models.CancelReservationResponse{
		ReservationID:     resp.GetReservationId(),
		ReservationStatus: resp.GetReservationStatus(),
		CanceledAt:        canceledAt,
	}, nil
}

// Reservation Negotiation operations

// RequestReservationTermsNegotiation requests negotiation for a reservation
func (r OfferRepository) RequestReservationTermsNegotiation(ctx context.Context, reservationID string, newLockedPrice, newReservationFee int64, message string) (*models.NegotiationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.RequestReservationNegotiation(ctx, &offerspb.RequestReservationNegotiationRequest{
		ReservationId:     reservationID,
		NewLockedPrice:    newLockedPrice,
		NewReservationFee: newReservationFee,
		Message:           message,
	})
	if err != nil {
		return nil, fmt.Errorf("RequestReservationNegotiation RPC failed: %w", err)
	}

	return &models.NegotiationResponse{
		ID:                resp.GetReservationId(),
		NegotiationStatus: resp.GetNegotiationStatus(),
	}, nil
}

// AcceptNegotiatedReservationTerms accepts a reservation negotiation
func (r OfferRepository) AcceptNegotiatedReservationTerms(ctx context.Context, reservationID string, agreedLockedPrice, agreedReservationFee int64) (*models.NegotiationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.AcceptReservationNegotiation(ctx, &offerspb.AcceptReservationNegotiationRequest{
		ReservationId:        reservationID,
		AgreedLockedPrice:    agreedLockedPrice,
		AgreedReservationFee: agreedReservationFee,
	})
	if err != nil {
		return nil, fmt.Errorf("AcceptReservationNegotiation RPC failed: %w", err)
	}

	return &models.NegotiationResponse{
		ID:                resp.GetReservationId(),
		NegotiationStatus: resp.GetNegotiationStatus(),
	}, nil
}

// DeclineReservationTermsNegotiation declines a reservation negotiation
func (r OfferRepository) DeclineReservationTermsNegotiation(ctx context.Context, reservationID, reason string) (*models.NegotiationResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := offerspb.NewOffersServiceClient(conn)
	resp, err := client.DeclineReservationNegotiation(ctx, &offerspb.DeclineReservationNegotiationRequest{
		ReservationId: reservationID,
		Reason:        reason,
	})
	if err != nil {
		return nil, fmt.Errorf("DeclineReservationNegotiation RPC failed: %w", err)
	}

	return &models.NegotiationResponse{
		ID:                resp.GetReservationId(),
		NegotiationStatus: resp.GetNegotiationStatus(),
	}, nil
}

// Additional query methods for AI tooling - These are simplified implementations

// GetAllOffersForProduct retrieves offers for a specific product
func (r OfferRepository) GetAllOffersForProduct(ctx context.Context, productID string, limit int64) ([]*models.Offer, error) {
	// Use ListOffersWithFilters with product filter - this is a simplified implementation
	resp, err := r.ListOffersWithFilters(ctx, "", "", "", 1, limit)
	if err != nil {
		return nil, err
	}

	// Filter by productID (in a real implementation, this would be done server-side)
	var offers []*models.Offer
	for _, offer := range resp.Offers {
		if offer.ProductID == productID {
			offerCopy := offer
			offers = append(offers, &offerCopy)
		}
	}

	return offers, nil
}

// GetAllOffersFromUser retrieves offers for a specific user
func (r OfferRepository) GetAllOffersFromUser(ctx context.Context, userID string, limit int64) ([]*models.Offer, error) {
	// Use ListOffersWithFilters with user filter
	resp, err := r.ListOffersWithFilters(ctx, userID, "", "", 1, limit)
	if err != nil {
		return nil, err
	}

	var offers []*models.Offer
	for _, offer := range resp.Offers {
		offerCopy := offer
		offers = append(offers, &offerCopy)
	}

	return offers, nil
}

// SearchOffersByKeyword searches offers by query - simplified implementation
func (r OfferRepository) SearchOffersByKeyword(ctx context.Context, query string, limit int64) ([]*models.Offer, error) {
	// This is a simplified implementation - in practice, you'd need a search endpoint
	resp, err := r.ListOffersWithFilters(ctx, "", "", "", 1, limit)
	if err != nil {
		return nil, err
	}

	var offers []*models.Offer
	for _, offer := range resp.Offers {
		offerCopy := offer
		offers = append(offers, &offerCopy)
	}

	return offers, nil
}

// GetActiveLeaseAgreements retrieves active leases for a user - simplified implementation
func (r OfferRepository) GetActiveLeaseAgreements(ctx context.Context, userID string, limit int64) ([]*models.Lease, error) {
	// This would require a separate endpoint in the real service
	log.Warn().Msg("GetActiveLeases: simplified implementation - returns empty slice")
	return []*models.Lease{}, nil
}

// GetActiveBuyBackAgreements retrieves active buy-backs for a user - simplified implementation
func (r OfferRepository) GetActiveBuyBackAgreements(ctx context.Context, userID string, limit int64) ([]*models.BuyBack, error) {
	// This would require a separate endpoint in the real service
	log.Warn().Msg("GetActiveBuyBacks: simplified implementation - returns empty slice")
	return []*models.BuyBack{}, nil
}

// GetActiveOfferReservations retrieves active reservations for a user - simplified implementation
func (r OfferRepository) GetActiveOfferReservations(ctx context.Context, userID string, limit int64) ([]*models.Reservation, error) {
	// This would require a separate endpoint in the real service
	log.Warn().Msg("GetActiveReservations: simplified implementation - returns empty slice")
	return []*models.Reservation{}, nil
}
