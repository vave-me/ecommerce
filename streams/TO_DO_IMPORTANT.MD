Implement streams with variants 

`package main

import (
"fmt"
)

// Stream struct representing the main product
type Stream struct {
StreamID       string      `json:"product_id"`
Name            string      `json:"name"`
SKU             string      `json:"sku"`
Description     string      `json:"description"`
Category        string      `json:"category"`
Brand           string      `json:"brand"`
Price           float64     `json:"price"`
StockQuantity   int         `json:"stock_quantity"`
MainImages      []string    `json:"main_images"`   // Multiple images
Tags            []string    `json:"tags"`
URLSlug         string      `json:"url_slug"`
Weight          float64     `json:"weight"`
Dimensions      Dimensions  `json:"dimensions"`
MetaTitle       string      `json:"meta_title"`
MetaDescription string      `json:"meta_description"`
Variants        []Variant   `json:"variants"`
Discounts       []Discount  `json:"discounts,omitempty"`
}

// Dimensions struct for product dimensions
type Dimensions struct {
Height float64 `json:"height"`
Width  float64 `json:"width"`
Depth  float64 `json:"depth"`
}

// Variant struct representing a specific product variant
type Variant struct {
ParentStreamID string      `json:"parent_product_id"`
VariantSKU      string      `json:"variant_sku"`
VariantName     string      `json:"variant_name"`
Option1Name     string      `json:"option_1_name"`
Option1Value    string      `json:"option_1_value"`
Option2Name     string      `json:"option_2_name"`
Option2Value    string      `json:"option_2_value"`
Option3Name     string      `json:"option_3_name,omitempty"` // Optional
Option3Value    string      `json:"option_3_value,omitempty"` // Optional
VariantPrice    float64     `json:"variant_price"`
StockQuantity   int         `json:"stock_quantity"`
VariantImages   []string    `json:"variant_images"`   // Multiple images for each variant
Weight          float64     `json:"weight"`
Dimensions      Dimensions  `json:"dimensions"`
Barcode         string      `json:"barcode,omitempty"`        // Optional
TaxClass        string      `json:"tax_class,omitempty"`      // Optional
ShippingClass   string      `json:"shipping_class,omitempty"` // Optional
Visibility      bool        `json:"visibility"`
MetaTitle       string      `json:"meta_title,omitempty"`     // Optional
MetaDescription string      `json:"meta_description,omitempty"` // Optional
Discounts       []Discount  `json:"discounts,omitempty"`      // Optional
}

// Discount struct to represent any discounts applied to the product or variant
type Discount struct {
Type         string  `json:"type"`   // E.g., "bulk", "percentage"
Amount       float64 `json:"amount"` // Discount amount
MinQuantity  int     `json:"min_quantity,omitempty"` // Minimum quantity to apply discount (for bulk discounts)
StartDate    string  `json:"start_date,omitempty"` // Discount start date
EndDate      string  `json:"end_date,omitempty"`   // Discount end date
}

func main() {
// Example product with variants, multiple images, and discounts
product := Stream{
StreamID:     "10001",
Name:          "T-Shirt",
SKU:           "TS-001",
Description:   "A comfortable cotton t-shirt",
Category:      "Apparel",
Brand:         "BrandName",
Price:         19.99,
StockQuantity: 150,
MainImages:    []string{"https://example.com/images/tshirt-1.jpg", "https://example.com/images/tshirt-2.jpg"},
Tags:          []string{"eco-friendly", "bestseller"},
URLSlug:       "t-shirt",
Weight:        0.2,
Dimensions: Dimensions{
Height: 1.0,
Width:  0.5,
Depth:  0.1,
},
Variants: []Variant{
{
ParentStreamID: "10001",
VariantSKU:      "TS-001-M-R",
VariantName:     "T-Shirt - Medium, Red",
Option1Name:     "Size",
Option1Value:    "Medium",
Option2Name:     "Color",
Option2Value:    "Red",
VariantPrice:    19.99,
StockQuantity:   100,
VariantImages:   []string{"https://example.com/images/tshirt-medium-red-1.jpg", "https://example.com/images/tshirt-medium-red-2.jpg"},
Weight:          0.2,
Dimensions: Dimensions{
Height: 1.0,
Width:  0.5,
Depth:  0.1,
},
Visibility: true,
},
{
ParentStreamID: "10001",
VariantSKU:      "TS-001-L-B",
VariantName:     "T-Shirt - Large, Blue",
Option1Name:     "Size",
Option1Value:    "Large",
Option2Name:     "Color",
Option2Value:    "Blue",
VariantPrice:    21.99,
StockQuantity:   50,
VariantImages:   []string{"https://example.com/images/tshirt-large-blue-1.jpg", "https://example.com/images/tshirt-large-blue-2.jpg"},
Weight:          0.2,
Dimensions: Dimensions{
Height: 1.0,
Width:  0.5,
Depth:  0.1,
},
Visibility: true,
},
},
Discounts: []Discount{
{
Type:        "percentage",
Amount:      10.0,
StartDate:   "2024-10-01",
EndDate:     "2024-10-31",
},
},
}

	fmt.Printf("%+v\n", product)
}
`