package grpc

import (
	"database/sql"
	"middleman/internal/di"
	"middleman/newsletters/newsletterspb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	newsletterspb.UnimplementedNewslettersServiceServer
}

var _ newsletterspb.NewslettersServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	newsletterspb.RegisterNewslettersServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

// Add transaction-wrapped methods here as needed
// Currently, the newsletter service doesn't require any special transaction handling

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