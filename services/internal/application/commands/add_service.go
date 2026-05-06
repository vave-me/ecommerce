package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/services/internal/domain"
)

type AddService struct {
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

type AddServiceHandler struct {
	services  domain.ServiceRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAddServiceHandler(
	services domain.ServiceRepository, publisher ddd.EventPublisher[ddd.Event]) AddServiceHandler {

	return AddServiceHandler{
		services:  services,
		publisher: publisher,
	}
}

func (h AddServiceHandler) AddService(ctx context.Context, cmd AddService) error {
	service, err := h.services.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error adding service")
	}

	event, err := service.InitService(
		cmd.ID,
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
		return errors.Wrap(err, "initializing service")
	}

	err = h.services.Save(ctx, service)
	if err != nil {
		return errors.Wrap(err, "error adding service")
	}

	return errors.Wrap(h.publisher.Publish(ctx, event), "publishing domain event")
}
