package auth

import (
	"context"
	"time"

	"auth/core/domain"
	"auth/core/domain/http"

	"go.uber.org/zap"
)

func (svc *Service) Signup(ctx context.Context, input *http.SignupInput) (string, error) {
	svc.log.Debug("svc.service Signup")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.config.ServiceTimeout)*time.Second)
	defer cancel()

	goClockUser, err := svc.authServerPort.CreateUser(ctx, input)
	if err != nil {
		svc.log.Error("create user error: ", zap.Error(err))
		return "", err
	}

	// create user entry inside database
	userUUID, err := svc.uuidPort.FromString(*goClockUser.ID)
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

	svc.log.Debug("user created", zap.String("userId", userId))
	return *goClockUser.ID, nil
}
