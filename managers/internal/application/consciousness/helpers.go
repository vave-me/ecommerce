package consciousness

import (
	"github.com/google/uuid"
)

// generateID generates a unique ID for decisions and actions
func generateID() string {
	return uuid.New().String()
}