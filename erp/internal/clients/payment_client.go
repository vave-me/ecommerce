package clients

import (
	"context"
	"fmt"
	"middleman/erp/internal/domain"

	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/payments/paymentspb"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// PaymentClient calls the remote payments service (gRPC).
type PaymentClient struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.PaymentRepository = (*PaymentClient)(nil)

// NewPaymentClient creates a new PaymentClient with JWT authentication support
func NewPaymentClient(endpoint string, authInstance *auth.Auth) PaymentClient {
	return PaymentClient{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// AuthorizePayment authorizes a payment for the given amount and customer
func (r PaymentClient) AuthorizePayment(ctx context.Context, userCustomerID string, amount int64) (*domain.AuthorizePaymentResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := paymentspb.NewPaymentsServiceClient(conn)
	resp, err := client.AuthorizePayment(ctx, &paymentspb.AuthorizePaymentRequest{
		UserCustomerId: userCustomerID,
		Amount:         amount,
	})
	if err != nil {
		return nil, fmt.Errorf("AuthorizePayment RPC failed: %w", err)
	}

	return &domain.AuthorizePaymentResponse{
		ID:           resp.GetId(),
		ClientSecret: resp.GetClientSecret(),
	}, nil
}

// ConfirmPayment confirms a previously authorized payment using a payment method ID
func (r PaymentClient) ConfirmPayment(ctx context.Context, paymentID, paymentMethodID string) (*domain.ConfirmPaymentResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := paymentspb.NewPaymentsServiceClient(conn)
	resp, err := client.ConfirmPayment(ctx, &paymentspb.ConfirmPaymentRequest{
		Id:              paymentID,
		PaymentMethodId: paymentMethodID,
	})
	if err != nil {
		return nil, fmt.Errorf("ConfirmPayment RPC failed: %w", err)
	}

	return &domain.ConfirmPaymentResponse{
		ID:            resp.GetId(),
		PaymentStatus: resp.GetPaymentStatus(),
	}, nil
}

// CapturePayment captures a payment, possibly partially
func (r PaymentClient) CapturePayment(ctx context.Context, paymentID string, amountToCapture int64) (*domain.CapturePaymentResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := paymentspb.NewPaymentsServiceClient(conn)
	resp, err := client.CapturePayment(ctx, &paymentspb.CapturePaymentRequest{
		PaymentId:       paymentID,
		AmountToCapture: amountToCapture,
	})
	if err != nil {
		return nil, fmt.Errorf("CapturePayment RPC failed: %w", err)
	}

	return &domain.CapturePaymentResponse{
		PaymentID:     resp.GetPaymentId(),
		PaymentStatus: resp.GetPaymentStatus(),
	}, nil
}

// CreateInvoice creates a new invoice
func (r PaymentClient) CreateInvoice(ctx context.Context, orderID, paymentID string, amount int64) (*domain.CreateInvoiceResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := paymentspb.NewPaymentsServiceClient(conn)
	resp, err := client.CreateInvoice(ctx, &paymentspb.CreateInvoiceRequest{
		OrderId:   orderID,
		PaymentId: paymentID,
		Amount:    amount,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateInvoice RPC failed: %w", err)
	}

	var createdAt time.Time
	if resp.GetCreatedAt() != nil {
		createdAt = resp.GetCreatedAt().AsTime()
	}

	return &domain.CreateInvoiceResponse{
		ID:        resp.GetId(),
		CreatedAt: createdAt,
	}, nil
}

// AdjustInvoice adjusts an existing invoice (e.g., changing the amount)
func (r PaymentClient) AdjustInvoice(ctx context.Context, invoiceID string, amount int64, reason string) (*domain.AdjustInvoiceResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := paymentspb.NewPaymentsServiceClient(conn)
	resp, err := client.AdjustInvoice(ctx, &paymentspb.AdjustInvoiceRequest{
		Id:     invoiceID,
		Amount: amount,
		Reason: reason,
	})
	if err != nil {
		return nil, fmt.Errorf("AdjustInvoice RPC failed: %w", err)
	}

	var adjustedAt time.Time
	if resp.GetAdjustedAt() != nil {
		adjustedAt = resp.GetAdjustedAt().AsTime()
	}

	return &domain.AdjustInvoiceResponse{
		ID:         resp.GetId(),
		AdjustedAt: adjustedAt,
		NewAmount:  resp.GetNewAmount(),
	}, nil
}

// PayInvoice pays an existing invoice with the specified payment method
func (r PaymentClient) PayInvoice(ctx context.Context, invoiceID, paymentMethodID string) (*domain.PayInvoiceResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := paymentspb.NewPaymentsServiceClient(conn)
	resp, err := client.PayInvoice(ctx, &paymentspb.PayInvoiceRequest{
		Id:              invoiceID,
		PaymentMethodId: paymentMethodID,
	})
	if err != nil {
		return nil, fmt.Errorf("PayInvoice RPC failed: %w", err)
	}

	var paidAt time.Time
	if resp.GetPaidAt() != nil {
		paidAt = resp.GetPaidAt().AsTime()
	}

	return &domain.PayInvoiceResponse{
		ID:            resp.GetId(),
		PaidAt:        paidAt,
		PaymentStatus: resp.GetPaymentStatus(),
	}, nil
}

// CancelInvoice cancels an existing invoice
func (r PaymentClient) CancelInvoice(ctx context.Context, invoiceID, reason string) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := paymentspb.NewPaymentsServiceClient(conn)
	_, err = client.CancelInvoice(ctx, &paymentspb.CancelInvoiceRequest{
		Id:     invoiceID,
		Reason: reason,
	})
	if err != nil {
		return fmt.Errorf("CancelInvoice RPC failed: %w", err)
	}

	return nil
}

// HandleWebhook handles a webhook event (e.g., from Stripe)
func (r PaymentClient) HandleWebhook(ctx context.Context, rawBody, signature string) (*domain.HandleWebhookResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := paymentspb.NewPaymentsServiceClient(conn)
	resp, err := client.HandleWebhook(ctx, &paymentspb.HandleWebhookRequest{
		RawBody:   rawBody,
		Signature: signature,
	})
	if err != nil {
		return nil, fmt.Errorf("HandleWebhook RPC failed: %w", err)
	}

	return &domain.HandleWebhookResponse{
		Success: resp.GetSuccess(),
		Message: resp.GetMessage(),
	}, nil
}

// GetPayment retrieves a payment by ID (mock implementation for AI tooling)
func (r PaymentClient) GetPayment(ctx context.Context, paymentID string) (*domain.Payment, error) {
	// Note: This would typically require a GetPayment RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetPayment called for ID: %s (mock implementation)", paymentID)

	return &domain.Payment{
		ID:              paymentID,
		UserCustomerID:  "mock_customer",
		Amount:          1000,
		PaymentMethodID: "mock_payment_method",
		Status:          domain.PaymentStatusPending,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

// GetInvoice retrieves an invoice by ID (mock implementation for AI tooling)
func (r PaymentClient) GetInvoice(ctx context.Context, invoiceID string) (*domain.Invoice, error) {
	// Note: This would typically require a GetInvoice RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetInvoice called for ID: %s (mock implementation)", invoiceID)

	invoice := &domain.Invoice{
		InvoiceNumber: invoiceID,
		OrderID:       "mock_order",
		CustomerID:    "mock_customer",
		Type:          domain.InvoiceTypeStandard,
		Status:        domain.InvoiceStatusDraft,
		IssueDate:     time.Now(),
		DueDate:       time.Now().Add(30 * 24 * time.Hour),
		Currency:      "USD",
		TotalAmount:   1000,
		BalanceDue:    1000,
	}
	// Set the aggregate ID
	invoice.Aggregate = es.NewAggregate(invoiceID, domain.InvoiceAggregate)
	return invoice, nil
}

// GetPaymentsByCustomer retrieves payments by customer ID (mock implementation for AI tooling)
func (r PaymentClient) GetPaymentsByCustomer(ctx context.Context, userCustomerID string) ([]*domain.Payment, error) {
	// Note: This would typically require a GetPaymentsByCustomer RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetPaymentsByCustomer called for customer: %s (mock implementation)", userCustomerID)

	return []*domain.Payment{
		{
			ID:              "payment_1",
			UserCustomerID:  userCustomerID,
			Amount:          1000,
			PaymentMethodID: "mock_payment_method",
			Status:          domain.PaymentStatusPending,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}, nil
}

// GetInvoicesByOrder retrieves invoices by order ID (mock implementation for AI tooling)
func (r PaymentClient) GetInvoicesByOrder(ctx context.Context, orderID string) ([]*domain.Invoice, error) {
	// Note: This would typically require a GetInvoicesByOrder RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetInvoicesByOrder called for order: %s (mock implementation)", orderID)

	invoice := &domain.Invoice{
		InvoiceNumber: "invoice_1",
		OrderID:       orderID,
		CustomerID:    "mock_customer",
		Type:          domain.InvoiceTypeStandard,
		Status:        domain.InvoiceStatusDraft,
		IssueDate:     time.Now(),
		DueDate:       time.Now().Add(30 * 24 * time.Hour),
		Currency:      "USD",
		TotalAmount:   1000,
		BalanceDue:    1000,
	}
	invoice.Aggregate = es.NewAggregate("invoice_1", domain.InvoiceAggregate)
	
	return []*domain.Invoice{invoice}, nil
}

// SearchPayments searches payments by status (mock implementation for AI tooling)
func (r PaymentClient) SearchPayments(ctx context.Context, status string, limit int64) ([]*domain.Payment, error) {
	// Note: This would typically require a SearchPayments RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("SearchPayments called with status: %s, limit: %d (mock implementation)", status, limit)

	payments := make([]*domain.Payment, 0, limit)
	for i := int64(0); i < limit && i < 5; i++ { // Mock max 5 results
		payments = append(payments, &domain.Payment{
			ID:              fmt.Sprintf("payment_%d", i+1),
			UserCustomerID:  "mock_customer",
			Amount:          1000 * (i + 1),
			PaymentMethodID: "mock_payment_method",
			Status:          status,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		})
	}

	return payments, nil
}

// SearchInvoices searches invoices by status (mock implementation for AI tooling)
func (r PaymentClient) SearchInvoices(ctx context.Context, status string, limit int64) ([]*domain.Invoice, error) {
	// Note: This would typically require a SearchInvoices RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("SearchInvoices called with status: %s, limit: %d (mock implementation)", status, limit)

	invoices := make([]*domain.Invoice, 0, limit)
	for i := int64(0); i < limit && i < 5; i++ { // Mock max 5 results
		invoice := &domain.Invoice{
			InvoiceNumber: fmt.Sprintf("invoice_%d", i+1),
			OrderID:       fmt.Sprintf("order_%d", i+1),
			CustomerID:    fmt.Sprintf("customer_%d", i+1),
			Type:          domain.InvoiceTypeStandard,
			Status:        domain.InvoiceStatus(status),
			IssueDate:     time.Now(),
			DueDate:       time.Now().Add(30 * 24 * time.Hour),
			Currency:      "USD",
			TotalAmount:   float64(1000 * (i + 1)),
			BalanceDue:    float64(1000 * (i + 1)),
		}
		invoice.Aggregate = es.NewAggregate(fmt.Sprintf("invoice_%d", i+1), domain.InvoiceAggregate)
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

// dial establishes a gRPC connection to the payments service
// dial sets up a gRPC connection with the microservice endpoint
func (r PaymentClient) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r PaymentClient) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
