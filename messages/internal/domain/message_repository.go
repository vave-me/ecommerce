package domain

import "context"

type MessageRepository interface {
	Load(ctx context.Context, id string) (*Message, error)
	Save(ctx context.Context, message *Message) error
}
