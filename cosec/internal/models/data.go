package models

type CheckoutData struct {
	BasketID  string
	OrderID   string
	PaymentID string
	UserID    string
	Total     int64 // cents
	Items     []CheckoutItem
}

type CheckoutItem struct {
	ProductID   string
	ProductName string
	SellerID    string
	SellerName  string
	Quantity    int64
	Price       int64
}
