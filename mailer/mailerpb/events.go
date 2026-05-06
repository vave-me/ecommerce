package mailerpb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	EmailAggregateChannel = "middleman.mailer.events.Email"
	EmailCreatedEvent     = "mailerapi.EmailCreated"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Store events
	if err := serde.Register(&EmailCreated{}); err != nil {
		return err
	}

	return nil
}

func (*EmailCreated) Key() string { return EmailCreatedEvent }
