package domain

type PostV1 struct {
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

func (PostV1) SnapshotName() string { return "posts.PostV1" }
