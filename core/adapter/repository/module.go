package repository

import (
	"auth/core/adapter/repository/subscription"
	"auth/core/adapter/repository/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	subscription.RepositoryModule,
	user.RepositoryModule,
)
