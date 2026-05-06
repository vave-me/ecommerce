package domain

type UserLanguage string

const (
	en      UserLanguage = "en"
	de      UserLanguage = "de"
	pl      UserLanguage = "pl"
	unknown UserLanguage = ""
)

func (s UserLanguage) String() string {
	switch s {
	case en, de, pl:
		return string(s)
	default:
		return ""
	}
}

func ToUserLanguage(status string) UserLanguage {
	switch status {
	case en.String():
		return en
	case de.String():
		return de
	case pl.String():
		return pl
	default:
		return unknown
	}
}
