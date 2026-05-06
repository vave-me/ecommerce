package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/streams/internal/domain"
	"middleman/streams/internal/infrastructure/streaming"
)

type ConfigureStreaming struct {
	StreamID           string
	Protocols          []domain.StreamingProtocol
	QualityProfiles    []domain.StreamingQualityProfile
	AdaptiveBitrate    bool
	LowLatencyMode     bool
	DVREnabled         bool
	DVRWindowMinutes   int
	CDNEndpoints       []domain.CDNEndpoint
	PrimaryCDN         string
	EnableFailover     bool
	EnableDRM          bool
	DRMProviders       map[string]domain.DRMConfig
}

type ConfigureStreamingHandler struct {
	liveStreams     domain.LiveStreamRepository
	publisher       ddd.EventPublisher
	streamingServer *streaming.StreamingServer
	cdnManager      *streaming.CDNManager
	drmManager      *streaming.DRMManager
}

func NewConfigureStreamingHandler(
	liveStreams domain.LiveStreamRepository,
	publisher ddd.EventPublisher,
	streamingServer *streaming.StreamingServer,
	cdnManager *streaming.CDNManager,
	drmManager *streaming.DRMManager,
) ConfigureStreamingHandler {
	return ConfigureStreamingHandler{
		liveStreams:     liveStreams,
		publisher:       publisher,
		streamingServer: streamingServer,
		cdnManager:      cdnManager,
		drmManager:      drmManager,
	}
}

func (h ConfigureStreamingHandler) ConfigureStreaming(ctx context.Context, cmd ConfigureStreaming) error {
	stream, err := h.liveStreams.Find(ctx, cmd.StreamID)
	if err != nil {
		return err
	}

	// Configure streaming protocols and quality
	if _, err := stream.ConfigureStreaming(
		cmd.Protocols,
		cmd.QualityProfiles,
		cmd.AdaptiveBitrate,
		cmd.LowLatencyMode,
		cmd.DVREnabled,
		cmd.DVRWindowMinutes,
	); err != nil {
		return err
	}

	// Configure CDN
	if len(cmd.CDNEndpoints) > 0 {
		if _, err := stream.ConfigureCDN(
			cmd.CDNEndpoints,
			cmd.PrimaryCDN,
			cmd.EnableFailover,
		); err != nil {
			return err
		}
	}

	// Configure DRM if enabled
	if cmd.EnableDRM && len(cmd.DRMProviders) > 0 {
		if _, err := stream.ConfigureDRM(cmd.DRMProviders, true); err != nil {
			return err
		}

		// Generate content key for DRM
		contentKey, err := h.drmManager.GenerateContentKey(cmd.StreamID)
		if err != nil {
			return err
		}

		// Store key ID with stream
		stream.DRMKeyID = contentKey.KeyID
	}

	if err := h.liveStreams.Save(ctx, stream); err != nil {
		return err
	}

	// Publish domain events
	if err := h.publisher.Publish(ctx, stream.Events()...); err != nil {
		return err
	}

	return nil
}