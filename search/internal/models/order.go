package models

import "time"

type Order struct {
	OrderID          string
	UserCustomerID   string
	UserCustomerName string
	Items            []Item
	Total            int64
	Status           string
	CreatedAt        time.Time
}

type Item struct {
	ProductID      string
	UserSellerID   string
	ProductName    string
	UserSellerName string
	Price          int64
	Quantity       int64
}
