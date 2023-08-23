package subscription

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain"

	"github.com/google/uuid"
)

func (svc Service) CreateSubscription(ctx context.Context, input *domain.Subscription) (uuid.UUID, error) {
	svc.log.Info("Creating subscription")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.subscription.create_subscription")
	defer span.End()

	return svc.repo.CreateSubscription(ctx, input)
}
