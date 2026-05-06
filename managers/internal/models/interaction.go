package models

type Interaction struct {
	ID         string
	ActivityID string
	ItemID     string // ID of the item being liked or disliked (e.g., listing, comment, product)
	ItemType   string
	ActionType string // The action, either "like" or "dislike"
}

type Activity struct {
	ID       string
	UserID   string
	Archived bool // Indicates whether the activity is archived
}

type MostReactionResult struct {
	ItemID   string
	ItemType string
	Action   string
	Count    int64
}
