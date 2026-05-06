package models

type UserType string

const (
	UserTypePrivate   UserType = "private"
	UserTypeBusiness  UserType = "business"
	UserTypeAnonymous UserType = "anonymous"
	UserTypeUnknown   UserType = ""
)

func (s UserType) String() string {
	switch s {
	case UserTypePrivate, UserTypeBusiness, UserTypeAnonymous:
		return string(s)
	default:
		return ""
	}
}

func ToUserType(s string) UserType {

	switch s {
	case UserTypePrivate.String():
		return UserTypePrivate
	case UserTypeBusiness.String():
		return UserTypeBusiness
	case UserTypeAnonymous.String():
		return UserTypeAnonymous
	default:
		return UserTypeUnknown
	}
}
