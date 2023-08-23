package subscription

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"auth/core/ports"
	"auth/internal"

	"go.uber.org/fx"
)

type Service struct {
	log         *otelzap.Logger
	config      *internal.AppConfig
	ctx         context.Context
	repo        ports.SubscriptionRepository
	uuidService ports.UUIDService
}

func NewSubscriptionService(ctx context.Context, log *otelzap.Logger, config *internal.AppConfig, repo ports.SubscriptionRepository, uuidService ports.UUIDService) (ports.SubscriptionService, error) {
	return &Service{
		log:         log,
		ctx:         ctx,
		config:      config,
		repo:        repo,
		uuidService: uuidService,
	}, nil
}

var ServiceModule = fx.Provide(NewSubscriptionService)
