package commands

import "time"

type SendNewsletterToBuyerCommand struct {
	NewsletterID string // ID of the newsletter being sent
	BuyerID      string // ID of the buyer receiving the newsletter
	SellerID     string // ID of the seller sending the newsletter
	SentAt       time.Time
}
