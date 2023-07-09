package auth

import (
	"context"
	"time"

	"auth/core/domain/http"
	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ctrl *Controller) Login(ctx *gin.Context) {
	ctrl.log.Debug("/ctrl/login")

	authCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Duration(ctrl.config.RequestTimeout)*time.Second)
	defer cancel()

	var loginInput http.LoginInput
	if err := ctx.BindJSON(&loginInput); err != nil {
		ctrl.log.Error("error binding json", zap.Error(err))

		err = ctx.Error(system.ErrInvalidInput)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	err := ctrl.validate.Struct(loginInput)
	if err != nil {
		ctrl.log.Error("error validating input", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	token, err := ctrl.service.Login(authCtx, loginInput.Email, loginInput.Password)
	if err != nil {
		ctrl.log.Error("error logging in", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	message := "user successfully logged in"
	res := system.NewHttpResponse(true, message, gin.H{
		"token": token,
	})
	ctx.JSON(200, res)
}
