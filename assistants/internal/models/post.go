package models

type Post struct {
	Name         string
	Description  string
	PostID       string
	UserID       string
	PostType     string
	CategoryID   string
	CategorySlug string
	UserType     string
	Tags         []string
	Status       string
	Thumbnail    string
	Lat          float64
	Lng          float64
	EntityType   EntityType
	Metrics      *ItemMetric
}
