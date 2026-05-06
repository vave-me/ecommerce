package domain

const (
	EmailCreatedEvent = "mailer.EmailCreated"
)

type EmailCreated struct {
	EmailID   string
	SenderID  string
	Recipient string
	Subject   string
	Body      string
	Status    EmailStatus
}

func (EmailCreated) Key() string { return EmailCreatedEvent }
