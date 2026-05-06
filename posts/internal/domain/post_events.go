package domain

// ------------------------------------------------------------------------------------
// 1. Post Event Names
// ------------------------------------------------------------------------------------

const (
	PostAddedEvent            = "posts.PostAdded"            // new post created
	PostUpdatedEvent          = "posts.PostUpdated"          // new post created
	PostRebrandedEvent        = "posts.PostRebranded"        // changes name, brand, model, or description
	PostRemovedEvent          = "posts.PostRemoved"          // post is removed from listing
	PostArchivedEvent         = "posts.PostArchived"         // post is pawned
	PostThumbnailAddedEvent   = "posts.PostThumbnailAdded"   // post is pawned
	PostThumbnailRemovedEvent = "posts.PostThumbnailRemoved" // post is pawned
	PostThumbnailUpdatedEvent = "posts.PostThumbnailUpdated" // post is pawned
)

type PostAdded struct {
	ID           string
	Name         string
	Description  string
	TypeOfPost   TypeOfPost
	UserID       string
	UserType     UserType
	CategoryID   string
	CategorySlug string
	Tags         []string
	Status       PostStatus
	Thumbnail    string
	Lat          float64
	Lng          float64
}

func (PostAdded) Key() string { return PostAddedEvent }

type PostUpdated struct {
	Name         string
	Description  string
	TypeOfPost   TypeOfPost
	UserID       string
	UserType     UserType
	CategoryID   string
	CategorySlug string
	Tags         []string
	Status       PostStatus
	Thumbnail    string
	Lat          float64
	Lng          float64
}

func (PostUpdated) Key() string { return PostUpdatedEvent }

type PostRebranded struct {
	Name        string
	Description string
	Brand       string
	Model       string
	// Possibly condition if rebranding includes changing condition
}

func (PostRebranded) Key() string { return PostRebrandedEvent }

type PostRemoved struct {
	PostID string
	UserID string // or reason
}

func (PostRemoved) Key() string { return PostRemovedEvent }

type PostArchived struct {
	PostID string
	UserID string
}

func (PostArchived) Key() string { return PostArchivedEvent }

type PostThumbnailAdded struct {
	PostID    string
	Thumbnail string
}

func (PostThumbnailAdded) Key() string { return PostThumbnailAddedEvent }

type PostThumbnailRemoved struct {
	PostID    string
	Thumbnail string
}

func (PostThumbnailRemoved) Key() string { return PostThumbnailRemovedEvent }

type PostThumbnailUpdated struct {
	PostID    string
	Thumbnail string
}

func (PostThumbnailUpdated) Key() string { return PostThumbnailUpdatedEvent }
