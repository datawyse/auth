package ports

import (
	"context"

	"auth/core/domain"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, input *domain.Subscription) (uuid.UUID, error)
	FindSubscriptionByID(ctx context.Context, id string) (*domain.Subscription, error)
	UpdateSubscription(ctx context.Context, input *domain.Subscription) (uuid.UUID, error)
}

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, input *domain.Subscription) (uuid.UUID, error)
	FindSubscriptionByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	UpdateSubscription(ctx context.Context, input *domain.Subscription) (uuid.UUID, error)
}
