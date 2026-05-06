package models

type SuggestionItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"` // "product", "post", "vehicle", etc.
}
