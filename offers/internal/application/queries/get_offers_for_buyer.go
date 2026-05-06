package queries

type GetOffersForBuyerQuery struct {
	BuyerID  string
	Page     int
	PageSize int
}
