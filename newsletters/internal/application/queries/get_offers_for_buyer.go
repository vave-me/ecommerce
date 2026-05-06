package queries

type GetNewslettersForBuyerQuery struct {
	BuyerID  string
	Page     int
	PageSize int
}
