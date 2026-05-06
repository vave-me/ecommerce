package domain

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"google.golang.org/api/content/v2.1"
)

// ProductValidator validates products according to Google Merchant Center requirements
type ProductValidator struct {
	// URL validation regex
	urlRegex *regexp.Regexp
	// GTIN validation regex (8, 12, 13, or 14 digits)
	gtinRegex *regexp.Regexp
	// Currency code validation
	currencyRegex *regexp.Regexp
}

// NewProductValidator creates a new product validator
func NewProductValidator() *ProductValidator {
	return &ProductValidator{
		urlRegex:      regexp.MustCompile(`^https?://`),
		gtinRegex:     regexp.MustCompile(`^(\d{8}|\d{12}|\d{13}|\d{14})$`),
		currencyRegex: regexp.MustCompile(`^[A-Z]{3}$`),
	}
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	
	var messages []string
	for _, e := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", e.Field, e.Message))
	}
	return strings.Join(messages, "; ")
}

// ValidateProduct validates a product according to Google Merchant Center requirements
func (v *ProductValidator) ValidateProduct(product *content.Product) ValidationErrors {
	var errors ValidationErrors
	
	// Required fields validation
	if product.OfferId == "" {
		errors = append(errors, ValidationError{
			Field:   "offerId",
			Message: "offer ID is required",
		})
	}
	
	if product.Title == "" {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "title is required",
		})
	} else if len(product.Title) > 150 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "title must not exceed 150 characters",
		})
	}
	
	if product.Description == "" {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "description is required",
		})
	} else if len(product.Description) > 5000 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "description must not exceed 5000 characters",
		})
	}
	
	// Link validation
	if product.Link == "" {
		errors = append(errors, ValidationError{
			Field:   "link",
			Message: "product link is required",
		})
	} else if !v.isValidURL(product.Link) {
		errors = append(errors, ValidationError{
			Field:   "link",
			Message: "product link must be a valid HTTP(S) URL",
		})
	}
	
	// Image link validation
	if product.ImageLink == "" {
		errors = append(errors, ValidationError{
			Field:   "imageLink",
			Message: "image link is required",
		})
	} else if !v.isValidURL(product.ImageLink) {
		errors = append(errors, ValidationError{
			Field:   "imageLink",
			Message: "image link must be a valid HTTP(S) URL",
		})
	}
	
	// Price validation
	if product.Price == nil {
		errors = append(errors, ValidationError{
			Field:   "price",
			Message: "price is required",
		})
	} else {
		if product.Price.Value == "" {
			errors = append(errors, ValidationError{
				Field:   "price.value",
				Message: "price value is required",
			})
		}
		
		if product.Price.Currency == "" {
			errors = append(errors, ValidationError{
				Field:   "price.currency",
				Message: "price currency is required",
			})
		} else if !v.currencyRegex.MatchString(product.Price.Currency) {
			errors = append(errors, ValidationError{
				Field:   "price.currency",
				Message: "price currency must be a valid 3-letter ISO 4217 code",
			})
		}
	}
	
	// Availability validation
	if product.Availability == "" {
		errors = append(errors, ValidationError{
			Field:   "availability",
			Message: "availability is required",
		})
	} else if !v.isValidAvailability(product.Availability) {
		errors = append(errors, ValidationError{
			Field:   "availability",
			Message: "availability must be one of: in_stock, out_of_stock, preorder, backorder",
		})
	}
	
	// Condition validation
	if product.Condition == "" {
		product.Condition = "new" // Default to new if not specified
	} else if !v.isValidCondition(product.Condition) {
		errors = append(errors, ValidationError{
			Field:   "condition",
			Message: "condition must be one of: new, refurbished, used",
		})
	}
	
	// Brand validation (required for most categories)
	if product.Brand == "" && !product.IdentifierExists {
		errors = append(errors, ValidationError{
			Field:   "brand",
			Message: "brand is required unless identifier_exists is set to false",
		})
	}
	
	// GTIN validation if provided
	if product.Gtin != "" && !v.gtinRegex.MatchString(product.Gtin) {
		errors = append(errors, ValidationError{
			Field:   "gtin",
			Message: "GTIN must be 8, 12, 13, or 14 digits",
		})
	}
	
	// Additional image links validation
	for i, imageLink := range product.AdditionalImageLinks {
		if !v.isValidURL(imageLink) {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("additionalImageLinks[%d]", i),
				Message: "additional image link must be a valid HTTP(S) URL",
			})
		}
	}
	
	// Shipping validation if provided
	for i, shipping := range product.Shipping {
		if shipping.Country == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("shipping[%d].country", i),
				Message: "shipping country is required",
			})
		}
		
		if shipping.Price == nil {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("shipping[%d].price", i),
				Message: "shipping price is required",
			})
		}
	}
	
	// Target country and language validation
	if product.TargetCountry == "" {
		errors = append(errors, ValidationError{
			Field:   "targetCountry",
			Message: "target country is required",
		})
	} else if len(product.TargetCountry) != 2 {
		errors = append(errors, ValidationError{
			Field:   "targetCountry",
			Message: "target country must be a 2-letter ISO 3166 country code",
		})
	}
	
	if product.ContentLanguage == "" {
		errors = append(errors, ValidationError{
			Field:   "contentLanguage",
			Message: "content language is required",
		})
	} else if len(product.ContentLanguage) != 2 {
		errors = append(errors, ValidationError{
			Field:   "contentLanguage",
			Message: "content language must be a 2-letter ISO 639 language code",
		})
	}
	
	return errors
}

// isValidURL checks if the URL is valid
func (v *ProductValidator) isValidURL(urlStr string) bool {
	if !v.urlRegex.MatchString(urlStr) {
		return false
	}
	
	_, err := url.Parse(urlStr)
	return err == nil
}

// isValidAvailability checks if the availability value is valid
func (v *ProductValidator) isValidAvailability(availability string) bool {
	validValues := []string{"in_stock", "out_of_stock", "preorder", "backorder"}
	for _, valid := range validValues {
		if availability == valid {
			return true
		}
	}
	return false
}

// isValidCondition checks if the condition value is valid
func (v *ProductValidator) isValidCondition(condition string) bool {
	validValues := []string{"new", "refurbished", "used"}
	for _, valid := range validValues {
		if condition == valid {
			return true
		}
	}
	return false
}