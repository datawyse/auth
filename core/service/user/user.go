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
	uuidPort ports.UUIDService
	repo     ports.UserRepository
	authPort ports.AuthServerService
}

func NewUserService(log *zap.Logger, ctx context.Context, config *internal.AppConfig, authPort ports.AuthServerService, repo ports.UserRepository, uuidPort ports.UUIDService) (ports.UserService, error) {
	return &Service{
		log:      log,
		ctx:      ctx,
		config:   config,
		repo:     repo,
		authPort: authPort,
		uuidPort: uuidPort,
	}, nil
}

func (svc *Service) CreateUser(ctx context.Context, input *domain.User) (string, error) {
	svc.log.Info("creating user")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
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

func (svc *Service) User(ctx context.Context, id string) (*domain.UserInfo, error) {
	ctx, cancel := context.WithTimeout(svc.ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	keycloakUser, err := svc.authPort.GetUserById(ctx, id)
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

func (svc *Service) UserByEmail(ctx context.Context, email string) (*domain.UserInfo, error) {
	svc.log.Info("getting user by email")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	keycloakRealm := svc.config.KeycloakRealm
	client := svc.authPort.NewClient(ctx)

	token, err := svc.authPort.AccessToken(ctx)
	if err != nil {
		return nil, err
	}

	keycloakUsers, err := client.GetUsers(ctx, token, keycloakRealm, gocloak.GetUsersParams{
		Email: &email,
	})
	if err != nil {
		svc.log.Error("error getting user", zap.Error(err))
		return nil, err
	}

	if len(keycloakUsers) == 0 {
		return nil, system.ErrUserNotFound
	}

	keycloakUser := keycloakUsers[0]

	useridUUID, err := system.ToUUID(*keycloakUser.ID)
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

func (svc *Service) UserByUsername(ctx context.Context, username string) (*domain.UserInfo, error) {
	svc.log.Info("getting user by email")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	keycloakRealm := svc.config.KeycloakRealm
	client := svc.authPort.NewClient(ctx)

	token, err := svc.authPort.AccessToken(ctx)
	if err != nil {
		return nil, err
	}

	keycloakUsers, err := client.GetUsers(ctx, token, keycloakRealm, gocloak.GetUsersParams{
		Username: &username,
	})
	if err != nil {
		svc.log.Error("error getting user", zap.Error(err))
		return nil, err
	}

	if len(keycloakUsers) == 0 {
		return nil, system.ErrUserNotFound
	}

	keycloakUser := keycloakUsers[0]

	useridUUID, err := system.ToUUID(*keycloakUser.ID)
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

func (svc *Service) Users(ctx context.Context) ([]*domain.UserInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	keycloakServer := svc.config.KeycloakServer
	keycloakRealm := svc.config.KeycloakRealm
	client := gocloak.NewClient(keycloakServer)

	token, err := svc.authPort.AccessToken(ctx)
	if err != nil {
		return nil, err
	}

	keycloakUsers, err := client.GetUsers(ctx, token, keycloakRealm, gocloak.GetUsersParams{})
	if err != nil {
		svc.log.Error("error getting users", zap.Error(err))
		return nil, err
	}

	userInfos := make([]*domain.UserInfo, 0)

	for _, keycloakUser := range keycloakUsers {
		useridUUID, err := system.ToUUID(*keycloakUser.ID)
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
		userInfos = append(userInfos, userInfo)
	}

	return userInfos, nil
}

var ServiceModule = fx.Provide(NewUserService)
