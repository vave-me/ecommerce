package domain

type ProductStatus string

const (
	ProductStatusActive    ProductStatus = "active"
	ProductStatusLocked    ProductStatus = "locked"
	ProductStatusSold      ProductStatus = "sold"
	ProductStatusLeased    ProductStatus = "leased"
	ProductStatusPaused    ProductStatus = "paused"
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusArchived  ProductStatus = "archived"
	ProductStatusReference ProductStatus = "reference"
	ProductStatusUnknown   ProductStatus = ""
)

func (s ProductStatus) String() string {
	switch s {
	case ProductStatusActive, ProductStatusLocked, ProductStatusSold, ProductStatusLeased, ProductStatusPaused, ProductStatusDraft, ProductStatusArchived, ProductStatusReference:
		return string(s)
	default:
		return ""
	}
}

func ToProductStatus(s string) ProductStatus {
	switch s {
	case ProductStatusActive.String():
		return ProductStatusActive
	case ProductStatusLocked.String():
		return ProductStatusLocked
	case ProductStatusSold.String():
		return ProductStatusSold
	case ProductStatusLeased.String():
		return ProductStatusLeased
	case ProductStatusPaused.String():
		return ProductStatusPaused
	case ProductStatusDraft.String():
		return ProductStatusDraft
	case ProductStatusArchived.String():
		return ProductStatusArchived
	case ProductStatusReference.String():
		return ProductStatusReference
	default:
		return ProductStatusUnknown
	}
}
