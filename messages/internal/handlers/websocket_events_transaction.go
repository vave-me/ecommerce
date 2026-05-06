package handlers

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog"
	"middleman/internal/am"
	"middleman/internal/di"
	"middleman/messages/internal/constants"
)

func RegisterWebsocketsCommandsHandlersTx(container di.Container, logger zerolog.Logger) error {
	rawMsgHandler := am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) (err error) {
		logger.Info().Str("messageID", msg.ID()).Msg("Handling incoming WebSocket message")

		ctx = container.Scoped(ctx)
		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
		logger.Info().Msgf("INFO ")
		return di.Get(ctx, constants.WebsocketsEventHandlersKey).(am.MessageHandler).HandleMessage(ctx, msg)
	})

	subscriber := container.Get(constants.WebSocketSubscriberKey).(am.MessageSubscriber)

	return RegisterWebsocketCommandHandlers(subscriber, rawMsgHandler, logger)
}
