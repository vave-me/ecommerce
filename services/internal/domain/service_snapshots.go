package domain

type ServiceV1 struct {
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

func (ServiceV1) SnapshotName() string { return "services.ServiceV1" }
