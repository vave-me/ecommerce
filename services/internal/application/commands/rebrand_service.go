package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/services/internal/domain"
)

type RebrandService struct {
	ID             string
	Name           string
	Description    string
	Qualifications []string
	Faq            string
	Tags           []string
}

type RebrandServiceHandler struct {
	services  domain.ServiceRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRebrandServiceHandler(services domain.ServiceRepository, publisher ddd.EventPublisher[ddd.Event]) RebrandServiceHandler {
	return RebrandServiceHandler{
		services:  services,
		publisher: publisher,
	}
}

func (h RebrandServiceHandler) RebrandService(ctx context.Context, cmd RebrandService) error {
	service, err := h.services.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := service.RebrandService(cmd.Name, cmd.Description, cmd.Faq, cmd.Qualifications, cmd.Tags)
	if err != nil {
		return err
	}

	err = h.services.Save(ctx, service)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
