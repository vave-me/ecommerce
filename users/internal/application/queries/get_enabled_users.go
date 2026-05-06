package queries

import (
	"context"
	"middleman/users/internal/domain"
)

type GetEnabledUsers struct{}

type GetEnabledUsersHandler struct {
	middleman domain.MiddlemanRepository
}

func NewGetEnabledUsersHandler(middleman domain.MiddlemanRepository) GetEnabledUsersHandler {
	return GetEnabledUsersHandler{middleman: middleman}
}

func (h GetEnabledUsersHandler) GetEnabledUsers(ctx context.Context, _ GetEnabledUsers) ([]*domain.MiddlemanUser, error) {
	return h.middleman.AllEnabled(ctx)
}