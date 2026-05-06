package domain

import (
	"context"
)

type FakeServiceRepository struct {
	services map[string]*Service
}

func NewFakeServiceRepository() *FakeServiceRepository {
	return &FakeServiceRepository{services: map[string]*Service{}}
}

var _ ServiceRepository = (*FakeServiceRepository)(nil)

func (r *FakeServiceRepository) Load(ctx context.Context, serviceID string) (*Service, error) {
	if service, exists := r.services[serviceID]; exists {
		return service, nil
	}

	return NewService(serviceID), nil
}

func (r *FakeServiceRepository) Save(ctx context.Context, service *Service) error {
	for _, event := range service.Events() {
		if err := service.ApplyEvent(event); err != nil {
			return err
		}
	}

	r.services[service.ID()] = service

	return nil
}

func (r *FakeServiceRepository) Reset(services ...*Service) {
	r.services = make(map[string]*Service)

	for _, service := range services {
		r.services[service.ID()] = service
	}
}
