package health

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain"
)

// ReadHealth returns the health.
func (svc *Service) ReadHealth(ctx context.Context) (*domain.Health, error) {
	svc.log.Info("checking health")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.uuid.new_uuid")
	defer span.End()

	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05")

	timezone := time.Now().Format("UTC")
	datetime := currentDate + " " + currentTime + " " + timezone

	isMongoHealthy, err := svc.mongo.IsHealthy()
	if err != nil {
		return nil, err
	}

	var mongoStatus = "ERROR"
	if isMongoHealthy {
		mongoStatus = "OK"
	}

	isRedisHealthy, err := svc.redis.IsHealthy()
	if err != nil {
		return nil, err
	}

	var redisStatus = "ERROR"
	if isRedisHealthy {
		redisStatus = "OK"
	}

	return &domain.Health{
		Status:   "OK",
		Message:  "System is healthy",
		Version:  "v1.0.0",
		Date:     currentDate,
		Datetime: datetime,
		Time:     currentTime,
		Timezone: timezone,
		Revision: "1",
		Redis:    mongoStatus,
		Mongo:    redisStatus,
	}, nil
}
