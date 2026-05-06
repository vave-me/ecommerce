package queries

type GetVisitedItemsForBuyerQuery struct {
	BuyerID   string
	MinVisits int // Minimum number of visits to count
}
