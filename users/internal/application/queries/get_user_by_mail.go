package queries

import (
	"context"
	"middleman/users/internal/domain"
)

type GetUserByMail struct {
	Email string
}

type GetUserByMailHandler struct {
	middleman domain.MiddlemanRepository
}

func NewGetUserByMailHandler(middleman domain.MiddlemanRepository) GetUserByMailHandler {
	return GetUserByMailHandler{middleman: middleman}
}

func (h GetUserByMailHandler) GetUserByMail(ctx context.Context, query GetUserByMail) (*domain.MiddlemanUser, error) {
	return h.middleman.FindByEmail(ctx, query.Email)
}
