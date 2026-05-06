package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/streams/internal/domain"
	"middleman/streams/internal/infrastructure/streaming"
)

type StartLiveStream struct {
	StreamID        string
	IngestProtocol  domain.StreamingProtocol
	IngestURL       string
	BackupIngestURL string
	StreamKey       string
}

type StartLiveStreamHandler struct {
	liveStreams     domain.LiveStreamRepository
	publisher       ddd.EventPublisher
	streamingServer *streaming.StreamingServer
}

func NewStartLiveStreamHandler(
	liveStreams domain.LiveStreamRepository,
	publisher ddd.EventPublisher,
	streamingServer *streaming.StreamingServer,
) StartLiveStreamHandler {
	return StartLiveStreamHandler{
		liveStreams:     liveStreams,
		publisher:       publisher,
		streamingServer: streamingServer,
	}
}

func (h StartLiveStreamHandler) StartLiveStream(ctx context.Context, cmd StartLiveStream) error {
	stream, err := h.liveStreams.Find(ctx, cmd.StreamID)
	if err != nil {
		return err
	}

	// Set ingestion configuration
	if _, err := stream.SetIngestionConfig(
		cmd.IngestProtocol,
		cmd.IngestURL,
		cmd.BackupIngestURL,
		cmd.StreamKey,
	); err != nil {
		return err
	}

	// Start the stream
	if _, err := stream.StartStream(); err != nil {
		return err
	}

	// Start streaming server processing
	if err := h.streamingServer.StartStream(
		cmd.StreamID,
		cmd.StreamKey,
		string(cmd.IngestProtocol),
	); err != nil {
		return err
	}

	// Generate manifest URLs
	hlsURL := "/hls/" + cmd.StreamID + "/master.m3u8"
	dashURL := "/dash/" + cmd.StreamID + "/manifest.mpd"

	if _, err := stream.SetManifestURL(domain.ProtocolHLS, hlsURL); err != nil {
		return err
	}

	if _, err := stream.SetManifestURL(domain.ProtocolDASH, dashURL); err != nil {
		return err
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