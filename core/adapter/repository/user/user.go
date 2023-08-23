package user

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"auth/core/ports"
	"auth/internal"
	"auth/internal/db/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
)

type Repository struct {
	log    *otelzap.Logger
	ctx    context.Context
	config *internal.AppConfig
	db     *mongodb.MongoDb
	*mongo.Collection
}

func NewUserRepository(log *otelzap.Logger, ctx context.Context, config *internal.AppConfig, db *mongodb.MongoDb) (ports.UserRepository, error) {
	return &Repository{
		db:         db,
		log:        log,
		ctx:        ctx,
		config:     config,
		Collection: db.Db.Collection("users"),
	}, nil
}

var RepositoryModule = fx.Provide(NewUserRepository)
