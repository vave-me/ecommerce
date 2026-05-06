package domain

// ------------------------------------------------------------------------------------
// 2. Category Event Names
// ------------------------------------------------------------------------------------

const (
	CategoryAddedEvent     = "categories.CategoryAdded"
	CategoryUpdatedEvent   = "categories.CategoryUpdated"
	CategoryRebrandedEvent = "categories.CategoryRebranded"
	CategoryRemovedEvent   = "categories.CategoryRemoved"
	CategoryArchivedEvent  = "categories.CategoryArchived"
)

// CategoryAdded represents creation of a new Category entity.
type CategoryAdded struct {
	CategoryID       string
	Description      string
	ParentID         string
	GoogleCategoryID string
	Tags             []string // optional set of tags
	IsActive         bool
	Slug             string
	SeoTitle         string
	SeoKeywords      []string
	SeoDesc          string
	CategoryType     string
	Lang             string
}

func (CategoryAdded) Key() string { return CategoryAddedEvent }

// CategoryUpdated updates the Category’s main fields.
type CategoryUpdated struct {
	CategoryID       string
	Description      string
	ParentID         string
	GoogleCategoryID string
	Tags             []string // optional set of tags
	IsActive         bool
	Slug             string
	SeoTitle         string
	SeoKeywords      []string
	SeoDesc          string
	CategoryType     string
	Lang             string
}

func (CategoryUpdated) Key() string { return CategoryUpdatedEvent }

// CategoryRebranded is for changing only name/description or “brand”-like fields.
type CategoryRebranded struct {
	Slug        string
	Description string
}

func (CategoryRebranded) Key() string { return CategoryRebrandedEvent }

// CategoryRemoved means the category is fully removed (deleted).
type CategoryRemoved struct {
	CategoryID string
	UserID     string // Possibly store the user or reason
}

func (CategoryRemoved) Key() string { return CategoryRemovedEvent }

// CategoryArchived marks the category as archived/inactive.
type CategoryArchived struct {
	CategoryID string
}

func (CategoryArchived) Key() string { return CategoryArchivedEvent }
