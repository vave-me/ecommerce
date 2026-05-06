package domain

type InteractionVi struct {
	ActivityID string
	ItemID     string // ID of the item being liked or disliked (e.g., listing, comment, product)
	ItemType   string
	ActionType string // The action, either "like" or "dislike"
}

func (InteractionVi) SnapshotName() string { return "activity.InteractionV1" }
