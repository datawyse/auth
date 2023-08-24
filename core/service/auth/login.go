package auth

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"

	"auth/core/domain"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func (svc *Service) Login(ctx context.Context, email, password string) (*domain.AuthToken, error) {
	svc.log.Info("login")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	span := trace.SpanFromContext(ctx)
	tracerProvider := span.TracerProvider()
	ctx, span = tracerProvider.Tracer(svc.config.ServiceName).Start(ctx, "service.auth.login")
	defer span.End()

	span.SetAttributes(attribute.String("service.name", "auth.Login"))

	resToken, err := svc.authServerPort.Login(ctx, email, password)
	if err != nil {
		svc.log.Error("login error: ", zap.Error(err))
		return nil, err
	}

	idToken := resToken.IDToken
	tokenType := resToken.TokenType
	expiredIn := resToken.ExpiresIn
	accessToken := resToken.AccessToken
	refreshToken := resToken.RefreshToken

	// perform user last login update
	user, err := svc.userPort.UserByEmail(ctx, email)
	if err != nil {
		svc.log.Error("error finding user by email", zap.Error(err))
		return nil, err
	}

	systemUser, er := user.ToUser()
	if er != nil {
		svc.log.Error("error converting user to system user", zap.Error(er))
		return nil, er
	}

	systemUser.LastSignInAt = time.Now()
	systemUser, err = svc.userPort.UpdateUser(ctx, systemUser)
	if err != nil {
		svc.log.Error("error updating user", zap.Error(err))
		return nil, err
	}

	token := domain.NewAuthToken(accessToken, refreshToken, tokenType, expiredIn, idToken)
	return token, nil
}
