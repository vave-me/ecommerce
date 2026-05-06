package am

import (
	"context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"middleman/internal/ddd"
	"middleman/internal/registry"
	"time"

	"github.com/rs/zerolog"
)

type (
	WebsocketMessage interface {
		MessageBase
		ddd.Websocket
	}

	IncomingWebsocketMessage interface {
		IncomingMessageBase
		ddd.Websocket
	}

	WebsocketPublisher interface {
		Publish(ctx context.Context, topicName string, websocket ddd.Websocket) error
	}

	websocketPublisher struct {
		reg       registry.Registry
		publisher MessagePublisher
		logger    zerolog.Logger
	}

	websocketMessage struct {
		id         string
		name       string
		payload    ddd.WebsocketPayload
		occurredAt time.Time
		msg        IncomingMessageBase
	}
)

var _ WebsocketMessage = (*websocketMessage)(nil)
var _ WebsocketPublisher = (*websocketPublisher)(nil)

func NewWebsocketPublisher(reg registry.Registry, msgPublisher MessagePublisher, logger zerolog.Logger, mws ...MessagePublisherMiddleware) WebsocketPublisher {
	logger.Info().Msg("Creating a new websocket publisher")
	return websocketPublisher{
		reg:       reg,
		publisher: MessagePublisherWithMiddleware(msgPublisher, mws...),
		logger:    logger,
	}
}

func (s websocketPublisher) Publish(ctx context.Context, topicName string, websocket ddd.Websocket) error {
	s.logger.Info().Str("websocketID", websocket.ID()).Str("websocketName", websocket.WebsocketName()).Msg("Publishing websocket")

	payload, err := s.reg.Serialize(websocket.WebsocketName(), websocket.Payload())
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to serialize websocket payload")
		return err
	}

	data, err := proto.Marshal(&WebsocketMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(websocket.OccurredAt()),
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to marshal websocket data")
		return err
	}

	err = s.publisher.Publish(ctx, topicName, message{
		id:       websocket.ID(),
		name:     websocket.WebsocketName(),
		subject:  topicName,
		data:     data,
		metadata: websocket.Metadata(),
		sentAt:   time.Now(),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("topic", topicName).Msg("Failed to publish message")
		return err
	}

	s.logger.Info().Str("websocketID", websocket.ID()).Str("topic", topicName).Msg("Websocket published successfully")
	return nil
}

func (e websocketMessage) ID() string                    { return e.id }
func (e websocketMessage) WebsocketName() string         { return e.name }
func (e websocketMessage) Payload() ddd.WebsocketPayload { return e.payload }
func (e websocketMessage) Metadata() ddd.Metadata        { return e.msg.Metadata() }
func (e websocketMessage) OccurredAt() time.Time         { return e.occurredAt }
func (e websocketMessage) Subject() string               { return e.msg.Subject() }
func (e websocketMessage) MessageName() string           { return e.msg.MessageName() }
func (e websocketMessage) SentAt() time.Time             { return e.msg.SentAt() }
func (e websocketMessage) ReceivedAt() time.Time         { return e.msg.ReceivedAt() }
func (e websocketMessage) Ack() error                    { return e.msg.Ack() }
func (e websocketMessage) NAck() error                   { return e.msg.NAck() }
func (e websocketMessage) Extend() error                 { return e.msg.Extend() }
func (e websocketMessage) Kill() error                   { return e.msg.Kill() }

type websocketMsgHandler struct {
	reg     registry.Registry
	handler ddd.WebsocketHandler[ddd.Websocket]
	logger  zerolog.Logger
}

func NewWebsocketHandler(reg registry.Registry, handler ddd.WebsocketHandler[ddd.Websocket], logger zerolog.Logger, mws ...MessageHandlerMiddleware) MessageHandler {
	logger.Info().Msg("Creating a new websocket handler")
	return MessageHandlerWithMiddleware(websocketMsgHandler{
		reg:     reg,
		handler: handler,
		logger:  logger,
	}, mws...)
}
func (h websocketMsgHandler) HandleMessage(ctx context.Context, msg IncomingMessage) error {
	var websocketData WebsocketMessageData

	err := proto.Unmarshal(msg.Data(), &websocketData)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to unmarshal WebsocketMessage")
		return err
	}
	websocketName := msg.MessageName()
	payload, err := h.reg.Deserialize(websocketName, websocketData.GetPayload())
	websocketMsg := websocketMessage{
		id:         msg.ID(),
		name:       websocketName,
		payload:    payload,
		occurredAt: websocketData.GetOccurredAt().AsTime(),
		msg:        msg,
	}

	return h.handler.HandleWebsocket(ctx, websocketMsg)
}
