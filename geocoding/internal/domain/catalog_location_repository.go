package domain

type CatalogLocation struct {
	Address string
	Lat     float64
	Lng     float64
}

type CatalogLocationRepository interface {
}
