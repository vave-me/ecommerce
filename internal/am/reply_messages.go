package am

import (
	"context"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"middleman/internal/ddd"
	"middleman/internal/registry"
	"time"
)

const (
	FailureReply = "am.Failure"
	SuccessReply = "am.Success"

	OutcomeSuccess = "SUCCESS"
	OutcomeFailure = "FAILURE"

	ReplyHdrPrefix  = "REPLY_"
	ReplyNameHdr    = ReplyHdrPrefix + "NAME"
	ReplyOutcomeHdr = ReplyHdrPrefix + "OUTCOME"
)

type (
	ReplyMessage interface {
		MessageBase
		ddd.Reply
	}

	IncomingReplyMessage interface {
		IncomingMessageBase
		ddd.Reply
	}

	ReplyPublisher interface {
		Publish(ctx context.Context, topicName string, reply ddd.Reply) error
	}

	// ReplyPublisher  = MessagePublisher[ddd.Reply]
	// ReplySubscriber = MessageSubscriber[IncomingReplyMessage]
	// ReplyStream     = MessageStream[ddd.Reply, IncomingReplyMessage]

	replyPublisher struct {
		reg       registry.Registry
		publisher MessagePublisher
		logger    zerolog.Logger
	}

	replyMessage struct {
		id         string
		name       string
		payload    ddd.ReplyPayload
		occurredAt time.Time
		msg        IncomingMessageBase
	}
)

var _ ReplyMessage = (*replyMessage)(nil)

var _ ReplyPublisher = (*replyPublisher)(nil)

func NewReplyPublisher(reg registry.Registry, msgPublisher MessagePublisher, logger zerolog.Logger, mws ...MessagePublisherMiddleware) ReplyPublisher {

	return &replyPublisher{
		reg:       reg,
		publisher: MessagePublisherWithMiddleware(msgPublisher, mws...),
	}
}

func (s replyPublisher) Publish(ctx context.Context, topicName string, reply ddd.Reply) error {

	var err error
	var payload []byte

	if reply.ReplyName() != SuccessReply && reply.ReplyName() != FailureReply {
		payload, err = s.reg.Serialize(reply.ReplyName(), reply.Payload())
		if err != nil {
			s.logger.Error().Err(err).Msg("Error serializing reply payload")
			return err
		}
	}

	data, err := proto.Marshal(&ReplyMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(reply.OccurredAt()),
	})
	if err != nil {
		return err
	}
	if err != nil {
		s.logger.Error().Err(err).Msg("Error marshalling reply data")
		return err
	}

	return s.publisher.Publish(ctx, topicName, message{
		id:       reply.ID(),
		name:     reply.ReplyName(),
		subject:  topicName,
		data:     data,
		metadata: reply.Metadata(),
		sentAt:   time.Now(),
	})
	if err != nil {
		return err
	}

	return nil

}

func (r replyMessage) ID() string                { return r.id }
func (r replyMessage) ReplyName() string         { return r.name }
func (r replyMessage) Payload() ddd.ReplyPayload { return r.payload }
func (r replyMessage) Metadata() ddd.Metadata    { return r.msg.Metadata() }
func (r replyMessage) OccurredAt() time.Time     { return r.occurredAt }
func (r replyMessage) Subject() string           { return r.msg.Subject() }
func (r replyMessage) MessageName() string       { return r.msg.MessageName() }
func (r replyMessage) SentAt() time.Time         { return r.msg.SentAt() }
func (r replyMessage) ReceivedAt() time.Time     { return r.msg.ReceivedAt() }
func (r replyMessage) Ack() error                { return r.msg.Ack() }
func (r replyMessage) NAck() error               { return r.msg.NAck() }
func (r replyMessage) Extend() error             { return r.msg.Extend() }
func (r replyMessage) Kill() error               { return r.msg.Kill() }

type replyMsgHandler struct {
	reg     registry.Registry
	handler ddd.ReplyHandler[ddd.Reply]
	logger  zerolog.Logger
}

func NewReplyHandler(reg registry.Registry, handler ddd.ReplyHandler[ddd.Reply], logger zerolog.Logger, mws ...MessageHandlerMiddleware) MessageHandler {

	return MessageHandlerWithMiddleware(replyMsgHandler{
		reg:     reg,
		handler: handler,
	}, mws...)
}

func (h replyMsgHandler) HandleMessage(ctx context.Context, msg IncomingMessage) error {

	var replyData ReplyMessageData

	err := proto.Unmarshal(msg.Data(), &replyData)
	if err != nil {

		return err
	}

	replyName := msg.MessageName()

	var payload any

	if replyName != SuccessReply && replyName != FailureReply {
		payload, err = h.reg.Deserialize(replyName, replyData.GetPayload())
		if err != nil {

			return err
		}
	}

	replyMsg := replyMessage{
		id:         msg.ID(),
		name:       replyName,
		payload:    payload,
		occurredAt: replyData.GetOccurredAt().AsTime(),
		msg:        msg,
	}

	return h.handler.HandleReply(ctx, replyMsg)
}
