package grpc

import (
	"context"
	"database/sql"
	"log"
	"middleman/internal/di"
	"middleman/users/internal/application"
	"middleman/users/internal/constants"
	"middleman/users/userspb"

	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	userspb.UsersServiceServer
}

var _ userspb.UsersServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	userspb.RegisterUsersServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateUser(ctx context.Context, request *userspb.CreateUserRequest) (resp *userspb.CreateUserResponse, err error) {
	log.Println("serverTx.CreateUser called")

	ctx = s.c.Scoped(ctx)
	log.Println("Context scoped")

	defer func(tx *sql.Tx) {
		if p := recover(); p != nil {
			log.Printf("Panic recovered in defer: %v", p)
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Printf("Error detected: %v. Rolling back transaction.", err)
			_ = tx.Rollback()
		} else {
			log.Println("Committing transaction.")
			err = tx.Commit()
			if err != nil {
				log.Printf("Error committing transaction: %v", err)
			}
		}
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	log.Println("server instance created, calling next.CreateUser")

	// Correctly assign resp and err
	resp, err = next.CreateUser(ctx, request)
	log.Printf("next.CreateUser returned, resp: %v, err: %v", resp, err)

	return resp, err
}
func (s serverTx) AddAdmin(ctx context.Context, request *userspb.AddAdminRequest) (resp *userspb.AddAdminResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.AddAdmin(ctx, request)
}
func (s serverTx) LoginUser(ctx context.Context, request *userspb.LoginUserRequest) (resp *userspb.LoginUserResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.LoginUser(ctx, request)
}
func (s serverTx) LogUserOut(ctx context.Context, request *userspb.LogUserOutRequest) (resp *userspb.LogUserOutResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.LogoutUser(ctx, request)
}
func (s serverTx) ForgotPassword(ctx context.Context, request *userspb.ForgotPasswordRequest) (resp *userspb.ForgotPasswordResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ForgotPassword(ctx, request)
}

func (s serverTx) ResetPassword(ctx context.Context, request *userspb.ResetPasswordRequest) (resp *userspb.ResetPasswordResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ResetPassword(ctx, request)
}

func (s serverTx) MobileLoginWithGoogle(ctx context.Context, request *userspb.MobileLoginWithGoogleRequest) (resp *userspb.LoginUserResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.MobileLoginWithGoogle(ctx, request)
}
func (s serverTx) WebLoginWithGoogle(ctx context.Context, request *userspb.WebLoginWithGoogleRequest) (resp *userspb.LoginUserResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.WebLoginWithGoogle(ctx, request)
}
func (s serverTx) EnableUser(ctx context.Context, request *userspb.EnableUserRequest) (resp *userspb.EnableUserResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))
	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.EnableUser(ctx, request)
}

func (s serverTx) DisableUser(ctx context.Context, request *userspb.DisableUserRequest) (resp *userspb.DisableUserResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.DisableUser(ctx, request)
}
func (s serverTx) RenameUser(ctx context.Context, request *userspb.RenameUserRequest) (resp *userspb.RenameUserResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RenameUser(ctx, request)
}

func (s serverTx) GetUser(ctx context.Context, request *userspb.GetUserRequest) (resp *userspb.GetUserResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetUser(ctx, request)
}
func (s serverTx) GetBaseUser(ctx context.Context, request *userspb.GetBaseUserRequest) (resp *userspb.GetBaseUserResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetBaseUser(ctx, request)
}
func (s serverTx) GetUsers(ctx context.Context, request *userspb.GetUsersRequest) (resp *userspb.GetUsersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetUsers(ctx, request)
}

func (s serverTx) ListEnabledUsers(ctx context.Context, request *userspb.ListEnabledUsersRequest) (resp *userspb.ListEnabledUsersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ListEnabledUsers(ctx, request)
}

//func (s serverTx) IncreaseProductPrice(ctx context.Context, request *userspb.IncreaseProductPriceRequest) (resp *userspb.IncreaseProductPriceResponse, err error) {
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
//func (s serverTx) DecreaseProductPrice(ctx context.Context, request *userspb.DecreaseProductPriceRequest) (resp *userspb.DecreaseProductPriceResponse, err error) {
//	ctx = s.c.Scoped(ctx)
//	defer func(tx *sql.Tx) {
//		err = s.closeTx(tx, err)
//	}(di.Get(ctx, "tx").(*sql.Tx))
//
//	next := server{app: di.Get(ctx, "app").(application.App)}
//
//	return next.DecreaseProductPrice(ctx, request)
//}

func (s serverTx) RefreshAuthToken(ctx context.Context, request *userspb.RefreshAuthTokenRequest) (resp *userspb.RefreshAuthTokenResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RefreshAuthToken(ctx, request)
}

func (s serverTx) ClearTokens(ctx context.Context, request *userspb.ClearTokensRequest) (resp *userspb.ClearTokensResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.ClearTokens(ctx, request)
}

func (s serverTx) UpdateUser(ctx context.Context, request *userspb.UpdateUserRequest) (resp *userspb.UpdateUserResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.UpdateUser(ctx, request)
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
