package health

import (
	"auth/core/ports"
	"auth/internal"
	"auth/internal/db/mongodb"
	"auth/internal/db/redis"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/fx"
)

type Service struct {
	log    *otelzap.Logger
	redis  *redis.RedisDb
	mongo  *mongodb.MongoDb
	config *internal.AppConfig
}

// NewHealthService returns a new HealthService.
func NewHealthService(log *otelzap.Logger, config *internal.AppConfig, redis *redis.RedisDb, mongo *mongodb.MongoDb) (ports.HealthService, error) {
	return &Service{
		log:    log,
		redis:  redis,
		mongo:  mongo,
		config: config,
	}, nil
}

var ServiceModule = fx.Provide(NewHealthService)
