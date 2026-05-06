package domain_test

import (
	"testing"

	"google.golang.org/api/content/v2.1"
	"middleman/merchant/internal/domain"
)

func TestProductValidator_ValidateProduct(t *testing.T) {
	validator := domain.NewProductValidator()

	tests := []struct {
		name        string
		product     *content.Product
		expectError bool
		errorCount  int
	}{
		{
			name: "valid product",
			product: &content.Product{
				OfferId:         "test-123",
				Title:           "Test Product",
				Description:     "A test product description",
				Link:            "https://example.com/product/123",
				ImageLink:       "https://example.com/images/product.jpg",
				ContentLanguage: "en",
				TargetCountry:   "US",
				Availability:    "in_stock",
				Condition:       "new",
				Brand:           "TestBrand",
				Price: &content.Price{
					Value:    "19.99",
					Currency: "USD",
				},
			},
			expectError: false,
			errorCount:  0,
		},
		{
			name: "missing required fields",
			product: &content.Product{
				OfferId: "test-123",
			},
			expectError: true,
			errorCount:  9, // title, description, link, image, price, availability, language, country, brand
		},
		{
			name: "invalid URLs",
			product: &content.Product{
				OfferId:         "test-123",
				Title:           "Test Product",
				Description:     "A test product description",
				Link:            "not-a-url",
				ImageLink:       "also-not-a-url",
				ContentLanguage: "en",
				TargetCountry:   "US",
				Availability:    "in_stock",
				Brand:           "TestBrand",
				Price: &content.Price{
					Value:    "19.99",
					Currency: "USD",
				},
			},
			expectError: true,
			errorCount:  2,
		},
		{
			name: "invalid availability",
			product: &content.Product{
				OfferId:         "test-123",
				Title:           "Test Product",
				Description:     "A test product description",
				Link:            "https://example.com/product/123",
				ImageLink:       "https://example.com/images/product.jpg",
				ContentLanguage: "en",
				TargetCountry:   "US",
				Availability:    "maybe_in_stock",
				Brand:           "TestBrand",
				Price: &content.Price{
					Value:    "19.99",
					Currency: "USD",
				},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "title too long",
			product: &content.Product{
				OfferId:         "test-123",
				Title:           string(make([]byte, 151)), // 151 characters
				Description:     "A test product description",
				Link:            "https://example.com/product/123",
				ImageLink:       "https://example.com/images/product.jpg",
				ContentLanguage: "en",
				TargetCountry:   "US",
				Availability:    "in_stock",
				Brand:           "TestBrand",
				Price: &content.Price{
					Value:    "19.99",
					Currency: "USD",
				},
			},
			expectError: true,
			errorCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateProduct(tt.product)
			
			if tt.expectError && len(errors) == 0 {
				t.Errorf("expected validation errors but got none")
			}
			
			if !tt.expectError && len(errors) > 0 {
				t.Errorf("unexpected validation errors: %v", errors)
			}
			
			if len(errors) != tt.errorCount {
				t.Errorf("expected %d errors but got %d: %v", tt.errorCount, len(errors), errors)
			}
		})
	}
}