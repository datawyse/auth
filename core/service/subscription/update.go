package subscription

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain"

	"github.com/google/uuid"
)

func (svc Service) UpdateSubscription(ctx context.Context, input *domain.Subscription) (uuid.UUID, error) {
	svc.log.Info("updating subscription")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.uuid.update_subscription")
	defer span.End()

	// update updated at
	input.UpdatedAt = time.Now()

	subscriptionResult, err := svc.repo.UpdateSubscription(ctx, input)
	if err != nil {
		return uuid.Nil, err
	}

	return subscriptionResult, nil
}
