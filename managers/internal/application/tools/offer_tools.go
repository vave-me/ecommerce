package tools

import ai2 "middleman/internal/ai"

func createOfferTools() []ai2.Tool {
	return []ai2.Tool{
		// Core Offer operations
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_create",
				Description: "Create a new seller offer for a product",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"seller_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the seller creating the offer",
							"minLength":   1,
						},
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the product being offered",
							"minLength":   1,
						},
						"price": map[string]interface{}{
							"type":        "integer",
							"description": "The offer price in cents",
							"minimum":     0,
						},
					},
					"required": []string{"seller_id", "product_id", "price"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_activate",
				Description: "Activate an existing offer",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"offer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the offer to activate",
							"minLength":   1,
						},
					},
					"required": []string{"offer_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_close",
				Description: "Close an offer with a reason",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"offer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the offer to close",
							"minLength":   1,
						},
						"reason": map[string]interface{}{
							"type":        "string",
							"description": "The reason for closing the offer",
							"minLength":   1,
							"maxLength":   500,
						},
					},
					"required": []string{"offer_id", "reason"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_accept",
				Description: "Accept an offer as a customer",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"offer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the offer to accept",
							"minLength":   1,
						},
						"customer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the customer accepting the offer",
							"minLength":   1,
						},
					},
					"required": []string{"offer_id", "customer_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_get_by_id",
				Description: "Get offer details by ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"offer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the offer to retrieve",
							"minLength":   1,
						},
					},
					"required": []string{"offer_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_list",
				Description: "List offers with optional filters",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"seller_id": map[string]interface{}{
							"type":        "string",
							"description": "Filter by seller ID",
						},
						"customer_id": map[string]interface{}{
							"type":        "string",
							"description": "Filter by customer ID",
						},
						"status": map[string]interface{}{
							"type":        "string",
							"description": "Filter by offer status",
							"enum":        []string{"active", "pending", "accepted", "rejected", "closed", "expired"},
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "Page number",
							"minimum":     1,
							"default":     1,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Number of results per page",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{},
				},
			},
		},

		// Negotiation operations
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_negotiate_price",
				Description: "Request price negotiation for an offer",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"offer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the offer to negotiate",
							"minLength":   1,
						},
						"proposed_price": map[string]interface{}{
							"type":        "integer",
							"description": "The proposed price in cents",
							"minimum":     0,
						},
						"message": map[string]interface{}{
							"type":        "string",
							"description": "Negotiation message",
							"maxLength":   1000,
						},
					},
					"required": []string{"offer_id", "proposed_price"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_accept_negotiation",
				Description: "Accept a negotiated price",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"offer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the offer",
							"minLength":   1,
						},
						"final_price": map[string]interface{}{
							"type":        "integer",
							"description": "The accepted final price in cents",
							"minimum":     0,
						},
					},
					"required": []string{"offer_id", "final_price"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_decline_negotiation",
				Description: "Decline a negotiation request",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"offer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the offer",
							"minLength":   1,
						},
						"reason": map[string]interface{}{
							"type":        "string",
							"description": "Reason for declining",
							"minLength":   1,
							"maxLength":   500,
						},
					},
					"required": []string{"offer_id", "reason"},
				},
			},
		},

		// BuyNow operations
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_create_buynow",
				Description: "Create a buy now transaction from an offer",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"offer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the offer",
							"minLength":   1,
						},
						"final_price": map[string]interface{}{
							"type":        "integer",
							"description": "The final buy now price in cents",
							"minimum":     0,
						},
					},
					"required": []string{"offer_id", "final_price"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_confirm_buynow",
				Description: "Confirm a buy now purchase",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"buynow_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the buy now transaction",
							"minLength":   1,
						},
					},
					"required": []string{"buynow_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_cancel_buynow",
				Description: "Cancel a buy now transaction",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"buynow_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the buy now transaction to cancel",
							"minLength":   1,
						},
					},
					"required": []string{"buynow_id"},
				},
			},
		},

		// Lease operations
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_create_lease",
				Description: "Create a lease agreement for an offer",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"offer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the offer",
							"minLength":   1,
						},
						"monthly_price": map[string]interface{}{
							"type":        "integer",
							"description": "Monthly lease payment in cents",
							"minimum":     0,
						},
						"lease_term_months": map[string]interface{}{
							"type":        "integer",
							"description": "Lease term in months",
							"minimum":     1,
							"maximum":     60,
							"default":     12,
						},
						"has_buyout": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether the lease has a buyout option",
							"default":     false,
						},
						"buyout_price": map[string]interface{}{
							"type":        "integer",
							"description": "Buyout price if has_buyout is true",
							"minimum":     0,
						},
					},
					"required": []string{"offer_id", "monthly_price"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_start_lease",
				Description: "Start a lease agreement",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"lease_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the lease to start",
							"minLength":   1,
						},
					},
					"required": []string{"lease_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_lease_payment",
				Description: "Record a monthly lease payment",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"lease_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the lease",
							"minLength":   1,
						},
						"amount": map[string]interface{}{
							"type":        "integer",
							"description": "Payment amount in cents",
							"minimum":     0,
						},
					},
					"required": []string{"lease_id", "amount"},
				},
			},
		},

		// Reservation operations
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_create_reservation",
				Description: "Create a reservation for an offer",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"offer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the offer to reserve",
							"minLength":   1,
						},
						"lock_buyer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the buyer making the reservation",
							"minLength":   1,
						},
						"locked_price": map[string]interface{}{
							"type":        "integer",
							"description": "The locked price in cents",
							"minimum":     0,
						},
						"reservation_fee": map[string]interface{}{
							"type":        "integer",
							"description": "Reservation fee in cents",
							"minimum":     0,
						},
						"lock_duration_days": map[string]interface{}{
							"type":        "integer",
							"description": "How many days to lock the reservation",
							"minimum":     1,
							"maximum":     30,
							"default":     7,
						},
					},
					"required": []string{"offer_id", "lock_buyer_id", "locked_price", "reservation_fee"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_redeem_reservation",
				Description: "Redeem an offer reservation",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"reservation_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the reservation to redeem",
							"minLength":   1,
						},
					},
					"required": []string{"reservation_id"},
				},
			},
		},

		// BuyBack operations
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_create_buyback",
				Description: "Create a buyback agreement for an offer",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"offer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the offer",
							"minLength":   1,
						},
						"lock_buyer_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the buyer",
							"minLength":   1,
						},
						"locked_price": map[string]interface{}{
							"type":        "integer",
							"description": "The locked buyback price in cents",
							"minimum":     0,
						},
						"redemption_fee": map[string]interface{}{
							"type":        "integer",
							"description": "Fee for redeeming the buyback in cents",
							"minimum":     0,
						},
						"lock_duration_days": map[string]interface{}{
							"type":        "integer",
							"description": "How many days the buyback is valid",
							"minimum":     1,
							"maximum":     365,
							"default":     30,
						},
					},
					"required": []string{"offer_id", "lock_buyer_id", "locked_price", "redemption_fee"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_redeem_buyback",
				Description: "Redeem a buyback option",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"buyback_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the buyback to redeem",
							"minLength":   1,
						},
					},
					"required": []string{"buyback_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_expire_buyback",
				Description: "Expire a buyback agreement",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"buyback_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the buyback to expire",
							"minLength":   1,
						},
					},
					"required": []string{"buyback_id"},
				},
			},
		},

		// Query operations
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_get_by_product",
				Description: "Get all offers for a specific product",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the product",
							"minLength":   1,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{"product_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_get_by_user",
				Description: "Get all offers from a specific user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the user",
							"minLength":   1,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_search",
				Description: "Search offers by keyword",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Search query",
							"minLength":   2,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_get_active_leases",
				Description: "Get active lease agreements for a user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the user",
							"minLength":   1,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_get_active_buybacks",
				Description: "Get active buyback agreements for a user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the user",
							"minLength":   1,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "offer_get_active_reservations",
				Description: "Get active offer reservations for a user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the user",
							"minLength":   1,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
	}
}