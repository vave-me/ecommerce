package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"middleman/internal/auth"
	"middleman/internal/errorsotel"
	"middleman/streams/internal/application"
	"middleman/streams/internal/application/commands"
	"middleman/streams/internal/application/queries"
	"middleman/streams/internal/domain"
	"middleman/streams/streamspb"
)

type server struct {
	app application.App
	streamspb.UnimplementedStreamsServiceServer
}

var _ streamspb.StreamsServiceServer = (*server)(nil)

// RegisterServer registers the gRPC server implementation
func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	streamspb.RegisterStreamsServiceServer(registrar, server{app: app})
	return nil
}

// Stream Management

func (s server) CreateStream(ctx context.Context, req *streamspb.CreateStreamRequest) (*streamspb.CreateStreamResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	streamID := uuid.New().String()
	
	cmd := commands.CreateStream{
		ID:             streamID,
		Title:          req.GetTitle(),
		Description:    req.GetDescription(),
		Synopsis:       req.GetSynopsis(),
		StreamType:     domain.StreamType(req.GetStreamType().String()),
		StreamURL:      req.GetStreamUrl(),
		ThumbnailURL:   req.GetThumbnailUrl(),
		TrailerURL:     req.GetTrailerUrl(),
		Duration:       req.GetDuration(),
		ContentRating:  domain.ContentRating(req.GetContentRating().String()),
		AccessType:     domain.AccessType(req.GetAccessType().String()),
		Genre:          req.GetGenre(),
		Language:       req.GetLanguage(),
		Country:        req.GetCountry(),
		Studio:         req.GetStudio(),
		ReleaseDate:    req.GetReleaseDate().AsTime(),
		CreatedBy:      claims.Subject,
	}

	if err := s.app.CreateStream(ctx, cmd); err != nil {
		return nil, err
	}

	return &streamspb.CreateStreamResponse{StreamId: streamID}, nil
}

func (s server) GetStream(ctx context.Context, req *streamspb.GetStreamRequest) (*streamspb.GetStreamResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("stream_id", req.GetStreamId()))

	// Optional user ID for access check
	var userID string
	if req.GetUserId() != "" {
		userID = req.GetUserId()
	} else if claims, ok := auth.ClaimsFromContext(ctx); ok {
		userID = claims.Subject
	}

	stream, err := s.app.GetStream(ctx, queries.GetStream{
		StreamID: req.GetStreamId(),
		UserID:   userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Check user access
	hasAccess := false
	var watchProgress int64
	if userID != "" && stream.UserAccess != nil {
		if access, ok := stream.UserAccess[userID]; ok {
			hasAccess = access.HasAccess()
			watchProgress = access.WatchProgress
		}
	}

	return &streamspb.GetStreamResponse{
		Stream:        s.streamFromDomain(stream),
		HasAccess:     hasAccess,
		WatchProgress: watchProgress,
	}, nil
}

func (s server) UpdateStream(ctx context.Context, req *streamspb.UpdateStreamRequest) (*streamspb.UpdateStreamResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("stream_id", req.GetStreamId()))

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	cmd := commands.UpdateStream{
		ID:            req.GetStreamId(),
		Title:         req.GetTitle(),
		Description:   req.GetDescription(),
		Synopsis:      req.GetSynopsis(),
		ThumbnailURL:  req.GetThumbnailUrl(),
		TrailerURL:    req.GetTrailerUrl(),
		ContentRating: domain.ContentRating(req.GetContentRating().String()),
		Genre:         req.GetGenre(),
		Tags:          req.GetTags(),
	}

	if err := s.app.UpdateStream(ctx, cmd); err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &streamspb.UpdateStreamResponse{Success: true}, nil
}

func (s server) PublishStream(ctx context.Context, req *streamspb.PublishStreamRequest) (*streamspb.PublishStreamResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("stream_id", req.GetStreamId()))

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	publishedAt, err := s.app.PublishStream(ctx, commands.PublishStream{
		ID: req.GetStreamId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &streamspb.PublishStreamResponse{
		Success:     true,
		PublishedAt: timestamppb.New(publishedAt),
	}, nil
}

func (s server) ArchiveStream(ctx context.Context, req *streamspb.ArchiveStreamRequest) (*streamspb.ArchiveStreamResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("stream_id", req.GetStreamId()))

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.ArchiveStream(ctx, commands.ArchiveStream{
		ID: req.GetStreamId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &streamspb.ArchiveStreamResponse{Success: true}, nil
}

// Pricing and Access

func (s server) SetStreamPricing(ctx context.Context, req *streamspb.SetStreamPricingRequest) (*streamspb.SetStreamPricingResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("stream_id", req.GetStreamId()))

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	cmd := commands.SetStreamPricing{
		ID:             req.GetStreamId(),
		RentalPrice:    req.GetRentalPrice(),
		RentalDuration: req.GetRentalDuration(),
		PurchasePrice:  req.GetPurchasePrice(),
		PPVPrice:       req.GetPpvPrice(),
	}

	if err := s.app.SetStreamPricing(ctx, cmd); err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &streamspb.SetStreamPricingResponse{Success: true}, nil
}

func (s server) GrantUserAccess(ctx context.Context, req *streamspb.GrantUserAccessRequest) (*streamspb.GrantUserAccessResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("stream_id", req.GetStreamId()),
		attribute.String("user_id", req.GetUserId()),
	)

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	var expiresAt time.Time
	if req.GetDuration() > 0 {
		expiresAt = time.Now().Add(time.Duration(req.GetDuration()) * time.Hour)
	}

	cmd := commands.GrantUserAccess{
		StreamID:   req.GetStreamId(),
		UserID:     req.GetUserId(),
		AccessType: domain.AccessType(req.GetAccessType().String()),
		ExpiresAt:  expiresAt,
	}

	if err := s.app.GrantUserAccess(ctx, cmd); err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &streamspb.GrantUserAccessResponse{
		Success:   true,
		ExpiresAt: timestamppb.New(expiresAt),
	}, nil
}

func (s server) RevokeUserAccess(ctx context.Context, req *streamspb.RevokeUserAccessRequest) (*streamspb.RevokeUserAccessResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("stream_id", req.GetStreamId()),
		attribute.String("user_id", req.GetUserId()),
	)

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	err := s.app.RevokeUserAccess(ctx, commands.RevokeUserAccess{
		StreamID: req.GetStreamId(),
		UserID:   req.GetUserId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &streamspb.RevokeUserAccessResponse{Success: true}, nil
}

// User Interaction

func (s server) UpdateWatchProgress(ctx context.Context, req *streamspb.UpdateWatchProgressRequest) (*streamspb.UpdateWatchProgressResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("stream_id", req.GetStreamId()),
		attribute.String("user_id", req.GetUserId()),
		attribute.Int64("progress", req.GetProgress()),
	)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Ensure the user can only update their own progress
	if req.GetUserId() != claims.Subject {
		return nil, status.Error(grpc_code.PermissionDenied, "cannot update other user's progress")
	}

	cmd := commands.UpdateWatchProgress{
		StreamID:  req.GetStreamId(),
		UserID:    req.GetUserId(),
		Progress:  req.GetProgress(),
		Completed: req.GetCompleted(),
	}

	if err := s.app.UpdateWatchProgress(ctx, cmd); err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &streamspb.UpdateWatchProgressResponse{Success: true}, nil
}

func (s server) RateStream(ctx context.Context, req *streamspb.RateStreamRequest) (*streamspb.RateStreamResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("stream_id", req.GetStreamId()),
		attribute.String("user_id", req.GetUserId()),
		attribute.Int64("rating", int64(req.GetRating())),
	)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Ensure the user can only rate as themselves
	if req.GetUserId() != claims.Subject {
		return nil, status.Error(grpc_code.PermissionDenied, "cannot rate as another user")
	}

	cmd := commands.RateStream{
		StreamID: req.GetStreamId(),
		UserID:   req.GetUserId(),
		Rating:   req.GetRating(),
		IsLike:   req.GetIsLike(),
	}

	if err := s.app.RateStream(ctx, cmd); err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &streamspb.RateStreamResponse{Success: true}, nil
}

// Discovery

func (s server) GetCatalog(ctx context.Context, req *streamspb.GetCatalogRequest) (*streamspb.GetCatalogResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}

	page := req.GetPage()
	if page <= 0 {
		page = 1
	}

	// Convert filters
	var filters *domain.StreamFilters
	if req.GetFilters() != nil {
		f := req.GetFilters()
		filters = &domain.StreamFilters{
			StreamType:  domain.StreamType(f.GetStreamType().String()),
			Genre:       f.GetGenre(),
			AccessType:  domain.AccessType(f.GetAccessType().String()),
			SortBy:      f.GetSortBy(),
			SearchQuery: f.GetSearchQuery(),
		}
	}

	catalog, err := s.app.GetCatalog(ctx, queries.GetCatalog{
		UserID:   req.GetUserId(),
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Convert streams to summaries
	summaries := make([]*streamspb.StreamSummary, len(catalog.Streams))
	for i, stream := range catalog.Streams {
		summaries[i] = s.streamSummaryFromDomain(stream)
	}

	// Convert categories
	categories := make([]*streamspb.Category, len(catalog.Categories))
	for i, cat := range catalog.Categories {
		categories[i] = &streamspb.Category{
			Id:          cat.ID,
			Name:        cat.Name,
			Slug:        cat.Slug,
			StreamCount: cat.StreamCount,
		}
	}

	totalPages := (catalog.TotalCount + int32(pageSize) - 1) / int32(pageSize)

	return &streamspb.GetCatalogResponse{
		Streams:     summaries,
		TotalCount:  catalog.TotalCount,
		Page:        page,
		PageSize:    pageSize,
		HasMore:     page < totalPages,
		Categories:  categories,
		FeaturedIds: catalog.FeaturedIDs,
	}, nil
}

func (s server) SearchStreams(ctx context.Context, req *streamspb.SearchStreamsRequest) (*streamspb.SearchStreamsResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("query", req.GetQuery()))

	var filters *domain.StreamFilters
	if req.GetFilters() != nil {
		f := req.GetFilters()
		filters = &domain.StreamFilters{
			StreamType:     domain.StreamType(f.GetStreamType().String()),
			MinDuration:    f.GetMinDuration(),
			MaxDuration:    f.GetMaxDuration(),
			ContentRating:  s.contentRatingsFromProto(f.GetContentRating()),
			Genre:          f.GetGenre(),
			Language:       f.GetLanguage(),
			Country:        f.GetCountry(),
			Studio:         f.GetStudio(),
			AccessType:     domain.AccessType(f.GetAccessType().String()),
			MinRating:      f.GetMinRating(),
			ReleasedAfter:  f.GetReleasedAfter(),
			ReleasedBefore: f.GetReleasedBefore(),
			SortBy:         f.GetSortBy(),
			SortOrder:      f.GetSortOrder(),
		}
	}

	streams, totalCount, err := s.app.SearchStreams(ctx, queries.SearchStreams{
		Query:   req.GetQuery(),
		Filters: filters,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	pbStreams := make([]*streamspb.Stream, len(streams))
	for i, stream := range streams {
		pbStreams[i] = s.streamFromDomain(stream)
	}

	return &streamspb.SearchStreamsResponse{
		Streams:    pbStreams,
		TotalCount: totalCount,
	}, nil
}

func (s server) GetUserStreams(ctx context.Context, req *streamspb.GetUserStreamsRequest) (*streamspb.GetUserStreamsResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("user_id", req.GetUserId()))

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Users can only get their own streams
	if req.GetUserId() != claims.Subject {
		return nil, status.Error(grpc_code.PermissionDenied, "cannot access other user's streams")
	}

	streams, err := s.app.GetUserStreams(ctx, queries.GetUserStreams{
		UserID: req.GetUserId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	pbStreams := make([]*streamspb.Stream, len(streams))
	for i, stream := range streams {
		pbStreams[i] = s.streamFromDomain(stream)
	}

	return &streamspb.GetUserStreamsResponse{Streams: pbStreams}, nil
}

func (s server) GetContinueWatching(ctx context.Context, req *streamspb.GetContinueWatchingRequest) (*streamspb.GetContinueWatchingResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("user_id", req.GetUserId()))

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	// Users can only get their own continue watching list
	if req.GetUserId() != claims.Subject {
		return nil, status.Error(grpc_code.PermissionDenied, "cannot access other user's watch list")
	}

	items, err := s.app.GetContinueWatching(ctx, queries.GetContinueWatching{
		UserID: req.GetUserId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	pbItems := make([]*streamspb.ContinueWatchingItem, len(items))
	for i, item := range items {
		pbItems[i] = &streamspb.ContinueWatchingItem{
			Stream:         s.streamFromDomain(item.Stream),
			WatchProgress:  item.WatchProgress,
			LastWatchedAt:  timestamppb.New(item.LastWatchedAt),
		}
	}

	return &streamspb.GetContinueWatchingResponse{Items: pbItems}, nil
}

// Series Management

func (s server) CreateSeries(ctx context.Context, req *streamspb.CreateSeriesRequest) (*streamspb.CreateSeriesResponse, error) {
	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	seriesID := uuid.New().String()

	cmd := commands.CreateSeries{
		ID:           seriesID,
		Title:        req.GetTitle(),
		Description:  req.GetDescription(),
		ThumbnailURL: req.GetThumbnailUrl(),
		Genre:        req.GetGenre(),
		Studio:       req.GetStudio(),
	}

	if err := s.app.CreateSeries(ctx, cmd); err != nil {
		return nil, err
	}

	return &streamspb.CreateSeriesResponse{SeriesId: seriesID}, nil
}

func (s server) GetSeries(ctx context.Context, req *streamspb.GetSeriesRequest) (*streamspb.GetSeriesResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("series_id", req.GetSeriesId()))

	series, err := s.app.GetSeries(ctx, queries.GetSeries{
		SeriesID: req.GetSeriesId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &streamspb.GetSeriesResponse{
		Series: s.seriesFromDomain(series),
	}, nil
}

func (s server) AddSeason(ctx context.Context, req *streamspb.AddSeasonRequest) (*streamspb.AddSeasonResponse, error) {
	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	seasonID := uuid.New().String()

	cmd := commands.AddSeason{
		SeasonID:     seasonID,
		SeriesID:     req.GetSeriesId(),
		SeasonNumber: req.GetSeasonNumber(),
		Title:        req.GetTitle(),
		Description:  req.GetDescription(),
		ThumbnailURL: req.GetThumbnailUrl(),
	}

	if err := s.app.AddSeason(ctx, cmd); err != nil {
		return nil, err
	}

	return &streamspb.AddSeasonResponse{SeasonId: seasonID}, nil
}

func (s server) AddEpisode(ctx context.Context, req *streamspb.AddEpisodeRequest) (*streamspb.AddEpisodeResponse, error) {
	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	episodeID := uuid.New().String()

	cmd := commands.AddEpisode{
		EpisodeID:     episodeID,
		SeriesID:      req.GetSeriesId(),
		SeasonNumber:  req.GetSeasonNumber(),
		EpisodeNumber: req.GetEpisodeNumber(),
		StreamID:      req.GetStreamId(),
		Title:         req.GetTitle(),
		Duration:      req.GetDuration(),
		AirDate:       req.GetAirDate().AsTime(),
	}

	if err := s.app.AddEpisode(ctx, cmd); err != nil {
		return nil, err
	}

	return &streamspb.AddEpisodeResponse{EpisodeId: episodeID}, nil
}

// Helper methods

func (s server) streamFromDomain(stream *domain.Stream) *streamspb.Stream {
	// Convert subtitles
	subtitles := make([]*streamspb.Subtitle, len(stream.Subtitles))
	for i, sub := range stream.Subtitles {
		subtitles[i] = &streamspb.Subtitle{
			Language: sub.Language,
			Url:      sub.URL,
			Default:  sub.Default,
		}
	}

	// Convert audio tracks
	audioTracks := make([]*streamspb.AudioTrack, len(stream.AudioTracks))
	for i, track := range stream.AudioTracks {
		audioTracks[i] = &streamspb.AudioTrack{
			Language: track.Language,
			Type:     track.Type,
			Default:  track.Default,
		}
	}

	// Convert cast
	cast := make([]*streamspb.CastMember, len(stream.Cast))
	for i, member := range stream.Cast {
		cast[i] = &streamspb.CastMember{
			Name:      member.Name,
			Role:      member.Role,
			Character: member.Character,
			ImageUrl:  member.ImageURL,
		}
	}

	// Convert qualities
	qualities := make([]streamspb.StreamQuality, len(stream.AvailableQualities))
	for i, q := range stream.AvailableQualities {
		qualities[i] = s.qualityToProto(q)
	}

	// Convert subscription tiers
	tiers := make([]string, len(stream.SubscriptionTiers))
	copy(tiers, stream.SubscriptionTiers)

	return &streamspb.Stream{
		Id:                  stream.ID,
		Title:               stream.Title,
		Description:         stream.Description,
		Synopsis:            stream.Synopsis,
		StreamType:          s.streamTypeToProto(stream.StreamType),
		Status:              s.streamStatusToProto(stream.Status),
		StreamUrl:           stream.StreamURL,
		ThumbnailUrl:        stream.ThumbnailURL,
		TrailerUrl:          stream.TrailerURL,
		Duration:            stream.Duration,
		ReleaseDate:         timestamppb.New(stream.ReleaseDate),
		ContentRating:       s.contentRatingToProto(stream.ContentRating),
		AvailableQualities:  qualities,
		DefaultQuality:      s.qualityToProto(stream.DefaultQuality),
		Subtitles:           subtitles,
		AudioTracks:         audioTracks,
		AccessType:          s.accessTypeToProto(stream.AccessType),
		SubscriptionTiers:   tiers,
		RentalPrice:         stream.RentalPrice,
		RentalDuration:      stream.RentalDuration,
		PurchasePrice:       stream.PurchasePrice,
		PpvPrice:            stream.PPVPrice,
		Genre:               stream.Genre,
		Tags:                stream.Tags,
		Cast:                cast,
		Directors:           stream.Directors,
		Producers:           stream.Producers,
		Studio:              stream.Studio,
		Language:            stream.Language,
		Country:             stream.Country,
		ViewCount:           stream.ViewCount,
		LikeCount:           stream.LikeCount,
		DislikeCount:        stream.DislikeCount,
		AverageRating:       stream.AverageRating,
		TotalRevenue:        stream.TotalRevenue,
		SeriesId:            stream.SeriesID,
		SeasonNumber:        stream.SeasonNumber,
		EpisodeNumber:       stream.EpisodeNumber,
		CreatedAt:           timestamppb.New(stream.CreatedAt),
		UpdatedAt:           timestamppb.New(stream.UpdatedAt),
		PublishedAt:         timestamppb.New(stream.PublishedAt),
	}
}

func (s server) streamSummaryFromDomain(stream *domain.Stream) *streamspb.StreamSummary {
	releaseYear := int32(stream.ReleaseDate.Year())
	isNew := time.Since(stream.PublishedAt) < 30*24*time.Hour // New if published in last 30 days

	return &streamspb.StreamSummary{
		Id:            stream.ID,
		Title:         stream.Title,
		ThumbnailUrl:  stream.ThumbnailURL,
		Duration:      stream.Duration,
		ContentRating: s.contentRatingToProto(stream.ContentRating),
		AccessType:    s.accessTypeToProto(stream.AccessType),
		Rating:        stream.AverageRating,
		ViewCount:     stream.ViewCount,
		ReleaseYear:   releaseYear,
		IsNew:         isNew,
		HasAccess:     false, // Will be set based on user context
	}
}

func (s server) seriesFromDomain(series *domain.Series) *streamspb.Series {
	// Convert seasons
	seasons := make([]*streamspb.Season, len(series.Seasons))
	for i, season := range series.Seasons {
		// Convert episodes
		episodes := make([]*streamspb.Episode, len(season.Episodes))
		for j, ep := range season.Episodes {
			episodes[j] = &streamspb.Episode{
				Id:            ep.ID,
				EpisodeNumber: ep.EpisodeNumber,
				StreamId:      ep.StreamID,
				Title:         ep.Title,
				Duration:      ep.Duration,
				AirDate:       timestamppb.New(ep.AirDate),
			}
		}

		seasons[i] = &streamspb.Season{
			Id:            season.ID,
			SeasonNumber:  season.SeasonNumber,
			Title:         season.Title,
			Description:   season.Description,
			ThumbnailUrl:  season.ThumbnailURL,
			Episodes:      episodes,
			TotalEpisodes: int32(len(episodes)),
			CreatedAt:     timestamppb.New(season.CreatedAt),
		}
	}

	return &streamspb.Series{
		Id:           series.ID,
		Title:        series.Title,
		Description:  series.Description,
		ThumbnailUrl: series.ThumbnailURL,
		Genre:        series.Genre,
		Studio:       series.Studio,
		Seasons:      seasons,
		TotalSeasons: int32(len(seasons)),
		Status:       series.Status,
		CreatedAt:    timestamppb.New(series.CreatedAt),
		UpdatedAt:    timestamppb.New(series.UpdatedAt),
	}
}

// Proto conversion helpers

func (s server) streamTypeToProto(st domain.StreamType) streamspb.StreamType {
	switch st {
	case domain.StreamTypeMovie:
		return streamspb.StreamType_STREAM_TYPE_MOVIE
	case domain.StreamTypeSeries:
		return streamspb.StreamType_STREAM_TYPE_SERIES
	case domain.StreamTypeDocumentary:
		return streamspb.StreamType_STREAM_TYPE_DOCUMENTARY
	case domain.StreamTypeLive:
		return streamspb.StreamType_STREAM_TYPE_LIVE
	case domain.StreamTypeEducational:
		return streamspb.StreamType_STREAM_TYPE_EDUCATIONAL
	case domain.StreamTypeMusic:
		return streamspb.StreamType_STREAM_TYPE_MUSIC
	case domain.StreamTypeSports:
		return streamspb.StreamType_STREAM_TYPE_SPORTS
	default:
		return streamspb.StreamType_STREAM_TYPE_UNSPECIFIED
	}
}

func (s server) streamStatusToProto(status domain.StreamStatus) streamspb.StreamStatus {
	switch status {
	case domain.StreamStatusDraft:
		return streamspb.StreamStatus_STREAM_STATUS_DRAFT
	case domain.StreamStatusProcessing:
		return streamspb.StreamStatus_STREAM_STATUS_PROCESSING
	case domain.StreamStatusPublished:
		return streamspb.StreamStatus_STREAM_STATUS_PUBLISHED
	case domain.StreamStatusArchived:
		return streamspb.StreamStatus_STREAM_STATUS_ARCHIVED
	case domain.StreamStatusDeleted:
		return streamspb.StreamStatus_STREAM_STATUS_DELETED
	default:
		return streamspb.StreamStatus_STREAM_STATUS_UNSPECIFIED
	}
}

func (s server) contentRatingToProto(rating domain.ContentRating) streamspb.ContentRating {
	switch rating {
	case domain.ContentRatingG:
		return streamspb.ContentRating_CONTENT_RATING_G
	case domain.ContentRatingPG:
		return streamspb.ContentRating_CONTENT_RATING_PG
	case domain.ContentRatingPG13:
		return streamspb.ContentRating_CONTENT_RATING_PG13
	case domain.ContentRatingR:
		return streamspb.ContentRating_CONTENT_RATING_R
	case domain.ContentRatingNC17:
		return streamspb.ContentRating_CONTENT_RATING_NC17
	case domain.ContentRatingUnrated:
		return streamspb.ContentRating_CONTENT_RATING_UNRATED
	default:
		return streamspb.ContentRating_CONTENT_RATING_UNSPECIFIED
	}
}

func (s server) accessTypeToProto(at domain.AccessType) streamspb.AccessType {
	switch at {
	case domain.AccessTypeFree:
		return streamspb.AccessType_ACCESS_TYPE_FREE
	case domain.AccessTypeSubscription:
		return streamspb.AccessType_ACCESS_TYPE_SUBSCRIPTION
	case domain.AccessTypeRental:
		return streamspb.AccessType_ACCESS_TYPE_RENTAL
	case domain.AccessTypePurchase:
		return streamspb.AccessType_ACCESS_TYPE_PURCHASE
	case domain.AccessTypePPV:
		return streamspb.AccessType_ACCESS_TYPE_PPV
	default:
		return streamspb.AccessType_ACCESS_TYPE_UNSPECIFIED
	}
}

func (s server) qualityToProto(q domain.StreamQuality) streamspb.StreamQuality {
	switch q {
	case domain.StreamQualitySD:
		return streamspb.StreamQuality_STREAM_QUALITY_SD
	case domain.StreamQualityHD:
		return streamspb.StreamQuality_STREAM_QUALITY_HD
	case domain.StreamQualityFullHD:
		return streamspb.StreamQuality_STREAM_QUALITY_FULL_HD
	case domain.StreamQuality4K:
		return streamspb.StreamQuality_STREAM_QUALITY_4K
	default:
		return streamspb.StreamQuality_STREAM_QUALITY_UNSPECIFIED
	}
}

func (s server) contentRatingsFromProto(ratings []streamspb.ContentRating) []domain.ContentRating {
	result := make([]domain.ContentRating, len(ratings))
	for i, r := range ratings {
		switch r {
		case streamspb.ContentRating_CONTENT_RATING_G:
			result[i] = domain.ContentRatingG
		case streamspb.ContentRating_CONTENT_RATING_PG:
			result[i] = domain.ContentRatingPG
		case streamspb.ContentRating_CONTENT_RATING_PG13:
			result[i] = domain.ContentRatingPG13
		case streamspb.ContentRating_CONTENT_RATING_R:
			result[i] = domain.ContentRatingR
		case streamspb.ContentRating_CONTENT_RATING_NC17:
			result[i] = domain.ContentRatingNC17
		case streamspb.ContentRating_CONTENT_RATING_UNRATED:
			result[i] = domain.ContentRatingUnrated
		}
	}
	return result
}