package commands

import (
	"context"
	"fmt"
	"middleman/internal/erp"
	"time"

	"middleman/erp/internal/domain"
)

// SyncCustomers command triggers customer master data synchronization
type SyncCustomers struct {
	ConnectorID   string
	CustomerIDs   []string // If empty, sync all customers
	Since         time.Time
	IncludeCredit bool // Include credit limit and payment terms
	BatchSize     int
}

// SyncCustomersHandler handles the SyncCustomers command
type SyncCustomersHandler struct {
	registry   erp.ConnectorRegistry
	repository domain.SyncLogRepository
}

// NewSyncCustomersHandler creates a new handler
func NewSyncCustomersHandler(
	registry erp.ConnectorRegistry,
	repository domain.SyncLogRepository,
) SyncCustomersHandler {
	return SyncCustomersHandler{
		registry:   registry,
		repository: repository,
	}
}

// SyncCustomers synchronizes customers from the ERP
func (h SyncCustomersHandler) SyncCustomers(ctx context.Context, cmd SyncCustomers) error {
	// Get connector
	connector, err := h.registry.GetConnector(cmd.ConnectorID)
	if err != nil {
		return fmt.Errorf("getting connector: %w", err)
	}

	// Create sync log
	syncLog := &domain.SyncLog{
		ID:          generateCustomerSyncID(),
		ConnectorID: cmd.ConnectorID,
		EntityType:  "customer",
		StartedAt:   time.Now(),
		Status:      domain.SyncStatusInProgress,
		Metadata: map[string]interface{}{
			"customerIds":   cmd.CustomerIDs,
			"includeCredit": cmd.IncludeCredit,
		},
	}

	if err := h.repository.Create(ctx, syncLog); err != nil {
		return fmt.Errorf("creating sync log: %w", err)
	}

	// Fetch customers
	var customers []*erp.CustomerPayload
	if len(cmd.CustomerIDs) > 0 {
		// Fetch specific customers
		for _, customerID := range cmd.CustomerIDs {
			customer, err := connector.FetchCustomer(ctx, customerID)
			if err != nil {
				syncLog.Error = fmt.Sprintf("failed to fetch customer %s: %v", customerID, err)
				continue
			}
			customers = append(customers, customer)
		}
	} else {
		// Sync all customers modified since the given time
		customers, err = connector.SyncCustomers(ctx, cmd.Since, cmd.BatchSize)
		if err != nil {
			syncLog.Status = domain.SyncStatusFailed
			syncLog.Error = err.Error()
			syncLog.CompletedAt = ptrTime(time.Now())
			h.repository.Update(ctx, syncLog)
			return fmt.Errorf("fetching customers: %w", err)
		}
	}

	// Process customers in batches
	batchSize := cmd.BatchSize
	if batchSize == 0 {
		batchSize = 50
	}

	processedCount := 0
	for i := 0; i < len(customers); i += batchSize {
		end := i + batchSize
		if end > len(customers) {
			end = len(customers)
		}

		batch := customers[i:end]
		for _, customer := range batch {
			// Validate customer data
			if err := validateCustomer(customer); err != nil {
				syncLog.Error = fmt.Sprintf("invalid customer %s: %v", customer.CustomerID, err)
				continue
			}

			// Publish customer event
			// Publishing events removed - handled separately
			// if err := h.publisher.Publish(ctx, event); err != nil {
			// 	// Log error but continue processing
			// 	syncLog.Error = fmt.Sprintf("failed to publish event for customer %s: %v", customer.CustomerID, err)
			// } else {
			processedCount++
		}
	}

	// Update sync log
	syncLog.Status = domain.SyncStatusCompleted
	syncLog.CompletedAt = ptrTime(time.Now())
	syncLog.RecordsProcessed = processedCount
	syncLog.RecordsTotal = len(customers)

	if err := h.repository.Update(ctx, syncLog); err != nil {
		return fmt.Errorf("updating sync log: %w", err)
	}

	return nil
}

func validateCustomer(customer *erp.CustomerPayload) error {
	if customer.CustomerID == "" {
		return fmt.Errorf("customer ID is required")
	}
	if customer.Email == "" && customer.Phone == "" {
		return fmt.Errorf("at least one contact method (email or phone) is required")
	}
	// Currency validation removed - CustomerPayload doesn't have currency field
	return nil
}

func generateCustomerSyncID() string {
	return fmt.Sprintf("customer_sync_%d", time.Now().UnixNano())
}
