package auth

import (
	"context"
	"strings"
	"time"

	"auth/core/domain"
	"auth/core/domain/http"
	"auth/core/domain/system"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (auth *Service) Signup(ctx context.Context, input *http.SignupInput) (string, error) {
	auth.log.Debug("auth.service Signup")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(auth.config.ServiceTimeout)*time.Second)
	defer cancel()

	token, err := auth.authServerPort.AccessToken()
	if err != nil {
		auth.log.Error("access token error: ", zap.Error(err))
		return "", err
	}

	keycloakServer := auth.config.KeycloakServer
	keycloakRealm := auth.config.KeycloakRealm

	client := gocloak.NewClient(keycloakServer)

	user, err := domain.NewUser(uuid.New())
	if err != nil {
		auth.log.Error("create user error: ", zap.Error(err))
		return "", err
	}

	// create user
	auth.log.Debug("creating user in keycloak", zap.String("userId", user.Id.String()))
	keycloakUserId, err := client.CreateUser(ctx, token, keycloakRealm, gocloak.User{
		ID:        gocloak.StringP(user.Id.String()),
		FirstName: gocloak.StringP(input.FirstName),
		LastName:  gocloak.StringP(input.LastName),
		Username:  gocloak.StringP(input.Username),
		Email:     gocloak.StringP(input.Email),
		Enabled:   gocloak.BoolP(true),
		Groups:    &[]string{"dashboard:admin", "project:admin", "organization:admin", "admin:user"},
		Access: &map[string]bool{
			"manageGroupMembership": true,
			"view":                  true,
			"mapRoles":              true,
			"impersonate":           true,
			"manage":                true,
		},
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP(input.Password),
			},
		},
		Attributes: &map[string][]string{
			"origin": []string{"local", "datawyse"},
		},
		RequiredActions: &[]string{
			"VERIFY_EMAIL",
		},
	})

	if err != nil {
		auth.log.Error("create user error: ", zap.Error(err))

		if strings.HasPrefix(err.Error(), "409") {
			return "", system.ErrUserAlreadyExists
		}
		if strings.HasPrefix(err.Error(), "400") {
			return "", system.ErrInvalidInput
		}

		return "", err
	}

	// send email
	if err := client.ExecuteActionsEmail(ctx, token, keycloakRealm, gocloak.ExecuteActionsEmail{
		UserID:   &keycloakUserId,
		ClientID: &auth.config.KeycloakClientId,
		Actions:  &[]string{"VERIFY_EMAIL"},
	}); err != nil {
		auth.log.Error("email verification email error: ", zap.Error(err))
	}

	// create subscription
	subscription, err := domain.NewSubscription()
	if err != nil {
		auth.log.Error("create subscription error: ", zap.Error(err))
		return "", err
	}

	subscriptionId, err := auth.subscriptionPort.CreateSubscription(subscription)
	if err != nil {
		auth.log.Error("create subscription error: ", zap.Error(err))
		return "", err
	}

	// create user entry inside database
	userUUID, err := system.ToUUID(keycloakUserId)
	if err != nil {
		auth.log.Error("create user error: ", zap.Error(err))
		return "", err
	}

	user.SetId(userUUID)
	user.SetSubscriptionId(subscriptionId)
	userId, err := auth.userPort.CreateUser(user)
	if err != nil {
		auth.log.Error("create user error: ", zap.Error(err))
		return "", err
	}

	auth.log.Debug("user created", zap.String("userId", userId))

	return keycloakUserId, nil
}
