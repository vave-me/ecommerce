package domain

import "context"

// OrderRepository defines the interface for order aggregate persistence
// This follows the Event Sourcing pattern for write operations
type OrderRepository interface {
	// Load retrieves an order aggregate by ID
	Load(ctx context.Context, orderID string) (*Order, error)
	
	// Save persists order events (for Event Sourcing)
	Save(ctx context.Context, order *Order) error
}
