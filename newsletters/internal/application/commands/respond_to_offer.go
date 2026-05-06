package commands

import "time"

type RespondToNewsletterCommand struct {
	NewsletterID string // ID of the newsletter being responded to
	BuyerID      string // ID of the buyer responding
	Response     string // Accept or Reject
	RespondedAt  time.Time
}
