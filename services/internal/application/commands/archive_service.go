package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/services/internal/domain"
)

// ArchiveService marks a service as archived or inactive.
type ArchiveService struct {
	ID           string
	UserSellerID string
}

type ArchiveServiceHandler struct {
	services  domain.ServiceRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewArchiveServiceHandler(
	services domain.ServiceRepository,
	publisher ddd.EventPublisher[ddd.Event],
) ArchiveServiceHandler {
	return ArchiveServiceHandler{
		services:  services,
		publisher: publisher,
	}
}

func (h ArchiveServiceHandler) ArchiveService(ctx context.Context, cmd ArchiveService) error {
	service, err := h.services.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := service.ArchiveService(cmd.UserSellerID) // domain: service.Archive(userSellerID)
	if err != nil {
		return err
	}

	if err = h.services.Save(ctx, service); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
