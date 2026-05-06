package handlers

import (
	"context"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/payments/internal/constants"

	"github.com/rs/zerolog/log"
)

func RegisterDomainEventHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) (err error) {
		log.Debug().Str("event", event.EventName()).Msg("DomainEventsTx: start handling")

		domainHandlers := di.Get(ctx, constants.DomainEventHandlersKey).(ddd.EventHandler[ddd.Event])

		err = domainHandlers.HandleEvent(ctx, event)
		if err != nil {
			log.Error().Err(err).Str("event", event.EventName()).Msg("DomainEventsTx: handler returned error")
		} else {
			log.Debug().Str("event", event.EventName()).Msg("DomainEventsTx: handler succeeded")
		}
		return err
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterDomainEventHandlers(subscriber, handlers)
}
