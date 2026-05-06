package domain

const (
	AddressCreatedEvent      = "geocoding.AddressCreated"
	AddressBatchDecodedEvent = "geocoding.AddressBatchDecoded"
)

type AddressCreated struct {
	Address   string
	Latitude  float64
	Longitude float64
}

func (AddressCreated) Key() string { return AddressCreatedEvent }

type AddressBatchDecoded struct {
	Address string
}
