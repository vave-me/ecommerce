package grpc

import (
	"context"
	"database/sql"
	"middleman/internal/di"
	"middleman/search/internal/application"
	"middleman/search/internal/constants"
	"middleman/search/searchpb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	searchpb.UnimplementedSearchServiceServer
}

var _ searchpb.SearchServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	searchpb.RegisterSearchServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) SearchOrders(ctx context.Context, request *searchpb.SearchOrdersRequest) (resp *searchpb.SearchOrdersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchOrders(ctx, request)
}
func (s serverTx) SearchProductsWithFilters(ctx context.Context, request *searchpb.SearchProductsWithFiltersRequest) (resp *searchpb.SearchProductsWithFiltersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchProductsWithFilters(ctx, request)
}
func (s serverTx) SearchProductsWithCategorySlug(ctx context.Context, request *searchpb.SearchProductsWithCategorySlugRequest) (resp *searchpb.SearchProductsWithCategorySlugResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchProductsWithCategorySlug(ctx, request)
}
func (s serverTx) SearchPostsWithFilters(ctx context.Context, request *searchpb.SearchPostsWithFiltersRequest) (resp *searchpb.SearchPostsWithFiltersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchPostsWithFilters(ctx, request)
}
func (s serverTx) SearchProductsWithTerm(ctx context.Context, request *searchpb.SearchProductsWithTermRequest) (resp *searchpb.SearchProductsWithTermResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchProductsWithTerm(ctx, request)
}
func (s serverTx) SearchPostsWithTerm(ctx context.Context, request *searchpb.SearchPostsWithTermRequest) (resp *searchpb.SearchPostsWithTermResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchPostsWithTerm(ctx, request)
}

func (s serverTx) GetOrder(ctx context.Context, request *searchpb.GetOrderRequest) (resp *searchpb.GetOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetOrder(ctx, request)
}
func (s serverTx) GetProduct(ctx context.Context, request *searchpb.GetProductRequest) (resp *searchpb.GetProductResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetProduct(ctx, request)
}

func (s serverTx) SuggestProducts(ctx context.Context, request *searchpb.SuggestProductsRequest) (resp *searchpb.SuggestProductsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SuggestProducts(ctx, request)
}
func (s serverTx) SuggestPosts(ctx context.Context, request *searchpb.SuggestPostsRequest) (resp *searchpb.SuggestPostsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SuggestPosts(ctx, request)
}

func (s serverTx) UnifiedSearch(ctx context.Context, request *searchpb.UnifiedSearchRequest) (resp *searchpb.UnifiedSearchResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.UnifiedSearch(ctx, request)
}

func (s serverTx) UnifiedFeed(ctx context.Context, request *searchpb.UnifiedFeedRequest) (resp *searchpb.UnifiedFeedResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.UnifiedFeed(ctx, request)
}

func (s serverTx) GetPost(ctx context.Context, request *searchpb.GetPostRequest) (resp *searchpb.GetPostResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetPost(ctx, request)
}

func (s serverTx) SearchPostsWithCategorySlug(ctx context.Context, request *searchpb.SearchPostsWithCategorySlugRequest) (resp *searchpb.SearchPostsWithCategorySlugResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchPostsWithCategorySlug(ctx, request)
}

func (s serverTx) SearchPostsWithCategory(ctx context.Context, request *searchpb.SearchPostsWithCategoryRequest) (resp *searchpb.SearchPostsWithCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchPostsWithCategory(ctx, request)
}

func (s serverTx) GetCatalog(ctx context.Context, request *searchpb.GetCatalogRequest) (resp *searchpb.UnifiedFeedResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetCatalog(ctx, request)
}

func (s serverTx) SearchUsers(ctx context.Context, request *searchpb.SearchUsersRequest) (resp *searchpb.SearchUsersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchUsers(ctx, request)
}

func (s serverTx) GetService(ctx context.Context, request *searchpb.GetServiceRequest) (resp *searchpb.GetServiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetService(ctx, request)
}

func (s serverTx) SearchServicesWithFilters(ctx context.Context, request *searchpb.SearchServicesWithFiltersRequest) (resp *searchpb.SearchServicesWithFiltersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchServicesWithFilters(ctx, request)
}

func (s serverTx) SearchServicesWithTerm(ctx context.Context, request *searchpb.SearchServicesWithTermRequest) (resp *searchpb.SearchServicesWithTermResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchServicesWithTerm(ctx, request)
}

func (s serverTx) SearchServicesWithCategorySlug(ctx context.Context, request *searchpb.SearchServicesWithCategorySlugRequest) (resp *searchpb.SearchServicesWithCategorySlugResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchServicesWithCategorySlug(ctx, request)
}

func (s serverTx) SearchServicesWithCategory(ctx context.Context, request *searchpb.SearchServicesWithCategoryRequest) (resp *searchpb.SearchServicesWithCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SearchServicesWithCategory(ctx, request)
}

func (s serverTx) SuggestServices(ctx context.Context, request *searchpb.SuggestServicesRequest) (resp *searchpb.SuggestServicesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.SuggestServices(ctx, request)
}

func (s serverTx) GetUser(ctx context.Context, request *searchpb.GetUserRequest) (resp *searchpb.GetUserResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.Application)}

	return next.GetUser(ctx, request)
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
