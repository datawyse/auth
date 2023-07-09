package user

import (
	"context"
	"time"

	"auth/core/domain"
	"auth/core/domain/system"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (repo Repository) CreateUser(ctx context.Context, input *domain.User) (uuid.UUID, error) {
	repo.log.Info("creating user")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(repo.config.DatabaseTimeout)*time.Second)
	defer cancel()

	userCollection := repo.db.Db.Collection("users")

	indexId, err := userCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"verified": 1},
		Options: options.Index().SetUnique(false),
	})
	if err != nil {
		repo.log.Error("error creating index", zap.Error(err))
		return uuid.Nil, err
	}

	repo.log.Info("created index", zap.Any("indexId", indexId))
	
	_, err = repo.db.Db.Collection("users").InsertOne(ctx, *input)
	if err != nil {
		repo.log.Error("error inserting user", zap.Error(err))

		if mongo.IsDuplicateKeyError(err) {
			repo.log.Error("user already exists", zap.Error(err))
			return uuid.Nil, system.ErrUserAlreadyExists
		}

		return uuid.Nil, err
	}

	return input.Id, nil
}
