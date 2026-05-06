package domain

type FollowV1 struct {
	ID               string
	UserID           string
	FollowedUserID   string
	FollowedUserType string
	CategoryID       string
	ParentID         string
	Content          string
	Approved         bool
	Flagged          bool
}

func (FollowV1) SnapshotName() string { return "following.FollowV1" }
