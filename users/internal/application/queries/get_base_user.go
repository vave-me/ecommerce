package queries

import (
	"context"
	"middleman/users/internal/domain"
)

type GetBaseUser struct {
	ID string
}

type GetBaseUserHandler struct {
	middleman domain.MiddlemanRepository
}

func NewGetBaseUserHandler(middleman domain.MiddlemanRepository) GetBaseUserHandler {
	return GetBaseUserHandler{middleman: middleman}
}

func (h GetBaseUserHandler) GetBaseUser(ctx context.Context, query GetBaseUser) (*domain.MiddlemanViewUser, error) {
	return h.middleman.FindSimple(ctx, query.ID)
}
