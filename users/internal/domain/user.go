package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"

	"crypto/rand"
	"encoding/base64"

	"github.com/nats-io/jwt/v2"
	"github.com/stackus/errors"
	"golang.org/x/crypto/bcrypt"
)

const UserAggregate = "users.User"

var (
	ErrUserIDCannotBeBlank        = errors.Wrap(errors.ErrBadRequest, "the user id cannot be blank")
	ErrEmailCannotBeBlank         = errors.Wrap(errors.ErrBadRequest, "the email cannot be blank")
	ErrPasswordCannotBeBlank      = errors.Wrap(errors.ErrBadRequest, "the password cannot be blank")
	ErrUserAlreadyEnabled         = errors.Wrap(errors.ErrBadRequest, "the user is already enabled")
	ErrUserAlreadyDisabled        = errors.Wrap(errors.ErrBadRequest, "the user is already disabled")
	ErrUserFirstNameCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the user first name cannot be blank")
	ErrUserLastNameCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "the user last name cannot be blank")
	ErrDuplicateEmail             = errors.Wrap(errors.ErrBadRequest, "a user with this email already exists")
)

type User struct {
	es.Aggregate
	Email                string
	Password             string
	Username             string
	Bio                  string
	Privacy              string
	FirstName            string
	LastName             string
	Location             string
	Enabled              bool
	KYCVerified          bool
	Thumbnail            string
	Background           string
	Lat                  float64
	Lang                 float64
	ResetToken           string
	TOTPEnabled          bool
	GoogleID             string
	ResetTokenExp        time.Time
	VerificationToken    string
	VerificationTokenExp time.Time
	Language             string
	Role                 UserRole
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*User)(nil)

func NewUser(id string) *User {
	return &User{
		Aggregate: es.NewAggregate(id, UserAggregate),
	}
}

func CreateUser(id, email, password, username, firstName, lastName string, lat, lng float64, thumbnail, token string, role UserRole) (ddd.Event, error) {
	if id == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	if email == "" {
		return nil, ErrEmailCannotBeBlank
	}

	if password == "" {
		return nil, ErrPasswordCannotBeBlank
	}

	if firstName == "" {
		return nil, ErrUserFirstNameCannotBeBlank
	}

	if lastName == "" {
		return nil, ErrUserLastNameCannotBeBlank
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to hash password")
	}

	user := NewUser(id)

	user.AddEvent(UserCreatedEvent, &UserCreated{
		Email:             email,
		Password:          string(hashedPassword),
		Username:          username,
		Firstname:         firstName,
		Lastname:          lastName,
		Enabled:           false,
		Lat:               lat,
		Lng:               lng,
		Thumbnail:         thumbnail,
		VerificationToken: token,
		Role:              role,
	})
	return ddd.NewEvent(UserCreatedEvent, user), nil
}

func (s *User) UpdateUser(username, bio, privacy, firstName, lastName, background string, latitude, longitude float64, thumbnail string) (ddd.Event, error) {

	if bio == "" {
		bio = s.Bio
	}
	if privacy == "" {
		privacy = s.Privacy
	}
	if firstName == "" {
		firstName = s.FirstName
	}

	if lastName == "" {
		lastName = s.LastName
	}
	if latitude == 0 {
		latitude = s.Lat
	}
	if longitude == 0 {
		longitude = s.Lang
	}
	if thumbnail == "" {
		thumbnail = s.Thumbnail
	}
	if background == "" {
		background = s.Background
	}
	if privacy == "" {
		privacy = s.Privacy
	}
	s.AddEvent(UserUpdatedEvent, &UserUpdated{
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		Thumbnail:  thumbnail,
		Privacy:    privacy,
		Bio:        bio,
		Background: background,
		Latitude:   latitude,
		Longitude:  longitude,
	})
	return ddd.NewEvent(UserUpdatedEvent, s), nil
}

func (s *User) InitUser(email, password, username, firstname, lastname string, lat, lng float64, thumbnail, token string, role UserRole) (ddd.Event, error) {
	if email == "" {
		return nil, ErrEmailCannotBeBlank
	}
	if password == "" {
		return nil, ErrPasswordCannotBeBlank
	}
	if firstname == "" {
		return nil, ErrUserFirstNameCannotBeBlank
	}
	if lastname == "" {
		return nil, ErrUserLastNameCannotBeBlank
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to hash password")
	}

	s.AddEvent(UserCreatedEvent, &UserCreated{
		Email:             email,
		Password:          string(hashedPassword),
		Username:          username,
		Firstname:         firstname,
		Lastname:          lastname,
		Enabled:           false,
		Lat:               lat,
		Lng:               lng,
		Thumbnail:         thumbnail,
		VerificationToken: token,
		Role:              role,
	})

	return ddd.NewEvent(UserCreatedEvent, s), nil
}

func (s *User) InitUserWithGoogle(email, firstname, lastname string, enabled bool, googleId, thumbnail string, lat, lng float64, role UserRole) (ddd.Event, error) {
	if email == "" {
		return nil, ErrEmailCannotBeBlank
	}
	if firstname == "" {
		return nil, ErrUserFirstNameCannotBeBlank
	}
	if lastname == "" {
		return nil, ErrUserLastNameCannotBeBlank
	}

	// Default to user role if not specified
	if role == "" {
		role = UserRoleUser
	}

	s.AddEvent(UserCreatedEvent, &UserCreated{
		Email:     email,
		Firstname: firstname,
		Lastname:  lastname,
		GoogleID:  googleId,
		Enabled:   enabled,
		Thumbnail: thumbnail,
		Lat:       lat, // Default to 0 for location if not provided
		Lng:       lng, // Default to 0 for location if not provided
		Role:      role,
	})

	return ddd.NewEvent(UserCreatedEvent, s), nil
}

func (s *User) InitAdminUser(email, password, username, firstname, lastname string, lat, lng float64, thumbnail, token string) (ddd.Event, error) {
	if email == "" {
		return nil, ErrEmailCannotBeBlank
	}
	if password == "" {
		return nil, ErrPasswordCannotBeBlank
	}
	if firstname == "" {
		return nil, ErrUserFirstNameCannotBeBlank
	}
	if lastname == "" {
		return nil, ErrUserLastNameCannotBeBlank
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to hash password")
	}

	s.AddEvent(UserCreatedEvent, &UserCreated{
		Email:             email,
		Password:          string(hashedPassword),
		Username:          username,
		Firstname:         firstname,
		Lastname:          lastname,
		Enabled:           true, // Admins are enabled by default
		Lat:               lat,
		Lng:               lng,
		Thumbnail:         thumbnail,
		VerificationToken: token,
		Role:              UserRoleAdmin, // Set role as admin
	})

	return ddd.NewEvent(UserCreatedEvent, s), nil
}

// Key implements registry.Registerable
func (User) Key() string { return UserAggregate }

func (s *User) Enable() (ddd.Event, error) {
	if s.Enabled {
		return nil, ErrUserAlreadyEnabled
	}

	s.AddEvent(UserEnabledEvent, &UserEnabledToggled{
		Enabled: true,
	})

	return ddd.NewEvent(UserEnabledEvent, s), nil
}

func (s *User) Disable() (ddd.Event, error) {
	if !s.Enabled {
		return nil, ErrUserAlreadyDisabled
	}

	s.AddEvent(UserDisabledEvent, &UserEnabledToggled{
		Enabled: false,
	})

	return ddd.NewEvent(UserDisabledEvent, s), nil
}

func (s *User) Rename(name string) (ddd.Event, error) {
	s.AddEvent(UserRenamedEvent, &UserRenamed{
		FirstName: name,
	})

	return ddd.NewEvent(UserRenamedEvent, s), nil
}

func (s *User) Login() (ddd.Event, error) {

	s.AddEvent(UserLoggedInEvent, &UserLoggedIn{
		UserID: s.ID(),
	})
	return ddd.NewEvent(UserLoggedInEvent, s), nil
}

func (s *User) Authorize(userID, token string) (ddd.Event, error) {
	// TODO decode token claims and implement this in the right way
	_, err := jwt.Decode(token)
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnauthorized, "invalid token")
	}

	return ddd.NewEvent(UserAuthorizedEvent, s), nil
}

func (s *User) KYCVerify(userID string) (ddd.Event, error) {

	return ddd.NewEvent(UserAuthorizedEvent, s), nil
}

// New method: Logs the user out, generating the UserLoggedOutEvent
func (s *User) Logout() (ddd.Event, error) {

	s.AddEvent(UserLoggedOutEvent, &UserLoggedOut{})
	return ddd.NewEvent(UserLoggedOutEvent, s), nil
}

func (s *User) RequestPasswordReset(token, email string, expirationTime time.Time) (ddd.Event, error) {
	s.ResetToken = token
	s.ResetTokenExp = expirationTime

	s.AddEvent(UserPasswordResetRequestedEvent, &UserPasswordResetRequested{
		Token:     token,
		Email:     email,
		ExpiresAt: expirationTime,
	})
	return ddd.NewEvent(UserPasswordResetRequestedEvent, s), nil
}

func (s *User) ResetPassword(newPassword string) (ddd.Event, error) {
	// Validate token expiration
	if s.ResetToken == "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "no reset token found")
	}

	if time.Now().After(s.ResetTokenExp) {
		return nil, errors.Wrap(errors.ErrBadRequest, "reset token has expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to hash password")
	}

	// Clear the reset token after use
	s.ResetToken = ""
	s.ResetTokenExp = time.Time{}

	s.AddEvent(UserPasswordResetEvent, &UserPasswordReset{
		NewPassword: string(hashedPassword),
	})
	return ddd.NewEvent(UserPasswordResetEvent, s), nil
}

// InvalidateTokens explicitly invalidates a user's tokens for security reasons
func (s *User) InvalidateTokens(tokenID, reason string) (ddd.Event, error) {
	s.AddEvent(UserTokenInvalidatedEvent, &UserTokenInvalidated{
		UserID:        s.ID(),
		TokenID:       tokenID,
		InvalidatedAt: time.Now().UTC(),
		Reason:        reason,
	})
	return ddd.NewEvent(UserTokenInvalidatedEvent, s), nil
}

// TokenRefreshed records a token refresh event
func (s *User) TokenRefreshed(oldTokenID, newTokenID string) (ddd.Event, error) {
	s.AddEvent(UserTokenRefreshedEvent, &UserTokenRefreshed{
		UserID:      s.ID(),
		OldTokenID:  oldTokenID,
		NewTokenID:  newTokenID,
		RefreshedAt: time.Now().UTC(),
	})
	return ddd.NewEvent(UserTokenRefreshedEvent, s), nil
}

// ApplyEvent implements es.EventApplier
func (s *User) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *UserCreated:
		s.Email = payload.Email
		s.Password = payload.Password
		s.Username = payload.Username
		s.FirstName = payload.Firstname
		s.LastName = payload.Lastname
		s.Enabled = payload.Enabled
		s.Lat = payload.Lat
		s.Lang = payload.Lng
		s.GoogleID = payload.GoogleID
		s.Thumbnail = payload.Thumbnail
		s.VerificationToken = payload.VerificationToken
		s.Role = payload.Role

	case *UserEnabledToggled:
		s.Enabled = payload.Enabled
	case *UserUpdated:
		s.Username = payload.Username
		s.FirstName = payload.FirstName
		s.LastName = payload.LastName
		s.Lat = payload.Latitude
		s.Lang = payload.Longitude
		s.Thumbnail = payload.Thumbnail
		s.Bio = payload.Bio
		s.Background = payload.Background
		s.Privacy = payload.Privacy
		if payload.Role != "" {
			s.Role = payload.Role
		}
	case *UserRenamed:
		s.FirstName = payload.FirstName

	case *UserLoggedIn:
		// No state changes needed for login event

	case *UserLoggedOut:
		// No state changes needed for logout event

	case *UserTokenInvalidated:
		// No state changes needed for token invalidation event

	case *UserTokenRefreshed:
		// No state changes needed for token refresh event

	case *UserPasswordResetRequested:
		s.ResetToken = payload.Token
		s.ResetTokenExp = payload.ExpiresAt

	case *UserPasswordReset:
		s.Password = payload.NewPassword
		// Clear the reset token
		s.ResetToken = ""
		s.ResetTokenExp = time.Time{}

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", s, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (s *User) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *UserV1:
		s.Email = ss.Email
		s.Password = ss.Password
		s.Username = ss.Username
		s.FirstName = ss.FirstName
		s.LastName = ss.LastName
		s.Enabled = ss.Enabled
		s.Lat = ss.Latitude
		s.Lang = ss.Longitude
		s.GoogleID = ss.GoogleID
		s.Thumbnail = ss.Thumbnail
		s.ResetToken = ss.ResetToken
		s.ResetTokenExp = ss.ResetTokenExp
		s.VerificationToken = ss.VerificationToken
		s.VerificationTokenExp = ss.VerificationTokenExp
		s.Language = ss.Language
		s.Role = ss.Role
		// Default role for existing users without a role
		if s.Role == "" {
			s.Role = UserRoleUser
		}
		s.Bio = ss.Bio
		s.Privacy = ss.Privacy
		s.Background = ss.Background
		s.TOTPEnabled = ss.TOTPEnabled

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", s, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (s User) ToSnapshot() es.Snapshot {
	return UserV1{
		Email:                s.Email,
		Password:             s.Password,
		Username:             s.Username,
		FirstName:            s.FirstName,
		LastName:             s.LastName,
		Enabled:              s.Enabled,
		Latitude:             s.Lat,
		Longitude:            s.Lang,
		GoogleID:             s.GoogleID,
		Thumbnail:            s.Thumbnail,
		ResetToken:           s.ResetToken,
		ResetTokenExp:        s.ResetTokenExp,
		VerificationToken:    s.VerificationToken,
		VerificationTokenExp: s.VerificationTokenExp,
		Language:             s.Language,
		Role:                 s.Role,
		Bio:                  s.Bio,
		Privacy:              s.Privacy,
		Background:           s.Background,
		TOTPEnabled:          s.TOTPEnabled,
	}
}

// GenerateResetToken generates a secure random token for password reset
func GenerateResetToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err) // This should never happen
	}
	return base64.URLEncoding.EncodeToString(b)
}
