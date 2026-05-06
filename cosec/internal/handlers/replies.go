package handlers

import (
	"middleman/cosec/internal"
	"middleman/cosec/internal/models"
	"middleman/internal/am"
	"middleman/internal/registry"
	"middleman/internal/sec"

	"github.com/rs/zerolog"
)

func NewReplyHandlers(reg registry.Registry, orchestrator sec.Orchestrator[*models.CheckoutData], logger zerolog.Logger, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewReplyHandler(reg, orchestrator, logger, mws...)
}

func RegisterReplyHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(internal.CheckoutSagaReplyChannel, handlers, am.GroupName("cosec-replies"))
	return err
}
