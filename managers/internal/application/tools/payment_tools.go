package tools

import ai2 "middleman/internal/ai"

func createPaymentTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_authorize",
				Description: "Authorize a payment for a customer",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_customer_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the customer",
						},
						"amount": map[string]interface{}{
							"type":        "integer",
							"description": "Amount to authorize in cents",
							"minimum":     0,
						},
					},
					"required": []string{"user_customer_id", "amount"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_confirm",
				Description: "Confirm a payment with a payment method",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"payment_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the payment to confirm",
						},
						"payment_method_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the payment method to use",
						},
					},
					"required": []string{"payment_id", "payment_method_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_capture",
				Description: "Capture an authorized payment amount",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"payment_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the payment to capture",
						},
						"amount_to_capture": map[string]interface{}{
							"type":        "integer",
							"description": "Amount to capture in cents (can be less than authorized)",
							"minimum":     0,
						},
					},
					"required": []string{"payment_id", "amount_to_capture"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_create_invoice",
				Description: "Create a new invoice for an order",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"order_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the order",
						},
						"payment_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the payment",
						},
						"amount": map[string]interface{}{
							"type":        "integer",
							"description": "Invoice amount in cents",
							"minimum":     0,
						},
					},
					"required": []string{"order_id", "payment_id", "amount"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_get",
				Description: "Get payment details by ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"payment_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the payment to retrieve",
						},
					},
					"required": []string{"payment_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_adjust_invoice",
				Description: "Adjust invoice amount with reason",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"invoice_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the invoice to adjust",
						},
						"amount": map[string]interface{}{
							"type":        "integer",
							"description": "New amount in cents",
							"minimum":     0,
						},
						"reason": map[string]interface{}{
							"type":        "string",
							"description": "Reason for adjustment",
						},
					},
					"required": []string{"invoice_id", "amount", "reason"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_get_customer_history",
				Description: "Get payment history for a customer",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_customer_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the customer",
						},
					},
					"required": []string{"user_customer_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_pay_invoice",
				Description: "Pay an invoice using a payment method",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"invoice_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the invoice to pay",
						},
						"payment_method_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the payment method to use",
						},
					},
					"required": []string{"invoice_id", "payment_method_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_cancel_invoice",
				Description: "Cancel an invoice with reason",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"invoice_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the invoice to cancel",
						},
						"reason": map[string]interface{}{
							"type":        "string",
							"description": "Reason for cancellation",
						},
					},
					"required": []string{"invoice_id", "reason"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_handle_webhook",
				Description: "Process payment webhook notification",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"raw_body": map[string]interface{}{
							"type":        "string",
							"description": "Raw webhook body",
						},
						"signature": map[string]interface{}{
							"type":        "string",
							"description": "Webhook signature for verification",
						},
					},
					"required": []string{"raw_body", "signature"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_get_invoice",
				Description: "Get invoice details by ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"invoice_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the invoice to retrieve",
						},
					},
					"required": []string{"invoice_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_get_invoices_by_order",
				Description: "Get all invoices for an order",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"order_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the order",
						},
					},
					"required": []string{"order_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_search_by_status",
				Description: "Search payments by status",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"status": map[string]interface{}{
							"type":        "string",
							"description": "Payment status to search for",
							"enum":        []string{"pending", "authorized", "captured", "failed", "refunded", "cancelled"},
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     50,
						},
					},
					"required": []string{"status"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "payment_search_invoices",
				Description: "Search invoices by status",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"status": map[string]interface{}{
							"type":        "string",
							"description": "Invoice status to search for",
							"enum":        []string{"draft", "sent", "paid", "overdue", "cancelled", "void"},
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     50,
						},
					},
					"required": []string{"status"},
				},
			},
		},
	}
}