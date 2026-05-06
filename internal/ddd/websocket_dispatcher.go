package ddd

import (
	"context"
	"sync"
)

type (
	WebsocketHandler[T Websocket] interface {
		HandleWebsocket(ctx context.Context, websocket T) error
	}

	WebsocketHandlerFunc[T Websocket] func(ctx context.Context, websocket T) error

	WebsocketSubscriber[T Websocket] interface {
		Subscribe(handler WebsocketHandler[T], websockets ...string)
	}

	WebsocketPublisher[T Websocket] interface {
		Publish(ctx context.Context, websockets ...T) error
	}

	WebsocketDispatcher[T Websocket] struct {
		handlers []websocketHandler[T]
		mu       sync.Mutex
	}

	websocketHandler[T Websocket] struct {
		h       WebsocketHandler[T]
		filters map[string]struct{}
	}
)

var _ interface {
	WebsocketSubscriber[Websocket]
	WebsocketPublisher[Websocket]
} = (*WebsocketDispatcher[Websocket])(nil)

func NewWebsocketDispatcher[T Websocket]() *WebsocketDispatcher[T] {
	return &WebsocketDispatcher[T]{
		handlers: make([]websocketHandler[T], 0),
	}
}

func (h *WebsocketDispatcher[T]) Subscribe(handler WebsocketHandler[T], websockets ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var filters map[string]struct{}
	if len(websockets) > 0 {
		filters = make(map[string]struct{})
		for _, websocket := range websockets {
			filters[websocket] = struct{}{}
		}
	}

	h.handlers = append(h.handlers, websocketHandler[T]{
		h:       handler,
		filters: filters,
	})
}

func (h *WebsocketDispatcher[T]) Publish(ctx context.Context, websockets ...T) error {
	for _, websocket := range websockets {
		for _, handler := range h.handlers {
			if handler.filters != nil {
				if _, exists := handler.filters[websocket.WebsocketName()]; !exists {
					continue
				}
			}
			err := handler.h.HandleWebsocket(ctx, websocket)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (f WebsocketHandlerFunc[T]) HandleWebsocket(ctx context.Context, websocket T) error {
	return f(ctx, websocket)
}
