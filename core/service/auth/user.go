package auth

import (
	"context"
	"time"

	"auth/core/domain"
	"auth/core/domain/http"
	"auth/core/domain/system"

	"github.com/Nerzal/gocloak/v13"
	"go.uber.org/zap"
)

func (auth *Service) User(userId string) (*domain.UserInfo, error) {
	auth.log.Info("auth.service.user")

	return auth.userPort.User(userId)
}

func (auth *Service) UpdateUser(ctx context.Context, input *http.UpdateUserInput) (*domain.User, error) {
	auth.log.Info("auth.service.user")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(auth.config.ServiceTimeout)*time.Second)
	defer cancel()

	user := &domain.User{}

	userId, err := system.ToUUID(input.Id)
	if err != nil {
		return nil, system.ErrInvalidInput
	}
	user.Id = userId

	for _, orgId := range input.Organizations {
		id, err := system.ToUUID(orgId)
		if err != nil {
			return nil, system.ErrInvalidInput
		}
		user.Organizations = append(user.Organizations, id)
	}

	for _, roleId := range input.Roles {
		id, err := system.ToUUID(roleId)
		if err != nil {
			return nil, system.ErrInvalidInput
		}
		user.Roles = append(user.Roles, id)
	}
	if input.Language != "" {
		user.Language = input.Language
	}
	if input.AccountType != "" {
		user.AccountType = input.AccountType
	}

	token, err := auth.authServerPort.AccessToken()
	if err != nil {
		auth.log.Error("access token error: ", zap.Error(err))
		return nil, err
	}

	keycloakServer := auth.config.KeycloakServer
	keycloakRealm := auth.config.KeycloakRealm

	client := gocloak.NewClient(keycloakServer)

	gocloakUser := gocloak.User{
		ID: gocloak.StringP(user.Id.String()),
	}
	if input.FirstName != "" {
		gocloakUser.FirstName = gocloak.StringP(input.FirstName)
	}
	if input.LastName != "" {
		gocloakUser.LastName = gocloak.StringP(input.LastName)
	}
	if input.Username != "" {
		gocloakUser.Username = gocloak.StringP(input.Username)
	}
	if input.Email != "" {
		gocloakUser.Email = gocloak.StringP(input.Email)
	}

	// updating user
	client.UpdateUser(ctx, token, keycloakRealm, gocloakUser)
	return auth.userPort.UpdateUser(user)
}
