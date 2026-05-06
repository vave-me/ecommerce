package domain

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusCreated   PostStatus = "created"
	PostStatusGhosted   PostStatus = "ghosted"
	PostStatusPaused    PostStatus = "paused"
	PostStatusArchived  PostStatus = "archived"
	PostStatusUnknown   PostStatus = ""
)

func (s PostStatus) String() string {
	switch s {
	case PostStatusPublished, PostStatusDraft, PostStatusArchived, PostStatusCreated, PostStatusGhosted, PostStatusPaused:
		return string(s)
	default:
		return ""
	}
}

func ToPostStatus(s string) PostStatus {
	switch s {
	case PostStatusCreated.String():
		return PostStatusCreated
	case PostStatusGhosted.String():
		return PostStatusGhosted
	case PostStatusPaused.String():
		return PostStatusPaused
	case PostStatusDraft.String():
		return PostStatusDraft
	case PostStatusArchived.String():
		return PostStatusArchived
	case PostStatusPublished.String():
		return PostStatusPublished
	default:
		return PostStatusUnknown
	}
}
