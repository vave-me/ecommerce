package grpc

import (
	"context"
	"fmt"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
	"middleman/payments/paymentspb"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// PaymentRepository calls the remote payments service (gRPC).
type PaymentRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.PaymentRepository = (*PaymentRepository)(nil)

// NewPaymentRepositoryWithAuth creates a new PaymentRepository with JWT authentication support
func NewPaymentRepository(endpoint string, authInstance *auth.Auth) PaymentRepository {
	return PaymentRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// AuthorizePayment authorizes a payment for the given amount and customer
func (r PaymentRepository) AuthorizePayment(ctx context.Context, userCustomerID string, amount int64) (*models.AuthorizePaymentResponse, error) {
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

	return &models.AuthorizePaymentResponse{
		ID:           resp.GetId(),
		ClientSecret: resp.GetClientSecret(),
	}, nil
}

// ConfirmPayment confirms a previously authorized payment using a payment method ID
func (r PaymentRepository) ConfirmPayment(ctx context.Context, paymentID, paymentMethodID string) (*models.ConfirmPaymentResponse, error) {
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

	return &models.ConfirmPaymentResponse{
		ID:            resp.GetId(),
		PaymentStatus: resp.GetPaymentStatus(),
	}, nil
}

// CapturePayment captures a payment, possibly partially
func (r PaymentRepository) CapturePayment(ctx context.Context, paymentID string, amountToCapture int64) (*models.CapturePaymentResponse, error) {
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

	return &models.CapturePaymentResponse{
		PaymentID:     resp.GetPaymentId(),
		PaymentStatus: resp.GetPaymentStatus(),
	}, nil
}

// CreateInvoice creates a new invoice
func (r PaymentRepository) CreateInvoice(ctx context.Context, orderID, paymentID string, amount int64) (*models.CreateInvoiceResponse, error) {
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

	return &models.CreateInvoiceResponse{
		ID:        resp.GetId(),
		CreatedAt: createdAt,
	}, nil
}

// AdjustInvoice adjusts an existing invoice (e.g., changing the amount)
func (r PaymentRepository) AdjustInvoice(ctx context.Context, invoiceID string, amount int64, reason string) (*models.AdjustInvoiceResponse, error) {
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

	return &models.AdjustInvoiceResponse{
		ID:         resp.GetId(),
		AdjustedAt: adjustedAt,
		NewAmount:  resp.GetNewAmount(),
	}, nil
}

// PayInvoice pays an existing invoice with the specified payment method
func (r PaymentRepository) PayInvoice(ctx context.Context, invoiceID, paymentMethodID string) (*models.PayInvoiceResponse, error) {
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

	return &models.PayInvoiceResponse{
		ID:            resp.GetId(),
		PaidAt:        paidAt,
		PaymentStatus: resp.GetPaymentStatus(),
	}, nil
}

// CancelInvoice cancels an existing invoice
func (r PaymentRepository) CancelInvoice(ctx context.Context, invoiceID, reason string) error {
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
func (r PaymentRepository) HandleWebhook(ctx context.Context, rawBody, signature string) (*models.HandleWebhookResponse, error) {
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

	return &models.HandleWebhookResponse{
		Success: resp.GetSuccess(),
		Message: resp.GetMessage(),
	}, nil
}

// GetPayment retrieves a payment by ID (mock implementation for AI tooling)
func (r PaymentRepository) GetPayment(ctx context.Context, paymentID string) (*models.Payment, error) {
	// Note: This would typically require a GetPayment RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetPayment called for ID: %s (mock implementation)", paymentID)

	return &models.Payment{
		ID:              paymentID,
		UserCustomerID:  "mock_customer",
		Amount:          1000,
		PaymentMethodID: "mock_payment_method",
		Status:          models.PaymentStatusPending,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

// GetInvoice retrieves an invoice by ID (mock implementation for AI tooling)
func (r PaymentRepository) GetInvoice(ctx context.Context, invoiceID string) (*models.Invoice, error) {
	// Note: This would typically require a GetInvoice RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetInvoice called for ID: %s (mock implementation)", invoiceID)

	return &models.Invoice{
		ID:            invoiceID,
		OrderID:       "mock_order",
		PaymentID:     "mock_payment",
		Amount:        1000,
		InvoiceStatus: models.InvoiceStatusDraft,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// GetPaymentsByCustomer retrieves payments by customer ID (mock implementation for AI tooling)
func (r PaymentRepository) GetPaymentsByCustomer(ctx context.Context, userCustomerID string) ([]*models.Payment, error) {
	// Note: This would typically require a GetPaymentsByCustomer RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetPaymentsByCustomer called for customer: %s (mock implementation)", userCustomerID)

	return []*models.Payment{
		{
			ID:              "payment_1",
			UserCustomerID:  userCustomerID,
			Amount:          1000,
			PaymentMethodID: "mock_payment_method",
			Status:          models.PaymentStatusPending,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}, nil
}

// GetInvoicesByOrder retrieves invoices by order ID (mock implementation for AI tooling)
func (r PaymentRepository) GetInvoicesByOrder(ctx context.Context, orderID string) ([]*models.Invoice, error) {
	// Note: This would typically require a GetInvoicesByOrder RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetInvoicesByOrder called for order: %s (mock implementation)", orderID)

	return []*models.Invoice{
		{
			ID:            "invoice_1",
			OrderID:       orderID,
			PaymentID:     "mock_payment",
			Amount:        1000,
			InvoiceStatus: models.InvoiceStatusDraft,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}, nil
}

// SearchPayments searches payments by status (mock implementation for AI tooling)
func (r PaymentRepository) SearchPayments(ctx context.Context, status string, limit int64) ([]*models.Payment, error) {
	// Note: This would typically require a SearchPayments RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("SearchPayments called with status: %s, limit: %d (mock implementation)", status, limit)

	payments := make([]*models.Payment, 0, limit)
	for i := int64(0); i < limit && i < 5; i++ { // Mock max 5 results
		payments = append(payments, &models.Payment{
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
func (r PaymentRepository) SearchInvoices(ctx context.Context, status string, limit int64) ([]*models.Invoice, error) {
	// Note: This would typically require a SearchInvoices RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("SearchInvoices called with status: %s, limit: %d (mock implementation)", status, limit)

	invoices := make([]*models.Invoice, 0, limit)
	for i := int64(0); i < limit && i < 5; i++ { // Mock max 5 results
		invoices = append(invoices, &models.Invoice{
			ID:            fmt.Sprintf("invoice_%d", i+1),
			OrderID:       fmt.Sprintf("order_%d", i+1),
			PaymentID:     fmt.Sprintf("payment_%d", i+1),
			Amount:        1000 * (i + 1),
			InvoiceStatus: status,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}

	return invoices, nil
}

// dial establishes a gRPC connection to the payments service
// dial sets up a gRPC connection with the microservice endpoint
func (r PaymentRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r PaymentRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
