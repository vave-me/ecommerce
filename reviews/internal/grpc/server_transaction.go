package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"middleman/internal/di"
	"middleman/reviews/internal/application"
	"middleman/reviews/internal/constants"
	"middleman/reviews/reviewspb"
)

type serverTx struct {
	c di.Container
	reviewspb.UnimplementedReviewsServiceServer
}

var _ reviewspb.ReviewsServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	reviewspb.RegisterReviewsServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) AddReview(ctx context.Context, request *reviewspb.AddReviewRequest) (resp *reviewspb.AddReviewResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddReview(ctx, request)
}
func (s serverTx) EditReview(ctx context.Context, request *reviewspb.EditReviewRequest) (resp *reviewspb.EditReviewResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.EditReview(ctx, request)
}
func (s serverTx) ApproveReview(ctx context.Context, request *reviewspb.ApproveReviewRequest) (resp *reviewspb.ApproveReviewResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ApproveReview(ctx, request)
}
func (s serverTx) FlagReview(ctx context.Context, request *reviewspb.FlagReviewRequest) (resp *reviewspb.FlagReviewResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.FlagReview(ctx, request)
}
func (s serverTx) RemoveReview(ctx context.Context, request *reviewspb.RemoveReviewRequest) (resp *reviewspb.RemoveReviewResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RemoveReview(ctx, request)
}
func (s serverTx) RejectReview(ctx context.Context, request *reviewspb.RejectReviewRequest) (resp *reviewspb.RejectReviewResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RejectReview(ctx, request)
}
func (s serverTx) GetReview(ctx context.Context, request *reviewspb.GetReviewRequest) (resp *reviewspb.GetReviewResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetReview(ctx, request)
}
func (s serverTx) GetReviews(ctx context.Context, request *reviewspb.GetReviewsRequest) (resp *reviewspb.GetReviewsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetReviews(ctx, request)
}

func (s serverTx) GetMostReviewedItems(ctx context.Context, request *reviewspb.GetMostReviewedRequest) (resp *reviewspb.GetMostReviewedResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetMostReviewed(ctx, request)
}

func (s serverTx) GetMostReviewedItemsByCategory(ctx context.Context, request *reviewspb.GetMostReviewedByCategoryRequest) (resp *reviewspb.GetMostReviewedByCategoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetMostReviewedByCategory(ctx, request)
}

func (s serverTx) GetReviewsBySender(ctx context.Context, request *reviewspb.GetReviewsBySenderRequest) (resp *reviewspb.GetReviewsBySenderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetReviewsBySender(ctx, request)
}
func (s serverTx) GetApprovedReviews(ctx context.Context, request *reviewspb.GetApprovedReviewsRequest) (resp *reviewspb.GetApprovedReviewsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetApprovedReviews(ctx, request)
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
