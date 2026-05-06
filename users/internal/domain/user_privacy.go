package domain

type UserPrivacy string

const (
	UserUnknown     UserPrivacy = ""
	UserIsPublic    UserPrivacy = "public"
	UserIsPrivate   UserPrivacy = "private"
	UserIsFollowers UserPrivacy = "followers"
)

func (s UserPrivacy) String() string {
	switch s {
	case UserIsPublic, UserIsPrivate, UserIsFollowers:
		return string(s)
	default:
		return ""
	}
}

func ToUserPrivacy(status string) UserPrivacy {
	switch status {
	case UserIsPublic.String():
		return UserIsPublic
	case UserIsPrivate.String():
		return UserIsPrivate
	case UserIsFollowers.String():
		return UserIsFollowers
	default:
		return UserUnknown
	}
}
