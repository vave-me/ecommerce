package grpc

import (
	"context"
	"database/sql"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"middleman/internal/di"
	"middleman/mailer/internal/application"
	"middleman/mailer/internal/constants"
	"middleman/mailer/mailerpb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	mailerpb.UnimplementedMailerServiceServer
	s3Client *s3.Client
	bucket   string
}

var _ mailerpb.MailerServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar,
) error {
	mailerpb.RegisterMailerServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateEmail(ctx context.Context, request *mailerpb.CreateEmailRequest) (resp *mailerpb.CreateEmailResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateEmail(ctx, request)
}

func (s serverTx) NotifyUserCreated(ctx context.Context, request *mailerpb.NotifyUserCreatedRequest) (resp *mailerpb.NotifyUserCreatedResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.NotifyUserCreated(ctx, request)
}

//func (s serverTx) RemoveImage(ctx context.Context, request *mailerpb.RemoveImageRequest) (resp *mailerpb.RemoveImageResponse, err error) {
//	ctx = s.c.Scoped(ctx)
//	defer func(tx *sql.Tx) {
//		err = s.closeTx(tx, err)
//	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
//
//	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
//
//	return next.RemoveImage(ctx, request)
//}

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
