package grpc

import (
	"context"
	"google.golang.org/grpc"
	"middleman/internal/rpc"
	"middleman/search/internal/application"
	"middleman/search/internal/models"
	"middleman/users/userspb"
)

type UserRepository struct {
	endpoint string
}

var _ application.UserRepository = (*UserRepository)(nil)

func NewUserRepository(endpoint string) UserRepository {
	return UserRepository{
		endpoint: endpoint}
}
func (r UserRepository) Find(ctx context.Context, userID string) (user *models.User, err error) {
	var conn *grpc.ClientConn
	conn, err = r.dial(ctx)

	if err != nil {
		return nil, err
	}
	defer conn.Close()

	resp, err := userspb.NewUsersServiceClient(conn).GetUser(ctx, &userspb.GetUserRequest{Id: userID})
	if err != nil {
		return nil, err
	}
	return r.userToDomain(resp.User), nil
}

func (r UserRepository) userToDomain(user *userspb.User) *models.User {
	return &models.User{
		ID:        user.GetId(),
		Email:     user.GetEmail(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Location:  user.GetLocation(),
		Enabled:   user.GetEnabled(),
	}
}

func (r UserRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}
