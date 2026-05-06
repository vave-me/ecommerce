package grpc

import (
	"context"
	"fmt"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"middleman/internal/auth"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"middleman/media/internal/application"
	"middleman/media/internal/application/commands"
	"middleman/media/internal/application/queries"
	"middleman/media/internal/constants"
	"middleman/media/internal/domain"
	"middleman/media/mediapb"
)

// server implements mediapb.MediaServiceServer
type server struct {
	app application.App
	mediapb.UnimplementedMediaServiceServer
}

var _ mediapb.MediaServiceServer = (*server)(nil)

// RegisterServer wires up the gRPC service.
func RegisterServer(ctx context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	// <<< CHANGED: Read environment variable for viewURLBase

	mediapb.RegisterMediaServiceServer(registrar, &server{
		app: app,
	})
	log.Println("Media service server registered successfully")
	return nil
}

// -----------------------------------------------------------------------------
// Generate a Presigned URL using MinIO (replacing old S3 logic).
// -----------------------------------------------------------------------------

func (s *server) generatePresignedUploadURL(ctx context.Context, mediaType, objectKey string) (string, error) {

	minioClient := di.Get(ctx, constants.MinioClient).(*minio.Client)

	// We can optionally figure out contentType if needed:
	// contentType := contentTypeForMediaType(mediaType)

	// Generate a presigned PUT URL valid for 15 minutes
	presignedURL, err := minioClient.PresignedPutObject(ctx, "classified", objectKey, 15*time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

// Optional helper if you want to handle content types in future
func contentTypeForMediaType(mediaType string) string {
	switch mediaType {
	case "image":
		return "image/jpeg"
	case "video":
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}

// -----------------------------------------------------------------------------
// Delete a single object in MinIO by key
// -----------------------------------------------------------------------------

func (s *server) deleteMinioObject(ctx context.Context, objectKey string) error {
	minioClient := di.Get(ctx, constants.MinioClient).(*minio.Client)

	err := minioClient.RemoveObject(ctx, "classified", objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("Failed to delete MinIO object %s: %v", objectKey, err)
		return err
	}
	log.Printf("Successfully deleted MinIO object: %s\n", objectKey)
	return nil
}

// -----------------------------------------------------------------------------
// Delete all objects matching a prefix in MinIO
// -----------------------------------------------------------------------------

func (s *server) deleteMinioPrefix(ctx context.Context, prefix string) error {
	minioClient := di.Get(ctx, constants.MinioClient).(*minio.Client)

	// 1) List objects with the given prefix
	objectCh := minioClient.ListObjects(ctx, "classified", minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var lastErr error
	// 2) Delete each object
	for object := range objectCh {
		if object.Err != nil {
			log.Printf("Error listing objects with prefix %q: %v", prefix, object.Err)
			return object.Err
		}

		err := minioClient.RemoveObject(ctx, "classified", object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("Error deleting object %s: %v", object.Key, err)
			lastErr = err
			continue
		}
		log.Printf("Deleted object: %s\n", object.Key)
	}
	return lastErr
}

// -----------------------------------------------------------------------------
// gRPC Methods
// -----------------------------------------------------------------------------

func (s server) CreateMedia(ctx context.Context, request *mediapb.CreateMediaRequest) (*mediapb.CreateMediaResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	mediaID := uuid.New().String()
	span.SetAttributes(attribute.String("MediaID", mediaID))

	err := s.app.CreateMedia(ctx, commands.CreateMedia{
		ID:       mediaID,
		ItemID:   request.GetItemId(),
		ItemType: domain.ItemType(request.GetItemType()),
		UserID:   userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &mediapb.CreateMediaResponse{Id: mediaID}, nil
}
func (s server) UpdateMedia(ctx context.Context, request *mediapb.UpdateMediaRequest) (*mediapb.UpdateMediaResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("MediaID", request.GetId()))

	err := s.app.UpdateMedia(ctx, commands.UpdateMedia{
		ID:       request.GetId(),
		ItemID:   request.GetItemId(),
		ItemType: domain.ItemType(request.GetItemType()),
		UserID:   userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &mediapb.UpdateMediaResponse{Id: request.GetId()}, nil
}
func (s server) AddImage(ctx context.Context, request *mediapb.AddImageRequest) (*mediapb.AddImageResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	span := trace.SpanFromContext(ctx)

	imageID := uuid.New().String()
	span.SetAttributes(attribute.String("ImageID", imageID))

	objectKey := fmt.Sprintf("%s/%s", request.GetMediaId(), imageID)
	presignedURL, err := s.generatePresignedUploadURL(ctx, "image", objectKey)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	//TODO uncomment when cdn launched
	// viewURL := fmt.Sprintf("https://d1yulz4hxrt30n.cloudfront.net/%s", objectKey)

	// <<< CHANGED: Use the configured viewURLBase
	viewURL := fmt.Sprintf("%s/%s", "http://192.168.178.84:9096/classified", objectKey)
	//viewURL := fmt.Sprintf("%s/%s", "https://minio-api.sfx-markt.de/classified", objectKey)
	err = s.app.AddImage(ctx, commands.AddImage{
		ID:           imageID,
		MediaID:      request.GetMediaId(),
		DisplayOrder: int(request.GetDisplayOrder()),
		IsMain:       request.GetIsMain(),
		Url:          viewURL,
		Metadata:     request.GetMetadata(),
		FileType:     request.GetFileType(),
		Thumbnail:    request.GetThumbnail(),
		UserID:       userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &mediapb.AddImageResponse{
		Url:     presignedURL,
		ViewUrl: viewURL,
	}, nil
}

func (s server) AddVideo(ctx context.Context, request *mediapb.AddVideoRequest) (*mediapb.AddVideoResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	videoID := uuid.New().String()
	span.SetAttributes(attribute.String("VideoID", videoID))

	objectKey := fmt.Sprintf("%s/%s.%s", request.GetMediaId(), videoID, request.GetFileType())
	presignedURL, err := s.generatePresignedUploadURL(ctx, "video", objectKey)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	//viewURL := fmt.Sprintf("https://d1yulz4hxrt30n.cloudfront.net/%s", objectKey)
	// <<< CHANGED: Use the configured viewURLBase
	viewURL := fmt.Sprintf("%s/%s", "http://192.168.178.84:9096/classified", objectKey)
	//viewURL := fmt.Sprintf("%s/%s", "https://minio-api.sfx-markt.de/classified", objectKey)

	err = s.app.AddVideo(ctx, commands.AddVideo{
		ID:           videoID,
		MediaID:      request.GetMediaId(),
		DisplayOrder: int(request.GetDisplayOrder()),
		IsMain:       request.GetIsMain(),
		Url:          viewURL,
		Metadata:     request.GetMetadata(),
		FileType:     request.GetFileType(),
		Thumbnail:    request.GetThumbnail(),
		UserID:       userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &mediapb.AddVideoResponse{
		Url:     presignedURL,
		ViewUrl: viewURL,
	}, nil
}

func (s server) GetAllItemImages(ctx context.Context, request *mediapb.GetAllItemImagesRequest) (*mediapb.GetAllItemImagesResponse, error) {
	span := trace.SpanFromContext(ctx)

	images, err := s.app.GetAllItemImages(ctx, queries.GetAllItemImages{
		ItemID: request.GetItemId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoImages := make([]*mediapb.Image, len(images))
	for i, image := range images {
		protoImages[i] = s.imageFromDomain(image)
	}
	return &mediapb.GetAllItemImagesResponse{Images: protoImages}, nil
}

func (s server) GetAllMediaImages(ctx context.Context, request *mediapb.GetAllMediaImagesRequest) (*mediapb.GetAllMediaImagesResponse, error) {
	span := trace.SpanFromContext(ctx)

	images, err := s.app.GetAllMediaImages(ctx, queries.GetAllMediaImages{
		MediaID: request.GetMediaId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoImages := make([]*mediapb.Image, len(images))
	for i, image := range images {
		protoImages[i] = s.imageFromDomain(image)
	}
	return &mediapb.GetAllMediaImagesResponse{Images: protoImages}, nil
}

func (s server) GetAllItemVideos(ctx context.Context, request *mediapb.GetAllItemVideosRequest) (*mediapb.GetAllItemVideosResponse, error) {
	span := trace.SpanFromContext(ctx)

	videos, err := s.app.GetAllItemVideos(ctx, queries.GetAllItemVideos{
		ItemID: request.GetItemId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoVideos := make([]*mediapb.Video, len(videos))
	for i, video := range videos {
		protoVideos[i] = s.videoFromDomain(video)
	}
	return &mediapb.GetAllItemVideosResponse{Videos: protoVideos}, nil
}

func (s server) GetAllVideos(ctx context.Context, request *mediapb.GetAllVideosRequest) (*mediapb.GetAllVideosResponse, error) {
	span := trace.SpanFromContext(ctx)

	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}
	page := request.GetPage()
	if page <= 0 {
		page = 1
	}

	videos, totalCount, err := s.app.GetAllVideos(ctx, queries.GetAllVideos{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    request.GetSortBy(),
		SortOrder: request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	protoVideos := make([]*mediapb.Video, len(videos))
	for i, video := range videos {
		protoVideos[i] = s.videoFromDomain(video)
	}
	return &mediapb.GetAllVideosResponse{
		Videos:      protoVideos,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}

func (s server) GetAllVideosByAuthor(ctx context.Context, request *mediapb.GetAllVideosByAuthorRequest) (*mediapb.GetAllVideosByAuthorResponse, error) {
	span := trace.SpanFromContext(ctx)

	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}
	page := request.GetPage()
	if page <= 0 {
		page = 1
	}

	videos, totalCount, err := s.app.GetAllVideosByAuthor(ctx, queries.GetAllVideosByAuthor{
		UserID:    request.GetUserId(),
		Page:      page,
		PageSize:  pageSize,
		SortBy:    request.GetSortBy(),
		SortOrder: request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	protoVideos := make([]*mediapb.Video, len(videos))
	for i, video := range videos {
		protoVideos[i] = s.videoFromDomain(video)
	}
	return &mediapb.GetAllVideosByAuthorResponse{
		Videos:      protoVideos,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}

func (s server) GetAllImagesByAuthor(ctx context.Context, request *mediapb.GetAllImagesByAuthorRequest) (*mediapb.GetAllImagesByAuthorResponse, error) {
	span := trace.SpanFromContext(ctx)

	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}
	page := request.GetPage()
	if page <= 0 {
		page = 1
	}

	images, totalCount, err := s.app.GetAllImagesByAuthor(ctx, queries.GetAllImagesByAuthor{
		UserID:    request.GetUserId(),
		Page:      page,
		PageSize:  pageSize,
		SortBy:    request.GetSortBy(),
		SortOrder: request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	protoImages := make([]*mediapb.Image, len(images))
	for i, image := range images {
		protoImages[i] = s.imageFromDomain(image)
	}
	return &mediapb.GetAllImagesByAuthorResponse{
		Images:      protoImages,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}

func (s server) GetAllMediaVideos(ctx context.Context, request *mediapb.GetAllMediaVideosRequest) (*mediapb.GetAllMediaVideosResponse, error) {
	span := trace.SpanFromContext(ctx)

	videos, err := s.app.GetAllMediaVideos(ctx, queries.GetAllMediaVideos{
		MediaID: request.GetMediaId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoVideos := make([]*mediapb.Video, len(videos))
	for i, video := range videos {
		protoVideos[i] = s.videoFromDomain(video)
	}
	return &mediapb.GetAllMediaVideosResponse{Videos: protoVideos}, nil
}

func (s server) GetMedia(ctx context.Context, request *mediapb.GetMediaRequest) (*mediapb.GetMediaResponse, error) {
	span := trace.SpanFromContext(ctx)

	media, err := s.app.GetMedia(ctx, queries.GetMedia{
		MediaID: request.GetMediaId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoMedia := s.mediaFromDomain(media)
	return &mediapb.GetMediaResponse{Media: protoMedia}, nil
}

func (s server) GetMediaByItem(ctx context.Context, request *mediapb.GetMediaByItemRequest) (*mediapb.GetMediaByItemResponse, error) {
	span := trace.SpanFromContext(ctx)

	media, err := s.app.GetMediaByItem(ctx, queries.GetMediaByItem{
		ItemID: request.GetItemId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoMedia := s.mediaFromDomain(media)
	return &mediapb.GetMediaByItemResponse{Media: protoMedia}, nil
}

func (s server) RemoveMedia(ctx context.Context, request *mediapb.RemoveMediaRequest) (*mediapb.RemoveMediaResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	span := trace.SpanFromContext(ctx)

	// 1) Remove all objects in MinIO with prefix: mediaId/
	prefix := fmt.Sprintf("%s/", request.GetMediaId())
	if err := s.deleteMinioPrefix(ctx, prefix); err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// 2) Remove from DB/domain
	err := s.app.RemoveMedia(ctx, commands.RemoveMedia{
		ID:     request.GetMediaId(),
		ItemID: request.GetItemId(),
		UserID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &mediapb.RemoveMediaResponse{}, nil
}

func (s server) RemoveImage(ctx context.Context, request *mediapb.RemoveImageRequest) (*mediapb.RemoveImageResponse, error) {
	span := trace.SpanFromContext(ctx)

	objectKey := fmt.Sprintf("%s/images/%s.jpg", request.GetMediaId(), request.GetImageId())

	// 1) Remove from MinIO
	if err := s.deleteMinioObject(ctx, objectKey); err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// 2) Remove from DB/domain
	err := s.app.RemoveImage(ctx, commands.RemoveImage{
		ID: request.GetImageId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &mediapb.RemoveImageResponse{}, nil
}

func (s server) RemoveVideo(ctx context.Context, request *mediapb.RemoveVideoRequest) (*mediapb.RemoveVideoResponse, error) {
	span := trace.SpanFromContext(ctx)

	objectKey := fmt.Sprintf("%s/videos/%s.mp4", request.GetMediaId(), request.GetVideoId())

	// 1) Remove from MinIO
	if err := s.deleteMinioObject(ctx, objectKey); err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// 2) Remove from DB/domain
	err := s.app.RemoveVideo(ctx, commands.RemoveVideo{
		ID: request.GetVideoId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &mediapb.RemoveVideoResponse{}, nil
}

func (s server) StartBulkImport(ctx context.Context, request *mediapb.StartBulkImportRequest) (*mediapb.ImportSession, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	sessionID := uuid.New().String()
	span.SetAttributes(attribute.String("SessionID", sessionID))

	err := s.app.StartBulkImport(ctx, commands.StartBulkImport{
		SessionID:          sessionID,
		ExternalSystemID:   request.GetExternalSystemId(),
		ExternalSystemType: request.GetExternalSystemType(),
		EstimatedCount:     request.GetEstimatedCount(),
		Options:            request.GetOptions(),
		UserID:             userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &mediapb.ImportSession{
		Id:                 sessionID,
		ExternalSystemId:   request.GetExternalSystemId(),
		ExternalSystemType: request.GetExternalSystemType(),
		TotalImages:        request.GetEstimatedCount(),
		ProcessedImages:    0,
		FailedImages:       0,
		Status:             "pending",
		StartedAt:          time.Now().Unix(),
	}, nil
}

func (s server) AddImportBatch(ctx context.Context, request *mediapb.AddImportBatchRequest) (*mediapb.BatchResult, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("SessionID", request.GetSessionId()))

	var items []commands.ImportItem
	for _, item := range request.GetItems() {
		items = append(items, commands.ImportItem{
			ExternalID:   item.GetExternalId(),
			SKU:          item.GetSku(),
			ImageURL:     item.GetImageUrl(),
			Metadata:     item.GetMetadata(),
			DisplayOrder: item.GetDisplayOrder(),
		})
	}

	err := s.app.AddImportBatch(ctx, commands.AddImportBatch{
		SessionID: request.GetSessionId(),
		Items:     items,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &mediapb.BatchResult{
		SessionId:      request.GetSessionId(),
		AcceptedItems:  int32(len(items)),
		RejectedItems:  0,
		Errors:         nil,
	}, nil
}

func (s server) GetImportStatus(ctx context.Context, request *mediapb.GetImportStatusRequest) (*mediapb.ImportStatus, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("SessionID", request.GetSessionId()))

	status, err := s.app.GetImportStatus(ctx, queries.GetImportStatus{
		SessionID: request.GetSessionId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var recentErrors []*mediapb.ImportError
	for _, err := range status.RecentErrors {
		recentErrors = append(recentErrors, &mediapb.ImportError{
			ExternalId:   err.ExternalID,
			ErrorCode:    err.ErrorCode,
			ErrorMessage: err.ErrorMessage,
		})
	}

	return &mediapb.ImportStatus{
		SessionId:               status.SessionID,
		Status:                  status.Status,
		TotalImages:             status.TotalImages,
		ProcessedImages:         status.ProcessedImages,
		FailedImages:            status.FailedImages,
		RecentErrors:            recentErrors,
		ProgressPercentage:      status.ProgressPercentage,
		EstimatedCompletionTime: status.EstimatedCompletionTime,
	}, nil
}

func (s server) CancelImport(ctx context.Context, request *mediapb.CancelImportRequest) (*mediapb.Empty, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("SessionID", request.GetSessionId()))

	err := s.app.CancelImport(ctx, commands.CancelImport{
		SessionID: request.GetSessionId(),
		Reason:    request.GetReason(),
		UserID:    userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &mediapb.Empty{}, nil
}

// -----------------------------------------------------------------------------
// Domain -> Proto converters
// -----------------------------------------------------------------------------

func (s server) imageFromDomain(p *domain.MiddlemanImage) *mediapb.Image {
	return &mediapb.Image{
		Id:           p.ID,
		MediaId:      p.MediaID,
		DisplayOrder: int32(p.DisplayOrder),
		IsMain:       p.IsMain,
		Url:          p.URL,
		Metadata:     p.Metadata,
		Thumbnail:    p.Thumbnail,
		UserId:       p.UserID,
	}
}

func (s server) videoFromDomain(p *domain.MiddlemanVideo) *mediapb.Video {
	return &mediapb.Video{
		Id:           p.ID,
		MediaId:      p.MediaID,
		DisplayOrder: int32(p.DisplayOrder),
		IsMain:       p.IsMain,
		Url:          p.Url,
		Metadata:     p.Metadata,
		Thumbnail:    p.Thumbnail,
		UserId:       p.UserID,
	}
}

func (s server) mediaFromDomain(p *domain.MiddlemanMedia) *mediapb.Media {
	// Gather the domain's map into a slice
	medias := make([]*mediapb.MediaOrder, 0, len(p.MediaOrder))

	for _, item := range p.MediaOrder {
		medias = append(medias, &mediapb.MediaOrder{
			MediaItemId: item.MediaItemID,
			Url:         item.URL,
		})
	}

	return &mediapb.Media{
		Id:         p.ID,
		ItemId:     p.ItemID,
		ItemType:   string(p.ItemType),
		UserId:     p.UserID,
		FileType:   "",
		MediaOrder: medias,
	}
}
