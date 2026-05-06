package models

type ServiceStatus string

const (
	ServiceStatusActive    ServiceStatus = "active"
	ServiceStatusLocked    ServiceStatus = "locked"
	ServiceStatusSold      ServiceStatus = "sold"
	ServiceStatusLeased    ServiceStatus = "leased"
	ServiceStatusPaused    ServiceStatus = "paused"
	ServiceStatusDraft     ServiceStatus = "draft"
	ServiceStatusArchived  ServiceStatus = "archived"
	ServiceStatusReference ServiceStatus = "reference"
	ServiceStatusUnknown   ServiceStatus = ""
)

func (s ServiceStatus) String() string {
	switch s {
	case ServiceStatusActive, ServiceStatusLocked, ServiceStatusSold, ServiceStatusLeased, ServiceStatusPaused, ServiceStatusDraft, ServiceStatusArchived, ServiceStatusReference:
		return string(s)
	default:
		return ""
	}
}

func ToServiceStatus(s string) ServiceStatus {
	switch s {
	case ServiceStatusActive.String():
		return ServiceStatusActive
	case ServiceStatusLocked.String():
		return ServiceStatusLocked
	case ServiceStatusSold.String():
		return ServiceStatusSold
	case ServiceStatusLeased.String():
		return ServiceStatusLeased
	case ServiceStatusPaused.String():
		return ServiceStatusPaused
	case ServiceStatusDraft.String():
		return ServiceStatusDraft
	case ServiceStatusArchived.String():
		return ServiceStatusArchived
	case ServiceStatusReference.String():
		return ServiceStatusReference
	default:
		return ServiceStatusUnknown
	}
}
