package domain

import (
	"context"
)

type FakeUserRepository struct {
	users map[string]*User
}

func NewFakeUserRepository() *FakeUserRepository {
	return &FakeUserRepository{users: map[string]*User{}}
}

var _ UserRepository = (*FakeUserRepository)(nil)

func (r *FakeUserRepository) Load(ctx context.Context, userID string) (*User, error) {
	if user, exists := r.users[userID]; exists {
		return user, nil
	}

	return NewUser(userID), nil
}

func (r *FakeUserRepository) Save(ctx context.Context, user *User) error {
	for _, event := range user.Events() {
		if err := user.ApplyEvent(event); err != nil {
			return err
		}
	}

	r.users[user.ID()] = user

	return nil
}

func (r *FakeUserRepository) Reset(users ...*User) {
	r.users = make(map[string]*User)

	for _, user := range users {
		r.users[user.ID()] = user
	}
}
