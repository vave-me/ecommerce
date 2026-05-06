package models

type TypeOfPost string

const (
	TypeOfPostPost        TypeOfPost = "post"
	TypeOfPostArticle     TypeOfPost = "article"
	TypeOfPostShort       TypeOfPost = "short"
	TypeOfPostSponsored   TypeOfPost = "sponsored"
	TypeOfPostAdvertise   TypeOfPost = "advertise"
	TypeOfPostAiGenerated TypeOfPost = "ai-generated"
	TypeOfPostUnknown     TypeOfPost = ""
)

func (s TypeOfPost) String() string {
	switch s {
	case TypeOfPostPost, TypeOfPostArticle, TypeOfPostShort, TypeOfPostSponsored, TypeOfPostAdvertise, TypeOfPostAiGenerated:
		return string(s)
	default:
		return ""
	}
}

func ToTypeOfPost(s string) TypeOfPost {

	switch s {
	case TypeOfPostPost.String():
		return TypeOfPostPost
	case TypeOfPostArticle.String():
		return TypeOfPostArticle
	case TypeOfPostShort.String():
		return TypeOfPostShort
	case TypeOfPostSponsored.String():
		return TypeOfPostSponsored
	case TypeOfPostAdvertise.String():
		return TypeOfPostAdvertise
	case TypeOfPostAiGenerated.String():
		return TypeOfPostAiGenerated
	default:
		return TypeOfPostUnknown
	}
}
