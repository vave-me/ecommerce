package jetstream

import (
	"github.com/rs/zerolog"
	"time"

	"middleman/internal/am"
	"middleman/internal/ddd"
)

type (
	rawMessage struct {
		id         string
		name       string
		subject    string
		data       []byte
		metadata   ddd.Metadata
		sentAt     time.Time
		receivedAt time.Time
		acked      bool
		ackFn      func() error
		nackFn     func() error
		extendFn   func() error
		killFn     func() error
		logger     zerolog.Logger
	}
)

var _ am.Message = (*rawMessage)(nil)

// Implement am.Message interface methods

func (m *rawMessage) ID() string {
	return m.id
}

func (m *rawMessage) Subject() string {
	return m.subject
}

func (m *rawMessage) MessageName() string {
	return m.name
}

func (m *rawMessage) Data() []byte {
	return m.data
}

func (m *rawMessage) Metadata() ddd.Metadata {
	return m.metadata
}

func (m *rawMessage) SentAt() time.Time {
	return m.sentAt
}

func (m *rawMessage) ReceivedAt() time.Time {
	return m.receivedAt
}

// Ack acknowledges the message
func (m *rawMessage) Ack() error {
	m.logger.With().
		Str("component", "rawMessage.Ack").
		Str("message_id", m.id).
		Logger()

	if m.acked {
		m.logger.Warn().
			Msg("Attempted to Ack an already acknowledged message")
		return nil
	}

	m.logger.Info().
		Msg("Acknowledging message")
	m.acked = true
	if err := m.ackFn(); err != nil {
		m.logger.Error().
			Err(err).
			Msg("Failed to Ack message")
		return err
	}

	m.logger.Info().
		Msg("Message acknowledged successfully")
	return nil
}

// NAck negatively acknowledges the message
func (m *rawMessage) NAck() error {
	m.logger.With().
		Str("component", "rawMessage.NAck").
		Str("message_id", m.id).
		Logger()

	if m.acked {
		m.logger.Warn().
			Msg("Attempted to NAck an already acknowledged message")
		return nil
	}

	m.logger.Info().
		Msg("Negatively acknowledging message")
	m.acked = true
	if err := m.nackFn(); err != nil {
		m.logger.Error().
			Err(err).
			Msg("Failed to NAck message")
		return err
	}

	m.logger.Info().
		Msg("Message negatively acknowledged successfully")
	return nil
}

// Extend extends the message processing time
func (m *rawMessage) Extend() error {
	m.logger.With().
		Str("component", "rawMessage.Extend").
		Str("message_id", m.id).
		Logger()

	m.logger.Info().
		Msg("Extending message processing time")

	if err := m.extendFn(); err != nil {
		m.logger.Error().
			Err(err).
			Msg("Failed to extend message processing time")
		return err
	}

	m.logger.Info().
		Msg("Message processing time extended successfully")
	return nil
}

// Kill terminates the message processing
func (m *rawMessage) Kill() error {
	m.logger.With().
		Str("component", "rawMessage.Kill").
		Str("message_id", m.id).
		Logger()

	if m.acked {
		m.logger.Warn().
			Msg("Attempted to Kill an already acknowledged message")
		return nil
	}

	m.logger.Info().
		Msg("Killing message processing")
	m.acked = true
	if err := m.killFn(); err != nil {
		m.logger.Error().
			Err(err).
			Msg("Failed to Kill message processing")
		return err
	}

	m.logger.Info().
		Msg("Message processing killed successfully")
	return nil
}
