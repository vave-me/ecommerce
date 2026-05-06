package domain

type CommentV1 struct {
	ID         string
	SenderID   string
	ItemID     string
	ItemType   string
	CategoryID string
	ParentID   string
	Content    string
	Approved   bool
	Flagged    bool
}

func (CommentV1) SnapshotName() string { return "comments.CommentV1" }
