package grpc

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
	"middleman/internal/auth"
	"middleman/internal/rpc"
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

// GetUserByID retrieves a user by ID
func (r UserRepository) GetUserByID(ctx context.Context, userID string) (user *models.User, err error) {
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

// GetBaseUserByID retrieves base user information by ID
func (r UserRepository) GetBaseUserByID(ctx context.Context, userID string) (*models.BaseUser, error) {
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

// GetMultipleUsersByIDs retrieves multiple users by their IDs
func (r UserRepository) GetMultipleUsersByIDs(ctx context.Context, userIDs []string) ([]*models.User, error) {
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

// GetAllParticipatingUsers retrieves all participating users
func (r UserRepository) GetAllParticipatingUsers(ctx context.Context) ([]*models.User, error) {
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

// CreateNewUser creates a new user
func (r UserRepository) CreateNewUser(ctx context.Context, email, password, username, firstName, lastName, location string, lat, lng float32, thumbnail, language string) (string, error) {
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

// UpdateUserProfile updates user information
func (r UserRepository) UpdateUserProfile(ctx context.Context, id, username, firstName, lastName, bio, privacy, background, location string, lat, lng float32, thumbnail string) (string, error) {
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

// ChangeUsername renames a user
func (r UserRepository) ChangeUsername(ctx context.Context, id, username string) (*models.User, error) {
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

// ActivateUserAccount enables a user account
func (r UserRepository) ActivateUserAccount(ctx context.Context, id, verificationToken string) error {
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

// DeactivateUserAccount disables a user account
func (r UserRepository) DeactivateUserAccount(ctx context.Context, id string) error {
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

// AuthenticateUser authenticates a user with email and password
func (r UserRepository) AuthenticateUser(ctx context.Context, email, password string) (*models.LoginResponse, error) {
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

// AuthenticateWithGoogleWeb authenticates a user with Google (web)
func (r UserRepository) AuthenticateWithGoogleWeb(ctx context.Context, idToken string) (*models.LoginResponse, error) {
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

// AuthenticateWithGoogleMobile authenticates a user with Google (mobile)
func (r UserRepository) AuthenticateWithGoogleMobile(ctx context.Context, idToken string) (*models.LoginResponse, error) {
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

// LogoutUser logs out a user
func (r UserRepository) LogoutUser(ctx context.Context, id, authToken, refreshToken string) error {
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

// RefreshUserAuthToken refreshes authentication tokens
func (r UserRepository) RefreshUserAuthToken(ctx context.Context, refreshToken, userID string) (*models.TokenResponse, error) {
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

// RevokeUserTokens clears user tokens
func (r UserRepository) RevokeUserTokens(ctx context.Context, userID, tokenID, refreshToken, reason string) (*models.ClearTokensResponse, error) {
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

// SendPasswordResetEmail initiates password reset
func (r UserRepository) SendPasswordResetEmail(ctx context.Context, email string) (*models.MessageResponse, error) {
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

// ResetUserPassword resets user password
func (r UserRepository) ResetUserPassword(ctx context.Context, token, newPassword string) (*models.MessageResponse, error) {
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
