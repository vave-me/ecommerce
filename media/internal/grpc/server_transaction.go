package grpc

import (
	"context"
	"database/sql"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"middleman/internal/di"
	"middleman/media/internal/application"
	"middleman/media/internal/constants"
	"middleman/media/mediapb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	mediapb.UnimplementedMediaServiceServer
	s3Client *s3.Client
	bucket   string
}

var _ mediapb.MediaServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar,
) error {
	mediapb.RegisterMediaServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateMedia(ctx context.Context, request *mediapb.CreateMediaRequest) (resp *mediapb.CreateMediaResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateMedia(ctx, request)
}

func (s serverTx) UpdateMedia(ctx context.Context, request *mediapb.UpdateMediaRequest) (resp *mediapb.UpdateMediaResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.UpdateMedia(ctx, request)
}
func (s serverTx) AddImage(ctx context.Context, request *mediapb.AddImageRequest) (resp *mediapb.AddImageResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddImage(ctx, request)
}

func (s serverTx) AddVideo(ctx context.Context, request *mediapb.AddVideoRequest) (resp *mediapb.AddVideoResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddVideo(ctx, request)
}

func (s serverTx) GetAllMediaImages(ctx context.Context, request *mediapb.GetAllMediaImagesRequest) (resp *mediapb.GetAllMediaImagesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetAllMediaImages(ctx, request)
}
func (s serverTx) GetAllItemImages(ctx context.Context, request *mediapb.GetAllItemImagesRequest) (resp *mediapb.GetAllItemImagesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetAllItemImages(ctx, request)
}
func (s serverTx) GetAllItemVideos(ctx context.Context, request *mediapb.GetAllItemVideosRequest) (resp *mediapb.GetAllItemVideosResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetAllItemVideos(ctx, request)
}
func (s serverTx) GetAllVideos(ctx context.Context, request *mediapb.GetAllVideosRequest) (resp *mediapb.GetAllVideosResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetAllVideos(ctx, request)
}
func (s serverTx) GetAllVideosByAuthor(ctx context.Context, request *mediapb.GetAllVideosByAuthorRequest) (resp *mediapb.GetAllVideosByAuthorResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetAllVideosByAuthor(ctx, request)
}
func (s serverTx) GetAllImagesByAuthor(ctx context.Context, request *mediapb.GetAllImagesByAuthorRequest) (resp *mediapb.GetAllImagesByAuthorResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetAllImagesByAuthor(ctx, request)
}
func (s serverTx) GetAllMediaVideos(ctx context.Context, request *mediapb.GetAllMediaVideosRequest) (resp *mediapb.GetAllMediaVideosResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetAllMediaVideos(ctx, request)
}

func (s serverTx) GetMedia(ctx context.Context, request *mediapb.GetMediaRequest) (resp *mediapb.GetMediaResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetMedia(ctx, request)
}
func (s serverTx) GetMediaByItem(ctx context.Context, request *mediapb.GetMediaByItemRequest) (resp *mediapb.GetMediaByItemResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetMediaByItem(ctx, request)
}
func (s serverTx) RemoveImage(ctx context.Context, request *mediapb.RemoveImageRequest) (resp *mediapb.RemoveImageResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RemoveImage(ctx, request)
}
func (s serverTx) RemoveVideo(ctx context.Context, request *mediapb.RemoveVideoRequest) (resp *mediapb.RemoveVideoResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RemoveVideo(ctx, request)
}
func (s serverTx) RemoveMedia(ctx context.Context, request *mediapb.RemoveMediaRequest) (resp *mediapb.RemoveMediaResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RemoveMedia(ctx, request)
}

func (s serverTx) StartBulkImport(ctx context.Context, request *mediapb.StartBulkImportRequest) (resp *mediapb.ImportSession, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.StartBulkImport(ctx, request)
}

func (s serverTx) AddImportBatch(ctx context.Context, request *mediapb.AddImportBatchRequest) (resp *mediapb.BatchResult, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddImportBatch(ctx, request)
}

func (s serverTx) GetImportStatus(ctx context.Context, request *mediapb.GetImportStatusRequest) (resp *mediapb.ImportStatus, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetImportStatus(ctx, request)
}

func (s serverTx) CancelImport(ctx context.Context, request *mediapb.CancelImportRequest) (resp *mediapb.Empty, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CancelImport(ctx, request)
}

func (s serverTx) closeTx(tx *sql.Tx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}
