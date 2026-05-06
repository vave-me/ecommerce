package domain

type CategoryV1 struct {
	ID               string
	Description      string
	ParentID         string
	GoogleCategoryID string
	Tags             []string
	IsActive         bool
	Slug             string
	SeoTitle         string
	SeoKeywords      []string
	SeoDesc          string
}

func (CategoryV1) SnapshotName() string { return "categories.CategoryV1" }
