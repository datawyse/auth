package auth_server

import (
	"context"
	"strings"
	"time"

	"auth/core/domain/system"
	"auth/core/ports"
	"auth/internal"

	"github.com/Nerzal/gocloak/v13"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service struct {
	ctx    context.Context
	log    *zap.Logger
	config *internal.AppConfig
}

func NewAuthServerService(log *zap.Logger, ctx context.Context, config *internal.AppConfig) ports.AuthServerService {
	return &Service{
		ctx:    ctx,
		log:    log,
		config: config,
	}
}

func (authServer *Service) AccessToken() (string, error) {
	authServer.log.Info("getting access token")

	ctx, cancel := context.WithTimeout(authServer.ctx, time.Duration(authServer.config.ServiceTimeout)*time.Second)
	defer cancel()

	keycloakServer := authServer.config.KeycloakServer
	keycloakRealm := authServer.config.KeycloakRealm
	keycloakClientId := authServer.config.KeycloakClientId
	keycloakClientSecret := authServer.config.KeycloakClientSecret

	client := gocloak.NewClient(keycloakServer)
	token, err := client.LoginClient(ctx, keycloakClientId, keycloakClientSecret, keycloakRealm)
	if err != nil {
		authServer.log.Error("login client error: ", zap.Error(err))

		// if strings.HasPrefix(err.Error(), "401") {
		// 	return "", domain.ErrServer
		// }
		if strings.HasPrefix(err.Error(), "400") {
			return "", system.ErrInvalidInput
		}

		return "", system.ErrServer
	}

	return token.AccessToken, nil
}

var ServiceModule = fx.Provide(NewAuthServerService)
