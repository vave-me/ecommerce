package models

import "encoding/json"

type Option struct {
	Name  string
	Value string
	Price float64
}
type Options []Option

// ToJSON returns a JSON string for the entire slice.
func (opts Options) ToJSON() (string, error) {
	b, err := json.Marshal(opts)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromJSON updates the slice from a JSON array string.
func (opts *Options) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), opts)
}
