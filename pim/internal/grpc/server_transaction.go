package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/internal/di"
	"middleman/products/internal/application"
	"middleman/products/internal/constants"
	"middleman/products/productspb"
)

type serverTx struct {
	c di.Container
	productspb.UnimplementedProductsServiceServer
}

var _ productspb.ProductsServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	productspb.RegisterProductsServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}
func (s serverTx) AddProduct(ctx context.Context, request *productspb.AddProductRequest) (resp *productspb.AddProductResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddProduct(ctx, request)
}

func (s serverTx) RebrandProduct(ctx context.Context, request *productspb.RebrandProductRequest) (resp *productspb.RebrandProductResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RebrandProduct(ctx, request)
}

// func (s serverTx) IncreaseProductPrice(
//
//	ctx context.Context,
//	request *productspb.IncreaseProductPriceRequest,
//
// ) (resp *productspb.IncreaseProductPriceResponse, err error) {
//
//		ctx = s.c.Scoped(ctx)
//		defer func(tx *sql.Tx) {
//			err = s.closeTx(tx, err)
//		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//		next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//		return next.IncreaseProductPrice(ctx, request)
//	}
//
// func (s serverTx) DecreaseProductPrice(
//
//	ctx context.Context,
//	request *productspb.DecreaseProductPriceRequest,
//
// ) (resp *productspb.DecreaseProductPriceResponse, err error) {
//
//		ctx = s.c.Scoped(ctx)
//		defer func(tx *sql.Tx) {
//			err = s.closeTx(tx, err)
//		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//		next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//		return next.DecreaseProductPrice(ctx, request)
//	}
func (s serverTx) UpdateProduct(
	ctx context.Context,
	request *productspb.UpdateProductRequest,
) (resp *productspb.UpdateProductResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateProduct(ctx, request)
}
func (s serverTx) AddProductThumbnail(
	ctx context.Context,
	request *productspb.AddProductThumbnailRequest,
) (resp *productspb.AddProductThumbnailResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AddProductThumbnail(ctx, request)
}
func (s serverTx) UpdateProductThumbnail(
	ctx context.Context,
	request *productspb.UpdateProductThumbnailRequest,
) (resp *productspb.UpdateProductThumbnailResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateProductThumbnail(ctx, request)
}

func (s serverTx) UpdateProductPrice(
	ctx context.Context,
	request *productspb.UpdateProductPriceRequest,
) (resp *productspb.UpdateProductPriceResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateProductPrice(ctx, request)
}

//func (s serverTx) IncreaseProductPrice(ctx context.Context, request *productspb.IncreaseProductPriceRequest) (resp *productspb.IncreaseProductPriceResponse, err error) {
//	ctx = s.c.Scoped(ctx)
//	defer func(tx *sql.Tx) {
//		err = s.closeTx(tx, err)
//	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//
//	return next.IncreaseProductPrice(ctx, request)
//}
//
//func (s serverTx) DecreaseProductPrice(ctx context.Context, request *productspb.DecreaseProductPriceRequest) (resp *productspb.DecreaseProductPriceResponse, err error) {
//	ctx = s.c.Scoped(ctx)
//	defer func(tx *sql.Tx) {
//		err = s.closeTx(tx, err)
//	}(di.Get(ctx, "tx").(*sql.Tx))
//
//	next := server{app: di.Get(ctx, "app").(application.App)}
//
//	return next.DecreaseProductPrice(ctx, request)
//}

func (s serverTx) RemoveProduct(ctx context.Context, request *productspb.RemoveProductRequest) (resp *productspb.RemoveProductResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.RemoveProduct(ctx, request)
}

func (s serverTx) GetProduct(ctx context.Context, request *productspb.GetProductRequest) (resp *productspb.GetProductResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetProduct(ctx, request)
}

func (s serverTx) GetCatalog(ctx context.Context, request *productspb.GetCatalogRequest) (resp *productspb.GetCatalogResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetCatalog(ctx, request)
}
func (s serverTx) GetPublicCatalog(ctx context.Context, request *productspb.GetPublicCatalogRequest) (resp *productspb.GetPublicCatalogResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetPublicCatalog(ctx, request)
}
func (s serverTx) GetProducts(ctx context.Context, request *productspb.GetProductsRequest) (resp *productspb.GetProductsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetProducts(ctx, request)
}
func (s serverTx) GetProductsWithFilters(ctx context.Context, request *productspb.GetProductsWithFiltersRequest) (resp *productspb.GetProductsWithFiltersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetProductsWithFilters(ctx, request)
}
func (s serverTx) GetProductsByCategory(ctx context.Context, request *productspb.GetProductsByCategoryRequest) (resp *productspb.GetProductsByCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetProductsByCategory(ctx, request)
}
func (s serverTx) GetProductsByCategorySlug(ctx context.Context, request *productspb.GetProductsByCategorySlugRequest) (resp *productspb.GetProductsByCategorySlugResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetProductsByCategorySlug(ctx, request)
}
func (s serverTx) MarkProductPawned(
	ctx context.Context,
	request *productspb.MarkProductPawnedRequest,
) (resp *productspb.MarkProductPawnedResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.MarkProductPawned(ctx, request)
}
func (s serverTx) AdjustProductStock(
	ctx context.Context,
	request *productspb.AdjustProductStockRequest,
) (resp *productspb.AdjustProductStockResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AdjustProductStock(ctx, request)
}

func (s serverTx) ArchiveProduct(
	ctx context.Context,
	request *productspb.ArchiveProductRequest,
) (resp *productspb.ArchiveProductResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ArchiveProduct(ctx, request)
}

func (s serverTx) MarkProductSold(
	ctx context.Context,
	request *productspb.MarkProductSoldRequest,
) (resp *productspb.MarkProductSoldResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.MarkProductSold(ctx, request)
}

func (s serverTx) MarkProductLeased(
	ctx context.Context,
	request *productspb.MarkProductLeasedRequest,
) (resp *productspb.MarkProductLeasedResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.MarkProductLeased(ctx, request)
}

func (s serverTx) AddVariant(
	ctx context.Context,
	request *productspb.AddVariantRequest,
) (resp *productspb.AddVariantResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AddVariant(ctx, request)
}

func (s serverTx) GetVariant(
	ctx context.Context,
	request *productspb.GetVariantRequest,
) (resp *productspb.GetVariantResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.GetVariant(ctx, request)
}

func (s serverTx) GetVariants(
	ctx context.Context,
	request *productspb.GetVariantsRequest,
) (resp *productspb.GetVariantsResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.GetVariants(ctx, request)
}

func (s serverTx) IncreaseVariantPrice(
	ctx context.Context,
	request *productspb.IncreaseVariantPriceRequest,
) (resp *productspb.IncreaseVariantPriceResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.IncreaseVariantPrice(ctx, request)
}

func (s serverTx) DecreaseVariantPrice(
	ctx context.Context,
	request *productspb.DecreaseVariantPriceRequest,
) (resp *productspb.DecreaseVariantPriceResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.DecreaseVariantPrice(ctx, request)
}

func (s serverTx) AdjustVariantStock(
	ctx context.Context,
	request *productspb.AdjustVariantStockRequest,
) (resp *productspb.AdjustVariantStockResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.AdjustVariantStock(ctx, request)
}

func (s serverTx) ArchiveVariant(
	ctx context.Context,
	request *productspb.ArchiveVariantRequest,
) (resp *productspb.ArchiveVariantResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.ArchiveVariant(ctx, request)
}

func (s serverTx) RemoveVariant(
	ctx context.Context,
	request *productspb.RemoveVariantRequest,
) (resp *productspb.RemoveVariantResponse, err error) {

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
