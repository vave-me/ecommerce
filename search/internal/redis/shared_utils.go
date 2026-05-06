// File: search/internal/redis/shared_utils.go
package redis

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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

// boolToInt converts a boolean value to an integer (1 for true, 0 for false)
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// safeString returns the string as-is, ensuring it's not nil
// This is mainly for clarity and consistency in the codebase
func safeString(s string) string {
	if s == "" {
		return ""
	}
	return s
}

// strVal safely converts an interface{} to string
func strVal(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

// parseInt64 parses an interface{} value to int64 with proper error handling and logging
func parseInt64(value interface{}, fieldName, docID string) (int64, error) {
	str := strings.TrimSpace(strVal(value))
	if str == "" {
		return 0, nil
	}
	
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Printf("parseInt64: Failed to parse field %s='%s' for doc %s: %v", fieldName, str, docID, err)
		return 0, err
	}
	return val, nil
}
