package handlers

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/streams/internal/domain"
	"middleman/streams/streamspb"
)

type liveStreamingDomainHandlers[T ddd.Event] struct {
	publisher         am.EventPublisher
	webhookDispatcher *WebhookDispatcher
}

var _ ddd.EventHandler[ddd.Event] = (*liveStreamingDomainHandlers[ddd.Event])(nil)

func NewLiveStreamingDomainEventHandlers(publisher am.EventPublisher, webhookDispatcher *WebhookDispatcher) ddd.EventHandler[ddd.Event] {
	return &liveStreamingDomainHandlers[ddd.Event]{
		publisher:         publisher,
		webhookDispatcher: webhookDispatcher,
	}
}

func RegisterLiveStreamingDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.LiveStreamCreatedEvent,
		domain.StreamingConfiguredEvent,
		domain.LiveStreamStartedEvent,
		domain.LiveStreamStoppedEvent,
		domain.ViewerJoinedEvent,
		domain.ViewerLeftEvent,
		domain.StreamQualityChangedEvent,
		domain.StreamHealthUpdatedEvent,
		domain.CDNEndpointAddedEvent,
		domain.DRMConfiguredEvent,
	)
}

func (h liveStreamingDomainHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling domain event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled domain event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling live streaming domain event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	// First, dispatch webhooks for this event
	if h.webhookDispatcher != nil {
		go func() {
			// Use a new context to avoid blocking the main event handling
			webhookCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			
			if err := h.webhookDispatcher.HandleEvent(webhookCtx, event); err != nil {
				// Log error but don't fail the main event handling
				span.AddEvent("Failed to dispatch webhook", trace.WithAttributes(
					attribute.String("error", err.Error()),
				))
			}
		}()
	}

	// Then handle the specific event
	switch event.EventName() {
	case domain.LiveStreamCreatedEvent:
		return h.onLiveStreamCreated(ctx, event)
	case domain.StreamingConfiguredEvent:
		return h.onStreamingConfigured(ctx, event)
	case domain.LiveStreamStartedEvent:
		return h.onLiveStreamStarted(ctx, event)
	case domain.LiveStreamStoppedEvent:
		return h.onLiveStreamStopped(ctx, event)
	case domain.ViewerJoinedEvent:
		return h.onViewerJoined(ctx, event)
	case domain.ViewerLeftEvent:
		return h.onViewerLeft(ctx, event)
	case domain.StreamQualityChangedEvent:
		return h.onStreamQualityChanged(ctx, event)
	case domain.StreamHealthUpdatedEvent:
		return h.onStreamHealthUpdated(ctx, event)
	case domain.CDNEndpointAddedEvent:
		return h.onCDNEndpointAdded(ctx, event)
	case domain.DRMConfiguredEvent:
		return h.onDRMConfigured(ctx, event)
	}
	return nil
}

func (h liveStreamingDomainHandlers[T]) onLiveStreamCreated(ctx context.Context, event ddd.Event) error {
	stream := event.Payload().(*domain.LiveStream)
	return h.publisher.Publish(ctx, streamspb.LiveStreamAggregateChannel,
		ddd.NewEvent(streamspb.LiveStreamCreatedEvent, &streamspb.LiveStreamCreated{
			Id:           stream.ID(),
			Title:        stream.Title,
			Description:  stream.Description,
			HomeTeam:     stream.HomeTeam,
			AwayTeam:     stream.AwayTeam,
			ScheduledAt:  stream.ScheduledStartTime.Unix(),
			UserSellerId: stream.UserSellerID,
			CategoryId:   stream.CategoryID,
			Thumbnail:    stream.Thumbnail,
		}),
	)
}

func (h liveStreamingDomainHandlers[T]) onStreamingConfigured(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.StreamingConfigured)
	
	// Convert protocols to protobuf format
	protocols := make([]string, len(payload.Protocols))
	for i, p := range payload.Protocols {
		protocols[i] = string(p)
	}
	
	// Convert CDN endpoints
	cdnEndpoints := make([]*streamspb.CDNEndpoint, len(payload.CDNEndpoints))
	for i, cdn := range payload.CDNEndpoints {
		cdnEndpoints[i] = &streamspb.CDNEndpoint{
			Provider: cdn.Provider,
			Url:      cdn.URL,
			Region:   cdn.Region,
		}
	}
	
	return h.publisher.Publish(ctx, streamspb.LiveStreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamingConfiguredEvent, &streamspb.StreamingConfigured{
			StreamId:        payload.StreamID,
			Protocols:       protocols,
			CdnEndpoints:    cdnEndpoints,
			BitrateSettings: payload.BitrateSettings,
		}),
	)
}

func (h liveStreamingDomainHandlers[T]) onLiveStreamStarted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.LiveStreamStarted)
	return h.publisher.Publish(ctx, streamspb.LiveStreamAggregateChannel,
		ddd.NewEvent(streamspb.LiveStreamStartedEvent, &streamspb.LiveStreamStarted{
			StreamId:  payload.StreamID,
			StartedAt: payload.StartedAt.Unix(),
		}),
	)
}

func (h liveStreamingDomainHandlers[T]) onLiveStreamStopped(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.LiveStreamStopped)
	return h.publisher.Publish(ctx, streamspb.LiveStreamAggregateChannel,
		ddd.NewEvent(streamspb.LiveStreamStoppedEvent, &streamspb.LiveStreamStopped{
			StreamId:       payload.StreamID,
			StoppedAt:      payload.StoppedAt.Unix(),
			DurationSec:    int32(payload.Duration.Seconds()),
			TotalViewers:   int32(payload.TotalViewers),
			PeakViewers:    int32(payload.PeakViewers),
			AverageQuality: payload.AverageQuality,
		}),
	)
}

func (h liveStreamingDomainHandlers[T]) onViewerJoined(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ViewerJoined)
	return h.publisher.Publish(ctx, streamspb.LiveStreamAggregateChannel,
		ddd.NewEvent(streamspb.ViewerJoinedEvent, &streamspb.ViewerJoined{
			StreamId:  payload.StreamID,
			ViewerId:  payload.ViewerID,
			JoinedAt:  payload.JoinedAt.Unix(),
			IpAddress: payload.IPAddress,
			UserAgent: payload.UserAgent,
			Country:   payload.Country,
		}),
	)
}

func (h liveStreamingDomainHandlers[T]) onViewerLeft(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ViewerLeft)
	return h.publisher.Publish(ctx, streamspb.LiveStreamAggregateChannel,
		ddd.NewEvent(streamspb.ViewerLeftEvent, &streamspb.ViewerLeft{
			StreamId:       payload.StreamID,
			ViewerId:       payload.ViewerID,
			LeftAt:         payload.LeftAt.Unix(),
			WatchDuration:  int32(payload.WatchDuration.Seconds()),
			AverageQuality: payload.AverageQuality,
		}),
	)
}

func (h liveStreamingDomainHandlers[T]) onStreamQualityChanged(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.StreamQualityChanged)
	return h.publisher.Publish(ctx, streamspb.LiveStreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamQualityChangedEvent, &streamspb.StreamQualityChanged{
			StreamId:    payload.StreamID,
			ViewerId:    payload.ViewerID,
			OldQuality:  payload.OldQuality,
			NewQuality:  payload.NewQuality,
			ChangedAt:   payload.ChangedAt.Unix(),
			Reason:      payload.Reason,
		}),
	)
}

func (h liveStreamingDomainHandlers[T]) onStreamHealthUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.StreamHealthUpdated)
	return h.publisher.Publish(ctx, streamspb.LiveStreamAggregateChannel,
		ddd.NewEvent(streamspb.StreamHealthUpdatedEvent, &streamspb.StreamHealthUpdated{
			StreamId:         payload.StreamID,
			BitrateKbps:      payload.Bitrate,
			FrameRate:        payload.FrameRate,
			PacketLossRate:   payload.PacketLossRate,
			Latency:          payload.Latency,
			BufferHealth:     payload.BufferHealth,
			ConnectionStatus: string(payload.Status),
			UpdatedAt:        payload.UpdatedAt.Unix(),
		}),
	)
}

func (h liveStreamingDomainHandlers[T]) onCDNEndpointAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.CDNEndpointAdded)
	return h.publisher.Publish(ctx, streamspb.LiveStreamAggregateChannel,
		ddd.NewEvent(streamspb.CDNEndpointAddedEvent, &streamspb.CDNEndpointAdded{
			StreamId: payload.StreamID,
			Endpoint: &streamspb.CDNEndpoint{
				Provider: payload.Provider,
				Url:      payload.URL,
				Region:   payload.Region,
			},
		}),
	)
}

func (h liveStreamingDomainHandlers[T]) onDRMConfigured(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.DRMConfigured)
	return h.publisher.Publish(ctx, streamspb.LiveStreamAggregateChannel,
		ddd.NewEvent(streamspb.DRMConfiguredEvent, &streamspb.DRMConfigured{
			StreamId:      payload.StreamID,
			DrmProvider:   payload.Provider,
			ContentId:     payload.ContentID,
			LicenseServer: payload.LicenseServer,
		}),
	)
}