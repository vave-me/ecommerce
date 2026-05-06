package models

import "time"

type Service struct {
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
	Status           string
	UserType         string
	ShippingCost     int64
	HasVariants      bool
	MiddlemanService bool
	Negotiable       bool
	Attributes       []Attribute
	Options          []Option
	Thumbnail        string
	Lat              float64
	Lng              float64
	EntityType       EntityType
	Metrics          *ItemMetric
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
