package tools

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// MultiValidationError represents multiple validation errors
type MultiValidationError struct {
	Errors []ValidationError
}

func (e *MultiValidationError) Error() string {
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// Validator provides validation methods for common parameter types
type Validator struct {
	errors []ValidationError
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		errors: make([]ValidationError, 0),
	}
}

// ValidateRequired checks if a required field is present and non-empty
func (v *Validator) ValidateRequired(field string, value interface{}) *Validator {
	if value == nil {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "required field is missing",
		})
		return v
	}

	// Check for empty strings
	if str, ok := value.(string); ok && strings.TrimSpace(str) == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "required field cannot be empty",
		})
	}

	return v
}

// ValidateMinLength validates minimum string length
func (v *Validator) ValidateMinLength(field, value string, minLength int) *Validator {
	if len(value) < minLength {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at least %d characters long", minLength),
		})
	}
	return v
}

// ValidateMaxLength validates maximum string length
func (v *Validator) ValidateMaxLength(field, value string, maxLength int) *Validator {
	if len(value) > maxLength {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must not exceed %d characters", maxLength),
		})
	}
	return v
}

// ValidateEmail validates email format
func (v *Validator) ValidateEmail(field, email string) *Validator {
	_, err := mail.ParseAddress(email)
	if err != nil {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "invalid email format",
		})
	}
	return v
}

// ValidateURL validates URL format
func (v *Validator) ValidateURL(field, urlStr string) *Validator {
	_, err := url.Parse(urlStr)
	if err != nil {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "invalid URL format",
		})
	}
	return v
}

// ValidatePattern validates string against regex pattern
func (v *Validator) ValidatePattern(field, value, pattern string) *Validator {
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("invalid pattern: %v", err),
		})
		return v
	}
	if !matched {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("does not match required pattern: %s", pattern),
		})
	}
	return v
}

// ValidateEnum validates value is in allowed set
func (v *Validator) ValidateEnum(field, value string, allowedValues []string) *Validator {
	found := false
	for _, allowed := range allowedValues {
		if value == allowed {
			found = true
			break
		}
	}
	if !found {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be one of: %s", strings.Join(allowedValues, ", ")),
		})
	}
	return v
}

// ValidateMinimum validates minimum numeric value
func (v *Validator) ValidateMinimum(field string, value, minimum float64) *Validator {
	if value < minimum {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at least %v", minimum),
		})
	}
	return v
}

// ValidateMaximum validates maximum numeric value
func (v *Validator) ValidateMaximum(field string, value, maximum float64) *Validator {
	if value > maximum {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must not exceed %v", maximum),
		})
	}
	return v
}

// ValidateRange validates numeric value is within range
func (v *Validator) ValidateRange(field string, value, min, max float64) *Validator {
	return v.ValidateMinimum(field, value, min).ValidateMaximum(field, value, max)
}

// ValidateArrayMinItems validates minimum array length
func (v *Validator) ValidateArrayMinItems(field string, items []interface{}, minItems int) *Validator {
	if len(items) < minItems {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must contain at least %d items", minItems),
		})
	}
	return v
}

// ValidateArrayMaxItems validates maximum array length
func (v *Validator) ValidateArrayMaxItems(field string, items []interface{}, maxItems int) *Validator {
	if len(items) > maxItems {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must not contain more than %d items", maxItems),
		})
	}
	return v
}

// ValidateLatitude validates latitude value
func (v *Validator) ValidateLatitude(field string, lat float64) *Validator {
	return v.ValidateRange(field, lat, -90, 90)
}

// ValidateLongitude validates longitude value
func (v *Validator) ValidateLongitude(field string, lng float64) *Validator {
	return v.ValidateRange(field, lng, -180, 180)
}

// HasErrors returns true if validation errors exist
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// GetError returns a single error if any validation failed
func (v *Validator) GetError() error {
	if !v.HasErrors() {
		return nil
	}
	if len(v.errors) == 1 {
		return &v.errors[0]
	}
	return &MultiValidationError{Errors: v.errors}
}

// Reset clears all validation errors
func (v *Validator) Reset() {
	v.errors = v.errors[:0]
}

// ValidateProductSearchParams validates product search parameters
func ValidateProductSearchParams(params map[string]interface{}) error {
	v := NewValidator()

	// Validate price range
	if minPrice := getInt64Param(params, "min_price", 0); minPrice < 0 {
		v.errors = append(v.errors, ValidationError{
			Field:   "min_price",
			Message: "must be non-negative",
		})
	}

	if maxPrice := getInt64Param(params, "max_price", 0); maxPrice < 0 {
		v.errors = append(v.errors, ValidationError{
			Field:   "max_price",
			Message: "must be non-negative",
		})
	}

	// Validate coordinates if provided
	if lat, exists := params["lat"]; exists {
		if latFloat, ok := lat.(float64); ok {
			v.ValidateLatitude("lat", latFloat)
		}
	}

	if lng, exists := params["lng"]; exists {
		if lngFloat, ok := lng.(float64); ok {
			v.ValidateLongitude("lng", lngFloat)
		}
	}

	// Validate radius
	if radius := getFloat64Param(params, "radius", 0); radius > 0 {
		v.ValidateRange("radius", radius, 0, 1000)
	}

	// Validate page/limit
	if page := getInt64Param(params, "page", 1); page < 1 {
		v.errors = append(v.errors, ValidationError{
			Field:   "page",
			Message: "must be at least 1",
		})
	}

	if limit := getInt64Param(params, "limit", 20); limit < 1 || limit > 100 {
		v.errors = append(v.errors, ValidationError{
			Field:   "limit",
			Message: "must be between 1 and 100",
		})
	}

	return v.GetError()
}

// ValidateEmailParam validates an email parameter
func ValidateEmailParam(email string) error {
	v := NewValidator()
	v.ValidateRequired("email", email).ValidateEmail("email", email)
	return v.GetError()
}

// ValidateIDParam validates an ID parameter
func ValidateIDParam(field, id string) error {
	v := NewValidator()
	v.ValidateRequired(field, id).ValidateMinLength(field, id, 1)
	return v.GetError()
}

// ValidatePaginationParams validates pagination parameters
func ValidatePaginationParams(page, limit int64) error {
	v := NewValidator()
	v.ValidateMinimum("page", float64(page), 1)
	v.ValidateRange("limit", float64(limit), 1, 100)
	return v.GetError()
}

// SanitizeString removes potentially dangerous characters from strings
func SanitizeString(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	// Trim excessive whitespace
	input = strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	space := regexp.MustCompile(`\s+`)
	input = space.ReplaceAllString(input, " ")
	
	return input
}

// SanitizeHTML removes HTML tags from input
func SanitizeHTML(input string) string {
	// Simple HTML tag removal - for production use a proper HTML sanitizer
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}

// SanitizeFilename removes dangerous characters from filenames
func SanitizeFilename(filename string) string {
	// Remove path separators and other dangerous characters
	dangerous := []string{"/", "\\", "..", ":", "*", "?", "\"", "<", ">", "|", "\x00"}
	result := filename
	for _, char := range dangerous {
		result = strings.ReplaceAll(result, char, "")
	}
	return result
}