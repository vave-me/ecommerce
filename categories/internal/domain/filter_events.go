package domain

// ------------------------------------------------------------------------------------
// 1. Filter Event Names
// ------------------------------------------------------------------------------------

const (
	FilterAddedEvent     = "categories.FilterAdded"
	FilterArchivedEvent  = "categories.FilterArchived"
	FilterRemovedEvent   = "categories.FilterRemoved"
	FilterRebrandedEvent = "categories.FilterRebranded"
	FilterUpdatedEvent   = "categories.FilterUpdated" // optional
)

// FilterAdded is used when a new Filter is created.
type FilterAdded struct {
	CategoryID string
	Name       string
	FilterType FilterType
	Values     []string
	IsActive   bool
}

func (FilterAdded) Key() string { return FilterAddedEvent }

// FilterUpdated modifies the Filter’s core fields (optional).
type FilterUpdated struct {
	Name       string
	FilterType FilterType
	Values     []string
}

func (FilterUpdated) Key() string { return FilterUpdatedEvent }

// FilterRebranded can rename the Filter or update textual fields (like name).
type FilterRebranded struct {
	FilterID string
	Name     string
}

func (FilterRebranded) Key() string { return FilterRebrandedEvent }

// FilterArchived marks the filter as no longer active.
type FilterArchived struct {
	FilterID string
}

func (FilterArchived) Key() string { return FilterArchivedEvent }

// FilterRemoved means the filter is removed from the domain entirely.
type FilterRemoved struct {
	FilterID string
}

func (FilterRemoved) Key() string { return FilterRemovedEvent }
