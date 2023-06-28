package subscription

import (
	"context"

	"auth/core/ports"
	"auth/internal"
	"auth/internal/db/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Repository struct {
	log    *zap.Logger
	ctx    context.Context
	config *internal.AppConfig
	db     *mongodb.MongoDb
	*mongo.Collection
}

func NewSubscriptionRepository(ctx context.Context, log *zap.Logger, config *internal.AppConfig, db *mongodb.MongoDb) (ports.SubscriptionRepository, error) {
	return &Repository{
		db:         db,
		log:        log,
		ctx:        ctx,
		config:     config,
		Collection: db.Db.Collection("subscriptions"),
	}, nil
}

var RepositoryModule = fx.Provide(NewSubscriptionRepository)
