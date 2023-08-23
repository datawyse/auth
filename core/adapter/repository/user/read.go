package user

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (repo Repository) User(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	repo.log.Info("finding user")

	ctx, cancel := context.WithTimeout(repo.ctx, time.Duration(repo.config.DatabaseTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(repo.config.ServiceName).Start(ctx, "repository.user.user")
	defer span.End()

	user := domain.User{}
	repo.log.Info("reading user from keycloak id ", zap.String("id", id.String()))

	err := repo.db.Db.Collection("users").FindOne(ctx, bson.D{{"_id", id}}).Decode(&user)
	if err != nil {
		repo.log.Error("error finding user", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (repo Repository) Users(ctx context.Context) ([]*domain.User, error) {
	repo.log.Info("reading users")

	ctx, cancel := context.WithTimeout(repo.ctx, time.Duration(repo.config.DatabaseTimeout)*time.Second)
	defer cancel()

	var users []*domain.User
	cursor, err := repo.db.Db.Collection("users").Find(ctx, bson.D{})
	if err != nil {
		repo.log.Error("error finding users", zap.Error(err))
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			repo.log.Error("error closing cursor", zap.Error(err))
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var user domain.User
		if err = cursor.Decode(&user); err != nil {
			repo.log.Error("error decoding user", zap.Error(err))
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		repo.log.Error("error finding users", zap.Error(err))
		return nil, err
	}

	return users, nil
}
