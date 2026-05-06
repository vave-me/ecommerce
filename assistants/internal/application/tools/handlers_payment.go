package tools

import (
	"context"
	"fmt"
)

// ==================== PAYMENT HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializePaymentHandlers() {
	r.handlers["payment_authorize"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userCustomerID := getStringParam(params, "user_customer_id")
		amount := getInt64Param(params, "amount", 0)
		if err := ValidateIDParam("user_customer_id", userCustomerID); err != nil {
			return nil, fmt.Errorf("invalid user_customer_id: %w", err)
		}
		if amount <= 0 {
			return nil, fmt.Errorf("amount must be greater than zero")
		}
		return reg.paymentRepo.AuthorizePaymentForCustomer(ctx, userCustomerID, amount)
	}

	r.handlers["payment_confirm"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		paymentID := getStringParam(params, "payment_id")
		paymentMethodID := getStringParam(params, "payment_method_id")
		if err := ValidateIDParam("payment_id", paymentID); err != nil {
			return nil, fmt.Errorf("invalid payment_id: %w", err)
		}
		if err := ValidateIDParam("payment_method_id", paymentMethodID); err != nil {
			return nil, fmt.Errorf("invalid payment_method_id: %w", err)
		}
		return reg.paymentRepo.ConfirmPaymentWithMethod(ctx, paymentID, paymentMethodID)
	}

	r.handlers["payment_capture"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		paymentID := getStringParam(params, "payment_id")
		amountToCapture := getInt64Param(params, "amount_to_capture", 0)
		if err := ValidateIDParam("payment_id", paymentID); err != nil {
			return nil, fmt.Errorf("invalid payment_id: %w", err)
		}
		if amountToCapture <= 0 {
			return nil, fmt.Errorf("amount_to_capture must be greater than zero")
		}
		return reg.paymentRepo.CaptureAuthorizedPaymentAmount(ctx, paymentID, amountToCapture)
	}

	r.handlers["payment_create_invoice"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		orderID := getStringParam(params, "order_id")
		paymentID := getStringParam(params, "payment_id")
		amount := getInt64Param(params, "amount", 0)
		if err := ValidateIDParam("order_id", orderID); err != nil {
			return nil, fmt.Errorf("invalid order_id: %w", err)
		}
		if err := ValidateIDParam("payment_id", paymentID); err != nil {
			return nil, fmt.Errorf("invalid payment_id: %w", err)
		}
		if amount <= 0 {
			return nil, fmt.Errorf("amount must be greater than zero")
		}
		return reg.paymentRepo.CreateNewInvoiceForOrder(ctx, orderID, paymentID, amount)
	}

	r.handlers["payment_get"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		paymentID := getStringParam(params, "payment_id")
		if err := ValidateIDParam("payment_id", paymentID); err != nil {
			return nil, fmt.Errorf("invalid payment_id: %w", err)
		}
		return reg.paymentRepo.GetPaymentDetailsByID(ctx, paymentID)
	}

	r.handlers["payment_adjust_invoice"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		invoiceID := getStringParam(params, "invoice_id")
		amount := getInt64Param(params, "amount", 0)
		reason := getStringParam(params, "reason")
		if err := ValidateIDParam("invoice_id", invoiceID); err != nil {
			return nil, fmt.Errorf("invalid invoice_id: %w", err)
		}
		if amount == 0 {
			return nil, fmt.Errorf("amount cannot be zero")
		}
		if reason == "" {
			return nil, fmt.Errorf("adjustment reason is required")
		}
		reason = SanitizeString(reason)
		return reg.paymentRepo.AdjustInvoiceAmountWithReason(ctx, invoiceID, amount, reason)
	}

	r.handlers["payment_get_customer_history"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userCustomerID := getStringParam(params, "user_customer_id")
		if err := ValidateIDParam("user_customer_id", userCustomerID); err != nil {
			return nil, fmt.Errorf("invalid user_customer_id: %w", err)
		}
		return reg.paymentRepo.GetCustomerPaymentHistory(ctx, userCustomerID)
	}

	r.handlers["payment_pay_invoice"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		invoiceID := getStringParam(params, "invoice_id")
		paymentMethodID := getStringParam(params, "payment_method_id")
		if err := ValidateIDParam("invoice_id", invoiceID); err != nil {
			return nil, fmt.Errorf("invalid invoice_id: %w", err)
		}
		if err := ValidateIDParam("payment_method_id", paymentMethodID); err != nil {
			return nil, fmt.Errorf("invalid payment_method_id: %w", err)
		}
		return reg.paymentRepo.PayInvoiceUsingPaymentMethod(ctx, invoiceID, paymentMethodID)
	}

	r.handlers["payment_cancel_invoice"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		invoiceID := getStringParam(params, "invoice_id")
		reason := getStringParam(params, "reason")
		if err := ValidateIDParam("invoice_id", invoiceID); err != nil {
			return nil, fmt.Errorf("invalid invoice_id: %w", err)
		}
		if reason == "" {
			return nil, fmt.Errorf("cancellation reason is required")
		}
		reason = SanitizeString(reason)
		return nil, reg.paymentRepo.CancelInvoiceWithReason(ctx, invoiceID, reason)
	}

	r.handlers["payment_handle_webhook"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		rawBody := getStringParam(params, "raw_body")
		signature := getStringParam(params, "signature")
		if rawBody == "" {
			return nil, fmt.Errorf("raw_body is required")
		}
		if signature == "" {
			return nil, fmt.Errorf("signature is required")
		}
		return reg.paymentRepo.ProcessWebhookNotification(ctx, rawBody, signature)
	}

	r.handlers["payment_get_invoice"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		invoiceID := getStringParam(params, "invoice_id")
		if err := ValidateIDParam("invoice_id", invoiceID); err != nil {
			return nil, fmt.Errorf("invalid invoice_id: %w", err)
		}
		return reg.paymentRepo.GetInvoiceDetailsByID(ctx, invoiceID)
	}

	r.handlers["payment_get_invoices_by_order"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		orderID := getStringParam(params, "order_id")
		if err := ValidateIDParam("order_id", orderID); err != nil {
			return nil, fmt.Errorf("invalid order_id: %w", err)
		}
		return reg.paymentRepo.GetAllInvoicesForOrder(ctx, orderID)
	}

	r.handlers["payment_search_by_status"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		status := getStringParam(params, "status")
		limit := getInt64Param(params, "limit", 50)
		if status == "" {
			return nil, fmt.Errorf("status is required")
		}
		if err := ValidatePaginationParams(1, limit); err != nil {
			return nil, fmt.Errorf("invalid limit: %w", err)
		}
		return reg.paymentRepo.SearchPaymentsByStatus(ctx, status, limit)
	}

	r.handlers["payment_search_invoices"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		status := getStringParam(params, "status")
		limit := getInt64Param(params, "limit", 50)
		if status == "" {
			return nil, fmt.Errorf("status is required")
		}
		if err := ValidatePaginationParams(1, limit); err != nil {
			return nil, fmt.Errorf("invalid limit: %w", err)
		}
		return reg.paymentRepo.SearchInvoicesByStatus(ctx, status, limit)
	}
}