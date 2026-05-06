package application

import (
	"context"
	"middleman/search/internal/models"
)

type UserRepository interface {
	Find(ctx context.Context, userID string) (*models.User, error)
}

type UserCacheRepository interface {
	Add(ctx context.Context, userID, email, username, firstName, lastName, location string, enabled bool) error
	Rename(ctx context.Context, userID string, firstName string) error
	Update(ctx context.Context, user *models.User) error
	//Search(ctx context.Context, search SearchUsers) ([]*models.User, error)
	Find(ctx context.Context, userID string) (*models.User, error)
	UserRepository
}
