package subscription

import (
	"context"

	"auth/core/domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (repo Repository) FindSubscriptionByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	repo.log.Debug("FindSubscriptionByID")

	subscription := domain.Subscription{}
	err := repo.db.Db.Collection("subscriptions").FindOne(ctx, bson.D{{"_id", id}}).Decode(&subscription)
	if err != nil {
		repo.log.Error("error finding subscription", zap.Error(err))
		return nil, err
	}

	return &subscription, nil
}
