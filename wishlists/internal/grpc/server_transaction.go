package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/internal/di"
	"middleman/wishlists/internal/application"
	"middleman/wishlists/internal/constants"
	"middleman/wishlists/wishlistspb"
)

type serverTx struct {
	c di.Container
	wishlistspb.UnimplementedWishlistServiceServer
}

var _ wishlistspb.WishlistServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	wishlistspb.RegisterWishlistServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateWishlist(ctx context.Context, request *wishlistspb.CreateWishlistRequest) (resp *wishlistspb.CreateWishlistResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateWishlist(ctx, request)
}
func (s serverTx) RemoveWishlist(ctx context.Context, request *wishlistspb.RemoveWishlistRequest) (resp *wishlistspb.RemoveWishlistResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RemoveWishlist(ctx, request)
}

func (s serverTx) AddWishlistItem(ctx context.Context, request *wishlistspb.AddWishlistItemRequest) (resp *wishlistspb.AddWishlistItemResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddWishlistItem(ctx, request)
}
func (s serverTx) RemoveWishlistItem(ctx context.Context, request *wishlistspb.RemoveWishlistItemRequest) (resp *wishlistspb.RemoveWishlistItemResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RemoveWishlistItem(ctx, request)
}

func (s serverTx) GetWishlist(ctx context.Context, request *wishlistspb.GetWishlistRequest) (resp *wishlistspb.GetWishlistResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetWishlist(ctx, request)
}

func (s serverTx) GetWishlistItem(ctx context.Context, request *wishlistspb.GetWishlistItemRequest) (resp *wishlistspb.GetWishlistItemResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetWishlistItem(ctx, request)
}

func (s serverTx) GetWishlistItems(ctx context.Context, request *wishlistspb.GetWishlistItemsRequest) (resp *wishlistspb.GetWishlistItemsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetWishlistItems(ctx, request)
}
func (s serverTx) GetWishlists(ctx context.Context, request *wishlistspb.GetWishlistsRequest) (resp *wishlistspb.GetWishlistsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetWishlists(ctx, request)
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
