package auth_server

import (
	"context"
	"strings"
	"time"

	"auth/core/domain"
	"auth/core/domain/http"
	"auth/core/domain/system"
	"auth/core/ports"
	"auth/internal"

	"github.com/Nerzal/gocloak/v13"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service struct {
	log              *zap.Logger
	config           *internal.AppConfig
	uuidPort         ports.UUIDService
	subscriptionPort ports.SubscriptionService
}

func NewAuthServerService(log *zap.Logger, ctx context.Context, config *internal.AppConfig, subscriptionPort ports.SubscriptionService, uuidPort ports.UUIDService) (ports.AuthServerService, error) {
	return &Service{
		log:              log,
		config:           config,
		uuidPort:         uuidPort,
		subscriptionPort: subscriptionPort,
	}, nil
}

func (svc *Service) NewClient(ctx context.Context) *gocloak.GoCloak {
	svc.log.Info("creating new client")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	client := gocloak.NewClient(svc.config.KeycloakServer)
	return client
}

func (svc *Service) Login(ctx context.Context, email, password string) (*gocloak.JWT, error) {
	svc.log.Info("login")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	client := svc.NewClient(ctx)

	keycloakRealm := svc.config.KeycloakRealm
	keycloakClientId := svc.config.KeycloakClientId
	keycloakClientSecret := svc.config.KeycloakClientSecret

	token, err := client.Login(ctx, keycloakClientId, keycloakClientSecret, keycloakRealm, email, password)
	if err != nil {
		svc.log.Error("login error: ", zap.Error(err))

		if strings.HasPrefix(err.Error(), "401") {
			return nil, system.ErrAuthorization
		}
		if strings.HasPrefix(err.Error(), "400") {
			return nil, system.ErrInvalidInput
		}

		return nil, system.ErrServer
	}

	return token, nil
}

func (svc *Service) AccessToken(ctx context.Context) (string, error) {
	svc.log.Info("getting access token")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	client := svc.NewClient(ctx)

	keycloakRealm := svc.config.KeycloakRealm
	keycloakClientId := svc.config.KeycloakClientId
	keycloakClientSecret := svc.config.KeycloakClientSecret

	token, err := client.LoginClient(ctx, keycloakClientId, keycloakClientSecret, keycloakRealm)
	if err != nil {
		svc.log.Error("login client error: ", zap.Error(err))

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

func (svc *Service) RetrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	client := svc.NewClient(ctx)
	keycloakRealm := svc.config.KeycloakRealm
	keycloakClientId := svc.config.KeycloakClientId
	keycloakClientSecret := svc.config.KeycloakClientSecret

	rpt, err := client.GetRequestingPartyToken(ctx, accessToken, keycloakRealm, gocloak.RequestingPartyTokenOptions{
		Audience: &keycloakClientId,
	})
	if err != nil {
		svc.log.Error("Get requesting party token error: ", zap.Error(err))
		return nil, err
	}

	rptResult, err := client.RetrospectToken(ctx, rpt.AccessToken, keycloakClientId, keycloakClientSecret, keycloakRealm)
	if err != nil {
		svc.log.Error("Retrospect token error: ", zap.Error(err))
		return nil, err
	}

	return rptResult, nil
}

// CreateUser creates a new user
func (svc *Service) CreateUser(ctx context.Context, input *http.SignupInput) (*gocloak.User, error) {
	svc.log.Info("creating user")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	// create subscription
	subscription, err := domain.NewSubscription("free")
	if err != nil {
		svc.log.Error("create subscription error: ", zap.Error(err))
		return nil, err
	}

	subscriptionId, err := svc.subscriptionPort.CreateSubscription(ctx, subscription)
	if err != nil {
		svc.log.Error("create subscription error: ", zap.Error(err))
		return nil, err
	}

	client := svc.NewClient(ctx)

	token, err := svc.AccessToken(ctx)
	if err != nil {
		svc.log.Error("access token error: ", zap.Error(err))
		return nil, err
	}

	// create user
	svc.log.Debug("creating user in keycloak")
	userAttributes := &map[string][]string{
		"origin":       []string{"datawyse.io"},
		"language":     []string{"en"},
		"subscription": []string{subscriptionId.String()},
		"accountType":  []string{domain.ACCOUNT_ADMIN.String()},
	}
	goClockUser := gocloak.User{
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
		Attributes:      userAttributes,
		RequiredActions: &[]string{"VERIFY_EMAIL"},
	}
	keycloakUserId, err := client.CreateUser(ctx, token, svc.config.KeycloakRealm, goClockUser)
	goClockUser.ID = &keycloakUserId

	if err != nil {
		svc.log.Error("create user error: ", zap.Error(err))

		if strings.HasPrefix(err.Error(), "409") {
			return nil, system.ErrUserAlreadyExists
		}
		if strings.HasPrefix(err.Error(), "400") {
			return nil, system.ErrInvalidInput
		}

		return nil, err
	}

	// send email
	if err := client.ExecuteActionsEmail(ctx, token, svc.config.KeycloakRealm, gocloak.ExecuteActionsEmail{
		UserID:   &keycloakUserId,
		ClientID: &svc.config.KeycloakClientId,
		Actions:  &[]string{"VERIFY_EMAIL"},
	}); err != nil {
		svc.log.Error("email verification email error: ", zap.Error(err))
	}

	return &goClockUser, nil
}

func (svc *Service) GetUserById(ctx context.Context, id string) (*gocloak.User, error) {
	svc.log.Info("getting user by id")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	client := svc.NewClient(ctx)

	token, err := svc.AccessToken(ctx)
	if err != nil {
		svc.log.Error("access token error: ", zap.Error(err))
		return nil, err
	}

	user, err := client.GetUserByID(ctx, token, svc.config.KeycloakRealm, id)
	if err != nil {
		svc.log.Error("get user by id error: ", zap.Error(err))
		return nil, err
	}

	return user, nil
}

var ServiceModule = fx.Provide(NewAuthServerService)
