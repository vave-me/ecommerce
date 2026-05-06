package domain

type Product struct {
	ID           string
	Name         string
	Description  string
	BasePrice    int64
	UserSellerID string
	Stock        int64
	SKU          string
	CategoryID   string
}
