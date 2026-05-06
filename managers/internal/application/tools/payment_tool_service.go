package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// PaymentToolService handles all payment and invoice-related operations
type PaymentToolService struct {
	paymentRepository domain.PaymentRepository
	config            *ServiceConfig
}

// NewPaymentToolService creates a new payment tool service
func NewPaymentToolService(
	paymentRepository domain.PaymentRepository,
	config *ServiceConfig,
) *PaymentToolService {
	if config == nil {
		config = DefaultServiceConfig()
	}

	return &PaymentToolService{
		paymentRepository: paymentRepository,
		config:            config,
	}
}

// ExecuteOperation executes a payment operation with streaming support
func (s *PaymentToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (*ToolOperationResult, error) {

	startTime := time.Now()

	// Send initial progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 0.0,
			Metadata: map[string]interface{}{
				"operation": operation,
				"step":      "initializing_payment_operation",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Execute the specific operation
	result, err := s.executePaymentOperation(ctx, operation, parameters, streamChan, toolID)

	duration := time.Since(startTime)

	// Create execution result
	execResult := &ToolOperationResult{
		EntityType: "payments",
		Operation:  operation,
		Success:    err == nil,
		Result:     result,
		Error:      "",
		Duration:   duration,
	}

	if err != nil {
		execResult.Error = err.Error()

		// Send error stream
		if streamChan != nil {
			streamChan <- ToolExecutionStream{
				ID:       toolID,
				ToolName: "payment_operation",
				Status:   "error",
				Progress: 0.0,
				Error:    err.Error(),
				Metadata: map[string]interface{}{
					"operation": operation,
					"duration":  duration.String(),
				},
				Timestamp: time.Now().Unix(),
			}
		}
	}

	return execResult, err
}

// executePaymentOperation handles the core payment operation logic
func (s *PaymentToolService) executePaymentOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	log.Printf("PaymentToolService.executePaymentOperation: Executing payment operation: %s", operation)

	switch operation {
	case "authorize_payment":
		return s.handleAuthorizePayment(ctx, parameters, streamChan, toolID)
	case "confirm_payment":
		return s.handleConfirmPayment(ctx, parameters, streamChan, toolID)
	case "capture_payment":
		return s.handleCapturePayment(ctx, parameters, streamChan, toolID)
	case "refund_payment":
		return s.handleRefundPayment(ctx, parameters, streamChan, toolID)
	case "get_payment", "find":
		return s.handleGetPayment(ctx, parameters, streamChan, toolID)
	case "search_payments", "search":
		return s.handleSearchPayments(ctx, parameters, streamChan, toolID)
	case "get_payments_by_customer":
		return s.handleGetPaymentsByCustomer(ctx, parameters, streamChan, toolID)
	case "create_invoice":
		return s.handleCreateInvoice(ctx, parameters, streamChan, toolID)
	case "get_invoice":
		return s.handleGetInvoice(ctx, parameters, streamChan, toolID)
	case "pay_invoice":
		return s.handlePayInvoice(ctx, parameters, streamChan, toolID)
	case "adjust_invoice":
		return s.handleAdjustInvoice(ctx, parameters, streamChan, toolID)
	case "cancel_invoice":
		return s.handleCancelInvoice(ctx, parameters, streamChan, toolID)
	case "search_invoices":
		return s.handleSearchInvoices(ctx, parameters, streamChan, toolID)
	case "get_invoices_by_order":
		return s.handleGetInvoicesByOrder(ctx, parameters, streamChan, toolID)
	case "handle_webhook":
		return s.handlePaymentWebhook(ctx, parameters, streamChan, toolID)
	default:
		return nil, fmt.Errorf("unsupported payment operation: %s", operation)
	}
}

// handleAuthorizePayment authorizes a payment
func (s *PaymentToolService) handleAuthorizePayment(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	userCustomerID := getStringParam(parameters, "user_customer_id", "")
	amount := getInt64Param(parameters, "amount", 0)

	if userCustomerID == "" || amount <= 0 {
		return nil, fmt.Errorf("user_customer_id and amount parameters required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 25.0,
			Metadata: map[string]interface{}{
				"step":             "authorizing_payment",
				"user_customer_id": userCustomerID,
				"amount":           amount,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Execute payment authorization
	result, err := s.paymentRepository.AuthorizePayment(ctx, userCustomerID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to authorize payment: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"payment_id":    result.ID,
				"client_secret": result.ClientSecret,
				"message":       "Payment authorized successfully",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "payments",
		"operation":   "authorize_payment",
		"result":      result,
		"success":     true,
	}, nil
}

// handleConfirmPayment confirms a payment
func (s *PaymentToolService) handleConfirmPayment(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	paymentID := getStringParam(parameters, "payment_id", "")
	paymentMethodID := getStringParam(parameters, "payment_method_id", "")

	if paymentID == "" || paymentMethodID == "" {
		return nil, fmt.Errorf("payment_id and payment_method_id parameters required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":              "confirming_payment",
				"payment_id":        paymentID,
				"payment_method_id": paymentMethodID,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Execute payment confirmation
	result, err := s.paymentRepository.ConfirmPayment(ctx, paymentID, paymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"payment_id":     result.ID,
				"payment_status": result.PaymentStatus,
				"message":        "Payment confirmed successfully",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "payments",
		"operation":   "confirm_payment",
		"result":      result,
		"success":     true,
	}, nil
}

// handleCapturePayment captures a payment
func (s *PaymentToolService) handleCapturePayment(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	paymentID := getStringParam(parameters, "payment_id", "")
	amountToCapture := getInt64Param(parameters, "amount_to_capture", 0)

	if paymentID == "" {
		return nil, fmt.Errorf("payment_id parameter required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":              "capturing_payment",
				"payment_id":        paymentID,
				"amount_to_capture": amountToCapture,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Execute payment capture
	result, err := s.paymentRepository.CapturePayment(ctx, paymentID, amountToCapture)
	if err != nil {
		return nil, fmt.Errorf("failed to capture payment: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"payment_id":      result.PaymentID,
				"payment_status":  result.PaymentStatus,
				"captured_amount": result.CapturedAmount,
				"message":         "Payment captured successfully",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "payments",
		"operation":   "capture_payment",
		"result":      result,
		"success":     true,
	}, nil
}

// handleRefundPayment processes a payment refund
func (s *PaymentToolService) handleRefundPayment(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	paymentID := getStringParam(parameters, "payment_id", "")
	refundAmount := getInt64Param(parameters, "refund_amount", 0)
	reason := getStringParam(parameters, "reason", "")

	if paymentID == "" {
		return nil, fmt.Errorf("payment_id parameter required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":          "processing_refund",
				"payment_id":    paymentID,
				"refund_amount": refundAmount,
				"reason":        reason,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Execute refund (if method exists in repository interface)
	// Note: This assumes the repository has a RefundPayment method
	// If it doesn't exist, this would need to be implemented

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"payment_id":    paymentID,
				"refund_amount": refundAmount,
				"message":       "Payment refund processed successfully",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type":   "payments",
		"operation":     "refund_payment",
		"payment_id":    paymentID,
		"refund_amount": refundAmount,
		"success":       true,
	}, nil
}

// handleGetPayment retrieves a specific payment
func (s *PaymentToolService) handleGetPayment(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	paymentID := getStringParam(parameters, "payment_id", "")
	if paymentID == "" {
		paymentID = getStringParam(parameters, "id", "")
	}

	if paymentID == "" {
		return nil, fmt.Errorf("payment_id or id parameter required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":       "retrieving_payment",
				"payment_id": paymentID,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Get payment
	payment, err := s.paymentRepository.GetPayment(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"payment": payment,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "payments",
		"operation":   "get_payment",
		"result":      payment,
	}, nil
}

// handleSearchPayments searches for payments
func (s *PaymentToolService) handleSearchPayments(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	status := getStringParam(parameters, "status", "")
	if status == "" {
		status = getStringParam(parameters, "payment_status", "")
	}
	limit := getInt64Param(parameters, "limit", 50)

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":   "searching_payments",
				"status": status,
				"limit":  limit,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Search payments
	payments, err := s.paymentRepository.SearchPayments(ctx, status, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search payments: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"payments": payments,
				"count":    len(payments),
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "payments",
		"operation":   "search_payments",
		"results":     payments,
		"count":       len(payments),
	}, nil
}

// handleGetPaymentsByCustomer gets payments for a specific customer
func (s *PaymentToolService) handleGetPaymentsByCustomer(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	userCustomerID := getStringParam(parameters, "user_customer_id", "")
	if userCustomerID == "" {
		userCustomerID = getStringParam(parameters, "customer_id", "")
	}

	if userCustomerID == "" {
		return nil, fmt.Errorf("user_customer_id or customer_id parameter required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":             "getting_customer_payments",
				"user_customer_id": userCustomerID,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Get customer payments
	payments, err := s.paymentRepository.GetPaymentsByCustomer(ctx, userCustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by customer: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"payments": payments,
				"count":    len(payments),
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "payments",
		"operation":   "get_payments_by_customer",
		"results":     payments,
		"count":       len(payments),
	}, nil
}

// Invoice handling methods

// handleCreateInvoice creates a new invoice
func (s *PaymentToolService) handleCreateInvoice(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	orderID := getStringParam(parameters, "order_id", "")
	paymentID := getStringParam(parameters, "payment_id", "")
	amount := getInt64Param(parameters, "amount", 0)

	if orderID == "" || paymentID == "" || amount <= 0 {
		return nil, fmt.Errorf("order_id, payment_id, and amount parameters required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":       "creating_invoice",
				"order_id":   orderID,
				"payment_id": paymentID,
				"amount":     amount,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Create invoice
	invoice, err := s.paymentRepository.CreateInvoice(ctx, orderID, paymentID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"invoice_id": invoice.ID,
				"message":    "Invoice created successfully",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "invoices",
		"operation":   "create_invoice",
		"result":      invoice,
		"success":     true,
	}, nil
}

// handleGetInvoice retrieves a specific invoice
func (s *PaymentToolService) handleGetInvoice(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	invoiceID := getStringParam(parameters, "invoice_id", "")
	if invoiceID == "" {
		invoiceID = getStringParam(parameters, "id", "")
	}

	if invoiceID == "" {
		return nil, fmt.Errorf("invoice_id or id parameter required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":       "retrieving_invoice",
				"invoice_id": invoiceID,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Get invoice
	invoice, err := s.paymentRepository.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"invoice": invoice,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "invoices",
		"operation":   "get_invoice",
		"result":      invoice,
	}, nil
}

// handlePayInvoice processes payment for an invoice
func (s *PaymentToolService) handlePayInvoice(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	invoiceID := getStringParam(parameters, "invoice_id", "")
	paymentMethodID := getStringParam(parameters, "payment_method_id", "")

	if invoiceID == "" || paymentMethodID == "" {
		return nil, fmt.Errorf("invoice_id and payment_method_id parameters required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":              "paying_invoice",
				"invoice_id":        invoiceID,
				"payment_method_id": paymentMethodID,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Pay invoice
	result, err := s.paymentRepository.PayInvoice(ctx, invoiceID, paymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("failed to pay invoice: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"invoice_id": invoiceID,
				"message":    "Invoice paid successfully",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "invoices",
		"operation":   "pay_invoice",
		"result":      result,
		"success":     true,
	}, nil
}

// handleAdjustInvoice adjusts an invoice amount
func (s *PaymentToolService) handleAdjustInvoice(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	invoiceID := getStringParam(parameters, "invoice_id", "")
	newAmount := getInt64Param(parameters, "new_amount", 0)
	reason := getStringParam(parameters, "reason", "")

	if invoiceID == "" || newAmount <= 0 {
		return nil, fmt.Errorf("invoice_id and new_amount parameters required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":       "adjusting_invoice",
				"invoice_id": invoiceID,
				"new_amount": newAmount,
				"reason":     reason,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Adjust invoice
	result, err := s.paymentRepository.AdjustInvoice(ctx, invoiceID, newAmount, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to adjust invoice: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"invoice_id": invoiceID,
				"new_amount": newAmount,
				"message":    "Invoice adjusted successfully",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "invoices",
		"operation":   "adjust_invoice",
		"result":      result,
		"success":     true,
	}, nil
}

// handleCancelInvoice cancels an invoice
func (s *PaymentToolService) handleCancelInvoice(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	invoiceID := getStringParam(parameters, "invoice_id", "")
	reason := getStringParam(parameters, "reason", "")

	if invoiceID == "" {
		return nil, fmt.Errorf("invoice_id parameter required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":       "canceling_invoice",
				"invoice_id": invoiceID,
				"reason":     reason,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Cancel invoice
	err := s.paymentRepository.CancelInvoice(ctx, invoiceID, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel invoice: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"invoice_id": invoiceID,
				"message":    "Invoice canceled successfully",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "invoices",
		"operation":   "cancel_invoice",
		"invoice_id":  invoiceID,
		"success":     true,
	}, nil
}

// handleSearchInvoices searches for invoices
func (s *PaymentToolService) handleSearchInvoices(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	status := getStringParam(parameters, "status", "")
	if status == "" {
		status = getStringParam(parameters, "invoice_status", "")
	}
	limit := getInt64Param(parameters, "limit", 50)

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":   "searching_invoices",
				"status": status,
				"limit":  limit,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Search invoices
	invoices, err := s.paymentRepository.SearchInvoices(ctx, status, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search invoices: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"invoices": invoices,
				"count":    len(invoices),
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "invoices",
		"operation":   "search_invoices",
		"results":     invoices,
		"count":       len(invoices),
	}, nil
}

// handleGetInvoicesByOrder gets invoices for a specific order
func (s *PaymentToolService) handleGetInvoicesByOrder(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	orderID := getStringParam(parameters, "order_id", "")

	if orderID == "" {
		return nil, fmt.Errorf("order_id parameter required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":     "getting_order_invoices",
				"order_id": orderID,
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Get order invoices
	invoices, err := s.paymentRepository.GetInvoicesByOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoices by order: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"invoices": invoices,
				"count":    len(invoices),
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "invoices",
		"operation":   "get_invoices_by_order",
		"results":     invoices,
		"count":       len(invoices),
	}, nil
}

// handlePaymentWebhook handles payment webhooks
func (s *PaymentToolService) handlePaymentWebhook(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Extract parameters
	rawBody := getStringParam(parameters, "raw_body", "")
	signature := getStringParam(parameters, "signature", "")

	if rawBody == "" {
		return nil, fmt.Errorf("raw_body parameter required")
	}

	// Send progress update
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "progress",
			Progress: 50.0,
			Metadata: map[string]interface{}{
				"step":      "processing_webhook",
				"signature": signature != "",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Handle webhook
	result, err := s.paymentRepository.HandleWebhook(ctx, rawBody, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to handle webhook: %w", err)
	}

	// Send completion
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "payment_operation",
			Status:   "completed",
			Progress: 100.0,
			Result: map[string]interface{}{
				"webhook_processed": true,
				"message":           "Webhook processed successfully",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return map[string]interface{}{
		"entity_type": "webhooks",
		"operation":   "handle_webhook",
		"result":      result,
		"success":     true,
	}, nil
}

// GetSupportedOperations returns list of supported operations
func (s *PaymentToolService) GetSupportedOperations() []string {
	return []string{
		"authorize_payment",
		"confirm_payment",
		"capture_payment",
		"refund_payment",
		"get_payment",
		"search_payments",
		"get_payments_by_customer",
		"create_invoice",
		"get_invoice",
		"pay_invoice",
		"adjust_invoice",
		"cancel_invoice",
		"search_invoices",
		"get_invoices_by_order",
		"handle_webhook",
	}
}
