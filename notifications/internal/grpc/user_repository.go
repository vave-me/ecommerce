package grpc

import (
	"context"
	"middleman/internal/rpc"
	"middleman/notifications/internal/domain"
	"middleman/users/userspb"

	"google.golang.org/grpc"
)

type UserRepository struct {
	endpoint string
}

var _ domain.UserRepository = (*UserRepository)(nil)

func NewUserRepository(endpoint string) UserRepository {
	return UserRepository{
		endpoint: endpoint,
	}
}

func (r UserRepository) Find(ctx context.Context, userID string) (user *domain.User, err error) {
	var conn *grpc.ClientConn
	conn, err = r.dial(ctx)
	if err != nil {
		return nil, err
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	resp, err := userspb.NewUsersServiceClient(conn).GetUser(ctx, &userspb.GetUserRequest{Id: userID})
	if err != nil {
		return nil, err
	}

	return r.userToDomain(resp.User), nil
}
func (r UserRepository) userToDomain(user *userspb.User) *domain.User {
	return &domain.User{
		ID:        user.GetId(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Email:     user.GetEmail(),
	}
}

func (r UserRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}
