package user

import (
	"context"
	"time"

	"auth/core/domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (repo Repository) User(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	repo.log.Info("reading user")

	ctx, cancel := context.WithTimeout(repo.ctx, time.Duration(repo.config.DatabaseTimeout)*time.Second)
	defer cancel()

	user := domain.User{}
	repo.log.Info("reading user from keycloak id ", zap.String("id", id.String()))

	err := repo.db.Db.Collection("users").FindOne(ctx, bson.D{{"_id", id}}).Decode(&user)
	if err != nil {
		repo.log.Error("error finding user", zap.Error(err))
		return nil, err
	}

	return &user, nil
}
