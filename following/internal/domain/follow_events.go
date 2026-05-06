package domain

type FollowAdded struct {
	UserID           string
	FollowedUserID   string
	FollowedUserType FollowedUserType
	Content          string
	CategoryID       string
	ParentID         string
	Approved         bool
	Flagged          bool
}

// Key implements registry.Registerable

type FollowEdited struct {
	FollowID string
	Content  string
}

type FollowRemoved struct {
	FollowID string
}

type FollowApproved struct {
	Approved bool
}

type FollowFlagged struct {
	FollowID string
	Flagged  bool
}

type FollowRejected struct {
	FollowID string
	Approve  bool
}

// Key implements registry.Registerable
