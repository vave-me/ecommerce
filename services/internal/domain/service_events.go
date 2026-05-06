package domain

// ------------------------------------------------------------------------------------
// 1. Service Event Names
// ------------------------------------------------------------------------------------

const (
	ServiceAddedEvent             = "services.ServiceAdded"             // new service created
	ServiceUpdatedEvent           = "services.ServiceUpdated"           // new service created
	ServiceRebrandedEvent         = "services.ServiceRebranded"         // changes name, brand, model, or description
	ServicePriceIncreasedEvent    = "services.ServicePriceIncreased"    // service base price goes up
	ServicePriceDecreasedEvent    = "services.ServicePriceDecreased"    // service base price goes down
	ServiceStockAdjustedEvent     = "services.ServiceStockAdjusted"     // service stock changes (if manageStock is true)
	ServiceNegotiableToggledEvent = "services.ServiceNegotiableToggled" // toggles the negotiable field
	ServiceRemovedEvent           = "services.ServiceRemoved"           // service is removed from listing
	ServiceSoldEvent              = "services.ServiceSold"              // service is sold
	ServiceLeasedEvent            = "services.ServiceLeased"            // service is leased
	ServicePawnedEvent            = "services.ServicePawned"
	ServiceArchivedEvent          = "services.ServiceArchived"         // service is pawned
	ServiceThumbnailAddedEvent    = "services.ServiceThumbnailAdded"   // service is pawned
	ServiceThumbnailRemovedEvent  = "services.ServiceThumbnailRemoved" // service is pawned
	ServiceThumbnailUpdatedEvent  = "services.ServiceThumbnailUpdated" // service is pawned
)

// ------------------------------------------------------------------------------------
// 2. ServiceAdded
// This event sets the essential fields for a newly created service
// ------------------------------------------------------------------------------------

type ServiceAdded struct {
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
	Status           ServiceStatus
	UserType         UserType
	MiddlemanService bool
	ShippingCost     int64
	HasVariants      bool
	Negotiable       bool
	Attributes       []Attribute
	Options          []Option
	Thumbnail        string
	Lat              float64
	Lng              float64
	// If you want to store top-level service options, you could do:
	// Options []Option
}

func (ServiceAdded) Key() string { return ServiceAddedEvent }

type ServiceUpdated struct {
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
	Status           ServiceStatus
	UserType         UserType
	MiddlemanService bool
	ShippingCost     int64
	HasVariants      bool
	Negotiable       bool
	Attributes       []Attribute
	Options          []Option
	Thumbnail        string
	Lat              float64
	Lng              float64
	// If you want to store top-level service options, you could do:
	// Options []Option
}

func (ServiceUpdated) Key() string { return ServiceUpdatedEvent }

type ServiceRebranded struct {
	Name           string
	Description    string
	Qualifications []string
	Tags           []string
	// Possibly condition if rebranding includes changing condition
}

func (ServiceRebranded) Key() string { return ServiceRebrandedEvent }

// If you prefer storing a delta:
type ServicePriceIncreased struct {
	ServiceID string
	OldPrice  int64
	NewPrice  int64
}

func (ServicePriceIncreased) Key() string { return ServicePriceIncreasedEvent }

type ServicePriceDecreased struct {
	ServiceID string
	OldPrice  int64
	NewPrice  int64
}

func (ServicePriceDecreased) Key() string { return ServicePriceDecreasedEvent }

type ServiceStockAdjusted struct {
	ServiceID string
	OldStock  int64
	NewStock  int64
}

func (ServiceStockAdjusted) Key() string { return ServiceStockAdjustedEvent }

// ------------------------------------------------------------------------------------
// 6. ServiceNegotiableToggled
// Switches the 'Negotiable' field from true->false or false->true
// ------------------------------------------------------------------------------------

type ServiceNegotiableToggled struct {
	ServiceID string
	OldValue  bool
	NewValue  bool
}

func (ServiceNegotiableToggled) Key() string { return ServiceNegotiableToggledEvent }

// ------------------------------------------------------------------------------------
// 7. ServiceRemoved
// Service is removed from listing or from system
// ------------------------------------------------------------------------------------

type ServiceRemoved struct {
	ServiceID string
	UserID    string // or reason
}

func (ServiceRemoved) Key() string { return ServiceRemovedEvent }

type ServiceSold struct {
	ServiceID  string
	UserID     string
	BuyerID    string // if needed
	FinalPrice int64
	SoldAt     string // or a timestamp
}

func (ServiceSold) Key() string { return ServiceSoldEvent }

type ServiceLeased struct {
	ServiceID string
	UserID    string
	LesseeID  string
}

func (ServiceLeased) Key() string { return ServiceLeasedEvent }

type ServicePawned struct {
	ServiceID string
	UserID    string
	PawnID    string
}

func (ServicePawned) Key() string { return ServicePawnedEvent }

type ServiceArchived struct {
	ServiceID string
	UserID    string
}

func (ServiceArchived) Key() string { return ServiceArchivedEvent }

type ServiceThumbnailAdded struct {
	ServiceID string
	Thumbnail string
}

func (ServiceThumbnailAdded) Key() string { return ServiceThumbnailAddedEvent }

type ServiceThumbnailRemoved struct {
	ServiceID string
	Thumbnail string
}

func (ServiceThumbnailRemoved) Key() string { return ServiceThumbnailRemovedEvent }

type ServiceThumbnailUpdated struct {
	ServiceID string
	Thumbnail string
}

func (ServiceThumbnailUpdated) Key() string { return ServiceThumbnailUpdatedEvent }
