package domain

import "context"

type CommentRepository interface {
	Load(ctx context.Context, commentID string) (*Comment, error)
	Save(ctx context.Context, comment *Comment) error
}
