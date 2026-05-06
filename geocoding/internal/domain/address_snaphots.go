package domain

type AddressV1 struct {
	ID               string
	FormattedAddress string
	UserID           string
	StreetName       string
	StreetNumber     int
	StreetSuffix     string
	City             string
	State            string
	Zip              string
	Latitude         float64
	Longitude        float64
}

func (AddressV1) SnapshotName() string { return "geocoding.AddressV1" }
