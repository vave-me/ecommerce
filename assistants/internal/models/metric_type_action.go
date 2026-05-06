package models

type MetricTypeAction string

const (
	MetricTypeActionAdd     MetricTypeAction = "add"
	MetricTypeActionRemove  MetricTypeAction = "remove"
	MetricTypeActionToggle  MetricTypeAction = "toggle"
	MetricTypeActionUnknown MetricTypeAction = ""
)

func (s MetricTypeAction) String() string {
	switch s {
	case MetricTypeActionAdd, MetricTypeActionRemove, MetricTypeActionToggle:
		return string(s)
	default:
		return ""
	}
}

func ToMetricTypeAction(s string) MetricTypeAction {

	switch s {
	case MetricTypeActionAdd.String():
		return MetricTypeActionAdd
	case MetricTypeActionRemove.String():
		return MetricTypeActionRemove
	case MetricTypeActionToggle.String():
		return MetricTypeActionToggle
	default:
		return MetricTypeActionUnknown
	}
}
