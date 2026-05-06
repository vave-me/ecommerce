package domain

// UserAlertAdded logs when a user-related notification is sent.
type UserAlertAdded struct {
	ID        string
	AlertType string
	UserID    string // ID of the user receiving the notification
	Payload   map[string]interface{}
	Message   string // Content of the notification
	IsRead    bool
}

// ProductAlertAdded logs when a product-related notification is sent.
type ProductAlertAdded struct {
	ID        string
	AlertType string
	UserID    string // ID of the user receiving the notification
	ProductID string // ID of the product in the notification
	Payload   map[string]interface{}
	Message   string
	IsRead    bool // Content of the notification
}

// BasketAlertAdded logs when a product-related notification is sent.
type BasketAlertAdded struct {
	ID        string
	AlertType string
	BasketID  string
	UserID    string // ID of the user receiving the notification
	ProductID string // ID of the product in the notification
	Payload   map[string]interface{}
	Message   string // Content of the notification
	IsRead    bool
}

// OrderAlertAdded logs when an order-related notification is sent.
type OrderAlertAdded struct {
	ID        string
	AlertType string
	UserID    string // ID of the user receiving the notification
	OrderID   string // ID of the order in the notification
	ProductID string
	Payload   map[string]interface{}
	Message   string // Content of the notification
	IsRead    bool
}

// WishlistAlertAdded logs when an order-related notification is sent.
type WishlistAlertAdded struct {
	ID         string
	AlertType  string
	WishlistID string
	UserID     string // ID of the user receiving the notification
	ProductID  string // ID of the order in the notification
	Payload    map[string]interface{}
	Message    string
	IsRead     bool
}

// MessageAlertAdded logs when an order-related notification is sent.
type MessageAlertAdded struct {
	ID        string
	AlertType string
	MessageID string
	UserID    string // ID of the user receiving the notification
	SenderID  string
	ProductID string // ID of the order in the notification
	Payload   map[string]interface{}
	Message   string // Content of the notification
	IsRead    bool
}

// InteractionAlertAdded logs when an order-related notification is sent.
type InteractionAlertAdded struct {
	ID        string
	AlertType string
	UserID    string // ID of the user receiving the notification
	ProductID string // ID of the order in the notification
	Payload   map[string]interface{}
	Message   string // Content of the notification
	IsRead    bool
}

// CommentAlertAdded logs when an order-related notification is sent.
type CommentAlertAdded struct {
	ID        string
	AlertType string
	CommentID string
	UserID    string // ID of the user receiving the notification
	SenderID  string
	ProductID string // ID of the order in the notification
	Payload   map[string]interface{}
	Message   string // Content of the notification
	IsRead    bool
}

// OfferAlertAdded logs when an order-related notification is sent.
type OfferAlertAdded struct {
	ID        string
	AlertType string
	OfferID   string
	UserID    string // ID of the user receiving the notification
	SenderID  string
	ProductID string // ID of the order in the notification
	Payload   map[string]interface{}
	Message   string // Content of the notification
	IsRead    bool
}

// SupportAlertAdded logs when an order-related notification is sent.
type SupportAlertAdded struct {
	ID        string
	AlertType string
	UserID    string // ID of the user receiving the notification
	TicketID  string
	Payload   map[string]interface{}
	Message   string // Content of the notification
	IsRead    bool
}

type AlertRead struct {
	ID     string
	IsRead bool
}

// ReviewAlertAdded logs when a review-related notification is sent.
type ReviewAlertAdded struct {
	ID        string
	AlertType string
	ReviewID  string
	UserID    string // ID of the user receiving the notification
	ProductID string // ID of the product being reviewed
	Payload   map[string]interface{}
	Message   string // Content of the notification
	IsRead    bool
}

// PaymentAlertAdded logs when a payment-related notification is sent.
type PaymentAlertAdded struct {
	ID        string
	AlertType string
	PaymentID string
	UserID    string // ID of the user receiving the notification
	OrderID   string // ID of the related order
	Payload   map[string]interface{}
	Message   string // Content of the notification
	IsRead    bool
}

// FollowingAlertAdded logs when a following-related notification is sent.
type FollowingAlertAdded struct {
	ID         string
	AlertType  string
	FollowerID string // ID of the user who followed
	UserID     string // ID of the user receiving the notification
	Payload    map[string]interface{}
	Message    string // Content of the notification
	IsRead     bool
}
