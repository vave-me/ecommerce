package models

import "time"

// Offer represents an offer for a product
type Offer struct {
	ID             string    `json:"id"`
	UserSellerID   string    `json:"user_seller_id"`
	UserCustomerID string    `json:"user_customer_id"`
	ProductID      string    `json:"product_id"`
	Price          int64     `json:"price"`
	OfferStatus    string    `json:"offer_status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BuyNow represents a buy-now transaction
type BuyNow struct {
	ID          string     `json:"id"`
	OfferID     string     `json:"offer_id"`
	FinalPrice  int64      `json:"final_price"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
	CanceledAt  *time.Time `json:"canceled_at,omitempty"`
}

// Lease represents a lease agreement
type Lease struct {
	ID              string     `json:"id"`
	OfferID         string     `json:"offer_id"`
	MonthlyPrice    int64      `json:"monthly_price"`
	LeaseTermMonths int64      `json:"lease_term_months"`
	HasBuyout       bool       `json:"has_buyout"`
	BuyoutPrice     int64      `json:"buyout_price"`
	LeaseStatus     string     `json:"lease_status"`
	CreatedAt       time.Time  `json:"created_at"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
	CanceledAt      *time.Time `json:"canceled_at,omitempty"`
	DefaultedAt     *time.Time `json:"defaulted_at,omitempty"`
	ExecutedAt      *time.Time `json:"executed_at,omitempty"`
}

// BuyBack represents a buy-back agreement
type BuyBack struct {
	ID               string     `json:"id"`
	OfferID          string     `json:"offer_id"`
	LockedPrice      int64      `json:"locked_price"`
	RedemptionFee    int64      `json:"redemption_fee"`
	LockDurationDays int32      `json:"lock_duration_days"`
	LockBuyerID      string     `json:"lock_buyer_id"`
	BuyBackStatus    string     `json:"buyback_status"`
	CreatedAt        time.Time  `json:"created_at"`
	RedeemedAt       *time.Time `json:"redeemed_at,omitempty"`
	ExpiredAt        *time.Time `json:"expired_at,omitempty"`
	CanceledAt       *time.Time `json:"canceled_at,omitempty"`
}

// Reservation represents a reservation agreement
type Reservation struct {
	ID                string     `json:"id"`
	OfferID           string     `json:"offer_id"`
	LockedPrice       int64      `json:"locked_price"`
	ReservationFee    int64      `json:"reservation_fee"`
	LockDurationDays  int32      `json:"lock_duration_days"`
	LockBuyerID       string     `json:"lock_buyer_id"`
	ReservationStatus string     `json:"reservation_status"`
	IsFree            bool       `json:"is_free"`
	CreatedAt         time.Time  `json:"created_at"`
	RedeemedAt        *time.Time `json:"redeemed_at,omitempty"`
	ExpiredAt         *time.Time `json:"expired_at,omitempty"`
	CanceledAt        *time.Time `json:"canceled_at,omitempty"`
}

// Response models for Offer operations
type CreateOfferResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type ActivateOfferResponse struct {
	OfferID     string `json:"offer_id"`
	OfferStatus string `json:"offer_status"`
}

type CloseOfferResponse struct {
	OfferID     string `json:"offer_id"`
	OfferStatus string `json:"offer_status"`
}

type AcceptOfferResponse struct {
	OfferID     string `json:"offer_id"`
	OfferStatus string `json:"offer_status"`
}

type GetOfferResponse struct {
	Offer Offer `json:"offer"`
}

type ListOffersResponse struct {
	Offers []Offer `json:"offers"`
	Total  int64   `json:"total"`
	Page   int64   `json:"page"`
	Limit  int64   `json:"limit"`
}

// BuyNow response models
type CreateBuyNowResponse struct {
	BuyNowID  string    `json:"buy_now_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ConfirmBuyNowResponse struct {
	BuyNowID    string    `json:"buy_now_id"`
	Status      string    `json:"status"`
	ConfirmedAt time.Time `json:"confirmed_at"`
}

type CancelBuyNowResponse struct {
	BuyNowID     string    `json:"buy_now_id"`
	BuyNowStatus string    `json:"buy_now_status"`
	CanceledAt   time.Time `json:"canceled_at"`
}

// Lease response models
type CreateLeaseResponse struct {
	LeaseID     string    `json:"lease_id"`
	LeaseStatus string    `json:"lease_status"`
	CreatedAt   time.Time `json:"created_at"`
}

type StartLeaseResponse struct {
	LeaseID     string    `json:"lease_id"`
	LeaseStatus string    `json:"lease_status"`
	StartedAt   time.Time `json:"started_at"`
	EndDate     time.Time `json:"end_date"`
}

type MakeLeasePaymentResponse struct {
	LeaseID     string    `json:"lease_id"`
	PaymentDate time.Time `json:"payment_date"`
}

type ExecuteLeaseBuyoutResponse struct {
	LeaseID     string    `json:"lease_id"`
	LeaseStatus string    `json:"lease_status"`
	ExecutedAt  time.Time `json:"executed_at"`
}

type EndLeaseResponse struct {
	LeaseID     string    `json:"lease_id"`
	LeaseStatus string    `json:"lease_status"`
	EndedAt     time.Time `json:"ended_at"`
}

type CancelLeaseResponse struct {
	LeaseID     string    `json:"lease_id"`
	LeaseStatus string    `json:"lease_status"`
	CanceledAt  time.Time `json:"canceled_at"`
}

type DefaultLeaseResponse struct {
	LeaseID     string    `json:"lease_id"`
	LeaseStatus string    `json:"lease_status"`
	DefaultedAt time.Time `json:"defaulted_at"`
}

// BuyBack response models
type CreateBuyBackResponse struct {
	BuyBackID     string    `json:"buyback_id"`
	BuyBackStatus string    `json:"buyback_status"`
	CreatedAt     time.Time `json:"created_at"`
}

type RedeemBuyBackResponse struct {
	BuyBackID     string    `json:"buyback_id"`
	BuyBackStatus string    `json:"buyback_status"`
	RedeemedAt    time.Time `json:"redeemed_at"`
}

type ExpireBuyBackResponse struct {
	BuyBackID     string    `json:"buyback_id"`
	BuyBackStatus string    `json:"buyback_status"`
	ExpiredAt     time.Time `json:"expired_at"`
}

type CancelBuyBackResponse struct {
	BuyBackID     string    `json:"buyback_id"`
	BuyBackStatus string    `json:"buyback_status"`
	CanceledAt    time.Time `json:"canceled_at"`
}

// Reservation response models
type CreateReservationResponse struct {
	ReservationID     string    `json:"reservation_id"`
	ReservationStatus string    `json:"reservation_status"`
	CreatedAt         time.Time `json:"created_at"`
	IsFree            bool      `json:"is_free"`
}

type RedeemReservationResponse struct {
	ReservationID     string    `json:"reservation_id"`
	ReservationStatus string    `json:"reservation_status"`
	RedeemedAt        time.Time `json:"redeemed_at"`
}

type ExpireReservationResponse struct {
	ReservationID     string    `json:"reservation_id"`
	ReservationStatus string    `json:"reservation_status"`
	ExpiredAt         time.Time `json:"expired_at"`
}

type CancelReservationResponse struct {
	ReservationID     string    `json:"reservation_id"`
	ReservationStatus string    `json:"reservation_status"`
	CanceledAt        time.Time `json:"canceled_at"`
}

// Negotiation response models
type NegotiationResponse struct {
	ID                string `json:"id"`
	NegotiationStatus string `json:"negotiation_status"`
}

// Status constants
const (
	// Offer statuses
	OfferStatusDraft    = "draft"
	OfferStatusActive   = "active"
	OfferStatusAccepted = "accepted"
	OfferStatusClosed   = "closed"

	// BuyNow statuses
	BuyNowStatusPending   = "pending"
	BuyNowStatusConfirmed = "confirmed"
	BuyNowStatusCanceled  = "canceled"

	// Lease statuses
	LeaseStatusCreated   = "created"
	LeaseStatusActive    = "active"
	LeaseStatusEnded     = "ended"
	LeaseStatusCanceled  = "canceled"
	LeaseStatusDefaulted = "defaulted"
	LeaseStatusBoughtOut = "bought_out"

	// BuyBack statuses
	BuyBackStatusActive   = "active"
	BuyBackStatusRedeemed = "redeemed"
	BuyBackStatusExpired  = "expired"
	BuyBackStatusCanceled = "canceled"

	// Reservation statuses
	ReservationStatusActive   = "active"
	ReservationStatusRedeemed = "redeemed"
	ReservationStatusExpired  = "expired"
	ReservationStatusCanceled = "canceled"

	// Negotiation statuses
	NegotiationStatusRequested = "requested"
	NegotiationStatusAccepted  = "accepted"
	NegotiationStatusDeclined  = "declined"
)
