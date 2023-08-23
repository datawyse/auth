package subscription

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain"

	"github.com/google/uuid"
)

func (repo Repository) CreateSubscription(ctx context.Context, input *domain.Subscription) (uuid.UUID, error) {
	repo.log.Info("creating subscription")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(repo.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(repo.config.ServiceName).Start(ctx, "repository.subscription.create_subscription")
	defer span.End()

	subscriptionResult, err := repo.InsertOne(ctx, *input)
	if err != nil {
		return uuid.Nil, err
	}

	if subscriptionID, ok := subscriptionResult.InsertedID.(uuid.UUID); ok {
		return subscriptionID, nil
	}

	return input.Id, nil
}
