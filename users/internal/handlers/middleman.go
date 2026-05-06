package handlers

import (
	"context"
	"log"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"middleman/users/internal/domain"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type middlemanHandlers[T ddd.Event] struct {
	middleman domain.MiddlemanRepository
}

var _ ddd.EventHandler[ddd.Event] = (*middlemanHandlers[ddd.Event])(nil)

func NewMiddlemanHandlers(middleman domain.MiddlemanRepository) ddd.EventHandler[ddd.Event] {
	return middlemanHandlers[ddd.Event]{
		middleman: middleman,
	}
}
func RegisterMiddlemanHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.UserCreatedEvent,
		domain.UserUpdatedEvent,
		domain.UserDisabledEvent,
		domain.UserEnabledEvent,
		domain.UserRenamedEvent,
		domain.UserLoggedInEvent,
		domain.UserLoggedOutEvent,
	)
}
func RegisterMiddlemanHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		middlemanHandlers := di.Get(ctx, "middlemanHandlers").(ddd.EventHandler[ddd.Event])

		return middlemanHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get("domainDispatcher").(*ddd.EventDispatcher[ddd.Event])

	RegisterMiddlemanHandlers(subscriber, handlers)
}
func (h middlemanHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {

	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling mall event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled mall event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling mall event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))
	switch event.EventName() {
	case domain.UserCreatedEvent:
		return h.onUserCreated(ctx, event)
	case domain.UserUpdatedEvent:
		return h.onUserUpdated(ctx, event)
	case domain.UserEnabledEvent:
		return h.onUserEnabled(ctx, event)
	case domain.UserDisabledEvent:
		return h.onUserDisabled(ctx, event)
	case domain.UserRenamedEvent:
		return h.onUserRenamed(ctx, event)
	case domain.UserLoggedInEvent:
		return h.onUserLoggedIn(ctx, event)
	case domain.UserLoggedOutEvent:
		return h.onUserLoggedOut(ctx, event)
	}
	return nil
}

func (h middlemanHandlers[T]) onUserCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.User)
	log.Printf("Adding user to database: ID=%s, Lat=%f, Lng=%f", payload.ID(), payload.Lat, payload.Lang)
	return h.middleman.AddUser(
		ctx,
		payload.ID(),
		payload.Email,
		payload.Username,
		payload.FirstName,
		payload.LastName,
		payload.GoogleID,
		payload.Enabled,
		payload.Lat,
		payload.Lang,
		payload.Thumbnail,
		string(payload.Role),
	)
}
func (h middlemanHandlers[T]) onUserUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.User)
	log.Printf("Update user to database: ID=%s, Lat=%f, Lng=%f", payload.ID(), payload.Lat, payload.Lang)
	return h.middleman.UpdateUser(
		ctx,
		payload.ID(),
		payload.Username,
		payload.FirstName,
		payload.LastName,
		payload.Bio,
		payload.Privacy,
		payload.Background,
		payload.Lat,
		payload.Lang,
		payload.Thumbnail,
		string(payload.Role),
	)
}

func (h middlemanHandlers[T]) onUserEnabled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.User)
	return h.middleman.EnableUser(ctx, payload.ID(), true)
}

func (h middlemanHandlers[T]) onUserDisabled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.User)
	return h.middleman.EnableUser(ctx, payload.ID(), false)
}

func (h middlemanHandlers[T]) onUserRenamed(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.User)
	return h.middleman.RenameUser(ctx, payload.ID(), payload.FirstName)
}

func (h middlemanHandlers[T]) onUserLoggedIn(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.User)
	return h.middleman.LogUserIn(ctx, payload.ID())
}

func (h middlemanHandlers[T]) onUserLoggedOut(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.User)
	return h.middleman.LogUserOut(ctx, payload.ID())
}
