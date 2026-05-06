package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/categories/categoriespb"
	"middleman/categories/internal/application"
	"middleman/categories/internal/constants"
	"middleman/internal/di"
)

type serverTx struct {
	c di.Container
	categoriespb.UnimplementedCategoriesServiceServer
}

var _ categoriespb.CategoriesServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	categoriespb.RegisterCategoriesServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}
func (s serverTx) AddCategory(ctx context.Context, request *categoriespb.AddCategoryRequest) (resp *categoriespb.AddCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddCategory(ctx, request)
}

func (s serverTx) RebrandCategory(ctx context.Context, request *categoriespb.RebrandCategoryRequest) (resp *categoriespb.RebrandCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RebrandCategory(ctx, request)
}

func (s serverTx) UpdateCategory(
	ctx context.Context,
	request *categoriespb.UpdateCategoryRequest,
) (resp *categoriespb.UpdateCategoryResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateCategory(ctx, request)
}

func (s serverTx) RemoveCategory(ctx context.Context, request *categoriespb.RemoveCategoryRequest) (resp *categoriespb.RemoveCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.RemoveCategory(ctx, request)
}

func (s serverTx) GetCategory(ctx context.Context, request *categoriespb.GetCategoryRequest) (resp *categoriespb.GetCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetCategory(ctx, request)
}
func (s serverTx) GetCategoryBySlug(ctx context.Context, request *categoriespb.GetCategoryBySlugRequest) (resp *categoriespb.GetCategoryBySlugResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetCategoryBySlug(ctx, request)
}
func (s serverTx) GetCategories(ctx context.Context, request *categoriespb.GetCategoriesRequest) (resp *categoriespb.GetCategoriesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetCategories(ctx, request)
}

func (s serverTx) GetMainCategories(ctx context.Context, request *categoriespb.GetMainCategoriesRequest) (resp *categoriespb.GetMainCategoriesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetMainCategories(ctx, request)
}
func (s serverTx) GetAllMainCategories(ctx context.Context, request *categoriespb.GetAllMainCategoriesRequest) (resp *categoriespb.GetAllMainCategoriesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetAllMainCategories(ctx, request)
}

func (s serverTx) GetSubCategories(ctx context.Context, request *categoriespb.GetSubCategoriesRequest) (resp *categoriespb.GetSubCategoriesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetSubCategories(ctx, request)
}

func (s serverTx) ArchiveCategory(
	ctx context.Context,
	request *categoriespb.ArchiveCategoryRequest,
) (resp *categoriespb.ArchiveCategoryResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ArchiveCategory(ctx, request)
}

func (s serverTx) AddFilter(
	ctx context.Context,
	request *categoriespb.AddFilterRequest,
) (resp *categoriespb.AddFilterResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AddFilter(ctx, request)
}

func (s serverTx) GetFilter(
	ctx context.Context,
	request *categoriespb.GetFilterRequest,
) (resp *categoriespb.GetFilterResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.GetFilter(ctx, request)
}

func (s serverTx) GetFilters(
	ctx context.Context,
	request *categoriespb.GetFiltersRequest,
) (resp *categoriespb.GetFiltersResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.GetFilters(ctx, request)
}

func (s serverTx) ArchiveFilter(
	ctx context.Context,
	request *categoriespb.ArchiveFilterRequest,
) (resp *categoriespb.ArchiveFilterResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.ArchiveFilter(ctx, request)
}

func (s serverTx) RemoveFilter(
	ctx context.Context,
	request *categoriespb.RemoveFilterRequest,
) (resp *categoriespb.RemoveFilterResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}
	return next.RemoveFilter(ctx, request)
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
