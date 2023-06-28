package uuid

import (
	"context"

	"auth/core/domain/system"
	"auth/core/ports"
	"auth/internal"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service struct {
	log    *zap.Logger
	config *internal.AppConfig
	ctx    context.Context
}

func NewUUIDService(log *zap.Logger, config *internal.AppConfig, ctx context.Context) ports.UUIDService {
	return &Service{
		log:    log,
		ctx:    ctx,
		config: config,
	}
}

func (svc *Service) FromString(id string) (uuid.UUID, error) {
	return system.ToUUID(id)
}

func (svc *Service) NewUUID() (uuid.UUID, error) {
	return system.NewUUID(), nil
}

func (svc *Service) IsValidUUID(id string) (bool, error) {
	return system.IsValidUUID(id), nil
}

var ServiceModule = fx.Provide(NewUUIDService)