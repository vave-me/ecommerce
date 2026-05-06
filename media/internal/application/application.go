package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/media/internal/application/commands"
	"middleman/media/internal/application/queries"
	"middleman/media/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		CreateMedia(ctx context.Context, cmd commands.CreateMedia) error
		UpdateMedia(ctx context.Context, cmd commands.UpdateMedia) error
		AddImage(ctx context.Context, cmd commands.AddImage) error
		AddVideo(ctx context.Context, cmd commands.AddVideo) error
		RemoveMedia(ctx context.Context, cmd commands.RemoveMedia) error
		RemoveImage(ctx context.Context, cmd commands.RemoveImage) error
		RemoveVideo(ctx context.Context, cmd commands.RemoveVideo) error
		StartBulkImport(ctx context.Context, cmd commands.StartBulkImport) error
		AddImportBatch(ctx context.Context, cmd commands.AddImportBatch) error
		CancelImport(ctx context.Context, cmd commands.CancelImport) error
	}
	Queries interface {
		GetMedia(ctx context.Context, query queries.GetMedia) (*domain.MiddlemanMedia, error)
		GetMediaByItem(ctx context.Context, query queries.GetMediaByItem) (*domain.MiddlemanMedia, error)
		GetAllItemImages(ctx context.Context, query queries.GetAllItemImages) ([]*domain.MiddlemanImage, error)
		GetAllItemVideos(ctx context.Context, query queries.GetAllItemVideos) ([]*domain.MiddlemanVideo, error)
		GetAllMediaImages(ctx context.Context, query queries.GetAllMediaImages) ([]*domain.MiddlemanImage, error)
		GetAllMediaVideos(ctx context.Context, query queries.GetAllMediaVideos) ([]*domain.MiddlemanVideo, error)
		GetAllVideos(ctx context.Context, query queries.GetAllVideos) ([]*domain.MiddlemanVideo, int64, error)
		GetAllVideosByAuthor(ctx context.Context, query queries.GetAllVideosByAuthor) ([]*domain.MiddlemanVideo, int64, error)
		GetAllImagesByAuthor(ctx context.Context, query queries.GetAllImagesByAuthor) ([]*domain.MiddlemanImage, int64, error)
		GetImportStatus(ctx context.Context, query queries.GetImportStatus) (*queries.ImportStatus, error)
	}

	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.CreateMediaHandler
		commands.UpdateMediaHandler
		commands.AddImageHandler
		commands.AddVideoHandler
		commands.RemoveMediaHandler
		commands.RemoveImageHandler
		commands.RemoveVideoHandler
		commands.StartBulkImportHandler
		commands.AddImportBatchHandler
		commands.CancelImportHandler
	}
	appQueries struct {
		queries.GetMediaHandler
		queries.GetMediaByItemHandler
		queries.GetAllItemImagesHandler
		queries.GetAllItemVideosHandler
		queries.GetAllMediaImagesHandler
		queries.GetAllMediaVideosHandler
		queries.GetAllVideosHandler
		queries.GetAllVideosByAuthorHandler
		queries.GetAllImagesByAuthorHandler
		queries.GetImportStatusHandler
	}
)

var _ App = (*Application)(nil)

func New(media domain.MediaRepository, videos domain.VideoRepository, images domain.ImageRepository, catalogMedia domain.MiddlemanMediaRepository, catalogImage domain.MiddlemanImageRepository, catalogVideo domain.MiddlemanVideoRepository, importers domain.ImporterRepository, importSessions domain.ImportSessionRepository, importItems domain.ImportItemRepository, productRepo commands.ProductRepository, publisher ddd.EventPublisher[ddd.Event]) *Application {

	return &Application{
		appCommands: appCommands{
			CreateMediaHandler:      commands.NewCreateMediaHandler(media, publisher),
			UpdateMediaHandler:      commands.NewUpdateMediaHandler(media, publisher),
			AddImageHandler:         commands.NewAddImageHandler(images, publisher),
			AddVideoHandler:         commands.NewAddVideoHandler(videos, publisher),
			RemoveMediaHandler:      commands.NewRemoveMediaHandler(media, publisher),
			RemoveImageHandler:      commands.NewRemoveImageHandler(images, publisher),
			RemoveVideoHandler:      commands.NewRemoveVideoHandler(videos, publisher),
			StartBulkImportHandler:  commands.NewStartBulkImportHandler(importers, publisher),
			AddImportBatchHandler:   commands.NewAddImportBatchHandler(importSessions, importItems, productRepo, publisher),
			CancelImportHandler:     commands.NewCancelImportHandler(importSessions, publisher)},
		appQueries: appQueries{
			GetMediaHandler:             queries.NewGetMediaHandler(catalogMedia),
			GetMediaByItemHandler:       queries.NewGetMediaByItemHandler(catalogMedia),
			GetAllItemImagesHandler:     queries.NewGetAllItemImagesHandler(catalogImage),
			GetAllItemVideosHandler:     queries.NewGetAllItemVideosHandler(catalogVideo),
			GetAllMediaImagesHandler:    queries.NewGetAllMediaImagesHandler(catalogImage),
			GetAllMediaVideosHandler:    queries.NewGetAllMediaVideosHandler(catalogVideo),
			GetAllVideosHandler:         queries.NewGetAllVideosHandler(catalogVideo),
			GetAllVideosByAuthorHandler: queries.NewGetAllVideosByAuthorHandler(catalogVideo),
			GetAllImagesByAuthorHandler: queries.NewGetAllImagesByAuthorHandler(catalogImage),
			GetImportStatusHandler:      queries.NewGetImportStatusHandler(importSessions, importItems),
		},
	}
}
