package models

import "encoding/json"

type Attribute struct {
	Key   string
	Value string
}

type Attributes []Attribute

func (attr Attributes) ToJSON() (string, error) {
	b, err := json.Marshal(attr)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromJSON updates the slice from a JSON array string.
func (attrs *Attributes) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), attrs)
}
