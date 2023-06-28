package auth

import (
	"context"

	"auth/core/ports"
	"auth/internal"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service struct {
	log              *zap.Logger
	config           *internal.AppConfig
	ctx              context.Context
	userPort         ports.UserService
	authServerPort   ports.AuthServerService
	subscriptionPort ports.SubscriptionService
}

func NewAuthService(log *zap.Logger, config *internal.AppConfig, ctx context.Context, authServerPort ports.AuthServerService, userServicePort ports.UserService, subscriptionPort ports.SubscriptionService) ports.AuthService {
	return &Service{
		log:              log,
		ctx:              ctx,
		config:           config,
		userPort:         userServicePort,
		authServerPort:   authServerPort,
		subscriptionPort: subscriptionPort,
	}
}

var ServiceModule = fx.Provide(NewAuthService)
