package models

type User struct {
	ID        string
	Email     string
	Username  string
	FirstName string
	LastName  string
	Location  string
	Enabled   bool
}
