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

func (auth *Service) Login(email string, password string) (*domain.AuthToken, error) {
	auth.log.Debug("auth.service Login")

	ctx, cancel := context.WithTimeout(auth.ctx, time.Duration(auth.config.ServiceTimeout)*time.Second)
	defer cancel()

	keycloakServer := auth.config.KeycloakServer
	keycloakRealm := auth.config.KeycloakRealm
	keycloakClientId := auth.config.KeycloakClientId
	keycloakClientSecret := auth.config.KeycloakClientSecret

	client := gocloak.NewClient(keycloakServer)
	resToken, err := client.Login(ctx, keycloakClientId, keycloakClientSecret, keycloakRealm, email, password)
	if err != nil {
		auth.log.Error("login error: ", zap.Error(err))

		if strings.HasPrefix(err.Error(), "401") {
			return nil, system.ErrInvalidCredentials
		}
		if strings.HasPrefix(err.Error(), "400") {
			return nil, system.ErrInvalidInput
		}

		return nil, err
	}

	accessToken := resToken.AccessToken
	refreshToken := resToken.RefreshToken
	tokenType := resToken.TokenType
	expiredIn := resToken.ExpiresIn
	idToken := resToken.IDToken

	token := domain.NewAuthToken(accessToken, refreshToken, tokenType, expiredIn, idToken)
	return token, nil
}
