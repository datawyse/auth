package auth

import (
	"context"
	"time"

	"auth/core/domain"
	"auth/core/domain/http"

	"github.com/Nerzal/gocloak/v13"
	"go.uber.org/zap"
)

func (svc *Service) User(ctx context.Context, userId string) (*domain.UserInfo, error) {
	svc.log.Info("svc.service.user")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	return svc.userPort.User(ctx, userId)
}

func (svc *Service) UpdateUser(ctx context.Context, input *http.UpdateUserInput) (*domain.UserInfo, error) {
	svc.log.Info("svc.service.user")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	user, err := svc.User(ctx, input.Id)
	if err != nil {
		svc.log.Error("error getting user", zap.Error(err))
		return nil, err
	}

	// // for _, orgId := range input.Organizations {
	// // 	id, err := system.ToUUID(orgId)
	// // 	if err != nil {
	// // 		return nil, system.ErrInvalidInput
	// // 	}
	// // 	user.Organizations = append(user.Organizations, id)
	// // }
	//
	// // for _, roleId := range input.Roles {
	// // 	id, err := system.ToUUID(roleId)
	// // 	if err != nil {
	// // 		return nil, system.ErrInvalidInput
	// // 	}
	// // 	user.Roles = append(user.Roles, id)
	// // }
	// // if input.Language != "" {
	// // 	user.Language = input.Language
	// // }
	// // if input.AccountType != "" {
	// // 	user.AccountType = input.AccountType
	// // }
	//

	if input.FirstName != "" {
		user.KeycloakUser.FirstName = gocloak.StringP(input.FirstName)
	}
	if input.LastName != "" {
		user.KeycloakUser.LastName = gocloak.StringP(input.LastName)
	}
	if input.Username != "" {
		user.KeycloakUser.Username = gocloak.StringP(input.Username)
	}
	if input.Email != "" {
		user.KeycloakUser.Email = gocloak.StringP(input.Email)
	}

	attributes := *user.Attributes
	for name, attr := range input.Attributes {
		// if name is not in attributes, add it
		if _, ok := attributes[name]; !ok {
			attributes[name] = attr
		}
		// if name is in attributes, update it
		attributes[name] = attr
	}
	// loop through all the attributes and remove the ones that are not in input.Attributes
	for name := range attributes {
		if _, ok := input.Attributes[name]; !ok {
			delete(attributes, name)
		}
	}
	user.KeycloakUser.Attributes = &attributes

	// updating user
	client := svc.authServerPort.NewClient(ctx)
	token, err := svc.authServerPort.AccessToken(ctx)
	if err != nil {
		svc.log.Error("access token error: ", zap.Error(err))
		return nil, err
	}

	err = client.UpdateUser(ctx, token, svc.config.KeycloakRealm, user.KeycloakUser)
	if err != nil {
		svc.log.Error("update user error: ", zap.Error(err))
		return nil, err
	}

	systemUser, err := user.ToUser()
	if err != nil {
		svc.log.Error("error converting user to system user", zap.Error(err))
		return nil, err
	}

	systemUser, err = svc.userPort.UpdateUser(systemUser)
	if err != nil {
		svc.log.Error("error updating user", zap.Error(err))
		return nil, err
	}

	user.SystemUser = *systemUser
	return user, nil
}
