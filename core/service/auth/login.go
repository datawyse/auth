package auth

import (
	"context"
	"time"

	"auth/core/domain"

	"go.uber.org/zap"
)

func (svc *Service) Login(ctx context.Context, email, password string) (*domain.AuthToken, error) {
	svc.log.Debug("svc.service Login")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	resToken, err := svc.authServerPort.Login(ctx, email, password)
	if err != nil {
		svc.log.Error("login error: ", zap.Error(err))
		return nil, err
	}

	accessToken := resToken.AccessToken
	refreshToken := resToken.RefreshToken
	tokenType := resToken.TokenType
	expiredIn := resToken.ExpiresIn
	idToken := resToken.IDToken

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
	systemUser, err = svc.userPort.UpdateUser(systemUser)
	if err != nil {
		svc.log.Error("error updating user", zap.Error(err))
		return nil, err
	}

	token := domain.NewAuthToken(accessToken, refreshToken, tokenType, expiredIn, idToken)
	return token, nil
}
