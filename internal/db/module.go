package db

import (
	"auth/internal/db/mongodb"
	"auth/internal/db/redis"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(mongodb.Module),
	fx.Provide(redis.Module),
)
