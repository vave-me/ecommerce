package domain

// You may define an enum or type alias for filter type
type FilterType string

// Example FilterType constants
const (
	FilterTypeSelect   FilterType = "select"
	FilterTypeRange    FilterType = "range"
	FilterTypeCheckbox FilterType = "checkbox"
	FilterTypeUnknown  FilterType = ""
)

func (s FilterType) String() string {
	switch s {
	case FilterTypeSelect, FilterTypeRange, FilterTypeCheckbox:
		return string(s)
	default:
		return ""
	}
}

func ToFilterType(s string) FilterType {

	switch s {
	case FilterTypeSelect.String():
		return FilterTypeSelect
	case FilterTypeRange.String():
		return FilterTypeRange
	case FilterTypeCheckbox.String():
		return FilterTypeCheckbox
	default:
		return FilterTypeUnknown
	}
}
