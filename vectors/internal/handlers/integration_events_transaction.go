package handlers

import (
	"context"
	"middleman/internal/am"
	"middleman/internal/di"
	"middleman/vectors/internal/constants"
)

func RegisterIntegrationEventHandlersTx(container di.Container) (err error) {
	rawMsgHandler := am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) (err error) {
		ctx = container.Scoped(ctx)

		return di.Get(ctx, constants.IntegrationEventHandlersKey).(am.MessageHandler).HandleMessage(ctx, msg)
	})

	subscriber := container.Get(constants.MessageSubscriberKey).(am.MessageSubscriber)

	return RegisterVectorIntegrationEventHandlers(subscriber, rawMsgHandler)
}
