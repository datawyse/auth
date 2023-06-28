package subscription

import (
	"context"

	"auth/core/ports"
	"auth/internal"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service struct {
	log         *zap.Logger
	config      *internal.AppConfig
	ctx         context.Context
	repo        ports.SubscriptionRepository
	uuidService ports.UUIDService
}

func NewSubscriptionService(ctx context.Context, log *zap.Logger, config *internal.AppConfig, repo ports.SubscriptionRepository, uuidService ports.UUIDService) (ports.SubscriptionService, error) {
	return &Service{
		log:         log,
		ctx:         ctx,
		config:      config,
		repo:        repo,
		uuidService: uuidService,
	}, nil
}

var ServiceModule = fx.Provide(NewSubscriptionService)
