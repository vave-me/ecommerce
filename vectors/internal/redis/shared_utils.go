// File: vector/internal/redis/shared_utils.go
package redis

import (
	"log"
)

// parseBool handles both string boolean values ("true"/"false") and numeric values (0/1)
// This is the CRITICAL FIX for property boolean parsing errors
func parseBool(value interface{}, fieldName, docID string) bool {
	if value == nil {
		return false
	}

	switch vv := value.(type) {
	case bool:
		return vv
	case string:
		switch vv {
		case "true", "True", "TRUE", "1":
			return true
		case "false", "False", "FALSE", "0", "":
			return false
		default:
			log.Printf("parseBool: Warning - unexpected string value '%s' for field %s in doc %s. Treating as false.", vv, fieldName, docID)
			return false
		}
	case int:
		return vv != 0
	case int64:
		return vv != 0
	case float64:
		return vv != 0
	default:
		log.Printf("parseBool: Warning - unsupported type %T for field %s in doc %s. Treating as false.", vv, fieldName, docID)
		return false
	}
}

// parseBoolVal is a convenience wrapper for parseBool with consistent logging
func parseBoolVal(value interface{}, fieldName, docID string) bool {
	return parseBool(value, fieldName, docID)
}
