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

func NewRedisDb(lifecycle fx.Lifecycle, ctx context.Context, config *internal.AppConfig, log *zap.Logger) (*redis.Client, error) {
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

	return client, nil
}

var RedisModule = NewRedisDb
