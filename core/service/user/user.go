package user

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain"
	"auth/core/domain/system"
	"auth/core/ports"
	"auth/internal"

	"github.com/Nerzal/gocloak/v13"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service struct {
	log      *otelzap.Logger
	ctx      context.Context
	config   *internal.AppConfig
	uuidPort ports.UUIDService
	repo     ports.UserRepository
	authPort ports.AuthServerService
}

func NewUserService(log *otelzap.Logger, ctx context.Context, config *internal.AppConfig, authPort ports.AuthServerService, repo ports.UserRepository, uuidPort ports.UUIDService) (ports.UserService, error) {
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

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.user.create_user")
	defer span.End()

	userId, err := svc.repo.CreateUser(ctx, input)
	if err != nil {
		svc.log.Error("error creating user", zap.Error(err))
		return "", err
	}

	return userId.String(), err
}

func (svc *Service) UpdateUser(ctx context.Context, input *domain.User) (*domain.User, error) {
	svc.log.Info("updating user")

	ctx, cancel := context.WithTimeout(svc.ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "repository.user.update_user")
	defer span.End()

	user, err := svc.repo.UpdateUser(ctx, input)
	if err != nil {
		svc.log.Error("error creating user", zap.Error(err))
		return nil, err
	}

	return user, err
}

func (svc *Service) User(ctx context.Context, id string) (*domain.UserInfo, error) {
	svc.log.Info("finding user")

	ctx, cancel := context.WithTimeout(svc.ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.user.user")
	defer span.End()

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
	svc.log.Info("finding user by email")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.user.user_by_email")
	defer span.End()

	svc.log.Info("getting user by email")
	span.SetAttributes(attribute.String("service.name", "user.user.UserByEmail"))

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
	svc.log.Info("finding user by username")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.user.user_by_username")
	defer span.End()

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
	svc.log.Info("finding users")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.user.users")
	defer span.End()

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
