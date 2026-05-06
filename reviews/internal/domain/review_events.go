package domain

type ReviewAdded struct {
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

type ReviewEdited struct {
	ReviewID string
	Content  string
}

type ReviewRemoved struct {
	ReviewID string
}

type ReviewApproved struct {
	Approved bool
}

type ReviewFlagged struct {
	ReviewID string
	Flagged  bool
}

type ReviewRejected struct {
	ReviewID string
	Approve  bool
}

// Key implements registry.Registerable
