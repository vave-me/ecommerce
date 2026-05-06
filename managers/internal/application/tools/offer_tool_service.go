package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// OfferToolService handles all offer-related operations including offers, buyNow, lease, buyback, and reservations
type OfferToolService struct {
	offerRepo domain.OfferRepository
}

// NewOfferToolService creates a new offer tool service
func NewOfferToolService(offerRepo domain.OfferRepository) *OfferToolService {
	return &OfferToolService{
		offerRepo: offerRepo,
	}
}

// ExecuteOperation handles offer-related operations with streaming support
func (s *OfferToolService) ExecuteOperation(ctx context.Context, operation string, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	log.Printf("OfferToolService.ExecuteOperation: Executing offer operation: %s", operation)

	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "started",
		Progress: 0,
		Metadata: map[string]interface{}{
			"operation": operation,
			"message":   fmt.Sprintf("Starting offer operation: %s", operation),
		},
		Timestamp: time.Now().Unix(),
	}

	var result interface{}
	var err error

	switch operation {
	// Basic Offer Operations
	case "create_offer", "create":
		result, err = s.handleCreateOffer(ctx, parameters, streamChan, toolID)
	case "activate_offer", "activate":
		result, err = s.handleActivateOffer(ctx, parameters, streamChan, toolID)
	case "close_offer", "close":
		result, err = s.handleCloseOffer(ctx, parameters, streamChan, toolID)
	case "accept_offer", "accept":
		result, err = s.handleAcceptOffer(ctx, parameters, streamChan, toolID)
	case "get_offer", "find":
		result, err = s.handleGetOffer(ctx, parameters, streamChan, toolID)
	case "list_offers", "get_offers":
		result, err = s.handleListOffers(ctx, parameters, streamChan, toolID)
	case "get_offers_by_product":
		result, err = s.handleGetOffersByProduct(ctx, parameters, streamChan, toolID)
	case "get_offers_by_user":
		result, err = s.handleGetOffersByUser(ctx, parameters, streamChan, toolID)
	case "search_offers", "search":
		result, err = s.handleSearchOffers(ctx, parameters, streamChan, toolID)

	// Offer Negotiation Operations
	case "request_offer_negotiation":
		result, err = s.handleRequestOfferNegotiation(ctx, parameters, streamChan, toolID)
	case "accept_offer_negotiation":
		result, err = s.handleAcceptOfferNegotiation(ctx, parameters, streamChan, toolID)
	case "decline_offer_negotiation":
		result, err = s.handleDeclineOfferNegotiation(ctx, parameters, streamChan, toolID)

	// BuyNow Operations
	case "create_buynow", "create_buy_now":
		result, err = s.handleCreateBuyNow(ctx, parameters, streamChan, toolID)
	case "confirm_buynow", "confirm_buy_now":
		result, err = s.handleConfirmBuyNow(ctx, parameters, streamChan, toolID)
	case "cancel_buynow", "cancel_buy_now":
		result, err = s.handleCancelBuyNow(ctx, parameters, streamChan, toolID)
	case "request_buynow_negotiation", "request_buy_now_negotiation":
		result, err = s.handleRequestBuyNowNegotiation(ctx, parameters, streamChan, toolID)
	case "accept_buynow_negotiation", "accept_buy_now_negotiation":
		result, err = s.handleAcceptBuyNowNegotiation(ctx, parameters, streamChan, toolID)
	case "decline_buynow_negotiation", "decline_buy_now_negotiation":
		result, err = s.handleDeclineBuyNowNegotiation(ctx, parameters, streamChan, toolID)

	// Lease Operations
	case "create_lease":
		result, err = s.handleCreateLease(ctx, parameters, streamChan, toolID)
	case "start_lease":
		result, err = s.handleStartLease(ctx, parameters, streamChan, toolID)
	case "make_lease_payment":
		result, err = s.handleMakeLeasePayment(ctx, parameters, streamChan, toolID)
	case "execute_lease_buyout":
		result, err = s.handleExecuteLeaseBuyout(ctx, parameters, streamChan, toolID)
	case "end_lease":
		result, err = s.handleEndLease(ctx, parameters, streamChan, toolID)
	case "cancel_lease":
		result, err = s.handleCancelLease(ctx, parameters, streamChan, toolID)
	case "default_lease":
		result, err = s.handleDefaultLease(ctx, parameters, streamChan, toolID)
	case "get_active_leases":
		result, err = s.handleGetActiveLeases(ctx, parameters, streamChan, toolID)
	case "request_lease_negotiation":
		result, err = s.handleRequestLeaseNegotiation(ctx, parameters, streamChan, toolID)
	case "accept_lease_negotiation":
		result, err = s.handleAcceptLeaseNegotiation(ctx, parameters, streamChan, toolID)
	case "decline_lease_negotiation":
		result, err = s.handleDeclineLeaseNegotiation(ctx, parameters, streamChan, toolID)

	// BuyBack Operations
	case "create_buyback", "create_buy_back":
		result, err = s.handleCreateBuyBack(ctx, parameters, streamChan, toolID)
	case "redeem_buyback", "redeem_buy_back":
		result, err = s.handleRedeemBuyBack(ctx, parameters, streamChan, toolID)
	case "expire_buyback", "expire_buy_back":
		result, err = s.handleExpireBuyBack(ctx, parameters, streamChan, toolID)
	case "cancel_buyback", "cancel_buy_back":
		result, err = s.handleCancelBuyBack(ctx, parameters, streamChan, toolID)
	case "get_active_buybacks", "get_active_buy_backs":
		result, err = s.handleGetActiveBuyBacks(ctx, parameters, streamChan, toolID)
	case "request_buyback_negotiation", "request_buy_back_negotiation":
		result, err = s.handleRequestBuyBackNegotiation(ctx, parameters, streamChan, toolID)
	case "accept_buyback_negotiation", "accept_buy_back_negotiation":
		result, err = s.handleAcceptBuyBackNegotiation(ctx, parameters, streamChan, toolID)
	case "decline_buyback_negotiation", "decline_buy_back_negotiation":
		result, err = s.handleDeclineBuyBackNegotiation(ctx, parameters, streamChan, toolID)

	// Reservation Operations
	case "create_reservation":
		result, err = s.handleCreateReservation(ctx, parameters, streamChan, toolID)
	case "redeem_reservation":
		result, err = s.handleRedeemReservation(ctx, parameters, streamChan, toolID)
	case "expire_reservation":
		result, err = s.handleExpireReservation(ctx, parameters, streamChan, toolID)
	case "cancel_reservation":
		result, err = s.handleCancelReservation(ctx, parameters, streamChan, toolID)
	case "get_active_reservations":
		result, err = s.handleGetActiveReservations(ctx, parameters, streamChan, toolID)
	case "request_reservation_negotiation":
		result, err = s.handleRequestReservationNegotiation(ctx, parameters, streamChan, toolID)
	case "accept_reservation_negotiation":
		result, err = s.handleAcceptReservationNegotiation(ctx, parameters, streamChan, toolID)
	case "decline_reservation_negotiation":
		result, err = s.handleDeclineReservationNegotiation(ctx, parameters, streamChan, toolID)

	default:
		err = fmt.Errorf("unsupported offer operation: %s", operation)
	}

	// Handle result
	if err != nil {
		streamChan <- ToolExecutionStream{
			ID:        toolID,
			ToolName:  "offer_operation",
			Status:    "error",
			Progress:  100,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
		}
		return nil, err
	}

	// Send success
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "completed",
		Progress: 100,
		Result:   result,
		Metadata: map[string]interface{}{
			"operation": operation,
			"message":   fmt.Sprintf("Offer operation %s completed successfully", operation),
		},
		Timestamp: time.Now().Unix(),
	}

	return result, nil
}

// Basic Offer Operations
func (s *OfferToolService) handleCreateOffer(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	productID := getStringParam(parameters, "product_id", "")
	buyerUserID := getStringParam(parameters, "buyer_user_id", "")
	if buyerUserID == "" {
		buyerUserID = getStringParam(parameters, "user_id", "")
	}
	offerPrice := getInt64Param(parameters, "offer_price", 0)

	if productID == "" || buyerUserID == "" || offerPrice == 0 {
		return nil, fmt.Errorf("product_id, buyer_user_id, and offer_price are required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":       "creating_offer",
			"product_id": productID,
		},
		Timestamp: time.Now().Unix(),
	}

	offerID, err := s.offerRepo.CreateOffer(ctx, productID, buyerUserID, offerPrice)
	if err != nil {
		return nil, fmt.Errorf("create offer failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "create_offer",
		"result":      "Offer created successfully",
		"offer_id":    offerID,
		"product_id":  productID,
		"offer_price": offerPrice,
	}, nil
}

func (s *OfferToolService) handleActivateOffer(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	offerID := getStringParam(parameters, "offer_id", "")
	if offerID == "" {
		offerID = getStringParam(parameters, "id", "")
	}

	if offerID == "" {
		return nil, fmt.Errorf("offer_id parameter required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":     "activating_offer",
			"offer_id": offerID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.ActivateOffer(ctx, offerID)
	if err != nil {
		return nil, fmt.Errorf("activate offer failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "activate_offer",
		"result":      "Offer activated successfully",
		"offer_id":    response.OfferID,
	}, nil
}

func (s *OfferToolService) handleCloseOffer(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	offerID := getStringParam(parameters, "offer_id", "")
	if offerID == "" {
		offerID = getStringParam(parameters, "id", "")
	}
	reason := getStringParam(parameters, "reason", "closed by user")

	if offerID == "" {
		return nil, fmt.Errorf("offer_id parameter required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":     "closing_offer",
			"offer_id": offerID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.CloseOffer(ctx, offerID, reason)
	if err != nil {
		return nil, fmt.Errorf("close offer failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "close_offer",
		"result":      "Offer closed successfully",
		"offer_id":    response.OfferID,
	}, nil
}

func (s *OfferToolService) handleAcceptOffer(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	offerID := getStringParam(parameters, "offer_id", "")
	if offerID == "" {
		offerID = getStringParam(parameters, "id", "")
	}
	userCustomerID := getStringParam(parameters, "user_customer_id", "")
	if userCustomerID == "" {
		userCustomerID = getStringParam(parameters, "user_id", "")
	}

	if offerID == "" || userCustomerID == "" {
		return nil, fmt.Errorf("offer_id and user_customer_id parameters required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":     "accepting_offer",
			"offer_id": offerID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.AcceptOffer(ctx, offerID, userCustomerID)
	if err != nil {
		return nil, fmt.Errorf("accept offer failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "accept_offer",
		"result":      "Offer accepted successfully",
		"offer_id":    response.OfferID,
	}, nil
}

func (s *OfferToolService) handleGetOffer(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	offerID := getStringParam(parameters, "offer_id", "")
	if offerID == "" {
		offerID = getStringParam(parameters, "id", "")
	}

	if offerID == "" {
		return nil, fmt.Errorf("offer_id parameter required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":     "getting_offer",
			"offer_id": offerID,
		},
		Timestamp: time.Now().Unix(),
	}

	offer, err := s.offerRepo.GetOffer(ctx, offerID)
	if err != nil {
		return nil, fmt.Errorf("get offer failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "get_offer",
		"result":      offer,
		"offer_id":    offerID,
	}, nil
}

func (s *OfferToolService) handleListOffers(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userSellerID := getStringParam(parameters, "user_seller_id", "")
	userCustomerID := getStringParam(parameters, "user_customer_id", "")
	offerStatus := getStringParam(parameters, "offer_status", "")
	page := getInt64Param(parameters, "page", 1)
	limit := getInt64Param(parameters, "limit", 20)

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step": "listing_offers",
		},
		Timestamp: time.Now().Unix(),
	}

	offers, err := s.offerRepo.ListOffers(ctx, userSellerID, userCustomerID, offerStatus, page, limit)
	if err != nil {
		return nil, fmt.Errorf("list offers failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "list_offers",
		"offers":      offers.Offers,
		"count":       len(offers.Offers),
		"page":        page,
		"limit":       limit,
	}, nil
}

func (s *OfferToolService) handleGetOffersByProduct(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	productID := getStringParam(parameters, "product_id", "")
	if productID == "" {
		return nil, fmt.Errorf("product_id parameter required")
	}

	limit := getInt64Param(parameters, "limit", 20)

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":       "getting_offers_by_product",
			"product_id": productID,
		},
		Timestamp: time.Now().Unix(),
	}

	offers, err := s.offerRepo.GetOffersByProduct(ctx, productID, limit)
	if err != nil {
		return nil, fmt.Errorf("get offers by product failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "get_offers_by_product",
		"results":     offers,
		"count":       len(offers),
		"product_id":  productID,
		"limit":       limit,
	}, nil
}

func (s *OfferToolService) handleGetOffersByUser(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	limit := getInt64Param(parameters, "limit", 20)

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "getting_offers_by_user",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	offers, err := s.offerRepo.GetOffersByUser(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get offers by user failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "get_offers_by_user",
		"results":     offers,
		"count":       len(offers),
		"user_id":     userID,
		"limit":       limit,
	}, nil
}

func (s *OfferToolService) handleSearchOffers(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	searchTerm := getStringParam(parameters, "search_term", "")
	if searchTerm == "" {
		searchTerm = getStringParam(parameters, "name", "")
	}
	if searchTerm == "" {
		return nil, fmt.Errorf("search_term parameter required")
	}

	limit := getInt64Param(parameters, "limit", 20)

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":        "searching_offers",
			"search_term": searchTerm,
		},
		Timestamp: time.Now().Unix(),
	}

	offers, err := s.offerRepo.SearchOffers(ctx, searchTerm, limit)
	if err != nil {
		return nil, fmt.Errorf("search offers failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "search_offers",
		"results":     offers,
		"count":       len(offers),
		"search_term": searchTerm,
		"limit":       limit,
	}, nil
}

// Offer Negotiation Operations
func (s *OfferToolService) handleRequestOfferNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	offerID := getStringParam(parameters, "offer_id", "")
	if offerID == "" {
		offerID = getStringParam(parameters, "id", "")
	}
	proposedPrice := getInt64Param(parameters, "proposed_price", 0)
	message := getStringParam(parameters, "message", "")

	if offerID == "" || proposedPrice == 0 {
		return nil, fmt.Errorf("offer_id and proposed_price are required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":     "requesting_offer_negotiation",
			"offer_id": offerID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.RequestOfferNegotiation(ctx, offerID, proposedPrice, message)
	if err != nil {
		return nil, fmt.Errorf("request offer negotiation failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type":    "offers",
		"operation":      "request_offer_negotiation",
		"result":         "Offer negotiation requested successfully",
		"offer_id":       offerID,
		"proposed_price": proposedPrice,
		"response":       response,
	}, nil
}

func (s *OfferToolService) handleAcceptOfferNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	offerID := getStringParam(parameters, "offer_id", "")
	if offerID == "" {
		offerID = getStringParam(parameters, "id", "")
	}
	finalPrice := getInt64Param(parameters, "final_price", 0)

	if offerID == "" || finalPrice == 0 {
		return nil, fmt.Errorf("offer_id and final_price are required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":     "accepting_offer_negotiation",
			"offer_id": offerID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.AcceptOfferNegotiation(ctx, offerID, finalPrice)
	if err != nil {
		return nil, fmt.Errorf("accept offer negotiation failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "accept_offer_negotiation",
		"result":      "Offer negotiation accepted successfully",
		"offer_id":    offerID,
		"final_price": finalPrice,
		"response":    response,
	}, nil
}

func (s *OfferToolService) handleDeclineOfferNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	offerID := getStringParam(parameters, "offer_id", "")
	if offerID == "" {
		offerID = getStringParam(parameters, "id", "")
	}
	reason := getStringParam(parameters, "reason", "declined by user")

	if offerID == "" {
		return nil, fmt.Errorf("offer_id is required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":     "declining_offer_negotiation",
			"offer_id": offerID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.DeclineOfferNegotiation(ctx, offerID, reason)
	if err != nil {
		return nil, fmt.Errorf("decline offer negotiation failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "decline_offer_negotiation",
		"result":      "Offer negotiation declined successfully",
		"offer_id":    offerID,
		"reason":      reason,
		"response":    response,
	}, nil
}

// BuyNow Operations
func (s *OfferToolService) handleCreateBuyNow(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	offerID := getStringParam(parameters, "offer_id", "")
	finalPrice := getInt64Param(parameters, "final_price", 0)

	if offerID == "" || finalPrice == 0 {
		return nil, fmt.Errorf("offer_id and final_price parameters required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":        "creating_buynow",
			"offer_id":    offerID,
			"final_price": finalPrice,
		},
		Timestamp: time.Now().Unix(),
	}

	buyNowID, err := s.offerRepo.CreateBuyNow(ctx, offerID, finalPrice)
	if err != nil {
		return nil, fmt.Errorf("create buy now failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "buynow",
		"operation":   "create_buynow",
		"result":      "BuyNow created successfully",
		"buynow_id":   buyNowID,
		"offer_id":    offerID,
		"final_price": finalPrice,
	}, nil
}

func (s *OfferToolService) handleConfirmBuyNow(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	buyNowID := getStringParam(parameters, "buynow_id", "")
	if buyNowID == "" {
		buyNowID = getStringParam(parameters, "id", "")
	}

	if buyNowID == "" {
		return nil, fmt.Errorf("buynow_id parameter required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "confirming_buynow",
			"buynow_id": buyNowID,
		},
		Timestamp: time.Now().Unix(),
	}

	result, err := s.offerRepo.ConfirmBuyNow(ctx, buyNowID)
	if err != nil {
		return nil, fmt.Errorf("confirm buy now failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "buynow",
		"operation":   "confirm_buynow",
		"result":      result,
		"buynow_id":   buyNowID,
	}, nil
}

func (s *OfferToolService) handleCancelBuyNow(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	buyNowID := getStringParam(parameters, "buynow_id", "")
	if buyNowID == "" {
		buyNowID = getStringParam(parameters, "id", "")
	}

	if buyNowID == "" {
		return nil, fmt.Errorf("buynow_id parameter required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "cancelling_buynow",
			"buynow_id": buyNowID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.CancelBuyNow(ctx, buyNowID)
	if err != nil {
		return nil, fmt.Errorf("cancel buy now failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "cancel_buy_now",
		"result":      "BuyNow offer cancelled successfully",
		"buynow_id":   response.BuyNowID,
	}, nil
}

// BuyNow Negotiation Operations
func (s *OfferToolService) handleRequestBuyNowNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	buyNowID := getStringParam(parameters, "buynow_id", "")
	if buyNowID == "" {
		buyNowID = getStringParam(parameters, "id", "")
	}
	newPrice := getInt64Param(parameters, "new_price", 0)
	message := getStringParam(parameters, "message", "")

	if buyNowID == "" || newPrice == 0 {
		return nil, fmt.Errorf("buynow_id and new_price are required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "requesting_buynow_negotiation",
			"buynow_id": buyNowID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.RequestBuyNowNegotiation(ctx, buyNowID, newPrice, message)
	if err != nil {
		return nil, fmt.Errorf("request buy now negotiation failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "request_buy_now_negotiation",
		"result":      "BuyNow negotiation requested successfully",
		"buynow_id":   buyNowID,
		"new_price":   newPrice,
		"response":    response,
	}, nil
}

func (s *OfferToolService) handleAcceptBuyNowNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	buyNowID := getStringParam(parameters, "buynow_id", "")
	if buyNowID == "" {
		buyNowID = getStringParam(parameters, "id", "")
	}
	finalPrice := getInt64Param(parameters, "final_price", 0)

	if buyNowID == "" || finalPrice == 0 {
		return nil, fmt.Errorf("buynow_id and final_price are required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "accepting_buynow_negotiation",
			"buynow_id": buyNowID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.AcceptBuyNowNegotiation(ctx, buyNowID, finalPrice)
	if err != nil {
		return nil, fmt.Errorf("accept buy now negotiation failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "accept_buy_now_negotiation",
		"result":      "BuyNow negotiation accepted successfully",
		"buynow_id":   buyNowID,
		"final_price": finalPrice,
		"response":    response,
	}, nil
}

func (s *OfferToolService) handleDeclineBuyNowNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	buyNowID := getStringParam(parameters, "buynow_id", "")
	if buyNowID == "" {
		buyNowID = getStringParam(parameters, "id", "")
	}
	reason := getStringParam(parameters, "reason", "declined by user")

	if buyNowID == "" {
		return nil, fmt.Errorf("buynow_id is required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "declining_buynow_negotiation",
			"buynow_id": buyNowID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.DeclineBuyNowNegotiation(ctx, buyNowID, reason)
	if err != nil {
		return nil, fmt.Errorf("decline buy now negotiation failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "offers",
		"operation":   "decline_buy_now_negotiation",
		"result":      "BuyNow negotiation declined successfully",
		"buynow_id":   buyNowID,
		"reason":      reason,
		"response":    response,
	}, nil
}

// Lease Operations
func (s *OfferToolService) handleCreateLease(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	offerID := getStringParam(parameters, "offer_id", "")
	monthlyPrice := getInt64Param(parameters, "monthly_price", 0)
	leaseTermMonths := getInt64Param(parameters, "lease_term_months", 0)
	hasBuyout := getBoolParam(parameters, "has_buyout", false)
	buyoutPrice := getInt64Param(parameters, "buyout_price", 0)

	if offerID == "" || monthlyPrice == 0 || leaseTermMonths == 0 {
		return nil, fmt.Errorf("offer_id, monthly_price, and lease_term_months are required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":     "creating_lease",
			"offer_id": offerID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.CreateLease(ctx, offerID, monthlyPrice, leaseTermMonths, hasBuyout, buyoutPrice)
	if err != nil {
		return nil, fmt.Errorf("create lease failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type":       "leases",
		"operation":         "create_lease",
		"result":            "Lease created successfully",
		"lease_id":          response.LeaseID,
		"offer_id":          offerID,
		"monthly_price":     monthlyPrice,
		"lease_term_months": leaseTermMonths,
	}, nil
}

func (s *OfferToolService) handleStartLease(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	leaseID := getStringParam(parameters, "lease_id", "")
	if leaseID == "" {
		leaseID = getStringParam(parameters, "id", "")
	}

	if leaseID == "" {
		return nil, fmt.Errorf("lease_id parameter required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":     "starting_lease",
			"lease_id": leaseID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.StartLease(ctx, leaseID)
	if err != nil {
		return nil, fmt.Errorf("start lease failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "leases",
		"operation":   "start_lease",
		"result":      "Lease started successfully",
		"lease_id":    response.LeaseID,
	}, nil
}

func (s *OfferToolService) handleMakeLeasePayment(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	leaseID := getStringParam(parameters, "lease_id", "")
	if leaseID == "" {
		leaseID = getStringParam(parameters, "id", "")
	}
	amount := getInt64Param(parameters, "amount", 0)

	if leaseID == "" || amount == 0 {
		return nil, fmt.Errorf("lease_id and amount are required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "offer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":     "making_lease_payment",
			"lease_id": leaseID,
		},
		Timestamp: time.Now().Unix(),
	}

	response, err := s.offerRepo.MakeLeasePayment(ctx, leaseID, amount)
	if err != nil {
		return nil, fmt.Errorf("make lease payment failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "leases",
		"operation":   "make_lease_payment",
		"result":      "Lease payment made successfully",
		"lease_id":    response.LeaseID,
		"amount":      amount,
	}, nil
}

// Placeholder implementations for the remaining methods to make the file compile
func (s *OfferToolService) handleExecuteLeaseBuyout(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "execute_lease_buyout", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleEndLease(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "end_lease", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleCancelLease(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "cancel_lease", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleDefaultLease(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "default_lease", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleGetActiveLeases(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "get_active_leases", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleRequestLeaseNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "request_lease_negotiation", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleAcceptLeaseNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "accept_lease_negotiation", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleDeclineLeaseNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "decline_lease_negotiation", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleCreateBuyBack(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "create_buy_back", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleRedeemBuyBack(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "redeem_buy_back", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleExpireBuyBack(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "expire_buy_back", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleCancelBuyBack(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "cancel_buy_back", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleGetActiveBuyBacks(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "get_active_buy_backs", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleRequestBuyBackNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "request_buy_back_negotiation", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleAcceptBuyBackNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "accept_buy_back_negotiation", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleDeclineBuyBackNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "decline_buy_back_negotiation", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleCreateReservation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "create_reservation", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleRedeemReservation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "redeem_reservation", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleExpireReservation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "expire_reservation", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleCancelReservation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "cancel_reservation", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleGetActiveReservations(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "get_active_reservations", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleRequestReservationNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "request_reservation_negotiation", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleAcceptReservationNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "accept_reservation_negotiation", "status": "not_implemented"}, nil
}

func (s *OfferToolService) handleDeclineReservationNegotiation(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	return map[string]interface{}{"operation": "decline_reservation_negotiation", "status": "not_implemented"}, nil
}
