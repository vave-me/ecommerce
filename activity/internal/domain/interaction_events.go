package domain

type InteractionAdded struct {
	ActivityID string
	ItemID     string
	ItemType   string
	ActionType string
}

type InteractionRemoved struct {
	ActivityID string
	ItemID     string
	ItemType   string
	ActionType string
}

type InteractionUpdated struct {
	ActivityID string
	ItemID     string
	ItemType   string
	ActionType string
}
