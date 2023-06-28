package auth

import (
	"context"
	"fmt"
	"time"

	"auth/core/domain/http"
	"auth/core/domain/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (ctrl *Controller) UpdateProfile(ctx *gin.Context) {
	ctrl.log.Info("updating profile up")

	authCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Duration(ctrl.RequestTimeout)*time.Second)
	defer cancel()

	var updateUserInput http.UpdateUserInput
	if err := ctx.ShouldBindJSON(&updateUserInput); err != nil {
		ctrl.log.Error("error binding json", zap.Error(err))

		err = ctx.Error(system.ErrInvalidInput)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	err := ctrl.validate.Struct(updateUserInput)
	if err != nil {
		ctrl.log.Error("error validating input", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	user, err := ctrl.service.UpdateUser(authCtx, &updateUserInput)
	if err != nil {
		ctrl.log.Error("error signing up", zap.Error(err))

		err = ctx.Error(err)
		if err != nil {
			ctrl.log.Error("error aborting with error", zap.Error(err))
		}
		return
	}

	message := fmt.Sprintf("user successfully created")
	ctx.JSON(201, system.NewHttpResponse(true, message, gin.H{"user": user}))
}
