package auth

import (
	"auth/core/ports"
	"auth/internal"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/fx"
)

type Service struct {
	log              *otelzap.Logger
	config           *internal.AppConfig
	userPort         ports.UserService
	uuidPort         ports.UUIDService
	authServerPort   ports.AuthServerService
	subscriptionPort ports.SubscriptionService
}

func NewAuthService(log *otelzap.Logger, config *internal.AppConfig, authServerPort ports.AuthServerService, userPort ports.UserService, subscriptionPort ports.SubscriptionService, uuidPort ports.UUIDService) (ports.AuthService, error) {
	return &Service{
		log:              log,
		config:           config,
		uuidPort:         uuidPort,
		userPort:         userPort,
		authServerPort:   authServerPort,
		subscriptionPort: subscriptionPort,
	}, nil
}

var ServiceModule = fx.Provide(NewAuthService)
