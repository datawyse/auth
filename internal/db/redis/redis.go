package redis

import (
	"context"
	"strconv"
	"time"

	"auth/internal"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RedisDb struct {
	*redis.Client
}

func NewRedisDb(lifecycle fx.Lifecycle, ctx context.Context, config *internal.AppConfig, log *zap.Logger) (*RedisDb, error) {
	db, err := strconv.Atoi(config.RedisDb)
	if err != nil {
		db = 0
	}

	log.Info("Connecting to RedisDB", zap.String("url", config.RedisURI), zap.String("database", strconv.Itoa(db)))

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisURI,
		Password: config.RedisPassword,
		DB:       db,
	})

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			redisCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if config.AppMode != "production" {
				res, err := client.Ping(redisCtx).Result()
				if err != nil {
					return err
				}

				log.Info("Successfully connected and pinged.", zap.String("result", res))
			}

			log.Info("Successfully connected to redis.")
			return err
		},
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})

	return &RedisDb{client}, nil
}

// IsHealthy returns health status
func (database *RedisDb) IsHealthy() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := database.Ping(ctx).Err(); err != nil {
		return false, err
	}

	return true, nil
}

var Module = NewRedisDb
