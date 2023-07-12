package subscription

import (
	"context"
	"time"

	"auth/core/domain"

	"github.com/google/uuid"
)

func (svc Service) UpdateSubscription(ctx context.Context, input *domain.Subscription) (uuid.UUID, error) {
	svc.log.Info("updating subscription")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	// update updated at
	input.UpdatedAt = time.Now()

	subscriptionResult, err := svc.repo.UpdateSubscription(ctx, input)
	if err != nil {
		return uuid.Nil, err
	}

	return subscriptionResult, nil
}
