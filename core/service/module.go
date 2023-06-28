package service

import (
	"auth/core/service/app-validator"
	"auth/core/service/auth"
	"auth/core/service/auth-server"
	"auth/core/service/subscription"
	"auth/core/service/user"
	"auth/core/service/uuid"

	"go.uber.org/fx"
)

// Module - inject the application services
var Module = fx.Options(
	user.ServiceModule,
	auth.ServiceModule,
	auth_server.ServiceModule,
	app_validator.ServiceModule,
	subscription.ServiceModule,
	uuid.ServiceModule,
)
