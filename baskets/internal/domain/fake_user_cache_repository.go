package domain

import (
	"context"
	"github.com/stackus/errors"
)

type FakeUserCacheRepository struct {
	users map[string]*User
}

var _ UserCacheRepository = (*FakeUserCacheRepository)(nil)

func NewFakeUserCacheRepository() *FakeUserCacheRepository {
	return &FakeUserCacheRepository{users: map[string]*User{}}
}

func (r *FakeUserCacheRepository) Add(ctx context.Context, userSellerID, email, username, firstName, lastName, location string, enabled bool) error {
	r.users[userSellerID] = &User{
		ID:        userSellerID,
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Location:  location,
		Enabled:   enabled,
	}

	return nil
}
func (r *FakeUserCacheRepository) Rename(ctx context.Context, userSellerID, name string) error {
	if user, exists := r.users[userSellerID]; exists {
		user.FirstName = name
	}

	return nil
}
func (r *FakeUserCacheRepository) Find(ctx context.Context, userSellerID string) (*User, error) {
	if user, exists := r.users[userSellerID]; exists {
		return user, nil
	}

	return nil, errors.ErrNotFound.Msgf("user with id: `%s` does not exist", userSellerID)
}

func (r *FakeUserCacheRepository) Reset(users ...*User) {
	r.users = make(map[string]*User)

	for _, user := range users {
		r.users[user.ID] = user
	}
}
