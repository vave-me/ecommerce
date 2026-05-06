package domain

type CategoryCondition string

const (
	NewCondition         CategoryCondition = "new"
	UsedCondition        CategoryCondition = "used"
	BrokenCondition      CategoryCondition = "broken"
	RefurbishedCondition CategoryCondition = "refurbished"
	UnknownCondition     CategoryCondition = ""
)

func (c CategoryCondition) String() string {
	switch c {
	case NewCondition, UsedCondition, BrokenCondition, RefurbishedCondition:
		return string(c)
	default:
		return ""
	}
}

func ToCategoryCondition(s string) CategoryCondition {
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
