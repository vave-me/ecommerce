package domain

import "context"

type EmailRepository interface {
	Load(ctx context.Context, id string) (*Email, error)
	Save(ctx context.Context, mailer *Email) error
}
