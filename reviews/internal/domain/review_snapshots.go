package domain

type ReviewV1 struct {
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

func (ReviewV1) SnapshotName() string { return "reviews.ReviewV1" }
