package models

type Condition string

const (
	NewCondition         Condition = "new"
	UsedCondition        Condition = "used"
	BrokenCondition      Condition = "broken"
	RefurbishedCondition Condition = "refurbished"
	UnknownCondition     Condition = ""
)

func (c Condition) String() string {
	switch c {
	case NewCondition, UsedCondition, BrokenCondition, RefurbishedCondition:
		return string(c)
	default:
		return ""
	}
}

func ToCondition(s string) Condition {
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
