package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/internal/di"
	"middleman/posts/internal/application"
	"middleman/posts/internal/constants"
	"middleman/posts/postspb"
)

type serverTx struct {
	c di.Container
	postspb.UnimplementedPostsServiceServer
}

var _ postspb.PostsServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	postspb.RegisterPostsServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}
func (s serverTx) AddPost(ctx context.Context, request *postspb.AddPostRequest) (resp *postspb.AddPostResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddPost(ctx, request)
}

func (s serverTx) UpdatePost(
	ctx context.Context,
	request *postspb.UpdatePostRequest,
) (resp *postspb.UpdatePostResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdatePost(ctx, request)
}

func (s serverTx) AddPostThumbnail(
	ctx context.Context,
	request *postspb.AddPostThumbnailRequest,
) (resp *postspb.AddPostThumbnailResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AddPostThumbnail(ctx, request)
}
func (s serverTx) UpdatePostThumbnail(
	ctx context.Context,
	request *postspb.UpdatePostThumbnailRequest,
) (resp *postspb.UpdatePostThumbnailResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdatePostThumbnail(ctx, request)
}
func (s serverTx) RemovePost(ctx context.Context, request *postspb.RemovePostRequest) (resp *postspb.RemovePostResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.RemovePost(ctx, request)
}

func (s serverTx) GetPost(ctx context.Context, request *postspb.GetPostRequest) (resp *postspb.GetPostResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetPost(ctx, request)
}

func (s serverTx) GetPosts(ctx context.Context, request *postspb.GetPostsRequest) (resp *postspb.GetPostsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetPosts(ctx, request)
}
func (s serverTx) GetPostsWithFilters(ctx context.Context, request *postspb.GetPostsWithFiltersRequest) (resp *postspb.GetPostsWithFiltersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetPostsWithFilters(ctx, request)
}
func (s serverTx) GetPostsByCategory(ctx context.Context, request *postspb.GetPostsByCategoryRequest) (resp *postspb.GetPostsByCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetPostsByCategory(ctx, request)
}
func (s serverTx) GetPostsByCategorySlug(ctx context.Context, request *postspb.GetPostsByCategorySlugRequest) (resp *postspb.GetPostsByCategorySlugResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetPostsByCategorySlug(ctx, request)
}
func (s serverTx) GetUserPosts(ctx context.Context, request *postspb.GetUserPostsRequest) (resp *postspb.GetUserPostsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetUserPosts(ctx, request)
}
func (s serverTx) GetPublicCatalog(ctx context.Context, request *postspb.GetPublicCatalogRequest) (resp *postspb.GetPublicCatalogResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetPublicCatalog(ctx, request)
}
func (s serverTx) ArchivePost(
	ctx context.Context,
	request *postspb.ArchivePostRequest,
) (resp *postspb.ArchivePostResponse, err error) {

	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ArchivePost(ctx, request)
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
