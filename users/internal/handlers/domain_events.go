package handlers

import (
	"context"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/users/internal/domain"
	"middleman/users/userspb"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandlers[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.UserCreatedEvent,
		domain.UserUpdatedEvent,
		domain.UserEnabledEvent,
		domain.UserDisabledEvent,
		domain.UserRenamedEvent,
		domain.UserLoggedInEvent,
		domain.UserLoggedOutEvent,
		domain.UserAuthorizedEvent,
		domain.UserPasswordForgotEvent,
		domain.UserPasswordResetEvent,
		domain.UserPasswordResetRequestedEvent,
		domain.UserTokenInvalidatedEvent,
		domain.UserTokenRefreshedEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {

	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling domain event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled domain event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling domain event", trace.WithAttributes(
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
	case domain.UserAuthorizedEvent:
		return h.onUserAuthorized(ctx, event)
	case domain.UserPasswordResetEvent:
		return h.onUserPasswordReset(ctx, event)
	case domain.UserPasswordResetRequestedEvent:
		return h.onUserPasswordReset(ctx, event)
	case domain.UserPasswordForgotEvent:
		return h.onUserPasswordForgotten(ctx, event)
	case domain.UserTokenInvalidatedEvent:
		return h.onUserTokenInvalidated(ctx, event)
	case domain.UserTokenRefreshedEvent:
		return h.onUserTokenRefreshed(ctx, event)
	}
	return nil
}

// Implement missing event handling methods

func (h domainHandlers[T]) onUserLoggedIn(ctx context.Context, event ddd.Event) error {
	user := event.Payload().(*domain.User)
	return h.publisher.Publish(ctx, userspb.UserAggregateChannel,
		ddd.NewEvent(userspb.UserLoggedInEvent, &userspb.UserLoggedIn{
			Id: user.ID(),
		}),
	)
}

func (h domainHandlers[T]) onUserLoggedOut(ctx context.Context, event ddd.Event) error {
	user := event.Payload().(*domain.User)
	return h.publisher.Publish(ctx, userspb.UserAggregateChannel,
		ddd.NewEvent(userspb.UserLoggedOutEvent, &userspb.UserLoggedOut{
			Id: user.ID(),
		}),
	)
}

func (h domainHandlers[T]) onUserAuthorized(ctx context.Context, event ddd.Event) error {
	user := event.Payload().(*domain.User)
	return h.publisher.Publish(ctx, userspb.UserAggregateChannel,
		ddd.NewEvent(userspb.UserAuthorizedEvent, &userspb.UserAuthorized{
			Id: user.ID(),
		}),
	)
}

func (h domainHandlers[T]) onUserCreated(ctx context.Context, event ddd.Event) error {
	user := event.Payload().(*domain.User)
	return h.publisher.Publish(ctx, userspb.UserAggregateChannel,
		ddd.NewEvent(userspb.UserCreatedEvent, &userspb.UserCreated{
			Id:                user.ID(),
			Email:             user.Email,
			Password:          user.Password,
			UserName:          user.Username,
			GoogleId:          user.GoogleID,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Lat:               float32(user.Lat),
			Lng:               float32(user.Lang),
			Thumbnail:         user.Thumbnail,
			VerificationToken: user.VerificationToken,
			Role:              string(user.Role),
		}),
	)
}
func (h domainHandlers[T]) onUserUpdated(ctx context.Context, event ddd.Event) error {
	user := event.Payload().(*domain.User)
	return h.publisher.Publish(ctx, userspb.UserAggregateChannel,
		ddd.NewEvent(userspb.UserUpdatedEvent, &userspb.UserUpdated{
			Id:         user.ID(),
			UserName:   user.Username,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Lat:        float32(user.Lat),
			Lng:        float32(user.Lang),
			Thumbnail:  user.Thumbnail,
			Bio:        user.Bio,
			Background: user.Background,
			Privacy:    user.Privacy,
			Role:       string(user.Role),
		}),
	)
}

func (h domainHandlers[T]) onUserRenamed(ctx context.Context, event ddd.Event) error {
	user := event.Payload().(*domain.User)
	return h.publisher.Publish(ctx, userspb.UserAggregateChannel,
		ddd.NewEvent(userspb.UserRenamedEvent, &userspb.UserRenamed{
			Id:   user.ID(),
			Name: user.Username,
		}),
	)
}
func (h domainHandlers[T]) onUserEnabled(ctx context.Context, event ddd.Event) error {
	user := event.Payload().(*domain.User)
	return h.publisher.Publish(ctx, userspb.UserAggregateChannel,
		ddd.NewEvent(userspb.UserEnabledToggledEvent, &userspb.UserEnabledToggled{
			Id:      user.ID(),
			Enabled: true,
		}),
	)
}

func (h domainHandlers[T]) onUserDisabled(ctx context.Context, event ddd.Event) error {
	user := event.Payload().(*domain.User)
	return h.publisher.Publish(ctx, userspb.UserAggregateChannel,
		ddd.NewEvent(userspb.UserEnabledToggledEvent, &userspb.UserEnabledToggled{
			Id:      user.ID(),
			Enabled: false,
		}),
	)
}

func (h domainHandlers[T]) onUserPasswordReset(ctx context.Context, event ddd.Event) error {
	user := event.Payload().(*domain.User)
	return h.publisher.Publish(ctx, userspb.UserAggregateChannel,
		ddd.NewEvent(userspb.UserPasswordResetRequestedEvent, &userspb.UserPasswordResetRequested{
			Email: user.Email,
		}),
	)
}

func (h domainHandlers[T]) onUserPasswordForgotten(ctx context.Context, event ddd.Event) error {
	//user := event.Payload().(*domain.User)
	return h.publisher.Publish(ctx, userspb.UserAggregateChannel,
		ddd.NewEvent(userspb.UserPasswordForgottenEvent, &userspb.UserPasswordForgotten{}),
	)
}

// Handle token invalidation events
func (h domainHandlers[T]) onUserTokenInvalidated(ctx context.Context, event ddd.Event) error {
	tokenEvent := event.Payload().(*domain.UserTokenInvalidated)

	// Convert domain event to integration event
	invalidatedAt := tokenEvent.InvalidatedAt

	// Create protobuf timestamp
	pbTimestamp := &timestamppb.Timestamp{
		Seconds: invalidatedAt.Unix(),
		Nanos:   int32(invalidatedAt.Nanosecond()),
	}

	return h.publisher.Publish(ctx, userspb.UserAggregateChannel,
		ddd.NewEvent(userspb.UserTokenInvalidatedEvent, &userspb.UserTokenInvalidated{
			UserId:        tokenEvent.UserID,
			TokenId:       tokenEvent.TokenID,
			InvalidatedAt: pbTimestamp,
			Reason:        tokenEvent.Reason,
		}),
	)
}

// Handle token refresh events
func (h domainHandlers[T]) onUserTokenRefreshed(ctx context.Context, event ddd.Event) error {
	tokenEvent := event.Payload().(*domain.UserTokenRefreshed)

	// Convert domain event to integration event
	refreshedAt := tokenEvent.RefreshedAt

	// Create protobuf timestamp
	pbTimestamp := &timestamppb.Timestamp{
		Seconds: refreshedAt.Unix(),
		Nanos:   int32(refreshedAt.Nanosecond()),
	}

	return h.publisher.Publish(ctx, userspb.UserAggregateChannel,
		ddd.NewEvent(userspb.UserTokenRefreshedEvent, &userspb.UserTokenRefreshed{
			UserId:      tokenEvent.UserID,
			OldTokenId:  tokenEvent.OldTokenID,
			NewTokenId:  tokenEvent.NewTokenID,
			RefreshedAt: pbTimestamp,
		}),
	)
}
