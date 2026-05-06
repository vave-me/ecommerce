package domain

type CategoryStatus string

const (
	CategoryStatusActive    CategoryStatus = "active"
	CategoryStatusLocked    CategoryStatus = "locked"
	CategoryStatusSold      CategoryStatus = "sold"
	CategoryStatusLeased    CategoryStatus = "leased"
	CategoryStatusPaused    CategoryStatus = "paused"
	CategoryStatusDraft     CategoryStatus = "draft"
	CategoryStatusArchived  CategoryStatus = "archived"
	CategoryStatusReference CategoryStatus = "reference"
	CategoryStatusUnknown   CategoryStatus = ""
)

func (s CategoryStatus) String() string {
	switch s {
	case CategoryStatusActive, CategoryStatusLocked, CategoryStatusSold, CategoryStatusLeased, CategoryStatusPaused, CategoryStatusDraft, CategoryStatusArchived, CategoryStatusReference:
		return string(s)
	default:
		return ""
	}
}

func ToCategoryStatus(s string) CategoryStatus {
	switch s {
	case CategoryStatusActive.String():
		return CategoryStatusActive
	case CategoryStatusLocked.String():
		return CategoryStatusLocked
	case CategoryStatusSold.String():
		return CategoryStatusSold
	case CategoryStatusLeased.String():
		return CategoryStatusLeased
	case CategoryStatusPaused.String():
		return CategoryStatusPaused
	case CategoryStatusDraft.String():
		return CategoryStatusDraft
	case CategoryStatusArchived.String():
		return CategoryStatusArchived
	case CategoryStatusReference.String():
		return CategoryStatusReference
	default:
		return CategoryStatusUnknown
	}
}
