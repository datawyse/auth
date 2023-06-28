package subscription

import (
	"context"

	"auth/core/domain"

	"github.com/google/uuid"
)

func (repo Repository) CreateSubscription(ctx context.Context, input *domain.Subscription) (uuid.UUID, error) {
	repo.log.Info("creating subscription")

	subscriptionResult, err := repo.InsertOne(ctx, *input)
	if err != nil {
		return uuid.Nil, err
	}

	if subscriptionID, ok := subscriptionResult.InsertedID.(uuid.UUID); ok {
		return subscriptionID, nil
	}

	return input.Id, nil
}
