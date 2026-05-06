package models

import "time"

type Post struct {
	Name         string
	Description  string
	PostID       string
	UserID       string
	TypeOfPost   string
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
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
