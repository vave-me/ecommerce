package models

type FuelType string

const (
	Petrol   FuelType = "petrol"
	Electric FuelType = "electric"
	Diesel   FuelType = "diesel"
	Hybrid   FuelType = "hybrid"
	Other    FuelType = "other"
)

func (t FuelType) String() string {
	switch t {
	case Petrol, Electric, Diesel, Hybrid:
		return string(t)
	default:
		return ""
	}
}

func ToFuelType(s string) FuelType {
	switch s {
	case Petrol.String():
		return Petrol
	case Electric.String():
		return Electric
	case Diesel.String():
		return Diesel
	default:
		return Other
	}
}
