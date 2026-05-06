package commands

import "time"

type CancelNewsletterCommand struct {
	NewsletterID string // ID of the newsletter to cancel
	SellerID     string // ID of the seller canceling the newsletter
	CanceledAt   time.Time
}
