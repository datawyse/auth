package uuid

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain/system"
	"auth/core/ports"
	"auth/internal"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type Service struct {
	log    *otelzap.Logger
	config *internal.AppConfig
	ctx    context.Context
}

func NewUUIDService(log *otelzap.Logger, config *internal.AppConfig, ctx context.Context) (ports.UUIDService, error) {
	return &Service{
		log:    log,
		ctx:    ctx,
		config: config,
	}, nil
}

func (svc *Service) FromString(ctx context.Context, id string) (uuid.UUID, error) {
	svc.log.Info("uuid from string")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.uuid.from_string")
	defer span.End()

	return system.ToUUID(id)
}

func (svc *Service) NewUUID(ctx context.Context) uuid.UUID {
	svc.log.Info("new uuid")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.uuid.new_uuid")
	defer span.End()

	return system.NewUUID()
}

func (svc *Service) IsValidUUID(ctx context.Context, id string) (bool, error) {
	svc.log.Info("is valid uuid")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.uuid.is_valid_uuid")
	defer span.End()

	return system.IsValidUUID(id), nil
}

var ServiceModule = fx.Provide(NewUUIDService)
