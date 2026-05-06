package models // In a file like domain/product.go or domain/filters.go

// ProductFilter defines the criteria for filtering products.
type ProductFilter struct {
	Name       string   // Name of the product (or partial name for searching)
	CategoryID string   // ID of the category to filter by
	MinPrice   int64    // Minimum price in cents
	MaxPrice   int64    // Maximum price in cents
	Brand      string   // Brand of the product
	Condition  string   // Condition of the product (e.g., "new", "used")
	Tags       []string // List of tags to filter by
	Status     string   // Status of the product (e.g., "active", "sold")

	// Pagination
	Limit  int // Number of items per page
	Offset int // Number of items to skip

	// Sorting
	SortBy    string // Field to sort by (e.g., "created_at", "price")
	SortOrder string // Sort order ("asc" or "desc")

	// Location-based filtering
	Latitude  float64 // Latitude for geo-search
	Longitude float64 // Longitude for geo-search
	RadiusKM  float64 // Radius in kilometers for geo-search

	// User context
	UserID string // ID of the user performing the search (for personalization, history, etc.)

	// Add any other common filter fields you might need, e.g.:
	// SellerID string
	// HasDiscount bool
	// MinRating float64
}
