package queries

import (
	"context"
	"middleman/users/internal/domain"
)

type GetUserByGoogleID struct {
	GoogleID string
}

type GetUserByGoogleIDHandler struct {
	middleman domain.MiddlemanRepository
}

func NewGetUserByGoogleIDHandler(middleman domain.MiddlemanRepository) GetUserByGoogleIDHandler {
	return GetUserByGoogleIDHandler{middleman: middleman}
}

func (h GetUserByGoogleIDHandler) GetUserByGoogleID(ctx context.Context, query GetUserByGoogleID) (*domain.MiddlemanUser, error) {
	return h.middleman.FindByGoogleID(ctx, query.GoogleID)
}
