package domain

type EmailV1 struct {
	EmailID   string
	SenderID  string
	Recipient string
	Subject   string
	Body      string
	Status    EmailStatus
}

func (EmailV1) SnapshotName() string { return "mailer.EmailV1" }
