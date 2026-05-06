package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type KycVerify struct {
	UserID string
}

type KycVerifyHandler struct {
	users     domain.UserRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewKycVerifyHandler(users domain.UserRepository, publisher ddd.EventPublisher[ddd.Event],
) KycVerifyHandler {
	return KycVerifyHandler{
		users:     users,
		publisher: publisher,
	}
}
func (h KycVerifyHandler) KycVerify(ctx context.Context, cmd KycVerify) error {

	user, err := h.users.Load(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	event, err := user.KYCVerify(cmd.UserID)

	if err != nil {
		return err
	}
	return h.publisher.Publish(ctx, event)
}
