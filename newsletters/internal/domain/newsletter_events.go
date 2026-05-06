package domain

const (
	// Newsletter events
	NewsletterCreatedEvent     = "newsletters.NewsletterCreated"
	NewsletterUpdatedEvent     = "newsletters.NewsletterUpdated"
	NewsletterActivatedEvent   = "newsletters.NewsletterActivated"
	NewsletterDeactivatedEvent = "newsletters.NewsletterDeactivated"
	NewsletterDeletedEvent     = "newsletters.NewsletterDeleted"

	// Subscription events
	SubscribedEvent            = "newsletters.Subscribed"
	UnsubscribedEvent          = "newsletters.Unsubscribed"
	PreferencesUpdatedEvent    = "newsletters.PreferencesUpdated"
	SubscriptionPausedEvent    = "newsletters.SubscriptionPaused"
	SubscriptionResumedEvent   = "newsletters.SubscriptionResumed"

	// Edition events
	EditionCreatedEvent   = "newsletters.EditionCreated"
	EditionUpdatedEvent   = "newsletters.EditionUpdated"
	EditionScheduledEvent = "newsletters.EditionScheduled"
	EditionSendingEvent   = "newsletters.EditionSending"
	EditionSentEvent      = "newsletters.EditionSent"

	// Template events
	TemplateCreatedEvent = "newsletters.TemplateCreated"
	TemplateUpdatedEvent = "newsletters.TemplateUpdated"
	TemplateDeletedEvent = "newsletters.TemplateDeleted"
)

// Newsletter Events
type NewsletterCreated struct {
	UserID      string
	Name        string
	Description string
	Frequency   string
	Category    string
	TemplateID  string
}

func (NewsletterCreated) Key() string { return NewsletterCreatedEvent }

type NewsletterUpdated struct {
	Name        string
	Description string
	Frequency   string
	Category    string
	TemplateID  string
}

func (NewsletterUpdated) Key() string { return NewsletterUpdatedEvent }

type NewsletterActivated struct{}

func (NewsletterActivated) Key() string { return NewsletterActivatedEvent }

type NewsletterDeactivated struct{}

func (NewsletterDeactivated) Key() string { return NewsletterDeactivatedEvent }

type NewsletterDeleted struct{}

func (NewsletterDeleted) Key() string { return NewsletterDeletedEvent }

// Subscription Events
type Subscribed struct {
	UserID            string
	NewsletterID      string
	FrequencyOverride string
	Topics            []string
	Format            string
}

func (Subscribed) Key() string { return SubscribedEvent }

type Unsubscribed struct {
	Reason string
}

func (Unsubscribed) Key() string { return UnsubscribedEvent }

type PreferencesUpdated struct {
	FrequencyOverride string
	Topics            []string
	Format            string
}

func (PreferencesUpdated) Key() string { return PreferencesUpdatedEvent }

type SubscriptionPaused struct{}

func (SubscriptionPaused) Key() string { return SubscriptionPausedEvent }

type SubscriptionResumed struct{}

func (SubscriptionResumed) Key() string { return SubscriptionResumedEvent }

// Edition Events
type EditionCreated struct {
	NewsletterID string
	Subject      string
	ContentHTML  string
	ContentText  string
	TemplateData map[string]string
	ScheduledAt  int64 // Unix timestamp
	CreatedBy    string
}

func (EditionCreated) Key() string { return EditionCreatedEvent }

type EditionUpdated struct {
	Subject      string
	ContentHTML  string
	ContentText  string
	TemplateData map[string]string
	ScheduledAt  int64 // Unix timestamp
}

func (EditionUpdated) Key() string { return EditionUpdatedEvent }

type EditionScheduled struct {
	ScheduledAt int64 // Unix timestamp
}

func (EditionScheduled) Key() string { return EditionScheduledEvent }

type EditionSending struct{}

func (EditionSending) Key() string { return EditionSendingEvent }

type EditionSent struct {
	EditionID      string
	RecipientCount int
}

func (EditionSent) Key() string { return EditionSentEvent }

// Template Events
type TemplateCreated struct {
	UserID       string
	Name         string
	Description  string
	HTMLTemplate string
	TextTemplate string
	Variables    map[string]string
	PreviewData  map[string]string
	IsPublic     bool
}

func (TemplateCreated) Key() string { return TemplateCreatedEvent }

type TemplateUpdated struct {
	Name         string
	Description  string
	HTMLTemplate string
	TextTemplate string
	Variables    map[string]string
	PreviewData  map[string]string
	IsPublic     bool
}

func (TemplateUpdated) Key() string { return TemplateUpdatedEvent }

type TemplateDeleted struct{}

func (TemplateDeleted) Key() string { return TemplateDeletedEvent }