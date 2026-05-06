package userspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	UserAggregateChannel = "middleman.users.events.User"

	UserAuthorizedEvent             = "usersapi.UserAuthorized"
	UserCreatedEvent                = "usersapi.UserCreated"
	UserUpdatedEvent                = "usersapi.UserUpdated"
	UserEnabledToggledEvent         = "usersapi.UserEnabledToggled"
	UserRenamedEvent                = "usersapi.UserRenamed"
	UserLoggedInEvent               = "usersapi.UserLoggedIn"
	UserLoggedOutEvent              = "usersapi.UserLoggedOut"
	UserPasswordResetEvent          = "usersapi.UserPasswordReset"
	UserPasswordResetRequestedEvent = "usersapi.UserPasswordResetRequested"
	UserPasswordForgottenEvent      = "usersapi.UserPasswordForgotten"
	UserTokenInvalidatedEvent       = "usersapi.UserTokenInvalidated"
	UserTokenRefreshedEvent         = "usersapi.UserTokenRefreshed"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {

	//user events
	if err := serde.Register(&UserCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&UserUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&UserAuthorized{}); err != nil {
		return err
	}
	if err := serde.Register(&UserLoggedIn{}); err != nil {
		return err
	}
	if err := serde.Register(&UserLoggedOut{}); err != nil {
		return err
	}
	if err := serde.Register(&UserPasswordReset{}); err != nil {
		return err
	}
	if err := serde.Register(&UserPasswordForgotten{}); err != nil {
		return err
	}
	if err := serde.Register(&UserPasswordResetRequested{}); err != nil {
		return err
	}
	if err := serde.Register(&UserTokenInvalidated{}); err != nil {
		return err
	}
	if err := serde.Register(&UserTokenRefreshed{}); err != nil {
		return err
	}
	if err := serde.Register(&UserEnabledToggled{}); err != nil {
		return err
	}
	if err := serde.Register(&UserRenamed{}); err != nil {
		return err
	}

	return nil
}
func (*UserAuthorized) Key() string             { return UserAuthorizedEvent }
func (*UserCreated) Key() string                { return UserCreatedEvent }
func (*UserUpdated) Key() string                { return UserUpdatedEvent }
func (*UserEnabledToggled) Key() string         { return UserEnabledToggledEvent }
func (*UserRenamed) Key() string                { return UserRenamedEvent }
func (*UserLoggedIn) Key() string               { return UserLoggedInEvent }
func (*UserLoggedOut) Key() string              { return UserLoggedOutEvent }
func (*UserPasswordReset) Key() string          { return UserPasswordResetEvent }
func (*UserPasswordResetRequested) Key() string { return UserPasswordResetRequestedEvent }
func (*UserPasswordForgotten) Key() string      { return UserPasswordForgottenEvent }
func (*UserTokenInvalidated) Key() string       { return UserTokenInvalidatedEvent }
func (*UserTokenRefreshed) Key() string         { return UserTokenRefreshedEvent }
