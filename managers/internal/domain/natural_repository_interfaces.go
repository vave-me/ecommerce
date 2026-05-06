package domain

import (
	"context"
	"middleman/managers/internal/models"
	"time"
)

// NaturalUserRepository provides natural language methods for user operations
type NaturalUserRepository interface {
	UserRepository
	
	// Natural language methods optimized for LLM communication
	FindPersonByEmail(ctx context.Context, email string) (*models.User, error)
	FindPersonByUsername(ctx context.Context, username string) (*models.User, error)
	GetPersonDetails(ctx context.Context, personID string) (*models.User, error)
	CreateNewPerson(ctx context.Context, email, password, name, firstName, lastName, location string, lat, lng float32, profilePicture, language string) (string, error)
	AuthenticatePerson(ctx context.Context, email, password string) (*models.LoginResponse, error)
	UpdatePersonProfile(ctx context.Context, personID string, updates map[string]interface{}) error
	EnablePersonAccount(ctx context.Context, personID, verificationToken string) error
	DisablePersonAccount(ctx context.Context, personID string) error
}

// NaturalProductRepository provides natural language methods for product operations
type NaturalProductRepository interface {
	ProductRepository
	
	// Natural language methods (excluding methods that conflict with base ProductRepository)
	SearchProductsByDescription(ctx context.Context, description string, limit int) ([]*models.Product, error)
	FindProductsUnderPrice(ctx context.Context, maxPrice float64, limit int) ([]*models.Product, error)
	FindProductsByBrand(ctx context.Context, brand string, limit int) ([]*models.Product, error)
	FindProductsInCategory(ctx context.Context, categoryName string, limit int) ([]*models.Product, error)
	GetProductDetails(ctx context.Context, productID string) (*models.Product, error)
	CreateNewProduct(ctx context.Context, details map[string]interface{}) (string, error)
	// UpdateProductPrice is already defined in ProductRepository
	UpdateProductStock(ctx context.Context, productID string, newStock int) error
	MarkProductAsSold(ctx context.Context, productID string) error
	GetTrendingProducts(ctx context.Context, limit int) ([]*models.Product, error)
	GetProductsByCondition(ctx context.Context, condition string, limit int) ([]*models.Product, error)
}

// NaturalOrderRepository provides natural language methods for order operations
type NaturalOrderRepository interface {
	OrderRepository
	
	// Natural language methods
	CreatePurchaseOrder(ctx context.Context, customerID string, items []interface{}) (*models.Order, error)
	GetOrderDetails(ctx context.Context, orderID string) (*models.Order, error)
	GetCustomerOrders(ctx context.Context, customerID string, limit int) ([]*models.Order, error)
	GetPendingOrders(ctx context.Context, limit int) ([]*models.Order, error)
	GetOrdersByDateRange(ctx context.Context, startDate, endDate time.Time, limit int) ([]*models.Order, error)
	MarkOrderAsShipped(ctx context.Context, orderID string, trackingNumber string) error
	MarkOrderAsDelivered(ctx context.Context, orderID string) error
	CancelOrder(ctx context.Context, orderID string, reason string) error
	GetOrdersAwaitingPayment(ctx context.Context, limit int) ([]*models.Order, error)
	GetOrdersReadyToShip(ctx context.Context, limit int) ([]*models.Order, error)
}

// NaturalPaymentRepository provides natural language methods for payment operations
type NaturalPaymentRepository interface {
	PaymentRepository
	
	// Natural language methods (excluding methods that conflict with base PaymentRepository)
	ProcessPaymentForOrder(ctx context.Context, orderID string, amount float64, method string) (*models.Payment, error)
	GetPaymentDetails(ctx context.Context, paymentID string) (*models.Payment, error)
	GetCustomerPayments(ctx context.Context, customerID string, limit int) ([]*models.Payment, error)
	RefundPayment(ctx context.Context, paymentID string, amount float64, reason string) error
	GetPendingPayments(ctx context.Context, limit int) ([]*models.Payment, error)
	// ConfirmPayment is already defined in PaymentRepository
	GetPaymentsByDateRange(ctx context.Context, startDate, endDate time.Time, limit int) ([]*models.Payment, error)
	CreateInvoiceForOrder(ctx context.Context, orderID string) (*models.Invoice, error)
	SendPaymentReceipt(ctx context.Context, paymentID string, email string) error
}

// NaturalBasketRepository provides natural language methods for shopping cart operations
type NaturalBasketRepository interface {
	BasketRepository
	
	// Natural language methods
	AddProductToCart(ctx context.Context, customerID, productID string, quantity int) error
	RemoveProductFromCart(ctx context.Context, customerID, productID string) error
	GetShoppingCart(ctx context.Context, customerID string) (*models.Basket, error)
	UpdateCartItemQuantity(ctx context.Context, customerID, productID string, newQuantity int) error
	EmptyShoppingCart(ctx context.Context, customerID string) error
	GetCartTotal(ctx context.Context, customerID string) (float64, error)
	ConvertCartToOrder(ctx context.Context, customerID string) (string, error)
	SaveCartForLater(ctx context.Context, customerID string) error
	GetAbandonedCarts(ctx context.Context, hours int) ([]*models.Basket, error)
}

// NaturalNotificationRepository provides natural language methods for notifications
type NaturalNotificationRepository interface {
	NotificationRepository
	
	// Natural language methods
	NotifyPerson(ctx context.Context, personID string, title, message, notificationType string) error
	NotifyAboutOrder(ctx context.Context, personID string, orderID string, status string) error
	NotifyAboutPayment(ctx context.Context, personID string, paymentID string, status string) error
	NotifyAboutNewMessage(ctx context.Context, personID string, senderName string) error
	GetPersonNotifications(ctx context.Context, personID string, limit int) ([]*models.Alert, error)
	GetUnreadNotifications(ctx context.Context, personID string) ([]*models.Alert, error)
	MarkNotificationAsRead(ctx context.Context, notificationID string) error
	MarkAllNotificationsAsRead(ctx context.Context, personID string) error
	SendBulkNotification(ctx context.Context, personIDs []string, title, message string) error
}

// NaturalMessageRepository provides natural language methods for messaging
type NaturalMessageRepository interface {
	MessagesRepository
	
	// Natural language methods
	SendMessageToPerson(ctx context.Context, senderID, recipientID, message string) error
	GetConversationBetween(ctx context.Context, person1ID, person2ID string, limit int) ([]*models.Message, error)
	GetPersonMessages(ctx context.Context, personID string, limit int) ([]*models.Message, error)
	GetUnreadMessages(ctx context.Context, personID string) ([]*models.Message, error)
	// MarkMessageAsRead is already defined in MessagesRepository
	DeleteMessage(ctx context.Context, messageID string) error
	SearchMessagesContaining(ctx context.Context, personID, searchText string, limit int) ([]*models.Message, error)
	GetMessageThread(ctx context.Context, messageID string) ([]*models.Message, error)
	SendBulkMessage(ctx context.Context, senderID string, recipientIDs []string, message string) error
}

// NaturalSupportRepository provides natural language methods for support operations
type NaturalSupportRepository interface {
	SupportRepository
	
	// Natural language methods
	CreateSupportTicket(ctx context.Context, customerID, subject, description, priority, category string) (string, error)
	GetTicketDetails(ctx context.Context, ticketID string) (*models.Ticket, error)
	GetCustomerTickets(ctx context.Context, customerID string, limit int) ([]*models.Ticket, error)
	UpdateTicketStatus(ctx context.Context, ticketID, status string) error
	AddTicketResponse(ctx context.Context, ticketID, responderID, response string) error
	// GetOpenTickets is already defined in SupportRepository
	GetHighPriorityTickets(ctx context.Context, limit int) ([]*models.Ticket, error)
	AssignTicketToAgent(ctx context.Context, ticketID, agentID string) error
	EscalateTicket(ctx context.Context, ticketID string, reason string) error
	ResolveTicket(ctx context.Context, ticketID string, resolution string) error
}

// NaturalReviewRepository provides natural language methods for reviews
type NaturalReviewRepository interface {
	ReviewRepository
	
	// Natural language methods
	AddProductReview(ctx context.Context, productID, reviewerID string, rating int, comment string) error
	GetProductReviews(ctx context.Context, productID string, limit int) ([]*models.Review, error)
	GetReviewerHistory(ctx context.Context, reviewerID string, limit int) ([]*models.Review, error)
	GetHighRatedProducts(ctx context.Context, minRating int, limit int) ([]*models.Product, error)
	// ApproveReview is already defined in ReviewRepository
	// RejectReview is already defined in ReviewRepository
	FlagInappropriateReview(ctx context.Context, reviewID string, reason string) error
	GetPendingReviews(ctx context.Context, limit int) ([]*models.Review, error)
	GetAverageRating(ctx context.Context, productID string) (float64, error)
	GetMostHelpfulReviews(ctx context.Context, productID string, limit int) ([]*models.Review, error)
}

// NaturalCategoryRepository provides natural language methods for categories
type NaturalCategoryRepository interface {
	CategoryRepository
	
	// Natural language methods
	// GetAllMainCategories is already defined in CategoryRepository
	GetSubcategoriesOf(ctx context.Context, parentCategoryName string) ([]*models.Category, error)
	FindCategoryByName(ctx context.Context, name string) (*models.Category, error)
	GetCategoryPath(ctx context.Context, categoryID string) ([]string, error)
	GetPopularCategories(ctx context.Context, limit int) ([]*models.Category, error)
	// CreateNewCategory is already defined in CategoryRepository
	UpdateCategoryName(ctx context.Context, categoryID, newName string) error
	GetCategoryProductCount(ctx context.Context, categoryID string) (int64, error)
	SearchCategories(ctx context.Context, searchTerm string) ([]*models.Category, error)
}

// NaturalWishlistRepository provides natural language methods for wishlists
type NaturalWishlistRepository interface {
	WishlistRepository
	
	// Natural language methods
	AddToWishlist(ctx context.Context, personID, productID string) error
	// RemoveFromWishlist is already defined in WishlistRepository
	GetPersonWishlist(ctx context.Context, personID string) ([]*models.Product, error)
	CheckIfInWishlist(ctx context.Context, personID, productID string) (bool, error)
	GetMostWishlisted(ctx context.Context, limit int) ([]*models.Product, error)
	NotifyWhenAvailable(ctx context.Context, personID, productID string) error
	ShareWishlist(ctx context.Context, personID string, recipientEmails []string) error
	// ClearWishlist is already defined in WishlistRepository
	// GetWishlistCount is already defined in WishlistRepository
}

// NaturalMetricRepository provides natural language methods for analytics
type NaturalMetricRepository interface {
	MetricRepository
	
	// Natural language methods
	RecordProductView(ctx context.Context, productID, viewerID string) error
	RecordProductPurchase(ctx context.Context, productID, buyerID string, amount float64) error
	GetProductMetrics(ctx context.Context, productID string) (*models.ItemMetric, error)
	GetTopSellingProducts(ctx context.Context, days int, limit int) ([]*models.Product, error)
	GetMostViewedProducts(ctx context.Context, days int, limit int) ([]*models.Product, error)
	GetUserEngagementMetrics(ctx context.Context, userID string) (*models.UserMetric, error)
	GetSalesMetricsByDateRange(ctx context.Context, startDate, endDate time.Time) (*models.SalesMetrics, error)
	GetConversionRate(ctx context.Context, productID string, days int) (float64, error)
	GetAverageOrderValue(ctx context.Context, days int) (float64, error)
	GenerateDashboardMetrics(ctx context.Context) (*models.DashboardMetrics, error)
}

// NaturalVectorRepository provides natural language methods for AI search
type NaturalVectorRepository interface {
	VectorRepository
	
	// Natural language methods
	SearchProductsWithAI(ctx context.Context, description string, limit int) ([]*models.Product, error)
	FindSimilarProducts(ctx context.Context, productID string, limit int) ([]*models.Product, error)
	SearchPeopleByInterests(ctx context.Context, interests string, limit int) ([]*models.User, error)
	SearchContentByMeaning(ctx context.Context, query string, contentType string, limit int) ([]interface{}, error)
	GetRecommendedProducts(ctx context.Context, personID string, limit int) ([]*models.Product, error)
	IndexNewContent(ctx context.Context, contentID, contentType, content string) error
	UpdateContentEmbedding(ctx context.Context, contentID, contentType string) error
	FindExpertiseMatch(ctx context.Context, skillDescription string, limit int) ([]*models.User, error)
	SearchByImage(ctx context.Context, imageURL string, limit int) ([]*models.Product, error)
}