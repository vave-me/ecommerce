package rest

import (
	"context"
	"database/sql"
	"middleman/users/internal/constants"
	"middleman/users/userspb"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

func RegisterGateway(ctx context.Context, mux *chi.Mux, grpcAddr string) error {
	const apiRoot = "/api/users"

	gateway := runtime.NewServeMux(
		// Configure gateway to handle cookie metadata from gRPC responses
		runtime.WithForwardResponseOption(handleCookieMetadata),
	)
	err := userspb.RegisterUsersServiceHandlerFromEndpoint(ctx, gateway, grpcAddr, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		return err
	}

	// mount the GRPC gateway
	mux.Mount(apiRoot, gateway)

	return nil
}

// handleCookieMetadata processes gRPC metadata to set HTTP cookies
// This enables gRPC server methods to set HTTP-only cookies via metadata
func handleCookieMetadata(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// Handle set-cookie metadata
	if cookies := md.HeaderMD.Get("set-cookie"); len(cookies) > 0 {
		for _, cookie := range cookies {
			w.Header().Add("Set-Cookie", cookie)
		}
		// Clean up metadata to avoid exposing it as regular headers
		delete(md.HeaderMD, "set-cookie")
	}

	return nil
}

func TransactionUnaryServerInterceptor(db *sql.DB) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Begin a new transaction
		tx, err := db.Begin()
		if err != nil {
			return nil, err
		}
		// Add the transaction to the context
		ctx = context.WithValue(ctx, constants.DatabaseTransactionKey, tx)
		defer func() {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}()
		// Call the handler
		resp, err = handler(ctx, req)
		return resp, err
	}
}
