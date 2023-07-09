package health

import (
	"auth/core/ports"
	"auth/internal/db/mongodb"
	"auth/internal/db/redis"

	"go.uber.org/fx"
)

type Service struct {
	redis *redis.RedisDb
	mongo *mongodb.MongoDb
}

// NewHealthService returns a new HealthService.
func NewHealthService(redis *redis.RedisDb, mongo *mongodb.MongoDb) (ports.HealthService, error) {
	return &Service{
		redis: redis,
		mongo: mongo,
	}, nil
}

var ServiceModule = fx.Provide(NewHealthService)
