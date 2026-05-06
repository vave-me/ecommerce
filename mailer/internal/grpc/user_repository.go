package grpc

import (
	"context"
	"google.golang.org/grpc"
	"middleman/internal/rpc"
	"middleman/mailer/internal/application"
	"middleman/mailer/internal/models"
	"middleman/users/userspb"
)

type UserRepository struct {
	endpoint string
}

var _ application.UserRepository = (*UserRepository)(nil)

func NewUserRepository(endpoint string) UserRepository {
	return UserRepository{
		endpoint: endpoint,
	}
}

func (r UserRepository) Find(ctx context.Context, userID string) (user *models.User, err error) {
	var conn *grpc.ClientConn
	conn, err = r.dial(ctx)
	if err != nil {
		return nil, err
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	resp, err := userspb.NewUsersServiceClient(conn).GetBaseUser(ctx, &userspb.GetBaseUserRequest{Id: userID})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        resp.GetUser().GetId(),
		Username:  "test",
		FirstName: "test",
		Email:     "redacted-email@example.com",
		LastName:  "test",
		Enabled:   true,
	}, nil
}

func (r UserRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}
