package auth

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain"
	"auth/core/domain/http"

	"go.uber.org/zap"
)

func (svc *Service) Signup(ctx context.Context, input *http.SignupInput) (string, error) {
	svc.log.Info("signup")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.auth.signup")
	defer span.End()

	goClockUser, err := svc.authServerPort.CreateUser(ctx, input)
	if err != nil {
		svc.log.Error("create user error: ", zap.Error(err))
		return "", err
	}

	// create user entry inside database
	userUUID, err := svc.uuidPort.FromString(ctx, *goClockUser.ID)
	if err != nil {
		svc.log.Error("create user error: ", zap.Error(err))
		return "", err
	}

	user, err := domain.NewUser(userUUID)
	if err != nil {
		svc.log.Error("create user error: ", zap.Error(err))
		return "", err
	}

	user.SetId(userUUID)
	userId, err := svc.userPort.CreateUser(ctx, user)
	if err != nil {
		svc.log.Error("create user error: ", zap.Error(err))
		return "", err
	}

	svc.log.Info("user created", zap.String("userId", userId))
	return *goClockUser.ID, nil
}
