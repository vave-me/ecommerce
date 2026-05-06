package commands

import (
	"context"
	"log"

	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"

	"github.com/stackus/errors"
	"golang.org/x/crypto/bcrypt"
)

type LoginUser struct {
	UserID   string
	Email    string
	Password string
}

type LoginUserHandler struct {
	users     domain.UserRepository
	auth      auth.Authenticator
	publisher ddd.EventPublisher[ddd.Event]
}

func NewLoginUserHandler(
	middleman domain.MiddlemanRepository,
	users domain.UserRepository,
	auth auth.Authenticator,
	publisher ddd.EventPublisher[ddd.Event],
) LoginUserHandler {
	return LoginUserHandler{
		users:     users,
		auth:      auth,
		publisher: publisher,
	}
}

func (h LoginUserHandler) LoginUser(ctx context.Context, cmd LoginUser) (string, string, string, error) {

	user, err := h.users.Load(ctx, cmd.UserID)

	if err != nil {
		log.Printf("Load user error: %v, user firstname %v", err, user)
		return "", "", "", errors.Wrap(errors.ErrUnauthorized, "invalid email or password 3")
	}

	//if user.Enabled == false {
	//	log.Printf("Load user error: %v, user firstname %v", err, user)
	//	return "", "", "", errors.Wrap(errors.ErrUnauthorized, "Please confirm your mail and activate account")
	//}
	log.Printf("Loaded domain user: ID=%s, PasswordHash=%s", user.ID(), user.FirstName)

	// Verify password

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cmd.Password))
	if err != nil {
		log.Printf("Password comparison failed for user %s: %v", user.ID(), err)
		return "", "", "", errors.Wrap(errors.ErrUnauthorized, "invalid email or password 5")
	}

	log.Printf("Password verified for user: %s", user.ID())

	tokens, err := h.auth.GenerateTokenPair(&auth.JwtUser{
		ID:        user.ID(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Lat:       user.Lat,
		Long:      user.Lang,
		Role:      user.Role.String(),
	})

	event, err := user.Login()

	err = h.publisher.Publish(ctx, event)
	if err != nil {
		return "", "", "", err
	}

	return tokens.AccessToken, tokens.RefreshToken, user.Username, nil
}
