package grpc

import (
	"context"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
	"middleman/users/userspb"

	"google.golang.org/grpc"
)

type UserRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.UserRepository = (*UserRepository)(nil)

// NewUserRepositoryWithAuth creates a new UserRepository with JWT authentication support
func NewUserRepository(endpoint string, authInstance *auth.Auth) UserRepository {
	return UserRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// Find retrieves a user by ID
func (r UserRepository) Find(ctx context.Context, userID string) (user *models.User, err error) {
	var conn *grpc.ClientConn
	conn, err = r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := userspb.NewUsersServiceClient(conn).GetUser(ctx, &userspb.GetUserRequest{Id: userID})
	if err != nil {
		return nil, err
	}
	return r.userToDomain(resp.User), nil
}

// GetBaseUser retrieves base user information by ID
func (r UserRepository) GetBaseUser(ctx context.Context, userID string) (*models.BaseUser, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := userspb.NewUsersServiceClient(conn).GetBaseUser(ctx, &userspb.GetBaseUserRequest{Id: userID})
	if err != nil {
		return nil, err
	}

	return r.baseUserToDomain(resp.GetUser()), nil
}

// ListUsers retrieves multiple users by their IDs
func (r UserRepository) ListUsers(ctx context.Context, userIDs []string) ([]*models.User, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := userspb.NewUsersServiceClient(conn).GetUsers(ctx, &userspb.GetUsersRequest{UserIds: userIDs})
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(resp.GetUsers()))
	for i, user := range resp.GetUsers() {
		users[i] = r.userToDomain(user)
	}
	return users, nil
}

// ListParticipatingUsers retrieves all participating users
func (r UserRepository) ListParticipatingUsers(ctx context.Context) ([]*models.User, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := userspb.NewUsersServiceClient(conn).ListEnabledUsers(ctx, &userspb.ListEnabledUsersRequest{})
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(resp.GetUsers()))
	for i, user := range resp.GetUsers() {
		users[i] = r.userToDomain(user)
	}
	return users, nil
}

// CreateUser creates a new user
func (r UserRepository) CreateUser(ctx context.Context, email, password, username, firstName, lastName, location string, lat, lng float32, thumbnail, language string) (string, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return "", err
	}

	resp, err := userspb.NewUsersServiceClient(conn).CreateUser(ctx, &userspb.CreateUserRequest{
		Email:     email,
		Password:  password,
		UserName:  username,
		FirstName: firstName,
		LastName:  lastName,
		Location:  location,
		Lat:       lat,
		Lng:       lng,
		Thumbnail: thumbnail,
		Language:  language,
	})
	if err != nil {
		return "", err
	}

	return resp.GetId(), nil
}

// UpdateUser updates user information
func (r UserRepository) UpdateUser(ctx context.Context, id, username, firstName, lastName, bio, privacy, background, location string, lat, lng float32, thumbnail string) (string, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return "", err
	}

	resp, err := userspb.NewUsersServiceClient(conn).UpdateUser(ctx, &userspb.UpdateUserRequest{
		Id:         id,
		UserName:   username,
		FirstName:  firstName,
		LastName:   lastName,
		Bio:        bio,
		Privacy:    privacy,
		Background: background,
		Location:   location,
		Lat:        lat,
		Lng:        lng,
		Thumbnail:  thumbnail,
	})
	if err != nil {
		return "", err
	}

	return resp.GetUserId(), nil
}

// RenameUser renames a user
func (r UserRepository) RenameUser(ctx context.Context, id, username string) (*models.User, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := userspb.NewUsersServiceClient(conn).RenameUser(ctx, &userspb.RenameUserRequest{
		Id:       id,
		UserName: username,
	})
	if err != nil {
		return nil, err
	}

	return r.userToDomain(resp.GetUser()), nil
}

// EnableUser enables a user account
func (r UserRepository) EnableUser(ctx context.Context, id, verificationToken string) error {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return err
	}

	_, err = userspb.NewUsersServiceClient(conn).EnableUser(ctx, &userspb.EnableUserRequest{
		Id:                id,
		VerificationToken: verificationToken,
	})
	return err
}

// DisableUser disables a user account
func (r UserRepository) DisableUser(ctx context.Context, id string) error {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return err
	}

	_, err = userspb.NewUsersServiceClient(conn).DisableUser(ctx, &userspb.DisableUserRequest{
		Id: id,
	})
	return err
}

// LoginUser authenticates a user with email and password
func (r UserRepository) LoginUser(ctx context.Context, email, password string) (*models.LoginResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := userspb.NewUsersServiceClient(conn).LoginUser(ctx, &userspb.LoginUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		AccessToken: resp.GetAccessToken(),
		Token:       resp.GetToken(),
		Username:    resp.GetUserName(),
		Lat:         resp.GetLat(),
		Lng:         resp.GetLng(),
	}, nil
}

// WebLoginWithGoogle authenticates a user with Google (web)
func (r UserRepository) WebLoginWithGoogle(ctx context.Context, idToken string) (*models.LoginResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := userspb.NewUsersServiceClient(conn).WebLoginWithGoogle(ctx, &userspb.WebLoginWithGoogleRequest{
		IdToken: idToken,
	})
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		AccessToken: resp.GetAccessToken(),
		Token:       resp.GetToken(),
		Username:    resp.GetUserName(),
		Lat:         resp.GetLat(),
		Lng:         resp.GetLng(),
	}, nil
}

// MobileLoginWithGoogle authenticates a user with Google (mobile)
func (r UserRepository) MobileLoginWithGoogle(ctx context.Context, idToken string) (*models.LoginResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := userspb.NewUsersServiceClient(conn).MobileLoginWithGoogle(ctx, &userspb.MobileLoginWithGoogleRequest{
		IdToken: idToken,
	})
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		AccessToken: resp.GetAccessToken(),
		Token:       resp.GetToken(),
		Username:    resp.GetUserName(),
		Lat:         resp.GetLat(),
		Lng:         resp.GetLng(),
	}, nil
}

// LogUserOut logs out a user
func (r UserRepository) LogUserOut(ctx context.Context, id, authToken, refreshToken string) error {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return err
	}

	_, err = userspb.NewUsersServiceClient(conn).LogUserOut(ctx, &userspb.LogUserOutRequest{
		Id:           id,
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	})
	return err
}

// RefreshAuthToken refreshes authentication tokens
func (r UserRepository) RefreshAuthToken(ctx context.Context, refreshToken, userID string) (*models.TokenResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := userspb.NewUsersServiceClient(conn).RefreshAuthToken(ctx, &userspb.RefreshAuthTokenRequest{
		RefreshToken: refreshToken,
		UserId:       userID,
	})
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		Token:        resp.GetToken(),
		RefreshToken: resp.GetRefreshToken(),
	}, nil
}

// ClearTokens clears user tokens
func (r UserRepository) ClearTokens(ctx context.Context, userID, tokenID, refreshToken, reason string) (*models.ClearTokensResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := userspb.NewUsersServiceClient(conn).ClearTokens(ctx, &userspb.ClearTokensRequest{
		UserId:       userID,
		TokenId:      tokenID,
		RefreshToken: refreshToken,
		Reason:       reason,
	})
	if err != nil {
		return nil, err
	}

	return &models.ClearTokensResponse{
		Success: resp.GetSuccess(),
		Message: resp.GetMessage(),
	}, nil
}

// ForgotPassword initiates password reset
func (r UserRepository) ForgotPassword(ctx context.Context, email string) (*models.MessageResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := userspb.NewUsersServiceClient(conn).ForgotPassword(ctx, &userspb.ForgotPasswordRequest{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return &models.MessageResponse{
		Message: resp.GetMessage(),
	}, nil
}

// ResetPassword resets user password
func (r UserRepository) ResetPassword(ctx context.Context, token, newPassword string) (*models.MessageResponse, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)

	if err != nil {
		return nil, err
	}

	resp, err := userspb.NewUsersServiceClient(conn).ResetPassword(ctx, &userspb.ResetPasswordRequest{
		Token:       token,
		NewPassword: newPassword,
	})
	if err != nil {
		return nil, err
	}

	return &models.MessageResponse{
		Message: resp.GetMessage(),
	}, nil
}

// Domain conversion methods

func (r UserRepository) userToDomain(user *userspb.User) *models.User {
	return &models.User{
		ID:         user.GetId(),
		Email:      user.GetEmail(),
		Username:   user.GetUserName(),
		FirstName:  user.GetFirstName(),
		LastName:   user.GetLastName(),
		Enabled:    user.GetEnabled(),
		GoogleID:   user.GetGoogleId(),
		Location:   user.GetLocation(),
		Lat:        user.GetLat(),
		Lng:        user.GetLng(),
		Thumbnail:  user.GetThumbnail(),
		Background: user.GetBackground(),
		Bio:        user.GetBio(),
		Privacy:    user.GetPrivacy(),
	}
}

func (r UserRepository) baseUserToDomain(user *userspb.BaseUser) *models.BaseUser {
	return &models.BaseUser{
		ID:         user.GetId(),
		Username:   user.GetUserName(),
		Thumbnail:  user.GetThumbnail(),
		Lat:        user.GetLat(),
		Lng:        user.GetLng(),
		Location:   user.GetLocation(),
		Bio:        user.GetBio(),
		Privacy:    user.GetPrivacy(),
		Background: user.GetBackground(),
	}
}

// dial sets up a gRPC connection with the microservice endpoint
func (r UserRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r UserRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
