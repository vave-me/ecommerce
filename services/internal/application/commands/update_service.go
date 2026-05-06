package commands

import (
	"context"

	"middleman/internal/ddd"
	"middleman/services/internal/domain"
)

type UpdateService struct {
	ID               string
	Name             string
	Description      string
	ServiceType      string
	BasePrice        int64
	Pricing          []string
	Availability     string
	ProviderName     string
	UserID           string
	CategoryID       string
	CategorySlug     string
	DescriptionShort string
	DescriptionLong  string
	Qualifications   []string
	Contact          string
	Faq              string
	Tags             []string
	Status           domain.ServiceStatus
	UserType         domain.UserType
	ShippingCost     int64
	Negotiable       bool
	HasVariants      bool
	MiddlemanService bool
	Attributes       []domain.Attribute
	Options          []domain.Option
	Thumbnail        string
	Lat              float64
	Lng              float64
}

type UpdateServiceHandler struct {
	services  domain.ServiceRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewUpdateServiceHandler(
	services domain.ServiceRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UpdateServiceHandler {
	return UpdateServiceHandler{
		services:  services,
		publisher: publisher,
	}
}

func (h UpdateServiceHandler) UpdateService(ctx context.Context, cmd UpdateService) error {
	service, err := h.services.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := service.UpdateService(
		cmd.Name,
		cmd.Description,
		cmd.ServiceType,
		cmd.BasePrice,
		cmd.Pricing,
		cmd.Availability,
		cmd.ProviderName,
		cmd.UserID,
		cmd.CategoryID,
		cmd.CategorySlug,
		cmd.DescriptionShort,
		cmd.DescriptionLong,
		cmd.Qualifications,
		cmd.Contact,
		cmd.Faq,
		cmd.Tags,
		cmd.Status,
		cmd.UserType,
		cmd.MiddlemanService,
		cmd.ShippingCost,
		cmd.Attributes,
		cmd.Negotiable,
		cmd.HasVariants,
		cmd.Options,
		cmd.Thumbnail,
		cmd.Lat,
		cmd.Lng,
	)
	if err != nil {
		return err
	}

	if err = h.services.Save(ctx, service); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
