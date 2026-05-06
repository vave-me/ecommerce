package application

import (
	"context"
	"middleman/mailer/internal/models"
)

type UserRepository interface {
	Find(ctx context.Context, userId string) (*models.User, error)
}
