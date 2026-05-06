package domain

import "context"

type ServiceRepository interface {
	Load(ctx context.Context, id string) (*Service, error)
	Save(ctx context.Context, service *Service) error
}
