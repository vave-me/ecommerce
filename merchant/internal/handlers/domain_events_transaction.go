package handlers

import (
	"context"
	"database/sql"
	"middleman/internal/di"
	"middleman/internal/ddd"
	"middleman/merchant/internal/constants"
)

func RegisterDomainEventHandlersTx(container di.Container) {
	handlers := container.Get(constants.DomainEventHandlersKey).(ddd.EventHandler[ddd.Event])
	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterDomainEventHandlers(subscriber, ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		ctx = container.Scoped(ctx)
		
		err := handlers.HandleEvent(ctx, event)
		
		tx := di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		
		return tx.Commit()
	}))
}