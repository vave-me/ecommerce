package queries

type GetNewslettersForSellerQuery struct {
	SellerID string
	Page     int
	PageSize int
}
