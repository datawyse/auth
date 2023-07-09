package subscription

import (
	"context"
	"time"

	"auth/core/domain"
)

func (svc Service) FindSubscriptionByID(ctx context.Context, id string) (*domain.Subscription, error) {
	svc.log.Debug("FindSubscriptionByID")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	subscriptionId, err := svc.uuidService.FromString(id)
	if err != nil {
		return nil, err
	}

	return svc.repo.FindSubscriptionByID(ctx, subscriptionId)
}
