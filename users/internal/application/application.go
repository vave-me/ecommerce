package application

import (
	"context"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/users/internal/application/commands"
	"middleman/users/internal/application/queries"
	"middleman/users/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		AuthorizeUser(ctx context.Context, cmd commands.AuthorizeUser) error
		CreateUser(ctx context.Context, cmd commands.CreateUser) error
		AddAdmin(ctx context.Context, cmd commands.AddAdmin) error
		UpdateUser(ctx context.Context, cmd commands.UpdateUser) error
		CreateUserFromGoogle(ctx context.Context, cmd commands.CreateUserFromGoogle) error
		EnableUser(ctx context.Context, cmd commands.EnableUser) error
		DisableUser(ctx context.Context, cmd commands.DisableUser) error
		RenameUser(ctx context.Context, cmd commands.RenameUser) error
		LoginUser(ctx context.Context, cmd commands.LoginUser) (string, string, string, error)
		MobileLoginWithGoogle(ctx context.Context, cmd commands.MobileLoginWithGoogle) (string, string, string, error)
		WebLoginWithGoogle(ctx context.Context, cmd commands.WebLoginWithGoogle) (string, string, string, error)
		LogoutUser(ctx context.Context, cmd commands.LogoutUser) error
		ClearTokens(ctx context.Context, cmd commands.ClearTokens) error
		KycVerify(ctx context.Context, cmd commands.KycVerify) error
		ForgotPassword(ctx context.Context, cmd commands.ForgotPassword) error
		ResetPassword(ctx context.Context, cmd commands.ResetPassword) error
		RefreshAuthToken(ctx context.Context, cmd commands.RefreshAuthToken) (string, string, error)
	}

	Queries interface {
		GetUser(ctx context.Context, query queries.GetUser) (*domain.MiddlemanUser, error)
		GetBaseUser(ctx context.Context, query queries.GetBaseUser) (*domain.MiddlemanViewUser, error)
		GetUserByMail(ctx context.Context, query queries.GetUserByMail) (*domain.MiddlemanUser, error)
		GetUserByGoogleID(ctx context.Context, query queries.GetUserByGoogleID) (*domain.MiddlemanUser, error)
		GetUsers(ctx context.Context, query queries.GetUsers) ([]*domain.MiddlemanUser, error)
		GetEnabledUsers(ctx context.Context, query queries.GetEnabledUsers) ([]*domain.MiddlemanUser, error)
	}

	Application struct {
		auth *auth.Auth
		appCommands
		appQueries
	}

	appCommands struct {
		commands.AuthorizeUserHandler
		commands.CreateUserHandler
		commands.AddAdminHandler
		commands.UpdateUserHandler
		commands.CreateUserFromGoogleHandler
		commands.EnableUserHandler
		commands.DisableUserHandler
		commands.RenameUserHandler
		commands.LoginUserHandler
		commands.WebLoginWithGoogleHandler
		commands.MobileLoginWithGoogleHandler
		commands.LogoutUserHandler
		commands.ClearTokensHandler
		commands.KycVerifyHandler
		commands.ForgotPasswordHandler
		commands.ResetPasswordHandler
		commands.RefreshAuthTokenHandler
	}

	appQueries struct {
		queries.GetUserHandler
		queries.GetBaseUserHandler
		queries.GetUserByMailHandler
		queries.GetUserByGoogleIDHandler
		queries.GetUsersHandler
		queries.GetEnabledUsersHandler
	}
)

var _ App = (*Application)(nil)

func New(users domain.UserRepository,
	middleman domain.MiddlemanRepository, publisher ddd.EventPublisher[ddd.Event], auth *auth.Auth) *Application {
	return &Application{
		auth: auth,
		appCommands: appCommands{
			AuthorizeUserHandler:         commands.NewAuthorizeUserHandler(users, publisher),
			CreateUserHandler:            commands.NewCreateUserHandler(users, middleman, auth, publisher),
			AddAdminHandler:              commands.NewAddAdminHandler(users, middleman, auth, publisher),
			UpdateUserHandler:            commands.NewUpdateUserHandler(users, publisher),
			CreateUserFromGoogleHandler:  commands.NewCreateUserFromGoogleHandler(users, middleman, publisher),
			EnableUserHandler:   commands.NewEnableUserHandler(users, publisher, auth),
			DisableUserHandler:  commands.NewDisableUserHandler(users, publisher),
			RenameUserHandler:            commands.NewRenameUserHandler(users, publisher),
			LogoutUserHandler:            commands.NewLogoutUserHandler(users, auth, publisher),
			ClearTokensHandler:           commands.NewClearTokensHandler(users, auth, publisher),
			LoginUserHandler:             commands.NewLoginUserHandler(middleman, users, auth, publisher),
			WebLoginWithGoogleHandler:    commands.NewWebLoginWithGoogleHandler(middleman, users, auth, publisher),
			MobileLoginWithGoogleHandler: commands.NewMobileLoginWithGoogleHandler(middleman, users, auth, publisher),
			KycVerifyHandler:             commands.NewKycVerifyHandler(users, publisher),
			ForgotPasswordHandler:        commands.NewForgotPasswordHandler(users, middleman, publisher, auth),
			ResetPasswordHandler:         commands.NewResetPasswordHandler(users, publisher, auth),
			RefreshAuthTokenHandler:      commands.NewRefreshAuthTokenHandler(users, middleman, auth, publisher),
		},
		appQueries: appQueries{
			GetUserHandler:               queries.NewGetUserHandler(middleman),
			GetBaseUserHandler:           queries.NewGetBaseUserHandler(middleman),
			GetUserByMailHandler:         queries.NewGetUserByMailHandler(middleman),
			GetUserByGoogleIDHandler:     queries.NewGetUserByGoogleIDHandler(middleman),
			GetUsersHandler:              queries.NewGetUsersHandler(middleman),
			GetEnabledUsersHandler: queries.NewGetEnabledUsersHandler(middleman),
		},
	}
}
