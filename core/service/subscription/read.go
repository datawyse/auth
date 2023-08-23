package subscription

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain"
)

func (svc Service) FindSubscriptionByID(ctx context.Context, id string) (*domain.Subscription, error) {
	svc.log.Debug("finding subscription by id")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.uuid.find_subscription_by_id")
	defer span.End()

	subscriptionId, err := svc.uuidService.FromString(ctx, id)
	if err != nil {
		return nil, err
	}

	return svc.repo.FindSubscriptionByID(ctx, subscriptionId)
}
