package tools

import (
	"context"
)

// ==================== OFFER HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeOfferHandlers() {
	// Core Offer operations
	r.handlers["offer_create"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userSellerID := getStringParam(params, "seller_id")
		productID := getStringParam(params, "product_id")
		price := getInt64Param(params, "price", 0)
		
		// Validate required parameters
		v := NewValidator()
		v.ValidateRequired("seller_id", userSellerID)
		v.ValidateRequired("product_id", productID)
		v.ValidateMinimum("price", float64(price), 0)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.CreateNewSellerOffer(ctx, userSellerID, productID, price)
	}

	r.handlers["offer_activate"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		offerID := getStringParam(params, "offer_id")
		
		if err := ValidateIDParam("offer_id", offerID); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.ActivateExistingOffer(ctx, offerID)
	}

	r.handlers["offer_close"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		offerID := getStringParam(params, "offer_id")
		reason := getStringParam(params, "reason")
		
		v := NewValidator()
		v.ValidateRequired("offer_id", offerID)
		v.ValidateRequired("reason", reason).ValidateMaxLength("reason", reason, 500)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		reason = SanitizeString(reason)
		return reg.offerRepo.CloseOfferWithReason(ctx, offerID, reason)
	}

	r.handlers["offer_accept"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		offerID := getStringParam(params, "offer_id")
		userCustomerID := getStringParam(params, "customer_id")
		
		v := NewValidator()
		v.ValidateRequired("offer_id", offerID)
		v.ValidateRequired("customer_id", userCustomerID)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.AcceptOfferByCustomer(ctx, offerID, userCustomerID)
	}

	r.handlers["offer_get_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		offerID := getStringParam(params, "offer_id")
		
		if err := ValidateIDParam("offer_id", offerID); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.GetOfferDetailsByID(ctx, offerID)
	}

	r.handlers["offer_list"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userSellerID := getStringParam(params, "seller_id")
		userCustomerID := getStringParam(params, "customer_id")
		offerStatus := getStringParam(params, "status")
		page := getInt64Param(params, "page", 1)
		limit := getInt64Param(params, "limit", 20)
		
		if err := ValidatePaginationParams(page, limit); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.ListOffersWithFilters(ctx, userSellerID, userCustomerID, offerStatus, page, limit)
	}

	// Offer Negotiation operations
	r.handlers["offer_negotiate_price"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		offerID := getStringParam(params, "offer_id")
		proposedPrice := getInt64Param(params, "proposed_price", 0)
		message := getStringParam(params, "message")
		
		v := NewValidator()
		v.ValidateRequired("offer_id", offerID)
		v.ValidateMinimum("proposed_price", float64(proposedPrice), 0)
		v.ValidateMaxLength("message", message, 1000)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		message = SanitizeString(message)
		return reg.offerRepo.RequestOfferPriceNegotiation(ctx, offerID, proposedPrice, message)
	}

	r.handlers["offer_accept_negotiation"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		offerID := getStringParam(params, "offer_id")
		finalPrice := getInt64Param(params, "final_price", 0)
		
		v := NewValidator()
		v.ValidateRequired("offer_id", offerID)
		v.ValidateMinimum("final_price", float64(finalPrice), 0)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.AcceptNegotiatedOfferPrice(ctx, offerID, finalPrice)
	}

	r.handlers["offer_decline_negotiation"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		offerID := getStringParam(params, "offer_id")
		reason := getStringParam(params, "reason")
		
		v := NewValidator()
		v.ValidateRequired("offer_id", offerID)
		v.ValidateRequired("reason", reason).ValidateMaxLength("reason", reason, 500)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		reason = SanitizeString(reason)
		return reg.offerRepo.DeclineOfferNegotiationRequest(ctx, offerID, reason)
	}

	// BuyNow operations
	r.handlers["offer_create_buynow"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		offerID := getStringParam(params, "offer_id")
		finalPrice := getInt64Param(params, "final_price", 0)
		
		v := NewValidator()
		v.ValidateRequired("offer_id", offerID)
		v.ValidateMinimum("final_price", float64(finalPrice), 0)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.CreateBuyNowTransaction(ctx, offerID, finalPrice)
	}

	r.handlers["offer_confirm_buynow"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		buyNowID := getStringParam(params, "buynow_id")
		
		if err := ValidateIDParam("buynow_id", buyNowID); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.ConfirmBuyNowPurchase(ctx, buyNowID)
	}

	r.handlers["offer_cancel_buynow"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		buyNowID := getStringParam(params, "buynow_id")
		
		if err := ValidateIDParam("buynow_id", buyNowID); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.CancelBuyNowTransaction(ctx, buyNowID)
	}

	// Lease operations
	r.handlers["offer_create_lease"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		offerID := getStringParam(params, "offer_id")
		monthlyPrice := getInt64Param(params, "monthly_price", 0)
		leaseTermMonths := getInt64Param(params, "lease_term_months", 12)
		hasBuyout := getBoolParam(params, "has_buyout", false)
		buyoutPrice := getInt64Param(params, "buyout_price", 0)
		
		v := NewValidator()
		v.ValidateRequired("offer_id", offerID)
		v.ValidateMinimum("monthly_price", float64(monthlyPrice), 0)
		v.ValidateRange("lease_term_months", float64(leaseTermMonths), 1, 60)
		
		if hasBuyout {
			v.ValidateMinimum("buyout_price", float64(buyoutPrice), 0)
		}
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.CreateLeaseAgreementForOffer(ctx, offerID, monthlyPrice, leaseTermMonths, hasBuyout, buyoutPrice)
	}

	r.handlers["offer_start_lease"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		leaseID := getStringParam(params, "lease_id")
		
		if err := ValidateIDParam("lease_id", leaseID); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.StartActiveLeaseAgreement(ctx, leaseID)
	}

	r.handlers["offer_lease_payment"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		leaseID := getStringParam(params, "lease_id")
		amount := getInt64Param(params, "amount", 0)
		
		v := NewValidator()
		v.ValidateRequired("lease_id", leaseID)
		v.ValidateMinimum("amount", float64(amount), 0)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.RecordMonthlyLeasePayment(ctx, leaseID, amount)
	}

	// Reservation operations
	r.handlers["offer_create_reservation"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		offerID := getStringParam(params, "offer_id")
		lockedPrice := getInt64Param(params, "locked_price", 0)
		reservationFee := getInt64Param(params, "reservation_fee", 0)
		lockDurationDays := int32(getInt64Param(params, "lock_duration_days", 7))
		lockBuyerID := getStringParam(params, "lock_buyer_id")
		
		v := NewValidator()
		v.ValidateRequired("offer_id", offerID)
		v.ValidateRequired("lock_buyer_id", lockBuyerID)
		v.ValidateMinimum("locked_price", float64(lockedPrice), 0)
		v.ValidateMinimum("reservation_fee", float64(reservationFee), 0)
		v.ValidateRange("lock_duration_days", float64(lockDurationDays), 1, 30)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.CreateOfferReservation(ctx, offerID, lockedPrice, reservationFee, lockDurationDays, lockBuyerID)
	}

	r.handlers["offer_redeem_reservation"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		reservationID := getStringParam(params, "reservation_id")
		
		if err := ValidateIDParam("reservation_id", reservationID); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.RedeemOfferReservation(ctx, reservationID)
	}

	// Query operations
	r.handlers["offer_get_by_product"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		limit := getInt64Param(params, "limit", 20)
		
		v := NewValidator()
		v.ValidateRequired("product_id", productID)
		v.ValidateRange("limit", float64(limit), 1, 100)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.GetAllOffersForProduct(ctx, productID, limit)
	}

	r.handlers["offer_get_by_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		limit := getInt64Param(params, "limit", 20)
		
		v := NewValidator()
		v.ValidateRequired("user_id", userID)
		v.ValidateRange("limit", float64(limit), 1, 100)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.GetAllOffersFromUser(ctx, userID, limit)
	}

	r.handlers["offer_search"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		query := getStringParam(params, "query")
		limit := getInt64Param(params, "limit", 20)
		
		v := NewValidator()
		v.ValidateRequired("query", query).ValidateMinLength("query", query, 2)
		v.ValidateRange("limit", float64(limit), 1, 100)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		query = SanitizeString(query)
		return reg.offerRepo.SearchOffersByKeyword(ctx, query, limit)
	}

	r.handlers["offer_get_active_leases"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		limit := getInt64Param(params, "limit", 20)
		
		v := NewValidator()
		v.ValidateRequired("user_id", userID)
		v.ValidateRange("limit", float64(limit), 1, 100)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.GetActiveLeaseAgreements(ctx, userID, limit)
	}

	r.handlers["offer_get_active_buybacks"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		limit := getInt64Param(params, "limit", 20)
		
		v := NewValidator()
		v.ValidateRequired("user_id", userID)
		v.ValidateRange("limit", float64(limit), 1, 100)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.GetActiveBuyBackAgreements(ctx, userID, limit)
	}

	r.handlers["offer_get_active_reservations"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		limit := getInt64Param(params, "limit", 20)
		
		v := NewValidator()
		v.ValidateRequired("user_id", userID)
		v.ValidateRange("limit", float64(limit), 1, 100)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.GetActiveOfferReservations(ctx, userID, limit)
	}

	// BuyBack operations
	r.handlers["offer_create_buyback"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		offerID := getStringParam(params, "offer_id")
		lockedPrice := getInt64Param(params, "locked_price", 0)
		redemptionFee := getInt64Param(params, "redemption_fee", 0)
		lockDurationDays := int32(getInt64Param(params, "lock_duration_days", 30))
		lockBuyerID := getStringParam(params, "lock_buyer_id")
		
		v := NewValidator()
		v.ValidateRequired("offer_id", offerID)
		v.ValidateRequired("lock_buyer_id", lockBuyerID)
		v.ValidateMinimum("locked_price", float64(lockedPrice), 0)
		v.ValidateMinimum("redemption_fee", float64(redemptionFee), 0)
		v.ValidateRange("lock_duration_days", float64(lockDurationDays), 1, 365)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.CreateBuyBackAgreement(ctx, offerID, lockedPrice, redemptionFee, lockDurationDays, lockBuyerID)
	}

	r.handlers["offer_redeem_buyback"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		buyBackID := getStringParam(params, "buyback_id")
		
		if err := ValidateIDParam("buyback_id", buyBackID); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.RedeemBuyBackOption(ctx, buyBackID)
	}

	r.handlers["offer_expire_buyback"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		buyBackID := getStringParam(params, "buyback_id")
		
		if err := ValidateIDParam("buyback_id", buyBackID); err != nil {
			return nil, err
		}
		
		return reg.offerRepo.ExpireBuyBackAgreement(ctx, buyBackID)
	}

	// Note: Some handlers from the original file don't have corresponding repository methods
	// These would need to be implemented in the repository or removed from the handlers
	// Examples: offer_get_history, offer_get_user_offers, offer_counter, etc.
}