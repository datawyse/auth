package user

import (
	"context"
	"time"

	"auth/core/domain"
	"auth/core/domain/system"
	"auth/core/ports"
	"auth/internal"

	"github.com/Nerzal/gocloak/v13"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service struct {
	log      *zap.Logger
	ctx      context.Context
	config   *internal.AppConfig
	repo     ports.UserRepository
	authPort ports.AuthServerService
}

func NewUserService(log *zap.Logger, ctx context.Context, config *internal.AppConfig, authPort ports.AuthServerService, repo ports.UserRepository) ports.UserService {
	return &Service{
		log:      log,
		ctx:      ctx,
		config:   config,
		repo:     repo,
		authPort: authPort,
	}
}

func (svc *Service) CreateUser(input *domain.User) (string, error) {
	svc.log.Info("creating user")

	ctx, cancel := context.WithTimeout(svc.ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	userId, err := svc.repo.CreateUser(ctx, input)
	if err != nil {
		svc.log.Error("error creating user", zap.Error(err))
		return "", err
	}

	return userId.String(), err
}

func (svc *Service) UpdateUser(input *domain.User) (*domain.User, error) {
	svc.log.Info("creating user")

	ctx, cancel := context.WithTimeout(svc.ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	user, err := svc.repo.UpdateUser(ctx, input)
	if err != nil {
		svc.log.Error("error creating user", zap.Error(err))
		return nil, err
	}

	return user, err
}

func (svc *Service) User(id string) (*domain.UserInfo, error) {
	ctx, cancel := context.WithTimeout(svc.ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	keycloakServer := svc.config.KeycloakServer
	keycloakRealm := svc.config.KeycloakRealm
	client := gocloak.NewClient(keycloakServer)

	token, err := svc.authPort.AccessToken()
	if err != nil {
		return nil, err
	}

	keycloakUser, err := client.GetUserByID(ctx, token, keycloakRealm, id)
	if err != nil {
		svc.log.Error("error getting user", zap.Error(err))
		return nil, err
	}

	useridUUID, err := system.ToUUID(id)
	if err != nil {
		svc.log.Error("error getting user", zap.Error(err))
		return nil, err
	}

	systemUser, err := svc.repo.User(ctx, useridUUID)
	if err != nil {
		svc.log.Error("error getting user", zap.Error(err))
		return nil, err
	}

	userInfo := domain.NewUserInfo(keycloakUser, systemUser)
	return userInfo, nil
}

var ServiceModule = fx.Provide(NewUserService)
