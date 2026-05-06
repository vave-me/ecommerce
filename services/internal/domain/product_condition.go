package domain

type ServiceCondition string

const (
	NewCondition         ServiceCondition = "new"
	UsedCondition        ServiceCondition = "used"
	BrokenCondition      ServiceCondition = "broken"
	RefurbishedCondition ServiceCondition = "refurbished"
	UnknownCondition     ServiceCondition = ""
)

func (c ServiceCondition) String() string {
	switch c {
	case NewCondition, UsedCondition, BrokenCondition, RefurbishedCondition:
		return string(c)
	default:
		return ""
	}
}

func ToServiceCondition(s string) ServiceCondition {
	switch s {
	case NewCondition.String():
		return NewCondition
	case UsedCondition.String():
		return UsedCondition
	case BrokenCondition.String():
		return BrokenCondition
	case RefurbishedCondition.String():
		return RefurbishedCondition
	default:
		return UnknownCondition
	}
}
