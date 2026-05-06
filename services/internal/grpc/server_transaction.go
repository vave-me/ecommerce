package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/internal/di"
	"middleman/services/internal/application"
	"middleman/services/internal/constants"
	"middleman/services/servicespb"
)

type serverTx struct {
	c di.Container
	servicespb.UnimplementedServicesServiceServer
}

var _ servicespb.ServicesServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	servicespb.RegisterServicesServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}
func (s serverTx) AddService(ctx context.Context, request *servicespb.AddServiceRequest) (resp *servicespb.AddServiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddService(ctx, request)
}

func (s serverTx) RebrandService(ctx context.Context, request *servicespb.RebrandServiceRequest) (resp *servicespb.RebrandServiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RebrandService(ctx, request)
}

// func (s serverTx) IncreaseServicePrice(
//
//	ctx context.Context,
//	request *servicespb.IncreaseServicePriceRequest,
//
// ) (resp *servicespb.IncreaseServicePriceResponse, err error) {
//
//		ctx = s.c.Scoped(ctx)
//		defer func(tx *sql.Tx) {
//			err = s.closeTx(tx, err)
//		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//		next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//		return next.IncreaseServicePrice(ctx, request)
//	}
//
// func (s serverTx) DecreaseServicePrice(
//
//	ctx context.Context,
//	request *servicespb.DecreaseServicePriceRequest,
//
// ) (resp *servicespb.DecreaseServicePriceResponse, err error) {
//
//		ctx = s.c.Scoped(ctx)
//		defer func(tx *sql.Tx) {
//			err = s.closeTx(tx, err)
//		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//		next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//		return next.DecreaseServicePrice(ctx, request)
//	}
func (s serverTx) UpdateService(
	ctx context.Context,
	request *servicespb.UpdateServiceRequest,
) (resp *servicespb.UpdateServiceResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateService(ctx, request)
}

func (s serverTx) UpdateServicePrice(
	ctx context.Context,
	request *servicespb.UpdateServicePriceRequest,
) (resp *servicespb.UpdateServicePriceResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateServicePrice(ctx, request)
}

//func (s serverTx) IncreaseServicePrice(ctx context.Context, request *servicespb.IncreaseServicePriceRequest) (resp *servicespb.IncreaseServicePriceResponse, err error) {
//	ctx = s.c.Scoped(ctx)
//	defer func(tx *sql.Tx) {
//		err = s.closeTx(tx, err)
//	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//
//	return next.IncreaseServicePrice(ctx, request)
//}
//
//func (s serverTx) DecreaseServicePrice(ctx context.Context, request *servicespb.DecreaseServicePriceRequest) (resp *servicespb.DecreaseServicePriceResponse, err error) {
//	ctx = s.c.Scoped(ctx)
//	defer func(tx *sql.Tx) {
//		err = s.closeTx(tx, err)
//	}(di.Get(ctx, "tx").(*sql.Tx))
//
//	next := server{app: di.Get(ctx, "app").(application.App)}
//
//	return next.DecreaseServicePrice(ctx, request)
//}

func (s serverTx) RemoveService(ctx context.Context, request *servicespb.RemoveServiceRequest) (resp *servicespb.RemoveServiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.RemoveService(ctx, request)
}

func (s serverTx) GetService(ctx context.Context, request *servicespb.GetServiceRequest) (resp *servicespb.GetServiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetService(ctx, request)
}

func (s serverTx) GetCatalog(ctx context.Context, request *servicespb.GetCatalogRequest) (resp *servicespb.GetCatalogResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetCatalog(ctx, request)
}
func (s serverTx) GetPublicCatalog(ctx context.Context, request *servicespb.GetPublicCatalogRequest) (resp *servicespb.GetPublicCatalogResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetPublicCatalog(ctx, request)
}
func (s serverTx) GetServices(ctx context.Context, request *servicespb.GetServicesRequest) (resp *servicespb.GetServicesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetServices(ctx, request)
}
func (s serverTx) GetServicesByCategory(ctx context.Context, request *servicespb.GetServicesByCategoryRequest) (resp *servicespb.GetServicesByCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetServicesByCategory(ctx, request)
}
func (s serverTx) GetServicesWithFilter(ctx context.Context, request *servicespb.GetServicesWithFilterRequest) (resp *servicespb.GetServicesWithFilterResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetServicesWithFilter(ctx, request)
}

func (s serverTx) AdjustServiceStock(
	ctx context.Context,
	request *servicespb.AdjustServiceStockRequest,
) (resp *servicespb.AdjustServiceStockResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AdjustServiceStock(ctx, request)
}

func (s serverTx) ArchiveService(
	ctx context.Context,
	request *servicespb.ArchiveServiceRequest,
) (resp *servicespb.ArchiveServiceResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ArchiveService(ctx, request)
}

func (s serverTx) MarkServiceSold(
	ctx context.Context,
	request *servicespb.MarkServiceSoldRequest,
) (resp *servicespb.MarkServiceSoldResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.MarkServiceSold(ctx, request)
}

func (s serverTx) MarkServiceLeased(
	ctx context.Context,
	request *servicespb.MarkServiceLeasedRequest,
) (resp *servicespb.MarkServiceLeasedResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.MarkServiceLeased(ctx, request)
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
