package ddd

import (
	"time"

	"github.com/google/uuid"
)

type (
	WebsocketOption interface {
		configureWebsocket(*websocket)
	}

	WebsocketPayload interface{}

	Websocket interface {
		IDer
		WebsocketName() string
		Payload() WebsocketPayload
		Metadata() Metadata
		OccurredAt() time.Time
	}

	websocket struct {
		Entity
		payload    WebsocketPayload
		metadata   Metadata
		occurredAt time.Time
	}
)

var _ Websocket = (*websocket)(nil)

func NewWebsocket(name string, payload WebsocketPayload, options ...WebsocketOption) Websocket {
	return newWebsocket(name, payload, options...)
}

func newWebsocket(name string, payload WebsocketPayload, options ...WebsocketOption) websocket {
	evt := websocket{
		Entity:     NewEntity(uuid.New().String(), name),
		payload:    payload,
		metadata:   make(Metadata),
		occurredAt: time.Now(),
	}

	for _, option := range options {
		option.configureWebsocket(&evt)
	}

	return evt
}

func (e websocket) WebsocketName() string     { return e.EntityName() }
func (e websocket) Payload() WebsocketPayload { return e.payload }
func (e websocket) Metadata() Metadata        { return e.metadata }
func (e websocket) OccurredAt() time.Time     { return e.occurredAt }
