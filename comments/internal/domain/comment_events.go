package domain

type CommentAdded struct {
	SenderID   string
	ItemID     string
	ItemType   ItemType
	Content    string
	CategoryID string
	ParentID   string
	Approved   bool
	Flagged    bool
}

// Key implements registry.Registerable

type CommentEdited struct {
	CommentID string
	Content   string
}

type CommentRemoved struct {
	CommentID string
}

type CommentApproved struct {
	Approved bool
}

type CommentFlagged struct {
	CommentID string
	Flagged   bool
}

type CommentRejected struct {
	CommentID string
	Approve   bool
}

// Key implements registry.Registerable
