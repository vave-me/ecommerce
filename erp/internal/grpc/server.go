package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"middleman/erp/erppb"
	"middleman/erp/internal/application"
	"middleman/erp/internal/application/commands"
	"middleman/erp/internal/application/queries"
	"middleman/erp/internal/domain"
	"middleman/internal/auth"
	"middleman/internal/erp"
)

type server struct {
	app application.App
	erppb.UnimplementedERPServiceServer
}

var _ erppb.ERPServiceServer = (*server)(nil)

// RegisterServer registers the gRPC server implementation
func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	erppb.RegisterERPServiceServer(registrar, server{app: app})
	return nil
}

// Connector Management Methods

// AddConnector creates a new ERP connector with secure credential storage
func (s server) AddConnector(ctx context.Context, req *erppb.AddConnectorRequest) (*erppb.AddConnectorResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConnectorName", req.GetName()),
		attribute.String("ConnectorType", req.GetType()),
	)

	// Validate authentication
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "missing authentication")
	}

	// Convert auth config
	authConfig := make(map[string]interface{})
	if req.AuthConfig != nil {
		authConfig = req.AuthConfig.AsMap()
	}
	authConfig["type"] = req.AuthType

	// Convert sync entities
	var syncEntities []commands.SyncEntityConfig
	for _, entity := range req.SyncEntities {
		syncEntities = append(syncEntities, commands.SyncEntityConfig{
			EntityType:    entity.EntityType,
			Enabled:       entity.Enabled,
			SyncDirection: entity.SyncDirection,
			Filters:       entity.Filters.AsMap(),
			FieldMapping:  entity.FieldMapping,
		})
	}

	// Execute command
	err := s.app.AddConnector(ctx, commands.AddConnector{
		Name:                       req.Name,
		Type:                       req.Type,
		Environment:                req.Environment,
		BaseURL:                    req.BaseUrl,
		AuthType:                   req.AuthType,
		AuthConfig:                 authConfig,
		WebhookEnabled:             req.WebhookEnabled,
		WebhookSecret:              req.WebhookSecret,
		WebhookEvents:              req.WebhookEvents,
		SyncEnabled:                req.SyncEnabled,
		SyncIntervalSeconds:        int(req.SyncIntervalSeconds),
		BatchSize:                  int(req.BatchSize),
		SyncEntities:               syncEntities,
		RateLimitRequestsPerSecond: int(req.RateLimitRequestsPerSecond),
		RateLimitBurst:             int(req.RateLimitBurst),
		RetryMaxAttempts:           int(req.RetryMaxAttempts),
		RetryInitialDelayMs:        int(req.RetryInitialDelayMs),
		RetryMaxDelayMs:            int(req.RetryMaxDelayMs),
		RetryMultiplier:            req.RetryMultiplier,
		CustomHeaders:              req.CustomHeaders,
		TimeoutSeconds:             int(req.TimeoutSeconds),
		CreatedBy:                  claims.Subject,
	})

	if err != nil {
		return nil, status.Errorf(grpc_code.Internal, "failed to add connector: %v", err)
	}

	// Get the created connector
	connectorStatus, err := s.app.GetConnectorStatus(ctx, queries.GetConnectorStatus{
		ConnectorID: req.Name, // Using name temporarily, should use returned ID
	})

	if err != nil {
		return nil, status.Errorf(grpc_code.Internal, "failed to get connector status: %v", err)
	}

	return &erppb.AddConnectorResponse{
		ConnectorId: connectorStatus.ConnectorID,
		Connector:   connectorToProto(connectorStatus),
	}, nil
}

// UpdateConnector updates an existing ERP connector configuration
func (s server) UpdateConnector(ctx context.Context, req *erppb.UpdateConnectorRequest) (*erppb.UpdateConnectorResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConnectorID", req.GetConnectorId()),
	)

	// Validate authentication
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "missing authentication")
	}

	// Convert optional fields
	cmd := commands.UpdateConnector{
		ConnectorID:   req.ConnectorId,
		Name:          req.Name,
		Environment:   req.Environment,
		BaseURL:       req.BaseUrl,
		CustomHeaders: req.CustomHeaders,
		UpdatedBy:     claims.Subject,
	}

	// Handle optional auth config
	if req.AuthConfig != nil && len(req.AuthConfig.AsMap()) > 0 {
		authConfig := req.AuthConfig.AsMap()
		cmd.AuthConfig = &authConfig
	}

	// Handle optional webhook fields
	if req.WebhookEnabled != nil {
		enabled := req.WebhookEnabled.GetBoolValue()
		cmd.WebhookEnabled = &enabled
	}
	if req.WebhookSecret != nil {
		secret := req.WebhookSecret.GetStringValue()
		cmd.WebhookSecret = &secret
	}
	if len(req.WebhookEvents) > 0 {
		cmd.WebhookEvents = req.WebhookEvents
	}

	// Handle optional sync fields
	if req.SyncEnabled != nil {
		enabled := req.SyncEnabled.GetBoolValue()
		cmd.SyncEnabled = &enabled
	}
	if req.SyncIntervalSeconds != nil {
		interval := int(req.SyncIntervalSeconds.GetNumberValue())
		cmd.SyncIntervalSeconds = &interval
	}
	if req.BatchSize != nil {
		batch := int(req.BatchSize.GetNumberValue())
		cmd.BatchSize = &batch
	}

	// Handle optional rate limit fields
	if req.RateLimitRequestsPerSecond != nil {
		limit := int(req.RateLimitRequestsPerSecond.GetNumberValue())
		cmd.RateLimitRequestsPerSecond = &limit
	}
	if req.RateLimitBurst != nil {
		burst := int(req.RateLimitBurst.GetNumberValue())
		cmd.RateLimitBurst = &burst
	}

	// Handle optional retry fields
	if req.RetryMaxAttempts != nil {
		attempts := int(req.RetryMaxAttempts.GetNumberValue())
		cmd.RetryMaxAttempts = &attempts
	}
	if req.RetryInitialDelayMs != nil {
		delay := int(req.RetryInitialDelayMs.GetNumberValue())
		cmd.RetryInitialDelayMs = &delay
	}
	if req.RetryMaxDelayMs != nil {
		delay := int(req.RetryMaxDelayMs.GetNumberValue())
		cmd.RetryMaxDelayMs = &delay
	}
	if req.RetryMultiplier != nil {
		multiplier := req.RetryMultiplier.GetNumberValue()
		cmd.RetryMultiplier = &multiplier
	}

	// Handle optional timeout
	if req.TimeoutSeconds != nil {
		timeout := int(req.TimeoutSeconds.GetNumberValue())
		cmd.TimeoutSeconds = &timeout
	}

	// Execute command
	err := s.app.UpdateConnector(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(grpc_code.Internal, "failed to update connector: %v", err)
	}

	// Get the updated connector
	connectorStatus, err := s.app.GetConnectorStatus(ctx, queries.GetConnectorStatus{
		ConnectorID: req.ConnectorId,
	})

	if err != nil {
		return nil, status.Errorf(grpc_code.Internal, "failed to get connector status: %v", err)
	}

	return &erppb.UpdateConnectorResponse{
		Connector: connectorToProto(connectorStatus),
	}, nil
}

// RemoveConnector removes an ERP connector and its associated data
func (s server) RemoveConnector(ctx context.Context, req *erppb.RemoveConnectorRequest) (*erppb.RemoveConnectorResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConnectorID", req.GetConnectorId()),
		attribute.Bool("Force", req.GetForce()),
	)

	// Validate authentication
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "missing authentication")
	}

	// Execute command
	err := s.app.RemoveConnector(ctx, commands.RemoveConnector{
		ConnectorID: req.ConnectorId,
		RemoveBy:    claims.Subject,
		Force:       req.Force,
	})

	if err != nil {
		return &erppb.RemoveConnectorResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &erppb.RemoveConnectorResponse{
		Success: true,
		Message: "Connector removed successfully",
	}, nil
}

// ToggleConnector activates or deactivates a connector
func (s server) ToggleConnector(ctx context.Context, req *erppb.ToggleConnectorRequest) (*erppb.ToggleConnectorResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConnectorID", req.GetConnectorId()),
		attribute.Bool("Activate", req.GetActivate()),
	)

	// Validate authentication
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "missing authentication")
	}

	// Get current status
	connectorStatus, err := s.app.GetConnectorStatus(ctx, queries.GetConnectorStatus{
		ConnectorID: req.ConnectorId,
	})
	if err != nil {
		return nil, status.Errorf(grpc_code.NotFound, "connector not found: %v", err)
	}

	previousStatus := connectorStatus.Status

	// Execute command
	err = s.app.ToggleConnector(ctx, commands.ToggleConnector{
		ConnectorID: req.ConnectorId,
		Activate:    req.Activate,
		ChangedBy:   claims.Subject,
		Reason:      req.Reason,
	})

	if err != nil {
		return nil, status.Errorf(grpc_code.Internal, "failed to toggle connector: %v", err)
	}

	// Get updated connector
	updatedStatus, err := s.app.GetConnectorStatus(ctx, queries.GetConnectorStatus{
		ConnectorID: req.ConnectorId,
	})
	if err != nil {
		return nil, status.Errorf(grpc_code.Internal, "failed to get updated connector status: %v", err)
	}

	return &erppb.ToggleConnectorResponse{
		Connector:      connectorToProto(updatedStatus),
		PreviousStatus: string(previousStatus),
		NewStatus:      string(updatedStatus.Status),
	}, nil
}

// Invoice Granular Commands

// CreateInvoice creates a new invoice
func (s server) CreateInvoice(ctx context.Context, req *erppb.CreateInvoiceRequest) (*erppb.CreateInvoiceResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("InvoiceID", req.GetInvoiceId()),
		attribute.String("OrderID", req.GetOrderId()),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Convert invoice lines
	lines := make([]domain.InvoiceLine, len(req.GetLines()))
	for i, line := range req.GetLines() {
		lines[i] = domain.InvoiceLine{
			SKU:         line.GetSku(),
			ProductName: line.GetDescription(), // Map description to ProductName
			Description: line.GetDescription(),
			Quantity:    int(line.GetQuantity()),
			UnitPrice:   float64(line.GetUnitPrice()),
			LineTotal:   float64(line.GetLineTotal()),
			// Set defaults for other fields
			TaxRate:     0,
			TaxAmount:   0,
			Discount:    0,
			AccountCode: "",
		}
	}

	// Create the invoice
	err := s.app.CreateInvoice(ctx, commands.CreateInvoice{
		InvoiceID:      req.GetInvoiceId(),
		InvoiceNumber:  req.GetInvoiceNumber(),
		OrderID:        req.GetOrderId(),
		CustomerID:     req.GetCustomerId(),
		Type:           domain.InvoiceType(req.GetType()),
		IssueDate:      req.GetIssueDate().AsTime(),
		DueDate:        req.GetDueDate().AsTime(),
		Currency:       req.GetCurrency(),
		Lines:          lines,
		SubTotal:       float64(req.GetSubTotal()),
		TaxAmount:      float64(req.GetTaxAmount()),
		DiscountAmount: 0, // Default value, not in proto
		ShippingAmount: 0, // Default value, not in proto
		TotalAmount:    float64(req.GetTotalAmount()),
		PaymentTerms:   req.GetPaymentTerms(),
		BillingAddress: erp.Address{}, // Default empty address
		TaxID:          "",            // Default empty
		PONumber:       "",            // Default empty
		Notes:          req.GetNotes(),
		ConnectorID:    req.GetConnectorId(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.CreateInvoiceResponse{
		InvoiceId:     req.GetInvoiceId(),
		InvoiceNumber: req.GetInvoiceNumber(),
		Status:        "created",
	}, nil
}

// ApproveInvoice approves an invoice
func (s server) ApproveInvoice(ctx context.Context, req *erppb.ApproveInvoiceRequest) (*erppb.ApproveInvoiceResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("InvoiceID", req.GetInvoiceId()))

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.ApproveInvoice(ctx, commands.ApproveInvoice{
		InvoiceID:  req.GetInvoiceId(),
		ApprovedBy: req.GetApprovedBy(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.ApproveInvoiceResponse{
		InvoiceId:  req.GetInvoiceId(),
		Status:     "approved",
		ApprovedAt: timestamppb.Now(),
	}, nil
}

// SendInvoice sends an invoice
func (s server) SendInvoice(ctx context.Context, req *erppb.SendInvoiceRequest) (*erppb.SendInvoiceResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("InvoiceID", req.GetInvoiceId()))

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.SendInvoice(ctx, commands.SendInvoice{
		InvoiceID: req.GetInvoiceId(),
		SentTo:    []string{req.GetRecipientEmail()},
		SentBy:    "system", // Default sender
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.SendInvoiceResponse{
		InvoiceId: req.GetInvoiceId(),
		Status:    "sent",
		SentAt:    timestamppb.Now(),
	}, nil
}

// VoidInvoice voids an invoice
func (s server) VoidInvoice(ctx context.Context, req *erppb.VoidInvoiceRequest) (*erppb.VoidInvoiceResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("InvoiceID", req.GetInvoiceId()))

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.VoidInvoice(ctx, commands.VoidInvoice{
		InvoiceID: req.GetInvoiceId(),
		Reason:    req.GetReason(),
		VoidedBy:  req.GetVoidedBy(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.VoidInvoiceResponse{
		InvoiceId: req.GetInvoiceId(),
		Status:    "voided",
		VoidedAt:  timestamppb.Now(),
	}, nil
}

// RecordInvoicePayment records a payment against an invoice
func (s server) RecordInvoicePayment(ctx context.Context, req *erppb.RecordInvoicePaymentRequest) (*erppb.RecordInvoicePaymentResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("InvoiceID", req.GetInvoiceId()))

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.RecordInvoicePayment(ctx, commands.RecordInvoicePayment{
		InvoiceID:     req.GetInvoiceId(),
		Amount:        float64(req.GetAmount()),
		PaymentMethod: req.GetPaymentMethod(),
		TransactionID: req.GetTransactionId(),
		PaymentDate:   req.GetPaymentDate().AsTime(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.RecordInvoicePaymentResponse{
		InvoiceId:        req.GetInvoiceId(),
		PaymentId:        uuid.New().String(),
		RemainingBalance: 0, // This would come from the aggregate
		PaymentStatus:    "paid",
	}, nil
}

// ProcessInvoice processes an invoice in the ERP system (DEPRECATED)
// This method is deprecated and not implemented in the application layer.
// Use the granular invoice commands instead (CreateInvoice, ApproveInvoice, etc.)
func (s server) ProcessInvoice(ctx context.Context, req *erppb.ProcessInvoiceRequest) (*erppb.ProcessInvoiceResponse, error) {
	return nil, status.Error(grpc_code.Unimplemented, "ProcessInvoice is deprecated, use granular invoice commands")
}

// Return Granular Commands

// CreateReturn creates a new return
func (s server) CreateReturn(ctx context.Context, req *erppb.CreateReturnRequest) (*erppb.CreateReturnResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ReturnID", req.GetReturnId()),
		attribute.String("OrderID", req.GetOriginalOrderId()),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Convert return items
	items := make([]domain.ReturnItem, len(req.GetItems()))
	for i, item := range req.GetItems() {
		items[i] = domain.ReturnItem{
			SKU:             item.GetSku(),
			ProductName:     item.GetProductName(),
			Quantity:        int(item.GetQuantity()),
			Condition:       domain.ItemCondition(item.GetCondition()),
			RestockingFee:   float64(item.GetRestockingFee()),
			RefundAmount:    float64(item.GetRefundAmount()),
			ExchangeForSKU:  item.GetExchangeForSku(),
			SerialNumbers:   item.GetSerialNumbers(),
			InspectionNotes: item.GetInspectionNotes(),
		}
	}

	err := s.app.CreateReturn(ctx, commands.CreateReturn{
		ReturnID:        req.GetReturnId(),
		ReturnNumber:    req.GetReturnNumber(),
		OriginalOrderID: req.GetOriginalOrderId(),
		CustomerID:      req.GetCustomerId(),
		CustomerName:    req.GetCustomerName(),
		CustomerEmail:   req.GetCustomerEmail(),
		Reason:          domain.ReturnReason(req.GetReason()),
		Items:           items,
		RefundMethod:    domain.RefundMethod(req.GetRefundMethod()),
		RefundAmount:    float64(req.GetRefundAmount()),
		WarehouseID:     req.GetWarehouseId(),
		Notes:           req.GetNotes(),
		ConnectorID:     req.GetConnectorId(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.CreateReturnResponse{
		ReturnId:     req.GetReturnId(),
		ReturnNumber: req.GetReturnNumber(),
		Status:       "created",
	}, nil
}

// ApproveReturn approves a return
func (s server) ApproveReturn(ctx context.Context, req *erppb.ApproveReturnRequest) (*erppb.ApproveReturnResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ReturnID", req.GetReturnId()))

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.ApproveReturn(ctx, commands.ApproveReturn{
		ReturnID:   req.GetReturnId(),
		ApprovedBy: req.GetApprovedBy(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.ApproveReturnResponse{
		ReturnId:   req.GetReturnId(),
		Status:     "approved",
		ApprovedAt: timestamppb.Now(),
	}, nil
}

// RejectReturn rejects a return
func (s server) RejectReturn(ctx context.Context, req *erppb.RejectReturnRequest) (*erppb.RejectReturnResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ReturnID", req.GetReturnId()))

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.RejectReturn(ctx, commands.RejectReturn{
		ReturnID: req.GetReturnId(),
		Reason:   req.GetReason(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.RejectReturnResponse{
		ReturnId:   req.GetReturnId(),
		Status:     "rejected",
		RejectedAt: timestamppb.Now(),
	}, nil
}

// ProcessReturnStart starts processing a return
func (s server) ProcessReturnStart(ctx context.Context, req *erppb.ProcessReturnStartRequest) (*erppb.ProcessReturnStartResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ReturnID", req.GetReturnId()))

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.ProcessReturnStart(ctx, commands.ProcessReturnStart{
		ReturnID:    req.GetReturnId(),
		ERPReturnID: req.GetErpReturnId(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.ProcessReturnStartResponse{
		ReturnId:    req.GetReturnId(),
		Status:      "processing",
		ProcessedAt: timestamppb.Now(),
	}, nil
}

// CompleteReturn completes a return
func (s server) CompleteReturn(ctx context.Context, req *erppb.CompleteReturnRequest) (*erppb.CompleteReturnResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ReturnID", req.GetReturnId()))

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.CompleteReturn(ctx, commands.CompleteReturn{
		ReturnID:            req.GetReturnId(),
		RefundTransactionID: req.GetRefundTransactionId(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.CompleteReturnResponse{
		ReturnId:    req.GetReturnId(),
		Status:      "completed",
		CompletedAt: timestamppb.Now(),
	}, nil
}

// RestockReturnItems restocks items from a return
func (s server) RestockReturnItems(ctx context.Context, req *erppb.RestockReturnItemsRequest) (*erppb.RestockReturnItemsResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ReturnID", req.GetReturnId()))

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Convert restocked items
	items := make([]domain.RestockedItem, len(req.GetItems()))
	for i, item := range req.GetItems() {
		items[i] = domain.RestockedItem{
			SKU:        item.GetSku(),
			Quantity:   int(item.GetQuantity()),
			LocationID: item.GetLocationId(),
		}
	}

	err := s.app.RestockReturnItems(ctx, commands.RestockReturnItems{
		ReturnID: req.GetReturnId(),
		Items:    items,
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.RestockReturnItemsResponse{
		ReturnId:       req.GetReturnId(),
		ItemsRestocked: int32(len(items)),
		RestockedAt:    timestamppb.Now(),
	}, nil
}

// ProcessReturn processes a return in the ERP system (DEPRECATED)
// This method is deprecated and not implemented in the application layer.
// Use the granular return commands instead (CreateReturn, ApproveReturn, etc.)
func (s server) ProcessReturn(ctx context.Context, req *erppb.ProcessReturnRequest) (*erppb.ProcessReturnResponse, error) {
	return nil, status.Error(grpc_code.Unimplemented, "ProcessReturn is deprecated, use granular return commands")
}

// ProcessWebhook is deprecated - webhooks are handled via REST endpoints
func (s server) ProcessWebhook(ctx context.Context, req *erppb.ProcessWebhookRequest) (*erppb.ProcessWebhookResponse, error) {
	return nil, status.Error(grpc_code.Unimplemented, "ProcessWebhook is deprecated, use REST webhook endpoints at /api/erp/webhook/{connectorId}")
}

// RegisterConnector registers a new ERP connector
func (s server) RegisterConnector(ctx context.Context, req *erppb.RegisterConnectorRequest) (*erppb.RegisterConnectorResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("Name", req.GetName()),
		attribute.String("Type", req.GetType()),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Generate connector ID
	connectorID := uuid.New().String()

	// Convert proto config to ERPConfig
	configMap := req.GetConfig().AsMap()
	erpConfig := erp.ERPConfig{
		Type:     erp.ERPType(req.GetType()),
		Endpoint: getStringFromConfigMap(configMap, "endpoint"),
		URL:      getStringFromConfigMap(configMap, "url"),
		Metadata: configMap,
	}

	// Parse auth config if present
	if authMap, ok := configMap["auth"].(map[string]interface{}); ok {
		erpConfig.Auth = erp.AuthConfig{
			Type:         getStringFromConfigMap(authMap, "type"),
			ClientID:     getStringFromConfigMap(authMap, "clientId"),
			ClientSecret: getStringFromConfigMap(authMap, "clientSecret"),
			APIKey:       getStringFromConfigMap(authMap, "apiKey"),
			Username:     getStringFromConfigMap(authMap, "username"),
			Password:     getStringFromConfigMap(authMap, "password"),
		}
	}

	// Parse webhook config if present
	if webhookMap, ok := configMap["webhook"].(map[string]interface{}); ok {
		erpConfig.Webhook = erp.WebhookConfig{
			Enabled: getBoolFromConfigMap(webhookMap, "enabled"),
			URL:     getStringFromConfigMap(webhookMap, "url"),
			Secret:  getStringFromConfigMap(webhookMap, "secret"),
		}
	}

	// Register the connector
	err := s.app.RegisterConnector(ctx, commands.RegisterConnector{
		ConnectorID: connectorID,
		Type:        erp.ERPType(req.GetType()),
		Config:      erpConfig,
	})
	if err != nil {
		return nil, handleError(err)
	}

	// Create response with connector details
	now := timestamppb.Now()
	return &erppb.RegisterConnectorResponse{
		ConnectorId: connectorID,
		Connector: &erppb.Connector{
			Id:        connectorID,
			Name:      req.GetName(),
			Type:      req.GetType(),
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil
}

// SendOrder sends an order to the ERP system
func (s server) SendOrder(ctx context.Context, req *erppb.SendOrderRequest) (*erppb.SendOrderResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConnectorID", req.GetConnectorId()),
		attribute.String("OrderID", req.GetOrder().GetId()),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	order := req.GetOrder()
	if order == nil {
		return nil, status.Error(grpc_code.InvalidArgument, "order is required")
	}

	// Convert order items
	items := make([]erp.OrderItem, len(order.GetItems()))
	for i, item := range order.GetItems() {
		items[i] = erp.OrderItem{
			Name:        item.GetName(),
			SKU:         item.GetSku(),
			Quantity:    int(item.GetQuantity()),
			Price:       float64(item.GetUnitPrice()) / 100.0,  // Convert cents to dollars
			TotalAmount: float64(item.GetTotalPrice()) / 100.0, // Convert cents to dollars
			TaxRate:     0,                                     // Default value
		}
	}

	// Convert to OrderPayload
	orderPayload := &erp.OrderPayload{
		OrderID:     order.GetId(),
		CustomerID:  order.GetCustomerId(),
		Items:       items,
		TotalAmount: float64(order.GetTotalAmount()) / 100.0, // Convert cents to dollars
		Currency:    order.GetCurrency(),
		Status:      "pending", // Default status
		CreatedAt:   time.Now(),
		Attributes:  order.GetMetadata().AsMap(),
	}

	// Send the order
	err := s.app.SendOrder(ctx, commands.SendOrder{
		ConnectorID: req.GetConnectorId(),
		Order:       orderPayload,
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.SendOrderResponse{
		OrderId:    order.GetId(),
		ExternalId: fmt.Sprintf("ORD-%s", order.GetId()),
		Status:     "sent",
		Result:     &structpb.Struct{},
	}, nil
}

// SyncCustomers synchronizes customers from the ERP system
func (s server) SyncCustomers(ctx context.Context, req *erppb.SyncCustomersRequest) (*erppb.SyncCustomersResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConnectorID", req.GetConnectorId()),
		attribute.Int("BatchSize", int(req.GetBatchSize())),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Parse since time
	var since time.Time
	if req.GetSince() != nil {
		since = req.GetSince().AsTime()
	}

	// Sync customers
	err := s.app.SyncCustomers(ctx, commands.SyncCustomers{
		ConnectorID:   req.GetConnectorId(),
		CustomerIDs:   []string{}, // Empty means all customers
		Since:         since,
		IncludeCredit: false, // Default to false
		BatchSize:     int(req.GetBatchSize()),
	})
	if err != nil {
		return nil, handleError(err)
	}

	// This would normally return actual synced customers
	return &erppb.SyncCustomersResponse{
		SyncId:          uuid.New().String(),
		CustomersSynced: 0,
		CustomersFailed: 0,
		Customers:       []*erppb.Customer{},
		Metadata:        &structpb.Struct{},
	}, nil
}

// SyncPrices synchronizes prices from the ERP system
func (s server) SyncPrices(ctx context.Context, req *erppb.SyncPricesRequest) (*erppb.SyncPricesResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConnectorID", req.GetConnectorId()),
		attribute.Int("BatchSize", int(req.GetBatchSize())),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Sync prices
	err := s.app.SyncPrices(ctx, commands.SyncPrices{
		ConnectorID:  req.GetConnectorId(),
		ProductIDs:   req.GetProductIds(), // Map proto field to command field
		PriceListIDs: []string{},          // Empty means all price lists
		Since:        time.Time{},         // No since filter in proto
		BatchSize:    int(req.GetBatchSize()),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.SyncPricesResponse{
		SyncId:       uuid.New().String(),
		PricesSynced: 0,
		PricesFailed: 0,
		PriceUpdates: []*erppb.PriceUpdate{},
		Metadata:     &structpb.Struct{},
	}, nil
}

// SyncProducts synchronizes products from the ERP system
func (s server) SyncProducts(ctx context.Context, req *erppb.SyncProductsRequest) (*erppb.SyncProductsResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConnectorID", req.GetConnectorId()),
		attribute.Int("BatchSize", int(req.GetBatchSize())),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Parse since time
	var since time.Time
	if req.GetSince() != nil {
		since = req.GetSince().AsTime()
	}

	// Sync products
	err := s.app.SyncProducts(ctx, commands.SyncProducts{
		ConnectorID: req.GetConnectorId(),
		Since:       since,
		BatchSize:   int(req.GetBatchSize()),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.SyncProductsResponse{
		SyncId:         uuid.New().String(),
		ProductsSynced: 0,
		ProductsFailed: 0,
		Products:       []*erppb.Product{},
		Metadata:       &structpb.Struct{},
	}, nil
}

// SyncStock synchronizes stock levels from the ERP system
func (s server) SyncStock(ctx context.Context, req *erppb.SyncStockRequest) (*erppb.SyncStockResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConnectorID", req.GetConnectorId()),
		attribute.Int("BatchSize", int(req.GetBatchSize())),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Sync stock
	err := s.app.SyncStock(ctx, commands.SyncStock{
		ConnectorID: req.GetConnectorId(),
		ProductIDs:  req.GetProductIds(),
		Since:       time.Time{}, // No since filter in proto
		BatchSize:   int(req.GetBatchSize()),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.SyncStockResponse{
		SyncId:       uuid.New().String(),
		StockSynced:  0,
		StockFailed:  0,
		StockUpdates: []*erppb.StockUpdate{},
		Metadata:     &structpb.Struct{},
	}, nil
}

// Inventory Reservation Granular Commands

// CreateInventoryReservation creates a new inventory reservation
func (s server) CreateInventoryReservation(ctx context.Context, req *erppb.CreateInventoryReservationRequest) (*erppb.CreateInventoryReservationResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ReservationID", req.GetReservationId()),
		attribute.String("OrderID", req.GetOrderId()),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		t := req.GetExpiresAt().AsTime()
		expiresAt = &t
	}

	err := s.app.CreateInventoryReservation(ctx, commands.CreateInventoryReservation{
		ReservationID: req.GetReservationId(),
		OrderID:       req.GetOrderId(),
		SKU:           req.GetSku(),
		WarehouseID:   req.GetWarehouseId(),
		Quantity:      int(req.GetQuantity()),
		ExpiresAt:     expiresAt,
		ConnectorID:   req.GetConnectorId(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.CreateInventoryReservationResponse{
		ReservationId: req.GetReservationId(),
		Status:        "active",
		CreatedAt:     timestamppb.Now(),
	}, nil
}

// ReleaseInventoryReservation releases an inventory reservation
func (s server) ReleaseInventoryReservation(ctx context.Context, req *erppb.ReleaseInventoryReservationRequest) (*erppb.ReleaseInventoryReservationResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ReservationID", req.GetReservationId()))

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.ReleaseInventoryReservation(ctx, commands.ReleaseInventoryReservation{
		ReservationID: req.GetReservationId(),
		Reason:        req.GetReason(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.ReleaseInventoryReservationResponse{
		ReservationId: req.GetReservationId(),
		Status:        "released",
		ReleasedAt:    timestamppb.Now(),
	}, nil
}

// FulfillInventoryReservation fulfills an inventory reservation
func (s server) FulfillInventoryReservation(ctx context.Context, req *erppb.FulfillInventoryReservationRequest) (*erppb.FulfillInventoryReservationResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ReservationID", req.GetReservationId()))

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.FulfillInventoryReservation(ctx, commands.FulfillInventoryReservation{
		ReservationID: req.GetReservationId(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.FulfillInventoryReservationResponse{
		ReservationId: req.GetReservationId(),
		Status:        "fulfilled",
		FulfilledAt:   timestamppb.Now(),
	}, nil
}

// TransferInventoryReservation transfers an inventory reservation
func (s server) TransferInventoryReservation(ctx context.Context, req *erppb.TransferInventoryReservationRequest) (*erppb.TransferInventoryReservationResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ReservationID", req.GetReservationId()),
		attribute.String("ToWarehouseID", req.GetToWarehouseId()),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.TransferInventoryReservation(ctx, commands.TransferInventoryReservation{
		ReservationID: req.GetReservationId(),
		ToWarehouseID: req.GetToWarehouseId(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &erppb.TransferInventoryReservationResponse{
		ReservationId:   req.GetReservationId(),
		FromWarehouseId: "", // This would come from the aggregate
		ToWarehouseId:   req.GetToWarehouseId(),
		Status:          "transferred",
		TransferredAt:   timestamppb.Now(),
	}, nil
}

// UpdateInventoryReservation updates inventory reservations in the ERP system (DEPRECATED)
func (s server) UpdateInventoryReservation(ctx context.Context, req *erppb.UpdateInventoryReservationRequest) (*erppb.UpdateInventoryReservationResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConnectorID", req.GetConnectorId()),
		attribute.String("OrderID", req.GetOrderId()),
		attribute.String("Action", req.GetAction()),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Convert reservation items
	reservations := make([]commands.InventoryReservation, len(req.GetItems()))
	for i, item := range req.GetItems() {
		reservations[i] = commands.InventoryReservation{
			SKU:           item.GetSku(),
			WarehouseID:   item.GetLocation(), // Location maps to WarehouseID
			Quantity:      int(item.GetQuantity()),
			ReservationID: fmt.Sprintf("res_%s_%d", req.GetOrderId(), i),
		}
	}

	// Update inventory reservation
	err := s.app.UpdateInventoryReservation(ctx, commands.UpdateInventoryReservation{
		ConnectorID:  req.GetConnectorId(),
		OrderID:      req.GetOrderId(),
		Action:       commands.ReservationAction(req.GetAction()),
		Reservations: reservations,
	})
	if err != nil {
		return nil, handleError(err)
	}

	// Create response with results
	results := make([]*erppb.ReservationResult, len(reservations))
	for i, res := range reservations {
		results[i] = &erppb.ReservationResult{
			ProductId:         "", // We don't have product ID, only SKU
			Sku:               res.SKU,
			QuantityReserved:  int64(int32(res.Quantity)),
			QuantityAvailable: 100, // This would come from actual inventory check
			Success:           true,
			Error:             "",
		}
	}

	return &erppb.UpdateInventoryReservationResponse{
		ReservationId: uuid.New().String(),
		Status:        "completed",
		Results:       results,
	}, nil
}

// GetConnectorStatus retrieves the status of a specific connector
func (s server) GetConnectorStatus(ctx context.Context, req *erppb.GetConnectorStatusRequest) (*erppb.GetConnectorStatusResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ConnectorID", req.GetConnectorId()))

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Get connector status
	status, err := s.app.GetConnectorStatus(ctx, queries.GetConnectorStatus{
		ConnectorID: req.GetConnectorId(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	// Create sync statuses from last sync summary
	syncStatuses := make([]*erppb.SyncStatus, 0)
	if status.LastSync != nil {
		if status.LastSync.LastProductSync != nil {
			syncStatuses = append(syncStatuses, &erppb.SyncStatus{
				EntityType:   "products",
				Status:       "completed",
				LastSyncedAt: timestamppb.New(*status.LastSync.LastProductSync),
			})
		}
		if status.LastSync.LastStockSync != nil {
			syncStatuses = append(syncStatuses, &erppb.SyncStatus{
				EntityType:   "stock",
				Status:       "completed",
				LastSyncedAt: timestamppb.New(*status.LastSync.LastStockSync),
			})
		}
		if status.LastSync.LastPriceSync != nil {
			syncStatuses = append(syncStatuses, &erppb.SyncStatus{
				EntityType:   "prices",
				Status:       "completed",
				LastSyncedAt: timestamppb.New(*status.LastSync.LastPriceSync),
			})
		}
	}

	// Create metadata from details
	metadata, _ := structpb.NewStruct(status.Details)

	return &erppb.GetConnectorStatusResponse{
		Connector: &erppb.Connector{
			Id:        status.ConnectorID,
			Name:      status.ConnectorID, // Use ID as name for now
			Type:      string(status.Type),
			Status:    string(status.Status),
			CreatedAt: timestamppb.Now(), // Default to now
			UpdatedAt: timestamppb.Now(), // Default to now
		},
		HealthStatus:    string(status.Status),
		LastSync:        timestamppb.Now(), // Default to now
		PendingWebhooks: 0,
		SyncStatuses:    syncStatuses,
		Metadata:        metadata,
	}, nil
}

// GetSyncHistory retrieves the synchronization history for a connector
func (s server) GetSyncHistory(ctx context.Context, req *erppb.GetSyncHistoryRequest) (*erppb.GetSyncHistoryResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("ConnectorID", req.GetConnectorId()),
		attribute.String("EntityType", req.GetEntityType()),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Get sync history
	since := time.Time{}
	if req.GetSince() != nil {
		since = req.GetSince().AsTime()
	}

	limit := 100 // Default limit
	if req.GetPageSize() > 0 {
		limit = int(req.GetPageSize())
	}

	history, err := s.app.GetSyncHistory(ctx, queries.GetSyncHistory{
		ConnectorID: req.GetConnectorId(),
		EntityType:  req.GetEntityType(),
		Since:       since,
		Limit:       limit,
	})
	if err != nil {
		return nil, handleError(err)
	}

	// Convert sync logs
	syncLogs := make([]*erppb.SyncLog, len(history))
	for i, log := range history {
		metadata, _ := structpb.NewStruct(log.Metadata)
		var completedAt *timestamppb.Timestamp
		if log.CompletedAt != nil {
			completedAt = timestamppb.New(*log.CompletedAt)
		}

		syncLogs[i] = &erppb.SyncLog{
			Id:               log.ID,
			ConnectorId:      log.ConnectorID,
			EntityType:       log.EntityType,
			StartedAt:        timestamppb.New(log.StartedAt),
			CompletedAt:      completedAt,
			Status:           string(log.Status),
			RecordsProcessed: int32(log.RecordsProcessed),
			RecordsTotal:     int32(log.RecordsTotal),
			Error:            log.Error,
			Metadata:         metadata,
		}
	}

	// Create summary
	summary, _ := structpb.NewStruct(map[string]interface{}{
		"total_syncs":      len(history),
		"successful":       0, // Would calculate from actual data
		"failed":           0,
		"average_duration": "0s",
	})

	return &erppb.GetSyncHistoryResponse{
		SyncLogs:    syncLogs,
		TotalCount:  int64(len(history)),
		TotalPages:  1,
		CurrentPage: 1,
		Summary:     summary,
	}, nil
}

// ListConnectors lists all registered connectors
func (s server) ListConnectors(ctx context.Context, req *erppb.ListConnectorsRequest) (*erppb.ListConnectorsResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("Type", req.GetType()),
		attribute.String("Status", req.GetStatus()),
	)

	// Validate authentication
	if _, ok := auth.ClaimsFromContext(ctx); !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// List connectors
	connectors, err := s.app.ListConnectors(ctx, queries.ListConnectors{
		Type:   erp.ERPType(req.GetType()),
		Status: req.GetStatus(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	// Convert to proto
	protoConnectors := make([]*erppb.Connector, len(connectors))
	for i, conn := range connectors {
		// Create a minimal config struct with endpoint
		_, _ = structpb.NewStruct(map[string]interface{}{
			"endpoint": conn.Endpoint,
		})

		var lastChecked *timestamppb.Timestamp
		if conn.LastChecked != nil {
			lastChecked = timestamppb.New(*conn.LastChecked)
		}

		protoConnectors[i] = &erppb.Connector{
			Id:        conn.ID,
			Name:      conn.ID, // Use ID as name since ConnectorListItem doesn't have a name
			Type:      string(conn.Type),
			Status:    string(conn.Status),
			CreatedAt: lastChecked, // Use last checked as created time
			UpdatedAt: lastChecked, // Use last checked as updated time
		}
	}

	return &erppb.ListConnectorsResponse{
		Connectors:  protoConnectors,
		TotalCount:  int64(len(connectors)),
		TotalPages:  1,
		CurrentPage: 1,
	}, nil
}

// connectorToProto converts a ConnectorStatus to proto Connector
func connectorToProto(status *queries.ConnectorStatus) *erppb.Connector {
	// Default timestamps
	var createdAt, updatedAt, lastHealthCheckAt *timestamppb.Timestamp
	if status.Details != nil {
		// Try to extract timestamps from details if available
		if created, ok := status.Details["created_at"].(time.Time); ok {
			createdAt = timestamppb.New(created)
		}
		if updated, ok := status.Details["updated_at"].(time.Time); ok {
			updatedAt = timestamppb.New(updated)
		}
	}

	// Use current time as fallback
	if createdAt == nil {
		createdAt = timestamppb.Now()
	}
	if updatedAt == nil {
		updatedAt = timestamppb.Now()
	}

	// Extract health check info
	var healthStatus, healthError string
	if status.LastSync != nil && status.LastSync.LastProductSync != nil {
		lastHealthCheckAt = timestamppb.New(*status.LastSync.LastProductSync)
		healthStatus = string(status.Status)
		healthError = status.Message
	}

	// Extract other fields from details
	var environment, baseURL, createdBy, updatedBy string
	var webhookEnabled, syncEnabled bool
	var syncInterval, batchSize, version int32
	var webhookEvents []string

	if status.Details != nil {
		if env, ok := status.Details["environment"].(string); ok {
			environment = env
		}
		if url, ok := status.Details["base_url"].(string); ok {
			baseURL = url
		}
		if by, ok := status.Details["created_by"].(string); ok {
			createdBy = by
		}
		if by, ok := status.Details["updated_by"].(string); ok {
			updatedBy = by
		}
		if enabled, ok := status.Details["webhook_enabled"].(bool); ok {
			webhookEnabled = enabled
		}
		if enabled, ok := status.Details["sync_enabled"].(bool); ok {
			syncEnabled = enabled
		}
		if interval, ok := status.Details["sync_interval_seconds"].(float64); ok {
			syncInterval = int32(interval)
		}
		if size, ok := status.Details["batch_size"].(float64); ok {
			batchSize = int32(size)
		}
		if v, ok := status.Details["version"].(float64); ok {
			version = int32(v)
		}
		if events, ok := status.Details["webhook_events"].([]interface{}); ok {
			for _, event := range events {
				if str, ok := event.(string); ok {
					webhookEvents = append(webhookEvents, str)
				}
			}
		}
	}

	// Set defaults
	if environment == "" {
		environment = "production"
	}
	if syncInterval == 0 {
		syncInterval = 300 // 5 minutes default
	}
	if batchSize == 0 {
		batchSize = 100
	}

	return &erppb.Connector{
		Id:                    status.ConnectorID,
		Name:                  status.ConnectorID, // Use ID as name if not available
		Type:                  string(status.Type),
		Status:                string(status.Status),
		Environment:           environment,
		BaseUrl:               baseURL,
		WebhookEnabled:        webhookEnabled,
		WebhookUrl:            fmt.Sprintf("/api/erp/webhook/%s", status.ConnectorID),
		WebhookEvents:         webhookEvents,
		SyncEnabled:           syncEnabled,
		SyncIntervalSeconds:   syncInterval,
		BatchSize:             batchSize,
		LastHealthCheckAt:     lastHealthCheckAt,
		LastHealthCheckStatus: healthStatus,
		LastHealthCheckError:  healthError,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
		CreatedBy:             createdBy,
		UpdatedBy:             updatedBy,
		Version:               version,
	}
}

// handleError converts application errors to gRPC status errors
func handleError(err error) error {
	if err == nil {
		return nil
	}

	// Map specific errors to gRPC codes
	switch err.Error() {
	case "connector not found":
		return status.Error(grpc_code.NotFound, err.Error())
	case "invalid configuration":
		return status.Error(grpc_code.InvalidArgument, err.Error())
	case "authentication required":
		return status.Error(grpc_code.Unauthenticated, err.Error())
	case "permission denied":
		return status.Error(grpc_code.PermissionDenied, err.Error())
	default:
		return status.Error(grpc_code.Internal, err.Error())
	}
}

// Helper functions for config parsing
func getStringFromConfigMap(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getBoolFromConfigMap(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}
