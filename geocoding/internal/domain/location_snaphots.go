package domain

type LocationV1 struct {
	ID        string
	ProductID string
	Latitude  float64
	Longitude float64
}

func (LocationV1) SnapshotName() string { return "geocoding.LocationV1" }
