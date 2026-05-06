package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/following/followingpb"
	"middleman/following/internal/application"
	"middleman/following/internal/constants"
	"middleman/internal/di"
)

type serverTx struct {
	c di.Container
	followingpb.UnimplementedFollowingServiceServer
}

var _ followingpb.FollowingServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	followingpb.RegisterFollowingServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) AddFollow(ctx context.Context, request *followingpb.AddFollowRequest) (resp *followingpb.AddFollowResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddFollow(ctx, request)
}

func (s serverTx) ApproveFollow(ctx context.Context, request *followingpb.ApproveFollowRequest) (resp *followingpb.ApproveFollowResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ApproveFollow(ctx, request)
}
func (s serverTx) FlagFollow(ctx context.Context, request *followingpb.FlagFollowRequest) (resp *followingpb.FlagFollowResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.FlagFollow(ctx, request)
}
func (s serverTx) RemoveFollow(ctx context.Context, request *followingpb.RemoveFollowRequest) (resp *followingpb.RemoveFollowResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RemoveFollow(ctx, request)
}
func (s serverTx) RejectFollow(ctx context.Context, request *followingpb.RejectFollowRequest) (resp *followingpb.RejectFollowResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RejectFollow(ctx, request)
}
func (s serverTx) GetFollow(ctx context.Context, request *followingpb.GetFollowRequest) (resp *followingpb.GetFollowResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetFollow(ctx, request)
}
func (s serverTx) GetFollowing(ctx context.Context, request *followingpb.GetFollowingRequest) (resp *followingpb.GetFollowingResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetFollowing(ctx, request)
}

func (s serverTx) GetMostFollowedItems(ctx context.Context, request *followingpb.GetMostFollowedRequest) (resp *followingpb.GetMostFollowedResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetMostFollowed(ctx, request)
}

func (s serverTx) GetMostFollowedItemsByCategory(ctx context.Context, request *followingpb.GetMostFollowedByCategoryRequest) (resp *followingpb.GetMostFollowedByCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetMostFollowedByCategory(ctx, request)
}

func (s serverTx) GetFollowingBySender(ctx context.Context, request *followingpb.GetFollowingBySenderRequest) (resp *followingpb.GetFollowingBySenderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetFollowingBySender(ctx, request)
}
func (s serverTx) GetApprovedFollowing(ctx context.Context, request *followingpb.GetApprovedFollowingRequest) (resp *followingpb.GetApprovedFollowingResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetApprovedFollowing(ctx, request)
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
