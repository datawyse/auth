package subscription

import (
	"context"
	"time"

	"auth/core/domain"

	"github.com/google/uuid"
)

func (svc Service) CreateSubscription(ctx context.Context, input *domain.Subscription) (uuid.UUID, error) {
	svc.log.Info("Creating subscription")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	return svc.repo.CreateSubscription(ctx, input)
}
