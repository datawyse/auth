package subscription

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (repo Repository) FindSubscriptionByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	repo.log.Info("finding subscription by id")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(repo.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(repo.config.ServiceName).Start(ctx, "repository.subscription.find_subscription_by_id")
	defer span.End()

	subscription := domain.Subscription{}
	err := repo.db.Db.Collection("subscriptions").FindOne(ctx, bson.D{{"_id", id}}).Decode(&subscription)
	if err != nil {
		repo.log.Error("error finding subscription", zap.Error(err))
		return nil, err
	}

	return &subscription, nil
}
