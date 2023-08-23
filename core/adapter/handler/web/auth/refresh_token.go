package auth

import (
	"context"
	"time"

	"auth/core/domain/http"
	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ctrl *Controller) RefreshToken(ctx *gin.Context) {
	ctrl.log.Debug("/ctrl/login")

	authCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Duration(ctrl.config.RequestTimeout)*time.Second)
	defer cancel()

	var refreshTokenInput http.RefreshTokenInput
	if err := ctx.BindJSON(&refreshTokenInput); err != nil {
		ctrl.log.Error("error binding json", zap.Error(err))

		err = ctx.Error(system.ErrInvalidInput)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	if err := ctrl.validate.Validate(authCtx, refreshTokenInput); err != nil {
		ctrl.log.Error("error validating input", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	token, err := ctrl.service.RefreshToken(authCtx, refreshTokenInput.RefreshToken)
	if err != nil {
		ctrl.log.Error("error logging in", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	message := "refresh token successfully refreshed"
	res := system.NewHttpResponse(true, message, gin.H{"token": token}, 200)
	ctx.JSON(200, res)
}
