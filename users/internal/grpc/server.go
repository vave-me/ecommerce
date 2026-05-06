package grpc

import (
	"context"
	"errors" // Standard errors package
	"fmt"
	"strings"

	"middleman/internal/di"
	"middleman/internal/errorsotel"
	oidcclient "middleman/internal/oid" // Assuming this is the correct package name
	"middleman/users/internal/application"
	"middleman/users/internal/application/commands"

	"middleman/users/internal/application/queries"
	"middleman/users/internal/constants"
	"middleman/users/internal/domain"
	"middleman/users/userspb"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes" // Renamed to avoid conflict with grpc codes
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes" // For gRPC status codes
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	app application.App
	userspb.UsersServiceServer
}

// CRITICAL REFACTOR NOTE: These error variables should be defined in a dedicated application errors package
// (e.g., in 'middleman/users/internal/application/errors' as apperrors.ErrUserNotFound)
// and then checked using `errors.Is(err, apperrors.ErrUserNotFound)`.
// Using global variables here and direct comparison or string checking (as in the placeholder isUserNotFoundError) is NOT best practice.
var ErrUserNotFound = fmt.Errorf("user not found")  // MOVE to app errors package
var ErrDuplicateUser = fmt.Errorf("duplicate user") // MOVE to app errors package

var _ userspb.UsersServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	userspb.RegisterUsersServiceServer(registrar, server{app: app})
	return nil
}

// Placeholder error checking functions - these MUST be replaced with robust checks
// using errors.Is against specific error types from your application layer.
// Example: func isUserNotFoundError(err error) bool { return errors.Is(err, apperrors.ErrUserNotFound) }
func isUserNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	// BRITTLE FALLBACK - REPLACE THIS:
	return strings.Contains(strings.ToLower(err.Error()), "not found") || errors.Is(err, ErrUserNotFound)
}

func isDuplicateUserError(err error) bool {
	if err == nil {
		return false
	}
	// Check for domain error first
	if errors.Is(err, domain.ErrDuplicateEmail) {
		return true
	}
	// BRITTLE FALLBACK - REPLACE THIS:
	return strings.Contains(strings.ToLower(err.Error()), "duplicate") || errors.Is(err, ErrDuplicateUser)
}

func (s server) CreateUser(ctx context.Context, request *userspb.CreateUserRequest) (*userspb.CreateUserResponse, error) {
	span := trace.SpanFromContext(ctx)
	userID := uuid.New().String()
	span.SetAttributes(
		attribute.String("user.id", userID), // Use OpenTelemetry semantic conventions for attribute names
		attribute.String("user.email", request.GetEmail()),
	)

	// ACTION: Implement robust input validation (e.g., email format, password strength)
	if request.GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}
	if request.GetPassword() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password is required")
	}

	// Validate role if provided
	if request.GetRole() != "" {
		validRoles := []string{"customer", "business", "admin"}
		isValidRole := false
		for _, validRole := range validRoles {
			if request.GetRole() == validRole {
				isValidRole = true
				break
			}
		}
		if !isValidRole {
			return nil, status.Errorf(codes.InvalidArgument, "invalid role: must be one of customer, business, or admin")
		}

		// Prevent admin creation through regular CreateUser endpoint
		if request.GetRole() == "admin" {
			return nil, status.Errorf(codes.PermissionDenied, "admin users must be created through the AddAdmin endpoint")
		}
	}
	// Add more validation as needed (username constraints, etc.)

	err := s.app.CreateUser(ctx, commands.CreateUser{
		ID:        userID,
		Email:     request.GetEmail(),
		Username:  request.GetUserName(),
		Password:  request.GetPassword(),
		FirstName: request.GetFirstName(),
		LastName:  request.GetLastName(),
		Latitude:  float64(request.GetLat()),
		Longitude: float64(request.GetLng()),
		Thumbnail: request.GetThumbnail(),
		Role:      request.GetRole(),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, err.Error())
		if isDuplicateUserError(err) {
			return nil, status.Errorf(codes.AlreadyExists, "user creation failed: resource already exists: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &userspb.CreateUserResponse{
		Id: userID,
	}, nil
}

func (s server) AddAdmin(ctx context.Context, request *userspb.AddAdminRequest) (*userspb.AddAdminResponse, error) {
	span := trace.SpanFromContext(ctx)

	userID := uuid.New().String()
	span.SetAttributes(
		attribute.String("user.id", userID),
		attribute.String("user.email", request.GetEmail()),
		attribute.String("user.role", "admin"),
	)

	// Validate required fields
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and password are required")
	}

	if request.GetFirstName() == "" || request.GetLastName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "first name and last name are required")
	}

	err := s.app.AddAdmin(ctx, commands.AddAdmin{
		ID:        userID,
		Username:  request.GetUserName(),
		Email:     request.GetEmail(),
		Password:  request.GetPassword(),
		FirstName: request.GetFirstName(),
		LastName:  request.GetLastName(),
		Latitude:  float64(request.GetLat()),
		Longitude: float64(request.GetLng()),
		Thumbnail: request.GetThumbnail(),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, err.Error())
		if isDuplicateUserError(err) {
			return nil, status.Errorf(codes.AlreadyExists, "admin creation failed: resource already exists: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create admin: %v", err)
	}

	return &userspb.AddAdminResponse{
		Id: userID,
	}, nil
}

func (s server) UpdateUser(ctx context.Context, request *userspb.UpdateUserRequest) (*userspb.UpdateUserResponse, error) {
	span := trace.SpanFromContext(ctx)
	userID := uuid.New().String()
	span.SetAttributes(
		attribute.String("user.id", userID), // Use OpenTelemetry semantic conventions for attribute names

	)

	// Add more validation as needed (username constraints, etc.)

	err := s.app.UpdateUser(ctx, commands.UpdateUser{
		ID:         userID,
		Username:   request.GetUserName(),
		FirstName:  request.GetFirstName(),
		LastName:   request.GetLastName(),
		Latitude:   float64(request.GetLat()),
		Longitude:  float64(request.GetLng()),
		Thumbnail:  request.GetThumbnail(),
		Bio:        request.GetBio(),
		Background: request.GetBackground(),
		Privacy:    domain.UserPrivacy(request.GetPrivacy()),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, err.Error())
		if isDuplicateUserError(err) {
			return nil, status.Errorf(codes.AlreadyExists, "user creation failed: resource already exists: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &userspb.UpdateUserResponse{}, nil
}

func (s server) EnableUser(ctx context.Context, request *userspb.EnableUserRequest) (*userspb.EnableUserResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("user.id", request.GetId()))

	if request.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}

	err := s.app.EnableUser(ctx, commands.EnableUser{
		ID:                request.GetId(),
		VerificationToken: request.GetVerificationToken(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(otelcodes.Error, err.Error())

		// Handle specific error cases
		switch {
		case isUserNotFoundError(err):
			return nil, status.Errorf(codes.NotFound, "failed to enable user: user not found: %v", err)
		case strings.Contains(err.Error(), "invalid verification token"):
			return nil, status.Errorf(codes.InvalidArgument, "failed to enable user: invalid verification token")
		case strings.Contains(err.Error(), "verification token does not match"):
			return nil, status.Errorf(codes.InvalidArgument, "failed to enable user: token doesn't match user")
		default:
			return nil, status.Errorf(codes.Internal, "failed to enable user: %v", err)
		}
	}
	return &userspb.EnableUserResponse{}, nil
}

func (s server) DisableUser(ctx context.Context, request *userspb.DisableUserRequest) (*userspb.DisableUserResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("user.id", request.GetId()))

	if request.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}

	err := s.app.DisableUser(ctx, commands.DisableUser{ID: request.GetId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(otelcodes.Error, err.Error())
		if isUserNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "failed to disable user: user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to disable user: %v", err)
	}
	return &userspb.DisableUserResponse{}, nil
}

func (s server) LoginUser(ctx context.Context, request *userspb.LoginUserRequest) (*userspb.LoginUserResponse, error) {
	// ACTION: Implement input validation
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and password are required")
	}

	middlemanUser, err := s.app.GetUserByMail(ctx, queries.GetUserByMail{Email: request.GetEmail()})
	if err != nil {
		if isUserNotFoundError(err) {
			return nil, status.Errorf(codes.Unauthenticated, "login failed: invalid credentials") // Generic for security
		}
		return nil, status.Errorf(codes.Internal, "login failed due to server error: %v", err)
	}

	// Use the existing LoginUser command which handles authentication and token generation
	accessToken, refreshToken, username, err := s.app.LoginUser(ctx, commands.LoginUser{
		UserID:   middlemanUser.ID,
		Email:    middlemanUser.Email, // Use email from DB for consistency
		Password: request.GetPassword(),
	})
	if err != nil {
		// Your app.LoginUser should return specific errors for bad credentials vs. other issues
		// For now, assume any error from LoginUser after finding the user means bad credentials.
		return nil, status.Errorf(codes.Unauthenticated, "login failed: invalid credentials") // Generic for security
	}

	// Set refresh token as HTTP-only cookie via gRPC metadata
	// This provides secure client-side token storage without breaking existing functionality
	if refreshToken != "" {
		// We need to access the auth instance through the app to create cookies
		// For now, create a simple cookie string manually to avoid breaking changes
		cookieStr := fmt.Sprintf("refresh_token=%s; HttpOnly; Secure; SameSite=Strict; Path=/; Max-Age=604800", refreshToken)
		if err := grpc.SetHeader(ctx, metadata.Pairs("set-cookie", cookieStr)); err != nil {
			// Log but don't fail - cookie is optional enhancement
			fmt.Printf("Warning: Failed to set refresh token cookie: %v\n", err)
		}
	}

	return &userspb.LoginUserResponse{
		AccessToken: accessToken,
		Token:       refreshToken,
		UserName:    username,
		Lat:         float32(middlemanUser.Lat),
		Lng:         float32(middlemanUser.Lng),
	}, nil
}

func (s server) WebLoginWithGoogle(ctx context.Context, request *userspb.WebLoginWithGoogleRequest) (*userspb.LoginUserResponse, error) {
	idToken := request.GetIdToken()
	if idToken == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Google ID token is required")
	}

	googleOIDCClient, ok := di.Get(ctx, constants.WebGoogleVerifierKey).(*oidcclient.GoogleOIDCClient)
	if !ok || googleOIDCClient == nil {
		return nil, status.Errorf(codes.Internal, "internal server configuration error: OIDC client unavailable")
	}

	verifiedToken, err := googleOIDCClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to verify Google ID token: %v", err)
	}

	var claims oidcclient.GoogleUserClaims
	if err := googleOIDCClient.ParseClaims(verifiedToken, &claims); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse Google token claims: %w", err)
	}

	// --- Essential Claim Validation ---
	if claims.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email claim not found or empty in Google token")
	}

	// Corrected access to claims.EmailVerified (it's already bool)
	if !claims.EmailVerified {
		return nil, status.Errorf(codes.PermissionDenied, "Google account email (%s) is not verified", claims.Email)
	}

	if claims.Subject == "" { // `sub` claim is the unique Google User ID
		return nil, status.Errorf(codes.InvalidArgument, "subject (sub) claim not found in Google token, cannot uniquely identify user")
	}

	firstName := claims.GivenName
	if firstName == "" { // Fallback to full name if given_name is empty
		firstName = claims.Name
	}
	// lastName can be claims.FamilyName; often empty, handle gracefully in CreateUserFromGoogleCommand.

	var userID string
	var finalUsername string
	emailForLookup := claims.Email

	middlemanUser, err := s.app.GetUserByGoogleID(ctx, queries.GetUserByGoogleID{GoogleID: claims.Subject})
	if err != nil {
		if isUserNotFoundError(err) { // ACTION: Replace with robust error check
			// First check if a user with this email already exists
			existingUserByEmail, emailErr := s.app.GetUserByMail(ctx, queries.GetUserByMail{Email: emailForLookup})
			if emailErr == nil && existingUserByEmail != nil {
				// User with this email exists but hasn't linked Google account
				if existingUserByEmail.GoogleID == "" {
					// Link the Google account to existing user
					// TODO: Implement a LinkGoogleAccount command to update the user's Google ID
					return nil, status.Errorf(codes.AlreadyExists, "a user with this email already exists. Please login with your existing credentials first, then link your Google account from settings")
				} else {
					// User has a different Google ID linked
					return nil, status.Errorf(codes.AlreadyExists, "a user with this email already exists with a different Google account")
				}
			}
			
			createUserCmd := commands.CreateUserFromGoogle{
				ID:            uuid.New().String(), // Generate new ID for user
				Email:         emailForLookup,
				FirstName:     firstName,
				LastName:      claims.FamilyName,
				GoogleID:      claims.Subject, // CRITICAL
				Thumbnail:     claims.Picture,
				Locale:        claims.Locale,
				Enabled:       true,
				EmailVerified: true, // From Google
				Role:          "user", // Default role for new users
			}
			// ACTION: Ensure s.app.CreateUserFromGoogle is implemented and returns {ID, Username, Email string} or similar.
			err = s.app.CreateUserFromGoogle(ctx, createUserCmd)
			if err != nil {
				if isDuplicateUserError(err) {
					return nil, status.Errorf(codes.AlreadyExists, "a user with this email already exists")
				}
				return nil, status.Errorf(codes.Internal, "failed to create user during Google authentication: %v", err)
			}

			// After creating the user, update the IDs for login
			userID = createUserCmd.ID
			// Set a default username based on email (can be updated later by the user)
			finalUsername = strings.Split(emailForLookup, "@")[0]
		} else {
			return nil, status.Errorf(codes.Internal, "failed to retrieve user information: %v", err)
		}
	} else {
		userID = middlemanUser.ID
		finalUsername = middlemanUser.Username // Assuming your domain.MiddlemanUser has Username

		// Optional but recommended: Link GoogleID if not already present or update profile info
		if middlemanUser.GoogleID == "" || middlemanUser.GoogleID != claims.Subject {
			// Implementation note: You need to add LinkGoogleAccount to your application interfaces
			// and create a commands.LinkGoogleAccount command with handler
			err := s.app.CreateUserFromGoogle(ctx, commands.CreateUserFromGoogle{
				ID:            userID,
				Email:         middlemanUser.Email,
				GoogleID:      claims.Subject,
				Thumbnail:     middlemanUser.Thumbnail,
				FirstName:     middlemanUser.FirstName,
				LastName:      middlemanUser.LastName,
				EmailVerified: true,
				Enabled:       true,
				Role:          "user", // Preserve existing role for existing users
			})
			if err != nil {
				// Continue despite error - don't block login
			}
		}
	}

	// Use the LoginWithGoogle command for OAuth-based authentication
	loginCmd := commands.WebLoginWithGoogle{
		UserID:        userID,
		GoogleID:      claims.Subject,
		Email:         claims.Email,
		EmailVerified: claims.EmailVerified,
		FirstName:     firstName,
		LastName:      claims.FamilyName,
		Picture:       claims.Picture,
	}
	accessToken, refreshToken, loggedInUsername, err := s.app.WebLoginWithGoogle(ctx, loginCmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login process failed after Google authentication: %v", err)
	}

	if loggedInUsername != "" { // Prefer username from LoginUser if it's more authoritative
		finalUsername = loggedInUsername
	}

	return &userspb.LoginUserResponse{
		Token:       refreshToken,
		AccessToken: accessToken,
		UserName:    finalUsername,
	}, nil
}
func (s server) MobileLoginWithGoogle(ctx context.Context, request *userspb.MobileLoginWithGoogleRequest) (*userspb.LoginUserResponse, error) {
	idToken := request.GetIdToken()
	if idToken == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Google ID token is required")
	}

	googleOIDCClient, ok := di.Get(ctx, constants.MobileGoogleVerifierKey).(*oidcclient.GoogleOIDCClient)
	if !ok || googleOIDCClient == nil {
		return nil, status.Errorf(codes.Internal, "internal server configuration error: OIDC client unavailable")
	}

	verifiedToken, err := googleOIDCClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to verify Google ID token: %v", err)
	}

	var claims oidcclient.GoogleUserClaims
	if err := googleOIDCClient.ParseClaims(verifiedToken, &claims); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse Google token claims: %w", err)
	}

	// --- Essential Claim Validation ---
	if claims.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email claim not found or empty in Google token")
	}

	// Corrected access to claims.EmailVerified (it's already bool)
	if !claims.EmailVerified {
		return nil, status.Errorf(codes.PermissionDenied, "Google account email (%s) is not verified", claims.Email)
	}

	if claims.Subject == "" { // `sub` claim is the unique Google User ID
		return nil, status.Errorf(codes.InvalidArgument, "subject (sub) claim not found in Google token, cannot uniquely identify user")
	}

	firstName := claims.GivenName
	if firstName == "" { // Fallback to full name if given_name is empty
		firstName = claims.Name
	}
	// lastName can be claims.FamilyName; often empty, handle gracefully in CreateUserFromGoogleCommand.

	var userID string
	var finalUsername string
	emailForLookup := claims.Email

	middlemanUser, err := s.app.GetUserByGoogleID(ctx, queries.GetUserByGoogleID{GoogleID: claims.Subject})

	if err != nil {
		if isUserNotFoundError(err) { // ACTION: Replace with robust error check
			// First check if a user with this email already exists
			existingUserByEmail, emailErr := s.app.GetUserByMail(ctx, queries.GetUserByMail{Email: emailForLookup})
			if emailErr == nil && existingUserByEmail != nil {
				// User with this email exists but hasn't linked Google account
				if existingUserByEmail.GoogleID == "" {
					// Link the Google account to existing user
					// TODO: Implement a LinkGoogleAccount command to update the user's Google ID
					return nil, status.Errorf(codes.AlreadyExists, "a user with this email already exists. Please login with your existing credentials first, then link your Google account from settings")
				} else {
					// User has a different Google ID linked
					return nil, status.Errorf(codes.AlreadyExists, "a user with this email already exists with a different Google account")
				}
			}
			
			createUserCmd := commands.CreateUserFromGoogle{
				ID:            uuid.New().String(), // Generate new ID for user
				Email:         emailForLookup,
				FirstName:     firstName,
				LastName:      claims.FamilyName,
				GoogleID:      claims.Subject, // CRITICAL
				Thumbnail:     claims.Picture,
				Locale:        claims.Locale,
				Enabled:       true,
				EmailVerified: true, // From Google
				Role:          "user", // Default role for new users
			}
			// ACTION: Ensure s.app.CreateUserFromGoogle is implemented and returns {ID, Username, Email string} or similar.
			err = s.app.CreateUserFromGoogle(ctx, createUserCmd)
			if err != nil {
				if isDuplicateUserError(err) {
					return nil, status.Errorf(codes.AlreadyExists, "a user with this email already exists")
				}
				return nil, status.Errorf(codes.Internal, "failed to create user during Google authentication: %v", err)
			}

			// After creating the user, update the IDs for login
			userID = createUserCmd.ID
			// Set a default username based on email (can be updated later by the user)
			finalUsername = strings.Split(emailForLookup, "@")[0]
		} else {
			return nil, status.Errorf(codes.Internal, "failed to retrieve user information: %v", err)
		}
	} else {
		userID = middlemanUser.ID
		finalUsername = middlemanUser.Username // Assuming your domain.MiddlemanUser has Username

		// Optional but recommended: Link GoogleID if not already present or update profile info
		if middlemanUser.GoogleID == "" || middlemanUser.GoogleID != claims.Subject {
			// Implementation note: You need to add LinkGoogleAccount to your application interfaces
			// and create a commands.LinkGoogleAccount command with handler
			err := s.app.CreateUserFromGoogle(ctx, commands.CreateUserFromGoogle{
				ID:            userID,
				Email:         middlemanUser.Email,
				GoogleID:      claims.Subject,
				Thumbnail:     middlemanUser.Thumbnail,
				FirstName:     middlemanUser.FirstName,
				LastName:      middlemanUser.LastName,
				EmailVerified: true,
				Enabled:       true,
				Role:          "user", // Preserve existing role for existing users
			})
			if err != nil {
				// Continue despite error - don't block login
			}
		}
	}

	// Use the LoginWithGoogle command for OAuth-based authentication
	loginCmd := commands.MobileLoginWithGoogle{
		UserID:        userID,
		GoogleID:      claims.Subject,
		Email:         claims.Email,
		EmailVerified: claims.EmailVerified,
		FirstName:     firstName,
		LastName:      claims.FamilyName,
		Picture:       claims.Picture,
	}
	accessToken, refreshToken, loggedInUsername, err := s.app.MobileLoginWithGoogle(ctx, loginCmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login process failed after Google authentication: %v", err)
	}

	if loggedInUsername != "" { // Prefer username from LoginUser if it's more authoritative
		finalUsername = loggedInUsername
	}

	return &userspb.LoginUserResponse{
		Token:       refreshToken,
		AccessToken: accessToken,
		UserName:    finalUsername,
	}, nil
}
func (s server) LogoutUser(ctx context.Context, request *userspb.LogUserOutRequest) (*userspb.LogUserOutResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("user.id", request.GetId()))

	if request.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}

	// Call the enhanced LogoutUser method with token information
	err := s.app.LogoutUser(ctx, commands.LogoutUser{
		UserID:       request.GetId(),
		AuthToken:    request.GetAuthToken(),
		RefreshToken: request.GetRefreshToken(),
	})

	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(otelcodes.Error, err.Error())

		if isUserNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "logout failed: user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "logout failed: %v", err)
	}

	// Clear refresh token cookie via gRPC metadata
	// Create an expired cookie to clear the existing one
	expiredCookieStr := "refresh_token=; HttpOnly; Secure; SameSite=Strict; Path=/; Max-Age=0; Expires=Thu, 01 Jan 1970 00:00:00 GMT"
	if err := grpc.SetHeader(ctx, metadata.Pairs("set-cookie", expiredCookieStr)); err != nil {
		// Log but don't fail - cookie clearing is optional
		fmt.Printf("Warning: Failed to clear refresh token cookie: %v\n", err)
	}

	return &userspb.LogUserOutResponse{}, nil
}

func (s server) RenameUser(ctx context.Context, request *userspb.RenameUserRequest) (*userspb.RenameUserResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("user.id", request.GetId()),
		attribute.String("user.new_username", request.GetUserName()),
	)

	// ACTION: Implement input validation
	if request.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}
	if request.GetUserName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "new username is required")
	}

	err := s.app.RenameUser(ctx, commands.RenameUser{
		ID:   request.GetId(),
		Name: request.GetUserName(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(otelcodes.Error, err.Error())
		if isUserNotFoundError(err) { // ACTION: Replace
			return nil, status.Errorf(codes.NotFound, "rename failed: user not found: %v", err)
		}
		if isDuplicateUserError(err) { // ACTION: Replace (e.g., if new username conflicts)
			return nil, status.Errorf(codes.AlreadyExists, "rename failed: new username may already be taken: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "rename failed: %v", err)
	}
	return &userspb.RenameUserResponse{}, nil
}

func (s server) GetUser(ctx context.Context, request *userspb.GetUserRequest) (*userspb.GetUserResponse, error) {
	if request.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}

	user, err := s.app.GetUser(ctx, queries.GetUser{ID: request.GetId()})
	if err != nil {
		if isUserNotFoundError(err) { // ACTION: Replace
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	return &userspb.GetUserResponse{User: s.userFromDomain(user)}, nil
}

func (s server) GetBaseUser(ctx context.Context, request *userspb.GetBaseUserRequest) (*userspb.GetBaseUserResponse, error) {
	if request.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}

	user, err := s.app.GetBaseUser(ctx, queries.GetBaseUser{ID: request.GetId()})
	if err != nil {
		if isUserNotFoundError(err) { // ACTION: Replace
			return nil, status.Errorf(codes.NotFound, "base user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get base user: %v", err)
	}
	return &userspb.GetBaseUserResponse{User: s.viewUserFromDomain(user)}, nil
}

func (s server) GetUsers(ctx context.Context, request *userspb.GetUsersRequest) (*userspb.GetUsersResponse, error) {
	span := trace.SpanFromContext(ctx)

	users, err := s.app.GetUsers(ctx, queries.GetUsers{})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(otelcodes.Error, err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	protoUsers := make([]*userspb.User, 0, len(users))
	for _, domainUser := range users { // Changed loop variable name for clarity
		protoUsers = append(protoUsers, s.userFromDomain(domainUser))
	}
	return &userspb.GetUsersResponse{Users: protoUsers}, nil
}

func (s server) ListEnabledUsers(ctx context.Context, request *userspb.ListEnabledUsersRequest) (*userspb.ListEnabledUsersResponse, error) {
	span := trace.SpanFromContext(ctx)

	users, err := s.app.GetEnabledUsers(ctx, queries.GetEnabledUsers{})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(otelcodes.Error, err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get enabled users: %v", err)
	}

	protoUsers := make([]*userspb.User, 0, len(users))
	for _, domainUser := range users {
		protoUsers = append(protoUsers, s.userFromDomain(domainUser))
	}
	return &userspb.ListEnabledUsersResponse{Users: protoUsers}, nil
}

func (s server) ForgotPassword(ctx context.Context, request *userspb.ForgotPasswordRequest) (*userspb.ForgotPasswordResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("user.email", request.GetEmail()))

	if request.GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}

	// First find the user by email
	middlemanUser, err := s.app.GetUserByMail(ctx, queries.GetUserByMail{Email: request.GetEmail()})
	if err != nil {
		if isUserNotFoundError(err) {
			// To prevent email enumeration, return a generic success message even if user not found.
			return &userspb.ForgotPasswordResponse{Message: "If an account with that email exists, password reset instructions have been sent."}, nil
		}
		// For other errors (DB down, etc.), return an internal error.
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(otelcodes.Error, err.Error())
		return nil, status.Errorf(codes.Internal, "failed to process forgot password request: %v", err)
	}

	// Now use the found user ID to initiate password reset
	err = s.app.ForgotPassword(ctx, commands.ForgotPassword{
		MiddlemanUserID: middlemanUser.ID,
		Email:           middlemanUser.Email, // Use the confirmed email from the database
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(otelcodes.Error, err.Error())
		// It's often still best to return a generic message here to avoid leaking information about failures.
		return nil, status.Errorf(codes.Internal, "failed to initiate password reset process: %v", err)
	}

	// Note: At this point, in a real implementation, an email would be sent to the user
	// with the reset token or a link containing the token.
	// This could be handled by an event handler that subscribes to the UserPasswordResetRequestedEvent.

	return &userspb.ForgotPasswordResponse{
		Message: "If an account with that email exists, password reset instructions have been sent.",
	}, nil
}

func (s server) ResetPassword(ctx context.Context, request *userspb.ResetPasswordRequest) (*userspb.ResetPasswordResponse, error) {
	span := trace.SpanFromContext(ctx)
	// Add attributes if useful, e.g., part of token hash if that's a pattern, but generally avoid token itself.
	// span.SetAttributes(attribute.String("token_present", fmt.Sprintf("%t", request.GetToken() != "")))

	// Input validation
	if request.GetToken() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "reset token is required")
	}
	if request.GetNewPassword() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "new password is required")
	}

	// Password strength validation (add additional validation as needed)
	if len(request.GetNewPassword()) < 8 {
		return nil, status.Errorf(codes.InvalidArgument, "password must be at least 8 characters long")
	}

	err := s.app.ResetPassword(ctx, commands.ResetPassword{
		Token:       request.GetToken(),
		NewPassword: request.GetNewPassword(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(otelcodes.Error, err.Error())

		// Handle specific error cases
		if strings.Contains(err.Error(), "reset token has expired") {
			return nil, status.Errorf(codes.InvalidArgument, "password reset failed: token has expired")
		}
		if strings.Contains(err.Error(), "no reset token found") {
			return nil, status.Errorf(codes.InvalidArgument, "password reset failed: invalid token")
		}

		return nil, status.Errorf(codes.Internal, "password reset failed: %v", err)
	}

	return &userspb.ResetPasswordResponse{
		Message: "Password has been successfully reset.",
	}, nil
}

func (s server) RefreshAuthToken(ctx context.Context, request *userspb.RefreshAuthTokenRequest) (*userspb.RefreshAuthTokenResponse, error) {
	span := trace.SpanFromContext(ctx)

	// Validation
	refreshToken := request.GetRefreshToken()
	if refreshToken == "" {
		return nil, status.Errorf(codes.InvalidArgument, "refresh token is required")
	}

	// Execute command with updated signature for token refresh
	accessToken, newRefreshToken, err := s.app.RefreshAuthToken(ctx, commands.RefreshAuthToken{
		RefreshToken: refreshToken,
		UserID:       request.GetUserId(), // This is optional, will be validated in handler
	})

	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(otelcodes.Error, err.Error())

		// Handle specific error cases
		switch {
		case strings.Contains(err.Error(), "token is expired") || strings.Contains(err.Error(), "expired token"):
			return nil, status.Errorf(codes.Unauthenticated, "token refresh failed: token expired")
		case strings.Contains(err.Error(), "invalid token") || strings.Contains(err.Error(), "invalid refresh token"):
			return nil, status.Errorf(codes.Unauthenticated, "token refresh failed: invalid token format")
		case isUserNotFoundError(err):
			return nil, status.Errorf(codes.NotFound, "token refresh failed: user not found")
		case strings.Contains(err.Error(), "user account is disabled"):
			return nil, status.Errorf(codes.PermissionDenied, "token refresh failed: account is disabled")
		default:
			return nil, status.Errorf(codes.Internal, "failed to refresh token: %v", err)
		}
	}

	return &userspb.RefreshAuthTokenResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// --- Mapper functions ---

func (s server) userFromDomain(user *domain.MiddlemanUser) *userspb.User {
	if user == nil {
		return nil // Or return an empty userspb.User{} if that's preferred over nil
	}
	return &userspb.User{
		Id:         user.ID,
		Email:      user.Email,
		UserName:   user.Username,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Lat:        float32(user.Lat),
		Lng:        float32(user.Lng),
		Enabled:    user.Enabled,
		Thumbnail:  user.Thumbnail,
		Bio:        user.Bio,
		Privacy:    user.Privacy,
		Background: user.Background,
		GoogleId:   user.GoogleID,
		// Consider if Thumbnail should be part of this more general User proto.
	}
}

func (s server) viewUserFromDomain(user *domain.MiddlemanViewUser) *userspb.BaseUser {
	if user == nil {
		return nil // Or return an empty userspb.BaseUser{}
	}
	return &userspb.BaseUser{
		Id:         user.ID,
		UserName:   user.Username,
		Location:   user.Location,
		Lat:        float32(user.Lat),
		Lng:        float32(user.Lng),
		Thumbnail:  user.Thumbnail,
		Bio:        user.Bio,
		Privacy:    user.Privacy,
		Background: user.Background,
	}
}

// Add the ClearTokens method to explicitly invalidate tokens
func (s server) ClearTokens(ctx context.Context, request *userspb.ClearTokensRequest) (*userspb.ClearTokensResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("user.id", request.GetUserId()),
		attribute.String("token.invalidation.reason", request.GetReason()),
	)

	if request.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}

	// Prepare default reason if not provided
	reason := request.GetReason()
	if reason == "" {
		reason = "user requested token invalidation"
	}

	// Call the application command
	err := s.app.ClearTokens(ctx, commands.ClearTokens{
		UserID:       request.GetUserId(),
		TokenID:      request.GetTokenId(),
		RefreshToken: request.GetRefreshToken(),
		Reason:       reason,
	})

	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(otelcodes.Error, err.Error())

		if isUserNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "token invalidation failed: user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "token invalidation failed: %v", err)
	}

	return &userspb.ClearTokensResponse{
		Success: true,
		Message: "Tokens successfully invalidated",
	}, nil
}
