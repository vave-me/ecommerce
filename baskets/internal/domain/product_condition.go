package domain

type ProductCondition string

const (
	NewCondition         ProductCondition = "new"
	UsedCondition        ProductCondition = "used"
	BrokenCondition      ProductCondition = "broken"
	RefurbishedCondition ProductCondition = "refurbished"
	UnknownCondition     ProductCondition = ""
)

func (c ProductCondition) String() string {
	switch c {
	case NewCondition, UsedCondition, BrokenCondition, RefurbishedCondition:
		return string(c)
	default:
		return ""
	}
}

func ToProductCondition(s string) ProductCondition {
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
