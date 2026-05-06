package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/internal/di"
	"middleman/tickets/internal/application"
	"middleman/tickets/internal/constants"
	"middleman/tickets/ticketspb"
)

type serverTx struct {
	c di.Container
	ticketspb.UnimplementedTicketsServiceServer
}

var _ ticketspb.TicketsServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	ticketspb.RegisterTicketsServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}
func (s serverTx) AddTicket(ctx context.Context, request *ticketspb.AddTicketRequest) (resp *ticketspb.AddTicketResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddTicket(ctx, request)
}

func (s serverTx) RebrandTicket(ctx context.Context, request *ticketspb.RebrandTicketRequest) (resp *ticketspb.RebrandTicketResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RebrandTicket(ctx, request)
}

// func (s serverTx) IncreaseTicketPrice(
//
//	ctx context.Context,
//	request *ticketspb.IncreaseTicketPriceRequest,
//
// ) (resp *ticketspb.IncreaseTicketPriceResponse, err error) {
//
//		ctx = s.c.Scoped(ctx)
//		defer func(tx *sql.Tx) {
//			err = s.closeTx(tx, err)
//		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//		next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//		return next.IncreaseTicketPrice(ctx, request)
//	}
//
// func (s serverTx) DecreaseTicketPrice(
//
//	ctx context.Context,
//	request *ticketspb.DecreaseTicketPriceRequest,
//
// ) (resp *ticketspb.DecreaseTicketPriceResponse, err error) {
//
//		ctx = s.c.Scoped(ctx)
//		defer func(tx *sql.Tx) {
//			err = s.closeTx(tx, err)
//		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//		next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//		return next.DecreaseTicketPrice(ctx, request)
//	}
func (s serverTx) UpdateTicket(
	ctx context.Context,
	request *ticketspb.UpdateTicketRequest,
) (resp *ticketspb.UpdateTicketResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateTicket(ctx, request)
}
func (s serverTx) AddTicketThumbnail(
	ctx context.Context,
	request *ticketspb.AddTicketThumbnailRequest,
) (resp *ticketspb.AddTicketThumbnailResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AddTicketThumbnail(ctx, request)
}
func (s serverTx) UpdateTicketThumbnail(
	ctx context.Context,
	request *ticketspb.UpdateTicketThumbnailRequest,
) (resp *ticketspb.UpdateTicketThumbnailResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateTicketThumbnail(ctx, request)
}

func (s serverTx) UpdateTicketPrice(
	ctx context.Context,
	request *ticketspb.UpdateTicketPriceRequest,
) (resp *ticketspb.UpdateTicketPriceResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateTicketPrice(ctx, request)
}

//func (s serverTx) IncreaseTicketPrice(ctx context.Context, request *ticketspb.IncreaseTicketPriceRequest) (resp *ticketspb.IncreaseTicketPriceResponse, err error) {
//	ctx = s.c.Scoped(ctx)
//	defer func(tx *sql.Tx) {
//		err = s.closeTx(tx, err)
//	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//
//	return next.IncreaseTicketPrice(ctx, request)
//}
//
//func (s serverTx) DecreaseTicketPrice(ctx context.Context, request *ticketspb.DecreaseTicketPriceRequest) (resp *ticketspb.DecreaseTicketPriceResponse, err error) {
//	ctx = s.c.Scoped(ctx)
//	defer func(tx *sql.Tx) {
//		err = s.closeTx(tx, err)
//	}(di.Get(ctx, "tx").(*sql.Tx))
//
//	next := server{app: di.Get(ctx, "app").(application.App)}
//
//	return next.DecreaseTicketPrice(ctx, request)
//}

func (s serverTx) RemoveTicket(ctx context.Context, request *ticketspb.RemoveTicketRequest) (resp *ticketspb.RemoveTicketResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.RemoveTicket(ctx, request)
}

func (s serverTx) GetTicket(ctx context.Context, request *ticketspb.GetTicketRequest) (resp *ticketspb.GetTicketResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetTicket(ctx, request)
}

func (s serverTx) GetCatalog(ctx context.Context, request *ticketspb.GetCatalogRequest) (resp *ticketspb.GetCatalogResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetCatalog(ctx, request)
}
func (s serverTx) GetPublicCatalog(ctx context.Context, request *ticketspb.GetPublicCatalogRequest) (resp *ticketspb.GetPublicCatalogResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetPublicCatalog(ctx, request)
}
func (s serverTx) GetTickets(ctx context.Context, request *ticketspb.GetTicketsRequest) (resp *ticketspb.GetTicketsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetTickets(ctx, request)
}
func (s serverTx) GetTicketsWithFilters(ctx context.Context, request *ticketspb.GetTicketsWithFiltersRequest) (resp *ticketspb.GetTicketsWithFiltersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetTicketsWithFilters(ctx, request)
}
func (s serverTx) GetTicketsByCategory(ctx context.Context, request *ticketspb.GetTicketsByCategoryRequest) (resp *ticketspb.GetTicketsByCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetTicketsByCategory(ctx, request)
}
func (s serverTx) GetTicketsByCategorySlug(ctx context.Context, request *ticketspb.GetTicketsByCategorySlugRequest) (resp *ticketspb.GetTicketsByCategorySlugResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetTicketsByCategorySlug(ctx, request)
}
func (s serverTx) MarkTicketPawned(
	ctx context.Context,
	request *ticketspb.MarkTicketPawnedRequest,
) (resp *ticketspb.MarkTicketPawnedResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.MarkTicketPawned(ctx, request)
}
func (s serverTx) AdjustTicketStock(
	ctx context.Context,
	request *ticketspb.AdjustTicketStockRequest,
) (resp *ticketspb.AdjustTicketStockResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AdjustTicketStock(ctx, request)
}

func (s serverTx) ArchiveTicket(
	ctx context.Context,
	request *ticketspb.ArchiveTicketRequest,
) (resp *ticketspb.ArchiveTicketResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ArchiveTicket(ctx, request)
}

func (s serverTx) MarkTicketSold(
	ctx context.Context,
	request *ticketspb.MarkTicketSoldRequest,
) (resp *ticketspb.MarkTicketSoldResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.MarkTicketSold(ctx, request)
}

func (s serverTx) MarkTicketLeased(
	ctx context.Context,
	request *ticketspb.MarkTicketLeasedRequest,
) (resp *ticketspb.MarkTicketLeasedResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.MarkTicketLeased(ctx, request)
}

func (s serverTx) AddVariant(
	ctx context.Context,
	request *ticketspb.AddVariantRequest,
) (resp *ticketspb.AddVariantResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AddVariant(ctx, request)
}

func (s serverTx) GetVariant(
	ctx context.Context,
	request *ticketspb.GetVariantRequest,
) (resp *ticketspb.GetVariantResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.GetVariant(ctx, request)
}

func (s serverTx) GetVariants(
	ctx context.Context,
	request *ticketspb.GetVariantsRequest,
) (resp *ticketspb.GetVariantsResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.GetVariants(ctx, request)
}

func (s serverTx) IncreaseVariantPrice(
	ctx context.Context,
	request *ticketspb.IncreaseVariantPriceRequest,
) (resp *ticketspb.IncreaseVariantPriceResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.IncreaseVariantPrice(ctx, request)
}

func (s serverTx) DecreaseVariantPrice(
	ctx context.Context,
	request *ticketspb.DecreaseVariantPriceRequest,
) (resp *ticketspb.DecreaseVariantPriceResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.DecreaseVariantPrice(ctx, request)
}

func (s serverTx) AdjustVariantStock(
	ctx context.Context,
	request *ticketspb.AdjustVariantStockRequest,
) (resp *ticketspb.AdjustVariantStockResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.AdjustVariantStock(ctx, request)
}

func (s serverTx) ArchiveVariant(
	ctx context.Context,
	request *ticketspb.ArchiveVariantRequest,
) (resp *ticketspb.ArchiveVariantResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.ArchiveVariant(ctx, request)
}

func (s serverTx) RemoveVariant(
	ctx context.Context,
	request *ticketspb.RemoveVariantRequest,
) (resp *ticketspb.RemoveVariantResponse, err error) {

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
