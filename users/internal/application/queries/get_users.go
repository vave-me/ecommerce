package queries

import (
	"context"
	"middleman/users/internal/domain"
)

type GetUsers struct{}

type GetUsersHandler struct {
	middleman domain.MiddlemanRepository
}

func NewGetUsersHandler(middleman domain.MiddlemanRepository) GetUsersHandler {
	return GetUsersHandler{middleman: middleman}
}

func (h GetUsersHandler) GetUsers(ctx context.Context, _ GetUsers) ([]*domain.MiddlemanUser, error) {
	return h.middleman.All(ctx)
}
