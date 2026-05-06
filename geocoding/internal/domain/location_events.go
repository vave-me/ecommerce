package domain

const (
	LocationAddedEvent = "geocoding.LocationAdded"
)

type LocationAdded struct {
	Address string
}

func (LocationAdded) Key() string { return LocationAddedEvent }
