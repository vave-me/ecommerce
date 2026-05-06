package domain

import "time"

type UserV1 struct {
	Email                string
	Password             string
	Username             string
	FirstName            string
	LastName             string
	Location             string
	Enabled              bool
	Latitude             float64
	Longitude            float64
	GoogleID             string
	Thumbnail            string
	ResetToken           string
	ResetTokenExp        time.Time
	VerificationToken    string
	VerificationTokenExp time.Time
	Language             string
	Role                 UserRole
	Bio                  string
	Privacy              string
	Background           string
	TOTPEnabled          bool
}

func (UserV1) SnapshotName() string {
	return "users.UserV1"
}
