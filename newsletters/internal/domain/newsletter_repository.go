package domain

import (
	"middleman/internal/es"
)

type NewsletterRepository interface {
	es.AggregateRepository[*Newsletter]
}

type SubscriptionRepository interface {
	es.AggregateRepository[*Subscription]
}

type EditionRepository interface {
	es.AggregateRepository[*Edition]
}

type TemplateRepository interface {
	es.AggregateRepository[*Template]
}
