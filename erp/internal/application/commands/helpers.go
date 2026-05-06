package commands

import (
	"fmt"
	"time"
)

// Helper functions shared across command handlers

func ptrTime(t time.Time) *time.Time {
	return &t
}

func ptrString(s string) *string {
	return &s
}

func stringPtr(s string) *string {
	return &s
}

func generateID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func generateWebhookID() string {
	return generateID("webhook")
}

func generateOrderSyncID() string {
	return generateID("order_sync")
}