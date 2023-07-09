package auth

import (
	"context"
	"strings"
	"time"

	"auth/core/domain"
	"auth/core/domain/system"

	"github.com/Nerzal/gocloak/v13"
	"go.uber.org/zap"
)

func (svc *Service) RefreshToken(ctx context.Context, token string) (*domain.AuthToken, error) {
	svc.log.Debug("svc.service RefreshToken")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	keycloakServer := svc.config.KeycloakServer
	keycloakRealm := svc.config.KeycloakRealm
	keycloakClientId := svc.config.KeycloakClientId
	keycloakClientSecret := svc.config.KeycloakClientSecret

	client := gocloak.NewClient(keycloakServer)
	resToken, err := client.RefreshToken(ctx, token, keycloakClientId, keycloakClientSecret, keycloakRealm)
	if err != nil {
		svc.log.Error("login error: ", zap.Error(err))

		if strings.HasPrefix(err.Error(), "401") {
			return nil, system.ErrInvalidCredentials
		}
		if strings.HasPrefix(err.Error(), "400") {
			return nil, system.ErrInvalidToken
		}

		return nil, err
	}

	newToken := domain.NewAuthToken(resToken.AccessToken, resToken.RefreshToken, resToken.TokenType, resToken.ExpiresIn, resToken.IDToken)
	return newToken, nil
}
