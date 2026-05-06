package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const ServiceAggregate = "services.Service"

// Domain-level errors
var (
	ErrServiceNameBlank         = errors.Wrap(errors.ErrBadRequest, "the offer name cannot be blank")
	ErrServicePriceNegative     = errors.Wrap(errors.ErrBadRequest, "the offer price cannot be negative")
	ErrServicePriceNotIncreased = errors.Wrap(errors.ErrBadRequest, "the price change would not be an increase")
	ErrServicePriceNotDecreased = errors.Wrap(errors.ErrBadRequest, "the price change would not be a decrease")
)

type Service struct {
	es.Aggregate
	Name             string
	Description      string
	ServiceType      string
	Availability     string
	ProviderName     string
	UserID           string
	CategoryID       string
	CategorySlug     string
	DescriptionShort string
	DescriptionLong  string
	Attributes       []Attribute
	BasePrice        int64
	Pricing          []string
	Qualifications   []string
	Contact          string
	Faq              string
	Tags             []string
	Status           ServiceStatus
	UserType         UserType
	MiddlemanService bool
	Negotiable       bool
	HasVariants      bool
	ShippingCost     int64
	Options          []Option
	Thumbnail        string
	Lat              float64
	Lng              float64
}

// Ensure Service implements the needed interfaces
var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Service)(nil)

// Constructor
func NewService(id string) *Service {
	return &Service{
		Aggregate: es.NewAggregate(id, ServiceAggregate),
	}
}

// Key implements registry.Registerable (if you have such a registry)
func (Service) Key() string {
	return ServiceAggregate
}

// --- Command-like Methods ---

func (o *Service) InitService(
	id, name, description string,
	serviceType string,
	basePrice int64,
	pricing []string,
	availability, providerName, userID, categoryID, categorySlug, descriptionShort, descriptionLong string,
	qualifications []string,
	contact string,
	faq string,
	tags []string,
	status ServiceStatus,
	userType UserType,
	middlemanService bool,
	shippingCost int64,
	attributes []Attribute,
	negotiable bool,
	hasVariants bool,
	options []Option,
	thumbnail string,
	lat, lng float64,
) (ddd.Event, error) {

	if name == "" {
		return nil, ErrServiceNameBlank
	}
	if basePrice < 0 {
		return nil, ErrServicePriceNegative
	}

	o.AddEvent(ServiceAddedEvent, &ServiceAdded{
		Name:             name,
		Description:      description,
		ServiceType:      serviceType,
		BasePrice:        basePrice,
		Pricing:          pricing,
		Availability:     availability,
		ProviderName:     providerName,
		UserID:           userID,
		CategoryID:       categoryID,
		CategorySlug:     categorySlug,
		DescriptionShort: descriptionShort,
		DescriptionLong:  descriptionLong,
		Qualifications:   qualifications,
		Contact:          contact,
		Faq:              faq,
		Tags:             tags,
		Attributes:       attributes,
		Status:           status,
		Negotiable:       negotiable,
		UserType:         userType,
		MiddlemanService: middlemanService,
		ShippingCost:     shippingCost,
		HasVariants:      hasVariants,
		Options:          options,
		Thumbnail:        thumbnail,
		Lat:              lat,
		Lng:              lng,
	})
	return ddd.NewEvent(ServiceAddedEvent, o), nil
}

func (o *Service) UpdateService(
	name, description string,
	serviceType string,
	basePrice int64,
	pricing []string,
	availability, providerName, userID, categoryID, categorySlug, descriptionShort, descriptionLong string,
	qualifications []string,
	contact string,
	faq string,
	tags []string,
	status ServiceStatus,
	userType UserType,
	middlemanService bool,
	shippingCost int64,
	attributes []Attribute,
	negotiable bool,
	hasVariants bool,
	options []Option,
	thumbnail string,
	lat, lng float64,
) (ddd.Event, error) {

	if name == "" {
		name = o.Name
	}
	if description == "" {
		description = o.Description
	}
	if serviceType == "" {
		serviceType = o.ServiceType
	}
	if basePrice < 0 {
		basePrice = o.BasePrice
	}
	if pricing == nil {
		pricing = o.Pricing
	}
	if availability == "" {
		availability = o.Availability
	}
	if providerName == "" {
		providerName = o.ProviderName
	}
	if userID == "" {
		userID = o.UserID
	}
	if categoryID == "" {
		categoryID = o.CategoryID
	}
	if categorySlug == "" {
		categorySlug = o.CategorySlug
	}
	if descriptionShort == "" {
		descriptionShort = o.DescriptionShort
	}
	if descriptionLong == "" {
		descriptionLong = o.DescriptionLong
	}
	if qualifications == nil {
		qualifications = o.Qualifications
	}
	if contact == "" {
		contact = o.Contact
	}
	if faq == "" {
		faq = o.Faq
	}
	if tags == nil {
		tags = o.Tags
	}
	if status == "" {
		status = o.Status
	}
	if userType == "" {
		userType = o.UserType
	}
	if middlemanService {
		middlemanService = o.MiddlemanService
	}
	if shippingCost == 0 {
		shippingCost = o.ShippingCost
	}
	if attributes == nil {
		attributes = o.Attributes
	}
	if hasVariants {
		hasVariants = o.HasVariants
	}
	if options == nil {
		options = o.Options
	}
	if thumbnail == "" {
		thumbnail = o.Thumbnail
	}
	if lat == 0 && lng == 0 {
		lat, lng = o.Lat, o.Lng
	}
	o.AddEvent(ServiceUpdatedEvent, &ServiceUpdated{
		Name:             name,
		Description:      description,
		ServiceType:      serviceType,
		BasePrice:        basePrice,
		Pricing:          pricing,
		Availability:     availability,
		ProviderName:     providerName,
		UserID:           userID,
		CategoryID:       categoryID,
		CategorySlug:     categorySlug,
		DescriptionShort: descriptionShort,
		DescriptionLong:  descriptionLong,
		Qualifications:   qualifications,
		Contact:          contact,
		Faq:              faq,
		Tags:             tags,
		Attributes:       attributes,
		Status:           status,
		Negotiable:       negotiable,
		UserType:         userType,
		MiddlemanService: middlemanService,
		ShippingCost:     shippingCost,
		HasVariants:      hasVariants,
		Options:          options,
		Thumbnail:        thumbnail,
		Lat:              lat,
		Lng:              lng,
	})
	return ddd.NewEvent(ServiceUpdatedEvent, o), nil
}

func (o *Service) RebrandService(name, description, faq string, qualifications, tags []string) (ddd.Event, error) {
	o.AddEvent(ServiceRebrandedEvent, &ServiceRebranded{
		Name:           name,
		Description:    description,
		Qualifications: qualifications,
		Tags:           tags,
	})
	return ddd.NewEvent(ServiceRebrandedEvent, o), nil
}

func (o *Service) IncreaseServicePrice(newPrice int64) (ddd.Event, error) {
	if newPrice < o.BasePrice {
		return nil, ErrServicePriceNotIncreased
	}
	o.AddEvent(ServicePriceIncreasedEvent, &ServicePriceIncreased{
		ServiceID: o.ID(),
		OldPrice:  o.BasePrice,
		NewPrice:  newPrice,
	})
	return ddd.NewEvent(ServicePriceIncreasedEvent, o), nil
}

func (o *Service) DecreaseServicePrice(newPrice int64) (ddd.Event, error) {
	if newPrice > o.BasePrice {
		return nil, ErrServicePriceNotDecreased
	}
	o.AddEvent(ServicePriceDecreasedEvent, &ServicePriceDecreased{
		ServiceID: o.ID(),
		OldPrice:  o.BasePrice,
		NewPrice:  newPrice,
	})
	return ddd.NewEvent(ServicePriceDecreasedEvent, o), nil
}

func (o *Service) RemoveService(id, userSellerID string) (ddd.Event, error) {
	o.AddEvent(ServiceRemovedEvent, &ServiceRemoved{
		ServiceID: id,
		UserID:    userSellerID,
	})
	return ddd.NewEvent(ServiceRemovedEvent, o), nil
}

func (o *Service) ArchiveService(userSellerID string) (ddd.Event, error) {
	o.AddEvent(ServiceArchivedEvent, &ServiceArchived{
		ServiceID: o.ID(),
		UserID:    userSellerID,
	})
	return ddd.NewEvent(ServiceArchivedEvent, o), nil
}

func (o *Service) MarkServiceSold(userSellerID string, finalPrice int64) (ddd.Event, error) {
	o.AddEvent(ServiceSoldEvent, &ServiceSold{
		ServiceID:  o.ID(),
		UserID:     userSellerID,
		FinalPrice: finalPrice,
	})
	return ddd.NewEvent(ServiceSoldEvent, o), nil
}

func (o *Service) MarkServiceLeased(userSellerID string, monthlyPrice, leaseTermMonths int64) (ddd.Event, error) {
	o.AddEvent(ServiceLeasedEvent, &ServiceLeased{
		ServiceID: o.ID(),
		UserID:    userSellerID,
	})
	return ddd.NewEvent(ServiceLeasedEvent, o), nil
}

func (o *Service) MarkServicePawned(userSellerID string, lockedPrice, redemptionFee int64) (ddd.Event, error) {
	o.AddEvent(ServicePawnedEvent, &ServicePawned{
		ServiceID: o.ID(),
		UserID:    userSellerID,
	})
	return ddd.NewEvent(ServicePawnedEvent, o), nil
}

// --- Event Applier ---

func (o *Service) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {

	case *ServiceAdded:
		o.Name = e.Name
		o.Description = e.Description
		o.BasePrice = e.BasePrice
		o.ServiceType = e.ServiceType
		o.Pricing = e.Pricing
		o.Availability = e.Availability
		o.ProviderName = e.ProviderName
		o.UserID = e.UserID
		o.CategoryID = e.CategoryID
		o.CategorySlug = e.CategorySlug
		o.DescriptionShort = e.DescriptionShort
		o.DescriptionLong = e.DescriptionLong
		o.Attributes = e.Attributes
		o.Qualifications = e.Qualifications
		o.Contact = e.Contact
		o.Faq = e.Faq
		o.Tags = e.Tags
		o.Status = e.Status
		o.UserType = e.UserType
		o.Negotiable = e.Negotiable
		o.MiddlemanService = e.MiddlemanService
		o.ShippingCost = e.ShippingCost
		o.HasVariants = e.HasVariants
		o.Options = e.Options
		o.Thumbnail = e.Thumbnail
		o.Lat = e.Lat
		o.Lng = e.Lng

	case *ServiceUpdated:
		o.Name = e.Name
		o.Description = e.Description
		o.BasePrice = e.BasePrice
		o.ServiceType = e.ServiceType
		o.Pricing = e.Pricing
		o.Availability = e.Availability
		o.ProviderName = e.ProviderName
		o.UserID = e.UserID
		o.CategoryID = e.CategoryID
		o.CategorySlug = e.CategorySlug
		o.DescriptionShort = e.DescriptionShort
		o.DescriptionLong = e.DescriptionLong
		o.Attributes = e.Attributes
		o.Qualifications = e.Qualifications
		o.Contact = e.Contact
		o.Faq = e.Faq
		o.Tags = e.Tags
		o.Status = e.Status
		o.UserType = e.UserType
		o.Negotiable = e.Negotiable
		o.MiddlemanService = e.MiddlemanService
		o.ShippingCost = e.ShippingCost
		o.HasVariants = e.HasVariants
		o.Options = e.Options
		o.Thumbnail = e.Thumbnail
		o.Lat = e.Lat
		o.Lng = e.Lng

	case *ServiceRebranded:
		o.Name = e.Name
		o.Description = e.Description

	case *ServicePriceIncreased:
		o.BasePrice = e.NewPrice

	case *ServicePriceDecreased:
		o.BasePrice = e.NewPrice

	case *ServiceStockAdjusted:

	case *ServiceArchived:
		// e.g. o.Status = "Archived"

	case *ServiceSold:
		// e.g. o.Status = "Sold"
		// o.BasePrice = e.FinalPrice (if you track final sale price)

	case *ServiceLeased:
		// e.g. o.Status = "Leased"

	case *ServicePawned:
		// e.g. o.Status = "Pawned"

	case *ServiceRemoved:
		// Possibly mark internal "removed" = true or do nothing if ephemeral

	default:
		return errors.ErrInternal.Msgf(
			"%T received event %s with unexpected payload %T",
			o, event.EventName(), e,
		)
	}
	return nil
}

// --- Snapshotter Methods ---

func (o *Service) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ServiceV1:
		o.Name = ss.Name
		o.Description = ss.Description
		o.BasePrice = ss.BasePrice
		o.ServiceType = ss.ServiceType
		o.Pricing = ss.Pricing
		o.Availability = ss.Availability
		o.ProviderName = ss.ProviderName
		o.UserID = ss.UserID
		o.CategoryID = ss.CategoryID
		o.CategorySlug = ss.CategorySlug
		o.DescriptionShort = ss.DescriptionShort
		o.DescriptionLong = ss.DescriptionLong
		o.Attributes = ss.Attributes
		o.Qualifications = ss.Qualifications
		o.Contact = ss.Contact
		o.Faq = ss.Faq
		o.Tags = ss.Tags
		o.Attributes = ss.Attributes
		o.Status = ss.Status
		o.Negotiable = ss.Negotiable
		o.MiddlemanService = ss.MiddlemanService
		o.ShippingCost = ss.ShippingCost
		o.HasVariants = ss.HasVariants
		o.Options = ss.Options
		o.Thumbnail = ss.Thumbnail
		o.Lat = ss.Lat
		o.Lng = ss.Lng

	default:
		return errors.ErrInternal.Msgf(
			"%T received an unexpected snapshot %T",
			o, snapshot,
		)
	}
	return nil
}

func (o Service) ToSnapshot() es.Snapshot {
	// Return a snapshot of the current fields
	return &ServiceV1{
		Name:             o.Name,
		Description:      o.Description,
		ServiceType:      o.ServiceType,
		Availability:     o.Availability,
		ProviderName:     o.ProviderName,
		UserID:           o.UserID,
		BasePrice:        o.BasePrice,
		Pricing:          o.Pricing,
		CategoryID:       o.CategoryID,
		CategorySlug:     o.CategorySlug,
		DescriptionShort: o.DescriptionShort,
		DescriptionLong:  o.DescriptionLong,
		Tags:             o.Tags,
		Attributes:       o.Attributes,
		Status:           o.Status,
		Negotiable:       o.Negotiable,
		Qualifications:   o.Qualifications,
		Contact:          o.Contact,
		Faq:              o.Faq,
		UserType:         o.UserType,
		MiddlemanService: o.MiddlemanService,
		ShippingCost:     o.ShippingCost,
		HasVariants:      o.HasVariants,
		Options:          o.Options,
		Thumbnail:        o.Thumbnail,
		Lat:              o.Lat,
		Lng:              o.Lng,
	}
}
