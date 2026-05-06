package domain

import (
	"time"
)

const (
	UserAuthorizedEvent             = "users.UserAuthorized"
	UserCreatedEvent                = "users.UserCreated"
	UserUpdatedEvent                = "users.UserUpdated"
	UserEnabledEvent                = "users.UserEnabled"
	UserDisabledEvent               = "users.UserDisabled"
	UserRenamedEvent                = "users.UserRenamed"
	UserLoggedInEvent               = "users.UserLoggedIn"
	UserLoggedOutEvent              = "users.UserLoggedOut"
	UserKYCVerifiedEvent            = "users.UserKYCVerified"
	UserKYCRejectedEvent            = "users.UserKYCRejectedEvent"
	UserLocationAddedEvent          = "users.UserLocationAddedEvent"
	UserLocationUpdatedEvent        = "users.UserLocationUpdatedEvent"
	UserLocationRemovedEvent        = "users.UserLocationRemovedEvent"
	UserPasswordResetRequestedEvent = "users.UserPasswordResetRequested"
	UserPasswordResetEvent          = "users.UserPasswordReset"
	UserPasswordForgotEvent         = "users.UserPasswordForgot"
	UserTokenInvalidatedEvent       = "users.UserTokenInvalidated"
	UserTokenRefreshedEvent         = "users.UserTokenRefreshed"
	UserTOTPSetupInitiatedEvent     = "users.UserTOTPSetupInitiated"
	UserTOTPEnabledEvent            = "users.UserTOTPEnabled"
	UserTOTPDisabledEvent           = "users.UserTOTPDisabled"
	UserTOTPVerifiedEvent           = "users.UserTOTPVerified"
	UserBackupCodeUsedEvent         = "users.UserBackupCodeUsed"
	UserSSOLinkedEvent              = "users.UserSSOLinked"
)

type UserAuthorized struct {
	Id string
}

// Key implements registry.Registerable
func (UserAuthorized) Key() string { return UserAuthorizedEvent }

type UserCreated struct {
	Email             string
	Password          string
	Username          string
	Firstname         string
	Lastname          string
	Location          string
	Enabled           bool
	GoogleID          string
	Lat               float64
	Lng               float64
	Thumbnail         string
	VerificationToken string
	Role              UserRole
}

type UserUpdated struct {
	ID         string
	Username   string
	Bio        string
	Privacy    string
	Background string
	FirstName  string
	LastName   string
	Latitude   float64
	Longitude  float64
	Thumbnail  string
	Role       UserRole
}

// Key implements registry.Registerable
func (UserCreated) Key() string { return UserCreatedEvent }
func (UserUpdated) Key() string { return UserUpdatedEvent }

type UserEnabledToggled struct {
	Enabled bool
}

type UserRenamed struct {
	FirstName string
}

// Key implements registry.Registerable
func (UserRenamed) Key() string { return UserRenamedEvent }

type UserLoggedIn struct {
	UserID string
	Token  string
}

func (UserLoggedIn) Key() string { return UserLoggedInEvent }

type UserLoggedOut struct {
	UserID string
	Token  string
}

func (UserLoggedOut) Key() string { return UserLoggedOutEvent }

type UserKYCVerified struct {
	UserID string
	Token  string
}

func (UserKYCVerified) Key() string { return UserKYCVerifiedEvent }

// Example event payload for KYC rejected
type UserKYCRejected struct {
	OfferID      string
	RejectReason string
}

func (UserKYCRejected) Key() string { return UserKYCRejectedEvent }

// Example event payload for KYC rejected
type UserLocationAdded struct {
	Lat string
	Lng string
}

func (UserLocationAdded) Key() string { return UserLocationAddedEvent }

type UserLocationRemoved struct {
	Lat string
	Lng string
}

func (UserLocationRemoved) Key() string { return UserLocationRemovedEvent }

type UserLocationUpdated struct {
	Lat string
	Lng string
}

func (UserLocationUpdated) Key() string { return UserLocationUpdatedEvent }

type UserPasswordResetRequested struct {
	Token     string
	Email     string
	ExpiresAt time.Time
}

func (UserPasswordResetRequested) Key() string { return UserPasswordResetRequestedEvent }

type UserPasswordReset struct {
	NewPassword string
}

func (UserPasswordReset) Key() string { return UserPasswordResetEvent }

// UserTokenInvalidated event is triggered when a user's tokens are explicitly invalidated
type UserTokenInvalidated struct {
	UserID        string
	TokenID       string
	InvalidatedAt time.Time
	Reason        string
}

func (UserTokenInvalidated) Key() string { return UserTokenInvalidatedEvent }

// UserTokenRefreshed event is triggered when a user's tokens are refreshed
type UserTokenRefreshed struct {
	UserID      string
	OldTokenID  string
	NewTokenID  string
	RefreshedAt time.Time
}

func (UserTokenRefreshed) Key() string { return UserTokenRefreshedEvent }

// UserTOTPSetupInitiated event is triggered when a user starts TOTP setup
type UserTOTPSetupInitiated struct {
	UserID       string
	InitiatedAt  time.Time
}

func (UserTOTPSetupInitiated) Key() string { return UserTOTPSetupInitiatedEvent }

// UserTOTPEnabled event is triggered when a user enables TOTP
type UserTOTPEnabled struct {
	UserID           string
	EnabledAt        time.Time
	BackupCodesCount int
}

func (UserTOTPEnabled) Key() string { return UserTOTPEnabledEvent }

// UserTOTPDisabled event is triggered when a user disables TOTP
type UserTOTPDisabled struct {
	UserID     string
	DisabledAt time.Time
	Reason     string
}

func (UserTOTPDisabled) Key() string { return UserTOTPDisabledEvent }

// UserTOTPVerified event is triggered when a user successfully verifies with TOTP
type UserTOTPVerified struct {
	UserID     string
	VerifiedAt time.Time
	IP         string
	UserAgent  string
}

func (UserTOTPVerified) Key() string { return UserTOTPVerifiedEvent }

// UserBackupCodeUsed event is triggered when a user uses a backup code
type UserBackupCodeUsed struct {
	UserID         string
	UsedAt         time.Time
	RemainingCodes int
}

func (UserBackupCodeUsed) Key() string { return UserBackupCodeUsedEvent }

// UserSSOLinked event is triggered when a user links an SSO provider
type UserSSOLinked struct {
	UserID       string
	Provider     string
	SubjectID    string
	LinkedAt     time.Time
}

func (UserSSOLinked) Key() string { return UserSSOLinkedEvent }
