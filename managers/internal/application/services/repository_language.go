package services

import (
	"context"
	"middleman/managers/internal/models"
)

// RepositoryQuery represents a standardized query that can be executed against any repository
type RepositoryQuery struct {
	// Entity type being queried
	EntityType models.EntityType `json:"entity_type"`

	// Operation type
	Operation OperationType `json:"operation"`

	// Query parameters
	Parameters QueryParameters `json:"parameters"`

	// Pagination
	Pagination *PaginationParams `json:"pagination,omitempty"`

	// Sorting
	Sorting *SortingParams `json:"sorting,omitempty"`

	// Location-based filtering
	Location *LocationParams `json:"location,omitempty"`
}

// OperationType defines the type of repository operation
type OperationType string

const (
	OpSearch           OperationType = "search"
	OpFind             OperationType = "find"
	OpAdd              OperationType = "add"
	OpUpdate           OperationType = "update"
	OpRemove           OperationType = "remove"
	OnSend             OperationType = "send"
	OnReceive          OperationType = "receive"
	OpSuggest          OperationType = "suggest"
	OpFilter           OperationType = "filter"
	OpSearchByTerm     OperationType = "search_by_term"
	OpSearchByCategory OperationType = "search_by_category"
	OpGetMetrics       OperationType = "get_metrics"

	// Service-specific operations
	OpGetServices          OperationType = "get_services"
	OpGetCatalog           OperationType = "get_catalog"
	OpGetPublicCatalog     OperationType = "get_public_catalog"
	OpSearchByCategorySlug OperationType = "search_by_category_slug"
	OpUpdatePrice          OperationType = "update_price"
	OpIncreasePrice        OperationType = "increase_price"
	OpDecreasePrice        OperationType = "decrease_price"
	OpRebrand              OperationType = "rebrand"
	OpAdjustStock          OperationType = "adjust_stock"
	OpArchive              OperationType = "archive"
	OpMarkSold             OperationType = "mark_sold"
	OpMarkLeased           OperationType = "mark_leased"

	// Review-specific operations
	OpGetReviews                OperationType = "get_reviews"
	OpGetReviewsBySender        OperationType = "get_reviews_by_sender"
	OpGetApprovedReviews        OperationType = "get_approved_reviews"
	OpGetMostReviewed           OperationType = "get_most_reviewed"
	OpGetMostReviewedByCategory OperationType = "get_most_reviewed_by_category"
	OpApproveReview             OperationType = "approve_review"
	OpRejectReview              OperationType = "reject_review"
	OpFlagReview                OperationType = "flag_review"
	OpEditReview                OperationType = "edit_review"

	// Comment-specific operations
	OpGetComments                OperationType = "get_comments"
	OpGetCommentsBySender        OperationType = "get_comments_by_sender"
	OpGetApprovedComments        OperationType = "get_approved_comments"
	OpGetMostCommented           OperationType = "get_most_commented"
	OpGetMostCommentedByCategory OperationType = "get_most_commented_by_category"
	OpApproveComment             OperationType = "approve_comment"
	OpRejectComment              OperationType = "reject_comment"
	OpFlagComment                OperationType = "flag_comment"
	OpEditComment                OperationType = "edit_comment"

	// Property-specific operations
	OpGetProperties               OperationType = "get_properties"
	OpGetPropertiesByCategory     OperationType = "get_properties_by_category"
	OpGetPropertiesByCategorySlug OperationType = "get_properties_by_category_slug"
	OpGetPropertiesWithFilter     OperationType = "get_properties_with_filter"
	OpGetPropertyCatalog          OperationType = "get_property_catalog"
	OpGetPublicPropertyCatalog    OperationType = "get_public_property_catalog"
	OpUpdatePropertyPrice         OperationType = "update_property_price"
	OpRenameProperty              OperationType = "rename_property"
	OpAdjustPropertySquareFootage OperationType = "adjust_property_square_footage"
	OpArchiveProperty             OperationType = "archive_property"
	OpMarkPropertySold            OperationType = "mark_property_sold"
	OpMarkPropertyRent            OperationType = "mark_property_rent"
	OpIncreasePropertyPrice       OperationType = "increase_property_price"
	OpDecreasePropertyPrice       OperationType = "decrease_property_price"

	// Payment-specific operations
	OpAuthorizePayment      OperationType = "authorize_payment"
	OpConfirmPayment        OperationType = "confirm_payment"
	OpCapturePayment        OperationType = "capture_payment"
	OpCreateInvoice         OperationType = "create_invoice"
	OpAdjustInvoice         OperationType = "adjust_invoice"
	OpPayInvoice            OperationType = "pay_invoice"
	OpCancelInvoice         OperationType = "cancel_invoice"
	OpHandleWebhook         OperationType = "handle_webhook"
	OpGetPayment            OperationType = "get_payment"
	OpGetInvoice            OperationType = "get_invoice"
	OpGetPaymentsByCustomer OperationType = "get_payments_by_customer"
	OpGetInvoicesByOrder    OperationType = "get_invoices_by_order"
	OpSearchPayments        OperationType = "search_payments"
	OpSearchInvoices        OperationType = "search_invoices"

	// Order-specific operations
	OpCreateOrder         OperationType = "create_order"
	OpGetOrder            OperationType = "get_order"
	OpCancelOrder         OperationType = "cancel_order"
	OpReadyOrder          OperationType = "ready_order"
	OpCompleteOrder       OperationType = "complete_order"
	OpApproveOrder        OperationType = "approve_order"
	OpRejectOrder         OperationType = "reject_order"
	OpShipOrder           OperationType = "ship_order"
	OpDeliverOrder        OperationType = "deliver_order"
	OpGetOrdersByCustomer OperationType = "get_orders_by_customer"
	OpGetOrdersByStatus   OperationType = "get_orders_by_status"
	OpSearchOrders        OperationType = "search_orders"

	// Offer operations
	OpCreateOffer             OperationType = "create_offer"
	OpActivateOffer           OperationType = "activate_offer"
	OpCloseOffer              OperationType = "close_offer"
	OpAcceptOffer             OperationType = "accept_offer"
	OpGetOffer                OperationType = "get_offer"
	OpListOffers              OperationType = "list_offers"
	OpGetOffersByProduct      OperationType = "get_offers_by_product"
	OpGetOffersByUser         OperationType = "get_offers_by_user"
	OpSearchOffers            OperationType = "search_offers"
	OpRequestOfferNegotiation OperationType = "request_offer_negotiation"
	OpAcceptOfferNegotiation  OperationType = "accept_offer_negotiation"
	OpDeclineOfferNegotiation OperationType = "decline_offer_negotiation"

	// Notification operations
	OpListAlerts        OperationType = "list_alerts"
	OpGetAlertsByType   OperationType = "get_alerts_by_type"
	OpGetAlert          OperationType = "get_alert"
	OpGetAlertsByUser   OperationType = "get_alerts_by_user"
	OpGetUnreadAlerts   OperationType = "get_unread_alerts"
	OpGetReadAlerts     OperationType = "get_read_alerts"
	OpSearchAlerts      OperationType = "search_alerts"
	OpCountAlerts       OperationType = "count_alerts"
	OpCountUnreadAlerts OperationType = "count_unread_alerts"

	// Newsletter operations
	OpSubscribeNewsletter            OperationType = "subscribe_newsletter"
	OpUnsubscribeNewsletter          OperationType = "unsubscribe_newsletter"
	OpGetSubscription                OperationType = "get_subscription"
	OpListSubscriptions              OperationType = "list_subscriptions"
	OpUpdateSubscription             OperationType = "update_subscription"
	OpSendNewsletter                 OperationType = "send_newsletter"
	OpGetSubscriptionByID            OperationType = "get_subscription_by_id"
	OpGetSubscriptionsByUser         OperationType = "get_subscriptions_by_user"
	OpGetSubscriptionsByNewsletter   OperationType = "get_subscriptions_by_newsletter"
	OpGetActiveSubscriptions         OperationType = "get_active_subscriptions"
	OpGetInactiveSubscriptions       OperationType = "get_inactive_subscriptions"
	OpSearchSubscriptions            OperationType = "search_subscriptions"
	OpCountSubscriptions             OperationType = "count_subscriptions"
	OpCountActiveSubscriptions       OperationType = "count_active_subscriptions"
	OpCountSubscriptionsByNewsletter OperationType = "count_subscriptions_by_newsletter"
	OpGetNewsletter                  OperationType = "get_newsletter"
	OpGetNewsletters                 OperationType = "get_newsletters"
	OpSearchNewsletters              OperationType = "search_newsletters"

	// BuyNow operations
	OpCreateBuyNow             OperationType = "create_buy_now"
	OpConfirmBuyNow            OperationType = "confirm_buy_now"
	OpCancelBuyNow             OperationType = "cancel_buy_now"
	OpRequestBuyNowNegotiation OperationType = "request_buy_now_negotiation"
	OpAcceptBuyNowNegotiation  OperationType = "accept_buy_now_negotiation"
	OpDeclineBuyNowNegotiation OperationType = "decline_buy_now_negotiation"

	// Lease operations
	OpCreateLease             OperationType = "create_lease"
	OpStartLease              OperationType = "start_lease"
	OpMakeLeasePayment        OperationType = "make_lease_payment"
	OpExecuteLeaseBuyout      OperationType = "execute_lease_buyout"
	OpEndLease                OperationType = "end_lease"
	OpCancelLease             OperationType = "cancel_lease"
	OpDefaultLease            OperationType = "default_lease"
	OpGetActiveLeases         OperationType = "get_active_leases"
	OpRequestLeaseNegotiation OperationType = "request_lease_negotiation"
	OpAcceptLeaseNegotiation  OperationType = "accept_lease_negotiation"
	OpDeclineLeaseNegotiation OperationType = "decline_lease_negotiation"

	// BuyBack operations
	OpCreateBuyBack             OperationType = "create_buy_back"
	OpRedeemBuyBack             OperationType = "redeem_buy_back"
	OpExpireBuyBack             OperationType = "expire_buy_back"
	OpCancelBuyBack             OperationType = "cancel_buy_back"
	OpGetActiveBuyBacks         OperationType = "get_active_buy_backs"
	OpRequestBuyBackNegotiation OperationType = "request_buy_back_negotiation"
	OpAcceptBuyBackNegotiation  OperationType = "accept_buy_back_negotiation"
	OpDeclineBuyBackNegotiation OperationType = "decline_buy_back_negotiation"

	// Reservation operations
	OpCreateReservation             OperationType = "create_reservation"
	OpRedeemReservation             OperationType = "redeem_reservation"
	OpExpireReservation             OperationType = "expire_reservation"
	OpCancelReservation             OperationType = "cancel_reservation"
	OpGetActiveReservations         OperationType = "get_active_reservations"
	OpRequestReservationNegotiation OperationType = "request_reservation_negotiation"
	OpAcceptReservationNegotiation  OperationType = "accept_reservation_negotiation"

	// Metric operations
	OpUpdateItemMetric              OperationType = "update_item_metric"
	OpShareItem                     OperationType = "share_item"
	OpVisitItem                     OperationType = "visit_item"
	OpUpdateUserMetric              OperationType = "update_user_metric"
	OpGetUserMetric                 OperationType = "get_user_metric"
	OpGetItemMetric                 OperationType = "get_item_metric"
	OpGetItemsMetric                OperationType = "get_items_metric"
	OpGetHighestMetricsByType       OperationType = "get_highest_metrics_by_type"
	OpGetLowestMetricsByType        OperationType = "get_lowest_metrics_by_type"
	OpGetItemMetricByType           OperationType = "get_item_metric_by_type"
	OpGetUserMetricByType           OperationType = "get_user_metric_by_type"
	OpGetItemMetricsByCategory      OperationType = "get_item_metrics_by_category"
	OpGetUserMetricsByCategory      OperationType = "get_user_metrics_by_category"
	OpGetTopItemsByMetric           OperationType = "get_top_items_by_metric"
	OpGetTopUsersByMetric           OperationType = "get_top_users_by_metric"
	OpGetItemMetricsInRange         OperationType = "get_item_metrics_in_range"
	OpGetMetricsSummary             OperationType = "get_metrics_summary"
	OpSearchItemMetrics             OperationType = "search_item_metrics"
	OpSearchUserMetrics             OperationType = "search_user_metrics"
	OpGetRecentlyUpdatedMetrics     OperationType = "get_recently_updated_metrics"
	OpGetMetricsByPriceRange        OperationType = "get_metrics_by_price_range"
	OpGetMetricsByRatingRange       OperationType = "get_metrics_by_rating_range"
	OpGetTrendingItems              OperationType = "get_trending_items"
	OpGetActiveUsers                OperationType = "get_active_users"
	OpCompareItemMetrics            OperationType = "compare_item_metrics"
	OpCompareUserMetrics            OperationType = "compare_user_metrics"
	OpGetMetricsAnalytics           OperationType = "get_metrics_analytics"
	OpBulkUpdateItemMetrics         OperationType = "bulk_update_item_metrics"
	OpBulkUpdateUserMetrics         OperationType = "bulk_update_user_metrics"
	OpResetItemMetrics              OperationType = "reset_item_metrics"
	OpResetUserMetrics              OperationType = "reset_user_metrics"
	OpDeclineReservationNegotiation OperationType = "decline_reservation_negotiation"

	// Category-specific operations
	OpGetCategories        OperationType = "get_categories"
	OpGetMainCategories    OperationType = "get_main_categories"
	OpGetAllMainCategories OperationType = "get_all_main_categories"
	OpGetSubCategories     OperationType = "get_sub_categories"
	OpGetCategoryBySlug    OperationType = "get_category_by_slug"

	// Filter-specific operations
	OpAddFilter     OperationType = "add_filter"
	OpGetFilter     OperationType = "get_filter"
	OpGetFilters    OperationType = "get_filters"
	OpArchiveFilter OperationType = "archive_filter"
	OpRemoveFilter  OperationType = "remove_filter"

	// Following/Social operations
	OpApprove OperationType = "approve"
	OpReject  OperationType = "reject"

	// Additional operations
	OpGet         OperationType = "get"
	OpSend        OperationType = "send"
	OpValidate    OperationType = "validate"
	OpTrack       OperationType = "track"
	OpList        OperationType = "list"
	OpSubscribe   OperationType = "subscribe"
	OpUnsubscribe OperationType = "unsubscribe"
	OpCancel      OperationType = "cancel"
	OpCheckout    OperationType = "checkout"
	OpStart       OperationType = "start"
	OpComplete    OperationType = "complete"
	OpReady       OperationType = "ready"
	OpShip        OperationType = "ship"
	OpDeliver     OperationType = "deliver"
	OpCapture     OperationType = "capture"
	OpConfirm     OperationType = "confirm"
	OpPay         OperationType = "pay"
	OpAdjust      OperationType = "adjust"
	OpAuthorize   OperationType = "authorize"
	OpHandle      OperationType = "handle"

	OpUnknown OperationType = ""
)

// QueryParameters holds all possible query parameters
type QueryParameters struct {
	// Common identifiers
	ID                string   `json:"id,omitempty"`
	UserID            string   `json:"user_id,omitempty"`
	ItemID            string   `json:"item_id,omitempty"`
	ItemType          string   `json:"item_type,omitempty"`
	SearchTerm        string   `json:"search_term,omitempty"`
	Name              string   `json:"name,omitempty"`
	Description       string   `json:"description,omitempty"`
	CategoryID        string   `json:"category_id,omitempty"`
	CategorySlug      string   `json:"category_slug,omitempty"`
	GoogleCategoryID  string   `json:"google_category_id,omitempty"`
	CategoryType      string   `json:"category_type,omitempty"`
	SEOTitle          string   `json:"seo_title,omitempty"`
	SEOKeywords       []string `json:"seo_keywords,omitempty"`
	SEODescription    string   `json:"seo_description,omitempty"`
	MinPrice          int64    `json:"min_price,omitempty"`
	MaxPrice          int64    `json:"max_price,omitempty"`
	Brand             string   `json:"brand,omitempty"`
	Condition         string   `json:"condition,omitempty"`
	Model             string   `json:"model,omitempty"`
	Status            string   `json:"status,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	ManageStock       bool     `json:"manage_stock,omitempty"`
	MinStock          int64    `json:"min_stock,omitempty"`
	MaxStock          int64    `json:"max_stock,omitempty"`
	SKU               string   `json:"sku,omitempty"`
	UserType          string   `json:"user_type,omitempty"`
	Negotiable        bool     `json:"negotiable,omitempty"`
	HasVariants       bool     `json:"has_variants,omitempty"`
	FuelType          string   `json:"fuel_type,omitempty"`
	TransmissionType  string   `json:"transmission_type,omitempty"`
	MinYear           int64    `json:"min_year,omitempty"`
	MaxYear           int64    `json:"max_year,omitempty"`
	MinMileage        int64    `json:"min_mileage,omitempty"`
	MaxMileage        int64    `json:"max_mileage,omitempty"`
	AccidentFree      bool     `json:"accident_free,omitempty"`
	NumberOfOwners    int64    `json:"number_of_owners,omitempty"`
	VIN               string   `json:"vin,omitempty"`
	MinSalary         int64    `json:"min_salary,omitempty"`
	MaxSalary         int64    `json:"max_salary,omitempty"`
	EmploymentType    string   `json:"employment_type,omitempty"`
	SeniorityLevel    string   `json:"seniority_level,omitempty"`
	RelocationSupport bool     `json:"relocation_support,omitempty"`
	CompanyName       string   `json:"company_name,omitempty"`
	PropertyType      string   `json:"property_type,omitempty"`
	MinSquareFootage  int64    `json:"min_square_footage,omitempty"`
	MaxSquareFootage  int64    `json:"max_square_footage,omitempty"`
	MinBedrooms       int64    `json:"min_bedrooms,omitempty"`
	MaxBedrooms       int64    `json:"max_bedrooms,omitempty"`
	MinBathrooms      int64    `json:"min_bathrooms,omitempty"`
	MaxBathrooms      int64    `json:"max_bathrooms,omitempty"`
	MinYearBuilt      int64    `json:"min_year_built,omitempty"`
	MaxYearBuilt      int64    `json:"max_year_built,omitempty"`

	// Service-specific
	ServiceType      string   `json:"service_type,omitempty"`
	ProviderName     string   `json:"provider_name,omitempty"`
	Qualifications   []string `json:"qualifications,omitempty"`
	Availability     string   `json:"availability,omitempty"`
	DescriptionShort string   `json:"description_short,omitempty"`
	DescriptionLong  string   `json:"description_long,omitempty"`
	Contact          string   `json:"contact,omitempty"`
	FAQ              string   `json:"faq,omitempty"`
	ShippingCost     int64    `json:"shipping_cost,omitempty"`
	MiddlemanService bool     `json:"middleman_service,omitempty"`
	Attributes       []string `json:"attributes,omitempty"`
	Options          []string `json:"options,omitempty"`
	Pricing          []string `json:"pricing,omitempty"`
	BasePrice        int64    `json:"base_price,omitempty"`
	OldPrice         int64    `json:"old_price,omitempty"`
	NewPrice         int64    `json:"new_price,omitempty"`
	Price            int64    `json:"price,omitempty"`
	NewStock         int64    `json:"new_stock,omitempty"`
	MonthlyPrice     int64    `json:"monthly_price,omitempty"`
	LeaseTermMonths  int64    `json:"lease_term_months,omitempty"`
	Page             int64    `json:"page,omitempty"`
	PageSize         int64    `json:"page_size,omitempty"`
	SortBy           string   `json:"sort_by,omitempty"`
	SortOrder        string   `json:"sort_order,omitempty"`

	// Review-specific fields
	ReviewID     string `json:"review_id,omitempty"`
	ReviewStatus string `json:"review_status,omitempty"`
	Flagged      bool   `json:"flagged,omitempty"`
	Offset       int64  `json:"offset,omitempty"`
	Limit        int64  `json:"limit,omitempty"`

	// Comment-specific fields
	CommentID     string `json:"comment_id,omitempty"`
	CommentStatus string `json:"comment_status,omitempty"`

	// Deal-specific
	DealType     string `json:"deal_type,omitempty"`
	MerchantName string `json:"merchant_name,omitempty"`
	Rabatt       int64  `json:"rabatt,omitempty"`

	// Dimensions
	MinWeight int64 `json:"min_weight,omitempty"`
	MaxWeight int64 `json:"max_weight,omitempty"`
	MinHeight int64 `json:"min_height,omitempty"`
	MaxHeight int64 `json:"max_height,omitempty"`
	MinWidth  int64 `json:"min_width,omitempty"`
	MaxWidth  int64 `json:"max_width,omitempty"`
	MinDepth  int64 `json:"min_depth,omitempty"`
	MaxDepth  int64 `json:"max_depth,omitempty"`

	// Additional fields for specific operations
	Content        string `json:"content,omitempty"`
	ParentID       string `json:"parent_id,omitempty"`
	ReceiverID     string `json:"receiver_id,omitempty"`
	RecipientID    string `json:"recipient_id,omitempty"`
	SenderID       string `json:"sender_id,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	Body           string `json:"body,omitempty"`
	IsRead         bool   `json:"is_read,omitempty"`
	ActionType     string `json:"action_type,omitempty"`
	MetricType     string `json:"metric_type,omitempty"`

	// User-specific fields
	Email      string  `json:"email,omitempty"`
	Password   string  `json:"password,omitempty"`
	FirstName  string  `json:"first_name,omitempty"`
	LastName   string  `json:"last_name,omitempty"`
	Location   string  `json:"location,omitempty"`
	Lat        float64 `json:"lat,omitempty"`
	Lng        float64 `json:"lng,omitempty"`
	Thumbnail  string  `json:"thumbnail,omitempty"`
	Language   string  `json:"language,omitempty"`
	Bio        string  `json:"bio,omitempty"`
	Privacy    string  `json:"privacy,omitempty"`
	Background string  `json:"background,omitempty"`

	// Support-specific fields
	SupportID string `json:"support_id,omitempty"`
	Title     string `json:"title,omitempty"`

	// Shipping-specific fields
	ProductID       string `json:"product_id,omitempty"`
	TrackingNumber  string `json:"tracking_number,omitempty"`
	LabelURL        string `json:"label_url,omitempty"`
	SenderName      string `json:"sender_name,omitempty"`
	SenderAddress   string `json:"sender_address,omitempty"`
	ReceiverName    string `json:"receiver_name,omitempty"`
	ReceiverAddress string `json:"receiver_address,omitempty"`
	Weight          string `json:"weight,omitempty"`
	Dimensions      string `json:"dimensions,omitempty"`
	ServiceTypes    string `json:"service_types,omitempty"`

	// Additional Property-specific fields
	PropertyID            string   `json:"property_id,omitempty"`
	ListingType           string   `json:"listing_type,omitempty"`
	Address               string   `json:"address,omitempty"`
	City                  string   `json:"city,omitempty"`
	State                 string   `json:"state,omitempty"`
	PostalCode            string   `json:"postal_code,omitempty"`
	Country               string   `json:"country,omitempty"`
	LotSize               int64    `json:"lot_size,omitempty"`
	EstimatedClosingCosts int64    `json:"estimated_closing_costs,omitempty"`
	UrlReference          string   `json:"url_reference,omitempty"`
	TypeOfOffer           string   `json:"type_of_offer,omitempty"`
	PropertyTypes         []string `json:"property_types,omitempty"`
	ListingTypes          []string `json:"listing_types,omitempty"`
	Conditions            string   `json:"conditions,omitempty"`
	Statuses              []string `json:"statuses,omitempty"`
	UserTypes             []string `json:"user_types,omitempty"`
	CategoryIDs           []string `json:"category_ids,omitempty"`
	CategorySlugs         []string `json:"category_slugs,omitempty"`
	UserIDs               []string `json:"user_ids,omitempty"`
	RadiusMeters          float64  `json:"radius_meters,omitempty"`
	NewSquareFootage      int64    `json:"new_square_footage,omitempty"`

	// Payment-specific fields
	PaymentID       string `json:"payment_id,omitempty"`
	UserCustomerID  string `json:"user_customer_id,omitempty"`
	PaymentMethodID string `json:"payment_method_id,omitempty"`
	PaymentStatus   string `json:"payment_status,omitempty"`
	Amount          int64  `json:"amount,omitempty"`
	AmountToCapture int64  `json:"amount_to_capture,omitempty"`
	ClientSecret    string `json:"client_secret,omitempty"`
	InvoiceID       string `json:"invoice_id,omitempty"`
	OrderID         string `json:"order_id,omitempty"`
	InvoiceStatus   string `json:"invoice_status,omitempty"`
	Reason          string `json:"reason,omitempty"`
	RawBody         string `json:"raw_body,omitempty"`
	Signature       string `json:"signature,omitempty"`
	NewAmount       int64  `json:"new_amount,omitempty"`

	// Order-specific fields
	OrderStatus    string `json:"order_status,omitempty"`
	ShoppingID     string `json:"shopping_id,omitempty"`
	Items          string `json:"items,omitempty"`
	OrderReason    string `json:"order_reason,omitempty"`
	ReadyAt        string `json:"ready_at,omitempty"`
	CompletedAt    string `json:"completed_at,omitempty"`
	ApprovedAt     string `json:"approved_at,omitempty"`
	RejectedAt     string `json:"rejected_at,omitempty"`
	ShippedAt      string `json:"shipped_at,omitempty"`
	DeliveredAt    string `json:"delivered_at,omitempty"`
	OrderCreatedAt string `json:"order_created_at,omitempty"`
	OrderUpdatedAt string `json:"order_updated_at,omitempty"`

	// Offer-specific fields
	OfferID           string `json:"offer_id,omitempty"`
	UserSellerID      string `json:"user_seller_id,omitempty"`
	OfferStatus       string `json:"offer_status,omitempty"`
	ProposedPrice     int64  `json:"proposed_price,omitempty"`
	FinalPrice        int64  `json:"final_price,omitempty"`
	Message           string `json:"message,omitempty"`
	NegotiationStatus string `json:"negotiation_status,omitempty"`

	// BuyNow-specific fields
	BuyNowID     string `json:"buy_now_id,omitempty"`
	BuyNowStatus string `json:"buy_now_status,omitempty"`

	// Lease-specific fields
	LeaseID              string `json:"lease_id,omitempty"`
	LeaseStatus          string `json:"lease_status,omitempty"`
	ProposedMonthlyPrice int64  `json:"proposed_monthly_price,omitempty"`
	ProposedTermMonths   int64  `json:"proposed_term_months,omitempty"`
	FinalMonthlyPrice    int64  `json:"final_monthly_price,omitempty"`
	FinalTermMonths      int64  `json:"final_term_months,omitempty"`
	HasBuyout            bool   `json:"has_buyout,omitempty"`
	BuyoutPrice          int64  `json:"buyout_price,omitempty"`
	BuyoutAmount         int64  `json:"buyout_amount,omitempty"`

	// BuyBack-specific fields
	BuyBackID           string `json:"buy_back_id,omitempty"`
	BuyBackStatus       string `json:"buy_back_status,omitempty"`
	LockedPrice         int64  `json:"locked_price,omitempty"`
	RedemptionFee       int64  `json:"redemption_fee,omitempty"`
	LockDurationDays    int32  `json:"lock_duration_days,omitempty"`
	LockBuyerID         string `json:"lock_buyer_id,omitempty"`
	NewLockedPrice      int64  `json:"new_locked_price,omitempty"`
	NewRedemptionFee    int64  `json:"new_redemption_fee,omitempty"`
	AgreedLockedPrice   int64  `json:"agreed_locked_price,omitempty"`
	AgreedRedemptionFee int64  `json:"agreed_redemption_fee,omitempty"`

	// Reservation-specific fields
	ReservationID        string `json:"reservation_id,omitempty"`
	ReservationStatus    string `json:"reservation_status,omitempty"`
	ReservationFee       int64  `json:"reservation_fee,omitempty"`
	NewReservationFee    int64  `json:"new_reservation_fee,omitempty"`
	AgreedReservationFee int64  `json:"agreed_reservation_fee,omitempty"`
	IsFree               bool   `json:"is_free,omitempty"`

	// Notification-specific fields
	AlertID     string            `json:"alert_id,omitempty"`
	AlertType   string            `json:"alert_type,omitempty"`
	AlertStatus string            `json:"alert_status,omitempty"`
	IsReadAlert bool              `json:"is_read_alert,omitempty"`
	Payload     map[string]string `json:"payload,omitempty"`
	Query       string            `json:"query,omitempty"`
	Count       int64             `json:"count,omitempty"`

	// Newsletter-specific fields
	SubscriptionID          string `json:"subscription_id,omitempty"`
	NewsletterID            string `json:"newsletter_id,omitempty"`
	SubscriptionPreferences string `json:"subscription_preferences,omitempty"`
	SubscriptionStatus      string `json:"subscription_status,omitempty"`
	NewsletterContent       string `json:"newsletter_content,omitempty"`
	NewsletterTitle         string `json:"newsletter_title,omitempty"`
	NewsletterStatus        string `json:"newsletter_status,omitempty"`

	// Metric-specific fields
	MetricTypeAction string   `json:"metric_type_action,omitempty"`
	EntityTypes      []string `json:"entity_types,omitempty"`
	MinRating        int64    `json:"min_rating,omitempty"`
	MaxRating        int64    `json:"max_rating,omitempty"`
	Days             int32    `json:"days,omitempty"`
	TimeRange        string   `json:"time_range,omitempty"`
	ItemID1          string   `json:"item_id1,omitempty"`
	ItemID2          string   `json:"item_id2,omitempty"`
	UserID1          string   `json:"user_id1,omitempty"`
	UserID2          string   `json:"user_id2,omitempty"`
	MetricTypes      []string `json:"metric_types,omitempty"`
	Updates          []string `json:"updates,omitempty"`
	Radius           float32  `json:"radius,omitempty"`
	LatFloat         float32  `json:"lat_float,omitempty"`
	LngFloat         float32  `json:"lng_float,omitempty"`

	// Following-specific fields
	FollowedUserID   string `json:"followed_user_id,omitempty"`
	FollowedUserType string `json:"followed_user_type,omitempty"`

	// Filter-specific fields
	FilterType    string   `json:"filter_type,omitempty"`
	FilterOptions []string `json:"filter_options,omitempty"`
	Required      bool     `json:"required,omitempty"`
	Searchable    bool     `json:"searchable,omitempty"`

	// Basket-specific fields
	BasketID string `json:"basket_id,omitempty"`
	Quantity int64  `json:"quantity,omitempty"`

	// Subject for emails/messages
	Subject string `json:"subject,omitempty"`
}

// PaginationParams defines pagination parameters
type PaginationParams struct {
	Offset   int64 `json:"offset,omitempty"`
	Limit    int64 `json:"limit,omitempty"`
	Page     int64 `json:"page,omitempty"`
	PageSize int64 `json:"page_size,omitempty"`
}

// SortingParams defines sorting parameters
type SortingParams struct {
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// LocationParams defines location-based filtering
type LocationParams struct {
	Lat    float64 `json:"lat,omitempty"`
	Lng    float64 `json:"lng,omitempty"`
	Radius float64 `json:"radius,omitempty"`
}

// RepositoryResponse represents a standardized response from any repository
type RepositoryResponse struct {
	// Success indicator
	Success bool `json:"success"`

	// Error information
	Error string `json:"error,omitempty"`

	// Result data
	Data interface{} `json:"data,omitempty"`

	// Metadata
	Metadata ResponseMetadata `json:"metadata"`
}

// ResponseMetadata contains metadata about the response
type ResponseMetadata struct {
	EntityType    models.EntityType `json:"entity_type"`
	Operation     OperationType     `json:"operation"`
	TotalCount    int64             `json:"total_count,omitempty"`
	ReturnedCount int               `json:"returned_count,omitempty"`
	Page          int64             `json:"page,omitempty"`
	PageSize      int64             `json:"page_size,omitempty"`
	ExecutionTime int64             `json:"execution_time_ms,omitempty"`
}

// UnifiedRepositoryInterface provides a unified interface for all repository operations
type UnifiedRepositoryInterface interface {
	// Execute performs a repository operation based on the query
	Execute(ctx context.Context, query RepositoryQuery) (*RepositoryResponse, error)

	// ValidateQuery validates if a query is valid for the given entity type
	ValidateQuery(query RepositoryQuery) error

	// GetSupportedOperations returns supported operations for an entity type
	GetSupportedOperations(entityType models.EntityType) []OperationType

	// TranslateAIRequest translates an AI request into a repository query
	TranslateAIRequest(aiRequest map[string]interface{}) (*RepositoryQuery, error)
}

// RepositoryLanguageProcessor processes natural language into repository queries
type RepositoryLanguageProcessor interface {
	// ParseNaturalLanguage converts natural language into a repository query
	ParseNaturalLanguage(ctx context.Context, input string) (*RepositoryQuery, error)

	// GenerateNaturalResponse converts repository response into natural language
	GenerateNaturalResponse(ctx context.Context, response *RepositoryResponse) (string, error)

	// SuggestQueries suggests possible queries based on partial input
	SuggestQueries(ctx context.Context, partialInput string) ([]RepositoryQuery, error)
}
