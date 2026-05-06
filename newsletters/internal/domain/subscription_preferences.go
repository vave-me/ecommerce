package domain

// Newsletter frequency options
type NewsletterFrequency string

const (
	UnknownFrequency NewsletterFrequency = ""
	Daily            NewsletterFrequency = "daily"
	Weekly           NewsletterFrequency = "weekly"
	Monthly          NewsletterFrequency = "monthly"
)

func (f NewsletterFrequency) String() string {
	switch f {
	case Daily, Weekly, Monthly:
		return string(f)
	default:
		return ""
	}
}

func ToNewsletterFrequency(frequency string) NewsletterFrequency {
	switch frequency {
	case Daily.String():
		return Daily
	case Weekly.String():
		return Weekly
	case Monthly.String():
		return Monthly
	default:
		return UnknownFrequency
	}
}

func ToNewsletterFrequencyPtr(frequency string) *NewsletterFrequency {
	if frequency == "" {
		return nil
	}
	f := ToNewsletterFrequency(frequency)
	return &f
}

// Subscription status
type SubscriptionStatus string

const (
	ActiveStatus       SubscriptionStatus = "active"
	PausedStatus       SubscriptionStatus = "paused"
	UnsubscribedStatus SubscriptionStatus = "unsubscribed"
)

func (s SubscriptionStatus) String() string {
	switch s {
	case ActiveStatus, PausedStatus, UnsubscribedStatus:
		return string(s)
	default:
		return ""
	}
}

func ToSubscriptionStatus(status string) SubscriptionStatus {
	switch status {
	case ActiveStatus.String():
		return ActiveStatus
	case PausedStatus.String():
		return PausedStatus
	case UnsubscribedStatus.String():
		return UnsubscribedStatus
	default:
		return ActiveStatus
	}
}

// Content format preferences
type ContentFormat string

const (
	HTMLFormat ContentFormat = "html"
	TextFormat ContentFormat = "text"
	BothFormat ContentFormat = "both"
)

func (f ContentFormat) String() string {
	switch f {
	case HTMLFormat, TextFormat, BothFormat:
		return string(f)
	default:
		return "both"
	}
}

func ToContentFormat(format string) ContentFormat {
	switch format {
	case HTMLFormat.String():
		return HTMLFormat
	case TextFormat.String():
		return TextFormat
	case BothFormat.String():
		return BothFormat
	default:
		return BothFormat
	}
}

// Edition status
type EditionStatus string

const (
	DraftStatus     EditionStatus = "draft"
	ScheduledStatus EditionStatus = "scheduled"
	SendingStatus   EditionStatus = "sending"
	SentStatus      EditionStatus = "sent"
)

func (s EditionStatus) String() string {
	switch s {
	case DraftStatus, ScheduledStatus, SendingStatus, SentStatus:
		return string(s)
	default:
		return ""
	}
}

func ToEditionStatus(status string) EditionStatus {
	switch status {
	case DraftStatus.String():
		return DraftStatus
	case ScheduledStatus.String():
		return ScheduledStatus
	case SendingStatus.String():
		return SendingStatus
	case SentStatus.String():
		return SentStatus
	default:
		return DraftStatus
	}
}

// Newsletter categories
type NewsletterCategory string

const (
	NewsCategory        NewsletterCategory = "news"
	ProductsCategory    NewsletterCategory = "products"
	PromotionalCategory NewsletterCategory = "promotional"
	EducationalCategory NewsletterCategory = "educational"
	GeneralCategory     NewsletterCategory = "general"
)

func (c NewsletterCategory) String() string {
	switch c {
	case NewsCategory, ProductsCategory, PromotionalCategory, EducationalCategory, GeneralCategory:
		return string(c)
	default:
		return "general"
	}
}

func ToNewsletterCategory(category string) NewsletterCategory {
	switch category {
	case NewsCategory.String():
		return NewsCategory
	case ProductsCategory.String():
		return ProductsCategory
	case PromotionalCategory.String():
		return PromotionalCategory
	case EducationalCategory.String():
		return EducationalCategory
	case GeneralCategory.String():
		return GeneralCategory
	default:
		return GeneralCategory
	}
}