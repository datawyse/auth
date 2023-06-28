package db

import (
	"auth/internal/db/mongodb"
	"auth/internal/db/redis"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(mongodb.MongoModule),
	fx.Provide(redis.RedisModule),
)
