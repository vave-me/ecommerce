package utils

import "strconv"

// StringVal tries a string cast or returns ""
func StringVal(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// IntVal tries to convert a value to int, returns 0 if not possible
func IntVal(v interface{}) int {
	if i, ok := v.(int); ok {
		return i
	}
	if s, ok := v.(string); ok {
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
	}
	return 0
}

// BoolVal tries to convert a value to bool using various type checks
func BoolVal(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	if s, ok := v.(string); ok {
		return s == "true" || s == "1"
	}
	if i, ok := v.(int); ok {
		return i > 0
	}
	return false
}
