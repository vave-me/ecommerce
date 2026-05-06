package handlers

import (
	"context"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/mailer/internal/application"
	"middleman/users/userspb"
	"time"

	"github.com/rs/zerolog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type integrationHandlers[T ddd.Event] struct {
	app   application.App
	users application.UserCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(reg registry.Registry, app application.App, users application.UserCacheRepository, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		app:   app,
		users: users,
	}, zerolog.Logger{}, mws...)
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	_, err = subscriber.Subscribe(userspb.UserAggregateChannel, handlers, am.MessageFilter{
		userspb.UserCreatedEvent,
		userspb.UserPasswordResetRequestedEvent,
	}, am.GroupName("notification-users"))
	if err != nil {
		return err
	}

	return err
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling integration event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled integration event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling integration event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case userspb.UserCreatedEvent:
		return h.onUserRegistered(ctx, event)
	case userspb.UserPasswordResetRequestedEvent:
		return h.onPasswordResetRequested(ctx, event)

	}

	return nil
}

func (h integrationHandlers[T]) onUserRegistered(ctx context.Context, event T) error {

	payload := event.Payload().(*userspb.UserCreated)
	//if payload.GetEnabled() == true {
	//	return nil
	//}
	cmd := application.UserCreated{
		UserID:    payload.GetId(),
		Email:     payload.GetEmail(),
		FirstName: payload.GetFirstName(),
		LastName:  payload.GetLastName(),
		Enabled:   payload.GetEnabled(),
	}
	return h.app.NotifyUserCreated(ctx, cmd)

}

func (h integrationHandlers[T]) onPasswordResetRequested(ctx context.Context, event T) error {

	payload := event.Payload().(*userspb.UserPasswordResetRequested)

	cmd := application.PasswordReset{
		Email:          payload.GetEmail(),
		ResetToken:     payload.GetToken(),
		ExpirationDate: application.ConvertProtoTimestamp(payload.GetExpirationTime()),
	}
	return h.app.ResetPassword(ctx, cmd)

}
