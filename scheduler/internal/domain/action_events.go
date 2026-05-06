package domain

import "time"

type ActionAdded struct {
	ID                  string
	SchedulerID         string
	NaturalLanguageTask string
	ExecutionTime       time.Time
	Status              string
	CreatedAt           time.Time
}

type ActionRemoved struct {
	SchedulerID string
	ActionID    string
}

type ActionUpdated struct {
	ActionID     string
	SchedulerID  string
	Status       string
	ExecutedAt   *time.Time
	Result       string
	ErrorMessage string
}
