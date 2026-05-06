package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/internal/di"
	"middleman/streams/internal/application"
	"middleman/streams/internal/constants"
	"middleman/streams/streamspb"
)

type serverTx struct {
	c di.Container
	streamspb.UnimplementedStreamsServiceServer
}

var _ streamspb.StreamsServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	streamspb.RegisterStreamsServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}
func (s serverTx) AddStream(ctx context.Context, request *streamspb.AddStreamRequest) (resp *streamspb.AddStreamResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddStream(ctx, request)
}

func (s serverTx) RebrandStream(ctx context.Context, request *streamspb.RebrandStreamRequest) (resp *streamspb.RebrandStreamResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RebrandStream(ctx, request)
}

// func (s serverTx) IncreaseStreamPrice(
//
//	ctx context.Context,
//	request *streamspb.IncreaseStreamPriceRequest,
//
// ) (resp *streamspb.IncreaseStreamPriceResponse, err error) {
//
//		ctx = s.c.Scoped(ctx)
//		defer func(tx *sql.Tx) {
//			err = s.closeTx(tx, err)
//		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//		next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//		return next.IncreaseStreamPrice(ctx, request)
//	}
//
// func (s serverTx) DecreaseStreamPrice(
//
//	ctx context.Context,
//	request *streamspb.DecreaseStreamPriceRequest,
//
// ) (resp *streamspb.DecreaseStreamPriceResponse, err error) {
//
//		ctx = s.c.Scoped(ctx)
//		defer func(tx *sql.Tx) {
//			err = s.closeTx(tx, err)
//		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//		next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//		return next.DecreaseStreamPrice(ctx, request)
//	}
func (s serverTx) UpdateStream(
	ctx context.Context,
	request *streamspb.UpdateStreamRequest,
) (resp *streamspb.UpdateStreamResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateStream(ctx, request)
}
func (s serverTx) AddStreamThumbnail(
	ctx context.Context,
	request *streamspb.AddStreamThumbnailRequest,
) (resp *streamspb.AddStreamThumbnailResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AddStreamThumbnail(ctx, request)
}
func (s serverTx) UpdateStreamThumbnail(
	ctx context.Context,
	request *streamspb.UpdateStreamThumbnailRequest,
) (resp *streamspb.UpdateStreamThumbnailResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateStreamThumbnail(ctx, request)
}

func (s serverTx) UpdateStreamPrice(
	ctx context.Context,
	request *streamspb.UpdateStreamPriceRequest,
) (resp *streamspb.UpdateStreamPriceResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateStreamPrice(ctx, request)
}

//func (s serverTx) IncreaseStreamPrice(ctx context.Context, request *streamspb.IncreaseStreamPriceRequest) (resp *streamspb.IncreaseStreamPriceResponse, err error) {
//	ctx = s.c.Scoped(ctx)
//	defer func(tx *sql.Tx) {
//		err = s.closeTx(tx, err)
//	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//
//	return next.IncreaseStreamPrice(ctx, request)
//}
//
//func (s serverTx) DecreaseStreamPrice(ctx context.Context, request *streamspb.DecreaseStreamPriceRequest) (resp *streamspb.DecreaseStreamPriceResponse, err error) {
//	ctx = s.c.Scoped(ctx)
//	defer func(tx *sql.Tx) {
//		err = s.closeTx(tx, err)
//	}(di.Get(ctx, "tx").(*sql.Tx))
//
//	next := server{app: di.Get(ctx, "app").(application.App)}
//
//	return next.DecreaseStreamPrice(ctx, request)
//}

func (s serverTx) RemoveStream(ctx context.Context, request *streamspb.RemoveStreamRequest) (resp *streamspb.RemoveStreamResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.RemoveStream(ctx, request)
}

func (s serverTx) GetStream(ctx context.Context, request *streamspb.GetStreamRequest) (resp *streamspb.GetStreamResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetStream(ctx, request)
}

func (s serverTx) GetCatalog(ctx context.Context, request *streamspb.GetCatalogRequest) (resp *streamspb.GetCatalogResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetCatalog(ctx, request)
}
func (s serverTx) GetPublicCatalog(ctx context.Context, request *streamspb.GetPublicCatalogRequest) (resp *streamspb.GetPublicCatalogResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetPublicCatalog(ctx, request)
}
func (s serverTx) GetStreams(ctx context.Context, request *streamspb.GetStreamsRequest) (resp *streamspb.GetStreamsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetStreams(ctx, request)
}
func (s serverTx) GetStreamsWithFilters(ctx context.Context, request *streamspb.GetStreamsWithFiltersRequest) (resp *streamspb.GetStreamsWithFiltersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetStreamsWithFilters(ctx, request)
}
func (s serverTx) GetStreamsByCategory(ctx context.Context, request *streamspb.GetStreamsByCategoryRequest) (resp *streamspb.GetStreamsByCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetStreamsByCategory(ctx, request)
}
func (s serverTx) GetStreamsByCategorySlug(ctx context.Context, request *streamspb.GetStreamsByCategorySlugRequest) (resp *streamspb.GetStreamsByCategorySlugResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetStreamsByCategorySlug(ctx, request)
}
func (s serverTx) MarkStreamPawned(
	ctx context.Context,
	request *streamspb.MarkStreamPawnedRequest,
) (resp *streamspb.MarkStreamPawnedResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.MarkStreamPawned(ctx, request)
}
func (s serverTx) AdjustStreamStock(
	ctx context.Context,
	request *streamspb.AdjustStreamStockRequest,
) (resp *streamspb.AdjustStreamStockResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AdjustStreamStock(ctx, request)
}

func (s serverTx) ArchiveStream(
	ctx context.Context,
	request *streamspb.ArchiveStreamRequest,
) (resp *streamspb.ArchiveStreamResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ArchiveStream(ctx, request)
}

func (s serverTx) MarkStreamSold(
	ctx context.Context,
	request *streamspb.MarkStreamSoldRequest,
) (resp *streamspb.MarkStreamSoldResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.MarkStreamSold(ctx, request)
}

func (s serverTx) MarkStreamLeased(
	ctx context.Context,
	request *streamspb.MarkStreamLeasedRequest,
) (resp *streamspb.MarkStreamLeasedResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.MarkStreamLeased(ctx, request)
}

func (s serverTx) AddVariant(
	ctx context.Context,
	request *streamspb.AddVariantRequest,
) (resp *streamspb.AddVariantResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AddVariant(ctx, request)
}

func (s serverTx) GetVariant(
	ctx context.Context,
	request *streamspb.GetVariantRequest,
) (resp *streamspb.GetVariantResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.GetVariant(ctx, request)
}

func (s serverTx) GetVariants(
	ctx context.Context,
	request *streamspb.GetVariantsRequest,
) (resp *streamspb.GetVariantsResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.GetVariants(ctx, request)
}

func (s serverTx) IncreaseVariantPrice(
	ctx context.Context,
	request *streamspb.IncreaseVariantPriceRequest,
) (resp *streamspb.IncreaseVariantPriceResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.IncreaseVariantPrice(ctx, request)
}

func (s serverTx) DecreaseVariantPrice(
	ctx context.Context,
	request *streamspb.DecreaseVariantPriceRequest,
) (resp *streamspb.DecreaseVariantPriceResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.DecreaseVariantPrice(ctx, request)
}

func (s serverTx) AdjustVariantStock(
	ctx context.Context,
	request *streamspb.AdjustVariantStockRequest,
) (resp *streamspb.AdjustVariantStockResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.AdjustVariantStock(ctx, request)
}

func (s serverTx) ArchiveVariant(
	ctx context.Context,
	request *streamspb.ArchiveVariantRequest,
) (resp *streamspb.ArchiveVariantResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.ArchiveVariant(ctx, request)
}

func (s serverTx) RemoveVariant(
	ctx context.Context,
	request *streamspb.RemoveVariantRequest,
) (resp *streamspb.RemoveVariantResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.RemoveVariant(ctx, request)
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
