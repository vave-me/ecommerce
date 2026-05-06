package queries

import (
	"context"
	"middleman/users/internal/domain"
)

type GetUser struct {
	ID string
}

type GetUserHandler struct {
	middleman domain.MiddlemanRepository
}

func NewGetUserHandler(middleman domain.MiddlemanRepository) GetUserHandler {
	return GetUserHandler{middleman: middleman}
}

func (h GetUserHandler) GetUser(ctx context.Context, query GetUser) (*domain.MiddlemanUser, error) {
	return h.middleman.Find(ctx, query.ID)
}
